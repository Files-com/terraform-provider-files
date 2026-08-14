package provider

import (
	"context"
	"fmt"

	files_sdk "github.com/Files-com/files-sdk-go/v3"
	schedule "github.com/Files-com/files-sdk-go/v3/schedule"
	"github.com/Files-com/terraform-provider-files/lib"

	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/datasource/schema"
	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

var (
	_ datasource.DataSource              = &scheduleDataSource{}
	_ datasource.DataSourceWithConfigure = &scheduleDataSource{}
)

func NewScheduleDataSource() datasource.DataSource {
	return &scheduleDataSource{}
}

type scheduleDataSource struct {
	client *schedule.Client
}

type scheduleDataSourceModel struct {
	Id                    types.Int64  `tfsdk:"id"`
	Name                  types.String `tfsdk:"name"`
	ScheduleDaysOfWeek    types.List   `tfsdk:"schedule_days_of_week"`
	ScheduleTimesOfDay    types.List   `tfsdk:"schedule_times_of_day"`
	ScheduleTimeZone      types.String `tfsdk:"schedule_time_zone"`
	HolidayRegion         types.String `tfsdk:"holiday_region"`
	HumanReadableSchedule types.String `tfsdk:"human_readable_schedule"`
	CreatedAt             types.String `tfsdk:"created_at"`
	UpdatedAt             types.String `tfsdk:"updated_at"`
}

func (r *scheduleDataSource) Configure(_ context.Context, req datasource.ConfigureRequest, resp *datasource.ConfigureResponse) {
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

	r.client = &schedule.Client{Config: sdk_config}
}

func (r *scheduleDataSource) Metadata(_ context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_schedule"
}

func (r *scheduleDataSource) Schema(_ context.Context, _ datasource.SchemaRequest, resp *datasource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Description: "A Schedule is a named, reusable weekday-and-time schedule shared by scheduled resources across a Site.",
		Attributes: map[string]schema.Attribute{
			"id": schema.Int64Attribute{
				Description: "Schedule ID.",
				Required:    true,
			},
			"name": schema.StringAttribute{
				Description: "Schedule name.",
				Computed:    true,
			},
			"schedule_days_of_week": schema.ListAttribute{
				Description: "0-based weekdays used by the Schedule. 0 is Sunday.",
				Computed:    true,
				ElementType: types.Int64Type,
			},
			"schedule_times_of_day": schema.ListAttribute{
				Description: "Times of day in HH:MM format (24-hour).",
				Computed:    true,
				ElementType: types.StringType,
			},
			"schedule_time_zone": schema.StringAttribute{
				Description: "Time zone for scheduled times. If not set, times are interpreted as UTC.",
				Computed:    true,
			},
			"holiday_region": schema.StringAttribute{
				Description: "Optional holiday region on which linked resources do not run.",
				Computed:    true,
			},
			"human_readable_schedule": schema.StringAttribute{
				Description: "Human-readable Schedule description.",
				Computed:    true,
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

func (r *scheduleDataSource) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
	var data scheduleDataSourceModel
	diags := req.Config.Get(ctx, &data)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	paramsScheduleFind := files_sdk.ScheduleFindParams{}
	paramsScheduleFind.Id = data.Id.ValueInt64()

	schedule, err := r.client.Find(paramsScheduleFind, files_sdk.WithContext(ctx))
	if err != nil {
		resp.Diagnostics.AddError(
			"Error Reading Files Schedule",
			"Could not read schedule id "+fmt.Sprint(data.Id.ValueInt64())+": "+err.Error(),
		)
		return
	}

	diags = r.populateDataSourceModel(ctx, schedule, &data)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	diags = resp.State.Set(ctx, data)
	resp.Diagnostics.Append(diags...)
}

func (r *scheduleDataSource) populateDataSourceModel(ctx context.Context, schedule files_sdk.Schedule, state *scheduleDataSourceModel) (diags diag.Diagnostics) {
	var propDiags diag.Diagnostics

	state.Id = types.Int64Value(schedule.Id)
	state.Name = types.StringValue(schedule.Name)
	state.ScheduleDaysOfWeek, propDiags = types.ListValueFrom(ctx, types.Int64Type, schedule.ScheduleDaysOfWeek)
	diags.Append(propDiags...)
	state.ScheduleTimesOfDay, propDiags = types.ListValueFrom(ctx, types.StringType, schedule.ScheduleTimesOfDay)
	diags.Append(propDiags...)
	state.ScheduleTimeZone = types.StringValue(schedule.ScheduleTimeZone)
	state.HolidayRegion = types.StringValue(schedule.HolidayRegion)
	state.HumanReadableSchedule = types.StringValue(schedule.HumanReadableSchedule)
	if err := lib.TimeToStringType(ctx, path.Root("created_at"), schedule.CreatedAt, &state.CreatedAt); err != nil {
		diags.AddError(
			"Error Creating Files Schedule",
			"Could not convert state created_at to string: "+err.Error(),
		)
	}
	if err := lib.TimeToStringType(ctx, path.Root("updated_at"), schedule.UpdatedAt, &state.UpdatedAt); err != nil {
		diags.AddError(
			"Error Creating Files Schedule",
			"Could not convert state updated_at to string: "+err.Error(),
		)
	}

	return
}
