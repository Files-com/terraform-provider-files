package provider

import (
	"context"
	"fmt"
	"strconv"
	"strings"

	files_sdk "github.com/Files-com/files-sdk-go/v3"
	certificate "github.com/Files-com/files-sdk-go/v3/certificate"
	"github.com/Files-com/terraform-provider-files/lib"

	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"

	"github.com/hashicorp/terraform-plugin-framework/resource/schema/int64planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/listplanmodifier"

	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"

	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

var (
	_ resource.Resource                = &certificateResource{}
	_ resource.ResourceWithConfigure   = &certificateResource{}
	_ resource.ResourceWithImportState = &certificateResource{}
)

func NewCertificateResource() resource.Resource {
	return &certificateResource{}
}

type certificateResource struct {
	client *certificate.Client
}

type certificateResourceModel struct {
	Name                       types.String `tfsdk:"name"`
	Certificate                types.String `tfsdk:"certificate"`
	Intermediates              types.String `tfsdk:"intermediates"`
	KeyType                    types.String `tfsdk:"key_type"`
	CertificateDomain          types.String `tfsdk:"certificate_domain"`
	CertificateExtraDomains    types.List   `tfsdk:"certificate_extra_domains"`
	CertificateCountry         types.String `tfsdk:"certificate_country"`
	CertificateStateOrProvince types.String `tfsdk:"certificate_state_or_province"`
	CertificateCityOrLocale    types.String `tfsdk:"certificate_city_or_locale"`
	CertificateCompanyName     types.String `tfsdk:"certificate_company_name"`
	CsrOu1                     types.String `tfsdk:"csr_ou1"`
	CsrOu2                     types.String `tfsdk:"csr_ou2"`
	CsrOu3                     types.String `tfsdk:"csr_ou3"`
	CertificateEmailAddress    types.String `tfsdk:"certificate_email_address"`
	PrivateKey                 types.String `tfsdk:"private_key"`
	Password                   types.String `tfsdk:"password"`
	Id                         types.Int64  `tfsdk:"id"`
	CustomDomainId             types.Int64  `tfsdk:"custom_domain_id"`
	CreatedAt                  types.String `tfsdk:"created_at"`
	DisplayStatus              types.String `tfsdk:"display_status"`
	Domains                    types.List   `tfsdk:"domains"`
	ExpiresAt                  types.String `tfsdk:"expires_at"`
	BrickManaged               types.Bool   `tfsdk:"brick_managed"`
	IpAddresses                types.List   `tfsdk:"ip_addresses"`
	Issuer                     types.String `tfsdk:"issuer"`
	Request                    types.String `tfsdk:"request"`
	Status                     types.String `tfsdk:"status"`
	Subject                    types.String `tfsdk:"subject"`
}

func (r *certificateResource) Configure(_ context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {
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

func (r *certificateResource) Metadata(_ context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_certificate"
}

func (r *certificateResource) Schema(_ context.Context, _ resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Description: "A Certificate is a way to store your custom SSL certificate for your Files.com site.\nFiles.com is also capable of registering SSL certificates itself, so this is only required if you have specific security requirements to use your own certificate.",
		Attributes: map[string]schema.Attribute{
			"name": schema.StringAttribute{
				Description: "Descriptive name of certificate",
				Required:    true,
			},
			"certificate": schema.StringAttribute{
				Description: "Full text of SSL certificate",
				Computed:    true,
				Optional:    true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
			},
			"intermediates": schema.StringAttribute{
				Description: "Intermediate certificates",
				Computed:    true,
				Optional:    true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
			},
			"key_type": schema.StringAttribute{
				Description: "Type of key",
				Computed:    true,
				Optional:    true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
					stringplanmodifier.RequiresReplace(),
				},
			},
			"certificate_domain": schema.StringAttribute{
				Description: "Domain for certificate.",
				Optional:    true,
				WriteOnly:   true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.RequiresReplace(),
				},
			},
			"certificate_extra_domains": schema.ListAttribute{
				Description: "Additional domains for the certificate CSR.",
				Optional:    true,
				WriteOnly:   true,
				ElementType: types.StringType,
				PlanModifiers: []planmodifier.List{
					listplanmodifier.RequiresReplace(),
				},
			},
			"certificate_country": schema.StringAttribute{
				Description: "Country.",
				Optional:    true,
				WriteOnly:   true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.RequiresReplace(),
				},
			},
			"certificate_state_or_province": schema.StringAttribute{
				Description: "State or province.",
				Optional:    true,
				WriteOnly:   true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.RequiresReplace(),
				},
			},
			"certificate_city_or_locale": schema.StringAttribute{
				Description: "City or locale.",
				Optional:    true,
				WriteOnly:   true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.RequiresReplace(),
				},
			},
			"certificate_company_name": schema.StringAttribute{
				Description: "Company name.",
				Optional:    true,
				WriteOnly:   true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.RequiresReplace(),
				},
			},
			"csr_ou1": schema.StringAttribute{
				Description: "Department name or organization unit 1",
				Optional:    true,
				WriteOnly:   true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.RequiresReplace(),
				},
			},
			"csr_ou2": schema.StringAttribute{
				Description: "Department name or organization unit 2",
				Optional:    true,
				WriteOnly:   true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.RequiresReplace(),
				},
			},
			"csr_ou3": schema.StringAttribute{
				Description: "Department name or organization unit 3",
				Optional:    true,
				WriteOnly:   true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.RequiresReplace(),
				},
			},
			"certificate_email_address": schema.StringAttribute{
				Description: "Email address for the certificate owner.",
				Optional:    true,
				WriteOnly:   true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.RequiresReplace(),
				},
			},
			"private_key": schema.StringAttribute{
				Description: "Certificate private key.",
				Optional:    true,
				WriteOnly:   true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.RequiresReplace(),
				},
			},
			"password": schema.StringAttribute{
				Description: "Certificate password.",
				Optional:    true,
				WriteOnly:   true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.RequiresReplace(),
				},
			},
			"id": schema.Int64Attribute{
				Description: "Certificate ID",
				Computed:    true,
				PlanModifiers: []planmodifier.Int64{
					int64planmodifier.UseStateForUnknown(),
				},
			},
			"custom_domain_id": schema.Int64Attribute{
				Description: "Custom Domain that owns this certificate, including while the certificate is inactive.",
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
			"ip_addresses": schema.ListAttribute{
				Description: "A list of IP addresses associated with this SSL Certificate",
				Computed:    true,
				ElementType: types.StringType,
			},
			"issuer": schema.StringAttribute{
				Description: "X509 issuer",
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

func (r *certificateResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	var plan certificateResourceModel
	diags := req.Plan.Get(ctx, &plan)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}
	var config certificateResourceModel
	diags = req.Config.Get(ctx, &config)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	paramsCertificateCreate := files_sdk.CertificateCreateParams{}
	paramsCertificateCreate.Name = plan.Name.ValueString()
	paramsCertificateCreate.CertificateDomain = config.CertificateDomain.ValueString()
	if !config.CertificateExtraDomains.IsNull() && !config.CertificateExtraDomains.IsUnknown() {
		diags = config.CertificateExtraDomains.ElementsAs(ctx, &paramsCertificateCreate.CertificateExtraDomains, false)
		resp.Diagnostics.Append(diags...)
	}
	paramsCertificateCreate.CertificateCountry = config.CertificateCountry.ValueString()
	paramsCertificateCreate.CertificateStateOrProvince = config.CertificateStateOrProvince.ValueString()
	paramsCertificateCreate.CertificateCityOrLocale = config.CertificateCityOrLocale.ValueString()
	paramsCertificateCreate.CertificateCompanyName = config.CertificateCompanyName.ValueString()
	paramsCertificateCreate.CsrOu1 = config.CsrOu1.ValueString()
	paramsCertificateCreate.CsrOu2 = config.CsrOu2.ValueString()
	paramsCertificateCreate.CsrOu3 = config.CsrOu3.ValueString()
	paramsCertificateCreate.CertificateEmailAddress = config.CertificateEmailAddress.ValueString()
	paramsCertificateCreate.KeyType = plan.KeyType.ValueString()
	paramsCertificateCreate.Certificate = plan.Certificate.ValueString()
	paramsCertificateCreate.PrivateKey = config.PrivateKey.ValueString()
	paramsCertificateCreate.Password = config.Password.ValueString()
	paramsCertificateCreate.Intermediates = plan.Intermediates.ValueString()

	if resp.Diagnostics.HasError() {
		return
	}

	certificate, err := r.client.Create(paramsCertificateCreate, files_sdk.WithContext(ctx))
	if err != nil {
		resp.Diagnostics.AddError(
			"Error Creating Files Certificate",
			"Could not create certificate, unexpected error: "+err.Error(),
		)
		return
	}

	diags = r.populateResourceModel(ctx, certificate, &plan)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	diags = resp.State.Set(ctx, plan)
	resp.Diagnostics.Append(diags...)
}

func (r *certificateResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	var state certificateResourceModel
	diags := req.State.Get(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	paramsCertificateFind := files_sdk.CertificateFindParams{}
	paramsCertificateFind.Id = state.Id.ValueInt64()

	certificate, err := r.client.Find(paramsCertificateFind, files_sdk.WithContext(ctx))
	if err != nil {
		if files_sdk.IsNotExist(err) {
			resp.State.RemoveResource(ctx)
			return
		}

		resp.Diagnostics.AddError(
			"Error Reading Files Certificate",
			"Could not read certificate id "+fmt.Sprint(state.Id.ValueInt64())+": "+err.Error(),
		)
		return
	}

	diags = r.populateResourceModel(ctx, certificate, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	diags = resp.State.Set(ctx, state)
	resp.Diagnostics.Append(diags...)
}

func (r *certificateResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	var plan certificateResourceModel
	diags := req.Plan.Get(ctx, &plan)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}
	var config certificateResourceModel
	diags = req.Config.Get(ctx, &config)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	paramsCertificateUpdate := map[string]interface{}{}
	if !plan.Id.IsNull() && !plan.Id.IsUnknown() {
		paramsCertificateUpdate["id"] = plan.Id.ValueInt64()
	}
	if !config.Name.IsNull() && !config.Name.IsUnknown() {
		paramsCertificateUpdate["name"] = config.Name.ValueString()
	}
	if !config.Intermediates.IsNull() && !config.Intermediates.IsUnknown() {
		paramsCertificateUpdate["intermediates"] = config.Intermediates.ValueString()
	}
	if !config.Certificate.IsNull() && !config.Certificate.IsUnknown() {
		paramsCertificateUpdate["certificate"] = config.Certificate.ValueString()
	}

	if resp.Diagnostics.HasError() {
		return
	}

	certificate, err := r.client.UpdateWithMap(paramsCertificateUpdate, files_sdk.WithContext(ctx))
	if err != nil {
		resp.Diagnostics.AddError(
			"Error Updating Files Certificate",
			"Could not update certificate, unexpected error: "+err.Error(),
		)
		return
	}

	diags = r.populateResourceModel(ctx, certificate, &plan)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	diags = resp.State.Set(ctx, plan)
	resp.Diagnostics.Append(diags...)
}

func (r *certificateResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	var state certificateResourceModel
	diags := req.State.Get(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	paramsCertificateDelete := files_sdk.CertificateDeleteParams{}
	paramsCertificateDelete.Id = state.Id.ValueInt64()

	err := r.client.Delete(paramsCertificateDelete, files_sdk.WithContext(ctx))
	if err != nil && !files_sdk.IsNotExist(err) {
		resp.Diagnostics.AddError(
			"Error Deleting Files Certificate",
			"Could not delete certificate id "+fmt.Sprint(state.Id.ValueInt64())+": "+err.Error(),
		)
	}
}

func (r *certificateResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	idParts := strings.SplitN(req.ID, ",", 1)

	if len(idParts) != 1 || idParts[0] == "" {
		resp.Diagnostics.AddError(
			"Unexpected Import Identifier",
			fmt.Sprintf("Expected import identifier with format: id. Got: %q", req.ID),
		)
		return
	}

	idPart, err := strconv.ParseFloat(idParts[0], 64)
	if err != nil {
		resp.Diagnostics.AddError(
			"Error Parsing ID",
			"Could not parse id: "+err.Error(),
		)
		return
	}
	resp.Diagnostics.Append(resp.State.SetAttribute(ctx, path.Root("id"), idPart)...)

}

func (r *certificateResource) populateResourceModel(ctx context.Context, certificate files_sdk.Certificate, state *certificateResourceModel) (diags diag.Diagnostics) {
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
