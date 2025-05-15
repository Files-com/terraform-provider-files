package provider

import (
	"context"
	"fmt"

	files_sdk "github.com/Files-com/files-sdk-go/v3"
	user_lifecycle_rule "github.com/Files-com/files-sdk-go/v3/userlifecyclerule"
	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/datasource/schema"
	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

var (
	_ datasource.DataSource              = &userLifecycleRuleDataSource{}
	_ datasource.DataSourceWithConfigure = &userLifecycleRuleDataSource{}
)

func NewUserLifecycleRuleDataSource() datasource.DataSource {
	return &userLifecycleRuleDataSource{}
}

type userLifecycleRuleDataSource struct {
	client *user_lifecycle_rule.Client
}

type userLifecycleRuleDataSourceModel struct {
	Id                   types.Int64  `tfsdk:"id"`
	AuthenticationMethod types.String `tfsdk:"authentication_method"`
	InactivityDays       types.Int64  `tfsdk:"inactivity_days"`
	IncludeFolderAdmins  types.Bool   `tfsdk:"include_folder_admins"`
	IncludeSiteAdmins    types.Bool   `tfsdk:"include_site_admins"`
	Action               types.String `tfsdk:"action"`
	SiteId               types.Int64  `tfsdk:"site_id"`
}

func (r *userLifecycleRuleDataSource) Configure(_ context.Context, req datasource.ConfigureRequest, resp *datasource.ConfigureResponse) {
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

	r.client = &user_lifecycle_rule.Client{Config: sdk_config}
}

func (r *userLifecycleRuleDataSource) Metadata(_ context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_user_lifecycle_rule"
}

func (r *userLifecycleRuleDataSource) Schema(_ context.Context, _ datasource.SchemaRequest, resp *datasource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Description: "A UserLifecycleRule represents a rule that applies to users based on their inactivity and authentication method.\n\n\n\nThe rule either disable or delete users who have been inactive for a specified number of days.\n\n\n\nThe authentication_method property specifies the authentication method for the rule, which can be set to \"all\" or other specific methods.\n\n\n\nThe rule can also include or exclude site and folder admins from the action.",
		Attributes: map[string]schema.Attribute{
			"id": schema.Int64Attribute{
				Description: "User Lifecycle Rule ID",
				Required:    true,
			},
			"authentication_method": schema.StringAttribute{
				Description: "User authentication method for the rule",
				Computed:    true,
			},
			"inactivity_days": schema.Int64Attribute{
				Description: "Number of days of inactivity before the rule applies",
				Computed:    true,
			},
			"include_folder_admins": schema.BoolAttribute{
				Description: "Include folder admins in the rule",
				Computed:    true,
			},
			"include_site_admins": schema.BoolAttribute{
				Description: "Include site admins in the rule",
				Computed:    true,
			},
			"action": schema.StringAttribute{
				Description: "Action to take on inactive users (disable or delete)",
				Computed:    true,
			},
			"site_id": schema.Int64Attribute{
				Description: "Site ID",
				Computed:    true,
			},
		},
	}
}

func (r *userLifecycleRuleDataSource) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
	var data userLifecycleRuleDataSourceModel
	diags := req.Config.Get(ctx, &data)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	paramsUserLifecycleRuleFind := files_sdk.UserLifecycleRuleFindParams{}
	paramsUserLifecycleRuleFind.Id = data.Id.ValueInt64()

	userLifecycleRule, err := r.client.Find(paramsUserLifecycleRuleFind, files_sdk.WithContext(ctx))
	if err != nil {
		resp.Diagnostics.AddError(
			"Error Reading Files UserLifecycleRule",
			"Could not read user_lifecycle_rule id "+fmt.Sprint(data.Id.ValueInt64())+": "+err.Error(),
		)
		return
	}

	diags = r.populateDataSourceModel(ctx, userLifecycleRule, &data)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	diags = resp.State.Set(ctx, data)
	resp.Diagnostics.Append(diags...)
}

func (r *userLifecycleRuleDataSource) populateDataSourceModel(ctx context.Context, userLifecycleRule files_sdk.UserLifecycleRule, state *userLifecycleRuleDataSourceModel) (diags diag.Diagnostics) {
	state.Id = types.Int64Value(userLifecycleRule.Id)
	state.AuthenticationMethod = types.StringValue(userLifecycleRule.AuthenticationMethod)
	state.InactivityDays = types.Int64Value(userLifecycleRule.InactivityDays)
	state.IncludeFolderAdmins = types.BoolPointerValue(userLifecycleRule.IncludeFolderAdmins)
	state.IncludeSiteAdmins = types.BoolPointerValue(userLifecycleRule.IncludeSiteAdmins)
	state.Action = types.StringValue(userLifecycleRule.Action)
	state.SiteId = types.Int64Value(userLifecycleRule.SiteId)

	return
}
