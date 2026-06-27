package provider

import (
	"context"
	"fmt"

	files_sdk "github.com/Files-com/files-sdk-go/v3"
	ai_task "github.com/Files-com/files-sdk-go/v3/aitask"
	"github.com/Files-com/terraform-provider-files/lib"
	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/datasource/schema"
	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

var (
	_ datasource.DataSource              = &aiTaskDataSource{}
	_ datasource.DataSourceWithConfigure = &aiTaskDataSource{}
)

func NewAiTaskDataSource() datasource.DataSource {
	return &aiTaskDataSource{}
}

type aiTaskDataSource struct {
	client *ai_task.Client
}

type aiTaskDataSourceModel struct {
	Id                    types.Int64  `tfsdk:"id"`
	WorkspaceId           types.Int64  `tfsdk:"workspace_id"`
	Name                  types.String `tfsdk:"name"`
	Description           types.String `tfsdk:"description"`
	Prompt                types.String `tfsdk:"prompt"`
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
	HumanReadableSchedule types.String `tfsdk:"human_readable_schedule"`
	LastRunAt             types.String `tfsdk:"last_run_at"`
	MasterAdminUserId     types.Int64  `tfsdk:"master_admin_user_id"`
	CreatedAt             types.String `tfsdk:"created_at"`
	UpdatedAt             types.String `tfsdk:"updated_at"`
}

func (r *aiTaskDataSource) Configure(_ context.Context, req datasource.ConfigureRequest, resp *datasource.ConfigureResponse) {
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

func (r *aiTaskDataSource) Metadata(_ context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_ai_task"
}

func (r *aiTaskDataSource) Schema(_ context.Context, _ datasource.SchemaRequest, resp *datasource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Description: "An AI Task defines a Files.com AI prompt that can run on a schedule or in response to file actions.",
		Attributes: map[string]schema.Attribute{
			"id": schema.Int64Attribute{
				Description: "AI Task ID.",
				Required:    true,
			},
			"workspace_id": schema.Int64Attribute{
				Description: "Workspace ID. `0` means the default workspace.",
				Computed:    true,
			},
			"name": schema.StringAttribute{
				Description: "AI Task name.",
				Computed:    true,
			},
			"description": schema.StringAttribute{
				Description: "AI Task description.",
				Computed:    true,
			},
			"prompt": schema.StringAttribute{
				Description: "Prompt sent when this AI Task is invoked.",
				Computed:    true,
			},
			"path": schema.StringAttribute{
				Description: "Path scope used for action-triggered AI Tasks. This must be slash-delimited, but it must neither start nor end with a slash. Maximum of 5000 characters.",
				Computed:    true,
			},
			"source": schema.StringAttribute{
				Description: "Source glob used with `path` for action-triggered AI Tasks.",
				Computed:    true,
			},
			"disabled": schema.BoolAttribute{
				Description: "If true, this AI Task will not run.",
				Computed:    true,
			},
			"trigger": schema.StringAttribute{
				Description: "How this AI Task is triggered.",
				Computed:    true,
			},
			"trigger_actions": schema.ListAttribute{
				Description: "If trigger is `action`, the file action types that invoke this AI Task. Valid actions are create, copy, move, archived_delete, update, read, destroy.",
				Computed:    true,
				ElementType: types.StringType,
			},
			"interval": schema.StringAttribute{
				Description: "If trigger is `daily`, this specifies how often to run the AI Task.",
				Computed:    true,
			},
			"recurring_day": schema.Int64Attribute{
				Description: "If trigger is `daily`, this selects the day number inside the chosen interval.",
				Computed:    true,
			},
			"schedule_days_of_week": schema.ListAttribute{
				Description: "If trigger is `custom_schedule`, the 0-based weekdays used by the schedule.",
				Computed:    true,
				ElementType: types.Int64Type,
			},
			"schedule_times_of_day": schema.ListAttribute{
				Description: "Times of day in HH:MM format for scheduled AI Tasks.",
				Computed:    true,
				ElementType: types.StringType,
			},
			"schedule_time_zone": schema.StringAttribute{
				Description: "Time zone used by the AI Task schedule.",
				Computed:    true,
			},
			"holiday_region": schema.StringAttribute{
				Description: "Optional holiday region used by scheduled AI Tasks.",
				Computed:    true,
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

func (r *aiTaskDataSource) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
	var data aiTaskDataSourceModel
	diags := req.Config.Get(ctx, &data)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	paramsAiTaskFind := files_sdk.AiTaskFindParams{}
	paramsAiTaskFind.Id = data.Id.ValueInt64()

	aiTask, err := r.client.Find(paramsAiTaskFind, files_sdk.WithContext(ctx))
	if err != nil {
		resp.Diagnostics.AddError(
			"Error Reading Files AiTask",
			"Could not read ai_task id "+fmt.Sprint(data.Id.ValueInt64())+": "+err.Error(),
		)
		return
	}

	diags = r.populateDataSourceModel(ctx, aiTask, &data)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	diags = resp.State.Set(ctx, data)
	resp.Diagnostics.Append(diags...)
}

func (r *aiTaskDataSource) populateDataSourceModel(ctx context.Context, aiTask files_sdk.AiTask, state *aiTaskDataSourceModel) (diags diag.Diagnostics) {
	var propDiags diag.Diagnostics

	state.Id = types.Int64Value(aiTask.Id)
	state.WorkspaceId = types.Int64Value(aiTask.WorkspaceId)
	state.Name = types.StringValue(aiTask.Name)
	state.Description = types.StringValue(aiTask.Description)
	state.Prompt = types.StringValue(aiTask.Prompt)
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
