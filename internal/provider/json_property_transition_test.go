package provider

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"slices"
	"strings"
	"sync"
	"testing"

	files_sdk "github.com/Files-com/files-sdk-go/v3"
	"github.com/hashicorp/terraform-plugin-framework/provider"
	"github.com/hashicorp/terraform-plugin-framework/providerserver"
	"github.com/hashicorp/terraform-plugin-go/tfprotov6"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/hashicorp/terraform-plugin-testing/knownvalue"
	"github.com/hashicorp/terraform-plugin-testing/plancheck"
	"github.com/hashicorp/terraform-plugin-testing/statecheck"
	"github.com/hashicorp/terraform-plugin-testing/tfjsonpath"
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

// Optional API behavior for a lifecycle case. normalize runs after every write, like a Rails callback.
// remoteChange edits the stored record before a final apply whose plan must update the resource. checks run after the native-HCL apply.
type recordBehavior struct {
	normalize, remoteChange func(map[string]any)
	checks                  []statecheck.StateCheck
}

// Rails replaces form fields sent without IDs, keeps earlier definitions in the response, and lists current IDs in form_layout.
func formFieldResponses() func(map[string]any) {
	var history []any
	nextID := float64(10)
	return func(stored map[string]any) {
		fields, sent := stored["form_fields"].([]any)
		if !sent {
			return
		}
		current, layout := []any{}, []any{}
		for _, field := range fields {
			object := field.(map[string]any)
			if _, existing := object["id"]; existing {
				return
			}
			nextID++
			record := map[string]any{"id": nextID, "required": true, "form_field_set_id": float64(1)}
			for key, value := range object {
				record[key] = value
			}
			if record["field_type"] == "text" {
				delete(record, "options_for_select")
				delete(record, "default_option")
			}
			current, layout = append(current, record), append(layout, nextID)
		}
		stored["form_fields"], stored["form_layout"] = append(slices.Clone(history), current...), layout
		history = append(history, current...)
	}
}

