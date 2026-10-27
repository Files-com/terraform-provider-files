package lib

import (
	"context"
	"encoding/json"
	"math/big"
	"strings"
	"testing"

	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/path"
	resourceschema "github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/schema/validator"
	"github.com/hashicorp/terraform-plugin-framework/tfsdk"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-go/tftypes"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

var behaviorValueVariants = []JSONSchemaVariant{
	{Name: "webhook", Value: "webhook", RootType: "object", LegacyProperty: "urls", LegacyList: true},
	{Name: "file_expiration", Value: "file_expiration", RootType: "object", LegacyProperty: "days_to_retain"},
	{Name: "storage_region", Value: "storage_region", RootType: "string"},
	{Name: "serve_publicly", Value: "serve_publicly", RootType: "object", Writable: []string{"key", "show_index", "password", "credentials"}, WriteOnly: []string{"password", "credentials.secret"}},
	{Name: "limit_file_regex", Value: "limit_file_regex", RootType: "array"},
	{Name: "malware_scanning", Value: "malware_scanning", RootType: "object"},
	{Name: "google_pub_sub", Value: "google_pub_sub", RootType: "object", WriteOnly: []string{"google_credentials.private_key"}},
}

func TestDynamicSiblingDiscriminatorToAPI(t *testing.T) {
	tests := []struct {
		name      string
		selector  string
		value     any
		expected  any
		expectErr bool
	}{
		{name: "legacy object", selector: "webhook", value: map[string]interface{}{"urls": []interface{}{"https://example.com"}}, expected: map[string]interface{}{"urls": []interface{}{"https://example.com"}}},
		{name: "legacy JSON object", selector: "webhook", value: `{"urls":["https://example.com"]}`, expected: map[string]interface{}{"urls": []interface{}{"https://example.com"}}},
		{name: "new wrapper", selector: "webhook", value: map[string]interface{}{"webhook": map[string]interface{}{"urls": []interface{}{"https://example.com"}}}, expected: map[string]interface{}{"urls": []interface{}{"https://example.com"}}},
		{name: "public null in a new wrapper", selector: "serve_publicly", value: map[string]interface{}{"serve_publicly": map[string]interface{}{"key": "Bar", "show_index": nil}}, expected: map[string]interface{}{"key": "Bar", "show_index": nil}},
		{name: "JSON encoded new wrapper", selector: "webhook", value: `{"webhook":{"urls":["https://example.com"]}}`, expected: map[string]interface{}{"urls": []interface{}{"https://example.com"}}},
		{name: "legacy JSON scalar", selector: "file_expiration", value: "14", expected: float64(14)},
		{name: "legacy JSON webhook string", selector: "webhook", value: `"https://example.com"`, expected: "https://example.com"},
		{name: "legacy scalar", selector: "storage_region", value: "us-east-1", expected: "us-east-1"},
		{name: "legacy array", selector: "limit_file_regex", value: []interface{}{`/Document-.*/`}, expected: []interface{}{`/Document-.*/`}},
		{name: "legacy empty object", selector: "malware_scanning", value: map[string]interface{}{}, expected: map[string]interface{}{}},
		{name: "wrapped file expiration scalar", selector: "file_expiration", value: map[string]interface{}{"file_expiration": float64(14)}, expectErr: true},
		{name: "wrapped webhook string", selector: "webhook", value: map[string]interface{}{"webhook": "https://example.com"}, expectErr: true},
		{name: "wrapped limit regex string", selector: "limit_file_regex", value: map[string]interface{}{"limit_file_regex": `/Document-.*/`}, expectErr: true},
		{name: "wrapped storage region array", selector: "storage_region", value: map[string]interface{}{"storage_region": []interface{}{"us-east-1"}}, expectErr: true},
		{name: "wrong wrapper", selector: "webhook", value: map[string]interface{}{"storage_region": "us-east-1"}, expectErr: true},
		{name: "multiple wrappers", selector: "webhook", value: map[string]interface{}{"webhook": map[string]interface{}{}, "storage_region": "us-east-1"}, expectErr: true},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			value := mustDynamicValue(t, test.value)
			result, diags := DynamicSiblingDiscriminatorToAPI(context.Background(), path.Root("value"), value, test.selector, behaviorValueVariants)
			assert.Equal(t, test.expectErr, diags.HasError())
			if !test.expectErr {
				assert.Equal(t, test.expected, result)
			}
		})
	}

	for _, value := range []types.Dynamic{types.DynamicNull(), types.DynamicUnknown()} {
		result, diags := DynamicSiblingDiscriminatorToAPI(context.Background(), path.Root("value"), value, "webhook", behaviorValueVariants)
		assert.False(t, diags.HasError())
		assert.Nil(t, result)
	}
}

