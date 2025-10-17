package provider

import (
	"context"
	"fmt"

	files_sdk "github.com/Files-com/files-sdk-go/v3"
	scim_log "github.com/Files-com/files-sdk-go/v3/scimlog"
	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/datasource/schema"
	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

var (
	_ datasource.DataSource              = &scimLogDataSource{}
	_ datasource.DataSourceWithConfigure = &scimLogDataSource{}
)

func NewScimLogDataSource() datasource.DataSource {
	return &scimLogDataSource{}
}

type scimLogDataSource struct {
	client *scim_log.Client
}

type scimLogDataSourceModel struct {
	Id               types.Int64  `tfsdk:"id"`
	CreatedAt        types.String `tfsdk:"created_at"`
	RequestPath      types.String `tfsdk:"request_path"`
	RequestMethod    types.String `tfsdk:"request_method"`
	HttpResponseCode types.String `tfsdk:"http_response_code"`
	UserAgent        types.String `tfsdk:"user_agent"`
	RequestJson      types.String `tfsdk:"request_json"`
	ResponseJson     types.String `tfsdk:"response_json"`
}

func (r *scimLogDataSource) Configure(_ context.Context, req datasource.ConfigureRequest, resp *datasource.ConfigureResponse) {
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

	r.client = &scim_log.Client{Config: sdk_config}
}

func (r *scimLogDataSource) Metadata(_ context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_scim_log"
}

func (r *scimLogDataSource) Schema(_ context.Context, _ datasource.SchemaRequest, resp *datasource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Description: "A SCIM log entry represents a single SCIM request made to the system. It includes the request made and response provided to the SCIM client.",
		Attributes: map[string]schema.Attribute{
			"id": schema.Int64Attribute{
				Description: "The unique ID of this SCIM request.",
				Required:    true,
			},
			"created_at": schema.StringAttribute{
				Description: "The date and time when this SCIM request occurred.",
				Computed:    true,
			},
			"request_path": schema.StringAttribute{
				Description: "The path portion of the URL requested.",
				Computed:    true,
			},
			"request_method": schema.StringAttribute{
				Description: "The HTTP method used for this request.",
				Computed:    true,
			},
			"http_response_code": schema.StringAttribute{
				Description: "The HTTP response code returned for this request.",
				Computed:    true,
			},
			"user_agent": schema.StringAttribute{
				Description: "The User-Agent header sent with the request.",
				Computed:    true,
			},
			"request_json": schema.StringAttribute{
				Description: "The JSON payload sent with the request.",
				Computed:    true,
			},
			"response_json": schema.StringAttribute{
				Description: "The JSON payload returned in the response.",
				Computed:    true,
			},
		},
	}
}

func (r *scimLogDataSource) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
	var data scimLogDataSourceModel
	diags := req.Config.Get(ctx, &data)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	paramsScimLogFind := files_sdk.ScimLogFindParams{}
	paramsScimLogFind.Id = data.Id.ValueInt64()

	scimLog, err := r.client.Find(paramsScimLogFind, files_sdk.WithContext(ctx))
	if err != nil {
		resp.Diagnostics.AddError(
			"Error Reading Files ScimLog",
			"Could not read scim_log id "+fmt.Sprint(data.Id.ValueInt64())+": "+err.Error(),
		)
		return
	}

	diags = r.populateDataSourceModel(ctx, scimLog, &data)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	diags = resp.State.Set(ctx, data)
	resp.Diagnostics.Append(diags...)
}

func (r *scimLogDataSource) populateDataSourceModel(ctx context.Context, scimLog files_sdk.ScimLog, state *scimLogDataSourceModel) (diags diag.Diagnostics) {
	state.Id = types.Int64Value(scimLog.Id)
	state.CreatedAt = types.StringValue(scimLog.CreatedAt)
	state.RequestPath = types.StringValue(scimLog.RequestPath)
	state.RequestMethod = types.StringValue(scimLog.RequestMethod)
	state.HttpResponseCode = types.StringValue(scimLog.HttpResponseCode)
	state.UserAgent = types.StringValue(scimLog.UserAgent)
	state.RequestJson = types.StringValue(scimLog.RequestJson)
	state.ResponseJson = types.StringValue(scimLog.ResponseJson)

	return
}
