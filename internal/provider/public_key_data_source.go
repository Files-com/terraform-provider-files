package provider

import (
	"context"
	"fmt"

	files_sdk "github.com/Files-com/files-sdk-go/v3"
	public_key "github.com/Files-com/files-sdk-go/v3/publickey"
	"github.com/Files-com/terraform-provider-files/lib"
	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/datasource/schema"
	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

var (
	_ datasource.DataSource              = &publicKeyDataSource{}
	_ datasource.DataSourceWithConfigure = &publicKeyDataSource{}
)

func NewPublicKeyDataSource() datasource.DataSource {
	return &publicKeyDataSource{}
}

type publicKeyDataSource struct {
	client *public_key.Client
}

type publicKeyDataSourceModel struct {
	Id                types.Int64  `tfsdk:"id"`
	Title             types.String `tfsdk:"title"`
	CreatedAt         types.String `tfsdk:"created_at"`
	Fingerprint       types.String `tfsdk:"fingerprint"`
	FingerprintSha256 types.String `tfsdk:"fingerprint_sha256"`
	Status            types.String `tfsdk:"status"`
	LastLoginAt       types.String `tfsdk:"last_login_at"`
	PrivateKey        types.String `tfsdk:"private_key"`
	PublicKey         types.String `tfsdk:"public_key"`
	Username          types.String `tfsdk:"username"`
	UserId            types.Int64  `tfsdk:"user_id"`
}

func (r *publicKeyDataSource) Configure(_ context.Context, req datasource.ConfigureRequest, resp *datasource.ConfigureResponse) {
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

	r.client = &public_key.Client{Config: sdk_config}
}

func (r *publicKeyDataSource) Metadata(_ context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_public_key"
}

func (r *publicKeyDataSource) Schema(_ context.Context, _ datasource.SchemaRequest, resp *datasource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Description: "A PublicKey is used to authenticate to Files.com via SFTP (SSH File Transfer Protocol). This method of authentication allows users to use their private key (which is never shared with Files.com) to authenticate themselves against the PublicKey stored on Files.com.\n\n\n\nWhen a user configures their PublicKey, it allows them to bypass traditional password-based authentication, leveraging the security of key-based authentication instead.\n\n\n\nNote that Files.com's SSH support is limited to file operations only. While users can securely transfer files and manage their data via SFTP, they do not have access to a full shell environment for executing arbitrary commands.\n\n\n\nWhen generating new SSH keys, here are the available options: Files.com supports multiple SSH key algorithms: RSA (default 4096 bits, range 1024-4096 in 8-bit increments), DSA (1024 bits only), Ed25519 (256 bits), and ECDSA (256, 384, or 521 bits). When generating keys, the system uses these default lengths unless a specific length is specified.\n\n\n\nFiles.com also supports importing additional key types that cannot be generated: security key types (sk-ecdsa-sha2-nistp256, sk-ssh-ed25519). RSA keys up to 8192 bits are also supported for import.",
		Attributes: map[string]schema.Attribute{
			"id": schema.Int64Attribute{
				Description: "Public key ID",
				Required:    true,
			},
			"title": schema.StringAttribute{
				Description: "Public key title",
				Computed:    true,
			},
			"created_at": schema.StringAttribute{
				Description: "Public key created at date/time",
				Computed:    true,
			},
			"fingerprint": schema.StringAttribute{
				Description: "Public key fingerprint (MD5)",
				Computed:    true,
			},
			"fingerprint_sha256": schema.StringAttribute{
				Description: "Public key fingerprint (SHA256)",
				Computed:    true,
			},
			"status": schema.StringAttribute{
				Description: "Only returned when generating keys. Can be invalid, not_generated, generating, complete",
				Computed:    true,
			},
			"last_login_at": schema.StringAttribute{
				Description: "Key's most recent login time via SFTP",
				Computed:    true,
			},
			"private_key": schema.StringAttribute{
				Description: "Only returned when generating keys. Private key generated for the user.",
				Computed:    true,
			},
			"public_key": schema.StringAttribute{
				Description: "Only returned when generating keys. Public key generated for the user.",
				Computed:    true,
			},
			"username": schema.StringAttribute{
				Description: "Username of the user this public key is associated with",
				Computed:    true,
			},
			"user_id": schema.Int64Attribute{
				Description: "User ID this public key is associated with",
				Computed:    true,
			},
		},
	}
}

func (r *publicKeyDataSource) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
	var data publicKeyDataSourceModel
	diags := req.Config.Get(ctx, &data)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	paramsPublicKeyFind := files_sdk.PublicKeyFindParams{}
	paramsPublicKeyFind.Id = data.Id.ValueInt64()

	publicKey, err := r.client.Find(paramsPublicKeyFind, files_sdk.WithContext(ctx))
	if err != nil {
		resp.Diagnostics.AddError(
			"Error Reading Files PublicKey",
			"Could not read public_key id "+fmt.Sprint(data.Id.ValueInt64())+": "+err.Error(),
		)
		return
	}

	diags = r.populateDataSourceModel(ctx, publicKey, &data)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	diags = resp.State.Set(ctx, data)
	resp.Diagnostics.Append(diags...)
}

func (r *publicKeyDataSource) populateDataSourceModel(ctx context.Context, publicKey files_sdk.PublicKey, state *publicKeyDataSourceModel) (diags diag.Diagnostics) {
	state.Id = types.Int64Value(publicKey.Id)
	state.Title = types.StringValue(publicKey.Title)
	if err := lib.TimeToStringType(ctx, path.Root("created_at"), publicKey.CreatedAt, &state.CreatedAt); err != nil {
		diags.AddError(
			"Error Creating Files PublicKey",
			"Could not convert state created_at to string: "+err.Error(),
		)
	}
	state.Fingerprint = types.StringValue(publicKey.Fingerprint)
	state.FingerprintSha256 = types.StringValue(publicKey.FingerprintSha256)
	state.Status = types.StringValue(publicKey.Status)
	if err := lib.TimeToStringType(ctx, path.Root("last_login_at"), publicKey.LastLoginAt, &state.LastLoginAt); err != nil {
		diags.AddError(
			"Error Creating Files PublicKey",
			"Could not convert state last_login_at to string: "+err.Error(),
		)
	}
	state.PrivateKey = types.StringValue(publicKey.PrivateKey)
	state.PublicKey = types.StringValue(publicKey.PublicKey)
	state.Username = types.StringValue(publicKey.Username)
	state.UserId = types.Int64Value(publicKey.UserId)

	return
}
