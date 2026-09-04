package provider

import (
	"context"
	"fmt"

	files_sdk "github.com/Files-com/files-sdk-go/v3"
	behavior "github.com/Files-com/files-sdk-go/v3/behavior"
	"github.com/Files-com/terraform-provider-files/lib"
	"github.com/hashicorp/terraform-plugin-framework-validators/stringvalidator"
	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/datasource/schema"
	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/schema/validator"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

var (
	_ datasource.DataSource              = &behaviorDataSource{}
	_ datasource.DataSourceWithConfigure = &behaviorDataSource{}
)

func NewBehaviorDataSource() datasource.DataSource {
	return &behaviorDataSource{}
}

type behaviorDataSource struct {
	client *behavior.Client
}

type behaviorDataSourceModel struct {
	Id                          types.Int64   `tfsdk:"id"`
	Path                        types.String  `tfsdk:"path"`
	AttachmentUrl               types.String  `tfsdk:"attachment_url"`
	Behavior                    types.String  `tfsdk:"behavior"`
	Name                        types.String  `tfsdk:"name"`
	Description                 types.String  `tfsdk:"description"`
	Value                       types.Dynamic `tfsdk:"value"`
	ValueFormat                 types.String  `tfsdk:"value_format"`
	PublicHostingUrl            types.String  `tfsdk:"public_hosting_url"`
	DisableParentFolderBehavior types.Bool    `tfsdk:"disable_parent_folder_behavior"`
	Recursive                   types.Bool    `tfsdk:"recursive"`
	Inherited                   types.Bool    `tfsdk:"inherited"`
	Managed                     types.Bool    `tfsdk:"managed"`
	RootBehaviorSiteAdminOnly   types.Bool    `tfsdk:"root_behavior_site_admin_only"`
}

func (r *behaviorDataSource) Configure(_ context.Context, req datasource.ConfigureRequest, resp *datasource.ConfigureResponse) {
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

	r.client = &behavior.Client{Config: sdk_config}
}

func (r *behaviorDataSource) Metadata(_ context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_behavior"
}

func (r *behaviorDataSource) Schema(_ context.Context, _ datasource.SchemaRequest, resp *datasource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Description: "A Behavior is an API resource for what are also known as Folder Settings. Every behavior is associated with a folder.\n\n\n\nDepending on the behavior, it may also operate on child folders. It may be overridable at the child folder level or maybe can be added to at the child folder level. The exact options for each behavior type are explained in the table below.\n\n\n\nEach behavior type also has a recursion mode in the behavior type documentation. `always` means the behavior is always recursive, `never` means it is never recursive, and `sometimes` means callers may choose the value of the `recursive` field.\n\n\n\nAdditionally, some behaviors are visible to non-admins, and others are even settable by non-admins. All the details are below.\n\n\n\nEach behavior uses a different format for its settings value. The accepted fields and an example are shown with each behavior type. In the REST API, send these settings as JSON within the `value` field.\n\n\n\nNote: Append Timestamp behavior removed. Check [Override Upload Filename](#override-upload-filename-behaviors) behavior which have even more functionality to modify name on upload.",
		Attributes: map[string]schema.Attribute{
			"id": schema.Int64Attribute{
				Description: "Folder behavior ID",
				Required:    true,
			},
			"path": schema.StringAttribute{
				Description: "Folder path.  Note that Behavior paths cannot be updated once initially set.  You will need to remove and re-create the behavior on the new path. This must be slash-delimited, but it must neither start nor end with a slash. Maximum of 5000 characters.",
				Computed:    true,
			},
			"attachment_url": schema.StringAttribute{
				Description: "URL for attached file",
				Computed:    true,
			},
			"behavior": schema.StringAttribute{
				Description: "Behavior type.",
				Computed:    true,
			},
			"name": schema.StringAttribute{
				Description: "Name for this behavior.",
				Computed:    true,
			},
			"description": schema.StringAttribute{
				Description: "Description for this behavior.",
				Computed:    true,
			},
			"value": schema.DynamicAttribute{
				Description: "Settings for this behavior. Set `value_format = \"typed\"` to return the future typed shape under the selected behavior name.",
				Computed:    true,
			},
			"value_format": schema.StringAttribute{
				Description: "Set to `typed` to return the future files_behavior.value output shape before it becomes the default on March 1, 2027. Omit this attribute to keep the current output until then.",
				Optional:    true,
				Validators: []validator.String{
					stringvalidator.OneOf("typed"),
				},
			},
			"public_hosting_url": schema.StringAttribute{
				Description: "Public URL for this publicly hosted folder when the `Serve Publicly` behavior has a key configured.  When a Custom Domain with `public_hosting` destination is attached to this behavior, the URL uses that domain.  Otherwise it uses the site's `subdomain.hosted-by-files.com` host.",
				Computed:    true,
			},
			"disable_parent_folder_behavior": schema.BoolAttribute{
				Description: "If true, the parent folder's behavior will be disabled for this folder and its children.",
				Computed:    true,
			},
			"recursive": schema.BoolAttribute{
				Description: "Whether this behavior is recursive for this record. `always` behaviors are always `true`, `never` behaviors are always `false`, and `sometimes` behaviors may be either value.",
				Computed:    true,
			},
			"inherited": schema.BoolAttribute{
				Description: "If true, this behavior is inherited from a higher scope rather than owned by the requested workspace.",
				Computed:    true,
			},
			"managed": schema.BoolAttribute{
				Description: "If true, this behavior is controlled by a parent-site policy and cannot be modified locally.",
				Computed:    true,
			},
			"root_behavior_site_admin_only": schema.BoolAttribute{
				Description: "If true, this behavior may only be modified by a site admin because it is at the site root or disables a root behavior.",
				Computed:    true,
			},
		},
	}
}

