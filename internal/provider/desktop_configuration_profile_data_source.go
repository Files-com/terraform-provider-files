package provider

import (
	"context"
	"fmt"

	files_sdk "github.com/Files-com/files-sdk-go/v3"
	desktop_configuration_profile "github.com/Files-com/files-sdk-go/v3/desktopconfigurationprofile"
	"github.com/Files-com/terraform-provider-files/lib"
	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/datasource/schema"
	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

var (
	_ datasource.DataSource              = &desktopConfigurationProfileDataSource{}
	_ datasource.DataSourceWithConfigure = &desktopConfigurationProfileDataSource{}
)

func NewDesktopConfigurationProfileDataSource() datasource.DataSource {
	return &desktopConfigurationProfileDataSource{}
}

type desktopConfigurationProfileDataSource struct {
	client *desktop_configuration_profile.Client
}

type desktopConfigurationProfileDataSourceModel struct {
	Id                   types.Int64   `tfsdk:"id"`
	Name                 types.String  `tfsdk:"name"`
	WorkspaceId          types.Int64   `tfsdk:"workspace_id"`
	UseForAllUsers       types.Bool    `tfsdk:"use_for_all_users"`
	DisableDriveMounting types.Bool    `tfsdk:"disable_drive_mounting"`
	MountMappings        types.Dynamic `tfsdk:"mount_mappings"`
}

func (r *desktopConfigurationProfileDataSource) Configure(_ context.Context, req datasource.ConfigureRequest, resp *datasource.ConfigureResponse) {
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

	r.client = &desktop_configuration_profile.Client{Config: sdk_config}
}

func (r *desktopConfigurationProfileDataSource) Metadata(_ context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_desktop_configuration_profile"
}

func (r *desktopConfigurationProfileDataSource) Schema(_ context.Context, _ datasource.SchemaRequest, resp *datasource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Description: "A Desktop Configuration Profile centrally defines desktop mount point mappings for users in a Site or Workspace.",
		Attributes: map[string]schema.Attribute{
			"id": schema.Int64Attribute{
				Description: "Desktop Configuration Profile ID",
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
			"disable_drive_mounting": schema.BoolAttribute{
				Description: "Whether the desktop app should hide drive mounting, prevent new drive mounts, and unmount active drive mounts for users with this profile",
				Computed:    true,
			},
			"mount_mappings": schema.DynamicAttribute{
				Description: "Mount point mappings for the desktop app. Keys must be a single uppercase Windows drive letter other than A, B, or C, and values are Files.com paths to mount there.",
				Computed:    true,
			},
		},
	}
}

func (r *desktopConfigurationProfileDataSource) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
	var data desktopConfigurationProfileDataSourceModel
	diags := req.Config.Get(ctx, &data)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	paramsDesktopConfigurationProfileFind := files_sdk.DesktopConfigurationProfileFindParams{}
	paramsDesktopConfigurationProfileFind.Id = data.Id.ValueInt64()

	desktopConfigurationProfile, err := r.client.Find(paramsDesktopConfigurationProfileFind, files_sdk.WithContext(ctx))
	if err != nil {
		resp.Diagnostics.AddError(
			"Error Reading Files DesktopConfigurationProfile",
			"Could not read desktop_configuration_profile id "+fmt.Sprint(data.Id.ValueInt64())+": "+err.Error(),
		)
		return
	}

	diags = r.populateDataSourceModel(ctx, desktopConfigurationProfile, &data)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	diags = resp.State.Set(ctx, data)
	resp.Diagnostics.Append(diags...)
}

func (r *desktopConfigurationProfileDataSource) populateDataSourceModel(ctx context.Context, desktopConfigurationProfile files_sdk.DesktopConfigurationProfile, state *desktopConfigurationProfileDataSourceModel) (diags diag.Diagnostics) {
	var propDiags diag.Diagnostics

	state.Id = types.Int64Value(desktopConfigurationProfile.Id)
	state.Name = types.StringValue(desktopConfigurationProfile.Name)
	state.WorkspaceId = types.Int64Value(desktopConfigurationProfile.WorkspaceId)
	state.UseForAllUsers = types.BoolPointerValue(desktopConfigurationProfile.UseForAllUsers)
	state.DisableDriveMounting = types.BoolPointerValue(desktopConfigurationProfile.DisableDriveMounting)
	state.MountMappings, propDiags = lib.ToDynamic(ctx, path.Root("mount_mappings"), desktopConfigurationProfile.MountMappings, state.MountMappings.UnderlyingValue())
	diags.Append(propDiags...)

	return
}
