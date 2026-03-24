package provider

import (
	"context"
	"fmt"
	"strconv"
	"strings"

	files_sdk "github.com/Files-com/files-sdk-go/v3"
	expectation "github.com/Files-com/files-sdk-go/v3/expectation"
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
	_ resource.Resource                = &expectationResource{}
	_ resource.ResourceWithConfigure   = &expectationResource{}
	_ resource.ResourceWithImportState = &expectationResource{}
)

func NewExpectationResource() resource.Resource {
	return &expectationResource{}
}

type expectationResource struct {
	client *expectation.Client
}

type expectationResourceModel struct {
	WorkspaceId            types.Int64   `tfsdk:"workspace_id"`
	Name                   types.String  `tfsdk:"name"`
	Description            types.String  `tfsdk:"description"`
	Path                   types.String  `tfsdk:"path"`
	Source                 types.String  `tfsdk:"source"`
	ExcludePattern         types.String  `tfsdk:"exclude_pattern"`
	Disabled               types.Bool    `tfsdk:"disabled"`
	Trigger                types.String  `tfsdk:"trigger"`
	Interval               types.String  `tfsdk:"interval"`
	RecurringDay           types.Int64   `tfsdk:"recurring_day"`
	ScheduleDaysOfWeek     types.List    `tfsdk:"schedule_days_of_week"`
	ScheduleTimesOfDay     types.List    `tfsdk:"schedule_times_of_day"`
	ScheduleTimeZone       types.String  `tfsdk:"schedule_time_zone"`
	HolidayRegion          types.String  `tfsdk:"holiday_region"`
	LookbackInterval       types.Int64   `tfsdk:"lookback_interval"`
	LateAcceptanceInterval types.Int64   `tfsdk:"late_acceptance_interval"`
	InactivityInterval     types.Int64   `tfsdk:"inactivity_interval"`
	MaxOpenInterval        types.Int64   `tfsdk:"max_open_interval"`
	Criteria               types.Dynamic `tfsdk:"criteria"`
	Id                     types.Int64   `tfsdk:"id"`
	ExpectationsVersion    types.Int64   `tfsdk:"expectations_version"`
	LastEvaluatedAt        types.String  `tfsdk:"last_evaluated_at"`
	LastSuccessAt          types.String  `tfsdk:"last_success_at"`
	LastFailureAt          types.String  `tfsdk:"last_failure_at"`
	LastResult             types.String  `tfsdk:"last_result"`
	CreatedAt              types.String  `tfsdk:"created_at"`
	UpdatedAt              types.String  `tfsdk:"updated_at"`
}

