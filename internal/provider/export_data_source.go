package provider

import (
	"context"
	"fmt"

	files_sdk "github.com/Files-com/files-sdk-go/v3"
	export "github.com/Files-com/files-sdk-go/v3/export"
	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/datasource/schema"
	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

var (
	_ datasource.DataSource              = &exportDataSource{}
	_ datasource.DataSourceWithConfigure = &exportDataSource{}
)

func NewExportDataSource() datasource.DataSource {
	return &exportDataSource{}
}

type exportDataSource struct {
	client *export.Client
}

type exportDataSourceModel struct {
	Id           types.Int64  `tfsdk:"id"`
	ExportStatus types.String `tfsdk:"export_status"`
	ExportType   types.String `tfsdk:"export_type"`
	DownloadUri  types.String `tfsdk:"download_uri"`
}

func (r *exportDataSource) Configure(_ context.Context, req datasource.ConfigureRequest, resp *datasource.ConfigureResponse) {
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

	r.client = &export.Client{Config: sdk_config}
}

func (r *exportDataSource) Metadata(_ context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_export"
}

func (r *exportDataSource) Schema(_ context.Context, _ datasource.SchemaRequest, resp *datasource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Description: "",
		Attributes: map[string]schema.Attribute{
			"id": schema.Int64Attribute{
				Description: "ID for this Export",
				Required:    true,
			},
			"export_status": schema.StringAttribute{
				Description: "Status of the Export",
				Computed:    true,
			},
			"export_type": schema.StringAttribute{
				Description: "Type of data being exported",
				Computed:    true,
			},
			"download_uri": schema.StringAttribute{
				Description: "Link to download Export file.",
				Computed:    true,
			},
		},
	}
}

func (r *exportDataSource) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
	var data exportDataSourceModel
	diags := req.Config.Get(ctx, &data)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	paramsExportFind := files_sdk.ExportFindParams{}
	paramsExportFind.Id = data.Id.ValueInt64()

	export, err := r.client.Find(paramsExportFind, files_sdk.WithContext(ctx))
	if err != nil {
		resp.Diagnostics.AddError(
			"Error Reading Files Export",
			"Could not read export id "+fmt.Sprint(data.Id.ValueInt64())+": "+err.Error(),
		)
		return
	}

	diags = r.populateDataSourceModel(ctx, export, &data)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	diags = resp.State.Set(ctx, data)
	resp.Diagnostics.Append(diags...)
}

func (r *exportDataSource) populateDataSourceModel(ctx context.Context, export files_sdk.Export, state *exportDataSourceModel) (diags diag.Diagnostics) {
	state.Id = types.Int64Value(export.Id)
	state.ExportStatus = types.StringValue(export.ExportStatus)
	state.ExportType = types.StringValue(export.ExportType)
	state.DownloadUri = types.StringValue(export.DownloadUri)

	return
}
