package provider

import (
	"context"
	"fmt"

	files_sdk "github.com/Files-com/files-sdk-go/v3"
	workspace "github.com/Files-com/files-sdk-go/v3/workspace"
	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/datasource/schema"
	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

var (
	_ datasource.DataSource              = &workspaceDataSource{}
	_ datasource.DataSourceWithConfigure = &workspaceDataSource{}
)

func NewWorkspaceDataSource() datasource.DataSource {
	return &workspaceDataSource{}
}

type workspaceDataSource struct {
	client *workspace.Client
}

type workspaceDataSourceModel struct {
	Id   types.Int64  `tfsdk:"id"`
	Name types.String `tfsdk:"name"`
}

func (r *workspaceDataSource) Configure(_ context.Context, req datasource.ConfigureRequest, resp *datasource.ConfigureResponse) {
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

	r.client = &workspace.Client{Config: sdk_config}
}

func (r *workspaceDataSource) Metadata(_ context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_workspace"
}

func (r *workspaceDataSource) Schema(_ context.Context, _ datasource.SchemaRequest, resp *datasource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Description: "A Workspace is a lightweight way to organize related resources inside a single Files.com Site.\n\n\n\nCustomers commonly group resources by project, department, client, or region. Workspaces provide a built-in structure for that grouping, so the UI can operate within a clear “workspace context” and admins can delegate management for a subset of resources without requiring full site-level isolation.\n\n\n\nEvery Site has an implicit Default workspace (ID 0). Resources that are not explicitly assigned to a named workspace are considered part of the Default workspace.",
		Attributes: map[string]schema.Attribute{
			"id": schema.Int64Attribute{
				Description: "Workspace ID",
				Required:    true,
			},
			"name": schema.StringAttribute{
				Description: "Workspace name",
				Computed:    true,
			},
		},
	}
}

func (r *workspaceDataSource) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
	var data workspaceDataSourceModel
	diags := req.Config.Get(ctx, &data)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	paramsWorkspaceFind := files_sdk.WorkspaceFindParams{}
	paramsWorkspaceFind.Id = data.Id.ValueInt64()

	workspace, err := r.client.Find(paramsWorkspaceFind, files_sdk.WithContext(ctx))
	if err != nil {
		resp.Diagnostics.AddError(
			"Error Reading Files Workspace",
			"Could not read workspace id "+fmt.Sprint(data.Id.ValueInt64())+": "+err.Error(),
		)
		return
	}

	diags = r.populateDataSourceModel(ctx, workspace, &data)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	diags = resp.State.Set(ctx, data)
	resp.Diagnostics.Append(diags...)
}

func (r *workspaceDataSource) populateDataSourceModel(ctx context.Context, workspace files_sdk.Workspace, state *workspaceDataSourceModel) (diags diag.Diagnostics) {
	state.Id = types.Int64Value(workspace.Id)
	state.Name = types.StringValue(workspace.Name)

	return
}
