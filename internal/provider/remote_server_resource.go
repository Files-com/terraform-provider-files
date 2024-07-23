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
	Hostname                              types.String `tfsdk:"hostname"`
	Name                                  types.String `tfsdk:"name"`
	Port                                  types.Int64  `tfsdk:"port"`
	MaxConnections                        types.Int64  `tfsdk:"max_connections"`
	PinToSiteRegion                       types.Bool   `tfsdk:"pin_to_site_region"`
	S3Bucket                              types.String `tfsdk:"s3_bucket"`
	S3Region                              types.String `tfsdk:"s3_region"`
	AwsAccessKey                          types.String `tfsdk:"aws_access_key"`
	ServerCertificate                     types.String `tfsdk:"server_certificate"`
	ServerHostKey                         types.String `tfsdk:"server_host_key"`
	ServerType                            types.String `tfsdk:"server_type"`
	Ssl                                   types.String `tfsdk:"ssl"`
	Username                              types.String `tfsdk:"username"`
	GoogleCloudStorageBucket              types.String `tfsdk:"google_cloud_storage_bucket"`
	GoogleCloudStorageProjectId           types.String `tfsdk:"google_cloud_storage_project_id"`
	BackblazeB2S3Endpoint                 types.String `tfsdk:"backblaze_b2_s3_endpoint"`
	BackblazeB2Bucket                     types.String `tfsdk:"backblaze_b2_bucket"`
	WasabiBucket                          types.String `tfsdk:"wasabi_bucket"`
	WasabiRegion                          types.String `tfsdk:"wasabi_region"`
	WasabiAccessKey                       types.String `tfsdk:"wasabi_access_key"`
	RackspaceUsername                     types.String `tfsdk:"rackspace_username"`
	RackspaceRegion                       types.String `tfsdk:"rackspace_region"`
	RackspaceContainer                    types.String `tfsdk:"rackspace_container"`
	OneDriveAccountType                   types.String `tfsdk:"one_drive_account_type"`
	AzureBlobStorageAccount               types.String `tfsdk:"azure_blob_storage_account"`
	AzureBlobStorageContainer             types.String `tfsdk:"azure_blob_storage_container"`
	AzureBlobStorageHierarchicalNamespace types.Bool   `tfsdk:"azure_blob_storage_hierarchical_namespace"`
	AzureFilesStorageAccount              types.String `tfsdk:"azure_files_storage_account"`
	AzureFilesStorageShareName            types.String `tfsdk:"azure_files_storage_share_name"`
	S3CompatibleBucket                    types.String `tfsdk:"s3_compatible_bucket"`
	S3CompatibleEndpoint                  types.String `tfsdk:"s3_compatible_endpoint"`
	S3CompatibleRegion                    types.String `tfsdk:"s3_compatible_region"`
	S3CompatibleAccessKey                 types.String `tfsdk:"s3_compatible_access_key"`
	EnableDedicatedIps                    types.Bool   `tfsdk:"enable_dedicated_ips"`
	FilesAgentPermissionSet               types.String `tfsdk:"files_agent_permission_set"`
	FilesAgentRoot                        types.String `tfsdk:"files_agent_root"`
	FilesAgentVersion                     types.String `tfsdk:"files_agent_version"`
	FilebaseBucket                        types.String `tfsdk:"filebase_bucket"`
	FilebaseAccessKey                     types.String `tfsdk:"filebase_access_key"`
	CloudflareBucket                      types.String `tfsdk:"cloudflare_bucket"`
	CloudflareAccessKey                   types.String `tfsdk:"cloudflare_access_key"`
	CloudflareEndpoint                    types.String `tfsdk:"cloudflare_endpoint"`
	DropboxTeams                          types.Bool   `tfsdk:"dropbox_teams"`
	LinodeBucket                          types.String `tfsdk:"linode_bucket"`
	LinodeAccessKey                       types.String `tfsdk:"linode_access_key"`
	LinodeRegion                          types.String `tfsdk:"linode_region"`
	AwsSecretKey                          types.String `tfsdk:"aws_secret_key"`
	Password                              types.String `tfsdk:"password"`
	PrivateKey                            types.String `tfsdk:"private_key"`
	PrivateKeyPassphrase                  types.String `tfsdk:"private_key_passphrase"`
	SslCertificate                        types.String `tfsdk:"ssl_certificate"`
	GoogleCloudStorageCredentialsJson     types.String `tfsdk:"google_cloud_storage_credentials_json"`
	WasabiSecretKey                       types.String `tfsdk:"wasabi_secret_key"`
	BackblazeB2KeyId                      types.String `tfsdk:"backblaze_b2_key_id"`
	BackblazeB2ApplicationKey             types.String `tfsdk:"backblaze_b2_application_key"`
	RackspaceApiKey                       types.String `tfsdk:"rackspace_api_key"`
	ResetAuthentication                   types.Bool   `tfsdk:"reset_authentication"`
	AzureBlobStorageAccessKey             types.String `tfsdk:"azure_blob_storage_access_key"`
	AzureFilesStorageAccessKey            types.String `tfsdk:"azure_files_storage_access_key"`
	AzureBlobStorageSasToken              types.String `tfsdk:"azure_blob_storage_sas_token"`
	AzureFilesStorageSasToken             types.String `tfsdk:"azure_files_storage_sas_token"`
	S3CompatibleSecretKey                 types.String `tfsdk:"s3_compatible_secret_key"`
	FilebaseSecretKey                     types.String `tfsdk:"filebase_secret_key"`
	CloudflareSecretKey                   types.String `tfsdk:"cloudflare_secret_key"`
	LinodeSecretKey                       types.String `tfsdk:"linode_secret_key"`
	Id                                    types.Int64  `tfsdk:"id"`
	Disabled                              types.Bool   `tfsdk:"disabled"`
	AuthenticationMethod                  types.String `tfsdk:"authentication_method"`
	RemoteHomePath                        types.String `tfsdk:"remote_home_path"`
	PinnedRegion                          types.String `tfsdk:"pinned_region"`
	AuthSetupLink                         types.String `tfsdk:"auth_setup_link"`
	AuthStatus                            types.String `tfsdk:"auth_status"`
	AuthAccountName                       types.String `tfsdk:"auth_account_name"`
	FilesAgentApiToken                    types.String `tfsdk:"files_agent_api_token"`
	SupportsVersioning                    types.Bool   `tfsdk:"supports_versioning"`
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
		Description: "Remote servers are used with the `remote_server_sync` Behavior.\n\n\n\nRemote Servers can be either an FTP server, SFTP server, S3 bucket, Google Cloud Storage, Wasabi, Backblaze B2 Cloud Storage, Rackspace Cloud Files container, WebDAV, Box, Dropbox, OneDrive, Google Drive, or Azure Blob Storage.\n\n\n\nNot every attribute will apply to every remote server.\n\n\n\nFTP Servers require that you specify their `hostname`, `port`, `username`, `password`, and a value for `ssl`. Optionally, provide `server_certificate`.\n\n\n\nSFTP Servers require that you specify their `hostname`, `port`, `username`, `password` or `private_key`, and a value for `ssl`. Optionally, provide `server_certificate`, `private_key_passphrase`.\n\n\n\nS3 Buckets require that you specify their `s3_bucket` name, and `s3_region`. Optionally provide a `aws_access_key`, and `aws_secret_key`. If you don't provide credentials, you will need to use AWS to grant us access to your bucket.\n\n\n\nS3-Compatible Buckets require that you specify `s3_compatible_bucket`, `s3_compatible_endpoint`, `s3_compatible_access_key`, and `s3_compatible_secret_key`.\n\n\n\nGoogle Cloud Storage requires that you specify `google_cloud_storage_bucket`, `google_cloud_storage_project_id`, and `google_cloud_storage_credentials_json`.\n\n\n\nWasabi requires `wasabi_bucket`, `wasabi_region`, `wasabi_access_key`, and `wasabi_secret_key`.\n\n\n\nBackblaze B2 Cloud Storage `backblaze_b2_bucket`, `backblaze_b2_s3_endpoint`, `backblaze_b2_application_key`, and `backblaze_b2_key_id`. (Requires S3 Compatible API) See https://help.backblaze.com/hc/en-us/articles/360047425453\n\n\n\nRackspace Cloud Files requires `rackspace_username`, `rackspace_api_key`, `rackspace_region`, and `rackspace_container`.\n\n\n\nWebDAV Servers require that you specify their `hostname`, `username`, and `password`.\n\n\n\nOneDrive follow the `auth_setup_link` and login with Microsoft.\n\n\n\nSharepoint follow the `auth_setup_link` and login with Microsoft.\n\n\n\nBox follow the `auth_setup_link` and login with Box.\n\n\n\nDropbox specify if `dropbox_teams` then follow the `auth_setup_link` and login with Dropbox.\n\n\n\nGoogle Drive follow the `auth_setup_link` and login with Google.\n\n\n\nAzure Blob Storage `azure_blob_storage_account`, `azure_blob_storage_container`, `azure_blob_storage_access_key`, `azure_blob_storage_sas_token`\n\n\n\nAzure File Storage `azure_files_storage_account`, `azure_files_storage_access_key`, `azure_files_storage_share_name`\n\n\n\nFilebase requires `filebase_bucket`, `filebase_access_key`, and `filebase_secret_key`.\n\n\n\nCloudflare requires `cloudflare_bucket`, `cloudflare_access_key`, `cloudflare_secret_key` and `cloudflare_endpoint`.\n\n\n\nLinode requires `linode_bucket`, `linode_access_key`, `linode_secret_key` and `linode_region`.",
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
			"max_connections": schema.Int64Attribute{
				Description: "Max number of parallel connections.  Ignored for S3 connections (we will parallelize these as much as possible).",
				Computed:    true,
				Optional:    true,
				PlanModifiers: []planmodifier.Int64{
					int64planmodifier.UseStateForUnknown(),
				},
			},
			"pin_to_site_region": schema.BoolAttribute{
				Description: "If true, we will ensure that all communications with this remote server are made through the primary region of the site.  This setting can also be overridden by a sitewide setting which will force it to true.",
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
					stringvalidator.OneOf("ftp", "sftp", "s3", "google_cloud_storage", "webdav", "wasabi", "backblaze_b2", "one_drive", "rackspace", "box", "dropbox", "google_drive", "azure", "sharepoint", "s3_compatible", "azure_files", "files_agent", "filebase", "cloudflare", "linode"),
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
				Description: "Google Cloud Storage bucket name",
				Computed:    true,
				Optional:    true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
			},
			"google_cloud_storage_project_id": schema.StringAttribute{
				Description: "Google Cloud Project ID",
				Computed:    true,
				Optional:    true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
			},
			"backblaze_b2_s3_endpoint": schema.StringAttribute{
				Description: "Backblaze B2 Cloud Storage S3 Endpoint",
				Computed:    true,
				Optional:    true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
			},
			"backblaze_b2_bucket": schema.StringAttribute{
				Description: "Backblaze B2 Cloud Storage Bucket name",
				Computed:    true,
				Optional:    true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
			},
			"wasabi_bucket": schema.StringAttribute{
				Description: "Wasabi Bucket name",
				Computed:    true,
				Optional:    true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
			},
			"wasabi_region": schema.StringAttribute{
				Description: "Wasabi region",
				Computed:    true,
				Optional:    true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
			},
			"wasabi_access_key": schema.StringAttribute{
				Description: "Wasabi access key.",
				Computed:    true,
				Optional:    true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
			},
			"rackspace_username": schema.StringAttribute{
				Description: "Rackspace username used to login to the Rackspace Cloud Control Panel.",
				Computed:    true,
				Optional:    true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
			},
			"rackspace_region": schema.StringAttribute{
				Description: "Three letter airport code for Rackspace region. See https://support.rackspace.com/how-to/about-regions/",
				Computed:    true,
				Optional:    true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
			},
			"rackspace_container": schema.StringAttribute{
				Description: "The name of the container (top level directory) where files will sync.",
				Computed:    true,
				Optional:    true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
			},
			"one_drive_account_type": schema.StringAttribute{
				Description: "Either personal or business_other account types",
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
				Description: "Azure Blob Storage Account name",
				Computed:    true,
				Optional:    true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
			},
			"azure_blob_storage_container": schema.StringAttribute{
				Description: "Azure Blob Storage Container name",
				Computed:    true,
				Optional:    true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
			},
			"azure_blob_storage_hierarchical_namespace": schema.BoolAttribute{
				Description: "Enable when storage account has hierarchical namespace feature enabled",
				Computed:    true,
				Optional:    true,
				PlanModifiers: []planmodifier.Bool{
					boolplanmodifier.UseStateForUnknown(),
				},
			},
			"azure_files_storage_account": schema.StringAttribute{
				Description: "Azure File Storage Account name",
				Computed:    true,
				Optional:    true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
			},
			"azure_files_storage_share_name": schema.StringAttribute{
				Description: "Azure File Storage Share name",
				Computed:    true,
				Optional:    true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
			},
			"s3_compatible_bucket": schema.StringAttribute{
				Description: "S3-compatible Bucket name",
				Computed:    true,
				Optional:    true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
			},
			"s3_compatible_endpoint": schema.StringAttribute{
				Description: "S3-compatible endpoint",
				Computed:    true,
				Optional:    true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
			},
			"s3_compatible_region": schema.StringAttribute{
				Description: "S3-compatible endpoint",
				Computed:    true,
				Optional:    true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
			},
			"s3_compatible_access_key": schema.StringAttribute{
				Description: "S3-compatible Access Key.",
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
				Description: "Filebase Bucket name",
				Computed:    true,
				Optional:    true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
			},
			"filebase_access_key": schema.StringAttribute{
				Description: "Filebase Access Key.",
				Computed:    true,
				Optional:    true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
			},
			"cloudflare_bucket": schema.StringAttribute{
				Description: "Cloudflare Bucket name",
				Computed:    true,
				Optional:    true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
			},
			"cloudflare_access_key": schema.StringAttribute{
				Description: "Cloudflare Access Key.",
				Computed:    true,
				Optional:    true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
			},
			"cloudflare_endpoint": schema.StringAttribute{
				Description: "Cloudflare endpoint",
				Computed:    true,
				Optional:    true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
			},
			"dropbox_teams": schema.BoolAttribute{
				Description: "List Team folders in root",
				Computed:    true,
				Optional:    true,
				PlanModifiers: []planmodifier.Bool{
					boolplanmodifier.UseStateForUnknown(),
				},
			},
			"linode_bucket": schema.StringAttribute{
				Description: "Linode Bucket name",
				Computed:    true,
				Optional:    true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
			},
			"linode_access_key": schema.StringAttribute{
				Description: "Linode Access Key.",
				Computed:    true,
				Optional:    true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
			},
			"linode_region": schema.StringAttribute{
				Description: "Linode region",
				Computed:    true,
				Optional:    true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
			},
			"aws_secret_key": schema.StringAttribute{
				Description: "AWS secret key.",
				Optional:    true,
			},
			"password": schema.StringAttribute{
				Description: "Password if needed.",
				Optional:    true,
			},
			"private_key": schema.StringAttribute{
				Description: "Private key if needed.",
				Optional:    true,
			},
			"private_key_passphrase": schema.StringAttribute{
				Description: "Passphrase for private key if needed.",
				Optional:    true,
			},
			"ssl_certificate": schema.StringAttribute{
				Description: "SSL client certificate.",
				Optional:    true,
			},
			"google_cloud_storage_credentials_json": schema.StringAttribute{
				Description: "A JSON file that contains the private key. To generate see https://cloud.google.com/storage/docs/json_api/v1/how-tos/authorizing#APIKey",
				Optional:    true,
			},
			"wasabi_secret_key": schema.StringAttribute{
				Description: "Wasabi secret key.",
				Optional:    true,
			},
			"backblaze_b2_key_id": schema.StringAttribute{
				Description: "Backblaze B2 Cloud Storage keyID.",
				Optional:    true,
			},
			"backblaze_b2_application_key": schema.StringAttribute{
				Description: "Backblaze B2 Cloud Storage applicationKey.",
				Optional:    true,
			},
			"rackspace_api_key": schema.StringAttribute{
				Description: "Rackspace API key from the Rackspace Cloud Control Panel.",
				Optional:    true,
			},
			"reset_authentication": schema.BoolAttribute{
				Description: "Reset authenticated account",
				Optional:    true,
			},
			"azure_blob_storage_access_key": schema.StringAttribute{
				Description: "Azure Blob Storage secret key.",
				Optional:    true,
			},
			"azure_files_storage_access_key": schema.StringAttribute{
				Description: "Azure File Storage access key.",
				Optional:    true,
			},
			"azure_blob_storage_sas_token": schema.StringAttribute{
				Description: "Shared Access Signature (SAS) token",
				Optional:    true,
			},
			"azure_files_storage_sas_token": schema.StringAttribute{
				Description: "Shared Access Signature (SAS) token",
				Optional:    true,
			},
			"s3_compatible_secret_key": schema.StringAttribute{
				Description: "S3-compatible secret key",
				Optional:    true,
			},
			"filebase_secret_key": schema.StringAttribute{
				Description: "Filebase secret key",
				Optional:    true,
			},
			"cloudflare_secret_key": schema.StringAttribute{
				Description: "Cloudflare secret key",
				Optional:    true,
			},
			"linode_secret_key": schema.StringAttribute{
				Description: "Linode secret key",
				Optional:    true,
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
				Description: "If set, all communciations with this remote server are made through the provided region.",
				Computed:    true,
			},
			"auth_setup_link": schema.StringAttribute{
				Description: "Returns link to login with an Oauth provider",
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

	paramsRemoteServerCreate := files_sdk.RemoteServerCreateParams{}
	paramsRemoteServerCreate.AwsAccessKey = plan.AwsAccessKey.ValueString()
	paramsRemoteServerCreate.AwsSecretKey = plan.AwsSecretKey.ValueString()
	paramsRemoteServerCreate.Password = plan.Password.ValueString()
	paramsRemoteServerCreate.PrivateKey = plan.PrivateKey.ValueString()
	paramsRemoteServerCreate.PrivateKeyPassphrase = plan.PrivateKeyPassphrase.ValueString()
	paramsRemoteServerCreate.SslCertificate = plan.SslCertificate.ValueString()
	paramsRemoteServerCreate.GoogleCloudStorageCredentialsJson = plan.GoogleCloudStorageCredentialsJson.ValueString()
	paramsRemoteServerCreate.WasabiAccessKey = plan.WasabiAccessKey.ValueString()
	paramsRemoteServerCreate.WasabiSecretKey = plan.WasabiSecretKey.ValueString()
	paramsRemoteServerCreate.BackblazeB2KeyId = plan.BackblazeB2KeyId.ValueString()
	paramsRemoteServerCreate.BackblazeB2ApplicationKey = plan.BackblazeB2ApplicationKey.ValueString()
	paramsRemoteServerCreate.RackspaceApiKey = plan.RackspaceApiKey.ValueString()
	paramsRemoteServerCreate.ResetAuthentication = plan.ResetAuthentication.ValueBoolPointer()
	paramsRemoteServerCreate.AzureBlobStorageAccessKey = plan.AzureBlobStorageAccessKey.ValueString()
	paramsRemoteServerCreate.AzureFilesStorageAccessKey = plan.AzureFilesStorageAccessKey.ValueString()
	paramsRemoteServerCreate.Hostname = plan.Hostname.ValueString()
	paramsRemoteServerCreate.Name = plan.Name.ValueString()
	paramsRemoteServerCreate.MaxConnections = plan.MaxConnections.ValueInt64()
	paramsRemoteServerCreate.PinToSiteRegion = plan.PinToSiteRegion.ValueBoolPointer()
	paramsRemoteServerCreate.Port = plan.Port.ValueInt64()
	paramsRemoteServerCreate.S3Bucket = plan.S3Bucket.ValueString()
	paramsRemoteServerCreate.S3Region = plan.S3Region.ValueString()
	paramsRemoteServerCreate.ServerCertificate = paramsRemoteServerCreate.ServerCertificate.Enum()[plan.ServerCertificate.ValueString()]
	paramsRemoteServerCreate.ServerHostKey = plan.ServerHostKey.ValueString()
	paramsRemoteServerCreate.ServerType = paramsRemoteServerCreate.ServerType.Enum()[plan.ServerType.ValueString()]
	paramsRemoteServerCreate.Ssl = paramsRemoteServerCreate.Ssl.Enum()[plan.Ssl.ValueString()]
	paramsRemoteServerCreate.Username = plan.Username.ValueString()
	paramsRemoteServerCreate.GoogleCloudStorageBucket = plan.GoogleCloudStorageBucket.ValueString()
	paramsRemoteServerCreate.GoogleCloudStorageProjectId = plan.GoogleCloudStorageProjectId.ValueString()
	paramsRemoteServerCreate.BackblazeB2Bucket = plan.BackblazeB2Bucket.ValueString()
	paramsRemoteServerCreate.BackblazeB2S3Endpoint = plan.BackblazeB2S3Endpoint.ValueString()
	paramsRemoteServerCreate.WasabiBucket = plan.WasabiBucket.ValueString()
	paramsRemoteServerCreate.WasabiRegion = plan.WasabiRegion.ValueString()
	paramsRemoteServerCreate.RackspaceUsername = plan.RackspaceUsername.ValueString()
	paramsRemoteServerCreate.RackspaceRegion = plan.RackspaceRegion.ValueString()
	paramsRemoteServerCreate.RackspaceContainer = plan.RackspaceContainer.ValueString()
	paramsRemoteServerCreate.OneDriveAccountType = paramsRemoteServerCreate.OneDriveAccountType.Enum()[plan.OneDriveAccountType.ValueString()]
	paramsRemoteServerCreate.AzureBlobStorageAccount = plan.AzureBlobStorageAccount.ValueString()
	paramsRemoteServerCreate.AzureBlobStorageContainer = plan.AzureBlobStorageContainer.ValueString()
	paramsRemoteServerCreate.AzureBlobStorageHierarchicalNamespace = plan.AzureBlobStorageHierarchicalNamespace.ValueBoolPointer()
	paramsRemoteServerCreate.AzureBlobStorageSasToken = plan.AzureBlobStorageSasToken.ValueString()
	paramsRemoteServerCreate.AzureFilesStorageAccount = plan.AzureFilesStorageAccount.ValueString()
	paramsRemoteServerCreate.AzureFilesStorageShareName = plan.AzureFilesStorageShareName.ValueString()
	paramsRemoteServerCreate.AzureFilesStorageSasToken = plan.AzureFilesStorageSasToken.ValueString()
	paramsRemoteServerCreate.S3CompatibleBucket = plan.S3CompatibleBucket.ValueString()
	paramsRemoteServerCreate.S3CompatibleEndpoint = plan.S3CompatibleEndpoint.ValueString()
	paramsRemoteServerCreate.S3CompatibleRegion = plan.S3CompatibleRegion.ValueString()
	paramsRemoteServerCreate.EnableDedicatedIps = plan.EnableDedicatedIps.ValueBoolPointer()
	paramsRemoteServerCreate.S3CompatibleAccessKey = plan.S3CompatibleAccessKey.ValueString()
	paramsRemoteServerCreate.S3CompatibleSecretKey = plan.S3CompatibleSecretKey.ValueString()
	paramsRemoteServerCreate.FilesAgentRoot = plan.FilesAgentRoot.ValueString()
	paramsRemoteServerCreate.FilesAgentPermissionSet = paramsRemoteServerCreate.FilesAgentPermissionSet.Enum()[plan.FilesAgentPermissionSet.ValueString()]
	paramsRemoteServerCreate.FilesAgentVersion = plan.FilesAgentVersion.ValueString()
	paramsRemoteServerCreate.FilebaseAccessKey = plan.FilebaseAccessKey.ValueString()
	paramsRemoteServerCreate.FilebaseSecretKey = plan.FilebaseSecretKey.ValueString()
	paramsRemoteServerCreate.FilebaseBucket = plan.FilebaseBucket.ValueString()
	paramsRemoteServerCreate.CloudflareAccessKey = plan.CloudflareAccessKey.ValueString()
	paramsRemoteServerCreate.CloudflareSecretKey = plan.CloudflareSecretKey.ValueString()
	paramsRemoteServerCreate.CloudflareBucket = plan.CloudflareBucket.ValueString()
	paramsRemoteServerCreate.CloudflareEndpoint = plan.CloudflareEndpoint.ValueString()
	paramsRemoteServerCreate.DropboxTeams = plan.DropboxTeams.ValueBoolPointer()
	paramsRemoteServerCreate.LinodeAccessKey = plan.LinodeAccessKey.ValueString()
	paramsRemoteServerCreate.LinodeSecretKey = plan.LinodeSecretKey.ValueString()
	paramsRemoteServerCreate.LinodeBucket = plan.LinodeBucket.ValueString()
	paramsRemoteServerCreate.LinodeRegion = plan.LinodeRegion.ValueString()

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

	paramsRemoteServerUpdate := files_sdk.RemoteServerUpdateParams{}
	paramsRemoteServerUpdate.Id = plan.Id.ValueInt64()
	paramsRemoteServerUpdate.AwsAccessKey = plan.AwsAccessKey.ValueString()
	paramsRemoteServerUpdate.AwsSecretKey = plan.AwsSecretKey.ValueString()
	paramsRemoteServerUpdate.Password = plan.Password.ValueString()
	paramsRemoteServerUpdate.PrivateKey = plan.PrivateKey.ValueString()
	paramsRemoteServerUpdate.PrivateKeyPassphrase = plan.PrivateKeyPassphrase.ValueString()
	paramsRemoteServerUpdate.SslCertificate = plan.SslCertificate.ValueString()
	paramsRemoteServerUpdate.GoogleCloudStorageCredentialsJson = plan.GoogleCloudStorageCredentialsJson.ValueString()
	paramsRemoteServerUpdate.WasabiAccessKey = plan.WasabiAccessKey.ValueString()
	paramsRemoteServerUpdate.WasabiSecretKey = plan.WasabiSecretKey.ValueString()
	paramsRemoteServerUpdate.BackblazeB2KeyId = plan.BackblazeB2KeyId.ValueString()
	paramsRemoteServerUpdate.BackblazeB2ApplicationKey = plan.BackblazeB2ApplicationKey.ValueString()
	paramsRemoteServerUpdate.RackspaceApiKey = plan.RackspaceApiKey.ValueString()
	paramsRemoteServerUpdate.ResetAuthentication = plan.ResetAuthentication.ValueBoolPointer()
	paramsRemoteServerUpdate.AzureBlobStorageAccessKey = plan.AzureBlobStorageAccessKey.ValueString()
	paramsRemoteServerUpdate.AzureFilesStorageAccessKey = plan.AzureFilesStorageAccessKey.ValueString()
	paramsRemoteServerUpdate.Hostname = plan.Hostname.ValueString()
	paramsRemoteServerUpdate.Name = plan.Name.ValueString()
	paramsRemoteServerUpdate.MaxConnections = plan.MaxConnections.ValueInt64()
	paramsRemoteServerUpdate.PinToSiteRegion = plan.PinToSiteRegion.ValueBoolPointer()
	paramsRemoteServerUpdate.Port = plan.Port.ValueInt64()
	paramsRemoteServerUpdate.S3Bucket = plan.S3Bucket.ValueString()
	paramsRemoteServerUpdate.S3Region = plan.S3Region.ValueString()
	paramsRemoteServerUpdate.ServerCertificate = paramsRemoteServerUpdate.ServerCertificate.Enum()[plan.ServerCertificate.ValueString()]
	paramsRemoteServerUpdate.ServerHostKey = plan.ServerHostKey.ValueString()
	paramsRemoteServerUpdate.ServerType = paramsRemoteServerUpdate.ServerType.Enum()[plan.ServerType.ValueString()]
	paramsRemoteServerUpdate.Ssl = paramsRemoteServerUpdate.Ssl.Enum()[plan.Ssl.ValueString()]
	paramsRemoteServerUpdate.Username = plan.Username.ValueString()
	paramsRemoteServerUpdate.GoogleCloudStorageBucket = plan.GoogleCloudStorageBucket.ValueString()
	paramsRemoteServerUpdate.GoogleCloudStorageProjectId = plan.GoogleCloudStorageProjectId.ValueString()
	paramsRemoteServerUpdate.BackblazeB2Bucket = plan.BackblazeB2Bucket.ValueString()
	paramsRemoteServerUpdate.BackblazeB2S3Endpoint = plan.BackblazeB2S3Endpoint.ValueString()
	paramsRemoteServerUpdate.WasabiBucket = plan.WasabiBucket.ValueString()
	paramsRemoteServerUpdate.WasabiRegion = plan.WasabiRegion.ValueString()
	paramsRemoteServerUpdate.RackspaceUsername = plan.RackspaceUsername.ValueString()
	paramsRemoteServerUpdate.RackspaceRegion = plan.RackspaceRegion.ValueString()
	paramsRemoteServerUpdate.RackspaceContainer = plan.RackspaceContainer.ValueString()
	paramsRemoteServerUpdate.OneDriveAccountType = paramsRemoteServerUpdate.OneDriveAccountType.Enum()[plan.OneDriveAccountType.ValueString()]
	paramsRemoteServerUpdate.AzureBlobStorageAccount = plan.AzureBlobStorageAccount.ValueString()
	paramsRemoteServerUpdate.AzureBlobStorageContainer = plan.AzureBlobStorageContainer.ValueString()
	paramsRemoteServerUpdate.AzureBlobStorageHierarchicalNamespace = plan.AzureBlobStorageHierarchicalNamespace.ValueBoolPointer()
	paramsRemoteServerUpdate.AzureBlobStorageSasToken = plan.AzureBlobStorageSasToken.ValueString()
	paramsRemoteServerUpdate.AzureFilesStorageAccount = plan.AzureFilesStorageAccount.ValueString()
	paramsRemoteServerUpdate.AzureFilesStorageShareName = plan.AzureFilesStorageShareName.ValueString()
	paramsRemoteServerUpdate.AzureFilesStorageSasToken = plan.AzureFilesStorageSasToken.ValueString()
	paramsRemoteServerUpdate.S3CompatibleBucket = plan.S3CompatibleBucket.ValueString()
	paramsRemoteServerUpdate.S3CompatibleEndpoint = plan.S3CompatibleEndpoint.ValueString()
	paramsRemoteServerUpdate.S3CompatibleRegion = plan.S3CompatibleRegion.ValueString()
	paramsRemoteServerUpdate.EnableDedicatedIps = plan.EnableDedicatedIps.ValueBoolPointer()
	paramsRemoteServerUpdate.S3CompatibleAccessKey = plan.S3CompatibleAccessKey.ValueString()
	paramsRemoteServerUpdate.S3CompatibleSecretKey = plan.S3CompatibleSecretKey.ValueString()
	paramsRemoteServerUpdate.FilesAgentRoot = plan.FilesAgentRoot.ValueString()
	paramsRemoteServerUpdate.FilesAgentPermissionSet = paramsRemoteServerUpdate.FilesAgentPermissionSet.Enum()[plan.FilesAgentPermissionSet.ValueString()]
	paramsRemoteServerUpdate.FilesAgentVersion = plan.FilesAgentVersion.ValueString()
	paramsRemoteServerUpdate.FilebaseAccessKey = plan.FilebaseAccessKey.ValueString()
	paramsRemoteServerUpdate.FilebaseSecretKey = plan.FilebaseSecretKey.ValueString()
	paramsRemoteServerUpdate.FilebaseBucket = plan.FilebaseBucket.ValueString()
	paramsRemoteServerUpdate.CloudflareAccessKey = plan.CloudflareAccessKey.ValueString()
	paramsRemoteServerUpdate.CloudflareSecretKey = plan.CloudflareSecretKey.ValueString()
	paramsRemoteServerUpdate.CloudflareBucket = plan.CloudflareBucket.ValueString()
	paramsRemoteServerUpdate.CloudflareEndpoint = plan.CloudflareEndpoint.ValueString()
	paramsRemoteServerUpdate.DropboxTeams = plan.DropboxTeams.ValueBoolPointer()
	paramsRemoteServerUpdate.LinodeAccessKey = plan.LinodeAccessKey.ValueString()
	paramsRemoteServerUpdate.LinodeSecretKey = plan.LinodeSecretKey.ValueString()
	paramsRemoteServerUpdate.LinodeBucket = plan.LinodeBucket.ValueString()
	paramsRemoteServerUpdate.LinodeRegion = plan.LinodeRegion.ValueString()

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
	state.BackblazeB2S3Endpoint = types.StringValue(remoteServer.BackblazeB2S3Endpoint)
	state.BackblazeB2Bucket = types.StringValue(remoteServer.BackblazeB2Bucket)
	state.WasabiBucket = types.StringValue(remoteServer.WasabiBucket)
	state.WasabiRegion = types.StringValue(remoteServer.WasabiRegion)
	state.WasabiAccessKey = types.StringValue(remoteServer.WasabiAccessKey)
	state.RackspaceUsername = types.StringValue(remoteServer.RackspaceUsername)
	state.RackspaceRegion = types.StringValue(remoteServer.RackspaceRegion)
	state.RackspaceContainer = types.StringValue(remoteServer.RackspaceContainer)
	state.AuthSetupLink = types.StringValue(remoteServer.AuthSetupLink)
	state.AuthStatus = types.StringValue(remoteServer.AuthStatus)
	state.AuthAccountName = types.StringValue(remoteServer.AuthAccountName)
	state.OneDriveAccountType = types.StringValue(remoteServer.OneDriveAccountType)
	state.AzureBlobStorageAccount = types.StringValue(remoteServer.AzureBlobStorageAccount)
	state.AzureBlobStorageContainer = types.StringValue(remoteServer.AzureBlobStorageContainer)
	state.AzureBlobStorageHierarchicalNamespace = types.BoolPointerValue(remoteServer.AzureBlobStorageHierarchicalNamespace)
	state.AzureFilesStorageAccount = types.StringValue(remoteServer.AzureFilesStorageAccount)
	state.AzureFilesStorageShareName = types.StringValue(remoteServer.AzureFilesStorageShareName)
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
