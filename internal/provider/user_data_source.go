package provider

import (
	"context"
	"fmt"

	files_sdk "github.com/Files-com/files-sdk-go/v3"
	user "github.com/Files-com/files-sdk-go/v3/user"
	"github.com/Files-com/terraform-provider-files/lib"
	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/datasource/schema"
	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

var (
	_ datasource.DataSource              = &userDataSource{}
	_ datasource.DataSourceWithConfigure = &userDataSource{}
)

func NewUserDataSource() datasource.DataSource {
	return &userDataSource{}
}

type userDataSource struct {
	client *user.Client
}

type userDataSourceModel struct {
	Id                               types.Int64  `tfsdk:"id"`
	Username                         types.String `tfsdk:"username"`
	AdminGroupIds                    types.List   `tfsdk:"admin_group_ids"`
	AllowedIps                       types.String `tfsdk:"allowed_ips"`
	AttachmentsPermission            types.Bool   `tfsdk:"attachments_permission"`
	ApiKeysCount                     types.Int64  `tfsdk:"api_keys_count"`
	AuthenticateUntil                types.String `tfsdk:"authenticate_until"`
	AuthenticationMethod             types.String `tfsdk:"authentication_method"`
	AvatarUrl                        types.String `tfsdk:"avatar_url"`
	Billable                         types.Bool   `tfsdk:"billable"`
	BillingPermission                types.Bool   `tfsdk:"billing_permission"`
	BypassSiteAllowedIps             types.Bool   `tfsdk:"bypass_site_allowed_ips"`
	BypassUserLifecycleRules         types.Bool   `tfsdk:"bypass_user_lifecycle_rules"`
	CreatedAt                        types.String `tfsdk:"created_at"`
	DavPermission                    types.Bool   `tfsdk:"dav_permission"`
	Disabled                         types.Bool   `tfsdk:"disabled"`
	DisabledExpiredOrInactive        types.Bool   `tfsdk:"disabled_expired_or_inactive"`
	Email                            types.String `tfsdk:"email"`
	FilesystemLayout                 types.String `tfsdk:"filesystem_layout"`
	FirstLoginAt                     types.String `tfsdk:"first_login_at"`
	FtpPermission                    types.Bool   `tfsdk:"ftp_permission"`
	GroupIds                         types.String `tfsdk:"group_ids"`
	HeaderText                       types.String `tfsdk:"header_text"`
	Language                         types.String `tfsdk:"language"`
	LastLoginAt                      types.String `tfsdk:"last_login_at"`
	LastWebLoginAt                   types.String `tfsdk:"last_web_login_at"`
	LastFtpLoginAt                   types.String `tfsdk:"last_ftp_login_at"`
	LastSftpLoginAt                  types.String `tfsdk:"last_sftp_login_at"`
	LastDavLoginAt                   types.String `tfsdk:"last_dav_login_at"`
	LastDesktopLoginAt               types.String `tfsdk:"last_desktop_login_at"`
	LastRestapiLoginAt               types.String `tfsdk:"last_restapi_login_at"`
	LastApiUseAt                     types.String `tfsdk:"last_api_use_at"`
	LastActiveAt                     types.String `tfsdk:"last_active_at"`
	LastProtocolCipher               types.String `tfsdk:"last_protocol_cipher"`
	LockoutExpires                   types.String `tfsdk:"lockout_expires"`
	Name                             types.String `tfsdk:"name"`
	Company                          types.String `tfsdk:"company"`
	Notes                            types.String `tfsdk:"notes"`
	NotificationDailySendTime        types.Int64  `tfsdk:"notification_daily_send_time"`
	OfficeIntegrationEnabled         types.Bool   `tfsdk:"office_integration_enabled"`
	PartnerAdmin                     types.Bool   `tfsdk:"partner_admin"`
	PartnerId                        types.Int64  `tfsdk:"partner_id"`
	PasswordSetAt                    types.String `tfsdk:"password_set_at"`
	PasswordValidityDays             types.Int64  `tfsdk:"password_validity_days"`
	PublicKeysCount                  types.Int64  `tfsdk:"public_keys_count"`
	ReceiveAdminAlerts               types.Bool   `tfsdk:"receive_admin_alerts"`
	Require2fa                       types.String `tfsdk:"require_2fa"`
	RequireLoginBy                   types.String `tfsdk:"require_login_by"`
	Active2fa                        types.Bool   `tfsdk:"active_2fa"`
	RequirePasswordChange            types.Bool   `tfsdk:"require_password_change"`
	PasswordExpired                  types.Bool   `tfsdk:"password_expired"`
	ReadonlySiteAdmin                types.Bool   `tfsdk:"readonly_site_admin"`
	RestapiPermission                types.Bool   `tfsdk:"restapi_permission"`
	SelfManaged                      types.Bool   `tfsdk:"self_managed"`
	SftpPermission                   types.Bool   `tfsdk:"sftp_permission"`
	SiteAdmin                        types.Bool   `tfsdk:"site_admin"`
	SiteId                           types.Int64  `tfsdk:"site_id"`
	SkipWelcomeScreen                types.Bool   `tfsdk:"skip_welcome_screen"`
	SslRequired                      types.String `tfsdk:"ssl_required"`
	SsoStrategyId                    types.Int64  `tfsdk:"sso_strategy_id"`
	SubscribeToNewsletter            types.Bool   `tfsdk:"subscribe_to_newsletter"`
	ExternallyManaged                types.Bool   `tfsdk:"externally_managed"`
	Tags                             types.String `tfsdk:"tags"`
	TimeZone                         types.String `tfsdk:"time_zone"`
	TypeOf2fa                        types.String `tfsdk:"type_of_2fa"`
	TypeOf2faForDisplay              types.String `tfsdk:"type_of_2fa_for_display"`
	UserRoot                         types.String `tfsdk:"user_root"`
	UserHome                         types.String `tfsdk:"user_home"`
	DaysRemainingUntilPasswordExpire types.Int64  `tfsdk:"days_remaining_until_password_expire"`
	PasswordExpireAt                 types.String `tfsdk:"password_expire_at"`
}

func (r *userDataSource) Configure(_ context.Context, req datasource.ConfigureRequest, resp *datasource.ConfigureResponse) {
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

	r.client = &user.Client{Config: sdk_config}
}

func (r *userDataSource) Metadata(_ context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_user"
}

func (r *userDataSource) Schema(_ context.Context, _ datasource.SchemaRequest, resp *datasource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Description: "A User represents a human or system/service user with the ability to connect to Files.com via any of the available connectivity methods (unless restricted to specific protocols).\n\n\n\nUsers are associated with API Keys, SSH (SFTP) Keys, Notifications, Permissions, and Group memberships.\n\n\n\n\n\n## Authentication\n\n\n\nThe `authentication_method` property on a User determines exactly how that user can login and authenticate to their Files.com account. Files.com offers a variety of authentication methods to ensure flexibility, security, migration, and compliance.\n\n\n\nThese authentication methods can be configured during user creation and can be modified at any time by site administrators. The meanings of the available values are as follows:\n\n\n\n* `password` - Allows authentication via a password. If API Keys or SSH (SFTP) Keys are also configured, those can be used *instead* of the password. If Two Factor Authentication (2FA) methods are also configured, a valid 2nd factor is required in addition to the password.\n\n* `email_signup` - When set upon user creation, an email will be sent to the new user with a link for them to create their password. Once the user has created their password, their authentication type will change to `password`.\n\n* `sso` - Allows authentication via a linked Single Sign On provider. If API Keys or SSH (SFTP) Keys are also configured, those can be used *instead* of Single Sign On. If Two Factor Authentication (2FA) methods are also configured, a valid 2nd factor is required in addition to Single Sign On. When using this method, you must also provide a valid `sso_strategy_id` to associate the User to the appropriate SSO provider.\n\n* `password_with_imported_hash` - Works like the `password` method but allows importing a hashed password in MD5, SHA-1, or SHA-256 format. Provide the imported hash in the field `imported_password_hash`. Upon first use, the password will be converted to Files.com's internal storage format and the authentication type will change to `password`. Typically only used when migrating to Files.com from another MFT solution.\n\n* `none` - Does not allow authentication via username and password, but does allow authentication via API Key or SSH (SFTP) Key. Typically only used for service users.\n\n* `password_and_ssh_key` - Allows authentication only by providing a password and also a valid SSH (SFTP) Key in a single attempt. If API Keys are also configured, those can be used *instead* of the password and key combination. This method only works with (typically enterprise) SSH/SFTP clients capable of sending both authentication methods at once. Typically only used for service users.",
		Attributes: map[string]schema.Attribute{
			"id": schema.Int64Attribute{
				Description: "User ID",
				Required:    true,
			},
			"username": schema.StringAttribute{
				Description: "User's username",
				Computed:    true,
			},
			"admin_group_ids": schema.ListAttribute{
				Description: "List of group IDs of which this user is an administrator",
				Computed:    true,
				ElementType: types.Int64Type,
			},
			"allowed_ips": schema.StringAttribute{
				Description: "A list of allowed IPs if applicable.  Newline delimited",
				Computed:    true,
			},
			"attachments_permission": schema.BoolAttribute{
				Description: "If `true`, the user can user create Bundles (aka Share Links). Use the bundle permission instead.",
				Computed:    true,
			},
			"api_keys_count": schema.Int64Attribute{
				Description: "Number of API keys associated with this user",
				Computed:    true,
			},
			"authenticate_until": schema.StringAttribute{
				Description: "Scheduled Date/Time at which user will be deactivated",
				Computed:    true,
			},
			"authentication_method": schema.StringAttribute{
				Description: "How is this user authenticated?",
				Computed:    true,
			},
			"avatar_url": schema.StringAttribute{
				Description: "URL holding the user's avatar",
				Computed:    true,
			},
			"billable": schema.BoolAttribute{
				Description: "Is this a billable user record?",
				Computed:    true,
			},
			"billing_permission": schema.BoolAttribute{
				Description: "Allow this user to perform operations on the account, payments, and invoices?",
				Computed:    true,
			},
			"bypass_site_allowed_ips": schema.BoolAttribute{
				Description: "Allow this user to skip site-wide IP blacklists?",
				Computed:    true,
			},
			"bypass_user_lifecycle_rules": schema.BoolAttribute{
				Description: "Exempt this user from user lifecycle rules?",
				Computed:    true,
			},
			"created_at": schema.StringAttribute{
				Description: "When this user was created",
				Computed:    true,
			},
			"dav_permission": schema.BoolAttribute{
				Description: "Can the user connect with WebDAV?",
				Computed:    true,
			},
			"disabled": schema.BoolAttribute{
				Description: "Is user disabled? Disabled users cannot log in, and do not count for billing purposes. Users can be automatically disabled after an inactivity period via a Site setting or schedule to be deactivated after specific date.",
				Computed:    true,
			},
			"disabled_expired_or_inactive": schema.BoolAttribute{
				Description: "Computed property that returns true if user disabled or expired or inactive.",
				Computed:    true,
			},
			"email": schema.StringAttribute{
				Description: "User email address",
				Computed:    true,
			},
			"filesystem_layout": schema.StringAttribute{
				Description: "File system layout",
				Computed:    true,
			},
			"first_login_at": schema.StringAttribute{
				Description: "User's first login time",
				Computed:    true,
			},
			"ftp_permission": schema.BoolAttribute{
				Description: "Can the user access with FTP/FTPS?",
				Computed:    true,
			},
			"group_ids": schema.StringAttribute{
				Description: "Comma-separated list of group IDs of which this user is a member",
				Computed:    true,
			},
			"header_text": schema.StringAttribute{
				Description: "Text to display to the user in the header of the UI",
				Computed:    true,
			},
			"language": schema.StringAttribute{
				Description: "Preferred language",
				Computed:    true,
			},
			"last_login_at": schema.StringAttribute{
				Description: "User's most recent login time via any protocol",
				Computed:    true,
			},
			"last_web_login_at": schema.StringAttribute{
				Description: "User's most recent login time via web",
				Computed:    true,
			},
			"last_ftp_login_at": schema.StringAttribute{
				Description: "User's most recent login time via FTP",
				Computed:    true,
			},
			"last_sftp_login_at": schema.StringAttribute{
				Description: "User's most recent login time via SFTP",
				Computed:    true,
			},
			"last_dav_login_at": schema.StringAttribute{
				Description: "User's most recent login time via WebDAV",
				Computed:    true,
			},
			"last_desktop_login_at": schema.StringAttribute{
				Description: "User's most recent login time via Desktop app",
				Computed:    true,
			},
			"last_restapi_login_at": schema.StringAttribute{
				Description: "User's most recent login time via Rest API",
				Computed:    true,
			},
			"last_api_use_at": schema.StringAttribute{
				Description: "User's most recent API use time",
				Computed:    true,
			},
			"last_active_at": schema.StringAttribute{
				Description: "User's most recent activity time, which is the latest of most recent login, most recent API use, enablement, or creation",
				Computed:    true,
			},
			"last_protocol_cipher": schema.StringAttribute{
				Description: "The most recent protocol and cipher used",
				Computed:    true,
			},
			"lockout_expires": schema.StringAttribute{
				Description: "Time in the future that the user will no longer be locked out if applicable",
				Computed:    true,
			},
			"name": schema.StringAttribute{
				Description: "User's full name",
				Computed:    true,
			},
			"company": schema.StringAttribute{
				Description: "User's company",
				Computed:    true,
			},
			"notes": schema.StringAttribute{
				Description: "Any internal notes on the user",
				Computed:    true,
			},
			"notification_daily_send_time": schema.Int64Attribute{
				Description: "Hour of the day at which daily notifications should be sent. Can be in range 0 to 23",
				Computed:    true,
			},
			"office_integration_enabled": schema.BoolAttribute{
				Description: "Enable integration with Office for the web?",
				Computed:    true,
			},
			"partner_admin": schema.BoolAttribute{
				Description: "Is this user a Partner administrator?",
				Computed:    true,
			},
			"partner_id": schema.Int64Attribute{
				Description: "Partner ID if this user belongs to a Partner",
				Computed:    true,
			},
			"password_set_at": schema.StringAttribute{
				Description: "Last time the user's password was set",
				Computed:    true,
			},
			"password_validity_days": schema.Int64Attribute{
				Description: "Number of days to allow user to use the same password",
				Computed:    true,
			},
			"public_keys_count": schema.Int64Attribute{
				Description: "Number of public keys associated with this user",
				Computed:    true,
			},
			"receive_admin_alerts": schema.BoolAttribute{
				Description: "Should the user receive admin alerts such a certificate expiration notifications and overages?",
				Computed:    true,
			},
			"require_2fa": schema.StringAttribute{
				Description: "2FA required setting",
				Computed:    true,
			},
			"require_login_by": schema.StringAttribute{
				Description: "Require user to login by specified date otherwise it will be disabled.",
				Computed:    true,
			},
			"active_2fa": schema.BoolAttribute{
				Description: "Is 2fa active for the user?",
				Computed:    true,
			},
			"require_password_change": schema.BoolAttribute{
				Description: "Is a password change required upon next user login?",
				Computed:    true,
			},
			"password_expired": schema.BoolAttribute{
				Description: "Is user's password expired?",
				Computed:    true,
			},
			"readonly_site_admin": schema.BoolAttribute{
				Description: "Is the user an allowed to view all (non-billing) site configuration for this site?",
				Computed:    true,
			},
			"restapi_permission": schema.BoolAttribute{
				Description: "Can this user access the Web app, Desktop app, SDKs, or REST API?  (All of these tools use the API internally, so this is one unified permission set.)",
				Computed:    true,
			},
			"self_managed": schema.BoolAttribute{
				Description: "Does this user manage it's own credentials or is it a shared/bot user?",
				Computed:    true,
			},
			"sftp_permission": schema.BoolAttribute{
				Description: "Can the user access with SFTP?",
				Computed:    true,
			},
			"site_admin": schema.BoolAttribute{
				Description: "Is the user an administrator for this site?",
				Computed:    true,
			},
			"site_id": schema.Int64Attribute{
				Description: "Site ID",
				Computed:    true,
			},
			"skip_welcome_screen": schema.BoolAttribute{
				Description: "Skip Welcome page in the UI?",
				Computed:    true,
			},
			"ssl_required": schema.StringAttribute{
				Description: "SSL required setting",
				Computed:    true,
			},
			"sso_strategy_id": schema.Int64Attribute{
				Description: "SSO (Single Sign On) strategy ID for the user, if applicable.",
				Computed:    true,
			},
			"subscribe_to_newsletter": schema.BoolAttribute{
				Description: "Is the user subscribed to the newsletter?",
				Computed:    true,
			},
			"externally_managed": schema.BoolAttribute{
				Description: "Is this user managed by a SsoStrategy?",
				Computed:    true,
			},
			"tags": schema.StringAttribute{
				Description: "Comma-separated list of Tags for this user. Tags are used for other features, such as UserLifecycleRules, which can target specific tags.  Tags must only contain lowercase letters, numbers, and hyphens.",
				Computed:    true,
			},
			"time_zone": schema.StringAttribute{
				Description: "User time zone",
				Computed:    true,
			},
			"type_of_2fa": schema.StringAttribute{
				Description: "Type(s) of 2FA methods in use, for programmatic use.  Will be either `sms`, `totp`, `webauthn`, `yubi`, `email`, or multiple values sorted alphabetically and joined by an underscore.  Does not specify whether user has more than one of a given method.",
				Computed:    true,
			},
			"type_of_2fa_for_display": schema.StringAttribute{
				Description: "Type(s) of 2FA methods in use, formatted for displaying in the UI.  Unlike `type_of_2fa`, this value will make clear when a user has more than 1 of the same type of method.",
				Computed:    true,
			},
			"user_root": schema.StringAttribute{
				Description: "Root folder for FTP (and optionally SFTP if the appropriate site-wide setting is set).  Note that this is not used for API, Desktop, or Web interface.",
				Computed:    true,
			},
			"user_home": schema.StringAttribute{
				Description: "Home folder for FTP/SFTP.  Note that this is not used for API, Desktop, or Web interface.",
				Computed:    true,
			},
			"days_remaining_until_password_expire": schema.Int64Attribute{
				Description: "Number of days remaining until password expires",
				Computed:    true,
			},
			"password_expire_at": schema.StringAttribute{
				Description: "Password expiration datetime",
				Computed:    true,
			},
		},
	}
}

