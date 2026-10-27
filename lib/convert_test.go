package lib

import (
	"context"
	"math/big"
	"strings"
	"testing"
	"time"

	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-framework/types/basetypes"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestTimeToStringType(t *testing.T) {
	tests := []struct {
		message  string
		source   string
		dest     string
		expected types.String
	}{
		{
			message:  "Destination is empty so source should be copied",
			source:   "2021-01-01T00:00:00-07:00",
			dest:     "",
			expected: types.StringValue("2021-01-01T00:00:00-07:00"),
		},
		{
			message:  "Source is empty so destination should be empty",
			source:   "",
			dest:     "",
			expected: types.StringValue(""),
		},
		{
			message:  "Destination does not match source so source should be copied",
			source:   "2021-01-01T00:00:00Z",
			dest:     "2021-01-01T01:01:01Z",
			expected: types.StringValue("2021-01-01T00:00:00Z"),
		},
		{
			message:  "Destination does not match source so source should be copied",
			source:   "2021-01-01T00:00:00-07:00",
			dest:     "2021-01-01T01:01:01-07:00",
			expected: types.StringValue("2021-01-01T00:00:00-07:00"),
		},
		{
			message:  "Destination matches source in UTC so destination should not be changed",
			source:   "2021-01-01T00:00:00-07:00",
			dest:     "2021-01-01T07:00:00Z",
			expected: types.StringValue("2021-01-01T07:00:00Z"),
		},
		{
			message:  "Destination in UTC matches source so destination should not be changed",
			source:   "2021-01-01T07:00:00Z",
			dest:     "2021-01-01T00:00:00-07:00",
			expected: types.StringValue("2021-01-01T00:00:00-07:00"),
		},
	}
	for _, c := range tests {
		var sourcePtr *time.Time

		if c.source != "" {
			source, err := time.Parse(time.RFC3339, c.source)
			assert.NoError(t, err, c.message)
			sourcePtr = &source
		}
		dest := types.StringValue(c.dest)

		err := TimeToStringType(context.Background(), path.Empty(), sourcePtr, &dest)
		assert.NoError(t, err, c.message)
		assert.Equal(t, c.expected, dest, c.message)
	}
}

func TestDynamicToStringMapSlice(t *testing.T) {
	tests := []struct {
		message   string
		source    types.Dynamic
		expected  []map[string]interface{}
		expectErr bool
	}{
		{
			message:  "Null dynamic should return nil",
			source:   types.DynamicNull(),
			expected: nil,
		},
		{
			message:  "Unknown dynamic should return nil",
			source:   types.DynamicUnknown(),
			expected: nil,
		},
		{
			message:  "Dynamic with null underlying value should return nil",
			source:   types.DynamicValue(types.StringNull()),
			expected: nil,
		},
		{
			message:  "Dynamic with unknown underlying value should return nil",
			source:   types.DynamicValue(types.StringUnknown()),
			expected: nil,
		},
		{
			message:  "Empty dynamic should return empty slice",
			source:   types.DynamicValue(basetypes.NewTupleValueMust([]attr.Type{}, nil)),
			expected: []map[string]interface{}{},
		},
		{
			message:   "Dynamic with a non-sequence underlying value should error",
			source:    types.DynamicValue(types.StringValue("foo")),
			expectErr: true,
		},
		{
			message: "Dynamic with a tuple underlying value should return string map slice",
			source: types.DynamicValue(basetypes.NewTupleValueMust(
				[]attr.Type{
					types.ObjectType{AttrTypes: map[string]attr.Type{"foo": types.StringType, "bar": types.StringType}},
					types.ObjectType{AttrTypes: map[string]attr.Type{"baz": types.BoolType, "qux": types.NumberType}},
					types.ObjectType{AttrTypes: map[string]attr.Type{"object": types.ObjectType{AttrTypes: map[string]attr.Type{"sub-key": types.StringType}}}},
					types.ObjectType{AttrTypes: map[string]attr.Type{"tuple": types.TupleType{ElemTypes: []attr.Type{types.StringType, types.BoolType}}}},
				},
				[]attr.Value{
					basetypes.NewObjectValueMust(
						map[string]attr.Type{"foo": types.StringType, "bar": types.StringType},
						map[string]attr.Value{"foo": types.StringValue("first"), "bar": types.StringValue("second")},
					),
					basetypes.NewObjectValueMust(
						map[string]attr.Type{"baz": types.BoolType, "qux": types.NumberType},
						map[string]attr.Value{"baz": types.BoolValue(true), "qux": types.NumberValue(big.NewFloat(123))},
					),
					basetypes.NewObjectValueMust(
						map[string]attr.Type{"object": types.ObjectType{AttrTypes: map[string]attr.Type{"sub-key": types.StringType}}},
						map[string]attr.Value{"object": basetypes.NewObjectValueMust(map[string]attr.Type{"sub-key": types.StringType}, map[string]attr.Value{"sub-key": types.StringValue("value")})},
					),
					basetypes.NewObjectValueMust(
						map[string]attr.Type{"tuple": types.TupleType{ElemTypes: []attr.Type{types.StringType, types.BoolType}}},
						map[string]attr.Value{"tuple": basetypes.NewTupleValueMust([]attr.Type{types.StringType, types.BoolType}, []attr.Value{types.StringValue("first"), types.BoolValue(true)})},
					),
				},
			)),
			expected: []map[string]interface{}{
				{"foo": "first", "bar": "second"},
				{"baz": true, "qux": float64(123)},
				{"object": map[string]interface{}{"sub-key": "value"}},
				{"tuple": []interface{}{"first", true}},
			},
		},
	}
	for _, c := range tests {
		result, diags := DynamicToStringMapSlice(context.Background(), path.Empty(), c.source)
		assert.Equal(t, c.expectErr, diags.HasError(), c.message)
		assert.Equal(t, c.expected, result, c.message)
	}
}

