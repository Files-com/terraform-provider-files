package provider

import (
	"context"
	"encoding/json"
	"fmt"

	files_sdk "github.com/Files-com/files-sdk-go/v3"
	bundle "github.com/Files-com/files-sdk-go/v3/bundle"
	"github.com/Files-com/terraform-provider-files/lib"
	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/datasource/schema"
	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

var (
	_ datasource.DataSource              = &bundleDataSource{}
	_ datasource.DataSourceWithConfigure = &bundleDataSource{}
)

func NewBundleDataSource() datasource.DataSource {
	return &bundleDataSource{}
}

type bundleDataSource struct {
	client *bundle.Client
}

type bundleDataSourceModel struct {
	Id                                           types.Int64   `tfsdk:"id"`
	Code                                         types.String  `tfsdk:"code"`
	ColorLeft                                    types.String  `tfsdk:"color_left"`
	ColorLink                                    types.String  `tfsdk:"color_link"`
	ColorText                                    types.String  `tfsdk:"color_text"`
	ColorTop                                     types.String  `tfsdk:"color_top"`
	ColorTopText                                 types.String  `tfsdk:"color_top_text"`
	Url                                          types.String  `tfsdk:"url"`
	Description                                  types.String  `tfsdk:"description"`
	ExpiresAt                                    types.String  `tfsdk:"expires_at"`
	PasswordProtected                            types.Bool    `tfsdk:"password_protected"`
	Permissions                                  types.String  `tfsdk:"permissions"`
	PreviewOnly                                  types.Bool    `tfsdk:"preview_only"`
	RequireRegistration                          types.Bool    `tfsdk:"require_registration"`
	RequireShareRecipient                        types.Bool    `tfsdk:"require_share_recipient"`
	RequireLogout                                types.Bool    `tfsdk:"require_logout"`
	ClickwrapBody                                types.String  `tfsdk:"clickwrap_body"`
	FormFieldSet                                 types.String  `tfsdk:"form_field_set"`
	SkipName                                     types.Bool    `tfsdk:"skip_name"`
	SkipEmail                                    types.Bool    `tfsdk:"skip_email"`
	StartAccessOnDate                            types.String  `tfsdk:"start_access_on_date"`
	SkipCompany                                  types.Bool    `tfsdk:"skip_company"`
	CreatedAt                                    types.String  `tfsdk:"created_at"`
	DontSeparateSubmissionsByFolder              types.Bool    `tfsdk:"dont_separate_submissions_by_folder"`
	MaxUses                                      types.Int64   `tfsdk:"max_uses"`
	Note                                         types.String  `tfsdk:"note"`
	PathTemplate                                 types.String  `tfsdk:"path_template"`
	PathTemplateTimeZone                         types.String  `tfsdk:"path_template_time_zone"`
	SendEmailReceiptToUploader                   types.Bool    `tfsdk:"send_email_receipt_to_uploader"`
	SnapshotId                                   types.Int64   `tfsdk:"snapshot_id"`
	UserId                                       types.Int64   `tfsdk:"user_id"`
	Username                                     types.String  `tfsdk:"username"`
	ClickwrapId                                  types.Int64   `tfsdk:"clickwrap_id"`
	InboxId                                      types.Int64   `tfsdk:"inbox_id"`
	WatermarkAttachment                          types.String  `tfsdk:"watermark_attachment"`
	WatermarkValue                               types.Dynamic `tfsdk:"watermark_value"`
	SendOneTimePasswordToRecipientAtRegistration types.Bool    `tfsdk:"send_one_time_password_to_recipient_at_registration"`
	HasInbox                                     types.Bool    `tfsdk:"has_inbox"`
	DontAllowFoldersInUploads                    types.Bool    `tfsdk:"dont_allow_folders_in_uploads"`
	Paths                                        types.List    `tfsdk:"paths"`
	Bundlepaths                                  types.Dynamic `tfsdk:"bundlepaths"`
}

func (r *bundleDataSource) Configure(_ context.Context, req datasource.ConfigureRequest, resp *datasource.ConfigureResponse) {
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

	r.client = &bundle.Client{Config: sdk_config}
}

func (r *bundleDataSource) Metadata(_ context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_bundle"
}

func (r *bundleDataSource) Schema(_ context.Context, _ datasource.SchemaRequest, resp *datasource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Description: "A Bundle is the API/SDK term for the feature called Share Links in the web interface.\n\nThe API provides the full set of actions related to Share Links, including sending them via E-Mail.\n\n\n\nPlease note that we very closely monitor the E-Mailing feature and any abuse will result in disabling of your site.",
		Attributes: map[string]schema.Attribute{
			"id": schema.Int64Attribute{
				Description: "Bundle ID",
				Required:    true,
			},
			"code": schema.StringAttribute{
				Description: "Bundle code.  This code forms the end part of the Public URL.",
				Computed:    true,
			},
			"color_left": schema.StringAttribute{
				Description: "Page link and button color",
				Computed:    true,
			},
			"color_link": schema.StringAttribute{
				Description: "Top bar link color",
				Computed:    true,
			},
			"color_text": schema.StringAttribute{
				Description: "Page link and button color",
				Computed:    true,
			},
			"color_top": schema.StringAttribute{
				Description: "Top bar background color",
				Computed:    true,
			},
			"color_top_text": schema.StringAttribute{
				Description: "Top bar text color",
				Computed:    true,
			},
			"url": schema.StringAttribute{
				Description: "Public URL of Share Link",
				Computed:    true,
			},
			"description": schema.StringAttribute{
				Description: "Public description",
				Computed:    true,
			},
			"expires_at": schema.StringAttribute{
				Description: "Bundle expiration date/time",
				Computed:    true,
			},
			"password_protected": schema.BoolAttribute{
				Description: "Is this bundle password protected?",
				Computed:    true,
			},
			"permissions": schema.StringAttribute{
				Description: "Permissions that apply to Folders in this Share Link.",
				Computed:    true,
			},
			"preview_only": schema.BoolAttribute{
				Computed: true,
			},
			"require_registration": schema.BoolAttribute{
				Description: "Show a registration page that captures the downloader's name and email address?",
				Computed:    true,
			},
			"require_share_recipient": schema.BoolAttribute{
				Description: "Only allow access to recipients who have explicitly received the share via an email sent through the Files.com UI?",
				Computed:    true,
			},
			"require_logout": schema.BoolAttribute{
				Description: "If true, we will hide the 'Remember Me' box on the Bundle registration page, requiring that the user logout and log back in every time they visit the page.",
				Computed:    true,
			},
			"clickwrap_body": schema.StringAttribute{
				Description: "Legal text that must be agreed to prior to accessing Bundle.",
				Computed:    true,
			},
			"form_field_set": schema.StringAttribute{
				Description: "Custom Form to use",
				Computed:    true,
			},
			"skip_name": schema.BoolAttribute{
				Description: "BundleRegistrations can be saved without providing name?",
				Computed:    true,
			},
			"skip_email": schema.BoolAttribute{
				Description: "BundleRegistrations can be saved without providing email?",
				Computed:    true,
			},
			"start_access_on_date": schema.StringAttribute{
				Description: "Date when share will start to be accessible. If `nil` access granted right after create.",
				Computed:    true,
			},
			"skip_company": schema.BoolAttribute{
				Description: "BundleRegistrations can be saved without providing company?",
				Computed:    true,
			},
			"created_at": schema.StringAttribute{
				Description: "Bundle created at date/time",
				Computed:    true,
			},
			"dont_separate_submissions_by_folder": schema.BoolAttribute{
				Description: "Do not create subfolders for files uploaded to this share. Note: there are subtle security pitfalls with allowing anonymous uploads from multiple users to live in the same folder. We strongly discourage use of this option unless absolutely required.",
				Computed:    true,
			},
			"max_uses": schema.Int64Attribute{
				Description: "Maximum number of times bundle can be accessed",
				Computed:    true,
			},
			"note": schema.StringAttribute{
				Description: "Bundle internal note",
				Computed:    true,
			},
			"path_template": schema.StringAttribute{
				Description: "Template for creating submission subfolders. Can use the uploader's name, email address, ip, company, `strftime` directives, and any custom form data.",
				Computed:    true,
			},
			"path_template_time_zone": schema.StringAttribute{
				Description: "Timezone to use when rendering timestamps in path templates.",
				Computed:    true,
			},
			"send_email_receipt_to_uploader": schema.BoolAttribute{
				Description: "Send delivery receipt to the uploader. Note: For writable share only",
				Computed:    true,
			},
			"snapshot_id": schema.Int64Attribute{
				Description: "ID of the snapshot containing this bundle's contents.",
				Computed:    true,
			},
			"user_id": schema.Int64Attribute{
				Description: "Bundle creator user ID",
				Computed:    true,
			},
			"username": schema.StringAttribute{
				Description: "Bundle creator username",
				Computed:    true,
			},
			"clickwrap_id": schema.Int64Attribute{
				Description: "ID of the clickwrap to use with this bundle.",
				Computed:    true,
			},
			"inbox_id": schema.Int64Attribute{
				Description: "ID of the associated inbox, if available.",
				Computed:    true,
			},
			"watermark_attachment": schema.StringAttribute{
				Description: "Preview watermark image applied to all bundle items.",
				Computed:    true,
			},
			"watermark_value": schema.DynamicAttribute{
				Description: "Preview watermark settings applied to all bundle items. Uses the same keys as Behavior.value",
				Computed:    true,
			},
			"send_one_time_password_to_recipient_at_registration": schema.BoolAttribute{
				Description: "If true, require_share_recipient bundles will send a one-time password to the recipient when they register. Cannot be enabled if the bundle has a password set.",
				Computed:    true,
			},
			"has_inbox": schema.BoolAttribute{
				Description: "Does this bundle have an associated inbox?",
				Computed:    true,
			},
			"dont_allow_folders_in_uploads": schema.BoolAttribute{
				Description: "Should folder uploads be prevented?",
				Computed:    true,
			},
			"paths": schema.ListAttribute{
				Description: "A list of paths in this bundle.  For performance reasons, this is not provided when listing bundles.",
				Computed:    true,
				ElementType: types.StringType,
			},
			"bundlepaths": schema.DynamicAttribute{
				Description: "A list of bundlepaths in this bundle.  For performance reasons, this is not provided when listing bundles.",
				Computed:    true,
			},
		},
	}
}