func TestDynamicSiblingDiscriminatorUpdateToAPI(t *testing.T) {
	for name, values := range map[string][]any{
		"new wrapper":   {map[string]interface{}{"serve_publicly": map[string]interface{}{"key": "Bar", "show_index": true, "configured": true}}, map[string]interface{}{"serve_publicly": map[string]interface{}{"key": "Bar"}}},
		"legacy native": {map[string]interface{}{"key": "Bar", "show_index": true, "configured": true}, map[string]interface{}{"key": "Bar"}},
		"legacy JSON":   {`{"key":"Bar","show_index":true,"configured":true}`, `{"key":"Bar"}`},
	} {
		t.Run(name, func(t *testing.T) {
			result, diags := DynamicSiblingDiscriminatorUpdateToAPI(context.Background(), path.Root("value"), mustDynamicValue(t, values[1]), mustDynamicValue(t, values[0]), "serve_publicly", behaviorValueVariants)
			require.False(t, diags.HasError(), diags)
			assert.Equal(t, map[string]interface{}{"key": "Bar", "show_index": nil}, result)
		})
	}
}

func TestAPIToDynamicSiblingDiscriminatorPreservesRepresentation(t *testing.T) {
	servePubliclyAPI := map[string]interface{}{"key": "Bar", "show_index": true, "cors_enabled": false}
	// The API returns Google credentials as an object without private_key.
	topics := []interface{}{map[string]interface{}{"project_id": "example", "topic_id": "events"}}
	googlePubSubAPI := map[string]interface{}{"projects_topics": topics, "google_credentials": map[string]interface{}{"client_email": "events@example.com"}}
	keyFile := `{"client_email": "events@example.com", "private_key": "example-key"}`
	keyFileValue := func(credentials any) map[string]interface{} {
		return map[string]interface{}{"projects_topics": topics, "google_credentials": credentials}
	}
	credentialsObject := map[string]interface{}{"client_email": "events@example.com", "private_key": "example-key"}
	formattedLegacyJSON := `{ "key": "Bar", "show_index": true }`
	tests := []struct {
		name       string
		selector   string
		apiValue   any
		prior      types.Dynamic
		expected   any
		jsonString bool
	}{
		{name: "new wrapper", selector: "serve_publicly", apiValue: servePubliclyAPI, prior: mustDynamicValue(t, map[string]interface{}{"serve_publicly": map[string]interface{}{"key": "Bar", "show_index": true}}), expected: map[string]interface{}{"serve_publicly": map[string]interface{}{"key": "Bar", "show_index": true}}},
		{name: "legacy object", selector: "serve_publicly", apiValue: servePubliclyAPI, prior: mustDynamicValue(t, map[string]interface{}{"key": "Bar", "show_index": true}), expected: map[string]interface{}{"key": "Bar", "show_index": true}},
		{name: "write only value", selector: "serve_publicly", apiValue: servePubliclyAPI, prior: mustDynamicValue(t, map[string]interface{}{"serve_publicly": map[string]interface{}{"key": "Bar", "password": "secret"}}), expected: map[string]interface{}{"serve_publicly": map[string]interface{}{"key": "Bar", "password": "secret"}}},
		{name: "nested write only value", selector: "serve_publicly", apiValue: map[string]interface{}{"credentials": map[string]interface{}{"name": "updated"}}, prior: mustDynamicValue(t, map[string]interface{}{"serve_publicly": map[string]interface{}{"credentials": map[string]interface{}{"name": "configured", "secret": "secret"}}}), expected: map[string]interface{}{"serve_publicly": map[string]interface{}{"credentials": map[string]interface{}{"name": "updated", "secret": "secret"}}}},
		{name: "removed configured value", selector: "serve_publicly", apiValue: map[string]interface{}{"key": "Bar"}, prior: mustDynamicValue(t, map[string]interface{}{"serve_publicly": map[string]interface{}{"key": "Bar", "show_index": true}}), expected: map[string]interface{}{"serve_publicly": map[string]interface{}{"key": "Bar"}}},
		{name: "removed legacy value", selector: "serve_publicly", apiValue: map[string]interface{}{"key": "Bar"}, prior: mustDynamicValue(t, map[string]interface{}{"key": "Bar", "show_index": true}), expected: map[string]interface{}{"key": "Bar"}},
		{name: "removed legacy JSON value", selector: "serve_publicly", apiValue: map[string]interface{}{"key": "Bar"}, prior: mustDynamicValue(t, `{"key":"Bar","show_index":true}`), expected: map[string]interface{}{"key": "Bar"}, jsonString: true},
		{name: "legacy JSON", selector: "serve_publicly", apiValue: servePubliclyAPI, prior: mustDynamicValue(t, formattedLegacyJSON), expected: map[string]interface{}{"key": "Bar", "show_index": true}, jsonString: true},
		{name: "file expiration JSON scalar", selector: "file_expiration", apiValue: map[string]interface{}{"days_to_retain": float64(14), "delete_empty_folders": false}, prior: mustDynamicValue(t, "14"), expected: float64(14), jsonString: true},
		{name: "webhook JSON scalar", selector: "webhook", apiValue: map[string]interface{}{"urls": []interface{}{"https://example.com"}}, prior: mustDynamicValue(t, `"https://example.com"`), expected: "https://example.com", jsonString: true},
		{name: "file expiration scalar", selector: "file_expiration", apiValue: map[string]interface{}{"days_to_retain": float64(14), "delete_empty_folders": false}, prior: types.DynamicValue(types.NumberValue(big.NewFloat(14))), expected: float64(14)},
		{name: "webhook scalar", selector: "webhook", apiValue: map[string]interface{}{"urls": []interface{}{"https://example.com"}}, prior: mustDynamicValue(t, "https://example.com"), expected: "https://example.com"},
		{name: "webhook array", selector: "webhook", apiValue: map[string]interface{}{"urls": []interface{}{"https://example.com"}}, prior: mustDynamicValue(t, []interface{}{"https://example.com"}), expected: []interface{}{"https://example.com"}},
		{name: "empty object", selector: "malware_scanning", apiValue: map[string]interface{}{}, prior: mustDynamicValue(t, map[string]interface{}{}), expected: map[string]interface{}{}},
		{name: "credentials key file text", selector: "google_pub_sub", apiValue: googlePubSubAPI, prior: mustDynamicValue(t, map[string]interface{}{"google_pub_sub": keyFileValue(keyFile)}), expected: map[string]interface{}{"google_pub_sub": keyFileValue(keyFile)}},
		{name: "legacy credentials key file text", selector: "google_pub_sub", apiValue: googlePubSubAPI, prior: mustDynamicValue(t, keyFileValue(keyFile)), expected: keyFileValue(keyFile)},
		{name: "credentials object", selector: "google_pub_sub", apiValue: googlePubSubAPI, prior: mustDynamicValue(t, map[string]interface{}{"google_pub_sub": keyFileValue(credentialsObject)}), expected: map[string]interface{}{"google_pub_sub": keyFileValue(credentialsObject)}},
		{name: "changed credentials", selector: "google_pub_sub", apiValue: keyFileValue(map[string]interface{}{"client_email": "other@example.com"}), prior: mustDynamicValue(t, map[string]interface{}{"google_pub_sub": keyFileValue(keyFile)}), expected: map[string]interface{}{"google_pub_sub": keyFileValue(map[string]interface{}{"client_email": "other@example.com"})}},
		{name: "body text", selector: "google_pub_sub", apiValue: map[string]interface{}{"body": `{"a":2}`}, prior: mustDynamicValue(t, map[string]interface{}{"google_pub_sub": map[string]interface{}{"body": `{"a": 1}`}}), expected: map[string]interface{}{"google_pub_sub": map[string]interface{}{"body": `{"a":2}`}}},
		{name: "null prior", selector: "serve_publicly", apiValue: servePubliclyAPI, prior: types.DynamicNull(), expected: servePubliclyAPI},
		{name: "unknown prior", selector: "serve_publicly", apiValue: servePubliclyAPI, prior: types.DynamicUnknown(), expected: servePubliclyAPI},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			result, diags := APIToDynamicSiblingDiscriminator(context.Background(), path.Root("value"), test.apiValue, test.selector, behaviorValueVariants, test.prior)
			require.False(t, diags.HasError(), diags)
			actual := dynamicValueInterface(t, result)
			if test.jsonString {
				var decoded any
				require.NoError(t, json.Unmarshal([]byte(actual.(string)), &decoded))
				actual = decoded
			}
			assert.Equal(t, test.expected, actual)
		})
	}
	formattedResult, diags := APIToDynamicSiblingDiscriminator(context.Background(), path.Root("value"), servePubliclyAPI, "serve_publicly", behaviorValueVariants, mustDynamicValue(t, formattedLegacyJSON))
	require.False(t, diags.HasError(), diags)
	assert.Equal(t, formattedLegacyJSON, dynamicValueInterface(t, formattedResult))

	configured := mustDynamicValue(t, map[string]interface{}{"serve_publicly": map[string]interface{}{"password": "secret"}})
	omittedResult, diags := APIToDynamicSiblingDiscriminator(context.Background(), path.Root("value"), nil, "serve_publicly", behaviorValueVariants, configured)
	require.False(t, diags.HasError(), diags)
	assert.Equal(t, dynamicValueInterface(t, configured), dynamicValueInterface(t, omittedResult))

	omittedResult, diags = APIToDynamicSiblingDiscriminator(context.Background(), path.Root("value"), nil, "serve_publicly", behaviorValueVariants, mustDynamicValue(t, map[string]interface{}{"serve_publicly": map[string]interface{}{"show_index": true}}))
	require.False(t, diags.HasError(), diags)
	assert.True(t, omittedResult.IsNull())
}

