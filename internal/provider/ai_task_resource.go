package provider

import (
	"context"
	"fmt"
	"strconv"
	"strings"

	files_sdk "github.com/Files-com/files-sdk-go/v3"
	ai_task "github.com/Files-com/files-sdk-go/v3/aitask"
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
	_ resource.Resource                = &aiTaskResource{}
	_ resource.ResourceWithConfigure   = &aiTaskResource{}
	_ resource.ResourceWithImportState = &aiTaskResource{}
)

func NewAiTaskResource() resource.Resource {
	return &aiTaskResource{}
}

type aiTaskResource struct {
	client *ai_task.Client
}

type aiTaskResourceModel struct {
	Name                  types.String `tfsdk:"name"`
	Prompt                types.String `tfsdk:"prompt"`
	WorkspaceId           types.Int64  `tfsdk:"workspace_id"`
	Description           types.String `tfsdk:"description"`
	PermissionSet         types.String `tfsdk:"permission_set"`
	Path                  types.String `tfsdk:"path"`
	Source                types.String `tfsdk:"source"`
	Disabled              types.Bool   `tfsdk:"disabled"`
	Trigger               types.String `tfsdk:"trigger"`
	TriggerActions        types.List   `tfsdk:"trigger_actions"`
	Interval              types.String `tfsdk:"interval"`
	RecurringDay          types.Int64  `tfsdk:"recurring_day"`
	ScheduleDaysOfWeek    types.List   `tfsdk:"schedule_days_of_week"`
	ScheduleTimesOfDay    types.List   `tfsdk:"schedule_times_of_day"`
	ScheduleTimeZone      types.String `tfsdk:"schedule_time_zone"`
	HolidayRegion         types.String `tfsdk:"holiday_region"`
	Id                    types.Int64  `tfsdk:"id"`
	HumanReadableSchedule types.String `tfsdk:"human_readable_schedule"`
	LastRunAt             types.String `tfsdk:"last_run_at"`
	MasterAdminUserId     types.Int64  `tfsdk:"master_admin_user_id"`
	CreatedAt             types.String `tfsdk:"created_at"`
	UpdatedAt             types.String `tfsdk:"updated_at"`
}

