package provider

import (
	"context"
	"fmt"
	"strconv"
	"strings"

	files_sdk "github.com/Files-com/files-sdk-go/v3"
	notification "github.com/Files-com/files-sdk-go/v3/notification"
	"github.com/hashicorp/terraform-plugin-framework-validators/stringvalidator"
	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/boolplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/int64planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/listplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/schema/validator"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

var (
	_ resource.Resource                = &notificationResource{}
	_ resource.ResourceWithConfigure   = &notificationResource{}
	_ resource.ResourceWithImportState = &notificationResource{}
)

func NewNotificationResource() resource.Resource {
	return &notificationResource{}
}

type notificationResource struct {
	client *notification.Client
}

type notificationResourceModel struct {
	Path                     types.String `tfsdk:"path"`
	GroupId                  types.Int64  `tfsdk:"group_id"`
	TriggeringGroupIds       types.List   `tfsdk:"triggering_group_ids"`
	TriggeringUserIds        types.List   `tfsdk:"triggering_user_ids"`
	TriggerByShareRecipients types.Bool   `tfsdk:"trigger_by_share_recipients"`
	NotifyUserActions        types.Bool   `tfsdk:"notify_user_actions"`
	NotifyOnCopy             types.Bool   `tfsdk:"notify_on_copy"`
	NotifyOnDelete           types.Bool   `tfsdk:"notify_on_delete"`
	NotifyOnDownload         types.Bool   `tfsdk:"notify_on_download"`
	NotifyOnMove             types.Bool   `tfsdk:"notify_on_move"`
	NotifyOnUpload           types.Bool   `tfsdk:"notify_on_upload"`
	Recursive                types.Bool   `tfsdk:"recursive"`
	SendInterval             types.String `tfsdk:"send_interval"`
	Message                  types.String `tfsdk:"message"`
	TriggeringFilenames      types.List   `tfsdk:"triggering_filenames"`
	UserId                   types.Int64  `tfsdk:"user_id"`
	Username                 types.String `tfsdk:"username"`
	Id                       types.Int64  `tfsdk:"id"`
	GroupName                types.String `tfsdk:"group_name"`
	Unsubscribed             types.Bool   `tfsdk:"unsubscribed"`
	UnsubscribedReason       types.String `tfsdk:"unsubscribed_reason"`
	SuppressedEmail          types.Bool   `tfsdk:"suppressed_email"`
}

func (r *notificationResource) Configure(_ context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {
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

	r.client = &notification.Client{Config: sdk_config}
}

func (r *notificationResource) Metadata(_ context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_notification"
}

func (r *notificationResource) Schema(_ context.Context, _ resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Description: "A Notification is our feature that sends E-Mails when specific actions occur in the folder.\n\n\n\nEmails are sent in batches, with email frequency options of every 5 minutes, every 15 minutes, hourly, or daily. They will include a list of the matching actions within the configured notification period, limited to the first 100.",
		Attributes: map[string]schema.Attribute{
			"path": schema.StringAttribute{
				Description: "Folder path to notify on. This must be slash-delimited, but it must neither start nor end with a slash. Maximum of 5000 characters.",
				Computed:    true,
				Optional:    true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
					stringplanmodifier.RequiresReplace(),
				},
			},
			"group_id": schema.Int64Attribute{
				Description: "ID of Group to receive notifications",
				Computed:    true,
				Optional:    true,
				PlanModifiers: []planmodifier.Int64{
					int64planmodifier.UseStateForUnknown(),
					int64planmodifier.RequiresReplace(),
				},
			},
			"triggering_group_ids": schema.ListAttribute{
				Description: "If set, will only notify on actions made by a member of one of the specified groups",
				Computed:    true,
				Optional:    true,
				ElementType: types.Int64Type,
				PlanModifiers: []planmodifier.List{
					listplanmodifier.UseStateForUnknown(),
				},
			},
			"triggering_user_ids": schema.ListAttribute{
				Description: "If set, will only notify on actions made one of the specified users",
				Computed:    true,
				Optional:    true,
				ElementType: types.Int64Type,
				PlanModifiers: []planmodifier.List{
					listplanmodifier.UseStateForUnknown(),
				},
			},
			"trigger_by_share_recipients": schema.BoolAttribute{
				Description: "Notify when actions are performed by a share recipient?",
				Computed:    true,
				Optional:    true,
				PlanModifiers: []planmodifier.Bool{
					boolplanmodifier.UseStateForUnknown(),
				},
			},
			"notify_user_actions": schema.BoolAttribute{
				Description: "If true, will send notifications about a user's own activity to that user.  If false, only activity performed by other users (or anonymous users) will be sent in notifications.",
				Computed:    true,
				Optional:    true,
				PlanModifiers: []planmodifier.Bool{
					boolplanmodifier.UseStateForUnknown(),
				},
			},
			"notify_on_copy": schema.BoolAttribute{
				Description: "Trigger on files copied to this path?",
				Computed:    true,
				Optional:    true,
				PlanModifiers: []planmodifier.Bool{
					boolplanmodifier.UseStateForUnknown(),
				},
			},
			"notify_on_delete": schema.BoolAttribute{
				Description: "Trigger on files deleted in this path?",
				Computed:    true,
				Optional:    true,
				PlanModifiers: []planmodifier.Bool{
					boolplanmodifier.UseStateForUnknown(),
				},
			},
			"notify_on_download": schema.BoolAttribute{
				Description: "Trigger on files downloaded in this path?",
				Computed:    true,
				Optional:    true,
				PlanModifiers: []planmodifier.Bool{
					boolplanmodifier.UseStateForUnknown(),
				},
			},
			"notify_on_move": schema.BoolAttribute{
				Description: "Trigger on files moved to this path?",
				Computed:    true,
				Optional:    true,
				PlanModifiers: []planmodifier.Bool{
					boolplanmodifier.UseStateForUnknown(),
				},
			},
			"notify_on_upload": schema.BoolAttribute{
				Description: "Trigger on files created/uploaded/updated/changed in this path?",
				Computed:    true,
				Optional:    true,
				PlanModifiers: []planmodifier.Bool{
					boolplanmodifier.UseStateForUnknown(),
				},
			},
			"recursive": schema.BoolAttribute{
				Description: "Apply notification recursively?  This will enable notifications for each subfolder.",
				Computed:    true,
				Optional:    true,
				PlanModifiers: []planmodifier.Bool{
					boolplanmodifier.UseStateForUnknown(),
				},
			},
			"send_interval": schema.StringAttribute{
				Description: "The time interval that notifications are aggregated to",
				Computed:    true,
				Optional:    true,
				Validators: []validator.String{
					stringvalidator.OneOf("five_minutes", "fifteen_minutes", "hourly", "daily"),
				},
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
			},
			"message": schema.StringAttribute{
				Description: "Custom message to include in notification emails",
				Computed:    true,
				Optional:    true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
			},
			"triggering_filenames": schema.ListAttribute{
				Description: "Array of filenames (possibly with wildcards) to scope trigger",
				Computed:    true,
				Optional:    true,
				ElementType: types.StringType,
				PlanModifiers: []planmodifier.List{
					listplanmodifier.UseStateForUnknown(),
				},
			},
			"user_id": schema.Int64Attribute{
				Description: "Notification user ID",
				Computed:    true,
				Optional:    true,
				PlanModifiers: []planmodifier.Int64{
					int64planmodifier.UseStateForUnknown(),
					int64planmodifier.RequiresReplace(),
				},
			},
			"username": schema.StringAttribute{
				Description: "Notification username",
				Computed:    true,
				Optional:    true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
					stringplanmodifier.RequiresReplace(),
				},
			},
			"id": schema.Int64Attribute{
				Description: "Notification ID",
				Computed:    true,
				PlanModifiers: []planmodifier.Int64{
					int64planmodifier.UseStateForUnknown(),
				},
			},
			"group_name": schema.StringAttribute{
				Description: "Group name, if a Group ID is set",
				Computed:    true,
			},
			"unsubscribed": schema.BoolAttribute{
				Description: "Is the user unsubscribed from this notification?",
				Computed:    true,
			},
			"unsubscribed_reason": schema.StringAttribute{
				Description: "The reason that the user unsubscribed",
				Computed:    true,
				Validators: []validator.String{
					stringvalidator.OneOf("none", "unsubscribe_link_clicked", "mail_bounced", "mail_marked_as_spam"),
				},
			},
			"suppressed_email": schema.BoolAttribute{
				Description: "If true, it means that the recipient at this user's email address has manually unsubscribed from all emails, or had their email \"hard bounce\", which means that we are unable to send mail to this user's current email address. Notifications will resume if the user changes their email address.",
				Computed:    true,
			},
		},
	}
}

