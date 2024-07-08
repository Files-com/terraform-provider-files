package provider

import (
	"context"
	"encoding/json"
	"fmt"
	"strconv"
	"strings"
	"time"

	files_sdk "github.com/Files-com/files-sdk-go/v3"
	bundle "github.com/Files-com/files-sdk-go/v3/bundle"
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
	_ resource.Resource                = &bundleResource{}
	_ resource.ResourceWithConfigure   = &bundleResource{}
	_ resource.ResourceWithImportState = &bundleResource{}
)

func NewBundleResource() resource.Resource {
	return &bundleResource{}
}

type bundleResource struct {
	client *bundle.Client
}

type bundleResourceModel struct {
	Code                            types.String  `tfsdk:"code"`
	ColorLeft                       types.String  `tfsdk:"color_left"`
	ColorLink                       types.String  `tfsdk:"color_link"`
	ColorText                       types.String  `tfsdk:"color_text"`
	ColorTop                        types.String  `tfsdk:"color_top"`
	ColorTopText                    types.String  `tfsdk:"color_top_text"`
	Url                             types.String  `tfsdk:"url"`
	Description                     types.String  `tfsdk:"description"`
	ExpiresAt                       types.String  `tfsdk:"expires_at"`
	PasswordProtected               types.Bool    `tfsdk:"password_protected"`
	Permissions                     types.String  `tfsdk:"permissions"`
	PreviewOnly                     types.Bool    `tfsdk:"preview_only"`
	RequireRegistration             types.Bool    `tfsdk:"require_registration"`
	RequireShareRecipient           types.Bool    `tfsdk:"require_share_recipient"`
	RequireLogout                   types.Bool    `tfsdk:"require_logout"`
	ClickwrapBody                   types.String  `tfsdk:"clickwrap_body"`
	FormFieldSet                    types.String  `tfsdk:"form_field_set"`
	SkipName                        types.Bool    `tfsdk:"skip_name"`
	SkipEmail                       types.Bool    `tfsdk:"skip_email"`
	StartAccessOnDate               types.String  `tfsdk:"start_access_on_date"`
	SkipCompany                     types.Bool    `tfsdk:"skip_company"`
	Id                              types.Int64   `tfsdk:"id"`
	CreatedAt                       types.String  `tfsdk:"created_at"`
	DontSeparateSubmissionsByFolder types.Bool    `tfsdk:"dont_separate_submissions_by_folder"`
	MaxUses                         types.Int64   `tfsdk:"max_uses"`
	Note                            types.String  `tfsdk:"note"`
	PathTemplate                    types.String  `tfsdk:"path_template"`
	PathTemplateTimeZone            types.String  `tfsdk:"path_template_time_zone"`
	SendEmailReceiptToUploader      types.Bool    `tfsdk:"send_email_receipt_to_uploader"`
	SnapshotId                      types.Int64   `tfsdk:"snapshot_id"`
	UserId                          types.Int64   `tfsdk:"user_id"`
	Username                        types.String  `tfsdk:"username"`
	ClickwrapId                     types.Int64   `tfsdk:"clickwrap_id"`
	InboxId                         types.Int64   `tfsdk:"inbox_id"`
	WatermarkAttachment             types.String  `tfsdk:"watermark_attachment"`
	WatermarkValue                  types.Dynamic `tfsdk:"watermark_value"`
	HasInbox                        types.Bool    `tfsdk:"has_inbox"`
	Paths                           types.List    `tfsdk:"paths"`
	Bundlepaths                     types.Dynamic `tfsdk:"bundlepaths"`
	Password                        types.String  `tfsdk:"password"`
	FormFieldSetId                  types.Int64   `tfsdk:"form_field_set_id"`
	CreateSnapshot                  types.Bool    `tfsdk:"create_snapshot"`
	FinalizeSnapshot                types.Bool    `tfsdk:"finalize_snapshot"`
}

