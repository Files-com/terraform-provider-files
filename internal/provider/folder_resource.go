package provider

import (
	"context"
	"encoding/json"
	"fmt"
	"time"

	files_sdk "github.com/Files-com/files-sdk-go/v3"
	"github.com/Files-com/files-sdk-go/v3/file"
	"github.com/Files-com/files-sdk-go/v3/folder"
	"github.com/Files-com/terraform-provider-files/lib"
	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/boolplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/dynamicplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-log/tflog"
)

var (
	_ resource.Resource                = &folderResource{}
	_ resource.ResourceWithConfigure   = &folderResource{}
	_ resource.ResourceWithImportState = &folderResource{}
)

func NewFolderResource() resource.Resource {
	return &folderResource{}
}

type folderResource struct {
	folderClient *folder.Client
	fileClient   *file.Client
}

type folderResourceModel struct {
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
	MimeType                           types.String  `tfsdk:"mime_type"`
	Region                             types.String  `tfsdk:"region"`
	Permissions                        types.String  `tfsdk:"permissions"`
	SubfoldersLocked                   types.Bool    `tfsdk:"subfolders_locked"`
	IsLocked                           types.Bool    `tfsdk:"is_locked"`
	DownloadUri                        types.String  `tfsdk:"download_uri"`
	PriorityColor                      types.String  `tfsdk:"priority_color"`
	PreviewId                          types.Int64   `tfsdk:"preview_id"`
	Preview                            types.String  `tfsdk:"preview"`
	MkdirParents                       types.Bool    `tfsdk:"mkdir_parents"`
}

