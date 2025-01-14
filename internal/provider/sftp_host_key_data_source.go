package provider

import (
	"context"
	"fmt"

	files_sdk "github.com/Files-com/files-sdk-go/v3"
	sftp_host_key "github.com/Files-com/files-sdk-go/v3/sftphostkey"
	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/datasource/schema"
	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

var (
	_ datasource.DataSource              = &sftpHostKeyDataSource{}
	_ datasource.DataSourceWithConfigure = &sftpHostKeyDataSource{}
)

func NewSftpHostKeyDataSource() datasource.DataSource {
	return &sftpHostKeyDataSource{}
}

type sftpHostKeyDataSource struct {
	client *sftp_host_key.Client
}

type sftpHostKeyDataSourceModel struct {
	Id                types.Int64  `tfsdk:"id"`
	Name              types.String `tfsdk:"name"`
	FingerprintMd5    types.String `tfsdk:"fingerprint_md5"`
	FingerprintSha256 types.String `tfsdk:"fingerprint_sha256"`
}

func (r *sftpHostKeyDataSource) Configure(_ context.Context, req datasource.ConfigureRequest, resp *datasource.ConfigureResponse) {
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

	r.client = &sftp_host_key.Client{Config: sdk_config}
}

func (r *sftpHostKeyDataSource) Metadata(_ context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_sftp_host_key"
}

func (r *sftpHostKeyDataSource) Schema(_ context.Context, _ datasource.SchemaRequest, resp *datasource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Description: "An SFTP Host Key is a cryptographic key used to verify the identity of the server during an SFTP connection. This allows the client to be sure that it is connecting to the intended server, preventing man-in-the-middle attacks and ensuring secure communication between the client and Files.com.\n\n\n\nFiles.com allows you to provide custom SFTP Host Keys, which is particularly useful when migrating to Files.com from an existing SFTP server, allowing the Files.com platform to match your previously-installed host key for a seamless transition.",
		Attributes: map[string]schema.Attribute{
			"id": schema.Int64Attribute{
				Description: "SFTP Host Key ID",
				Required:    true,
			},
			"name": schema.StringAttribute{
				Description: "The friendly name of this SFTP Host Key.",
				Computed:    true,
			},
			"fingerprint_md5": schema.StringAttribute{
				Description: "MD5 Fingerprint of the public key",
				Computed:    true,
			},
			"fingerprint_sha256": schema.StringAttribute{
				Description: "SHA256 Fingerprint of the public key",
				Computed:    true,
			},
		},
	}
}

func (r *sftpHostKeyDataSource) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
	var data sftpHostKeyDataSourceModel
	diags := req.Config.Get(ctx, &data)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	paramsSftpHostKeyFind := files_sdk.SftpHostKeyFindParams{}
	paramsSftpHostKeyFind.Id = data.Id.ValueInt64()

	sftpHostKey, err := r.client.Find(paramsSftpHostKeyFind, files_sdk.WithContext(ctx))
	if err != nil {
		resp.Diagnostics.AddError(
			"Error Reading Files SftpHostKey",
			"Could not read sftp_host_key id "+fmt.Sprint(data.Id.ValueInt64())+": "+err.Error(),
		)
		return
	}

	diags = r.populateDataSourceModel(ctx, sftpHostKey, &data)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	diags = resp.State.Set(ctx, data)
	resp.Diagnostics.Append(diags...)
}

func (r *sftpHostKeyDataSource) populateDataSourceModel(ctx context.Context, sftpHostKey files_sdk.SftpHostKey, state *sftpHostKeyDataSourceModel) (diags diag.Diagnostics) {
	state.Id = types.Int64Value(sftpHostKey.Id)
	state.Name = types.StringValue(sftpHostKey.Name)
	state.FingerprintMd5 = types.StringValue(sftpHostKey.FingerprintMd5)
	state.FingerprintSha256 = types.StringValue(sftpHostKey.FingerprintSha256)

	return
}
