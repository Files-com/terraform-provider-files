package provider

import (
	"context"
	"fmt"

	files_sdk "github.com/Files-com/files-sdk-go/v3"
	gpg_key "github.com/Files-com/files-sdk-go/v3/gpgkey"
	"github.com/Files-com/terraform-provider-files/lib"
	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/datasource/schema"
	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

var (
	_ datasource.DataSource              = &gpgKeyDataSource{}
	_ datasource.DataSourceWithConfigure = &gpgKeyDataSource{}
)

func NewGpgKeyDataSource() datasource.DataSource {
	return &gpgKeyDataSource{}
}

type gpgKeyDataSource struct {
	client *gpg_key.Client
}

type gpgKeyDataSourceModel struct {
	Id                    types.Int64  `tfsdk:"id"`
	ExpiresAt             types.String `tfsdk:"expires_at"`
	Name                  types.String `tfsdk:"name"`
	PartnerId             types.Int64  `tfsdk:"partner_id"`
	PartnerName           types.String `tfsdk:"partner_name"`
	UserId                types.Int64  `tfsdk:"user_id"`
	PublicKeyMd5          types.String `tfsdk:"public_key_md5"`
	PrivateKeyMd5         types.String `tfsdk:"private_key_md5"`
	GeneratedPublicKey    types.String `tfsdk:"generated_public_key"`
	GeneratedPrivateKey   types.String `tfsdk:"generated_private_key"`
	PrivateKeyPasswordMd5 types.String `tfsdk:"private_key_password_md5"`
}

func (r *gpgKeyDataSource) Configure(_ context.Context, req datasource.ConfigureRequest, resp *datasource.ConfigureResponse) {
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

	r.client = &gpg_key.Client{Config: sdk_config}
}

func (r *gpgKeyDataSource) Metadata(_ context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_gpg_key"
}

func (r *gpgKeyDataSource) Schema(_ context.Context, _ datasource.SchemaRequest, resp *datasource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Description: "A GPGKey object on Files.com is used to securely store both the private and public keys associated with a GPG (GNU Privacy Guard) encryption key pair. This object enables the encryption and decryption of data using GPG, allowing you to protect sensitive information.\n\n\n\nThe private key is kept confidential and is used for decrypting data or signing messages to prove authenticity, while the public key is used to encrypt messages that only the owner of the private key can decrypt.\n\n\n\nBy storing both keys together in a GPGKey object, Files.com makes it easier to understand encryption operations, ensuring secure and efficient handling of encrypted data within the platform.",
		Attributes: map[string]schema.Attribute{
			"id": schema.Int64Attribute{
				Description: "Your GPG key ID.",
				Required:    true,
			},
			"expires_at": schema.StringAttribute{
				Description: "Your GPG key expiration date.",
				Computed:    true,
			},
			"name": schema.StringAttribute{
				Description: "Your GPG key name.",
				Computed:    true,
			},
			"partner_id": schema.Int64Attribute{
				Description: "Partner ID who owns this GPG Key, if applicable.",
				Computed:    true,
			},
			"partner_name": schema.StringAttribute{
				Description: "Name of the Partner who owns this GPG Key, if applicable.",
				Computed:    true,
			},
			"user_id": schema.Int64Attribute{
				Description: "User ID who owns this GPG Key, if applicable.",
				Computed:    true,
			},
			"public_key_md5": schema.StringAttribute{
				Description: "MD5 hash of your GPG public key",
				Computed:    true,
			},
			"private_key_md5": schema.StringAttribute{
				Description: "MD5 hash of your GPG private key.",
				Computed:    true,
			},
			"generated_public_key": schema.StringAttribute{
				Description: "Your GPG public key",
				Computed:    true,
			},
			"generated_private_key": schema.StringAttribute{
				Description: "Your GPG private key.",
				Computed:    true,
			},
			"private_key_password_md5": schema.StringAttribute{
				Description: "Your GPG private key password. Only required for password protected keys.",
				Computed:    true,
			},
		},
	}
}

func (r *gpgKeyDataSource) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
	var data gpgKeyDataSourceModel
	diags := req.Config.Get(ctx, &data)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	paramsGpgKeyFind := files_sdk.GpgKeyFindParams{}
	paramsGpgKeyFind.Id = data.Id.ValueInt64()

	gpgKey, err := r.client.Find(paramsGpgKeyFind, files_sdk.WithContext(ctx))
	if err != nil {
		resp.Diagnostics.AddError(
			"Error Reading Files GpgKey",
			"Could not read gpg_key id "+fmt.Sprint(data.Id.ValueInt64())+": "+err.Error(),
		)
		return
	}

	diags = r.populateDataSourceModel(ctx, gpgKey, &data)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	diags = resp.State.Set(ctx, data)
	resp.Diagnostics.Append(diags...)
}

func (r *gpgKeyDataSource) populateDataSourceModel(ctx context.Context, gpgKey files_sdk.GpgKey, state *gpgKeyDataSourceModel) (diags diag.Diagnostics) {
	state.Id = types.Int64Value(gpgKey.Id)
	if err := lib.TimeToStringType(ctx, path.Root("expires_at"), gpgKey.ExpiresAt, &state.ExpiresAt); err != nil {
		diags.AddError(
			"Error Creating Files GpgKey",
			"Could not convert state expires_at to string: "+err.Error(),
		)
	}
	state.Name = types.StringValue(gpgKey.Name)
	state.PartnerId = types.Int64Value(gpgKey.PartnerId)
	state.PartnerName = types.StringValue(gpgKey.PartnerName)
	state.UserId = types.Int64Value(gpgKey.UserId)
	state.PublicKeyMd5 = types.StringValue(gpgKey.PublicKeyMd5)
	state.PrivateKeyMd5 = types.StringValue(gpgKey.PrivateKeyMd5)
	state.GeneratedPublicKey = types.StringValue(gpgKey.GeneratedPublicKey)
	state.GeneratedPrivateKey = types.StringValue(gpgKey.GeneratedPrivateKey)
	state.PrivateKeyPasswordMd5 = types.StringValue(gpgKey.PrivateKeyPasswordMd5)

	return
}