func (r *aiTaskResource) Configure(_ context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {
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

	r.client = &ai_task.Client{Config: sdk_config}
}

func (r *aiTaskResource) Metadata(_ context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_ai_task"
}

func (r *aiTaskResource) Schema(_ context.Context, _ resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Description: "An AI Task defines a Files.com AI prompt that can run on a schedule or in response to file actions.",
		Attributes: map[string]schema.Attribute{
			"name": schema.StringAttribute{
				Description: "AI Task name.",
				Required:    true,
			},
			"prompt": schema.StringAttribute{
				Description: "Prompt sent when this AI Task is invoked.",
				Required:    true,
			},
			"workspace_id": schema.Int64Attribute{
				Description: "Workspace ID. `0` means the default workspace.",
				Computed:    true,
				Optional:    true,
				PlanModifiers: []planmodifier.Int64{
					int64planmodifier.UseStateForUnknown(),
				},
			},
			"description": schema.StringAttribute{
				Description: "AI Task description.",
				Computed:    true,
				Optional:    true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
			},
			"permission_set": schema.StringAttribute{
				Description: "Permissions used by the internal API key for this AI Task. Valid values are `full` and `files_only`.",
				Computed:    true,
				Optional:    true,
				Validators: []validator.String{
					stringvalidator.OneOf("full", "files_only"),
				},
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
			},
			"path": schema.StringAttribute{
				Description: "Path scope used for action-triggered AI Tasks. This must be slash-delimited, but it must neither start nor end with a slash. Maximum of 5000 characters.",
				Computed:    true,
				Optional:    true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
			},
			"source": schema.StringAttribute{
				Description: "Source glob used with `path` for action-triggered AI Tasks.",
				Computed:    true,
				Optional:    true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
			},
			"disabled": schema.BoolAttribute{
				Description: "If true, this AI Task will not run.",
				Computed:    true,
				Optional:    true,
				PlanModifiers: []planmodifier.Bool{
					boolplanmodifier.UseStateForUnknown(),
				},
			},
			"trigger": schema.StringAttribute{
				Description: "How this AI Task is triggered.",
				Computed:    true,
				Optional:    true,
				Validators: []validator.String{
					stringvalidator.OneOf("manual", "daily", "custom_schedule", "action"),
				},
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
			},
			"trigger_actions": schema.ListAttribute{
				Description: "If trigger is `action`, the file action types that invoke this AI Task. Valid actions are create, copy, move, archived_delete, update, read, destroy.",
				Computed:    true,
				Optional:    true,
				ElementType: types.StringType,
				PlanModifiers: []planmodifier.List{
					listplanmodifier.UseStateForUnknown(),
				},
			},
			"interval": schema.StringAttribute{
				Description: "If trigger is `daily`, this specifies how often to run the AI Task.",
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
				Description: "Times of day in HH:MM format for scheduled AI Tasks.",
				Computed:    true,
				Optional:    true,
				ElementType: types.StringType,
				PlanModifiers: []planmodifier.List{
					listplanmodifier.UseStateForUnknown(),
				},
			},
			"schedule_time_zone": schema.StringAttribute{
				Description: "Time zone used by the AI Task schedule.",
				Computed:    true,
				Optional:    true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
			},
			"holiday_region": schema.StringAttribute{
				Description: "Optional holiday region used by scheduled AI Tasks.",
				Computed:    true,
				Optional:    true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
			},
			"id": schema.Int64Attribute{
				Description: "AI Task ID.",
				Computed:    true,
				PlanModifiers: []planmodifier.Int64{
					int64planmodifier.UseStateForUnknown(),
				},
			},
			"human_readable_schedule": schema.StringAttribute{
				Description: "Human-readable schedule description.",
				Computed:    true,
			},
			"last_run_at": schema.StringAttribute{
				Description: "Most recent successful invocation time.",
				Computed:    true,
			},
			"master_admin_user_id": schema.Int64Attribute{
				Description: "Master User ID used for AI Task invocations.",
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

func (r *aiTaskResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	var plan aiTaskResourceModel
	diags := req.Plan.Get(ctx, &plan)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}
	var config aiTaskResourceModel
	diags = req.Config.Get(ctx, &config)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	paramsAiTaskCreate := files_sdk.AiTaskCreateParams{}
	paramsAiTaskCreate.Description = plan.Description.ValueString()
	if !plan.Disabled.IsNull() && !plan.Disabled.IsUnknown() {
		paramsAiTaskCreate.Disabled = plan.Disabled.ValueBoolPointer()
	}
	paramsAiTaskCreate.HolidayRegion = plan.HolidayRegion.ValueString()
	paramsAiTaskCreate.Interval = plan.Interval.ValueString()
	paramsAiTaskCreate.Name = plan.Name.ValueString()
	paramsAiTaskCreate.Path = plan.Path.ValueString()
	paramsAiTaskCreate.PermissionSet = paramsAiTaskCreate.PermissionSet.Enum()[plan.PermissionSet.ValueString()]
	paramsAiTaskCreate.Prompt = plan.Prompt.ValueString()
	paramsAiTaskCreate.RecurringDay = plan.RecurringDay.ValueInt64()
	if !plan.ScheduleDaysOfWeek.IsNull() && !plan.ScheduleDaysOfWeek.IsUnknown() {
		diags = plan.ScheduleDaysOfWeek.ElementsAs(ctx, &paramsAiTaskCreate.ScheduleDaysOfWeek, false)
		resp.Diagnostics.Append(diags...)
	}
	paramsAiTaskCreate.ScheduleTimeZone = plan.ScheduleTimeZone.ValueString()
	if !plan.ScheduleTimesOfDay.IsNull() && !plan.ScheduleTimesOfDay.IsUnknown() {
		diags = plan.ScheduleTimesOfDay.ElementsAs(ctx, &paramsAiTaskCreate.ScheduleTimesOfDay, false)
		resp.Diagnostics.Append(diags...)
	}
	paramsAiTaskCreate.Source = plan.Source.ValueString()
	paramsAiTaskCreate.Trigger = paramsAiTaskCreate.Trigger.Enum()[plan.Trigger.ValueString()]
	if !plan.TriggerActions.IsNull() && !plan.TriggerActions.IsUnknown() {
		diags = plan.TriggerActions.ElementsAs(ctx, &paramsAiTaskCreate.TriggerActions, false)
		resp.Diagnostics.Append(diags...)
	}
	paramsAiTaskCreate.WorkspaceId = plan.WorkspaceId.ValueInt64()

	if resp.Diagnostics.HasError() {
		return
	}

	aiTask, err := r.client.Create(paramsAiTaskCreate, files_sdk.WithContext(ctx))
	if err != nil {
		resp.Diagnostics.AddError(
			"Error Creating Files AiTask",
			"Could not create ai_task, unexpected error: "+err.Error(),
		)
		return
	}

	diags = r.populateResourceModel(ctx, aiTask, &plan)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	diags = resp.State.Set(ctx, plan)
	resp.Diagnostics.Append(diags...)
}

func (r *aiTaskResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	var state aiTaskResourceModel
	diags := req.State.Get(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	paramsAiTaskFind := files_sdk.AiTaskFindParams{}
	paramsAiTaskFind.Id = state.Id.ValueInt64()

	aiTask, err := r.client.Find(paramsAiTaskFind, files_sdk.WithContext(ctx))
	if err != nil {
		if files_sdk.IsNotExist(err) {
			resp.State.RemoveResource(ctx)
			return
		}

		resp.Diagnostics.AddError(
			"Error Reading Files AiTask",
			"Could not read ai_task id "+fmt.Sprint(state.Id.ValueInt64())+": "+err.Error(),
		)
		return
	}

	diags = r.populateResourceModel(ctx, aiTask, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	diags = resp.State.Set(ctx, state)
	resp.Diagnostics.Append(diags...)
}

func (r *aiTaskResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	var plan aiTaskResourceModel
	diags := req.Plan.Get(ctx, &plan)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}
	var config aiTaskResourceModel
	diags = req.Config.Get(ctx, &config)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	paramsAiTaskUpdate := map[string]interface{}{}
	if !plan.Id.IsNull() && !plan.Id.IsUnknown() {
		paramsAiTaskUpdate["id"] = plan.Id.ValueInt64()
	}
	if !config.Description.IsNull() && !config.Description.IsUnknown() {
		paramsAiTaskUpdate["description"] = config.Description.ValueString()
	}
	if !config.Disabled.IsNull() && !config.Disabled.IsUnknown() {
		paramsAiTaskUpdate["disabled"] = config.Disabled.ValueBool()
	}
	if !config.HolidayRegion.IsNull() && !config.HolidayRegion.IsUnknown() {
		paramsAiTaskUpdate["holiday_region"] = config.HolidayRegion.ValueString()
	}
	if !config.Interval.IsNull() && !config.Interval.IsUnknown() {
		paramsAiTaskUpdate["interval"] = config.Interval.ValueString()
	}
	if !config.Name.IsNull() && !config.Name.IsUnknown() {
		paramsAiTaskUpdate["name"] = config.Name.ValueString()
	}
	if !config.Path.IsNull() && !config.Path.IsUnknown() {
		paramsAiTaskUpdate["path"] = config.Path.ValueString()
	}
	if !config.PermissionSet.IsNull() && !config.PermissionSet.IsUnknown() {
		paramsAiTaskUpdate["permission_set"] = config.PermissionSet.ValueString()
	}
	if !config.Prompt.IsNull() && !config.Prompt.IsUnknown() {
		paramsAiTaskUpdate["prompt"] = config.Prompt.ValueString()
	}
	if !config.RecurringDay.IsNull() && !config.RecurringDay.IsUnknown() {
		paramsAiTaskUpdate["recurring_day"] = config.RecurringDay.ValueInt64()
	}
	if !config.ScheduleDaysOfWeek.IsNull() && !config.ScheduleDaysOfWeek.IsUnknown() {
		var updateScheduleDaysOfWeek []int64
		diags = config.ScheduleDaysOfWeek.ElementsAs(ctx, &updateScheduleDaysOfWeek, false)
		resp.Diagnostics.Append(diags...)
		paramsAiTaskUpdate["schedule_days_of_week"] = updateScheduleDaysOfWeek
	}
	if !config.ScheduleTimeZone.IsNull() && !config.ScheduleTimeZone.IsUnknown() {
		paramsAiTaskUpdate["schedule_time_zone"] = config.ScheduleTimeZone.ValueString()
	}
	if !config.ScheduleTimesOfDay.IsNull() && !config.ScheduleTimesOfDay.IsUnknown() {
		var updateScheduleTimesOfDay []string
		diags = config.ScheduleTimesOfDay.ElementsAs(ctx, &updateScheduleTimesOfDay, false)
		resp.Diagnostics.Append(diags...)
		paramsAiTaskUpdate["schedule_times_of_day"] = updateScheduleTimesOfDay
	}
	if !config.Source.IsNull() && !config.Source.IsUnknown() {
		paramsAiTaskUpdate["source"] = config.Source.ValueString()
	}
	if !config.Trigger.IsNull() && !config.Trigger.IsUnknown() {
		paramsAiTaskUpdate["trigger"] = config.Trigger.ValueString()
	}
	if !config.TriggerActions.IsNull() && !config.TriggerActions.IsUnknown() {
		var updateTriggerActions []string
		diags = config.TriggerActions.ElementsAs(ctx, &updateTriggerActions, false)
		resp.Diagnostics.Append(diags...)
		paramsAiTaskUpdate["trigger_actions"] = updateTriggerActions
	}
	if !config.WorkspaceId.IsNull() && !config.WorkspaceId.IsUnknown() {
		paramsAiTaskUpdate["workspace_id"] = config.WorkspaceId.ValueInt64()
	}

	if resp.Diagnostics.HasError() {
		return
	}

	aiTask, err := r.client.UpdateWithMap(paramsAiTaskUpdate, files_sdk.WithContext(ctx))
	if err != nil {
		resp.Diagnostics.AddError(
			"Error Updating Files AiTask",
			"Could not update ai_task, unexpected error: "+err.Error(),
		)
		return
	}

	diags = r.populateResourceModel(ctx, aiTask, &plan)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	diags = resp.State.Set(ctx, plan)
	resp.Diagnostics.Append(diags...)
}

func (r *aiTaskResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	var state aiTaskResourceModel
	diags := req.State.Get(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	paramsAiTaskDelete := files_sdk.AiTaskDeleteParams{}
	paramsAiTaskDelete.Id = state.Id.ValueInt64()

	err := r.client.Delete(paramsAiTaskDelete, files_sdk.WithContext(ctx))
	if err != nil && !files_sdk.IsNotExist(err) {
		resp.Diagnostics.AddError(
			"Error Deleting Files AiTask",
			"Could not delete ai_task id "+fmt.Sprint(state.Id.ValueInt64())+": "+err.Error(),
		)
	}
}

func (r *aiTaskResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
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

func (r *aiTaskResource) populateResourceModel(ctx context.Context, aiTask files_sdk.AiTask, state *aiTaskResourceModel) (diags diag.Diagnostics) {
	var propDiags diag.Diagnostics

	state.Id = types.Int64Value(aiTask.Id)
	state.WorkspaceId = types.Int64Value(aiTask.WorkspaceId)
	state.Name = types.StringValue(aiTask.Name)
	state.Description = types.StringValue(aiTask.Description)
	state.Prompt = types.StringValue(aiTask.Prompt)
	state.PermissionSet = types.StringValue(aiTask.PermissionSet)
	state.Path = types.StringValue(aiTask.Path)
	state.Source = types.StringValue(aiTask.Source)
	state.Disabled = types.BoolPointerValue(aiTask.Disabled)
	state.Trigger = types.StringValue(aiTask.Trigger)
	state.TriggerActions, propDiags = types.ListValueFrom(ctx, types.StringType, aiTask.TriggerActions)
	diags.Append(propDiags...)
	state.Interval = types.StringValue(aiTask.Interval)
	state.RecurringDay = types.Int64Value(aiTask.RecurringDay)
	state.ScheduleDaysOfWeek, propDiags = types.ListValueFrom(ctx, types.Int64Type, aiTask.ScheduleDaysOfWeek)
	diags.Append(propDiags...)
	state.ScheduleTimesOfDay, propDiags = types.ListValueFrom(ctx, types.StringType, aiTask.ScheduleTimesOfDay)
	diags.Append(propDiags...)
	state.ScheduleTimeZone = types.StringValue(aiTask.ScheduleTimeZone)
	state.HolidayRegion = types.StringValue(aiTask.HolidayRegion)
	state.HumanReadableSchedule = types.StringValue(aiTask.HumanReadableSchedule)
	if err := lib.TimeToStringType(ctx, path.Root("last_run_at"), aiTask.LastRunAt, &state.LastRunAt); err != nil {
		diags.AddError(
			"Error Creating Files AiTask",
			"Could not convert state last_run_at to string: "+err.Error(),
		)
	}
	state.MasterAdminUserId = types.Int64Value(aiTask.MasterAdminUserId)
	if err := lib.TimeToStringType(ctx, path.Root("created_at"), aiTask.CreatedAt, &state.CreatedAt); err != nil {
		diags.AddError(
			"Error Creating Files AiTask",
			"Could not convert state created_at to string: "+err.Error(),
		)
	}
	if err := lib.TimeToStringType(ctx, path.Root("updated_at"), aiTask.UpdatedAt, &state.UpdatedAt); err != nil {
		diags.AddError(
			"Error Creating Files AiTask",
			"Could not convert state updated_at to string: "+err.Error(),
		)
	}

	return
}
