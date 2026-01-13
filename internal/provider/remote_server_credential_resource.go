package provider

import (
	"context"
	"fmt"
	"strconv"
	"strings"

	files_sdk "github.com/Files-com/files-sdk-go/v3"
	remote_server_credential "github.com/Files-com/files-sdk-go/v3/remoteservercredential"
	"github.com/hashicorp/terraform-plugin-framework-validators/stringvalidator"
	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/int64planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/schema/validator"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

var (
	_ resource.Resource                = &remoteServerCredentialResource{}
	_ resource.ResourceWithConfigure   = &remoteServerCredentialResource{}
	_ resource.ResourceWithImportState = &remoteServerCredentialResource{}
)

func NewRemoteServerCredentialResource() resource.Resource {
	return &remoteServerCredentialResource{}
}

type remoteServerCredentialResource struct {
	client *remote_server_credential.Client
}

type remoteServerCredentialResourceModel struct {
	WorkspaceId                             types.Int64  `tfsdk:"workspace_id"`
	Name                                    types.String `tfsdk:"name"`
	Description                             types.String `tfsdk:"description"`
	ServerType                              types.String `tfsdk:"server_type"`
	AwsAccessKey                            types.String `tfsdk:"aws_access_key"`
	GoogleCloudStorageS3CompatibleAccessKey types.String `tfsdk:"google_cloud_storage_s3_compatible_access_key"`
	WasabiAccessKey                         types.String `tfsdk:"wasabi_access_key"`
	S3CompatibleAccessKey                   types.String `tfsdk:"s3_compatible_access_key"`
	FilebaseAccessKey                       types.String `tfsdk:"filebase_access_key"`
	CloudflareAccessKey                     types.String `tfsdk:"cloudflare_access_key"`
	LinodeAccessKey                         types.String `tfsdk:"linode_access_key"`
	Username                                types.String `tfsdk:"username"`
	Password                                types.String `tfsdk:"password"`
	PrivateKey                              types.String `tfsdk:"private_key"`
	PrivateKeyPassphrase                    types.String `tfsdk:"private_key_passphrase"`
	AwsSecretKey                            types.String `tfsdk:"aws_secret_key"`
	AzureBlobStorageAccessKey               types.String `tfsdk:"azure_blob_storage_access_key"`
	AzureBlobStorageSasToken                types.String `tfsdk:"azure_blob_storage_sas_token"`
	AzureFilesStorageAccessKey              types.String `tfsdk:"azure_files_storage_access_key"`
	AzureFilesStorageSasToken               types.String `tfsdk:"azure_files_storage_sas_token"`
	BackblazeB2ApplicationKey               types.String `tfsdk:"backblaze_b2_application_key"`
	BackblazeB2KeyId                        types.String `tfsdk:"backblaze_b2_key_id"`
	CloudflareSecretKey                     types.String `tfsdk:"cloudflare_secret_key"`
	FilebaseSecretKey                       types.String `tfsdk:"filebase_secret_key"`
	GoogleCloudStorageCredentialsJson       types.String `tfsdk:"google_cloud_storage_credentials_json"`
	GoogleCloudStorageS3CompatibleSecretKey types.String `tfsdk:"google_cloud_storage_s3_compatible_secret_key"`
	LinodeSecretKey                         types.String `tfsdk:"linode_secret_key"`
	S3CompatibleSecretKey                   types.String `tfsdk:"s3_compatible_secret_key"`
	WasabiSecretKey                         types.String `tfsdk:"wasabi_secret_key"`
	Id                                      types.Int64  `tfsdk:"id"`
}

func (r *remoteServerCredentialResource) Configure(_ context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {
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

	r.client = &remote_server_credential.Client{Config: sdk_config}
}

func (r *remoteServerCredentialResource) Metadata(_ context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_remote_server_credential"
}

func (r *remoteServerCredentialResource) Schema(_ context.Context, _ resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Description: "A RemoteServerCredential is a way to store a credential for Remote Servers in a centralized vault and then reference it from Remote Server definitions.\n\n\n\nThis allows you to manage your credentials in one place and avoid duplicating them across multiple Remote Server configurations. It also enhances security by allowing you to use Terraform or APIs for Remote Server management without having to worry about credential exposure.",
		Attributes: map[string]schema.Attribute{
			"workspace_id": schema.Int64Attribute{
				Description: "Workspace ID (0 for default workspace)",
				Computed:    true,
				Optional:    true,
				PlanModifiers: []planmodifier.Int64{
					int64planmodifier.UseStateForUnknown(),
					int64planmodifier.RequiresReplace(),
				},
			},
			"name": schema.StringAttribute{
				Description: "Internal name for your reference",
				Computed:    true,
				Optional:    true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
			},
			"description": schema.StringAttribute{
				Description: "Internal description for your reference",
				Computed:    true,
				Optional:    true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
			},
			"server_type": schema.StringAttribute{
				Description: "Remote server type.  Remote Server Credentials are only valid for a single type of Remote Server.",
				Computed:    true,
				Optional:    true,
				Validators: []validator.String{
					stringvalidator.OneOf("ftp", "sftp", "s3", "google_cloud_storage", "webdav", "wasabi", "backblaze_b2", "one_drive", "box", "dropbox", "google_drive", "azure", "sharepoint", "s3_compatible", "azure_files", "files_agent", "filebase", "cloudflare", "linode"),
				},
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
			},
			"aws_access_key": schema.StringAttribute{
				Description: "AWS Access Key.",
				Computed:    true,
				Optional:    true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
			},
			"google_cloud_storage_s3_compatible_access_key": schema.StringAttribute{
				Description: "Google Cloud Storage: S3-compatible Access Key.",
				Computed:    true,
				Optional:    true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
			},
			"wasabi_access_key": schema.StringAttribute{
				Description: "Wasabi: Access Key.",
				Computed:    true,
				Optional:    true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
			},
			"s3_compatible_access_key": schema.StringAttribute{
				Description: "S3-compatible: Access Key",
				Computed:    true,
				Optional:    true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
			},
			"filebase_access_key": schema.StringAttribute{
				Description: "Filebase: Access Key.",
				Computed:    true,
				Optional:    true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
			},
			"cloudflare_access_key": schema.StringAttribute{
				Description: "Cloudflare: Access Key.",
				Computed:    true,
				Optional:    true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
			},
			"linode_access_key": schema.StringAttribute{
				Description: "Linode: Access Key",
				Computed:    true,
				Optional:    true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
			},
			"username": schema.StringAttribute{
				Description: "Remote server username.",
				Computed:    true,
				Optional:    true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
			},
			"password": schema.StringAttribute{
				Description: "Password, if needed.",
				Optional:    true,
				WriteOnly:   true,
			},
			"private_key": schema.StringAttribute{
				Description: "Private key, if needed.",
				Optional:    true,
				WriteOnly:   true,
			},
			"private_key_passphrase": schema.StringAttribute{
				Description: "Passphrase for private key if needed.",
				Optional:    true,
				WriteOnly:   true,
			},
			"aws_secret_key": schema.StringAttribute{
				Description: "AWS: secret key.",
				Optional:    true,
				WriteOnly:   true,
			},
			"azure_blob_storage_access_key": schema.StringAttribute{
				Description: "Azure Blob Storage: Access Key",
				Optional:    true,
				WriteOnly:   true,
			},
			"azure_blob_storage_sas_token": schema.StringAttribute{
				Description: "Azure Blob Storage: Shared Access Signature (SAS) token",
				Optional:    true,
				WriteOnly:   true,
			},
			"azure_files_storage_access_key": schema.StringAttribute{
				Description: "Azure File Storage: Access Key",
				Optional:    true,
				WriteOnly:   true,
			},
			"azure_files_storage_sas_token": schema.StringAttribute{
				Description: "Azure File Storage: Shared Access Signature (SAS) token",
				Optional:    true,
				WriteOnly:   true,
			},
			"backblaze_b2_application_key": schema.StringAttribute{
				Description: "Backblaze B2 Cloud Storage: applicationKey",
				Optional:    true,
				WriteOnly:   true,
			},
			"backblaze_b2_key_id": schema.StringAttribute{
				Description: "Backblaze B2 Cloud Storage: keyID",
				Optional:    true,
				WriteOnly:   true,
			},
			"cloudflare_secret_key": schema.StringAttribute{
				Description: "Cloudflare: Secret Key",
				Optional:    true,
				WriteOnly:   true,
			},
			"filebase_secret_key": schema.StringAttribute{
				Description: "Filebase: Secret Key",
				Optional:    true,
				WriteOnly:   true,
			},
			"google_cloud_storage_credentials_json": schema.StringAttribute{
				Description: "Google Cloud Storage: JSON file that contains the private key. To generate see https://cloud.google.com/storage/docs/json_api/v1/how-tos/authorizing#APIKey",
				Optional:    true,
				WriteOnly:   true,
			},
			"google_cloud_storage_s3_compatible_secret_key": schema.StringAttribute{
				Description: "Google Cloud Storage: S3-compatible secret key",
				Optional:    true,
				WriteOnly:   true,
			},
			"linode_secret_key": schema.StringAttribute{
				Description: "Linode: Secret Key",
				Optional:    true,
				WriteOnly:   true,
			},
			"s3_compatible_secret_key": schema.StringAttribute{
				Description: "S3-compatible: Secret Key",
				Optional:    true,
				WriteOnly:   true,
			},
			"wasabi_secret_key": schema.StringAttribute{
				Description: "Wasabi: Secret Key",
				Optional:    true,
				WriteOnly:   true,
			},
			"id": schema.Int64Attribute{
				Description: "Remote Server Credential ID",
				Computed:    true,
				PlanModifiers: []planmodifier.Int64{
					int64planmodifier.UseStateForUnknown(),
				},
			},
		},
	}
}

func (r *remoteServerCredentialResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	var plan remoteServerCredentialResourceModel
	diags := req.Plan.Get(ctx, &plan)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}
	var config remoteServerCredentialResourceModel
	diags = req.Config.Get(ctx, &config)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	paramsRemoteServerCredentialCreate := files_sdk.RemoteServerCredentialCreateParams{}
	paramsRemoteServerCredentialCreate.Name = plan.Name.ValueString()
	paramsRemoteServerCredentialCreate.Description = plan.Description.ValueString()
	paramsRemoteServerCredentialCreate.ServerType = paramsRemoteServerCredentialCreate.ServerType.Enum()[plan.ServerType.ValueString()]
	paramsRemoteServerCredentialCreate.AwsAccessKey = plan.AwsAccessKey.ValueString()
	paramsRemoteServerCredentialCreate.CloudflareAccessKey = plan.CloudflareAccessKey.ValueString()
	paramsRemoteServerCredentialCreate.FilebaseAccessKey = plan.FilebaseAccessKey.ValueString()
	paramsRemoteServerCredentialCreate.GoogleCloudStorageS3CompatibleAccessKey = plan.GoogleCloudStorageS3CompatibleAccessKey.ValueString()
	paramsRemoteServerCredentialCreate.LinodeAccessKey = plan.LinodeAccessKey.ValueString()
	paramsRemoteServerCredentialCreate.S3CompatibleAccessKey = plan.S3CompatibleAccessKey.ValueString()
	paramsRemoteServerCredentialCreate.Username = plan.Username.ValueString()
	paramsRemoteServerCredentialCreate.WasabiAccessKey = plan.WasabiAccessKey.ValueString()
	paramsRemoteServerCredentialCreate.Password = config.Password.ValueString()
	paramsRemoteServerCredentialCreate.PrivateKey = config.PrivateKey.ValueString()
	paramsRemoteServerCredentialCreate.PrivateKeyPassphrase = config.PrivateKeyPassphrase.ValueString()
	paramsRemoteServerCredentialCreate.AwsSecretKey = config.AwsSecretKey.ValueString()
	paramsRemoteServerCredentialCreate.AzureBlobStorageAccessKey = config.AzureBlobStorageAccessKey.ValueString()
	paramsRemoteServerCredentialCreate.AzureBlobStorageSasToken = config.AzureBlobStorageSasToken.ValueString()
	paramsRemoteServerCredentialCreate.AzureFilesStorageAccessKey = config.AzureFilesStorageAccessKey.ValueString()
	paramsRemoteServerCredentialCreate.AzureFilesStorageSasToken = config.AzureFilesStorageSasToken.ValueString()
	paramsRemoteServerCredentialCreate.BackblazeB2ApplicationKey = config.BackblazeB2ApplicationKey.ValueString()
	paramsRemoteServerCredentialCreate.BackblazeB2KeyId = config.BackblazeB2KeyId.ValueString()
	paramsRemoteServerCredentialCreate.CloudflareSecretKey = config.CloudflareSecretKey.ValueString()
	paramsRemoteServerCredentialCreate.FilebaseSecretKey = config.FilebaseSecretKey.ValueString()
	paramsRemoteServerCredentialCreate.GoogleCloudStorageCredentialsJson = config.GoogleCloudStorageCredentialsJson.ValueString()
	paramsRemoteServerCredentialCreate.GoogleCloudStorageS3CompatibleSecretKey = config.GoogleCloudStorageS3CompatibleSecretKey.ValueString()
	paramsRemoteServerCredentialCreate.LinodeSecretKey = config.LinodeSecretKey.ValueString()
	paramsRemoteServerCredentialCreate.S3CompatibleSecretKey = config.S3CompatibleSecretKey.ValueString()
	paramsRemoteServerCredentialCreate.WasabiSecretKey = config.WasabiSecretKey.ValueString()
	paramsRemoteServerCredentialCreate.WorkspaceId = plan.WorkspaceId.ValueInt64()

	if resp.Diagnostics.HasError() {
		return
	}

	remoteServerCredential, err := r.client.Create(paramsRemoteServerCredentialCreate, files_sdk.WithContext(ctx))
	if err != nil {
		resp.Diagnostics.AddError(
			"Error Creating Files RemoteServerCredential",
			"Could not create remote_server_credential, unexpected error: "+err.Error(),
		)
		return
	}

	diags = r.populateResourceModel(ctx, remoteServerCredential, &plan)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	diags = resp.State.Set(ctx, plan)
	resp.Diagnostics.Append(diags...)
}

