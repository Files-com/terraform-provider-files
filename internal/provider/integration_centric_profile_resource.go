package provider

import (
	"context"
	"fmt"
	"strconv"
	"strings"

	files_sdk "github.com/Files-com/files-sdk-go/v3"
	integration_centric_profile "github.com/Files-com/files-sdk-go/v3/integrationcentricprofile"
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
	_ resource.Resource                = &integrationCentricProfileResource{}
	_ resource.ResourceWithConfigure   = &integrationCentricProfileResource{}
	_ resource.ResourceWithImportState = &integrationCentricProfileResource{}
)

func NewIntegrationCentricProfileResource() resource.Resource {
	return &integrationCentricProfileResource{}
}

type integrationCentricProfileResource struct {
	client *integration_centric_profile.Client
}

type integrationCentricProfileResourceModel struct {
	Name                  types.String  `tfsdk:"name"`
	ExpectedRemoteServers types.Dynamic `tfsdk:"expected_remote_servers"`
	WorkspaceId           types.Int64   `tfsdk:"workspace_id"`
	UseForAllUsers        types.Bool    `tfsdk:"use_for_all_users"`
	Id                    types.Int64   `tfsdk:"id"`
}

func (r *integrationCentricProfileResource) Configure(_ context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {
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

func (r *integrationCentricProfileResource) Metadata(_ context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_integration_centric_profile"
}

func (r *integrationCentricProfileResource) Schema(_ context.Context, _ resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Description: "An Integration Centric Profile defines the Remote Server integrations a user is expected to add and connect during integration-centric onboarding.\n\n\n\nUse this to automate setup guidance for users who need access to multiple business systems without sending long manual instructions. Common scenarios include ongoing access to systems such as SharePoint, bridging Google, Microsoft, and Box environments after M&A activity, and migrations where users connect legacy EFSS accounts during transition work.",
		Attributes: map[string]schema.Attribute{
			"name": schema.StringAttribute{
				Description: "Profile name",
				Required:    true,
			},
			"expected_remote_servers": schema.DynamicAttribute{
				Description: "Remote Server integrations the user is expected to add and connect. Each entry requires `server_type` and may include a display `name`.",
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
				Description: "Integration Centric Profile ID",
				Computed:    true,
				PlanModifiers: []planmodifier.Int64{
					int64planmodifier.UseStateForUnknown(),
				},
			},
		},
	}
}

func (r *integrationCentricProfileResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	var plan integrationCentricProfileResourceModel
	diags := req.Plan.Get(ctx, &plan)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}
	var config integrationCentricProfileResourceModel
	diags = req.Config.Get(ctx, &config)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	paramsIntegrationCentricProfileCreate := files_sdk.IntegrationCentricProfileCreateParams{}
	paramsIntegrationCentricProfileCreate.Name = plan.Name.ValueString()
	paramsIntegrationCentricProfileCreate.ExpectedRemoteServers, diags = lib.DynamicToStringMapSlice(ctx, path.Root("expected_remote_servers"), plan.ExpectedRemoteServers)
	resp.Diagnostics.Append(diags...)
	paramsIntegrationCentricProfileCreate.WorkspaceId = plan.WorkspaceId.ValueInt64()
	if !plan.UseForAllUsers.IsNull() && !plan.UseForAllUsers.IsUnknown() {
		paramsIntegrationCentricProfileCreate.UseForAllUsers = plan.UseForAllUsers.ValueBoolPointer()
	}

	if resp.Diagnostics.HasError() {
		return
	}

	integrationCentricProfile, err := r.client.Create(paramsIntegrationCentricProfileCreate, files_sdk.WithContext(ctx))
	if err != nil {
		resp.Diagnostics.AddError(
			"Error Creating Files IntegrationCentricProfile",
			"Could not create integration_centric_profile, unexpected error: "+err.Error(),
		)
		return
	}

	diags = r.populateResourceModel(ctx, integrationCentricProfile, &plan)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	diags = resp.State.Set(ctx, plan)
	resp.Diagnostics.Append(diags...)
}

