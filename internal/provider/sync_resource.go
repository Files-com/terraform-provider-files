package provider

import (
	"context"
	"fmt"
	"strconv"
	"strings"

	files_sdk "github.com/Files-com/files-sdk-go/v3"
	sync "github.com/Files-com/files-sdk-go/v3/sync"
	"github.com/Files-com/terraform-provider-files/lib"
	"github.com/hashicorp/terraform-plugin-framework-validators/stringvalidator"
	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/boolplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/int64planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/listplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/schema/validator"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

var (
	_ resource.Resource                = &syncResource{}
	_ resource.ResourceWithConfigure   = &syncResource{}
	_ resource.ResourceWithImportState = &syncResource{}
)

func NewSyncResource() resource.Resource {
	return &syncResource{}
}

type syncResource struct {
	client *sync.Client
}

type syncResourceModel struct {
	Name                types.String `tfsdk:"name"`
	Description         types.String `tfsdk:"description"`
	SrcPath             types.String `tfsdk:"src_path"`
	DestPath            types.String `tfsdk:"dest_path"`
	SrcRemoteServerId   types.Int64  `tfsdk:"src_remote_server_id"`
	DestRemoteServerId  types.Int64  `tfsdk:"dest_remote_server_id"`
	TwoWay              types.Bool   `tfsdk:"two_way"`
	KeepAfterCopy       types.Bool   `tfsdk:"keep_after_copy"`
	DeleteEmptyFolders  types.Bool   `tfsdk:"delete_empty_folders"`
	Disabled            types.Bool   `tfsdk:"disabled"`
	Trigger             types.String `tfsdk:"trigger"`
	TriggerFile         types.String `tfsdk:"trigger_file"`
	Interval            types.String `tfsdk:"interval"`
	RecurringDay        types.Int64  `tfsdk:"recurring_day"`
	ScheduleDaysOfWeek  types.List   `tfsdk:"schedule_days_of_week"`
	ScheduleTimesOfDay  types.List   `tfsdk:"schedule_times_of_day"`
	ScheduleTimeZone    types.String `tfsdk:"schedule_time_zone"`
	Id                  types.Int64  `tfsdk:"id"`
	SiteId              types.Int64  `tfsdk:"site_id"`
	UserId              types.Int64  `tfsdk:"user_id"`
	IncludePatterns     types.List   `tfsdk:"include_patterns"`
	ExcludePatterns     types.List   `tfsdk:"exclude_patterns"`
	CreatedAt           types.String `tfsdk:"created_at"`
	UpdatedAt           types.String `tfsdk:"updated_at"`
	SyncIntervalMinutes types.Int64  `tfsdk:"sync_interval_minutes"`
	HolidayRegion       types.String `tfsdk:"holiday_region"`
}

func (r *syncResource) Configure(_ context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {
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

	r.client = &sync.Client{Config: sdk_config}
}

func (r *syncResource) Metadata(_ context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_sync"
}

func (r *syncResource) Schema(_ context.Context, _ resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Description: "A Sync represents a file synchronization job between two locations (local-remote, remote-remote, local-child_site, etc). \n\nIt can be scheduled, run manually, or triggered by custom logic. \n\nSyncs track their runs, status, and configuration.",
		Attributes: map[string]schema.Attribute{
			"name": schema.StringAttribute{
				Description: "Name for this sync job",
				Computed:    true,
				Optional:    true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
			},
			"description": schema.StringAttribute{
				Description: "Description for this sync job",
				Computed:    true,
				Optional:    true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
			},
			"src_path": schema.StringAttribute{
				Description: "Absolute source path for the sync",
				Computed:    true,
				Optional:    true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
			},
			"dest_path": schema.StringAttribute{
				Description: "Absolute destination path for the sync",
				Computed:    true,
				Optional:    true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
			},
			"src_remote_server_id": schema.Int64Attribute{
				Description: "Remote server ID for the source (if remote)",
				Computed:    true,
				Optional:    true,
				PlanModifiers: []planmodifier.Int64{
					int64planmodifier.UseStateForUnknown(),
				},
			},
			"dest_remote_server_id": schema.Int64Attribute{
				Description: "Remote server ID for the destination (if remote)",
				Computed:    true,
				Optional:    true,
				PlanModifiers: []planmodifier.Int64{
					int64planmodifier.UseStateForUnknown(),
				},
			},
			"two_way": schema.BoolAttribute{
				Description: "Is this a two-way sync?",
				Computed:    true,
				Optional:    true,
				PlanModifiers: []planmodifier.Bool{
					boolplanmodifier.UseStateForUnknown(),
				},
			},
			"keep_after_copy": schema.BoolAttribute{
				Description: "Keep files after copying?",
				Computed:    true,
				Optional:    true,
				PlanModifiers: []planmodifier.Bool{
					boolplanmodifier.UseStateForUnknown(),
				},
			},
			"delete_empty_folders": schema.BoolAttribute{
				Description: "Delete empty folders after sync?",
				Computed:    true,
				Optional:    true,
				PlanModifiers: []planmodifier.Bool{
					boolplanmodifier.UseStateForUnknown(),
				},
			},
			"disabled": schema.BoolAttribute{
				Description: "Is this sync disabled?",
				Computed:    true,
				Optional:    true,
				PlanModifiers: []planmodifier.Bool{
					boolplanmodifier.UseStateForUnknown(),
				},
			},
			"trigger": schema.StringAttribute{
				Description: "Trigger type: daily, custom_schedule, or manual",
				Computed:    true,
				Optional:    true,
				Validators: []validator.String{
					stringvalidator.OneOf("daily", "custom_schedule", "manual"),
				},
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
			},
			"trigger_file": schema.StringAttribute{
				Description: "Some MFT services request an empty file (known as a trigger file) to signal the sync is complete and they can begin further processing. If trigger_file is set, a zero-byte file will be sent at the end of the sync.",
				Computed:    true,
				Optional:    true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
			},
			"interval": schema.StringAttribute{
				Description: "If trigger is `daily`, this specifies how often to run this sync.  One of: `day`, `week`, `week_end`, `month`, `month_end`, `quarter`, `quarter_end`, `year`, `year_end`",
				Computed:    true,
				Optional:    true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
			},
			"recurring_day": schema.Int64Attribute{
				Description: "If trigger type is `daily`, this specifies a day number to run in one of the supported intervals: `week`, `month`, `quarter`, `year`.",
				Computed:    true,
				Optional:    true,
				PlanModifiers: []planmodifier.Int64{
					int64planmodifier.UseStateForUnknown(),
				},
			},
			"schedule_days_of_week": schema.ListAttribute{
				Description: "If trigger is `custom_schedule`, Custom schedule description for when the sync should be run. 0-based days of the week. 0 is Sunday, 1 is Monday, etc.",
				Computed:    true,
				Optional:    true,
				ElementType: types.Int64Type,
				PlanModifiers: []planmodifier.List{
					listplanmodifier.UseStateForUnknown(),
				},
			},
			"schedule_times_of_day": schema.ListAttribute{
				Description: "If trigger is `custom_schedule`, Custom schedule description for when the sync should be run. Times of day in HH:MM format.",
				Computed:    true,
				Optional:    true,
				ElementType: types.StringType,
				PlanModifiers: []planmodifier.List{
					listplanmodifier.UseStateForUnknown(),
				},
			},
			"schedule_time_zone": schema.StringAttribute{
				Description: "If trigger is `custom_schedule`, Custom schedule Time Zone for when the sync should be run.",
				Computed:    true,
				Optional:    true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
			},
			"id": schema.Int64Attribute{
				Description: "Sync ID",
				Computed:    true,
				PlanModifiers: []planmodifier.Int64{
					int64planmodifier.UseStateForUnknown(),
				},
			},
			"site_id": schema.Int64Attribute{
				Description: "Site ID this sync belongs to",
				Computed:    true,
			},
			"user_id": schema.Int64Attribute{
				Description: "User who created or owns this sync",
				Computed:    true,
			},
			"include_patterns": schema.ListAttribute{
				Description: "Array of glob patterns to include",
				Computed:    true,
				ElementType: types.StringType,
			},
			"exclude_patterns": schema.ListAttribute{
				Description: "Array of glob patterns to exclude",
				Computed:    true,
				ElementType: types.StringType,
			},
			"created_at": schema.StringAttribute{
				Description: "When this sync was created",
				Computed:    true,
			},
			"updated_at": schema.StringAttribute{
				Description: "When this sync was last updated",
				Computed:    true,
			},
			"sync_interval_minutes": schema.Int64Attribute{
				Description: "Frequency in minutes between syncs. If set, this value must be greater than or equal to the `remote_sync_interval` value for the site's plan. If left blank, the plan's `remote_sync_interval` will be used. This setting is only used if `trigger` is empty.",
				Computed:    true,
			},
			"holiday_region": schema.StringAttribute{
				Description: "If trigger is `custom_schedule`, the Automation will check if there is a formal, observed holiday for the region, and if so, it will not run.",
				Computed:    true,
			},
		},
	}
}

func (r *syncResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	var plan syncResourceModel
	diags := req.Plan.Get(ctx, &plan)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	paramsSyncCreate := files_sdk.SyncCreateParams{}
	paramsSyncCreate.Name = plan.Name.ValueString()
	paramsSyncCreate.Description = plan.Description.ValueString()
	paramsSyncCreate.SrcPath = plan.SrcPath.ValueString()
	paramsSyncCreate.DestPath = plan.DestPath.ValueString()
	paramsSyncCreate.SrcRemoteServerId = plan.SrcRemoteServerId.ValueInt64()
	paramsSyncCreate.DestRemoteServerId = plan.DestRemoteServerId.ValueInt64()
	if !plan.TwoWay.IsNull() && !plan.TwoWay.IsUnknown() {
		paramsSyncCreate.TwoWay = plan.TwoWay.ValueBoolPointer()
	}
	if !plan.KeepAfterCopy.IsNull() && !plan.KeepAfterCopy.IsUnknown() {
		paramsSyncCreate.KeepAfterCopy = plan.KeepAfterCopy.ValueBoolPointer()
	}
	if !plan.DeleteEmptyFolders.IsNull() && !plan.DeleteEmptyFolders.IsUnknown() {
		paramsSyncCreate.DeleteEmptyFolders = plan.DeleteEmptyFolders.ValueBoolPointer()
	}
	if !plan.Disabled.IsNull() && !plan.Disabled.IsUnknown() {
		paramsSyncCreate.Disabled = plan.Disabled.ValueBoolPointer()
	}
	paramsSyncCreate.Interval = plan.Interval.ValueString()
	paramsSyncCreate.Trigger = plan.Trigger.ValueString()
	paramsSyncCreate.TriggerFile = plan.TriggerFile.ValueString()
	paramsSyncCreate.RecurringDay = plan.RecurringDay.ValueInt64()
	paramsSyncCreate.ScheduleTimeZone = plan.ScheduleTimeZone.ValueString()
	if !plan.ScheduleDaysOfWeek.IsNull() && !plan.ScheduleDaysOfWeek.IsUnknown() {
		diags = plan.ScheduleDaysOfWeek.ElementsAs(ctx, &paramsSyncCreate.ScheduleDaysOfWeek, false)
		resp.Diagnostics.Append(diags...)
	}
	if !plan.ScheduleTimesOfDay.IsNull() && !plan.ScheduleTimesOfDay.IsUnknown() {
		diags = plan.ScheduleTimesOfDay.ElementsAs(ctx, &paramsSyncCreate.ScheduleTimesOfDay, false)
		resp.Diagnostics.Append(diags...)
	}

	if resp.Diagnostics.HasError() {
		return
	}

	sync, err := r.client.Create(paramsSyncCreate, files_sdk.WithContext(ctx))
	if err != nil {
		resp.Diagnostics.AddError(
			"Error Creating Files Sync",
			"Could not create sync, unexpected error: "+err.Error(),
		)
		return
	}

	diags = r.populateResourceModel(ctx, sync, &plan)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	diags = resp.State.Set(ctx, plan)
	resp.Diagnostics.Append(diags...)
}

func (r *syncResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	var state syncResourceModel
	diags := req.State.Get(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	paramsSyncFind := files_sdk.SyncFindParams{}
	paramsSyncFind.Id = state.Id.ValueInt64()

	sync, err := r.client.Find(paramsSyncFind, files_sdk.WithContext(ctx))
	if err != nil {
		if files_sdk.IsNotExist(err) {
			resp.State.RemoveResource(ctx)
			return
		}

		resp.Diagnostics.AddError(
			"Error Reading Files Sync",
			"Could not read sync id "+fmt.Sprint(state.Id.ValueInt64())+": "+err.Error(),
		)
		return
	}

	diags = r.populateResourceModel(ctx, sync, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	diags = resp.State.Set(ctx, state)
	resp.Diagnostics.Append(diags...)
}

func (r *syncResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	var plan syncResourceModel
	diags := req.Plan.Get(ctx, &plan)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	paramsSyncUpdate := files_sdk.SyncUpdateParams{}
	paramsSyncUpdate.Id = plan.Id.ValueInt64()
	paramsSyncUpdate.Name = plan.Name.ValueString()
	paramsSyncUpdate.Description = plan.Description.ValueString()
	paramsSyncUpdate.SrcPath = plan.SrcPath.ValueString()
	paramsSyncUpdate.DestPath = plan.DestPath.ValueString()
	paramsSyncUpdate.SrcRemoteServerId = plan.SrcRemoteServerId.ValueInt64()
	paramsSyncUpdate.DestRemoteServerId = plan.DestRemoteServerId.ValueInt64()
	if !plan.TwoWay.IsNull() && !plan.TwoWay.IsUnknown() {
		paramsSyncUpdate.TwoWay = plan.TwoWay.ValueBoolPointer()
	}
	if !plan.KeepAfterCopy.IsNull() && !plan.KeepAfterCopy.IsUnknown() {
		paramsSyncUpdate.KeepAfterCopy = plan.KeepAfterCopy.ValueBoolPointer()
	}
	if !plan.DeleteEmptyFolders.IsNull() && !plan.DeleteEmptyFolders.IsUnknown() {
		paramsSyncUpdate.DeleteEmptyFolders = plan.DeleteEmptyFolders.ValueBoolPointer()
	}
	if !plan.Disabled.IsNull() && !plan.Disabled.IsUnknown() {
		paramsSyncUpdate.Disabled = plan.Disabled.ValueBoolPointer()
	}
	paramsSyncUpdate.Interval = plan.Interval.ValueString()
	paramsSyncUpdate.Trigger = plan.Trigger.ValueString()
	paramsSyncUpdate.TriggerFile = plan.TriggerFile.ValueString()
	paramsSyncUpdate.RecurringDay = plan.RecurringDay.ValueInt64()
	paramsSyncUpdate.ScheduleTimeZone = plan.ScheduleTimeZone.ValueString()
	if !plan.ScheduleDaysOfWeek.IsNull() && !plan.ScheduleDaysOfWeek.IsUnknown() {
		diags = plan.ScheduleDaysOfWeek.ElementsAs(ctx, &paramsSyncUpdate.ScheduleDaysOfWeek, false)
		resp.Diagnostics.Append(diags...)
	}
	if !plan.ScheduleTimesOfDay.IsNull() && !plan.ScheduleTimesOfDay.IsUnknown() {
		diags = plan.ScheduleTimesOfDay.ElementsAs(ctx, &paramsSyncUpdate.ScheduleTimesOfDay, false)
		resp.Diagnostics.Append(diags...)
	}

	if resp.Diagnostics.HasError() {
		return
	}

	sync, err := r.client.Update(paramsSyncUpdate, files_sdk.WithContext(ctx))
	if err != nil {
		resp.Diagnostics.AddError(
			"Error Updating Files Sync",
			"Could not update sync, unexpected error: "+err.Error(),
		)
		return
	}

	diags = r.populateResourceModel(ctx, sync, &plan)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	diags = resp.State.Set(ctx, plan)
	resp.Diagnostics.Append(diags...)
}

func (r *syncResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	var state syncResourceModel
	diags := req.State.Get(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	paramsSyncDelete := files_sdk.SyncDeleteParams{}
	paramsSyncDelete.Id = state.Id.ValueInt64()

	err := r.client.Delete(paramsSyncDelete, files_sdk.WithContext(ctx))
	if err != nil && !files_sdk.IsNotExist(err) {
		resp.Diagnostics.AddError(
			"Error Deleting Files Sync",
			"Could not delete sync id "+fmt.Sprint(state.Id.ValueInt64())+": "+err.Error(),
		)
	}
}

func (r *syncResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
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

func (r *syncResource) populateResourceModel(ctx context.Context, sync files_sdk.Sync, state *syncResourceModel) (diags diag.Diagnostics) {
	var propDiags diag.Diagnostics

	state.Id = types.Int64Value(sync.Id)
	state.Name = types.StringValue(sync.Name)
	state.Description = types.StringValue(sync.Description)
	state.SiteId = types.Int64Value(sync.SiteId)
	state.UserId = types.Int64Value(sync.UserId)
	state.SrcPath = types.StringValue(sync.SrcPath)
	state.DestPath = types.StringValue(sync.DestPath)
	state.SrcRemoteServerId = types.Int64Value(sync.SrcRemoteServerId)
	state.DestRemoteServerId = types.Int64Value(sync.DestRemoteServerId)
	state.TwoWay = types.BoolPointerValue(sync.TwoWay)
	state.KeepAfterCopy = types.BoolPointerValue(sync.KeepAfterCopy)
	state.DeleteEmptyFolders = types.BoolPointerValue(sync.DeleteEmptyFolders)
	state.Disabled = types.BoolPointerValue(sync.Disabled)
	state.Trigger = types.StringValue(sync.Trigger)
	state.TriggerFile = types.StringValue(sync.TriggerFile)
	state.IncludePatterns, propDiags = types.ListValueFrom(ctx, types.StringType, sync.IncludePatterns)
	diags.Append(propDiags...)
	state.ExcludePatterns, propDiags = types.ListValueFrom(ctx, types.StringType, sync.ExcludePatterns)
	diags.Append(propDiags...)
	if err := lib.TimeToStringType(ctx, path.Root("created_at"), sync.CreatedAt, &state.CreatedAt); err != nil {
		diags.AddError(
			"Error Creating Files Sync",
			"Could not convert state created_at to string: "+err.Error(),
		)
	}
	if err := lib.TimeToStringType(ctx, path.Root("updated_at"), sync.UpdatedAt, &state.UpdatedAt); err != nil {
		diags.AddError(
			"Error Creating Files Sync",
			"Could not convert state updated_at to string: "+err.Error(),
		)
	}
	state.SyncIntervalMinutes = types.Int64Value(sync.SyncIntervalMinutes)
	state.Interval = types.StringValue(sync.Interval)
	state.RecurringDay = types.Int64Value(sync.RecurringDay)
	state.ScheduleDaysOfWeek, propDiags = types.ListValueFrom(ctx, types.Int64Type, sync.ScheduleDaysOfWeek)
	diags.Append(propDiags...)
	state.ScheduleTimesOfDay, propDiags = types.ListValueFrom(ctx, types.StringType, sync.ScheduleTimesOfDay)
	diags.Append(propDiags...)
	state.ScheduleTimeZone = types.StringValue(sync.ScheduleTimeZone)
	state.HolidayRegion = types.StringValue(sync.HolidayRegion)

	return
}