func (r *folderResource) Configure(_ context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {
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

func (r *folderResource) Metadata(_ context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_folder"
}

func (r *folderResource) Schema(_ context.Context, _ resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Description: "",
		Attributes: map[string]schema.Attribute{
			"path": schema.StringAttribute{
				Description: "File/Folder path This must be slash-delimited, but it must neither start nor end with a slash. Maximum of 5000 characters.",
				Required:    true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
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
				Optional:    true,
				PlanModifiers: []planmodifier.Dynamic{
					dynamicplanmodifier.UseStateForUnknown(),
				},
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
				Optional:    true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
			},
			"crc32": schema.StringAttribute{
				Description: "File CRC32 checksum. This is sometimes delayed, so if you get a blank response, wait and try again.",
				Computed:    true,
			},
			"md5": schema.StringAttribute{
				Description: "File MD5 checksum. This is sometimes delayed, so if you get a blank response, wait and try again.",
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
				Optional:    true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
			},
			"preview_id": schema.Int64Attribute{
				Description: "File preview ID",
				Computed:    true,
			},
			"preview": schema.StringAttribute{
				Description: "File preview",
				Computed:    true,
			},
			"mkdir_parents": schema.BoolAttribute{
				Description: "Create parent directories if they do not exist?",
				Optional:    true,
				PlanModifiers: []planmodifier.Bool{
					boolplanmodifier.RequiresReplace(),
				},
			},
		},
	}
}

func (r *folderResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	var plan folderResourceModel
	diags := req.Plan.Get(ctx, &plan)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	paramsFolderCreate := files_sdk.FolderCreateParams{}
	paramsFolderCreate.Path = plan.Path.ValueString()
	paramsFolderCreate.MkdirParents = plan.MkdirParents.ValueBoolPointer()
	if !plan.ProvidedMtime.IsNull() && plan.ProvidedMtime.ValueString() != "" {
		createProvidedMtime, err := time.Parse(time.RFC3339, plan.ProvidedMtime.ValueString())
		if err != nil {
			resp.Diagnostics.AddAttributeError(
				path.Root("provided_mtime"),
				"Error Parsing provided_mtime Time",
				"Could not parse provided_mtime time: "+err.Error(),
			)
		} else {
			paramsFolderCreate.ProvidedMtime = &createProvidedMtime
		}
	}

	if resp.Diagnostics.HasError() {
		return
	}

	_, err := r.folderClient.Create(paramsFolderCreate, files_sdk.WithContext(ctx))
	if err != nil {
		resp.Diagnostics.AddError(
			"Error Creating Files Folder",
			"Could not create folder, unexpected error: "+err.Error(),
		)
		return
	}

	paramsFolderUpdate := files_sdk.FileUpdateParams{}
	paramsFolderUpdate.Path = plan.Path.ValueString()
	updateCustomMetadata, diags := lib.DynamicToStringMap(ctx, path.Root("custom_metadata"), plan.CustomMetadata)
	resp.Diagnostics.Append(diags...)
	paramsFolderUpdate.CustomMetadata = updateCustomMetadata
	if !plan.ProvidedMtime.IsNull() && plan.ProvidedMtime.ValueString() != "" {
		updateProvidedMtime, err := time.Parse(time.RFC3339, plan.ProvidedMtime.ValueString())
		if err != nil {
			resp.Diagnostics.AddAttributeError(
				path.Root("provided_mtime"),
				"Error Parsing provided_mtime Time",
				"Could not parse provided_mtime time: "+err.Error(),
			)
		} else {
			paramsFolderUpdate.ProvidedMtime = &updateProvidedMtime
		}
	}
	paramsFolderUpdate.PriorityColor = plan.PriorityColor.ValueString()

	if resp.Diagnostics.HasError() {
		return
	}

	folder, err := r.fileClient.Update(paramsFolderUpdate, files_sdk.WithContext(ctx))
	if err != nil {
		resp.Diagnostics.AddError(
			"Error Updating Files Folder",
			"Could not update folder, unexpected error: "+err.Error(),
		)
		return
	}

	diags = r.populateResourceModel(ctx, folder, &plan)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	diags = resp.State.Set(ctx, plan)
	resp.Diagnostics.Append(diags...)
}

func (r *folderResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	var state folderResourceModel
	diags := req.State.Get(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	withPriorityColor := true
	paramsFolderFind := files_sdk.FileFindParams{
		Path:              state.Path.ValueString(),
		WithPriorityColor: &withPriorityColor,
	}

	folder, err := r.fileClient.Find(paramsFolderFind, files_sdk.WithContext(ctx))
	if err != nil {
		resp.Diagnostics.AddError(
			"Error Reading Files Folder",
			"Could not read folder path "+fmt.Sprint(state.Path.ValueString())+": "+err.Error(),
		)
		return
	}

	diags = r.populateResourceModel(ctx, folder, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	diags = resp.State.Set(ctx, state)
	resp.Diagnostics.Append(diags...)
}

func (r *folderResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	var plan folderResourceModel
	diags := req.Plan.Get(ctx, &plan)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	var state folderResourceModel
	diags = req.State.Get(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	if plan.Path.ValueString() != state.Path.ValueString() {
		tflog.Info(ctx, "Detected path change, moving folder", map[string]interface{}{
			"path":        state.Path.ValueString(),
			"destination": plan.Path.ValueString(),
		})
		paramsFolderMove := files_sdk.FileMoveParams{
			Path:        state.Path.ValueString(),
			Destination: plan.Path.ValueString(),
		}
		_, err := r.fileClient.Move(paramsFolderMove, files_sdk.WithContext(ctx))
		if err != nil {
			resp.Diagnostics.AddError(
				"Error Moving Files Folder",
				"Could not move folder path "+fmt.Sprint(state.Path.ValueString())+": "+err.Error(),
			)
			return
		}
	}

	paramsFolderUpdate := files_sdk.FileUpdateParams{}
	paramsFolderUpdate.Path = plan.Path.ValueString()
	updateCustomMetadata, diags := lib.DynamicToStringMap(ctx, path.Root("custom_metadata"), plan.CustomMetadata)
	resp.Diagnostics.Append(diags...)
	paramsFolderUpdate.CustomMetadata = updateCustomMetadata
	if !plan.ProvidedMtime.IsNull() && plan.ProvidedMtime.ValueString() != "" {
		updateProvidedMtime, err := time.Parse(time.RFC3339, plan.ProvidedMtime.ValueString())
		if err != nil {
			resp.Diagnostics.AddAttributeError(
				path.Root("provided_mtime"),
				"Error Parsing provided_mtime Time",
				"Could not parse provided_mtime time: "+err.Error(),
			)
		} else {
			paramsFolderUpdate.ProvidedMtime = &updateProvidedMtime
		}
	}
	paramsFolderUpdate.PriorityColor = plan.PriorityColor.ValueString()

	if resp.Diagnostics.HasError() {
		return
	}

	folder, err := r.fileClient.Update(paramsFolderUpdate, files_sdk.WithContext(ctx))
	if err != nil {
		resp.Diagnostics.AddError(
			"Error Updating Files Folder",
			"Could not update folder, unexpected error: "+err.Error(),
		)
		return
	}

	diags = r.populateResourceModel(ctx, folder, &plan)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	diags = resp.State.Set(ctx, plan)
	resp.Diagnostics.Append(diags...)
}

func (r *folderResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	var state folderResourceModel
	diags := req.State.Get(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	paramsFolderDelete := files_sdk.FileDeleteParams{
		Path: state.Path.ValueString(),
	}

	err := r.fileClient.Delete(paramsFolderDelete, files_sdk.WithContext(ctx))
	if err != nil {
		resp.Diagnostics.AddError(
			"Error Deleting Files Folder",
			"Could not delete folder path "+fmt.Sprint(state.Path.ValueString())+": "+err.Error(),
		)
	}
}

func (r *folderResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	resource.ImportStatePassthroughID(ctx, path.Root("path"), req, resp)
}

func (r *folderResource) populateResourceModel(ctx context.Context, folder files_sdk.File, state *folderResourceModel) (diags diag.Diagnostics) {
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
