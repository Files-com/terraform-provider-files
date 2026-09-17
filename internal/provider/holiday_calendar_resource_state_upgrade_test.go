package provider

import (
	"context"
	"encoding/json"
	"testing"

	"github.com/Files-com/terraform-provider-files/lib"
	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/tfsdk"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-go/tftypes"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestHolidayCalendarJSONStateUpgrades(t *testing.T) {
	ctx := context.Background()
	r := &holidayCalendarResource{}
	t.Run("v0/definition", func(t *testing.T) {
		example := "{\"months\":{\"0\":[{\"name\":\"Good Friday\",\"function\":\"easter(year)\",\"function_modifier\":-2}],\"1\":[{\"name\":\"New Year's Day\",\"mday\":1,\"observed\":\"to_weekday_if_weekend(date)\"},{\"name\":\"Third Monday\",\"week\":3,\"wday\":1}],\"11\":[{\"name\":\"Thanksgiving\",\"week\":4,\"wday\":4}],\"12\":[{\"name\":\"Christmas Eve Early Close\",\"mday\":24,\"start_time\":\"13:00\",\"end_time\":\"17:00\",\"year_ranges\":{\"from\":2026}}]}}"
		var apiValue any
		require.NoError(t, json.Unmarshal([]byte(example), &apiValue))
		native, diags := lib.ToDynamic(ctx, path.Root("definition"), apiValue, nil)
		require.False(t, diags.HasError(), diags)
		for _, value := range []types.Dynamic{native, types.DynamicValue(types.StringValue(example)), types.DynamicNull()} {
			upgrader := r.UpgradeState(ctx)[0]
			priorType := upgrader.PriorSchema.Type().TerraformType(ctx).(tftypes.Object)
			attributes := make(map[string]tftypes.Value, len(priorType.AttributeTypes))
			for name, attributeType := range priorType.AttributeTypes {
				attributes[name] = tftypes.NewValue(attributeType, nil)
			}
			prior := tfsdk.State{Schema: *upgrader.PriorSchema, Raw: tftypes.NewValue(priorType, attributes)}
			require.False(t, prior.SetAttribute(ctx, path.Root("definition"), value).HasError())
			response := resource.UpgradeStateResponse{State: tfsdk.State{Schema: r.resourceSchema()}}
			upgrader.StateUpgrader(ctx, resource.UpgradeStateRequest{State: &prior}, &response)
			require.False(t, response.Diagnostics.HasError(), response.Diagnostics)
			var upgraded types.Object
			require.False(t, response.State.GetAttribute(ctx, path.Root("definition"), &upgraded).HasError())
			if value.IsNull() {
				assert.True(t, upgraded.IsNull())
				continue
			}
			actual, diags := lib.SchemaAttributeToInterface(ctx, path.Root("definition"), upgraded)
			require.False(t, diags.HasError(), diags)
			actual, transformDiags0 := lib.UngroupStructuralUnionAtPath(ctx, path.Root("definition"), actual, []string{"months"}, []lib.JSONSchemaVariant{{Name: "calculated_rules", Required: []string{"function"}, Allowed: []string{"name", "observed", "start_time", "end_time", "year_ranges", "function", "function_modifier"}}, {Name: "fixed_rules", Required: []string{"mday"}, Allowed: []string{"name", "observed", "start_time", "end_time", "year_ranges", "mday"}}, {Name: "weekday_rules", Required: []string{"week", "wday"}, Allowed: []string{"name", "observed", "start_time", "end_time", "year_ranges", "week", "wday"}}})
			diags.Append(transformDiags0...)
			require.False(t, diags.HasError(), diags)
			encoded, err := json.Marshal(actual)
			require.NoError(t, err)
			assert.JSONEq(t, example, string(encoded))
		}
	})
}