func (r *behaviorDataSource) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
	var data behaviorDataSourceModel
	diags := req.Config.Get(ctx, &data)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	paramsBehaviorFind := files_sdk.BehaviorFindParams{}
	paramsBehaviorFind.Id = data.Id.ValueInt64()

	behavior, err := r.client.Find(paramsBehaviorFind, files_sdk.WithContext(ctx))
	if err != nil {
		resp.Diagnostics.AddError(
			"Error Reading Files Behavior",
			"Could not read behavior id "+fmt.Sprint(data.Id.ValueInt64())+": "+err.Error(),
		)
		return
	}

	diags = r.populateDataSourceModel(ctx, behavior, &data)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	diags = resp.State.Set(ctx, data)
	resp.Diagnostics.Append(diags...)
}

func (r *behaviorDataSource) populateDataSourceModel(ctx context.Context, behavior files_sdk.Behavior, state *behaviorDataSourceModel) (diags diag.Diagnostics) {
	var propDiags diag.Diagnostics

	state.Id = types.Int64Value(behavior.Id)
	state.Path = types.StringValue(behavior.Path)
	state.AttachmentUrl = types.StringValue(behavior.AttachmentUrl)
	state.Behavior = types.StringValue(behavior.Behavior)
	state.Name = types.StringValue(behavior.Name)
	state.Description = types.StringValue(behavior.Description)
	if state.ValueFormat.ValueString() == "typed" {
		valueValue, transformDiags := lib.WrapSiblingDiscriminator(ctx, path.Root("value"), behavior.Value, behavior.Behavior, []lib.JSONSchemaVariant{{Name: "webhook", Value: "webhook", RootType: "object", Writable: []string{"urls", "method", "triggers", "triggering_filenames", "exclude_filenames", "encoding", "headers", "body", "verification_token", "file_form_field", "file_as_body", "use_dedicated_ips"}, LegacyProperty: "urls", LegacyList: true}, {Name: "file_expiration", Value: "file_expiration", RootType: "object", Writable: []string{"days_to_retain", "delete_empty_folders"}, LegacyProperty: "days_to_retain"}, {Name: "auto_encrypt", Value: "auto_encrypt", RootType: "object", Writable: []string{"gpg_key_id", "gpg_key_ids", "algorithm", "signing_key_id", "suffix", "armor", "gpg_key_partner_id"}}, {Name: "lock_subfolders", Value: "lock_subfolders", RootType: "object", Writable: []string{"level"}}, {Name: "storage_region", Value: "storage_region", RootType: "string"}, {Name: "serve_publicly", Value: "serve_publicly", RootType: "object", Writable: []string{"key", "show_index", "force_download", "username", "password", "cors_enabled", "require_site_authentication"}, WriteOnly: []string{"password"}}, {Name: "create_user_folders", Value: "create_user_folders", RootType: "object", Writable: []string{"permission", "additional_permission", "existing_users", "group_id", "new_folder_name", "subfolders"}}, {Name: "inbox", Value: "inbox", RootType: "object", Writable: []string{"key", "dont_separate_submissions_by_folder", "dont_separate_submissions_by_folder_for_inbound_email", "dont_allow_folders_in_uploads", "require_inbox_recipient", "show_on_login_page", "clickwrap_id", "form_field_set_id", "title", "description", "help_text", "require_registration", "password", "path_template", "path_template_time_zone", "enable_inbound_email_address", "notify_senders_on_successful_uploads_via_email", "notify_senders_on_successful_uploads_via_web", "allow_whitelisting", "whitelist", "disable_web_upload", "capture_email_body_filename", "requested_upload_slots"}, WriteOnly: []string{"password", "enable_inbound_email_address"}}, {Name: "limit_file_extensions", Value: "limit_file_extensions", RootType: "object", Writable: []string{"extensions", "mode"}}, {Name: "limit_file_regex", Value: "limit_file_regex", RootType: "array"}, {Name: "amazon_sns", Value: "amazon_sns", RootType: "object", Writable: []string{"arns", "triggers", "aws_credentials", "body"}, WriteOnly: []string{"aws_credentials.secret_access_key"}}, {Name: "watermark", Value: "watermark", RootType: "object", Writable: []string{"gravity", "max_height_or_width", "transparency", "dynamic_text"}}, {Name: "remote_server_mount", Value: "remote_server_mount", RootType: "object", Writable: []string{"remote_server_id", "remote_path"}}, {Name: "slack_webhook", Value: "slack_webhook", RootType: "object", Writable: []string{"url", "username", "channel", "icon_emoji", "triggers"}}, {Name: "auto_decrypt", Value: "auto_decrypt", RootType: "object", Writable: []string{"gpg_key_id", "gpg_key_ids", "algorithm", "suffix", "ignore_mdc_error", "gpg_key_partner_id", "use_all_private_keys"}}, {Name: "override_upload_filename", Value: "override_upload_filename", RootType: "object", Writable: []string{"filename_override_pattern", "filename_replace_from", "filename_replace_to", "filename_regex_replace_from", "filename_regex_replace_to", "time_zone"}}, {Name: "permission_fence", Value: "permission_fence", RootType: "object", Writable: []string{"fenced_permissions"}}, {Name: "limit_filename_length", Value: "limit_filename_length", RootType: "object", Writable: []string{"max_length", "shorten"}}, {Name: "organize_files_into_subfolders", Value: "organize_files_into_subfolders", RootType: "object", Writable: []string{"subfolder_name_type", "regex", "strftime_format", "time_zone", "apply_behavior"}, WriteOnly: []string{"apply_behavior"}}, {Name: "teams_webhook", Value: "teams_webhook", RootType: "object", Writable: []string{"url", "triggers"}}, {Name: "google_pub_sub", Value: "google_pub_sub", RootType: "object", Writable: []string{"projects_topics", "triggers", "google_credentials", "body"}, WriteOnly: []string{"google_credentials.private_key"}}, {Name: "archive_overwritten_or_deleted_files", Value: "archive_overwritten_or_deleted_files", RootType: "object", Writable: []string{"archive_path"}}, {Name: "auto_recrypt", Value: "auto_recrypt", RootType: "object", Writable: []string{"decrypt_gpg_key_ids", "encrypt_gpg_key_ids", "decrypt_gpg_key_partner_id", "encrypt_gpg_key_partner_id", "ignore_mdc_error", "signing_key_id", "armor"}}, {Name: "metadata_category", Value: "metadata_category", RootType: "object", Writable: []string{"metadata_category_id"}}, {Name: "auto_unzip", Value: "auto_unzip", RootType: "object", Writable: []string{"destination_path", "path_time_zone"}}, {Name: "remote_server_metadata_index", Value: "remote_server_metadata_index", RootType: "object", Writable: []string{"interval_minutes"}}, {Name: "malware_scanning", Value: "malware_scanning", RootType: "object"}})
		diags.Append(transformDiags...)
		state.Value, propDiags = lib.ToDynamic(ctx, path.Root("value"), valueValue, nil)
	} else {
		state.Value, propDiags = lib.ToDynamic(ctx, path.Root("value"), behavior.Value, nil)
		diags.AddAttributeWarning(path.Root("value"), "Value Output Changes March 1, 2027: "+behavior.Behavior, fmt.Sprintf("On March 1, 2027, `value` output for behavior %q will switch to the typed shape. Set `value_format = \"typed\"` now to use the new output before the change. The new output nests the current value under `%s`. Read it from `.value.%s` instead of `.value`.", behavior.Behavior, behavior.Behavior, behavior.Behavior))
	}
	diags.Append(propDiags...)
	state.PublicHostingUrl = types.StringValue(behavior.PublicHostingUrl)
	state.DisableParentFolderBehavior = types.BoolPointerValue(behavior.DisableParentFolderBehavior)
	state.Recursive = types.BoolPointerValue(behavior.Recursive)
	state.Inherited = types.BoolPointerValue(behavior.Inherited)
	state.Managed = types.BoolPointerValue(behavior.Managed)
	state.RootBehaviorSiteAdminOnly = types.BoolPointerValue(behavior.RootBehaviorSiteAdminOnly)

	return
}
