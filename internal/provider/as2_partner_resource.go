package provider

import (
	"context"
	"fmt"
	"strconv"
	"strings"

	files_sdk "github.com/Files-com/files-sdk-go/v3"
	as2_partner "github.com/Files-com/files-sdk-go/v3/as2partner"
	"github.com/hashicorp/terraform-plugin-framework-validators/stringvalidator"
	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/boolplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/int64planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/schema/validator"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

var (
	_ resource.Resource                = &as2PartnerResource{}
	_ resource.ResourceWithConfigure   = &as2PartnerResource{}
	_ resource.ResourceWithImportState = &as2PartnerResource{}
)

func NewAs2PartnerResource() resource.Resource {
	return &as2PartnerResource{}
}

type as2PartnerResource struct {
	client *as2_partner.Client
}

type as2PartnerResourceModel struct {
	As2StationId               types.Int64  `tfsdk:"as2_station_id"`
	Name                       types.String `tfsdk:"name"`
	Uri                        types.String `tfsdk:"uri"`
	PublicCertificate          types.String `tfsdk:"public_certificate"`
	ServerCertificate          types.String `tfsdk:"server_certificate"`
	HttpAuthUsername           types.String `tfsdk:"http_auth_username"`
	MdnValidationLevel         types.String `tfsdk:"mdn_validation_level"`
	EnableDedicatedIps         types.Bool   `tfsdk:"enable_dedicated_ips"`
	HttpAuthPassword           types.String `tfsdk:"http_auth_password"`
	Id                         types.Int64  `tfsdk:"id"`
	HexPublicCertificateSerial types.String `tfsdk:"hex_public_certificate_serial"`
	PublicCertificateMd5       types.String `tfsdk:"public_certificate_md5"`
	PublicCertificateSubject   types.String `tfsdk:"public_certificate_subject"`
	PublicCertificateIssuer    types.String `tfsdk:"public_certificate_issuer"`
	PublicCertificateSerial    types.String `tfsdk:"public_certificate_serial"`
	PublicCertificateNotBefore types.String `tfsdk:"public_certificate_not_before"`
	PublicCertificateNotAfter  types.String `tfsdk:"public_certificate_not_after"`
}

func (r *as2PartnerResource) Configure(_ context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {
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

func (r *as2PartnerResource) Metadata(_ context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_as2_partner"
}

func (r *as2PartnerResource) Schema(_ context.Context, _ resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Description: "An AS2Partner is a counterparty of the Files.com site's AS2 connectivity. Generally you will have one AS2 Partner created for each counterparty with whom you send and/or receive files via AS2.",
		Attributes: map[string]schema.Attribute{
			"as2_station_id": schema.Int64Attribute{
				Description: "ID of the AS2 Station associated with this partner.",
				Required:    true,
				PlanModifiers: []planmodifier.Int64{
					int64planmodifier.RequiresReplace(),
				},
			},
			"name": schema.StringAttribute{
				Description: "The partner's formal AS2 name.",
				Required:    true,
			},
			"uri": schema.StringAttribute{
				Description: "Public URI where we will send the AS2 messages (via HTTP/HTTPS).",
				Required:    true,
			},
			"public_certificate": schema.StringAttribute{
				Description: "Public certificate for AS2 Partner.  Note: This is the certificate for AS2 message security, not a certificate used for HTTPS authentication.",
				Required:    true,
			},
			"server_certificate": schema.StringAttribute{
				Description: "Should we require that the remote HTTP server have a valid SSL Certificate for HTTPS?",
				Computed:    true,
				Optional:    true,
				Validators: []validator.String{
					stringvalidator.OneOf("require_match", "allow_any"),
				},
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
			},
			"http_auth_username": schema.StringAttribute{
				Description: "Username to send to server for HTTP Authentication.",
				Computed:    true,
				Optional:    true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
			},
			"mdn_validation_level": schema.StringAttribute{
				Description: "How should Files.com evaluate message transfer success based on a partner's MDN response?  This setting does not affect MDN storage; all MDNs received from a partner are always stored. `none`: MDN is stored for informational purposes only, a successful HTTPS transfer is a successful AS2 transfer. `weak`: Inspect the MDN for MIC and Disposition only. `normal`: `weak` plus validate MDN signature matches body, `strict`: `normal` but do not allow signatures from self-signed or incorrectly purposed certificates.",
				Computed:    true,
				Optional:    true,
				Validators: []validator.String{
					stringvalidator.OneOf("none", "weak", "normal", "strict"),
				},
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
			},
			"enable_dedicated_ips": schema.BoolAttribute{
				Description: "If `true`, we will use your site's dedicated IPs for all outbound connections to this AS2 PArtner.",
				Computed:    true,
				Optional:    true,
				PlanModifiers: []planmodifier.Bool{
					boolplanmodifier.UseStateForUnknown(),
				},
			},
			"http_auth_password": schema.StringAttribute{
				Description: "Password to send to server for HTTP Authentication.",
				Optional:    true,
			},
			"id": schema.Int64Attribute{
				Description: "ID of the AS2 Partner.",
				Computed:    true,
				PlanModifiers: []planmodifier.Int64{
					int64planmodifier.UseStateForUnknown(),
				},
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

func (r *as2PartnerResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	var plan as2PartnerResourceModel
	diags := req.Plan.Get(ctx, &plan)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	paramsAs2PartnerCreate := files_sdk.As2PartnerCreateParams{}
	if !plan.EnableDedicatedIps.IsNull() && !plan.EnableDedicatedIps.IsUnknown() {
		paramsAs2PartnerCreate.EnableDedicatedIps = plan.EnableDedicatedIps.ValueBoolPointer()
	}
	paramsAs2PartnerCreate.HttpAuthUsername = plan.HttpAuthUsername.ValueString()
	paramsAs2PartnerCreate.HttpAuthPassword = plan.HttpAuthPassword.ValueString()
	paramsAs2PartnerCreate.MdnValidationLevel = paramsAs2PartnerCreate.MdnValidationLevel.Enum()[plan.MdnValidationLevel.ValueString()]
	paramsAs2PartnerCreate.ServerCertificate = paramsAs2PartnerCreate.ServerCertificate.Enum()[plan.ServerCertificate.ValueString()]
	paramsAs2PartnerCreate.As2StationId = plan.As2StationId.ValueInt64()
	paramsAs2PartnerCreate.Name = plan.Name.ValueString()
	paramsAs2PartnerCreate.Uri = plan.Uri.ValueString()
	paramsAs2PartnerCreate.PublicCertificate = plan.PublicCertificate.ValueString()

	if resp.Diagnostics.HasError() {
		return
	}

	as2Partner, err := r.client.Create(paramsAs2PartnerCreate, files_sdk.WithContext(ctx))
	if err != nil {
		resp.Diagnostics.AddError(
			"Error Creating Files As2Partner",
			"Could not create as2_partner, unexpected error: "+err.Error(),
		)
		return
	}

	diags = r.populateResourceModel(ctx, as2Partner, &plan)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	diags = resp.State.Set(ctx, plan)
	resp.Diagnostics.Append(diags...)
}

func (r *as2PartnerResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	var state as2PartnerResourceModel
	diags := req.State.Get(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	paramsAs2PartnerFind := files_sdk.As2PartnerFindParams{}
	paramsAs2PartnerFind.Id = state.Id.ValueInt64()

	as2Partner, err := r.client.Find(paramsAs2PartnerFind, files_sdk.WithContext(ctx))
	if err != nil {
		if files_sdk.IsNotExist(err) {
			resp.State.RemoveResource(ctx)
			return
		}

		resp.Diagnostics.AddError(
			"Error Reading Files As2Partner",
			"Could not read as2_partner id "+fmt.Sprint(state.Id.ValueInt64())+": "+err.Error(),
		)
		return
	}

	diags = r.populateResourceModel(ctx, as2Partner, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	diags = resp.State.Set(ctx, state)
	resp.Diagnostics.Append(diags...)
}

func (r *as2PartnerResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	var plan as2PartnerResourceModel
	diags := req.Plan.Get(ctx, &plan)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	paramsAs2PartnerUpdate := files_sdk.As2PartnerUpdateParams{}
	paramsAs2PartnerUpdate.Id = plan.Id.ValueInt64()
	if !plan.EnableDedicatedIps.IsNull() && !plan.EnableDedicatedIps.IsUnknown() {
		paramsAs2PartnerUpdate.EnableDedicatedIps = plan.EnableDedicatedIps.ValueBoolPointer()
	}
	paramsAs2PartnerUpdate.HttpAuthUsername = plan.HttpAuthUsername.ValueString()
	paramsAs2PartnerUpdate.HttpAuthPassword = plan.HttpAuthPassword.ValueString()
	paramsAs2PartnerUpdate.MdnValidationLevel = paramsAs2PartnerUpdate.MdnValidationLevel.Enum()[plan.MdnValidationLevel.ValueString()]
	paramsAs2PartnerUpdate.ServerCertificate = paramsAs2PartnerUpdate.ServerCertificate.Enum()[plan.ServerCertificate.ValueString()]
	paramsAs2PartnerUpdate.Name = plan.Name.ValueString()
	paramsAs2PartnerUpdate.Uri = plan.Uri.ValueString()
	paramsAs2PartnerUpdate.PublicCertificate = plan.PublicCertificate.ValueString()

	if resp.Diagnostics.HasError() {
		return
	}

	as2Partner, err := r.client.Update(paramsAs2PartnerUpdate, files_sdk.WithContext(ctx))
	if err != nil {
		resp.Diagnostics.AddError(
			"Error Updating Files As2Partner",
			"Could not update as2_partner, unexpected error: "+err.Error(),
		)
		return
	}

	diags = r.populateResourceModel(ctx, as2Partner, &plan)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	diags = resp.State.Set(ctx, plan)
	resp.Diagnostics.Append(diags...)
}

func (r *as2PartnerResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	var state as2PartnerResourceModel
	diags := req.State.Get(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	paramsAs2PartnerDelete := files_sdk.As2PartnerDeleteParams{}
	paramsAs2PartnerDelete.Id = state.Id.ValueInt64()

	err := r.client.Delete(paramsAs2PartnerDelete, files_sdk.WithContext(ctx))
	if err != nil && !files_sdk.IsNotExist(err) {
		resp.Diagnostics.AddError(
			"Error Deleting Files As2Partner",
			"Could not delete as2_partner id "+fmt.Sprint(state.Id.ValueInt64())+": "+err.Error(),
		)
	}
}

func (r *as2PartnerResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
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

func (r *as2PartnerResource) populateResourceModel(ctx context.Context, as2Partner files_sdk.As2Partner, state *as2PartnerResourceModel) (diags diag.Diagnostics) {
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