func (r *notificationResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	var plan notificationResourceModel
	diags := req.Plan.Get(ctx, &plan)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}
	var config notificationResourceModel
	diags = req.Config.Get(ctx, &config)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	paramsNotificationCreate := files_sdk.NotificationCreateParams{}
	paramsNotificationCreate.UserId = plan.UserId.ValueInt64()
	if !plan.NotifyOnCopy.IsNull() && !plan.NotifyOnCopy.IsUnknown() {
		paramsNotificationCreate.NotifyOnCopy = plan.NotifyOnCopy.ValueBoolPointer()
	}
	if !plan.NotifyOnDelete.IsNull() && !plan.NotifyOnDelete.IsUnknown() {
		paramsNotificationCreate.NotifyOnDelete = plan.NotifyOnDelete.ValueBoolPointer()
	}
	if !plan.NotifyOnDownload.IsNull() && !plan.NotifyOnDownload.IsUnknown() {
		paramsNotificationCreate.NotifyOnDownload = plan.NotifyOnDownload.ValueBoolPointer()
	}
	if !plan.NotifyOnMove.IsNull() && !plan.NotifyOnMove.IsUnknown() {
		paramsNotificationCreate.NotifyOnMove = plan.NotifyOnMove.ValueBoolPointer()
	}
	if !plan.NotifyOnUpload.IsNull() && !plan.NotifyOnUpload.IsUnknown() {
		paramsNotificationCreate.NotifyOnUpload = plan.NotifyOnUpload.ValueBoolPointer()
	}
	if !plan.NotifyUserActions.IsNull() && !plan.NotifyUserActions.IsUnknown() {
		paramsNotificationCreate.NotifyUserActions = plan.NotifyUserActions.ValueBoolPointer()
	}
	if !plan.Recursive.IsNull() && !plan.Recursive.IsUnknown() {
		paramsNotificationCreate.Recursive = plan.Recursive.ValueBoolPointer()
	}
	paramsNotificationCreate.SendInterval = plan.SendInterval.ValueString()
	paramsNotificationCreate.Message = plan.Message.ValueString()
	if !plan.TriggeringFilenames.IsNull() && !plan.TriggeringFilenames.IsUnknown() {
		diags = plan.TriggeringFilenames.ElementsAs(ctx, &paramsNotificationCreate.TriggeringFilenames, false)
		resp.Diagnostics.Append(diags...)
	}
	if !plan.TriggeringGroupIds.IsNull() && !plan.TriggeringGroupIds.IsUnknown() {
		diags = plan.TriggeringGroupIds.ElementsAs(ctx, &paramsNotificationCreate.TriggeringGroupIds, false)
		resp.Diagnostics.Append(diags...)
	}
	if !plan.TriggeringUserIds.IsNull() && !plan.TriggeringUserIds.IsUnknown() {
		diags = plan.TriggeringUserIds.ElementsAs(ctx, &paramsNotificationCreate.TriggeringUserIds, false)
		resp.Diagnostics.Append(diags...)
	}
	if !plan.TriggerByShareRecipients.IsNull() && !plan.TriggerByShareRecipients.IsUnknown() {
		paramsNotificationCreate.TriggerByShareRecipients = plan.TriggerByShareRecipients.ValueBoolPointer()
	}
	paramsNotificationCreate.GroupId = plan.GroupId.ValueInt64()
	paramsNotificationCreate.Path = plan.Path.ValueString()
	paramsNotificationCreate.Username = plan.Username.ValueString()

	if resp.Diagnostics.HasError() {
		return
	}

	notification, err := r.client.Create(paramsNotificationCreate, files_sdk.WithContext(ctx))
	if err != nil {
		resp.Diagnostics.AddError(
			"Error Creating Files Notification",
			"Could not create notification, unexpected error: "+err.Error(),
		)
		return
	}

	diags = r.populateResourceModel(ctx, notification, &plan)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	diags = resp.State.Set(ctx, plan)
	resp.Diagnostics.Append(diags...)
}

