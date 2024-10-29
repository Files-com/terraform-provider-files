package provider

import (
	"context"
	"encoding/json"
	"fmt"

	files_sdk "github.com/Files-com/files-sdk-go/v3"
	"github.com/Files-com/files-sdk-go/v3/file"
	"github.com/Files-com/files-sdk-go/v3/folder"
	"github.com/Files-com/terraform-provider-files/lib"
	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/datasource/schema"
	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

var (
	_ datasource.DataSource              = &folderDataSource{}
	_ datasource.DataSourceWithConfigure = &folderDataSource{}
)

func NewFolderDataSource() datasource.DataSource {
	return &folderDataSource{}
}

type folderDataSource struct {
	folderClient *folder.Client
	fileClient   *file.Client
}

type folderDataSourceModel struct {
	Path                               types.String  `tfsdk:"path"`
	CreatedById                        types.Int64   `tfsdk:"created_by_id"`
	CreatedByApiKeyId                  types.Int64   `tfsdk:"created_by_api_key_id"`
	CreatedByAs2IncomingMessageId      types.Int64   `tfsdk:"created_by_as2_incoming_message_id"`
	CreatedByAutomationId              types.Int64   `tfsdk:"created_by_automation_id"`
	CreatedByBundleRegistrationId      types.Int64   `tfsdk:"created_by_bundle_registration_id"`
	CreatedByInboxId                   types.Int64   `tfsdk:"created_by_inbox_id"`
	CreatedByRemoteServerId            types.Int64   `tfsdk:"created_by_remote_server_id"`
	CreatedByRemoteServerSyncId        types.Int64   `tfsdk:"created_by_remote_server_sync_id"`
	CustomMetadata                     types.Dynamic `tfsdk:"custom_metadata"`
	DisplayName                        types.String  `tfsdk:"display_name"`
	Type                               types.String  `tfsdk:"type"`
	Size                               types.Int64   `tfsdk:"size"`
	CreatedAt                          types.String  `tfsdk:"created_at"`
	LastModifiedById                   types.Int64   `tfsdk:"last_modified_by_id"`
	LastModifiedByApiKeyId             types.Int64   `tfsdk:"last_modified_by_api_key_id"`
	LastModifiedByAutomationId         types.Int64   `tfsdk:"last_modified_by_automation_id"`
	LastModifiedByBundleRegistrationId types.Int64   `tfsdk:"last_modified_by_bundle_registration_id"`
	LastModifiedByRemoteServerId       types.Int64   `tfsdk:"last_modified_by_remote_server_id"`
	LastModifiedByRemoteServerSyncId   types.Int64   `tfsdk:"last_modified_by_remote_server_sync_id"`
	Mtime                              types.String  `tfsdk:"mtime"`
	ProvidedMtime                      types.String  `tfsdk:"provided_mtime"`
	Crc32                              types.String  `tfsdk:"crc32"`
	Md5                                types.String  `tfsdk:"md5"`
	Sha1                               types.String  `tfsdk:"sha1"`
	Sha256                             types.String  `tfsdk:"sha256"`
	MimeType                           types.String  `tfsdk:"mime_type"`
	Region                             types.String  `tfsdk:"region"`
	Permissions                        types.String  `tfsdk:"permissions"`
	SubfoldersLocked                   types.Bool    `tfsdk:"subfolders_locked"`
	IsLocked                           types.Bool    `tfsdk:"is_locked"`
	DownloadUri                        types.String  `tfsdk:"download_uri"`
	PriorityColor                      types.String  `tfsdk:"priority_color"`
	PreviewId                          types.Int64   `tfsdk:"preview_id"`
	Preview                            types.String  `tfsdk:"preview"`
}

func (r *folderDataSource) Configure(_ context.Context, req datasource.ConfigureRequest, resp *datasource.ConfigureResponse) {
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

	r.folderClient = &folder.Client{Config: sdk_config}
	r.fileClient = &file.Client{Config: sdk_config}
}

func (r *folderDataSource) Metadata(_ context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_folder"
}

func (r *folderDataSource) Schema(_ context.Context, _ datasource.SchemaRequest, resp *datasource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Description: "",
		Attributes: map[string]schema.Attribute{
			"path": schema.StringAttribute{
				Description: "File/Folder path. This must be slash-delimited, but it must neither start nor end with a slash. Maximum of 5000 characters.",
				Required:    true,
			},
			"created_by_id": schema.Int64Attribute{
				Description: "User ID of the User who created the file/folder",
				Computed:    true,
			},
			"created_by_api_key_id": schema.Int64Attribute{
				Description: "ID of the API key that created the file/folder",
				Computed:    true,
			},
			"created_by_as2_incoming_message_id": schema.Int64Attribute{
				Description: "ID of the AS2 Incoming Message that created the file/folder",
				Computed:    true,
			},
			"created_by_automation_id": schema.Int64Attribute{
				Description: "ID of the Automation that created the file/folder",
				Computed:    true,
			},
			"created_by_bundle_registration_id": schema.Int64Attribute{
				Description: "ID of the Bundle Registration that created the file/folder",
				Computed:    true,
			},
			"created_by_inbox_id": schema.Int64Attribute{
				Description: "ID of the Inbox that created the file/folder",
				Computed:    true,
			},
			"created_by_remote_server_id": schema.Int64Attribute{
				Description: "ID of the Remote Server that created the file/folder",
				Computed:    true,
			},
			"created_by_remote_server_sync_id": schema.Int64Attribute{
				Description: "ID of the Remote Server Sync that created the file/folder",
				Computed:    true,
			},
			"custom_metadata": schema.DynamicAttribute{
				Description: "Custom metadata map of keys and values. Limited to 32 keys, 256 characters per key and 1024 characters per value.",
				Computed:    true,
			},
			"display_name": schema.StringAttribute{
				Description: "File/Folder display name",
				Computed:    true,
			},
			"type": schema.StringAttribute{
				Description: "Type: `directory` or `file`.",
				Computed:    true,
			},
			"size": schema.Int64Attribute{
				Description: "File/Folder size",
				Computed:    true,
			},
			"created_at": schema.StringAttribute{
				Description: "File created date/time",
				Computed:    true,
			},
			"last_modified_by_id": schema.Int64Attribute{
				Description: "User ID of the User who last modified the file/folder",
				Computed:    true,
			},
			"last_modified_by_api_key_id": schema.Int64Attribute{
				Description: "ID of the API key that last modified the file/folder",
				Computed:    true,
			},
			"last_modified_by_automation_id": schema.Int64Attribute{
				Description: "ID of the Automation that last modified the file/folder",
				Computed:    true,
			},
			"last_modified_by_bundle_registration_id": schema.Int64Attribute{
				Description: "ID of the Bundle Registration that last modified the file/folder",
				Computed:    true,
			},
			"last_modified_by_remote_server_id": schema.Int64Attribute{
				Description: "ID of the Remote Server that last modified the file/folder",
				Computed:    true,
			},
			"last_modified_by_remote_server_sync_id": schema.Int64Attribute{
				Description: "ID of the Remote Server Sync that last modified the file/folder",
				Computed:    true,
			},
			"mtime": schema.StringAttribute{
				Description: "File last modified date/time, according to the server.  This is the timestamp of the last Files.com operation of the file, regardless of what modified timestamp was sent.",
				Computed:    true,
			},
			"provided_mtime": schema.StringAttribute{
				Description: "File last modified date/time, according to the client who set it.  Files.com allows desktop, FTP, SFTP, and WebDAV clients to set modified at times.  This allows Desktop<->Cloud syncing to preserve modified at times.",
				Computed:    true,
			},
			"crc32": schema.StringAttribute{
				Description: "File CRC32 checksum. This is sometimes delayed, so if you get a blank response, wait and try again.",
				Computed:    true,
			},
			"md5": schema.StringAttribute{
				Description: "File MD5 checksum. This is sometimes delayed, so if you get a blank response, wait and try again.",
				Computed:    true,
			},
			"sha1": schema.StringAttribute{
				Description: "File SHA1 checksum. This is sometimes delayed, so if you get a blank response, wait and try again.",
				Computed:    true,
			},
			"sha256": schema.StringAttribute{
				Description: "File SHA256 checksum. This is sometimes delayed, so if you get a blank response, wait and try again.",
				Computed:    true,
			},
			"mime_type": schema.StringAttribute{
				Description: "MIME Type.  This is determined by the filename extension and is not stored separately internally.",
				Computed:    true,
			},
			"region": schema.StringAttribute{
				Description: "Region location",
				Computed:    true,
			},
			"permissions": schema.StringAttribute{
				Description: "A short string representing the current user's permissions.  Can be `r` (Read),`w` (Write),`d` (Delete), `l` (List) or any combination",
				Computed:    true,
			},
			"subfolders_locked": schema.BoolAttribute{
				Description: "Are subfolders locked and unable to be modified?",
				Computed:    true,
			},
			"is_locked": schema.BoolAttribute{
				Description: "Is this folder locked and unable to be modified?",
				Computed:    true,
			},
			"download_uri": schema.StringAttribute{
				Description: "Link to download file. Provided only in response to a download request.",
				Computed:    true,
			},
			"priority_color": schema.StringAttribute{
				Description: "Bookmark/priority color of file/folder",
				Computed:    true,
			},
			"preview_id": schema.Int64Attribute{
				Description: "File preview ID",
				Computed:    true,
			},
			"preview": schema.StringAttribute{
				Description: "File preview",
				Computed:    true,
			},
		},
	}
}

