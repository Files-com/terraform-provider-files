package provider

import (
	"context"
	"fmt"

	files_sdk "github.com/Files-com/files-sdk-go/v3"
	expectation_evaluation "github.com/Files-com/files-sdk-go/v3/expectationevaluation"
	"github.com/Files-com/terraform-provider-files/lib"
	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/datasource/schema"
	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

var (
	_ datasource.DataSource              = &expectationEvaluationDataSource{}
	_ datasource.DataSourceWithConfigure = &expectationEvaluationDataSource{}
)

func NewExpectationEvaluationDataSource() datasource.DataSource {
	return &expectationEvaluationDataSource{}
}

type expectationEvaluationDataSource struct {
	client *expectation_evaluation.Client
}

type expectationEvaluationDataSourceModel struct {
	Id                     types.Int64   `tfsdk:"id"`
	WorkspaceId            types.Int64   `tfsdk:"workspace_id"`
	ExpectationId          types.Int64   `tfsdk:"expectation_id"`
	Status                 types.String  `tfsdk:"status"`
	OpenedVia              types.String  `tfsdk:"opened_via"`
	OpenedAt               types.String  `tfsdk:"opened_at"`
	WindowStartAt          types.String  `tfsdk:"window_start_at"`
	WindowEndAt            types.String  `tfsdk:"window_end_at"`
	DeadlineAt             types.String  `tfsdk:"deadline_at"`
	LateAcceptanceCutoffAt types.String  `tfsdk:"late_acceptance_cutoff_at"`
	HardCloseAt            types.String  `tfsdk:"hard_close_at"`
	ClosedAt               types.String  `tfsdk:"closed_at"`
	MatchedFiles           types.Dynamic `tfsdk:"matched_files"`
	MissingFiles           types.Dynamic `tfsdk:"missing_files"`
	CriteriaErrors         types.List    `tfsdk:"criteria_errors"`
	Summary                types.Dynamic `tfsdk:"summary"`
	CreatedAt              types.String  `tfsdk:"created_at"`
	UpdatedAt              types.String  `tfsdk:"updated_at"`
}

func (r *expectationEvaluationDataSource) Configure(_ context.Context, req datasource.ConfigureRequest, resp *datasource.ConfigureResponse) {
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

	r.client = &expectation_evaluation.Client{Config: sdk_config}
}

func (r *expectationEvaluationDataSource) Metadata(_ context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_expectation_evaluation"
}

func (r *expectationEvaluationDataSource) Schema(_ context.Context, _ datasource.SchemaRequest, resp *datasource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Description: "An ExpectationEvaluation records one open or closed window for an Expectation.",
		Attributes: map[string]schema.Attribute{
			"id": schema.Int64Attribute{
				Description: "ExpectationEvaluation ID",
				Required:    true,
			},
			"workspace_id": schema.Int64Attribute{
				Description: "Workspace ID. `0` means the default workspace.",
				Computed:    true,
			},
			"expectation_id": schema.Int64Attribute{
				Description: "Expectation ID.",
				Computed:    true,
			},
			"status": schema.StringAttribute{
				Description: "Evaluation status.",
				Computed:    true,
			},
			"opened_via": schema.StringAttribute{
				Description: "How the evaluation window was opened.",
				Computed:    true,
			},
			"opened_at": schema.StringAttribute{
				Description: "When the evaluation row was opened.",
				Computed:    true,
			},
			"window_start_at": schema.StringAttribute{
				Description: "Logical start of the candidate window.",
				Computed:    true,
			},
			"window_end_at": schema.StringAttribute{
				Description: "Actual candidate cutoff boundary for the window.",
				Computed:    true,
			},
			"deadline_at": schema.StringAttribute{
				Description: "Logical due boundary for schedule-driven windows.",
				Computed:    true,
			},
			"late_acceptance_cutoff_at": schema.StringAttribute{
				Description: "Logical cutoff for late acceptance, when present.",
				Computed:    true,
			},
			"hard_close_at": schema.StringAttribute{
				Description: "Hard stop after which the window may not remain open.",
				Computed:    true,
			},
			"closed_at": schema.StringAttribute{
				Description: "When the evaluation row was finalized.",
				Computed:    true,
			},
			"matched_files": schema.DynamicAttribute{
				Description: "Captured evidence for files that matched the window.",
				Computed:    true,
			},
			"missing_files": schema.DynamicAttribute{
				Description: "Captured evidence for required files that were missing.",
				Computed:    true,
			},
			"criteria_errors": schema.ListAttribute{
				Description: "Captured criteria failures for the window.",
				Computed:    true,
				ElementType: types.StringType,
			},
			"summary": schema.DynamicAttribute{
				Description: "Compact evaluator summary payload.",
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

func (r *expectationEvaluationDataSource) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
	var data expectationEvaluationDataSourceModel
	diags := req.Config.Get(ctx, &data)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	paramsExpectationEvaluationFind := files_sdk.ExpectationEvaluationFindParams{}
	paramsExpectationEvaluationFind.Id = data.Id.ValueInt64()

	expectationEvaluation, err := r.client.Find(paramsExpectationEvaluationFind, files_sdk.WithContext(ctx))
	if err != nil {
		resp.Diagnostics.AddError(
			"Error Reading Files ExpectationEvaluation",
			"Could not read expectation_evaluation id "+fmt.Sprint(data.Id.ValueInt64())+": "+err.Error(),
		)
		return
	}

	diags = r.populateDataSourceModel(ctx, expectationEvaluation, &data)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	diags = resp.State.Set(ctx, data)
	resp.Diagnostics.Append(diags...)
}

func (r *expectationEvaluationDataSource) populateDataSourceModel(ctx context.Context, expectationEvaluation files_sdk.ExpectationEvaluation, state *expectationEvaluationDataSourceModel) (diags diag.Diagnostics) {
	var propDiags diag.Diagnostics

	state.Id = types.Int64Value(expectationEvaluation.Id)
	state.WorkspaceId = types.Int64Value(expectationEvaluation.WorkspaceId)
	state.ExpectationId = types.Int64Value(expectationEvaluation.ExpectationId)
	state.Status = types.StringValue(expectationEvaluation.Status)
	state.OpenedVia = types.StringValue(expectationEvaluation.OpenedVia)
	if err := lib.TimeToStringType(ctx, path.Root("opened_at"), expectationEvaluation.OpenedAt, &state.OpenedAt); err != nil {
		diags.AddError(
			"Error Creating Files ExpectationEvaluation",
			"Could not convert state opened_at to string: "+err.Error(),
		)
	}
	if err := lib.TimeToStringType(ctx, path.Root("window_start_at"), expectationEvaluation.WindowStartAt, &state.WindowStartAt); err != nil {
		diags.AddError(
			"Error Creating Files ExpectationEvaluation",
			"Could not convert state window_start_at to string: "+err.Error(),
		)
	}
	if err := lib.TimeToStringType(ctx, path.Root("window_end_at"), expectationEvaluation.WindowEndAt, &state.WindowEndAt); err != nil {
		diags.AddError(
			"Error Creating Files ExpectationEvaluation",
			"Could not convert state window_end_at to string: "+err.Error(),
		)
	}
	if err := lib.TimeToStringType(ctx, path.Root("deadline_at"), expectationEvaluation.DeadlineAt, &state.DeadlineAt); err != nil {
		diags.AddError(
			"Error Creating Files ExpectationEvaluation",
			"Could not convert state deadline_at to string: "+err.Error(),
		)
	}
	if err := lib.TimeToStringType(ctx, path.Root("late_acceptance_cutoff_at"), expectationEvaluation.LateAcceptanceCutoffAt, &state.LateAcceptanceCutoffAt); err != nil {
		diags.AddError(
			"Error Creating Files ExpectationEvaluation",
			"Could not convert state late_acceptance_cutoff_at to string: "+err.Error(),
		)
	}
	if err := lib.TimeToStringType(ctx, path.Root("hard_close_at"), expectationEvaluation.HardCloseAt, &state.HardCloseAt); err != nil {
		diags.AddError(
			"Error Creating Files ExpectationEvaluation",
			"Could not convert state hard_close_at to string: "+err.Error(),
		)
	}
	if err := lib.TimeToStringType(ctx, path.Root("closed_at"), expectationEvaluation.ClosedAt, &state.ClosedAt); err != nil {
		diags.AddError(
			"Error Creating Files ExpectationEvaluation",
			"Could not convert state closed_at to string: "+err.Error(),
		)
	}
	state.MatchedFiles, propDiags = lib.ToDynamic(ctx, path.Root("matched_files"), expectationEvaluation.MatchedFiles, state.MatchedFiles.UnderlyingValue())
	diags.Append(propDiags...)
	state.MissingFiles, propDiags = lib.ToDynamic(ctx, path.Root("missing_files"), expectationEvaluation.MissingFiles, state.MissingFiles.UnderlyingValue())
	diags.Append(propDiags...)
	state.CriteriaErrors, propDiags = types.ListValueFrom(ctx, types.StringType, expectationEvaluation.CriteriaErrors)
	diags.Append(propDiags...)
	state.Summary, propDiags = lib.ToDynamic(ctx, path.Root("summary"), expectationEvaluation.Summary, state.Summary.UnderlyingValue())
	diags.Append(propDiags...)
	if err := lib.TimeToStringType(ctx, path.Root("created_at"), expectationEvaluation.CreatedAt, &state.CreatedAt); err != nil {
		diags.AddError(
			"Error Creating Files ExpectationEvaluation",
			"Could not convert state created_at to string: "+err.Error(),
		)
	}
	if err := lib.TimeToStringType(ctx, path.Root("updated_at"), expectationEvaluation.UpdatedAt, &state.UpdatedAt); err != nil {
		diags.AddError(
			"Error Creating Files ExpectationEvaluation",
			"Could not convert state updated_at to string: "+err.Error(),
		)
	}

	return
}
