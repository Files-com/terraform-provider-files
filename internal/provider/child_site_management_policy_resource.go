package provider

import (
	"context"
	"fmt"
	"strconv"
	"strings"

	files_sdk "github.com/Files-com/files-sdk-go/v3"
	child_site_management_policy "github.com/Files-com/files-sdk-go/v3/childsitemanagementpolicy"
	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/int64planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/listplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

var (
	_ resource.Resource                = &childSiteManagementPolicyResource{}
	_ resource.ResourceWithConfigure   = &childSiteManagementPolicyResource{}
	_ resource.ResourceWithImportState = &childSiteManagementPolicyResource{}
)

func NewChildSiteManagementPolicyResource() resource.Resource {
	return &childSiteManagementPolicyResource{}
}

type childSiteManagementPolicyResource struct {
	client *child_site_management_policy.Client
}

type childSiteManagementPolicyResourceModel struct {
	SiteSettingName  types.String `tfsdk:"site_setting_name"`
	ManagedValue     types.String `tfsdk:"managed_value"`
	SkipChildSiteIds types.List   `tfsdk:"skip_child_site_ids"`
	Id               types.Int64  `tfsdk:"id"`
	SiteId           types.Int64  `tfsdk:"site_id"`
}

func (r *childSiteManagementPolicyResource) Configure(_ context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {
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

	r.client = &child_site_management_policy.Client{Config: sdk_config}
}

func (r *childSiteManagementPolicyResource) Metadata(_ context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_child_site_management_policy"
}

func (r *childSiteManagementPolicyResource) Schema(_ context.Context, _ resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Description: "A ChildSiteManagementPolicyEntity is a policy object defined by a parent site that enforces a specific setting and its managed value across all child sites.\n\nThis setting remains locked on child sites unless the policy explicitly exempts them.",
		Attributes: map[string]schema.Attribute{
			"site_setting_name": schema.StringAttribute{
				Description: "The name of the setting that is managed by the policy",
				Required:    true,
			},
			"managed_value": schema.StringAttribute{
				Description: "The value for the setting that will be enforced for all child sites that are not exempt",
				Required:    true,
			},
			"skip_child_site_ids": schema.ListAttribute{
				Description: "The list of child site IDs that are exempt from this policy",
				Computed:    true,
				Optional:    true,
				ElementType: types.Int64Type,
				PlanModifiers: []planmodifier.List{
					listplanmodifier.UseStateForUnknown(),
				},
			},
			"id": schema.Int64Attribute{
				Description: "ChildSiteManagementPolicy ID",
				Computed:    true,
				PlanModifiers: []planmodifier.Int64{
					int64planmodifier.UseStateForUnknown(),
				},
			},
			"site_id": schema.Int64Attribute{
				Description: "ID of the Site managing the policy",
				Computed:    true,
			},
		},
	}
}

func (r *childSiteManagementPolicyResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	var plan childSiteManagementPolicyResourceModel
	diags := req.Plan.Get(ctx, &plan)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}
	var config childSiteManagementPolicyResourceModel
	diags = req.Config.Get(ctx, &config)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	paramsChildSiteManagementPolicyCreate := files_sdk.ChildSiteManagementPolicyCreateParams{}
	paramsChildSiteManagementPolicyCreate.SiteSettingName = plan.SiteSettingName.ValueString()
	paramsChildSiteManagementPolicyCreate.ManagedValue = plan.ManagedValue.ValueString()
	if !plan.SkipChildSiteIds.IsNull() && !plan.SkipChildSiteIds.IsUnknown() {
		diags = plan.SkipChildSiteIds.ElementsAs(ctx, &paramsChildSiteManagementPolicyCreate.SkipChildSiteIds, false)
		resp.Diagnostics.Append(diags...)
	}

	if resp.Diagnostics.HasError() {
		return
	}

	childSiteManagementPolicy, err := r.client.Create(paramsChildSiteManagementPolicyCreate, files_sdk.WithContext(ctx))
	if err != nil {
		resp.Diagnostics.AddError(
			"Error Creating Files ChildSiteManagementPolicy",
			"Could not create child_site_management_policy, unexpected error: "+err.Error(),
		)
		return
	}

	diags = r.populateResourceModel(ctx, childSiteManagementPolicy, &plan)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	diags = resp.State.Set(ctx, plan)
	resp.Diagnostics.Append(diags...)
}