func TestNullScalarConversion(t *testing.T) {
	value, diags := AttributeToInterface(context.Background(), path.Root("value"), types.BoolNull())
	assert.False(t, diags.HasError())
	assert.Equal(t, false, value)

	value, diags = SchemaAttributeToInterface(context.Background(), path.Root("value"), types.BoolNull())
	assert.False(t, diags.HasError())
	assert.Nil(t, value)
}

func TestSchemaInt64StringConversion(t *testing.T) {
	value, diags := schemaValue(context.Background(), path.Root("value"), "5", types.Int64Type)
	assert.False(t, diags.HasError())
	assert.Equal(t, types.Int64Value(5), value)

	value, diags = schemaValue(context.Background(), path.Root("value"), "", types.Int64Type)
	assert.False(t, diags.HasError())
	assert.Equal(t, types.Int64Null(), value)

	value, diags = schemaValue(context.Background(), path.Root("value"), "not-an-integer", types.Int64Type)
	assert.True(t, diags.HasError())
	assert.Nil(t, value)
}

func TestSchemaObjectRoundTrip(t *testing.T) {
	stepTypes := map[string]attr.Type{"path": types.StringType, "delay": types.Int64Type}
	definitionTypes := map[string]attr.Type{
		"version":  types.Int64Type,
		"labels":   types.MapType{ElemType: types.StringType},
		"steps":    types.ListType{ElemType: types.ObjectType{AttrTypes: stepTypes}},
		"optional": types.StringType,
	}
	source := map[string]interface{}{
		"version":  1,
		"labels":   map[string]interface{}{"team": "operations"},
		"steps":    []interface{}{map[string]interface{}{"path": "incoming/", "delay": 30}},
		"optional": nil,
	}

	value, diags := ToObject(context.Background(), path.Root("definition"), source, types.ObjectNull(definitionTypes))
	assert.False(t, diags.HasError())
	// A typed object cannot tell a field set to null from an omitted field, so the API value leaves both out.
	result, diags := SchemaAttributeToInterface(context.Background(), path.Root("definition"), value)
	assert.False(t, diags.HasError())
	assert.Equal(t, map[string]interface{}{
		"version": int64(1),
		"labels":  map[string]interface{}{"team": "operations"},
		"steps":   []interface{}{map[string]interface{}{"path": "incoming/", "delay": int64(30)}},
	}, result)

	// A Dynamic value keeps the null in the API value.
	dynamic, diags := JSONValueToAPI(context.Background(), path.Root("definition"), types.DynamicValue(nativeJSONValue(t, map[string]any{"optional": nil})))
	require.False(t, diags.HasError(), diags)
	assert.Equal(t, map[string]any{"optional": nil}, dynamic)
}

