package provider

import (
	"bytes"
	"context"
	"encoding/json"
	"io"
	"net/http"
	"net/http/httptest"
	"sync"
	"testing"

	"github.com/hashicorp/terraform-plugin-framework/providerserver"
	"github.com/hashicorp/terraform-plugin-go/tfprotov6"
	"github.com/hashicorp/terraform-plugin-go/tftypes"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// These tests drive the provider through its protocol server, so plans can
// hold unknown values directly and no Terraform CLI is needed.

const remoteMountBackendTypeName = "files_remote_mount_backend"

// State as earlier provider versions stored it: decimals are strings.
var remoteMountBackendState = map[string]tftypes.Value{
	"id":                     tftypes.NewValue(tftypes.Number, 7),
	"canary_file_path":       tftypes.NewValue(tftypes.String, "canary.txt"),
	"remote_server_mount_id": tftypes.NewValue(tftypes.Number, 2),
	"remote_server_id":       tftypes.NewValue(tftypes.Number, 3),
	"min_free_cpu":           tftypes.NewValue(tftypes.String, "1.50"),
	"min_free_mem":           tftypes.NewValue(tftypes.String, "0.0049999999999999999999999999"),
}

type remoteMountBackendRequest struct {
	method string
	body   map[string]interface{}
}

type remoteMountBackendHarness struct {
	server   tfprotov6.ProviderServer
	schema   *tfprotov6.Schema
	mu       sync.Mutex
	requests []remoteMountBackendRequest
}

// newRemoteMountBackendHarness configures the provider against a local API
// that stores backends and, like the real API, returns decimals as strings.
func newRemoteMountBackendHarness(t *testing.T) *remoteMountBackendHarness {
	h := &remoteMountBackendHarness{server: providerserver.NewProtocol6(New("test")())()}
	stored := map[string]interface{}{"id": 7, "canary_file_path": "canary.txt", "remote_server_mount_id": 2, "remote_server_id": 3, "min_free_cpu": "1.50", "min_free_mem": "0.0049999999999999999999999999"}
	api := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		h.mu.Lock()
		defer h.mu.Unlock()
		var body map[string]interface{}
		if data, _ := io.ReadAll(r.Body); len(data) > 0 {
			decoder := json.NewDecoder(bytes.NewReader(data))
			decoder.UseNumber() // compare number tokens, not float64 conversions
			assert.NoError(t, decoder.Decode(&body))
		}
		h.requests = append(h.requests, remoteMountBackendRequest{r.Method, body})
		for key, value := range body {
			if number, ok := value.(json.Number); ok && (key == "min_free_cpu" || key == "min_free_mem") {
				value = number.String()
			}
			stored[key] = value
		}
		w.Header().Set("Content-Type", "application/json")
		assert.NoError(t, json.NewEncoder(w).Encode(stored))
	}))
	t.Cleanup(api.Close)

	ctx := context.Background()
	schemas, err := h.server.GetProviderSchema(ctx, &tfprotov6.GetProviderSchemaRequest{})
	require.NoError(t, err)
	require.Empty(t, schemas.Diagnostics)
	h.schema = schemas.ResourceSchemas[remoteMountBackendTypeName]
	configured, err := h.server.ConfigureProvider(ctx, &tfprotov6.ConfigureProviderRequest{
		Config: dynamicValue(t, objectValue(schemas.Provider, map[string]tftypes.Value{
			"api_key":           tftypes.NewValue(tftypes.String, "test"),
			"endpoint_override": tftypes.NewValue(tftypes.String, api.URL),
		})),
	})
	require.NoError(t, err)
	require.Empty(t, configured.Diagnostics)
	return h
}

// decimalConfig sets the required attributes plus the two decimals; a nil
// decimal is left out of the configuration.
func decimalConfig(cpu, mem *string) map[string]tftypes.Value {
	config := map[string]tftypes.Value{
		"canary_file_path":       remoteMountBackendState["canary_file_path"],
		"remote_server_mount_id": remoteMountBackendState["remote_server_mount_id"],
		"remote_server_id":       remoteMountBackendState["remote_server_id"],
	}
	if cpu != nil {
		config["min_free_cpu"] = tftypes.NewValue(tftypes.String, *cpu)
	}
	if mem != nil {
		config["min_free_mem"] = tftypes.NewValue(tftypes.String, *mem)
	}
	return config
}

