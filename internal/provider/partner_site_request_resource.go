package provider

import (
	"context"
	"fmt"
	"strconv"
	"strings"

	files_sdk "github.com/Files-com/files-sdk-go/v3"
	partner_site_request "github.com/Files-com/files-sdk-go/v3/partnersiterequest"
	"github.com/Files-com/terraform-provider-files/lib"
	"github.com/hashicorp/terraform-plugin-framework-validators/stringvalidator"
	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/int64planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/schema/validator"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

var (
	_ resource.Resource                = &partnerSiteRequestResource{}
	_ resource.ResourceWithConfigure   = &partnerSiteRequestResource{}
	_ resource.ResourceWithImportState = &partnerSiteRequestResource{}
)

func NewPartnerSiteRequestResource() resource.Resource {
	return &partnerSiteRequestResource{}
}

type partnerSiteRequestResource struct {
	client *partner_site_request.Client
}

type partnerSiteRequestResourceModel struct {
	PartnerId    types.Int64  `tfsdk:"partner_id"`
	SiteUrl      types.String `tfsdk:"site_url"`
	Id           types.Int64  `tfsdk:"id"`
	LinkedSiteId types.Int64  `tfsdk:"linked_site_id"`
	Status       types.String `tfsdk:"status"`
	MainSiteName types.String `tfsdk:"main_site_name"`
	PairingKey   types.String `tfsdk:"pairing_key"`
	CreatedAt    types.String `tfsdk:"created_at"`
	UpdatedAt    types.String `tfsdk:"updated_at"`
}

func (r *partnerSiteRequestResource) Configure(_ context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {
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

	r.client = &partner_site_request.Client{Config: sdk_config}
}

func (r *partnerSiteRequestResource) Metadata(_ context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_partner_site_request"
}

func (r *partnerSiteRequestResource) Schema(_ context.Context, _ resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Description: "A PartnerSiteRequest represents a request to link a partner's Files.com site with another Files.com site.\n\n\n\nThe Site with the Partner can initiate a request, which generates a pairing key. The target site admin must then approve the request using the pairing key.",
		Attributes: map[string]schema.Attribute{
			"partner_id": schema.Int64Attribute{
				Description: "Partner ID",
				Required:    true,
				PlanModifiers: []planmodifier.Int64{
					int64planmodifier.RequiresReplace(),
				},
			},
			"site_url": schema.StringAttribute{
				Description: "Site URL to link to",
				Required:    true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.RequiresReplace(),
				},
			},
			"id": schema.Int64Attribute{
				Description: "Partner Site Request ID",
				Computed:    true,
				PlanModifiers: []planmodifier.Int64{
					int64planmodifier.UseStateForUnknown(),
				},
			},
			"linked_site_id": schema.Int64Attribute{
				Description: "Linked Site ID",
				Computed:    true,
			},
			"status": schema.StringAttribute{
				Description: "Request status (pending, approved, rejected)",
				Computed:    true,
				Validators: []validator.String{
					stringvalidator.OneOf("pending", "approved", "rejected"),
				},
			},
			"main_site_name": schema.StringAttribute{
				Description: "Main Site Name",
				Computed:    true,
			},
			"pairing_key": schema.StringAttribute{
				Description: "Pairing key used to approve this request on the target site",
				Computed:    true,
			},
			"created_at": schema.StringAttribute{
				Description: "Request creation date/time",
				Computed:    true,
			},
			"updated_at": schema.StringAttribute{
				Description: "Request last updated date/time",
				Computed:    true,
			},
		},
	}
}

func (r *partnerSiteRequestResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	var plan partnerSiteRequestResourceModel
	diags := req.Plan.Get(ctx, &plan)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}
	var config partnerSiteRequestResourceModel
	diags = req.Config.Get(ctx, &config)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	paramsPartnerSiteRequestCreate := files_sdk.PartnerSiteRequestCreateParams{}
	paramsPartnerSiteRequestCreate.PartnerId = plan.PartnerId.ValueInt64()
	paramsPartnerSiteRequestCreate.SiteUrl = plan.SiteUrl.ValueString()

	if resp.Diagnostics.HasError() {
		return
	}

	partnerSiteRequest, err := r.client.Create(paramsPartnerSiteRequestCreate, files_sdk.WithContext(ctx))
	if err != nil {
		resp.Diagnostics.AddError(
			"Error Creating Files PartnerSiteRequest",
			"Could not create partner_site_request, unexpected error: "+err.Error(),
		)
		return
	}

	diags = r.populateResourceModel(ctx, partnerSiteRequest, &plan)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	diags = resp.State.Set(ctx, plan)
	resp.Diagnostics.Append(diags...)
}

