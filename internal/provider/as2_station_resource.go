package provider

import (
	"context"
	"fmt"
	"strconv"
	"strings"

	files_sdk "github.com/Files-com/files-sdk-go/v3"
	as2_station "github.com/Files-com/files-sdk-go/v3/as2station"
	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/int64planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

var (
	_ resource.Resource                = &as2StationResource{}
	_ resource.ResourceWithConfigure   = &as2StationResource{}
	_ resource.ResourceWithImportState = &as2StationResource{}
)

func NewAs2StationResource() resource.Resource {
	return &as2StationResource{}
}

type as2StationResource struct {
	client *as2_station.Client
}

type as2StationResourceModel struct {
	Name                       types.String `tfsdk:"name"`
	PublicCertificate          types.String `tfsdk:"public_certificate"`
	PrivateKey                 types.String `tfsdk:"private_key"`
	PrivateKeyPassword         types.String `tfsdk:"private_key_password"`
	Id                         types.Int64  `tfsdk:"id"`
	Uri                        types.String `tfsdk:"uri"`
	Domain                     types.String `tfsdk:"domain"`
	HexPublicCertificateSerial types.String `tfsdk:"hex_public_certificate_serial"`
	PublicCertificateMd5       types.String `tfsdk:"public_certificate_md5"`
	PrivateKeyMd5              types.String `tfsdk:"private_key_md5"`
	PublicCertificateSubject   types.String `tfsdk:"public_certificate_subject"`
	PublicCertificateIssuer    types.String `tfsdk:"public_certificate_issuer"`
	PublicCertificateSerial    types.String `tfsdk:"public_certificate_serial"`
	PublicCertificateNotBefore types.String `tfsdk:"public_certificate_not_before"`
	PublicCertificateNotAfter  types.String `tfsdk:"public_certificate_not_after"`
	PrivateKeyPasswordMd5      types.String `tfsdk:"private_key_password_md5"`
}

func (r *as2StationResource) Configure(_ context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {
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

func (r *as2StationResource) Metadata(_ context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_as2_station"
}

func (r *as2StationResource) Schema(_ context.Context, _ resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Description: "An AS2Station is a remote AS2 server that can send data into Files.com and receive data from Files.com.",
		Attributes: map[string]schema.Attribute{
			"name": schema.StringAttribute{
				Description: "The station's formal AS2 name.",
				Required:    true,
			},
			"public_certificate": schema.StringAttribute{
				Description: "Public certificate used for message security.",
				Required:    true,
			},
			"private_key": schema.StringAttribute{
				Required:  true,
				WriteOnly: true,
			},
			"private_key_password": schema.StringAttribute{
				Optional:  true,
				WriteOnly: true,
			},
			"id": schema.Int64Attribute{
				Description: "Id of the AS2 Station.",
				Computed:    true,
				PlanModifiers: []planmodifier.Int64{
					int64planmodifier.UseStateForUnknown(),
				},
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

func (r *as2StationResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	var plan as2StationResourceModel
	diags := req.Plan.Get(ctx, &plan)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}
	var config as2StationResourceModel
	diags = req.Config.Get(ctx, &config)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	paramsAs2StationCreate := files_sdk.As2StationCreateParams{}
	paramsAs2StationCreate.Name = plan.Name.ValueString()
	paramsAs2StationCreate.PublicCertificate = plan.PublicCertificate.ValueString()
	paramsAs2StationCreate.PrivateKey = config.PrivateKey.ValueString()
	paramsAs2StationCreate.PrivateKeyPassword = config.PrivateKeyPassword.ValueString()

	if resp.Diagnostics.HasError() {
		return
	}

	as2Station, err := r.client.Create(paramsAs2StationCreate, files_sdk.WithContext(ctx))
	if err != nil {
		resp.Diagnostics.AddError(
			"Error Creating Files As2Station",
			"Could not create as2_station, unexpected error: "+err.Error(),
		)
		return
	}

	diags = r.populateResourceModel(ctx, as2Station, &plan)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	diags = resp.State.Set(ctx, plan)
	resp.Diagnostics.Append(diags...)
}

func (r *as2StationResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	var state as2StationResourceModel
	diags := req.State.Get(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	paramsAs2StationFind := files_sdk.As2StationFindParams{}
	paramsAs2StationFind.Id = state.Id.ValueInt64()

	as2Station, err := r.client.Find(paramsAs2StationFind, files_sdk.WithContext(ctx))
	if err != nil {
		if files_sdk.IsNotExist(err) {
			resp.State.RemoveResource(ctx)
			return
		}

		resp.Diagnostics.AddError(
			"Error Reading Files As2Station",
			"Could not read as2_station id "+fmt.Sprint(state.Id.ValueInt64())+": "+err.Error(),
		)
		return
	}

	diags = r.populateResourceModel(ctx, as2Station, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	diags = resp.State.Set(ctx, state)
	resp.Diagnostics.Append(diags...)
}

func (r *as2StationResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	var plan as2StationResourceModel
	diags := req.Plan.Get(ctx, &plan)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}
	var config as2StationResourceModel
	diags = req.Config.Get(ctx, &config)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	paramsAs2StationUpdate := files_sdk.As2StationUpdateParams{}
	paramsAs2StationUpdate.Id = plan.Id.ValueInt64()
	paramsAs2StationUpdate.Name = plan.Name.ValueString()
	paramsAs2StationUpdate.PublicCertificate = plan.PublicCertificate.ValueString()
	paramsAs2StationUpdate.PrivateKey = config.PrivateKey.ValueString()
	paramsAs2StationUpdate.PrivateKeyPassword = config.PrivateKeyPassword.ValueString()

	if resp.Diagnostics.HasError() {
		return
	}

	as2Station, err := r.client.Update(paramsAs2StationUpdate, files_sdk.WithContext(ctx))
	if err != nil {
		resp.Diagnostics.AddError(
			"Error Updating Files As2Station",
			"Could not update as2_station, unexpected error: "+err.Error(),
		)
		return
	}

	diags = r.populateResourceModel(ctx, as2Station, &plan)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	diags = resp.State.Set(ctx, plan)
	resp.Diagnostics.Append(diags...)
}

func (r *as2StationResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	var state as2StationResourceModel
	diags := req.State.Get(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	paramsAs2StationDelete := files_sdk.As2StationDeleteParams{}
	paramsAs2StationDelete.Id = state.Id.ValueInt64()

	err := r.client.Delete(paramsAs2StationDelete, files_sdk.WithContext(ctx))
	if err != nil && !files_sdk.IsNotExist(err) {
		resp.Diagnostics.AddError(
			"Error Deleting Files As2Station",
			"Could not delete as2_station id "+fmt.Sprint(state.Id.ValueInt64())+": "+err.Error(),
		)
	}
}

func (r *as2StationResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
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

func (r *as2StationResource) populateResourceModel(ctx context.Context, as2Station files_sdk.As2Station, state *as2StationResourceModel) (diags diag.Diagnostics) {
	state.Id = types.Int64Value(as2Station.Id)
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
