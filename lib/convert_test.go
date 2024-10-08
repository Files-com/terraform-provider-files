package lib

import (
	"context"
	"math/big"
	"testing"
	"time"

	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-framework/types/basetypes"
	"github.com/stretchr/testify/assert"
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