func TestDeprecatedSiblingDiscriminatorValue(t *testing.T) {
	valueValidator := DeprecatedSiblingDiscriminatorValue("files_behavior.value", "behavior", behaviorValueVariants, "legacy warning")
	tests := []struct {
		name         string
		value        types.Dynamic
		warnings     int
		expectErrors bool
	}{
		{name: "legacy native", value: mustDynamicValue(t, map[string]interface{}{"urls": []interface{}{"https://example.com"}}), warnings: 1},
		{name: "legacy JSON", value: mustDynamicValue(t, `{"urls":["https://example.com"]}`), warnings: 1},
		{name: "JSON-encoded new wrapper", value: mustDynamicValue(t, `{"webhook":{"urls":["https://example.com"]}}`), warnings: 1},
		{name: "JSON-encoded wrong wrapper", value: mustDynamicValue(t, `{"storage_region":"us-east-1"}`), expectErrors: true},
		{name: "new wrapper", value: mustDynamicValue(t, map[string]interface{}{"webhook": map[string]interface{}{"urls": []interface{}{"https://example.com"}}})},
		{name: "wrapped legacy value", value: mustDynamicValue(t, map[string]interface{}{"webhook": "https://example.com"}), expectErrors: true},
		{name: "wrong wrapper", value: mustDynamicValue(t, map[string]interface{}{"storage_region": "us-east-1"}), expectErrors: true},
		{name: "multiple wrappers", value: mustDynamicValue(t, map[string]interface{}{"webhook": map[string]interface{}{}, "storage_region": "us-east-1"}), expectErrors: true},
		{name: "null", value: types.DynamicNull()},
		{name: "unknown", value: types.DynamicUnknown()},
		{name: "partially unknown", value: types.DynamicValue(types.ObjectValueMust(
			map[string]attr.Type{"gpg_key_ids": types.TupleType{ElemTypes: []attr.Type{types.NumberType}}},
			map[string]attr.Value{"gpg_key_ids": types.TupleValueMust([]attr.Type{types.NumberType}, []attr.Value{types.NumberUnknown()})},
		))},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			request := validatorRequest(t, "webhook", test.value)
			response := &validator.DynamicResponse{}
			valueValidator.ValidateDynamic(context.Background(), request, response)
			assert.Equal(t, test.expectErrors, response.Diagnostics.HasError())
			assert.Equal(t, test.warnings, response.Diagnostics.WarningsCount())
		})
	}

	legacyValue := mustDynamicValue(t, map[string]interface{}{"urls": []interface{}{"https://example.com"}})
	response := &validator.DynamicResponse{}
	valueValidator.ValidateDynamic(context.Background(), validatorRequest(t, "webhook", legacyValue), response)
	assert.Equal(t, "Deprecated files_behavior.value format: webhook", response.Diagnostics[0].Summary())
	assert.Contains(t, response.Diagnostics[0].Detail(), `Wrap value under "webhook": value = { webhook = <documented webhook value> }. Remove jsonencode(...) if present.`)
	assert.Contains(t, response.Diagnostics[0].Detail(), "legacy warning")

	// A JSON string that already uses the new nesting only needs its encoding removed.
	response = &validator.DynamicResponse{}
	valueValidator.ValidateDynamic(context.Background(), validatorRequest(t, "webhook", mustDynamicValue(t, `{"webhook":{"urls":["https://example.com"]}}`)), response)
	assert.Equal(t, "Deprecated JSON encoding: files_behavior.value", response.Diagnostics[0].Summary())
	assert.Equal(t, `Remove jsonencode(...) and write the value as native HCL, keeping it under "webhook". legacy warning`, response.Diagnostics[0].Detail())
	regexVariants := []JSONSchemaVariant{{Name: "limit_file_regex", Value: "limit_file_regex", RootType: "array", ValueTypes: map[string]string{"": "list"}}}
	response = &validator.DynamicResponse{}
	DeprecatedSiblingDiscriminatorValue("files_behavior.value", "behavior", regexVariants, "legacy warning").ValidateDynamic(context.Background(), validatorRequest(t, "limit_file_regex", mustDynamicValue(t, `{"limit_file_regex":["/a/"]}`)), response)
	require.Equal(t, 1, response.Diagnostics.WarningsCount())
	assert.NotContains(t, response.Diagnostics[0].Detail(), "Write limit_file_regex as lists", "the wrapped value already is a list")

	response = &validator.DynamicResponse{}
	valueValidator.ValidateDynamic(context.Background(), validatorRequest(t, "webhook", mustDynamicValue(t, map[string]interface{}{"storage_region": "us-east-1"})), response)
	assert.Equal(t, "Invalid Value Wrapper", response.Diagnostics[0].Summary())
	response = &validator.DynamicResponse{}
	valueValidator.ValidateDynamic(context.Background(), validatorRequest(t, "webhook", mustDynamicValue(t, map[string]interface{}{"webhook": "https://example.com"})), response)
	assert.Contains(t, response.Diagnostics[0].Detail(), `value wrapper "webhook" must contain object, got string; keep the legacy value unwrapped or use the documented new format`)

	for _, selector := range []types.String{types.StringNull(), types.StringUnknown()} {
		response := &validator.DynamicResponse{}
		valueValidator.ValidateDynamic(context.Background(), validatorRequestWithSelector(t, selector, legacyValue), response)
		assert.Empty(t, response.Diagnostics)
	}
}

