package provider

import (
	"context"
	"fmt"

	files_sdk "github.com/Files-com/files-sdk-go/v3"
	priority "github.com/Files-com/files-sdk-go/v3/priority"
	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/datasource/schema"
	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

var (
	_ datasource.DataSource              = &priorityDataSource{}
	_ datasource.DataSourceWithConfigure = &priorityDataSource{}
)

func NewPriorityDataSource() datasource.DataSource {
	return &priorityDataSource{}
}

type priorityDataSource struct {
	client *priority.Client
}

type priorityDataSourceModel struct {
	Path  types.String `tfsdk:"path"`
	Color types.String `tfsdk:"color"`
}

func (r *priorityDataSource) Configure(_ context.Context, req datasource.ConfigureRequest, resp *datasource.ConfigureResponse) {
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

	r.client = &priority.Client{Config: sdk_config}
}

func (r *priorityDataSource) Metadata(_ context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_priority"
}

func (r *priorityDataSource) Schema(_ context.Context, _ datasource.SchemaRequest, resp *datasource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Description: "",
		Attributes: map[string]schema.Attribute{
			"path": schema.StringAttribute{
				Description: "The path corresponding to the priority color. This must be slash-delimited, but it must neither start nor end with a slash. Maximum of 5000 characters.",
				Required:    true,
			},
			"color": schema.StringAttribute{
				Description: "The priority color",
				Computed:    true,
			},
		},
	}
}

func (r *priorityDataSource) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
	var data priorityDataSourceModel
	diags := req.Config.Get(ctx, &data)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	paramsPriorityList := files_sdk.PriorityListParams{}
	paramsPriorityList.Path = data.Path.ValueString()

	priorityIt, err := r.client.List(paramsPriorityList, files_sdk.WithContext(ctx))
	if err != nil {
		resp.Diagnostics.AddError(
			"Error Reading Files Priority",
			"Could not read priority path "+fmt.Sprint(data.Path.ValueString())+": "+err.Error(),
		)
		return
	}

	var priority *files_sdk.Priority
	for priorityIt.Next() {
		entry := priorityIt.Priority()
		if entry.Path == data.Path.ValueString() {
			priority = &entry
			break
		}
	}

	if priority == nil {
		resp.Diagnostics.AddError(
			"Error Reading Files Priority",
			"Could not find priority path "+fmt.Sprint(data.Path.ValueString()),
		)
		return
	}

	diags = r.populateDataSourceModel(ctx, *priority, &data)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	diags = resp.State.Set(ctx, data)
	resp.Diagnostics.Append(diags...)
}

func (r *priorityDataSource) populateDataSourceModel(ctx context.Context, priority files_sdk.Priority, state *priorityDataSourceModel) (diags diag.Diagnostics) {
	state.Path = types.StringValue(priority.Path)
	state.Color = types.StringValue(priority.Color)

	return
}
