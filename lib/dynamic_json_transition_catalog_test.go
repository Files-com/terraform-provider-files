package lib

import (
	"context"
	"encoding/json"
	"testing"

	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/schema/validator"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

const typedSchemaDate = "March 1, 2027"

func TestPublishedJSONTransitions(t *testing.T) {
	for _, test := range []struct {
		name, example, normalization string
		maps, writeOnly              []string
		computed                     bool
		typed                        attr.Type
	}{
		{name: "as2_partner.additional_http_headers", example: "{\"X-Partner\":\"acme\"}",
			maps:          []string{""},
			normalization: "",
			computed:      false,
			typed:         types.MapType{ElemType: types.StringType}},
		{name: "automation.value", example: "{\"limit\":\"1\"}",
			maps:          []string{""},
			normalization: "",
			computed:      false,
			typed:         types.MapType{ElemType: types.StringType}},
		{name: "bundle.watermark_value", example: "{\"gravity\":\"SouthWest\",\"max_height_or_width\":20,\"transparency\":25}",
			maps:          []string{},
			normalization: "watermark",
			computed:      false,
			typed:         types.ObjectType{AttrTypes: map[string]attr.Type{"gravity": types.StringType, "max_height_or_width": types.Int64Type, "transparency": types.Int64Type, "dynamic_text": types.StringType}}},
		{name: "bundle.requested_upload_slots", example: "[{\"name\":\"Photo ID\"}]",
			maps:          []string{},
			normalization: "",
			computed:      true,
			typed:         types.ListType{ElemType: types.ObjectType{AttrTypes: map[string]attr.Type{"name": types.StringType}}}},
		{name: "desktop_configuration_profile.mount_mappings", example: "{\"W\":\"Americas\"}",
			maps:          []string{""},
			normalization: "mount_mappings",
			computed:      false,
			typed:         types.MapType{ElemType: types.StringType}},
		{name: "expectation.criteria", example: "{\"count\":{\"exact\":1},\"extensions\":[\"csv\"]}",
			maps:          []string{"required_files"},
			normalization: "empty_object",
			computed:      false,
			typed:         types.ObjectType{AttrTypes: map[string]attr.Type{"count": types.ObjectType{AttrTypes: map[string]attr.Type{"exact": types.Int64Type, "max": types.Int64Type, "min": types.Int64Type}}, "extensions": types.ListType{ElemType: types.StringType}, "filename_regex": types.StringType, "total_bytes": types.ObjectType{AttrTypes: map[string]attr.Type{"exact": types.Int64Type, "max": types.Int64Type, "min": types.Int64Type}}, "forbidden_files": types.ListType{ElemType: types.StringType}, "required_files": types.MapType{ElemType: types.ObjectType{AttrTypes: map[string]attr.Type{"count": types.ObjectType{AttrTypes: map[string]attr.Type{"exact": types.Int64Type, "max": types.Int64Type, "min": types.Int64Type}}, "extensions": types.ListType{ElemType: types.StringType}, "filename_regex": types.StringType, "size_bytes": types.ObjectType{AttrTypes: map[string]attr.Type{"exact": types.Int64Type, "max": types.Int64Type, "min": types.Int64Type}}}}}, "content_validation": types.ObjectType{AttrTypes: map[string]attr.Type{"mode": types.StringType, "fts": types.StringType}}}}},
		{name: "file.custom_metadata", example: "{\"department\":\"finance\"}",
			maps:          []string{""},
			normalization: "empty_object",
			computed:      false,
			typed:         types.MapType{ElemType: types.StringType}},
		{name: "folder.custom_metadata", example: "{\"department\":\"finance\"}",
			maps:          []string{""},
			normalization: "empty_object",
			computed:      false,
			typed:         types.MapType{ElemType: types.StringType}},
		{name: "integration_centric_profile.expected_remote_servers", example: "[{\"server_type\":\"dropbox\",\"name\":\"Dropbox\"}]",
			maps:          []string{},
			normalization: "expected_remote_servers",
			computed:      false,
			typed:         types.ListType{ElemType: types.ObjectType{AttrTypes: map[string]attr.Type{"server_type": types.StringType, "name": types.StringType}}}},
		{name: "secret.metadata", example: "{\"header_name\":\"Authorization\"}",
			maps:          []string{},
			normalization: "empty_object",
			computed:      false,
			typed:         types.ObjectType{AttrTypes: map[string]attr.Type{"username": types.StringType, "header_name": types.StringType, "query_parameter_name": types.StringType}}},
		{name: "siem_http_destination.additional_headers", example: "{\"Authorization\":\"Bearer YOUR_TOKEN\"}",
			maps:          []string{""},
			normalization: "",
			computed:      false,
			typed:         types.MapType{ElemType: types.StringType}},
		{name: "site.bundle_watermark_value", example: "{\"gravity\":\"SouthWest\",\"max_height_or_width\":20,\"transparency\":25}",
			maps:          []string{},
			normalization: "watermark",
			computed:      false,
			typed:         types.ObjectType{AttrTypes: map[string]attr.Type{"gravity": types.StringType, "max_height_or_width": types.Int64Type, "transparency": types.Int64Type, "dynamic_text": types.StringType}}},
	} {
		t.Run(test.name, func(t *testing.T) {
			ctx, attributePath := context.Background(), path.Root("value")
			var example any
			require.NoError(t, json.Unmarshal([]byte(test.example), &example))
			var empty any = map[string]any{}
			if _, list := example.([]any); list {
				empty = []any{}
			}
			for _, value := range []any{example, nil, empty} {
				for _, encoded := range []bool{false, true} {
					if test.computed {
						continue
					}
					prior := types.DynamicValue(nativeJSONValue(t, value))
					if value == nil {
						prior = types.DynamicNull()
					}
					if encoded {
						text, err := json.Marshal(value)
						require.NoError(t, err)
						prior = types.DynamicValue(types.StringValue(string(text)))
					}
					request, diags := JSONValueToAPI(ctx, attributePath, prior)
					require.False(t, diags.HasError(), diags)
					assert.Equal(t, value, request)
					validation := &validator.DynamicResponse{}
					DeprecatedJSONEncoding(test.name, test.example, typedSchemaDate).ValidateDynamic(ctx, validator.DynamicRequest{Path: attributePath, ConfigValue: prior}, validation)
					assert.Equal(t, encoded, validation.Diagnostics.WarningsCount() == 1)
					if encoded {
						assert.Contains(t, validation.Diagnostics[0].Detail(), typedSchemaDate)
					}
					state, diags := APIToDynamicJSON(ctx, attributePath, request, prior, test.maps, test.writeOnly, test.normalization)
					require.False(t, diags.HasError(), diags)
					assert.True(t, prior.Equal(state), "apply must preserve the HCL type and representation")
					refresh, diags := APIToDynamicJSON(ctx, attributePath, request, state, test.maps, test.writeOnly, test.normalization)
					require.False(t, diags.HasError(), diags)
					assert.True(t, state.Equal(refresh))
				}
			}
			for _, unknown := range []types.Dynamic{types.DynamicUnknown(), types.DynamicValue(types.StringUnknown())} {
				validation := &validator.DynamicResponse{}
				DeprecatedJSONEncoding(test.name, test.example, typedSchemaDate).ValidateDynamic(ctx, validator.DynamicRequest{Path: attributePath, ConfigValue: unknown}, validation)
				assert.Empty(t, validation.Diagnostics)
				value, diags := JSONValueToAPI(ctx, attributePath, unknown)
				require.False(t, diags.HasError(), diags)
				assert.Nil(t, value)
			}
			responses := []any{example, nil, empty}
			if test.name == "bundle.watermark_value" {
				responses = append(responses, map[string]any{"transparency": float64(25)}, map[string]any{"transparency": "25"}, map[string]any{"transparency": ""}, map[string]any{"transparency": nil}, map[string]any{"extra": "kept"})
			}
			for _, response := range responses {
				state, diags := APIToDynamicJSON(ctx, attributePath, response, types.DynamicNull(), test.maps, test.writeOnly, test.normalization)
				require.False(t, diags.HasError(), diags)
				value, diags := JSONValueToAPI(ctx, attributePath, state)
				require.False(t, diags.HasError(), diags)
				assert.Equal(t, response, value, "imports and computed reads must not add wrappers or coerce values")
				encoded, err := json.Marshal(response)
				require.NoError(t, err)
				for _, prior := range []types.Dynamic{state, types.DynamicValue(types.StringValue(string(encoded)))} {
					value, diags := JSONValueToAPI(ctx, attributePath, prior)
					require.False(t, diags.HasError(), diags)
					_, diags = schemaValue(ctx, attributePath, value, test.typed)
					require.False(t, diags.HasError(), diags)
				}
			}
		})
	}
}