func TestDiscriminatedUnionRoundTrip(t *testing.T) {
	variants := []JSONSchemaVariant{{Name: "trigger_manual", Value: "trigger_manual"}, {Name: "copy_file", Value: "copy_file"}}
	apiValue := map[string]interface{}{"nodes": []interface{}{
		map[string]interface{}{"id": "trigger", "type": "trigger_manual"},
		map[string]interface{}{"id": "archive", "type": "copy_file", "config": map[string]interface{}{"destinations": []interface{}{"archive/"}}},
	}}

	terraformValue, diags := WrapDiscriminatedUnionAtPath(context.Background(), path.Root("definition"), apiValue, []string{"nodes"}, "type", variants)
	assert.False(t, diags.HasError())
	assert.Equal(t, map[string]interface{}{"nodes": []interface{}{
		map[string]interface{}{"trigger_manual": map[string]interface{}{"id": "trigger"}},
		map[string]interface{}{"copy_file": map[string]interface{}{"id": "archive", "config": map[string]interface{}{"destinations": []interface{}{"archive/"}}}},
	}}, terraformValue)

	result, diags := UnwrapDiscriminatedUnionAtPath(context.Background(), path.Root("definition"), terraformValue, []string{"nodes"}, "type", variants)
	assert.False(t, diags.HasError())
	assert.Equal(t, apiValue, result)

	terraformValue, diags = WrapDiscriminatedUnionAtPath(context.Background(), path.Root("definition"), nil, []string{"nodes"}, "type", variants)
	assert.False(t, diags.HasError())
	assert.Nil(t, terraformValue)

	result, diags = UnwrapDiscriminatedUnionAtPath(context.Background(), path.Root("definition"), nil, []string{"nodes"}, "type", variants)
	assert.False(t, diags.HasError())
	assert.Nil(t, result)

	_, diags = UnwrapDiscriminatedUnionAtPath(context.Background(), path.Root("definition"), map[string]interface{}{"nodes": []interface{}{map[string]interface{}{"trigger_manual": map[string]interface{}{}, "copy_file": map[string]interface{}{}}}}, []string{"nodes"}, "type", variants)
	assert.True(t, diags.HasError())
}

func TestSiblingDiscriminatorRoundTrip(t *testing.T) {
	variants := []JSONSchemaVariant{{Name: "webhook", Value: "webhook"}, {Name: "storage_region", Value: "storage_region"}, {Name: "limit_file_regex", Value: "limit_file_regex"}}
	tests := []struct {
		selector string
		apiValue any
	}{
		{selector: "webhook", apiValue: map[string]interface{}{"urls": []interface{}{"https://example.com/webhook"}}},
		{selector: "storage_region", apiValue: "us-east-1"},
		{selector: "limit_file_regex", apiValue: []interface{}{`/Document-.*/`}},
	}

	for _, test := range tests {
		terraformValue, diags := WrapSiblingDiscriminator(context.Background(), path.Root("value"), test.apiValue, test.selector, variants)
		assert.False(t, diags.HasError())
		assert.Equal(t, map[string]interface{}{test.selector: test.apiValue}, terraformValue)

		result, diags := UnwrapSiblingDiscriminator(context.Background(), path.Root("value"), terraformValue, test.selector, variants)
		assert.False(t, diags.HasError())
		assert.Equal(t, test.apiValue, result)
	}

	_, diags := UnwrapSiblingDiscriminator(context.Background(), path.Root("value"), map[string]interface{}{"webhook": map[string]interface{}{}}, "storage_region", variants)
	assert.True(t, diags.HasError())
}

