package provider

import (
	"context"
	"fmt"

	files_sdk "github.com/Files-com/files-sdk-go/v3"
	notification "github.com/Files-com/files-sdk-go/v3/notification"
	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/datasource/schema"
	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

var (
	_ datasource.DataSource              = &notificationDataSource{}
	_ datasource.DataSourceWithConfigure = &notificationDataSource{}
)

func NewNotificationDataSource() datasource.DataSource {
	return &notificationDataSource{}
}

type notificationDataSource struct {
	client *notification.Client
}

type notificationDataSourceModel struct {
	Id                       types.Int64  `tfsdk:"id"`
	Path                     types.String `tfsdk:"path"`
	GroupId                  types.Int64  `tfsdk:"group_id"`
	GroupName                types.String `tfsdk:"group_name"`
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
	Unsubscribed             types.Bool   `tfsdk:"unsubscribed"`
	UnsubscribedReason       types.String `tfsdk:"unsubscribed_reason"`
	UserId                   types.Int64  `tfsdk:"user_id"`
	Username                 types.String `tfsdk:"username"`
	SuppressedEmail          types.Bool   `tfsdk:"suppressed_email"`
}

func (r *notificationDataSource) Configure(_ context.Context, req datasource.ConfigureRequest, resp *datasource.ConfigureResponse) {
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

func (r *notificationDataSource) Metadata(_ context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_notification"
}

func (r *notificationDataSource) Schema(_ context.Context, _ datasource.SchemaRequest, resp *datasource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Description: "Notifications are our feature that send E-Mails when new files are uploaded into a folder.",
		Attributes: map[string]schema.Attribute{
			"id": schema.Int64Attribute{
				Description: "Notification ID",
				Required:    true,
			},
			"path": schema.StringAttribute{
				Description: "Folder path to notify on This must be slash-delimited, but it must neither start nor end with a slash. Maximum of 5000 characters.",
				Computed:    true,
			},
			"group_id": schema.Int64Attribute{
				Description: "ID of Group to receive notifications",
				Computed:    true,
			},
			"group_name": schema.StringAttribute{
				Description: "Group name, if a Group ID is set",
				Computed:    true,
			},
			"triggering_group_ids": schema.ListAttribute{
				Description: "If set, will only notify on actions made by a member of one of the specified groups",
				Computed:    true,
				ElementType: types.Int64Type,
			},
			"triggering_user_ids": schema.ListAttribute{
				Description: "If set, will onlynotify on actions made one of the specified users",
				Computed:    true,
				ElementType: types.Int64Type,
			},
			"trigger_by_share_recipients": schema.BoolAttribute{
				Description: "Notify when actions are performed by a share recipient?",
				Computed:    true,
			},
			"notify_user_actions": schema.BoolAttribute{
				Description: "If true, will send notifications about a user's own activity to that user.  If false, only activity performed by other users (or anonymous users) will be sent in notifications.",
				Computed:    true,
			},
			"notify_on_copy": schema.BoolAttribute{
				Description: "Trigger on files copied to this path?",
				Computed:    true,
			},
			"notify_on_delete": schema.BoolAttribute{
				Description: "Trigger on files deleted in this path?",
				Computed:    true,
			},
			"notify_on_download": schema.BoolAttribute{
				Description: "Trigger on files downloaded in this path?",
				Computed:    true,
			},
			"notify_on_move": schema.BoolAttribute{
				Description: "Trigger on files moved to this path?",
				Computed:    true,
			},
			"notify_on_upload": schema.BoolAttribute{
				Description: "Trigger on files created/uploaded/updated/changed in this path?",
				Computed:    true,
			},
			"recursive": schema.BoolAttribute{
				Description: "Apply notification recursively?  This will enable notifications for each subfolder.",
				Computed:    true,
			},
			"send_interval": schema.StringAttribute{
				Description: "The time interval that notifications are aggregated to",
				Computed:    true,
			},
			"message": schema.StringAttribute{
				Description: "Custom message to include in notification emails",
				Computed:    true,
			},
			"triggering_filenames": schema.ListAttribute{
				Description: "Array of filenames (possibly with wildcards) to scope trigger",
				Computed:    true,
				ElementType: types.StringType,
			},
			"unsubscribed": schema.BoolAttribute{
				Description: "Is the user unsubscribed from this notification?",
				Computed:    true,
			},
			"unsubscribed_reason": schema.StringAttribute{
				Description: "The reason that the user unsubscribed",
				Computed:    true,
			},
			"user_id": schema.Int64Attribute{
				Description: "Notification user ID",
				Computed:    true,
			},
			"username": schema.StringAttribute{
				Description: "Notification username",
				Computed:    true,
			},
			"suppressed_email": schema.BoolAttribute{
				Description: "If true, it means that the recipient at this user's email address has manually unsubscribed from all emails, or had their email \"hard bounce\", which means that we are unable to send mail to this user's current email address. Notifications will resume if the user changes their email address.",
				Computed:    true,
			},
		},
	}
}

func (r *notificationDataSource) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
	var data notificationDataSourceModel
	diags := req.Config.Get(ctx, &data)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	paramsNotificationFind := files_sdk.NotificationFindParams{}
	paramsNotificationFind.Id = data.Id.ValueInt64()

	notification, err := r.client.Find(paramsNotificationFind, files_sdk.WithContext(ctx))
	if err != nil {
		resp.Diagnostics.AddError(
			"Error Reading Files Notification",
			"Could not read notification id "+fmt.Sprint(data.Id.ValueInt64())+": "+err.Error(),
		)
		return
	}

	diags = r.populateDataSourceModel(ctx, notification, &data)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	diags = resp.State.Set(ctx, data)
	resp.Diagnostics.Append(diags...)
}

func (r *notificationDataSource) populateDataSourceModel(ctx context.Context, notification files_sdk.Notification, state *notificationDataSourceModel) (diags diag.Diagnostics) {
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
