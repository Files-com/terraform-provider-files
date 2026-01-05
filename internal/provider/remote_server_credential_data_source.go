package provider

import (
	"context"
	"fmt"

	files_sdk "github.com/Files-com/files-sdk-go/v3"
	remote_server_credential "github.com/Files-com/files-sdk-go/v3/remoteservercredential"
	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/datasource/schema"
	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

var (
	_ datasource.DataSource              = &remoteServerCredentialDataSource{}
	_ datasource.DataSourceWithConfigure = &remoteServerCredentialDataSource{}
)

func NewRemoteServerCredentialDataSource() datasource.DataSource {
	return &remoteServerCredentialDataSource{}
}

type remoteServerCredentialDataSource struct {
	client *remote_server_credential.Client
}

type remoteServerCredentialDataSourceModel struct {
	Id                                      types.Int64  `tfsdk:"id"`
	WorkspaceId                             types.Int64  `tfsdk:"workspace_id"`
	Name                                    types.String `tfsdk:"name"`
	Description                             types.String `tfsdk:"description"`
	ServerType                              types.String `tfsdk:"server_type"`
	AwsAccessKey                            types.String `tfsdk:"aws_access_key"`
	GoogleCloudStorageS3CompatibleAccessKey types.String `tfsdk:"google_cloud_storage_s3_compatible_access_key"`
	WasabiAccessKey                         types.String `tfsdk:"wasabi_access_key"`
	AzureBlobStorageAccount                 types.String `tfsdk:"azure_blob_storage_account"`
	AzureFilesStorageAccount                types.String `tfsdk:"azure_files_storage_account"`
	S3CompatibleAccessKey                   types.String `tfsdk:"s3_compatible_access_key"`
	FilebaseAccessKey                       types.String `tfsdk:"filebase_access_key"`
	CloudflareAccessKey                     types.String `tfsdk:"cloudflare_access_key"`
	LinodeAccessKey                         types.String `tfsdk:"linode_access_key"`
	Username                                types.String `tfsdk:"username"`
}

func (r *remoteServerCredentialDataSource) Configure(_ context.Context, req datasource.ConfigureRequest, resp *datasource.ConfigureResponse) {
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

func (r *remoteServerCredentialDataSource) Metadata(_ context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_remote_server_credential"
}

func (r *remoteServerCredentialDataSource) Schema(_ context.Context, _ datasource.SchemaRequest, resp *datasource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Description: "A RemoteServerCredential is a way to store a credential for Remote Servers in a centralized vault and then reference it from Remote Server definitions.\n\n\n\nThis allows you to manage your credentials in one place and avoid duplicating them across multiple Remote Server configurations. It also enhances security by allowing you to use Terraform or APIs for Remote Server management without having to worry about credential exposure.",
		Attributes: map[string]schema.Attribute{
			"id": schema.Int64Attribute{
				Description: "Remote Server Credential ID",
				Required:    true,
			},
			"workspace_id": schema.Int64Attribute{
				Description: "Workspace ID (0 for default workspace)",
				Computed:    true,
			},
			"name": schema.StringAttribute{
				Description: "Internal name for your reference",
				Computed:    true,
			},
			"description": schema.StringAttribute{
				Description: "Internal description for your reference",
				Computed:    true,
			},
			"server_type": schema.StringAttribute{
				Description: "Remote server type.  Remote Server Credentials are only valid for a single type of Remote Server.",
				Computed:    true,
			},
			"aws_access_key": schema.StringAttribute{
				Description: "AWS Access Key.",
				Computed:    true,
			},
			"google_cloud_storage_s3_compatible_access_key": schema.StringAttribute{
				Description: "Google Cloud Storage: S3-compatible Access Key.",
				Computed:    true,
			},
			"wasabi_access_key": schema.StringAttribute{
				Description: "Wasabi: Access Key.",
				Computed:    true,
			},
			"azure_blob_storage_account": schema.StringAttribute{
				Description: "Azure Blob Storage: Account name",
				Computed:    true,
			},
			"azure_files_storage_account": schema.StringAttribute{
				Description: "Azure Files: Storage Account name",
				Computed:    true,
			},
			"s3_compatible_access_key": schema.StringAttribute{
				Description: "S3-compatible: Access Key",
				Computed:    true,
			},
			"filebase_access_key": schema.StringAttribute{
				Description: "Filebase: Access Key.",
				Computed:    true,
			},
			"cloudflare_access_key": schema.StringAttribute{
				Description: "Cloudflare: Access Key.",
				Computed:    true,
			},
			"linode_access_key": schema.StringAttribute{
				Description: "Linode: Access Key",
				Computed:    true,
			},
			"username": schema.StringAttribute{
				Description: "Remote server username.",
				Computed:    true,
			},
		},
	}
}

func (r *remoteServerCredentialDataSource) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
	var data remoteServerCredentialDataSourceModel
	diags := req.Config.Get(ctx, &data)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	paramsRemoteServerCredentialFind := files_sdk.RemoteServerCredentialFindParams{}
	paramsRemoteServerCredentialFind.Id = data.Id.ValueInt64()

	remoteServerCredential, err := r.client.Find(paramsRemoteServerCredentialFind, files_sdk.WithContext(ctx))
	if err != nil {
		resp.Diagnostics.AddError(
			"Error Reading Files RemoteServerCredential",
			"Could not read remote_server_credential id "+fmt.Sprint(data.Id.ValueInt64())+": "+err.Error(),
		)
		return
	}

	diags = r.populateDataSourceModel(ctx, remoteServerCredential, &data)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	diags = resp.State.Set(ctx, data)
	resp.Diagnostics.Append(diags...)
}

func (r *remoteServerCredentialDataSource) populateDataSourceModel(ctx context.Context, remoteServerCredential files_sdk.RemoteServerCredential, state *remoteServerCredentialDataSourceModel) (diags diag.Diagnostics) {
	state.Id = types.Int64Value(remoteServerCredential.Id)
	state.WorkspaceId = types.Int64Value(remoteServerCredential.WorkspaceId)
	state.Name = types.StringValue(remoteServerCredential.Name)
	state.Description = types.StringValue(remoteServerCredential.Description)
	state.ServerType = types.StringValue(remoteServerCredential.ServerType)
	state.AwsAccessKey = types.StringValue(remoteServerCredential.AwsAccessKey)
	state.GoogleCloudStorageS3CompatibleAccessKey = types.StringValue(remoteServerCredential.GoogleCloudStorageS3CompatibleAccessKey)
	state.WasabiAccessKey = types.StringValue(remoteServerCredential.WasabiAccessKey)
	state.AzureBlobStorageAccount = types.StringValue(remoteServerCredential.AzureBlobStorageAccount)
	state.AzureFilesStorageAccount = types.StringValue(remoteServerCredential.AzureFilesStorageAccount)
	state.S3CompatibleAccessKey = types.StringValue(remoteServerCredential.S3CompatibleAccessKey)
	state.FilebaseAccessKey = types.StringValue(remoteServerCredential.FilebaseAccessKey)
	state.CloudflareAccessKey = types.StringValue(remoteServerCredential.CloudflareAccessKey)
	state.LinodeAccessKey = types.StringValue(remoteServerCredential.LinodeAccessKey)
	state.Username = types.StringValue(remoteServerCredential.Username)

	return
}
