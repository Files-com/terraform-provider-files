package provider

import (
	"context"
	"fmt"

	files_sdk "github.com/Files-com/files-sdk-go/v3"
	site_subdomain_redirect "github.com/Files-com/files-sdk-go/v3/sitesubdomainredirect"
	"github.com/Files-com/terraform-provider-files/lib"
	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/datasource/schema"
	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

var (
	_ datasource.DataSource              = &siteSubdomainRedirectDataSource{}
	_ datasource.DataSourceWithConfigure = &siteSubdomainRedirectDataSource{}
)

func NewSiteSubdomainRedirectDataSource() datasource.DataSource {
	return &siteSubdomainRedirectDataSource{}
}

type siteSubdomainRedirectDataSource struct {
	client *site_subdomain_redirect.Client
}

type siteSubdomainRedirectDataSourceModel struct {
	Id        types.Int64  `tfsdk:"id"`
	Subdomain types.String `tfsdk:"subdomain"`
	CreatedAt types.String `tfsdk:"created_at"`
	UpdatedAt types.String `tfsdk:"updated_at"`
}

func (r *siteSubdomainRedirectDataSource) Configure(_ context.Context, req datasource.ConfigureRequest, resp *datasource.ConfigureResponse) {
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

	r.client = &site_subdomain_redirect.Client{Config: sdk_config}
}

func (r *siteSubdomainRedirectDataSource) Metadata(_ context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_site_subdomain_redirect"
}

func (r *siteSubdomainRedirectDataSource) Schema(_ context.Context, _ datasource.SchemaRequest, resp *datasource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Description: "A SiteSubdomainRedirect object represents an old Files.com subdomain that continues to work after the site's Files.com subdomain changes.\n\nHTTPS requests redirect to the current subdomain, and other protocols such as FTP and SFTP are routed through DNS.",
		Attributes: map[string]schema.Attribute{
			"id": schema.Int64Attribute{
				Description: "Site subdomain redirect ID.",
				Required:    true,
			},
			"subdomain": schema.StringAttribute{
				Description: "Files.com subdomain that continues to route to the current site subdomain.",
				Computed:    true,
			},
			"created_at": schema.StringAttribute{
				Description: "When this redirect was created.",
				Computed:    true,
			},
			"updated_at": schema.StringAttribute{
				Description: "When this redirect was last updated.",
				Computed:    true,
			},
		},
	}
}

func (r *siteSubdomainRedirectDataSource) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
	var data siteSubdomainRedirectDataSourceModel
	diags := req.Config.Get(ctx, &data)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	paramsSiteSubdomainRedirectFind := files_sdk.SiteSubdomainRedirectFindParams{}
	paramsSiteSubdomainRedirectFind.Id = data.Id.ValueInt64()

	siteSubdomainRedirect, err := r.client.Find(paramsSiteSubdomainRedirectFind, files_sdk.WithContext(ctx))
	if err != nil {
		resp.Diagnostics.AddError(
			"Error Reading Files SiteSubdomainRedirect",
			"Could not read site_subdomain_redirect id "+fmt.Sprint(data.Id.ValueInt64())+": "+err.Error(),
		)
		return
	}

	diags = r.populateDataSourceModel(ctx, siteSubdomainRedirect, &data)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	diags = resp.State.Set(ctx, data)
	resp.Diagnostics.Append(diags...)
}

func (r *siteSubdomainRedirectDataSource) populateDataSourceModel(ctx context.Context, siteSubdomainRedirect files_sdk.SiteSubdomainRedirect, state *siteSubdomainRedirectDataSourceModel) (diags diag.Diagnostics) {
	state.Id = types.Int64Value(siteSubdomainRedirect.Id)
	state.Subdomain = types.StringValue(siteSubdomainRedirect.Subdomain)
	if err := lib.TimeToStringType(ctx, path.Root("created_at"), siteSubdomainRedirect.CreatedAt, &state.CreatedAt); err != nil {
		diags.AddError(
			"Error Creating Files SiteSubdomainRedirect",
			"Could not convert state created_at to string: "+err.Error(),
		)
	}
	if err := lib.TimeToStringType(ctx, path.Root("updated_at"), siteSubdomainRedirect.UpdatedAt, &state.UpdatedAt); err != nil {
		diags.AddError(
			"Error Creating Files SiteSubdomainRedirect",
			"Could not convert state updated_at to string: "+err.Error(),
		)
	}

	return
}