func TestDeprecatedSiblingDiscriminatorValueTypeGuidance(t *testing.T) {
	variants := []JSONSchemaVariant{{Name: "file_expiration", Value: "file_expiration", RootType: "object", LegacyProperty: "days_to_retain", ValueTypes: map[string]string{"days_to_retain": "integer", "delete_empty_folders": "boolean"}}}
	valueValidator := DeprecatedSiblingDiscriminatorValue("files_behavior.value", "behavior", variants, "legacy warning")
	const legacy, typed = "Deprecated files_behavior.value format: file_expiration", "Typed values expected: files_behavior.value"
	const guidance = "Write days_to_retain as integers, for example days_to_retain = 123."
	for _, test := range []struct {
		name, summary, detail string
		value                 types.Dynamic
	}{
		{name: "wrapped typed values", value: mustDynamicValue(t, map[string]any{"file_expiration": map[string]any{"days_to_retain": float64(30), "delete_empty_folders": true}})},
		{name: "wrapped strings", value: mustDynamicValue(t, map[string]any{"file_expiration": map[string]any{"days_to_retain": "30", "delete_empty_folders": "true"}}),
			summary: typed, detail: guidance + " Write delete_empty_folders as true or false. legacy warning"},
		{name: "legacy object", value: mustDynamicValue(t, map[string]any{"days_to_retain": "30"}), summary: legacy, detail: "legacy warning " + guidance},
		{name: "legacy JSON", value: mustDynamicValue(t, `{"days_to_retain":"30"}`), summary: legacy, detail: "legacy warning " + guidance},
		{name: "legacy bare value", value: mustDynamicValue(t, "30"), summary: legacy, detail: "legacy warning"},
	} {
		t.Run(test.name, func(t *testing.T) {
			response := &validator.DynamicResponse{}
			valueValidator.ValidateDynamic(context.Background(), validatorRequest(t, "file_expiration", test.value), response)
			require.False(t, response.Diagnostics.HasError(), response.Diagnostics)
			if test.summary == "" {
				assert.Zero(t, response.Diagnostics.WarningsCount())
				return
			}
			require.Equal(t, 1, response.Diagnostics.WarningsCount())
			assert.Equal(t, test.summary, response.Diagnostics[0].Summary())
			assert.True(t, strings.HasSuffix(response.Diagnostics[0].Detail(), test.detail), response.Diagnostics[0].Detail())
			assert.NotContains(t, response.Diagnostics[0].Detail(), "30")
		})
	}
}

