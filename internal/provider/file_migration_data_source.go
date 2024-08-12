package provider

import (
	"context"
	"fmt"

	files_sdk "github.com/Files-com/files-sdk-go/v3"
	file_migration "github.com/Files-com/files-sdk-go/v3/filemigration"
	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/datasource/schema"
	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

var (
	_ datasource.DataSource              = &fileMigrationDataSource{}
	_ datasource.DataSourceWithConfigure = &fileMigrationDataSource{}
)

func NewFileMigrationDataSource() datasource.DataSource {
	return &fileMigrationDataSource{}
}

type fileMigrationDataSource struct {
	client *file_migration.Client
}

type fileMigrationDataSourceModel struct {
	Id         types.Int64  `tfsdk:"id"`
	Path       types.String `tfsdk:"path"`
	DestPath   types.String `tfsdk:"dest_path"`
	FilesMoved types.Int64  `tfsdk:"files_moved"`
	FilesTotal types.Int64  `tfsdk:"files_total"`
	Operation  types.String `tfsdk:"operation"`
	Region     types.String `tfsdk:"region"`
	Status     types.String `tfsdk:"status"`
	LogUrl     types.String `tfsdk:"log_url"`
}

func (r *fileMigrationDataSource) Configure(_ context.Context, req datasource.ConfigureRequest, resp *datasource.ConfigureResponse) {
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

	r.client = &file_migration.Client{Config: sdk_config}
}

func (r *fileMigrationDataSource) Metadata(_ context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_file_migration"
}

func (r *fileMigrationDataSource) Schema(_ context.Context, _ datasource.SchemaRequest, resp *datasource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Description: "A FileMigration is a background operation on one or more files, such as a copy or a region migration.\n\n\n\nIf no `operation` or `dest_path` is present, then the record represents a region migration.",
		Attributes: map[string]schema.Attribute{
			"id": schema.Int64Attribute{
				Description: "File migration ID",
				Required:    true,
			},
			"path": schema.StringAttribute{
				Description: "Source path. This must be slash-delimited, but it must neither start nor end with a slash. Maximum of 5000 characters.",
				Computed:    true,
			},
			"dest_path": schema.StringAttribute{
				Description: "Destination path",
				Computed:    true,
			},
			"files_moved": schema.Int64Attribute{
				Description: "Number of files processed",
				Computed:    true,
			},
			"files_total": schema.Int64Attribute{
				Computed: true,
			},
			"operation": schema.StringAttribute{
				Description: "The type of operation",
				Computed:    true,
			},
			"region": schema.StringAttribute{
				Description: "Region",
				Computed:    true,
			},
			"status": schema.StringAttribute{
				Description: "Status",
				Computed:    true,
			},
			"log_url": schema.StringAttribute{
				Description: "Link to download the log file for this migration.",
				Computed:    true,
			},
		},
	}
}

func (r *fileMigrationDataSource) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
	var data fileMigrationDataSourceModel
	diags := req.Config.Get(ctx, &data)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	paramsFileMigrationFind := files_sdk.FileMigrationFindParams{}
	paramsFileMigrationFind.Id = data.Id.ValueInt64()

	fileMigration, err := r.client.Find(paramsFileMigrationFind, files_sdk.WithContext(ctx))
	if err != nil {
		resp.Diagnostics.AddError(
			"Error Reading Files FileMigration",
			"Could not read file_migration id "+fmt.Sprint(data.Id.ValueInt64())+": "+err.Error(),
		)
		return
	}

	diags = r.populateDataSourceModel(ctx, fileMigration, &data)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	diags = resp.State.Set(ctx, data)
	resp.Diagnostics.Append(diags...)
}

func (r *fileMigrationDataSource) populateDataSourceModel(ctx context.Context, fileMigration files_sdk.FileMigration, state *fileMigrationDataSourceModel) (diags diag.Diagnostics) {
	state.Id = types.Int64Value(fileMigration.Id)
	state.Path = types.StringValue(fileMigration.Path)
	state.DestPath = types.StringValue(fileMigration.DestPath)
	state.FilesMoved = types.Int64Value(fileMigration.FilesMoved)
	state.FilesTotal = types.Int64Value(fileMigration.FilesTotal)
	state.Operation = types.StringValue(fileMigration.Operation)
	state.Region = types.StringValue(fileMigration.Region)
	state.Status = types.StringValue(fileMigration.Status)
	state.LogUrl = types.StringValue(fileMigration.LogUrl)

	return
}
