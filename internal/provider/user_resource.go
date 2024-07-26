package provider

import (
	"context"
	"fmt"
	"strconv"
	"strings"
	"time"

	files_sdk "github.com/Files-com/files-sdk-go/v3"
	user "github.com/Files-com/files-sdk-go/v3/user"
	"github.com/Files-com/terraform-provider-files/lib"
	"github.com/hashicorp/terraform-plugin-framework-validators/stringvalidator"
	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/boolplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/int64planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/schema/validator"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

var (
	_ resource.Resource                = &userResource{}
	_ resource.ResourceWithConfigure   = &userResource{}
	_ resource.ResourceWithImportState = &userResource{}
)

func NewUserResource() resource.Resource {
	return &userResource{}
}

type userResource struct {
	client *user.Client
}

type userResourceModel struct {
	Username                         types.String `tfsdk:"username"`
	AllowedIps                       types.String `tfsdk:"allowed_ips"`
	AttachmentsPermission            types.Bool   `tfsdk:"attachments_permission"`
	AuthenticateUntil                types.String `tfsdk:"authenticate_until"`
	AuthenticationMethod             types.String `tfsdk:"authentication_method"`
	BillingPermission                types.Bool   `tfsdk:"billing_permission"`
	BypassSiteAllowedIps             types.Bool   `tfsdk:"bypass_site_allowed_ips"`
	BypassInactiveDisable            types.Bool   `tfsdk:"bypass_inactive_disable"`
	DavPermission                    types.Bool   `tfsdk:"dav_permission"`
	Disabled                         types.Bool   `tfsdk:"disabled"`
	Email                            types.String `tfsdk:"email"`
	FtpPermission                    types.Bool   `tfsdk:"ftp_permission"`
	GroupIds                         types.String `tfsdk:"group_ids"`
	HeaderText                       types.String `tfsdk:"header_text"`
	Language                         types.String `tfsdk:"language"`
	Name                             types.String `tfsdk:"name"`
	Company                          types.String `tfsdk:"company"`
	Notes                            types.String `tfsdk:"notes"`
	NotificationDailySendTime        types.Int64  `tfsdk:"notification_daily_send_time"`
	OfficeIntegrationEnabled         types.Bool   `tfsdk:"office_integration_enabled"`
	PasswordValidityDays             types.Int64  `tfsdk:"password_validity_days"`
	ReceiveAdminAlerts               types.Bool   `tfsdk:"receive_admin_alerts"`
	Require2fa                       types.String `tfsdk:"require_2fa"`
	RequireLoginBy                   types.String `tfsdk:"require_login_by"`
	RequirePasswordChange            types.Bool   `tfsdk:"require_password_change"`
	RestapiPermission                types.Bool   `tfsdk:"restapi_permission"`
	SelfManaged                      types.Bool   `tfsdk:"self_managed"`
	SftpPermission                   types.Bool   `tfsdk:"sftp_permission"`
	SiteAdmin                        types.Bool   `tfsdk:"site_admin"`
	SkipWelcomeScreen                types.Bool   `tfsdk:"skip_welcome_screen"`
	SslRequired                      types.String `tfsdk:"ssl_required"`
	SsoStrategyId                    types.Int64  `tfsdk:"sso_strategy_id"`
	SubscribeToNewsletter            types.Bool   `tfsdk:"subscribe_to_newsletter"`
	TimeZone                         types.String `tfsdk:"time_zone"`
	UserRoot                         types.String `tfsdk:"user_root"`
	AvatarDelete                     types.Bool   `tfsdk:"avatar_delete"`
	ChangePassword                   types.String `tfsdk:"change_password"`
	ChangePasswordConfirmation       types.String `tfsdk:"change_password_confirmation"`
	GrantPermission                  types.String `tfsdk:"grant_permission"`
	GroupId                          types.Int64  `tfsdk:"group_id"`
	ImportedPasswordHash             types.String `tfsdk:"imported_password_hash"`
	Password                         types.String `tfsdk:"password"`
	PasswordConfirmation             types.String `tfsdk:"password_confirmation"`
	AnnouncementsRead                types.Bool   `tfsdk:"announcements_read"`
	Id                               types.Int64  `tfsdk:"id"`
	AdminGroupIds                    types.List   `tfsdk:"admin_group_ids"`
	ApiKeysCount                     types.Int64  `tfsdk:"api_keys_count"`
	AvatarUrl                        types.String `tfsdk:"avatar_url"`
	CreatedAt                        types.String `tfsdk:"created_at"`
	DisabledExpiredOrInactive        types.Bool   `tfsdk:"disabled_expired_or_inactive"`
	FirstLoginAt                     types.String `tfsdk:"first_login_at"`
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
	PasswordSetAt                    types.String `tfsdk:"password_set_at"`
	PublicKeysCount                  types.Int64  `tfsdk:"public_keys_count"`
	Active2fa                        types.Bool   `tfsdk:"active_2fa"`
	PasswordExpired                  types.Bool   `tfsdk:"password_expired"`
	ExternallyManaged                types.Bool   `tfsdk:"externally_managed"`
	TypeOf2fa                        types.String `tfsdk:"type_of_2fa"`
	TypeOf2faForDisplay              types.String `tfsdk:"type_of_2fa_for_display"`
	DaysRemainingUntilPasswordExpire types.Int64  `tfsdk:"days_remaining_until_password_expire"`
	PasswordExpireAt                 types.String `tfsdk:"password_expire_at"`
}

func (r *userResource) Configure(_ context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {
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

func (r *userResource) Metadata(_ context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_user"
}

func (r *userResource) Schema(_ context.Context, _ resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Description: "",
		Attributes: map[string]schema.Attribute{
			"username": schema.StringAttribute{
				Description: "User's username",
				Required:    true,
			},
			"allowed_ips": schema.StringAttribute{
				Description: "A list of allowed IPs if applicable.  Newline delimited",
				Computed:    true,
				Optional:    true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
			},
			"attachments_permission": schema.BoolAttribute{
				Description: "If `true`, the user can user create Bundles (aka Share Links). Use the bundle permission instead.",
				Computed:    true,
				Optional:    true,
				PlanModifiers: []planmodifier.Bool{
					boolplanmodifier.UseStateForUnknown(),
				},
			},
			"authenticate_until": schema.StringAttribute{
				Description: "Scheduled Date/Time at which user will be deactivated",
				Computed:    true,
				Optional:    true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
			},
			"authentication_method": schema.StringAttribute{
				Description: "How is this user authenticated?",
				Computed:    true,
				Optional:    true,
				Validators: []validator.String{
					stringvalidator.OneOf("password", "sso", "none", "email_signup", "password_with_imported_hash", "password_and_ssh_key"),
				},
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
			},
			"billing_permission": schema.BoolAttribute{
				Description: "Allow this user to perform operations on the account, payments, and invoices?",
				Computed:    true,
				Optional:    true,
				PlanModifiers: []planmodifier.Bool{
					boolplanmodifier.UseStateForUnknown(),
				},
			},
			"bypass_site_allowed_ips": schema.BoolAttribute{
				Description: "Allow this user to skip site-wide IP blacklists?",
				Computed:    true,
				Optional:    true,
				PlanModifiers: []planmodifier.Bool{
					boolplanmodifier.UseStateForUnknown(),
				},
			},
			"bypass_inactive_disable": schema.BoolAttribute{
				Description: "Exempt this user from being disabled based on inactivity?",
				Computed:    true,
				Optional:    true,
				PlanModifiers: []planmodifier.Bool{
					boolplanmodifier.UseStateForUnknown(),
				},
			},
			"dav_permission": schema.BoolAttribute{
				Description: "Can the user connect with WebDAV?",
				Computed:    true,
				Optional:    true,
				PlanModifiers: []planmodifier.Bool{
					boolplanmodifier.UseStateForUnknown(),
				},
			},
			"disabled": schema.BoolAttribute{
				Description: "Is user disabled? Disabled users cannot log in, and do not count for billing purposes. Users can be automatically disabled after an inactivity period via a Site setting or schedule to be deactivated after specific date.",
				Computed:    true,
				Optional:    true,
				PlanModifiers: []planmodifier.Bool{
					boolplanmodifier.UseStateForUnknown(),
				},
			},
			"email": schema.StringAttribute{
				Description: "User email address",
				Computed:    true,
				Optional:    true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
			},
			"ftp_permission": schema.BoolAttribute{
				Description: "Can the user access with FTP/FTPS?",
				Computed:    true,
				Optional:    true,
				PlanModifiers: []planmodifier.Bool{
					boolplanmodifier.UseStateForUnknown(),
				},
			},
			"group_ids": schema.StringAttribute{
				Description: "Comma-separated list of group IDs of which this user is a member",
				Computed:    true,
				Optional:    true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
			},
			"header_text": schema.StringAttribute{
				Description: "Text to display to the user in the header of the UI",
				Computed:    true,
				Optional:    true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
			},
			"language": schema.StringAttribute{
				Description: "Preferred language",
				Computed:    true,
				Optional:    true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
			},
			"name": schema.StringAttribute{
				Description: "User's full name",
				Computed:    true,
				Optional:    true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
			},
			"company": schema.StringAttribute{
				Description: "User's company",
				Computed:    true,
				Optional:    true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
			},
			"notes": schema.StringAttribute{
				Description: "Any internal notes on the user",
				Computed:    true,
				Optional:    true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
			},
			"notification_daily_send_time": schema.Int64Attribute{
				Description: "Hour of the day at which daily notifications should be sent. Can be in range 0 to 23",
				Computed:    true,
				Optional:    true,
				PlanModifiers: []planmodifier.Int64{
					int64planmodifier.UseStateForUnknown(),
				},
			},
			"office_integration_enabled": schema.BoolAttribute{
				Description: "Enable integration with Office for the web?",
				Computed:    true,
				Optional:    true,
				PlanModifiers: []planmodifier.Bool{
					boolplanmodifier.UseStateForUnknown(),
				},
			},
			"password_validity_days": schema.Int64Attribute{
				Description: "Number of days to allow user to use the same password",
				Computed:    true,
				Optional:    true,
				PlanModifiers: []planmodifier.Int64{
					int64planmodifier.UseStateForUnknown(),
				},
			},
			"receive_admin_alerts": schema.BoolAttribute{
				Description: "Should the user receive admin alerts such a certificate expiration notifications and overages?",
				Computed:    true,
				Optional:    true,
				PlanModifiers: []planmodifier.Bool{
					boolplanmodifier.UseStateForUnknown(),
				},
			},
			"require_2fa": schema.StringAttribute{
				Description: "2FA required setting",
				Computed:    true,
				Optional:    true,
				Validators: []validator.String{
					stringvalidator.OneOf("use_system_setting", "always_require", "never_require"),
				},
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
			},
			"require_login_by": schema.StringAttribute{
				Description: "Require user to login by specified date otherwise it will be disabled.",
				Computed:    true,
				Optional:    true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
			},
			"require_password_change": schema.BoolAttribute{
				Description: "Is a password change required upon next user login?",
				Computed:    true,
				Optional:    true,
				PlanModifiers: []planmodifier.Bool{
					boolplanmodifier.UseStateForUnknown(),
				},
			},
			"restapi_permission": schema.BoolAttribute{
				Description: "Can this user access the Web app, Desktop app, SDKs, or REST API?  (All of these tools use the API internally, so this is one unified permission set.)",
				Computed:    true,
				Optional:    true,
				PlanModifiers: []planmodifier.Bool{
					boolplanmodifier.UseStateForUnknown(),
				},
			},
			"self_managed": schema.BoolAttribute{
				Description: "Does this user manage it's own credentials or is it a shared/bot user?",
				Computed:    true,
				Optional:    true,
				PlanModifiers: []planmodifier.Bool{
					boolplanmodifier.UseStateForUnknown(),
				},
			},
			"sftp_permission": schema.BoolAttribute{
				Description: "Can the user access with SFTP?",
				Computed:    true,
				Optional:    true,
				PlanModifiers: []planmodifier.Bool{
					boolplanmodifier.UseStateForUnknown(),
				},
			},
			"site_admin": schema.BoolAttribute{
				Description: "Is the user an administrator for this site?",
				Computed:    true,
				Optional:    true,
				PlanModifiers: []planmodifier.Bool{
					boolplanmodifier.UseStateForUnknown(),
				},
			},
			"skip_welcome_screen": schema.BoolAttribute{
				Description: "Skip Welcome page in the UI?",
				Computed:    true,
				Optional:    true,
				PlanModifiers: []planmodifier.Bool{
					boolplanmodifier.UseStateForUnknown(),
				},
			},
			"ssl_required": schema.StringAttribute{
				Description: "SSL required setting",
				Computed:    true,
				Optional:    true,
				Validators: []validator.String{
					stringvalidator.OneOf("use_system_setting", "always_require", "never_require"),
				},
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
			},
			"sso_strategy_id": schema.Int64Attribute{
				Description: "SSO (Single Sign On) strategy ID for the user, if applicable.",
				Computed:    true,
				Optional:    true,
				PlanModifiers: []planmodifier.Int64{
					int64planmodifier.UseStateForUnknown(),
				},
			},
			"subscribe_to_newsletter": schema.BoolAttribute{
				Description: "Is the user subscribed to the newsletter?",
				Computed:    true,
				Optional:    true,
				PlanModifiers: []planmodifier.Bool{
					boolplanmodifier.UseStateForUnknown(),
				},
			},
			"time_zone": schema.StringAttribute{
				Description: "User time zone",
				Computed:    true,
				Optional:    true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
			},
			"user_root": schema.StringAttribute{
				Description: "Root folder for FTP (and optionally SFTP if the appropriate site-wide setting is set.)  Note that this is not used for API, Desktop, or Web interface.",
				Computed:    true,
				Optional:    true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
			},
			"avatar_delete": schema.BoolAttribute{
				Description: "If true, the avatar will be deleted.",
				Optional:    true,
			},
			"change_password": schema.StringAttribute{
				Description: "Used for changing a password on an existing user.",
				Optional:    true,
			},
			"change_password_confirmation": schema.StringAttribute{
				Description: "Optional, but if provided, we will ensure that it matches the value sent in `change_password`.",
				Optional:    true,
			},
			"grant_permission": schema.StringAttribute{
				Description: "Permission to grant on the user root.  Can be blank or `full`, `read`, `write`, `list`, `read+write`, or `list+write`",
				Optional:    true,
			},
			"group_id": schema.Int64Attribute{
				Description: "Group ID to associate this user with.",
				Optional:    true,
			},
			"imported_password_hash": schema.StringAttribute{
				Description: "Pre-calculated hash of the user's password. If supplied, this will be used to authenticate the user on first login. Supported hash menthods are MD5, SHA1, and SHA256.",
				Optional:    true,
			},
			"password": schema.StringAttribute{
				Description: "User password.",
				Optional:    true,
			},
			"password_confirmation": schema.StringAttribute{
				Description: "Optional, but if provided, we will ensure that it matches the value sent in `password`.",
				Optional:    true,
			},
			"announcements_read": schema.BoolAttribute{
				Description: "Signifies that the user has read all the announcements in the UI.",
				Optional:    true,
			},
			"id": schema.Int64Attribute{
				Description: "User ID",
				Computed:    true,
				PlanModifiers: []planmodifier.Int64{
					int64planmodifier.UseStateForUnknown(),
				},
			},
			"admin_group_ids": schema.ListAttribute{
				Description: "List of group IDs of which this user is an administrator",
				Computed:    true,
				ElementType: types.Int64Type,
			},
			"api_keys_count": schema.Int64Attribute{
				Description: "Number of API keys associated with this user",
				Computed:    true,
			},
			"avatar_url": schema.StringAttribute{
				Description: "URL holding the user's avatar",
				Computed:    true,
			},
			"created_at": schema.StringAttribute{
				Description: "When this user was created",
				Computed:    true,
			},
			"disabled_expired_or_inactive": schema.BoolAttribute{
				Description: "Computed property that returns true if user disabled or expired or inactive.",
				Computed:    true,
			},
			"first_login_at": schema.StringAttribute{
				Description: "User's first login time",
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
			"password_set_at": schema.StringAttribute{
				Description: "Last time the user's password was set",
				Computed:    true,
			},
			"public_keys_count": schema.Int64Attribute{
				Description: "Number of public keys associated with this user",
				Computed:    true,
			},
			"active_2fa": schema.BoolAttribute{
				Description: "Is 2fa active for the user?",
				Computed:    true,
			},
			"password_expired": schema.BoolAttribute{
				Description: "Is user's password expired?",
				Computed:    true,
			},
			"externally_managed": schema.BoolAttribute{
				Description: "Is this user managed by a SsoStrategy?",
				Computed:    true,
			},
			"type_of_2fa": schema.StringAttribute{
				Description: "Type(s) of 2FA methods in use, for programmatic use.  Will be either `sms`, `totp`, `u2f`, `yubi`, or multiple values sorted alphabetically and joined by an underscore.  Does not specify whether user has more than one of a given method.",
				Computed:    true,
			},
			"type_of_2fa_for_display": schema.StringAttribute{
				Description: "Type(s) of 2FA methods in use, formatted for displaying in the UI.  Unlike `type_of_2fa`, this value will make clear when a user has more than 1 of the same type of method.",
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

func (r *userResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	var plan userResourceModel
	diags := req.Plan.Get(ctx, &plan)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	paramsUserCreate := files_sdk.UserCreateParams{}
	if !plan.AvatarDelete.IsNull() && !plan.AvatarDelete.IsUnknown() {
		paramsUserCreate.AvatarDelete = plan.AvatarDelete.ValueBoolPointer()
	}
	paramsUserCreate.ChangePassword = plan.ChangePassword.ValueString()
	paramsUserCreate.ChangePasswordConfirmation = plan.ChangePasswordConfirmation.ValueString()
	paramsUserCreate.Email = plan.Email.ValueString()
	paramsUserCreate.GrantPermission = plan.GrantPermission.ValueString()
	paramsUserCreate.GroupId = plan.GroupId.ValueInt64()
	paramsUserCreate.GroupIds = plan.GroupIds.ValueString()
	paramsUserCreate.ImportedPasswordHash = plan.ImportedPasswordHash.ValueString()
	paramsUserCreate.Password = plan.Password.ValueString()
	paramsUserCreate.PasswordConfirmation = plan.PasswordConfirmation.ValueString()
	if !plan.AnnouncementsRead.IsNull() && !plan.AnnouncementsRead.IsUnknown() {
		paramsUserCreate.AnnouncementsRead = plan.AnnouncementsRead.ValueBoolPointer()
	}
	paramsUserCreate.AllowedIps = plan.AllowedIps.ValueString()
	if !plan.AttachmentsPermission.IsNull() && !plan.AttachmentsPermission.IsUnknown() {
		paramsUserCreate.AttachmentsPermission = plan.AttachmentsPermission.ValueBoolPointer()
	}
	if !plan.AuthenticateUntil.IsNull() && plan.AuthenticateUntil.ValueString() != "" {
		createAuthenticateUntil, err := time.Parse(time.RFC3339, plan.AuthenticateUntil.ValueString())
		if err != nil {
			resp.Diagnostics.AddAttributeError(
				path.Root("authenticate_until"),
				"Error Parsing authenticate_until Time",
				"Could not parse authenticate_until time: "+err.Error(),
			)
		} else {
			paramsUserCreate.AuthenticateUntil = &createAuthenticateUntil
		}
	}
	paramsUserCreate.AuthenticationMethod = paramsUserCreate.AuthenticationMethod.Enum()[plan.AuthenticationMethod.ValueString()]
	if !plan.BillingPermission.IsNull() && !plan.BillingPermission.IsUnknown() {
		paramsUserCreate.BillingPermission = plan.BillingPermission.ValueBoolPointer()
	}
	if !plan.BypassInactiveDisable.IsNull() && !plan.BypassInactiveDisable.IsUnknown() {
		paramsUserCreate.BypassInactiveDisable = plan.BypassInactiveDisable.ValueBoolPointer()
	}
	if !plan.BypassSiteAllowedIps.IsNull() && !plan.BypassSiteAllowedIps.IsUnknown() {
		paramsUserCreate.BypassSiteAllowedIps = plan.BypassSiteAllowedIps.ValueBoolPointer()
	}
	if !plan.DavPermission.IsNull() && !plan.DavPermission.IsUnknown() {
		paramsUserCreate.DavPermission = plan.DavPermission.ValueBoolPointer()
	}
	if !plan.Disabled.IsNull() && !plan.Disabled.IsUnknown() {
		paramsUserCreate.Disabled = plan.Disabled.ValueBoolPointer()
	}
	if !plan.FtpPermission.IsNull() && !plan.FtpPermission.IsUnknown() {
		paramsUserCreate.FtpPermission = plan.FtpPermission.ValueBoolPointer()
	}
	paramsUserCreate.HeaderText = plan.HeaderText.ValueString()
	paramsUserCreate.Language = plan.Language.ValueString()
	paramsUserCreate.NotificationDailySendTime = plan.NotificationDailySendTime.ValueInt64()
	paramsUserCreate.Name = plan.Name.ValueString()
	paramsUserCreate.Company = plan.Company.ValueString()
	paramsUserCreate.Notes = plan.Notes.ValueString()
	if !plan.OfficeIntegrationEnabled.IsNull() && !plan.OfficeIntegrationEnabled.IsUnknown() {
		paramsUserCreate.OfficeIntegrationEnabled = plan.OfficeIntegrationEnabled.ValueBoolPointer()
	}
	paramsUserCreate.PasswordValidityDays = plan.PasswordValidityDays.ValueInt64()
	if !plan.ReceiveAdminAlerts.IsNull() && !plan.ReceiveAdminAlerts.IsUnknown() {
		paramsUserCreate.ReceiveAdminAlerts = plan.ReceiveAdminAlerts.ValueBoolPointer()
	}
	if !plan.RequireLoginBy.IsNull() && plan.RequireLoginBy.ValueString() != "" {
		createRequireLoginBy, err := time.Parse(time.RFC3339, plan.RequireLoginBy.ValueString())
		if err != nil {
			resp.Diagnostics.AddAttributeError(
				path.Root("require_login_by"),
				"Error Parsing require_login_by Time",
				"Could not parse require_login_by time: "+err.Error(),
			)
		} else {
			paramsUserCreate.RequireLoginBy = &createRequireLoginBy
		}
	}
	if !plan.RequirePasswordChange.IsNull() && !plan.RequirePasswordChange.IsUnknown() {
		paramsUserCreate.RequirePasswordChange = plan.RequirePasswordChange.ValueBoolPointer()
	}
	if !plan.RestapiPermission.IsNull() && !plan.RestapiPermission.IsUnknown() {
		paramsUserCreate.RestapiPermission = plan.RestapiPermission.ValueBoolPointer()
	}
	if !plan.SelfManaged.IsNull() && !plan.SelfManaged.IsUnknown() {
		paramsUserCreate.SelfManaged = plan.SelfManaged.ValueBoolPointer()
	}
	if !plan.SftpPermission.IsNull() && !plan.SftpPermission.IsUnknown() {
		paramsUserCreate.SftpPermission = plan.SftpPermission.ValueBoolPointer()
	}
	if !plan.SiteAdmin.IsNull() && !plan.SiteAdmin.IsUnknown() {
		paramsUserCreate.SiteAdmin = plan.SiteAdmin.ValueBoolPointer()
	}
	if !plan.SkipWelcomeScreen.IsNull() && !plan.SkipWelcomeScreen.IsUnknown() {
		paramsUserCreate.SkipWelcomeScreen = plan.SkipWelcomeScreen.ValueBoolPointer()
	}
	paramsUserCreate.SslRequired = paramsUserCreate.SslRequired.Enum()[plan.SslRequired.ValueString()]
	paramsUserCreate.SsoStrategyId = plan.SsoStrategyId.ValueInt64()
	if !plan.SubscribeToNewsletter.IsNull() && !plan.SubscribeToNewsletter.IsUnknown() {
		paramsUserCreate.SubscribeToNewsletter = plan.SubscribeToNewsletter.ValueBoolPointer()
	}
	paramsUserCreate.Require2fa = paramsUserCreate.Require2fa.Enum()[plan.Require2fa.ValueString()]
	paramsUserCreate.TimeZone = plan.TimeZone.ValueString()
	paramsUserCreate.UserRoot = plan.UserRoot.ValueString()
	paramsUserCreate.Username = plan.Username.ValueString()

	if resp.Diagnostics.HasError() {
		return
	}

	user, err := r.client.Create(paramsUserCreate, files_sdk.WithContext(ctx))
	if err != nil {
		resp.Diagnostics.AddError(
			"Error Creating Files User",
			"Could not create user, unexpected error: "+err.Error(),
		)
		return
	}

	diags = r.populateResourceModel(ctx, user, &plan)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	diags = resp.State.Set(ctx, plan)
	resp.Diagnostics.Append(diags...)
}

func (r *userResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	var state userResourceModel
	diags := req.State.Get(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	paramsUserFind := files_sdk.UserFindParams{}
	paramsUserFind.Id = state.Id.ValueInt64()

	user, err := r.client.Find(paramsUserFind, files_sdk.WithContext(ctx))
	if err != nil {
		if files_sdk.IsNotExist(err) {
			resp.State.RemoveResource(ctx)
			return
		}

		resp.Diagnostics.AddError(
			"Error Reading Files User",
			"Could not read user id "+fmt.Sprint(state.Id.ValueInt64())+": "+err.Error(),
		)
		return
	}

	diags = r.populateResourceModel(ctx, user, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	diags = resp.State.Set(ctx, state)
	resp.Diagnostics.Append(diags...)
}

func (r *userResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	var plan userResourceModel
	diags := req.Plan.Get(ctx, &plan)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	paramsUserUpdate := files_sdk.UserUpdateParams{}
	paramsUserUpdate.Id = plan.Id.ValueInt64()
	if !plan.AvatarDelete.IsNull() && !plan.AvatarDelete.IsUnknown() {
		paramsUserUpdate.AvatarDelete = plan.AvatarDelete.ValueBoolPointer()
	}
	paramsUserUpdate.ChangePassword = plan.ChangePassword.ValueString()
	paramsUserUpdate.ChangePasswordConfirmation = plan.ChangePasswordConfirmation.ValueString()
	paramsUserUpdate.Email = plan.Email.ValueString()
	paramsUserUpdate.GrantPermission = plan.GrantPermission.ValueString()
	paramsUserUpdate.GroupId = plan.GroupId.ValueInt64()
	paramsUserUpdate.GroupIds = plan.GroupIds.ValueString()
	paramsUserUpdate.ImportedPasswordHash = plan.ImportedPasswordHash.ValueString()
	paramsUserUpdate.Password = plan.Password.ValueString()
	paramsUserUpdate.PasswordConfirmation = plan.PasswordConfirmation.ValueString()
	if !plan.AnnouncementsRead.IsNull() && !plan.AnnouncementsRead.IsUnknown() {
		paramsUserUpdate.AnnouncementsRead = plan.AnnouncementsRead.ValueBoolPointer()
	}
	paramsUserUpdate.AllowedIps = plan.AllowedIps.ValueString()
	if !plan.AttachmentsPermission.IsNull() && !plan.AttachmentsPermission.IsUnknown() {
		paramsUserUpdate.AttachmentsPermission = plan.AttachmentsPermission.ValueBoolPointer()
	}
	if !plan.AuthenticateUntil.IsNull() && plan.AuthenticateUntil.ValueString() != "" {
		updateAuthenticateUntil, err := time.Parse(time.RFC3339, plan.AuthenticateUntil.ValueString())
		if err != nil {
			resp.Diagnostics.AddAttributeError(
				path.Root("authenticate_until"),
				"Error Parsing authenticate_until Time",
				"Could not parse authenticate_until time: "+err.Error(),
			)
		} else {
			paramsUserUpdate.AuthenticateUntil = &updateAuthenticateUntil
		}
	}
	paramsUserUpdate.AuthenticationMethod = paramsUserUpdate.AuthenticationMethod.Enum()[plan.AuthenticationMethod.ValueString()]
	if !plan.BillingPermission.IsNull() && !plan.BillingPermission.IsUnknown() {
		paramsUserUpdate.BillingPermission = plan.BillingPermission.ValueBoolPointer()
	}
	if !plan.BypassInactiveDisable.IsNull() && !plan.BypassInactiveDisable.IsUnknown() {
		paramsUserUpdate.BypassInactiveDisable = plan.BypassInactiveDisable.ValueBoolPointer()
	}
	if !plan.BypassSiteAllowedIps.IsNull() && !plan.BypassSiteAllowedIps.IsUnknown() {
		paramsUserUpdate.BypassSiteAllowedIps = plan.BypassSiteAllowedIps.ValueBoolPointer()
	}
	if !plan.DavPermission.IsNull() && !plan.DavPermission.IsUnknown() {
		paramsUserUpdate.DavPermission = plan.DavPermission.ValueBoolPointer()
	}
	if !plan.Disabled.IsNull() && !plan.Disabled.IsUnknown() {
		paramsUserUpdate.Disabled = plan.Disabled.ValueBoolPointer()
	}
	if !plan.FtpPermission.IsNull() && !plan.FtpPermission.IsUnknown() {
		paramsUserUpdate.FtpPermission = plan.FtpPermission.ValueBoolPointer()
	}
	paramsUserUpdate.HeaderText = plan.HeaderText.ValueString()
	paramsUserUpdate.Language = plan.Language.ValueString()
	paramsUserUpdate.NotificationDailySendTime = plan.NotificationDailySendTime.ValueInt64()
	paramsUserUpdate.Name = plan.Name.ValueString()
	paramsUserUpdate.Company = plan.Company.ValueString()
	paramsUserUpdate.Notes = plan.Notes.ValueString()
	if !plan.OfficeIntegrationEnabled.IsNull() && !plan.OfficeIntegrationEnabled.IsUnknown() {
		paramsUserUpdate.OfficeIntegrationEnabled = plan.OfficeIntegrationEnabled.ValueBoolPointer()
	}
	paramsUserUpdate.PasswordValidityDays = plan.PasswordValidityDays.ValueInt64()
	if !plan.ReceiveAdminAlerts.IsNull() && !plan.ReceiveAdminAlerts.IsUnknown() {
		paramsUserUpdate.ReceiveAdminAlerts = plan.ReceiveAdminAlerts.ValueBoolPointer()
	}
	if !plan.RequireLoginBy.IsNull() && plan.RequireLoginBy.ValueString() != "" {
		updateRequireLoginBy, err := time.Parse(time.RFC3339, plan.RequireLoginBy.ValueString())
		if err != nil {
			resp.Diagnostics.AddAttributeError(
				path.Root("require_login_by"),
				"Error Parsing require_login_by Time",
				"Could not parse require_login_by time: "+err.Error(),
			)
		} else {
			paramsUserUpdate.RequireLoginBy = &updateRequireLoginBy
		}
	}
	if !plan.RequirePasswordChange.IsNull() && !plan.RequirePasswordChange.IsUnknown() {
		paramsUserUpdate.RequirePasswordChange = plan.RequirePasswordChange.ValueBoolPointer()
	}
	if !plan.RestapiPermission.IsNull() && !plan.RestapiPermission.IsUnknown() {
		paramsUserUpdate.RestapiPermission = plan.RestapiPermission.ValueBoolPointer()
	}
	if !plan.SelfManaged.IsNull() && !plan.SelfManaged.IsUnknown() {
		paramsUserUpdate.SelfManaged = plan.SelfManaged.ValueBoolPointer()
	}
	if !plan.SftpPermission.IsNull() && !plan.SftpPermission.IsUnknown() {
		paramsUserUpdate.SftpPermission = plan.SftpPermission.ValueBoolPointer()
	}
	if !plan.SiteAdmin.IsNull() && !plan.SiteAdmin.IsUnknown() {
		paramsUserUpdate.SiteAdmin = plan.SiteAdmin.ValueBoolPointer()
	}
	if !plan.SkipWelcomeScreen.IsNull() && !plan.SkipWelcomeScreen.IsUnknown() {
		paramsUserUpdate.SkipWelcomeScreen = plan.SkipWelcomeScreen.ValueBoolPointer()
	}
	paramsUserUpdate.SslRequired = paramsUserUpdate.SslRequired.Enum()[plan.SslRequired.ValueString()]
	paramsUserUpdate.SsoStrategyId = plan.SsoStrategyId.ValueInt64()
	if !plan.SubscribeToNewsletter.IsNull() && !plan.SubscribeToNewsletter.IsUnknown() {
		paramsUserUpdate.SubscribeToNewsletter = plan.SubscribeToNewsletter.ValueBoolPointer()
	}
	paramsUserUpdate.Require2fa = paramsUserUpdate.Require2fa.Enum()[plan.Require2fa.ValueString()]
	paramsUserUpdate.TimeZone = plan.TimeZone.ValueString()
	paramsUserUpdate.UserRoot = plan.UserRoot.ValueString()
	paramsUserUpdate.Username = plan.Username.ValueString()

	if resp.Diagnostics.HasError() {
		return
	}

	user, err := r.client.Update(paramsUserUpdate, files_sdk.WithContext(ctx))
	if err != nil {
		resp.Diagnostics.AddError(
			"Error Updating Files User",
			"Could not update user, unexpected error: "+err.Error(),
		)
		return
	}

	diags = r.populateResourceModel(ctx, user, &plan)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	diags = resp.State.Set(ctx, plan)
	resp.Diagnostics.Append(diags...)
}

func (r *userResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	var state userResourceModel
	diags := req.State.Get(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	paramsUserDelete := files_sdk.UserDeleteParams{}
	paramsUserDelete.Id = state.Id.ValueInt64()

	err := r.client.Delete(paramsUserDelete, files_sdk.WithContext(ctx))
	if err != nil && !files_sdk.IsNotExist(err) {
		resp.Diagnostics.AddError(
			"Error Deleting Files User",
			"Could not delete user id "+fmt.Sprint(state.Id.ValueInt64())+": "+err.Error(),
		)
	}
}

func (r *userResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
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

func (r *userResource) populateResourceModel(ctx context.Context, user files_sdk.User, state *userResourceModel) (diags diag.Diagnostics) {
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
	state.BillingPermission = types.BoolPointerValue(user.BillingPermission)
	state.BypassSiteAllowedIps = types.BoolPointerValue(user.BypassSiteAllowedIps)
	state.BypassInactiveDisable = types.BoolPointerValue(user.BypassInactiveDisable)
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
	state.RestapiPermission = types.BoolPointerValue(user.RestapiPermission)
	state.SelfManaged = types.BoolPointerValue(user.SelfManaged)
	state.SftpPermission = types.BoolPointerValue(user.SftpPermission)
	state.SiteAdmin = types.BoolPointerValue(user.SiteAdmin)
	state.SkipWelcomeScreen = types.BoolPointerValue(user.SkipWelcomeScreen)
	state.SslRequired = types.StringValue(user.SslRequired)
	state.SsoStrategyId = types.Int64Value(user.SsoStrategyId)
	state.SubscribeToNewsletter = types.BoolPointerValue(user.SubscribeToNewsletter)
	state.ExternallyManaged = types.BoolPointerValue(user.ExternallyManaged)
	state.TimeZone = types.StringValue(user.TimeZone)
	state.TypeOf2fa = types.StringValue(user.TypeOf2fa)
	state.TypeOf2faForDisplay = types.StringValue(user.TypeOf2faForDisplay)
	state.UserRoot = types.StringValue(user.UserRoot)
	state.DaysRemainingUntilPasswordExpire = types.Int64Value(user.DaysRemainingUntilPasswordExpire)
	if err := lib.TimeToStringType(ctx, path.Root("password_expire_at"), user.PasswordExpireAt, &state.PasswordExpireAt); err != nil {
		diags.AddError(
			"Error Creating Files User",
			"Could not convert state password_expire_at to string: "+err.Error(),
		)
	}

	return
}
