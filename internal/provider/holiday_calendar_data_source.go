package provider

import (
	"context"
	"fmt"

	files_sdk "github.com/Files-com/files-sdk-go/v3"
	holiday_calendar "github.com/Files-com/files-sdk-go/v3/holidaycalendar"
	"github.com/Files-com/terraform-provider-files/lib"
	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/datasource/schema"
	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

var (
	_ datasource.DataSource              = &holidayCalendarDataSource{}
	_ datasource.DataSourceWithConfigure = &holidayCalendarDataSource{}
)

func NewHolidayCalendarDataSource() datasource.DataSource {
	return &holidayCalendarDataSource{}
}

type holidayCalendarDataSource struct {
	client *holiday_calendar.Client
}

type holidayCalendarDataSourceModel struct {
	Id         types.Int64  `tfsdk:"id"`
	Name       types.String `tfsdk:"name"`
	Definition types.Object `tfsdk:"definition"`
	CreatedAt  types.String `tfsdk:"created_at"`
	UpdatedAt  types.String `tfsdk:"updated_at"`
}

func (r *holidayCalendarDataSource) Configure(_ context.Context, req datasource.ConfigureRequest, resp *datasource.ConfigureResponse) {
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

func (r *holidayCalendarDataSource) Metadata(_ context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_holiday_calendar"
}

func (r *holidayCalendarDataSource) Schema(_ context.Context, _ datasource.SchemaRequest, resp *datasource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Description: "A Holiday Calendar defines site-wide holiday dates and optional partial-day windows that scheduled resources skip.",
		Attributes: map[string]schema.Attribute{
			"id": schema.Int64Attribute{
				Description: "Holiday Calendar ID. Set a scheduled resource's `holiday_region` to `custom_` followed by this ID to make it skip the days in this calendar.",
				Required:    true,
			},
			"name": schema.StringAttribute{
				Description: "Holiday Calendar name.",
				Computed:    true,
			},
			"definition": schema.SingleNestedAttribute{
				Description: "Holiday rules for the calendar.",
				Computed:    true,
				Attributes: map[string]schema.Attribute{
					"months": schema.MapNestedAttribute{
						Description: "Keys 1 through 12 contain rules for that month. Key 0 contains calculated-date rules. At most 366 rules are allowed across all months.",
						Computed:    true,
						NestedObject: schema.NestedAttributeObject{Attributes: map[string]schema.Attribute{
							"calculated_rules": schema.ListNestedAttribute{
								Computed: true,
								NestedObject: schema.NestedAttributeObject{Attributes: map[string]schema.Attribute{
									"name": schema.StringAttribute{
										Description: "Optional rule name.",
										Computed:    true,
									},
									"observed": schema.StringAttribute{
										Description: "Optional function used to move the selected date when it falls on a weekend.",
										Computed:    true,
									},
									"start_time": schema.StringAttribute{
										Description: "Inclusive start of a partial-day window in 24-hour HH:MM or HH:MM:SS format. Must be earlier than end_time.",
										Computed:    true,
									},
									"end_time": schema.StringAttribute{
										Description: "Exclusive end of a partial-day window in 24-hour HH:MM or HH:MM:SS format. Must be later than start_time.",
										Computed:    true,
									},
									"year_ranges": schema.SingleNestedAttribute{
										Description: "Optional inclusive year restriction.",
										Computed:    true,
										Attributes: map[string]schema.Attribute{
											"from": schema.Int64Attribute{
												Description: "First year in which the rule applies.",
												Computed:    true,
											},
											"until": schema.Int64Attribute{
												Description: "Last year in which the rule applies.",
												Computed:    true,
											},
											"between": schema.SingleNestedAttribute{
												Description: "Required start and end years.",
												Computed:    true,
												Attributes: map[string]schema.Attribute{
													"start": schema.Int64Attribute{
														Computed: true,
													},
													"end": schema.Int64Attribute{
														Computed: true,
													},
												},
											},
											"limited": schema.ListAttribute{
												Description: "Years in which the rule applies.",
												Computed:    true,
												ElementType: types.Int64Type,
											},
										},
									},
									"function": schema.StringAttribute{
										Description: "Supported calculated-date function.",
										Computed:    true,
									},
									"function_modifier": schema.Int64Attribute{
										Description: "Number of days to add to or subtract from the calculated date.",
										Computed:    true,
									},
								},
								},
							},
							"fixed_rules": schema.ListNestedAttribute{
								Computed: true,
								NestedObject: schema.NestedAttributeObject{Attributes: map[string]schema.Attribute{
									"name": schema.StringAttribute{
										Description: "Optional rule name.",
										Computed:    true,
									},
									"observed": schema.StringAttribute{
										Description: "Optional function used to move the selected date when it falls on a weekend.",
										Computed:    true,
									},
									"start_time": schema.StringAttribute{
										Description: "Inclusive start of a partial-day window in 24-hour HH:MM or HH:MM:SS format. Must be earlier than end_time.",
										Computed:    true,
									},
									"end_time": schema.StringAttribute{
										Description: "Exclusive end of a partial-day window in 24-hour HH:MM or HH:MM:SS format. Must be later than start_time.",
										Computed:    true,
									},
									"year_ranges": schema.SingleNestedAttribute{
										Description: "Optional inclusive year restriction.",
										Computed:    true,
										Attributes: map[string]schema.Attribute{
											"from": schema.Int64Attribute{
												Description: "First year in which the rule applies.",
												Computed:    true,
											},
											"until": schema.Int64Attribute{
												Description: "Last year in which the rule applies.",
												Computed:    true,
											},
											"between": schema.SingleNestedAttribute{
												Description: "Required start and end years.",
												Computed:    true,
												Attributes: map[string]schema.Attribute{
													"start": schema.Int64Attribute{
														Computed: true,
													},
													"end": schema.Int64Attribute{
														Computed: true,
													},
												},
											},
											"limited": schema.ListAttribute{
												Description: "Years in which the rule applies.",
												Computed:    true,
												ElementType: types.Int64Type,
											},
										},
									},
									"mday": schema.Int64Attribute{
										Description: "Day of the month. Must be valid for the selected month.",
										Computed:    true,
									},
								},
								},
							},
							"weekday_rules": schema.ListNestedAttribute{
								Computed: true,
								NestedObject: schema.NestedAttributeObject{Attributes: map[string]schema.Attribute{
									"name": schema.StringAttribute{
										Description: "Optional rule name.",
										Computed:    true,
									},
									"observed": schema.StringAttribute{
										Description: "Optional function used to move the selected date when it falls on a weekend.",
										Computed:    true,
									},
									"start_time": schema.StringAttribute{
										Description: "Inclusive start of a partial-day window in 24-hour HH:MM or HH:MM:SS format. Must be earlier than end_time.",
										Computed:    true,
									},
									"end_time": schema.StringAttribute{
										Description: "Exclusive end of a partial-day window in 24-hour HH:MM or HH:MM:SS format. Must be later than start_time.",
										Computed:    true,
									},
									"year_ranges": schema.SingleNestedAttribute{
										Description: "Optional inclusive year restriction.",
										Computed:    true,
										Attributes: map[string]schema.Attribute{
											"from": schema.Int64Attribute{
												Description: "First year in which the rule applies.",
												Computed:    true,
											},
											"until": schema.Int64Attribute{
												Description: "Last year in which the rule applies.",
												Computed:    true,
											},
											"between": schema.SingleNestedAttribute{
												Description: "Required start and end years.",
												Computed:    true,
												Attributes: map[string]schema.Attribute{
													"start": schema.Int64Attribute{
														Computed: true,
													},
													"end": schema.Int64Attribute{
														Computed: true,
													},
												},
											},
											"limited": schema.ListAttribute{
												Description: "Years in which the rule applies.",
												Computed:    true,
												ElementType: types.Int64Type,
											},
										},
									},
									"week": schema.Int64Attribute{
										Description: "Occurrence of the weekday in the month. Negative values count from the end of the month.",
										Computed:    true,
									},
									"wday": schema.Int64Attribute{
										Description: "Day of the week, where 0 is Sunday and 6 is Saturday.",
										Computed:    true,
									},
								},
								},
							},
						},
						},
					},
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
	}
}

func (r *holidayCalendarDataSource) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
	var data holidayCalendarDataSourceModel
	diags := req.Config.Get(ctx, &data)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	paramsHolidayCalendarFind := files_sdk.HolidayCalendarFindParams{}
	paramsHolidayCalendarFind.Id = data.Id.ValueInt64()

	holidayCalendar, err := r.client.Find(paramsHolidayCalendarFind, files_sdk.WithContext(ctx))
	if err != nil {
		resp.Diagnostics.AddError(
			"Error Reading Files HolidayCalendar",
			"Could not read holiday_calendar id "+fmt.Sprint(data.Id.ValueInt64())+": "+err.Error(),
		)
		return
	}

	diags = r.populateDataSourceModel(ctx, holidayCalendar, &data)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	diags = resp.State.Set(ctx, data)
	resp.Diagnostics.Append(diags...)
}

func (r *holidayCalendarDataSource) populateDataSourceModel(ctx context.Context, holidayCalendar files_sdk.HolidayCalendar, state *holidayCalendarDataSourceModel) (diags diag.Diagnostics) {
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
