package provider

import (
	"context"
	"fmt"

	files_sdk "github.com/Files-com/files-sdk-go/v3"
	partner_site_request "github.com/Files-com/files-sdk-go/v3/partnersiterequest"
	"github.com/Files-com/terraform-provider-files/lib"
	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/datasource/schema"
	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

var (
	_ datasource.DataSource              = &partnerSiteRequestDataSource{}
	_ datasource.DataSourceWithConfigure = &partnerSiteRequestDataSource{}
)

func NewPartnerSiteRequestDataSource() datasource.DataSource {
	return &partnerSiteRequestDataSource{}
}

type partnerSiteRequestDataSource struct {
	client *partner_site_request.Client
}

type partnerSiteRequestDataSourceModel struct {
	Id           types.Int64  `tfsdk:"id"`
	PartnerId    types.Int64  `tfsdk:"partner_id"`
	LinkedSiteId types.Int64  `tfsdk:"linked_site_id"`
	Status       types.String `tfsdk:"status"`
	PairingKey   types.String `tfsdk:"pairing_key"`
	CreatedAt    types.String `tfsdk:"created_at"`
	UpdatedAt    types.String `tfsdk:"updated_at"`
}

func (r *partnerSiteRequestDataSource) Configure(_ context.Context, req datasource.ConfigureRequest, resp *datasource.ConfigureResponse) {
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

func (r *partnerSiteRequestDataSource) Metadata(_ context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_partner_site_request"
}

func (r *partnerSiteRequestDataSource) Schema(_ context.Context, _ datasource.SchemaRequest, resp *datasource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Description: "A PartnerSiteRequest represents a request to link a partner's Files.com site with another Files.com site.\n\n\n\nThe Site with the Partner can initiate a request, which generates a pairing key. The target site admin must then approve the request using the pairing key.",
		Attributes: map[string]schema.Attribute{
			"id": schema.Int64Attribute{
				Description: "Partner Site Request ID",
				Required:    true,
			},
			"partner_id": schema.Int64Attribute{
				Description: "Partner ID",
				Computed:    true,
			},
			"linked_site_id": schema.Int64Attribute{
				Description: "Linked Site ID",
				Computed:    true,
			},
			"status": schema.StringAttribute{
				Description: "Request status (pending, approved, rejected)",
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

func (r *partnerSiteRequestDataSource) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
	var data partnerSiteRequestDataSourceModel
	diags := req.Config.Get(ctx, &data)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	paramsPartnerSiteRequestList := files_sdk.PartnerSiteRequestListParams{}

	partnerSiteRequestIt, err := r.client.List(paramsPartnerSiteRequestList, files_sdk.WithContext(ctx))
	if err != nil {
		resp.Diagnostics.AddError(
			"Error Reading Files PartnerSiteRequest",
			"Could not read partner_site_request id "+fmt.Sprint(data.Id.ValueInt64())+": "+err.Error(),
		)
		return
	}

	var partnerSiteRequest *files_sdk.PartnerSiteRequest
	for partnerSiteRequestIt.Next() {
		entry := partnerSiteRequestIt.PartnerSiteRequest()
		if entry.Id == data.Id.ValueInt64() {
			partnerSiteRequest = &entry
			break
		}
	}

	if err = partnerSiteRequestIt.Err(); err != nil {
		resp.Diagnostics.AddError(
			"Error Reading Files PartnerSiteRequest",
			"Could not read partner_site_request id "+fmt.Sprint(data.Id.ValueInt64())+": "+err.Error(),
		)
	}

	if partnerSiteRequest == nil {
		resp.Diagnostics.AddError(
			"Error Reading Files PartnerSiteRequest",
			"Could not find partner_site_request id "+fmt.Sprint(data.Id.ValueInt64())+"",
		)
		return
	}

	diags = r.populateDataSourceModel(ctx, *partnerSiteRequest, &data)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	diags = resp.State.Set(ctx, data)
	resp.Diagnostics.Append(diags...)
}

func (r *partnerSiteRequestDataSource) populateDataSourceModel(ctx context.Context, partnerSiteRequest files_sdk.PartnerSiteRequest, state *partnerSiteRequestDataSourceModel) (diags diag.Diagnostics) {
	state.Id = types.Int64Value(partnerSiteRequest.Id)
	state.PartnerId = types.Int64Value(partnerSiteRequest.PartnerId)
	state.LinkedSiteId = types.Int64Value(partnerSiteRequest.LinkedSiteId)
	state.Status = types.StringValue(partnerSiteRequest.Status)
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
