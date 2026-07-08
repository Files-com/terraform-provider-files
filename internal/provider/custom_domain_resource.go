package provider

import (
	"context"
	"fmt"
	"strconv"
	"strings"

	files_sdk "github.com/Files-com/files-sdk-go/v3"
	custom_domain "github.com/Files-com/files-sdk-go/v3/customdomain"
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
	_ resource.Resource                = &customDomainResource{}
	_ resource.ResourceWithConfigure   = &customDomainResource{}
	_ resource.ResourceWithImportState = &customDomainResource{}
)

func NewCustomDomainResource() resource.Resource {
	return &customDomainResource{}
}

type customDomainResource struct {
	client *custom_domain.Client
}

type customDomainResourceModel struct {
	Domain           types.String `tfsdk:"domain"`
	Destination      types.String `tfsdk:"destination"`
	SslCertificateId types.Int64  `tfsdk:"ssl_certificate_id"`
	FolderBehaviorId types.Int64  `tfsdk:"folder_behavior_id"`
	Id               types.Int64  `tfsdk:"id"`
	DnsStatus        types.String `tfsdk:"dns_status"`
	BrickManaged     types.Bool   `tfsdk:"brick_managed"`
	CreatedAt        types.String `tfsdk:"created_at"`
	UpdatedAt        types.String `tfsdk:"updated_at"`
}

func (r *customDomainResource) Configure(_ context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {
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

	r.client = &custom_domain.Client{Config: sdk_config}
}

func (r *customDomainResource) Metadata(_ context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_custom_domain"
}

func (r *customDomainResource) Schema(_ context.Context, _ resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Description: "A CustomDomain object represents an additional customer-owned domain that routes to a Files.com site.",
		Attributes: map[string]schema.Attribute{
			"domain": schema.StringAttribute{
				Description: "Customer-owned domain name.",
				Required:    true,
			},
			"destination": schema.StringAttribute{
				Description: "Where this custom domain routes. Can be `site_alias`, `public_hosting`, `s3_endpoint`, or `unassigned` (not routing traffic). Set to `unassigned` automatically when a bound `public_hosting` folder behavior is deleted, and can be set manually via the API for any reason.",
				Computed:    true,
				Optional:    true,
				Validators: []validator.String{
					stringvalidator.OneOf("site_alias", "public_hosting", "s3_endpoint", "unassigned"),
				},
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
			},
			"ssl_certificate_id": schema.Int64Attribute{
				Description: "Current SSL certificate ID.",
				Computed:    true,
				Optional:    true,
				PlanModifiers: []planmodifier.Int64{
					int64planmodifier.UseStateForUnknown(),
				},
			},
			"folder_behavior_id": schema.Int64Attribute{
				Description: "Public Hosting behavior ID when this domain routes to a specific Public Hosting behavior.  Preserved as historical context when `destination` becomes `unassigned`.",
				Computed:    true,
				Optional:    true,
				PlanModifiers: []planmodifier.Int64{
					int64planmodifier.UseStateForUnknown(),
				},
			},
			"id": schema.Int64Attribute{
				Description: "Custom Domain ID.",
				Computed:    true,
				PlanModifiers: []planmodifier.Int64{
					int64planmodifier.UseStateForUnknown(),
				},
			},
			"dns_status": schema.StringAttribute{
				Description: "Current DNS verification status.",
				Computed:    true,
				Validators: []validator.String{
					stringvalidator.OneOf("disabled", "correct", "no_records", "cname_wrong", "a_record", "caa_conflict"),
				},
			},
			"brick_managed": schema.BoolAttribute{
				Description: "Is this domain's SSL certificate automatically managed and renewed by Files.com?",
				Computed:    true,
			},
			"created_at": schema.StringAttribute{
				Description: "When this Custom Domain was created.",
				Computed:    true,
			},
			"updated_at": schema.StringAttribute{
				Description: "When this Custom Domain was last updated.",
				Computed:    true,
			},
		},
	}
}

func (r *customDomainResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	var plan customDomainResourceModel
	diags := req.Plan.Get(ctx, &plan)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}
	var config customDomainResourceModel
	diags = req.Config.Get(ctx, &config)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	paramsCustomDomainCreate := files_sdk.CustomDomainCreateParams{}
	paramsCustomDomainCreate.Destination = paramsCustomDomainCreate.Destination.Enum()[plan.Destination.ValueString()]
	paramsCustomDomainCreate.FolderBehaviorId = plan.FolderBehaviorId.ValueInt64()
	paramsCustomDomainCreate.SslCertificateId = plan.SslCertificateId.ValueInt64()
	paramsCustomDomainCreate.Domain = plan.Domain.ValueString()

	if resp.Diagnostics.HasError() {
		return
	}

	customDomain, err := r.client.Create(paramsCustomDomainCreate, files_sdk.WithContext(ctx))
	if err != nil {
		resp.Diagnostics.AddError(
			"Error Creating Files CustomDomain",
			"Could not create custom_domain, unexpected error: "+err.Error(),
		)
		return
	}

	diags = r.populateResourceModel(ctx, customDomain, &plan)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	diags = resp.State.Set(ctx, plan)
	resp.Diagnostics.Append(diags...)
}

