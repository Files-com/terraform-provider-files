package lib

import (
	"context"
	"encoding/json"
	"errors"
	"io"
	"maps"
	"math/big"
	"net/http"
	"strconv"
	"strings"
	"testing"
	"testing/iotest"

	files_sdk "github.com/Files-com/files-sdk-go/v3"
	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/schema/validator"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// Native HCL literals use object/tuple/number types, not nested Dynamic values.
func nativeJSONValue(t *testing.T, value any) attr.Value {
	t.Helper()
	switch value := value.(type) {
	case map[string]any:
		attributes, attributeTypes := map[string]attr.Value{}, map[string]attr.Type{}
		for key, entry := range value {
			attributes[key] = nativeJSONValue(t, entry)
			attributeTypes[key] = attributes[key].Type(context.Background())
		}
		return types.ObjectValueMust(attributeTypes, attributes)
	case []any:
		values, valueTypes := []attr.Value{}, []attr.Type{}
		for _, entry := range value {
			item := nativeJSONValue(t, entry)
			values, valueTypes = append(values, item), append(valueTypes, item.Type(context.Background()))
		}
		return types.TupleValueMust(valueTypes, values)
	case float64:
		return types.NumberValue(big.NewFloat(value))
	case string:
		return types.StringValue(value)
	case bool:
		return types.BoolValue(value)
	case nil:
		return types.DynamicNull()
	default:
		t.Fatalf("unsupported test JSON value %T", value)
		return nil
	}
}

func TestJSONEncodingDiagnostics(t *testing.T) {
	for _, test := range []struct {
		name             string
		value            types.Dynamic
		warnings, errors int
	}{
		{"native", types.DynamicValue(nativeJSONValue(t, map[string]any{"Authorization": "test"})), 0, 0},
		{"JSON", types.DynamicValue(types.StringValue("{\"Authorization\":\"test\"}")), 1, 0},
		{"malformed JSON", types.DynamicValue(types.StringValue("{private-data")), 0, 1},
		{"null", types.DynamicNull(), 0, 0},
		{"unknown", types.DynamicUnknown(), 0, 0},
		{"unknown string", types.DynamicValue(types.StringUnknown()), 0, 0},
		{"null string", types.DynamicValue(types.StringNull()), 0, 0},
	} {
		t.Run(test.name, func(t *testing.T) {
			response := &validator.DynamicResponse{}
			DeprecatedJSONEncoding("files_siem_http_destination.additional_headers", "additional_headers = { Authorization = \"Bearer YOUR_TOKEN\" }", "March 1, 2027", nil, nil).ValidateDynamic(context.Background(), validator.DynamicRequest{Path: path.Root("additional_headers"), ConfigValue: test.value}, response)
			assert.Equal(t, test.warnings, response.Diagnostics.WarningsCount())
			assert.Equal(t, test.errors, response.Diagnostics.ErrorsCount())
			for _, diagnostic := range response.Diagnostics {
				assert.NotContains(t, diagnostic.Detail(), "private-data")
			}
			if test.warnings > 0 {
				assert.Equal(t, "Remove jsonencode(...) or replace the JSON string with native HCL. Keep the same keys and nesting. JSON strings stop being accepted on March 1, 2027, when the typed schema ships in a major provider release. No additional wrapper is needed.", response.Diagnostics[0].Detail())
			}

			// A property that stays Dynamic rejects the same malformed JSON, with no deprecation warning and no migration text.
			response = &validator.DynamicResponse{}
			ValidJSONEncoding().ValidateDynamic(context.Background(), validator.DynamicRequest{Path: path.Root("import_urls"), ConfigValue: test.value}, response)
			assert.Zero(t, response.Diagnostics.WarningsCount())
			assert.Equal(t, test.errors, response.Diagnostics.ErrorsCount())
			for _, diagnostic := range response.Diagnostics {
				assert.Equal(t, "Invalid JSON value", diagnostic.Summary())
				assert.NotContains(t, diagnostic.Detail(), "private-data")
				assert.NotContains(t, diagnostic.Detail(), "migrat")
			}
		})
	}
}

func TestJSONEncodingTypeGuidance(t *testing.T) {
	// Value types as the generator derives them from the published schemas of these properties.
	config := map[string]string{"user_ids": "list", "user_ids[]": "integer", "use_dedicated_ips": "boolean", "google_credentials": "object"}
	keyFile := `{"client_email":"events@example.com","private_key":"example-secret"}`
	deliveryPolicy := map[string]string{"max_attempts": "integer", "retry_delay_minutes": "list", "retry_delay_minutes[]": "integer"}
	aiFeatureAvailability := map[string]string{"in_app_ai_assistant.site_admins": "boolean", "in_app_ai_assistant.all_users": "boolean", "in_app_ai_assistant.group_ids": "list", "in_app_ai_assistant.group_ids[]": "integer"}
	criteria := map[string]string{"extensions": "list", "required_files.*.count.exact": "integer"}
	members := map[string]string{"": "list"}
	native := func(value any) types.Dynamic { return types.DynamicValue(nativeJSONValue(t, value)) }
	encoded := func(value string) types.Dynamic { return types.DynamicValue(types.StringValue(value)) }
	unknownID := types.TupleValueMust([]attr.Type{types.NumberType}, []attr.Value{types.NumberUnknown()})
	const typedSchema = " The typed schema that ships on March 1, 2027 uses these types."
	aiRoles := []string{"in_app_ai_assistant.site_admins", "in_app_ai_assistant.all_users"}
	policyAIRoles := []string{"ai_feature_availability.in_app_ai_assistant.site_admins"}
	const requiredNull = "The typed schema that ships on March 1, 2027 requires these fields and treats null as missing."
	unknownMember := types.ObjectValueMust(map[string]attr.Type{"email": types.StringType}, map[string]attr.Value{"email": types.StringUnknown()})
	unknownMembers := types.TupleValueMust([]attr.Type{unknownMember.Type(context.Background())}, []attr.Value{unknownMember})
	unknownRole := types.ObjectValueMust(map[string]attr.Type{"site_admins": types.BoolType, "all_users": types.BoolType}, map[string]attr.Value{"site_admins": types.BoolUnknown(), "all_users": types.BoolValue(false)})
	for _, test := range []struct {
		name, rule    string
		valueTypes    map[string]string
		requiredPaths []string
		value         types.Dynamic
		summary       string
		detail        string
		errors        int
	}{
		{name: "integer IDs and a boolean", rule: "files_event_target.config", valueTypes: config, value: native(map[string]any{"user_ids": []any{float64(101), float64(102)}, "use_dedicated_ips": true})},
		{name: "numeric string IDs", rule: "files_event_target.config", valueTypes: config, value: native(map[string]any{"user_ids": []any{"101", "102"}}),
			summary: "Typed values expected: files_event_target.config", detail: "Write user_ids as integers, for example user_ids = [123]." + typedSchema},
		{name: "integer and numeric string IDs", rule: "files_event_target.config", valueTypes: config, value: native(map[string]any{"user_ids": []any{float64(101), "102"}}),
			summary: "Typed values expected: files_event_target.config", detail: "Write user_ids as integers, for example user_ids = [123]." + typedSchema},
		{name: "fraction ID", rule: "files_event_target.config", valueTypes: config, value: native(map[string]any{"user_ids": []any{101.5}}),
			summary: "Typed values expected: files_event_target.config", detail: "Write user_ids as integers, for example user_ids = [123]." + typedSchema},
		{name: "string IDs and a string boolean, next to a secret", rule: "files_event_target.config", valueTypes: config, value: native(map[string]any{"user_ids": []any{"101"}, "use_dedicated_ips": "true", "aws_credentials": map[string]any{"secret_access_key": "example-secret"}}),
			summary: "Typed values expected: files_event_target.config", detail: "Write user_ids as integers, for example user_ids = [123]. Write use_dedicated_ips as true or false." + typedSchema},
		{name: "unknown ID", rule: "files_event_target.config", valueTypes: config, value: types.DynamicValue(types.ObjectValueMust(map[string]attr.Type{"user_ids": unknownID.Type(context.Background())}, map[string]attr.Value{"user_ids": unknownID}))},
		{name: "integer policy", rule: "files_event_target.delivery_policy", valueTypes: deliveryPolicy, value: native(map[string]any{"max_attempts": float64(3), "retry_delay_minutes": []any{float64(1), float64(5)}})},
		{name: "empty and null policy values", rule: "files_event_target.delivery_policy", valueTypes: deliveryPolicy, value: native(map[string]any{"max_attempts": nil, "retry_delay_minutes": []any{}})},
		{name: "numeric string attempts and one delay instead of a list", rule: "files_event_target.delivery_policy", valueTypes: deliveryPolicy, value: native(map[string]any{"max_attempts": "3", "retry_delay_minutes": float64(5)}),
			summary: "Typed values expected: files_event_target.delivery_policy", detail: "Write max_attempts as integers, for example max_attempts = 123. Write retry_delay_minutes as lists, for example retry_delay_minutes = [123]." + typedSchema},
		{name: "numeric string delays", rule: "files_event_target.delivery_policy", valueTypes: deliveryPolicy, value: native(map[string]any{"retry_delay_minutes": []any{"1", float64(5)}}),
			summary: "Typed values expected: files_event_target.delivery_policy", detail: "Write retry_delay_minutes as integers, for example retry_delay_minutes = [123]." + typedSchema},
		{name: "boolean options and null", rule: "files_site.ai_feature_availability", valueTypes: aiFeatureAvailability, value: native(map[string]any{"in_app_ai_assistant": map[string]any{"site_admins": true, "all_users": nil, "group_ids": []any{float64(1)}}})},
		{name: "string boolean options", rule: "files_site.ai_feature_availability", valueTypes: aiFeatureAvailability, value: native(map[string]any{"in_app_ai_assistant": map[string]any{"site_admins": "true", "all_users": "false"}}),
			summary: "Typed values expected: files_site.ai_feature_availability", detail: "Write in_app_ai_assistant.all_users and in_app_ai_assistant.site_admins as true or false." + typedSchema},
		{name: "required options set to null", rule: "files_site.ai_feature_availability", valueTypes: aiFeatureAvailability, requiredPaths: aiRoles, value: native(map[string]any{"in_app_ai_assistant": map[string]any{"site_admins": nil, "all_users": nil}}),
			summary: "Typed values expected: files_site.ai_feature_availability", detail: "Replace null in in_app_ai_assistant.site_admins and in_app_ai_assistant.all_users with true or false. " + requiredNull},
		{name: "required options set to false or omitted", rule: "files_site.ai_feature_availability", valueTypes: aiFeatureAvailability, requiredPaths: aiRoles, value: native(map[string]any{"in_app_ai_assistant": map[string]any{"site_admins": false}})},
		{name: "unknown required option", rule: "files_site.ai_feature_availability", valueTypes: aiFeatureAvailability, requiredPaths: aiRoles, value: types.DynamicValue(types.ObjectValueMust(map[string]attr.Type{"in_app_ai_assistant": unknownRole.Type(context.Background())}, map[string]attr.Value{"in_app_ai_assistant": unknownRole}))},
		{name: "string boolean option and a required null", rule: "files_site.ai_feature_availability", valueTypes: aiFeatureAvailability, requiredPaths: aiRoles, value: native(map[string]any{"in_app_ai_assistant": map[string]any{"site_admins": nil, "all_users": "false"}}),
			summary: "Typed values expected: files_site.ai_feature_availability", detail: "Write in_app_ai_assistant.all_users as true or false." + typedSchema + " Replace null in in_app_ai_assistant.site_admins with true or false. " + requiredNull},
		{name: "JSON-encoded required null", rule: "files_site.ai_feature_availability", valueTypes: aiFeatureAvailability, requiredPaths: aiRoles, value: encoded(`{"in_app_ai_assistant":{"site_admins":null,"all_users":false}}`),
			summary: "Deprecated JSON encoding: files_site.ai_feature_availability", detail: "No additional wrapper is needed. Replace null in in_app_ai_assistant.site_admins with true or false. " + requiredNull},
		{name: "required null in a root list item", rule: "files_share_group.members", valueTypes: members, requiredPaths: []string{"[].email"}, value: native([]any{map[string]any{"email": nil}}),
			summary: "Typed values expected: files_share_group.members", detail: "Replace null in email with values. " + requiredNull},
		{name: "unknown required field in a root list item", rule: "files_share_group.members", valueTypes: members, requiredPaths: []string{"[].email"}, value: types.DynamicValue(unknownMembers)},
		{name: "required null in a policy", rule: "files_child_site_management_policy.value", valueTypes: map[string]string{policyAIRoles[0]: "boolean"}, requiredPaths: policyAIRoles, value: native(map[string]any{"motd_text": nil, "ai_feature_availability": map[string]any{"in_app_ai_assistant": map[string]any{"site_admins": nil}}}),
			summary: "Typed values expected: files_child_site_management_policy.value", detail: "Replace null in ai_feature_availability.in_app_ai_assistant.site_admins with true or false. " + requiredNull},
		{name: "numeric string in a map value", rule: "files_expectation.criteria", valueTypes: criteria, value: native(map[string]any{"required_files": map[string]any{"report.csv": map[string]any{"count": map[string]any{"exact": "2"}}}}),
			summary: "Typed values expected: files_expectation.criteria", detail: "Write required_files.*.count.exact as integers, for example exact = 123." + typedSchema},
		{name: "one extension instead of a list", rule: "files_expectation.criteria", valueTypes: criteria, value: native(map[string]any{"extensions": "csv"}),
			summary: "Typed values expected: files_expectation.criteria", detail: "Write extensions as lists." + typedSchema},
		{name: "object instead of a root list", rule: "files_share_group.members", valueTypes: members, value: native(map[string]any{"email": "a@example.com"}),
			summary: "Typed values expected: files_share_group.members", detail: "Write members as lists." + typedSchema},
		{name: "JSON-encoded numeric strings", rule: "files_event_target.config", valueTypes: config, value: encoded(`{"user_ids":["101"],"use_dedicated_ips":"false"}`),
			summary: "Deprecated JSON encoding: files_event_target.config", detail: "Write user_ids as integers, for example user_ids = [123]. Write use_dedicated_ips as true or false."},
		{name: "JSON-encoded one delay", rule: "files_event_target.delivery_policy", valueTypes: deliveryPolicy, value: encoded(`{"retry_delay_minutes":5}`),
			summary: "Deprecated JSON encoding: files_event_target.delivery_policy", detail: "No additional wrapper is needed. Write retry_delay_minutes as lists, for example retry_delay_minutes = [123]."},
		{name: "credentials object", rule: "files_event_target.config", valueTypes: config, value: native(map[string]any{"google_credentials": map[string]any{"client_email": "events@example.com", "private_key": "example-secret"}})},
		{name: "credentials key file text", rule: "files_event_target.config", valueTypes: config, value: native(map[string]any{"google_credentials": keyFile}),
			summary: "Typed values expected: files_event_target.config", detail: "Write google_credentials as objects, for example google_credentials = jsondecode(...)." + typedSchema},
		{name: "JSON-encoded config with credentials key file text", rule: "files_event_target.config", valueTypes: config, value: encoded(`{"google_credentials":` + strconv.Quote(keyFile) + `}`),
			summary: "Deprecated JSON encoding: files_event_target.config", detail: "No additional wrapper is needed. Write google_credentials as objects, for example google_credentials = jsondecode(...)."},
		{name: "JSON-encoded integers", rule: "files_event_target.config", valueTypes: config, value: encoded(`{"user_ids":[101]}`),
			summary: "Deprecated JSON encoding: files_event_target.config", detail: "No additional wrapper is needed."},
		{name: "malformed JSON", rule: "files_event_target.config", valueTypes: config, value: encoded(`{"user_ids":["101"`), errors: 1},
		{name: "unknown", rule: "files_event_target.config", valueTypes: config, value: types.DynamicUnknown()},
		{name: "unknown string", rule: "files_event_target.config", valueTypes: config, value: types.DynamicValue(types.StringUnknown())},
		{name: "null", rule: "files_event_target.config", valueTypes: config, value: types.DynamicNull()},
	} {
		t.Run(test.name, func(t *testing.T) {
			response := &validator.DynamicResponse{}
			DeprecatedJSONEncoding(test.rule, "example = {}", "March 1, 2027", test.valueTypes, test.requiredPaths).ValidateDynamic(context.Background(), validator.DynamicRequest{Path: path.Root("value"), ConfigValue: test.value}, response)
			require.Equal(t, test.errors, response.Diagnostics.ErrorsCount())
			if test.summary == "" {
				assert.Zero(t, response.Diagnostics.WarningsCount())
				return
			}
			require.Equal(t, 1, response.Diagnostics.WarningsCount(), "one warning even when the value is also JSON-encoded")
			warning := response.Diagnostics.Warnings()[0]
			assert.Equal(t, test.summary, warning.Summary())
			if strings.HasPrefix(test.summary, "Typed") {
				assert.Equal(t, test.detail, warning.Detail())
			} else {
				assert.True(t, strings.HasSuffix(warning.Detail(), test.detail), warning.Detail())
			}
			for _, configured := range []string{"101", "102", "example-secret", "report.csv", "csv", "a@example.com"} {
				assert.NotContains(t, warning.Detail(), configured, "diagnostics must not repeat configured values or map keys")
			}
		})
	}
}

func TestJSONEncodingNestingGuidance(t *testing.T) {
	const typedExample = "value = { folder_behaviors = [{ file_expiration = { value = { days_to_retain = 30 } } }] }"
	valueTypes := map[string]string{"welcome_email_enabled": "boolean",
		"folder_behaviors[].file_expiration.value.days_to_retain": "integer", "folder_behaviors[].file_expiration.value.delete_empty_folders": "boolean",
		"folder_behaviors[behavior=file_expiration].value.days_to_retain": "integer", "folder_behaviors[behavior=file_expiration].value.delete_empty_folders": "boolean"}
	valueValidator := DeprecatedJSONEncoding("files_child_site_management_policy.value", typedExample, "March 1, 2027", valueTypes, nil, "folder_behaviors[].behavior")
	legacy := []any{map[string]any{"behavior": "file_expiration", "value": map[string]any{"days_to_retain": float64(45)}}}
	typed := []any{map[string]any{"file_expiration": map[string]any{"value": map[string]any{"days_to_retain": float64(45)}}}}
	typedStrings := []any{map[string]any{"file_expiration": map[string]any{"value": map[string]any{"days_to_retain": "45", "delete_empty_folders": "true"}}}}
	legacyString := []any{map[string]any{"behavior": "file_expiration", "value": map[string]any{"days_to_retain": "45"}}}
	const nesting = "Move each folder_behaviors item under the name in its behavior field, as in this example: " + typedExample + ". Both nestings work now."
	for _, test := range []struct {
		name, summary, guidance string
		value                   types.Dynamic
		nesting                 bool
	}{
		{name: "settings only", value: types.DynamicValue(nativeJSONValue(t, map[string]any{"welcome_email_enabled": true}))},
		{name: "no folder behaviors", value: types.DynamicValue(nativeJSONValue(t, map[string]any{"folder_behaviors": []any{}}))},
		{name: "typed nesting", value: types.DynamicValue(nativeJSONValue(t, map[string]any{"folder_behaviors": typed}))},
		{name: "old nesting", value: types.DynamicValue(nativeJSONValue(t, map[string]any{"folder_behaviors": legacy})), summary: "Deprecated nesting: files_child_site_management_policy.value", nesting: true},
		{name: "old nesting and a string boolean", value: types.DynamicValue(nativeJSONValue(t, map[string]any{"welcome_email_enabled": "false", "folder_behaviors": legacy})), summary: "Deprecated nesting: files_child_site_management_policy.value", nesting: true},
		{name: "typed nesting and a string boolean", value: types.DynamicValue(nativeJSONValue(t, map[string]any{"welcome_email_enabled": "false", "folder_behaviors": typed})), summary: "Typed values expected: files_child_site_management_policy.value"},
		{name: "typed nesting with string behavior settings", value: types.DynamicValue(nativeJSONValue(t, map[string]any{"folder_behaviors": typedStrings})), summary: "Typed values expected: files_child_site_management_policy.value",
			guidance: "Write folder_behaviors.file_expiration.value.days_to_retain as integers, for example days_to_retain = 123. Write folder_behaviors.file_expiration.value.delete_empty_folders as true or false."},
		{name: "old nesting with a string behavior setting", value: types.DynamicValue(nativeJSONValue(t, map[string]any{"folder_behaviors": legacyString})), summary: "Deprecated nesting: files_child_site_management_policy.value", nesting: true,
			guidance: "Write folder_behaviors.value.days_to_retain as integers, for example days_to_retain = 123."},
		{name: "JSON-encoded old nesting", value: types.DynamicValue(types.StringValue(`{"folder_behaviors":[{"behavior":"file_expiration","value":{"days_to_retain":45}}]}`)), summary: "Deprecated JSON encoding: files_child_site_management_policy.value", nesting: true},
		{name: "JSON-encoded typed nesting", value: types.DynamicValue(types.StringValue(`{"folder_behaviors":[{"file_expiration":{"value":{"days_to_retain":45}}}]}`)), summary: "Deprecated JSON encoding: files_child_site_management_policy.value"},
		{name: "unknown", value: types.DynamicUnknown()},
	} {
		t.Run(test.name, func(t *testing.T) {
			response := &validator.DynamicResponse{}
			valueValidator.ValidateDynamic(context.Background(), validator.DynamicRequest{Path: path.Root("value"), ConfigValue: test.value}, response)
			require.False(t, response.Diagnostics.HasError(), response.Diagnostics)
			if test.summary == "" {
				assert.Zero(t, response.Diagnostics.WarningsCount())
				return
			}
			require.Equal(t, 1, response.Diagnostics.WarningsCount())
			warning := response.Diagnostics.Warnings()[0]
			assert.Equal(t, test.summary, warning.Summary())
			assert.Equal(t, test.nesting, strings.Contains(warning.Detail(), nesting), warning.Detail())
			assert.Contains(t, warning.Detail(), test.guidance)
			assert.Equal(t, strings.HasPrefix(test.summary, "Deprecated JSON") && !test.nesting, strings.Contains(warning.Detail(), "Keep the same keys and nesting"), "only values without old nesting keep their nesting")
			assert.NotContains(t, warning.Detail(), "45", "diagnostics must not repeat configured values")
		})
	}
}

// During deprecation a Dynamic value with a list union accepts items in the API nesting and in the typed nesting. Both
// send API items, both read back in the configured nesting, and the typed state upgrade accepts both.
func TestDiscriminatedUnionNestingDuringDeprecation(t *testing.T) {
	ctx, attributePath := context.Background(), path.Root("value")
	names, discriminator := []string{"folder_behaviors"}, "behavior"
	variants := []JSONSchemaVariant{{Name: "malware_scanning", Value: "malware_scanning"}, {Name: "file_expiration", Value: "file_expiration"}, {Name: "limit_file_extensions", Value: "limit_file_extensions"}, {Name: "limit_file_regex", Value: "limit_file_regex"}}
	values := map[string]any{"malware_scanning": map[string]any{}, "file_expiration": map[string]any{"days_to_retain": float64(30)}, "limit_file_extensions": map[string]any{"extensions": []any{"txt"}}, "limit_file_regex": []any{"^report"}}
	order := []string{"malware_scanning", "file_expiration", "limit_file_extensions", "limit_file_regex"}
	items := func(typed ...bool) []any {
		result := []any{}
		for index, behavior := range order {
			if typed[index%len(typed)] {
				result = append(result, map[string]any{behavior: map[string]any{"value": values[behavior]}})
			} else {
				result = append(result, map[string]any{"behavior": behavior, "value": values[behavior]})
			}
		}
		return result
	}
	// An object-typed Terraform variable sets the unused variants of each typed item to null.
	withUnusedVariants := []any{}
	for _, item := range items(true) {
		object := maps.Clone(item.(map[string]any))
		for _, behavior := range order {
			if _, selected := object[behavior]; !selected {
				object[behavior] = nil
			}
		}
		withUnusedVariants = append(withUnusedVariants, object)
	}
	request := map[string]any{"welcome_email_enabled": true, "folder_behaviors": items(false)}
	upgraded := map[string]any{"welcome_email_enabled": true, "folder_behaviors": items(true)}
	// Rails stores the behavior values with their server defaults.
	stored := map[string]any{"welcome_email_enabled": true, "folder_behaviors": []any{
		map[string]any{"behavior": "malware_scanning", "value": map[string]any{}},
		map[string]any{"behavior": "file_expiration", "value": map[string]any{"days_to_retain": float64(30), "delete_empty_folders": false}},
		map[string]any{"behavior": "limit_file_extensions", "value": map[string]any{"extensions": []any{"txt"}, "mode": "whitelist"}},
		map[string]any{"behavior": "limit_file_regex", "value": []any{"^report"}},
	}}
	encode := func(value any) types.Dynamic {
		text, err := json.Marshal(value)
		require.NoError(t, err)
		return types.DynamicValue(types.StringValue(string(text)))
	}
	for _, test := range []struct {
		name       string
		configured types.Dynamic
		upgraded   map[string]any
	}{
		{"API nesting", types.DynamicValue(nativeJSONValue(t, map[string]any{"welcome_email_enabled": true, "folder_behaviors": items(false)})), upgraded},
		{"typed nesting", types.DynamicValue(nativeJSONValue(t, map[string]any{"welcome_email_enabled": true, "folder_behaviors": items(true)})), upgraded},
		{"typed nesting with unused variants null", types.DynamicValue(nativeJSONValue(t, map[string]any{"welcome_email_enabled": true, "folder_behaviors": withUnusedVariants})), map[string]any{"welcome_email_enabled": true, "folder_behaviors": withUnusedVariants}},
		{"mixed nesting", types.DynamicValue(nativeJSONValue(t, map[string]any{"welcome_email_enabled": true, "folder_behaviors": items(true, false)})), upgraded},
		{"JSON-encoded API nesting", encode(map[string]any{"welcome_email_enabled": true, "folder_behaviors": items(false)}), upgraded},
		{"JSON-encoded typed nesting", encode(map[string]any{"welcome_email_enabled": true, "folder_behaviors": items(true)}), upgraded},
	} {
		t.Run(test.name, func(t *testing.T) {
			sent, diags := JSONValueToAPI(ctx, attributePath, test.configured)
			require.False(t, diags.HasError(), diags)
			sent, diags = UnwrapDiscriminatedUnionAtPath(ctx, attributePath, sent, names, discriminator, variants)
			require.False(t, diags.HasError(), diags)
			assert.Equal(t, request, sent, "both nestings send API items")

			response, diags := WrapDiscriminatedUnionLikeConfigAtPath(ctx, attributePath, stored, test.configured, names, discriminator, variants)
			require.False(t, diags.HasError(), diags)
			state, diags := APIToDynamicJSON(ctx, attributePath, response, test.configured, nil, nil, "empty_object")
			require.False(t, diags.HasError(), diags)
			assert.True(t, test.configured.Equal(state), "apply keeps the configured nesting, types and omitted defaults: %v", state)
			response, diags = WrapDiscriminatedUnionLikeConfigAtPath(ctx, attributePath, stored, state, names, discriminator, variants)
			require.False(t, diags.HasError(), diags)
			refresh, diags := APIToDynamicJSON(ctx, attributePath, response, state, nil, nil, "empty_object")
			require.False(t, diags.HasError(), diags)
			assert.True(t, state.Equal(refresh))

			prior, diags := JSONValueToAPI(ctx, attributePath, state)
			require.False(t, diags.HasError(), diags)
			prior, diags = WrapDiscriminatedUnionAtPath(ctx, attributePath, prior, names, discriminator, variants)
			require.False(t, diags.HasError(), diags)
			assert.Equal(t, test.upgraded, prior, "the typed state upgrade wraps old items once and keeps typed items")
		})
	}

	for _, configured := range []types.Dynamic{types.DynamicNull(), types.DynamicUnknown()} {
		response, diags := WrapDiscriminatedUnionLikeConfigAtPath(ctx, attributePath, stored, configured, names, discriminator, variants)
		require.False(t, diags.HasError(), diags)
		assert.Equal(t, stored, response, "imports and unknown configuration read the API nesting")
	}
	selected := map[string]any{"value": map[string]any{"days_to_retain": float64(30)}}
	assert.True(t, typedDiscriminatedItem(map[string]any{"file_expiration": selected, "malware_scanning": nil}, variants))
	assert.False(t, typedDiscriminatedItem(map[string]any{"file_expiration": selected, "malware_scanning": map[string]any{"value": map[string]any{}}}, variants), "two selected variants")
	assert.False(t, typedDiscriminatedItem(map[string]any{"file_expiration": selected, "other": true}, variants), "a key that is not a variant")
	assert.False(t, typedDiscriminatedItem(map[string]any{"file_expiration": nil}, variants), "no selected variant")
	_, diags := WrapDiscriminatedUnionLikeConfigAtPath(ctx, attributePath, stored, types.DynamicValue(types.StringValue(`{"folder_behaviors":`)), names, discriminator, variants)
	assert.True(t, diags.HasError(), "malformed configured JSON is reported")
	_, diags = WrapDiscriminatedUnionLikeConfigAtPath(ctx, attributePath, stored, types.DynamicValue(nativeJSONValue(t, []any{"not an object"})), names, discriminator, variants)
	assert.True(t, diags.HasError(), "a configured value that is not an object is reported")
	masked := map[string]any{"welcome_email_enabled": nil, "folder_behaviors": nil}
	response, diags := WrapDiscriminatedUnionLikeConfigAtPath(ctx, attributePath, masked, types.DynamicValue(nativeJSONValue(t, upgraded)), names, discriminator, variants)
	require.False(t, diags.HasError(), diags)
	assert.Equal(t, masked, response, "values hidden from a Child Site stay null")
}

func TestDynamicJSONNormalizationAndDrift(t *testing.T) {
	for _, test := range []struct {
		name, configured, response, expected, normalization string
		maps, writeOnly                                     []string
	}{
		{"header arbitrary keys", `{"X-Customer":"a"}`, `{"X-Customer":"b","X-New":"c"}`, `{"X-Customer":"b","X-New":"c"}`, "", []string{""}, nil},
		{"map remote deletion", `{"department":"finance"}`, `{}`, `{}`, "", []string{""}, nil},
		{"object remote deletion", `{"header_name":"Authorization"}`, `{}`, `{"header_name":null}`, "", nil, nil},
		{"only preserve write-only", `{"password":"test","name":"old"}`, `{}`, `{"password":"test","name":null}`, "", nil, []string{"password"}},
		{"ignore object defaults", `{"count":{"exact":1}}`, `{"count":{"exact":1},"extensions":[]}`, `{"count":{"exact":1}}`, "", nil, nil},
		{"nested map defaults", `{"required_files":{"report.csv":{}}}`, `{"required_files":{"report.csv":{"count":{"exact":1}}}}`, `{"required_files":{"report.csv":{}}}`, "", []string{"required_files"}, nil},
		{"nested map deletion", `{"required_files":{"report.csv":{}}}`, `{"required_files":{}}`, `{"required_files":{}}`, "", []string{"required_files"}, nil},
		{"mount normalization", `{" W ":" /Americas/ "}`, `{"W":"Americas"}`, `{" W ":" /Americas/ "}`, "mount_mappings", []string{""}, nil},
		{"Rails path normalization", `{"W":"Reports/.././Inbox"}`, `{"W":"Reports/Inbox"}`, `{"W":"Reports/.././Inbox"}`, "mount_mappings", []string{""}, nil},
		{"mount normalization with drift", `{" W ":" /Americas/ ","Z":"Europe"}`, `{"W":"Americas","Z":"Asia"}`, `{" W ":" /Americas/ ","Z":"Asia"}`, "mount_mappings", []string{""}, nil},
		{"trimmed optional name", `[{"server_type":" dropbox ","name":" " }]`, `[{"server_type":"dropbox"}]`, `[{"server_type":" dropbox ","name":" " }]`, "expected_remote_servers", nil, nil},
		{"trim with drift", `[{"server_type":" dropbox ","name":"Old"}]`, `[{"server_type":"dropbox","name":"New"}]`, `[{"server_type":" dropbox ","name":"New"}]`, "expected_remote_servers", nil, nil},
		{"empty array", `[]`, `[]`, `[]`, "", nil, nil},
		{"watermark numeric string", `{"transparency":"25"}`, `{"transparency":25}`, `{"transparency":"25"}`, "watermark", nil, nil},
		{"watermark with drift", `{"transparency":"25","gravity":"Center"}`, `{"transparency":25,"gravity":"South"}`, `{"transparency":"25","gravity":"South"}`, "watermark", nil, nil},
		{"watermark numeric drift", `{"transparency":"25"}`, `{"transparency":30}`, `{"transparency":30}`, "watermark", nil, nil},
		{"boolean drift", `{"flag":"true"}`, `{"flag":false}`, `{"flag":false}`, "", []string{""}, nil},
		{"string drift", `{"flag":true}`, `{"flag":"false"}`, `{"flag":"false"}`, "", []string{""}, nil},
		{"numeric integration name", `[{"server_type":"dropbox","name":42}]`, `[{"server_type":"dropbox","name":"42"}]`, `[{"server_type":"dropbox","name":42}]`, "expected_remote_servers", nil, nil},
		{"boolean integration name", `[{"server_type":"dropbox","name":false}]`, `[{"server_type":"dropbox","name":"false"}]`, `[{"server_type":"dropbox","name":false}]`, "expected_remote_servers", nil, nil},
		{"watermark blank", `{"transparency":""}`, `{"transparency":""}`, `{"transparency":""}`, "watermark", nil, nil},
		{"Google credentials object without the key in the response", `{"project_id":"example","google_credentials":{"client_email":"events@example.com","private_key":"example-key"}}`, `{"project_id":"example","google_credentials":{"client_email":"events@example.com"}}`, `{"project_id":"example","google_credentials":{"client_email":"events@example.com","private_key":"example-key"}}`, "event_target_config", nil, []string{"aws_credentials.secret_access_key", "google_credentials.private_key"}},
		{"omitted event target secrets", `{"arn":"arn:aws:sns:us-east-1:123456789012:events","aws_credentials":{"access_key_id":"AKIAEXAMPLE","secret_access_key":"example-secret"}}`, `{"arn":"arn:aws:sns:us-east-1:123456789012:events","aws_credentials":{"access_key_id":"AKIAEXAMPLE"}}`, `{"arn":"arn:aws:sns:us-east-1:123456789012:events","aws_credentials":{"access_key_id":"AKIAEXAMPLE","secret_access_key":"example-secret"}}`, "event_target_config", nil, []string{"aws_credentials.secret_access_key", "google_credentials.private_key"}},
		{"Google credentials JSON string", `{"project_id":"example","google_credentials":"{\"client_email\":\"events@example.com\",\"private_key\":\"example-key\"}"}`, `{"project_id":"example","google_credentials":{"client_email":"events@example.com"}}`, `{"project_id":"example","google_credentials":"{\"client_email\":\"events@example.com\",\"private_key\":\"example-key\"}"}`, "event_target_config", nil, []string{"google_credentials.private_key"}},
		{"Google credentials drift", `{"google_credentials":"{\"client_email\":\"events@example.com\",\"private_key\":\"example-key\"}"}`, `{"google_credentials":{"client_email":"other@example.com"}}`, `{"google_credentials":{"client_email":"other@example.com"}}`, "event_target_config", nil, []string{"google_credentials.private_key"}},
		{"batch interval numeric string", `{"batch_interval":"900"}`, `{"batch_interval":900}`, `{"batch_interval":"900"}`, "event_target_delivery_policy", nil, nil},
		{"blank batch interval default", `{"batch_interval":"","max_attempts":3}`, `{"batch_interval":600,"max_attempts":3}`, `{"batch_interval":"","max_attempts":3}`, "event_target_delivery_policy", nil, nil},
		{"batch interval drift", `{"batch_interval":"900"}`, `{"batch_interval":1200}`, `{"batch_interval":1200}`, "event_target_delivery_policy", nil, nil},
		// Rails returns every form field key, and clears options for text fields.
		{"form field defaults and cleared options", `[{"label":"Name","field_type":"text","help_text":"","options_for_select":["Red"]}]`, `[{"id":7,"label":"Name","required":true,"help_text":"","field_type":"text","options_for_select":null,"default_option":null,"form_field_set_id":3}]`, `[{"label":"Name","field_type":"text","help_text":"","options_for_select":["Red"]}]`, "form_fields", nil, nil},
		{"form field empty help text drift", `[{"label":"Name","field_type":"text","help_text":""}]`, `[{"id":7,"label":"Name","required":true,"help_text":null,"field_type":"text","options_for_select":null,"default_option":null,"form_field_set_id":3}]`, `[{"label":"Name","field_type":"text","help_text":null}]`, "form_fields", nil, nil},
		{"form field boolean string", `[{"label":"Name","field_type":"text","required":"false"},{"label":"Team","field_type":"text","required":"1"}]`, `[{"id":7,"label":"Name","required":false,"help_text":null,"field_type":"text","options_for_select":null,"default_option":null,"form_field_set_id":3},{"id":8,"label":"Team","required":true,"help_text":null,"field_type":"text","options_for_select":null,"default_option":null,"form_field_set_id":3}]`, `[{"label":"Name","field_type":"text","required":"false"},{"label":"Team","field_type":"text","required":"1"}]`, "form_fields", nil, nil},
		{"form field boolean string drift", `[{"label":"Name","field_type":"text","required":"false"}]`, `[{"id":7,"label":"Name","required":true,"help_text":null,"field_type":"text","options_for_select":null,"default_option":null,"form_field_set_id":3}]`, `[{"label":"Name","field_type":"text","required":true}]`, "form_fields", nil, nil},
		{"form field drift", `[{"label":"Name","field_type":"text"}]`, `[{"id":8,"label":"Full name","field_type":"text"}]`, `[{"label":"Full name","field_type":"text"}]`, "form_fields", nil, nil},
		{"free-form additions with fixed defaults", `[{"name":"a.json","headers":{"X-A":"a"},"content":{"nested":{"a":1}}}]`, `[{"name":"a.json","method":"get","headers":{"X-A":"a","X-B":"b"},"content":{"nested":{"a":1,"b":2}}}]`, `[{"name":"a.json","headers":{"X-A":"a","X-B":"b"},"content":{"nested":{"a":1,"b":2}}}]`, "", []string{"headers.*", "content.*"}, nil},
		{"free-form root", `{"group_by":"user"}`, `{"group_by":"user","extra":{"added":true}}`, `{"group_by":"user","extra":{"added":true}}`, "empty_object", []string{"*"}, nil},
		{"free-form secret", `{"google_credentials":{"client_email":"events@example.com","private_key":"example-key"}}`, `{"google_credentials":{"client_email":"events@example.com","client_id":"added"}}`, `{"google_credentials":{"client_email":"events@example.com","client_id":"added","private_key":"example-key"}}`, "event_target_config", []string{"google_credentials.*"}, []string{"google_credentials.private_key"}},
	} {
		t.Run(test.name, func(t *testing.T) {
			var configured, response, expected any
			require.NoError(t, json.Unmarshal([]byte(test.configured), &configured))
			require.NoError(t, json.Unmarshal([]byte(test.response), &response))
			require.NoError(t, json.Unmarshal([]byte(test.expected), &expected))
			for _, encoded := range []bool{false, true} {
				prior := types.DynamicValue(nativeJSONValue(t, configured))
				if encoded {
					prior = types.DynamicValue(types.StringValue(test.configured))
				}
				state, diags := APIToDynamicJSON(context.Background(), path.Root("value"), response, prior, test.maps, test.writeOnly, test.normalization)
				require.False(t, diags.HasError(), diags)
				actual, diags := JSONValueToAPI(context.Background(), path.Root("value"), state)
				require.False(t, diags.HasError(), diags)
				if encoded && (test.name == "object remote deletion" || test.name == "only preserve write-only") {
					// JSON strings have no fixed object type requiring removed fields to remain null.
					object := expected.(map[string]any)
					for key, value := range object {
						if value == nil {
							delete(object, key)
						}
					}
				}
				assert.Equal(t, expected, actual)
				refreshed, diags := APIToDynamicJSON(context.Background(), path.Root("value"), response, state, test.maps, test.writeOnly, test.normalization)
				require.False(t, diags.HasError(), diags)
				assert.True(t, state.Equal(refreshed), "refresh should preserve the selected representation")
			}
		})
	}
}

func TestDynamicJSONMetadataUpdate(t *testing.T) {
	for _, value := range []attr.Value{types.DynamicValue(nativeJSONValue(t, map[string]any{"department": "legal"})), types.DynamicValue(types.StringValue("{\"department\":\"legal\"}")), types.MapValueMust(types.StringType, map[string]attr.Value{"department": types.StringValue("legal")})} {
		result, diags := JSONValueToAPI(context.Background(), path.Root("custom_metadata"), value)
		require.False(t, diags.HasError(), diags)
		assert.Equal(t, map[string]any{"department": "legal"}, result)
	}
}

func TestJSONNullMaterializedAsStoredDefault(t *testing.T) {
	for _, test := range []struct {
		normalization, response string
		unchanged               bool
	}{
		{"watermark", `{}`, true},
		{"empty_object", `{}`, true},
		{"event_target_delivery_policy", `{}`, true},
		{"event_target_delivery_policy", `{"batch_interval":600}`, true},
		{"event_target_delivery_policy", `{"batch_interval":900}`, false},
		{"event_target_delivery_policy", `{"batch_interval":600,"max_attempts":3}`, false},
	} {
		t.Run(test.normalization+" "+test.response, func(t *testing.T) {
			var response any
			require.NoError(t, json.Unmarshal([]byte(test.response), &response))
			prior := types.DynamicValue(types.StringValue("null"))
			state, diags := APIToDynamicJSON(context.Background(), path.Root("value"), response, prior, nil, nil, test.normalization)
			require.False(t, diags.HasError(), diags)
			assert.Equal(t, test.unchanged, prior.Equal(state), state)
		})
	}
}

// Malformed JSON is left to DeprecatedJSONEncoding.
func TestRejectJSONNull(t *testing.T) {
	for _, test := range []struct {
		name   string
		value  types.Dynamic
		errors int
	}{
		{"JSON null", types.DynamicValue(types.StringValue(" null ")), 1},
		{"native null", types.DynamicNull(), 0},
		{"unknown", types.DynamicUnknown(), 0},
		{"JSON empty list", types.DynamicValue(types.StringValue("[]")), 0},
		{"native empty list", types.DynamicValue(nativeJSONValue(t, []any{})), 0},
		{"malformed JSON", types.DynamicValue(types.StringValue("{")), 0},
	} {
		t.Run(test.name, func(t *testing.T) {
			response := &validator.DynamicResponse{}
			RejectJSONNull("Set members to the list of members.").ValidateDynamic(context.Background(), validator.DynamicRequest{Path: path.Root("members"), ConfigValue: test.value}, response)
			assert.Equal(t, test.errors, response.Diagnostics.ErrorsCount())
		})
	}
}

func TestTypedListFromSDKObjects(t *testing.T) {
	objectType := types.ObjectType{AttrTypes: map[string]attr.Type{"server_type": types.StringType, "name": types.StringType}}
	value, diags := ToList(context.Background(), path.Root("expected_remote_servers"), []map[string]any{{"server_type": "dropbox"}}, types.ListNull(objectType))
	assert.False(t, diags.HasError(), diags)
	assert.Len(t, value.Elements(), 1)
	assert.Equal(t, types.StringNull(), value.Elements()[0].(types.Object).Attributes()["name"])
}

// responseValue decodes a property of a response body the way the generated reads do.
func responseValue(t *testing.T, body, name string) any {
	t.Helper()
	value, diags := ResponseProperty(path.Root(name), []byte(body), name)
	require.False(t, diags.HasError(), diags)
	return value
}

func TestTypedCurrentFormFields(t *testing.T) {
	ctx, attributePath := context.Background(), path.Root("form_fields")
	fieldType := types.ObjectType{AttrTypes: map[string]attr.Type{"id": types.Int64Type, "label": types.StringType, "field_type": types.StringType, "required": types.BoolType, "help_text": types.StringType}}
	body := `{"form_layout":[8],"form_fields":[{"id":7,"label":"Old","field_type":"text","required":true,"help_text":null},{"id":8,"label":"Name","field_type":"text","required":false,"help_text":null}]}`
	formFieldSet := files_sdk.FormFieldSet{FormLayout: []int64{8}}

	value, diags := ToList(ctx, attributePath, CurrentFormFields(formFieldSet, responseValue(t, body, "form_fields")), types.ListNull(fieldType))
	require.False(t, diags.HasError(), diags)
	assert.Equal(t, types.ListValueMust(fieldType, []attr.Value{types.ObjectValueMust(fieldType.AttrTypes, map[string]attr.Value{
		"id": types.Int64Value(8), "label": types.StringValue("Name"), "field_type": types.StringValue("text"), "required": types.BoolValue(false), "help_text": types.StringNull(),
	})}), value)

	value, diags = ToList(ctx, attributePath, responseValue(t, `{"form_fields":null}`, "form_fields"), types.ListNull(fieldType))
	require.False(t, diags.HasError(), diags)
	assert.True(t, value.IsNull())
}

// The response body tells an empty string from null, so reads from it keep both exactly, with or without configuration,
// as resources, data sources and imports need. A configured "" never hides a null the API returns.
func TestManagedReadsFromResponseBody(t *testing.T) {
	ctx := context.Background()
	body := `{"id":4,"form_layout":[2,1],"form_fields":[{"id":1,"label":"Name","help_text":""},{"id":2,"label":"Team","help_text":null},{"id":3,"label":"Old","help_text":"x"}],` +
		`"members":[{"email":"a@example.com","name":"","company":null}]}`
	formFieldSet := files_sdk.FormFieldSet{FormLayout: []int64{2, 1}}
	fieldType := types.ObjectType{AttrTypes: map[string]attr.Type{"id": types.Int64Type, "label": types.StringType, "help_text": types.StringType}}
	field := func(id int64, label string, helpText types.String) attr.Value {
		return types.ObjectValueMust(fieldType.AttrTypes, map[string]attr.Value{"id": types.Int64Value(id), "label": types.StringValue(label), "help_text": helpText})
	}
	configured := types.ListValueMust(fieldType, []attr.Value{field(2, "Team", types.StringValue("")), field(1, "Name", types.StringValue(""))})
	fields, diags := ToList(ctx, path.Root("form_fields"), CurrentFormFields(formFieldSet, responseValue(t, body, "form_fields")), configured)
	require.False(t, diags.HasError(), diags)
	assert.Equal(t, types.ListValueMust(fieldType, []attr.Value{field(2, "Team", types.StringNull()), field(1, "Name", types.StringValue(""))}), fields,
		"current fields in form_layout order, and the API's null for Team shows as drift from the configured empty string")

	all, diags := ToList(ctx, path.Root("form_fields"), responseValue(t, body, "form_fields"), types.ListNull(fieldType))
	require.False(t, diags.HasError(), diags)
	assert.Equal(t, types.ListValueMust(fieldType, []attr.Value{field(1, "Name", types.StringValue("")), field(2, "Team", types.StringNull()), field(3, "Old", types.StringValue("x"))}), all, "data sources keep historical fields")

	memberType := types.ObjectType{AttrTypes: map[string]attr.Type{"email": types.StringType, "name": types.StringType, "company": types.StringType}}
	members, diags := ToList(ctx, path.Root("members"), AlignMembersByEmail(ctx, responseValue(t, body, "members"), types.ListNull(memberType)), types.ListNull(memberType))
	require.False(t, diags.HasError(), diags)
	assert.Equal(t, types.ListValueMust(memberType, []attr.Value{types.ObjectValueMust(memberType.AttrTypes, map[string]attr.Value{
		"email": types.StringValue("a@example.com"), "name": types.StringValue(""), "company": types.StringNull(),
	})}), members)

	assert.Nil(t, responseValue(t, `{"id":4}`, "members"), "a missing property reads as null")
	assert.Nil(t, responseValue(t, `{"members":null}`, "members"))
	for _, invalid := range []string{"", "null", "[1]", `{"members":`} {
		value, diags := ResponseProperty(path.Root("members"), []byte(invalid), "members")
		assert.True(t, diags.HasError(), "%q is not a JSON object, so the read fails instead of guessing", invalid)
		assert.Nil(t, value)
	}
}

func TestCaptureResponseBody(t *testing.T) {
	var body []byte
	response := &http.Response{StatusCode: 200, Body: io.NopCloser(strings.NewReader(`{"id":1,"name":""}`))}
	response, err := files_sdk.BuildResponse(response, CaptureResponseBody(&body))
	require.NoError(t, err)
	assert.Equal(t, `{"id":1,"name":""}`, string(body))
	decoded, err := io.ReadAll(response.Body)
	require.NoError(t, err)
	assert.Equal(t, body, decoded, "the SDK still reads the whole body")

	failing := &http.Response{StatusCode: 200, Body: io.NopCloser(iotest.ErrReader(errors.New("connection reset")))}
	_, err = files_sdk.BuildResponse(failing, CaptureResponseBody(&body))
	assert.ErrorContains(t, err, "connection reset")
	assert.Equal(t, http.NoBody, failing.Body, "a failed read leaves no closed body behind")
}

func TestAlignMembersByEmail(t *testing.T) {
	ctx, attributePath := context.Background(), path.Root("members")
	// Rails returns name, company and email for every member, null when unset.
	body := `{"members":[{"email":"a@example.com","name":null,"company":""},{"email":"b@example.com","name":"B","company":null}]}`

	for _, test := range []struct {
		name, configured, expected, body string
	}{
		{"configured order", `[{"email":"b@example.com","name":"B"},{"email":"A@example.com","company":""}]`, `[{"email":"b@example.com","name":"B"},{"email":"A@example.com","company":""}]`, body},
		{"member removed remotely", `[{"email":"a@example.com"},{"email":"b@example.com"}]`, `[{"email":"b@example.com"}]`, `{"members":[{"email":"b@example.com","name":"B","company":null}]}`},
	} {
		t.Run(test.name, func(t *testing.T) {
			var configured, expected any
			require.NoError(t, json.Unmarshal([]byte(test.configured), &configured))
			require.NoError(t, json.Unmarshal([]byte(test.expected), &expected))
			for _, prior := range []types.Dynamic{types.DynamicValue(nativeJSONValue(t, configured)), types.DynamicValue(types.StringValue(test.configured))} {
				state, diags := APIToDynamicJSON(ctx, attributePath, AlignMembersByEmail(ctx, responseValue(t, test.body, "members"), prior), prior, nil, nil, "")
				require.False(t, diags.HasError(), diags)
				actual, diags := JSONValueToAPI(ctx, attributePath, state)
				require.False(t, diags.HasError(), diags)
				assert.Equal(t, expected, actual)
			}
		})
	}

	memberType := types.ObjectType{AttrTypes: map[string]attr.Type{"email": types.StringType, "name": types.StringType, "company": types.StringType}}
	member := func(email string, name, company types.String) attr.Value {
		return types.ObjectValueMust(memberType.AttrTypes, map[string]attr.Value{"email": types.StringValue(email), "name": name, "company": company})
	}
	prior := types.ListValueMust(memberType, []attr.Value{member("B@example.com", types.StringValue("B"), types.StringNull()), member("a@example.com", types.StringNull(), types.StringValue(""))})
	state, diags := ToList(ctx, attributePath, AlignMembersByEmail(ctx, responseValue(t, body, "members"), prior), prior)
	require.False(t, diags.HasError(), diags)
	assert.Equal(t, prior, state, "configured order and email spelling, with the API's exact values")
}

func TestDynamicJSONNumericStateDrift(t *testing.T) {
	for _, priorNumber := range []attr.Value{types.Int64Value(1), types.Float64Value(1)} {
		prior := types.DynamicValue(types.ObjectValueMust(map[string]attr.Type{"limit": priorNumber.Type(context.Background())}, map[string]attr.Value{"limit": priorNumber}))
		state, diags := APIToDynamicJSON(context.Background(), path.Root("value"), map[string]any{"limit": "2"}, prior, []string{""}, nil, "")
		require.False(t, diags.HasError(), diags)
		assert.Equal(t, types.StringValue("2"), state.UnderlyingValue().(types.Object).Attributes()["limit"])
	}
}

func TestDynamicJSONMapScalarDrift(t *testing.T) {
	ctx, attributePath := context.Background(), path.Root("bundle_watermark_value")
	prior := types.DynamicValue(types.MapValueMust(types.StringType, map[string]attr.Value{"transparency": types.StringValue("25")}))
	response := map[string]any{"transparency": float64(30)}

	state, diags := APIToDynamicJSON(ctx, attributePath, response, prior, nil, nil, "watermark")
	require.False(t, diags.HasError(), diags)
	actual, diags := JSONValueToAPI(ctx, attributePath, state)
	require.False(t, diags.HasError(), diags)
	assert.Equal(t, response, actual)
}

// A Child Site policy changed outside Terraform can hold a list where prior state held an object at the same position,
// as when limit_file_regex replaces file_expiration. The read returns the API value, so the plan shows the drift
// instead of failing.
func TestDynamicJSONKindDrift(t *testing.T) {
	ctx, attributePath := context.Background(), path.Root("value")
	behavior := func(value attr.Value) attr.Value {
		return types.ObjectValueMust(map[string]attr.Type{"behavior": types.StringType, "value": value.Type(ctx)}, map[string]attr.Value{"behavior": types.StringValue("file_expiration"), "value": value})
	}
	expiration := types.ObjectValueMust(map[string]attr.Type{"days_to_retain": types.NumberType}, map[string]attr.Value{"days_to_retain": types.NumberValue(big.NewFloat(30))})
	prior := types.DynamicValue(types.ObjectValueMust(map[string]attr.Type{"folder_behaviors": types.TupleType{ElemTypes: []attr.Type{behavior(expiration).Type(ctx)}}},
		map[string]attr.Value{"folder_behaviors": types.TupleValueMust([]attr.Type{behavior(expiration).Type(ctx)}, []attr.Value{behavior(expiration)})}))
	for name, response := range map[string]map[string]any{
		"list where prior has an object":   {"folder_behaviors": []any{map[string]any{"behavior": "limit_file_regex", "value": []any{"/E2E-.*/"}}}},
		"object where prior has a list":    {"folder_behaviors": map[string]any{"behavior": "limit_file_regex"}},
		"string where prior has an object": {"folder_behaviors": []any{map[string]any{"behavior": "malware_scanning", "value": "on"}}},
	} {
		t.Run(name, func(t *testing.T) {
			state, diags := APIToDynamicJSON(ctx, attributePath, response, prior, nil, nil, "empty_object")
			require.False(t, diags.HasError(), diags)
			actual, diags := JSONValueToAPI(ctx, attributePath, state)
			require.False(t, diags.HasError(), diags)
			assert.Equal(t, any(response), actual)
		})
	}
}

func TestDynamicJSONTypedCollectionScalarDrift(t *testing.T) {
	ctx, attributePath := context.Background(), path.Root("value")
	for name, prior := range map[string]attr.Value{
		"list": types.ListValueMust(types.StringType, []attr.Value{types.StringValue("25")}),
		"set":  types.SetValueMust(types.StringType, []attr.Value{types.StringValue("25")}),
	} {
		t.Run(name, func(t *testing.T) {
			state, diags := APIToDynamicJSON(ctx, attributePath, []any{float64(30)}, types.DynamicValue(prior), nil, nil, "")
			require.False(t, diags.HasError(), diags)
			actual, diags := JSONValueToAPI(ctx, attributePath, state)
			require.False(t, diags.HasError(), diags)
			assert.Equal(t, []any{float64(30)}, actual)
			unchanged, diags := APIToDynamicJSON(ctx, attributePath, []any{"25"}, types.DynamicValue(prior), nil, nil, "")
			require.False(t, diags.HasError(), diags)
			assert.True(t, types.DynamicValue(prior).Equal(unchanged))
		})
	}
}