func (r *integrationCentricProfileResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	var state integrationCentricProfileResourceModel
	diags := req.State.Get(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	paramsIntegrationCentricProfileFind := files_sdk.IntegrationCentricProfileFindParams{}
	paramsIntegrationCentricProfileFind.Id = state.Id.ValueInt64()

	integrationCentricProfile, err := r.client.Find(paramsIntegrationCentricProfileFind, files_sdk.WithContext(ctx))
	if err != nil {
		if files_sdk.IsNotExist(err) {
			resp.State.RemoveResource(ctx)
			return
		}

		resp.Diagnostics.AddError(
			"Error Reading Files IntegrationCentricProfile",
			"Could not read integration_centric_profile id "+fmt.Sprint(state.Id.ValueInt64())+": "+err.Error(),
		)
		return
	}

	diags = r.populateResourceModel(ctx, integrationCentricProfile, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	diags = resp.State.Set(ctx, state)
	resp.Diagnostics.Append(diags...)
}

func (r *integrationCentricProfileResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	var plan integrationCentricProfileResourceModel
	diags := req.Plan.Get(ctx, &plan)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}
	var config integrationCentricProfileResourceModel
	diags = req.Config.Get(ctx, &config)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	paramsIntegrationCentricProfileUpdate := map[string]interface{}{}
	if !plan.Id.IsNull() && !plan.Id.IsUnknown() {
		paramsIntegrationCentricProfileUpdate["id"] = plan.Id.ValueInt64()
	}
	if !config.Name.IsNull() && !config.Name.IsUnknown() {
		paramsIntegrationCentricProfileUpdate["name"] = config.Name.ValueString()
	}
	if !config.WorkspaceId.IsNull() && !config.WorkspaceId.IsUnknown() {
		paramsIntegrationCentricProfileUpdate["workspace_id"] = config.WorkspaceId.ValueInt64()
	}
	if !config.ExpectedRemoteServers.IsNull() && !config.ExpectedRemoteServers.IsUnknown() {
		updateExpectedRemoteServers, diags := lib.DynamicToStringMapSlice(ctx, path.Root("expected_remote_servers"), config.ExpectedRemoteServers)
		resp.Diagnostics.Append(diags...)
		paramsIntegrationCentricProfileUpdate["expected_remote_servers"] = updateExpectedRemoteServers
	}
	if !config.UseForAllUsers.IsNull() && !config.UseForAllUsers.IsUnknown() {
		paramsIntegrationCentricProfileUpdate["use_for_all_users"] = config.UseForAllUsers.ValueBool()
	}

	if resp.Diagnostics.HasError() {
		return
	}

	integrationCentricProfile, err := r.client.UpdateWithMap(paramsIntegrationCentricProfileUpdate, files_sdk.WithContext(ctx))
	if err != nil {
		resp.Diagnostics.AddError(
			"Error Updating Files IntegrationCentricProfile",
			"Could not update integration_centric_profile, unexpected error: "+err.Error(),
		)
		return
	}

	diags = r.populateResourceModel(ctx, integrationCentricProfile, &plan)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	diags = resp.State.Set(ctx, plan)
	resp.Diagnostics.Append(diags...)
}

func (r *integrationCentricProfileResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	var state integrationCentricProfileResourceModel
	diags := req.State.Get(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	paramsIntegrationCentricProfileDelete := files_sdk.IntegrationCentricProfileDeleteParams{}
	paramsIntegrationCentricProfileDelete.Id = state.Id.ValueInt64()

	err := r.client.Delete(paramsIntegrationCentricProfileDelete, files_sdk.WithContext(ctx))
	if err != nil && !files_sdk.IsNotExist(err) {
		resp.Diagnostics.AddError(
			"Error Deleting Files IntegrationCentricProfile",
			"Could not delete integration_centric_profile id "+fmt.Sprint(state.Id.ValueInt64())+": "+err.Error(),
		)
	}
}

func (r *integrationCentricProfileResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
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

func (r *integrationCentricProfileResource) populateResourceModel(ctx context.Context, integrationCentricProfile files_sdk.IntegrationCentricProfile, state *integrationCentricProfileResourceModel) (diags diag.Diagnostics) {
	var propDiags diag.Diagnostics

	state.Id = types.Int64Value(integrationCentricProfile.Id)
	state.Name = types.StringValue(integrationCentricProfile.Name)
	state.WorkspaceId = types.Int64Value(integrationCentricProfile.WorkspaceId)
	state.UseForAllUsers = types.BoolPointerValue(integrationCentricProfile.UseForAllUsers)
	state.ExpectedRemoteServers, propDiags = lib.ToDynamic(ctx, path.Root("expected_remote_servers"), integrationCentricProfile.ExpectedRemoteServers, state.ExpectedRemoteServers.UnderlyingValue())
	diags.Append(propDiags...)

	return
}
