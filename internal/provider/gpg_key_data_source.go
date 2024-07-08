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
	Id                     types.Int64  `tfsdk:"id"`
	ExpiresAt              types.String `tfsdk:"expires_at"`
	Name                   types.String `tfsdk:"name"`
	UserId                 types.Int64  `tfsdk:"user_id"`
	PublicKey              types.String `tfsdk:"public_key"`
	PublicKeyHash          types.String `tfsdk:"public_key_hash"`
	PrivateKey             types.String `tfsdk:"private_key"`
	PrivateKeyHash         types.String `tfsdk:"private_key_hash"`
	PrivateKeyPassword     types.String `tfsdk:"private_key_password"`
	PrivateKeyPasswordHash types.String `tfsdk:"private_key_password_hash"`
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
		Description: "GPG keys for decrypt or encrypt behaviors.",
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
			"user_id": schema.Int64Attribute{
				Description: "GPG owner's user id",
				Computed:    true,
			},
			"public_key": schema.StringAttribute{
				Description: "Your GPG public key",
				Computed:    true,
			},
			"public_key_hash": schema.StringAttribute{
				Computed: true,
			},
			"private_key": schema.StringAttribute{
				Description: "Your GPG private key.",
				Computed:    true,
			},
			"private_key_hash": schema.StringAttribute{
				Computed: true,
			},
			"private_key_password": schema.StringAttribute{
				Description: "Your GPG private key password. Only required for password protected keys.",
				Computed:    true,
			},
			"private_key_password_hash": schema.StringAttribute{
				Computed: true,
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
	state.UserId = types.Int64Value(gpgKey.UserId)
	state.PublicKeyHash = types.StringValue(gpgKey.PublicKey)
	state.PrivateKeyHash = types.StringValue(gpgKey.PrivateKey)
	state.PrivateKeyPasswordHash = types.StringValue(gpgKey.PrivateKeyPassword)

	return
}
