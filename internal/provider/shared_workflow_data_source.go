package provider

import (
	"context"
	"fmt"

	files_sdk "github.com/Files-com/files-sdk-go/v3"
	shared_workflow "github.com/Files-com/files-sdk-go/v3/sharedworkflow"
	"github.com/Files-com/terraform-provider-files/lib"

	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/datasource/schema"
	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/path"

	"github.com/hashicorp/terraform-plugin-framework/types"
)

var (
	_ datasource.DataSource              = &sharedWorkflowDataSource{}
	_ datasource.DataSourceWithConfigure = &sharedWorkflowDataSource{}
)

func NewSharedWorkflowDataSource() datasource.DataSource {
	return &sharedWorkflowDataSource{}
}

type sharedWorkflowDataSource struct {
	client *shared_workflow.Client
}

type sharedWorkflowDataSourceModel struct {
	Id                 types.Int64   `tfsdk:"id"`
	Name               types.String  `tfsdk:"name"`
	Description        types.String  `tfsdk:"description"`
	SubmissionRevision types.String  `tfsdk:"submission_revision"`
	FormFields         types.Dynamic `tfsdk:"form_fields"`
}

func (r *sharedWorkflowDataSource) Configure(_ context.Context, req datasource.ConfigureRequest, resp *datasource.ConfigureResponse) {
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

	r.client = &shared_workflow.Client{Config: sdk_config}
}

func (r *sharedWorkflowDataSource) Metadata(_ context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_shared_workflow"
}

func (r *sharedWorkflowDataSource) Schema(_ context.Context, _ datasource.SchemaRequest, resp *datasource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Description: "",
		Attributes: map[string]schema.Attribute{
			"id": schema.Int64Attribute{
				Description: "Shared Workflow ID.",
				Required:    true,
			},
			"name": schema.StringAttribute{
				Description: "Name displayed to recipients.",
				Computed:    true,
			},
			"description": schema.StringAttribute{
				Description: "Description displayed to recipients.",
				Computed:    true,
			},
			"submission_revision": schema.StringAttribute{
				Description: "Revision required when submitting a Shared Workflow run.",
				Computed:    true,
			},
			"form_fields": schema.DynamicAttribute{
				Description: "Current custom questions in display order.",
				Computed:    true,
			},
		},
	}
}

func (r *sharedWorkflowDataSource) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
	var data sharedWorkflowDataSourceModel
	diags := req.Config.Get(ctx, &data)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	paramsSharedWorkflowFind := files_sdk.SharedWorkflowFindParams{}
	paramsSharedWorkflowFind.Id = data.Id.ValueInt64()

	sharedWorkflow, err := r.client.Find(paramsSharedWorkflowFind, files_sdk.WithContext(ctx))
	if err != nil {
		resp.Diagnostics.AddError(
			"Error Reading Files SharedWorkflow",
			"Could not read shared_workflow id "+fmt.Sprint(data.Id.ValueInt64())+": "+err.Error(),
		)
		return
	}

	diags = r.populateDataSourceModel(ctx, sharedWorkflow, &data)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	diags = resp.State.Set(ctx, data)
	resp.Diagnostics.Append(diags...)
}

func (r *sharedWorkflowDataSource) populateDataSourceModel(ctx context.Context, sharedWorkflow files_sdk.SharedWorkflow, state *sharedWorkflowDataSourceModel) (diags diag.Diagnostics) {
	var propDiags diag.Diagnostics

	state.Id = types.Int64Value(sharedWorkflow.Id)
	state.Name = types.StringValue(sharedWorkflow.Name)
	state.Description = types.StringValue(sharedWorkflow.Description)
	state.SubmissionRevision = types.StringValue(sharedWorkflow.SubmissionRevision)
	state.FormFields, propDiags = lib.ToDynamic(ctx, path.Root("form_fields"), sharedWorkflow.FormFields, state.FormFields.UnderlyingValue())
	diags.Append(propDiags...)

	return
}
