package provider

import (
	"context"
	"fmt"
	"strconv"
	"strings"

	files_sdk "github.com/Files-com/files-sdk-go/v3"
	scheduled_export "github.com/Files-com/files-sdk-go/v3/scheduledexport"
	"github.com/Files-com/terraform-provider-files/lib"
	"github.com/hashicorp/terraform-plugin-framework-validators/stringvalidator"
	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/boolplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/dynamicplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/int64planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/listplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/schema/validator"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

var (
	_ resource.Resource                = &scheduledExportResource{}
	_ resource.ResourceWithConfigure   = &scheduledExportResource{}
	_ resource.ResourceWithImportState = &scheduledExportResource{}
)

func NewScheduledExportResource() resource.Resource {
	return &scheduledExportResource{}
}

type scheduledExportResource struct {
	client *scheduled_export.Client
}

type scheduledExportResourceModel struct {
	Name                  types.String  `tfsdk:"name"`
	ExportType            types.String  `tfsdk:"export_type"`
	ExportOptions         types.Dynamic `tfsdk:"export_options"`
	UserId                types.Int64   `tfsdk:"user_id"`
	Disabled              types.Bool    `tfsdk:"disabled"`
	Trigger               types.String  `tfsdk:"trigger"`
	Interval              types.String  `tfsdk:"interval"`
	RecurringDay          types.Int64   `tfsdk:"recurring_day"`
	ScheduleDaysOfWeek    types.List    `tfsdk:"schedule_days_of_week"`
	ScheduleTimesOfDay    types.List    `tfsdk:"schedule_times_of_day"`
	ScheduleTimeZone      types.String  `tfsdk:"schedule_time_zone"`
	HolidayRegion         types.String  `tfsdk:"holiday_region"`
	Id                    types.Int64   `tfsdk:"id"`
	ReportName            types.String  `tfsdk:"report_name"`
	HumanReadableSchedule types.String  `tfsdk:"human_readable_schedule"`
	LastRunAt             types.String  `tfsdk:"last_run_at"`
	LastExportId          types.Int64   `tfsdk:"last_export_id"`
	CreatedAt             types.String  `tfsdk:"created_at"`
	UpdatedAt             types.String  `tfsdk:"updated_at"`
}

func (r *scheduledExportResource) Configure(_ context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {
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

	r.client = &scheduled_export.Client{Config: sdk_config}
}

func (r *scheduledExportResource) Metadata(_ context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_scheduled_export"
}

func (r *scheduledExportResource) Schema(_ context.Context, _ resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Description: "A Scheduled Export defines a recurring schedule for generating one of the built-in CSV exports and e-mailing it to a Site Admin recipient.",
		Attributes: map[string]schema.Attribute{
			"name": schema.StringAttribute{
				Description: "Name for this scheduled export.",
				Required:    true,
			},
			"export_type": schema.StringAttribute{
				Description: "Export report type. Valid values: folder_size_audit, group_membership_audit, permission_audit, share_link_audit",
				Required:    true,
			},
			"export_options": schema.DynamicAttribute{
				Description: "Report-specific options. `permission_audit` supports `group_by` with `user` or `path`.",
				Computed:    true,
				Optional:    true,
				PlanModifiers: []planmodifier.Dynamic{
					dynamicplanmodifier.UseStateForUnknown(),
				},
			},
			"user_id": schema.Int64Attribute{
				Description: "Site Admin user who receives the completed export e-mail.",
				Computed:    true,
				Optional:    true,
				PlanModifiers: []planmodifier.Int64{
					int64planmodifier.UseStateForUnknown(),
				},
			},
			"disabled": schema.BoolAttribute{
				Description: "If true, this scheduled export will not run.",
				Computed:    true,
				Optional:    true,
				PlanModifiers: []planmodifier.Bool{
					boolplanmodifier.UseStateForUnknown(),
				},
			},
			"trigger": schema.StringAttribute{
				Description: "Schedule trigger type: `daily` or `custom_schedule`.",
				Computed:    true,
				Optional:    true,
				Validators: []validator.String{
					stringvalidator.OneOf("daily", "custom_schedule"),
				},
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
			},
			"interval": schema.StringAttribute{
				Description: "If trigger is `daily`, this specifies how often to run the scheduled export.",
				Computed:    true,
				Optional:    true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
			},
			"recurring_day": schema.Int64Attribute{
				Description: "If trigger is `daily`, this selects the day number inside the chosen interval.",
				Computed:    true,
				Optional:    true,
				PlanModifiers: []planmodifier.Int64{
					int64planmodifier.UseStateForUnknown(),
				},
			},
			"schedule_days_of_week": schema.ListAttribute{
				Description: "If trigger is `custom_schedule`, the 0-based weekdays used by the schedule.",
				Computed:    true,
				Optional:    true,
				ElementType: types.Int64Type,
				PlanModifiers: []planmodifier.List{
					listplanmodifier.UseStateForUnknown(),
				},
			},
			"schedule_times_of_day": schema.ListAttribute{
				Description: "Times of day in HH:MM format for schedule-driven exports.",
				Computed:    true,
				Optional:    true,
				ElementType: types.StringType,
				PlanModifiers: []planmodifier.List{
					listplanmodifier.UseStateForUnknown(),
				},
			},
			"schedule_time_zone": schema.StringAttribute{
				Description: "Time zone used by the scheduled export.",
				Computed:    true,
				Optional:    true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
			},
			"holiday_region": schema.StringAttribute{
				Description: "Optional holiday region used by schedule-driven exports.",
				Computed:    true,
				Optional:    true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
			},
			"id": schema.Int64Attribute{
				Description: "Scheduled Export ID",
				Computed:    true,
				PlanModifiers: []planmodifier.Int64{
					int64planmodifier.UseStateForUnknown(),
				},
			},
			"report_name": schema.StringAttribute{
				Description: "Human-readable report name.",
				Computed:    true,
			},
			"human_readable_schedule": schema.StringAttribute{
				Description: "Human-readable schedule description.",
				Computed:    true,
			},
			"last_run_at": schema.StringAttribute{
				Description: "Most recent scheduled run time.",
				Computed:    true,
			},
			"last_export_id": schema.Int64Attribute{
				Description: "Most recent Export ID created by this schedule.",
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

func (r *scheduledExportResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	var plan scheduledExportResourceModel
	diags := req.Plan.Get(ctx, &plan)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}
	var config scheduledExportResourceModel
	diags = req.Config.Get(ctx, &config)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	paramsScheduledExportCreate := files_sdk.ScheduledExportCreateParams{}
	paramsScheduledExportCreate.Name = plan.Name.ValueString()
	paramsScheduledExportCreate.ExportType = plan.ExportType.ValueString()
	createExportOptions, diags := lib.DynamicToInterface(ctx, path.Root("export_options"), plan.ExportOptions)
	resp.Diagnostics.Append(diags...)
	paramsScheduledExportCreate.ExportOptions = createExportOptions
	paramsScheduledExportCreate.UserId = plan.UserId.ValueInt64()
	if !plan.Disabled.IsNull() && !plan.Disabled.IsUnknown() {
		paramsScheduledExportCreate.Disabled = plan.Disabled.ValueBoolPointer()
	}
	paramsScheduledExportCreate.Trigger = paramsScheduledExportCreate.Trigger.Enum()[plan.Trigger.ValueString()]
	paramsScheduledExportCreate.Interval = plan.Interval.ValueString()
	paramsScheduledExportCreate.RecurringDay = plan.RecurringDay.ValueInt64()
	if !plan.ScheduleDaysOfWeek.IsNull() && !plan.ScheduleDaysOfWeek.IsUnknown() {
		diags = plan.ScheduleDaysOfWeek.ElementsAs(ctx, &paramsScheduledExportCreate.ScheduleDaysOfWeek, false)
		resp.Diagnostics.Append(diags...)
	}
	if !plan.ScheduleTimesOfDay.IsNull() && !plan.ScheduleTimesOfDay.IsUnknown() {
		diags = plan.ScheduleTimesOfDay.ElementsAs(ctx, &paramsScheduledExportCreate.ScheduleTimesOfDay, false)
		resp.Diagnostics.Append(diags...)
	}
	paramsScheduledExportCreate.ScheduleTimeZone = plan.ScheduleTimeZone.ValueString()
	paramsScheduledExportCreate.HolidayRegion = plan.HolidayRegion.ValueString()

	if resp.Diagnostics.HasError() {
		return
	}

	scheduledExport, err := r.client.Create(paramsScheduledExportCreate, files_sdk.WithContext(ctx))
	if err != nil {
		resp.Diagnostics.AddError(
			"Error Creating Files ScheduledExport",
			"Could not create scheduled_export, unexpected error: "+err.Error(),
		)
		return
	}

	diags = r.populateResourceModel(ctx, scheduledExport, &plan)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	diags = resp.State.Set(ctx, plan)
	resp.Diagnostics.Append(diags...)
}

func (r *scheduledExportResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	var state scheduledExportResourceModel
	diags := req.State.Get(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	paramsScheduledExportFind := files_sdk.ScheduledExportFindParams{}
	paramsScheduledExportFind.Id = state.Id.ValueInt64()

	scheduledExport, err := r.client.Find(paramsScheduledExportFind, files_sdk.WithContext(ctx))
	if err != nil {
		if files_sdk.IsNotExist(err) {
			resp.State.RemoveResource(ctx)
			return
		}

		resp.Diagnostics.AddError(
			"Error Reading Files ScheduledExport",
			"Could not read scheduled_export id "+fmt.Sprint(state.Id.ValueInt64())+": "+err.Error(),
		)
		return
	}

	diags = r.populateResourceModel(ctx, scheduledExport, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	diags = resp.State.Set(ctx, state)
	resp.Diagnostics.Append(diags...)
}

func (r *scheduledExportResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	var plan scheduledExportResourceModel
	diags := req.Plan.Get(ctx, &plan)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}
	var config scheduledExportResourceModel
	diags = req.Config.Get(ctx, &config)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	paramsScheduledExportUpdate := map[string]interface{}{}
	if !plan.Id.IsNull() && !plan.Id.IsUnknown() {
		paramsScheduledExportUpdate["id"] = plan.Id.ValueInt64()
	}
	if !config.Name.IsNull() && !config.Name.IsUnknown() {
		paramsScheduledExportUpdate["name"] = config.Name.ValueString()
	}
	if !config.ExportType.IsNull() && !config.ExportType.IsUnknown() {
		paramsScheduledExportUpdate["export_type"] = config.ExportType.ValueString()
	}
	updateExportOptions, diags := lib.DynamicToInterface(ctx, path.Root("export_options"), config.ExportOptions)
	resp.Diagnostics.Append(diags...)
	paramsScheduledExportUpdate["export_options"] = updateExportOptions
	if !config.UserId.IsNull() && !config.UserId.IsUnknown() {
		paramsScheduledExportUpdate["user_id"] = config.UserId.ValueInt64()
	}
	if !config.Disabled.IsNull() && !config.Disabled.IsUnknown() {
		paramsScheduledExportUpdate["disabled"] = config.Disabled.ValueBool()
	}
	if !config.Trigger.IsNull() && !config.Trigger.IsUnknown() {
		paramsScheduledExportUpdate["trigger"] = config.Trigger.ValueString()
	}
	if !config.Interval.IsNull() && !config.Interval.IsUnknown() {
		paramsScheduledExportUpdate["interval"] = config.Interval.ValueString()
	}
	if !config.RecurringDay.IsNull() && !config.RecurringDay.IsUnknown() {
		paramsScheduledExportUpdate["recurring_day"] = config.RecurringDay.ValueInt64()
	}
	if !config.ScheduleDaysOfWeek.IsNull() && !config.ScheduleDaysOfWeek.IsUnknown() {
		var updateScheduleDaysOfWeek []int64
		diags = config.ScheduleDaysOfWeek.ElementsAs(ctx, &updateScheduleDaysOfWeek, false)
		resp.Diagnostics.Append(diags...)
		paramsScheduledExportUpdate["schedule_days_of_week"] = updateScheduleDaysOfWeek
	}
	if !config.ScheduleTimesOfDay.IsNull() && !config.ScheduleTimesOfDay.IsUnknown() {
		var updateScheduleTimesOfDay []string
		diags = config.ScheduleTimesOfDay.ElementsAs(ctx, &updateScheduleTimesOfDay, false)
		resp.Diagnostics.Append(diags...)
		paramsScheduledExportUpdate["schedule_times_of_day"] = updateScheduleTimesOfDay
	}
	if !config.ScheduleTimeZone.IsNull() && !config.ScheduleTimeZone.IsUnknown() {
		paramsScheduledExportUpdate["schedule_time_zone"] = config.ScheduleTimeZone.ValueString()
	}
	if !config.HolidayRegion.IsNull() && !config.HolidayRegion.IsUnknown() {
		paramsScheduledExportUpdate["holiday_region"] = config.HolidayRegion.ValueString()
	}

	if resp.Diagnostics.HasError() {
		return
	}

	scheduledExport, err := r.client.UpdateWithMap(paramsScheduledExportUpdate, files_sdk.WithContext(ctx))
	if err != nil {
		resp.Diagnostics.AddError(
			"Error Updating Files ScheduledExport",
			"Could not update scheduled_export, unexpected error: "+err.Error(),
		)
		return
	}

	diags = r.populateResourceModel(ctx, scheduledExport, &plan)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	diags = resp.State.Set(ctx, plan)
	resp.Diagnostics.Append(diags...)
}

func (r *scheduledExportResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	var state scheduledExportResourceModel
	diags := req.State.Get(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	paramsScheduledExportDelete := files_sdk.ScheduledExportDeleteParams{}
	paramsScheduledExportDelete.Id = state.Id.ValueInt64()

	err := r.client.Delete(paramsScheduledExportDelete, files_sdk.WithContext(ctx))
	if err != nil && !files_sdk.IsNotExist(err) {
		resp.Diagnostics.AddError(
			"Error Deleting Files ScheduledExport",
			"Could not delete scheduled_export id "+fmt.Sprint(state.Id.ValueInt64())+": "+err.Error(),
		)
	}
}

func (r *scheduledExportResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
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

func (r *scheduledExportResource) populateResourceModel(ctx context.Context, scheduledExport files_sdk.ScheduledExport, state *scheduledExportResourceModel) (diags diag.Diagnostics) {
	var propDiags diag.Diagnostics

	state.Id = types.Int64Value(scheduledExport.Id)
	state.Name = types.StringValue(scheduledExport.Name)
	state.ExportType = types.StringValue(scheduledExport.ExportType)
	state.ReportName = types.StringValue(scheduledExport.ReportName)
	state.ExportOptions, propDiags = lib.ToDynamic(ctx, path.Root("export_options"), scheduledExport.ExportOptions, state.ExportOptions.UnderlyingValue())
	diags.Append(propDiags...)
	state.UserId = types.Int64Value(scheduledExport.UserId)
	state.Disabled = types.BoolPointerValue(scheduledExport.Disabled)
	state.Trigger = types.StringValue(scheduledExport.Trigger)
	state.Interval = types.StringValue(scheduledExport.Interval)
	state.RecurringDay = types.Int64Value(scheduledExport.RecurringDay)
	state.ScheduleDaysOfWeek, propDiags = types.ListValueFrom(ctx, types.Int64Type, scheduledExport.ScheduleDaysOfWeek)
	diags.Append(propDiags...)
	state.ScheduleTimesOfDay, propDiags = types.ListValueFrom(ctx, types.StringType, scheduledExport.ScheduleTimesOfDay)
	diags.Append(propDiags...)
	state.ScheduleTimeZone = types.StringValue(scheduledExport.ScheduleTimeZone)
	state.HolidayRegion = types.StringValue(scheduledExport.HolidayRegion)
	state.HumanReadableSchedule = types.StringValue(scheduledExport.HumanReadableSchedule)
	if err := lib.TimeToStringType(ctx, path.Root("last_run_at"), scheduledExport.LastRunAt, &state.LastRunAt); err != nil {
		diags.AddError(
			"Error Creating Files ScheduledExport",
			"Could not convert state last_run_at to string: "+err.Error(),
		)
	}
	state.LastExportId = types.Int64Value(scheduledExport.LastExportId)
	if err := lib.TimeToStringType(ctx, path.Root("created_at"), scheduledExport.CreatedAt, &state.CreatedAt); err != nil {
		diags.AddError(
			"Error Creating Files ScheduledExport",
			"Could not convert state created_at to string: "+err.Error(),
		)
	}
	if err := lib.TimeToStringType(ctx, path.Root("updated_at"), scheduledExport.UpdatedAt, &state.UpdatedAt); err != nil {
		diags.AddError(
			"Error Creating Files ScheduledExport",
			"Could not convert state updated_at to string: "+err.Error(),
		)
	}

	return
}
