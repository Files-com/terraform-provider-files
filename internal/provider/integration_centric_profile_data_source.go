package provider

import (
	"context"
	"fmt"

	files_sdk "github.com/Files-com/files-sdk-go/v3"
	integration_centric_profile "github.com/Files-com/files-sdk-go/v3/integrationcentricprofile"
	"github.com/Files-com/terraform-provider-files/lib"
	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/datasource/schema"
	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

var (
	_ datasource.DataSource              = &integrationCentricProfileDataSource{}
	_ datasource.DataSourceWithConfigure = &integrationCentricProfileDataSource{}
)

func NewIntegrationCentricProfileDataSource() datasource.DataSource {
	return &integrationCentricProfileDataSource{}
}

type integrationCentricProfileDataSource struct {
	client *integration_centric_profile.Client
}

type integrationCentricProfileDataSourceModel struct {
	Id                    types.Int64   `tfsdk:"id"`
	Name                  types.String  `tfsdk:"name"`
	WorkspaceId           types.Int64   `tfsdk:"workspace_id"`
	UseForAllUsers        types.Bool    `tfsdk:"use_for_all_users"`
	ExpectedRemoteServers types.Dynamic `tfsdk:"expected_remote_servers"`
}

func (r *integrationCentricProfileDataSource) Configure(_ context.Context, req datasource.ConfigureRequest, resp *datasource.ConfigureResponse) {
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

	r.client = &integration_centric_profile.Client{Config: sdk_config}
}

func (r *integrationCentricProfileDataSource) Metadata(_ context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_integration_centric_profile"
}

func (r *integrationCentricProfileDataSource) Schema(_ context.Context, _ datasource.SchemaRequest, resp *datasource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Description: "An Integration Centric Profile defines the Remote Server integrations a user is expected to add and connect during integration-centric onboarding.\n\n\n\nUse this to automate setup guidance for users who need access to multiple business systems without sending long manual instructions. Common scenarios include ongoing access to systems such as SharePoint, bridging Google, Microsoft, and Box environments after M&A activity, and migrations where users connect legacy EFSS accounts during transition work.",
		Attributes: map[string]schema.Attribute{
			"id": schema.Int64Attribute{
				Description: "Integration Centric Profile ID",
				Required:    true,
			},
			"name": schema.StringAttribute{
				Description: "Profile name",
				Computed:    true,
			},
			"workspace_id": schema.Int64Attribute{
				Description: "Workspace ID",
				Computed:    true,
			},
			"use_for_all_users": schema.BoolAttribute{
				Description: "Whether this profile applies to all users in the Workspace by default",
				Computed:    true,
			},
			"expected_remote_servers": schema.DynamicAttribute{
				Description: "Remote Server integrations the user is expected to add and connect. Each entry requires `server_type` and may include a display `name`.",
				Computed:    true,
			},
		},
	}
}

func (r *integrationCentricProfileDataSource) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
	var data integrationCentricProfileDataSourceModel
	diags := req.Config.Get(ctx, &data)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	paramsIntegrationCentricProfileFind := files_sdk.IntegrationCentricProfileFindParams{}
	paramsIntegrationCentricProfileFind.Id = data.Id.ValueInt64()

	integrationCentricProfile, err := r.client.Find(paramsIntegrationCentricProfileFind, files_sdk.WithContext(ctx))
	if err != nil {
		resp.Diagnostics.AddError(
			"Error Reading Files IntegrationCentricProfile",
			"Could not read integration_centric_profile id "+fmt.Sprint(data.Id.ValueInt64())+": "+err.Error(),
		)
		return
	}

	diags = r.populateDataSourceModel(ctx, integrationCentricProfile, &data)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	diags = resp.State.Set(ctx, data)
	resp.Diagnostics.Append(diags...)
}

func (r *integrationCentricProfileDataSource) populateDataSourceModel(ctx context.Context, integrationCentricProfile files_sdk.IntegrationCentricProfile, state *integrationCentricProfileDataSourceModel) (diags diag.Diagnostics) {
	var propDiags diag.Diagnostics

	state.Id = types.Int64Value(integrationCentricProfile.Id)
	state.Name = types.StringValue(integrationCentricProfile.Name)
	state.WorkspaceId = types.Int64Value(integrationCentricProfile.WorkspaceId)
	state.UseForAllUsers = types.BoolPointerValue(integrationCentricProfile.UseForAllUsers)
	state.ExpectedRemoteServers, propDiags = lib.ToDynamic(ctx, path.Root("expected_remote_servers"), integrationCentricProfile.ExpectedRemoteServers, state.ExpectedRemoteServers.UnderlyingValue())
	diags.Append(propDiags...)

	return
}
