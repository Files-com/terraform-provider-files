package provider

import (
	"context"
	"fmt"

	files_sdk "github.com/Files-com/files-sdk-go/v3"
	partner_channel_template "github.com/Files-com/files-sdk-go/v3/partnerchanneltemplate"
	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/datasource/schema"
	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

var (
	_ datasource.DataSource              = &partnerChannelTemplateDataSource{}
	_ datasource.DataSourceWithConfigure = &partnerChannelTemplateDataSource{}
)

func NewPartnerChannelTemplateDataSource() datasource.DataSource {
	return &partnerChannelTemplateDataSource{}
}

type partnerChannelTemplateDataSource struct {
	client *partner_channel_template.Client
}

type partnerChannelTemplateDataSourceModel struct {
	Id                             types.Int64  `tfsdk:"id"`
	WorkspaceId                    types.Int64  `tfsdk:"workspace_id"`
	Name                           types.String `tfsdk:"name"`
	Path                           types.String `tfsdk:"path"`
	ToPartnerFolderName            types.String `tfsdk:"to_partner_folder_name"`
	FromPartnerFolderName          types.String `tfsdk:"from_partner_folder_name"`
	FromPartnerRoutePath           types.String `tfsdk:"from_partner_route_path"`
	ToPartnerRoutePath             types.String `tfsdk:"to_partner_route_path"`
	ToPartnerManagedFolderPaths    types.List   `tfsdk:"to_partner_managed_folder_paths"`
	FromPartnerManagedFolderPaths  types.List   `tfsdk:"from_partner_managed_folder_paths"`
	EffectiveToPartnerFolderName   types.String `tfsdk:"effective_to_partner_folder_name"`
	EffectiveFromPartnerFolderName types.String `tfsdk:"effective_from_partner_folder_name"`
}

func (r *partnerChannelTemplateDataSource) Configure(_ context.Context, req datasource.ConfigureRequest, resp *datasource.ConfigureResponse) {
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

	r.client = &partner_channel_template.Client{Config: sdk_config}
}

func (r *partnerChannelTemplateDataSource) Metadata(_ context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_partner_channel_template"
}

func (r *partnerChannelTemplateDataSource) Schema(_ context.Context, _ datasource.SchemaRequest, resp *datasource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Description: "A PartnerChannelTemplate defines reusable Partner Channel configuration that can be applied to Partners.",
		Attributes: map[string]schema.Attribute{
			"id": schema.Int64Attribute{
				Description: "The unique ID of the Partner Channel Template.",
				Required:    true,
			},
			"workspace_id": schema.Int64Attribute{
				Description: "ID of the Workspace associated with this Partner Channel Template.",
				Computed:    true,
			},
			"name": schema.StringAttribute{
				Description: "The name of the Partner Channel Template.",
				Computed:    true,
			},
			"path": schema.StringAttribute{
				Description: "Channel path relative to the Partner root folder. This must be slash-delimited, but it must neither start nor end with a slash. Maximum of 5000 characters.",
				Computed:    true,
			},
			"to_partner_folder_name": schema.StringAttribute{
				Description: "Optional Channel-level to-Partner folder name override.",
				Computed:    true,
			},
			"from_partner_folder_name": schema.StringAttribute{
				Description: "Optional Channel-level from-Partner folder name override.",
				Computed:    true,
			},
			"from_partner_route_path": schema.StringAttribute{
				Description: "Optional route path for files uploaded by the Partner.",
				Computed:    true,
			},
			"to_partner_route_path": schema.StringAttribute{
				Description: "Optional route path for files delivered to the Partner.",
				Computed:    true,
			},
			"to_partner_managed_folder_paths": schema.ListAttribute{
				Description: "Managed folder paths inside the to-Partner folder.",
				Computed:    true,
				ElementType: types.StringType,
			},
			"from_partner_managed_folder_paths": schema.ListAttribute{
				Description: "Managed folder paths inside the from-Partner folder.",
				Computed:    true,
				ElementType: types.StringType,
			},
			"effective_to_partner_folder_name": schema.StringAttribute{
				Description: "Resolved to-Partner folder name after Template override and default.",
				Computed:    true,
			},
			"effective_from_partner_folder_name": schema.StringAttribute{
				Description: "Resolved from-Partner folder name after Template override and default.",
				Computed:    true,
			},
		},
	}
}

func (r *partnerChannelTemplateDataSource) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
	var data partnerChannelTemplateDataSourceModel
	diags := req.Config.Get(ctx, &data)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	paramsPartnerChannelTemplateFind := files_sdk.PartnerChannelTemplateFindParams{}
	paramsPartnerChannelTemplateFind.Id = data.Id.ValueInt64()

	partnerChannelTemplate, err := r.client.Find(paramsPartnerChannelTemplateFind, files_sdk.WithContext(ctx))
	if err != nil {
		resp.Diagnostics.AddError(
			"Error Reading Files PartnerChannelTemplate",
			"Could not read partner_channel_template id "+fmt.Sprint(data.Id.ValueInt64())+": "+err.Error(),
		)
		return
	}

	diags = r.populateDataSourceModel(ctx, partnerChannelTemplate, &data)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	diags = resp.State.Set(ctx, data)
	resp.Diagnostics.Append(diags...)
}

func (r *partnerChannelTemplateDataSource) populateDataSourceModel(ctx context.Context, partnerChannelTemplate files_sdk.PartnerChannelTemplate, state *partnerChannelTemplateDataSourceModel) (diags diag.Diagnostics) {
	var propDiags diag.Diagnostics

	state.Id = types.Int64Value(partnerChannelTemplate.Id)
	state.WorkspaceId = types.Int64Value(partnerChannelTemplate.WorkspaceId)
	state.Name = types.StringValue(partnerChannelTemplate.Name)
	state.Path = types.StringValue(partnerChannelTemplate.Path)
	state.ToPartnerFolderName = types.StringValue(partnerChannelTemplate.ToPartnerFolderName)
	state.FromPartnerFolderName = types.StringValue(partnerChannelTemplate.FromPartnerFolderName)
	state.FromPartnerRoutePath = types.StringValue(partnerChannelTemplate.FromPartnerRoutePath)
	state.ToPartnerRoutePath = types.StringValue(partnerChannelTemplate.ToPartnerRoutePath)
	state.ToPartnerManagedFolderPaths, propDiags = types.ListValueFrom(ctx, types.StringType, partnerChannelTemplate.ToPartnerManagedFolderPaths)
	diags.Append(propDiags...)
	state.FromPartnerManagedFolderPaths, propDiags = types.ListValueFrom(ctx, types.StringType, partnerChannelTemplate.FromPartnerManagedFolderPaths)
	diags.Append(propDiags...)
	state.EffectiveToPartnerFolderName = types.StringValue(partnerChannelTemplate.EffectiveToPartnerFolderName)
	state.EffectiveFromPartnerFolderName = types.StringValue(partnerChannelTemplate.EffectiveFromPartnerFolderName)

	return
}