func (r *notificationResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	var state notificationResourceModel
	diags := req.State.Get(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	paramsNotificationFind := files_sdk.NotificationFindParams{}
	paramsNotificationFind.Id = state.Id.ValueInt64()

	notification, err := r.client.Find(paramsNotificationFind, files_sdk.WithContext(ctx))
	if err != nil {
		if files_sdk.IsNotExist(err) {
			resp.State.RemoveResource(ctx)
			return
		}

		resp.Diagnostics.AddError(
			"Error Reading Files Notification",
			"Could not read notification id "+fmt.Sprint(state.Id.ValueInt64())+": "+err.Error(),
		)
		return
	}

	diags = r.populateResourceModel(ctx, notification, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	diags = resp.State.Set(ctx, state)
	resp.Diagnostics.Append(diags...)
}

func (r *notificationResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	var plan notificationResourceModel
	diags := req.Plan.Get(ctx, &plan)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}
	var config notificationResourceModel
	diags = req.Config.Get(ctx, &config)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	paramsNotificationUpdate := files_sdk.NotificationUpdateParams{}
	paramsNotificationUpdate.Id = plan.Id.ValueInt64()
	if !plan.NotifyOnCopy.IsNull() && !plan.NotifyOnCopy.IsUnknown() {
		paramsNotificationUpdate.NotifyOnCopy = plan.NotifyOnCopy.ValueBoolPointer()
	}
	if !plan.NotifyOnDelete.IsNull() && !plan.NotifyOnDelete.IsUnknown() {
		paramsNotificationUpdate.NotifyOnDelete = plan.NotifyOnDelete.ValueBoolPointer()
	}
	if !plan.NotifyOnDownload.IsNull() && !plan.NotifyOnDownload.IsUnknown() {
		paramsNotificationUpdate.NotifyOnDownload = plan.NotifyOnDownload.ValueBoolPointer()
	}
	if !plan.NotifyOnMove.IsNull() && !plan.NotifyOnMove.IsUnknown() {
		paramsNotificationUpdate.NotifyOnMove = plan.NotifyOnMove.ValueBoolPointer()
	}
	if !plan.NotifyOnUpload.IsNull() && !plan.NotifyOnUpload.IsUnknown() {
		paramsNotificationUpdate.NotifyOnUpload = plan.NotifyOnUpload.ValueBoolPointer()
	}
	if !plan.NotifyUserActions.IsNull() && !plan.NotifyUserActions.IsUnknown() {
		paramsNotificationUpdate.NotifyUserActions = plan.NotifyUserActions.ValueBoolPointer()
	}
	if !plan.Recursive.IsNull() && !plan.Recursive.IsUnknown() {
		paramsNotificationUpdate.Recursive = plan.Recursive.ValueBoolPointer()
	}
	paramsNotificationUpdate.SendInterval = plan.SendInterval.ValueString()
	paramsNotificationUpdate.Message = plan.Message.ValueString()
	if !plan.TriggeringFilenames.IsNull() && !plan.TriggeringFilenames.IsUnknown() {
		diags = plan.TriggeringFilenames.ElementsAs(ctx, &paramsNotificationUpdate.TriggeringFilenames, false)
		resp.Diagnostics.Append(diags...)
	}
	if !plan.TriggeringGroupIds.IsNull() && !plan.TriggeringGroupIds.IsUnknown() {
		diags = plan.TriggeringGroupIds.ElementsAs(ctx, &paramsNotificationUpdate.TriggeringGroupIds, false)
		resp.Diagnostics.Append(diags...)
	}
	if !plan.TriggeringUserIds.IsNull() && !plan.TriggeringUserIds.IsUnknown() {
		diags = plan.TriggeringUserIds.ElementsAs(ctx, &paramsNotificationUpdate.TriggeringUserIds, false)
		resp.Diagnostics.Append(diags...)
	}
	if !plan.TriggerByShareRecipients.IsNull() && !plan.TriggerByShareRecipients.IsUnknown() {
		paramsNotificationUpdate.TriggerByShareRecipients = plan.TriggerByShareRecipients.ValueBoolPointer()
	}

	if resp.Diagnostics.HasError() {
		return
	}

	notification, err := r.client.Update(paramsNotificationUpdate, files_sdk.WithContext(ctx))
	if err != nil {
		resp.Diagnostics.AddError(
			"Error Updating Files Notification",
			"Could not update notification, unexpected error: "+err.Error(),
		)
		return
	}

	diags = r.populateResourceModel(ctx, notification, &plan)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	diags = resp.State.Set(ctx, plan)
	resp.Diagnostics.Append(diags...)
}

func (r *notificationResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	var state notificationResourceModel
	diags := req.State.Get(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	paramsNotificationDelete := files_sdk.NotificationDeleteParams{}
	paramsNotificationDelete.Id = state.Id.ValueInt64()

	err := r.client.Delete(paramsNotificationDelete, files_sdk.WithContext(ctx))
	if err != nil && !files_sdk.IsNotExist(err) {
		resp.Diagnostics.AddError(
			"Error Deleting Files Notification",
			"Could not delete notification id "+fmt.Sprint(state.Id.ValueInt64())+": "+err.Error(),
		)
	}
}

func (r *notificationResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
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

func (r *notificationResource) populateResourceModel(ctx context.Context, notification files_sdk.Notification, state *notificationResourceModel) (diags diag.Diagnostics) {
	var propDiags diag.Diagnostics

	state.Id = types.Int64Value(notification.Id)
	state.Path = types.StringValue(notification.Path)
	state.GroupId = types.Int64Value(notification.GroupId)
	state.GroupName = types.StringValue(notification.GroupName)
	state.TriggeringGroupIds, propDiags = types.ListValueFrom(ctx, types.Int64Type, notification.TriggeringGroupIds)
	diags.Append(propDiags...)
	state.TriggeringUserIds, propDiags = types.ListValueFrom(ctx, types.Int64Type, notification.TriggeringUserIds)
	diags.Append(propDiags...)
	state.TriggerByShareRecipients = types.BoolPointerValue(notification.TriggerByShareRecipients)
	state.NotifyUserActions = types.BoolPointerValue(notification.NotifyUserActions)
	state.NotifyOnCopy = types.BoolPointerValue(notification.NotifyOnCopy)
	state.NotifyOnDelete = types.BoolPointerValue(notification.NotifyOnDelete)
	state.NotifyOnDownload = types.BoolPointerValue(notification.NotifyOnDownload)
	state.NotifyOnMove = types.BoolPointerValue(notification.NotifyOnMove)
	state.NotifyOnUpload = types.BoolPointerValue(notification.NotifyOnUpload)
	state.Recursive = types.BoolPointerValue(notification.Recursive)
	state.SendInterval = types.StringValue(notification.SendInterval)
	state.Message = types.StringValue(notification.Message)
	state.TriggeringFilenames, propDiags = types.ListValueFrom(ctx, types.StringType, notification.TriggeringFilenames)
	diags.Append(propDiags...)
	state.Unsubscribed = types.BoolPointerValue(notification.Unsubscribed)
	state.UnsubscribedReason = types.StringValue(notification.UnsubscribedReason)
	state.UserId = types.Int64Value(notification.UserId)
	state.Username = types.StringValue(notification.Username)
	state.SuppressedEmail = types.BoolPointerValue(notification.SuppressedEmail)

	return
}
