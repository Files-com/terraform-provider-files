package provider

import (
	"context"
	"encoding/json"
	"fmt"

	files_sdk "github.com/Files-com/files-sdk-go/v3"
	site "github.com/Files-com/files-sdk-go/v3/site"
	"github.com/Files-com/terraform-provider-files/lib"
	"github.com/hashicorp/terraform-plugin-framework-validators/stringvalidator"
	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/boolplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/dynamicplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/int64planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/listplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/schema/validator"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

var (
	_ resource.Resource                = &siteResource{}
	_ resource.ResourceWithConfigure   = &siteResource{}
	_ resource.ResourceWithImportState = &siteResource{}
)

func NewSiteResource() resource.Resource {
	return &siteResource{}
}

type siteResource struct {
	client *site.Client
}

type siteResourceModel struct {
	Name                                     types.String  `tfsdk:"name"`
	AdditionalTextFileTypes                  types.List    `tfsdk:"additional_text_file_types"`
	Allowed2faMethodSms                      types.Bool    `tfsdk:"allowed_2fa_method_sms"`
	Allowed2faMethodTotp                     types.Bool    `tfsdk:"allowed_2fa_method_totp"`
	Allowed2faMethodWebauthn                 types.Bool    `tfsdk:"allowed_2fa_method_webauthn"`
	Allowed2faMethodYubi                     types.Bool    `tfsdk:"allowed_2fa_method_yubi"`
	Allowed2faMethodEmail                    types.Bool    `tfsdk:"allowed_2fa_method_email"`
	Allowed2faMethodStatic                   types.Bool    `tfsdk:"allowed_2fa_method_static"`
	Allowed2faMethodBypassForFtpSftpDav      types.Bool    `tfsdk:"allowed_2fa_method_bypass_for_ftp_sftp_dav"`
	AdminsBypassLockedSubfolders             types.Bool    `tfsdk:"admins_bypass_locked_subfolders"`
	AllowBundleNames                         types.Bool    `tfsdk:"allow_bundle_names"`
	AllowedCountries                         types.String  `tfsdk:"allowed_countries"`
	AllowedIps                               types.String  `tfsdk:"allowed_ips"`
	AlwaysMkdirParents                       types.Bool    `tfsdk:"always_mkdir_parents"`
	AskAboutOverwrites                       types.Bool    `tfsdk:"ask_about_overwrites"`
	BundleActivityNotifications              types.String  `tfsdk:"bundle_activity_notifications"`
	BundleExpiration                         types.Int64   `tfsdk:"bundle_expiration"`
	BundleNotFoundMessage                    types.String  `tfsdk:"bundle_not_found_message"`
	BundlePasswordRequired                   types.Bool    `tfsdk:"bundle_password_required"`
	BundleRecipientBlacklistDomains          types.List    `tfsdk:"bundle_recipient_blacklist_domains"`
	BundleRecipientBlacklistFreeEmailDomains types.Bool    `tfsdk:"bundle_recipient_blacklist_free_email_domains"`
	BundleRegistrationNotifications          types.String  `tfsdk:"bundle_registration_notifications"`
	BundleRequireRegistration                types.Bool    `tfsdk:"bundle_require_registration"`
	BundleRequireShareRecipient              types.Bool    `tfsdk:"bundle_require_share_recipient"`
	BundleRequireNote                        types.Bool    `tfsdk:"bundle_require_note"`
	BundleSendSharedReceipts                 types.Bool    `tfsdk:"bundle_send_shared_receipts"`
	BundleUploadReceiptNotifications         types.String  `tfsdk:"bundle_upload_receipt_notifications"`
	BundleWatermarkValue                     types.Dynamic `tfsdk:"bundle_watermark_value"`
	CalculateFileChecksumsCrc32              types.Bool    `tfsdk:"calculate_file_checksums_crc32"`
	CalculateFileChecksumsMd5                types.Bool    `tfsdk:"calculate_file_checksums_md5"`
	CalculateFileChecksumsSha1               types.Bool    `tfsdk:"calculate_file_checksums_sha1"`
	CalculateFileChecksumsSha256             types.Bool    `tfsdk:"calculate_file_checksums_sha256"`
	UploadsViaEmailAuthentication            types.Bool    `tfsdk:"uploads_via_email_authentication"`
	Color2Left                               types.String  `tfsdk:"color2_left"`
	Color2Link                               types.String  `tfsdk:"color2_link"`
	Color2Text                               types.String  `tfsdk:"color2_text"`
	Color2Top                                types.String  `tfsdk:"color2_top"`
	Color2TopText                            types.String  `tfsdk:"color2_top_text"`
	CustomNamespace                          types.Bool    `tfsdk:"custom_namespace"`
	DavEnabled                               types.Bool    `tfsdk:"dav_enabled"`
	DavUserRootEnabled                       types.Bool    `tfsdk:"dav_user_root_enabled"`
	DaysBeforeDeletingDisabledUsers          types.Int64   `tfsdk:"days_before_deleting_disabled_users"`
	DaysToRetainBackups                      types.Int64   `tfsdk:"days_to_retain_backups"`
	DocumentEditsInBundleAllowed             types.Bool    `tfsdk:"document_edits_in_bundle_allowed"`
	DefaultTimeZone                          types.String  `tfsdk:"default_time_zone"`
	DesktopApp                               types.Bool    `tfsdk:"desktop_app"`
	DesktopAppSessionIpPinning               types.Bool    `tfsdk:"desktop_app_session_ip_pinning"`
	DesktopAppSessionLifetime                types.Int64   `tfsdk:"desktop_app_session_lifetime"`
	LegacyChecksumsMode                      types.Bool    `tfsdk:"legacy_checksums_mode"`
	MobileApp                                types.Bool    `tfsdk:"mobile_app"`
	MobileAppSessionIpPinning                types.Bool    `tfsdk:"mobile_app_session_ip_pinning"`
	MobileAppSessionLifetime                 types.Int64   `tfsdk:"mobile_app_session_lifetime"`
	DisallowedCountries                      types.String  `tfsdk:"disallowed_countries"`
	DisableFilesCertificateGeneration        types.Bool    `tfsdk:"disable_files_certificate_generation"`
	DisablePasswordReset                     types.Bool    `tfsdk:"disable_password_reset"`
	Domain                                   types.String  `tfsdk:"domain"`
	DomainHstsHeader                         types.Bool    `tfsdk:"domain_hsts_header"`
	DomainLetsencryptChain                   types.String  `tfsdk:"domain_letsencrypt_chain"`
	Email                                    types.String  `tfsdk:"email"`
	FtpEnabled                               types.Bool    `tfsdk:"ftp_enabled"`
	ReplyToEmail                             types.String  `tfsdk:"reply_to_email"`
	NonSsoGroupsAllowed                      types.Bool    `tfsdk:"non_sso_groups_allowed"`
	NonSsoUsersAllowed                       types.Bool    `tfsdk:"non_sso_users_allowed"`
	FolderPermissionsGroupsOnly              types.Bool    `tfsdk:"folder_permissions_groups_only"`
	IncludePasswordInWelcomeEmail            types.Bool    `tfsdk:"include_password_in_welcome_email"`
	Language                                 types.String  `tfsdk:"language"`
	LdapBaseDn                               types.String  `tfsdk:"ldap_base_dn"`
	LdapDomain                               types.String  `tfsdk:"ldap_domain"`
	LdapEnabled                              types.Bool    `tfsdk:"ldap_enabled"`
	LdapGroupAction                          types.String  `tfsdk:"ldap_group_action"`
	LdapGroupExclusion                       types.String  `tfsdk:"ldap_group_exclusion"`
	LdapGroupInclusion                       types.String  `tfsdk:"ldap_group_inclusion"`
	LdapHost                                 types.String  `tfsdk:"ldap_host"`
	LdapHost2                                types.String  `tfsdk:"ldap_host_2"`
	LdapHost3                                types.String  `tfsdk:"ldap_host_3"`
	LdapPort                                 types.Int64   `tfsdk:"ldap_port"`
	LdapSecure                               types.Bool    `tfsdk:"ldap_secure"`
	LdapType                                 types.String  `tfsdk:"ldap_type"`
	LdapUserAction                           types.String  `tfsdk:"ldap_user_action"`
	LdapUserIncludeGroups                    types.String  `tfsdk:"ldap_user_include_groups"`
	LdapUsername                             types.String  `tfsdk:"ldap_username"`
	LdapUsernameField                        types.String  `tfsdk:"ldap_username_field"`
	LoginHelpText                            types.String  `tfsdk:"login_help_text"`
	MaxPriorPasswords                        types.Int64   `tfsdk:"max_prior_passwords"`
	MotdText                                 types.String  `tfsdk:"motd_text"`
	MotdUseForFtp                            types.Bool    `tfsdk:"motd_use_for_ftp"`
	MotdUseForSftp                           types.Bool    `tfsdk:"motd_use_for_sftp"`
	OfficeIntegrationAvailable               types.Bool    `tfsdk:"office_integration_available"`
	OfficeIntegrationType                    types.String  `tfsdk:"office_integration_type"`
	OptOutGlobal                             types.Bool    `tfsdk:"opt_out_global"`
	PasswordMinLength                        types.Int64   `tfsdk:"password_min_length"`
	PasswordRequireLetter                    types.Bool    `tfsdk:"password_require_letter"`
	PasswordRequireMixed                     types.Bool    `tfsdk:"password_require_mixed"`
	PasswordRequireNumber                    types.Bool    `tfsdk:"password_require_number"`
	PasswordRequireSpecial                   types.Bool    `tfsdk:"password_require_special"`
	PasswordRequireUnbreached                types.Bool    `tfsdk:"password_require_unbreached"`
	PasswordRequirementsApplyToBundles       types.Bool    `tfsdk:"password_requirements_apply_to_bundles"`
	PasswordValidityDays                     types.Int64   `tfsdk:"password_validity_days"`
	PinAllRemoteServersToSiteRegion          types.Bool    `tfsdk:"pin_all_remote_servers_to_site_region"`
	PreventRootPermissionsForNonSiteAdmins   types.Bool    `tfsdk:"prevent_root_permissions_for_non_site_admins"`
	ProtocolAccessGroupsOnly                 types.Bool    `tfsdk:"protocol_access_groups_only"`
	Require2fa                               types.Bool    `tfsdk:"require_2fa"`
	RevokeBundleAccessOnDisableOrDelete      types.Bool    `tfsdk:"revoke_bundle_access_on_disable_or_delete"`
	Require2faUserType                       types.String  `tfsdk:"require_2fa_user_type"`
	RequireLogoutFromBundlesAndInboxes       types.Bool    `tfsdk:"require_logout_from_bundles_and_inboxes"`
	SftpEnabled                              types.Bool    `tfsdk:"sftp_enabled"`
	SftpHostKeyType                          types.String  `tfsdk:"sftp_host_key_type"`
	ActiveSftpHostKeyId                      types.Int64   `tfsdk:"active_sftp_host_key_id"`
	SftpInsecureCiphers                      types.Bool    `tfsdk:"sftp_insecure_ciphers"`
	SftpInsecureDiffieHellman                types.Bool    `tfsdk:"sftp_insecure_diffie_hellman"`
	SftpUserRootEnabled                      types.Bool    `tfsdk:"sftp_user_root_enabled"`
	SharingEnabled                           types.Bool    `tfsdk:"sharing_enabled"`
	ShowUserNotificationsLogInLink           types.Bool    `tfsdk:"show_user_notifications_log_in_link"`
	ShowRequestAccessLink                    types.Bool    `tfsdk:"show_request_access_link"`
	SiteFooter                               types.String  `tfsdk:"site_footer"`
	SiteHeader                               types.String  `tfsdk:"site_header"`
	SmtpAddress                              types.String  `tfsdk:"smtp_address"`
	SmtpAuthentication                       types.String  `tfsdk:"smtp_authentication"`
	SmtpFrom                                 types.String  `tfsdk:"smtp_from"`
	SmtpPort                                 types.Int64   `tfsdk:"smtp_port"`
	SmtpUsername                             types.String  `tfsdk:"smtp_username"`
	SessionExpiry                            types.String  `tfsdk:"session_expiry"`
	SessionExpiryMinutes                     types.Int64   `tfsdk:"session_expiry_minutes"`
	SnapshotSharingEnabled                   types.Bool    `tfsdk:"snapshot_sharing_enabled"`
	SslRequired                              types.Bool    `tfsdk:"ssl_required"`
	Subdomain                                types.String  `tfsdk:"subdomain"`
	TlsDisabled                              types.Bool    `tfsdk:"tls_disabled"`
	UseDedicatedIpsForSmtp                   types.Bool    `tfsdk:"use_dedicated_ips_for_smtp"`
	UseProvidedModifiedAt                    types.Bool    `tfsdk:"use_provided_modified_at"`
	UserLockout                              types.Bool    `tfsdk:"user_lockout"`
	UserLockoutLockPeriod                    types.Int64   `tfsdk:"user_lockout_lock_period"`
	UserLockoutTries                         types.Int64   `tfsdk:"user_lockout_tries"`
	UserLockoutWithin                        types.Int64   `tfsdk:"user_lockout_within"`
	UserRequestsEnabled                      types.Bool    `tfsdk:"user_requests_enabled"`
	UserRequestsNotifyAdmins                 types.Bool    `tfsdk:"user_requests_notify_admins"`
	UsersCanCreateApiKeys                    types.Bool    `tfsdk:"users_can_create_api_keys"`
	UsersCanCreateSshKeys                    types.Bool    `tfsdk:"users_can_create_ssh_keys"`
	WelcomeCustomText                        types.String  `tfsdk:"welcome_custom_text"`
	WelcomeEmailCc                           types.String  `tfsdk:"welcome_email_cc"`
	WelcomeEmailSubject                      types.String  `tfsdk:"welcome_email_subject"`
	WelcomeEmailEnabled                      types.Bool    `tfsdk:"welcome_email_enabled"`
	WelcomeScreen                            types.String  `tfsdk:"welcome_screen"`
	WindowsModeFtp                           types.Bool    `tfsdk:"windows_mode_ftp"`
	GroupAdminsCanSetUserPassword            types.Bool    `tfsdk:"group_admins_can_set_user_password"`
	Id                                       types.Int64   `tfsdk:"id"`
	AdminUserId                              types.Int64   `tfsdk:"admin_user_id"`
	BundleWatermarkAttachment                types.String  `tfsdk:"bundle_watermark_attachment"`
	ContactName                              types.String  `tfsdk:"contact_name"`
	CreatedAt                                types.String  `tfsdk:"created_at"`
	Currency                                 types.String  `tfsdk:"currency"`
	DisableNotifications                     types.Bool    `tfsdk:"disable_notifications"`
	Hipaa                                    types.Bool    `tfsdk:"hipaa"`
	Icon128                                  types.String  `tfsdk:"icon128"`
	Icon16                                   types.String  `tfsdk:"icon16"`
	Icon32                                   types.String  `tfsdk:"icon32"`
	Icon48                                   types.String  `tfsdk:"icon48"`
	ImmutableFilesSetAt                      types.String  `tfsdk:"immutable_files_set_at"`
	Logo                                     types.String  `tfsdk:"logo"`
	LoginPageBackgroundImage                 types.String  `tfsdk:"login_page_background_image"`
	NextBillingAmount                        types.String  `tfsdk:"next_billing_amount"`
	NextBillingDate                          types.String  `tfsdk:"next_billing_date"`
	OncehubLink                              types.String  `tfsdk:"oncehub_link"`
	Overdue                                  types.Bool    `tfsdk:"overdue"`
	Phone                                    types.String  `tfsdk:"phone"`
	Require2faStopTime                       types.String  `tfsdk:"require_2fa_stop_time"`
	Session                                  types.String  `tfsdk:"session"`
	SwitchToPlanDate                         types.String  `tfsdk:"switch_to_plan_date"`
	TrialDaysLeft                            types.Int64   `tfsdk:"trial_days_left"`
	TrialUntil                               types.String  `tfsdk:"trial_until"`
	User                                     types.String  `tfsdk:"user"`
}

func (r *siteResource) Configure(_ context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {
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

	r.client = &site.Client{Config: sdk_config}
}

func (r *siteResource) Metadata(_ context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_site"
}

func (r *siteResource) Schema(_ context.Context, _ resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Description: "A Site is the place you'll come to update site settings, as well as manage site-wide API keys.\n\n\n\nMost site settings can be set via the API.",
		Attributes: map[string]schema.Attribute{
			"name": schema.StringAttribute{
				Description: "Site name",
				Computed:    true,
				Optional:    true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
			},
			"additional_text_file_types": schema.ListAttribute{
				Description: "Additional extensions that are considered text files",
				Computed:    true,
				Optional:    true,
				ElementType: types.StringType,
				PlanModifiers: []planmodifier.List{
					listplanmodifier.UseStateForUnknown(),
				},
			},
			"allowed_2fa_method_sms": schema.BoolAttribute{
				Description: "Is SMS two factor authentication allowed?",
				Computed:    true,
				Optional:    true,
				PlanModifiers: []planmodifier.Bool{
					boolplanmodifier.UseStateForUnknown(),
				},
			},
			"allowed_2fa_method_totp": schema.BoolAttribute{
				Description: "Is TOTP two factor authentication allowed?",
				Computed:    true,
				Optional:    true,
				PlanModifiers: []planmodifier.Bool{
					boolplanmodifier.UseStateForUnknown(),
				},
			},
			"allowed_2fa_method_webauthn": schema.BoolAttribute{
				Description: "Is WebAuthn two factor authentication allowed?",
				Computed:    true,
				Optional:    true,
				PlanModifiers: []planmodifier.Bool{
					boolplanmodifier.UseStateForUnknown(),
				},
			},
			"allowed_2fa_method_yubi": schema.BoolAttribute{
				Description: "Is yubikey two factor authentication allowed?",
				Computed:    true,
				Optional:    true,
				PlanModifiers: []planmodifier.Bool{
					boolplanmodifier.UseStateForUnknown(),
				},
			},
			"allowed_2fa_method_email": schema.BoolAttribute{
				Description: "Is OTP via email two factor authentication allowed?",
				Computed:    true,
				Optional:    true,
				PlanModifiers: []planmodifier.Bool{
					boolplanmodifier.UseStateForUnknown(),
				},
			},
			"allowed_2fa_method_static": schema.BoolAttribute{
				Description: "Is OTP via static codes for two factor authentication allowed?",
				Computed:    true,
				Optional:    true,
				PlanModifiers: []planmodifier.Bool{
					boolplanmodifier.UseStateForUnknown(),
				},
			},
			"allowed_2fa_method_bypass_for_ftp_sftp_dav": schema.BoolAttribute{
				Description: "Are users allowed to configure their two factor authentication to be bypassed for FTP/SFTP/WebDAV?",
				Computed:    true,
				Optional:    true,
				PlanModifiers: []planmodifier.Bool{
					boolplanmodifier.UseStateForUnknown(),
				},
			},
			"admins_bypass_locked_subfolders": schema.BoolAttribute{
				Description: "Allow admins to bypass the locked subfolders setting.",
				Computed:    true,
				Optional:    true,
				PlanModifiers: []planmodifier.Bool{
					boolplanmodifier.UseStateForUnknown(),
				},
			},
			"allow_bundle_names": schema.BoolAttribute{
				Description: "Are manual Bundle names allowed?",
				Computed:    true,
				Optional:    true,
				PlanModifiers: []planmodifier.Bool{
					boolplanmodifier.UseStateForUnknown(),
				},
			},
			"allowed_countries": schema.StringAttribute{
				Description: "Comma separated list of allowed Country codes",
				Computed:    true,
				Optional:    true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
			},
			"allowed_ips": schema.StringAttribute{
				Description: "List of allowed IP addresses",
				Computed:    true,
				Optional:    true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
			},
			"always_mkdir_parents": schema.BoolAttribute{
				Description: "Create parent directories if they do not exist during uploads?  This is primarily used to work around broken upload clients that assume servers will perform this step.",
				Computed:    true,
				Optional:    true,
				PlanModifiers: []planmodifier.Bool{
					boolplanmodifier.UseStateForUnknown(),
				},
			},
			"ask_about_overwrites": schema.BoolAttribute{
				Description: "If false, rename conflicting files instead of asking for overwrite confirmation.  Only applies to web interface.",
				Computed:    true,
				Optional:    true,
				PlanModifiers: []planmodifier.Bool{
					boolplanmodifier.UseStateForUnknown(),
				},
			},
			"bundle_activity_notifications": schema.StringAttribute{
				Description: "Do Bundle owners receive activity notifications?",
				Computed:    true,
				Optional:    true,
				Validators: []validator.String{
					stringvalidator.OneOf("never", "always", "per_bundle_setting"),
				},
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
			},
			"bundle_expiration": schema.Int64Attribute{
				Description: "Site-wide Bundle expiration in days",
				Computed:    true,
				Optional:    true,
				PlanModifiers: []planmodifier.Int64{
					int64planmodifier.UseStateForUnknown(),
				},
			},
			"bundle_not_found_message": schema.StringAttribute{
				Description: "Custom error message to show when bundle is not found.",
				Computed:    true,
				Optional:    true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
			},
			"bundle_password_required": schema.BoolAttribute{
				Description: "Do Bundles require password protection?",
				Computed:    true,
				Optional:    true,
				PlanModifiers: []planmodifier.Bool{
					boolplanmodifier.UseStateForUnknown(),
				},
			},
			"bundle_recipient_blacklist_domains": schema.ListAttribute{
				Description: "List of email domains to disallow when entering a Bundle/Inbox recipients",
				Computed:    true,
				Optional:    true,
				ElementType: types.StringType,
				PlanModifiers: []planmodifier.List{
					listplanmodifier.UseStateForUnknown(),
				},
			},
			"bundle_recipient_blacklist_free_email_domains": schema.BoolAttribute{
				Description: "Disallow free email domains for Bundle/Inbox recipients?",
				Computed:    true,
				Optional:    true,
				PlanModifiers: []planmodifier.Bool{
					boolplanmodifier.UseStateForUnknown(),
				},
			},
			"bundle_registration_notifications": schema.StringAttribute{
				Description: "Do Bundle owners receive registration notification?",
				Computed:    true,
				Optional:    true,
				Validators: []validator.String{
					stringvalidator.OneOf("never", "always", "per_bundle_setting"),
				},
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
			},
			"bundle_require_registration": schema.BoolAttribute{
				Description: "Do Bundles require registration?",
				Computed:    true,
				Optional:    true,
				PlanModifiers: []planmodifier.Bool{
					boolplanmodifier.UseStateForUnknown(),
				},
			},
			"bundle_require_share_recipient": schema.BoolAttribute{
				Description: "Do Bundles require recipients for sharing?",
				Computed:    true,
				Optional:    true,
				PlanModifiers: []planmodifier.Bool{
					boolplanmodifier.UseStateForUnknown(),
				},
			},
			"bundle_require_note": schema.BoolAttribute{
				Description: "Do Bundles require internal notes?",
				Computed:    true,
				Optional:    true,
				PlanModifiers: []planmodifier.Bool{
					boolplanmodifier.UseStateForUnknown(),
				},
			},
			"bundle_send_shared_receipts": schema.BoolAttribute{
				Description: "Do Bundle creators receive receipts of invitations?",
				Computed:    true,
				Optional:    true,
				PlanModifiers: []planmodifier.Bool{
					boolplanmodifier.UseStateForUnknown(),
				},
			},
			"bundle_upload_receipt_notifications": schema.StringAttribute{
				Description: "Do Bundle uploaders receive upload confirmation notifications?",
				Computed:    true,
				Optional:    true,
				Validators: []validator.String{
					stringvalidator.OneOf("never", "always", "per_bundle_setting"),
				},
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
			},
			"bundle_watermark_value": schema.DynamicAttribute{
				Description: "Preview watermark settings applied to all bundle items. Uses the same keys as Behavior.value",
				Computed:    true,
				Optional:    true,
				PlanModifiers: []planmodifier.Dynamic{
					dynamicplanmodifier.UseStateForUnknown(),
				},
			},
			"calculate_file_checksums_crc32": schema.BoolAttribute{
				Description: "Calculate CRC32 checksums for files?",
				Computed:    true,
				Optional:    true,
				PlanModifiers: []planmodifier.Bool{
					boolplanmodifier.UseStateForUnknown(),
				},
			},
			"calculate_file_checksums_md5": schema.BoolAttribute{
				Description: "Calculate MD5 checksums for files?",
				Computed:    true,
				Optional:    true,
				PlanModifiers: []planmodifier.Bool{
					boolplanmodifier.UseStateForUnknown(),
				},
			},
			"calculate_file_checksums_sha1": schema.BoolAttribute{
				Description: "Calculate SHA1 checksums for files?",
				Computed:    true,
				Optional:    true,
				PlanModifiers: []planmodifier.Bool{
					boolplanmodifier.UseStateForUnknown(),
				},
			},
			"calculate_file_checksums_sha256": schema.BoolAttribute{
				Description: "Calculate SHA256 checksums for files?",
				Computed:    true,
				Optional:    true,
				PlanModifiers: []planmodifier.Bool{
					boolplanmodifier.UseStateForUnknown(),
				},
			},
			"uploads_via_email_authentication": schema.BoolAttribute{
				Description: "Do incoming emails in the Inboxes require checking for SPF/DKIM/DMARC?",
				Computed:    true,
				Optional:    true,
				PlanModifiers: []planmodifier.Bool{
					boolplanmodifier.UseStateForUnknown(),
				},
			},
			"color2_left": schema.StringAttribute{
				Description: "Page link and button color",
				Computed:    true,
				Optional:    true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
			},
			"color2_link": schema.StringAttribute{
				Description: "Top bar link color",
				Computed:    true,
				Optional:    true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
			},
			"color2_text": schema.StringAttribute{
				Description: "Page link and button color",
				Computed:    true,
				Optional:    true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
			},
			"color2_top": schema.StringAttribute{
				Description: "Top bar background color",
				Computed:    true,
				Optional:    true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
			},
			"color2_top_text": schema.StringAttribute{
				Description: "Top bar text color",
				Computed:    true,
				Optional:    true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
			},
			"custom_namespace": schema.BoolAttribute{
				Description: "Is this site using a custom namespace for users?",
				Computed:    true,
				Optional:    true,
				PlanModifiers: []planmodifier.Bool{
					boolplanmodifier.UseStateForUnknown(),
				},
			},
			"dav_enabled": schema.BoolAttribute{
				Description: "Is WebDAV enabled?",
				Computed:    true,
				Optional:    true,
				PlanModifiers: []planmodifier.Bool{
					boolplanmodifier.UseStateForUnknown(),
				},
			},
			"dav_user_root_enabled": schema.BoolAttribute{
				Description: "Use user FTP roots also for WebDAV?",
				Computed:    true,
				Optional:    true,
				PlanModifiers: []planmodifier.Bool{
					boolplanmodifier.UseStateForUnknown(),
				},
			},
			"days_before_deleting_disabled_users": schema.Int64Attribute{
				Description: "Number of days to keep disabled users before deleting them. If set to 0, disabled users will not be deleted.",
				Computed:    true,
				Optional:    true,
				PlanModifiers: []planmodifier.Int64{
					int64planmodifier.UseStateForUnknown(),
				},
			},
			"days_to_retain_backups": schema.Int64Attribute{
				Description: "Number of days to keep deleted files",
				Computed:    true,
				Optional:    true,
				PlanModifiers: []planmodifier.Int64{
					int64planmodifier.UseStateForUnknown(),
				},
			},
			"document_edits_in_bundle_allowed": schema.BoolAttribute{
				Description: "If true, allow public viewers of Bundles with full permissions to use document editing integrations.",
				Computed:    true,
				Optional:    true,
				PlanModifiers: []planmodifier.Bool{
					boolplanmodifier.UseStateForUnknown(),
				},
			},
			"default_time_zone": schema.StringAttribute{
				Description: "Site default time zone",
				Computed:    true,
				Optional:    true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
			},
			"desktop_app": schema.BoolAttribute{
				Description: "Is the desktop app enabled?",
				Computed:    true,
				Optional:    true,
				PlanModifiers: []planmodifier.Bool{
					boolplanmodifier.UseStateForUnknown(),
				},
			},
			"desktop_app_session_ip_pinning": schema.BoolAttribute{
				Description: "Is desktop app session IP pinning enabled?",
				Computed:    true,
				Optional:    true,
				PlanModifiers: []planmodifier.Bool{
					boolplanmodifier.UseStateForUnknown(),
				},
			},
			"desktop_app_session_lifetime": schema.Int64Attribute{
				Description: "Desktop app session lifetime (in hours)",
				Computed:    true,
				Optional:    true,
				PlanModifiers: []planmodifier.Int64{
					int64planmodifier.UseStateForUnknown(),
				},
			},
			"legacy_checksums_mode": schema.BoolAttribute{
				Description: "Use legacy checksums mode?",
				Computed:    true,
				Optional:    true,
				PlanModifiers: []planmodifier.Bool{
					boolplanmodifier.UseStateForUnknown(),
				},
			},
			"mobile_app": schema.BoolAttribute{
				Description: "Is the mobile app enabled?",
				Computed:    true,
				Optional:    true,
				PlanModifiers: []planmodifier.Bool{
					boolplanmodifier.UseStateForUnknown(),
				},
			},
			"mobile_app_session_ip_pinning": schema.BoolAttribute{
				Description: "Is mobile app session IP pinning enabled?",
				Computed:    true,
				Optional:    true,
				PlanModifiers: []planmodifier.Bool{
					boolplanmodifier.UseStateForUnknown(),
				},
			},
			"mobile_app_session_lifetime": schema.Int64Attribute{
				Description: "Mobile app session lifetime (in hours)",
				Computed:    true,
				Optional:    true,
				PlanModifiers: []planmodifier.Int64{
					int64planmodifier.UseStateForUnknown(),
				},
			},
			"disallowed_countries": schema.StringAttribute{
				Description: "Comma separated list of disallowed Country codes",
				Computed:    true,
				Optional:    true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
			},
			"disable_files_certificate_generation": schema.BoolAttribute{
				Description: "If set, Files.com will not set the CAA records required to generate future SSL certificates for this domain.",
				Computed:    true,
				Optional:    true,
				PlanModifiers: []planmodifier.Bool{
					boolplanmodifier.UseStateForUnknown(),
				},
			},
			"disable_password_reset": schema.BoolAttribute{
				Description: "Is password reset disabled?",
				Computed:    true,
				Optional:    true,
				PlanModifiers: []planmodifier.Bool{
					boolplanmodifier.UseStateForUnknown(),
				},
			},
			"domain": schema.StringAttribute{
				Description: "Custom domain",
				Computed:    true,
				Optional:    true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
			},
			"domain_hsts_header": schema.BoolAttribute{
				Description: "Send HSTS (HTTP Strict Transport Security) header when visitors access the site via a custom domain?",
				Computed:    true,
				Optional:    true,
				PlanModifiers: []planmodifier.Bool{
					boolplanmodifier.UseStateForUnknown(),
				},
			},
			"domain_letsencrypt_chain": schema.StringAttribute{
				Description: "Letsencrypt chain to use when registering SSL Certificate for domain.",
				Computed:    true,
				Optional:    true,
				Validators: []validator.String{
					stringvalidator.OneOf("default", "isrg_root_x1", "dst_root_ca_x3"),
				},
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
			},
			"email": schema.StringAttribute{
				Description: "Main email for this site",
				Computed:    true,
				Optional:    true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
			},
			"ftp_enabled": schema.BoolAttribute{
				Description: "Is FTP enabled?",
				Computed:    true,
				Optional:    true,
				PlanModifiers: []planmodifier.Bool{
					boolplanmodifier.UseStateForUnknown(),
				},
			},
			"reply_to_email": schema.StringAttribute{
				Description: "Reply-to email for this site",
				Computed:    true,
				Optional:    true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
			},
			"non_sso_groups_allowed": schema.BoolAttribute{
				Description: "If true, groups can be manually created / modified / deleted by Site Admins. Otherwise, groups can only be managed via your SSO provider.",
				Computed:    true,
				Optional:    true,
				PlanModifiers: []planmodifier.Bool{
					boolplanmodifier.UseStateForUnknown(),
				},
			},
			"non_sso_users_allowed": schema.BoolAttribute{
				Description: "If true, users can be manually created / modified / deleted by Site Admins. Otherwise, users can only be managed via your SSO provider.",
				Computed:    true,
				Optional:    true,
				PlanModifiers: []planmodifier.Bool{
					boolplanmodifier.UseStateForUnknown(),
				},
			},
			"folder_permissions_groups_only": schema.BoolAttribute{
				Description: "If true, permissions for this site must be bound to a group (not a user).",
				Computed:    true,
				Optional:    true,
				PlanModifiers: []planmodifier.Bool{
					boolplanmodifier.UseStateForUnknown(),
				},
			},
			"include_password_in_welcome_email": schema.BoolAttribute{
				Description: "Include password in emails to new users?",
				Computed:    true,
				Optional:    true,
				PlanModifiers: []planmodifier.Bool{
					boolplanmodifier.UseStateForUnknown(),
				},
			},
			"language": schema.StringAttribute{
				Description: "Site default language",
				Computed:    true,
				Optional:    true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
			},
			"ldap_base_dn": schema.StringAttribute{
				Description: "Base DN for looking up users in LDAP server",
				Computed:    true,
				Optional:    true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
			},
			"ldap_domain": schema.StringAttribute{
				Description: "Domain name that will be appended to usernames",
				Computed:    true,
				Optional:    true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
			},
			"ldap_enabled": schema.BoolAttribute{
				Description: "Main LDAP setting: is LDAP enabled?",
				Computed:    true,
				Optional:    true,
				PlanModifiers: []planmodifier.Bool{
					boolplanmodifier.UseStateForUnknown(),
				},
			},
			"ldap_group_action": schema.StringAttribute{
				Description: "Should we sync groups from LDAP server?",
				Computed:    true,
				Optional:    true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
			},
			"ldap_group_exclusion": schema.StringAttribute{
				Description: "Comma or newline separated list of group names (with optional wildcards) to exclude when syncing.",
				Computed:    true,
				Optional:    true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
			},
			"ldap_group_inclusion": schema.StringAttribute{
				Description: "Comma or newline separated list of group names (with optional wildcards) to include when syncing.",
				Computed:    true,
				Optional:    true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
			},
			"ldap_host": schema.StringAttribute{
				Description: "LDAP host",
				Computed:    true,
				Optional:    true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
			},
			"ldap_host_2": schema.StringAttribute{
				Description: "LDAP backup host",
				Computed:    true,
				Optional:    true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
			},
			"ldap_host_3": schema.StringAttribute{
				Description: "LDAP backup host",
				Computed:    true,
				Optional:    true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
			},
			"ldap_port": schema.Int64Attribute{
				Description: "LDAP port",
				Computed:    true,
				Optional:    true,
				PlanModifiers: []planmodifier.Int64{
					int64planmodifier.UseStateForUnknown(),
				},
			},
			"ldap_secure": schema.BoolAttribute{
				Description: "Use secure LDAP?",
				Computed:    true,
				Optional:    true,
				PlanModifiers: []planmodifier.Bool{
					boolplanmodifier.UseStateForUnknown(),
				},
			},
			"ldap_type": schema.StringAttribute{
				Description: "LDAP type",
				Computed:    true,
				Optional:    true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
			},
			"ldap_user_action": schema.StringAttribute{
				Description: "Should we sync users from LDAP server?",
				Computed:    true,
				Optional:    true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
			},
			"ldap_user_include_groups": schema.StringAttribute{
				Description: "Comma or newline separated list of group names (with optional wildcards) - if provided, only users in these groups will be added or synced.",
				Computed:    true,
				Optional:    true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
			},
			"ldap_username": schema.StringAttribute{
				Description: "Username for signing in to LDAP server.",
				Computed:    true,
				Optional:    true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
			},
			"ldap_username_field": schema.StringAttribute{
				Description: "LDAP username field",
				Computed:    true,
				Optional:    true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
			},
			"login_help_text": schema.StringAttribute{
				Description: "Login help text",
				Computed:    true,
				Optional:    true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
			},
			"max_prior_passwords": schema.Int64Attribute{
				Description: "Number of prior passwords to disallow",
				Computed:    true,
				Optional:    true,
				PlanModifiers: []planmodifier.Int64{
					int64planmodifier.UseStateForUnknown(),
				},
			},
			"motd_text": schema.StringAttribute{
				Description: "A message to show users when they connect via FTP or SFTP.",
				Computed:    true,
				Optional:    true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
			},
			"motd_use_for_ftp": schema.BoolAttribute{
				Description: "Show message to users connecting via FTP",
				Computed:    true,
				Optional:    true,
				PlanModifiers: []planmodifier.Bool{
					boolplanmodifier.UseStateForUnknown(),
				},
			},
			"motd_use_for_sftp": schema.BoolAttribute{
				Description: "Show message to users connecting via SFTP",
				Computed:    true,
				Optional:    true,
				PlanModifiers: []planmodifier.Bool{
					boolplanmodifier.UseStateForUnknown(),
				},
			},
			"office_integration_available": schema.BoolAttribute{
				Description: "If true, allows users to use a document editing integration.",
				Computed:    true,
				Optional:    true,
				PlanModifiers: []planmodifier.Bool{
					boolplanmodifier.UseStateForUnknown(),
				},
			},
			"office_integration_type": schema.StringAttribute{
				Description: "Which document editing integration to support. Files.com Editor or Microsoft Office for the Web.",
				Computed:    true,
				Optional:    true,
				Validators: []validator.String{
					stringvalidator.OneOf("only_office", "office_365", "disabled"),
				},
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
			},
			"opt_out_global": schema.BoolAttribute{
				Description: "Use servers in the USA only?",
				Computed:    true,
				Optional:    true,
				PlanModifiers: []planmodifier.Bool{
					boolplanmodifier.UseStateForUnknown(),
				},
			},
			"password_min_length": schema.Int64Attribute{
				Description: "Shortest password length for users",
				Computed:    true,
				Optional:    true,
				PlanModifiers: []planmodifier.Int64{
					int64planmodifier.UseStateForUnknown(),
				},
			},
			"password_require_letter": schema.BoolAttribute{
				Description: "Require a letter in passwords?",
				Computed:    true,
				Optional:    true,
				PlanModifiers: []planmodifier.Bool{
					boolplanmodifier.UseStateForUnknown(),
				},
			},
			"password_require_mixed": schema.BoolAttribute{
				Description: "Require lower and upper case letters in passwords?",
				Computed:    true,
				Optional:    true,
				PlanModifiers: []planmodifier.Bool{
					boolplanmodifier.UseStateForUnknown(),
				},
			},
			"password_require_number": schema.BoolAttribute{
				Description: "Require a number in passwords?",
				Computed:    true,
				Optional:    true,
				PlanModifiers: []planmodifier.Bool{
					boolplanmodifier.UseStateForUnknown(),
				},
			},
			"password_require_special": schema.BoolAttribute{
				Description: "Require special characters in password?",
				Computed:    true,
				Optional:    true,
				PlanModifiers: []planmodifier.Bool{
					boolplanmodifier.UseStateForUnknown(),
				},
			},
			"password_require_unbreached": schema.BoolAttribute{
				Description: "Require passwords that have not been previously breached? (see https://haveibeenpwned.com/)",
				Computed:    true,
				Optional:    true,
				PlanModifiers: []planmodifier.Bool{
					boolplanmodifier.UseStateForUnknown(),
				},
			},
			"password_requirements_apply_to_bundles": schema.BoolAttribute{
				Description: "Require bundles' passwords, and passwords for other items (inboxes, public shares, etc.) to conform to the same requirements as users' passwords?",
				Computed:    true,
				Optional:    true,
				PlanModifiers: []planmodifier.Bool{
					boolplanmodifier.UseStateForUnknown(),
				},
			},
			"password_validity_days": schema.Int64Attribute{
				Description: "Number of days password is valid",
				Computed:    true,
				Optional:    true,
				PlanModifiers: []planmodifier.Int64{
					int64planmodifier.UseStateForUnknown(),
				},
			},
			"pin_all_remote_servers_to_site_region": schema.BoolAttribute{
				Description: "If true, we will ensure that all internal communications with any remote server are made through the primary region of the site. This setting overrides individual remote server settings.",
				Computed:    true,
				Optional:    true,
				PlanModifiers: []planmodifier.Bool{
					boolplanmodifier.UseStateForUnknown(),
				},
			},
			"prevent_root_permissions_for_non_site_admins": schema.BoolAttribute{
				Description: "If true, we will prevent non-administrators from receiving any permissions directly on the root folder.  This is commonly used to prevent the accidental application of permissions.",
				Computed:    true,
				Optional:    true,
				PlanModifiers: []planmodifier.Bool{
					boolplanmodifier.UseStateForUnknown(),
				},
			},
			"protocol_access_groups_only": schema.BoolAttribute{
				Description: "If true, protocol access permissions on users will be ignored, and only protocol access permissions set on Groups will be honored.  Make sure that your current user is a member of a group with API permission when changing this value to avoid locking yourself out of your site.",
				Computed:    true,
				Optional:    true,
				PlanModifiers: []planmodifier.Bool{
					boolplanmodifier.UseStateForUnknown(),
				},
			},
			"require_2fa": schema.BoolAttribute{
				Description: "Require two-factor authentication for all users?",
				Computed:    true,
				Optional:    true,
				PlanModifiers: []planmodifier.Bool{
					boolplanmodifier.UseStateForUnknown(),
				},
			},
			"revoke_bundle_access_on_disable_or_delete": schema.BoolAttribute{
				Description: "Auto-removes bundles for disabled/deleted users and enforces bundle expiry within user access period.",
				Computed:    true,
				Optional:    true,
				PlanModifiers: []planmodifier.Bool{
					boolplanmodifier.UseStateForUnknown(),
				},
			},
			"require_2fa_user_type": schema.StringAttribute{
				Description: "What type of user is required to use two-factor authentication (when require_2fa is set to `true` for this site)?",
				Computed:    true,
				Optional:    true,
				Validators: []validator.String{
					stringvalidator.OneOf("all", "folder_and_site_admins", "site_admins"),
				},
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
			},
			"require_logout_from_bundles_and_inboxes": schema.BoolAttribute{
				Description: "If true, we will hide the 'Remember Me' box on Inbox and Bundle registration pages, requiring that the user logout and log back in every time they visit the page.",
				Computed:    true,
				Optional:    true,
				PlanModifiers: []planmodifier.Bool{
					boolplanmodifier.UseStateForUnknown(),
				},
			},
			"sftp_enabled": schema.BoolAttribute{
				Description: "Is SFTP enabled?",
				Computed:    true,
				Optional:    true,
				PlanModifiers: []planmodifier.Bool{
					boolplanmodifier.UseStateForUnknown(),
				},
			},
			"sftp_host_key_type": schema.StringAttribute{
				Description: "Sftp Host Key Type",
				Computed:    true,
				Optional:    true,
				Validators: []validator.String{
					stringvalidator.OneOf("default", "exavault", "custom", "smartfile"),
				},
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
			},
			"active_sftp_host_key_id": schema.Int64Attribute{
				Description: "Id of the currently selected custom SFTP Host Key",
				Computed:    true,
				Optional:    true,
				PlanModifiers: []planmodifier.Int64{
					int64planmodifier.UseStateForUnknown(),
				},
			},
			"sftp_insecure_ciphers": schema.BoolAttribute{
				Description: "If true, we will allow weak and known insecure ciphers to be used for SFTP connections.  Enabling this setting severely weakens the security of your site and it is not recommend, except as a last resort for compatibility.",
				Computed:    true,
				Optional:    true,
				PlanModifiers: []planmodifier.Bool{
					boolplanmodifier.UseStateForUnknown(),
				},
			},
			"sftp_insecure_diffie_hellman": schema.BoolAttribute{
				Description: "If true, we will allow weak Diffie Hellman parameters to be used within ciphers for SFTP that are otherwise on our secure list.  This has the effect of making the cipher weaker than our normal threshold for security, but is required to support certain legacy or broken SSH and MFT clients.  Enabling this weakens security, but not nearly as much as enabling the full `sftp_insecure_ciphers` option.",
				Computed:    true,
				Optional:    true,
				PlanModifiers: []planmodifier.Bool{
					boolplanmodifier.UseStateForUnknown(),
				},
			},
			"sftp_user_root_enabled": schema.BoolAttribute{
				Description: "Use user FTP roots also for SFTP?",
				Computed:    true,
				Optional:    true,
				PlanModifiers: []planmodifier.Bool{
					boolplanmodifier.UseStateForUnknown(),
				},
			},
			"sharing_enabled": schema.BoolAttribute{
				Description: "Allow bundle creation",
				Computed:    true,
				Optional:    true,
				PlanModifiers: []planmodifier.Bool{
					boolplanmodifier.UseStateForUnknown(),
				},
			},
			"show_user_notifications_log_in_link": schema.BoolAttribute{
				Description: "Show log in link in user notifications?",
				Computed:    true,
				Optional:    true,
				PlanModifiers: []planmodifier.Bool{
					boolplanmodifier.UseStateForUnknown(),
				},
			},
			"show_request_access_link": schema.BoolAttribute{
				Description: "Show request access link for users without access?  Currently unused.",
				Computed:    true,
				Optional:    true,
				PlanModifiers: []planmodifier.Bool{
					boolplanmodifier.UseStateForUnknown(),
				},
			},
			"site_footer": schema.StringAttribute{
				Description: "Custom site footer text",
				Computed:    true,
				Optional:    true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
			},
			"site_header": schema.StringAttribute{
				Description: "Custom site header text",
				Computed:    true,
				Optional:    true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
			},
			"smtp_address": schema.StringAttribute{
				Description: "SMTP server hostname or IP",
				Computed:    true,
				Optional:    true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
			},
			"smtp_authentication": schema.StringAttribute{
				Description: "SMTP server authentication type",
				Computed:    true,
				Optional:    true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
			},
			"smtp_from": schema.StringAttribute{
				Description: "From address to use when mailing through custom SMTP",
				Computed:    true,
				Optional:    true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
			},
			"smtp_port": schema.Int64Attribute{
				Description: "SMTP server port",
				Computed:    true,
				Optional:    true,
				PlanModifiers: []planmodifier.Int64{
					int64planmodifier.UseStateForUnknown(),
				},
			},
			"smtp_username": schema.StringAttribute{
				Description: "SMTP server username",
				Computed:    true,
				Optional:    true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
			},
			"session_expiry": schema.StringAttribute{
				Description: "Session expiry in hours",
				Computed:    true,
				Optional:    true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
			},
			"session_expiry_minutes": schema.Int64Attribute{
				Description: "Session expiry in minutes",
				Computed:    true,
				Optional:    true,
				PlanModifiers: []planmodifier.Int64{
					int64planmodifier.UseStateForUnknown(),
				},
			},
			"snapshot_sharing_enabled": schema.BoolAttribute{
				Description: "Allow snapshot share links creation",
				Computed:    true,
				Optional:    true,
				PlanModifiers: []planmodifier.Bool{
					boolplanmodifier.UseStateForUnknown(),
				},
			},
			"ssl_required": schema.BoolAttribute{
				Description: "Is SSL required?  Disabling this is insecure.",
				Computed:    true,
				Optional:    true,
				PlanModifiers: []planmodifier.Bool{
					boolplanmodifier.UseStateForUnknown(),
				},
			},
			"subdomain": schema.StringAttribute{
				Description: "Site subdomain",
				Computed:    true,
				Optional:    true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
			},
			"tls_disabled": schema.BoolAttribute{
				Description: "DO NOT ENABLE. This setting allows TLSv1.0 and TLSv1.1 to be used on your site.  We intend to remove this capability entirely in early 2024.  If set, the `sftp_insecure_ciphers` flag will be automatically set to true.",
				Computed:    true,
				Optional:    true,
				PlanModifiers: []planmodifier.Bool{
					boolplanmodifier.UseStateForUnknown(),
				},
			},
			"use_dedicated_ips_for_smtp": schema.BoolAttribute{
				Description: "If using custom SMTP, should we use dedicated IPs to deliver emails?",
				Computed:    true,
				Optional:    true,
				PlanModifiers: []planmodifier.Bool{
					boolplanmodifier.UseStateForUnknown(),
				},
			},
			"use_provided_modified_at": schema.BoolAttribute{
				Description: "Allow uploaders to set `provided_modified_at` for uploaded files?",
				Computed:    true,
				Optional:    true,
				PlanModifiers: []planmodifier.Bool{
					boolplanmodifier.UseStateForUnknown(),
				},
			},
			"user_lockout": schema.BoolAttribute{
				Description: "Will users be locked out after incorrect login attempts?",
				Computed:    true,
				Optional:    true,
				PlanModifiers: []planmodifier.Bool{
					boolplanmodifier.UseStateForUnknown(),
				},
			},
			"user_lockout_lock_period": schema.Int64Attribute{
				Description: "How many hours to lock user out for failed password?",
				Computed:    true,
				Optional:    true,
				PlanModifiers: []planmodifier.Int64{
					int64planmodifier.UseStateForUnknown(),
				},
			},
			"user_lockout_tries": schema.Int64Attribute{
				Description: "Number of login tries within `user_lockout_within` hours before users are locked out",
				Computed:    true,
				Optional:    true,
				PlanModifiers: []planmodifier.Int64{
					int64planmodifier.UseStateForUnknown(),
				},
			},
			"user_lockout_within": schema.Int64Attribute{
				Description: "Number of hours for user lockout window",
				Computed:    true,
				Optional:    true,
				PlanModifiers: []planmodifier.Int64{
					int64planmodifier.UseStateForUnknown(),
				},
			},
			"user_requests_enabled": schema.BoolAttribute{
				Description: "Enable User Requests feature",
				Computed:    true,
				Optional:    true,
				PlanModifiers: []planmodifier.Bool{
					boolplanmodifier.UseStateForUnknown(),
				},
			},
			"user_requests_notify_admins": schema.BoolAttribute{
				Description: "Send email to site admins when a user request is received?",
				Computed:    true,
				Optional:    true,
				PlanModifiers: []planmodifier.Bool{
					boolplanmodifier.UseStateForUnknown(),
				},
			},
			"users_can_create_api_keys": schema.BoolAttribute{
				Description: "Allow users to create their own API keys?",
				Computed:    true,
				Optional:    true,
				PlanModifiers: []planmodifier.Bool{
					boolplanmodifier.UseStateForUnknown(),
				},
			},
			"users_can_create_ssh_keys": schema.BoolAttribute{
				Description: "Allow users to create their own SSH keys?",
				Computed:    true,
				Optional:    true,
				PlanModifiers: []planmodifier.Bool{
					boolplanmodifier.UseStateForUnknown(),
				},
			},
			"welcome_custom_text": schema.StringAttribute{
				Description: "Custom text send in user welcome email",
				Computed:    true,
				Optional:    true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
			},
			"welcome_email_cc": schema.StringAttribute{
				Description: "Include this email in welcome emails if enabled",
				Computed:    true,
				Optional:    true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
			},
			"welcome_email_subject": schema.StringAttribute{
				Description: "Include this email subject in welcome emails if enabled",
				Computed:    true,
				Optional:    true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
			},
			"welcome_email_enabled": schema.BoolAttribute{
				Description: "Will the welcome email be sent to new users?",
				Computed:    true,
				Optional:    true,
				PlanModifiers: []planmodifier.Bool{
					boolplanmodifier.UseStateForUnknown(),
				},
			},
			"welcome_screen": schema.StringAttribute{
				Description: "Does the welcome screen appear?",
				Computed:    true,
				Optional:    true,
				Validators: []validator.String{
					stringvalidator.OneOf("enabled", "hidden", "disabled"),
				},
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
			},
			"windows_mode_ftp": schema.BoolAttribute{
				Description: "Does FTP user Windows emulation mode?",
				Computed:    true,
				Optional:    true,
				PlanModifiers: []planmodifier.Bool{
					boolplanmodifier.UseStateForUnknown(),
				},
			},
			"group_admins_can_set_user_password": schema.BoolAttribute{
				Description: "Allow group admins set password authentication method",
				Computed:    true,
				Optional:    true,
				PlanModifiers: []planmodifier.Bool{
					boolplanmodifier.UseStateForUnknown(),
				},
			},
			"id": schema.Int64Attribute{
				Description: "Site Id",
				Computed:    true,
				PlanModifiers: []planmodifier.Int64{
					int64planmodifier.UseStateForUnknown(),
				},
			},
			"admin_user_id": schema.Int64Attribute{
				Description: "User ID for the main site administrator",
				Computed:    true,
			},
			"bundle_watermark_attachment": schema.StringAttribute{
				Description: "Preview watermark image applied to all bundle items.",
				Computed:    true,
			},
			"contact_name": schema.StringAttribute{
				Description: "Site main contact name",
				Computed:    true,
			},
			"created_at": schema.StringAttribute{
				Description: "Time this site was created",
				Computed:    true,
			},
			"currency": schema.StringAttribute{
				Description: "Preferred currency",
				Computed:    true,
			},
			"disable_notifications": schema.BoolAttribute{
				Description: "Are notifications disabled?",
				Computed:    true,
			},
			"hipaa": schema.BoolAttribute{
				Description: "Is there a signed HIPAA BAA between Files.com and this site?",
				Computed:    true,
			},
			"icon128": schema.StringAttribute{
				Description: "Branded icon 128x128",
				Computed:    true,
			},
			"icon16": schema.StringAttribute{
				Description: "Branded icon 16x16",
				Computed:    true,
			},
			"icon32": schema.StringAttribute{
				Description: "Branded icon 32x32",
				Computed:    true,
			},
			"icon48": schema.StringAttribute{
				Description: "Branded icon 48x48",
				Computed:    true,
			},
			"immutable_files_set_at": schema.StringAttribute{
				Description: "Can files be modified?",
				Computed:    true,
			},
			"logo": schema.StringAttribute{
				Description: "Branded logo",
				Computed:    true,
			},
			"login_page_background_image": schema.StringAttribute{
				Description: "Branded login page background",
				Computed:    true,
			},
			"next_billing_amount": schema.StringAttribute{
				Description: "Next billing amount",
				Computed:    true,
			},
			"next_billing_date": schema.StringAttribute{
				Description: "Next billing date",
				Computed:    true,
			},
			"oncehub_link": schema.StringAttribute{
				Description: "Link to scheduling a meeting with our Sales team",
				Computed:    true,
			},
			"overdue": schema.BoolAttribute{
				Description: "Is this site's billing overdue?",
				Computed:    true,
			},
			"phone": schema.StringAttribute{
				Description: "Site phone number",
				Computed:    true,
			},
			"require_2fa_stop_time": schema.StringAttribute{
				Description: "If set, requirement for two-factor authentication has been scheduled to end on this date-time.",
				Computed:    true,
			},
			"session": schema.StringAttribute{
				Description: "Current session",
				Computed:    true,
			},
			"switch_to_plan_date": schema.StringAttribute{
				Description: "If switching plans, when does the new plan take effect?",
				Computed:    true,
			},
			"trial_days_left": schema.Int64Attribute{
				Description: "Number of days left in trial",
				Computed:    true,
			},
			"trial_until": schema.StringAttribute{
				Description: "When does this Site trial expire?",
				Computed:    true,
			},
			"user": schema.StringAttribute{
				Description: "User of current session",
				Computed:    true,
			},
		},
	}
}

func (r *siteResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	resp.Diagnostics.AddError(
		"Resource Create Not Implemented",
		"This resource does not support creation. Please import the resource.",
	)
}

func (r *siteResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	var state siteResourceModel
	diags := req.State.Get(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	site, err := r.client.Get(files_sdk.WithContext(ctx))
	if err != nil {
		if files_sdk.IsNotExist(err) {
			resp.State.RemoveResource(ctx)
			return
		}

		resp.Diagnostics.AddError(
			"Error Reading Files Site",
			"Could not read site: "+err.Error(),
		)
		return
	}

	diags = r.populateResourceModel(ctx, site, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	diags = resp.State.Set(ctx, state)
	resp.Diagnostics.Append(diags...)
}

func (r *siteResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	var plan siteResourceModel
	diags := req.Plan.Get(ctx, &plan)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	paramsSiteUpdate := files_sdk.SiteUpdateParams{}
	paramsSiteUpdate.Name = plan.Name.ValueString()
	paramsSiteUpdate.Subdomain = plan.Subdomain.ValueString()
	paramsSiteUpdate.Domain = plan.Domain.ValueString()
	if !plan.DomainHstsHeader.IsNull() && !plan.DomainHstsHeader.IsUnknown() {
		paramsSiteUpdate.DomainHstsHeader = plan.DomainHstsHeader.ValueBoolPointer()
	}
	paramsSiteUpdate.DomainLetsencryptChain = plan.DomainLetsencryptChain.ValueString()
	paramsSiteUpdate.Email = plan.Email.ValueString()
	paramsSiteUpdate.ReplyToEmail = plan.ReplyToEmail.ValueString()
	if !plan.AllowBundleNames.IsNull() && !plan.AllowBundleNames.IsUnknown() {
		paramsSiteUpdate.AllowBundleNames = plan.AllowBundleNames.ValueBoolPointer()
	}
	paramsSiteUpdate.BundleExpiration = plan.BundleExpiration.ValueInt64()
	if !plan.WelcomeEmailEnabled.IsNull() && !plan.WelcomeEmailEnabled.IsUnknown() {
		paramsSiteUpdate.WelcomeEmailEnabled = plan.WelcomeEmailEnabled.ValueBoolPointer()
	}
	if !plan.AskAboutOverwrites.IsNull() && !plan.AskAboutOverwrites.IsUnknown() {
		paramsSiteUpdate.AskAboutOverwrites = plan.AskAboutOverwrites.ValueBoolPointer()
	}
	if !plan.ShowRequestAccessLink.IsNull() && !plan.ShowRequestAccessLink.IsUnknown() {
		paramsSiteUpdate.ShowRequestAccessLink = plan.ShowRequestAccessLink.ValueBoolPointer()
	}
	if !plan.AlwaysMkdirParents.IsNull() && !plan.AlwaysMkdirParents.IsUnknown() {
		paramsSiteUpdate.AlwaysMkdirParents = plan.AlwaysMkdirParents.ValueBoolPointer()
	}
	paramsSiteUpdate.WelcomeEmailCc = plan.WelcomeEmailCc.ValueString()
	paramsSiteUpdate.WelcomeEmailSubject = plan.WelcomeEmailSubject.ValueString()
	paramsSiteUpdate.WelcomeCustomText = plan.WelcomeCustomText.ValueString()
	paramsSiteUpdate.Language = plan.Language.ValueString()
	if !plan.WindowsModeFtp.IsNull() && !plan.WindowsModeFtp.IsUnknown() {
		paramsSiteUpdate.WindowsModeFtp = plan.WindowsModeFtp.ValueBoolPointer()
	}
	paramsSiteUpdate.DefaultTimeZone = plan.DefaultTimeZone.ValueString()
	if !plan.DesktopApp.IsNull() && !plan.DesktopApp.IsUnknown() {
		paramsSiteUpdate.DesktopApp = plan.DesktopApp.ValueBoolPointer()
	}
	if !plan.DesktopAppSessionIpPinning.IsNull() && !plan.DesktopAppSessionIpPinning.IsUnknown() {
		paramsSiteUpdate.DesktopAppSessionIpPinning = plan.DesktopAppSessionIpPinning.ValueBoolPointer()
	}
	paramsSiteUpdate.DesktopAppSessionLifetime = plan.DesktopAppSessionLifetime.ValueInt64()
	if !plan.MobileApp.IsNull() && !plan.MobileApp.IsUnknown() {
		paramsSiteUpdate.MobileApp = plan.MobileApp.ValueBoolPointer()
	}
	if !plan.MobileAppSessionIpPinning.IsNull() && !plan.MobileAppSessionIpPinning.IsUnknown() {
		paramsSiteUpdate.MobileAppSessionIpPinning = plan.MobileAppSessionIpPinning.ValueBoolPointer()
	}
	paramsSiteUpdate.MobileAppSessionLifetime = plan.MobileAppSessionLifetime.ValueInt64()
	if !plan.FolderPermissionsGroupsOnly.IsNull() && !plan.FolderPermissionsGroupsOnly.IsUnknown() {
		paramsSiteUpdate.FolderPermissionsGroupsOnly = plan.FolderPermissionsGroupsOnly.ValueBoolPointer()
	}
	paramsSiteUpdate.WelcomeScreen = plan.WelcomeScreen.ValueString()
	if !plan.OfficeIntegrationAvailable.IsNull() && !plan.OfficeIntegrationAvailable.IsUnknown() {
		paramsSiteUpdate.OfficeIntegrationAvailable = plan.OfficeIntegrationAvailable.ValueBoolPointer()
	}
	paramsSiteUpdate.OfficeIntegrationType = plan.OfficeIntegrationType.ValueString()
	if !plan.PinAllRemoteServersToSiteRegion.IsNull() && !plan.PinAllRemoteServersToSiteRegion.IsUnknown() {
		paramsSiteUpdate.PinAllRemoteServersToSiteRegion = plan.PinAllRemoteServersToSiteRegion.ValueBoolPointer()
	}
	paramsSiteUpdate.MotdText = plan.MotdText.ValueString()
	if !plan.MotdUseForFtp.IsNull() && !plan.MotdUseForFtp.IsUnknown() {
		paramsSiteUpdate.MotdUseForFtp = plan.MotdUseForFtp.ValueBoolPointer()
	}
	if !plan.MotdUseForSftp.IsNull() && !plan.MotdUseForSftp.IsUnknown() {
		paramsSiteUpdate.MotdUseForSftp = plan.MotdUseForSftp.ValueBoolPointer()
	}
	if !plan.AdditionalTextFileTypes.IsNull() && !plan.AdditionalTextFileTypes.IsUnknown() {
		diags = plan.AdditionalTextFileTypes.ElementsAs(ctx, &paramsSiteUpdate.AdditionalTextFileTypes, false)
		resp.Diagnostics.Append(diags...)
	}
	if !plan.BundleRequireNote.IsNull() && !plan.BundleRequireNote.IsUnknown() {
		paramsSiteUpdate.BundleRequireNote = plan.BundleRequireNote.ValueBoolPointer()
	}
	if !plan.BundleSendSharedReceipts.IsNull() && !plan.BundleSendSharedReceipts.IsUnknown() {
		paramsSiteUpdate.BundleSendSharedReceipts = plan.BundleSendSharedReceipts.ValueBoolPointer()
	}
	if !plan.CalculateFileChecksumsCrc32.IsNull() && !plan.CalculateFileChecksumsCrc32.IsUnknown() {
		paramsSiteUpdate.CalculateFileChecksumsCrc32 = plan.CalculateFileChecksumsCrc32.ValueBoolPointer()
	}
	if !plan.CalculateFileChecksumsMd5.IsNull() && !plan.CalculateFileChecksumsMd5.IsUnknown() {
		paramsSiteUpdate.CalculateFileChecksumsMd5 = plan.CalculateFileChecksumsMd5.ValueBoolPointer()
	}
	if !plan.CalculateFileChecksumsSha1.IsNull() && !plan.CalculateFileChecksumsSha1.IsUnknown() {
		paramsSiteUpdate.CalculateFileChecksumsSha1 = plan.CalculateFileChecksumsSha1.ValueBoolPointer()
	}
	if !plan.CalculateFileChecksumsSha256.IsNull() && !plan.CalculateFileChecksumsSha256.IsUnknown() {
		paramsSiteUpdate.CalculateFileChecksumsSha256 = plan.CalculateFileChecksumsSha256.ValueBoolPointer()
	}
	if !plan.LegacyChecksumsMode.IsNull() && !plan.LegacyChecksumsMode.IsUnknown() {
		paramsSiteUpdate.LegacyChecksumsMode = plan.LegacyChecksumsMode.ValueBoolPointer()
	}
	paramsSiteUpdate.SessionExpiry = plan.SessionExpiry.ValueString()
	if !plan.SslRequired.IsNull() && !plan.SslRequired.IsUnknown() {
		paramsSiteUpdate.SslRequired = plan.SslRequired.ValueBoolPointer()
	}
	if !plan.TlsDisabled.IsNull() && !plan.TlsDisabled.IsUnknown() {
		paramsSiteUpdate.TlsDisabled = plan.TlsDisabled.ValueBoolPointer()
	}
	if !plan.SftpInsecureCiphers.IsNull() && !plan.SftpInsecureCiphers.IsUnknown() {
		paramsSiteUpdate.SftpInsecureCiphers = plan.SftpInsecureCiphers.ValueBoolPointer()
	}
	if !plan.SftpInsecureDiffieHellman.IsNull() && !plan.SftpInsecureDiffieHellman.IsUnknown() {
		paramsSiteUpdate.SftpInsecureDiffieHellman = plan.SftpInsecureDiffieHellman.ValueBoolPointer()
	}
	if !plan.DisableFilesCertificateGeneration.IsNull() && !plan.DisableFilesCertificateGeneration.IsUnknown() {
		paramsSiteUpdate.DisableFilesCertificateGeneration = plan.DisableFilesCertificateGeneration.ValueBoolPointer()
	}
	if !plan.UserLockout.IsNull() && !plan.UserLockout.IsUnknown() {
		paramsSiteUpdate.UserLockout = plan.UserLockout.ValueBoolPointer()
	}
	paramsSiteUpdate.UserLockoutTries = plan.UserLockoutTries.ValueInt64()
	paramsSiteUpdate.UserLockoutWithin = plan.UserLockoutWithin.ValueInt64()
	paramsSiteUpdate.UserLockoutLockPeriod = plan.UserLockoutLockPeriod.ValueInt64()
	if !plan.IncludePasswordInWelcomeEmail.IsNull() && !plan.IncludePasswordInWelcomeEmail.IsUnknown() {
		paramsSiteUpdate.IncludePasswordInWelcomeEmail = plan.IncludePasswordInWelcomeEmail.ValueBoolPointer()
	}
	paramsSiteUpdate.AllowedCountries = plan.AllowedCountries.ValueString()
	paramsSiteUpdate.AllowedIps = plan.AllowedIps.ValueString()
	paramsSiteUpdate.DisallowedCountries = plan.DisallowedCountries.ValueString()
	paramsSiteUpdate.DaysBeforeDeletingDisabledUsers = plan.DaysBeforeDeletingDisabledUsers.ValueInt64()
	paramsSiteUpdate.DaysToRetainBackups = plan.DaysToRetainBackups.ValueInt64()
	paramsSiteUpdate.MaxPriorPasswords = plan.MaxPriorPasswords.ValueInt64()
	paramsSiteUpdate.PasswordValidityDays = plan.PasswordValidityDays.ValueInt64()
	paramsSiteUpdate.PasswordMinLength = plan.PasswordMinLength.ValueInt64()
	if !plan.PasswordRequireLetter.IsNull() && !plan.PasswordRequireLetter.IsUnknown() {
		paramsSiteUpdate.PasswordRequireLetter = plan.PasswordRequireLetter.ValueBoolPointer()
	}
	if !plan.PasswordRequireMixed.IsNull() && !plan.PasswordRequireMixed.IsUnknown() {
		paramsSiteUpdate.PasswordRequireMixed = plan.PasswordRequireMixed.ValueBoolPointer()
	}
	if !plan.PasswordRequireSpecial.IsNull() && !plan.PasswordRequireSpecial.IsUnknown() {
		paramsSiteUpdate.PasswordRequireSpecial = plan.PasswordRequireSpecial.ValueBoolPointer()
	}
	if !plan.PasswordRequireNumber.IsNull() && !plan.PasswordRequireNumber.IsUnknown() {
		paramsSiteUpdate.PasswordRequireNumber = plan.PasswordRequireNumber.ValueBoolPointer()
	}
	if !plan.PasswordRequireUnbreached.IsNull() && !plan.PasswordRequireUnbreached.IsUnknown() {
		paramsSiteUpdate.PasswordRequireUnbreached = plan.PasswordRequireUnbreached.ValueBoolPointer()
	}
	if !plan.RequireLogoutFromBundlesAndInboxes.IsNull() && !plan.RequireLogoutFromBundlesAndInboxes.IsUnknown() {
		paramsSiteUpdate.RequireLogoutFromBundlesAndInboxes = plan.RequireLogoutFromBundlesAndInboxes.ValueBoolPointer()
	}
	if !plan.DavUserRootEnabled.IsNull() && !plan.DavUserRootEnabled.IsUnknown() {
		paramsSiteUpdate.DavUserRootEnabled = plan.DavUserRootEnabled.ValueBoolPointer()
	}
	if !plan.SftpUserRootEnabled.IsNull() && !plan.SftpUserRootEnabled.IsUnknown() {
		paramsSiteUpdate.SftpUserRootEnabled = plan.SftpUserRootEnabled.ValueBoolPointer()
	}
	if !plan.DisablePasswordReset.IsNull() && !plan.DisablePasswordReset.IsUnknown() {
		paramsSiteUpdate.DisablePasswordReset = plan.DisablePasswordReset.ValueBoolPointer()
	}
	paramsSiteUpdate.BundleNotFoundMessage = plan.BundleNotFoundMessage.ValueString()
	if !plan.BundlePasswordRequired.IsNull() && !plan.BundlePasswordRequired.IsUnknown() {
		paramsSiteUpdate.BundlePasswordRequired = plan.BundlePasswordRequired.ValueBoolPointer()
	}
	if !plan.BundleRequireRegistration.IsNull() && !plan.BundleRequireRegistration.IsUnknown() {
		paramsSiteUpdate.BundleRequireRegistration = plan.BundleRequireRegistration.ValueBoolPointer()
	}
	if !plan.BundleRequireShareRecipient.IsNull() && !plan.BundleRequireShareRecipient.IsUnknown() {
		paramsSiteUpdate.BundleRequireShareRecipient = plan.BundleRequireShareRecipient.ValueBoolPointer()
	}
	paramsSiteUpdate.BundleRegistrationNotifications = plan.BundleRegistrationNotifications.ValueString()
	paramsSiteUpdate.BundleActivityNotifications = plan.BundleActivityNotifications.ValueString()
	paramsSiteUpdate.BundleUploadReceiptNotifications = plan.BundleUploadReceiptNotifications.ValueString()
	if !plan.DocumentEditsInBundleAllowed.IsNull() && !plan.DocumentEditsInBundleAllowed.IsUnknown() {
		paramsSiteUpdate.DocumentEditsInBundleAllowed = plan.DocumentEditsInBundleAllowed.ValueBoolPointer()
	}
	if !plan.PasswordRequirementsApplyToBundles.IsNull() && !plan.PasswordRequirementsApplyToBundles.IsUnknown() {
		paramsSiteUpdate.PasswordRequirementsApplyToBundles = plan.PasswordRequirementsApplyToBundles.ValueBoolPointer()
	}
	if !plan.PreventRootPermissionsForNonSiteAdmins.IsNull() && !plan.PreventRootPermissionsForNonSiteAdmins.IsUnknown() {
		paramsSiteUpdate.PreventRootPermissionsForNonSiteAdmins = plan.PreventRootPermissionsForNonSiteAdmins.ValueBoolPointer()
	}
	if !plan.OptOutGlobal.IsNull() && !plan.OptOutGlobal.IsUnknown() {
		paramsSiteUpdate.OptOutGlobal = plan.OptOutGlobal.ValueBoolPointer()
	}
	if !plan.UseProvidedModifiedAt.IsNull() && !plan.UseProvidedModifiedAt.IsUnknown() {
		paramsSiteUpdate.UseProvidedModifiedAt = plan.UseProvidedModifiedAt.ValueBoolPointer()
	}
	if !plan.CustomNamespace.IsNull() && !plan.CustomNamespace.IsUnknown() {
		paramsSiteUpdate.CustomNamespace = plan.CustomNamespace.ValueBoolPointer()
	}
	if !plan.NonSsoGroupsAllowed.IsNull() && !plan.NonSsoGroupsAllowed.IsUnknown() {
		paramsSiteUpdate.NonSsoGroupsAllowed = plan.NonSsoGroupsAllowed.ValueBoolPointer()
	}
	if !plan.NonSsoUsersAllowed.IsNull() && !plan.NonSsoUsersAllowed.IsUnknown() {
		paramsSiteUpdate.NonSsoUsersAllowed = plan.NonSsoUsersAllowed.ValueBoolPointer()
	}
	if !plan.SharingEnabled.IsNull() && !plan.SharingEnabled.IsUnknown() {
		paramsSiteUpdate.SharingEnabled = plan.SharingEnabled.ValueBoolPointer()
	}
	if !plan.SnapshotSharingEnabled.IsNull() && !plan.SnapshotSharingEnabled.IsUnknown() {
		paramsSiteUpdate.SnapshotSharingEnabled = plan.SnapshotSharingEnabled.ValueBoolPointer()
	}
	if !plan.UserRequestsEnabled.IsNull() && !plan.UserRequestsEnabled.IsUnknown() {
		paramsSiteUpdate.UserRequestsEnabled = plan.UserRequestsEnabled.ValueBoolPointer()
	}
	if !plan.UserRequestsNotifyAdmins.IsNull() && !plan.UserRequestsNotifyAdmins.IsUnknown() {
		paramsSiteUpdate.UserRequestsNotifyAdmins = plan.UserRequestsNotifyAdmins.ValueBoolPointer()
	}
	if !plan.DavEnabled.IsNull() && !plan.DavEnabled.IsUnknown() {
		paramsSiteUpdate.DavEnabled = plan.DavEnabled.ValueBoolPointer()
	}
	if !plan.FtpEnabled.IsNull() && !plan.FtpEnabled.IsUnknown() {
		paramsSiteUpdate.FtpEnabled = plan.FtpEnabled.ValueBoolPointer()
	}
	if !plan.SftpEnabled.IsNull() && !plan.SftpEnabled.IsUnknown() {
		paramsSiteUpdate.SftpEnabled = plan.SftpEnabled.ValueBoolPointer()
	}
	if !plan.UsersCanCreateApiKeys.IsNull() && !plan.UsersCanCreateApiKeys.IsUnknown() {
		paramsSiteUpdate.UsersCanCreateApiKeys = plan.UsersCanCreateApiKeys.ValueBoolPointer()
	}
	if !plan.UsersCanCreateSshKeys.IsNull() && !plan.UsersCanCreateSshKeys.IsUnknown() {
		paramsSiteUpdate.UsersCanCreateSshKeys = plan.UsersCanCreateSshKeys.ValueBoolPointer()
	}
	if !plan.ShowUserNotificationsLogInLink.IsNull() && !plan.ShowUserNotificationsLogInLink.IsUnknown() {
		paramsSiteUpdate.ShowUserNotificationsLogInLink = plan.ShowUserNotificationsLogInLink.ValueBoolPointer()
	}
	paramsSiteUpdate.SftpHostKeyType = plan.SftpHostKeyType.ValueString()
	paramsSiteUpdate.ActiveSftpHostKeyId = plan.ActiveSftpHostKeyId.ValueInt64()
	if !plan.ProtocolAccessGroupsOnly.IsNull() && !plan.ProtocolAccessGroupsOnly.IsUnknown() {
		paramsSiteUpdate.ProtocolAccessGroupsOnly = plan.ProtocolAccessGroupsOnly.ValueBoolPointer()
	}
	if !plan.RevokeBundleAccessOnDisableOrDelete.IsNull() && !plan.RevokeBundleAccessOnDisableOrDelete.IsUnknown() {
		paramsSiteUpdate.RevokeBundleAccessOnDisableOrDelete = plan.RevokeBundleAccessOnDisableOrDelete.ValueBoolPointer()
	}
	updateBundleWatermarkValue, diags := lib.DynamicToStringMap(ctx, path.Root("bundle_watermark_value"), plan.BundleWatermarkValue)
	resp.Diagnostics.Append(diags...)
	paramsSiteUpdate.BundleWatermarkValue = updateBundleWatermarkValue
	if !plan.GroupAdminsCanSetUserPassword.IsNull() && !plan.GroupAdminsCanSetUserPassword.IsUnknown() {
		paramsSiteUpdate.GroupAdminsCanSetUserPassword = plan.GroupAdminsCanSetUserPassword.ValueBoolPointer()
	}
	if !plan.BundleRecipientBlacklistFreeEmailDomains.IsNull() && !plan.BundleRecipientBlacklistFreeEmailDomains.IsUnknown() {
		paramsSiteUpdate.BundleRecipientBlacklistFreeEmailDomains = plan.BundleRecipientBlacklistFreeEmailDomains.ValueBoolPointer()
	}
	if !plan.BundleRecipientBlacklistDomains.IsNull() && !plan.BundleRecipientBlacklistDomains.IsUnknown() {
		diags = plan.BundleRecipientBlacklistDomains.ElementsAs(ctx, &paramsSiteUpdate.BundleRecipientBlacklistDomains, false)
		resp.Diagnostics.Append(diags...)
	}
	if !plan.AdminsBypassLockedSubfolders.IsNull() && !plan.AdminsBypassLockedSubfolders.IsUnknown() {
		paramsSiteUpdate.AdminsBypassLockedSubfolders = plan.AdminsBypassLockedSubfolders.ValueBoolPointer()
	}
	if !plan.Allowed2faMethodSms.IsNull() && !plan.Allowed2faMethodSms.IsUnknown() {
		paramsSiteUpdate.Allowed2faMethodSms = plan.Allowed2faMethodSms.ValueBoolPointer()
	}
	if !plan.Allowed2faMethodTotp.IsNull() && !plan.Allowed2faMethodTotp.IsUnknown() {
		paramsSiteUpdate.Allowed2faMethodTotp = plan.Allowed2faMethodTotp.ValueBoolPointer()
	}
	if !plan.Allowed2faMethodWebauthn.IsNull() && !plan.Allowed2faMethodWebauthn.IsUnknown() {
		paramsSiteUpdate.Allowed2faMethodWebauthn = plan.Allowed2faMethodWebauthn.ValueBoolPointer()
	}
	if !plan.Allowed2faMethodYubi.IsNull() && !plan.Allowed2faMethodYubi.IsUnknown() {
		paramsSiteUpdate.Allowed2faMethodYubi = plan.Allowed2faMethodYubi.ValueBoolPointer()
	}
	if !plan.Allowed2faMethodEmail.IsNull() && !plan.Allowed2faMethodEmail.IsUnknown() {
		paramsSiteUpdate.Allowed2faMethodEmail = plan.Allowed2faMethodEmail.ValueBoolPointer()
	}
	if !plan.Allowed2faMethodStatic.IsNull() && !plan.Allowed2faMethodStatic.IsUnknown() {
		paramsSiteUpdate.Allowed2faMethodStatic = plan.Allowed2faMethodStatic.ValueBoolPointer()
	}
	if !plan.Allowed2faMethodBypassForFtpSftpDav.IsNull() && !plan.Allowed2faMethodBypassForFtpSftpDav.IsUnknown() {
		paramsSiteUpdate.Allowed2faMethodBypassForFtpSftpDav = plan.Allowed2faMethodBypassForFtpSftpDav.ValueBoolPointer()
	}
	if !plan.Require2fa.IsNull() && !plan.Require2fa.IsUnknown() {
		paramsSiteUpdate.Require2fa = plan.Require2fa.ValueBoolPointer()
	}
	paramsSiteUpdate.Require2faUserType = plan.Require2faUserType.ValueString()
	paramsSiteUpdate.Color2Top = plan.Color2Top.ValueString()
	paramsSiteUpdate.Color2Left = plan.Color2Left.ValueString()
	paramsSiteUpdate.Color2Link = plan.Color2Link.ValueString()
	paramsSiteUpdate.Color2Text = plan.Color2Text.ValueString()
	paramsSiteUpdate.Color2TopText = plan.Color2TopText.ValueString()
	paramsSiteUpdate.SiteHeader = plan.SiteHeader.ValueString()
	paramsSiteUpdate.SiteFooter = plan.SiteFooter.ValueString()
	paramsSiteUpdate.LoginHelpText = plan.LoginHelpText.ValueString()
	if !plan.UseDedicatedIpsForSmtp.IsNull() && !plan.UseDedicatedIpsForSmtp.IsUnknown() {
		paramsSiteUpdate.UseDedicatedIpsForSmtp = plan.UseDedicatedIpsForSmtp.ValueBoolPointer()
	}
	paramsSiteUpdate.SmtpAddress = plan.SmtpAddress.ValueString()
	paramsSiteUpdate.SmtpAuthentication = plan.SmtpAuthentication.ValueString()
	paramsSiteUpdate.SmtpFrom = plan.SmtpFrom.ValueString()
	paramsSiteUpdate.SmtpUsername = plan.SmtpUsername.ValueString()
	paramsSiteUpdate.SmtpPort = plan.SmtpPort.ValueInt64()
	if !plan.LdapEnabled.IsNull() && !plan.LdapEnabled.IsUnknown() {
		paramsSiteUpdate.LdapEnabled = plan.LdapEnabled.ValueBoolPointer()
	}
	paramsSiteUpdate.LdapType = plan.LdapType.ValueString()
	paramsSiteUpdate.LdapHost = plan.LdapHost.ValueString()
	paramsSiteUpdate.LdapHost2 = plan.LdapHost2.ValueString()
	paramsSiteUpdate.LdapHost3 = plan.LdapHost3.ValueString()
	paramsSiteUpdate.LdapPort = plan.LdapPort.ValueInt64()
	if !plan.LdapSecure.IsNull() && !plan.LdapSecure.IsUnknown() {
		paramsSiteUpdate.LdapSecure = plan.LdapSecure.ValueBoolPointer()
	}
	paramsSiteUpdate.LdapUsername = plan.LdapUsername.ValueString()
	paramsSiteUpdate.LdapUsernameField = plan.LdapUsernameField.ValueString()
	paramsSiteUpdate.LdapDomain = plan.LdapDomain.ValueString()
	paramsSiteUpdate.LdapUserAction = plan.LdapUserAction.ValueString()
	paramsSiteUpdate.LdapGroupAction = plan.LdapGroupAction.ValueString()
	paramsSiteUpdate.LdapUserIncludeGroups = plan.LdapUserIncludeGroups.ValueString()
	paramsSiteUpdate.LdapGroupExclusion = plan.LdapGroupExclusion.ValueString()
	paramsSiteUpdate.LdapGroupInclusion = plan.LdapGroupInclusion.ValueString()
	paramsSiteUpdate.LdapBaseDn = plan.LdapBaseDn.ValueString()
	if !plan.UploadsViaEmailAuthentication.IsNull() && !plan.UploadsViaEmailAuthentication.IsUnknown() {
		paramsSiteUpdate.UploadsViaEmailAuthentication = plan.UploadsViaEmailAuthentication.ValueBoolPointer()
	}
	paramsSiteUpdate.SessionExpiryMinutes = plan.SessionExpiryMinutes.ValueInt64()

	if resp.Diagnostics.HasError() {
		return
	}

	site, err := r.client.Update(paramsSiteUpdate, files_sdk.WithContext(ctx))
	if err != nil {
		resp.Diagnostics.AddError(
			"Error Updating Files Site",
			"Could not update site, unexpected error: "+err.Error(),
		)
		return
	}

	diags = r.populateResourceModel(ctx, site, &plan)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	diags = resp.State.Set(ctx, plan)
	resp.Diagnostics.Append(diags...)
}

func (r *siteResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	resp.Diagnostics.AddError(
		"Resource Delete Not Implemented",
		"This resource does not support deletion.",
	)
}

func (r *siteResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	resp.Diagnostics.Append(resp.State.SetAttribute(ctx, path.Root("id"), 0)...)
}

func (r *siteResource) populateResourceModel(ctx context.Context, site files_sdk.Site, state *siteResourceModel) (diags diag.Diagnostics) {
	var propDiags diag.Diagnostics

	state.Id = types.Int64Value(site.Id)
	state.Name = types.StringValue(site.Name)
	state.AdditionalTextFileTypes, propDiags = types.ListValueFrom(ctx, types.StringType, site.AdditionalTextFileTypes)
	diags.Append(propDiags...)
	state.Allowed2faMethodSms = types.BoolPointerValue(site.Allowed2faMethodSms)
	state.Allowed2faMethodTotp = types.BoolPointerValue(site.Allowed2faMethodTotp)
	state.Allowed2faMethodWebauthn = types.BoolPointerValue(site.Allowed2faMethodWebauthn)
	state.Allowed2faMethodYubi = types.BoolPointerValue(site.Allowed2faMethodYubi)
	state.Allowed2faMethodEmail = types.BoolPointerValue(site.Allowed2faMethodEmail)
	state.Allowed2faMethodStatic = types.BoolPointerValue(site.Allowed2faMethodStatic)
	state.Allowed2faMethodBypassForFtpSftpDav = types.BoolPointerValue(site.Allowed2faMethodBypassForFtpSftpDav)
	state.AdminUserId = types.Int64Value(site.AdminUserId)
	state.AdminsBypassLockedSubfolders = types.BoolPointerValue(site.AdminsBypassLockedSubfolders)
	state.AllowBundleNames = types.BoolPointerValue(site.AllowBundleNames)
	state.AllowedCountries = types.StringValue(site.AllowedCountries)
	state.AllowedIps = types.StringValue(site.AllowedIps)
	state.AlwaysMkdirParents = types.BoolPointerValue(site.AlwaysMkdirParents)
	state.AskAboutOverwrites = types.BoolPointerValue(site.AskAboutOverwrites)
	state.BundleActivityNotifications = types.StringValue(site.BundleActivityNotifications)
	state.BundleExpiration = types.Int64Value(site.BundleExpiration)
	state.BundleNotFoundMessage = types.StringValue(site.BundleNotFoundMessage)
	state.BundlePasswordRequired = types.BoolPointerValue(site.BundlePasswordRequired)
	state.BundleRecipientBlacklistDomains, propDiags = types.ListValueFrom(ctx, types.StringType, site.BundleRecipientBlacklistDomains)
	diags.Append(propDiags...)
	state.BundleRecipientBlacklistFreeEmailDomains = types.BoolPointerValue(site.BundleRecipientBlacklistFreeEmailDomains)
	state.BundleRegistrationNotifications = types.StringValue(site.BundleRegistrationNotifications)
	state.BundleRequireRegistration = types.BoolPointerValue(site.BundleRequireRegistration)
	state.BundleRequireShareRecipient = types.BoolPointerValue(site.BundleRequireShareRecipient)
	state.BundleRequireNote = types.BoolPointerValue(site.BundleRequireNote)
	state.BundleSendSharedReceipts = types.BoolPointerValue(site.BundleSendSharedReceipts)
	state.BundleUploadReceiptNotifications = types.StringValue(site.BundleUploadReceiptNotifications)
	respBundleWatermarkAttachment, err := json.Marshal(site.BundleWatermarkAttachment)
	if err != nil {
		diags.AddError(
			"Error Creating Files Site",
			"Could not marshal bundle_watermark_attachment to JSON: "+err.Error(),
		)
	}
	state.BundleWatermarkAttachment = types.StringValue(string(respBundleWatermarkAttachment))
	state.BundleWatermarkValue, propDiags = lib.ToDynamic(ctx, path.Root("bundle_watermark_value"), site.BundleWatermarkValue, state.BundleWatermarkValue.UnderlyingValue())
	diags.Append(propDiags...)
	state.CalculateFileChecksumsCrc32 = types.BoolPointerValue(site.CalculateFileChecksumsCrc32)
	state.CalculateFileChecksumsMd5 = types.BoolPointerValue(site.CalculateFileChecksumsMd5)
	state.CalculateFileChecksumsSha1 = types.BoolPointerValue(site.CalculateFileChecksumsSha1)
	state.CalculateFileChecksumsSha256 = types.BoolPointerValue(site.CalculateFileChecksumsSha256)
	state.UploadsViaEmailAuthentication = types.BoolPointerValue(site.UploadsViaEmailAuthentication)
	state.Color2Left = types.StringValue(site.Color2Left)
	state.Color2Link = types.StringValue(site.Color2Link)
	state.Color2Text = types.StringValue(site.Color2Text)
	state.Color2Top = types.StringValue(site.Color2Top)
	state.Color2TopText = types.StringValue(site.Color2TopText)
	state.ContactName = types.StringValue(site.ContactName)
	if err := lib.TimeToStringType(ctx, path.Root("created_at"), site.CreatedAt, &state.CreatedAt); err != nil {
		diags.AddError(
			"Error Creating Files Site",
			"Could not convert state created_at to string: "+err.Error(),
		)
	}
	state.Currency = types.StringValue(site.Currency)
	state.CustomNamespace = types.BoolPointerValue(site.CustomNamespace)
	state.DavEnabled = types.BoolPointerValue(site.DavEnabled)
	state.DavUserRootEnabled = types.BoolPointerValue(site.DavUserRootEnabled)
	state.DaysBeforeDeletingDisabledUsers = types.Int64Value(site.DaysBeforeDeletingDisabledUsers)
	state.DaysToRetainBackups = types.Int64Value(site.DaysToRetainBackups)
	state.DocumentEditsInBundleAllowed = types.BoolPointerValue(site.DocumentEditsInBundleAllowed)
	state.DefaultTimeZone = types.StringValue(site.DefaultTimeZone)
	state.DesktopApp = types.BoolPointerValue(site.DesktopApp)
	state.DesktopAppSessionIpPinning = types.BoolPointerValue(site.DesktopAppSessionIpPinning)
	state.DesktopAppSessionLifetime = types.Int64Value(site.DesktopAppSessionLifetime)
	state.LegacyChecksumsMode = types.BoolPointerValue(site.LegacyChecksumsMode)
	state.MobileApp = types.BoolPointerValue(site.MobileApp)
	state.MobileAppSessionIpPinning = types.BoolPointerValue(site.MobileAppSessionIpPinning)
	state.MobileAppSessionLifetime = types.Int64Value(site.MobileAppSessionLifetime)
	state.DisallowedCountries = types.StringValue(site.DisallowedCountries)
	state.DisableFilesCertificateGeneration = types.BoolPointerValue(site.DisableFilesCertificateGeneration)
	state.DisableNotifications = types.BoolPointerValue(site.DisableNotifications)
	state.DisablePasswordReset = types.BoolPointerValue(site.DisablePasswordReset)
	state.Domain = types.StringValue(site.Domain)
	state.DomainHstsHeader = types.BoolPointerValue(site.DomainHstsHeader)
	state.DomainLetsencryptChain = types.StringValue(site.DomainLetsencryptChain)
	state.Email = types.StringValue(site.Email)
	state.FtpEnabled = types.BoolPointerValue(site.FtpEnabled)
	state.ReplyToEmail = types.StringValue(site.ReplyToEmail)
	state.NonSsoGroupsAllowed = types.BoolPointerValue(site.NonSsoGroupsAllowed)
	state.NonSsoUsersAllowed = types.BoolPointerValue(site.NonSsoUsersAllowed)
	state.FolderPermissionsGroupsOnly = types.BoolPointerValue(site.FolderPermissionsGroupsOnly)
	state.Hipaa = types.BoolPointerValue(site.Hipaa)
	respIcon128, err := json.Marshal(site.Icon128)
	if err != nil {
		diags.AddError(
			"Error Creating Files Site",
			"Could not marshal icon128 to JSON: "+err.Error(),
		)
	}
	state.Icon128 = types.StringValue(string(respIcon128))
	respIcon16, err := json.Marshal(site.Icon16)
	if err != nil {
		diags.AddError(
			"Error Creating Files Site",
			"Could not marshal icon16 to JSON: "+err.Error(),
		)
	}
	state.Icon16 = types.StringValue(string(respIcon16))
	respIcon32, err := json.Marshal(site.Icon32)
	if err != nil {
		diags.AddError(
			"Error Creating Files Site",
			"Could not marshal icon32 to JSON: "+err.Error(),
		)
	}
	state.Icon32 = types.StringValue(string(respIcon32))
	respIcon48, err := json.Marshal(site.Icon48)
	if err != nil {
		diags.AddError(
			"Error Creating Files Site",
			"Could not marshal icon48 to JSON: "+err.Error(),
		)
	}
	state.Icon48 = types.StringValue(string(respIcon48))
	if err := lib.TimeToStringType(ctx, path.Root("immutable_files_set_at"), site.ImmutableFilesSetAt, &state.ImmutableFilesSetAt); err != nil {
		diags.AddError(
			"Error Creating Files Site",
			"Could not convert state immutable_files_set_at to string: "+err.Error(),
		)
	}
	state.IncludePasswordInWelcomeEmail = types.BoolPointerValue(site.IncludePasswordInWelcomeEmail)
	state.Language = types.StringValue(site.Language)
	state.LdapBaseDn = types.StringValue(site.LdapBaseDn)
	state.LdapDomain = types.StringValue(site.LdapDomain)
	state.LdapEnabled = types.BoolPointerValue(site.LdapEnabled)
	state.LdapGroupAction = types.StringValue(site.LdapGroupAction)
	state.LdapGroupExclusion = types.StringValue(site.LdapGroupExclusion)
	state.LdapGroupInclusion = types.StringValue(site.LdapGroupInclusion)
	state.LdapHost = types.StringValue(site.LdapHost)
	state.LdapHost2 = types.StringValue(site.LdapHost2)
	state.LdapHost3 = types.StringValue(site.LdapHost3)
	state.LdapPort = types.Int64Value(site.LdapPort)
	state.LdapSecure = types.BoolPointerValue(site.LdapSecure)
	state.LdapType = types.StringValue(site.LdapType)
	state.LdapUserAction = types.StringValue(site.LdapUserAction)
	state.LdapUserIncludeGroups = types.StringValue(site.LdapUserIncludeGroups)
	state.LdapUsername = types.StringValue(site.LdapUsername)
	state.LdapUsernameField = types.StringValue(site.LdapUsernameField)
	state.LoginHelpText = types.StringValue(site.LoginHelpText)
	respLogo, err := json.Marshal(site.Logo)
	if err != nil {
		diags.AddError(
			"Error Creating Files Site",
			"Could not marshal logo to JSON: "+err.Error(),
		)
	}
	state.Logo = types.StringValue(string(respLogo))
	respLoginPageBackgroundImage, err := json.Marshal(site.LoginPageBackgroundImage)
	if err != nil {
		diags.AddError(
			"Error Creating Files Site",
			"Could not marshal login_page_background_image to JSON: "+err.Error(),
		)
	}
	state.LoginPageBackgroundImage = types.StringValue(string(respLoginPageBackgroundImage))
	state.MaxPriorPasswords = types.Int64Value(site.MaxPriorPasswords)
	state.MotdText = types.StringValue(site.MotdText)
	state.MotdUseForFtp = types.BoolPointerValue(site.MotdUseForFtp)
	state.MotdUseForSftp = types.BoolPointerValue(site.MotdUseForSftp)
	state.NextBillingAmount = types.StringValue(site.NextBillingAmount)
	state.NextBillingDate = types.StringValue(site.NextBillingDate)
	state.OfficeIntegrationAvailable = types.BoolPointerValue(site.OfficeIntegrationAvailable)
	state.OfficeIntegrationType = types.StringValue(site.OfficeIntegrationType)
	state.OncehubLink = types.StringValue(site.OncehubLink)
	state.OptOutGlobal = types.BoolPointerValue(site.OptOutGlobal)
	state.Overdue = types.BoolPointerValue(site.Overdue)
	state.PasswordMinLength = types.Int64Value(site.PasswordMinLength)
	state.PasswordRequireLetter = types.BoolPointerValue(site.PasswordRequireLetter)
	state.PasswordRequireMixed = types.BoolPointerValue(site.PasswordRequireMixed)
	state.PasswordRequireNumber = types.BoolPointerValue(site.PasswordRequireNumber)
	state.PasswordRequireSpecial = types.BoolPointerValue(site.PasswordRequireSpecial)
	state.PasswordRequireUnbreached = types.BoolPointerValue(site.PasswordRequireUnbreached)
	state.PasswordRequirementsApplyToBundles = types.BoolPointerValue(site.PasswordRequirementsApplyToBundles)
	state.PasswordValidityDays = types.Int64Value(site.PasswordValidityDays)
	state.Phone = types.StringValue(site.Phone)
	state.PinAllRemoteServersToSiteRegion = types.BoolPointerValue(site.PinAllRemoteServersToSiteRegion)
	state.PreventRootPermissionsForNonSiteAdmins = types.BoolPointerValue(site.PreventRootPermissionsForNonSiteAdmins)
	state.ProtocolAccessGroupsOnly = types.BoolPointerValue(site.ProtocolAccessGroupsOnly)
	state.Require2fa = types.BoolPointerValue(site.Require2fa)
	if err := lib.TimeToStringType(ctx, path.Root("require_2fa_stop_time"), site.Require2faStopTime, &state.Require2faStopTime); err != nil {
		diags.AddError(
			"Error Creating Files Site",
			"Could not convert state require_2fa_stop_time to string: "+err.Error(),
		)
	}
	state.RevokeBundleAccessOnDisableOrDelete = types.BoolPointerValue(site.RevokeBundleAccessOnDisableOrDelete)
	state.Require2faUserType = types.StringValue(site.Require2faUserType)
	state.RequireLogoutFromBundlesAndInboxes = types.BoolPointerValue(site.RequireLogoutFromBundlesAndInboxes)
	respSession, err := json.Marshal(site.Session)
	if err != nil {
		diags.AddError(
			"Error Creating Files Site",
			"Could not marshal session to JSON: "+err.Error(),
		)
	}
	state.Session = types.StringValue(string(respSession))
	state.SftpEnabled = types.BoolPointerValue(site.SftpEnabled)
	state.SftpHostKeyType = types.StringValue(site.SftpHostKeyType)
	state.ActiveSftpHostKeyId = types.Int64Value(site.ActiveSftpHostKeyId)
	state.SftpInsecureCiphers = types.BoolPointerValue(site.SftpInsecureCiphers)
	state.SftpInsecureDiffieHellman = types.BoolPointerValue(site.SftpInsecureDiffieHellman)
	state.SftpUserRootEnabled = types.BoolPointerValue(site.SftpUserRootEnabled)
	state.SharingEnabled = types.BoolPointerValue(site.SharingEnabled)
	state.ShowUserNotificationsLogInLink = types.BoolPointerValue(site.ShowUserNotificationsLogInLink)
	state.ShowRequestAccessLink = types.BoolPointerValue(site.ShowRequestAccessLink)
	state.SiteFooter = types.StringValue(site.SiteFooter)
	state.SiteHeader = types.StringValue(site.SiteHeader)
	state.SmtpAddress = types.StringValue(site.SmtpAddress)
	state.SmtpAuthentication = types.StringValue(site.SmtpAuthentication)
	state.SmtpFrom = types.StringValue(site.SmtpFrom)
	state.SmtpPort = types.Int64Value(site.SmtpPort)
	state.SmtpUsername = types.StringValue(site.SmtpUsername)
	state.SessionExpiry = types.StringValue(site.SessionExpiry)
	state.SessionExpiryMinutes = types.Int64Value(site.SessionExpiryMinutes)
	state.SnapshotSharingEnabled = types.BoolPointerValue(site.SnapshotSharingEnabled)
	state.SslRequired = types.BoolPointerValue(site.SslRequired)
	state.Subdomain = types.StringValue(site.Subdomain)
	if err := lib.TimeToStringType(ctx, path.Root("switch_to_plan_date"), site.SwitchToPlanDate, &state.SwitchToPlanDate); err != nil {
		diags.AddError(
			"Error Creating Files Site",
			"Could not convert state switch_to_plan_date to string: "+err.Error(),
		)
	}
	state.TlsDisabled = types.BoolPointerValue(site.TlsDisabled)
	state.TrialDaysLeft = types.Int64Value(site.TrialDaysLeft)
	if err := lib.TimeToStringType(ctx, path.Root("trial_until"), site.TrialUntil, &state.TrialUntil); err != nil {
		diags.AddError(
			"Error Creating Files Site",
			"Could not convert state trial_until to string: "+err.Error(),
		)
	}
	state.UseDedicatedIpsForSmtp = types.BoolPointerValue(site.UseDedicatedIpsForSmtp)
	state.UseProvidedModifiedAt = types.BoolPointerValue(site.UseProvidedModifiedAt)
	respUser, err := json.Marshal(site.User)
	if err != nil {
		diags.AddError(
			"Error Creating Files Site",
			"Could not marshal user to JSON: "+err.Error(),
		)
	}
	state.User = types.StringValue(string(respUser))
	state.UserLockout = types.BoolPointerValue(site.UserLockout)
	state.UserLockoutLockPeriod = types.Int64Value(site.UserLockoutLockPeriod)
	state.UserLockoutTries = types.Int64Value(site.UserLockoutTries)
	state.UserLockoutWithin = types.Int64Value(site.UserLockoutWithin)
	state.UserRequestsEnabled = types.BoolPointerValue(site.UserRequestsEnabled)
	state.UserRequestsNotifyAdmins = types.BoolPointerValue(site.UserRequestsNotifyAdmins)
	state.UsersCanCreateApiKeys = types.BoolPointerValue(site.UsersCanCreateApiKeys)
	state.UsersCanCreateSshKeys = types.BoolPointerValue(site.UsersCanCreateSshKeys)
	state.WelcomeCustomText = types.StringValue(site.WelcomeCustomText)
	state.WelcomeEmailCc = types.StringValue(site.WelcomeEmailCc)
	state.WelcomeEmailSubject = types.StringValue(site.WelcomeEmailSubject)
	state.WelcomeEmailEnabled = types.BoolPointerValue(site.WelcomeEmailEnabled)
	state.WelcomeScreen = types.StringValue(site.WelcomeScreen)
	state.WindowsModeFtp = types.BoolPointerValue(site.WindowsModeFtp)
	state.GroupAdminsCanSetUserPassword = types.BoolPointerValue(site.GroupAdminsCanSetUserPassword)

	return
}