func (r *remoteServerCredentialResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	var state remoteServerCredentialResourceModel
	diags := req.State.Get(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	paramsRemoteServerCredentialFind := files_sdk.RemoteServerCredentialFindParams{}
	paramsRemoteServerCredentialFind.Id = state.Id.ValueInt64()

	remoteServerCredential, err := r.client.Find(paramsRemoteServerCredentialFind, files_sdk.WithContext(ctx))
	if err != nil {
		if files_sdk.IsNotExist(err) {
			resp.State.RemoveResource(ctx)
			return
		}

		resp.Diagnostics.AddError(
			"Error Reading Files RemoteServerCredential",
			"Could not read remote_server_credential id "+fmt.Sprint(state.Id.ValueInt64())+": "+err.Error(),
		)
		return
	}

	diags = r.populateResourceModel(ctx, remoteServerCredential, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	diags = resp.State.Set(ctx, state)
	resp.Diagnostics.Append(diags...)
}

func (r *remoteServerCredentialResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	var plan remoteServerCredentialResourceModel
	diags := req.Plan.Get(ctx, &plan)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}
	var config remoteServerCredentialResourceModel
	diags = req.Config.Get(ctx, &config)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	paramsRemoteServerCredentialUpdate := map[string]interface{}{}
	if !plan.Id.IsNull() && !plan.Id.IsUnknown() {
		paramsRemoteServerCredentialUpdate["id"] = plan.Id.ValueInt64()
	}
	if !config.Name.IsNull() && !config.Name.IsUnknown() {
		paramsRemoteServerCredentialUpdate["name"] = config.Name.ValueString()
	}
	if !config.Description.IsNull() && !config.Description.IsUnknown() {
		paramsRemoteServerCredentialUpdate["description"] = config.Description.ValueString()
	}
	if !config.ServerType.IsNull() && !config.ServerType.IsUnknown() {
		paramsRemoteServerCredentialUpdate["server_type"] = config.ServerType.ValueString()
	}
	if !config.AwsAccessKey.IsNull() && !config.AwsAccessKey.IsUnknown() {
		paramsRemoteServerCredentialUpdate["aws_access_key"] = config.AwsAccessKey.ValueString()
	}
	if !config.CloudflareAccessKey.IsNull() && !config.CloudflareAccessKey.IsUnknown() {
		paramsRemoteServerCredentialUpdate["cloudflare_access_key"] = config.CloudflareAccessKey.ValueString()
	}
	if !config.FilebaseAccessKey.IsNull() && !config.FilebaseAccessKey.IsUnknown() {
		paramsRemoteServerCredentialUpdate["filebase_access_key"] = config.FilebaseAccessKey.ValueString()
	}
	if !config.GoogleCloudStorageS3CompatibleAccessKey.IsNull() && !config.GoogleCloudStorageS3CompatibleAccessKey.IsUnknown() {
		paramsRemoteServerCredentialUpdate["google_cloud_storage_s3_compatible_access_key"] = config.GoogleCloudStorageS3CompatibleAccessKey.ValueString()
	}
	if !config.LinodeAccessKey.IsNull() && !config.LinodeAccessKey.IsUnknown() {
		paramsRemoteServerCredentialUpdate["linode_access_key"] = config.LinodeAccessKey.ValueString()
	}
	if !config.S3CompatibleAccessKey.IsNull() && !config.S3CompatibleAccessKey.IsUnknown() {
		paramsRemoteServerCredentialUpdate["s3_compatible_access_key"] = config.S3CompatibleAccessKey.ValueString()
	}
	if !config.Username.IsNull() && !config.Username.IsUnknown() {
		paramsRemoteServerCredentialUpdate["username"] = config.Username.ValueString()
	}
	if !config.WasabiAccessKey.IsNull() && !config.WasabiAccessKey.IsUnknown() {
		paramsRemoteServerCredentialUpdate["wasabi_access_key"] = config.WasabiAccessKey.ValueString()
	}
	if !config.Password.IsNull() && !config.Password.IsUnknown() {
		paramsRemoteServerCredentialUpdate["password"] = config.Password.ValueString()
	}
	if !config.PrivateKey.IsNull() && !config.PrivateKey.IsUnknown() {
		paramsRemoteServerCredentialUpdate["private_key"] = config.PrivateKey.ValueString()
	}
	if !config.PrivateKeyPassphrase.IsNull() && !config.PrivateKeyPassphrase.IsUnknown() {
		paramsRemoteServerCredentialUpdate["private_key_passphrase"] = config.PrivateKeyPassphrase.ValueString()
	}
	if !config.AwsSecretKey.IsNull() && !config.AwsSecretKey.IsUnknown() {
		paramsRemoteServerCredentialUpdate["aws_secret_key"] = config.AwsSecretKey.ValueString()
	}
	if !config.AzureBlobStorageAccessKey.IsNull() && !config.AzureBlobStorageAccessKey.IsUnknown() {
		paramsRemoteServerCredentialUpdate["azure_blob_storage_access_key"] = config.AzureBlobStorageAccessKey.ValueString()
	}
	if !config.AzureBlobStorageSasToken.IsNull() && !config.AzureBlobStorageSasToken.IsUnknown() {
		paramsRemoteServerCredentialUpdate["azure_blob_storage_sas_token"] = config.AzureBlobStorageSasToken.ValueString()
	}
	if !config.AzureFilesStorageAccessKey.IsNull() && !config.AzureFilesStorageAccessKey.IsUnknown() {
		paramsRemoteServerCredentialUpdate["azure_files_storage_access_key"] = config.AzureFilesStorageAccessKey.ValueString()
	}
	if !config.AzureFilesStorageSasToken.IsNull() && !config.AzureFilesStorageSasToken.IsUnknown() {
		paramsRemoteServerCredentialUpdate["azure_files_storage_sas_token"] = config.AzureFilesStorageSasToken.ValueString()
	}
	if !config.BackblazeB2ApplicationKey.IsNull() && !config.BackblazeB2ApplicationKey.IsUnknown() {
		paramsRemoteServerCredentialUpdate["backblaze_b2_application_key"] = config.BackblazeB2ApplicationKey.ValueString()
	}
	if !config.BackblazeB2KeyId.IsNull() && !config.BackblazeB2KeyId.IsUnknown() {
		paramsRemoteServerCredentialUpdate["backblaze_b2_key_id"] = config.BackblazeB2KeyId.ValueString()
	}
	if !config.CloudflareSecretKey.IsNull() && !config.CloudflareSecretKey.IsUnknown() {
		paramsRemoteServerCredentialUpdate["cloudflare_secret_key"] = config.CloudflareSecretKey.ValueString()
	}
	if !config.FilebaseSecretKey.IsNull() && !config.FilebaseSecretKey.IsUnknown() {
		paramsRemoteServerCredentialUpdate["filebase_secret_key"] = config.FilebaseSecretKey.ValueString()
	}
	if !config.GoogleCloudStorageCredentialsJson.IsNull() && !config.GoogleCloudStorageCredentialsJson.IsUnknown() {
		paramsRemoteServerCredentialUpdate["google_cloud_storage_credentials_json"] = config.GoogleCloudStorageCredentialsJson.ValueString()
	}
	if !config.GoogleCloudStorageS3CompatibleSecretKey.IsNull() && !config.GoogleCloudStorageS3CompatibleSecretKey.IsUnknown() {
		paramsRemoteServerCredentialUpdate["google_cloud_storage_s3_compatible_secret_key"] = config.GoogleCloudStorageS3CompatibleSecretKey.ValueString()
	}
	if !config.LinodeSecretKey.IsNull() && !config.LinodeSecretKey.IsUnknown() {
		paramsRemoteServerCredentialUpdate["linode_secret_key"] = config.LinodeSecretKey.ValueString()
	}
	if !config.S3CompatibleSecretKey.IsNull() && !config.S3CompatibleSecretKey.IsUnknown() {
		paramsRemoteServerCredentialUpdate["s3_compatible_secret_key"] = config.S3CompatibleSecretKey.ValueString()
	}
	if !config.WasabiSecretKey.IsNull() && !config.WasabiSecretKey.IsUnknown() {
		paramsRemoteServerCredentialUpdate["wasabi_secret_key"] = config.WasabiSecretKey.ValueString()
	}

	if resp.Diagnostics.HasError() {
		return
	}

	remoteServerCredential, err := r.client.UpdateWithMap(paramsRemoteServerCredentialUpdate, files_sdk.WithContext(ctx))
	if err != nil {
		resp.Diagnostics.AddError(
			"Error Updating Files RemoteServerCredential",
			"Could not update remote_server_credential, unexpected error: "+err.Error(),
		)
		return
	}

	diags = r.populateResourceModel(ctx, remoteServerCredential, &plan)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	diags = resp.State.Set(ctx, plan)
	resp.Diagnostics.Append(diags...)
}

func (r *remoteServerCredentialResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	var state remoteServerCredentialResourceModel
	diags := req.State.Get(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	paramsRemoteServerCredentialDelete := files_sdk.RemoteServerCredentialDeleteParams{}
	paramsRemoteServerCredentialDelete.Id = state.Id.ValueInt64()

	err := r.client.Delete(paramsRemoteServerCredentialDelete, files_sdk.WithContext(ctx))
	if err != nil && !files_sdk.IsNotExist(err) {
		resp.Diagnostics.AddError(
			"Error Deleting Files RemoteServerCredential",
			"Could not delete remote_server_credential id "+fmt.Sprint(state.Id.ValueInt64())+": "+err.Error(),
		)
	}
}

func (r *remoteServerCredentialResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
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

func (r *remoteServerCredentialResource) populateResourceModel(ctx context.Context, remoteServerCredential files_sdk.RemoteServerCredential, state *remoteServerCredentialResourceModel) (diags diag.Diagnostics) {
	state.Id = types.Int64Value(remoteServerCredential.Id)
	state.WorkspaceId = types.Int64Value(remoteServerCredential.WorkspaceId)
	state.Name = types.StringValue(remoteServerCredential.Name)
	state.Description = types.StringValue(remoteServerCredential.Description)
	state.ServerType = types.StringValue(remoteServerCredential.ServerType)
	state.AwsAccessKey = types.StringValue(remoteServerCredential.AwsAccessKey)
	state.GoogleCloudStorageS3CompatibleAccessKey = types.StringValue(remoteServerCredential.GoogleCloudStorageS3CompatibleAccessKey)
	state.WasabiAccessKey = types.StringValue(remoteServerCredential.WasabiAccessKey)
	state.S3CompatibleAccessKey = types.StringValue(remoteServerCredential.S3CompatibleAccessKey)
	state.FilebaseAccessKey = types.StringValue(remoteServerCredential.FilebaseAccessKey)
	state.CloudflareAccessKey = types.StringValue(remoteServerCredential.CloudflareAccessKey)
	state.LinodeAccessKey = types.StringValue(remoteServerCredential.LinodeAccessKey)
	state.Username = types.StringValue(remoteServerCredential.Username)

	return
}
