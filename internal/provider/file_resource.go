package provider

import (
	"context"
	"encoding/json"
	"fmt"
	"time"

	files_sdk "github.com/Files-com/files-sdk-go/v3"
	"github.com/Files-com/files-sdk-go/v3/file"
	"github.com/Files-com/terraform-provider-files/lib"
	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/dynamicplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/int64planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-log/tflog"
)

var (
	_ resource.Resource                = &fileResource{}
	_ resource.ResourceWithConfigure   = &fileResource{}
	_ resource.ResourceWithImportState = &fileResource{}
)

func NewFileResource() resource.Resource {
	return &fileResource{}
}

type fileResource struct {
	client *file.Client
}

type fileResourceModel struct {
	Source                             types.String  `tfsdk:"source"`
	Md5                                types.String  `tfsdk:"md5"`
	Path                               types.String  `tfsdk:"path"`
	CustomMetadata                     types.Dynamic `tfsdk:"custom_metadata"`
	Size                               types.Int64   `tfsdk:"size"`
	ProvidedMtime                      types.String  `tfsdk:"provided_mtime"`
	PriorityColor                      types.String  `tfsdk:"priority_color"`
	CreatedById                        types.Int64   `tfsdk:"created_by_id"`
	CreatedByApiKeyId                  types.Int64   `tfsdk:"created_by_api_key_id"`
	CreatedByAs2IncomingMessageId      types.Int64   `tfsdk:"created_by_as2_incoming_message_id"`
	CreatedByAutomationId              types.Int64   `tfsdk:"created_by_automation_id"`
	CreatedByBundleRegistrationId      types.Int64   `tfsdk:"created_by_bundle_registration_id"`
	CreatedByInboxId                   types.Int64   `tfsdk:"created_by_inbox_id"`
	CreatedByRemoteServerId            types.Int64   `tfsdk:"created_by_remote_server_id"`
	CreatedByRemoteServerSyncId        types.Int64   `tfsdk:"created_by_remote_server_sync_id"`
	DisplayName                        types.String  `tfsdk:"display_name"`
	Type                               types.String  `tfsdk:"type"`
	CreatedAt                          types.String  `tfsdk:"created_at"`
	LastModifiedById                   types.Int64   `tfsdk:"last_modified_by_id"`
	LastModifiedByApiKeyId             types.Int64   `tfsdk:"last_modified_by_api_key_id"`
	LastModifiedByAutomationId         types.Int64   `tfsdk:"last_modified_by_automation_id"`
	LastModifiedByBundleRegistrationId types.Int64   `tfsdk:"last_modified_by_bundle_registration_id"`
	LastModifiedByRemoteServerId       types.Int64   `tfsdk:"last_modified_by_remote_server_id"`
	LastModifiedByRemoteServerSyncId   types.Int64   `tfsdk:"last_modified_by_remote_server_sync_id"`
	Mtime                              types.String  `tfsdk:"mtime"`
	Crc32                              types.String  `tfsdk:"crc32"`
	MimeType                           types.String  `tfsdk:"mime_type"`
	Region                             types.String  `tfsdk:"region"`
	Permissions                        types.String  `tfsdk:"permissions"`
	SubfoldersLocked                   types.Bool    `tfsdk:"subfolders_locked"`
	IsLocked                           types.Bool    `tfsdk:"is_locked"`
	DownloadUri                        types.String  `tfsdk:"download_uri"`
	PreviewId                          types.Int64   `tfsdk:"preview_id"`
	Preview                            types.String  `tfsdk:"preview"`
}

