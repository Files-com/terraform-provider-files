package provider

import (
	"context"
	"fmt"

	files_sdk "github.com/Files-com/files-sdk-go/v3"
	request "github.com/Files-com/files-sdk-go/v3/request"
	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/datasource/schema"
	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

var (
	_ datasource.DataSource              = &requestDataSource{}
	_ datasource.DataSourceWithConfigure = &requestDataSource{}
)

func NewRequestDataSource() datasource.DataSource {
	return &requestDataSource{}
}

type requestDataSource struct {
	client *request.Client
}

type requestDataSourceModel struct {
	Id              types.Int64  `tfsdk:"id"`
	Path            types.String `tfsdk:"path"`
	Source          types.String `tfsdk:"source"`
	Destination     types.String `tfsdk:"destination"`
	AutomationId    types.String `tfsdk:"automation_id"`
	UserDisplayName types.String `tfsdk:"user_display_name"`
}

func (r *requestDataSource) Configure(_ context.Context, req datasource.ConfigureRequest, resp *datasource.ConfigureResponse) {
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

	r.client = &request.Client{Config: sdk_config}
}

func (r *requestDataSource) Metadata(_ context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_request"
}

func (r *requestDataSource) Schema(_ context.Context, _ datasource.SchemaRequest, resp *datasource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Description: "A Request represents a file that *should* be uploaded by a specific user or group.\n\n\n\nRequests can either be manually created and managed, or managed automatically by an Automation.",
		Attributes: map[string]schema.Attribute{
			"id": schema.Int64Attribute{
				Description: "Request ID",
				Required:    true,
			},
			"path": schema.StringAttribute{
				Description: "Folder path. This must be slash-delimited, but it must neither start nor end with a slash. Maximum of 5000 characters.",
				Computed:    true,
			},
			"source": schema.StringAttribute{
				Description: "Source filename, if applicable",
				Computed:    true,
			},
			"destination": schema.StringAttribute{
				Description: "Destination filename",
				Computed:    true,
			},
			"automation_id": schema.StringAttribute{
				Description: "ID of automation that created request",
				Computed:    true,
			},
			"user_display_name": schema.StringAttribute{
				Description: "User making the request (if applicable)",
				Computed:    true,
			},
		},
	}
}

func (r *requestDataSource) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
	var data requestDataSourceModel
	diags := req.Config.Get(ctx, &data)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	paramsRequestList := files_sdk.RequestListParams{}

	requestIt, err := r.client.List(paramsRequestList, files_sdk.WithContext(ctx))
	if err != nil {
		resp.Diagnostics.AddError(
			"Error Reading Files Request",
			"Could not read request id "+fmt.Sprint(data.Id.ValueInt64())+": "+err.Error(),
		)
		return
	}

	var request *files_sdk.Request
	for requestIt.Next() {
		entry := requestIt.Request()
		if entry.Id == data.Id.ValueInt64() {
			request = &entry
			break
		}
	}

	if request == nil {
		resp.Diagnostics.AddError(
			"Error Reading Files Request",
			"Could not find request id "+fmt.Sprint(data.Id.ValueInt64()),
		)
		return
	}

	diags = r.populateDataSourceModel(ctx, *request, &data)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	diags = resp.State.Set(ctx, data)
	resp.Diagnostics.Append(diags...)
}

func (r *requestDataSource) populateDataSourceModel(ctx context.Context, request files_sdk.Request, state *requestDataSourceModel) (diags diag.Diagnostics) {
	state.Id = types.Int64Value(request.Id)
	state.Path = types.StringValue(request.Path)
	state.Source = types.StringValue(request.Source)
	state.Destination = types.StringValue(request.Destination)
	state.AutomationId = types.StringValue(request.AutomationId)
	state.UserDisplayName = types.StringValue(request.UserDisplayName)

	return
}
