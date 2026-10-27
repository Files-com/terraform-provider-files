package lib

import (
	"context"
	"encoding/json"
	"testing"

	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/schema/validator"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

const typedSchemaDate = "March 1, 2027"

func TestPublishedJSONTransitions(t *testing.T) {
	for _, test := range []struct {
		name, example, normalization   string
		maps, writeOnly, requiredPaths []string
		computed, deprecated           bool
		typed                          attr.Type
		valueTypes                     map[string]string
	}{
		{name: "as2_partner.additional_http_headers", example: "{\"X-Partner\":\"acme\"}",
			maps:          []string{""},
			requiredPaths: nil,
			normalization: "",
			computed:      false,
			deprecated:    true,
			typed:         types.MapType{ElemType: types.StringType},
			valueTypes:    nil},
		{name: "automation.import_urls", example: "[{\"name\":\"users.json\",\"url\":\"https://example.com/users\",\"method\":\"post\",\"headers\":{\"Content-Type\":\"application/json\"},\"content\":{\"group\":\"support\"}}]",
			maps:          []string{"*"},
			requiredPaths: nil,
			normalization: "",
			computed:      false,
			deprecated:    false,
			typed:         types.DynamicType,
			valueTypes:    nil},
		{name: "automation.value", example: "{\"limit\":\"1\"}",
			maps:          []string{""},
			requiredPaths: nil,
			normalization: "",
			computed:      false,
			deprecated:    true,
			typed:         types.MapType{ElemType: types.StringType},
			valueTypes:    nil},
		{name: "bundle.watermark_value", example: "{\"gravity\":\"SouthWest\",\"max_height_or_width\":20,\"transparency\":25}",
			maps:          []string{},
			requiredPaths: nil,
			normalization: "watermark",
			computed:      false,
			deprecated:    true,
			typed:         types.ObjectType{AttrTypes: map[string]attr.Type{"gravity": types.StringType, "max_height_or_width": types.Int64Type, "transparency": types.Int64Type, "dynamic_text": types.StringType}},
			valueTypes:    map[string]string{"max_height_or_width": "integer", "transparency": "integer"}},
		{name: "bundle.requested_upload_slots", example: "[{\"name\":\"Photo ID\"}]",
			maps:          []string{},
			requiredPaths: nil,
			normalization: "",
			computed:      true,
			deprecated:    true,
			typed:         types.ListType{ElemType: types.ObjectType{AttrTypes: map[string]attr.Type{"name": types.StringType}}},
			valueTypes:    map[string]string{"": "list", "[]": "object"}},
		{name: "child_site_management_policy.value", example: "{\"color2_left\":\"#000000\",\"folder_behaviors\":[{\"behavior\":\"file_expiration\",\"value\":{\"days_to_retain\":30}}]}",
			maps:          []string{},
			requiredPaths: []string{"ai_feature_availability.in_app_ai_assistant.site_admins", "ai_feature_availability.in_app_ai_assistant.workspace_admins", "ai_feature_availability.in_app_ai_assistant.folder_admins", "ai_feature_availability.in_app_ai_assistant.all_users"},
			normalization: "empty_object",
			computed:      false,
			deprecated:    true,
			typed:         types.ObjectType{AttrTypes: map[string]attr.Type{"domain_hsts_header": types.BoolType, "domain_letsencrypt_chain": types.StringType, "email": types.StringType, "reply_to_email": types.StringType, "allow_bundle_names": types.BoolType, "bundle_expiration": types.Int64Type, "welcome_email_enabled": types.BoolType, "ask_about_overwrites": types.BoolType, "always_mkdir_parents": types.BoolType, "welcome_email_cc": types.StringType, "welcome_email_subject": types.StringType, "welcome_custom_text": types.StringType, "language": types.StringType, "windows_mode_ftp": types.BoolType, "default_time_zone": types.StringType, "desktop_app": types.BoolType, "desktop_app_session_ip_pinning": types.BoolType, "desktop_app_session_lifetime": types.Int64Type, "mobile_app": types.BoolType, "mobile_app_session_ip_pinning": types.BoolType, "mobile_app_session_lifetime": types.Int64Type, "folder_permissions_groups_only": types.BoolType, "office_integration_available": types.BoolType, "office_integration_type": types.StringType, "pin_all_remote_servers_to_site_region": types.BoolType, "motd_text": types.StringType, "motd_use_for_ftp": types.BoolType, "motd_use_for_sftp": types.BoolType, "left_navigation_visibility": types.ObjectType{AttrTypes: map[string]attr.Type{"sharing": types.ObjectType{AttrTypes: map[string]attr.Type{"admin": types.BoolType, "all": types.BoolType}}, "automation_flow": types.ObjectType{AttrTypes: map[string]attr.Type{"admin": types.BoolType, "all": types.BoolType}}, "integrations": types.ObjectType{AttrTypes: map[string]attr.Type{"admin": types.BoolType, "all": types.BoolType}}, "clients_protocols": types.ObjectType{AttrTypes: map[string]attr.Type{"admin": types.BoolType, "all": types.BoolType}}}}, "disable_all_ai_features": types.BoolType, "ai_feature_availability": types.ObjectType{AttrTypes: map[string]attr.Type{"in_app_ai_assistant": types.ObjectType{AttrTypes: map[string]attr.Type{"site_admins": types.BoolType, "workspace_admins": types.BoolType, "folder_admins": types.BoolType, "all_users": types.BoolType, "selected_group_members": types.BoolType, "group_ids": types.ListType{ElemType: types.Int64Type}}}}}, "mcp_dcr_enabled": types.BoolType, "additional_text_file_types": types.ListType{ElemType: types.StringType}, "bundle_require_note": types.BoolType, "bundle_send_shared_receipts": types.BoolType, "bundles_default_owned_by_primary_group": types.BoolType, "calculate_file_checksums_crc32": types.BoolType, "calculate_file_checksums_md5": types.BoolType, "calculate_file_checksums_sha1": types.BoolType, "calculate_file_checksums_sha256": types.BoolType, "migrate_remote_server_sync_to_sync": types.BoolType, "as2_message_retention_days": types.Int64Type, "file_activity_events_enabled": types.BoolType, "username_display": types.StringType, "session_expiry_minutes": types.Int64Type, "ssl_required": types.BoolType, "sftp_insecure_ciphers": types.BoolType, "sftp_insecure_diffie_hellman": types.BoolType, "fedramp": types.BoolType, "user_lockout": types.BoolType, "user_lockout_tries": types.Int64Type, "user_lockout_within": types.Int64Type, "user_lockout_lock_period": types.Int64Type, "include_password_in_welcome_email": types.BoolType, "allowed_countries": types.StringType, "allowed_ips": types.StringType, "allow_user_level_2fa_override": types.BoolType, "allow_user_level_allowed_ip_override": types.BoolType, "allow_user_level_ssl_override": types.BoolType, "enforce_parent_site_user_ip_restrictions": types.BoolType, "disallowed_countries": types.StringType, "days_to_retain_backups": types.Int64Type, "max_prior_passwords": types.Int64Type, "password_validity_days": types.Int64Type, "password_min_length": types.Int64Type, "password_require_letter": types.BoolType, "password_require_mixed": types.BoolType, "password_require_special": types.BoolType, "password_require_number": types.BoolType, "password_require_unbreached": types.BoolType, "require_logout_from_bundles_and_inboxes": types.BoolType, "dav_user_root_enabled": types.BoolType, "sftp_user_root_enabled": types.BoolType, "disable_password_reset": types.BoolType, "immutable_files": types.BoolType, "bundle_not_found_message": types.StringType, "bundle_password_required": types.BoolType, "bundle_require_registration": types.BoolType, "bundle_require_share_recipient": types.BoolType, "bundle_send_one_time_password_to_recipient_at_registration": types.BoolType, "bundle_registration_notifications": types.StringType, "bundle_activity_notifications": types.StringType, "bundle_upload_receipt_notifications": types.StringType, "document_edits_in_bundle_allowed": types.BoolType, "password_requirements_apply_to_bundles": types.BoolType, "prevent_root_permissions_for_non_site_admins": types.BoolType, "restrict_root_folder_behaviors_to_site_admins": types.BoolType, "root_folder_behaviors_apply_to_workspaces": types.BoolType, "opt_out_global": types.BoolType, "use_provided_modified_at": types.BoolType, "non_sso_groups_allowed": types.BoolType, "non_sso_users_allowed": types.BoolType, "sharing_enabled": types.BoolType, "snapshot_sharing_enabled": types.BoolType, "user_requests_enabled": types.BoolType, "user_requests_notify_admins": types.BoolType, "dav_enabled": types.BoolType, "files_com_remote_server_enabled": types.BoolType, "ftp_enabled": types.BoolType, "s3_compatible_endpoint_enabled": types.BoolType, "sftp_enabled": types.BoolType, "sftp_finalize_partial_uploads": types.BoolType, "users_can_create_api_keys": types.BoolType, "users_can_create_ssh_keys": types.BoolType, "show_user_notifications_log_in_link": types.BoolType, "sftp_host_key_type": types.StringType, "active_sftp_host_key_id": types.Int64Type, "active_sftp_host_key_ids": types.ListType{ElemType: types.Int64Type}, "protocol_access_groups_only": types.BoolType, "revoke_bundle_access_on_disable_or_delete": types.BoolType, "bundle_watermark_value": types.ObjectType{AttrTypes: map[string]attr.Type{"gravity": types.StringType, "max_height_or_width": types.Int64Type, "transparency": types.Int64Type, "dynamic_text": types.StringType}}, "group_admins_can_add_users": types.BoolType, "group_admins_can_manage_group_memberships": types.BoolType, "group_admins_can_delete_users": types.BoolType, "group_admins_can_enable_disable_users": types.BoolType, "group_admins_can_modify_users": types.BoolType, "group_admins_can_bypass_user_lifecycle_rules": types.BoolType, "group_admins_can_reset_passwords": types.BoolType, "group_admins_can_set_user_password": types.BoolType, "child_site_folder_admins_can_see_parent_site_directory": types.BoolType, "workspace_folder_admins_can_see_default_workspace_directory": types.BoolType, "bundle_recipient_blacklist_free_email_domains": types.BoolType, "bundle_recipient_blacklist_domains": types.ListType{ElemType: types.StringType}, "admins_bypass_locked_subfolders": types.BoolType, "allowed_2fa_method_sms": types.BoolType, "allowed_2fa_method_totp": types.BoolType, "allowed_2fa_method_webauthn": types.BoolType, "allowed_2fa_method_yubi": types.BoolType, "allowed_2fa_method_email": types.BoolType, "allowed_2fa_method_static": types.BoolType, "allowed_2fa_method_bypass_for_ftp_sftp_dav": types.BoolType, "require_2fa": types.BoolType, "require_2fa_exempt_all_sso_users": types.BoolType, "require_2fa_user_type": types.StringType, "color2_top": types.StringType, "color2_left": types.StringType, "color2_link": types.StringType, "color2_text": types.StringType, "color2_top_text": types.StringType, "site_header": types.StringType, "site_footer": types.StringType, "site_public_header": types.StringType, "site_public_footer": types.StringType, "login_help_text": types.StringType, "use_dedicated_ips_for_smtp": types.BoolType, "email_footer_custom_text": types.StringType, "smtp_address": types.StringType, "smtp_authentication": types.StringType, "smtp_from": types.StringType, "smtp_username": types.StringType, "smtp_password": types.StringType, "smtp_port": types.Int64Type, "smtp_ssl": types.StringType, "ldap_enabled": types.BoolType, "ldap_type": types.StringType, "ldap_host": types.StringType, "ldap_host_2": types.StringType, "ldap_host_3": types.StringType, "ldap_port": types.Int64Type, "ldap_secure": types.BoolType, "ldap_server_certificate": types.StringType, "ldap_username": types.StringType, "ldap_username_field": types.StringType, "ldap_domain": types.StringType, "ldap_user_action": types.StringType, "ldap_group_action": types.StringType, "ldap_user_include_groups": types.StringType, "ldap_group_exclusion": types.StringType, "ldap_group_inclusion": types.StringType, "ldap_base_dn": types.StringType, "ldap_password_change": types.StringType, "ldap_password_change_confirmation": types.StringType, "uploads_via_email_authentication": types.BoolType, "folder_behaviors": types.ListType{ElemType: types.ObjectType{AttrTypes: map[string]attr.Type{"malware_scanning": types.ObjectType{AttrTypes: map[string]attr.Type{"value": types.ObjectType{AttrTypes: map[string]attr.Type{}}}}, "file_expiration": types.ObjectType{AttrTypes: map[string]attr.Type{"value": types.ObjectType{AttrTypes: map[string]attr.Type{"days_to_retain": types.Int64Type, "delete_empty_folders": types.BoolType}}}}, "limit_file_extensions": types.ObjectType{AttrTypes: map[string]attr.Type{"value": types.ObjectType{AttrTypes: map[string]attr.Type{"extensions": types.ListType{ElemType: types.StringType}, "mode": types.StringType}}}}, "limit_file_regex": types.ObjectType{AttrTypes: map[string]attr.Type{"value": types.ListType{ElemType: types.StringType}}}}}}}},
			valueTypes:    map[string]string{"domain_hsts_header": "boolean", "allow_bundle_names": "boolean", "bundle_expiration": "integer", "welcome_email_enabled": "boolean", "ask_about_overwrites": "boolean", "always_mkdir_parents": "boolean", "windows_mode_ftp": "boolean", "desktop_app": "boolean", "desktop_app_session_ip_pinning": "boolean", "desktop_app_session_lifetime": "integer", "mobile_app": "boolean", "mobile_app_session_ip_pinning": "boolean", "mobile_app_session_lifetime": "integer", "folder_permissions_groups_only": "boolean", "office_integration_available": "boolean", "pin_all_remote_servers_to_site_region": "boolean", "motd_use_for_ftp": "boolean", "motd_use_for_sftp": "boolean", "left_navigation_visibility": "object", "left_navigation_visibility.sharing": "object", "left_navigation_visibility.sharing.admin": "boolean", "left_navigation_visibility.sharing.all": "boolean", "left_navigation_visibility.automation_flow": "object", "left_navigation_visibility.automation_flow.admin": "boolean", "left_navigation_visibility.automation_flow.all": "boolean", "left_navigation_visibility.integrations": "object", "left_navigation_visibility.integrations.admin": "boolean", "left_navigation_visibility.integrations.all": "boolean", "left_navigation_visibility.clients_protocols": "object", "left_navigation_visibility.clients_protocols.admin": "boolean", "left_navigation_visibility.clients_protocols.all": "boolean", "disable_all_ai_features": "boolean", "ai_feature_availability": "object", "ai_feature_availability.in_app_ai_assistant": "object", "ai_feature_availability.in_app_ai_assistant.site_admins": "boolean", "ai_feature_availability.in_app_ai_assistant.workspace_admins": "boolean", "ai_feature_availability.in_app_ai_assistant.folder_admins": "boolean", "ai_feature_availability.in_app_ai_assistant.all_users": "boolean", "ai_feature_availability.in_app_ai_assistant.selected_group_members": "boolean", "ai_feature_availability.in_app_ai_assistant.group_ids": "list", "ai_feature_availability.in_app_ai_assistant.group_ids[]": "integer", "mcp_dcr_enabled": "boolean", "additional_text_file_types": "list", "bundle_require_note": "boolean", "bundle_send_shared_receipts": "boolean", "bundles_default_owned_by_primary_group": "boolean", "calculate_file_checksums_crc32": "boolean", "calculate_file_checksums_md5": "boolean", "calculate_file_checksums_sha1": "boolean", "calculate_file_checksums_sha256": "boolean", "migrate_remote_server_sync_to_sync": "boolean", "as2_message_retention_days": "integer", "file_activity_events_enabled": "boolean", "session_expiry_minutes": "integer", "ssl_required": "boolean", "sftp_insecure_ciphers": "boolean", "sftp_insecure_diffie_hellman": "boolean", "fedramp": "boolean", "user_lockout": "boolean", "user_lockout_tries": "integer", "user_lockout_within": "integer", "user_lockout_lock_period": "integer", "include_password_in_welcome_email": "boolean", "allow_user_level_2fa_override": "boolean", "allow_user_level_allowed_ip_override": "boolean", "allow_user_level_ssl_override": "boolean", "enforce_parent_site_user_ip_restrictions": "boolean", "days_to_retain_backups": "integer", "max_prior_passwords": "integer", "password_validity_days": "integer", "password_min_length": "integer", "password_require_letter": "boolean", "password_require_mixed": "boolean", "password_require_special": "boolean", "password_require_number": "boolean", "password_require_unbreached": "boolean", "require_logout_from_bundles_and_inboxes": "boolean", "dav_user_root_enabled": "boolean", "sftp_user_root_enabled": "boolean", "disable_password_reset": "boolean", "immutable_files": "boolean", "bundle_password_required": "boolean", "bundle_require_registration": "boolean", "bundle_require_share_recipient": "boolean", "bundle_send_one_time_password_to_recipient_at_registration": "boolean", "document_edits_in_bundle_allowed": "boolean", "password_requirements_apply_to_bundles": "boolean", "prevent_root_permissions_for_non_site_admins": "boolean", "restrict_root_folder_behaviors_to_site_admins": "boolean", "root_folder_behaviors_apply_to_workspaces": "boolean", "opt_out_global": "boolean", "use_provided_modified_at": "boolean", "non_sso_groups_allowed": "boolean", "non_sso_users_allowed": "boolean", "sharing_enabled": "boolean", "snapshot_sharing_enabled": "boolean", "user_requests_enabled": "boolean", "user_requests_notify_admins": "boolean", "dav_enabled": "boolean", "files_com_remote_server_enabled": "boolean", "ftp_enabled": "boolean", "s3_compatible_endpoint_enabled": "boolean", "sftp_enabled": "boolean", "sftp_finalize_partial_uploads": "boolean", "users_can_create_api_keys": "boolean", "users_can_create_ssh_keys": "boolean", "show_user_notifications_log_in_link": "boolean", "active_sftp_host_key_id": "integer", "active_sftp_host_key_ids": "list", "active_sftp_host_key_ids[]": "integer", "protocol_access_groups_only": "boolean", "revoke_bundle_access_on_disable_or_delete": "boolean", "bundle_watermark_value": "object", "bundle_watermark_value.max_height_or_width": "integer", "bundle_watermark_value.transparency": "integer", "group_admins_can_add_users": "boolean", "group_admins_can_manage_group_memberships": "boolean", "group_admins_can_delete_users": "boolean", "group_admins_can_enable_disable_users": "boolean", "group_admins_can_modify_users": "boolean", "group_admins_can_bypass_user_lifecycle_rules": "boolean", "group_admins_can_reset_passwords": "boolean", "group_admins_can_set_user_password": "boolean", "child_site_folder_admins_can_see_parent_site_directory": "boolean", "workspace_folder_admins_can_see_default_workspace_directory": "boolean", "bundle_recipient_blacklist_free_email_domains": "boolean", "bundle_recipient_blacklist_domains": "list", "admins_bypass_locked_subfolders": "boolean", "allowed_2fa_method_sms": "boolean", "allowed_2fa_method_totp": "boolean", "allowed_2fa_method_webauthn": "boolean", "allowed_2fa_method_yubi": "boolean", "allowed_2fa_method_email": "boolean", "allowed_2fa_method_static": "boolean", "allowed_2fa_method_bypass_for_ftp_sftp_dav": "boolean", "require_2fa": "boolean", "require_2fa_exempt_all_sso_users": "boolean", "use_dedicated_ips_for_smtp": "boolean", "smtp_port": "integer", "ldap_enabled": "boolean", "ldap_port": "integer", "ldap_secure": "boolean", "uploads_via_email_authentication": "boolean", "folder_behaviors": "list", "folder_behaviors[]": "object", "folder_behaviors[].malware_scanning": "object", "folder_behaviors[].malware_scanning.value": "object", "folder_behaviors[behavior=malware_scanning]": "object", "folder_behaviors[behavior=malware_scanning].value": "object", "folder_behaviors[].file_expiration": "object", "folder_behaviors[].file_expiration.value": "object", "folder_behaviors[].file_expiration.value.days_to_retain": "integer", "folder_behaviors[].file_expiration.value.delete_empty_folders": "boolean", "folder_behaviors[behavior=file_expiration]": "object", "folder_behaviors[behavior=file_expiration].value": "object", "folder_behaviors[behavior=file_expiration].value.days_to_retain": "integer", "folder_behaviors[behavior=file_expiration].value.delete_empty_folders": "boolean", "folder_behaviors[].limit_file_extensions": "object", "folder_behaviors[].limit_file_extensions.value": "object", "folder_behaviors[].limit_file_extensions.value.extensions": "list", "folder_behaviors[behavior=limit_file_extensions]": "object", "folder_behaviors[behavior=limit_file_extensions].value": "object", "folder_behaviors[behavior=limit_file_extensions].value.extensions": "list", "folder_behaviors[].limit_file_regex": "object", "folder_behaviors[].limit_file_regex.value": "list", "folder_behaviors[behavior=limit_file_regex]": "object", "folder_behaviors[behavior=limit_file_regex].value": "list"}},
		{name: "desktop_configuration_profile.mount_mappings", example: "{\"W\":\"Americas\"}",
			maps:          []string{""},
			requiredPaths: nil,
			normalization: "mount_mappings",
			computed:      false,
			deprecated:    true,
			typed:         types.MapType{ElemType: types.StringType},
			valueTypes:    nil},
		{name: "event_subscription.filter", example: "{\"payload.sync_name\":{\"prefix\":\"Nightly\"},\"payload.event_errors\":{\"exists\":true}}",
			maps:          []string{"*"},
			requiredPaths: nil,
			normalization: "",
			computed:      false,
			deprecated:    false,
			typed:         types.DynamicType,
			valueTypes:    nil},
		{name: "event_subscription.delivery_policy", example: "{\"max_attempts\":3,\"retry_delay_minutes\":[5,30]}",
			maps:          []string{},
			requiredPaths: nil,
			normalization: "empty_object",
			computed:      false,
			deprecated:    true,
			typed:         types.ObjectType{AttrTypes: map[string]attr.Type{"max_attempts": types.Int64Type, "retry_delay_minutes": types.ListType{ElemType: types.Int64Type}}},
			valueTypes:    map[string]string{"max_attempts": "integer", "retry_delay_minutes": "list", "retry_delay_minutes[]": "integer"}},
		{name: "event_target.config", example: "{\"url\":\"https://example.com/files-events\",\"headers\":{\"X-Example\":\"value\"},\"use_dedicated_ips\":false}",
			maps:          []string{"headers"},
			requiredPaths: nil,
			normalization: "event_target_config",
			computed:      false,
			deprecated:    true,
			typed:         types.ObjectType{AttrTypes: map[string]attr.Type{"user_ids": types.ListType{ElemType: types.Int64Type}, "url": types.StringType, "headers": types.MapType{ElemType: types.StringType}, "signing_secret": types.StringType, "use_dedicated_ips": types.BoolType, "username": types.StringType, "channel": types.StringType, "icon_emoji": types.StringType, "arn": types.StringType, "aws_credentials": types.ObjectType{AttrTypes: map[string]attr.Type{"access_key_id": types.StringType, "region": types.StringType, "secret_access_key": types.StringType}}, "project_id": types.StringType, "topic_id": types.StringType, "google_credentials": types.ObjectType{AttrTypes: map[string]attr.Type{"type": types.StringType, "project_id": types.StringType, "private_key_id": types.StringType, "private_key": types.StringType, "client_email": types.StringType, "client_id": types.StringType, "auth_uri": types.StringType, "token_uri": types.StringType, "auth_provider_x509_cert_url": types.StringType, "client_x509_cert_url": types.StringType, "universe_domain": types.StringType, "quota_project_id": types.StringType}}, "path": types.StringType, "format": types.StringType}},
			valueTypes:    map[string]string{"user_ids": "list", "user_ids[]": "integer", "use_dedicated_ips": "boolean", "aws_credentials": "object", "google_credentials": "object"}},
		{name: "event_target.delivery_policy", example: "{\"max_attempts\":3,\"retry_delay_minutes\":[1,10]}",
			maps:          []string{},
			requiredPaths: nil,
			normalization: "event_target_delivery_policy",
			computed:      false,
			deprecated:    true,
			typed:         types.ObjectType{AttrTypes: map[string]attr.Type{"max_attempts": types.Int64Type, "retry_delay_minutes": types.ListType{ElemType: types.Int64Type}, "batch_interval": types.Int64Type, "delivery_mode": types.StringType, "summary_path": types.StringType, "summary_delay": types.Int64Type}},
			valueTypes:    map[string]string{"max_attempts": "integer", "retry_delay_minutes": "list", "retry_delay_minutes[]": "integer", "batch_interval": "integer", "summary_delay": "integer"}},
		{name: "expectation.criteria", example: "{\"count\":{\"exact\":1},\"extensions\":[\"csv\"]}",
			maps:          []string{"required_files"},
			requiredPaths: nil,
			normalization: "empty_object",
			computed:      false,
			deprecated:    true,
			typed:         types.ObjectType{AttrTypes: map[string]attr.Type{"count": types.ObjectType{AttrTypes: map[string]attr.Type{"exact": types.Int64Type, "max": types.Int64Type, "min": types.Int64Type}}, "extensions": types.ListType{ElemType: types.StringType}, "filename_regex": types.StringType, "total_bytes": types.ObjectType{AttrTypes: map[string]attr.Type{"exact": types.Int64Type, "max": types.Int64Type, "min": types.Int64Type}}, "forbidden_files": types.ListType{ElemType: types.StringType}, "required_files": types.MapType{ElemType: types.ObjectType{AttrTypes: map[string]attr.Type{"count": types.ObjectType{AttrTypes: map[string]attr.Type{"exact": types.Int64Type, "max": types.Int64Type, "min": types.Int64Type}}, "extensions": types.ListType{ElemType: types.StringType}, "filename_regex": types.StringType, "size_bytes": types.ObjectType{AttrTypes: map[string]attr.Type{"exact": types.Int64Type, "max": types.Int64Type, "min": types.Int64Type}}}}}, "content_validation": types.ObjectType{AttrTypes: map[string]attr.Type{"mode": types.StringType, "fts": types.StringType}}}},
			valueTypes:    map[string]string{"count": "object", "count.exact": "integer", "count.max": "integer", "count.min": "integer", "extensions": "list", "total_bytes": "object", "total_bytes.exact": "integer", "total_bytes.max": "integer", "total_bytes.min": "integer", "forbidden_files": "list", "required_files.*": "object", "required_files.*.count": "object", "required_files.*.count.exact": "integer", "required_files.*.count.max": "integer", "required_files.*.count.min": "integer", "required_files.*.extensions": "list", "required_files.*.size_bytes": "object", "required_files.*.size_bytes.exact": "integer", "required_files.*.size_bytes.max": "integer", "required_files.*.size_bytes.min": "integer", "content_validation": "object"}},
		{name: "file.custom_metadata", example: "{\"department\":\"finance\"}",
			maps:          []string{""},
			requiredPaths: nil,
			normalization: "empty_object",
			computed:      false,
			deprecated:    true,
			typed:         types.MapType{ElemType: types.StringType},
			valueTypes:    nil},
		{name: "folder.custom_metadata", example: "{\"department\":\"finance\"}",
			maps:          []string{""},
			requiredPaths: nil,
			normalization: "empty_object",
			computed:      false,
			deprecated:    true,
			typed:         types.MapType{ElemType: types.StringType},
			valueTypes:    nil},
		{name: "form_field_set.form_fields", example: "[{\"label\":\"Department\",\"field_type\":\"dropdown\",\"required\":true,\"options_for_select\":[\"Finance\",\"Sales\"]}]",
			maps:          []string{},
			requiredPaths: nil,
			normalization: "form_fields",
			computed:      false,
			deprecated:    true,
			typed:         types.ListType{ElemType: types.ObjectType{AttrTypes: map[string]attr.Type{"id": types.Int64Type, "label": types.StringType, "field_type": types.StringType, "help_text": types.StringType, "required": types.BoolType, "options_for_select": types.ListType{ElemType: types.StringType}, "default_option": types.StringType, "form_field_set_id": types.Int64Type}}},
			valueTypes:    map[string]string{"": "list", "[]": "object", "[].id": "integer", "[].required": "boolean", "[].options_for_select": "list"}},
		{name: "integration_centric_profile.expected_remote_servers", example: "[{\"server_type\":\"dropbox\",\"name\":\"Dropbox\"}]",
			maps:          []string{},
			requiredPaths: nil,
			normalization: "expected_remote_servers",
			computed:      false,
			deprecated:    true,
			typed:         types.ListType{ElemType: types.ObjectType{AttrTypes: map[string]attr.Type{"server_type": types.StringType, "name": types.StringType}}},
			valueTypes:    map[string]string{"": "list", "[]": "object"}},
		{name: "scheduled_export.export_options", example: "{\"group_by\":\"user\"}",
			maps:          []string{},
			requiredPaths: nil,
			normalization: "empty_object",
			computed:      false,
			deprecated:    true,
			typed:         types.ObjectType{AttrTypes: map[string]attr.Type{"path": types.StringType, "max_depth": types.Int64Type, "size_threshold": types.Int64Type, "all_folders": types.BoolType, "group_by": types.StringType}},
			valueTypes:    map[string]string{"max_depth": "integer", "size_threshold": "integer", "all_folders": "boolean"}},
		{name: "secret.metadata", example: "{\"header_name\":\"Authorization\"}",
			maps:          []string{},
			requiredPaths: nil,
			normalization: "empty_object",
			computed:      false,
			deprecated:    true,
			typed:         types.ObjectType{AttrTypes: map[string]attr.Type{"username": types.StringType, "header_name": types.StringType, "query_parameter_name": types.StringType}},
			valueTypes:    nil},
		{name: "share_group.members", example: "[{\"name\":\"John Doe\",\"company\":\"Acme Ltd\",\"email\":\"johndoe@example.com\"}]",
			maps:          []string{},
			requiredPaths: nil,
			normalization: "",
			computed:      false,
			deprecated:    true,
			typed:         types.ListType{ElemType: types.ObjectType{AttrTypes: map[string]attr.Type{"email": types.StringType, "name": types.StringType, "company": types.StringType}}},
			valueTypes:    map[string]string{"": "list", "[]": "object"}},
		{name: "siem_http_destination.additional_headers", example: "{\"Authorization\":\"Bearer YOUR_TOKEN\"}",
			maps:          []string{""},
			requiredPaths: nil,
			normalization: "",
			computed:      false,
			deprecated:    true,
			typed:         types.MapType{ElemType: types.StringType},
			valueTypes:    nil},
		{name: "site.ai_feature_availability", example: "{\"in_app_ai_assistant\":{\"site_admins\":true,\"workspace_admins\":true,\"folder_admins\":true,\"all_users\":false,\"selected_group_members\":false}}",
			maps:          []string{},
			requiredPaths: []string{"in_app_ai_assistant.site_admins", "in_app_ai_assistant.workspace_admins", "in_app_ai_assistant.folder_admins", "in_app_ai_assistant.all_users"},
			normalization: "",
			computed:      false,
			deprecated:    true,
			typed:         types.ObjectType{AttrTypes: map[string]attr.Type{"in_app_ai_assistant": types.ObjectType{AttrTypes: map[string]attr.Type{"site_admins": types.BoolType, "workspace_admins": types.BoolType, "folder_admins": types.BoolType, "all_users": types.BoolType, "selected_group_members": types.BoolType, "group_ids": types.ListType{ElemType: types.Int64Type}}}}},
			valueTypes:    map[string]string{"in_app_ai_assistant": "object", "in_app_ai_assistant.site_admins": "boolean", "in_app_ai_assistant.workspace_admins": "boolean", "in_app_ai_assistant.folder_admins": "boolean", "in_app_ai_assistant.all_users": "boolean", "in_app_ai_assistant.selected_group_members": "boolean", "in_app_ai_assistant.group_ids": "list", "in_app_ai_assistant.group_ids[]": "integer"}},
		{name: "site.bundle_watermark_value", example: "{\"gravity\":\"SouthWest\",\"max_height_or_width\":20,\"transparency\":25}",
			maps:          []string{},
			requiredPaths: nil,
			normalization: "watermark",
			computed:      false,
			deprecated:    true,
			typed:         types.ObjectType{AttrTypes: map[string]attr.Type{"gravity": types.StringType, "max_height_or_width": types.Int64Type, "transparency": types.Int64Type, "dynamic_text": types.StringType}},
			valueTypes:    map[string]string{"max_height_or_width": "integer", "transparency": "integer"}},
	} {
		t.Run(test.name, func(t *testing.T) {
			ctx, attributePath := context.Background(), path.Root("value")
			var valueValidator validator.Dynamic = ValidJSONEncoding()
			if test.deprecated {
				valueValidator = DeprecatedJSONEncoding(test.name, test.example, typedSchemaDate, test.valueTypes, test.requiredPaths)
			}
			var example any
			require.NoError(t, json.Unmarshal([]byte(test.example), &example))
			var empty any = map[string]any{}
			if _, list := example.([]any); list {
				empty = []any{}
			}
			for _, value := range []any{example, nil, empty} {
				for _, encoded := range []bool{false, true} {
					if test.computed {
						continue
					}
					prior := types.DynamicValue(nativeJSONValue(t, value))
					if value == nil {
						prior = types.DynamicNull()
					}
					if encoded {
						text, err := json.Marshal(value)
						require.NoError(t, err)
						prior = types.DynamicValue(types.StringValue(string(text)))
					}
					request, diags := JSONValueToAPI(ctx, attributePath, prior)
					require.False(t, diags.HasError(), diags)
					assert.Equal(t, value, request)
					validation := &validator.DynamicResponse{}
					valueValidator.ValidateDynamic(ctx, validator.DynamicRequest{Path: attributePath, ConfigValue: prior}, validation)
					assert.False(t, validation.Diagnostics.HasError(), validation.Diagnostics)
					assert.Equal(t, encoded && test.deprecated, validation.Diagnostics.WarningsCount() == 1)
					if encoded && test.deprecated {
						assert.Contains(t, validation.Diagnostics[0].Detail(), typedSchemaDate)
					}
					state, diags := APIToDynamicJSON(ctx, attributePath, request, prior, test.maps, test.writeOnly, test.normalization)
					require.False(t, diags.HasError(), diags)
					assert.True(t, prior.Equal(state), "apply must preserve the HCL type and representation")
					refresh, diags := APIToDynamicJSON(ctx, attributePath, request, state, test.maps, test.writeOnly, test.normalization)
					require.False(t, diags.HasError(), diags)
					assert.True(t, state.Equal(refresh))
				}
			}
			for _, unknown := range []types.Dynamic{types.DynamicUnknown(), types.DynamicValue(types.StringUnknown())} {
				validation := &validator.DynamicResponse{}
				valueValidator.ValidateDynamic(ctx, validator.DynamicRequest{Path: attributePath, ConfigValue: unknown}, validation)
				assert.Empty(t, validation.Diagnostics)
				value, diags := JSONValueToAPI(ctx, attributePath, unknown)
				require.False(t, diags.HasError(), diags)
				assert.Nil(t, value)
			}
			responses := []any{example, nil, empty}
			if test.name == "bundle.watermark_value" {
				responses = append(responses, map[string]any{"transparency": float64(25)}, map[string]any{"transparency": "25"}, map[string]any{"transparency": ""}, map[string]any{"transparency": nil}, map[string]any{"extra": "kept"})
			}
			for _, response := range responses {
				state, diags := APIToDynamicJSON(ctx, attributePath, response, types.DynamicNull(), test.maps, test.writeOnly, test.normalization)
				require.False(t, diags.HasError(), diags)
				value, diags := JSONValueToAPI(ctx, attributePath, state)
				require.False(t, diags.HasError(), diags)
				assert.Equal(t, response, value, "imports and computed reads must not add wrappers or coerce values")
				encoded, err := json.Marshal(response)
				require.NoError(t, err)
				for _, prior := range []types.Dynamic{state, types.DynamicValue(types.StringValue(string(encoded)))} {
					value, diags := JSONValueToAPI(ctx, attributePath, prior)
					require.False(t, diags.HasError(), diags)
					_, diags = schemaValue(ctx, attributePath, value, test.typed)
					require.False(t, diags.HasError(), diags)
				}
			}
		})
	}
}
