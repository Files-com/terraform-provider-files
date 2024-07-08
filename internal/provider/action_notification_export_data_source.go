package provider

import (
	"context"
	"fmt"

	files_sdk "github.com/Files-com/files-sdk-go/v3"
	action_notification_export "github.com/Files-com/files-sdk-go/v3/actionnotificationexport"
	"github.com/Files-com/terraform-provider-files/lib"
	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/datasource/schema"
	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

var (
	_ datasource.DataSource              = &actionNotificationExportDataSource{}
	_ datasource.DataSourceWithConfigure = &actionNotificationExportDataSource{}
)

func NewActionNotificationExportDataSource() datasource.DataSource {
	return &actionNotificationExportDataSource{}
}

type actionNotificationExportDataSource struct {
	client *action_notification_export.Client
}

type actionNotificationExportDataSourceModel struct {
	Id                 types.Int64  `tfsdk:"id"`
	ExportVersion      types.String `tfsdk:"export_version"`
	StartAt            types.String `tfsdk:"start_at"`
	EndAt              types.String `tfsdk:"end_at"`
	Status             types.String `tfsdk:"status"`
	QueryPath          types.String `tfsdk:"query_path"`
	QueryFolder        types.String `tfsdk:"query_folder"`
	QueryMessage       types.String `tfsdk:"query_message"`
	QueryRequestMethod types.String `tfsdk:"query_request_method"`
	QueryRequestUrl    types.String `tfsdk:"query_request_url"`
	QueryStatus        types.String `tfsdk:"query_status"`
	QuerySuccess       types.Bool   `tfsdk:"query_success"`
	ResultsUrl         types.String `tfsdk:"results_url"`
}

func (r *actionNotificationExportDataSource) Configure(_ context.Context, req datasource.ConfigureRequest, resp *datasource.ConfigureResponse) {
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

	r.client = &action_notification_export.Client{Config: sdk_config}
}

func (r *actionNotificationExportDataSource) Metadata(_ context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_action_notification_export"
}

func (r *actionNotificationExportDataSource) Schema(_ context.Context, _ datasource.SchemaRequest, resp *datasource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Description: "The Action Notification Export API provides access to outgoing webhook logs. Querying webhook logs is a little different than other APIs.\n\n\n\nAll queries against the archive must be submitted as Exports. (Even our Web UI creates an Export behind the scenes.)\n\n\n\nIn any query field in this API, you may specify multiple values separated by commas. That means that commas\n\ncannot be searched for themselves, and neither can single quotation marks.\n\n\n\nUse the following steps to complete an export:\n\n\n\n1. Initiate the export by using the Create Action Notification Export endpoint. Non Site Admins must query by folder or path.\n\n2. Using the `id` from the response to step 1, poll the Show Action Notification Export endpoint. Check the `status` field until it is `ready`.\n\n3. You can download the results of the export as a CSV file using the `results_url` field in the response from step 2. If you want to page through the records in JSON format, use the List Action Notification Export Results endpoint, passing the `id` that you got in step 1 as the `action_notification_export_id` parameter. Check the `X-Files-Cursor-Next` header to see if there are more records available, and resubmit the same request with a `cursor` parameter to fetch the next page of results. Unlike most API Endpoints, this endpoint does not provide `X-Files-Cursor-Prev` cursors allowing reverse pagination through the results. This is due to limitations in Amazon Athena, the underlying data lake for these records.\n\n\n\nIf you intend to use this API for high volume or automated use, please contact us with more information\n\nabout your use case.\n\n\n\n## Example Queries\n\n\n\n* History for a folder: `{ \"query_folder\": \"path/to/folder\" }`\n\n* History for a range of time: `{ \"start_at\": \"2021-03-18 12:00:00\", \"end_at\": \"2021-03-19 12:00:00\" }`\n\n* History of all notifications that used GET or POST: `{ \"query_request_method\": \"GET,POST\" }`",
		Attributes: map[string]schema.Attribute{
			"id": schema.Int64Attribute{
				Description: "History Export ID",
				Required:    true,
			},
			"export_version": schema.StringAttribute{
				Description: "Version of the underlying records for the export.",
				Computed:    true,
			},
			"start_at": schema.StringAttribute{
				Description: "Start date/time of export range.",
				Computed:    true,
			},
			"end_at": schema.StringAttribute{
				Description: "End date/time of export range.",
				Computed:    true,
			},
			"status": schema.StringAttribute{
				Description: "Status of export.  Valid values: `building`, `ready`, or `failed`",
				Computed:    true,
			},
			"query_path": schema.StringAttribute{
				Description: "Return notifications that were triggered by actions on this specific path.",
				Computed:    true,
			},
			"query_folder": schema.StringAttribute{
				Description: "Return notifications that were triggered by actions in this folder.",
				Computed:    true,
			},
			"query_message": schema.StringAttribute{
				Description: "Error message associated with the request, if any.",
				Computed:    true,
			},
			"query_request_method": schema.StringAttribute{
				Description: "The HTTP request method used by the webhook.",
				Computed:    true,
			},
			"query_request_url": schema.StringAttribute{
				Description: "The target webhook URL.",
				Computed:    true,
			},
			"query_status": schema.StringAttribute{
				Description: "The HTTP status returned from the server in response to the webhook request.",
				Computed:    true,
			},
			"query_success": schema.BoolAttribute{
				Description: "true if the webhook request succeeded (i.e. returned a 200 or 204 response status). false otherwise.",
				Computed:    true,
			},
			"results_url": schema.StringAttribute{
				Description: "If `status` is `ready`, this will be a URL where all the results can be downloaded at once as a CSV.",
				Computed:    true,
			},
		},
	}
}

