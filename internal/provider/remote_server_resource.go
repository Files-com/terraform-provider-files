package provider

import (
	"context"
	"fmt"
	"strconv"
	"strings"

	files_sdk "github.com/Files-com/files-sdk-go/v3"
	remote_server "github.com/Files-com/files-sdk-go/v3/remoteserver"
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
	_ resource.Resource                = &remoteServerResource{}
	_ resource.ResourceWithConfigure   = &remoteServerResource{}
	_ resource.ResourceWithImportState = &remoteServerResource{}
)

func NewRemoteServerResource() resource.Resource {
	return &remoteServerResource{}
}

type remoteServerResource struct {
	client *remote_server.Client
}

type remoteServerResourceModel struct {
	Hostname                                types.String `tfsdk:"hostname"`
	Name                                    types.String `tfsdk:"name"`
	Port                                    types.Int64  `tfsdk:"port"`
	BufferUploads                           types.String `tfsdk:"buffer_uploads"`
	MaxConnections                          types.Int64  `tfsdk:"max_connections"`
	PinToSiteRegion                         types.Bool   `tfsdk:"pin_to_site_region"`
	S3Bucket                                types.String `tfsdk:"s3_bucket"`
	S3Region                                types.String `tfsdk:"s3_region"`
	AwsAccessKey                            types.String `tfsdk:"aws_access_key"`
	ServerCertificate                       types.String `tfsdk:"server_certificate"`
	ServerHostKey                           types.String `tfsdk:"server_host_key"`
	ServerType                              types.String `tfsdk:"server_type"`
	Ssl                                     types.String `tfsdk:"ssl"`
	Username                                types.String `tfsdk:"username"`
	GoogleCloudStorageBucket                types.String `tfsdk:"google_cloud_storage_bucket"`
	GoogleCloudStorageProjectId             types.String `tfsdk:"google_cloud_storage_project_id"`
	GoogleCloudStorageS3CompatibleAccessKey types.String `tfsdk:"google_cloud_storage_s3_compatible_access_key"`
	BackblazeB2S3Endpoint                   types.String `tfsdk:"backblaze_b2_s3_endpoint"`
	BackblazeB2Bucket                       types.String `tfsdk:"backblaze_b2_bucket"`
	WasabiBucket                            types.String `tfsdk:"wasabi_bucket"`
	WasabiRegion                            types.String `tfsdk:"wasabi_region"`
	WasabiAccessKey                         types.String `tfsdk:"wasabi_access_key"`
	OneDriveAccountType                     types.String `tfsdk:"one_drive_account_type"`
	AzureBlobStorageAccount                 types.String `tfsdk:"azure_blob_storage_account"`
	AzureBlobStorageContainer               types.String `tfsdk:"azure_blob_storage_container"`
	AzureBlobStorageHierarchicalNamespace   types.Bool   `tfsdk:"azure_blob_storage_hierarchical_namespace"`
	AzureBlobStorageDnsSuffix               types.String `tfsdk:"azure_blob_storage_dns_suffix"`
	AzureFilesStorageAccount                types.String `tfsdk:"azure_files_storage_account"`
	AzureFilesStorageShareName              types.String `tfsdk:"azure_files_storage_share_name"`
	AzureFilesStorageDnsSuffix              types.String `tfsdk:"azure_files_storage_dns_suffix"`
	S3CompatibleBucket                      types.String `tfsdk:"s3_compatible_bucket"`
	S3CompatibleEndpoint                    types.String `tfsdk:"s3_compatible_endpoint"`
	S3CompatibleRegion                      types.String `tfsdk:"s3_compatible_region"`
	S3CompatibleAccessKey                   types.String `tfsdk:"s3_compatible_access_key"`
	EnableDedicatedIps                      types.Bool   `tfsdk:"enable_dedicated_ips"`
	FilesAgentPermissionSet                 types.String `tfsdk:"files_agent_permission_set"`
	FilesAgentRoot                          types.String `tfsdk:"files_agent_root"`
	FilesAgentVersion                       types.String `tfsdk:"files_agent_version"`
	FilebaseBucket                          types.String `tfsdk:"filebase_bucket"`
	FilebaseAccessKey                       types.String `tfsdk:"filebase_access_key"`
	CloudflareBucket                        types.String `tfsdk:"cloudflare_bucket"`
	CloudflareAccessKey                     types.String `tfsdk:"cloudflare_access_key"`
	CloudflareEndpoint                      types.String `tfsdk:"cloudflare_endpoint"`
	DropboxTeams                            types.Bool   `tfsdk:"dropbox_teams"`
	LinodeBucket                            types.String `tfsdk:"linode_bucket"`
	LinodeAccessKey                         types.String `tfsdk:"linode_access_key"`
	LinodeRegion                            types.String `tfsdk:"linode_region"`
	Password                                types.String `tfsdk:"password"`
	PrivateKey                              types.String `tfsdk:"private_key"`
	PrivateKeyPassphrase                    types.String `tfsdk:"private_key_passphrase"`
	ResetAuthentication                     types.Bool   `tfsdk:"reset_authentication"`
	SslCertificate                          types.String `tfsdk:"ssl_certificate"`
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
	Disabled                                types.Bool   `tfsdk:"disabled"`
	AuthenticationMethod                    types.String `tfsdk:"authentication_method"`
	RemoteHomePath                          types.String `tfsdk:"remote_home_path"`
	PinnedRegion                            types.String `tfsdk:"pinned_region"`
	AuthStatus                              types.String `tfsdk:"auth_status"`
	AuthAccountName                         types.String `tfsdk:"auth_account_name"`
	FilesAgentApiToken                      types.String `tfsdk:"files_agent_api_token"`
	SupportsVersioning                      types.Bool   `tfsdk:"supports_versioning"`
}

