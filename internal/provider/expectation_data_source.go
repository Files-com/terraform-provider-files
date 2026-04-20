package provider

import (
	"context"
	"fmt"

	files_sdk "github.com/Files-com/files-sdk-go/v3"
	expectation "github.com/Files-com/files-sdk-go/v3/expectation"
	"github.com/Files-com/terraform-provider-files/lib"
	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/datasource/schema"
	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

var (
	_ datasource.DataSource              = &expectationDataSource{}
	_ datasource.DataSourceWithConfigure = &expectationDataSource{}
)

func NewExpectationDataSource() datasource.DataSource {
	return &expectationDataSource{}
}

type expectationDataSource struct {
	client *expectation.Client
}

type expectationDataSourceModel struct {
	Id                     types.Int64   `tfsdk:"id"`
	WorkspaceId            types.Int64   `tfsdk:"workspace_id"`
	Name                   types.String  `tfsdk:"name"`
	Description            types.String  `tfsdk:"description"`
	Path                   types.String  `tfsdk:"path"`
	Source                 types.String  `tfsdk:"source"`
	ExcludePattern         types.String  `tfsdk:"exclude_pattern"`
	Disabled               types.Bool    `tfsdk:"disabled"`
	ExpectationsVersion    types.Int64   `tfsdk:"expectations_version"`
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
	LastEvaluatedAt        types.String  `tfsdk:"last_evaluated_at"`
	LastSuccessAt          types.String  `tfsdk:"last_success_at"`
	LastFailureAt          types.String  `tfsdk:"last_failure_at"`
	LastResult             types.String  `tfsdk:"last_result"`
	CreatedAt              types.String  `tfsdk:"created_at"`
	UpdatedAt              types.String  `tfsdk:"updated_at"`
}

func (r *expectationDataSource) Configure(_ context.Context, req datasource.ConfigureRequest, resp *datasource.ConfigureResponse) {
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

func (r *expectationDataSource) Metadata(_ context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_expectation"
}

func (r *expectationDataSource) Schema(_ context.Context, _ datasource.SchemaRequest, resp *datasource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Description: "Expectations let your Files.com site define what “correct” file delivery looks like, continuously evaluate whether it happened, and keep history when it did not.\n\n\n\nExpectations are meant to answer operational questions like:\n\n\n\n* Did the expected file arrive?\n\n* Was it on time?\n\n* Did it meet the required shape and count rules?\n\n* Is there an active issue someone needs to acknowledge?\n\n\n\nExpectations are different from Automations and Syncs. Automations and Syncs act on files; Expectations monitor whether expected files arrived on time, in the right place, and in the right shape. In practice, Expectations are the sensor and Automations are the actuator.\n\n\n\nAn Expectation combines four concepts:\n\n\n\n1. **Scope**: where to look for candidate files, using `path`, `source`, and optional `exclude_pattern`.\n\n2. **Trigger / timing**: when a window opens and how long it stays eligible, using `trigger`, schedule fields, `lookback_interval`, `late_acceptance_interval`, `inactivity_interval`, and `max_open_interval`.\n\n3. **Criteria**: what must be true for the window to succeed, using the structured `criteria` JSON document.\n\n4. **Outcome history**: what happened over time, exposed through `ExpectationEvaluation` history and `ExpectationIncident` lifecycle records.\n\n\n\n## Scope and matching\n\n\n\nExpectations reuse the familiar Files.com path-plus-glob model.\n\n\n\nThe `path` field identifies the folder scope, while `source` identifies which files within that scope are candidates. `exclude_pattern` removes files from consideration.\n\n\n\nLike Automations, these fields support glob-style matching. Expectations treat those matches as one logical candidate set for each window. A single Expectation does not implicitly fan out into separate per-customer or per-folder evaluations just because the path contains wildcards.\n\n\n\n## Expectation windows\n\n\n\nExpectations are evaluated in windows.\n\n\n\nEach window is persisted as an `ExpectationEvaluation` record. A window opens, remains `open` while evidence can still arrive, and then closes into a terminal result such as `success`, `late`, `missing`, or `invalid`.\n\n\n\nAn Expectation has only one open window at a time.\n\n\n\n## Trigger modes\n\n\n\nExpectations can open windows in three ways:\n\n\n\n* `daily`: run on a recurring daily/weekly/monthly/quarterly/yearly cadence using `interval` and `recurring_day`.\n\n* `custom_schedule`: run on specific weekdays and times using `schedule_days_of_week`, `schedule_times_of_day`, and optional `schedule_time_zone` / `holiday_region`.\n\n* `manual`: an operator explicitly opens the window.\n\n\n\nSchedule-driven expectations define an on-time deadline and may optionally remain eligible to close as `late` during `late_acceptance_interval`.\n\n\n\nManual expectations have no concept of `late`; they open when triggered and close based on inactivity or hard-stop timing.\n\n\n\n## Success criteria\n\n\n\nThe `criteria` field is a structured JSON object describing what counts as success for the window.\n\n\n\nIn criteria v1, this can express things like:\n\n\n\n* file count constraints\n\n* total byte constraints\n\n* allowed extensions\n\n* filename regex validation\n\n* forbidden files\n\n* required named or globbed files with their own per-file constraints\n\n\n\nRequired file rule keys may also include standard strftime-style date/time tokens like `%Y`, `%m`, and `%d`. Those tokens are resolved at evaluation time using a stable window anchor: schedule-driven expectations use the window's `deadline_at`, while manual and upload expectations use the window's `opened_at`.\n\n\n\nThis is intentionally structured rather than scriptable, so it stays safe, explainable, and versionable.\n\n\n\n## History and incidents\n\n\n\nThe Expectation itself stores summary state like `last_evaluated_at`, `last_success_at`, `last_failure_at`, and `last_result`.\n\n\n\nFor deeper inspection:\n\n\n\n* `ExpectationEvaluation` history shows each open or closed window and the evidence captured for it.\n\n* `ExpectationIncident` records track ongoing failure situations over time, including acknowledge, snooze, and resolve actions.\n\n\n\nManual windows do not open incidents in v1. Schedule-driven failures can open incidents, and later qualifying success can resolve them.",
		Attributes: map[string]schema.Attribute{
			"id": schema.Int64Attribute{
				Description: "Expectation ID",
				Required:    true,
			},
			"workspace_id": schema.Int64Attribute{
				Description: "Workspace ID. `0` means the default workspace.",
				Computed:    true,
			},
			"name": schema.StringAttribute{
				Description: "Expectation name.",
				Computed:    true,
			},
			"description": schema.StringAttribute{
				Description: "Expectation description.",
				Computed:    true,
			},
			"path": schema.StringAttribute{
				Description: "Path scope for the expectation. Supports workspace-relative presentation. This must be slash-delimited, but it must neither start nor end with a slash. Maximum of 5000 characters.",
				Computed:    true,
			},
			"source": schema.StringAttribute{
				Description: "Source glob used to select candidate files.",
				Computed:    true,
			},
			"exclude_pattern": schema.StringAttribute{
				Description: "Optional source exclusion glob.",
				Computed:    true,
			},
			"disabled": schema.BoolAttribute{
				Description: "If true, the expectation is disabled.",
				Computed:    true,
			},
			"expectations_version": schema.Int64Attribute{
				Description: "Criteria schema version for this expectation.",
				Computed:    true,
			},
			"trigger": schema.StringAttribute{
				Description: "How this expectation opens windows.",
				Computed:    true,
			},
			"interval": schema.StringAttribute{
				Description: "If trigger is `daily`, this specifies how often to run the expectation.",
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
				Description: "Times of day in HH:MM format for schedule-driven expectations.",
				Computed:    true,
				ElementType: types.StringType,
			},
			"schedule_time_zone": schema.StringAttribute{
				Description: "Time zone used by the expectation schedule.",
				Computed:    true,
			},
			"holiday_region": schema.StringAttribute{
				Description: "Optional holiday region used by schedule-driven expectations.",
				Computed:    true,
			},
			"lookback_interval": schema.Int64Attribute{
				Description: "How many seconds before the due boundary the window starts.",
				Computed:    true,
			},
			"late_acceptance_interval": schema.Int64Attribute{
				Description: "How many seconds a schedule-driven window may remain eligible to close as late.",
				Computed:    true,
			},
			"inactivity_interval": schema.Int64Attribute{
				Description: "How many quiet seconds are required before final closure.",
				Computed:    true,
			},
			"max_open_interval": schema.Int64Attribute{
				Description: "Hard-stop duration in seconds for unscheduled expectations.",
				Computed:    true,
			},
			"criteria": schema.DynamicAttribute{
				Description: "Structured criteria v1 definition for the expectation.",
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

func (r *expectationDataSource) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
	var data expectationDataSourceModel
	diags := req.Config.Get(ctx, &data)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	paramsExpectationFind := files_sdk.ExpectationFindParams{}
	paramsExpectationFind.Id = data.Id.ValueInt64()

	expectation, err := r.client.Find(paramsExpectationFind, files_sdk.WithContext(ctx))
	if err != nil {
		resp.Diagnostics.AddError(
			"Error Reading Files Expectation",
			"Could not read expectation id "+fmt.Sprint(data.Id.ValueInt64())+": "+err.Error(),
		)
		return
	}

	diags = r.populateDataSourceModel(ctx, expectation, &data)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	diags = resp.State.Set(ctx, data)
	resp.Diagnostics.Append(diags...)
}

func (r *expectationDataSource) populateDataSourceModel(ctx context.Context, expectation files_sdk.Expectation, state *expectationDataSourceModel) (diags diag.Diagnostics) {
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
