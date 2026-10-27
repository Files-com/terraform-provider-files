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
	ChildSiteIds        types.List    `tfsdk:"child_site_ids"`
	DefaultPolicy       types.Bool    `tfsdk:"default_policy"`
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
		Description: "A Child Site Management Policy is a centralized policy defined by a parent site to enforce consistent configurations across child sites. These policies allow parent sites to maintain control over specific aspects of their child sites' functionality and appearance.\n\nNon-default policies apply only to the child sites listed in `child_site_ids`, and each child site can be explicitly assigned to only one policy.\n\nOne policy can be designated as the default policy. It applies to every child site not explicitly assigned to another policy or listed in its `skip_child_site_ids`, including newly created child sites. Only a default policy can exclude child sites. The `value` field contains the policy configuration data, with the format varying based on the policy type. When a policy is active, its managed configurations are automatically enforced on applicable child sites, and attribute modifications are not permitted.\n\nOn a Child Site, listing policies returns only the policy covering that site. Each key in `value` identifies a managed setting, but its value is always `null`. Read non-secret effective setting values from the Child Site's own Site settings response.",
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
				Description: "Policy configuration data. On a Child Site, contains only managed setting names with null values. Settings policies accept site settings plus an optional `folder_behaviors` array for parent-managed root behaviors on child sites. For more information, refer to the Value Hash section of the developer documentation.",
				Computed:    true,
				Optional:    true,
				Sensitive:   true,
				Validators: []validator.Dynamic{
					lib.DeprecatedJSONEncoding("files_child_site_management_policy.value", "value = {\n  color2_left      = \"#000000\"\n  folder_behaviors = [\n    {\n      file_expiration = {\n        value = {\n          days_to_retain = 30\n        }\n      }\n    }\n  ]\n}", "March 1, 2027", map[string]string{"domain_hsts_header": "boolean", "allow_bundle_names": "boolean", "bundle_expiration": "integer", "welcome_email_enabled": "boolean", "ask_about_overwrites": "boolean", "always_mkdir_parents": "boolean", "windows_mode_ftp": "boolean", "desktop_app": "boolean", "desktop_app_session_ip_pinning": "boolean", "desktop_app_session_lifetime": "integer", "mobile_app": "boolean", "mobile_app_session_ip_pinning": "boolean", "mobile_app_session_lifetime": "integer", "folder_permissions_groups_only": "boolean", "office_integration_available": "boolean", "pin_all_remote_servers_to_site_region": "boolean", "motd_use_for_ftp": "boolean", "motd_use_for_sftp": "boolean", "left_navigation_visibility": "object", "left_navigation_visibility.sharing": "object", "left_navigation_visibility.sharing.admin": "boolean", "left_navigation_visibility.sharing.all": "boolean", "left_navigation_visibility.automation_flow": "object", "left_navigation_visibility.automation_flow.admin": "boolean", "left_navigation_visibility.automation_flow.all": "boolean", "left_navigation_visibility.integrations": "object", "left_navigation_visibility.integrations.admin": "boolean", "left_navigation_visibility.integrations.all": "boolean", "left_navigation_visibility.clients_protocols": "object", "left_navigation_visibility.clients_protocols.admin": "boolean", "left_navigation_visibility.clients_protocols.all": "boolean", "disable_all_ai_features": "boolean", "ai_feature_availability": "object", "ai_feature_availability.in_app_ai_assistant": "object", "ai_feature_availability.in_app_ai_assistant.site_admins": "boolean", "ai_feature_availability.in_app_ai_assistant.workspace_admins": "boolean", "ai_feature_availability.in_app_ai_assistant.folder_admins": "boolean", "ai_feature_availability.in_app_ai_assistant.all_users": "boolean", "ai_feature_availability.in_app_ai_assistant.selected_group_members": "boolean", "ai_feature_availability.in_app_ai_assistant.group_ids": "list", "ai_feature_availability.in_app_ai_assistant.group_ids[]": "integer", "mcp_dcr_enabled": "boolean", "additional_text_file_types": "list", "bundle_require_note": "boolean", "bundle_send_shared_receipts": "boolean", "bundles_default_owned_by_primary_group": "boolean", "calculate_file_checksums_crc32": "boolean", "calculate_file_checksums_md5": "boolean", "calculate_file_checksums_sha1": "boolean", "calculate_file_checksums_sha256": "boolean", "migrate_remote_server_sync_to_sync": "boolean", "as2_message_retention_days": "integer", "file_activity_events_enabled": "boolean", "session_expiry_minutes": "integer", "ssl_required": "boolean", "sftp_insecure_ciphers": "boolean", "sftp_insecure_diffie_hellman": "boolean", "fedramp": "boolean", "user_lockout": "boolean", "user_lockout_tries": "integer", "user_lockout_within": "integer", "user_lockout_lock_period": "integer", "include_password_in_welcome_email": "boolean", "allow_user_level_2fa_override": "boolean", "allow_user_level_allowed_ip_override": "boolean", "allow_user_level_ssl_override": "boolean", "enforce_parent_site_user_ip_restrictions": "boolean", "days_to_retain_backups": "integer", "max_prior_passwords": "integer", "password_validity_days": "integer", "password_min_length": "integer", "password_require_letter": "boolean", "password_require_mixed": "boolean", "password_require_special": "boolean", "password_require_number": "boolean", "password_require_unbreached": "boolean", "require_logout_from_bundles_and_inboxes": "boolean", "dav_user_root_enabled": "boolean", "sftp_user_root_enabled": "boolean", "disable_password_reset": "boolean", "immutable_files": "boolean", "bundle_password_required": "boolean", "bundle_require_registration": "boolean", "bundle_require_share_recipient": "boolean", "bundle_send_one_time_password_to_recipient_at_registration": "boolean", "document_edits_in_bundle_allowed": "boolean", "password_requirements_apply_to_bundles": "boolean", "prevent_root_permissions_for_non_site_admins": "boolean", "restrict_root_folder_behaviors_to_site_admins": "boolean", "root_folder_behaviors_apply_to_workspaces": "boolean", "opt_out_global": "boolean", "use_provided_modified_at": "boolean", "non_sso_groups_allowed": "boolean", "non_sso_users_allowed": "boolean", "sharing_enabled": "boolean", "snapshot_sharing_enabled": "boolean", "user_requests_enabled": "boolean", "user_requests_notify_admins": "boolean", "dav_enabled": "boolean", "files_com_remote_server_enabled": "boolean", "ftp_enabled": "boolean", "s3_compatible_endpoint_enabled": "boolean", "sftp_enabled": "boolean", "sftp_finalize_partial_uploads": "boolean", "users_can_create_api_keys": "boolean", "users_can_create_ssh_keys": "boolean", "show_user_notifications_log_in_link": "boolean", "active_sftp_host_key_id": "integer", "active_sftp_host_key_ids": "list", "active_sftp_host_key_ids[]": "integer", "protocol_access_groups_only": "boolean", "revoke_bundle_access_on_disable_or_delete": "boolean", "bundle_watermark_value": "object", "bundle_watermark_value.max_height_or_width": "integer", "bundle_watermark_value.transparency": "integer", "group_admins_can_add_users": "boolean", "group_admins_can_manage_group_memberships": "boolean", "group_admins_can_delete_users": "boolean", "group_admins_can_enable_disable_users": "boolean", "group_admins_can_modify_users": "boolean", "group_admins_can_bypass_user_lifecycle_rules": "boolean", "group_admins_can_reset_passwords": "boolean", "group_admins_can_set_user_password": "boolean", "child_site_folder_admins_can_see_parent_site_directory": "boolean", "workspace_folder_admins_can_see_default_workspace_directory": "boolean", "bundle_recipient_blacklist_free_email_domains": "boolean", "bundle_recipient_blacklist_domains": "list", "admins_bypass_locked_subfolders": "boolean", "allowed_2fa_method_sms": "boolean", "allowed_2fa_method_totp": "boolean", "allowed_2fa_method_webauthn": "boolean", "allowed_2fa_method_yubi": "boolean", "allowed_2fa_method_email": "boolean", "allowed_2fa_method_static": "boolean", "allowed_2fa_method_bypass_for_ftp_sftp_dav": "boolean", "require_2fa": "boolean", "require_2fa_exempt_all_sso_users": "boolean", "use_dedicated_ips_for_smtp": "boolean", "smtp_port": "integer", "ldap_enabled": "boolean", "ldap_port": "integer", "ldap_secure": "boolean", "uploads_via_email_authentication": "boolean", "folder_behaviors": "list", "folder_behaviors[]": "object", "folder_behaviors[].malware_scanning": "object", "folder_behaviors[].malware_scanning.value": "object", "folder_behaviors[behavior=malware_scanning]": "object", "folder_behaviors[behavior=malware_scanning].value": "object", "folder_behaviors[].file_expiration": "object", "folder_behaviors[].file_expiration.value": "object", "folder_behaviors[].file_expiration.value.days_to_retain": "integer", "folder_behaviors[].file_expiration.value.delete_empty_folders": "boolean", "folder_behaviors[behavior=file_expiration]": "object", "folder_behaviors[behavior=file_expiration].value": "object", "folder_behaviors[behavior=file_expiration].value.days_to_retain": "integer", "folder_behaviors[behavior=file_expiration].value.delete_empty_folders": "boolean", "folder_behaviors[].limit_file_extensions": "object", "folder_behaviors[].limit_file_extensions.value": "object", "folder_behaviors[].limit_file_extensions.value.extensions": "list", "folder_behaviors[behavior=limit_file_extensions]": "object", "folder_behaviors[behavior=limit_file_extensions].value": "object", "folder_behaviors[behavior=limit_file_extensions].value.extensions": "list", "folder_behaviors[].limit_file_regex": "object", "folder_behaviors[].limit_file_regex.value": "list", "folder_behaviors[behavior=limit_file_regex]": "object", "folder_behaviors[behavior=limit_file_regex].value": "list"}, []string{"ai_feature_availability.in_app_ai_assistant.site_admins", "ai_feature_availability.in_app_ai_assistant.workspace_admins", "ai_feature_availability.in_app_ai_assistant.folder_admins", "ai_feature_availability.in_app_ai_assistant.all_users"}, "folder_behaviors[].behavior"),
					lib.KnownKeys("Files.com removes keys that are not policy settings when it saves the policy, so the stored value would not match. Remove them, or check their spelling.", "domain_hsts_header", "domain_letsencrypt_chain", "email", "reply_to_email", "allow_bundle_names", "bundle_expiration", "welcome_email_enabled", "ask_about_overwrites", "always_mkdir_parents", "welcome_email_cc", "welcome_email_subject", "welcome_custom_text", "language", "windows_mode_ftp", "default_time_zone", "desktop_app", "desktop_app_session_ip_pinning", "desktop_app_session_lifetime", "mobile_app", "mobile_app_session_ip_pinning", "mobile_app_session_lifetime", "folder_permissions_groups_only", "office_integration_available", "office_integration_type", "pin_all_remote_servers_to_site_region", "motd_text", "motd_use_for_ftp", "motd_use_for_sftp", "left_navigation_visibility", "disable_all_ai_features", "ai_feature_availability", "mcp_dcr_enabled", "additional_text_file_types", "bundle_require_note", "bundle_send_shared_receipts", "bundles_default_owned_by_primary_group", "calculate_file_checksums_crc32", "calculate_file_checksums_md5", "calculate_file_checksums_sha1", "calculate_file_checksums_sha256", "migrate_remote_server_sync_to_sync", "as2_message_retention_days", "file_activity_events_enabled", "username_display", "session_expiry_minutes", "ssl_required", "sftp_insecure_ciphers", "sftp_insecure_diffie_hellman", "fedramp", "user_lockout", "user_lockout_tries", "user_lockout_within", "user_lockout_lock_period", "include_password_in_welcome_email", "allowed_countries", "allowed_ips", "allow_user_level_2fa_override", "allow_user_level_allowed_ip_override", "allow_user_level_ssl_override", "enforce_parent_site_user_ip_restrictions", "disallowed_countries", "days_to_retain_backups", "max_prior_passwords", "password_validity_days", "password_min_length", "password_require_letter", "password_require_mixed", "password_require_special", "password_require_number", "password_require_unbreached", "require_logout_from_bundles_and_inboxes", "dav_user_root_enabled", "sftp_user_root_enabled", "disable_password_reset", "immutable_files", "bundle_not_found_message", "bundle_password_required", "bundle_require_registration", "bundle_require_share_recipient", "bundle_send_one_time_password_to_recipient_at_registration", "bundle_registration_notifications", "bundle_activity_notifications", "bundle_upload_receipt_notifications", "document_edits_in_bundle_allowed", "password_requirements_apply_to_bundles", "prevent_root_permissions_for_non_site_admins", "restrict_root_folder_behaviors_to_site_admins", "root_folder_behaviors_apply_to_workspaces", "opt_out_global", "use_provided_modified_at", "non_sso_groups_allowed", "non_sso_users_allowed", "sharing_enabled", "snapshot_sharing_enabled", "user_requests_enabled", "user_requests_notify_admins", "dav_enabled", "files_com_remote_server_enabled", "ftp_enabled", "s3_compatible_endpoint_enabled", "sftp_enabled", "sftp_finalize_partial_uploads", "users_can_create_api_keys", "users_can_create_ssh_keys", "show_user_notifications_log_in_link", "sftp_host_key_type", "active_sftp_host_key_id", "active_sftp_host_key_ids", "protocol_access_groups_only", "revoke_bundle_access_on_disable_or_delete", "bundle_watermark_value", "group_admins_can_add_users", "group_admins_can_manage_group_memberships", "group_admins_can_delete_users", "group_admins_can_enable_disable_users", "group_admins_can_modify_users", "group_admins_can_bypass_user_lifecycle_rules", "group_admins_can_reset_passwords", "group_admins_can_set_user_password", "child_site_folder_admins_can_see_parent_site_directory", "workspace_folder_admins_can_see_default_workspace_directory", "bundle_recipient_blacklist_free_email_domains", "bundle_recipient_blacklist_domains", "admins_bypass_locked_subfolders", "allowed_2fa_method_sms", "allowed_2fa_method_totp", "allowed_2fa_method_webauthn", "allowed_2fa_method_yubi", "allowed_2fa_method_email", "allowed_2fa_method_static", "allowed_2fa_method_bypass_for_ftp_sftp_dav", "require_2fa", "require_2fa_exempt_all_sso_users", "require_2fa_user_type", "color2_top", "color2_left", "color2_link", "color2_text", "color2_top_text", "site_header", "site_footer", "site_public_header", "site_public_footer", "login_help_text", "use_dedicated_ips_for_smtp", "email_footer_custom_text", "smtp_address", "smtp_authentication", "smtp_from", "smtp_username", "smtp_password", "smtp_port", "smtp_ssl", "ldap_enabled", "ldap_type", "ldap_host", "ldap_host_2", "ldap_host_3", "ldap_port", "ldap_secure", "ldap_server_certificate", "ldap_username", "ldap_username_field", "ldap_domain", "ldap_user_action", "ldap_group_action", "ldap_user_include_groups", "ldap_group_exclusion", "ldap_group_inclusion", "ldap_base_dn", "ldap_password_change", "ldap_password_change_confirmation", "uploads_via_email_authentication", "folder_behaviors"),
				},
				PlanModifiers: []planmodifier.Dynamic{
					dynamicplanmodifier.UseStateForUnknown(),
				},
			},
			"skip_child_site_ids": schema.ListAttribute{
				Description: "IDs of child sites excluded from this default policy.",
				Computed:    true,
				Optional:    true,
				ElementType: types.Int64Type,
				PlanModifiers: []planmodifier.List{
					listplanmodifier.UseStateForUnknown(),
				},
			},
			"child_site_ids": schema.ListAttribute{
				Description: "IDs of child sites explicitly assigned to this non-default policy.",
				Computed:    true,
				Optional:    true,
				ElementType: types.Int64Type,
				PlanModifiers: []planmodifier.List{
					listplanmodifier.UseStateForUnknown(),
				},
			},
			"default_policy": schema.BoolAttribute{
				Description: "Whether this policy applies to child sites not explicitly assigned to another policy.",
				Computed:    true,
				Optional:    true,
				PlanModifiers: []planmodifier.Bool{
					boolplanmodifier.UseStateForUnknown(),
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
	createValue, diags := lib.JSONValueToAPI(ctx, path.Root("value"), config.Value)
	resp.Diagnostics.Append(diags...)
	createValue, createValueTransformDiags0 := lib.UnwrapDiscriminatedUnionAtPath(ctx, path.Root("value"), createValue, []string{"folder_behaviors"}, "behavior", []lib.JSONSchemaVariant{{Name: "malware_scanning", Value: "malware_scanning"}, {Name: "file_expiration", Value: "file_expiration"}, {Name: "limit_file_extensions", Value: "limit_file_extensions"}, {Name: "limit_file_regex", Value: "limit_file_regex"}})
	resp.Diagnostics.Append(createValueTransformDiags0...)
	paramsChildSiteManagementPolicyCreate.Value = createValue
	if !plan.SkipChildSiteIds.IsNull() && !plan.SkipChildSiteIds.IsUnknown() {
		diags = plan.SkipChildSiteIds.ElementsAs(ctx, &paramsChildSiteManagementPolicyCreate.SkipChildSiteIds, false)
		resp.Diagnostics.Append(diags...)
	}
	if !plan.ChildSiteIds.IsNull() && !plan.ChildSiteIds.IsUnknown() {
		diags = plan.ChildSiteIds.ElementsAs(ctx, &paramsChildSiteManagementPolicyCreate.ChildSiteIds, false)
		resp.Diagnostics.Append(diags...)
	}
	if !plan.DefaultPolicy.IsNull() && !plan.DefaultPolicy.IsUnknown() {
		paramsChildSiteManagementPolicyCreate.DefaultPolicy = plan.DefaultPolicy.ValueBoolPointer()
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
	var state childSiteManagementPolicyResourceModel
	diags = req.State.Get(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	paramsChildSiteManagementPolicyUpdate := map[string]interface{}{}
	if !plan.Id.IsNull() && !plan.Id.IsUnknown() {
		paramsChildSiteManagementPolicyUpdate["id"] = plan.Id.ValueInt64()
	}
	if !config.Value.IsNull() && !config.Value.IsUnknown() {
		updateValue, diags := lib.JSONValueToAPI(ctx, path.Root("value"), config.Value)
		resp.Diagnostics.Append(diags...)
		updateValue, updateValueTransformDiags0 := lib.UnwrapDiscriminatedUnionAtPath(ctx, path.Root("value"), updateValue, []string{"folder_behaviors"}, "behavior", []lib.JSONSchemaVariant{{Name: "malware_scanning", Value: "malware_scanning"}, {Name: "file_expiration", Value: "file_expiration"}, {Name: "limit_file_extensions", Value: "limit_file_extensions"}, {Name: "limit_file_regex", Value: "limit_file_regex"}})
		resp.Diagnostics.Append(updateValueTransformDiags0...)
		paramsChildSiteManagementPolicyUpdate["value"] = updateValue
	}
	if !config.SkipChildSiteIds.IsNull() && !config.SkipChildSiteIds.IsUnknown() {
		var updateSkipChildSiteIds []int64
		diags = config.SkipChildSiteIds.ElementsAs(ctx, &updateSkipChildSiteIds, false)
		resp.Diagnostics.Append(diags...)
		paramsChildSiteManagementPolicyUpdate["skip_child_site_ids"] = updateSkipChildSiteIds
	}
	if !config.ChildSiteIds.IsNull() && !config.ChildSiteIds.IsUnknown() {
		var updateChildSiteIds []int64
		diags = config.ChildSiteIds.ElementsAs(ctx, &updateChildSiteIds, false)
		resp.Diagnostics.Append(diags...)
		paramsChildSiteManagementPolicyUpdate["child_site_ids"] = updateChildSiteIds
	}
	if !config.DefaultPolicy.IsNull() && !config.DefaultPolicy.IsUnknown() {
		paramsChildSiteManagementPolicyUpdate["default_policy"] = config.DefaultPolicy.ValueBool()
	}
	if !config.PolicyType.IsNull() && !config.PolicyType.IsUnknown() {
		paramsChildSiteManagementPolicyUpdate["policy_type"] = config.PolicyType.ValueString()
	}
	if !config.Name.IsNull() && !config.Name.IsUnknown() {
		paramsChildSiteManagementPolicyUpdate["name"] = config.Name.ValueString()
	}
	if !config.Description.IsNull() && !config.Description.IsUnknown() {
		paramsChildSiteManagementPolicyUpdate["description"] = config.Description.ValueString()
	}

	if resp.Diagnostics.HasError() {
		return
	}

	childSiteManagementPolicy, err := r.client.UpdateWithMap(paramsChildSiteManagementPolicyUpdate, files_sdk.WithContext(ctx))
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
	valueValue := interface{}(childSiteManagementPolicy.Value)
	valueValue, valueValueNestingDiags0 := lib.WrapDiscriminatedUnionLikeConfigAtPath(ctx, path.Root("value"), valueValue, state.Value, []string{"folder_behaviors"}, "behavior", []lib.JSONSchemaVariant{{Name: "malware_scanning", Value: "malware_scanning"}, {Name: "file_expiration", Value: "file_expiration"}, {Name: "limit_file_extensions", Value: "limit_file_extensions"}, {Name: "limit_file_regex", Value: "limit_file_regex"}})
	diags.Append(valueValueNestingDiags0...)
	state.Value, propDiags = lib.APIToDynamicJSON(ctx, path.Root("value"), valueValue, state.Value, []string{}, []string{}, "empty_object")
	diags.Append(propDiags...)
	state.AppliedChildSiteIds, propDiags = types.ListValueFrom(ctx, types.Int64Type, childSiteManagementPolicy.AppliedChildSiteIds)
	diags.Append(propDiags...)
	state.SkipChildSiteIds, propDiags = types.ListValueFrom(ctx, types.Int64Type, childSiteManagementPolicy.SkipChildSiteIds)
	diags.Append(propDiags...)
	state.ChildSiteIds, propDiags = types.ListValueFrom(ctx, types.Int64Type, childSiteManagementPolicy.ChildSiteIds)
	diags.Append(propDiags...)
	state.DefaultPolicy = types.BoolPointerValue(childSiteManagementPolicy.DefaultPolicy)
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