func (r *expectationResource) Configure(_ context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {
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

	r.client = &expectation.Client{Config: sdk_config}
}

func (r *expectationResource) Metadata(_ context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_expectation"
}

func (r *expectationResource) Schema(_ context.Context, _ resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Description: "Expectations let your Files.com site define what “correct” file delivery looks like, continuously evaluate whether it happened, and keep history when it did not.\n\n\n\nExpectations are meant to answer operational questions like:\n\n\n\n* Did the expected file arrive?\n\n* Was it on time?\n\n* Did it meet the required shape and count rules?\n\n* Is there an active issue someone needs to acknowledge?\n\n\n\nExpectations are different from Automations and Syncs. Automations and Syncs act on files; Expectations monitor whether expected files arrived on time, in the right place, and in the right shape. In practice, Expectations are the sensor and Automations are the actuator.\n\n\n\nAn Expectation combines four concepts:\n\n\n\n1. **Scope**: where to look for candidate files, using `path`, `source`, and optional `exclude_pattern`.\n\n2. **Trigger / timing**: when a window opens and how long it stays eligible, using `trigger`, schedule fields, `lookback_interval`, `late_acceptance_interval`, `inactivity_interval`, and `max_open_interval`.\n\n3. **Criteria**: what must be true for the window to succeed, using the structured `criteria` JSON document.\n\n4. **Outcome history**: what happened over time, exposed through `ExpectationEvaluation` history and `ExpectationIncident` lifecycle records.\n\n\n\n## Scope and matching\n\n\n\nExpectations reuse the familiar Files.com path-plus-glob model.\n\n\n\nThe `path` field identifies the folder scope, while `source` identifies which files within that scope are candidates. `exclude_pattern` removes files from consideration.\n\n\n\nLike Automations, these fields support glob-style matching. Expectations treat those matches as one logical candidate set for each window. A single Expectation does not implicitly fan out into separate per-customer or per-folder evaluations just because the path contains wildcards.\n\n\n\n## Expectation windows\n\n\n\nExpectations are evaluated in windows.\n\n\n\nEach window is persisted as an `ExpectationEvaluation` record. A window opens, remains `open` while evidence can still arrive, and then closes into a terminal result such as `success`, `late`, `missing`, or `invalid`.\n\n\n\nAn Expectation has only one open window at a time.\n\n\n\n## Trigger modes\n\n\n\nExpectations can open windows in three ways:\n\n\n\n* `daily`: run on a recurring daily/weekly/monthly/quarterly/yearly cadence using `interval` and `recurring_day`.\n\n* `custom_schedule`: run on specific weekdays and times using `schedule_days_of_week`, `schedule_times_of_day`, and optional `schedule_time_zone` / `holiday_region`.\n\n* `manual`: an operator explicitly opens the window.\n\n\n\nSchedule-driven expectations define an on-time deadline and may optionally remain eligible to close as `late` during `late_acceptance_interval`.\n\n\n\nManual expectations have no concept of `late`; they open when triggered and close based on inactivity or hard-stop timing.\n\n\n\n## Success criteria\n\n\n\nThe `criteria` field is a structured JSON object describing what counts as success for the window.\n\n\n\nIn criteria v1, this can express things like:\n\n\n\n* file count constraints\n\n* total byte constraints\n\n* allowed extensions\n\n* filename regex validation\n\n* forbidden files\n\n* required named or globbed files with their own per-file constraints\n\n\n\nThis is intentionally structured rather than scriptable, so it stays safe, explainable, and versionable.\n\n\n\n## History and incidents\n\n\n\nThe Expectation itself stores summary state like `last_evaluated_at`, `last_success_at`, `last_failure_at`, and `last_result`.\n\n\n\nFor deeper inspection:\n\n\n\n* `ExpectationEvaluation` history shows each open or closed window and the evidence captured for it.\n\n* `ExpectationIncident` records track ongoing failure situations over time, including acknowledge, snooze, and resolve actions.\n\n\n\nManual windows do not open incidents in v1. Schedule-driven failures can open incidents, and later qualifying success can resolve them.",
		Attributes: map[string]schema.Attribute{
			"workspace_id": schema.Int64Attribute{
				Description: "Workspace ID. `0` means the default workspace.",
				Computed:    true,
				Optional:    true,
				PlanModifiers: []planmodifier.Int64{
					int64planmodifier.UseStateForUnknown(),
				},
			},
			"name": schema.StringAttribute{
				Description: "Expectation name.",
				Computed:    true,
				Optional:    true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
			},
			"description": schema.StringAttribute{
				Description: "Expectation description.",
				Computed:    true,
				Optional:    true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
			},
			"path": schema.StringAttribute{
				Description: "Path scope for the expectation. Supports workspace-relative presentation. This must be slash-delimited, but it must neither start nor end with a slash. Maximum of 5000 characters.",
				Computed:    true,
				Optional:    true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
			},
			"source": schema.StringAttribute{
				Description: "Source glob used to select candidate files.",
				Computed:    true,
				Optional:    true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
			},
			"exclude_pattern": schema.StringAttribute{
				Description: "Optional source exclusion glob.",
				Computed:    true,
				Optional:    true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
			},
			"disabled": schema.BoolAttribute{
				Description: "If true, the expectation is disabled.",
				Computed:    true,
				Optional:    true,
				PlanModifiers: []planmodifier.Bool{
					boolplanmodifier.UseStateForUnknown(),
				},
			},
			"trigger": schema.StringAttribute{
				Description: "How this expectation opens windows.",
				Computed:    true,
				Optional:    true,
				Validators: []validator.String{
					stringvalidator.OneOf("manual", "upload", "daily", "custom_schedule"),
				},
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
			},
			"interval": schema.StringAttribute{
				Description: "If trigger is `daily`, this specifies how often to run the expectation.",
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
				Description: "Times of day in HH:MM format for schedule-driven expectations.",
				Computed:    true,
				Optional:    true,
				ElementType: types.StringType,
				PlanModifiers: []planmodifier.List{
					listplanmodifier.UseStateForUnknown(),
				},
			},
			"schedule_time_zone": schema.StringAttribute{
				Description: "Time zone used by the expectation schedule.",
				Computed:    true,
				Optional:    true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
			},
			"holiday_region": schema.StringAttribute{
				Description: "Optional holiday region used by schedule-driven expectations.",
				Computed:    true,
				Optional:    true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
			},
			"lookback_interval": schema.Int64Attribute{
				Description: "How many seconds before the due boundary the window starts.",
				Computed:    true,
				Optional:    true,
				PlanModifiers: []planmodifier.Int64{
					int64planmodifier.UseStateForUnknown(),
				},
			},
			"late_acceptance_interval": schema.Int64Attribute{
				Description: "How many seconds a schedule-driven window may remain eligible to close as late.",
				Computed:    true,
				Optional:    true,
				PlanModifiers: []planmodifier.Int64{
					int64planmodifier.UseStateForUnknown(),
				},
			},
			"inactivity_interval": schema.Int64Attribute{
				Description: "How many quiet seconds are required before final closure.",
				Computed:    true,
				Optional:    true,
				PlanModifiers: []planmodifier.Int64{
					int64planmodifier.UseStateForUnknown(),
				},
			},
			"max_open_interval": schema.Int64Attribute{
				Description: "Hard-stop duration in seconds for unscheduled expectations.",
				Computed:    true,
				Optional:    true,
				PlanModifiers: []planmodifier.Int64{
					int64planmodifier.UseStateForUnknown(),
				},
			},
			"criteria": schema.DynamicAttribute{
				Description: "Structured criteria v1 definition for the expectation.",
				Computed:    true,
				Optional:    true,
				PlanModifiers: []planmodifier.Dynamic{
					dynamicplanmodifier.UseStateForUnknown(),
				},
			},
			"id": schema.Int64Attribute{
				Description: "Expectation ID",
				Computed:    true,
				PlanModifiers: []planmodifier.Int64{
					int64planmodifier.UseStateForUnknown(),
				},
			},
			"expectations_version": schema.Int64Attribute{
				Description: "Criteria schema version for this expectation.",
				Computed:    true,
			},
			"last_evaluated_at": schema.StringAttribute{
				Description: "Last time this expectation was evaluated.",
				Computed:    true,
			},
			"last_success_at": schema.StringAttribute{
				Description: "Last time this expectation closed successfully.",
				Computed:    true,
			},
			"last_failure_at": schema.StringAttribute{
				Description: "Last time this expectation closed with a failure result.",
				Computed:    true,
			},
			"last_result": schema.StringAttribute{
				Description: "Most recent terminal result for this expectation.",
				Computed:    true,
				Validators: []validator.String{
					stringvalidator.OneOf("success", "late", "missing", "invalid"),
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

func (r *expectationResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	var plan expectationResourceModel
	diags := req.Plan.Get(ctx, &plan)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}
	var config expectationResourceModel
	diags = req.Config.Get(ctx, &config)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	paramsExpectationCreate := files_sdk.ExpectationCreateParams{}
	paramsExpectationCreate.Name = plan.Name.ValueString()
	paramsExpectationCreate.Description = plan.Description.ValueString()
	paramsExpectationCreate.Path = plan.Path.ValueString()
	paramsExpectationCreate.Source = plan.Source.ValueString()
	paramsExpectationCreate.ExcludePattern = plan.ExcludePattern.ValueString()
	if !plan.Disabled.IsNull() && !plan.Disabled.IsUnknown() {
		paramsExpectationCreate.Disabled = plan.Disabled.ValueBoolPointer()
	}
	paramsExpectationCreate.Trigger = paramsExpectationCreate.Trigger.Enum()[plan.Trigger.ValueString()]
	paramsExpectationCreate.Interval = plan.Interval.ValueString()
	paramsExpectationCreate.RecurringDay = plan.RecurringDay.ValueInt64()
	if !plan.ScheduleDaysOfWeek.IsNull() && !plan.ScheduleDaysOfWeek.IsUnknown() {
		diags = plan.ScheduleDaysOfWeek.ElementsAs(ctx, &paramsExpectationCreate.ScheduleDaysOfWeek, false)
		resp.Diagnostics.Append(diags...)
	}
	if !plan.ScheduleTimesOfDay.IsNull() && !plan.ScheduleTimesOfDay.IsUnknown() {
		diags = plan.ScheduleTimesOfDay.ElementsAs(ctx, &paramsExpectationCreate.ScheduleTimesOfDay, false)
		resp.Diagnostics.Append(diags...)
	}
	paramsExpectationCreate.ScheduleTimeZone = plan.ScheduleTimeZone.ValueString()
	paramsExpectationCreate.HolidayRegion = plan.HolidayRegion.ValueString()
	paramsExpectationCreate.LookbackInterval = plan.LookbackInterval.ValueInt64()
	paramsExpectationCreate.LateAcceptanceInterval = plan.LateAcceptanceInterval.ValueInt64()
	paramsExpectationCreate.InactivityInterval = plan.InactivityInterval.ValueInt64()
	paramsExpectationCreate.MaxOpenInterval = plan.MaxOpenInterval.ValueInt64()
	createCriteria, diags := lib.DynamicToInterface(ctx, path.Root("criteria"), plan.Criteria)
	resp.Diagnostics.Append(diags...)
	paramsExpectationCreate.Criteria = createCriteria
	paramsExpectationCreate.WorkspaceId = plan.WorkspaceId.ValueInt64()

	if resp.Diagnostics.HasError() {
		return
	}

	expectation, err := r.client.Create(paramsExpectationCreate, files_sdk.WithContext(ctx))
	if err != nil {
		resp.Diagnostics.AddError(
			"Error Creating Files Expectation",
			"Could not create expectation, unexpected error: "+err.Error(),
		)
		return
	}

	diags = r.populateResourceModel(ctx, expectation, &plan)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	diags = resp.State.Set(ctx, plan)
	resp.Diagnostics.Append(diags...)
}

func (r *expectationResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	var state expectationResourceModel
	diags := req.State.Get(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	paramsExpectationFind := files_sdk.ExpectationFindParams{}
	paramsExpectationFind.Id = state.Id.ValueInt64()

	expectation, err := r.client.Find(paramsExpectationFind, files_sdk.WithContext(ctx))
	if err != nil {
		if files_sdk.IsNotExist(err) {
			resp.State.RemoveResource(ctx)
			return
		}

		resp.Diagnostics.AddError(
			"Error Reading Files Expectation",
			"Could not read expectation id "+fmt.Sprint(state.Id.ValueInt64())+": "+err.Error(),
		)
		return
	}

	diags = r.populateResourceModel(ctx, expectation, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	diags = resp.State.Set(ctx, state)
	resp.Diagnostics.Append(diags...)
}

func (r *expectationResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	var plan expectationResourceModel
	diags := req.Plan.Get(ctx, &plan)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}
	var config expectationResourceModel
	diags = req.Config.Get(ctx, &config)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	paramsExpectationUpdate := map[string]interface{}{}
	if !plan.Id.IsNull() && !plan.Id.IsUnknown() {
		paramsExpectationUpdate["id"] = plan.Id.ValueInt64()
	}
	if !config.Name.IsNull() && !config.Name.IsUnknown() {
		paramsExpectationUpdate["name"] = config.Name.ValueString()
	}
	if !config.Description.IsNull() && !config.Description.IsUnknown() {
		paramsExpectationUpdate["description"] = config.Description.ValueString()
	}
	if !config.Path.IsNull() && !config.Path.IsUnknown() {
		paramsExpectationUpdate["path"] = config.Path.ValueString()
	}
	if !config.Source.IsNull() && !config.Source.IsUnknown() {
		paramsExpectationUpdate["source"] = config.Source.ValueString()
	}
	if !config.ExcludePattern.IsNull() && !config.ExcludePattern.IsUnknown() {
		paramsExpectationUpdate["exclude_pattern"] = config.ExcludePattern.ValueString()
	}
	if !config.Disabled.IsNull() && !config.Disabled.IsUnknown() {
		paramsExpectationUpdate["disabled"] = config.Disabled.ValueBool()
	}
	if !config.Trigger.IsNull() && !config.Trigger.IsUnknown() {
		paramsExpectationUpdate["trigger"] = config.Trigger.ValueString()
	}
	if !config.Interval.IsNull() && !config.Interval.IsUnknown() {
		paramsExpectationUpdate["interval"] = config.Interval.ValueString()
	}
	if !config.RecurringDay.IsNull() && !config.RecurringDay.IsUnknown() {
		paramsExpectationUpdate["recurring_day"] = config.RecurringDay.ValueInt64()
	}
	if !config.ScheduleDaysOfWeek.IsNull() && !config.ScheduleDaysOfWeek.IsUnknown() {
		var updateScheduleDaysOfWeek []int64
		diags = config.ScheduleDaysOfWeek.ElementsAs(ctx, &updateScheduleDaysOfWeek, false)
		resp.Diagnostics.Append(diags...)
		paramsExpectationUpdate["schedule_days_of_week"] = updateScheduleDaysOfWeek
	}
	if !config.ScheduleTimesOfDay.IsNull() && !config.ScheduleTimesOfDay.IsUnknown() {
		var updateScheduleTimesOfDay []string
		diags = config.ScheduleTimesOfDay.ElementsAs(ctx, &updateScheduleTimesOfDay, false)
		resp.Diagnostics.Append(diags...)
		paramsExpectationUpdate["schedule_times_of_day"] = updateScheduleTimesOfDay
	}
	if !config.ScheduleTimeZone.IsNull() && !config.ScheduleTimeZone.IsUnknown() {
		paramsExpectationUpdate["schedule_time_zone"] = config.ScheduleTimeZone.ValueString()
	}
	if !config.HolidayRegion.IsNull() && !config.HolidayRegion.IsUnknown() {
		paramsExpectationUpdate["holiday_region"] = config.HolidayRegion.ValueString()
	}
	if !config.LookbackInterval.IsNull() && !config.LookbackInterval.IsUnknown() {
		paramsExpectationUpdate["lookback_interval"] = config.LookbackInterval.ValueInt64()
	}
	if !config.LateAcceptanceInterval.IsNull() && !config.LateAcceptanceInterval.IsUnknown() {
		paramsExpectationUpdate["late_acceptance_interval"] = config.LateAcceptanceInterval.ValueInt64()
	}
	if !config.InactivityInterval.IsNull() && !config.InactivityInterval.IsUnknown() {
		paramsExpectationUpdate["inactivity_interval"] = config.InactivityInterval.ValueInt64()
	}
	if !config.MaxOpenInterval.IsNull() && !config.MaxOpenInterval.IsUnknown() {
		paramsExpectationUpdate["max_open_interval"] = config.MaxOpenInterval.ValueInt64()
	}
	updateCriteria, diags := lib.DynamicToInterface(ctx, path.Root("criteria"), config.Criteria)
	resp.Diagnostics.Append(diags...)
	paramsExpectationUpdate["criteria"] = updateCriteria
	if !config.WorkspaceId.IsNull() && !config.WorkspaceId.IsUnknown() {
		paramsExpectationUpdate["workspace_id"] = config.WorkspaceId.ValueInt64()
	}

	if resp.Diagnostics.HasError() {
		return
	}

	expectation, err := r.client.UpdateWithMap(paramsExpectationUpdate, files_sdk.WithContext(ctx))
	if err != nil {
		resp.Diagnostics.AddError(
			"Error Updating Files Expectation",
			"Could not update expectation, unexpected error: "+err.Error(),
		)
		return
	}

	diags = r.populateResourceModel(ctx, expectation, &plan)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	diags = resp.State.Set(ctx, plan)
	resp.Diagnostics.Append(diags...)
}

func (r *expectationResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	var state expectationResourceModel
	diags := req.State.Get(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	paramsExpectationDelete := files_sdk.ExpectationDeleteParams{}
	paramsExpectationDelete.Id = state.Id.ValueInt64()

	err := r.client.Delete(paramsExpectationDelete, files_sdk.WithContext(ctx))
	if err != nil && !files_sdk.IsNotExist(err) {
		resp.Diagnostics.AddError(
			"Error Deleting Files Expectation",
			"Could not delete expectation id "+fmt.Sprint(state.Id.ValueInt64())+": "+err.Error(),
		)
	}
}

func (r *expectationResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
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

func (r *expectationResource) populateResourceModel(ctx context.Context, expectation files_sdk.Expectation, state *expectationResourceModel) (diags diag.Diagnostics) {
	var propDiags diag.Diagnostics

	state.Id = types.Int64Value(expectation.Id)
	state.WorkspaceId = types.Int64Value(expectation.WorkspaceId)
	state.Name = types.StringValue(expectation.Name)
	state.Description = types.StringValue(expectation.Description)
	state.Path = types.StringValue(expectation.Path)
	state.Source = types.StringValue(expectation.Source)
	state.ExcludePattern = types.StringValue(expectation.ExcludePattern)
	state.Disabled = types.BoolPointerValue(expectation.Disabled)
	state.ExpectationsVersion = types.Int64Value(expectation.ExpectationsVersion)
	state.Trigger = types.StringValue(expectation.Trigger)
	state.Interval = types.StringValue(expectation.Interval)
	state.RecurringDay = types.Int64Value(expectation.RecurringDay)
	state.ScheduleDaysOfWeek, propDiags = types.ListValueFrom(ctx, types.Int64Type, expectation.ScheduleDaysOfWeek)
	diags.Append(propDiags...)
	state.ScheduleTimesOfDay, propDiags = types.ListValueFrom(ctx, types.StringType, expectation.ScheduleTimesOfDay)
	diags.Append(propDiags...)
	state.ScheduleTimeZone = types.StringValue(expectation.ScheduleTimeZone)
	state.HolidayRegion = types.StringValue(expectation.HolidayRegion)
	state.LookbackInterval = types.Int64Value(expectation.LookbackInterval)
	state.LateAcceptanceInterval = types.Int64Value(expectation.LateAcceptanceInterval)
	state.InactivityInterval = types.Int64Value(expectation.InactivityInterval)
	state.MaxOpenInterval = types.Int64Value(expectation.MaxOpenInterval)
	state.Criteria, propDiags = lib.ToDynamic(ctx, path.Root("criteria"), expectation.Criteria, state.Criteria.UnderlyingValue())
	diags.Append(propDiags...)
	if err := lib.TimeToStringType(ctx, path.Root("last_evaluated_at"), expectation.LastEvaluatedAt, &state.LastEvaluatedAt); err != nil {
		diags.AddError(
			"Error Creating Files Expectation",
			"Could not convert state last_evaluated_at to string: "+err.Error(),
		)
	}
	if err := lib.TimeToStringType(ctx, path.Root("last_success_at"), expectation.LastSuccessAt, &state.LastSuccessAt); err != nil {
		diags.AddError(
			"Error Creating Files Expectation",
			"Could not convert state last_success_at to string: "+err.Error(),
		)
	}
	if err := lib.TimeToStringType(ctx, path.Root("last_failure_at"), expectation.LastFailureAt, &state.LastFailureAt); err != nil {
		diags.AddError(
			"Error Creating Files Expectation",
			"Could not convert state last_failure_at to string: "+err.Error(),
		)
	}
	state.LastResult = types.StringValue(expectation.LastResult)
	if err := lib.TimeToStringType(ctx, path.Root("created_at"), expectation.CreatedAt, &state.CreatedAt); err != nil {
		diags.AddError(
			"Error Creating Files Expectation",
			"Could not convert state created_at to string: "+err.Error(),
		)
	}
	if err := lib.TimeToStringType(ctx, path.Root("updated_at"), expectation.UpdatedAt, &state.UpdatedAt); err != nil {
		diags.AddError(
			"Error Creating Files Expectation",
			"Could not convert state updated_at to string: "+err.Error(),
		)
	}

	return
}
