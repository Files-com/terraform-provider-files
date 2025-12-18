package provider

import (
	"context"
	"fmt"

	files_sdk "github.com/Files-com/files-sdk-go/v3"
	key_lifecycle_rule "github.com/Files-com/files-sdk-go/v3/keylifecyclerule"
	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/datasource/schema"
	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

var (
	_ datasource.DataSource              = &keyLifecycleRuleDataSource{}
	_ datasource.DataSourceWithConfigure = &keyLifecycleRuleDataSource{}
)

func NewKeyLifecycleRuleDataSource() datasource.DataSource {
	return &keyLifecycleRuleDataSource{}
}

type keyLifecycleRuleDataSource struct {
	client *key_lifecycle_rule.Client
}

type keyLifecycleRuleDataSourceModel struct {
	Id             types.Int64  `tfsdk:"id"`
	KeyType        types.String `tfsdk:"key_type"`
	InactivityDays types.Int64  `tfsdk:"inactivity_days"`
	Name           types.String `tfsdk:"name"`
}

func (r *keyLifecycleRuleDataSource) Configure(_ context.Context, req datasource.ConfigureRequest, resp *datasource.ConfigureResponse) {
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

	r.client = &key_lifecycle_rule.Client{Config: sdk_config}
}

func (r *keyLifecycleRuleDataSource) Metadata(_ context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_key_lifecycle_rule"
}

func (r *keyLifecycleRuleDataSource) Schema(_ context.Context, _ datasource.SchemaRequest, resp *datasource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Description: "A KeyLifecycleRule represents a rule that applies to GPG keys and SSH keys (also called User Public Keys) based on their inactivity.\n\n\n\nKeys that have been unused for the specified number of days will be deleted.",
		Attributes: map[string]schema.Attribute{
			"id": schema.Int64Attribute{
				Description: "Key Lifecycle Rule ID",
				Required:    true,
			},
			"key_type": schema.StringAttribute{
				Description: "Key type for which the rule will apply (gpg or ssh).",
				Computed:    true,
			},
			"inactivity_days": schema.Int64Attribute{
				Description: "Number of days of inactivity before the rule applies.",
				Computed:    true,
			},
			"name": schema.StringAttribute{
				Description: "Key Lifecycle Rule name",
				Computed:    true,
			},
		},
	}
}

func (r *keyLifecycleRuleDataSource) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
	var data keyLifecycleRuleDataSourceModel
	diags := req.Config.Get(ctx, &data)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	paramsKeyLifecycleRuleFind := files_sdk.KeyLifecycleRuleFindParams{}
	paramsKeyLifecycleRuleFind.Id = data.Id.ValueInt64()

	keyLifecycleRule, err := r.client.Find(paramsKeyLifecycleRuleFind, files_sdk.WithContext(ctx))
	if err != nil {
		resp.Diagnostics.AddError(
			"Error Reading Files KeyLifecycleRule",
			"Could not read key_lifecycle_rule id "+fmt.Sprint(data.Id.ValueInt64())+": "+err.Error(),
		)
		return
	}

	diags = r.populateDataSourceModel(ctx, keyLifecycleRule, &data)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	diags = resp.State.Set(ctx, data)
	resp.Diagnostics.Append(diags...)
}

func (r *keyLifecycleRuleDataSource) populateDataSourceModel(ctx context.Context, keyLifecycleRule files_sdk.KeyLifecycleRule, state *keyLifecycleRuleDataSourceModel) (diags diag.Diagnostics) {
	state.Id = types.Int64Value(keyLifecycleRule.Id)
	state.KeyType = types.StringValue(keyLifecycleRule.KeyType)
	state.InactivityDays = types.Int64Value(keyLifecycleRule.InactivityDays)
	state.Name = types.StringValue(keyLifecycleRule.Name)

	return
}
