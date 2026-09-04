package provider

import (
	"context"
	"fmt"
	"strconv"
	"strings"

	files_sdk "github.com/Files-com/files-sdk-go/v3"
	holiday_calendar "github.com/Files-com/files-sdk-go/v3/holidaycalendar"
	"github.com/Files-com/terraform-provider-files/lib"
	"github.com/hashicorp/terraform-plugin-framework-validators/int64validator"
	"github.com/hashicorp/terraform-plugin-framework-validators/listvalidator"
	"github.com/hashicorp/terraform-plugin-framework-validators/mapvalidator"
	"github.com/hashicorp/terraform-plugin-framework-validators/stringvalidator"
	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"

	"github.com/hashicorp/terraform-plugin-framework/resource/schema/int64planmodifier"

	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"

	"github.com/hashicorp/terraform-plugin-framework/schema/validator"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

var (
	_ resource.Resource                 = &holidayCalendarResource{}
	_ resource.ResourceWithConfigure    = &holidayCalendarResource{}
	_ resource.ResourceWithImportState  = &holidayCalendarResource{}
	_ resource.ResourceWithUpgradeState = &holidayCalendarResource{}
)

func NewHolidayCalendarResource() resource.Resource {
	return &holidayCalendarResource{}
}

type holidayCalendarResource struct {
	client *holiday_calendar.Client
}

type holidayCalendarResourceModel struct {
	Name       types.String `tfsdk:"name"`
	Definition types.Object `tfsdk:"definition"`
	Id         types.Int64  `tfsdk:"id"`
	CreatedAt  types.String `tfsdk:"created_at"`
	UpdatedAt  types.String `tfsdk:"updated_at"`
}

type holidayCalendarResourceModelV0 struct {
	Name       types.String  `tfsdk:"name"`
	Definition types.Dynamic `tfsdk:"definition"`
	Id         types.Int64   `tfsdk:"id"`
	CreatedAt  types.String  `tfsdk:"created_at"`
	UpdatedAt  types.String  `tfsdk:"updated_at"`
}

func (r *holidayCalendarResource) Configure(_ context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {
	if req.ProviderData == nil {
		return
	}

	sdk_config, ok := req.ProviderData.(files_sdk.Config)

	if !ok {
		resp.Diagnostics.AddError(
			"Unexpected Data Source Configure Type",
			fmt.Sprintf("Expected files_sdk.Config, got: %T. Please report this issue to the provider developers.", req.ProviderData),
		)

		return
	}

	r.client = &holiday_calendar.Client{Config: sdk_config}
}

func (r *holidayCalendarResource) Metadata(_ context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_holiday_calendar"
}

func (r *holidayCalendarResource) Schema(_ context.Context, _ resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = r.resourceSchema()
}

func (r *holidayCalendarResource) resourceSchema() schema.Schema {
	return schema.Schema{
		Description: "A Holiday Calendar defines site-wide holiday dates and optional partial-day windows that scheduled resources skip.",
		Attributes: map[string]schema.Attribute{
			"name": schema.StringAttribute{
				Description: "Holiday Calendar name.",
				Required:    true,
			},
			"definition": schema.SingleNestedAttribute{
				Description: "Holiday rules for the calendar.",
				Required:    true,
				Attributes: map[string]schema.Attribute{
					"months": schema.MapNestedAttribute{
						Description: "Keys 1 through 12 contain rules for that month. Key 0 contains calculated-date rules. At most 366 rules are allowed across all months.",
						Required:    true,
						NestedObject: schema.NestedAttributeObject{Attributes: map[string]schema.Attribute{
							"calculated_rules": schema.ListNestedAttribute{
								Optional: true,
								NestedObject: schema.NestedAttributeObject{Attributes: map[string]schema.Attribute{
									"name": schema.StringAttribute{
										Description: "Optional rule name.",
										Optional:    true,
									},
									"observed": schema.StringAttribute{
										Description: "Optional function used to move the selected date when it falls on a weekend.",
										Optional:    true,
										Validators:  []validator.String{stringvalidator.OneOf("to_monday_if_sunday(date)", "to_monday_if_weekend(date)", "to_tuesday_if_sunday_or_monday_if_saturday(date)", "to_weekday_if_boxing_weekend(date)", "to_weekday_if_weekend(date)")},
									},
									"start_time": schema.StringAttribute{
										Description: "Inclusive start of a partial-day window in 24-hour HH:MM or HH:MM:SS format. Must be earlier than end_time.",
										Optional:    true,
									},
									"end_time": schema.StringAttribute{
										Description: "Exclusive end of a partial-day window in 24-hour HH:MM or HH:MM:SS format. Must be later than start_time.",
										Optional:    true,
									},
									"year_ranges": schema.SingleNestedAttribute{
										Description: "Optional inclusive year restriction.",
										Optional:    true,
										Attributes: map[string]schema.Attribute{
											"from": schema.Int64Attribute{
												Description: "First year in which the rule applies.",
												Optional:    true,
											},
											"until": schema.Int64Attribute{
												Description: "Last year in which the rule applies.",
												Optional:    true,
											},
											"between": schema.SingleNestedAttribute{
												Description: "Required start and end years.",
												Optional:    true,
												Attributes: map[string]schema.Attribute{
													"start": schema.Int64Attribute{
														Required: true,
													},
													"end": schema.Int64Attribute{
														Required: true,
													},
												},
											},
											"limited": schema.ListAttribute{
												Description: "Years in which the rule applies.",
												Optional:    true,
												ElementType: types.Int64Type,
												Validators:  []validator.List{listvalidator.SizeAtLeast(1)},
											},
										},
										Validators: []validator.Object{lib.ExactlyOneOfAttributes("from", "until", "between", "limited")},
									},
									"function": schema.StringAttribute{
										Description: "Supported calculated-date function.",
										Required:    true,
										Validators:  []validator.String{stringvalidator.OneOf("easter(year)", "orthodox_easter(year)", "orthodox_easter_julian(year)", "to_weekday_if_boxing_weekend_from_year(year)", "to_weekday_if_boxing_weekend_from_year_or_to_tuesday_if_monday(year)")},
									},
									"function_modifier": schema.Int64Attribute{
										Description: "Number of days to add to or subtract from the calculated date.",
										Optional:    true,
										Validators:  []validator.Int64{int64validator.Between(-366, 366)},
									},
								},
								},
								Validators: []validator.List{listvalidator.SizeAtMost(366)},
							},
							"fixed_rules": schema.ListNestedAttribute{
								Optional: true,
								NestedObject: schema.NestedAttributeObject{Attributes: map[string]schema.Attribute{
									"name": schema.StringAttribute{
										Description: "Optional rule name.",
										Optional:    true,
									},
									"observed": schema.StringAttribute{
										Description: "Optional function used to move the selected date when it falls on a weekend.",
										Optional:    true,
										Validators:  []validator.String{stringvalidator.OneOf("to_monday_if_sunday(date)", "to_monday_if_weekend(date)", "to_tuesday_if_sunday_or_monday_if_saturday(date)", "to_weekday_if_boxing_weekend(date)", "to_weekday_if_weekend(date)")},
									},
									"start_time": schema.StringAttribute{
										Description: "Inclusive start of a partial-day window in 24-hour HH:MM or HH:MM:SS format. Must be earlier than end_time.",
										Optional:    true,
									},
									"end_time": schema.StringAttribute{
										Description: "Exclusive end of a partial-day window in 24-hour HH:MM or HH:MM:SS format. Must be later than start_time.",
										Optional:    true,
									},
									"year_ranges": schema.SingleNestedAttribute{
										Description: "Optional inclusive year restriction.",
										Optional:    true,
										Attributes: map[string]schema.Attribute{
											"from": schema.Int64Attribute{
												Description: "First year in which the rule applies.",
												Optional:    true,
											},
											"until": schema.Int64Attribute{
												Description: "Last year in which the rule applies.",
												Optional:    true,
											},
											"between": schema.SingleNestedAttribute{
												Description: "Required start and end years.",
												Optional:    true,
												Attributes: map[string]schema.Attribute{
													"start": schema.Int64Attribute{
														Required: true,
													},
													"end": schema.Int64Attribute{
														Required: true,
													},
												},
											},
											"limited": schema.ListAttribute{
												Description: "Years in which the rule applies.",
												Optional:    true,
												ElementType: types.Int64Type,
												Validators:  []validator.List{listvalidator.SizeAtLeast(1)},
											},
										},
										Validators: []validator.Object{lib.ExactlyOneOfAttributes("from", "until", "between", "limited")},
									},
									"mday": schema.Int64Attribute{
										Description: "Day of the month. Must be valid for the selected month.",
										Required:    true,
										Validators:  []validator.Int64{int64validator.Between(1, 31)},
									},
								},
								},
								Validators: []validator.List{listvalidator.SizeAtMost(366)},
							},
							"weekday_rules": schema.ListNestedAttribute{
								Optional: true,
								NestedObject: schema.NestedAttributeObject{Attributes: map[string]schema.Attribute{
									"name": schema.StringAttribute{
										Description: "Optional rule name.",
										Optional:    true,
									},
									"observed": schema.StringAttribute{
										Description: "Optional function used to move the selected date when it falls on a weekend.",
										Optional:    true,
										Validators:  []validator.String{stringvalidator.OneOf("to_monday_if_sunday(date)", "to_monday_if_weekend(date)", "to_tuesday_if_sunday_or_monday_if_saturday(date)", "to_weekday_if_boxing_weekend(date)", "to_weekday_if_weekend(date)")},
									},
									"start_time": schema.StringAttribute{
										Description: "Inclusive start of a partial-day window in 24-hour HH:MM or HH:MM:SS format. Must be earlier than end_time.",
										Optional:    true,
									},
									"end_time": schema.StringAttribute{
										Description: "Exclusive end of a partial-day window in 24-hour HH:MM or HH:MM:SS format. Must be later than start_time.",
										Optional:    true,
									},
									"year_ranges": schema.SingleNestedAttribute{
										Description: "Optional inclusive year restriction.",
										Optional:    true,
										Attributes: map[string]schema.Attribute{
											"from": schema.Int64Attribute{
												Description: "First year in which the rule applies.",
												Optional:    true,
											},
											"until": schema.Int64Attribute{
												Description: "Last year in which the rule applies.",
												Optional:    true,
											},
											"between": schema.SingleNestedAttribute{
												Description: "Required start and end years.",
												Optional:    true,
												Attributes: map[string]schema.Attribute{
													"start": schema.Int64Attribute{
														Required: true,
													},
													"end": schema.Int64Attribute{
														Required: true,
													},
												},
											},
											"limited": schema.ListAttribute{
												Description: "Years in which the rule applies.",
												Optional:    true,
												ElementType: types.Int64Type,
												Validators:  []validator.List{listvalidator.SizeAtLeast(1)},
											},
										},
										Validators: []validator.Object{lib.ExactlyOneOfAttributes("from", "until", "between", "limited")},
									},
									"week": schema.Int64Attribute{
										Description: "Occurrence of the weekday in the month. Negative values count from the end of the month.",
										Required:    true,
										Validators:  []validator.Int64{int64validator.OneOf(1, 2, 3, 4, 5, -1, -2, -3)},
									},
									"wday": schema.Int64Attribute{
										Description: "Day of the week, where 0 is Sunday and 6 is Saturday.",
										Required:    true,
										Validators:  []validator.Int64{int64validator.Between(0, 6)},
									},
								},
								},
								Validators: []validator.List{listvalidator.SizeAtMost(366)},
							},
						},
							Validators: []validator.Object{lib.AtLeastOneOfAttributes("calculated_rules", "fixed_rules", "weekday_rules")},
						},
						Validators: []validator.Map{mapvalidator.KeysAre(stringvalidator.OneOf("0", "1", "2", "3", "4", "5", "6", "7", "8", "9", "10", "11", "12"))},
					},
				},
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
		Version: 1,
	}
}

func (r *holidayCalendarResource) UpgradeState(ctx context.Context) map[int64]resource.StateUpgrader {
	return map[int64]resource.StateUpgrader{
		0: {
			PriorSchema: &schema.Schema{
				Description: "A Holiday Calendar defines site-wide holiday dates and optional partial-day windows that scheduled resources skip.",
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
				definitionValue, conversionDiags := lib.DynamicToInterface(ctx, path.Root("definition"), priorState.Definition)
				resp.Diagnostics.Append(conversionDiags...)
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

func (r *holidayCalendarResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	var plan holidayCalendarResourceModel
	diags := req.Plan.Get(ctx, &plan)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}
	var config holidayCalendarResourceModel
	diags = req.Config.Get(ctx, &config)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	paramsHolidayCalendarCreate := files_sdk.HolidayCalendarCreateParams{}
	createDefinition, diags := lib.SchemaAttributeToInterface(ctx, path.Root("definition"), plan.Definition)
	resp.Diagnostics.Append(diags...)
	createDefinition, transformDiags0 := lib.UngroupStructuralUnionAtPath(ctx, path.Root("definition"), createDefinition, []string{"months"}, []lib.JSONSchemaVariant{{Name: "calculated_rules", Required: []string{"function"}, Allowed: []string{"name", "observed", "start_time", "end_time", "year_ranges", "function", "function_modifier"}}, {Name: "fixed_rules", Required: []string{"mday"}, Allowed: []string{"name", "observed", "start_time", "end_time", "year_ranges", "mday"}}, {Name: "weekday_rules", Required: []string{"week", "wday"}, Allowed: []string{"name", "observed", "start_time", "end_time", "year_ranges", "week", "wday"}}})
	resp.Diagnostics.Append(transformDiags0...)
	paramsHolidayCalendarCreate.Definition = createDefinition
	paramsHolidayCalendarCreate.Name = plan.Name.ValueString()

	if resp.Diagnostics.HasError() {
		return
	}

	holidayCalendar, err := r.client.Create(paramsHolidayCalendarCreate, files_sdk.WithContext(ctx))
	if err != nil {
		resp.Diagnostics.AddError(
			"Error Creating Files HolidayCalendar",
			"Could not create holiday_calendar, unexpected error: "+err.Error(),
		)
		return
	}

	diags = r.populateResourceModel(ctx, holidayCalendar, &plan)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	diags = resp.State.Set(ctx, plan)
	resp.Diagnostics.Append(diags...)
}

func (r *holidayCalendarResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	var state holidayCalendarResourceModel
	diags := req.State.Get(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	paramsHolidayCalendarFind := files_sdk.HolidayCalendarFindParams{}
	paramsHolidayCalendarFind.Id = state.Id.ValueInt64()

	holidayCalendar, err := r.client.Find(paramsHolidayCalendarFind, files_sdk.WithContext(ctx))
	if err != nil {
		if files_sdk.IsNotExist(err) {
			resp.State.RemoveResource(ctx)
			return
		}

		resp.Diagnostics.AddError(
			"Error Reading Files HolidayCalendar",
			"Could not read holiday_calendar id "+fmt.Sprint(state.Id.ValueInt64())+": "+err.Error(),
		)
		return
	}

	diags = r.populateResourceModel(ctx, holidayCalendar, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	diags = resp.State.Set(ctx, state)
	resp.Diagnostics.Append(diags...)
}

func (r *holidayCalendarResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	var plan holidayCalendarResourceModel
	diags := req.Plan.Get(ctx, &plan)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}
	var config holidayCalendarResourceModel
	diags = req.Config.Get(ctx, &config)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	paramsHolidayCalendarUpdate := map[string]interface{}{}
	if !plan.Id.IsNull() && !plan.Id.IsUnknown() {
		paramsHolidayCalendarUpdate["id"] = plan.Id.ValueInt64()
	}
	if !config.Definition.IsNull() && !config.Definition.IsUnknown() {
		updateDefinition, diags := lib.SchemaAttributeToInterface(ctx, path.Root("definition"), config.Definition)
		resp.Diagnostics.Append(diags...)
		updateDefinition, transformDiags0 := lib.UngroupStructuralUnionAtPath(ctx, path.Root("definition"), updateDefinition, []string{"months"}, []lib.JSONSchemaVariant{{Name: "calculated_rules", Required: []string{"function"}, Allowed: []string{"name", "observed", "start_time", "end_time", "year_ranges", "function", "function_modifier"}}, {Name: "fixed_rules", Required: []string{"mday"}, Allowed: []string{"name", "observed", "start_time", "end_time", "year_ranges", "mday"}}, {Name: "weekday_rules", Required: []string{"week", "wday"}, Allowed: []string{"name", "observed", "start_time", "end_time", "year_ranges", "week", "wday"}}})
		resp.Diagnostics.Append(transformDiags0...)
		paramsHolidayCalendarUpdate["definition"] = updateDefinition
	}
	if !config.Name.IsNull() && !config.Name.IsUnknown() {
		paramsHolidayCalendarUpdate["name"] = config.Name.ValueString()
	}

	if resp.Diagnostics.HasError() {
		return
	}

	holidayCalendar, err := r.client.UpdateWithMap(paramsHolidayCalendarUpdate, files_sdk.WithContext(ctx))
	if err != nil {
		resp.Diagnostics.AddError(
			"Error Updating Files HolidayCalendar",
			"Could not update holiday_calendar, unexpected error: "+err.Error(),
		)
		return
	}

	diags = r.populateResourceModel(ctx, holidayCalendar, &plan)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	diags = resp.State.Set(ctx, plan)
	resp.Diagnostics.Append(diags...)
}

func (r *holidayCalendarResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	var state holidayCalendarResourceModel
	diags := req.State.Get(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	paramsHolidayCalendarDelete := files_sdk.HolidayCalendarDeleteParams{}
	paramsHolidayCalendarDelete.Id = state.Id.ValueInt64()

	err := r.client.Delete(paramsHolidayCalendarDelete, files_sdk.WithContext(ctx))
	if err != nil && !files_sdk.IsNotExist(err) {
		resp.Diagnostics.AddError(
			"Error Deleting Files HolidayCalendar",
			"Could not delete holiday_calendar id "+fmt.Sprint(state.Id.ValueInt64())+": "+err.Error(),
		)
	}
}

func (r *holidayCalendarResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	idParts := strings.SplitN(req.ID, ",", 1)

	if len(idParts) != 1 || idParts[0] == "" {
		resp.Diagnostics.AddError(
			"Unexpected Import Identifier",
			fmt.Sprintf("Expected import identifier with format: id. Got: %q", req.ID),
		)
		return
	}

	idPart, err := strconv.ParseFloat(idParts[0], 64)
	if err != nil {
		resp.Diagnostics.AddError(
			"Error Parsing ID",
			"Could not parse id: "+err.Error(),
		)
		return
	}
	resp.Diagnostics.Append(resp.State.SetAttribute(ctx, path.Root("id"), idPart)...)

}

func (r *holidayCalendarResource) populateResourceModel(ctx context.Context, holidayCalendar files_sdk.HolidayCalendar, state *holidayCalendarResourceModel) (diags diag.Diagnostics) {
	var propDiags diag.Diagnostics

	state.Id = types.Int64Value(holidayCalendar.Id)
	state.Name = types.StringValue(holidayCalendar.Name)
	definitionValue := interface{}(holidayCalendar.Definition)
	definitionValue, transformDiags0 := lib.GroupStructuralUnionAtPath(ctx, path.Root("definition"), definitionValue, []string{"months"}, []lib.JSONSchemaVariant{{Name: "calculated_rules", Required: []string{"function"}, Allowed: []string{"name", "observed", "start_time", "end_time", "year_ranges", "function", "function_modifier"}}, {Name: "fixed_rules", Required: []string{"mday"}, Allowed: []string{"name", "observed", "start_time", "end_time", "year_ranges", "mday"}}, {Name: "weekday_rules", Required: []string{"week", "wday"}, Allowed: []string{"name", "observed", "start_time", "end_time", "year_ranges", "week", "wday"}}})
	diags.Append(transformDiags0...)
	state.Definition, propDiags = lib.ToObject(ctx, path.Root("definition"), definitionValue, state.Definition)
	diags.Append(propDiags...)
	if err := lib.TimeToStringType(ctx, path.Root("created_at"), holidayCalendar.CreatedAt, &state.CreatedAt); err != nil {
		diags.AddError(
			"Error Creating Files HolidayCalendar",
			"Could not convert state created_at to string: "+err.Error(),
		)
	}
	if err := lib.TimeToStringType(ctx, path.Root("updated_at"), holidayCalendar.UpdatedAt, &state.UpdatedAt); err != nil {
		diags.AddError(
			"Error Creating Files HolidayCalendar",
			"Could not convert state updated_at to string: "+err.Error(),
		)
	}

	return
}
