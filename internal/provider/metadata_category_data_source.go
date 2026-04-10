package provider

import (
	"context"
	"fmt"

	files_sdk "github.com/Files-com/files-sdk-go/v3"
	metadata_category "github.com/Files-com/files-sdk-go/v3/metadatacategory"
	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/datasource/schema"
	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

var (
	_ datasource.DataSource              = &metadataCategoryDataSource{}
	_ datasource.DataSourceWithConfigure = &metadataCategoryDataSource{}
)

func NewMetadataCategoryDataSource() datasource.DataSource {
	return &metadataCategoryDataSource{}
}

type metadataCategoryDataSource struct {
	client *metadata_category.Client
}

type metadataCategoryDataSourceModel struct {
	Id             types.Int64  `tfsdk:"id"`
	Name           types.String `tfsdk:"name"`
	Definitions    types.Map    `tfsdk:"definitions"`
	DefaultColumns types.List   `tfsdk:"default_columns"`
}

func (r *metadataCategoryDataSource) Configure(_ context.Context, req datasource.ConfigureRequest, resp *datasource.ConfigureResponse) {
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

	r.client = &metadata_category.Client{Config: sdk_config}
}

func (r *metadataCategoryDataSource) Metadata(_ context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_metadata_category"
}

func (r *metadataCategoryDataSource) Schema(_ context.Context, _ datasource.SchemaRequest, resp *datasource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Description: "A MetadataCategory defines a reusable set of Custom Metadata rules that can be assigned to folders\n\nvia a folder behavior. Each category specifies named metadata keys with optional allowed-value\n\nconstraints, and a set of default columns to display in the UI.\n\n\n\nIf a key's `allowed_values` array is empty, it is treated as a free-form text field.\n\nIf the array is non-empty, the key is constrained to those values in the Web UI.",
		Attributes: map[string]schema.Attribute{
			"id": schema.Int64Attribute{
				Description: "Metadata Category ID",
				Required:    true,
			},
			"name": schema.StringAttribute{
				Description: "Name of the metadata category.",
				Computed:    true,
			},
			"definitions": schema.MapAttribute{
				Description: "Map of key names to arrays of allowed values. An empty array means free-form text.",
				Computed:    true,
				ElementType: types.ListType{ElemType: types.StringType},
			},
			"default_columns": schema.ListAttribute{
				Description: "Metadata keys that should appear as columns in the UI by default.",
				Computed:    true,
				ElementType: types.StringType,
			},
		},
	}
}

func (r *metadataCategoryDataSource) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
	var data metadataCategoryDataSourceModel
	diags := req.Config.Get(ctx, &data)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	paramsMetadataCategoryFind := files_sdk.MetadataCategoryFindParams{}
	paramsMetadataCategoryFind.Id = data.Id.ValueInt64()

	metadataCategory, err := r.client.Find(paramsMetadataCategoryFind, files_sdk.WithContext(ctx))
	if err != nil {
		resp.Diagnostics.AddError(
			"Error Reading Files MetadataCategory",
			"Could not read metadata_category id "+fmt.Sprint(data.Id.ValueInt64())+": "+err.Error(),
		)
		return
	}

	diags = r.populateDataSourceModel(ctx, metadataCategory, &data)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	diags = resp.State.Set(ctx, data)
	resp.Diagnostics.Append(diags...)
}

func (r *metadataCategoryDataSource) populateDataSourceModel(ctx context.Context, metadataCategory files_sdk.MetadataCategory, state *metadataCategoryDataSourceModel) (diags diag.Diagnostics) {
	var propDiags diag.Diagnostics

	state.Id = types.Int64Value(metadataCategory.Id)
	state.Name = types.StringValue(metadataCategory.Name)
	state.Definitions, propDiags = types.MapValueFrom(ctx, types.ListType{ElemType: types.StringType}, metadataCategory.Definitions)
	diags.Append(propDiags...)
	state.DefaultColumns, propDiags = types.ListValueFrom(ctx, types.StringType, metadataCategory.DefaultColumns)
	diags.Append(propDiags...)

	return
}