func (r *bundleResource) Configure(_ context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {
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

func (r *bundleResource) Metadata(_ context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_bundle"
}

func (r *bundleResource) Schema(_ context.Context, _ resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Description: "Bundles are the API/SDK term for the feature called Share Links in the web interface.\n\nThe API provides the full set of actions related to Share Links, including sending them via E-Mail.\n\n\n\nPlease note that we very closely monitor the E-Mailing feature and any abuse will result in disabling of your site.",
		Attributes: map[string]schema.Attribute{
			"code": schema.StringAttribute{
				Description: "Bundle code.  This code forms the end part of the Public URL.",
				Computed:    true,
				Optional:    true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
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
				Optional:    true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
			},
			"expires_at": schema.StringAttribute{
				Description: "Bundle expiration date/time",
				Computed:    true,
				Optional:    true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
			},
			"password_protected": schema.BoolAttribute{
				Description: "Is this bundle password protected?",
				Computed:    true,
			},
			"permissions": schema.StringAttribute{
				Description: "Permissions that apply to Folders in this Share Link.",
				Computed:    true,
				Optional:    true,
				Validators: []validator.String{
					stringvalidator.OneOf("read", "write", "read_write", "full", "none", "preview_only"),
				},
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
			},
			"preview_only": schema.BoolAttribute{
				Computed: true,
			},
			"require_registration": schema.BoolAttribute{
				Description: "Show a registration page that captures the downloader's name and email address?",
				Computed:    true,
				Optional:    true,
				PlanModifiers: []planmodifier.Bool{
					boolplanmodifier.UseStateForUnknown(),
				},
			},
			"require_share_recipient": schema.BoolAttribute{
				Description: "Only allow access to recipients who have explicitly received the share via an email sent through the Files.com UI?",
				Computed:    true,
				Optional:    true,
				PlanModifiers: []planmodifier.Bool{
					boolplanmodifier.UseStateForUnknown(),
				},
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
				Optional:    true,
				PlanModifiers: []planmodifier.Bool{
					boolplanmodifier.UseStateForUnknown(),
				},
			},
			"skip_email": schema.BoolAttribute{
				Description: "BundleRegistrations can be saved without providing email?",
				Computed:    true,
				Optional:    true,
				PlanModifiers: []planmodifier.Bool{
					boolplanmodifier.UseStateForUnknown(),
				},
			},
			"start_access_on_date": schema.StringAttribute{
				Description: "Date when share will start to be accessible. If `nil` access granted right after create.",
				Computed:    true,
				Optional:    true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
			},
			"skip_company": schema.BoolAttribute{
				Description: "BundleRegistrations can be saved without providing company?",
				Computed:    true,
				Optional:    true,
				PlanModifiers: []planmodifier.Bool{
					boolplanmodifier.UseStateForUnknown(),
				},
			},
			"id": schema.Int64Attribute{
				Description: "Bundle ID",
				Computed:    true,
				PlanModifiers: []planmodifier.Int64{
					int64planmodifier.UseStateForUnknown(),
				},
			},
			"created_at": schema.StringAttribute{
				Description: "Bundle created at date/time",
				Computed:    true,
			},
			"dont_separate_submissions_by_folder": schema.BoolAttribute{
				Description: "Do not create subfolders for files uploaded to this share. Note: there are subtle security pitfalls with allowing anonymous uploads from multiple users to live in the same folder. We strongly discourage use of this option unless absolutely required.",
				Computed:    true,
				Optional:    true,
				PlanModifiers: []planmodifier.Bool{
					boolplanmodifier.UseStateForUnknown(),
				},
			},
			"max_uses": schema.Int64Attribute{
				Description: "Maximum number of times bundle can be accessed",
				Computed:    true,
				Optional:    true,
				PlanModifiers: []planmodifier.Int64{
					int64planmodifier.UseStateForUnknown(),
				},
			},
			"note": schema.StringAttribute{
				Description: "Bundle internal note",
				Computed:    true,
				Optional:    true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
			},
			"path_template": schema.StringAttribute{
				Description: "Template for creating submission subfolders. Can use the uploader's name, email address, ip, company, `strftime` directives, and any custom form data.",
				Computed:    true,
				Optional:    true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
			},
			"path_template_time_zone": schema.StringAttribute{
				Description: "Timezone to use when rendering timestamps in path templates.",
				Computed:    true,
				Optional:    true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
			},
			"send_email_receipt_to_uploader": schema.BoolAttribute{
				Description: "Send delivery receipt to the uploader. Note: For writable share only",
				Computed:    true,
				Optional:    true,
				PlanModifiers: []planmodifier.Bool{
					boolplanmodifier.UseStateForUnknown(),
				},
			},
			"snapshot_id": schema.Int64Attribute{
				Description: "ID of the snapshot containing this bundle's contents.",
				Computed:    true,
				Optional:    true,
				PlanModifiers: []planmodifier.Int64{
					int64planmodifier.UseStateForUnknown(),
					int64planmodifier.RequiresReplace(),
				},
			},
			"user_id": schema.Int64Attribute{
				Description: "Bundle creator user ID",
				Computed:    true,
				Optional:    true,
				PlanModifiers: []planmodifier.Int64{
					int64planmodifier.UseStateForUnknown(),
					int64planmodifier.RequiresReplace(),
				},
			},
			"username": schema.StringAttribute{
				Description: "Bundle creator username",
				Computed:    true,
			},
			"clickwrap_id": schema.Int64Attribute{
				Description: "ID of the clickwrap to use with this bundle.",
				Computed:    true,
				Optional:    true,
				PlanModifiers: []planmodifier.Int64{
					int64planmodifier.UseStateForUnknown(),
				},
			},
			"inbox_id": schema.Int64Attribute{
				Description: "ID of the associated inbox, if available.",
				Computed:    true,
				Optional:    true,
				PlanModifiers: []planmodifier.Int64{
					int64planmodifier.UseStateForUnknown(),
				},
			},
			"watermark_attachment": schema.StringAttribute{
				Description: "Preview watermark image applied to all bundle items.",
				Computed:    true,
			},
			"watermark_value": schema.DynamicAttribute{
				Description: "Preview watermark settings applied to all bundle items. Uses the same keys as Behavior.value",
				Computed:    true,
			},
			"has_inbox": schema.BoolAttribute{
				Description: "Does this bundle have an associated inbox?",
				Computed:    true,
			},
			"paths": schema.ListAttribute{
				Description: "A list of paths in this bundle.  For performance reasons, this is not provided when listing bundles.",
				Required:    true,
				ElementType: types.StringType,
			},
			"bundlepaths": schema.DynamicAttribute{
				Description: "A list of bundlepaths in this bundle.  For performance reasons, this is not provided when listing bundles.",
				Computed:    true,
			},
			"password": schema.StringAttribute{
				Description: "Password for this bundle.",
				Optional:    true,
			},
			"form_field_set_id": schema.Int64Attribute{
				Description: "Id of Form Field Set to use with this bundle",
				Optional:    true,
			},
			"create_snapshot": schema.BoolAttribute{
				Description: "If true, create a snapshot of this bundle's contents.",
				Optional:    true,
			},
			"finalize_snapshot": schema.BoolAttribute{
				Description: "If true, finalize the snapshot of this bundle's contents. Note that `create_snapshot` must also be true.",
				Optional:    true,
			},
		},
	}
}

func (r *bundleResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	var plan bundleResourceModel
	diags := req.Plan.Get(ctx, &plan)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	paramsBundleCreate := files_sdk.BundleCreateParams{}
	paramsBundleCreate.UserId = plan.UserId.ValueInt64()
	if !plan.Paths.IsNull() && !plan.Paths.IsUnknown() {
		diags = plan.Paths.ElementsAs(ctx, &paramsBundleCreate.Paths, false)
		resp.Diagnostics.Append(diags...)
	}
	paramsBundleCreate.Password = plan.Password.ValueString()
	paramsBundleCreate.FormFieldSetId = plan.FormFieldSetId.ValueInt64()
	paramsBundleCreate.CreateSnapshot = plan.CreateSnapshot.ValueBoolPointer()
	paramsBundleCreate.DontSeparateSubmissionsByFolder = plan.DontSeparateSubmissionsByFolder.ValueBoolPointer()
	if !plan.ExpiresAt.IsNull() && plan.ExpiresAt.ValueString() != "" {
		createExpiresAt, err := time.Parse(time.RFC3339, plan.ExpiresAt.ValueString())
		if err != nil {
			resp.Diagnostics.AddAttributeError(
				path.Root("expires_at"),
				"Error Parsing expires_at Time",
				"Could not parse expires_at time: "+err.Error(),
			)
		} else {
			paramsBundleCreate.ExpiresAt = &createExpiresAt
		}
	}
	paramsBundleCreate.FinalizeSnapshot = plan.FinalizeSnapshot.ValueBoolPointer()
	paramsBundleCreate.MaxUses = plan.MaxUses.ValueInt64()
	paramsBundleCreate.Description = plan.Description.ValueString()
	paramsBundleCreate.Note = plan.Note.ValueString()
	paramsBundleCreate.Code = plan.Code.ValueString()
	paramsBundleCreate.PathTemplate = plan.PathTemplate.ValueString()
	paramsBundleCreate.PathTemplateTimeZone = plan.PathTemplateTimeZone.ValueString()
	paramsBundleCreate.Permissions = paramsBundleCreate.Permissions.Enum()[plan.Permissions.ValueString()]
	paramsBundleCreate.RequireRegistration = plan.RequireRegistration.ValueBoolPointer()
	paramsBundleCreate.ClickwrapId = plan.ClickwrapId.ValueInt64()
	paramsBundleCreate.InboxId = plan.InboxId.ValueInt64()
	paramsBundleCreate.RequireShareRecipient = plan.RequireShareRecipient.ValueBoolPointer()
	paramsBundleCreate.SendEmailReceiptToUploader = plan.SendEmailReceiptToUploader.ValueBoolPointer()
	paramsBundleCreate.SkipEmail = plan.SkipEmail.ValueBoolPointer()
	paramsBundleCreate.SkipName = plan.SkipName.ValueBoolPointer()
	paramsBundleCreate.SkipCompany = plan.SkipCompany.ValueBoolPointer()
	if !plan.StartAccessOnDate.IsNull() && plan.StartAccessOnDate.ValueString() != "" {
		createStartAccessOnDate, err := time.Parse(time.RFC3339, plan.StartAccessOnDate.ValueString())
		if err != nil {
			resp.Diagnostics.AddAttributeError(
				path.Root("start_access_on_date"),
				"Error Parsing start_access_on_date Time",
				"Could not parse start_access_on_date time: "+err.Error(),
			)
		} else {
			paramsBundleCreate.StartAccessOnDate = &createStartAccessOnDate
		}
	}
	paramsBundleCreate.SnapshotId = plan.SnapshotId.ValueInt64()

	if resp.Diagnostics.HasError() {
		return
	}

	bundle, err := r.client.Create(paramsBundleCreate, files_sdk.WithContext(ctx))
	if err != nil {
		resp.Diagnostics.AddError(
			"Error Creating Files Bundle",
			"Could not create bundle, unexpected error: "+err.Error(),
		)
		return
	}

	diags = r.populateResourceModel(ctx, bundle, &plan)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	diags = resp.State.Set(ctx, plan)
	resp.Diagnostics.Append(diags...)
}

func (r *bundleResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	var state bundleResourceModel
	diags := req.State.Get(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	paramsBundleFind := files_sdk.BundleFindParams{}
	paramsBundleFind.Id = state.Id.ValueInt64()

	bundle, err := r.client.Find(paramsBundleFind, files_sdk.WithContext(ctx))
	if err != nil {
		resp.Diagnostics.AddError(
			"Error Reading Files Bundle",
			"Could not read bundle id "+fmt.Sprint(state.Id.ValueInt64())+": "+err.Error(),
		)
		return
	}

	diags = r.populateResourceModel(ctx, bundle, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	diags = resp.State.Set(ctx, state)
	resp.Diagnostics.Append(diags...)
}

func (r *bundleResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	var plan bundleResourceModel
	diags := req.Plan.Get(ctx, &plan)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	paramsBundleUpdate := files_sdk.BundleUpdateParams{}
	paramsBundleUpdate.Id = plan.Id.ValueInt64()
	if !plan.Paths.IsNull() && !plan.Paths.IsUnknown() {
		diags = plan.Paths.ElementsAs(ctx, &paramsBundleUpdate.Paths, false)
		resp.Diagnostics.Append(diags...)
	}
	paramsBundleUpdate.Password = plan.Password.ValueString()
	paramsBundleUpdate.FormFieldSetId = plan.FormFieldSetId.ValueInt64()
	paramsBundleUpdate.ClickwrapId = plan.ClickwrapId.ValueInt64()
	paramsBundleUpdate.Code = plan.Code.ValueString()
	paramsBundleUpdate.CreateSnapshot = plan.CreateSnapshot.ValueBoolPointer()
	paramsBundleUpdate.Description = plan.Description.ValueString()
	paramsBundleUpdate.DontSeparateSubmissionsByFolder = plan.DontSeparateSubmissionsByFolder.ValueBoolPointer()
	if !plan.ExpiresAt.IsNull() && plan.ExpiresAt.ValueString() != "" {
		updateExpiresAt, err := time.Parse(time.RFC3339, plan.ExpiresAt.ValueString())
		if err != nil {
			resp.Diagnostics.AddAttributeError(
				path.Root("expires_at"),
				"Error Parsing expires_at Time",
				"Could not parse expires_at time: "+err.Error(),
			)
		} else {
			paramsBundleUpdate.ExpiresAt = &updateExpiresAt
		}
	}
	paramsBundleUpdate.FinalizeSnapshot = plan.FinalizeSnapshot.ValueBoolPointer()
	paramsBundleUpdate.InboxId = plan.InboxId.ValueInt64()
	paramsBundleUpdate.MaxUses = plan.MaxUses.ValueInt64()
	paramsBundleUpdate.Note = plan.Note.ValueString()
	paramsBundleUpdate.PathTemplate = plan.PathTemplate.ValueString()
	paramsBundleUpdate.PathTemplateTimeZone = plan.PathTemplateTimeZone.ValueString()
	paramsBundleUpdate.Permissions = paramsBundleUpdate.Permissions.Enum()[plan.Permissions.ValueString()]
	paramsBundleUpdate.RequireRegistration = plan.RequireRegistration.ValueBoolPointer()
	paramsBundleUpdate.RequireShareRecipient = plan.RequireShareRecipient.ValueBoolPointer()
	paramsBundleUpdate.SendEmailReceiptToUploader = plan.SendEmailReceiptToUploader.ValueBoolPointer()
	paramsBundleUpdate.SkipCompany = plan.SkipCompany.ValueBoolPointer()
	if !plan.StartAccessOnDate.IsNull() && plan.StartAccessOnDate.ValueString() != "" {
		updateStartAccessOnDate, err := time.Parse(time.RFC3339, plan.StartAccessOnDate.ValueString())
		if err != nil {
			resp.Diagnostics.AddAttributeError(
				path.Root("start_access_on_date"),
				"Error Parsing start_access_on_date Time",
				"Could not parse start_access_on_date time: "+err.Error(),
			)
		} else {
			paramsBundleUpdate.StartAccessOnDate = &updateStartAccessOnDate
		}
	}
	paramsBundleUpdate.SkipEmail = plan.SkipEmail.ValueBoolPointer()
	paramsBundleUpdate.SkipName = plan.SkipName.ValueBoolPointer()

	if resp.Diagnostics.HasError() {
		return
	}

	bundle, err := r.client.Update(paramsBundleUpdate, files_sdk.WithContext(ctx))
	if err != nil {
		resp.Diagnostics.AddError(
			"Error Updating Files Bundle",
			"Could not update bundle, unexpected error: "+err.Error(),
		)
		return
	}

	diags = r.populateResourceModel(ctx, bundle, &plan)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	diags = resp.State.Set(ctx, plan)
	resp.Diagnostics.Append(diags...)
}

func (r *bundleResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	var state bundleResourceModel
	diags := req.State.Get(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	paramsBundleDelete := files_sdk.BundleDeleteParams{}
	paramsBundleDelete.Id = state.Id.ValueInt64()

	err := r.client.Delete(paramsBundleDelete, files_sdk.WithContext(ctx))
	if err != nil {
		resp.Diagnostics.AddError(
			"Error Deleting Files Bundle",
			"Could not delete bundle id "+fmt.Sprint(state.Id.ValueInt64())+": "+err.Error(),
		)
	}
}

func (r *bundleResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
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

func (r *bundleResource) populateResourceModel(ctx context.Context, bundle files_sdk.Bundle, state *bundleResourceModel) (diags diag.Diagnostics) {
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
	state.HasInbox = types.BoolPointerValue(bundle.HasInbox)
	state.Paths, propDiags = types.ListValueFrom(ctx, types.StringType, bundle.Paths)
	diags.Append(propDiags...)
	state.Bundlepaths, propDiags = lib.ToDynamic(ctx, path.Root("bundlepaths"), bundle.Bundlepaths, state.Bundlepaths.UnderlyingValue())
	diags.Append(propDiags...)

	return
}