func (r *fileResource) Configure(_ context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {
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

	r.client = &file.Client{Config: sdk_config}
}

func (r *fileResource) Metadata(_ context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_file"
}

func (r *fileResource) Schema(_ context.Context, _ resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Description: "",
		Attributes: map[string]schema.Attribute{
			"source": schema.StringAttribute{
				Description: "Path to a file that will be read and uploaded.",
				Required:    true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.RequiresReplace(),
				},
			},
			"md5": schema.StringAttribute{
				Description: "File MD5 checksum.",
				Computed:    true,
				Optional:    true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.RequiresReplace(),
					stringplanmodifier.UseStateForUnknown(),
				},
			},
			"path": schema.StringAttribute{
				Description: "File/Folder path. This must be slash-delimited, but it must neither start nor end with a slash. Maximum of 5000 characters.",
				Required:    true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
			},
			"custom_metadata": schema.DynamicAttribute{
				Description: "Custom metadata map of keys and values. Limited to 32 keys, 256 characters per key and 1024 characters per value.",
				Computed:    true,
				Optional:    true,
				PlanModifiers: []planmodifier.Dynamic{
					dynamicplanmodifier.UseStateForUnknown(),
				},
			},
			"size": schema.Int64Attribute{
				Description: "File/Folder size",
				Computed:    true,
				Optional:    true,
				PlanModifiers: []planmodifier.Int64{
					int64planmodifier.UseStateForUnknown(),
					int64planmodifier.RequiresReplace(),
				},
			},
			"provided_mtime": schema.StringAttribute{
				Description: "File last modified date/time, according to the client who set it.  Files.com allows desktop, FTP, SFTP, and WebDAV clients to set modified at times.  This allows Desktop<->Cloud syncing to preserve modified at times.",
				Computed:    true,
				Optional:    true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
			},
			"priority_color": schema.StringAttribute{
				Description: "Bookmark/priority color of file/folder",
				Computed:    true,
				Optional:    true,
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
			"display_name": schema.StringAttribute{
				Description: "File/Folder display name",
				Computed:    true,
			},
			"type": schema.StringAttribute{
				Description: "Type: `directory` or `file`.",
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
			"crc32": schema.StringAttribute{
				Description: "File CRC32 checksum. This is sometimes delayed, so if you get a blank response, wait and try again.",
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

func (r *fileResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	var plan fileResourceModel
	diags := req.Plan.Get(ctx, &plan)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	err := r.client.Upload(file.UploadWithContext(ctx), file.UploadWithFile(plan.Source.ValueString()), file.UploadWithDestinationPath(plan.Path.ValueString()))
	if err != nil {
		resp.Diagnostics.AddError(
			"Error Creating Files File",
			"Could not create file, unexpected error: "+err.Error(),
		)
		return
	}

	paramsFileUpdate := files_sdk.FileUpdateParams{}
	paramsFileUpdate.Path = plan.Path.ValueString()
	updateCustomMetadata, diags := lib.DynamicToStringMap(ctx, path.Root("custom_metadata"), plan.CustomMetadata)
	resp.Diagnostics.Append(diags...)
	paramsFileUpdate.CustomMetadata = updateCustomMetadata
	if !plan.ProvidedMtime.IsNull() {
		if plan.ProvidedMtime.ValueString() == "" {
			paramsFileUpdate.ProvidedMtime = new(time.Time)
		} else {
			updateProvidedMtime, err := time.Parse(time.RFC3339, plan.ProvidedMtime.ValueString())
			if err != nil {
				resp.Diagnostics.AddAttributeError(
					path.Root("provided_mtime"),
					"Error Parsing provided_mtime Time",
					"Could not parse provided_mtime time: "+err.Error(),
				)
			} else {
				paramsFileUpdate.ProvidedMtime = &updateProvidedMtime
			}
		}
	}
	paramsFileUpdate.PriorityColor = plan.PriorityColor.ValueString()

	if resp.Diagnostics.HasError() {
		return
	}

	file, err := r.client.Update(paramsFileUpdate, files_sdk.WithContext(ctx))
	if err != nil {
		resp.Diagnostics.AddError(
			"Error Updating Files File",
			"Could not update file, unexpected error: "+err.Error(),
		)
		return
	}

	diags = r.populateResourceModel(ctx, file, &plan)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	diags = resp.State.Set(ctx, plan)
	resp.Diagnostics.Append(diags...)
}

func (r *fileResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	var state fileResourceModel
	diags := req.State.Get(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	withPriorityColor := true
	paramsFileFind := files_sdk.FileFindParams{
		Path:              state.Path.ValueString(),
		WithPriorityColor: &withPriorityColor,
	}

	file, err := r.client.Find(paramsFileFind, files_sdk.WithContext(ctx))
	if err != nil {
		if files_sdk.IsNotExist(err) {
			resp.State.RemoveResource(ctx)
			return
		}

		resp.Diagnostics.AddError(
			"Error Reading Files File",
			"Could not read file path "+fmt.Sprint(state.Path.ValueString())+": "+err.Error(),
		)
		return
	}

	diags = r.populateResourceModel(ctx, file, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	diags = resp.State.Set(ctx, state)
	resp.Diagnostics.Append(diags...)
}

func (r *fileResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	var plan fileResourceModel
	diags := req.Plan.Get(ctx, &plan)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	var state fileResourceModel
	diags = req.State.Get(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	if plan.Path.ValueString() != state.Path.ValueString() {
		tflog.Info(ctx, "Detected path change, moving file", map[string]interface{}{
			"path":        state.Path.ValueString(),
			"destination": plan.Path.ValueString(),
		})
		paramsFileMove := files_sdk.FileMoveParams{
			Path:        state.Path.ValueString(),
			Destination: plan.Path.ValueString(),
		}
		_, err := r.client.Move(paramsFileMove, files_sdk.WithContext(ctx))
		if err != nil {
			resp.Diagnostics.AddError(
				"Error Moving Files File",
				"Could not move file path "+fmt.Sprint(state.Path.ValueString())+": "+err.Error(),
			)
			return
		}
	}

	paramsFileUpdate := files_sdk.FileUpdateParams{}
	paramsFileUpdate.Path = plan.Path.ValueString()
	updateCustomMetadata, diags := lib.DynamicToStringMap(ctx, path.Root("custom_metadata"), plan.CustomMetadata)
	resp.Diagnostics.Append(diags...)
	paramsFileUpdate.CustomMetadata = updateCustomMetadata
	if !plan.ProvidedMtime.IsNull() {
		if plan.ProvidedMtime.ValueString() == "" {
			paramsFileUpdate.ProvidedMtime = new(time.Time)
		} else {
			updateProvidedMtime, err := time.Parse(time.RFC3339, plan.ProvidedMtime.ValueString())
			if err != nil {
				resp.Diagnostics.AddAttributeError(
					path.Root("provided_mtime"),
					"Error Parsing provided_mtime Time",
					"Could not parse provided_mtime time: "+err.Error(),
				)
			} else {
				paramsFileUpdate.ProvidedMtime = &updateProvidedMtime
			}
		}
	}
	paramsFileUpdate.PriorityColor = plan.PriorityColor.ValueString()

	if resp.Diagnostics.HasError() {
		return
	}

	file, err := r.client.Update(paramsFileUpdate, files_sdk.WithContext(ctx))
	if err != nil {
		resp.Diagnostics.AddError(
			"Error Updating Files File",
			"Could not update file, unexpected error: "+err.Error(),
		)
		return
	}

	diags = r.populateResourceModel(ctx, file, &plan)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	diags = resp.State.Set(ctx, plan)
	resp.Diagnostics.Append(diags...)
}

func (r *fileResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	var state fileResourceModel
	diags := req.State.Get(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	paramsFileDelete := files_sdk.FileDeleteParams{
		Path: state.Path.ValueString(),
	}

	err := r.client.Delete(paramsFileDelete, files_sdk.WithContext(ctx))
	if err != nil && !files_sdk.IsNotExist(err) {
		resp.Diagnostics.AddError(
			"Error Deleting Files File",
			"Could not delete file path "+fmt.Sprint(state.Path.ValueString())+": "+err.Error(),
		)
	}
}

func (r *fileResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	resource.ImportStatePassthroughID(ctx, path.Root("path"), req, resp)
}

func (r *fileResource) populateResourceModel(ctx context.Context, file files_sdk.File, state *fileResourceModel) (diags diag.Diagnostics) {
	var propDiags diag.Diagnostics

	state.Path = types.StringValue(file.Path)
	state.CreatedById = types.Int64Value(file.CreatedById)
	state.CreatedByApiKeyId = types.Int64Value(file.CreatedByApiKeyId)
	state.CreatedByAs2IncomingMessageId = types.Int64Value(file.CreatedByAs2IncomingMessageId)
	state.CreatedByAutomationId = types.Int64Value(file.CreatedByAutomationId)
	state.CreatedByBundleRegistrationId = types.Int64Value(file.CreatedByBundleRegistrationId)
	state.CreatedByInboxId = types.Int64Value(file.CreatedByInboxId)
	state.CreatedByRemoteServerId = types.Int64Value(file.CreatedByRemoteServerId)
	state.CreatedByRemoteServerSyncId = types.Int64Value(file.CreatedByRemoteServerSyncId)
	state.CustomMetadata, propDiags = lib.ToDynamic(ctx, path.Root("custom_metadata"), file.CustomMetadata, state.CustomMetadata.UnderlyingValue())
	diags.Append(propDiags...)
	state.DisplayName = types.StringValue(file.DisplayName)
	state.Type = types.StringValue(file.Type)
	state.Size = types.Int64Value(file.Size)
	if err := lib.TimeToStringType(ctx, path.Root("created_at"), file.CreatedAt, &state.CreatedAt); err != nil {
		diags.AddError(
			"Error Creating Files File",
			"Could not convert state created_at to string: "+err.Error(),
		)
	}
	state.LastModifiedById = types.Int64Value(file.LastModifiedById)
	state.LastModifiedByApiKeyId = types.Int64Value(file.LastModifiedByApiKeyId)
	state.LastModifiedByAutomationId = types.Int64Value(file.LastModifiedByAutomationId)
	state.LastModifiedByBundleRegistrationId = types.Int64Value(file.LastModifiedByBundleRegistrationId)
	state.LastModifiedByRemoteServerId = types.Int64Value(file.LastModifiedByRemoteServerId)
	state.LastModifiedByRemoteServerSyncId = types.Int64Value(file.LastModifiedByRemoteServerSyncId)
	if err := lib.TimeToStringType(ctx, path.Root("mtime"), file.Mtime, &state.Mtime); err != nil {
		diags.AddError(
			"Error Creating Files File",
			"Could not convert state mtime to string: "+err.Error(),
		)
	}
	if err := lib.TimeToStringType(ctx, path.Root("provided_mtime"), file.ProvidedMtime, &state.ProvidedMtime); err != nil {
		diags.AddError(
			"Error Creating Files File",
			"Could not convert state provided_mtime to string: "+err.Error(),
		)
	}
	state.Crc32 = types.StringValue(file.Crc32)
	state.MimeType = types.StringValue(file.MimeType)
	state.Region = types.StringValue(file.Region)
	state.Permissions = types.StringValue(file.Permissions)
	state.SubfoldersLocked = types.BoolPointerValue(file.SubfoldersLocked)
	state.IsLocked = types.BoolPointerValue(file.IsLocked)
	state.DownloadUri = types.StringValue(file.DownloadUri)
	state.PriorityColor = types.StringValue(file.PriorityColor)
	state.PreviewId = types.Int64Value(file.PreviewId)
	respPreview, err := json.Marshal(file.Preview)
	if err != nil {
		diags.AddError(
			"Error Creating Files File",
			"Could not marshal preview to JSON: "+err.Error(),
		)
	}
	state.Preview = types.StringValue(string(respPreview))

	if state.Md5.IsNull() || state.Md5.IsUnknown() || file.Md5 != "" {
		state.Md5 = types.StringValue(file.Md5)
	}

	return
}
