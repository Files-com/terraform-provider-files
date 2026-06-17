package provider

import (
	"context"
	"fmt"

	files_sdk "github.com/Files-com/files-sdk-go/v3"
	partner_channel "github.com/Files-com/files-sdk-go/v3/partnerchannel"
	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/datasource/schema"
	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

var (
	_ datasource.DataSource              = &partnerChannelDataSource{}
	_ datasource.DataSourceWithConfigure = &partnerChannelDataSource{}
)

func NewPartnerChannelDataSource() datasource.DataSource {
	return &partnerChannelDataSource{}
}

type partnerChannelDataSource struct {
	client *partner_channel.Client
}

type partnerChannelDataSourceModel struct {
	Id                             types.Int64  `tfsdk:"id"`
	WorkspaceId                    types.Int64  `tfsdk:"workspace_id"`
	PartnerId                      types.Int64  `tfsdk:"partner_id"`
	Path                           types.String `tfsdk:"path"`
	ToPartnerFolderName            types.String `tfsdk:"to_partner_folder_name"`
	FromPartnerFolderName          types.String `tfsdk:"from_partner_folder_name"`
	FromPartnerRoutePath           types.String `tfsdk:"from_partner_route_path"`
	ToPartnerRoutePath             types.String `tfsdk:"to_partner_route_path"`
	EffectiveToPartnerFolderName   types.String `tfsdk:"effective_to_partner_folder_name"`
	EffectiveFromPartnerFolderName types.String `tfsdk:"effective_from_partner_folder_name"`
	ChannelPath                    types.String `tfsdk:"channel_path"`
	ToPartnerFolderPath            types.String `tfsdk:"to_partner_folder_path"`
	FromPartnerFolderPath          types.String `tfsdk:"from_partner_folder_path"`
}

func (r *partnerChannelDataSource) Configure(_ context.Context, req datasource.ConfigureRequest, resp *datasource.ConfigureResponse) {
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

	r.client = &partner_channel.Client{Config: sdk_config}
}

func (r *partnerChannelDataSource) Metadata(_ context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_partner_channel"
}

func (r *partnerChannelDataSource) Schema(_ context.Context, _ datasource.SchemaRequest, resp *datasource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Description: "A PartnerChannel defines a structured communication path within a Partner root folder, including directional folder names and partner-scoped routing configuration.",
		Attributes: map[string]schema.Attribute{
			"id": schema.Int64Attribute{
				Description: "The unique ID of the Partner Channel.",
				Required:    true,
			},
			"workspace_id": schema.Int64Attribute{
				Description: "ID of the Workspace associated with this Partner Channel.",
				Computed:    true,
			},
			"partner_id": schema.Int64Attribute{
				Description: "ID of the Partner this Channel belongs to.",
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
			"effective_to_partner_folder_name": schema.StringAttribute{
				Description: "Resolved to-Partner folder name after Channel override and default.",
				Computed:    true,
			},
			"effective_from_partner_folder_name": schema.StringAttribute{
				Description: "Resolved from-Partner folder name after Channel override and default.",
				Computed:    true,
			},
			"channel_path": schema.StringAttribute{
				Description: "Resolved Channel folder path.",
				Computed:    true,
			},
			"to_partner_folder_path": schema.StringAttribute{
				Description: "Resolved to-Partner folder path.",
				Computed:    true,
			},
			"from_partner_folder_path": schema.StringAttribute{
				Description: "Resolved from-Partner folder path.",
				Computed:    true,
			},
		},
	}
}

func (r *partnerChannelDataSource) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
	var data partnerChannelDataSourceModel
	diags := req.Config.Get(ctx, &data)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	paramsPartnerChannelFind := files_sdk.PartnerChannelFindParams{}
	paramsPartnerChannelFind.Id = data.Id.ValueInt64()

	partnerChannel, err := r.client.Find(paramsPartnerChannelFind, files_sdk.WithContext(ctx))
	if err != nil {
		resp.Diagnostics.AddError(
			"Error Reading Files PartnerChannel",
			"Could not read partner_channel id "+fmt.Sprint(data.Id.ValueInt64())+": "+err.Error(),
		)
		return
	}

	diags = r.populateDataSourceModel(ctx, partnerChannel, &data)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	diags = resp.State.Set(ctx, data)
	resp.Diagnostics.Append(diags...)
}

func (r *partnerChannelDataSource) populateDataSourceModel(ctx context.Context, partnerChannel files_sdk.PartnerChannel, state *partnerChannelDataSourceModel) (diags diag.Diagnostics) {
	state.Id = types.Int64Value(partnerChannel.Id)
	state.WorkspaceId = types.Int64Value(partnerChannel.WorkspaceId)
	state.PartnerId = types.Int64Value(partnerChannel.PartnerId)
	state.Path = types.StringValue(partnerChannel.Path)
	state.ToPartnerFolderName = types.StringValue(partnerChannel.ToPartnerFolderName)
	state.FromPartnerFolderName = types.StringValue(partnerChannel.FromPartnerFolderName)
	state.FromPartnerRoutePath = types.StringValue(partnerChannel.FromPartnerRoutePath)
	state.ToPartnerRoutePath = types.StringValue(partnerChannel.ToPartnerRoutePath)
	state.EffectiveToPartnerFolderName = types.StringValue(partnerChannel.EffectiveToPartnerFolderName)
	state.EffectiveFromPartnerFolderName = types.StringValue(partnerChannel.EffectiveFromPartnerFolderName)
	state.ChannelPath = types.StringValue(partnerChannel.ChannelPath)
	state.ToPartnerFolderPath = types.StringValue(partnerChannel.ToPartnerFolderPath)
	state.FromPartnerFolderPath = types.StringValue(partnerChannel.FromPartnerFolderPath)

	return
}
