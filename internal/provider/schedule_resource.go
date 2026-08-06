package provider

import (
	"context"
	"fmt"
	"strconv"
	"strings"

	files_sdk "github.com/Files-com/files-sdk-go/v3"
	schedule "github.com/Files-com/files-sdk-go/v3/schedule"
	"github.com/Files-com/terraform-provider-files/lib"
	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/int64planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

var (
	_ resource.Resource                = &scheduleResource{}
	_ resource.ResourceWithConfigure   = &scheduleResource{}
	_ resource.ResourceWithImportState = &scheduleResource{}
)

func NewScheduleResource() resource.Resource {
	return &scheduleResource{}
}

type scheduleResource struct {
	client *schedule.Client
}

type scheduleResourceModel struct {
	Name                  types.String `tfsdk:"name"`
	ScheduleDaysOfWeek    types.List   `tfsdk:"schedule_days_of_week"`
	ScheduleTimesOfDay    types.List   `tfsdk:"schedule_times_of_day"`
	ScheduleTimeZone      types.String `tfsdk:"schedule_time_zone"`
	HolidayRegion         types.String `tfsdk:"holiday_region"`
	Id                    types.Int64  `tfsdk:"id"`
	HumanReadableSchedule types.String `tfsdk:"human_readable_schedule"`
	CreatedAt             types.String `tfsdk:"created_at"`
	UpdatedAt             types.String `tfsdk:"updated_at"`
}

func (r *scheduleResource) Configure(_ context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {
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

func (r *scheduleResource) Metadata(_ context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_schedule"
}

func (r *scheduleResource) Schema(_ context.Context, _ resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Description: "A Schedule is a named, reusable weekday-and-time schedule shared by scheduled resources across a Site.",
		Attributes: map[string]schema.Attribute{
			"name": schema.StringAttribute{
				Description: "Schedule name.",
				Required:    true,
			},
			"schedule_days_of_week": schema.ListAttribute{
				Description: "0-based weekdays used by the Schedule. 0 is Sunday.",
				Required:    true,
				ElementType: types.Int64Type,
			},
			"schedule_times_of_day": schema.ListAttribute{
				Description: "Times of day in HH:MM format (24-hour).",
				Required:    true,
				ElementType: types.StringType,
			},
			"schedule_time_zone": schema.StringAttribute{
				Description: "Time zone for scheduled times. If not set, times are interpreted as UTC.",
				Computed:    true,
				Optional:    true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
			},
			"holiday_region": schema.StringAttribute{
				Description: "Optional holiday region on which linked resources do not run.",
				Computed:    true,
				Optional:    true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
			},
			"id": schema.Int64Attribute{
				Description: "Schedule ID.",
				Computed:    true,
				PlanModifiers: []planmodifier.Int64{
					int64planmodifier.UseStateForUnknown(),
				},
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

func (r *scheduleResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	var plan scheduleResourceModel
	diags := req.Plan.Get(ctx, &plan)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}
	var config scheduleResourceModel
	diags = req.Config.Get(ctx, &config)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	paramsScheduleCreate := files_sdk.ScheduleCreateParams{}
	paramsScheduleCreate.Name = plan.Name.ValueString()
	if !plan.ScheduleDaysOfWeek.IsNull() && !plan.ScheduleDaysOfWeek.IsUnknown() {
		diags = plan.ScheduleDaysOfWeek.ElementsAs(ctx, &paramsScheduleCreate.ScheduleDaysOfWeek, false)
		resp.Diagnostics.Append(diags...)
	}
	if !plan.ScheduleTimesOfDay.IsNull() && !plan.ScheduleTimesOfDay.IsUnknown() {
		diags = plan.ScheduleTimesOfDay.ElementsAs(ctx, &paramsScheduleCreate.ScheduleTimesOfDay, false)
		resp.Diagnostics.Append(diags...)
	}
	paramsScheduleCreate.ScheduleTimeZone = plan.ScheduleTimeZone.ValueString()
	paramsScheduleCreate.HolidayRegion = plan.HolidayRegion.ValueString()

	if resp.Diagnostics.HasError() {
		return
	}

	schedule, err := r.client.Create(paramsScheduleCreate, files_sdk.WithContext(ctx))
	if err != nil {
		resp.Diagnostics.AddError(
			"Error Creating Files Schedule",
			"Could not create schedule, unexpected error: "+err.Error(),
		)
		return
	}

	diags = r.populateResourceModel(ctx, schedule, &plan)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	diags = resp.State.Set(ctx, plan)
	resp.Diagnostics.Append(diags...)
}

func (r *scheduleResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	var state scheduleResourceModel
	diags := req.State.Get(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	paramsScheduleFind := files_sdk.ScheduleFindParams{}
	paramsScheduleFind.Id = state.Id.ValueInt64()

	schedule, err := r.client.Find(paramsScheduleFind, files_sdk.WithContext(ctx))
	if err != nil {
		if files_sdk.IsNotExist(err) {
			resp.State.RemoveResource(ctx)
			return
		}

		resp.Diagnostics.AddError(
			"Error Reading Files Schedule",
			"Could not read schedule id "+fmt.Sprint(state.Id.ValueInt64())+": "+err.Error(),
		)
		return
	}

	diags = r.populateResourceModel(ctx, schedule, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	diags = resp.State.Set(ctx, state)
	resp.Diagnostics.Append(diags...)
}

func (r *scheduleResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	var plan scheduleResourceModel
	diags := req.Plan.Get(ctx, &plan)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}
	var config scheduleResourceModel
	diags = req.Config.Get(ctx, &config)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	paramsScheduleUpdate := map[string]interface{}{}
	if !plan.Id.IsNull() && !plan.Id.IsUnknown() {
		paramsScheduleUpdate["id"] = plan.Id.ValueInt64()
	}
	if !config.Name.IsNull() && !config.Name.IsUnknown() {
		paramsScheduleUpdate["name"] = config.Name.ValueString()
	}
	if !config.ScheduleDaysOfWeek.IsNull() && !config.ScheduleDaysOfWeek.IsUnknown() {
		var updateScheduleDaysOfWeek []int64
		diags = config.ScheduleDaysOfWeek.ElementsAs(ctx, &updateScheduleDaysOfWeek, false)
		resp.Diagnostics.Append(diags...)
		paramsScheduleUpdate["schedule_days_of_week"] = updateScheduleDaysOfWeek
	}
	if !config.ScheduleTimesOfDay.IsNull() && !config.ScheduleTimesOfDay.IsUnknown() {
		var updateScheduleTimesOfDay []string
		diags = config.ScheduleTimesOfDay.ElementsAs(ctx, &updateScheduleTimesOfDay, false)
		resp.Diagnostics.Append(diags...)
		paramsScheduleUpdate["schedule_times_of_day"] = updateScheduleTimesOfDay
	}
	if !config.ScheduleTimeZone.IsNull() && !config.ScheduleTimeZone.IsUnknown() {
		paramsScheduleUpdate["schedule_time_zone"] = config.ScheduleTimeZone.ValueString()
	}
	if !config.HolidayRegion.IsNull() && !config.HolidayRegion.IsUnknown() {
		paramsScheduleUpdate["holiday_region"] = config.HolidayRegion.ValueString()
	}

	if resp.Diagnostics.HasError() {
		return
	}

	schedule, err := r.client.UpdateWithMap(paramsScheduleUpdate, files_sdk.WithContext(ctx))
	if err != nil {
		resp.Diagnostics.AddError(
			"Error Updating Files Schedule",
			"Could not update schedule, unexpected error: "+err.Error(),
		)
		return
	}

	diags = r.populateResourceModel(ctx, schedule, &plan)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	diags = resp.State.Set(ctx, plan)
	resp.Diagnostics.Append(diags...)
}

func (r *scheduleResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	var state scheduleResourceModel
	diags := req.State.Get(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	paramsScheduleDelete := files_sdk.ScheduleDeleteParams{}
	paramsScheduleDelete.Id = state.Id.ValueInt64()

	err := r.client.Delete(paramsScheduleDelete, files_sdk.WithContext(ctx))
	if err != nil && !files_sdk.IsNotExist(err) {
		resp.Diagnostics.AddError(
			"Error Deleting Files Schedule",
			"Could not delete schedule id "+fmt.Sprint(state.Id.ValueInt64())+": "+err.Error(),
		)
	}
}

func (r *scheduleResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
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

func (r *scheduleResource) populateResourceModel(ctx context.Context, schedule files_sdk.Schedule, state *scheduleResourceModel) (diags diag.Diagnostics) {
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