// Rails merges a Behavior value on update, so a typed update sends null for a writable field that state sets and the
// configuration removed. Computed fields, fields of another behavior and unchanged fields are left alone.
func TestNullRemovedSiblingFields(t *testing.T) {
	variants := []JSONSchemaVariant{
		{Name: "webhook", Value: "webhook", Writable: []string{"urls", "method", "headers"}},
		{Name: "file_expiration", Value: "file_expiration", Writable: []string{"days_to_retain"}},
	}
	webhookType := types.ObjectType{AttrTypes: map[string]attr.Type{"urls": types.ListType{ElemType: types.StringType}, "method": types.StringType, "headers": types.MapType{ElemType: types.StringType}, "configured": types.BoolType}}
	expirationType := types.ObjectType{AttrTypes: map[string]attr.Type{"days_to_retain": types.Int64Type}}
	valueType := map[string]attr.Type{"webhook": webhookType, "file_expiration": expirationType}
	urls := types.ListValueMust(types.StringType, []attr.Value{types.StringValue("https://example.com/hook")})
	prior := types.ObjectValueMust(valueType, map[string]attr.Value{
		"webhook": types.ObjectValueMust(webhookType.AttrTypes, map[string]attr.Value{
			"urls": urls, "method": types.StringValue("POST"), "headers": types.MapValueMust(types.StringType, map[string]attr.Value{"X-Example": types.StringValue("1")}), "configured": types.BoolValue(true),
		}),
		"file_expiration": types.ObjectNull(expirationType.AttrTypes),
	})
	request := map[string]interface{}{"urls": []interface{}{"https://example.com/hook"}, "method": "GET"}

	result, diags := NullRemovedSiblingFields(context.Background(), path.Root("value"), request, prior, "webhook", variants)
	require.False(t, diags.HasError(), diags)
	assert.Equal(t, map[string]interface{}{"urls": []interface{}{"https://example.com/hook"}, "method": "GET", "headers": nil}, result)
	assert.NotContains(t, request, "headers", "the request value is not changed in place")

	result, diags = NullRemovedSiblingFields(context.Background(), path.Root("value"), map[string]interface{}{"days_to_retain": float64(7)}, prior, "file_expiration", variants)
	require.False(t, diags.HasError(), diags)
	assert.Equal(t, map[string]interface{}{"days_to_retain": float64(7)}, result, "prior state under another behavior removes nothing")

	for _, empty := range []attr.Value{types.ObjectNull(valueType), types.ObjectUnknown(valueType)} {
		result, diags = NullRemovedSiblingFields(context.Background(), path.Root("value"), request, empty, "webhook", variants)
		require.False(t, diags.HasError(), diags)
		assert.Equal(t, request, result)
	}
}

func TestSiblingDiscriminatorWrapsLegacyValues(t *testing.T) {
	variants := []JSONSchemaVariant{
		{Name: "file_expiration", Value: "file_expiration", LegacyProperty: "days_to_retain"},
		{Name: "webhook", Value: "webhook", LegacyProperty: "urls", LegacyList: true},
	}

	result, diags := WrapSiblingDiscriminator(context.Background(), path.Root("value"), float64(14), "file_expiration", variants)
	assert.False(t, diags.HasError())
	assert.Equal(t, map[string]interface{}{"file_expiration": map[string]interface{}{"days_to_retain": float64(14)}}, result)

	result, diags = WrapSiblingDiscriminator(context.Background(), path.Root("value"), "https://example.com/webhook", "webhook", variants)
	assert.False(t, diags.HasError())
	assert.Equal(t, map[string]interface{}{"webhook": map[string]interface{}{"urls": []interface{}{"https://example.com/webhook"}}}, result)

	result, diags = WrapSiblingDiscriminator(context.Background(), path.Root("value"), "", "file_expiration", variants)
	assert.False(t, diags.HasError())
	assert.Nil(t, result)
}

// A typed object also reads its JSON text, as older records and Dynamic state may hold a Google key file. A response
// without private_key leaves it null. Text that is not a JSON object fails instead of becoming an empty object.
func TestSchemaValueReadsObjectJSONText(t *testing.T) {
	credentials := types.ObjectType{AttrTypes: map[string]attr.Type{"client_email": types.StringType, "private_key": types.StringType}}
	config := types.ObjectType{AttrTypes: map[string]attr.Type{"google_credentials": credentials}}
	for _, test := range []struct {
		name        string
		credentials any
		want        attr.Value
	}{
		{"object without private_key", map[string]any{"client_email": "events@example.com"}, types.ObjectValueMust(credentials.AttrTypes, map[string]attr.Value{"client_email": types.StringValue("events@example.com"), "private_key": types.StringNull()})},
		{"key file text", `{"client_email":"events@example.com","private_key":"example-key"}`, types.ObjectValueMust(credentials.AttrTypes, map[string]attr.Value{"client_email": types.StringValue("events@example.com"), "private_key": types.StringValue("example-key")})},
		{"malformed text", `{"client_email":`, nil},
		{"list text", `["events@example.com"]`, nil},
		{"null text", `null`, nil},
	} {
		t.Run(test.name, func(t *testing.T) {
			value, diags := schemaValue(context.Background(), path.Root("config"), map[string]any{"google_credentials": test.credentials}, config)
			if test.want == nil {
				require.True(t, diags.HasError())
				assert.Equal(t, "Failed to convert API value", diags[0].Summary())
				assert.Nil(t, value)
				return
			}
			require.False(t, diags.HasError(), diags)
			assert.Equal(t, types.ObjectValueMust(config.AttrTypes, map[string]attr.Value{"google_credentials": test.want}), value)
		})
	}
}

