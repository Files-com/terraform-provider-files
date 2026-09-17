package provider

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"sync"
	"testing"

	files_sdk "github.com/Files-com/files-sdk-go/v3"
	"github.com/hashicorp/terraform-plugin-framework/provider"
	"github.com/hashicorp/terraform-plugin-framework/providerserver"
	"github.com/hashicorp/terraform-plugin-go/tfprotov6"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

type jsonTransitionTestProvider struct {
	filesProvider
	endpoint string
}

func (p *jsonTransitionTestProvider) Configure(ctx context.Context, req provider.ConfigureRequest, resp *provider.ConfigureResponse) {
	p.filesProvider.Configure(ctx, req, resp)
	if resp.Diagnostics.HasError() {
		return
	}
	config := files_sdk.Config{APIKey: "test", EndpointOverride: p.endpoint}.Init()
	resp.ResourceData, resp.DataSourceData = config, config
}

func TestJSONPropertyResourceTransitions(t *testing.T) {
	for _, test := range []struct {
		name, model, property, base, value string
		normalize                          func(any) any
	}{
		{"site/numeric string watermark", "site", "bundle_watermark_value", "", `{ transparency = "25", gravity = "SouthWest" }`, func(value any) any {
			value.(map[string]any)["transparency"] = float64(25)
			return value
		}},
		{"bundle/watermark value", "bundle", "watermark_value", `paths = ["reports"]`, `{ gravity = "Center", transparency = "25", dynamic_text = "Confidential" }`, func(value any) any {
			value.(map[string]any)["transparency"] = float64(25)
			return value
		}},
		{"siem_http_destination/headers", "siem_http_destination", "additional_headers", "name = \"test\"\ndestination_type = \"generic\"", `{ Authorization = "Bearer test", "X-Customer" = "acme" }`, nil},
		{"desktop_configuration_profile/normalized path", "desktop_configuration_profile", "mount_mappings", `name = "test"`, `{ W = "/Americas/" }`, func(value any) any {
			value.(map[string]any)["W"] = "Americas"
			return value
		}},
		{"integration_centric_profile/blank name", "integration_centric_profile", "expected_remote_servers", `name = "test"`, `[{ server_type = "dropbox", name = " " }]`, func(value any) any {
			delete(value.([]any)[0].(map[string]any), "name")
			return value
		}},
		{"integration_centric_profile/numeric name", "integration_centric_profile", "expected_remote_servers", `name = "numeric integration name"`, `[{ server_type = "dropbox", name = 42 }]`, func(value any) any {
			value.([]any)[0].(map[string]any)["name"] = "42"
			return value
		}},
		{"expectation/null criteria", "expectation", "criteria", "name = \"test\"\npath = \"incoming\"\ntrigger = \"manual\"", `null`, func(value any) any {
			if value == nil {
				return map[string]any{}
			}
			return value
		}},
		{"expectation/required file default", "expectation", "criteria", "name = \"test\"\npath = \"incoming\"\ntrigger = \"manual\"", `{ required_files = { "report.csv" = {} } }`, func(value any) any {
			value.(map[string]any)["required_files"].(map[string]any)["report.csv"] = map[string]any{"count": map[string]any{"exact": 1}}
			return value
		}},
	} {
		t.Run(test.name, func(t *testing.T) {
			var lock sync.Mutex
			var stored map[string]any
			if test.model == "site" {
				stored = map[string]any{"id": float64(0), "bundle_watermark_value": map[string]any{}}
			}
			writes := 0
			server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, req *http.Request) {
				lock.Lock()
				defer lock.Unlock()
				w.Header().Set("Content-Type", "application/json")
				switch req.Method {
				case http.MethodPost, http.MethodPatch:
					var body map[string]any
					if err := json.NewDecoder(req.Body).Decode(&body); err != nil {
						t.Error(err)
						w.WriteHeader(400)
						return
					}
					if value, exists := body[test.property]; exists {
						_, encoded := value.(string)
						assert.False(t, encoded, "API requests must contain unwrapped JSON, not a JSON string")
					}
					if stored == nil {
						stored = map[string]any{"id": float64(1)}
					}
					for key, value := range body {
						stored[key] = value
					}
					if test.normalize != nil {
						stored[test.property] = test.normalize(stored[test.property])
					}
					writes++
				case http.MethodDelete:
					stored = nil
					w.WriteHeader(http.StatusNoContent)
					return
				}
				if stored == nil {
					w.WriteHeader(http.StatusNotFound)
					fmt.Fprint(w, `{"error":"Not Found","http-code":404}`)
					return
				}
				if err := json.NewEncoder(w).Encode(stored); err != nil {
					t.Error(err)
				}
			}))
			defer server.Close()
			factories := map[string]func() (tfprotov6.ProviderServer, error){
				"files": providerserver.NewProtocol6WithError(&jsonTransitionTestProvider{filesProvider: filesProvider{version: "test"}, endpoint: server.URL}),
			}
			address := "files_" + test.model + ".example"
			config := func(encoded bool, read bool) string {
				value := test.value
				if encoded {
					value = "jsonencode(" + value + ")"
				}
				result := providerConfig + fmt.Sprintf("resource \"files_%s\" \"example\" {\n%s\n%s = %s\n}\n", test.model, test.base, test.property, value)
				if test.model == "site" {
					result += "import {\nto = files_site.example\nid = \"0\"\n}\n"
				}
				if read {
					selector := "id = " + address + ".id"
					if test.model == "site" {
						selector = "depends_on = [files_site.example]"
					}
					result += fmt.Sprintf("data \"files_%s\" \"example\" { %s }\noutput \"read_value\" {\nvalue = data.files_%s.example.%s\nsensitive = true\n}", test.model, selector, test.model, test.property)
				}
				return result
			}
			steps := []resource.TestStep{
				{Config: config(true, false)},
				{Config: config(true, false), PlanOnly: true},
				{Config: config(false, true)},
				{Config: config(false, true), PlanOnly: true},
				{ResourceName: address, ImportState: true, ImportStateVerify: true, ImportStateVerifyIgnore: []string{test.property}},
			}
			if test.model == "site" {
				steps = append(steps, resource.TestStep{Config: providerConfig + "removed {\nfrom = files_site.example\nlifecycle { destroy = false }\n}\n"})
			}
			resource.UnitTest(t, resource.TestCase{
				ProtoV6ProviderFactories: factories,
				Steps:                    steps,
			})
			require.GreaterOrEqual(t, writes, 1)
		})
	}
}