func TestJSONPropertyResourceTransitions(t *testing.T) {
	for _, test := range []struct {
		name, model, property, base, value string
		normalize                          func(any) any
		record                             *recordBehavior
	}{
		{"event_target/omitted AWS secret", "event_target", "config", "name = \"test\"\ntarget_type = \"amazon_sns\"", `{ arn = "arn:aws:sns:us-east-1:123456789012:events", aws_credentials = { access_key_id = "AKIAEXAMPLE", secret_access_key = "example-secret", region = "us-east-1" } }`, func(value any) any {
			delete(value.(map[string]any)["aws_credentials"].(map[string]any), "secret_access_key")
			return value
		}, nil},
		{"event_target/Google key file string", "event_target", "config", "name = \"test\"\ntarget_type = \"google_pubsub\"", `{ project_id = "example", topic_id = "events", google_credentials = jsonencode({ client_email = "events@example.com", private_key = "example-key" }) }`, func(value any) any {
			config := value.(map[string]any)
			if text, ok := config["google_credentials"].(string); ok {
				var credentials map[string]any
				_ = json.Unmarshal([]byte(text), &credentials)
				config["google_credentials"] = credentials
			}
			delete(config["google_credentials"].(map[string]any), "private_key")
			return value
		}, nil},
		{"event_target/stored batch interval", "event_target", "delivery_policy", "name = \"test\"\ntarget_type = \"email\"\nconfig = { user_ids = [1] }", `{ batch_interval = "900", max_attempts = 3 }`, func(value any) any {
			value.(map[string]any)["batch_interval"] = float64(900)
			return value
		}, nil},
		{"event_target/workspace folder path", "event_target", "config", "name = \"test\"\ntarget_type = \"folder\"", `{ path = "events", format = "json" }`, nil, &recordBehavior{normalize: func(stored map[string]any) {
			stored["delivery_policy"] = map[string]any{"batch_interval": float64(600)}
		}}},
		{"event_target/workspace summary path", "event_target", "delivery_policy", "name = \"test\"\ntarget_type = \"slack_webhook\"\nconfig = { url = \"https://example.com/hook\" }", `{ delivery_mode = "subfolder_summary", summary_path = "Uploads", summary_delay = 30 }`, nil, nil},
		{"automation/legacy import URL", "automation", "import_urls", `automation = "import_file"`, `[{ name = "users.json", url = "https://example.com/users", method = "GET", headers = [], content_type = "application/json" }]`, nil, nil},
		{"automation/import URL drift", "automation", "import_urls", `automation = "import_file"`, `[{ name = "users.json", url = "https://example.com/users", method = "POST", headers = { "X-A" = "a" }, content = { nested = { a = 1 } } }]`, nil, &recordBehavior{remoteChange: func(stored map[string]any) {
			importURL := stored["import_urls"].([]any)[0].(map[string]any)
			importURL["headers"].(map[string]any)["X-B"] = "b"
			importURL["content"].(map[string]any)["nested"].(map[string]any)["b"] = float64(2)
		}}},
		{"scheduled_export/integer all_folders and max_depth drift", "scheduled_export", "export_options", "name = \"test\"\nexport_type = \"folder_size_audit\"", `{ all_folders = 1, max_depth = 2 }`, nil, &recordBehavior{remoteChange: func(stored map[string]any) {
			stored["export_options"].(map[string]any)["max_depth"] = float64(3)
		}}},
		{"site/string boolean AI options", "site", "ai_feature_availability", "", `{ in_app_ai_assistant = { site_admins = "true", workspace_admins = true, folder_admins = "false", all_users = false } }`, nil, nil},
		{"child_site_management_policy/null value", "child_site_management_policy", "value", `policy_type = "settings"`, `null`, func(value any) any {
			if value == nil {
				return map[string]any{}
			}
			return value
		}, nil},
		{"child_site_management_policy/behavior defaults", "child_site_management_policy", "value", `policy_type = "settings"`, `{ welcome_email_enabled = false, folder_behaviors = [{ behavior = "file_expiration", value = { days_to_retain = 30 } }] }`, func(value any) any {
			behavior := value.(map[string]any)["folder_behaviors"].([]any)[0].(map[string]any)["value"].(map[string]any)
			behavior["delete_empty_folders"] = false
			return value
		}, nil},
		{"form_field_set/replaced fields", "form_field_set", "form_fields", `title = "test"`, `[{ label = "Name", field_type = "text", options_for_select = ["Unused"] }, { label = "Team", field_type = "dropdown", options_for_select = ["Finance", "Sales"] }]`, nil, &recordBehavior{normalize: formFieldResponses(), checks: []statecheck.StateCheck{
			statecheck.ExpectKnownValue("files_form_field_set.example", tfjsonpath.New("form_fields"), knownvalue.ListSizeExact(2)),
			statecheck.ExpectKnownOutputValue("read_value", knownvalue.ListSizeExact(4)),
		}}},
		{"share_group/members by email", "share_group", "members", `name = "test"`, `[{ email = "b@example.com", name = "B" }, { email = "A@example.com", company = "" }]`, func(value any) any {
			members := value.([]any)
			slices.SortFunc(members, func(a, b any) int {
				return strings.Compare(strings.ToLower(a.(map[string]any)["email"].(string)), strings.ToLower(b.(map[string]any)["email"].(string)))
			})
			return members
		}, nil},
		{"site/numeric string watermark", "site", "bundle_watermark_value", "", `{ transparency = "25", gravity = "SouthWest" }`, func(value any) any {
			value.(map[string]any)["transparency"] = float64(25)
			return value
		}, nil},
		{"bundle/watermark value", "bundle", "watermark_value", `paths = ["reports"]`, `{ gravity = "Center", transparency = "25", dynamic_text = "Confidential" }`, func(value any) any {
			value.(map[string]any)["transparency"] = float64(25)
			return value
		}, nil},
		{"siem_http_destination/headers", "siem_http_destination", "additional_headers", "name = \"test\"\ndestination_type = \"generic\"", `{ Authorization = "Bearer test", "X-Customer" = "acme" }`, nil, nil},
		{"desktop_configuration_profile/normalized path", "desktop_configuration_profile", "mount_mappings", `name = "test"`, `{ W = "/Americas/" }`, func(value any) any {
			value.(map[string]any)["W"] = "Americas"
			return value
		}, nil},
		{"integration_centric_profile/blank name", "integration_centric_profile", "expected_remote_servers", `name = "test"`, `[{ server_type = "dropbox", name = " " }]`, func(value any) any {
			delete(value.([]any)[0].(map[string]any), "name")
			return value
		}, nil},
		{"integration_centric_profile/numeric name", "integration_centric_profile", "expected_remote_servers", `name = "numeric integration name"`, `[{ server_type = "dropbox", name = 42 }]`, func(value any) any {
			value.([]any)[0].(map[string]any)["name"] = "42"
			return value
		}, nil},
		{"expectation/null criteria", "expectation", "criteria", "name = \"test\"\npath = \"incoming\"\ntrigger = \"manual\"", `null`, func(value any) any {
			if value == nil {
				return map[string]any{}
			}
			return value
		}, nil},
		{"expectation/required file default", "expectation", "criteria", "name = \"test\"\npath = \"incoming\"\ntrigger = \"manual\"", `{ required_files = { "report.csv" = {} } }`, func(value any) any {
			value.(map[string]any)["required_files"].(map[string]any)["report.csv"] = map[string]any{"count": map[string]any{"exact": 1}}
			return value
		}, nil},
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
					if test.record != nil && test.record.normalize != nil {
						test.record.normalize(stored)
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
			native := resource.TestStep{Config: config(false, true)}
			if test.record != nil {
				native.ConfigStateChecks = test.record.checks
			}
			steps := []resource.TestStep{
				{Config: config(true, false)},
				{Config: config(true, false), PlanOnly: true},
				native,
				{Config: config(false, true), PlanOnly: true},
				{ResourceName: address, ImportState: true, ImportStateVerify: true, ImportStateVerifyIgnore: []string{test.property}},
			}
			if test.record != nil && test.record.remoteChange != nil {
				steps = append(steps, resource.TestStep{
					PreConfig: func() {
						lock.Lock()
						defer lock.Unlock()
						test.record.remoteChange(stored)
					},
					Config:           config(false, true),
					ConfigPlanChecks: resource.ConfigPlanChecks{PreApply: []plancheck.PlanCheck{plancheck.ExpectResourceAction(address, plancheck.ResourceActionUpdate)}},
				})
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

// Rails keeps a v2 Automation's definition when an update leaves it out. Removing definition from the configuration must
// plan no change, and a real change must still plan and apply an update.
func TestAutomationOmittedDefinitionPlan(t *testing.T) {
	var lock sync.Mutex
	stored := map[string]any{"id": float64(1)}
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, req *http.Request) {
		lock.Lock()
		defer lock.Unlock()
		switch req.Method {
		case http.MethodPost, http.MethodPatch:
			// Decoding into the stored map keeps every field that the request leaves out.
			assert.NoError(t, json.NewDecoder(req.Body).Decode(&stored))
		case http.MethodDelete:
			w.WriteHeader(http.StatusNoContent)
			return
		}
		w.Header().Set("Content-Type", "application/json")
		assert.NoError(t, json.NewEncoder(w).Encode(stored))
	}))
	defer server.Close()
	address := "files_automation.example"
	definition := `definition = { schema_version = 1, nodes = [{ trigger_manual = { id = "trigger" } }], edges = [] }`
	config := func(name, attributes string) string {
		return providerConfig + fmt.Sprintf("resource \"files_automation\" \"example\" {\nautomation = \"v2\"\nname = %q\n%s\n}\n", name, attributes)
	}
	resource.UnitTest(t, resource.TestCase{
		ProtoV6ProviderFactories: map[string]func() (tfprotov6.ProviderServer, error){
			"files": providerserver.NewProtocol6WithError(&jsonTransitionTestProvider{filesProvider: filesProvider{version: "test"}, endpoint: server.URL}),
		},
		Steps: []resource.TestStep{
			{Config: config("first", definition)},
			{Config: config("first", definition), PlanOnly: true},
			{Config: config("first", ""), PlanOnly: true},
			{
				Config:           config("second", ""),
				ConfigPlanChecks: resource.ConfigPlanChecks{PreApply: []plancheck.PlanCheck{plancheck.ExpectResourceAction(address, plancheck.ResourceActionUpdate)}},
				ConfigStateChecks: []statecheck.StateCheck{
					statecheck.ExpectKnownValue(address, tfjsonpath.New("name"), knownvalue.StringExact("second")),
					statecheck.ExpectKnownValue(address, tfjsonpath.New("definition").AtMapKey("nodes"), knownvalue.ListSizeExact(1)),
				},
			},
		},
	})
}
