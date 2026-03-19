package provider

import (
	"context"
	"fmt"
	"strconv"
	"strings"

	files_sdk "github.com/Files-com/files-sdk-go/v3"
	desktop_configuration_profile "github.com/Files-com/files-sdk-go/v3/desktopconfigurationprofile"
	"github.com/Files-com/terraform-provider-files/lib"
	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/boolplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/int64planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

var (
	_ resource.Resource                = &desktopConfigurationProfileResource{}
	_ resource.ResourceWithConfigure   = &desktopConfigurationProfileResource{}
	_ resource.ResourceWithImportState = &desktopConfigurationProfileResource{}
)

func NewDesktopConfigurationProfileResource() resource.Resource {
	return &desktopConfigurationProfileResource{}
}

type desktopConfigurationProfileResource struct {
	client *desktop_configuration_profile.Client
}

type desktopConfigurationProfileResourceModel struct {
	Name           types.String  `tfsdk:"name"`
	MountMappings  types.Dynamic `tfsdk:"mount_mappings"`
	WorkspaceId    types.Int64   `tfsdk:"workspace_id"`
	UseForAllUsers types.Bool    `tfsdk:"use_for_all_users"`
	Id             types.Int64   `tfsdk:"id"`
}

func (r *desktopConfigurationProfileResource) Configure(_ context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {
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

func (r *desktopConfigurationProfileResource) Metadata(_ context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_desktop_configuration_profile"
}

func (r *desktopConfigurationProfileResource) Schema(_ context.Context, _ resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Description: "A Desktop Configuration Profile centrally defines desktop mount point mappings for users in a Site or Workspace.",
		Attributes: map[string]schema.Attribute{
			"name": schema.StringAttribute{
				Description: "Profile name",
				Required:    true,
			},
			"mount_mappings": schema.DynamicAttribute{
				Description: "Mount point mappings for the desktop app. Keys are mount points (e.g. drive letters) and values are paths in Files.com that the mount points map to.",
				Required:    true,
			},
			"workspace_id": schema.Int64Attribute{
				Description: "Workspace ID",
				Computed:    true,
				Optional:    true,
				PlanModifiers: []planmodifier.Int64{
					int64planmodifier.UseStateForUnknown(),
				},
			},
			"use_for_all_users": schema.BoolAttribute{
				Description: "Whether this profile applies to all users in the Workspace by default",
				Computed:    true,
				Optional:    true,
				PlanModifiers: []planmodifier.Bool{
					boolplanmodifier.UseStateForUnknown(),
				},
			},
			"id": schema.Int64Attribute{
				Description: "Desktop Configuration Profile ID",
				Computed:    true,
				PlanModifiers: []planmodifier.Int64{
					int64planmodifier.UseStateForUnknown(),
				},
			},
		},
	}
}

func (r *desktopConfigurationProfileResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	var plan desktopConfigurationProfileResourceModel
	diags := req.Plan.Get(ctx, &plan)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}
	var config desktopConfigurationProfileResourceModel
	diags = req.Config.Get(ctx, &config)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	paramsDesktopConfigurationProfileCreate := files_sdk.DesktopConfigurationProfileCreateParams{}
	paramsDesktopConfigurationProfileCreate.Name = plan.Name.ValueString()
	createMountMappings, diags := lib.DynamicToInterface(ctx, path.Root("mount_mappings"), plan.MountMappings)
	resp.Diagnostics.Append(diags...)
	paramsDesktopConfigurationProfileCreate.MountMappings = createMountMappings
	paramsDesktopConfigurationProfileCreate.WorkspaceId = plan.WorkspaceId.ValueInt64()
	if !plan.UseForAllUsers.IsNull() && !plan.UseForAllUsers.IsUnknown() {
		paramsDesktopConfigurationProfileCreate.UseForAllUsers = plan.UseForAllUsers.ValueBoolPointer()
	}

	if resp.Diagnostics.HasError() {
		return
	}

	desktopConfigurationProfile, err := r.client.Create(paramsDesktopConfigurationProfileCreate, files_sdk.WithContext(ctx))
	if err != nil {
		resp.Diagnostics.AddError(
			"Error Creating Files DesktopConfigurationProfile",
			"Could not create desktop_configuration_profile, unexpected error: "+err.Error(),
		)
		return
	}

	diags = r.populateResourceModel(ctx, desktopConfigurationProfile, &plan)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	diags = resp.State.Set(ctx, plan)
	resp.Diagnostics.Append(diags...)
}

func (r *desktopConfigurationProfileResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	var state desktopConfigurationProfileResourceModel
	diags := req.State.Get(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	paramsDesktopConfigurationProfileFind := files_sdk.DesktopConfigurationProfileFindParams{}
	paramsDesktopConfigurationProfileFind.Id = state.Id.ValueInt64()

	desktopConfigurationProfile, err := r.client.Find(paramsDesktopConfigurationProfileFind, files_sdk.WithContext(ctx))
	if err != nil {
		if files_sdk.IsNotExist(err) {
			resp.State.RemoveResource(ctx)
			return
		}

		resp.Diagnostics.AddError(
			"Error Reading Files DesktopConfigurationProfile",
			"Could not read desktop_configuration_profile id "+fmt.Sprint(state.Id.ValueInt64())+": "+err.Error(),
		)
		return
	}

	diags = r.populateResourceModel(ctx, desktopConfigurationProfile, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	diags = resp.State.Set(ctx, state)
	resp.Diagnostics.Append(diags...)
}

func (r *desktopConfigurationProfileResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	var plan desktopConfigurationProfileResourceModel
	diags := req.Plan.Get(ctx, &plan)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}
	var config desktopConfigurationProfileResourceModel
	diags = req.Config.Get(ctx, &config)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	paramsDesktopConfigurationProfileUpdate := map[string]interface{}{}
	if !plan.Id.IsNull() && !plan.Id.IsUnknown() {
		paramsDesktopConfigurationProfileUpdate["id"] = plan.Id.ValueInt64()
	}
	if !config.Name.IsNull() && !config.Name.IsUnknown() {
		paramsDesktopConfigurationProfileUpdate["name"] = config.Name.ValueString()
	}
	if !config.WorkspaceId.IsNull() && !config.WorkspaceId.IsUnknown() {
		paramsDesktopConfigurationProfileUpdate["workspace_id"] = config.WorkspaceId.ValueInt64()
	}
	updateMountMappings, diags := lib.DynamicToInterface(ctx, path.Root("mount_mappings"), config.MountMappings)
	resp.Diagnostics.Append(diags...)
	paramsDesktopConfigurationProfileUpdate["mount_mappings"] = updateMountMappings
	if !config.UseForAllUsers.IsNull() && !config.UseForAllUsers.IsUnknown() {
		paramsDesktopConfigurationProfileUpdate["use_for_all_users"] = config.UseForAllUsers.ValueBool()
	}

	if resp.Diagnostics.HasError() {
		return
	}

	desktopConfigurationProfile, err := r.client.UpdateWithMap(paramsDesktopConfigurationProfileUpdate, files_sdk.WithContext(ctx))
	if err != nil {
		resp.Diagnostics.AddError(
			"Error Updating Files DesktopConfigurationProfile",
			"Could not update desktop_configuration_profile, unexpected error: "+err.Error(),
		)
		return
	}

	diags = r.populateResourceModel(ctx, desktopConfigurationProfile, &plan)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	diags = resp.State.Set(ctx, plan)
	resp.Diagnostics.Append(diags...)
}

func (r *desktopConfigurationProfileResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	var state desktopConfigurationProfileResourceModel
	diags := req.State.Get(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	paramsDesktopConfigurationProfileDelete := files_sdk.DesktopConfigurationProfileDeleteParams{}
	paramsDesktopConfigurationProfileDelete.Id = state.Id.ValueInt64()

	err := r.client.Delete(paramsDesktopConfigurationProfileDelete, files_sdk.WithContext(ctx))
	if err != nil && !files_sdk.IsNotExist(err) {
		resp.Diagnostics.AddError(
			"Error Deleting Files DesktopConfigurationProfile",
			"Could not delete desktop_configuration_profile id "+fmt.Sprint(state.Id.ValueInt64())+": "+err.Error(),
		)
	}
}

func (r *desktopConfigurationProfileResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
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

func (r *desktopConfigurationProfileResource) populateResourceModel(ctx context.Context, desktopConfigurationProfile files_sdk.DesktopConfigurationProfile, state *desktopConfigurationProfileResourceModel) (diags diag.Diagnostics) {
	var propDiags diag.Diagnostics

	state.Id = types.Int64Value(desktopConfigurationProfile.Id)
	state.Name = types.StringValue(desktopConfigurationProfile.Name)
	state.WorkspaceId = types.Int64Value(desktopConfigurationProfile.WorkspaceId)
	state.UseForAllUsers = types.BoolPointerValue(desktopConfigurationProfile.UseForAllUsers)
	state.MountMappings, propDiags = lib.ToDynamic(ctx, path.Root("mount_mappings"), desktopConfigurationProfile.MountMappings, state.MountMappings.UnderlyingValue())
	diags.Append(propDiags...)

	return
}