func (r *folderDataSource) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
	var data folderDataSourceModel
	diags := req.Config.Get(ctx, &data)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	withPriorityColor := true
	paramsFolderFind := files_sdk.FileFindParams{
		Path:              data.Path.ValueString(),
		WithPriorityColor: &withPriorityColor,
	}

	folder, err := r.fileClient.Find(paramsFolderFind, files_sdk.WithContext(ctx))
	if err != nil {
		resp.Diagnostics.AddError(
			"Error Reading Files Folder",
			"Could not read folder path "+fmt.Sprint(data.Path.ValueString())+": "+err.Error(),
		)
		return
	}

	diags = r.populateDataSourceModel(ctx, folder, &data)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	diags = resp.State.Set(ctx, data)
	resp.Diagnostics.Append(diags...)
}

func (r *folderDataSource) populateDataSourceModel(ctx context.Context, folder files_sdk.File, state *folderDataSourceModel) (diags diag.Diagnostics) {
	var propDiags diag.Diagnostics

	state.Path = types.StringValue(folder.Path)
	state.CreatedById = types.Int64Value(folder.CreatedById)
	state.CreatedByApiKeyId = types.Int64Value(folder.CreatedByApiKeyId)
	state.CreatedByAs2IncomingMessageId = types.Int64Value(folder.CreatedByAs2IncomingMessageId)
	state.CreatedByAutomationId = types.Int64Value(folder.CreatedByAutomationId)
	state.CreatedByBundleRegistrationId = types.Int64Value(folder.CreatedByBundleRegistrationId)
	state.CreatedByInboxId = types.Int64Value(folder.CreatedByInboxId)
	state.CreatedByRemoteServerId = types.Int64Value(folder.CreatedByRemoteServerId)
	state.CreatedByRemoteServerSyncId = types.Int64Value(folder.CreatedByRemoteServerSyncId)
	state.CustomMetadata, propDiags = lib.ToDynamic(ctx, path.Root("custom_metadata"), folder.CustomMetadata, state.CustomMetadata.UnderlyingValue())
	diags.Append(propDiags...)
	state.DisplayName = types.StringValue(folder.DisplayName)
	state.Type = types.StringValue(folder.Type)
	state.Size = types.Int64Value(folder.Size)
	if err := lib.TimeToStringType(ctx, path.Root("created_at"), folder.CreatedAt, &state.CreatedAt); err != nil {
		diags.AddError(
			"Error Creating Files Folder",
			"Could not convert state created_at to string: "+err.Error(),
		)
	}
	state.LastModifiedById = types.Int64Value(folder.LastModifiedById)
	state.LastModifiedByApiKeyId = types.Int64Value(folder.LastModifiedByApiKeyId)
	state.LastModifiedByAutomationId = types.Int64Value(folder.LastModifiedByAutomationId)
	state.LastModifiedByBundleRegistrationId = types.Int64Value(folder.LastModifiedByBundleRegistrationId)
	state.LastModifiedByRemoteServerId = types.Int64Value(folder.LastModifiedByRemoteServerId)
	state.LastModifiedByRemoteServerSyncId = types.Int64Value(folder.LastModifiedByRemoteServerSyncId)
	if err := lib.TimeToStringType(ctx, path.Root("mtime"), folder.Mtime, &state.Mtime); err != nil {
		diags.AddError(
			"Error Creating Files Folder",
			"Could not convert state mtime to string: "+err.Error(),
		)
	}
	if err := lib.TimeToStringType(ctx, path.Root("provided_mtime"), folder.ProvidedMtime, &state.ProvidedMtime); err != nil {
		diags.AddError(
			"Error Creating Files Folder",
			"Could not convert state provided_mtime to string: "+err.Error(),
		)
	}
	state.Crc32 = types.StringValue(folder.Crc32)
	state.Md5 = types.StringValue(folder.Md5)
	state.Sha1 = types.StringValue(folder.Sha1)
	state.Sha256 = types.StringValue(folder.Sha256)
	state.MimeType = types.StringValue(folder.MimeType)
	state.Region = types.StringValue(folder.Region)
	state.Permissions = types.StringValue(folder.Permissions)
	state.SubfoldersLocked = types.BoolPointerValue(folder.SubfoldersLocked)
	state.IsLocked = types.BoolPointerValue(folder.IsLocked)
	state.DownloadUri = types.StringValue(folder.DownloadUri)
	state.PriorityColor = types.StringValue(folder.PriorityColor)
	state.PreviewId = types.Int64Value(folder.PreviewId)
	respPreview, err := json.Marshal(folder.Preview)
	if err != nil {
		diags.AddError(
			"Error Creating Files Folder",
			"Could not marshal preview to JSON: "+err.Error(),
		)
	}
	state.Preview = types.StringValue(string(respPreview))

	return
}