func (r *childSiteManagementPolicyResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	var state childSiteManagementPolicyResourceModel
	diags := req.State.Get(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	paramsChildSiteManagementPolicyFind := files_sdk.ChildSiteManagementPolicyFindParams{}
	paramsChildSiteManagementPolicyFind.Id = state.Id.ValueInt64()

	childSiteManagementPolicy, err := r.client.Find(paramsChildSiteManagementPolicyFind, files_sdk.WithContext(ctx))
	if err != nil {
		if files_sdk.IsNotExist(err) {
			resp.State.RemoveResource(ctx)
			return
		}

		resp.Diagnostics.AddError(
			"Error Reading Files ChildSiteManagementPolicy",
			"Could not read child_site_management_policy id "+fmt.Sprint(state.Id.ValueInt64())+": "+err.Error(),
		)
		return
	}

	diags = r.populateResourceModel(ctx, childSiteManagementPolicy, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	diags = resp.State.Set(ctx, state)
	resp.Diagnostics.Append(diags...)
}

func (r *childSiteManagementPolicyResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	var plan childSiteManagementPolicyResourceModel
	diags := req.Plan.Get(ctx, &plan)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}
	var config childSiteManagementPolicyResourceModel
	diags = req.Config.Get(ctx, &config)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	paramsChildSiteManagementPolicyUpdate := files_sdk.ChildSiteManagementPolicyUpdateParams{}
	paramsChildSiteManagementPolicyUpdate.Id = plan.Id.ValueInt64()
	paramsChildSiteManagementPolicyUpdate.SiteSettingName = plan.SiteSettingName.ValueString()
	paramsChildSiteManagementPolicyUpdate.ManagedValue = plan.ManagedValue.ValueString()
	if !plan.SkipChildSiteIds.IsNull() && !plan.SkipChildSiteIds.IsUnknown() {
		diags = plan.SkipChildSiteIds.ElementsAs(ctx, &paramsChildSiteManagementPolicyUpdate.SkipChildSiteIds, false)
		resp.Diagnostics.Append(diags...)
	}

	if resp.Diagnostics.HasError() {
		return
	}

	childSiteManagementPolicy, err := r.client.Update(paramsChildSiteManagementPolicyUpdate, files_sdk.WithContext(ctx))
	if err != nil {
		resp.Diagnostics.AddError(
			"Error Updating Files ChildSiteManagementPolicy",
			"Could not update child_site_management_policy, unexpected error: "+err.Error(),
		)
		return
	}

	diags = r.populateResourceModel(ctx, childSiteManagementPolicy, &plan)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	diags = resp.State.Set(ctx, plan)
	resp.Diagnostics.Append(diags...)
}

func (r *childSiteManagementPolicyResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	var state childSiteManagementPolicyResourceModel
	diags := req.State.Get(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	paramsChildSiteManagementPolicyDelete := files_sdk.ChildSiteManagementPolicyDeleteParams{}
	paramsChildSiteManagementPolicyDelete.Id = state.Id.ValueInt64()

	err := r.client.Delete(paramsChildSiteManagementPolicyDelete, files_sdk.WithContext(ctx))
	if err != nil && !files_sdk.IsNotExist(err) {
		resp.Diagnostics.AddError(
			"Error Deleting Files ChildSiteManagementPolicy",
			"Could not delete child_site_management_policy id "+fmt.Sprint(state.Id.ValueInt64())+": "+err.Error(),
		)
	}
}

func (r *childSiteManagementPolicyResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
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

func (r *childSiteManagementPolicyResource) populateResourceModel(ctx context.Context, childSiteManagementPolicy files_sdk.ChildSiteManagementPolicy, state *childSiteManagementPolicyResourceModel) (diags diag.Diagnostics) {
	var propDiags diag.Diagnostics

	state.Id = types.Int64Value(childSiteManagementPolicy.Id)
	state.SiteId = types.Int64Value(childSiteManagementPolicy.SiteId)
	state.SiteSettingName = types.StringValue(childSiteManagementPolicy.SiteSettingName)
	state.ManagedValue = types.StringValue(childSiteManagementPolicy.ManagedValue)
	state.SkipChildSiteIds, propDiags = types.ListValueFrom(ctx, types.Int64Type, childSiteManagementPolicy.SkipChildSiteIds)
	diags.Append(propDiags...)

	return
}
