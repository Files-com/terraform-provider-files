package lib

import (
	"context"
	"encoding/json"
	"testing"

	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/schema/validator"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

var generatedBehaviorValueVariants = []JSONSchemaVariant{{Name: "webhook", Value: "webhook", RootType: "object", Writable: []string{"urls", "method", "triggers", "triggering_filenames", "exclude_filenames", "encoding", "headers", "body", "verification_token", "file_form_field", "file_as_body", "use_dedicated_ips"}, LegacyProperty: "urls", LegacyList: true}, {Name: "file_expiration", Value: "file_expiration", RootType: "object", Writable: []string{"days_to_retain", "delete_empty_folders"}, LegacyProperty: "days_to_retain"}, {Name: "auto_encrypt", Value: "auto_encrypt", RootType: "object", Writable: []string{"gpg_key_id", "gpg_key_ids", "algorithm", "signing_key_id", "suffix", "armor", "gpg_key_partner_id"}}, {Name: "lock_subfolders", Value: "lock_subfolders", RootType: "object", Writable: []string{"level"}}, {Name: "storage_region", Value: "storage_region", RootType: "string"}, {Name: "serve_publicly", Value: "serve_publicly", RootType: "object", Writable: []string{"key", "show_index", "force_download", "username", "password", "cors_enabled", "require_site_authentication"}, WriteOnly: []string{"password"}}, {Name: "create_user_folders", Value: "create_user_folders", RootType: "object", Writable: []string{"permission", "additional_permission", "existing_users", "group_id", "new_folder_name", "subfolders"}}, {Name: "inbox", Value: "inbox", RootType: "object", Writable: []string{"key", "dont_separate_submissions_by_folder", "dont_separate_submissions_by_folder_for_inbound_email", "dont_allow_folders_in_uploads", "require_inbox_recipient", "show_on_login_page", "clickwrap_id", "form_field_set_id", "title", "description", "help_text", "require_registration", "password", "path_template", "path_template_time_zone", "enable_inbound_email_address", "notify_senders_on_successful_uploads_via_email", "notify_senders_on_successful_uploads_via_web", "allow_whitelisting", "whitelist", "disable_web_upload", "capture_email_body_filename", "requested_upload_slots"}, WriteOnly: []string{"password", "enable_inbound_email_address"}}, {Name: "limit_file_extensions", Value: "limit_file_extensions", RootType: "object", Writable: []string{"extensions", "mode"}}, {Name: "limit_file_regex", Value: "limit_file_regex", RootType: "array"}, {Name: "amazon_sns", Value: "amazon_sns", RootType: "object", Writable: []string{"arns", "triggers", "aws_credentials", "body"}, WriteOnly: []string{"aws_credentials.secret_access_key"}}, {Name: "watermark", Value: "watermark", RootType: "object", Writable: []string{"gravity", "max_height_or_width", "transparency", "dynamic_text"}}, {Name: "remote_server_mount", Value: "remote_server_mount", RootType: "object", Writable: []string{"remote_server_id", "remote_path"}}, {Name: "slack_webhook", Value: "slack_webhook", RootType: "object", Writable: []string{"url", "username", "channel", "icon_emoji", "triggers"}}, {Name: "auto_decrypt", Value: "auto_decrypt", RootType: "object", Writable: []string{"gpg_key_id", "gpg_key_ids", "algorithm", "suffix", "ignore_mdc_error", "gpg_key_partner_id", "use_all_private_keys"}}, {Name: "override_upload_filename", Value: "override_upload_filename", RootType: "object", Writable: []string{"filename_override_pattern", "filename_replace_from", "filename_replace_to", "filename_regex_replace_from", "filename_regex_replace_to", "time_zone"}}, {Name: "permission_fence", Value: "permission_fence", RootType: "object", Writable: []string{"fenced_permissions"}}, {Name: "limit_filename_length", Value: "limit_filename_length", RootType: "object", Writable: []string{"max_length", "shorten"}}, {Name: "organize_files_into_subfolders", Value: "organize_files_into_subfolders", RootType: "object", Writable: []string{"subfolder_name_type", "regex", "strftime_format", "time_zone", "apply_behavior"}, WriteOnly: []string{"apply_behavior"}}, {Name: "teams_webhook", Value: "teams_webhook", RootType: "object", Writable: []string{"url", "triggers"}}, {Name: "google_pub_sub", Value: "google_pub_sub", RootType: "object", Writable: []string{"projects_topics", "triggers", "google_credentials", "body"}, WriteOnly: []string{"google_credentials.private_key"}}, {Name: "archive_overwritten_or_deleted_files", Value: "archive_overwritten_or_deleted_files", RootType: "object", Writable: []string{"archive_path"}}, {Name: "auto_recrypt", Value: "auto_recrypt", RootType: "object", Writable: []string{"decrypt_gpg_key_ids", "encrypt_gpg_key_ids", "decrypt_gpg_key_partner_id", "encrypt_gpg_key_partner_id", "ignore_mdc_error", "signing_key_id", "armor"}}, {Name: "metadata_category", Value: "metadata_category", RootType: "object", Writable: []string{"metadata_category_id"}}, {Name: "auto_unzip", Value: "auto_unzip", RootType: "object", Writable: []string{"destination_path", "path_time_zone"}}, {Name: "remote_server_metadata_index", Value: "remote_server_metadata_index", RootType: "object", Writable: []string{"interval_minutes"}}, {Name: "malware_scanning", Value: "malware_scanning", RootType: "object"}}

func TestGeneratedBehaviorValueTransitionMatrix(t *testing.T) {
	tests := []struct {
		name     string
		apiValue string
	}{
		{name: "webhook", apiValue: "{\"urls\":[\"https://example.com/webhook\"],\"method\":\"POST\",\"encoding\":\"JSON\"}"},
		{name: "file_expiration", apiValue: "{\"days_to_retain\":30,\"delete_empty_folders\":false}"},
		{name: "auto_encrypt", apiValue: "{\"gpg_key_ids\":[1],\"algorithm\":\"PGP/GPG\",\"suffix\":\".gpg\",\"armor\":false}"},
		{name: "lock_subfolders", apiValue: "{\"level\":\"children_recursive\"}"},
		{name: "storage_region", apiValue: "\"us-east-1\""},
		{name: "serve_publicly", apiValue: "{\"key\":\"public-files\",\"show_index\":true,\"force_download\":false}"},
		{name: "create_user_folders", apiValue: "{\"permission\":\"full\",\"existing_users\":false,\"new_folder_name\":\"name\"}"},
		{name: "inbox", apiValue: "{\"key\":\"application-forms\",\"dont_separate_submissions_by_folder\":false,\"show_on_login_page\":false,\"title\":\"Application Forms\",\"require_registration\":false,\"disable_web_upload\":false,\"requested_upload_slots\":[{\"name\":\"Photo ID\"}]}"},
		{name: "limit_file_extensions", apiValue: "{\"extensions\":[\"pdf\",\"csv\"],\"mode\":\"whitelist\"}"},
		{name: "limit_file_regex", apiValue: "[\"/Document-.*/\"]"},
		{name: "amazon_sns", apiValue: "{\"arns\":[\"arn:aws:sns:us-east-1:123456789012:files-events\"],\"aws_credentials\":{\"access_key_id\":\"ACCESS_KEY_ID\",\"region\":\"us-east-1\",\"secret_access_key\":\"SECRET_ACCESS_KEY\"}}"},
		{name: "watermark", apiValue: "{\"gravity\":\"SouthWest\",\"max_height_or_width\":20,\"transparency\":25}"},
		{name: "remote_server_mount", apiValue: "{\"remote_server_id\":1,\"remote_path\":\"shared/files\"}"},
		{name: "slack_webhook", apiValue: "{\"url\":\"https://hooks.slack.com/services/example\",\"triggers\":[\"create\"]}"},
		{name: "auto_decrypt", apiValue: "{\"gpg_key_ids\":[1],\"algorithm\":\"PGP/GPG\",\"suffix\":\".gpg\",\"ignore_mdc_error\":false}"},
		{name: "override_upload_filename", apiValue: "{\"filename_override_pattern\":\"%Fb_uploaded%Fe\"}"},
		{name: "permission_fence", apiValue: "{\"fenced_permissions\":\"all\"}"},
		{name: "limit_filename_length", apiValue: "{\"max_length\":30,\"shorten\":true}"},
		{name: "organize_files_into_subfolders", apiValue: "{\"subfolder_name_type\":\"extension\"}"},
		{name: "teams_webhook", apiValue: "{\"url\":\"https://example.webhook.office.com/webhook\",\"triggers\":[\"create\"]}"},
		{name: "google_pub_sub", apiValue: "{\"projects_topics\":[{\"project_id\":\"my-project\",\"topic_id\":\"files-events\"}],\"google_credentials\":{\"type\":\"service_account\",\"project_id\":\"your-project-id\",\"private_key_id\":\"your-private-key-id\",\"private_key\":\"-----BEGIN PRIVATE KEY-----\\\\nMIIC...\",\"client_email\":\"your-service-account@your-project-id.iam.gserviceaccount.com\",\"client_id\":\"your-client-id\",\"auth_uri\":\"https://accounts.google.com/o/oauth2/auth\",\"token_uri\":\"https://oauth2.googleapis.com/token\",\"auth_provider_x509_cert_url\":\"https://www.googleapis.com/oauth2/v1/certs\",\"client_x509_cert_url\":\"https://www.googleapis.com/robot/v1/metadata/x509/your-service-account%40your-project-id.iam.gserviceaccount.com\"}}"},
		{name: "archive_overwritten_or_deleted_files", apiValue: "{\"archive_path\":\"/Archive\"}"},
		{name: "auto_recrypt", apiValue: "{\"decrypt_gpg_key_ids\":[1],\"encrypt_gpg_key_ids\":[2],\"ignore_mdc_error\":false,\"armor\":false}"},
		{name: "metadata_category", apiValue: "{\"metadata_category_id\":1}"},
		{name: "auto_unzip", apiValue: "{\"destination_path\":\"/Uploads/Unzipped/%Y/%m/%d\"}"},
		{name: "remote_server_metadata_index", apiValue: "{\"interval_minutes\":1440}"},
		{name: "malware_scanning", apiValue: "{}"},
	}
	require.Len(t, generatedBehaviorValueVariants, len(tests))
	valueValidator := DeprecatedSiblingDiscriminatorValue("files_behavior.value", "behavior", generatedBehaviorValueVariants, "legacy warning")

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			var apiValue any
			require.NoError(t, json.Unmarshal([]byte(test.apiValue), &apiValue))

			legacyNative := mustDynamicValue(t, apiValue)
			legacyJSON := types.DynamicValue(types.StringValue(test.apiValue))
			wrappedValue := map[string]interface{}{test.name: apiValue}
			wrapped := mustDynamicValue(t, wrappedValue)

			for name, value := range map[string]types.Dynamic{"legacy native": legacyNative, "legacy JSON": legacyJSON} {
				t.Run(name, func(t *testing.T) {
					response := &validator.DynamicResponse{}
					valueValidator.ValidateDynamic(context.Background(), validatorRequest(t, test.name, value), response)
					require.False(t, response.Diagnostics.HasError(), response.Diagnostics)
					assert.Equal(t, 1, response.Diagnostics.WarningsCount())

					result, diags := DynamicSiblingDiscriminatorToAPI(context.Background(), path.Root("value"), value, test.name, generatedBehaviorValueVariants)
					require.False(t, diags.HasError(), diags)
					assert.Equal(t, apiValue, result)
				})
			}

			response := &validator.DynamicResponse{}
			valueValidator.ValidateDynamic(context.Background(), validatorRequest(t, test.name, wrapped), response)
			require.False(t, response.Diagnostics.HasError(), response.Diagnostics)
			assert.Zero(t, response.Diagnostics.WarningsCount())
			result, diags := DynamicSiblingDiscriminatorToAPI(context.Background(), path.Root("value"), wrapped, test.name, generatedBehaviorValueVariants)
			require.False(t, diags.HasError(), diags)
			assert.Equal(t, apiValue, result)

			for name, prior := range map[string]types.Dynamic{"legacy native": legacyNative, "new wrapper": wrapped} {
				t.Run(name+" state", func(t *testing.T) {
					state, stateDiags := APIToDynamicSiblingDiscriminator(context.Background(), path.Root("value"), apiValue, test.name, generatedBehaviorValueVariants, prior)
					require.False(t, stateDiags.HasError(), stateDiags)
					assert.Equal(t, dynamicValueInterface(t, prior), dynamicValueInterface(t, state))
				})
			}

			jsonState, stateDiags := APIToDynamicSiblingDiscriminator(context.Background(), path.Root("value"), apiValue, test.name, generatedBehaviorValueVariants, legacyJSON)
			require.False(t, stateDiags.HasError(), stateDiags)
			var decodedState any
			require.NoError(t, json.Unmarshal([]byte(dynamicValueInterface(t, jsonState).(string)), &decodedState))
			assert.Equal(t, apiValue, decodedState)

			withoutPrior, stateDiags := APIToDynamicSiblingDiscriminator(context.Background(), path.Root("value"), apiValue, test.name, generatedBehaviorValueVariants, types.DynamicNull())
			require.False(t, stateDiags.HasError(), stateDiags)
			assert.Equal(t, apiValue, dynamicValueInterface(t, withoutPrior))
		})
	}
}
