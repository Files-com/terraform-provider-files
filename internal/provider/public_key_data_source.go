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
		Description: "Public keys are used by Users who want to connect via SFTP/SSH.\n\n(Note that our SSH support is limited to file operations only, no shell is provided.)",
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
	state.Username = types.StringValue(publicKey.Username)
	state.UserId = types.Int64Value(publicKey.UserId)

	return
}
