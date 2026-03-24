package provider

import (
	"context"
	"fmt"

	files_sdk "github.com/Files-com/files-sdk-go/v3"
	expectation_incident "github.com/Files-com/files-sdk-go/v3/expectationincident"
	"github.com/Files-com/terraform-provider-files/lib"
	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/datasource/schema"
	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

var (
	_ datasource.DataSource              = &expectationIncidentDataSource{}
	_ datasource.DataSourceWithConfigure = &expectationIncidentDataSource{}
)

func NewExpectationIncidentDataSource() datasource.DataSource {
	return &expectationIncidentDataSource{}
}

type expectationIncidentDataSource struct {
	client *expectation_incident.Client
}

type expectationIncidentDataSourceModel struct {
	Id                     types.Int64   `tfsdk:"id"`
	WorkspaceId            types.Int64   `tfsdk:"workspace_id"`
	ExpectationId          types.Int64   `tfsdk:"expectation_id"`
	Status                 types.String  `tfsdk:"status"`
	OpenedAt               types.String  `tfsdk:"opened_at"`
	LastFailedAt           types.String  `tfsdk:"last_failed_at"`
	AcknowledgedAt         types.String  `tfsdk:"acknowledged_at"`
	SnoozedUntil           types.String  `tfsdk:"snoozed_until"`
	ResolvedAt             types.String  `tfsdk:"resolved_at"`
	OpenedByEvaluationId   types.Int64   `tfsdk:"opened_by_evaluation_id"`
	LastEvaluationId       types.Int64   `tfsdk:"last_evaluation_id"`
	ResolvedByEvaluationId types.Int64   `tfsdk:"resolved_by_evaluation_id"`
	Summary                types.Dynamic `tfsdk:"summary"`
	CreatedAt              types.String  `tfsdk:"created_at"`
	UpdatedAt              types.String  `tfsdk:"updated_at"`
}

func (r *expectationIncidentDataSource) Configure(_ context.Context, req datasource.ConfigureRequest, resp *datasource.ConfigureResponse) {
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

	r.client = &expectation_incident.Client{Config: sdk_config}
}

func (r *expectationIncidentDataSource) Metadata(_ context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_expectation_incident"
}

func (r *expectationIncidentDataSource) Schema(_ context.Context, _ datasource.SchemaRequest, resp *datasource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Description: "An ExpectationIncident groups ongoing failure behavior for an Expectation over time.",
		Attributes: map[string]schema.Attribute{
			"id": schema.Int64Attribute{
				Description: "Expectation Incident ID",
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
				Description: "Incident status.",
				Computed:    true,
			},
			"opened_at": schema.StringAttribute{
				Description: "When the incident was opened.",
				Computed:    true,
			},
			"last_failed_at": schema.StringAttribute{
				Description: "When the most recent failing evaluation contributing to the incident occurred.",
				Computed:    true,
			},
			"acknowledged_at": schema.StringAttribute{
				Description: "When the incident was acknowledged.",
				Computed:    true,
			},
			"snoozed_until": schema.StringAttribute{
				Description: "When the current snooze expires.",
				Computed:    true,
			},
			"resolved_at": schema.StringAttribute{
				Description: "When the incident was resolved.",
				Computed:    true,
			},
			"opened_by_evaluation_id": schema.Int64Attribute{
				Description: "Evaluation that first opened the incident.",
				Computed:    true,
			},
			"last_evaluation_id": schema.Int64Attribute{
				Description: "Most recent evaluation linked to the incident.",
				Computed:    true,
			},
			"resolved_by_evaluation_id": schema.Int64Attribute{
				Description: "Evaluation that resolved the incident.",
				Computed:    true,
			},
			"summary": schema.DynamicAttribute{
				Description: "Compact incident summary payload.",
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

func (r *expectationIncidentDataSource) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
	var data expectationIncidentDataSourceModel
	diags := req.Config.Get(ctx, &data)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	paramsExpectationIncidentFind := files_sdk.ExpectationIncidentFindParams{}
	paramsExpectationIncidentFind.Id = data.Id.ValueInt64()

	expectationIncident, err := r.client.Find(paramsExpectationIncidentFind, files_sdk.WithContext(ctx))
	if err != nil {
		resp.Diagnostics.AddError(
			"Error Reading Files ExpectationIncident",
			"Could not read expectation_incident id "+fmt.Sprint(data.Id.ValueInt64())+": "+err.Error(),
		)
		return
	}

	diags = r.populateDataSourceModel(ctx, expectationIncident, &data)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	diags = resp.State.Set(ctx, data)
	resp.Diagnostics.Append(diags...)
}

func (r *expectationIncidentDataSource) populateDataSourceModel(ctx context.Context, expectationIncident files_sdk.ExpectationIncident, state *expectationIncidentDataSourceModel) (diags diag.Diagnostics) {
	var propDiags diag.Diagnostics

	state.Id = types.Int64Value(expectationIncident.Id)
	state.WorkspaceId = types.Int64Value(expectationIncident.WorkspaceId)
	state.ExpectationId = types.Int64Value(expectationIncident.ExpectationId)
	state.Status = types.StringValue(expectationIncident.Status)
	if err := lib.TimeToStringType(ctx, path.Root("opened_at"), expectationIncident.OpenedAt, &state.OpenedAt); err != nil {
		diags.AddError(
			"Error Creating Files ExpectationIncident",
			"Could not convert state opened_at to string: "+err.Error(),
		)
	}
	if err := lib.TimeToStringType(ctx, path.Root("last_failed_at"), expectationIncident.LastFailedAt, &state.LastFailedAt); err != nil {
		diags.AddError(
			"Error Creating Files ExpectationIncident",
			"Could not convert state last_failed_at to string: "+err.Error(),
		)
	}
	if err := lib.TimeToStringType(ctx, path.Root("acknowledged_at"), expectationIncident.AcknowledgedAt, &state.AcknowledgedAt); err != nil {
		diags.AddError(
			"Error Creating Files ExpectationIncident",
			"Could not convert state acknowledged_at to string: "+err.Error(),
		)
	}
	if err := lib.TimeToStringType(ctx, path.Root("snoozed_until"), expectationIncident.SnoozedUntil, &state.SnoozedUntil); err != nil {
		diags.AddError(
			"Error Creating Files ExpectationIncident",
			"Could not convert state snoozed_until to string: "+err.Error(),
		)
	}
	if err := lib.TimeToStringType(ctx, path.Root("resolved_at"), expectationIncident.ResolvedAt, &state.ResolvedAt); err != nil {
		diags.AddError(
			"Error Creating Files ExpectationIncident",
			"Could not convert state resolved_at to string: "+err.Error(),
		)
	}
	state.OpenedByEvaluationId = types.Int64Value(expectationIncident.OpenedByEvaluationId)
	state.LastEvaluationId = types.Int64Value(expectationIncident.LastEvaluationId)
	state.ResolvedByEvaluationId = types.Int64Value(expectationIncident.ResolvedByEvaluationId)
	state.Summary, propDiags = lib.ToDynamic(ctx, path.Root("summary"), expectationIncident.Summary, state.Summary.UnderlyingValue())
	diags.Append(propDiags...)
	if err := lib.TimeToStringType(ctx, path.Root("created_at"), expectationIncident.CreatedAt, &state.CreatedAt); err != nil {
		diags.AddError(
			"Error Creating Files ExpectationIncident",
			"Could not convert state created_at to string: "+err.Error(),
		)
	}
	if err := lib.TimeToStringType(ctx, path.Root("updated_at"), expectationIncident.UpdatedAt, &state.UpdatedAt); err != nil {
		diags.AddError(
			"Error Creating Files ExpectationIncident",
			"Could not convert state updated_at to string: "+err.Error(),
		)
	}

	return
}