// apply runs Create when prior is nil and Update otherwise. Unconfigured
// computed attributes are planned as unknown on Create and as their prior
// value on Update, as their UseStateForUnknown plan modifiers do.
func (h *remoteMountBackendHarness) apply(t *testing.T, prior, config map[string]tftypes.Value) *tfprotov6.ApplyResourceChangeResponse {
	priorState := tftypes.NewValue(h.schema.ValueType(), nil)
	priorValues := map[string]tftypes.Value{}
	if prior != nil {
		priorState = objectValue(h.schema, prior)
		require.NoError(t, priorState.As(&priorValues))
	}
	planned := map[string]tftypes.Value{}
	for _, attribute := range h.schema.Block.Attributes {
		value, configured := config[attribute.Name]
		switch {
		case configured:
			planned[attribute.Name] = value
		case !attribute.Computed:
			// stays null
		case prior == nil:
			planned[attribute.Name] = tftypes.NewValue(attribute.ValueType(), tftypes.UnknownValue)
		default:
			planned[attribute.Name] = priorValues[attribute.Name]
		}
	}
	response, err := h.server.ApplyResourceChange(context.Background(), &tfprotov6.ApplyResourceChangeRequest{
		TypeName:     remoteMountBackendTypeName,
		PriorState:   dynamicValue(t, priorState),
		PlannedState: dynamicValue(t, objectValue(h.schema, planned)),
		Config:       dynamicValue(t, objectValue(h.schema, config)),
	})
	require.NoError(t, err)
	return response
}

func (h *remoteMountBackendHarness) takeRequests() []remoteMountBackendRequest {
	h.mu.Lock()
	defer h.mu.Unlock()
	requests := h.requests
	h.requests = nil
	return requests
}

func TestRemoteMountBackendDecimalsSendDecimalTextExactlyAndLegacySpellingsAsNumbers(t *testing.T) {
	exact := "1.0049999999999999999999999999"
	for _, test := range []struct {
		name     string
		cpu, mem *string
		// A string is sent as exact text; a json.Number as the former float64.
		create, update map[string]interface{}
	}{
		{
			"exact text and hexadecimal float", &exact, text("0x1.0000000000001p0"),
			map[string]interface{}{"min_free_cpu": exact, "min_free_mem": json.Number("1.0000000000000002")},
			map[string]interface{}{"min_free_cpu": exact, "min_free_mem": json.Number("1.0000000000000002")},
		},
		{
			"underscore float and exact zero", text("1.2345678901234567_89"), text("0"),
			map[string]interface{}{"min_free_cpu": json.Number("1.2345678901234567"), "min_free_mem": "0"},
			map[string]interface{}{"min_free_cpu": json.Number("1.2345678901234567"), "min_free_mem": "0"},
		},
		{
			"signed legacy spellings", text("+0x1.8p+2"), text("-1_2.3_4e+0_2"),
			map[string]interface{}{"min_free_cpu": json.Number("6"), "min_free_mem": json.Number("-1234")},
			map[string]interface{}{"min_free_cpu": json.Number("6"), "min_free_mem": json.Number("-1234")},
		},
		{
			// Create still omits a zero float64, as before; Update sends it.
			"legacy zeros", text("0x0p0"), text("-0x0p0"),
			map[string]interface{}{},
			map[string]interface{}{"min_free_cpu": json.Number("0"), "min_free_mem": json.Number("-0")},
		},
		{
			"exact text beyond float64 range and unconfigured", text("1e9999"), nil,
			map[string]interface{}{"min_free_cpu": "1e9999"},
			map[string]interface{}{"min_free_cpu": "1e9999"},
		},
	} {
		t.Run(test.name, func(t *testing.T) {
			h := newRemoteMountBackendHarness(t)

			created := h.apply(t, nil, decimalConfig(test.cpu, test.mem))
			require.Empty(t, created.Diagnostics)
			updated := h.apply(t, remoteMountBackendState, decimalConfig(test.cpu, test.mem))
			require.Empty(t, updated.Diagnostics)

			requests := h.takeRequests()
			require.Len(t, requests, 2)
			for n, want := range []struct {
				method string
				values map[string]interface{}
			}{{http.MethodPost, test.create}, {http.MethodPatch, test.update}} {
				assert.Equal(t, want.method, requests[n].method)
				for _, key := range []string{"min_free_cpu", "min_free_mem"} {
					value, sent := requests[n].body[key]
					wantValue, wantSent := want.values[key]
					assert.Equal(t, wantSent, sent, "%s %s sent", want.method, key)
					assert.Equal(t, wantValue, value, "%s %s", want.method, key)
				}
			}
		})
	}
}