func mustDynamicValue(t *testing.T, value any) types.Dynamic {
	t.Helper()
	dynamic, diags := ToDynamic(context.Background(), path.Root("value"), value, nil)
	require.False(t, diags.HasError(), diags)
	return dynamic
}

func dynamicValueInterface(t *testing.T, value types.Dynamic) any {
	t.Helper()
	result, diags := DynamicToInterface(context.Background(), path.Root("value"), value)
	require.False(t, diags.HasError(), diags)
	return result
}

func validatorRequest(t *testing.T, selector string, value types.Dynamic) validator.DynamicRequest {
	return validatorRequestWithSelector(t, types.StringValue(selector), value)
}

func validatorRequestWithSelector(t *testing.T, selector types.String, value types.Dynamic) validator.DynamicRequest {
	t.Helper()
	ctx := context.Background()
	schema := resourceschema.Schema{Attributes: map[string]resourceschema.Attribute{
		"behavior": resourceschema.StringAttribute{Required: true},
		"value":    resourceschema.DynamicAttribute{Optional: true},
	}}
	selectorValue, err := selector.ToTerraformValue(ctx)
	require.NoError(t, err)
	valueValue, err := value.ToTerraformValue(ctx)
	require.NoError(t, err)
	config := tfsdk.Config{
		Schema: schema,
		Raw: tftypes.NewValue(schema.Type().TerraformType(ctx), map[string]tftypes.Value{
			"behavior": selectorValue,
			"value":    valueValue,
		}),
	}
	return validator.DynamicRequest{Path: path.Root("value"), Config: config, ConfigValue: value}
}