func (r *userDataSource) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
	var data userDataSourceModel
	diags := req.Config.Get(ctx, &data)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	paramsUserFind := files_sdk.UserFindParams{}
	paramsUserFind.Id = data.Id.ValueInt64()

	user, err := r.client.Find(paramsUserFind, files_sdk.WithContext(ctx))
	if err != nil {
		resp.Diagnostics.AddError(
			"Error Reading Files User",
			"Could not read user id "+fmt.Sprint(data.Id.ValueInt64())+": "+err.Error(),
		)
		return
	}

	diags = r.populateDataSourceModel(ctx, user, &data)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	diags = resp.State.Set(ctx, data)
	resp.Diagnostics.Append(diags...)
}

func (r *userDataSource) populateDataSourceModel(ctx context.Context, user files_sdk.User, state *userDataSourceModel) (diags diag.Diagnostics) {
	var propDiags diag.Diagnostics

	state.Id = types.Int64Value(user.Id)
	state.Username = types.StringValue(user.Username)
	state.AdminGroupIds, propDiags = types.ListValueFrom(ctx, types.Int64Type, user.AdminGroupIds)
	diags.Append(propDiags...)
	state.AllowedIps = types.StringValue(user.AllowedIps)
	state.AttachmentsPermission = types.BoolPointerValue(user.AttachmentsPermission)
	state.ApiKeysCount = types.Int64Value(user.ApiKeysCount)
	if err := lib.TimeToStringType(ctx, path.Root("authenticate_until"), user.AuthenticateUntil, &state.AuthenticateUntil); err != nil {
		diags.AddError(
			"Error Creating Files User",
			"Could not convert state authenticate_until to string: "+err.Error(),
		)
	}
	state.AuthenticationMethod = types.StringValue(user.AuthenticationMethod)
	state.AvatarUrl = types.StringValue(user.AvatarUrl)
	state.Billable = types.BoolPointerValue(user.Billable)
	state.BillingPermission = types.BoolPointerValue(user.BillingPermission)
	state.BypassSiteAllowedIps = types.BoolPointerValue(user.BypassSiteAllowedIps)
	state.BypassUserLifecycleRules = types.BoolPointerValue(user.BypassUserLifecycleRules)
	if err := lib.TimeToStringType(ctx, path.Root("created_at"), user.CreatedAt, &state.CreatedAt); err != nil {
		diags.AddError(
			"Error Creating Files User",
			"Could not convert state created_at to string: "+err.Error(),
		)
	}
	state.DavPermission = types.BoolPointerValue(user.DavPermission)
	state.Disabled = types.BoolPointerValue(user.Disabled)
	state.DisabledExpiredOrInactive = types.BoolPointerValue(user.DisabledExpiredOrInactive)
	state.Email = types.StringValue(user.Email)
	state.FilesystemLayout = types.StringValue(user.FilesystemLayout)
	if err := lib.TimeToStringType(ctx, path.Root("first_login_at"), user.FirstLoginAt, &state.FirstLoginAt); err != nil {
		diags.AddError(
			"Error Creating Files User",
			"Could not convert state first_login_at to string: "+err.Error(),
		)
	}
	state.FtpPermission = types.BoolPointerValue(user.FtpPermission)
	state.GroupIds = types.StringValue(user.GroupIds)
	state.HeaderText = types.StringValue(user.HeaderText)
	state.Language = types.StringValue(user.Language)
	if err := lib.TimeToStringType(ctx, path.Root("last_login_at"), user.LastLoginAt, &state.LastLoginAt); err != nil {
		diags.AddError(
			"Error Creating Files User",
			"Could not convert state last_login_at to string: "+err.Error(),
		)
	}
	if err := lib.TimeToStringType(ctx, path.Root("last_web_login_at"), user.LastWebLoginAt, &state.LastWebLoginAt); err != nil {
		diags.AddError(
			"Error Creating Files User",
			"Could not convert state last_web_login_at to string: "+err.Error(),
		)
	}
	if err := lib.TimeToStringType(ctx, path.Root("last_ftp_login_at"), user.LastFtpLoginAt, &state.LastFtpLoginAt); err != nil {
		diags.AddError(
			"Error Creating Files User",
			"Could not convert state last_ftp_login_at to string: "+err.Error(),
		)
	}
	if err := lib.TimeToStringType(ctx, path.Root("last_sftp_login_at"), user.LastSftpLoginAt, &state.LastSftpLoginAt); err != nil {
		diags.AddError(
			"Error Creating Files User",
			"Could not convert state last_sftp_login_at to string: "+err.Error(),
		)
	}
	if err := lib.TimeToStringType(ctx, path.Root("last_dav_login_at"), user.LastDavLoginAt, &state.LastDavLoginAt); err != nil {
		diags.AddError(
			"Error Creating Files User",
			"Could not convert state last_dav_login_at to string: "+err.Error(),
		)
	}
	if err := lib.TimeToStringType(ctx, path.Root("last_desktop_login_at"), user.LastDesktopLoginAt, &state.LastDesktopLoginAt); err != nil {
		diags.AddError(
			"Error Creating Files User",
			"Could not convert state last_desktop_login_at to string: "+err.Error(),
		)
	}
	if err := lib.TimeToStringType(ctx, path.Root("last_restapi_login_at"), user.LastRestapiLoginAt, &state.LastRestapiLoginAt); err != nil {
		diags.AddError(
			"Error Creating Files User",
			"Could not convert state last_restapi_login_at to string: "+err.Error(),
		)
	}
	if err := lib.TimeToStringType(ctx, path.Root("last_api_use_at"), user.LastApiUseAt, &state.LastApiUseAt); err != nil {
		diags.AddError(
			"Error Creating Files User",
			"Could not convert state last_api_use_at to string: "+err.Error(),
		)
	}
	if err := lib.TimeToStringType(ctx, path.Root("last_active_at"), user.LastActiveAt, &state.LastActiveAt); err != nil {
		diags.AddError(
			"Error Creating Files User",
			"Could not convert state last_active_at to string: "+err.Error(),
		)
	}
	state.LastProtocolCipher = types.StringValue(user.LastProtocolCipher)
	if err := lib.TimeToStringType(ctx, path.Root("lockout_expires"), user.LockoutExpires, &state.LockoutExpires); err != nil {
		diags.AddError(
			"Error Creating Files User",
			"Could not convert state lockout_expires to string: "+err.Error(),
		)
	}
	state.Name = types.StringValue(user.Name)
	state.Company = types.StringValue(user.Company)
	state.Notes = types.StringValue(user.Notes)
	state.NotificationDailySendTime = types.Int64Value(user.NotificationDailySendTime)
	state.OfficeIntegrationEnabled = types.BoolPointerValue(user.OfficeIntegrationEnabled)
	state.PartnerAdmin = types.BoolPointerValue(user.PartnerAdmin)
	state.PartnerId = types.Int64Value(user.PartnerId)
	if err := lib.TimeToStringType(ctx, path.Root("password_set_at"), user.PasswordSetAt, &state.PasswordSetAt); err != nil {
		diags.AddError(
			"Error Creating Files User",
			"Could not convert state password_set_at to string: "+err.Error(),
		)
	}
	state.PasswordValidityDays = types.Int64Value(user.PasswordValidityDays)
	state.PublicKeysCount = types.Int64Value(user.PublicKeysCount)
	state.ReceiveAdminAlerts = types.BoolPointerValue(user.ReceiveAdminAlerts)
	state.Require2fa = types.StringValue(user.Require2fa)
	if err := lib.TimeToStringType(ctx, path.Root("require_login_by"), user.RequireLoginBy, &state.RequireLoginBy); err != nil {
		diags.AddError(
			"Error Creating Files User",
			"Could not convert state require_login_by to string: "+err.Error(),
		)
	}
	state.Active2fa = types.BoolPointerValue(user.Active2fa)
	state.RequirePasswordChange = types.BoolPointerValue(user.RequirePasswordChange)
	state.PasswordExpired = types.BoolPointerValue(user.PasswordExpired)
	state.ReadonlySiteAdmin = types.BoolPointerValue(user.ReadonlySiteAdmin)
	state.RestapiPermission = types.BoolPointerValue(user.RestapiPermission)
	state.SelfManaged = types.BoolPointerValue(user.SelfManaged)
	state.SftpPermission = types.BoolPointerValue(user.SftpPermission)
	state.SiteAdmin = types.BoolPointerValue(user.SiteAdmin)
	state.SiteId = types.Int64Value(user.SiteId)
	state.SkipWelcomeScreen = types.BoolPointerValue(user.SkipWelcomeScreen)
	state.SslRequired = types.StringValue(user.SslRequired)
	state.SsoStrategyId = types.Int64Value(user.SsoStrategyId)
	state.SubscribeToNewsletter = types.BoolPointerValue(user.SubscribeToNewsletter)
	state.ExternallyManaged = types.BoolPointerValue(user.ExternallyManaged)
	state.Tags = types.StringValue(user.Tags)
	state.TimeZone = types.StringValue(user.TimeZone)
	state.TypeOf2fa = types.StringValue(user.TypeOf2fa)
	state.TypeOf2faForDisplay = types.StringValue(user.TypeOf2faForDisplay)
	state.UserRoot = types.StringValue(user.UserRoot)
	state.UserHome = types.StringValue(user.UserHome)
	state.DaysRemainingUntilPasswordExpire = types.Int64Value(user.DaysRemainingUntilPasswordExpire)
	if err := lib.TimeToStringType(ctx, path.Root("password_expire_at"), user.PasswordExpireAt, &state.PasswordExpireAt); err != nil {
		diags.AddError(
			"Error Creating Files User",
			"Could not convert state password_expire_at to string: "+err.Error(),
		)
	}

	return
}