func TestRemoteMountBackendDecimalsRejectInvalidTextBeforeRequest(t *testing.T) {
	for _, test := range []struct {
		name, attribute string
		cpu, mem        *string
	}{
		{"repeated underscores", "min_free_cpu", text("1__0"), nil},
		{"hexadecimal beyond float64 range", "min_free_mem", nil, text("0x1p1024")},
		{"NaN", "min_free_cpu", text("NaN"), nil},
		{"infinity", "min_free_mem", nil, text("+Infinity")},
		{"empty text", "min_free_cpu", text(""), nil},
		{"surrounding space", "min_free_mem", text("2.5"), text(" 1.5")},
	} {
		t.Run(test.name, func(t *testing.T) {
			h := newRemoteMountBackendHarness(t)
			for _, prior := range []map[string]tftypes.Value{nil, remoteMountBackendState} {
				response := h.apply(t, prior, decimalConfig(test.cpu, test.mem))

				require.Len(t, response.Diagnostics, 1)
				diagnostic := response.Diagnostics[0]
				assert.Equal(t, tfprotov6.DiagnosticSeverityError, diagnostic.Severity)
				assert.True(t, diagnostic.Attribute.Equal(tftypes.NewAttributePath().WithAttributeName(test.attribute)), "diagnostic on %s, got %s", test.attribute, diagnostic.Attribute)
				assert.Equal(t, "Could not parse "+test.attribute+": expected a decimal number such as 1.5 or 2e-3", diagnostic.Detail)
			}
			assert.Empty(t, h.takeRequests(), "invalid input must not reach the API")
		})
	}
}

func TestRemoteMountBackendDecimalStateStaysString(t *testing.T) {
	h := newRemoteMountBackendHarness(t)
	ctx := context.Background()
	resourceType := h.schema.ValueType()

	assert.Zero(t, h.schema.Version, "no state migration")
	for _, attribute := range h.schema.Block.Attributes {
		if attribute.Name == "min_free_cpu" || attribute.Name == "min_free_mem" {
			assert.True(t, attribute.Type.Is(tftypes.String), "%s stays a string attribute", attribute.Name)
		}
	}

	stateJSON, err := json.Marshal(map[string]interface{}{
		"id": 7, "canary_file_path": "canary.txt", "remote_server_mount_id": 2, "remote_server_id": 3,
		"min_free_cpu": "1.50", "min_free_mem": "0.0049999999999999999999999999",
	})
	require.NoError(t, err)
	upgraded, err := h.server.UpgradeResourceState(ctx, &tfprotov6.UpgradeResourceStateRequest{
		TypeName: remoteMountBackendTypeName,
		Version:  0,
		RawState: &tfprotov6.RawState{JSON: stateJSON},
	})
	require.NoError(t, err)
	require.Empty(t, upgraded.Diagnostics)
	assertDecimalState(t, resourceType, upgraded.UpgradedState)

	read, err := h.server.ReadResource(ctx, &tfprotov6.ReadResourceRequest{
		TypeName:     remoteMountBackendTypeName,
		CurrentState: dynamicValue(t, objectValue(h.schema, remoteMountBackendState)),
	})
	require.NoError(t, err)
	require.Empty(t, read.Diagnostics)
	assertDecimalState(t, resourceType, read.NewState)

	imported, err := h.server.ImportResourceState(ctx, &tfprotov6.ImportResourceStateRequest{TypeName: remoteMountBackendTypeName, ID: "7"})
	require.NoError(t, err)
	require.Empty(t, imported.Diagnostics)
	require.Len(t, imported.ImportedResources, 1)
	refreshed, err := h.server.ReadResource(ctx, &tfprotov6.ReadResourceRequest{
		TypeName:     remoteMountBackendTypeName,
		CurrentState: imported.ImportedResources[0].State,
	})
	require.NoError(t, err)
	require.Empty(t, refreshed.Diagnostics)
	assertDecimalState(t, resourceType, refreshed.NewState)
}

// assertDecimalState checks that the decimals keep the API's text, such as
// "1.50", rather than a reformatted number.
func assertDecimalState(t *testing.T, resourceType tftypes.Type, state *tfprotov6.DynamicValue) {
	t.Helper()
	value, err := state.Unmarshal(resourceType)
	require.NoError(t, err)
	var attributes map[string]tftypes.Value
	require.NoError(t, value.As(&attributes))
	for key, want := range map[string]string{"min_free_cpu": "1.50", "min_free_mem": "0.0049999999999999999999999999"} {
		var got string
		require.NoError(t, attributes[key].As(&got))
		assert.Equal(t, want, got, key)
	}
}

func objectValue(schema *tfprotov6.Schema, values map[string]tftypes.Value) tftypes.Value {
	objectType := schema.ValueType().(tftypes.Object)
	attributes := make(map[string]tftypes.Value, len(objectType.AttributeTypes))
	for name, attributeType := range objectType.AttributeTypes {
		if value, ok := values[name]; ok {
			attributes[name] = value
		} else {
			attributes[name] = tftypes.NewValue(attributeType, nil)
		}
	}
	return tftypes.NewValue(objectType, attributes)
}

func dynamicValue(t *testing.T, value tftypes.Value) *tfprotov6.DynamicValue {
	t.Helper()
	dynamic, err := tfprotov6.NewDynamicValue(value.Type(), value)
	require.NoError(t, err)
	return &dynamic
}

func text(value string) *string {
	return &value
}