func TestSchemaValueReadsLegacyBooleanStrings(t *testing.T) {
	value, diags := schemaValue(context.Background(), path.Root("armor"), "true", types.BoolType)
	assert.False(t, diags.HasError())
	assert.Equal(t, types.BoolValue(true), value)
}

// Stored legacy values of the typed delivery_policy and ai_feature_availability fields convert, or fail with an error.
// They never become 0 or false, and are never dropped. Only a list marked x-files-legacy-single-item, as
// retry_delay_minutes is, reads one stored item as a one-item list.
func TestSchemaValueReadsLegacyDeliveryPolicyAndAIOptions(t *testing.T) {
	delays := types.ListType{ElemType: types.Int64Type}
	policy := types.ObjectType{AttrTypes: map[string]attr.Type{"max_attempts": types.Int64Type, "retry_delay_minutes": delays}}
	options := types.ObjectType{AttrTypes: map[string]attr.Type{"site_admins": types.BoolType, "all_users": types.BoolType}}
	delayList := func(values ...int64) map[string]attr.Value {
		items := []attr.Value{}
		for _, value := range values {
			items = append(items, types.Int64Value(value))
		}
		return map[string]attr.Value{"max_attempts": types.Int64Null(), "retry_delay_minutes": types.ListValueMust(types.Int64Type, items)}
	}
	for _, test := range []struct {
		name   string
		source map[string]any
		target types.ObjectType
		marked bool
		want   map[string]attr.Value
		error  string
	}{
		{name: "numeric strings", source: map[string]any{"max_attempts": "3", "retry_delay_minutes": []any{"1", float64(5)}}, target: policy, marked: true,
			want: map[string]attr.Value{"max_attempts": types.Int64Value(3), "retry_delay_minutes": types.ListValueMust(types.Int64Type, []attr.Value{types.Int64Value(1), types.Int64Value(5)})}},
		{name: "null and an empty list", source: map[string]any{"max_attempts": nil, "retry_delay_minutes": []any{}}, target: policy, marked: true, want: delayList()},
		{name: "null delays", source: map[string]any{"retry_delay_minutes": nil}, target: policy, marked: true,
			want: map[string]attr.Value{"max_attempts": types.Int64Null(), "retry_delay_minutes": types.ListNull(types.Int64Type)}},
		{name: "one delay in a marked list", source: map[string]any{"retry_delay_minutes": float64(5)}, target: policy, marked: true, want: delayList(5)},
		{name: "one numeric string delay in a marked list", source: map[string]any{"retry_delay_minutes": "5"}, target: policy, marked: true, want: delayList(5)},
		{name: "one non-numeric delay in a marked list", source: map[string]any{"retry_delay_minutes": "soon"}, target: policy, marked: true, error: "retry_delay_minutes[0]: Expected basetypes.Int64Type, got string"},
		{name: "one fraction delay in a marked list", source: map[string]any{"retry_delay_minutes": 1.5}, target: policy, marked: true, error: "retry_delay_minutes[0]: Expected basetypes.Int64Type, got float64"},
		{name: "one delay in an unmarked list", source: map[string]any{"retry_delay_minutes": float64(5)}, target: policy, error: "retry_delay_minutes: Expected types.ListType[basetypes.Int64Type], got float64"},
		{name: "one numeric string delay in an unmarked list", source: map[string]any{"retry_delay_minutes": "5"}, target: policy, error: "retry_delay_minutes: Expected types.ListType[basetypes.Int64Type], got string"},
		{name: "non-numeric attempts", source: map[string]any{"max_attempts": "three"}, target: policy, marked: true, error: "max_attempts: Expected basetypes.Int64Type, got string"},
		{name: "fraction delay", source: map[string]any{"retry_delay_minutes": []any{1.5}}, target: policy, marked: true, error: "retry_delay_minutes[0]: Expected basetypes.Int64Type, got float64"},
		{name: "string booleans and null", source: map[string]any{"site_admins": "true", "all_users": nil}, target: options,
			want: map[string]attr.Value{"site_admins": types.BoolValue(true), "all_users": types.BoolNull()}},
		{name: "other boolean string", source: map[string]any{"site_admins": "yes"}, target: options, error: "site_admins: Expected basetypes.BoolType, got string"},
	} {
		t.Run(test.name, func(t *testing.T) {
			ctx, source := context.Background(), any(test.source)
			if test.marked {
				var diags diag.Diagnostics
				source, diags = WrapLegacySingleItemAtPath(ctx, path.Root("value"), source, []string{"retry_delay_minutes"})
				require.False(t, diags.HasError(), diags)
			}
			value, diags := schemaValue(ctx, path.Root("value"), source, test.target)
			if test.error == "" {
				require.False(t, diags.HasError(), diags)
				assert.Equal(t, types.ObjectValueMust(test.target.AttrTypes, test.want), value)
				return
			}
			require.Equal(t, 1, diags.ErrorsCount(), diags)
			assert.Nil(t, value)
			attributePath := diags[0].(diag.DiagnosticWithPath).Path().String()
			assert.Equal(t, test.error, strings.TrimPrefix(attributePath, "value.")+": "+diags[0].Detail())
			assert.Equal(t, "Failed to convert API value", diags[0].Summary())
		})
	}
}

