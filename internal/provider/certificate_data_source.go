package provider

import (
	"context"
	"fmt"

	files_sdk "github.com/Files-com/files-sdk-go/v3"
	certificate "github.com/Files-com/files-sdk-go/v3/certificate"
	"github.com/Files-com/terraform-provider-files/lib"

	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/datasource/schema"
	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/path"

	"github.com/hashicorp/terraform-plugin-framework/types"
)

var (
	_ datasource.DataSource              = &certificateDataSource{}
	_ datasource.DataSourceWithConfigure = &certificateDataSource{}
)

func NewCertificateDataSource() datasource.DataSource {
	return &certificateDataSource{}
}

type certificateDataSource struct {
	client *certificate.Client
}

type certificateDataSourceModel struct {
	Id             types.Int64  `tfsdk:"id"`
	Name           types.String `tfsdk:"name"`
	CustomDomainId types.Int64  `tfsdk:"custom_domain_id"`
	Certificate    types.String `tfsdk:"certificate"`
	CreatedAt      types.String `tfsdk:"created_at"`
	DisplayStatus  types.String `tfsdk:"display_status"`
	Domains        types.List   `tfsdk:"domains"`
	ExpiresAt      types.String `tfsdk:"expires_at"`
	BrickManaged   types.Bool   `tfsdk:"brick_managed"`
	Intermediates  types.String `tfsdk:"intermediates"`
	IpAddresses    types.List   `tfsdk:"ip_addresses"`
	Issuer         types.String `tfsdk:"issuer"`
	KeyType        types.String `tfsdk:"key_type"`
	Request        types.String `tfsdk:"request"`
	Status         types.String `tfsdk:"status"`
	Subject        types.String `tfsdk:"subject"`
}

func (r *certificateDataSource) Configure(_ context.Context, req datasource.ConfigureRequest, resp *datasource.ConfigureResponse) {
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

	r.client = &certificate.Client{Config: sdk_config}
}

func (r *certificateDataSource) Metadata(_ context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_certificate"
}

func (r *certificateDataSource) Schema(_ context.Context, _ datasource.SchemaRequest, resp *datasource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Description: "A Certificate is a way to store your custom SSL certificate for your Files.com site.\nFiles.com is also capable of registering SSL certificates itself, so this is only required if you have specific security requirements to use your own certificate.",
		Attributes: map[string]schema.Attribute{
			"id": schema.Int64Attribute{
				Description: "Certificate ID",
				Required:    true,
			},
			"name": schema.StringAttribute{
				Description: "Descriptive name of certificate",
				Computed:    true,
			},
			"custom_domain_id": schema.Int64Attribute{
				Description: "Custom Domain that owns this certificate, including while the certificate is inactive.",
				Computed:    true,
			},
			"certificate": schema.StringAttribute{
				Description: "Full text of SSL certificate",
				Computed:    true,
			},
			"created_at": schema.StringAttribute{
				Description: "Certificate created at",
				Computed:    true,
			},
			"display_status": schema.StringAttribute{
				Description: "Certificate status. (One of `Request Generated`, `New`, `Active`, `Active/Expired`, `Expired`, or `Available`)",
				Computed:    true,
			},
			"domains": schema.ListAttribute{
				Description: "Domains on this certificate other than files.com domains",
				Computed:    true,
				ElementType: types.StringType,
			},
			"expires_at": schema.StringAttribute{
				Description: "Certificate expiration date/time",
				Computed:    true,
			},
			"brick_managed": schema.BoolAttribute{
				Description: "Is this certificate automatically managed and renewed by Files.com?",
				Computed:    true,
			},
			"intermediates": schema.StringAttribute{
				Description: "Intermediate certificates",
				Computed:    true,
			},
			"ip_addresses": schema.ListAttribute{
				Description: "A list of IP addresses associated with this SSL Certificate",
				Computed:    true,
				ElementType: types.StringType,
			},
			"issuer": schema.StringAttribute{
				Description: "X509 issuer",
				Computed:    true,
			},
			"key_type": schema.StringAttribute{
				Description: "Type of key",
				Computed:    true,
			},
			"request": schema.StringAttribute{
				Description: "Certificate signing request text",
				Computed:    true,
			},
			"status": schema.StringAttribute{
				Description: "Certificate status (Request Generated, New, Active, Active/Expired, Expired, or Available)",
				Computed:    true,
			},
			"subject": schema.StringAttribute{
				Description: "X509 Subject name",
				Computed:    true,
			},
		},
	}
}

func (r *certificateDataSource) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
	var data certificateDataSourceModel
	diags := req.Config.Get(ctx, &data)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	paramsCertificateFind := files_sdk.CertificateFindParams{}
	paramsCertificateFind.Id = data.Id.ValueInt64()

	certificate, err := r.client.Find(paramsCertificateFind, files_sdk.WithContext(ctx))
	if err != nil {
		resp.Diagnostics.AddError(
			"Error Reading Files Certificate",
			"Could not read certificate id "+fmt.Sprint(data.Id.ValueInt64())+": "+err.Error(),
		)
		return
	}

	diags = r.populateDataSourceModel(ctx, certificate, &data)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	diags = resp.State.Set(ctx, data)
	resp.Diagnostics.Append(diags...)
}

func (r *certificateDataSource) populateDataSourceModel(ctx context.Context, certificate files_sdk.Certificate, state *certificateDataSourceModel) (diags diag.Diagnostics) {
	var propDiags diag.Diagnostics

	state.Id = types.Int64Value(certificate.Id)
	state.Name = types.StringValue(certificate.Name)
	state.CustomDomainId = types.Int64Value(certificate.CustomDomainId)
	state.Certificate = types.StringValue(certificate.Certificate)
	if err := lib.TimeToStringType(ctx, path.Root("created_at"), certificate.CreatedAt, &state.CreatedAt); err != nil {
		diags.AddError(
			"Error Creating Files Certificate",
			"Could not convert state created_at to string: "+err.Error(),
		)
	}
	state.DisplayStatus = types.StringValue(certificate.DisplayStatus)
	state.Domains, propDiags = types.ListValueFrom(ctx, types.StringType, certificate.Domains)
	diags.Append(propDiags...)
	if err := lib.TimeToStringType(ctx, path.Root("expires_at"), certificate.ExpiresAt, &state.ExpiresAt); err != nil {
		diags.AddError(
			"Error Creating Files Certificate",
			"Could not convert state expires_at to string: "+err.Error(),
		)
	}
	state.BrickManaged = types.BoolPointerValue(certificate.BrickManaged)
	state.Intermediates = types.StringValue(certificate.Intermediates)
	state.IpAddresses, propDiags = types.ListValueFrom(ctx, types.StringType, certificate.IpAddresses)
	diags.Append(propDiags...)
	state.Issuer = types.StringValue(certificate.Issuer)
	state.KeyType = types.StringValue(certificate.KeyType)
	state.Request = types.StringValue(certificate.Request)
	state.Status = types.StringValue(certificate.Status)
	state.Subject = types.StringValue(certificate.Subject)

	return
}
