package provider

import (
	"context"
	"fmt"

	files_sdk "github.com/Files-com/files-sdk-go/v3"
	sso_strategy "github.com/Files-com/files-sdk-go/v3/ssostrategy"
	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/datasource/schema"
	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

var (
	_ datasource.DataSource              = &ssoStrategyDataSource{}
	_ datasource.DataSourceWithConfigure = &ssoStrategyDataSource{}
)

func NewSsoStrategyDataSource() datasource.DataSource {
	return &ssoStrategyDataSource{}
}

type ssoStrategyDataSource struct {
	client *sso_strategy.Client
}

type ssoStrategyDataSourceModel struct {
	Id                             types.Int64  `tfsdk:"id"`
	Protocol                       types.String `tfsdk:"protocol"`
	Provider_                      types.String `tfsdk:"provider_"`
	Label                          types.String `tfsdk:"label"`
	LogoUrl                        types.String `tfsdk:"logo_url"`
	UserCount                      types.Int64  `tfsdk:"user_count"`
	SamlProviderCertFingerprint    types.String `tfsdk:"saml_provider_cert_fingerprint"`
	SamlProviderIssuerUrl          types.String `tfsdk:"saml_provider_issuer_url"`
	SamlProviderMetadataContent    types.String `tfsdk:"saml_provider_metadata_content"`
	SamlProviderMetadataUrl        types.String `tfsdk:"saml_provider_metadata_url"`
	SamlProviderSloTargetUrl       types.String `tfsdk:"saml_provider_slo_target_url"`
	SamlProviderSsoTargetUrl       types.String `tfsdk:"saml_provider_sso_target_url"`
	ScimAuthenticationMethod       types.String `tfsdk:"scim_authentication_method"`
	ScimUsername                   types.String `tfsdk:"scim_username"`
	ScimOauthAccessToken           types.String `tfsdk:"scim_oauth_access_token"`
	ScimOauthAccessTokenExpiresAt  types.String `tfsdk:"scim_oauth_access_token_expires_at"`
	Subdomain                      types.String `tfsdk:"subdomain"`
	ProvisionUsers                 types.Bool   `tfsdk:"provision_users"`
	ProvisionGroups                types.Bool   `tfsdk:"provision_groups"`
	DeprovisionUsers               types.Bool   `tfsdk:"deprovision_users"`
	DeprovisionGroups              types.Bool   `tfsdk:"deprovision_groups"`
	DeprovisionBehavior            types.String `tfsdk:"deprovision_behavior"`
	ProvisionGroupDefault          types.String `tfsdk:"provision_group_default"`
	ProvisionGroupExclusion        types.String `tfsdk:"provision_group_exclusion"`
	ProvisionGroupInclusion        types.String `tfsdk:"provision_group_inclusion"`
	ProvisionGroupRequired         types.String `tfsdk:"provision_group_required"`
	ProvisionEmailSignupGroups     types.String `tfsdk:"provision_email_signup_groups"`
	ProvisionSiteAdminGroups       types.String `tfsdk:"provision_site_admin_groups"`
	ProvisionGroupAdminGroups      types.String `tfsdk:"provision_group_admin_groups"`
	ProvisionAttachmentsPermission types.Bool   `tfsdk:"provision_attachments_permission"`
	ProvisionDavPermission         types.Bool   `tfsdk:"provision_dav_permission"`
	ProvisionFtpPermission         types.Bool   `tfsdk:"provision_ftp_permission"`
	ProvisionSftpPermission        types.Bool   `tfsdk:"provision_sftp_permission"`
	ProvisionTimeZone              types.String `tfsdk:"provision_time_zone"`
	ProvisionCompany               types.String `tfsdk:"provision_company"`
	ProvisionRequire2fa            types.String `tfsdk:"provision_require_2fa"`
	ProviderIdentifier             types.String `tfsdk:"provider_identifier"`
	LdapBaseDn                     types.String `tfsdk:"ldap_base_dn"`
	LdapDomain                     types.String `tfsdk:"ldap_domain"`
	Enabled                        types.Bool   `tfsdk:"enabled"`
	LdapHost                       types.String `tfsdk:"ldap_host"`
	LdapHost2                      types.String `tfsdk:"ldap_host_2"`
	LdapHost3                      types.String `tfsdk:"ldap_host_3"`
	LdapPort                       types.Int64  `tfsdk:"ldap_port"`
	LdapSecure                     types.Bool   `tfsdk:"ldap_secure"`
	LdapUsername                   types.String `tfsdk:"ldap_username"`
	LdapUsernameField              types.String `tfsdk:"ldap_username_field"`
}

func (r *ssoStrategyDataSource) Configure(_ context.Context, req datasource.ConfigureRequest, resp *datasource.ConfigureResponse) {
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

	r.client = &sso_strategy.Client{Config: sdk_config}
}

func (r *ssoStrategyDataSource) Metadata(_ context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_sso_strategy"
}

func (r *ssoStrategyDataSource) Schema(_ context.Context, _ datasource.SchemaRequest, resp *datasource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Description: "An SSOStrategy is a way for users to sign in via another identity provider, such as Okta or Auth0.\n\n\n\nIt is rare that you will need to use API endpoints for managing these, and we recommend instead managing these via the web interface.\n\nNevertheless, we share the API documentation here.",
		Attributes: map[string]schema.Attribute{
			"id": schema.Int64Attribute{
				Description: "ID",
				Required:    true,
			},
			"protocol": schema.StringAttribute{
				Description: "SSO Protocol",
				Computed:    true,
			},
			"provider_": schema.StringAttribute{
				Description: "Provider name",
				Computed:    true,
			},
			"label": schema.StringAttribute{
				Description: "Custom label for the SSO provider on the login page.",
				Computed:    true,
			},
			"logo_url": schema.StringAttribute{
				Description: "URL holding a custom logo for the SSO provider on the login page.",
				Computed:    true,
			},
			"user_count": schema.Int64Attribute{
				Description: "Count of users with this SSO Strategy",
				Computed:    true,
			},
			"saml_provider_cert_fingerprint": schema.StringAttribute{
				Description: "Identity provider sha256 cert fingerprint if saml_provider_metadata_url is not available.",
				Computed:    true,
			},
			"saml_provider_issuer_url": schema.StringAttribute{
				Description: "Identity provider issuer url",
				Computed:    true,
			},
			"saml_provider_metadata_content": schema.StringAttribute{
				Description: "Custom identity provider metadata",
				Computed:    true,
			},
			"saml_provider_metadata_url": schema.StringAttribute{
				Description: "Metadata URL for the SAML identity provider",
				Computed:    true,
			},
			"saml_provider_slo_target_url": schema.StringAttribute{
				Description: "Identity provider SLO endpoint",
				Computed:    true,
			},
			"saml_provider_sso_target_url": schema.StringAttribute{
				Description: "Identity provider SSO endpoint if saml_provider_metadata_url is not available.",
				Computed:    true,
			},
			"scim_authentication_method": schema.StringAttribute{
				Description: "SCIM authentication type.",
				Computed:    true,
			},
			"scim_username": schema.StringAttribute{
				Description: "SCIM username.",
				Computed:    true,
			},
			"scim_oauth_access_token": schema.StringAttribute{
				Description: "SCIM OAuth Access Token.",
				Computed:    true,
			},
			"scim_oauth_access_token_expires_at": schema.StringAttribute{
				Description: "SCIM OAuth Access Token Expiration Time.",
				Computed:    true,
			},
			"subdomain": schema.StringAttribute{
				Description: "Subdomain",
				Computed:    true,
			},
			"provision_users": schema.BoolAttribute{
				Description: "Auto-provision users?",
				Computed:    true,
			},
			"provision_groups": schema.BoolAttribute{
				Description: "Auto-provision group membership based on group memberships on the SSO side?",
				Computed:    true,
			},
			"deprovision_users": schema.BoolAttribute{
				Description: "Auto-deprovision users?",
				Computed:    true,
			},
			"deprovision_groups": schema.BoolAttribute{
				Description: "Auto-deprovision group membership based on group memberships on the SSO side?",
				Computed:    true,
			},
			"deprovision_behavior": schema.StringAttribute{
				Description: "Method used for deprovisioning users.",
				Computed:    true,
			},
			"provision_group_default": schema.StringAttribute{
				Description: "Comma-separated list of group names for groups to automatically add all auto-provisioned users to.",
				Computed:    true,
			},
			"provision_group_exclusion": schema.StringAttribute{
				Description: "Comma-separated list of group names for groups (with optional wildcards) that will be excluded from auto-provisioning.",
				Computed:    true,
			},
			"provision_group_inclusion": schema.StringAttribute{
				Description: "Comma-separated list of group names for groups (with optional wildcards) that will be auto-provisioned.",
				Computed:    true,
			},
			"provision_group_required": schema.StringAttribute{
				Description: "Comma or newline separated list of group names (with optional wildcards) to require membership for user provisioning.",
				Computed:    true,
			},
			"provision_email_signup_groups": schema.StringAttribute{
				Description: "Comma-separated list of group names whose members will be created with email_signup authentication.",
				Computed:    true,
			},
			"provision_site_admin_groups": schema.StringAttribute{
				Description: "Comma-separated list of group names whose members will be created as Site Admins.",
				Computed:    true,
			},
			"provision_group_admin_groups": schema.StringAttribute{
				Description: "Comma-separated list of group names whose members will be provisioned as Group Admins.",
				Computed:    true,
			},
			"provision_attachments_permission": schema.BoolAttribute{
				Computed: true,
			},
			"provision_dav_permission": schema.BoolAttribute{
				Description: "Auto-provisioned users get WebDAV permission?",
				Computed:    true,
			},
			"provision_ftp_permission": schema.BoolAttribute{
				Description: "Auto-provisioned users get FTP permission?",
				Computed:    true,
			},
			"provision_sftp_permission": schema.BoolAttribute{
				Description: "Auto-provisioned users get SFTP permission?",
				Computed:    true,
			},
			"provision_time_zone": schema.StringAttribute{
				Description: "Default time zone for auto provisioned users.",
				Computed:    true,
			},
			"provision_company": schema.StringAttribute{
				Description: "Default company for auto provisioned users.",
				Computed:    true,
			},
			"provision_require_2fa": schema.StringAttribute{
				Description: "2FA required setting for auto provisioned users.",
				Computed:    true,
			},
			"provider_identifier": schema.StringAttribute{
				Description: "URL-friendly, unique identifier for Azure SAML configuration",
				Computed:    true,
			},
			"ldap_base_dn": schema.StringAttribute{
				Description: "Base DN for looking up users in LDAP server",
				Computed:    true,
			},
			"ldap_domain": schema.StringAttribute{
				Description: "Domain name that will be appended to LDAP usernames",
				Computed:    true,
			},
			"enabled": schema.BoolAttribute{
				Description: "Is strategy enabled?  This may become automatically set to `false` after a high number and duration of failures.",
				Computed:    true,
			},
			"ldap_host": schema.StringAttribute{
				Description: "LDAP host",
				Computed:    true,
			},
			"ldap_host_2": schema.StringAttribute{
				Description: "LDAP backup host",
				Computed:    true,
			},
			"ldap_host_3": schema.StringAttribute{
				Description: "LDAP backup host",
				Computed:    true,
			},
			"ldap_port": schema.Int64Attribute{
				Description: "LDAP port",
				Computed:    true,
			},
			"ldap_secure": schema.BoolAttribute{
				Description: "Use secure LDAP?",
				Computed:    true,
			},
			"ldap_username": schema.StringAttribute{
				Description: "Username for signing in to LDAP server.",
				Computed:    true,
			},
			"ldap_username_field": schema.StringAttribute{
				Description: "LDAP username field",
				Computed:    true,
			},
		},
	}
}

