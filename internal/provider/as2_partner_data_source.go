package provider

import (
	"context"
	"fmt"

	files_sdk "github.com/Files-com/files-sdk-go/v3"
	as2_partner "github.com/Files-com/files-sdk-go/v3/as2partner"
	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/datasource/schema"
	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

var (
	_ datasource.DataSource              = &as2PartnerDataSource{}
	_ datasource.DataSourceWithConfigure = &as2PartnerDataSource{}
)

func NewAs2PartnerDataSource() datasource.DataSource {
	return &as2PartnerDataSource{}
}

type as2PartnerDataSource struct {
	client *as2_partner.Client
}

type as2PartnerDataSourceModel struct {
	Id                         types.Int64  `tfsdk:"id"`
	As2StationId               types.Int64  `tfsdk:"as2_station_id"`
	Name                       types.String `tfsdk:"name"`
	Uri                        types.String `tfsdk:"uri"`
	ServerCertificate          types.String `tfsdk:"server_certificate"`
	HttpAuthUsername           types.String `tfsdk:"http_auth_username"`
	MdnValidationLevel         types.String `tfsdk:"mdn_validation_level"`
	EnableDedicatedIps         types.Bool   `tfsdk:"enable_dedicated_ips"`
	HexPublicCertificateSerial types.String `tfsdk:"hex_public_certificate_serial"`
	PublicCertificateMd5       types.String `tfsdk:"public_certificate_md5"`
	PublicCertificateSubject   types.String `tfsdk:"public_certificate_subject"`
	PublicCertificateIssuer    types.String `tfsdk:"public_certificate_issuer"`
	PublicCertificateSerial    types.String `tfsdk:"public_certificate_serial"`
	PublicCertificateNotBefore types.String `tfsdk:"public_certificate_not_before"`
	PublicCertificateNotAfter  types.String `tfsdk:"public_certificate_not_after"`
}

func (r *as2PartnerDataSource) Configure(_ context.Context, req datasource.ConfigureRequest, resp *datasource.ConfigureResponse) {
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

	r.client = &as2_partner.Client{Config: sdk_config}
}

func (r *as2PartnerDataSource) Metadata(_ context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_as2_partner"
}

func (r *as2PartnerDataSource) Schema(_ context.Context, _ datasource.SchemaRequest, resp *datasource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Description: "An AS2Partner is a counterparty of the Files.com site's AS2 connectivity. Generally you will have one AS2 Partner created for each counterparty with whom you send and/or receive files via AS2.",
		Attributes: map[string]schema.Attribute{
			"id": schema.Int64Attribute{
				Description: "ID of the AS2 Partner.",
				Required:    true,
			},
			"as2_station_id": schema.Int64Attribute{
				Description: "ID of the AS2 Station associated with this partner.",
				Computed:    true,
			},
			"name": schema.StringAttribute{
				Description: "The partner's formal AS2 name.",
				Computed:    true,
			},
			"uri": schema.StringAttribute{
				Description: "Public URI where we will send the AS2 messages (via HTTP/HTTPS).",
				Computed:    true,
			},
			"server_certificate": schema.StringAttribute{
				Description: "Should we require that the remote HTTP server have a valid SSL Certificate for HTTPS?",
				Computed:    true,
			},
			"http_auth_username": schema.StringAttribute{
				Description: "Username to send to server for HTTP Authentication.",
				Computed:    true,
			},
			"mdn_validation_level": schema.StringAttribute{
				Description: "How should Files.com evaluate message transfer success based on a partner's MDN response?  This setting does not affect MDN storage; all MDNs received from a partner are always stored. `none`: MDN is stored for informational purposes only, a successful HTTPS transfer is a successful AS2 transfer. `weak`: Inspect the MDN for MIC and Disposition only. `normal`: `weak` plus validate MDN signature matches body, `strict`: `normal` but do not allow signatures from self-signed or incorrectly purposed certificates.",
				Computed:    true,
			},
			"enable_dedicated_ips": schema.BoolAttribute{
				Description: "If `true`, we will use your site's dedicated IPs for all outbound connections to this AS2 PArtner.",
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
		},
	}
}

func (r *as2PartnerDataSource) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
	var data as2PartnerDataSourceModel
	diags := req.Config.Get(ctx, &data)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	paramsAs2PartnerFind := files_sdk.As2PartnerFindParams{}
	paramsAs2PartnerFind.Id = data.Id.ValueInt64()

	as2Partner, err := r.client.Find(paramsAs2PartnerFind, files_sdk.WithContext(ctx))
	if err != nil {
		resp.Diagnostics.AddError(
			"Error Reading Files As2Partner",
			"Could not read as2_partner id "+fmt.Sprint(data.Id.ValueInt64())+": "+err.Error(),
		)
		return
	}

	diags = r.populateDataSourceModel(ctx, as2Partner, &data)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	diags = resp.State.Set(ctx, data)
	resp.Diagnostics.Append(diags...)
}

func (r *as2PartnerDataSource) populateDataSourceModel(ctx context.Context, as2Partner files_sdk.As2Partner, state *as2PartnerDataSourceModel) (diags diag.Diagnostics) {
	state.Id = types.Int64Value(as2Partner.Id)
	state.As2StationId = types.Int64Value(as2Partner.As2StationId)
	state.Name = types.StringValue(as2Partner.Name)
	state.Uri = types.StringValue(as2Partner.Uri)
	state.ServerCertificate = types.StringValue(as2Partner.ServerCertificate)
	state.HttpAuthUsername = types.StringValue(as2Partner.HttpAuthUsername)
	state.MdnValidationLevel = types.StringValue(as2Partner.MdnValidationLevel)
	state.EnableDedicatedIps = types.BoolPointerValue(as2Partner.EnableDedicatedIps)
	state.HexPublicCertificateSerial = types.StringValue(as2Partner.HexPublicCertificateSerial)
	state.PublicCertificateMd5 = types.StringValue(as2Partner.PublicCertificateMd5)
	state.PublicCertificateSubject = types.StringValue(as2Partner.PublicCertificateSubject)
	state.PublicCertificateIssuer = types.StringValue(as2Partner.PublicCertificateIssuer)
	state.PublicCertificateSerial = types.StringValue(as2Partner.PublicCertificateSerial)
	state.PublicCertificateNotBefore = types.StringValue(as2Partner.PublicCertificateNotBefore)
	state.PublicCertificateNotAfter = types.StringValue(as2Partner.PublicCertificateNotAfter)

	return
}