func (r *actionNotificationExportDataSource) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
	var data actionNotificationExportDataSourceModel
	diags := req.Config.Get(ctx, &data)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	paramsActionNotificationExportFind := files_sdk.ActionNotificationExportFindParams{}
	paramsActionNotificationExportFind.Id = data.Id.ValueInt64()

	actionNotificationExport, err := r.client.Find(paramsActionNotificationExportFind, files_sdk.WithContext(ctx))
	if err != nil {
		resp.Diagnostics.AddError(
			"Error Reading Files ActionNotificationExport",
			"Could not read action_notification_export id "+fmt.Sprint(data.Id.ValueInt64())+": "+err.Error(),
		)
		return
	}

	diags = r.populateDataSourceModel(ctx, actionNotificationExport, &data)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	diags = resp.State.Set(ctx, data)
	resp.Diagnostics.Append(diags...)
}

func (r *actionNotificationExportDataSource) populateDataSourceModel(ctx context.Context, actionNotificationExport files_sdk.ActionNotificationExport, state *actionNotificationExportDataSourceModel) (diags diag.Diagnostics) {
	state.Id = types.Int64Value(actionNotificationExport.Id)
	state.ExportVersion = types.StringValue(actionNotificationExport.ExportVersion)
	if err := lib.TimeToStringType(ctx, path.Root("start_at"), actionNotificationExport.StartAt, &state.StartAt); err != nil {
		diags.AddError(
			"Error Creating Files ActionNotificationExport",
			"Could not convert state start_at to string: "+err.Error(),
		)
	}
	if err := lib.TimeToStringType(ctx, path.Root("end_at"), actionNotificationExport.EndAt, &state.EndAt); err != nil {
		diags.AddError(
			"Error Creating Files ActionNotificationExport",
			"Could not convert state end_at to string: "+err.Error(),
		)
	}
	state.Status = types.StringValue(actionNotificationExport.Status)
	state.QueryPath = types.StringValue(actionNotificationExport.QueryPath)
	state.QueryFolder = types.StringValue(actionNotificationExport.QueryFolder)
	state.QueryMessage = types.StringValue(actionNotificationExport.QueryMessage)
	state.QueryRequestMethod = types.StringValue(actionNotificationExport.QueryRequestMethod)
	state.QueryRequestUrl = types.StringValue(actionNotificationExport.QueryRequestUrl)
	state.QueryStatus = types.StringValue(actionNotificationExport.QueryStatus)
	state.QuerySuccess = types.BoolPointerValue(actionNotificationExport.QuerySuccess)
	state.ResultsUrl = types.StringValue(actionNotificationExport.ResultsUrl)

	return
}