func (r *ssoStrategyDataSource) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
	var data ssoStrategyDataSourceModel
	diags := req.Config.Get(ctx, &data)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	paramsSsoStrategyFind := files_sdk.SsoStrategyFindParams{}
	paramsSsoStrategyFind.Id = data.Id.ValueInt64()

	ssoStrategy, err := r.client.Find(paramsSsoStrategyFind, files_sdk.WithContext(ctx))
	if err != nil {
		resp.Diagnostics.AddError(
			"Error Reading Files SsoStrategy",
			"Could not read sso_strategy id "+fmt.Sprint(data.Id.ValueInt64())+": "+err.Error(),
		)
		return
	}

	diags = r.populateDataSourceModel(ctx, ssoStrategy, &data)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	diags = resp.State.Set(ctx, data)
	resp.Diagnostics.Append(diags...)
}

func (r *ssoStrategyDataSource) populateDataSourceModel(ctx context.Context, ssoStrategy files_sdk.SsoStrategy, state *ssoStrategyDataSourceModel) (diags diag.Diagnostics) {
	state.Protocol = types.StringValue(ssoStrategy.Protocol)
	state.Provider_ = types.StringValue(ssoStrategy.Provider)
	state.Label = types.StringValue(ssoStrategy.Label)
	state.LogoUrl = types.StringValue(ssoStrategy.LogoUrl)
	state.Id = types.Int64Value(ssoStrategy.Id)
	state.UserCount = types.Int64Value(ssoStrategy.UserCount)
	state.SamlProviderCertFingerprint = types.StringValue(ssoStrategy.SamlProviderCertFingerprint)
	state.SamlProviderIssuerUrl = types.StringValue(ssoStrategy.SamlProviderIssuerUrl)
	state.SamlProviderMetadataContent = types.StringValue(ssoStrategy.SamlProviderMetadataContent)
	state.SamlProviderMetadataUrl = types.StringValue(ssoStrategy.SamlProviderMetadataUrl)
	state.SamlProviderSloTargetUrl = types.StringValue(ssoStrategy.SamlProviderSloTargetUrl)
	state.SamlProviderSsoTargetUrl = types.StringValue(ssoStrategy.SamlProviderSsoTargetUrl)
	state.ScimAuthenticationMethod = types.StringValue(ssoStrategy.ScimAuthenticationMethod)
	state.ScimUsername = types.StringValue(ssoStrategy.ScimUsername)
	state.ScimOauthAccessToken = types.StringValue(ssoStrategy.ScimOauthAccessToken)
	state.ScimOauthAccessTokenExpiresAt = types.StringValue(ssoStrategy.ScimOauthAccessTokenExpiresAt)
	state.Subdomain = types.StringValue(ssoStrategy.Subdomain)
	state.ProvisionUsers = types.BoolPointerValue(ssoStrategy.ProvisionUsers)
	state.ProvisionGroups = types.BoolPointerValue(ssoStrategy.ProvisionGroups)
	state.DeprovisionUsers = types.BoolPointerValue(ssoStrategy.DeprovisionUsers)
	state.DeprovisionGroups = types.BoolPointerValue(ssoStrategy.DeprovisionGroups)
	state.DeprovisionBehavior = types.StringValue(ssoStrategy.DeprovisionBehavior)
	state.ProvisionGroupDefault = types.StringValue(ssoStrategy.ProvisionGroupDefault)
	state.ProvisionGroupExclusion = types.StringValue(ssoStrategy.ProvisionGroupExclusion)
	state.ProvisionGroupInclusion = types.StringValue(ssoStrategy.ProvisionGroupInclusion)
	state.ProvisionGroupRequired = types.StringValue(ssoStrategy.ProvisionGroupRequired)
	state.ProvisionEmailSignupGroups = types.StringValue(ssoStrategy.ProvisionEmailSignupGroups)
	state.ProvisionSiteAdminGroups = types.StringValue(ssoStrategy.ProvisionSiteAdminGroups)
	state.ProvisionGroupAdminGroups = types.StringValue(ssoStrategy.ProvisionGroupAdminGroups)
	state.ProvisionAttachmentsPermission = types.BoolPointerValue(ssoStrategy.ProvisionAttachmentsPermission)
	state.ProvisionDavPermission = types.BoolPointerValue(ssoStrategy.ProvisionDavPermission)
	state.ProvisionFtpPermission = types.BoolPointerValue(ssoStrategy.ProvisionFtpPermission)
	state.ProvisionSftpPermission = types.BoolPointerValue(ssoStrategy.ProvisionSftpPermission)
	state.ProvisionTimeZone = types.StringValue(ssoStrategy.ProvisionTimeZone)
	state.ProvisionCompany = types.StringValue(ssoStrategy.ProvisionCompany)
	state.ProvisionRequire2fa = types.StringValue(ssoStrategy.ProvisionRequire2fa)
	state.ProviderIdentifier = types.StringValue(ssoStrategy.ProviderIdentifier)
	state.LdapBaseDn = types.StringValue(ssoStrategy.LdapBaseDn)
	state.LdapDomain = types.StringValue(ssoStrategy.LdapDomain)
	state.Enabled = types.BoolPointerValue(ssoStrategy.Enabled)
	state.LdapHost = types.StringValue(ssoStrategy.LdapHost)
	state.LdapHost2 = types.StringValue(ssoStrategy.LdapHost2)
	state.LdapHost3 = types.StringValue(ssoStrategy.LdapHost3)
	state.LdapPort = types.Int64Value(ssoStrategy.LdapPort)
	state.LdapSecure = types.BoolPointerValue(ssoStrategy.LdapSecure)
	state.LdapUsername = types.StringValue(ssoStrategy.LdapUsername)
	state.LdapUsernameField = types.StringValue(ssoStrategy.LdapUsernameField)

	return
}
