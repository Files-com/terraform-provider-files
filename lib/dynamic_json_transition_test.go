package lib

import (
	"context"
	"encoding/json"
	"math/big"
	"testing"

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
			DeprecatedJSONEncoding("files_siem_http_destination.additional_headers", "additional_headers = { Authorization = \"Bearer YOUR_TOKEN\" }", "March 1, 2027").ValidateDynamic(context.Background(), validator.DynamicRequest{Path: path.Root("additional_headers"), ConfigValue: test.value}, response)
			assert.Equal(t, test.warnings, response.Diagnostics.WarningsCount())
			assert.Equal(t, test.errors, response.Diagnostics.ErrorsCount())
			for _, diagnostic := range response.Diagnostics {
				assert.NotContains(t, diagnostic.Detail(), "private-data")
			}
			if test.warnings > 0 {
				assert.Contains(t, response.Diagnostics[0].Detail(), "stop being accepted on March 1, 2027")
			}
		})
	}
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

func TestJSONNullMaterializedAsEmptyObject(t *testing.T) {
	for _, normalization := range []string{"watermark", "empty_object"} {
		prior := types.DynamicValue(types.StringValue("null"))
		state, diags := APIToDynamicJSON(context.Background(), path.Root("value"), map[string]any{}, prior, nil, nil, normalization)
		require.False(t, diags.HasError(), diags)
		assert.True(t, prior.Equal(state))
	}
}

func TestTypedListFromSDKObjects(t *testing.T) {
	objectType := types.ObjectType{AttrTypes: map[string]attr.Type{"server_type": types.StringType, "name": types.StringType}}
	value, diags := ToList(context.Background(), path.Root("expected_remote_servers"), []map[string]any{{"server_type": "dropbox"}}, types.ListNull(objectType))
	assert.False(t, diags.HasError(), diags)
	assert.Len(t, value.Elements(), 1)
	assert.Equal(t, types.StringNull(), value.Elements()[0].(types.Object).Attributes()["name"])
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
