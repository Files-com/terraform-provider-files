package provider

import (
	"context"
	"fmt"

	files_sdk "github.com/Files-com/files-sdk-go/v3"
	child_site_management_policy "github.com/Files-com/files-sdk-go/v3/childsitemanagementpolicy"
	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/datasource/schema"
	"github.com/hashicorp/terraform-plugin-framework/diag"
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
	Id               types.Int64  `tfsdk:"id"`
	SiteId           types.Int64  `tfsdk:"site_id"`
	SiteSettingName  types.String `tfsdk:"site_setting_name"`
	ManagedValue     types.String `tfsdk:"managed_value"`
	SkipChildSiteIds types.List   `tfsdk:"skip_child_site_ids"`
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
		Description: "A ChildSiteManagementPolicyEntity is a policy object defined by a parent site that enforces a specific setting and its managed value across all child sites.\n\nThis setting remains locked on child sites unless the policy explicitly exempts them.",
		Attributes: map[string]schema.Attribute{
			"id": schema.Int64Attribute{
				Description: "ChildSiteManagementPolicy ID",
				Required:    true,
			},
			"site_id": schema.Int64Attribute{
				Description: "ID of the Site managing the policy",
				Computed:    true,
			},
			"site_setting_name": schema.StringAttribute{
				Description: "The name of the setting that is managed by the policy",
				Computed:    true,
			},
			"managed_value": schema.StringAttribute{
				Description: "The value for the setting that will be enforced for all child sites that are not exempt",
				Computed:    true,
			},
			"skip_child_site_ids": schema.ListAttribute{
				Description: "The list of child site IDs that are exempt from this policy",
				Computed:    true,
				ElementType: types.Int64Type,
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
	state.SiteId = types.Int64Value(childSiteManagementPolicy.SiteId)
	state.SiteSettingName = types.StringValue(childSiteManagementPolicy.SiteSettingName)
	state.ManagedValue = types.StringValue(childSiteManagementPolicy.ManagedValue)
	state.SkipChildSiteIds, propDiags = types.ListValueFrom(ctx, types.Int64Type, childSiteManagementPolicy.SkipChildSiteIds)
	diags.Append(propDiags...)

	return
}