func (r *bundleDataSource) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
	var data bundleDataSourceModel
	diags := req.Config.Get(ctx, &data)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	paramsBundleFind := files_sdk.BundleFindParams{}
	paramsBundleFind.Id = data.Id.ValueInt64()

	bundle, err := r.client.Find(paramsBundleFind, files_sdk.WithContext(ctx))
	if err != nil {
		resp.Diagnostics.AddError(
			"Error Reading Files Bundle",
			"Could not read bundle id "+fmt.Sprint(data.Id.ValueInt64())+": "+err.Error(),
		)
		return
	}

	diags = r.populateDataSourceModel(ctx, bundle, &data)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	diags = resp.State.Set(ctx, data)
	resp.Diagnostics.Append(diags...)
}

func (r *bundleDataSource) populateDataSourceModel(ctx context.Context, bundle files_sdk.Bundle, state *bundleDataSourceModel) (diags diag.Diagnostics) {
	var propDiags diag.Diagnostics

	state.Code = types.StringValue(bundle.Code)
	state.ColorLeft = types.StringValue(bundle.ColorLeft)
	state.ColorLink = types.StringValue(bundle.ColorLink)
	state.ColorText = types.StringValue(bundle.ColorText)
	state.ColorTop = types.StringValue(bundle.ColorTop)
	state.ColorTopText = types.StringValue(bundle.ColorTopText)
	state.Url = types.StringValue(bundle.Url)
	state.Description = types.StringValue(bundle.Description)
	if err := lib.TimeToStringType(ctx, path.Root("expires_at"), bundle.ExpiresAt, &state.ExpiresAt); err != nil {
		diags.AddError(
			"Error Creating Files Bundle",
			"Could not convert state expires_at to string: "+err.Error(),
		)
	}
	state.PasswordProtected = types.BoolPointerValue(bundle.PasswordProtected)
	state.Permissions = types.StringValue(bundle.Permissions)
	state.PreviewOnly = types.BoolPointerValue(bundle.PreviewOnly)
	state.RequireRegistration = types.BoolPointerValue(bundle.RequireRegistration)
	state.RequireShareRecipient = types.BoolPointerValue(bundle.RequireShareRecipient)
	state.RequireLogout = types.BoolPointerValue(bundle.RequireLogout)
	state.ClickwrapBody = types.StringValue(bundle.ClickwrapBody)
	respFormFieldSet, err := json.Marshal(bundle.FormFieldSet)
	if err != nil {
		diags.AddError(
			"Error Creating Files Bundle",
			"Could not marshal form_field_set to JSON: "+err.Error(),
		)
	}
	state.FormFieldSet = types.StringValue(string(respFormFieldSet))
	state.SkipName = types.BoolPointerValue(bundle.SkipName)
	state.SkipEmail = types.BoolPointerValue(bundle.SkipEmail)
	if err := lib.TimeToStringType(ctx, path.Root("start_access_on_date"), bundle.StartAccessOnDate, &state.StartAccessOnDate); err != nil {
		diags.AddError(
			"Error Creating Files Bundle",
			"Could not convert state start_access_on_date to string: "+err.Error(),
		)
	}
	state.SkipCompany = types.BoolPointerValue(bundle.SkipCompany)
	state.Id = types.Int64Value(bundle.Id)
	if err := lib.TimeToStringType(ctx, path.Root("created_at"), bundle.CreatedAt, &state.CreatedAt); err != nil {
		diags.AddError(
			"Error Creating Files Bundle",
			"Could not convert state created_at to string: "+err.Error(),
		)
	}
	state.DontSeparateSubmissionsByFolder = types.BoolPointerValue(bundle.DontSeparateSubmissionsByFolder)
	state.MaxUses = types.Int64Value(bundle.MaxUses)
	state.Note = types.StringValue(bundle.Note)
	state.PathTemplate = types.StringValue(bundle.PathTemplate)
	state.PathTemplateTimeZone = types.StringValue(bundle.PathTemplateTimeZone)
	state.SendEmailReceiptToUploader = types.BoolPointerValue(bundle.SendEmailReceiptToUploader)
	state.SnapshotId = types.Int64Value(bundle.SnapshotId)
	state.UserId = types.Int64Value(bundle.UserId)
	state.Username = types.StringValue(bundle.Username)
	state.ClickwrapId = types.Int64Value(bundle.ClickwrapId)
	state.InboxId = types.Int64Value(bundle.InboxId)
	respWatermarkAttachment, err := json.Marshal(bundle.WatermarkAttachment)
	if err != nil {
		diags.AddError(
			"Error Creating Files Bundle",
			"Could not marshal watermark_attachment to JSON: "+err.Error(),
		)
	}
	state.WatermarkAttachment = types.StringValue(string(respWatermarkAttachment))
	state.WatermarkValue, propDiags = lib.ToDynamic(ctx, path.Root("watermark_value"), bundle.WatermarkValue, state.WatermarkValue.UnderlyingValue())
	diags.Append(propDiags...)
	state.SendOneTimePasswordToRecipientAtRegistration = types.BoolPointerValue(bundle.SendOneTimePasswordToRecipientAtRegistration)
	state.HasInbox = types.BoolPointerValue(bundle.HasInbox)
	state.DontAllowFoldersInUploads = types.BoolPointerValue(bundle.DontAllowFoldersInUploads)
	state.Paths, propDiags = types.ListValueFrom(ctx, types.StringType, bundle.Paths)
	diags.Append(propDiags...)
	state.Bundlepaths, propDiags = lib.ToDynamic(ctx, path.Root("bundlepaths"), bundle.Bundlepaths, state.Bundlepaths.UnderlyingValue())
	diags.Append(propDiags...)

	return
}
