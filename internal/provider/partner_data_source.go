package provider

import (
	"context"
	"fmt"

	files_sdk "github.com/Files-com/files-sdk-go/v3"
	partner "github.com/Files-com/files-sdk-go/v3/partner"
	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/datasource/schema"
	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

var (
	_ datasource.DataSource              = &partnerDataSource{}
	_ datasource.DataSourceWithConfigure = &partnerDataSource{}
)

func NewPartnerDataSource() datasource.DataSource {
	return &partnerDataSource{}
}

type partnerDataSource struct {
	client *partner.Client
}

type partnerDataSourceModel struct {
	Id                        types.Int64  `tfsdk:"id"`
	AllowBypassing2faPolicies types.Bool   `tfsdk:"allow_bypassing_2fa_policies"`
	AllowCredentialChanges    types.Bool   `tfsdk:"allow_credential_changes"`
	AllowProvidingGpgKeys     types.Bool   `tfsdk:"allow_providing_gpg_keys"`
	AllowUserCreation         types.Bool   `tfsdk:"allow_user_creation"`
	WorkspaceId               types.Int64  `tfsdk:"workspace_id"`
	Name                      types.String `tfsdk:"name"`
	Notes                     types.String `tfsdk:"notes"`
	PartnerAdminIds           types.List   `tfsdk:"partner_admin_ids"`
	RootFolder                types.String `tfsdk:"root_folder"`
	Tags                      types.String `tfsdk:"tags"`
	UserIds                   types.List   `tfsdk:"user_ids"`
}

func (r *partnerDataSource) Configure(_ context.Context, req datasource.ConfigureRequest, resp *datasource.ConfigureResponse) {
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

	r.client = &partner.Client{Config: sdk_config}
}

func (r *partnerDataSource) Metadata(_ context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_partner"
}

func (r *partnerDataSource) Schema(_ context.Context, _ datasource.SchemaRequest, resp *datasource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Description: "A Partner is a first-class entity that cleanly represents an external organization, enables delegated administration, and enforces strict boundaries.",
		Attributes: map[string]schema.Attribute{
			"id": schema.Int64Attribute{
				Description: "The unique ID of the Partner.",
				Required:    true,
			},
			"allow_bypassing_2fa_policies": schema.BoolAttribute{
				Description: "Allow Partner Admins to change Two-Factor Authentication requirements for Partner Users.",
				Computed:    true,
			},
			"allow_credential_changes": schema.BoolAttribute{
				Description: "Allow Partner Admins to change or reset credentials for users belonging to this Partner.",
				Computed:    true,
			},
			"allow_providing_gpg_keys": schema.BoolAttribute{
				Description: "Allow Partner Admins to provide GPG keys.",
				Computed:    true,
			},
			"allow_user_creation": schema.BoolAttribute{
				Description: "Allow Partner Admins to create users.",
				Computed:    true,
			},
			"workspace_id": schema.Int64Attribute{
				Description: "ID of the Workspace associated with this Partner.",
				Computed:    true,
			},
			"name": schema.StringAttribute{
				Description: "The name of the Partner.",
				Computed:    true,
			},
			"notes": schema.StringAttribute{
				Description: "Notes about this Partner.",
				Computed:    true,
			},
			"partner_admin_ids": schema.ListAttribute{
				Description: "Array of User IDs that are Partner Admins for this Partner.",
				Computed:    true,
				ElementType: types.Int64Type,
			},
			"root_folder": schema.StringAttribute{
				Description: "The root folder path for this Partner.",
				Computed:    true,
			},
			"tags": schema.StringAttribute{
				Description: "Comma-separated list of Tags for this Partner. Tags are used for other features, such as UserLifecycleRules, which can target specific tags.  Tags must only contain lowercase letters, numbers, and hyphens.",
				Computed:    true,
			},
			"user_ids": schema.ListAttribute{
				Description: "Array of User IDs that belong to this Partner.",
				Computed:    true,
				ElementType: types.Int64Type,
			},
		},
	}
}

func (r *partnerDataSource) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
	var data partnerDataSourceModel
	diags := req.Config.Get(ctx, &data)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	paramsPartnerFind := files_sdk.PartnerFindParams{}
	paramsPartnerFind.Id = data.Id.ValueInt64()

	partner, err := r.client.Find(paramsPartnerFind, files_sdk.WithContext(ctx))
	if err != nil {
		resp.Diagnostics.AddError(
			"Error Reading Files Partner",
			"Could not read partner id "+fmt.Sprint(data.Id.ValueInt64())+": "+err.Error(),
		)
		return
	}

	diags = r.populateDataSourceModel(ctx, partner, &data)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	diags = resp.State.Set(ctx, data)
	resp.Diagnostics.Append(diags...)
}

func (r *partnerDataSource) populateDataSourceModel(ctx context.Context, partner files_sdk.Partner, state *partnerDataSourceModel) (diags diag.Diagnostics) {
	var propDiags diag.Diagnostics

	state.AllowBypassing2faPolicies = types.BoolPointerValue(partner.AllowBypassing2faPolicies)
	state.AllowCredentialChanges = types.BoolPointerValue(partner.AllowCredentialChanges)
	state.AllowProvidingGpgKeys = types.BoolPointerValue(partner.AllowProvidingGpgKeys)
	state.AllowUserCreation = types.BoolPointerValue(partner.AllowUserCreation)
	state.Id = types.Int64Value(partner.Id)
	state.WorkspaceId = types.Int64Value(partner.WorkspaceId)
	state.Name = types.StringValue(partner.Name)
	state.Notes = types.StringValue(partner.Notes)
	state.PartnerAdminIds, propDiags = types.ListValueFrom(ctx, types.Int64Type, partner.PartnerAdminIds)
	diags.Append(propDiags...)
	state.RootFolder = types.StringValue(partner.RootFolder)
	state.Tags = types.StringValue(partner.Tags)
	state.UserIds, propDiags = types.ListValueFrom(ctx, types.Int64Type, partner.UserIds)
	diags.Append(propDiags...)

	return
}
