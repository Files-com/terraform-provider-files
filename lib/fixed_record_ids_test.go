package lib

import (
	"context"
	"math/big"
	"testing"

	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// Rails keeps a form field whose request names its ID, with the stored values the request leaves out, and stores a
// changed field under a new ID.
func TestFixedRecordIDs(t *testing.T) {
	ctx := context.Background()
	fieldType := types.ObjectType{AttrTypes: map[string]attr.Type{"id": types.Int64Type, "label": types.StringType, "help_text": types.StringType, "required": types.BoolType}}
	field := func(id, label, helpText, required attr.Value) attr.Value {
		return types.ObjectValueMust(fieldType.AttrTypes, map[string]attr.Value{"id": id, "label": label, "help_text": helpText, "required": required})
	}
	list := func(fields ...attr.Value) types.List { return types.ListValueMust(fieldType, fields) }
	str, null, unknown := types.StringValue, types.StringNull(), types.StringUnknown()
	name := field(types.Int64Value(7), str("Name"), str("a"), types.BoolValue(false))
	team := field(types.Int64Value(8), str("Team"), null, types.BoolValue(true))
	stored := list(name, team)
	modifier := FixedRecordIDs("id", []string{"required"})
	// plan stands for the framework's plan: computed fields that the configuration leaves out are unknown.
	plan := func(config types.List, state types.List) (types.List, int) {
		elements := make([]attr.Value, len(config.Elements()))
		for index, element := range config.Elements() {
			attributes := element.(types.Object).Attributes()
			for key, value := range attributes {
				if value.IsNull() {
					attributes[key] = map[attr.Type]attr.Value{types.StringType: types.StringUnknown(), types.BoolType: types.BoolUnknown(), types.Int64Type: types.Int64Unknown()}[value.Type(ctx)]
				}
			}
			elements[index] = types.ObjectValueMust(fieldType.AttrTypes, attributes)
		}
		resp := &planmodifier.ListResponse{PlanValue: types.ListValueMust(fieldType, elements)}
		modifier.PlanModifyList(ctx, planmodifier.ListRequest{Path: path.Root("form_fields"), ConfigValue: config, StateValue: state, PlanValue: resp.PlanValue}, resp)
		return resp.PlanValue, resp.Diagnostics.ErrorsCount()
	}

	planned, errors := plan(list(field(types.Int64Value(8), null, null, types.BoolNull()), field(types.Int64Value(7), null, null, types.BoolNull())), stored)
	require.Zero(t, errors)
	assert.Equal(t, list(team, name), planned, "ID-only references in a new order plan the stored fields of each ID")

	planned, errors = plan(list(field(types.Int64Value(7), str("Name"), null, types.BoolNull())), stored)
	require.Zero(t, errors)
	assert.Equal(t, list(name), planned, "a partial unchanged definition keeps the left-out help_text and required")

	_, errors = plan(list(field(types.Int64Value(7), str("Renamed"), null, types.BoolNull())), stored)
	assert.Equal(t, 1, errors, "a changed label would get a new ID")
	_, errors = plan(list(field(types.Int64Value(7), null, null, types.BoolValue(true))), stored)
	assert.Equal(t, 1, errors, "a configured required that differs would get a new ID")

	planned, errors = plan(list(field(types.Int64Value(7), unknown, null, types.BoolNull())), stored)
	require.Zero(t, errors, "an unknown label is checked when Terraform plans again with the final value")
	assert.True(t, planned.Elements()[0].(types.Object).Attributes()["label"].IsUnknown())
	assert.Equal(t, str("a"), planned.Elements()[0].(types.Object).Attributes()["help_text"])

	planned, errors = plan(list(field(types.Int64Null(), str("New"), null, types.BoolNull())), stored)
	require.Zero(t, errors)
	added := planned.Elements()[0].(types.Object).Attributes()
	assert.Equal(t, null, added["help_text"], "a new field plans a left-out help_text as null")
	assert.True(t, added["required"].IsUnknown(), "the API fills required for a new field")
	assert.True(t, added["id"].IsUnknown(), "the API assigns the ID")

	planned, errors = plan(list(field(types.Int64Value(99), str("Other"), null, types.BoolNull())), stored)
	require.Zero(t, errors, "an ID that is not stored is left to the API, which rejects it before any change")
	assert.True(t, planned.Elements()[0].(types.Object).Attributes()["help_text"].IsUnknown())

	planned, errors = plan(list(field(types.Int64Null(), str("New"), null, types.BoolNull())), types.ListNull(fieldType))
	require.Zero(t, errors)
	assert.Equal(t, null, planned.Elements()[0].(types.Object).Attributes()["help_text"], "a create plans left-out fields as null too")

	// A Dynamic record holds only the configured keys, and the API casts boolean and number strings.
	dynamicStored := types.DynamicValue(types.StringValue(`[{"id":7,"label":"Name","required":false}]`))
	dynamicErrors := func(config types.Dynamic) int {
		resp := &planmodifier.DynamicResponse{PlanValue: config}
		modifier.PlanModifyDynamic(ctx, planmodifier.DynamicRequest{Path: path.Root("form_fields"), ConfigValue: config, StateValue: dynamicStored, PlanValue: config}, resp)
		return resp.Diagnostics.ErrorsCount()
	}
	assert.Zero(t, dynamicErrors(types.DynamicValue(types.StringValue(`[{"id":7}]`))), "an ID-only reference")
	assert.Zero(t, dynamicErrors(types.DynamicValue(types.StringValue(`[{"id":7,"label":"Name","required":"false"}]`))), `"false" is the same stored value`)
	assert.Zero(t, dynamicErrors(types.DynamicValue(types.StringValue(`[{"id":7,"label":"Name","help_text":"x"}]`))), "a key that state never tracked is not compared")
	assert.Equal(t, 1, dynamicErrors(types.DynamicValue(types.StringValue(`[{"id":7,"label":"Renamed"}]`))))
	objectType := map[string]attr.Type{"id": types.NumberType, "label": types.StringType}
	pending := types.DynamicValue(types.TupleValueMust([]attr.Type{types.ObjectType{AttrTypes: objectType}}, []attr.Value{
		types.ObjectValueMust(objectType, map[string]attr.Value{"id": types.NumberValue(big.NewFloat(7)), "label": types.StringUnknown()}),
	}))
	assert.Zero(t, dynamicErrors(pending), "an unknown label is checked when Terraform plans again with the final value")
}
