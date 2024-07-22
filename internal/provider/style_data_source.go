package provider

import (
	"context"
	"encoding/json"
	"fmt"

	files_sdk "github.com/Files-com/files-sdk-go/v3"
	style "github.com/Files-com/files-sdk-go/v3/style"
	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/datasource/schema"
	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

var (
	_ datasource.DataSource              = &styleDataSource{}
	_ datasource.DataSourceWithConfigure = &styleDataSource{}
)

func NewStyleDataSource() datasource.DataSource {
	return &styleDataSource{}
}

type styleDataSource struct {
	client *style.Client
}

type styleDataSourceModel struct {
	Path      types.String `tfsdk:"path"`
	Id        types.Int64  `tfsdk:"id"`
	Logo      types.String `tfsdk:"logo"`
	Thumbnail types.String `tfsdk:"thumbnail"`
}

func (r *styleDataSource) Configure(_ context.Context, req datasource.ConfigureRequest, resp *datasource.ConfigureResponse) {
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

	r.client = &style.Client{Config: sdk_config}
}

func (r *styleDataSource) Metadata(_ context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_style"
}

func (r *styleDataSource) Schema(_ context.Context, _ datasource.SchemaRequest, resp *datasource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Description: "Styles are custom sets of branding that can be applied on a per-folder basis.\n\nCurrently these only support Logos per folder, but in the future we may extend these to also support colors.\n\nIf you want to see that, please let us know so we can add your vote to the list.",
		Attributes: map[string]schema.Attribute{
			"path": schema.StringAttribute{
				Description: "Folder path This must be slash-delimited, but it must neither start nor end with a slash. Maximum of 5000 characters.",
				Required:    true,
			},
			"id": schema.Int64Attribute{
				Description: "Style ID",
				Computed:    true,
			},
			"logo": schema.StringAttribute{
				Description: "Logo",
				Computed:    true,
			},
			"thumbnail": schema.StringAttribute{
				Description: "Logo thumbnail",
				Computed:    true,
			},
		},
	}
}

func (r *styleDataSource) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
	var data styleDataSourceModel
	diags := req.Config.Get(ctx, &data)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	paramsStyleFind := files_sdk.StyleFindParams{}
	paramsStyleFind.Path = data.Path.ValueString()

	style, err := r.client.Find(paramsStyleFind, files_sdk.WithContext(ctx))
	if err != nil {
		resp.Diagnostics.AddError(
			"Error Reading Files Style",
			"Could not read style path "+fmt.Sprint(data.Path.ValueString())+": "+err.Error(),
		)
		return
	}

	diags = r.populateDataSourceModel(ctx, style, &data)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	diags = resp.State.Set(ctx, data)
	resp.Diagnostics.Append(diags...)
}

func (r *styleDataSource) populateDataSourceModel(ctx context.Context, style files_sdk.Style, state *styleDataSourceModel) (diags diag.Diagnostics) {
	state.Id = types.Int64Value(style.Id)
	state.Path = types.StringValue(style.Path)
	respLogo, err := json.Marshal(style.Logo)
	if err != nil {
		diags.AddError(
			"Error Creating Files Style",
			"Could not marshal logo to JSON: "+err.Error(),
		)
	}
	state.Logo = types.StringValue(string(respLogo))
	respThumbnail, err := json.Marshal(style.Thumbnail)
	if err != nil {
		diags.AddError(
			"Error Creating Files Style",
			"Could not marshal thumbnail to JSON: "+err.Error(),
		)
	}
	state.Thumbnail = types.StringValue(string(respThumbnail))

	return
}
