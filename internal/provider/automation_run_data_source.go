package provider

import (
	"context"
	"fmt"

	files_sdk "github.com/Files-com/files-sdk-go/v3"
	automation_run "github.com/Files-com/files-sdk-go/v3/automationrun"
	"github.com/Files-com/terraform-provider-files/lib"
	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/datasource/schema"
	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

var (
	_ datasource.DataSource              = &automationRunDataSource{}
	_ datasource.DataSourceWithConfigure = &automationRunDataSource{}
)

func NewAutomationRunDataSource() datasource.DataSource {
	return &automationRunDataSource{}
}

type automationRunDataSource struct {
	client *automation_run.Client
}

type automationRunDataSourceModel struct {
	Id                   types.Int64  `tfsdk:"id"`
	AutomationId         types.Int64  `tfsdk:"automation_id"`
	CompletedAt          types.String `tfsdk:"completed_at"`
	CreatedAt            types.String `tfsdk:"created_at"`
	RetryAt              types.String `tfsdk:"retry_at"`
	RetriedAt            types.String `tfsdk:"retried_at"`
	RetriedInRunId       types.Int64  `tfsdk:"retried_in_run_id"`
	RetryOfRunId         types.Int64  `tfsdk:"retry_of_run_id"`
	Runtime              types.String `tfsdk:"runtime"`
	Status               types.String `tfsdk:"status"`
	SuccessfulOperations types.Int64  `tfsdk:"successful_operations"`
	FailedOperations     types.Int64  `tfsdk:"failed_operations"`
	StatusMessagesUrl    types.String `tfsdk:"status_messages_url"`
}

func (r *automationRunDataSource) Configure(_ context.Context, req datasource.ConfigureRequest, resp *datasource.ConfigureResponse) {
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

	r.client = &automation_run.Client{Config: sdk_config}
}

func (r *automationRunDataSource) Metadata(_ context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_automation_run"
}

func (r *automationRunDataSource) Schema(_ context.Context, _ datasource.SchemaRequest, resp *datasource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Description: "An AutomationRun is a record with information about a single execution of a given Automation.\n\n\n\nAutomation Runs can be retried upon `failure` or `partial_failure` by setting the `retry_on_failure_interval_in_minutes` and `retry_on_failure_number_of_attempts` settings on the associated Automation.\n\n\n\nWhen retries occur, a new AutomationRun will be created for each retry. The property `retry_at` on the original run, if set, represents when that run will be retried next. The properties `retried_at` and `retried_in_run_id` will be set in the original run that failed, at the time of retry. The property `retry_of_run_id` will be set in the new run.",
		Attributes: map[string]schema.Attribute{
			"id": schema.Int64Attribute{
				Description: "ID.",
				Required:    true,
			},
			"automation_id": schema.Int64Attribute{
				Description: "ID of the associated Automation.",
				Computed:    true,
			},
			"completed_at": schema.StringAttribute{
				Description: "Automation run completion/failure date/time.",
				Computed:    true,
			},
			"created_at": schema.StringAttribute{
				Description: "Automation run start date/time.",
				Computed:    true,
			},
			"retry_at": schema.StringAttribute{
				Description: "If set, this automation will be retried at this date/time due to `failure` or `partial_failure`.",
				Computed:    true,
			},
			"retried_at": schema.StringAttribute{
				Description: "If set, this Automation run was retried due to `failure` or `partial_failure`.",
				Computed:    true,
			},
			"retried_in_run_id": schema.Int64Attribute{
				Description: "ID of the run that is or will be retrying this run.",
				Computed:    true,
			},
			"retry_of_run_id": schema.Int64Attribute{
				Description: "ID of the original run that this run is retrying.",
				Computed:    true,
			},
			"runtime": schema.StringAttribute{
				Description: "Automation run runtime.",
				Computed:    true,
			},
			"status": schema.StringAttribute{
				Description: "The success status of the AutomationRun. One of `running`, `success`, `partial_failure`, or `failure`.",
				Computed:    true,
			},
			"successful_operations": schema.Int64Attribute{
				Description: "Count of successful operations.",
				Computed:    true,
			},
			"failed_operations": schema.Int64Attribute{
				Description: "Count of failed operations.",
				Computed:    true,
			},
			"status_messages_url": schema.StringAttribute{
				Description: "Link to status messages log file.",
				Computed:    true,
			},
		},
	}
}

func (r *automationRunDataSource) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
	var data automationRunDataSourceModel
	diags := req.Config.Get(ctx, &data)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	paramsAutomationRunFind := files_sdk.AutomationRunFindParams{}
	paramsAutomationRunFind.Id = data.Id.ValueInt64()

	automationRun, err := r.client.Find(paramsAutomationRunFind, files_sdk.WithContext(ctx))
	if err != nil {
		resp.Diagnostics.AddError(
			"Error Reading Files AutomationRun",
			"Could not read automation_run id "+fmt.Sprint(data.Id.ValueInt64())+": "+err.Error(),
		)
		return
	}

	diags = r.populateDataSourceModel(ctx, automationRun, &data)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	diags = resp.State.Set(ctx, data)
	resp.Diagnostics.Append(diags...)
}

func (r *automationRunDataSource) populateDataSourceModel(ctx context.Context, automationRun files_sdk.AutomationRun, state *automationRunDataSourceModel) (diags diag.Diagnostics) {
	state.Id = types.Int64Value(automationRun.Id)
	state.AutomationId = types.Int64Value(automationRun.AutomationId)
	if err := lib.TimeToStringType(ctx, path.Root("completed_at"), automationRun.CompletedAt, &state.CompletedAt); err != nil {
		diags.AddError(
			"Error Creating Files AutomationRun",
			"Could not convert state completed_at to string: "+err.Error(),
		)
	}
	if err := lib.TimeToStringType(ctx, path.Root("created_at"), automationRun.CreatedAt, &state.CreatedAt); err != nil {
		diags.AddError(
			"Error Creating Files AutomationRun",
			"Could not convert state created_at to string: "+err.Error(),
		)
	}
	if err := lib.TimeToStringType(ctx, path.Root("retry_at"), automationRun.RetryAt, &state.RetryAt); err != nil {
		diags.AddError(
			"Error Creating Files AutomationRun",
			"Could not convert state retry_at to string: "+err.Error(),
		)
	}
	if err := lib.TimeToStringType(ctx, path.Root("retried_at"), automationRun.RetriedAt, &state.RetriedAt); err != nil {
		diags.AddError(
			"Error Creating Files AutomationRun",
			"Could not convert state retried_at to string: "+err.Error(),
		)
	}
	state.RetriedInRunId = types.Int64Value(automationRun.RetriedInRunId)
	state.RetryOfRunId = types.Int64Value(automationRun.RetryOfRunId)
	state.Runtime = types.StringValue(automationRun.Runtime)
	state.Status = types.StringValue(automationRun.Status)
	state.SuccessfulOperations = types.Int64Value(automationRun.SuccessfulOperations)
	state.FailedOperations = types.Int64Value(automationRun.FailedOperations)
	state.StatusMessagesUrl = types.StringValue(automationRun.StatusMessagesUrl)

	return
}