func (r *partnerSiteRequestResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	var state partnerSiteRequestResourceModel
	diags := req.State.Get(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	paramsPartnerSiteRequestList := files_sdk.PartnerSiteRequestListParams{}

	partnerSiteRequestIt, err := r.client.List(paramsPartnerSiteRequestList, files_sdk.WithContext(ctx))
	if err != nil {
		if files_sdk.IsNotExist(err) {
			resp.State.RemoveResource(ctx)
			return
		}

		resp.Diagnostics.AddError(
			"Error Reading Files PartnerSiteRequest",
			"Could not read partner_site_request id "+fmt.Sprint(state.Id.ValueInt64())+": "+err.Error(),
		)
		return
	}

	var partnerSiteRequest *files_sdk.PartnerSiteRequest
	for partnerSiteRequestIt.Next() {
		entry := partnerSiteRequestIt.PartnerSiteRequest()
		if entry.Id == state.Id.ValueInt64() {
			partnerSiteRequest = &entry
			break
		}
	}

	if err = partnerSiteRequestIt.Err(); err != nil {
		if files_sdk.IsNotExist(err) {
			resp.State.RemoveResource(ctx)
			return
		}

		resp.Diagnostics.AddError(
			"Error Reading Files PartnerSiteRequest",
			"Could not read partner_site_request id "+fmt.Sprint(state.Id.ValueInt64())+": "+err.Error(),
		)
	}

	if partnerSiteRequest == nil {
		resp.State.RemoveResource(ctx)
		return
	}

	diags = r.populateResourceModel(ctx, *partnerSiteRequest, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	diags = resp.State.Set(ctx, state)
	resp.Diagnostics.Append(diags...)
}

func (r *partnerSiteRequestResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	resp.Diagnostics.AddError(
		"Resource Update Not Implemented",
		"This resource does not support updates.",
	)
}

func (r *partnerSiteRequestResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	var state partnerSiteRequestResourceModel
	diags := req.State.Get(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	paramsPartnerSiteRequestDelete := files_sdk.PartnerSiteRequestDeleteParams{}
	paramsPartnerSiteRequestDelete.Id = state.Id.ValueInt64()

	err := r.client.Delete(paramsPartnerSiteRequestDelete, files_sdk.WithContext(ctx))
	if err != nil && !files_sdk.IsNotExist(err) {
		resp.Diagnostics.AddError(
			"Error Deleting Files PartnerSiteRequest",
			"Could not delete partner_site_request id "+fmt.Sprint(state.Id.ValueInt64())+": "+err.Error(),
		)
	}
}

func (r *partnerSiteRequestResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
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

func (r *partnerSiteRequestResource) populateResourceModel(ctx context.Context, partnerSiteRequest files_sdk.PartnerSiteRequest, state *partnerSiteRequestResourceModel) (diags diag.Diagnostics) {
	state.Id = types.Int64Value(partnerSiteRequest.Id)
	state.PartnerId = types.Int64Value(partnerSiteRequest.PartnerId)
	state.LinkedSiteId = types.Int64Value(partnerSiteRequest.LinkedSiteId)
	state.Status = types.StringValue(partnerSiteRequest.Status)
	state.MainSiteName = types.StringValue(partnerSiteRequest.MainSiteName)
	state.PairingKey = types.StringValue(partnerSiteRequest.PairingKey)
	if err := lib.TimeToStringType(ctx, path.Root("created_at"), partnerSiteRequest.CreatedAt, &state.CreatedAt); err != nil {
		diags.AddError(
			"Error Creating Files PartnerSiteRequest",
			"Could not convert state created_at to string: "+err.Error(),
		)
	}
	if err := lib.TimeToStringType(ctx, path.Root("updated_at"), partnerSiteRequest.UpdatedAt, &state.UpdatedAt); err != nil {
		diags.AddError(
			"Error Creating Files PartnerSiteRequest",
			"Could not convert state updated_at to string: "+err.Error(),
		)
	}

	return
}
