package provider

import (
	"context"
	"fmt"

	files_sdk "github.com/Files-com/files-sdk-go/v3"
	as2_station "github.com/Files-com/files-sdk-go/v3/as2station"
	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/datasource/schema"
	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

var (
	_ datasource.DataSource              = &as2StationDataSource{}
	_ datasource.DataSourceWithConfigure = &as2StationDataSource{}
)

func NewAs2StationDataSource() datasource.DataSource {
	return &as2StationDataSource{}
}

type as2StationDataSource struct {
	client *as2_station.Client
}

type as2StationDataSourceModel struct {
	Id                         types.Int64  `tfsdk:"id"`
	WorkspaceId                types.Int64  `tfsdk:"workspace_id"`
	Name                       types.String `tfsdk:"name"`
	Uri                        types.String `tfsdk:"uri"`
	Domain                     types.String `tfsdk:"domain"`
	HexPublicCertificateSerial types.String `tfsdk:"hex_public_certificate_serial"`
	PublicCertificateMd5       types.String `tfsdk:"public_certificate_md5"`
	PublicCertificate          types.String `tfsdk:"public_certificate"`
	PrivateKeyMd5              types.String `tfsdk:"private_key_md5"`
	PublicCertificateSubject   types.String `tfsdk:"public_certificate_subject"`
	PublicCertificateIssuer    types.String `tfsdk:"public_certificate_issuer"`
	PublicCertificateSerial    types.String `tfsdk:"public_certificate_serial"`
	PublicCertificateNotBefore types.String `tfsdk:"public_certificate_not_before"`
	PublicCertificateNotAfter  types.String `tfsdk:"public_certificate_not_after"`
	PrivateKeyPasswordMd5      types.String `tfsdk:"private_key_password_md5"`
}

func (r *as2StationDataSource) Configure(_ context.Context, req datasource.ConfigureRequest, resp *datasource.ConfigureResponse) {
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

	r.client = &as2_station.Client{Config: sdk_config}
}

func (r *as2StationDataSource) Metadata(_ context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_as2_station"
}

func (r *as2StationDataSource) Schema(_ context.Context, _ datasource.SchemaRequest, resp *datasource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Description: "An AS2Station is a remote AS2 server that can send data into Files.com and receive data from Files.com.",
		Attributes: map[string]schema.Attribute{
			"id": schema.Int64Attribute{
				Description: "Id of the AS2 Station.",
				Required:    true,
			},
			"workspace_id": schema.Int64Attribute{
				Description: "ID of the Workspace associated with this AS2 Station.",
				Computed:    true,
			},
			"name": schema.StringAttribute{
				Description: "The station's formal AS2 name.",
				Computed:    true,
			},
			"uri": schema.StringAttribute{
				Description: "Public URI for sending AS2 message to.",
				Computed:    true,
			},
			"domain": schema.StringAttribute{
				Description: "The station's AS2 domain name.",
				Computed:    true,
			},
			"hex_public_certificate_serial": schema.StringAttribute{
				Description: "Serial of public certificate used for message security in hex format.",
				Computed:    true,
			},
			"public_certificate_md5": schema.StringAttribute{
				Description: "MD5 hash of public certificate used for message security.",
				Computed:    true,
			},
			"public_certificate": schema.StringAttribute{
				Description: "Public certificate used for message security.",
				Computed:    true,
			},
			"private_key_md5": schema.StringAttribute{
				Description: "MD5 hash of private key used for message security.",
				Computed:    true,
			},
			"public_certificate_subject": schema.StringAttribute{
				Description: "Subject of public certificate used for message security.",
				Computed:    true,
			},
			"public_certificate_issuer": schema.StringAttribute{
				Description: "Issuer of public certificate used for message security.",
				Computed:    true,
			},
			"public_certificate_serial": schema.StringAttribute{
				Description: "Serial of public certificate used for message security.",
				Computed:    true,
			},
			"public_certificate_not_before": schema.StringAttribute{
				Description: "Not before value of public certificate used for message security.",
				Computed:    true,
			},
			"public_certificate_not_after": schema.StringAttribute{
				Description: "Not after value of public certificate used for message security.",
				Computed:    true,
			},
			"private_key_password_md5": schema.StringAttribute{
				Description: "MD5 hash of private key password used for message security.",
				Computed:    true,
			},
		},
	}
}

func (r *as2StationDataSource) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
	var data as2StationDataSourceModel
	diags := req.Config.Get(ctx, &data)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	paramsAs2StationFind := files_sdk.As2StationFindParams{}
	paramsAs2StationFind.Id = data.Id.ValueInt64()

	as2Station, err := r.client.Find(paramsAs2StationFind, files_sdk.WithContext(ctx))
	if err != nil {
		resp.Diagnostics.AddError(
			"Error Reading Files As2Station",
			"Could not read as2_station id "+fmt.Sprint(data.Id.ValueInt64())+": "+err.Error(),
		)
		return
	}

	diags = r.populateDataSourceModel(ctx, as2Station, &data)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	diags = resp.State.Set(ctx, data)
	resp.Diagnostics.Append(diags...)
}

func (r *as2StationDataSource) populateDataSourceModel(ctx context.Context, as2Station files_sdk.As2Station, state *as2StationDataSourceModel) (diags diag.Diagnostics) {
	state.Id = types.Int64Value(as2Station.Id)
	state.WorkspaceId = types.Int64Value(as2Station.WorkspaceId)
	state.Name = types.StringValue(as2Station.Name)
	state.Uri = types.StringValue(as2Station.Uri)
	state.Domain = types.StringValue(as2Station.Domain)
	state.HexPublicCertificateSerial = types.StringValue(as2Station.HexPublicCertificateSerial)
	state.PublicCertificateMd5 = types.StringValue(as2Station.PublicCertificateMd5)
	state.PublicCertificate = types.StringValue(as2Station.PublicCertificate)
	state.PrivateKeyMd5 = types.StringValue(as2Station.PrivateKeyMd5)
	state.PublicCertificateSubject = types.StringValue(as2Station.PublicCertificateSubject)
	state.PublicCertificateIssuer = types.StringValue(as2Station.PublicCertificateIssuer)
	state.PublicCertificateSerial = types.StringValue(as2Station.PublicCertificateSerial)
	state.PublicCertificateNotBefore = types.StringValue(as2Station.PublicCertificateNotBefore)
	state.PublicCertificateNotAfter = types.StringValue(as2Station.PublicCertificateNotAfter)
	state.PrivateKeyPasswordMd5 = types.StringValue(as2Station.PrivateKeyPasswordMd5)

	return
}
