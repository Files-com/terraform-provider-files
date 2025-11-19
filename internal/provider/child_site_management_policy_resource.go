package provider

import (
	"context"
	"fmt"
	"strconv"
	"strings"

	files_sdk "github.com/Files-com/files-sdk-go/v3"
	child_site_management_policy "github.com/Files-com/files-sdk-go/v3/childsitemanagementpolicy"
	"github.com/Files-com/terraform-provider-files/lib"
	"github.com/hashicorp/terraform-plugin-framework-validators/stringvalidator"
	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/dynamicplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/int64planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/listplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/schema/validator"
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
	PolicyType          types.String  `tfsdk:"policy_type"`
	Name                types.String  `tfsdk:"name"`
	Description         types.String  `tfsdk:"description"`
	Value               types.Dynamic `tfsdk:"value"`
	SkipChildSiteIds    types.List    `tfsdk:"skip_child_site_ids"`
	Id                  types.Int64   `tfsdk:"id"`
	AppliedChildSiteIds types.List    `tfsdk:"applied_child_site_ids"`
	CreatedAt           types.String  `tfsdk:"created_at"`
	UpdatedAt           types.String  `tfsdk:"updated_at"`
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
		Description: "A Child Site Management Policy is a centralized policy defined by a parent site to enforce consistent configurations across child sites. These policies allow parent sites to maintain control over specific aspects of their child sites' functionality and appearance.\n\n\n\nPolicies can be applied to all child sites, or specific sites can be exempted from policy management by adding their site ID to the `skip_child_site_ids` parameter.\n\n\n\nThe `value` field contains the policy configuration data, with the format varying based on the policy type. When a policy is active, its managed configurations are automatically enforced on applicable child sites, and attribute modifications are not permitted.",
		Attributes: map[string]schema.Attribute{
			"policy_type": schema.StringAttribute{
				Description: "Type of policy.  Valid values: `settings`.",
				Required:    true,
				Validators: []validator.String{
					stringvalidator.OneOf("settings"),
				},
			},
			"name": schema.StringAttribute{
				Description: "Name for this policy.",
				Computed:    true,
				Optional:    true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
			},
			"description": schema.StringAttribute{
				Description: "Description for this policy.",
				Computed:    true,
				Optional:    true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
			},
			"value": schema.DynamicAttribute{
				Description: "Policy configuration data. Attributes differ by policy type. For more information, refer to the Value Hash section of the developer documentation.",
				Computed:    true,
				Optional:    true,
				PlanModifiers: []planmodifier.Dynamic{
					dynamicplanmodifier.UseStateForUnknown(),
				},
			},
			"skip_child_site_ids": schema.ListAttribute{
				Description: "IDs of child sites that this policy has been exempted from. If `skip_child_site_ids` is empty, the policy will be applied to all child sites. To apply a policy to a child site that has been exempted, remove it from `skip_child_site_ids` or set it to an empty array (`[]`).",
				Computed:    true,
				Optional:    true,
				ElementType: types.Int64Type,
				PlanModifiers: []planmodifier.List{
					listplanmodifier.UseStateForUnknown(),
				},
			},
			"id": schema.Int64Attribute{
				Description: "Policy ID.",
				Computed:    true,
				PlanModifiers: []planmodifier.Int64{
					int64planmodifier.UseStateForUnknown(),
				},
			},
			"applied_child_site_ids": schema.ListAttribute{
				Description: "IDs of child sites that this policy has been applied to. This field is read-only.",
				Computed:    true,
				ElementType: types.Int64Type,
			},
			"created_at": schema.StringAttribute{
				Description: "When this policy was created.",
				Computed:    true,
			},
			"updated_at": schema.StringAttribute{
				Description: "When this policy was last updated.",
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
	createValue, diags := lib.DynamicToInterface(ctx, path.Root("value"), plan.Value)
	resp.Diagnostics.Append(diags...)
	paramsChildSiteManagementPolicyCreate.Value = createValue
	if !plan.SkipChildSiteIds.IsNull() && !plan.SkipChildSiteIds.IsUnknown() {
		diags = plan.SkipChildSiteIds.ElementsAs(ctx, &paramsChildSiteManagementPolicyCreate.SkipChildSiteIds, false)
		resp.Diagnostics.Append(diags...)
	}
	paramsChildSiteManagementPolicyCreate.PolicyType = paramsChildSiteManagementPolicyCreate.PolicyType.Enum()[plan.PolicyType.ValueString()]
	paramsChildSiteManagementPolicyCreate.Name = plan.Name.ValueString()
	paramsChildSiteManagementPolicyCreate.Description = plan.Description.ValueString()

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
	updateValue, diags := lib.DynamicToInterface(ctx, path.Root("value"), plan.Value)
	resp.Diagnostics.Append(diags...)
	paramsChildSiteManagementPolicyUpdate.Value = updateValue
	if !plan.SkipChildSiteIds.IsNull() && !plan.SkipChildSiteIds.IsUnknown() {
		diags = plan.SkipChildSiteIds.ElementsAs(ctx, &paramsChildSiteManagementPolicyUpdate.SkipChildSiteIds, false)
		resp.Diagnostics.Append(diags...)
	}
	paramsChildSiteManagementPolicyUpdate.PolicyType = paramsChildSiteManagementPolicyUpdate.PolicyType.Enum()[plan.PolicyType.ValueString()]
	paramsChildSiteManagementPolicyUpdate.Name = plan.Name.ValueString()
	paramsChildSiteManagementPolicyUpdate.Description = plan.Description.ValueString()

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
	state.PolicyType = types.StringValue(childSiteManagementPolicy.PolicyType)
	state.Name = types.StringValue(childSiteManagementPolicy.Name)
	state.Description = types.StringValue(childSiteManagementPolicy.Description)
	state.Value, propDiags = lib.ToDynamic(ctx, path.Root("value"), childSiteManagementPolicy.Value, state.Value.UnderlyingValue())
	diags.Append(propDiags...)
	state.AppliedChildSiteIds, propDiags = types.ListValueFrom(ctx, types.Int64Type, childSiteManagementPolicy.AppliedChildSiteIds)
	diags.Append(propDiags...)
	state.SkipChildSiteIds, propDiags = types.ListValueFrom(ctx, types.Int64Type, childSiteManagementPolicy.SkipChildSiteIds)
	diags.Append(propDiags...)
	if err := lib.TimeToStringType(ctx, path.Root("created_at"), childSiteManagementPolicy.CreatedAt, &state.CreatedAt); err != nil {
		diags.AddError(
			"Error Creating Files ChildSiteManagementPolicy",
			"Could not convert state created_at to string: "+err.Error(),
		)
	}
	if err := lib.TimeToStringType(ctx, path.Root("updated_at"), childSiteManagementPolicy.UpdatedAt, &state.UpdatedAt); err != nil {
		diags.AddError(
			"Error Creating Files ChildSiteManagementPolicy",
			"Could not convert state updated_at to string: "+err.Error(),
		)
	}

	return
}
