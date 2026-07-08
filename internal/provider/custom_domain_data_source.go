package provider

import (
	"context"
	"fmt"

	files_sdk "github.com/Files-com/files-sdk-go/v3"
	custom_domain "github.com/Files-com/files-sdk-go/v3/customdomain"
	"github.com/Files-com/terraform-provider-files/lib"
	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/datasource/schema"
	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

var (
	_ datasource.DataSource              = &customDomainDataSource{}
	_ datasource.DataSourceWithConfigure = &customDomainDataSource{}
)

func NewCustomDomainDataSource() datasource.DataSource {
	return &customDomainDataSource{}
}

type customDomainDataSource struct {
	client *custom_domain.Client
}

type customDomainDataSourceModel struct {
	Id               types.Int64  `tfsdk:"id"`
	Domain           types.String `tfsdk:"domain"`
	Destination      types.String `tfsdk:"destination"`
	DnsStatus        types.String `tfsdk:"dns_status"`
	SslCertificateId types.Int64  `tfsdk:"ssl_certificate_id"`
	BrickManaged     types.Bool   `tfsdk:"brick_managed"`
	FolderBehaviorId types.Int64  `tfsdk:"folder_behavior_id"`
	CreatedAt        types.String `tfsdk:"created_at"`
	UpdatedAt        types.String `tfsdk:"updated_at"`
}

func (r *customDomainDataSource) Configure(_ context.Context, req datasource.ConfigureRequest, resp *datasource.ConfigureResponse) {
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

func (r *customDomainDataSource) Metadata(_ context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_custom_domain"
}

func (r *customDomainDataSource) Schema(_ context.Context, _ datasource.SchemaRequest, resp *datasource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Description: "A CustomDomain object represents an additional customer-owned domain that routes to a Files.com site.",
		Attributes: map[string]schema.Attribute{
			"id": schema.Int64Attribute{
				Description: "Custom Domain ID.",
				Required:    true,
			},
			"domain": schema.StringAttribute{
				Description: "Customer-owned domain name.",
				Computed:    true,
			},
			"destination": schema.StringAttribute{
				Description: "Where this custom domain routes. Can be `site_alias`, `public_hosting`, `s3_endpoint`, or `unassigned` (not routing traffic). Set to `unassigned` automatically when a bound `public_hosting` folder behavior is deleted, and can be set manually via the API for any reason.",
				Computed:    true,
			},
			"dns_status": schema.StringAttribute{
				Description: "Current DNS verification status.",
				Computed:    true,
			},
			"ssl_certificate_id": schema.Int64Attribute{
				Description: "Current SSL certificate ID.",
				Computed:    true,
			},
			"brick_managed": schema.BoolAttribute{
				Description: "Is this domain's SSL certificate automatically managed and renewed by Files.com?",
				Computed:    true,
			},
			"folder_behavior_id": schema.Int64Attribute{
				Description: "Public Hosting behavior ID when this domain routes to a specific Public Hosting behavior.  Preserved as historical context when `destination` becomes `unassigned`.",
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

func (r *customDomainDataSource) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
	var data customDomainDataSourceModel
	diags := req.Config.Get(ctx, &data)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	paramsCustomDomainFind := files_sdk.CustomDomainFindParams{}
	paramsCustomDomainFind.Id = data.Id.ValueInt64()

	customDomain, err := r.client.Find(paramsCustomDomainFind, files_sdk.WithContext(ctx))
	if err != nil {
		resp.Diagnostics.AddError(
			"Error Reading Files CustomDomain",
			"Could not read custom_domain id "+fmt.Sprint(data.Id.ValueInt64())+": "+err.Error(),
		)
		return
	}

	diags = r.populateDataSourceModel(ctx, customDomain, &data)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	diags = resp.State.Set(ctx, data)
	resp.Diagnostics.Append(diags...)
}

func (r *customDomainDataSource) populateDataSourceModel(ctx context.Context, customDomain files_sdk.CustomDomain, state *customDomainDataSourceModel) (diags diag.Diagnostics) {
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
