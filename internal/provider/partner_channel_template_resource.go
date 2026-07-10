package provider

import (
	"context"
	"fmt"
	"strconv"
	"strings"

	files_sdk "github.com/Files-com/files-sdk-go/v3"
	partner_channel_template "github.com/Files-com/files-sdk-go/v3/partnerchanneltemplate"
	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/int64planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/listplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

var (
	_ resource.Resource                = &partnerChannelTemplateResource{}
	_ resource.ResourceWithConfigure   = &partnerChannelTemplateResource{}
	_ resource.ResourceWithImportState = &partnerChannelTemplateResource{}
)

func NewPartnerChannelTemplateResource() resource.Resource {
	return &partnerChannelTemplateResource{}
}

type partnerChannelTemplateResource struct {
	client *partner_channel_template.Client
}

type partnerChannelTemplateResourceModel struct {
	Name                           types.String `tfsdk:"name"`
	Path                           types.String `tfsdk:"path"`
	WorkspaceId                    types.Int64  `tfsdk:"workspace_id"`
	ToPartnerFolderName            types.String `tfsdk:"to_partner_folder_name"`
	FromPartnerFolderName          types.String `tfsdk:"from_partner_folder_name"`
	FromPartnerRoutePathPattern    types.String `tfsdk:"from_partner_route_path_pattern"`
	ToPartnerRoutePathPattern      types.String `tfsdk:"to_partner_route_path_pattern"`
	ToPartnerManagedFolderPaths    types.List   `tfsdk:"to_partner_managed_folder_paths"`
	FromPartnerManagedFolderPaths  types.List   `tfsdk:"from_partner_managed_folder_paths"`
	Id                             types.Int64  `tfsdk:"id"`
	EffectiveToPartnerFolderName   types.String `tfsdk:"effective_to_partner_folder_name"`
	EffectiveFromPartnerFolderName types.String `tfsdk:"effective_from_partner_folder_name"`
}

func (r *partnerChannelTemplateResource) Configure(_ context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {
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

func (r *partnerChannelTemplateResource) Metadata(_ context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_partner_channel_template"
}

func (r *partnerChannelTemplateResource) Schema(_ context.Context, _ resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Description: "A PartnerChannelTemplate defines reusable Partner Channel configuration that can be applied to Partners.",
		Attributes: map[string]schema.Attribute{
			"name": schema.StringAttribute{
				Description: "The name of the Partner Channel Template.",
				Required:    true,
			},
			"path": schema.StringAttribute{
				Description: "Channel path relative to the Partner root folder. This must be slash-delimited, but it must neither start nor end with a slash. Maximum of 5000 characters.",
				Required:    true,
			},
			"workspace_id": schema.Int64Attribute{
				Description: "ID of the Workspace associated with this Partner Channel Template.",
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
			"from_partner_route_path_pattern": schema.StringAttribute{
				Description: "Optional route path pattern for files uploaded by the Partner. Supports {{partner_name}}.",
				Computed:    true,
				Optional:    true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
			},
			"to_partner_route_path_pattern": schema.StringAttribute{
				Description: "Optional route path pattern for files delivered to the Partner. Supports {{partner_name}}.",
				Computed:    true,
				Optional:    true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
			},
			"to_partner_managed_folder_paths": schema.ListAttribute{
				Description: "Managed folder paths inside the to-Partner folder.",
				Computed:    true,
				Optional:    true,
				ElementType: types.StringType,
				PlanModifiers: []planmodifier.List{
					listplanmodifier.UseStateForUnknown(),
				},
			},
			"from_partner_managed_folder_paths": schema.ListAttribute{
				Description: "Managed folder paths inside the from-Partner folder.",
				Computed:    true,
				Optional:    true,
				ElementType: types.StringType,
				PlanModifiers: []planmodifier.List{
					listplanmodifier.UseStateForUnknown(),
				},
			},
			"id": schema.Int64Attribute{
				Description: "The unique ID of the Partner Channel Template.",
				Computed:    true,
				PlanModifiers: []planmodifier.Int64{
					int64planmodifier.UseStateForUnknown(),
				},
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

func (r *partnerChannelTemplateResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	var plan partnerChannelTemplateResourceModel
	diags := req.Plan.Get(ctx, &plan)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}
	var config partnerChannelTemplateResourceModel
	diags = req.Config.Get(ctx, &config)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	paramsPartnerChannelTemplateCreate := files_sdk.PartnerChannelTemplateCreateParams{}
	paramsPartnerChannelTemplateCreate.FromPartnerFolderName = plan.FromPartnerFolderName.ValueString()
	if !plan.FromPartnerManagedFolderPaths.IsNull() && !plan.FromPartnerManagedFolderPaths.IsUnknown() {
		diags = plan.FromPartnerManagedFolderPaths.ElementsAs(ctx, &paramsPartnerChannelTemplateCreate.FromPartnerManagedFolderPaths, false)
		resp.Diagnostics.Append(diags...)
	}
	paramsPartnerChannelTemplateCreate.FromPartnerRoutePathPattern = plan.FromPartnerRoutePathPattern.ValueString()
	paramsPartnerChannelTemplateCreate.ToPartnerFolderName = plan.ToPartnerFolderName.ValueString()
	if !plan.ToPartnerManagedFolderPaths.IsNull() && !plan.ToPartnerManagedFolderPaths.IsUnknown() {
		diags = plan.ToPartnerManagedFolderPaths.ElementsAs(ctx, &paramsPartnerChannelTemplateCreate.ToPartnerManagedFolderPaths, false)
		resp.Diagnostics.Append(diags...)
	}
	paramsPartnerChannelTemplateCreate.ToPartnerRoutePathPattern = plan.ToPartnerRoutePathPattern.ValueString()
	paramsPartnerChannelTemplateCreate.Name = plan.Name.ValueString()
	paramsPartnerChannelTemplateCreate.Path = plan.Path.ValueString()
	paramsPartnerChannelTemplateCreate.WorkspaceId = plan.WorkspaceId.ValueInt64()

	if resp.Diagnostics.HasError() {
		return
	}

	partnerChannelTemplate, err := r.client.Create(paramsPartnerChannelTemplateCreate, files_sdk.WithContext(ctx))
	if err != nil {
		resp.Diagnostics.AddError(
			"Error Creating Files PartnerChannelTemplate",
			"Could not create partner_channel_template, unexpected error: "+err.Error(),
		)
		return
	}

	diags = r.populateResourceModel(ctx, partnerChannelTemplate, &plan)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	diags = resp.State.Set(ctx, plan)
	resp.Diagnostics.Append(diags...)
}

func (r *partnerChannelTemplateResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	var state partnerChannelTemplateResourceModel
	diags := req.State.Get(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	paramsPartnerChannelTemplateFind := files_sdk.PartnerChannelTemplateFindParams{}
	paramsPartnerChannelTemplateFind.Id = state.Id.ValueInt64()

	partnerChannelTemplate, err := r.client.Find(paramsPartnerChannelTemplateFind, files_sdk.WithContext(ctx))
	if err != nil {
		if files_sdk.IsNotExist(err) {
			resp.State.RemoveResource(ctx)
			return
		}

		resp.Diagnostics.AddError(
			"Error Reading Files PartnerChannelTemplate",
			"Could not read partner_channel_template id "+fmt.Sprint(state.Id.ValueInt64())+": "+err.Error(),
		)
		return
	}

	diags = r.populateResourceModel(ctx, partnerChannelTemplate, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	diags = resp.State.Set(ctx, state)
	resp.Diagnostics.Append(diags...)
}

func (r *partnerChannelTemplateResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	var plan partnerChannelTemplateResourceModel
	diags := req.Plan.Get(ctx, &plan)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}
	var config partnerChannelTemplateResourceModel
	diags = req.Config.Get(ctx, &config)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	paramsPartnerChannelTemplateUpdate := map[string]interface{}{}
	if !plan.Id.IsNull() && !plan.Id.IsUnknown() {
		paramsPartnerChannelTemplateUpdate["id"] = plan.Id.ValueInt64()
	}
	if !config.FromPartnerFolderName.IsNull() && !config.FromPartnerFolderName.IsUnknown() {
		paramsPartnerChannelTemplateUpdate["from_partner_folder_name"] = config.FromPartnerFolderName.ValueString()
	}
	if !config.FromPartnerManagedFolderPaths.IsNull() && !config.FromPartnerManagedFolderPaths.IsUnknown() {
		var updateFromPartnerManagedFolderPaths []string
		diags = config.FromPartnerManagedFolderPaths.ElementsAs(ctx, &updateFromPartnerManagedFolderPaths, false)
		resp.Diagnostics.Append(diags...)
		paramsPartnerChannelTemplateUpdate["from_partner_managed_folder_paths"] = updateFromPartnerManagedFolderPaths
	}
	if !config.FromPartnerRoutePathPattern.IsNull() && !config.FromPartnerRoutePathPattern.IsUnknown() {
		paramsPartnerChannelTemplateUpdate["from_partner_route_path_pattern"] = config.FromPartnerRoutePathPattern.ValueString()
	}
	if !config.ToPartnerFolderName.IsNull() && !config.ToPartnerFolderName.IsUnknown() {
		paramsPartnerChannelTemplateUpdate["to_partner_folder_name"] = config.ToPartnerFolderName.ValueString()
	}
	if !config.ToPartnerManagedFolderPaths.IsNull() && !config.ToPartnerManagedFolderPaths.IsUnknown() {
		var updateToPartnerManagedFolderPaths []string
		diags = config.ToPartnerManagedFolderPaths.ElementsAs(ctx, &updateToPartnerManagedFolderPaths, false)
		resp.Diagnostics.Append(diags...)
		paramsPartnerChannelTemplateUpdate["to_partner_managed_folder_paths"] = updateToPartnerManagedFolderPaths
	}
	if !config.ToPartnerRoutePathPattern.IsNull() && !config.ToPartnerRoutePathPattern.IsUnknown() {
		paramsPartnerChannelTemplateUpdate["to_partner_route_path_pattern"] = config.ToPartnerRoutePathPattern.ValueString()
	}
	if !config.Name.IsNull() && !config.Name.IsUnknown() {
		paramsPartnerChannelTemplateUpdate["name"] = config.Name.ValueString()
	}
	if !config.Path.IsNull() && !config.Path.IsUnknown() {
		paramsPartnerChannelTemplateUpdate["path"] = config.Path.ValueString()
	}

	if resp.Diagnostics.HasError() {
		return
	}

	partnerChannelTemplate, err := r.client.UpdateWithMap(paramsPartnerChannelTemplateUpdate, files_sdk.WithContext(ctx))
	if err != nil {
		resp.Diagnostics.AddError(
			"Error Updating Files PartnerChannelTemplate",
			"Could not update partner_channel_template, unexpected error: "+err.Error(),
		)
		return
	}

	diags = r.populateResourceModel(ctx, partnerChannelTemplate, &plan)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	diags = resp.State.Set(ctx, plan)
	resp.Diagnostics.Append(diags...)
}

func (r *partnerChannelTemplateResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	var state partnerChannelTemplateResourceModel
	diags := req.State.Get(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	paramsPartnerChannelTemplateDelete := files_sdk.PartnerChannelTemplateDeleteParams{}
	paramsPartnerChannelTemplateDelete.Id = state.Id.ValueInt64()

	err := r.client.Delete(paramsPartnerChannelTemplateDelete, files_sdk.WithContext(ctx))
	if err != nil && !files_sdk.IsNotExist(err) {
		resp.Diagnostics.AddError(
			"Error Deleting Files PartnerChannelTemplate",
			"Could not delete partner_channel_template id "+fmt.Sprint(state.Id.ValueInt64())+": "+err.Error(),
		)
	}
}

func (r *partnerChannelTemplateResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
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

func (r *partnerChannelTemplateResource) populateResourceModel(ctx context.Context, partnerChannelTemplate files_sdk.PartnerChannelTemplate, state *partnerChannelTemplateResourceModel) (diags diag.Diagnostics) {
	var propDiags diag.Diagnostics

	state.Id = types.Int64Value(partnerChannelTemplate.Id)
	state.WorkspaceId = types.Int64Value(partnerChannelTemplate.WorkspaceId)
	state.Name = types.StringValue(partnerChannelTemplate.Name)
	state.Path = types.StringValue(partnerChannelTemplate.Path)
	state.ToPartnerFolderName = types.StringValue(partnerChannelTemplate.ToPartnerFolderName)
	state.FromPartnerFolderName = types.StringValue(partnerChannelTemplate.FromPartnerFolderName)
	state.FromPartnerRoutePathPattern = types.StringValue(partnerChannelTemplate.FromPartnerRoutePathPattern)
	state.ToPartnerRoutePathPattern = types.StringValue(partnerChannelTemplate.ToPartnerRoutePathPattern)
	state.ToPartnerManagedFolderPaths, propDiags = types.ListValueFrom(ctx, types.StringType, partnerChannelTemplate.ToPartnerManagedFolderPaths)
	diags.Append(propDiags...)
	state.FromPartnerManagedFolderPaths, propDiags = types.ListValueFrom(ctx, types.StringType, partnerChannelTemplate.FromPartnerManagedFolderPaths)
	diags.Append(propDiags...)
	state.EffectiveToPartnerFolderName = types.StringValue(partnerChannelTemplate.EffectiveToPartnerFolderName)
	state.EffectiveFromPartnerFolderName = types.StringValue(partnerChannelTemplate.EffectiveFromPartnerFolderName)

	return
}