func (r *remoteServerResource) Configure(_ context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {
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

	r.client = &remote_server.Client{Config: sdk_config}
}

func (r *remoteServerResource) Metadata(_ context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_remote_server"
}

func (r *remoteServerResource) Schema(_ context.Context, _ resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Description: "A RemoteServer is a specific type of Behavior called `remote_server_sync`.\n\n\n\nRemote Servers can be either an FTP server, SFTP server, S3 bucket, Google Cloud Storage, Wasabi, Backblaze B2 Cloud Storage, Rackspace Cloud Files container, WebDAV, Box, Dropbox, OneDrive, Google Drive, or Azure Blob Storage.\n\n\n\nNot every attribute will apply to every remote server.\n\n\n\nFTP Servers require that you specify their `hostname`, `port`, `username`, `password`, and a value for `ssl`. Optionally, provide `server_certificate`.\n\n\n\nSFTP Servers require that you specify their `hostname`, `port`, `username`, `password` or `private_key`, and a value for `ssl`. Optionally, provide `server_certificate`, `private_key_passphrase`.\n\n\n\nS3 Buckets require that you specify their `s3_bucket` name, and `s3_region`. Optionally provide a `aws_access_key`, and `aws_secret_key`. If you don't provide credentials, you will need to use AWS to grant us access to your bucket.\n\n\n\nS3-Compatible Buckets require that you specify `s3_compatible_bucket`, `s3_compatible_endpoint`, `s3_compatible_access_key`, and `s3_compatible_secret_key`.\n\n\n\nGoogle Cloud Storage requires that you specify `google_cloud_storage_bucket`, and then one of the following sets of authentication credentials:\n\n - for JSON authentcation: `google_cloud_storage_project_id`, and `google_cloud_storage_credentials_json`\n\n - for HMAC (S3-Compatible) authentication: `google_cloud_storage_s3_compatible_access_key`, and `google_cloud_storage_s3_compatible_secret_key`\n\n\n\nWasabi requires `wasabi_bucket`, `wasabi_region`, `wasabi_access_key`, and `wasabi_secret_key`.\n\n\n\nBackblaze B2 Cloud Storage `backblaze_b2_bucket`, `backblaze_b2_s3_endpoint`, `backblaze_b2_application_key`, and `backblaze_b2_key_id`. (Requires S3 Compatible API) See https://help.backblaze.com/hc/en-us/articles/360047425453\n\n\n\nWebDAV Servers require that you specify their `hostname`, `username`, and `password`.\n\n\n\nOneDrive follow the `auth_setup_link` and login with Microsoft.\n\n\n\nSharepoint follow the `auth_setup_link` and login with Microsoft.\n\n\n\nBox follow the `auth_setup_link` and login with Box.\n\n\n\nDropbox specify if `dropbox_teams` then follow the `auth_setup_link` and login with Dropbox.\n\n\n\nGoogle Drive follow the `auth_setup_link` and login with Google.\n\n\n\nAzure Blob Storage `azure_blob_storage_account`, `azure_blob_storage_container`, `azure_blob_storage_access_key`, `azure_blob_storage_sas_token`, `azure_blob_storage_dns_suffix`\n\n\n\nAzure File Storage `azure_files_storage_account`, `azure_files_storage_access_key`, `azure_files_storage_share_name`, `azure_files_storage_dns_suffix`\n\n\n\nFilebase requires `filebase_bucket`, `filebase_access_key`, and `filebase_secret_key`.\n\n\n\nCloudflare requires `cloudflare_bucket`, `cloudflare_access_key`, `cloudflare_secret_key` and `cloudflare_endpoint`.\n\n\n\nLinode requires `linode_bucket`, `linode_access_key`, `linode_secret_key` and `linode_region`.",
		Attributes: map[string]schema.Attribute{
			"hostname": schema.StringAttribute{
				Description: "Hostname or IP address",
				Computed:    true,
				Optional:    true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
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
			"port": schema.Int64Attribute{
				Description: "Port for remote server.  Not needed for S3.",
				Computed:    true,
				Optional:    true,
				PlanModifiers: []planmodifier.Int64{
					int64planmodifier.UseStateForUnknown(),
				},
			},
			"buffer_uploads": schema.StringAttribute{
				Description: "If set to always, uploads to this server will be uploaded first to Files.com before being sent to the remote server. This can improve performance in certain access patterns, such as high-latency connections.  It will cause data to be temporarily stored in Files.com. If set to auto, we will perform this optimization if we believe it to be a benefit in a given situation.",
				Computed:    true,
				Optional:    true,
				Validators: []validator.String{
					stringvalidator.OneOf("auto", "always", "never"),
				},
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
			},
			"max_connections": schema.Int64Attribute{
				Description: "Max number of parallel connections.  Ignored for S3 connections (we will parallelize these as much as possible).",
				Computed:    true,
				Optional:    true,
				PlanModifiers: []planmodifier.Int64{
					int64planmodifier.UseStateForUnknown(),
				},
			},
			"pin_to_site_region": schema.BoolAttribute{
				Description: "If true, we will ensure that all communications with this remote server are made through the primary region of the site.  This setting can also be overridden by a site-wide setting which will force it to true.",
				Computed:    true,
				Optional:    true,
				PlanModifiers: []planmodifier.Bool{
					boolplanmodifier.UseStateForUnknown(),
				},
			},
			"s3_bucket": schema.StringAttribute{
				Description: "S3 bucket name",
				Computed:    true,
				Optional:    true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
			},
			"s3_region": schema.StringAttribute{
				Description: "S3 region",
				Computed:    true,
				Optional:    true,
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
			"server_certificate": schema.StringAttribute{
				Description: "Remote server certificate",
				Computed:    true,
				Optional:    true,
				Validators: []validator.String{
					stringvalidator.OneOf("require_match", "allow_any"),
				},
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
			},
			"server_host_key": schema.StringAttribute{
				Description: "Remote server SSH Host Key. If provided, we will require that the server host key matches the provided key. Uses OpenSSH format similar to what would go into ~/.ssh/known_hosts",
				Computed:    true,
				Optional:    true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
			},
			"server_type": schema.StringAttribute{
				Description: "Remote server type.",
				Computed:    true,
				Optional:    true,
				Validators: []validator.String{
					stringvalidator.OneOf("ftp", "sftp", "s3", "google_cloud_storage", "webdav", "wasabi", "backblaze_b2", "one_drive", "box", "dropbox", "google_drive", "azure", "sharepoint", "s3_compatible", "azure_files", "files_agent", "filebase", "cloudflare", "linode"),
				},
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
			},
			"ssl": schema.StringAttribute{
				Description: "Should we require SSL?",
				Computed:    true,
				Optional:    true,
				Validators: []validator.String{
					stringvalidator.OneOf("if_available", "require", "require_implicit", "never"),
				},
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
			},
			"username": schema.StringAttribute{
				Description: "Remote server username.  Not needed for S3 buckets.",
				Computed:    true,
				Optional:    true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
			},
			"google_cloud_storage_bucket": schema.StringAttribute{
				Description: "Google Cloud Storage: Bucket Name",
				Computed:    true,
				Optional:    true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
			},
			"google_cloud_storage_project_id": schema.StringAttribute{
				Description: "Google Cloud Storage: Project ID",
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
			"backblaze_b2_s3_endpoint": schema.StringAttribute{
				Description: "Backblaze B2 Cloud Storage: S3 Endpoint",
				Computed:    true,
				Optional:    true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
			},
			"backblaze_b2_bucket": schema.StringAttribute{
				Description: "Backblaze B2 Cloud Storage: Bucket name",
				Computed:    true,
				Optional:    true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
			},
			"wasabi_bucket": schema.StringAttribute{
				Description: "Wasabi: Bucket name",
				Computed:    true,
				Optional:    true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
			},
			"wasabi_region": schema.StringAttribute{
				Description: "Wasabi: Region",
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
			"one_drive_account_type": schema.StringAttribute{
				Description: "OneDrive: Either personal or business_other account types",
				Computed:    true,
				Optional:    true,
				Validators: []validator.String{
					stringvalidator.OneOf("personal", "business_other"),
				},
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
			},
			"azure_blob_storage_account": schema.StringAttribute{
				Description: "Azure Blob Storage: Account name",
				Computed:    true,
				Optional:    true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
			},
			"azure_blob_storage_container": schema.StringAttribute{
				Description: "Azure Blob Storage: Container name",
				Computed:    true,
				Optional:    true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
			},
			"azure_blob_storage_hierarchical_namespace": schema.BoolAttribute{
				Description: "Azure Blob Storage: Does the storage account has hierarchical namespace feature enabled?",
				Computed:    true,
				Optional:    true,
				PlanModifiers: []planmodifier.Bool{
					boolplanmodifier.UseStateForUnknown(),
				},
			},
			"azure_blob_storage_dns_suffix": schema.StringAttribute{
				Description: "Azure Blob Storage: Custom DNS suffix",
				Computed:    true,
				Optional:    true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
			},
			"azure_files_storage_account": schema.StringAttribute{
				Description: "Azure Files: Storage Account name",
				Computed:    true,
				Optional:    true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
			},
			"azure_files_storage_share_name": schema.StringAttribute{
				Description: "Azure Files:  Storage Share name",
				Computed:    true,
				Optional:    true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
			},
			"azure_files_storage_dns_suffix": schema.StringAttribute{
				Description: "Azure Files: Custom DNS suffix",
				Computed:    true,
				Optional:    true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
			},
			"s3_compatible_bucket": schema.StringAttribute{
				Description: "S3-compatible: Bucket name",
				Computed:    true,
				Optional:    true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
			},
			"s3_compatible_endpoint": schema.StringAttribute{
				Description: "S3-compatible: endpoint",
				Computed:    true,
				Optional:    true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
			},
			"s3_compatible_region": schema.StringAttribute{
				Description: "S3-compatible: region",
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
			"enable_dedicated_ips": schema.BoolAttribute{
				Description: "`true` if remote server only accepts connections from dedicated IPs",
				Computed:    true,
				Optional:    true,
				PlanModifiers: []planmodifier.Bool{
					boolplanmodifier.UseStateForUnknown(),
				},
			},
			"files_agent_permission_set": schema.StringAttribute{
				Description: "Local permissions for files agent. read_only, write_only, or read_write",
				Computed:    true,
				Optional:    true,
				Validators: []validator.String{
					stringvalidator.OneOf("read_write", "read_only", "write_only"),
				},
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
			},
			"files_agent_root": schema.StringAttribute{
				Description: "Agent local root path",
				Computed:    true,
				Optional:    true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
			},
			"files_agent_version": schema.StringAttribute{
				Description: "Files Agent version",
				Computed:    true,
				Optional:    true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
			},
			"filebase_bucket": schema.StringAttribute{
				Description: "Filebase: Bucket name",
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
			"cloudflare_bucket": schema.StringAttribute{
				Description: "Cloudflare: Bucket name",
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
			"cloudflare_endpoint": schema.StringAttribute{
				Description: "Cloudflare: endpoint",
				Computed:    true,
				Optional:    true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
			},
			"dropbox_teams": schema.BoolAttribute{
				Description: "Dropbox: If true, list Team folders in root?",
				Computed:    true,
				Optional:    true,
				PlanModifiers: []planmodifier.Bool{
					boolplanmodifier.UseStateForUnknown(),
				},
			},
			"linode_bucket": schema.StringAttribute{
				Description: "Linode: Bucket name",
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
			"linode_region": schema.StringAttribute{
				Description: "Linode: region",
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
			"reset_authentication": schema.BoolAttribute{
				Description: "Reset authenticated account?",
				Optional:    true,
				WriteOnly:   true,
			},
			"ssl_certificate": schema.StringAttribute{
				Description: "SSL client certificate.",
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
				Description: "Remote server ID",
				Computed:    true,
				PlanModifiers: []planmodifier.Int64{
					int64planmodifier.UseStateForUnknown(),
				},
			},
			"disabled": schema.BoolAttribute{
				Description: "If true, this server has been disabled due to failures.  Make any change or set disabled to false to clear this flag.",
				Computed:    true,
			},
			"authentication_method": schema.StringAttribute{
				Description: "Type of authentication method",
				Computed:    true,
			},
			"remote_home_path": schema.StringAttribute{
				Description: "Initial home folder on remote server",
				Computed:    true,
			},
			"pinned_region": schema.StringAttribute{
				Description: "If set, all communications with this remote server are made through the provided region.",
				Computed:    true,
			},
			"auth_status": schema.StringAttribute{
				Description: "Either `in_setup` or `complete`",
				Computed:    true,
				Validators: []validator.String{
					stringvalidator.OneOf("not_applicable", "in_setup", "complete", "reauthenticate"),
				},
			},
			"auth_account_name": schema.StringAttribute{
				Description: "Describes the authorized account",
				Computed:    true,
			},
			"files_agent_api_token": schema.StringAttribute{
				Description: "Files Agent API Token",
				Computed:    true,
			},
			"supports_versioning": schema.BoolAttribute{
				Description: "If true, this remote server supports file versioning. This value is determined automatically by Files.com.",
				Computed:    true,
			},
		},
	}
}

func (r *remoteServerResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	var plan remoteServerResourceModel
	diags := req.Plan.Get(ctx, &plan)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}
	var config remoteServerResourceModel
	diags = req.Config.Get(ctx, &config)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	paramsRemoteServerCreate := files_sdk.RemoteServerCreateParams{}
	paramsRemoteServerCreate.Password = config.Password.ValueString()
	paramsRemoteServerCreate.PrivateKey = config.PrivateKey.ValueString()
	paramsRemoteServerCreate.PrivateKeyPassphrase = config.PrivateKeyPassphrase.ValueString()
	if !config.ResetAuthentication.IsNull() && !config.ResetAuthentication.IsUnknown() {
		paramsRemoteServerCreate.ResetAuthentication = config.ResetAuthentication.ValueBoolPointer()
	}
	paramsRemoteServerCreate.SslCertificate = config.SslCertificate.ValueString()
	paramsRemoteServerCreate.AwsSecretKey = config.AwsSecretKey.ValueString()
	paramsRemoteServerCreate.AzureBlobStorageAccessKey = config.AzureBlobStorageAccessKey.ValueString()
	paramsRemoteServerCreate.AzureBlobStorageSasToken = config.AzureBlobStorageSasToken.ValueString()
	paramsRemoteServerCreate.AzureFilesStorageAccessKey = config.AzureFilesStorageAccessKey.ValueString()
	paramsRemoteServerCreate.AzureFilesStorageSasToken = config.AzureFilesStorageSasToken.ValueString()
	paramsRemoteServerCreate.BackblazeB2ApplicationKey = config.BackblazeB2ApplicationKey.ValueString()
	paramsRemoteServerCreate.BackblazeB2KeyId = config.BackblazeB2KeyId.ValueString()
	paramsRemoteServerCreate.CloudflareSecretKey = config.CloudflareSecretKey.ValueString()
	paramsRemoteServerCreate.FilebaseSecretKey = config.FilebaseSecretKey.ValueString()
	paramsRemoteServerCreate.GoogleCloudStorageCredentialsJson = config.GoogleCloudStorageCredentialsJson.ValueString()
	paramsRemoteServerCreate.GoogleCloudStorageS3CompatibleSecretKey = config.GoogleCloudStorageS3CompatibleSecretKey.ValueString()
	paramsRemoteServerCreate.LinodeSecretKey = config.LinodeSecretKey.ValueString()
	paramsRemoteServerCreate.S3CompatibleSecretKey = config.S3CompatibleSecretKey.ValueString()
	paramsRemoteServerCreate.WasabiSecretKey = config.WasabiSecretKey.ValueString()
	paramsRemoteServerCreate.AwsAccessKey = plan.AwsAccessKey.ValueString()
	paramsRemoteServerCreate.AzureBlobStorageAccount = plan.AzureBlobStorageAccount.ValueString()
	paramsRemoteServerCreate.AzureBlobStorageContainer = plan.AzureBlobStorageContainer.ValueString()
	paramsRemoteServerCreate.AzureBlobStorageDnsSuffix = plan.AzureBlobStorageDnsSuffix.ValueString()
	if !plan.AzureBlobStorageHierarchicalNamespace.IsNull() && !plan.AzureBlobStorageHierarchicalNamespace.IsUnknown() {
		paramsRemoteServerCreate.AzureBlobStorageHierarchicalNamespace = plan.AzureBlobStorageHierarchicalNamespace.ValueBoolPointer()
	}
	paramsRemoteServerCreate.AzureFilesStorageAccount = plan.AzureFilesStorageAccount.ValueString()
	paramsRemoteServerCreate.AzureFilesStorageDnsSuffix = plan.AzureFilesStorageDnsSuffix.ValueString()
	paramsRemoteServerCreate.AzureFilesStorageShareName = plan.AzureFilesStorageShareName.ValueString()
	paramsRemoteServerCreate.BackblazeB2Bucket = plan.BackblazeB2Bucket.ValueString()
	paramsRemoteServerCreate.BackblazeB2S3Endpoint = plan.BackblazeB2S3Endpoint.ValueString()
	paramsRemoteServerCreate.BufferUploads = paramsRemoteServerCreate.BufferUploads.Enum()[plan.BufferUploads.ValueString()]
	paramsRemoteServerCreate.CloudflareAccessKey = plan.CloudflareAccessKey.ValueString()
	paramsRemoteServerCreate.CloudflareBucket = plan.CloudflareBucket.ValueString()
	paramsRemoteServerCreate.CloudflareEndpoint = plan.CloudflareEndpoint.ValueString()
	if !plan.DropboxTeams.IsNull() && !plan.DropboxTeams.IsUnknown() {
		paramsRemoteServerCreate.DropboxTeams = plan.DropboxTeams.ValueBoolPointer()
	}
	if !plan.EnableDedicatedIps.IsNull() && !plan.EnableDedicatedIps.IsUnknown() {
		paramsRemoteServerCreate.EnableDedicatedIps = plan.EnableDedicatedIps.ValueBoolPointer()
	}
	paramsRemoteServerCreate.FilebaseAccessKey = plan.FilebaseAccessKey.ValueString()
	paramsRemoteServerCreate.FilebaseBucket = plan.FilebaseBucket.ValueString()
	paramsRemoteServerCreate.FilesAgentPermissionSet = paramsRemoteServerCreate.FilesAgentPermissionSet.Enum()[plan.FilesAgentPermissionSet.ValueString()]
	paramsRemoteServerCreate.FilesAgentRoot = plan.FilesAgentRoot.ValueString()
	paramsRemoteServerCreate.FilesAgentVersion = plan.FilesAgentVersion.ValueString()
	paramsRemoteServerCreate.GoogleCloudStorageBucket = plan.GoogleCloudStorageBucket.ValueString()
	paramsRemoteServerCreate.GoogleCloudStorageProjectId = plan.GoogleCloudStorageProjectId.ValueString()
	paramsRemoteServerCreate.GoogleCloudStorageS3CompatibleAccessKey = plan.GoogleCloudStorageS3CompatibleAccessKey.ValueString()
	paramsRemoteServerCreate.Hostname = plan.Hostname.ValueString()
	paramsRemoteServerCreate.LinodeAccessKey = plan.LinodeAccessKey.ValueString()
	paramsRemoteServerCreate.LinodeBucket = plan.LinodeBucket.ValueString()
	paramsRemoteServerCreate.LinodeRegion = plan.LinodeRegion.ValueString()
	paramsRemoteServerCreate.MaxConnections = plan.MaxConnections.ValueInt64()
	paramsRemoteServerCreate.Name = plan.Name.ValueString()
	paramsRemoteServerCreate.OneDriveAccountType = paramsRemoteServerCreate.OneDriveAccountType.Enum()[plan.OneDriveAccountType.ValueString()]
	if !plan.PinToSiteRegion.IsNull() && !plan.PinToSiteRegion.IsUnknown() {
		paramsRemoteServerCreate.PinToSiteRegion = plan.PinToSiteRegion.ValueBoolPointer()
	}
	paramsRemoteServerCreate.Port = plan.Port.ValueInt64()
	paramsRemoteServerCreate.S3Bucket = plan.S3Bucket.ValueString()
	paramsRemoteServerCreate.S3CompatibleAccessKey = plan.S3CompatibleAccessKey.ValueString()
	paramsRemoteServerCreate.S3CompatibleBucket = plan.S3CompatibleBucket.ValueString()
	paramsRemoteServerCreate.S3CompatibleEndpoint = plan.S3CompatibleEndpoint.ValueString()
	paramsRemoteServerCreate.S3CompatibleRegion = plan.S3CompatibleRegion.ValueString()
	paramsRemoteServerCreate.S3Region = plan.S3Region.ValueString()
	paramsRemoteServerCreate.ServerCertificate = paramsRemoteServerCreate.ServerCertificate.Enum()[plan.ServerCertificate.ValueString()]
	paramsRemoteServerCreate.ServerHostKey = plan.ServerHostKey.ValueString()
	paramsRemoteServerCreate.ServerType = paramsRemoteServerCreate.ServerType.Enum()[plan.ServerType.ValueString()]
	paramsRemoteServerCreate.Ssl = paramsRemoteServerCreate.Ssl.Enum()[plan.Ssl.ValueString()]
	paramsRemoteServerCreate.Username = plan.Username.ValueString()
	paramsRemoteServerCreate.WasabiAccessKey = plan.WasabiAccessKey.ValueString()
	paramsRemoteServerCreate.WasabiBucket = plan.WasabiBucket.ValueString()
	paramsRemoteServerCreate.WasabiRegion = plan.WasabiRegion.ValueString()

	if resp.Diagnostics.HasError() {
		return
	}

	remoteServer, err := r.client.Create(paramsRemoteServerCreate, files_sdk.WithContext(ctx))
	if err != nil {
		resp.Diagnostics.AddError(
			"Error Creating Files RemoteServer",
			"Could not create remote_server, unexpected error: "+err.Error(),
		)
		return
	}

	diags = r.populateResourceModel(ctx, remoteServer, &plan)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	diags = resp.State.Set(ctx, plan)
	resp.Diagnostics.Append(diags...)
}

func (r *remoteServerResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	var state remoteServerResourceModel
	diags := req.State.Get(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	paramsRemoteServerFind := files_sdk.RemoteServerFindParams{}
	paramsRemoteServerFind.Id = state.Id.ValueInt64()

	remoteServer, err := r.client.Find(paramsRemoteServerFind, files_sdk.WithContext(ctx))
	if err != nil {
		if files_sdk.IsNotExist(err) {
			resp.State.RemoveResource(ctx)
			return
		}

		resp.Diagnostics.AddError(
			"Error Reading Files RemoteServer",
			"Could not read remote_server id "+fmt.Sprint(state.Id.ValueInt64())+": "+err.Error(),
		)
		return
	}

	diags = r.populateResourceModel(ctx, remoteServer, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	diags = resp.State.Set(ctx, state)
	resp.Diagnostics.Append(diags...)
}

func (r *remoteServerResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	var plan remoteServerResourceModel
	diags := req.Plan.Get(ctx, &plan)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}
	var config remoteServerResourceModel
	diags = req.Config.Get(ctx, &config)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	paramsRemoteServerUpdate := files_sdk.RemoteServerUpdateParams{}
	paramsRemoteServerUpdate.Id = plan.Id.ValueInt64()
	paramsRemoteServerUpdate.Password = config.Password.ValueString()
	paramsRemoteServerUpdate.PrivateKey = config.PrivateKey.ValueString()
	paramsRemoteServerUpdate.PrivateKeyPassphrase = config.PrivateKeyPassphrase.ValueString()
	if !config.ResetAuthentication.IsNull() && !config.ResetAuthentication.IsUnknown() {
		paramsRemoteServerUpdate.ResetAuthentication = config.ResetAuthentication.ValueBoolPointer()
	}
	paramsRemoteServerUpdate.SslCertificate = config.SslCertificate.ValueString()
	paramsRemoteServerUpdate.AwsSecretKey = config.AwsSecretKey.ValueString()
	paramsRemoteServerUpdate.AzureBlobStorageAccessKey = config.AzureBlobStorageAccessKey.ValueString()
	paramsRemoteServerUpdate.AzureBlobStorageSasToken = config.AzureBlobStorageSasToken.ValueString()
	paramsRemoteServerUpdate.AzureFilesStorageAccessKey = config.AzureFilesStorageAccessKey.ValueString()
	paramsRemoteServerUpdate.AzureFilesStorageSasToken = config.AzureFilesStorageSasToken.ValueString()
	paramsRemoteServerUpdate.BackblazeB2ApplicationKey = config.BackblazeB2ApplicationKey.ValueString()
	paramsRemoteServerUpdate.BackblazeB2KeyId = config.BackblazeB2KeyId.ValueString()
	paramsRemoteServerUpdate.CloudflareSecretKey = config.CloudflareSecretKey.ValueString()
	paramsRemoteServerUpdate.FilebaseSecretKey = config.FilebaseSecretKey.ValueString()
	paramsRemoteServerUpdate.GoogleCloudStorageCredentialsJson = config.GoogleCloudStorageCredentialsJson.ValueString()
	paramsRemoteServerUpdate.GoogleCloudStorageS3CompatibleSecretKey = config.GoogleCloudStorageS3CompatibleSecretKey.ValueString()
	paramsRemoteServerUpdate.LinodeSecretKey = config.LinodeSecretKey.ValueString()
	paramsRemoteServerUpdate.S3CompatibleSecretKey = config.S3CompatibleSecretKey.ValueString()
	paramsRemoteServerUpdate.WasabiSecretKey = config.WasabiSecretKey.ValueString()
	paramsRemoteServerUpdate.AwsAccessKey = plan.AwsAccessKey.ValueString()
	paramsRemoteServerUpdate.AzureBlobStorageAccount = plan.AzureBlobStorageAccount.ValueString()
	paramsRemoteServerUpdate.AzureBlobStorageContainer = plan.AzureBlobStorageContainer.ValueString()
	paramsRemoteServerUpdate.AzureBlobStorageDnsSuffix = plan.AzureBlobStorageDnsSuffix.ValueString()
	if !plan.AzureBlobStorageHierarchicalNamespace.IsNull() && !plan.AzureBlobStorageHierarchicalNamespace.IsUnknown() {
		paramsRemoteServerUpdate.AzureBlobStorageHierarchicalNamespace = plan.AzureBlobStorageHierarchicalNamespace.ValueBoolPointer()
	}
	paramsRemoteServerUpdate.AzureFilesStorageAccount = plan.AzureFilesStorageAccount.ValueString()
	paramsRemoteServerUpdate.AzureFilesStorageDnsSuffix = plan.AzureFilesStorageDnsSuffix.ValueString()
	paramsRemoteServerUpdate.AzureFilesStorageShareName = plan.AzureFilesStorageShareName.ValueString()
	paramsRemoteServerUpdate.BackblazeB2Bucket = plan.BackblazeB2Bucket.ValueString()
	paramsRemoteServerUpdate.BackblazeB2S3Endpoint = plan.BackblazeB2S3Endpoint.ValueString()
	paramsRemoteServerUpdate.BufferUploads = paramsRemoteServerUpdate.BufferUploads.Enum()[plan.BufferUploads.ValueString()]
	paramsRemoteServerUpdate.CloudflareAccessKey = plan.CloudflareAccessKey.ValueString()
	paramsRemoteServerUpdate.CloudflareBucket = plan.CloudflareBucket.ValueString()
	paramsRemoteServerUpdate.CloudflareEndpoint = plan.CloudflareEndpoint.ValueString()
	if !plan.DropboxTeams.IsNull() && !plan.DropboxTeams.IsUnknown() {
		paramsRemoteServerUpdate.DropboxTeams = plan.DropboxTeams.ValueBoolPointer()
	}
	if !plan.EnableDedicatedIps.IsNull() && !plan.EnableDedicatedIps.IsUnknown() {
		paramsRemoteServerUpdate.EnableDedicatedIps = plan.EnableDedicatedIps.ValueBoolPointer()
	}
	paramsRemoteServerUpdate.FilebaseAccessKey = plan.FilebaseAccessKey.ValueString()
	paramsRemoteServerUpdate.FilebaseBucket = plan.FilebaseBucket.ValueString()
	paramsRemoteServerUpdate.FilesAgentPermissionSet = paramsRemoteServerUpdate.FilesAgentPermissionSet.Enum()[plan.FilesAgentPermissionSet.ValueString()]
	paramsRemoteServerUpdate.FilesAgentRoot = plan.FilesAgentRoot.ValueString()
	paramsRemoteServerUpdate.FilesAgentVersion = plan.FilesAgentVersion.ValueString()
	paramsRemoteServerUpdate.GoogleCloudStorageBucket = plan.GoogleCloudStorageBucket.ValueString()
	paramsRemoteServerUpdate.GoogleCloudStorageProjectId = plan.GoogleCloudStorageProjectId.ValueString()
	paramsRemoteServerUpdate.GoogleCloudStorageS3CompatibleAccessKey = plan.GoogleCloudStorageS3CompatibleAccessKey.ValueString()
	paramsRemoteServerUpdate.Hostname = plan.Hostname.ValueString()
	paramsRemoteServerUpdate.LinodeAccessKey = plan.LinodeAccessKey.ValueString()
	paramsRemoteServerUpdate.LinodeBucket = plan.LinodeBucket.ValueString()
	paramsRemoteServerUpdate.LinodeRegion = plan.LinodeRegion.ValueString()
	paramsRemoteServerUpdate.MaxConnections = plan.MaxConnections.ValueInt64()
	paramsRemoteServerUpdate.Name = plan.Name.ValueString()
	paramsRemoteServerUpdate.OneDriveAccountType = paramsRemoteServerUpdate.OneDriveAccountType.Enum()[plan.OneDriveAccountType.ValueString()]
	if !plan.PinToSiteRegion.IsNull() && !plan.PinToSiteRegion.IsUnknown() {
		paramsRemoteServerUpdate.PinToSiteRegion = plan.PinToSiteRegion.ValueBoolPointer()
	}
	paramsRemoteServerUpdate.Port = plan.Port.ValueInt64()
	paramsRemoteServerUpdate.S3Bucket = plan.S3Bucket.ValueString()
	paramsRemoteServerUpdate.S3CompatibleAccessKey = plan.S3CompatibleAccessKey.ValueString()
	paramsRemoteServerUpdate.S3CompatibleBucket = plan.S3CompatibleBucket.ValueString()
	paramsRemoteServerUpdate.S3CompatibleEndpoint = plan.S3CompatibleEndpoint.ValueString()
	paramsRemoteServerUpdate.S3CompatibleRegion = plan.S3CompatibleRegion.ValueString()
	paramsRemoteServerUpdate.S3Region = plan.S3Region.ValueString()
	paramsRemoteServerUpdate.ServerCertificate = paramsRemoteServerUpdate.ServerCertificate.Enum()[plan.ServerCertificate.ValueString()]
	paramsRemoteServerUpdate.ServerHostKey = plan.ServerHostKey.ValueString()
	paramsRemoteServerUpdate.ServerType = paramsRemoteServerUpdate.ServerType.Enum()[plan.ServerType.ValueString()]
	paramsRemoteServerUpdate.Ssl = paramsRemoteServerUpdate.Ssl.Enum()[plan.Ssl.ValueString()]
	paramsRemoteServerUpdate.Username = plan.Username.ValueString()
	paramsRemoteServerUpdate.WasabiAccessKey = plan.WasabiAccessKey.ValueString()
	paramsRemoteServerUpdate.WasabiBucket = plan.WasabiBucket.ValueString()
	paramsRemoteServerUpdate.WasabiRegion = plan.WasabiRegion.ValueString()

	if resp.Diagnostics.HasError() {
		return
	}

	remoteServer, err := r.client.Update(paramsRemoteServerUpdate, files_sdk.WithContext(ctx))
	if err != nil {
		resp.Diagnostics.AddError(
			"Error Updating Files RemoteServer",
			"Could not update remote_server, unexpected error: "+err.Error(),
		)
		return
	}

	diags = r.populateResourceModel(ctx, remoteServer, &plan)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	diags = resp.State.Set(ctx, plan)
	resp.Diagnostics.Append(diags...)
}

func (r *remoteServerResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	var state remoteServerResourceModel
	diags := req.State.Get(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	paramsRemoteServerDelete := files_sdk.RemoteServerDeleteParams{}
	paramsRemoteServerDelete.Id = state.Id.ValueInt64()

	err := r.client.Delete(paramsRemoteServerDelete, files_sdk.WithContext(ctx))
	if err != nil && !files_sdk.IsNotExist(err) {
		resp.Diagnostics.AddError(
			"Error Deleting Files RemoteServer",
			"Could not delete remote_server id "+fmt.Sprint(state.Id.ValueInt64())+": "+err.Error(),
		)
	}
}

func (r *remoteServerResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
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

func (r *remoteServerResource) populateResourceModel(ctx context.Context, remoteServer files_sdk.RemoteServer, state *remoteServerResourceModel) (diags diag.Diagnostics) {
	state.Id = types.Int64Value(remoteServer.Id)
	state.Disabled = types.BoolPointerValue(remoteServer.Disabled)
	state.AuthenticationMethod = types.StringValue(remoteServer.AuthenticationMethod)
	state.Hostname = types.StringValue(remoteServer.Hostname)
	state.RemoteHomePath = types.StringValue(remoteServer.RemoteHomePath)
	state.Name = types.StringValue(remoteServer.Name)
	state.Port = types.Int64Value(remoteServer.Port)
	state.BufferUploads = types.StringValue(remoteServer.BufferUploads)
	state.MaxConnections = types.Int64Value(remoteServer.MaxConnections)
	state.PinToSiteRegion = types.BoolPointerValue(remoteServer.PinToSiteRegion)
	state.PinnedRegion = types.StringValue(remoteServer.PinnedRegion)
	state.S3Bucket = types.StringValue(remoteServer.S3Bucket)
	state.S3Region = types.StringValue(remoteServer.S3Region)
	state.AwsAccessKey = types.StringValue(remoteServer.AwsAccessKey)
	state.ServerCertificate = types.StringValue(remoteServer.ServerCertificate)
	state.ServerHostKey = types.StringValue(remoteServer.ServerHostKey)
	state.ServerType = types.StringValue(remoteServer.ServerType)
	state.Ssl = types.StringValue(remoteServer.Ssl)
	state.Username = types.StringValue(remoteServer.Username)
	state.GoogleCloudStorageBucket = types.StringValue(remoteServer.GoogleCloudStorageBucket)
	state.GoogleCloudStorageProjectId = types.StringValue(remoteServer.GoogleCloudStorageProjectId)
	state.GoogleCloudStorageS3CompatibleAccessKey = types.StringValue(remoteServer.GoogleCloudStorageS3CompatibleAccessKey)
	state.BackblazeB2S3Endpoint = types.StringValue(remoteServer.BackblazeB2S3Endpoint)
	state.BackblazeB2Bucket = types.StringValue(remoteServer.BackblazeB2Bucket)
	state.WasabiBucket = types.StringValue(remoteServer.WasabiBucket)
	state.WasabiRegion = types.StringValue(remoteServer.WasabiRegion)
	state.WasabiAccessKey = types.StringValue(remoteServer.WasabiAccessKey)
	state.AuthStatus = types.StringValue(remoteServer.AuthStatus)
	state.AuthAccountName = types.StringValue(remoteServer.AuthAccountName)
	state.OneDriveAccountType = types.StringValue(remoteServer.OneDriveAccountType)
	state.AzureBlobStorageAccount = types.StringValue(remoteServer.AzureBlobStorageAccount)
	state.AzureBlobStorageContainer = types.StringValue(remoteServer.AzureBlobStorageContainer)
	state.AzureBlobStorageHierarchicalNamespace = types.BoolPointerValue(remoteServer.AzureBlobStorageHierarchicalNamespace)
	state.AzureBlobStorageDnsSuffix = types.StringValue(remoteServer.AzureBlobStorageDnsSuffix)
	state.AzureFilesStorageAccount = types.StringValue(remoteServer.AzureFilesStorageAccount)
	state.AzureFilesStorageShareName = types.StringValue(remoteServer.AzureFilesStorageShareName)
	state.AzureFilesStorageDnsSuffix = types.StringValue(remoteServer.AzureFilesStorageDnsSuffix)
	state.S3CompatibleBucket = types.StringValue(remoteServer.S3CompatibleBucket)
	state.S3CompatibleEndpoint = types.StringValue(remoteServer.S3CompatibleEndpoint)
	state.S3CompatibleRegion = types.StringValue(remoteServer.S3CompatibleRegion)
	state.S3CompatibleAccessKey = types.StringValue(remoteServer.S3CompatibleAccessKey)
	state.EnableDedicatedIps = types.BoolPointerValue(remoteServer.EnableDedicatedIps)
	state.FilesAgentPermissionSet = types.StringValue(remoteServer.FilesAgentPermissionSet)
	state.FilesAgentRoot = types.StringValue(remoteServer.FilesAgentRoot)
	state.FilesAgentApiToken = types.StringValue(remoteServer.FilesAgentApiToken)
	state.FilesAgentVersion = types.StringValue(remoteServer.FilesAgentVersion)
	state.FilebaseBucket = types.StringValue(remoteServer.FilebaseBucket)
	state.FilebaseAccessKey = types.StringValue(remoteServer.FilebaseAccessKey)
	state.CloudflareBucket = types.StringValue(remoteServer.CloudflareBucket)
	state.CloudflareAccessKey = types.StringValue(remoteServer.CloudflareAccessKey)
	state.CloudflareEndpoint = types.StringValue(remoteServer.CloudflareEndpoint)
	state.DropboxTeams = types.BoolPointerValue(remoteServer.DropboxTeams)
	state.LinodeBucket = types.StringValue(remoteServer.LinodeBucket)
	state.LinodeAccessKey = types.StringValue(remoteServer.LinodeAccessKey)
	state.LinodeRegion = types.StringValue(remoteServer.LinodeRegion)
	state.SupportsVersioning = types.BoolPointerValue(remoteServer.SupportsVersioning)

	return
}