func TestWrapLegacySingleItemAtPathLeavesListsAndOtherFields(t *testing.T) {
	ctx := context.Background()
	source := map[string]any{"retry_delay_minutes": []any{float64(1), "5"}, "max_attempts": "3"}
	value, diags := WrapLegacySingleItemAtPath(ctx, path.Root("delivery_policy"), source, []string{"retry_delay_minutes"})
	require.False(t, diags.HasError(), diags)
	assert.Equal(t, source, value)

	value, diags = WrapLegacySingleItemAtPath(ctx, path.Root("delivery_policy"), map[string]any{"retry_delay_minutes": float64(5), "max_attempts": "3"}, []string{"retry_delay_minutes"})
	require.False(t, diags.HasError(), diags)
	assert.Equal(t, map[string]any{"retry_delay_minutes": []any{float64(5)}, "max_attempts": "3"}, value, "only the marked field changes, and its item is not converted yet")

	for _, empty := range []any{nil, map[string]any{}} {
		value, diags = WrapLegacySingleItemAtPath(ctx, path.Root("delivery_policy"), empty, []string{"retry_delay_minutes"})
		require.False(t, diags.HasError(), diags)
		assert.Equal(t, empty, value)
	}
}

func TestReadLegacyTruthyAtPath(t *testing.T) {
	ctx := context.Background()
	for _, test := range []struct {
		value any
		want  any
	}{
		{float64(0), true},
		{float64(1), true},
		{"false", true},
		{"", true},
		{true, true},
		{false, false},
		{nil, nil},
	} {
		value, diags := ReadLegacyTruthyAtPath(ctx, path.Root("export_options"), map[string]any{"all_folders": test.value}, []string{"all_folders"})
		require.False(t, diags.HasError(), diags)
		assert.Equal(t, map[string]any{"all_folders": test.want}, value)
	}
	empty := map[string]any{}
	value, diags := ReadLegacyTruthyAtPath(ctx, path.Root("export_options"), empty, []string{"all_folders"})
	require.False(t, diags.HasError(), diags)
	assert.Equal(t, empty, value)
}

