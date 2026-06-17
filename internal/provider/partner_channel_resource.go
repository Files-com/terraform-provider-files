package provider

import (
	"context"
	"fmt"
	"strconv"
	"strings"

	files_sdk "github.com/Files-com/files-sdk-go/v3"
	partner_channel "github.com/Files-com/files-sdk-go/v3/partnerchannel"
	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/int64planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

var (
	_ resource.Resource                = &partnerChannelResource{}
	_ resource.ResourceWithConfigure   = &partnerChannelResource{}
	_ resource.ResourceWithImportState = &partnerChannelResource{}
)

func NewPartnerChannelResource() resource.Resource {
	return &partnerChannelResource{}
}

type partnerChannelResource struct {
	client *partner_channel.Client
}

type partnerChannelResourceModel struct {
	PartnerId                      types.Int64  `tfsdk:"partner_id"`
	Path                           types.String `tfsdk:"path"`
	WorkspaceId                    types.Int64  `tfsdk:"workspace_id"`
	ToPartnerFolderName            types.String `tfsdk:"to_partner_folder_name"`
	FromPartnerFolderName          types.String `tfsdk:"from_partner_folder_name"`
	FromPartnerRoutePath           types.String `tfsdk:"from_partner_route_path"`
	ToPartnerRoutePath             types.String `tfsdk:"to_partner_route_path"`
	Id                             types.Int64  `tfsdk:"id"`
	EffectiveToPartnerFolderName   types.String `tfsdk:"effective_to_partner_folder_name"`
	EffectiveFromPartnerFolderName types.String `tfsdk:"effective_from_partner_folder_name"`
	ChannelPath                    types.String `tfsdk:"channel_path"`
	ToPartnerFolderPath            types.String `tfsdk:"to_partner_folder_path"`
	FromPartnerFolderPath          types.String `tfsdk:"from_partner_folder_path"`
}

func (r *partnerChannelResource) Configure(_ context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {
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

func (r *partnerChannelResource) Metadata(_ context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_partner_channel"
}

func (r *partnerChannelResource) Schema(_ context.Context, _ resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Description: "A PartnerChannel defines a structured communication path within a Partner root folder, including directional folder names and partner-scoped routing configuration.",
		Attributes: map[string]schema.Attribute{
			"partner_id": schema.Int64Attribute{
				Description: "ID of the Partner this Channel belongs to.",
				Required:    true,
				PlanModifiers: []planmodifier.Int64{
					int64planmodifier.RequiresReplace(),
				},
			},
			"path": schema.StringAttribute{
				Description: "Channel path relative to the Partner root folder. This must be slash-delimited, but it must neither start nor end with a slash. Maximum of 5000 characters.",
				Required:    true,
			},
			"workspace_id": schema.Int64Attribute{
				Description: "ID of the Workspace associated with this Partner Channel.",
				Computed:    true,
				Optional:    true,
				PlanModifiers: []planmodifier.Int64{
					int64planmodifier.UseStateForUnknown(),
					int64planmodifier.RequiresReplace(),
				},
			},
			"to_partner_folder_name": schema.StringAttribute{
				Description: "Optional Channel-level to-Partner folder name override.",
				Computed:    true,
				Optional:    true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
			},
			"from_partner_folder_name": schema.StringAttribute{
				Description: "Optional Channel-level from-Partner folder name override.",
				Computed:    true,
				Optional:    true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
			},
			"from_partner_route_path": schema.StringAttribute{
				Description: "Optional route path for files uploaded by the Partner.",
				Computed:    true,
				Optional:    true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
			},
			"to_partner_route_path": schema.StringAttribute{
				Description: "Optional route path for files delivered to the Partner.",
				Computed:    true,
				Optional:    true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
			},
			"id": schema.Int64Attribute{
				Description: "The unique ID of the Partner Channel.",
				Computed:    true,
				PlanModifiers: []planmodifier.Int64{
					int64planmodifier.UseStateForUnknown(),
				},
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

func (r *partnerChannelResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	var plan partnerChannelResourceModel
	diags := req.Plan.Get(ctx, &plan)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}
	var config partnerChannelResourceModel
	diags = req.Config.Get(ctx, &config)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	paramsPartnerChannelCreate := files_sdk.PartnerChannelCreateParams{}
	paramsPartnerChannelCreate.FromPartnerFolderName = plan.FromPartnerFolderName.ValueString()
	paramsPartnerChannelCreate.FromPartnerRoutePath = plan.FromPartnerRoutePath.ValueString()
	paramsPartnerChannelCreate.ToPartnerFolderName = plan.ToPartnerFolderName.ValueString()
	paramsPartnerChannelCreate.ToPartnerRoutePath = plan.ToPartnerRoutePath.ValueString()
	paramsPartnerChannelCreate.PartnerId = plan.PartnerId.ValueInt64()
	paramsPartnerChannelCreate.Path = plan.Path.ValueString()
	paramsPartnerChannelCreate.WorkspaceId = plan.WorkspaceId.ValueInt64()

	if resp.Diagnostics.HasError() {
		return
	}

	partnerChannel, err := r.client.Create(paramsPartnerChannelCreate, files_sdk.WithContext(ctx))
	if err != nil {
		resp.Diagnostics.AddError(
			"Error Creating Files PartnerChannel",
			"Could not create partner_channel, unexpected error: "+err.Error(),
		)
		return
	}

	diags = r.populateResourceModel(ctx, partnerChannel, &plan)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	diags = resp.State.Set(ctx, plan)
	resp.Diagnostics.Append(diags...)
}

func (r *partnerChannelResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	var state partnerChannelResourceModel
	diags := req.State.Get(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	paramsPartnerChannelFind := files_sdk.PartnerChannelFindParams{}
	paramsPartnerChannelFind.Id = state.Id.ValueInt64()

	partnerChannel, err := r.client.Find(paramsPartnerChannelFind, files_sdk.WithContext(ctx))
	if err != nil {
		if files_sdk.IsNotExist(err) {
			resp.State.RemoveResource(ctx)
			return
		}

		resp.Diagnostics.AddError(
			"Error Reading Files PartnerChannel",
			"Could not read partner_channel id "+fmt.Sprint(state.Id.ValueInt64())+": "+err.Error(),
		)
		return
	}

	diags = r.populateResourceModel(ctx, partnerChannel, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	diags = resp.State.Set(ctx, state)
	resp.Diagnostics.Append(diags...)
}

func (r *partnerChannelResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	var plan partnerChannelResourceModel
	diags := req.Plan.Get(ctx, &plan)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}
	var config partnerChannelResourceModel
	diags = req.Config.Get(ctx, &config)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	paramsPartnerChannelUpdate := map[string]interface{}{}
	if !plan.Id.IsNull() && !plan.Id.IsUnknown() {
		paramsPartnerChannelUpdate["id"] = plan.Id.ValueInt64()
	}
	if !config.FromPartnerFolderName.IsNull() && !config.FromPartnerFolderName.IsUnknown() {
		paramsPartnerChannelUpdate["from_partner_folder_name"] = config.FromPartnerFolderName.ValueString()
	}
	if !config.FromPartnerRoutePath.IsNull() && !config.FromPartnerRoutePath.IsUnknown() {
		paramsPartnerChannelUpdate["from_partner_route_path"] = config.FromPartnerRoutePath.ValueString()
	}
	if !config.ToPartnerFolderName.IsNull() && !config.ToPartnerFolderName.IsUnknown() {
		paramsPartnerChannelUpdate["to_partner_folder_name"] = config.ToPartnerFolderName.ValueString()
	}
	if !config.ToPartnerRoutePath.IsNull() && !config.ToPartnerRoutePath.IsUnknown() {
		paramsPartnerChannelUpdate["to_partner_route_path"] = config.ToPartnerRoutePath.ValueString()
	}
	if !config.Path.IsNull() && !config.Path.IsUnknown() {
		paramsPartnerChannelUpdate["path"] = config.Path.ValueString()
	}

	if resp.Diagnostics.HasError() {
		return
	}

	partnerChannel, err := r.client.UpdateWithMap(paramsPartnerChannelUpdate, files_sdk.WithContext(ctx))
	if err != nil {
		resp.Diagnostics.AddError(
			"Error Updating Files PartnerChannel",
			"Could not update partner_channel, unexpected error: "+err.Error(),
		)
		return
	}

	diags = r.populateResourceModel(ctx, partnerChannel, &plan)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	diags = resp.State.Set(ctx, plan)
	resp.Diagnostics.Append(diags...)
}

func (r *partnerChannelResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	var state partnerChannelResourceModel
	diags := req.State.Get(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	paramsPartnerChannelDelete := files_sdk.PartnerChannelDeleteParams{}
	paramsPartnerChannelDelete.Id = state.Id.ValueInt64()

	err := r.client.Delete(paramsPartnerChannelDelete, files_sdk.WithContext(ctx))
	if err != nil && !files_sdk.IsNotExist(err) {
		resp.Diagnostics.AddError(
			"Error Deleting Files PartnerChannel",
			"Could not delete partner_channel id "+fmt.Sprint(state.Id.ValueInt64())+": "+err.Error(),
		)
	}
}

func (r *partnerChannelResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	idParts := strings.SplitN(req.ID, ",", 1)

	if len(idParts) != 1 || idParts[0] == "" {
		resp.Diagnostics.AddError(
			"Unexpected Import Identifier",
			fmt.Sprintf("Expected import identifier with format: id. Got: %q", req.ID),
		)
		return
	}

	idPart, err := strconv.ParseFloat(idParts[0], 64)
	if err != nil {
		resp.Diagnostics.AddError(
			"Error Parsing ID",
			"Could not parse id: "+err.Error(),
		)
		return
	}
	resp.Diagnostics.Append(resp.State.SetAttribute(ctx, path.Root("id"), idPart)...)

}

func (r *partnerChannelResource) populateResourceModel(ctx context.Context, partnerChannel files_sdk.PartnerChannel, state *partnerChannelResourceModel) (diags diag.Diagnostics) {
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
