package provider

import (
	"context"
	"fmt"

	files_sdk "github.com/Files-com/files-sdk-go/v3"
	remote_server "github.com/Files-com/files-sdk-go/v3/remoteserver"
	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/datasource/schema"
	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

var (
	_ datasource.DataSource              = &remoteServerDataSource{}
	_ datasource.DataSourceWithConfigure = &remoteServerDataSource{}
)

func NewRemoteServerDataSource() datasource.DataSource {
	return &remoteServerDataSource{}
}

type remoteServerDataSource struct {
	client *remote_server.Client
}

type remoteServerDataSourceModel struct {
	Id                                      types.Int64  `tfsdk:"id"`
	Disabled                                types.Bool   `tfsdk:"disabled"`
	AuthenticationMethod                    types.String `tfsdk:"authentication_method"`
	Hostname                                types.String `tfsdk:"hostname"`
	RemoteHomePath                          types.String `tfsdk:"remote_home_path"`
	Name                                    types.String `tfsdk:"name"`
	Port                                    types.Int64  `tfsdk:"port"`
	MaxConnections                          types.Int64  `tfsdk:"max_connections"`
	PinToSiteRegion                         types.Bool   `tfsdk:"pin_to_site_region"`
	PinnedRegion                            types.String `tfsdk:"pinned_region"`
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
	AuthStatus                              types.String `tfsdk:"auth_status"`
	AuthAccountName                         types.String `tfsdk:"auth_account_name"`
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
	FilesAgentApiToken                      types.String `tfsdk:"files_agent_api_token"`
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
	SupportsVersioning                      types.Bool   `tfsdk:"supports_versioning"`
}

func (r *remoteServerDataSource) Configure(_ context.Context, req datasource.ConfigureRequest, resp *datasource.ConfigureResponse) {
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

func (r *remoteServerDataSource) Metadata(_ context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_remote_server"
}

func (r *remoteServerDataSource) Schema(_ context.Context, _ datasource.SchemaRequest, resp *datasource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Description: "A RemoteServer is a specific type of Behavior called `remote_server_sync`.\n\n\n\nRemote Servers can be either an FTP server, SFTP server, S3 bucket, Google Cloud Storage, Wasabi, Backblaze B2 Cloud Storage, Rackspace Cloud Files container, WebDAV, Box, Dropbox, OneDrive, Google Drive, or Azure Blob Storage.\n\n\n\nNot every attribute will apply to every remote server.\n\n\n\nFTP Servers require that you specify their `hostname`, `port`, `username`, `password`, and a value for `ssl`. Optionally, provide `server_certificate`.\n\n\n\nSFTP Servers require that you specify their `hostname`, `port`, `username`, `password` or `private_key`, and a value for `ssl`. Optionally, provide `server_certificate`, `private_key_passphrase`.\n\n\n\nS3 Buckets require that you specify their `s3_bucket` name, and `s3_region`. Optionally provide a `aws_access_key`, and `aws_secret_key`. If you don't provide credentials, you will need to use AWS to grant us access to your bucket.\n\n\n\nS3-Compatible Buckets require that you specify `s3_compatible_bucket`, `s3_compatible_endpoint`, `s3_compatible_access_key`, and `s3_compatible_secret_key`.\n\n\n\nGoogle Cloud Storage requires that you specify `google_cloud_storage_bucket`, and then one of the following sets of authentication credentials:\n\n - for JSON authentcation: `google_cloud_storage_project_id`, and `google_cloud_storage_credentials_json`\n\n - for HMAC (S3-Compatible) authentication: `google_cloud_storage_s3_compatible_access_key`, and `google_cloud_storage_s3_compatible_secret_key`\n\n\n\nWasabi requires `wasabi_bucket`, `wasabi_region`, `wasabi_access_key`, and `wasabi_secret_key`.\n\n\n\nBackblaze B2 Cloud Storage `backblaze_b2_bucket`, `backblaze_b2_s3_endpoint`, `backblaze_b2_application_key`, and `backblaze_b2_key_id`. (Requires S3 Compatible API) See https://help.backblaze.com/hc/en-us/articles/360047425453\n\n\n\nWebDAV Servers require that you specify their `hostname`, `username`, and `password`.\n\n\n\nOneDrive follow the `auth_setup_link` and login with Microsoft.\n\n\n\nSharepoint follow the `auth_setup_link` and login with Microsoft.\n\n\n\nBox follow the `auth_setup_link` and login with Box.\n\n\n\nDropbox specify if `dropbox_teams` then follow the `auth_setup_link` and login with Dropbox.\n\n\n\nGoogle Drive follow the `auth_setup_link` and login with Google.\n\n\n\nAzure Blob Storage `azure_blob_storage_account`, `azure_blob_storage_container`, `azure_blob_storage_access_key`, `azure_blob_storage_sas_token`, `azure_blob_storage_dns_suffix`\n\n\n\nAzure File Storage `azure_files_storage_account`, `azure_files_storage_access_key`, `azure_files_storage_share_name`, `azure_files_storage_dns_suffix`\n\n\n\nFilebase requires `filebase_bucket`, `filebase_access_key`, and `filebase_secret_key`.\n\n\n\nCloudflare requires `cloudflare_bucket`, `cloudflare_access_key`, `cloudflare_secret_key` and `cloudflare_endpoint`.\n\n\n\nLinode requires `linode_bucket`, `linode_access_key`, `linode_secret_key` and `linode_region`.",
		Attributes: map[string]schema.Attribute{
			"id": schema.Int64Attribute{
				Description: "Remote server ID",
				Required:    true,
			},
			"disabled": schema.BoolAttribute{
				Description: "If true, this server has been disabled due to failures.  Make any change or set disabled to false to clear this flag.",
				Computed:    true,
			},
			"authentication_method": schema.StringAttribute{
				Description: "Type of authentication method",
				Computed:    true,
			},
			"hostname": schema.StringAttribute{
				Description: "Hostname or IP address",
				Computed:    true,
			},
			"remote_home_path": schema.StringAttribute{
				Description: "Initial home folder on remote server",
				Computed:    true,
			},
			"name": schema.StringAttribute{
				Description: "Internal name for your reference",
				Computed:    true,
			},
			"port": schema.Int64Attribute{
				Description: "Port for remote server.  Not needed for S3.",
				Computed:    true,
			},
			"max_connections": schema.Int64Attribute{
				Description: "Max number of parallel connections.  Ignored for S3 connections (we will parallelize these as much as possible).",
				Computed:    true,
			},
			"pin_to_site_region": schema.BoolAttribute{
				Description: "If true, we will ensure that all communications with this remote server are made through the primary region of the site.  This setting can also be overridden by a site-wide setting which will force it to true.",
				Computed:    true,
			},
			"pinned_region": schema.StringAttribute{
				Description: "If set, all communications with this remote server are made through the provided region.",
				Computed:    true,
			},
			"s3_bucket": schema.StringAttribute{
				Description: "S3 bucket name",
				Computed:    true,
			},
			"s3_region": schema.StringAttribute{
				Description: "S3 region",
				Computed:    true,
			},
			"aws_access_key": schema.StringAttribute{
				Description: "AWS Access Key.",
				Computed:    true,
			},
			"server_certificate": schema.StringAttribute{
				Description: "Remote server certificate",
				Computed:    true,
			},
			"server_host_key": schema.StringAttribute{
				Description: "Remote server SSH Host Key. If provided, we will require that the server host key matches the provided key. Uses OpenSSH format similar to what would go into ~/.ssh/known_hosts",
				Computed:    true,
			},
			"server_type": schema.StringAttribute{
				Description: "Remote server type.",
				Computed:    true,
			},
			"ssl": schema.StringAttribute{
				Description: "Should we require SSL?",
				Computed:    true,
			},
			"username": schema.StringAttribute{
				Description: "Remote server username.  Not needed for S3 buckets.",
				Computed:    true,
			},
			"google_cloud_storage_bucket": schema.StringAttribute{
				Description: "Google Cloud Storage: Bucket Name",
				Computed:    true,
			},
			"google_cloud_storage_project_id": schema.StringAttribute{
				Description: "Google Cloud Storage: Project ID",
				Computed:    true,
			},
			"google_cloud_storage_s3_compatible_access_key": schema.StringAttribute{
				Description: "Google Cloud Storage: S3-compatible Access Key.",
				Computed:    true,
			},
			"backblaze_b2_s3_endpoint": schema.StringAttribute{
				Description: "Backblaze B2 Cloud Storage: S3 Endpoint",
				Computed:    true,
			},
			"backblaze_b2_bucket": schema.StringAttribute{
				Description: "Backblaze B2 Cloud Storage: Bucket name",
				Computed:    true,
			},
			"wasabi_bucket": schema.StringAttribute{
				Description: "Wasabi: Bucket name",
				Computed:    true,
			},
			"wasabi_region": schema.StringAttribute{
				Description: "Wasabi: Region",
				Computed:    true,
			},
			"wasabi_access_key": schema.StringAttribute{
				Description: "Wasabi: Access Key.",
				Computed:    true,
			},
			"auth_status": schema.StringAttribute{
				Description: "Either `in_setup` or `complete`",
				Computed:    true,
			},
			"auth_account_name": schema.StringAttribute{
				Description: "Describes the authorized account",
				Computed:    true,
			},
			"one_drive_account_type": schema.StringAttribute{
				Description: "OneDrive: Either personal or business_other account types",
				Computed:    true,
			},
			"azure_blob_storage_account": schema.StringAttribute{
				Description: "Azure Blob Storage: Account name",
				Computed:    true,
			},
			"azure_blob_storage_container": schema.StringAttribute{
				Description: "Azure Blob Storage: Container name",
				Computed:    true,
			},
			"azure_blob_storage_hierarchical_namespace": schema.BoolAttribute{
				Description: "Azure Blob Storage: Does the storage account has hierarchical namespace feature enabled?",
				Computed:    true,
			},
			"azure_blob_storage_dns_suffix": schema.StringAttribute{
				Description: "Azure Blob Storage: Custom DNS suffix",
				Computed:    true,
			},
			"azure_files_storage_account": schema.StringAttribute{
				Description: "Azure Files: Storage Account name",
				Computed:    true,
			},
			"azure_files_storage_share_name": schema.StringAttribute{
				Description: "Azure Files:  Storage Share name",
				Computed:    true,
			},
			"azure_files_storage_dns_suffix": schema.StringAttribute{
				Description: "Azure Files: Custom DNS suffix",
				Computed:    true,
			},
			"s3_compatible_bucket": schema.StringAttribute{
				Description: "S3-compatible: Bucket name",
				Computed:    true,
			},
			"s3_compatible_endpoint": schema.StringAttribute{
				Description: "S3-compatible: endpoint",
				Computed:    true,
			},
			"s3_compatible_region": schema.StringAttribute{
				Description: "S3-compatible: region",
				Computed:    true,
			},
			"s3_compatible_access_key": schema.StringAttribute{
				Description: "S3-compatible: Access Key",
				Computed:    true,
			},
			"enable_dedicated_ips": schema.BoolAttribute{
				Description: "`true` if remote server only accepts connections from dedicated IPs",
				Computed:    true,
			},
			"files_agent_permission_set": schema.StringAttribute{
				Description: "Local permissions for files agent. read_only, write_only, or read_write",
				Computed:    true,
			},
			"files_agent_root": schema.StringAttribute{
				Description: "Agent local root path",
				Computed:    true,
			},
			"files_agent_api_token": schema.StringAttribute{
				Description: "Files Agent API Token",
				Computed:    true,
			},
			"files_agent_version": schema.StringAttribute{
				Description: "Files Agent version",
				Computed:    true,
			},
			"filebase_bucket": schema.StringAttribute{
				Description: "Filebase: Bucket name",
				Computed:    true,
			},
			"filebase_access_key": schema.StringAttribute{
				Description: "Filebase: Access Key.",
				Computed:    true,
			},
			"cloudflare_bucket": schema.StringAttribute{
				Description: "Cloudflare: Bucket name",
				Computed:    true,
			},
			"cloudflare_access_key": schema.StringAttribute{
				Description: "Cloudflare: Access Key.",
				Computed:    true,
			},
			"cloudflare_endpoint": schema.StringAttribute{
				Description: "Cloudflare: endpoint",
				Computed:    true,
			},
			"dropbox_teams": schema.BoolAttribute{
				Description: "Dropbox: If true, list Team folders in root?",
				Computed:    true,
			},
			"linode_bucket": schema.StringAttribute{
				Description: "Linode: Bucket name",
				Computed:    true,
			},
			"linode_access_key": schema.StringAttribute{
				Description: "Linode: Access Key",
				Computed:    true,
			},
			"linode_region": schema.StringAttribute{
				Description: "Linode: region",
				Computed:    true,
			},
			"supports_versioning": schema.BoolAttribute{
				Description: "If true, this remote server supports file versioning. This value is determined automatically by Files.com.",
				Computed:    true,
			},
		},
	}
}

func (r *remoteServerDataSource) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
	var data remoteServerDataSourceModel
	diags := req.Config.Get(ctx, &data)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	paramsRemoteServerFind := files_sdk.RemoteServerFindParams{}
	paramsRemoteServerFind.Id = data.Id.ValueInt64()

	remoteServer, err := r.client.Find(paramsRemoteServerFind, files_sdk.WithContext(ctx))
	if err != nil {
		resp.Diagnostics.AddError(
			"Error Reading Files RemoteServer",
			"Could not read remote_server id "+fmt.Sprint(data.Id.ValueInt64())+": "+err.Error(),
		)
		return
	}

	diags = r.populateDataSourceModel(ctx, remoteServer, &data)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	diags = resp.State.Set(ctx, data)
	resp.Diagnostics.Append(diags...)
}

func (r *remoteServerDataSource) populateDataSourceModel(ctx context.Context, remoteServer files_sdk.RemoteServer, state *remoteServerDataSourceModel) (diags diag.Diagnostics) {
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