func TestStructuralUnionRoundTrip(t *testing.T) {
	variants := []JSONSchemaVariant{
		{Name: "calculated_rules", Required: []string{"function"}, Allowed: []string{"name", "function"}},
		{Name: "fixed_rules", Required: []string{"mday"}, Allowed: []string{"name", "mday"}},
		{Name: "weekday_rules", Required: []string{"week", "wday"}, Allowed: []string{"name", "week", "wday"}},
	}
	apiValue := map[string]interface{}{"months": map[string]interface{}{"1": []interface{}{
		map[string]interface{}{"name": "New Year's Day", "mday": int64(1)},
		map[string]interface{}{"name": "Third Monday", "week": int64(3), "wday": int64(1)},
	}}}

	terraformValue, diags := GroupStructuralUnionAtPath(context.Background(), path.Root("definition"), apiValue, []string{"months"}, variants)
	assert.False(t, diags.HasError())
	assert.Equal(t, map[string]interface{}{"months": map[string]interface{}{"1": map[string]interface{}{
		"fixed_rules":   []interface{}{map[string]interface{}{"name": "New Year's Day", "mday": int64(1)}},
		"weekday_rules": []interface{}{map[string]interface{}{"name": "Third Monday", "week": int64(3), "wday": int64(1)}},
	}}}, terraformValue)

	result, diags := UngroupStructuralUnionAtPath(context.Background(), path.Root("definition"), terraformValue, []string{"months"}, variants)
	assert.False(t, diags.HasError())
	assert.Equal(t, apiValue, result)

	_, diags = GroupStructuralUnionAtPath(context.Background(), path.Root("definition"), map[string]interface{}{"months": map[string]interface{}{"1": []interface{}{map[string]interface{}{"name": "No selector"}}}}, []string{"months"}, variants)
	assert.True(t, diags.HasError())
}

func TestListValueToString(t *testing.T) {
	tests := []struct {
		message   string
		source    types.List
		expected  string
		expectErr bool
	}{
		{
			message:  "Null list should return empty string",
			source:   types.ListNull(types.StringType),
			expected: "",
		},
		{
			message:  "Unknown list should return empty string",
			source:   types.ListUnknown(types.StringType),
			expected: "",
		},
		{
			message:  "Empty list should return empty string",
			source:   types.ListValueMust(types.StringType, []attr.Value{}),
			expected: "",
		},
		{
			message: "String list should return comma separated string",
			source: types.ListValueMust(types.StringType, []attr.Value{
				types.StringValue("foo"),
				types.StringValue("bar"),
				types.StringValue("baz"),
			}),
			expected: "foo,bar,baz",
		},
		{
			message: "Int list should return comma separated string",
			source: types.ListValueMust(types.Int64Type, []attr.Value{
				types.Int64Value(1),
				types.Int64Value(2),
				types.Int64Value(123),
			}),
			expected: "1,2,123",
		},
		{
			message:   "Map list should error",
			source:    types.ListValueMust(types.MapType{ElemType: types.StringType}, []attr.Value{basetypes.NewMapValueMust(types.StringType, map[string]attr.Value{"foo": types.StringValue("bar")})}),
			expectErr: true,
		},
	}
	for _, c := range tests {
		result, diags := ListValueToString(context.Background(), path.Empty(), c.source, ",")
		assert.Equal(t, c.expectErr, diags.HasError(), c.message)
		assert.Equal(t, c.expected, result, c.message)
	}
}

func TestMapValueFrom(t *testing.T) {
	stringMap, diags := types.MapValueFrom(context.Background(), types.StringType, map[string]string{"foo": "bar", "baz": "qux"})
	assert.False(t, diags.HasError(), "typed string map should convert to types.Map")

	var stringMapResult map[string]string
	diags = stringMap.ElementsAs(context.Background(), &stringMapResult, false)
	assert.False(t, diags.HasError(), "typed string map should round-trip from types.Map")
	assert.Equal(t, map[string]string{"foo": "bar", "baz": "qux"}, stringMapResult)

	listMap, diags := types.MapValueFrom(
		context.Background(),
		types.ListType{ElemType: types.StringType},
		map[string][]string{"foo": {"bar", "baz"}, "empty": {}},
	)
	assert.False(t, diags.HasError(), "typed string slice map should convert to types.Map")

	var listMapResult map[string][]string
	diags = listMap.ElementsAs(context.Background(), &listMapResult, false)
	assert.False(t, diags.HasError(), "typed string slice map should round-trip from types.Map")
	assert.Equal(t, map[string][]string{"foo": {"bar", "baz"}, "empty": {}}, listMapResult)
}
