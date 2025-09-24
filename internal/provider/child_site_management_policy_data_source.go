package provider

import (
	"context"
	"fmt"

	files_sdk "github.com/Files-com/files-sdk-go/v3"
	child_site_management_policy "github.com/Files-com/files-sdk-go/v3/childsitemanagementpolicy"
	"github.com/Files-com/terraform-provider-files/lib"
	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/datasource/schema"
	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

var (
	_ datasource.DataSource              = &childSiteManagementPolicyDataSource{}
	_ datasource.DataSourceWithConfigure = &childSiteManagementPolicyDataSource{}
)

func NewChildSiteManagementPolicyDataSource() datasource.DataSource {
	return &childSiteManagementPolicyDataSource{}
}

type childSiteManagementPolicyDataSource struct {
	client *child_site_management_policy.Client
}

type childSiteManagementPolicyDataSourceModel struct {
	Id                  types.Int64   `tfsdk:"id"`
	PolicyType          types.String  `tfsdk:"policy_type"`
	Name                types.String  `tfsdk:"name"`
	Description         types.String  `tfsdk:"description"`
	Value               types.Dynamic `tfsdk:"value"`
	AppliedChildSiteIds types.List    `tfsdk:"applied_child_site_ids"`
	SkipChildSiteIds    types.List    `tfsdk:"skip_child_site_ids"`
	CreatedAt           types.String  `tfsdk:"created_at"`
	UpdatedAt           types.String  `tfsdk:"updated_at"`
}

func (r *childSiteManagementPolicyDataSource) Configure(_ context.Context, req datasource.ConfigureRequest, resp *datasource.ConfigureResponse) {
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

func (r *childSiteManagementPolicyDataSource) Metadata(_ context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_child_site_management_policy"
}

func (r *childSiteManagementPolicyDataSource) Schema(_ context.Context, _ datasource.SchemaRequest, resp *datasource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Description: "A Child Site Management Policy is a centralized policy defined by a parent site to enforce consistent configurations across child sites. These policies allow parent sites to maintain control over specific aspects of their child sites' functionality and appearance.\n\n\n\nPolicies can be applied to all child sites, or specific sites can be exempted from policy management by adding their site ID to the `skip_child_site_ids` parameter.\n\n\n\nThe `value` field contains the policy configuration data, with the format varying based on the policy type. When a policy is active, its managed configurations are automatically enforced on applicable child sites, and attribute modifications are not permitted.",
		Attributes: map[string]schema.Attribute{
			"id": schema.Int64Attribute{
				Description: "Policy ID.",
				Required:    true,
			},
			"policy_type": schema.StringAttribute{
				Description: "Type of policy.  Valid values: `settings`.",
				Computed:    true,
			},
			"name": schema.StringAttribute{
				Description: "Name for this policy.",
				Computed:    true,
			},
			"description": schema.StringAttribute{
				Description: "Description for this policy.",
				Computed:    true,
			},
			"value": schema.DynamicAttribute{
				Description: "Policy configuration data. Attributes differ by policy type. For more information, refer to the Value Hash section of the developer documentation.",
				Computed:    true,
			},
			"applied_child_site_ids": schema.ListAttribute{
				Description: "IDs of child sites that this policy has been applied to. This field is read-only.",
				Computed:    true,
				ElementType: types.Int64Type,
			},
			"skip_child_site_ids": schema.ListAttribute{
				Description: "IDs of child sites that this policy has been exempted from. If `skip_child_site_ids` is empty, the policy will be applied to all child sites. To apply a policy to a child site that has been exempted, remove it from `skip_child_site_ids` or set it to an empty array (`[]`).",
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

func (r *childSiteManagementPolicyDataSource) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
	var data childSiteManagementPolicyDataSourceModel
	diags := req.Config.Get(ctx, &data)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	paramsChildSiteManagementPolicyFind := files_sdk.ChildSiteManagementPolicyFindParams{}
	paramsChildSiteManagementPolicyFind.Id = data.Id.ValueInt64()

	childSiteManagementPolicy, err := r.client.Find(paramsChildSiteManagementPolicyFind, files_sdk.WithContext(ctx))
	if err != nil {
		resp.Diagnostics.AddError(
			"Error Reading Files ChildSiteManagementPolicy",
			"Could not read child_site_management_policy id "+fmt.Sprint(data.Id.ValueInt64())+": "+err.Error(),
		)
		return
	}

	diags = r.populateDataSourceModel(ctx, childSiteManagementPolicy, &data)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	diags = resp.State.Set(ctx, data)
	resp.Diagnostics.Append(diags...)
}

func (r *childSiteManagementPolicyDataSource) populateDataSourceModel(ctx context.Context, childSiteManagementPolicy files_sdk.ChildSiteManagementPolicy, state *childSiteManagementPolicyDataSourceModel) (diags diag.Diagnostics) {
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
