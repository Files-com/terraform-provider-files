package provider

import (
	"context"
	"fmt"

	files_sdk "github.com/Files-com/files-sdk-go/v3"
	project "github.com/Files-com/files-sdk-go/v3/project"
	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/datasource/schema"
	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

var (
	_ datasource.DataSource              = &projectDataSource{}
	_ datasource.DataSourceWithConfigure = &projectDataSource{}
)

func NewProjectDataSource() datasource.DataSource {
	return &projectDataSource{}
}

type projectDataSource struct {
	client *project.Client
}

type projectDataSourceModel struct {
	Id           types.Int64  `tfsdk:"id"`
	GlobalAccess types.String `tfsdk:"global_access"`
}

func (r *projectDataSource) Configure(_ context.Context, req datasource.ConfigureRequest, resp *datasource.ConfigureResponse) {
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

	r.client = &project.Client{Config: sdk_config}
}

func (r *projectDataSource) Metadata(_ context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_project"
}

func (r *projectDataSource) Schema(_ context.Context, _ datasource.SchemaRequest, resp *datasource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Description: "Projects are associated with a folder and add project management features to that folder.",
		Attributes: map[string]schema.Attribute{
			"id": schema.Int64Attribute{
				Description: "Project ID",
				Required:    true,
			},
			"global_access": schema.StringAttribute{
				Description: "Global access settings",
				Computed:    true,
			},
		},
	}
}

func (r *projectDataSource) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
	var data projectDataSourceModel
	diags := req.Config.Get(ctx, &data)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	paramsProjectFind := files_sdk.ProjectFindParams{}
	paramsProjectFind.Id = data.Id.ValueInt64()

	project, err := r.client.Find(paramsProjectFind, files_sdk.WithContext(ctx))
	if err != nil {
		resp.Diagnostics.AddError(
			"Error Reading Files Project",
			"Could not read project id "+fmt.Sprint(data.Id.ValueInt64())+": "+err.Error(),
		)
		return
	}

	diags = r.populateDataSourceModel(ctx, project, &data)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	diags = resp.State.Set(ctx, data)
	resp.Diagnostics.Append(diags...)
}

func (r *projectDataSource) populateDataSourceModel(ctx context.Context, project files_sdk.Project, state *projectDataSourceModel) (diags diag.Diagnostics) {
	state.Id = types.Int64Value(project.Id)
	state.GlobalAccess = types.StringValue(project.GlobalAccess)

	return
}
