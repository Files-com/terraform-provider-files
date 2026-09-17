package provider

import (
	"context"

	"github.com/Files-com/terraform-provider-files/lib"
	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/int64planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

type holidayCalendarResourceModelV0 struct {
	Name       types.String  `tfsdk:"name"`
	Definition types.Dynamic `tfsdk:"definition"`
	Id         types.Int64   `tfsdk:"id"`
	CreatedAt  types.String  `tfsdk:"created_at"`
	UpdatedAt  types.String  `tfsdk:"updated_at"`
}

func (r *holidayCalendarResource) UpgradeState(ctx context.Context) map[int64]resource.StateUpgrader {
	return map[int64]resource.StateUpgrader{
		0: {
			PriorSchema: &schema.Schema{
				Attributes: map[string]schema.Attribute{

					"name": schema.StringAttribute{
						Description: "Holiday Calendar name.",
						Required:    true,
					},
					"definition": schema.DynamicAttribute{
						Description: "Holiday rules for the calendar.",
						Computed:    true,
					},
					"id": schema.Int64Attribute{
						Description: "Holiday Calendar ID. Set a scheduled resource's `holiday_region` to `custom_` followed by this ID to make it skip the days in this calendar.",
						Computed:    true,
						PlanModifiers: []planmodifier.Int64{
							int64planmodifier.UseStateForUnknown(),
						},
					},
					"created_at": schema.StringAttribute{
						Description: "Creation time.",
						Computed:    true,
					},
					"updated_at": schema.StringAttribute{
						Description: "Last update time.",
						Computed:    true,
					},
				},
				Version: 0,
			},
			StateUpgrader: func(ctx context.Context, req resource.UpgradeStateRequest, resp *resource.UpgradeStateResponse) {
				var priorState holidayCalendarResourceModelV0
				resp.Diagnostics.Append(req.State.Get(ctx, &priorState)...)
				if resp.Diagnostics.HasError() {
					return
				}
				upgradedState := holidayCalendarResourceModel{

					Name:      priorState.Name,
					Id:        priorState.Id,
					CreatedAt: priorState.CreatedAt,
					UpdatedAt: priorState.UpdatedAt,
				}
				currentSchema := r.resourceSchema()
				definitionValue, conversionDiags := lib.JSONValueToAPI(ctx, path.Root("definition"), priorState.Definition)
				resp.Diagnostics.Append(conversionDiags...)
				if resp.Diagnostics.HasError() {
					return
				}
				definitionValue, transformDiags0 := lib.GroupStructuralUnionAtPath(ctx, path.Root("definition"), definitionValue, []string{"months"}, []lib.JSONSchemaVariant{{Name: "calculated_rules", Required: []string{"function"}, Allowed: []string{"name", "observed", "start_time", "end_time", "year_ranges", "function", "function_modifier"}}, {Name: "fixed_rules", Required: []string{"mday"}, Allowed: []string{"name", "observed", "start_time", "end_time", "year_ranges", "mday"}}, {Name: "weekday_rules", Required: []string{"week", "wday"}, Allowed: []string{"name", "observed", "start_time", "end_time", "year_ranges", "week", "wday"}}})
				resp.Diagnostics.Append(transformDiags0...)
				definitionType := currentSchema.Attributes["definition"].GetType().(types.ObjectType)
				upgradedState.Definition, conversionDiags = lib.ToObject(ctx, path.Root("definition"), definitionValue, types.ObjectNull(definitionType.AttrTypes))
				resp.Diagnostics.Append(conversionDiags...)
				if resp.Diagnostics.HasError() {
					return
				}
				resp.Diagnostics.Append(resp.State.Set(ctx, upgradedState)...)
			},
		},
	}
}