func (r *customDomainResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	var state customDomainResourceModel
	diags := req.State.Get(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	paramsCustomDomainFind := files_sdk.CustomDomainFindParams{}
	paramsCustomDomainFind.Id = state.Id.ValueInt64()

	customDomain, err := r.client.Find(paramsCustomDomainFind, files_sdk.WithContext(ctx))
	if err != nil {
		if files_sdk.IsNotExist(err) {
			resp.State.RemoveResource(ctx)
			return
		}

		resp.Diagnostics.AddError(
			"Error Reading Files CustomDomain",
			"Could not read custom_domain id "+fmt.Sprint(state.Id.ValueInt64())+": "+err.Error(),
		)
		return
	}

	diags = r.populateResourceModel(ctx, customDomain, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	diags = resp.State.Set(ctx, state)
	resp.Diagnostics.Append(diags...)
}

func (r *customDomainResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	var plan customDomainResourceModel
	diags := req.Plan.Get(ctx, &plan)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}
	var config customDomainResourceModel
	diags = req.Config.Get(ctx, &config)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	paramsCustomDomainUpdate := map[string]interface{}{}
	if !plan.Id.IsNull() && !plan.Id.IsUnknown() {
		paramsCustomDomainUpdate["id"] = plan.Id.ValueInt64()
	}
	if !config.Destination.IsNull() && !config.Destination.IsUnknown() {
		paramsCustomDomainUpdate["destination"] = config.Destination.ValueString()
	}
	if !config.FolderBehaviorId.IsNull() && !config.FolderBehaviorId.IsUnknown() {
		paramsCustomDomainUpdate["folder_behavior_id"] = config.FolderBehaviorId.ValueInt64()
	}
	if !config.SslCertificateId.IsNull() && !config.SslCertificateId.IsUnknown() {
		paramsCustomDomainUpdate["ssl_certificate_id"] = config.SslCertificateId.ValueInt64()
	}
	if !config.Domain.IsNull() && !config.Domain.IsUnknown() {
		paramsCustomDomainUpdate["domain"] = config.Domain.ValueString()
	}

	if resp.Diagnostics.HasError() {
		return
	}

	customDomain, err := r.client.UpdateWithMap(paramsCustomDomainUpdate, files_sdk.WithContext(ctx))
	if err != nil {
		resp.Diagnostics.AddError(
			"Error Updating Files CustomDomain",
			"Could not update custom_domain, unexpected error: "+err.Error(),
		)
		return
	}

	diags = r.populateResourceModel(ctx, customDomain, &plan)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	diags = resp.State.Set(ctx, plan)
	resp.Diagnostics.Append(diags...)
}

func (r *customDomainResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	var state customDomainResourceModel
	diags := req.State.Get(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	paramsCustomDomainDelete := files_sdk.CustomDomainDeleteParams{}
	paramsCustomDomainDelete.Id = state.Id.ValueInt64()

	err := r.client.Delete(paramsCustomDomainDelete, files_sdk.WithContext(ctx))
	if err != nil && !files_sdk.IsNotExist(err) {
		resp.Diagnostics.AddError(
			"Error Deleting Files CustomDomain",
			"Could not delete custom_domain id "+fmt.Sprint(state.Id.ValueInt64())+": "+err.Error(),
		)
	}
}

func (r *customDomainResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
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

func (r *customDomainResource) populateResourceModel(ctx context.Context, customDomain files_sdk.CustomDomain, state *customDomainResourceModel) (diags diag.Diagnostics) {
	state.Id = types.Int64Value(customDomain.Id)
	state.Domain = types.StringValue(customDomain.Domain)
	state.Destination = types.StringValue(customDomain.Destination)
	state.DnsStatus = types.StringValue(customDomain.DnsStatus)
	state.SslCertificateId = types.Int64Value(customDomain.SslCertificateId)
	state.BrickManaged = types.BoolPointerValue(customDomain.BrickManaged)
	state.FolderBehaviorId = types.Int64Value(customDomain.FolderBehaviorId)
	if err := lib.TimeToStringType(ctx, path.Root("created_at"), customDomain.CreatedAt, &state.CreatedAt); err != nil {
		diags.AddError(
			"Error Creating Files CustomDomain",
			"Could not convert state created_at to string: "+err.Error(),
		)
	}
	if err := lib.TimeToStringType(ctx, path.Root("updated_at"), customDomain.UpdatedAt, &state.UpdatedAt); err != nil {
		diags.AddError(
			"Error Creating Files CustomDomain",
			"Could not convert state updated_at to string: "+err.Error(),
		)
	}

	return
}
