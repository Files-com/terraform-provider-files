package provider

import (
	"context"
	"fmt"

	files_sdk "github.com/Files-com/files-sdk-go/v3"
	history_export "github.com/Files-com/files-sdk-go/v3/historyexport"
	"github.com/Files-com/terraform-provider-files/lib"
	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/datasource/schema"
	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

var (
	_ datasource.DataSource              = &historyExportDataSource{}
	_ datasource.DataSourceWithConfigure = &historyExportDataSource{}
)

func NewHistoryExportDataSource() datasource.DataSource {
	return &historyExportDataSource{}
}

type historyExportDataSource struct {
	client *history_export.Client
}

type historyExportDataSourceModel struct {
	Id                       types.Int64  `tfsdk:"id"`
	HistoryVersion           types.String `tfsdk:"history_version"`
	StartAt                  types.String `tfsdk:"start_at"`
	EndAt                    types.String `tfsdk:"end_at"`
	Status                   types.String `tfsdk:"status"`
	QueryAction              types.String `tfsdk:"query_action"`
	QueryInterface           types.String `tfsdk:"query_interface"`
	QueryUserId              types.String `tfsdk:"query_user_id"`
	QueryFileId              types.String `tfsdk:"query_file_id"`
	QueryParentId            types.String `tfsdk:"query_parent_id"`
	QueryPath                types.String `tfsdk:"query_path"`
	QueryFolder              types.String `tfsdk:"query_folder"`
	QuerySrc                 types.String `tfsdk:"query_src"`
	QueryDestination         types.String `tfsdk:"query_destination"`
	QueryIp                  types.String `tfsdk:"query_ip"`
	QueryUsername            types.String `tfsdk:"query_username"`
	QueryFailureType         types.String `tfsdk:"query_failure_type"`
	QueryTargetId            types.String `tfsdk:"query_target_id"`
	QueryTargetName          types.String `tfsdk:"query_target_name"`
	QueryTargetPermission    types.String `tfsdk:"query_target_permission"`
	QueryTargetUserId        types.String `tfsdk:"query_target_user_id"`
	QueryTargetUsername      types.String `tfsdk:"query_target_username"`
	QueryTargetPlatform      types.String `tfsdk:"query_target_platform"`
	QueryTargetPermissionSet types.String `tfsdk:"query_target_permission_set"`
	ResultsUrl               types.String `tfsdk:"results_url"`
}

func (r *historyExportDataSource) Configure(_ context.Context, req datasource.ConfigureRequest, resp *datasource.ConfigureResponse) {
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

	r.client = &history_export.Client{Config: sdk_config}
}

func (r *historyExportDataSource) Metadata(_ context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_history_export"
}

func (r *historyExportDataSource) Schema(_ context.Context, _ datasource.SchemaRequest, resp *datasource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Description: "The History Export resource on the API is used to export historical action (history) logs.\n\n\n\nAll queries against the archive must be submitted as Exports. (Even our Web UI creates an Export behind\n\nthe scenes.)\n\n\n\nWe use Amazon Athena behind the scenes for processing these queries, and as such, have powerful\n\nsearch capabilities. We've done our best to expose search capabilities via this History Export API.\n\n\n\nIn any query field in this API, you may specify multiple values separated by commas. That means that commas\n\ncannot be searched for themselves, and neither can single quotation marks.\n\n\n\nWe do not currently partition data by date on the backend, so all queries result in a full scan of the entire\n\ndata lake. This means that all queries will take about the same amount of time to complete, and we incur about\n\nthe same cost per query internally. We don't typically bill our customers for these queries, assuming\n\nusage is occasional and manual.\n\n\n\nIf you intend to use this API for high volume or automated use, please contact us with more information\n\nabout your use case. We may decide to change the backend data schema to match your use case more closely, and\n\nwe may also need to charge an additional cost per query.\n\n\n\n## Example History Queries\n\n\n\n* History for a user: `{ \"query_user_id\": 123 }`\n\n* History for a range of time: `{ \"start_at\": \"2021-03-18 12:00:00\", \"end_at\": \"2021-03-19 12:00:00\" }`\n\n* History of logins and failed logins: `{ \"query_action\": \"login,failedlogin\" }`\n\n* A Complex query: `{ \"query_folder\": \"uploads\", \"query_action\": \"create,copy,move\", \"start_at\": \"2021-03-18 12:00:00\", \"end_at\": \"2021-03-19 12:00:00\" }`",
		Attributes: map[string]schema.Attribute{
			"id": schema.Int64Attribute{
				Description: "History Export ID",
				Required:    true,
			},
			"history_version": schema.StringAttribute{
				Description: "Version of the history for the export.",
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
				Description: "Status of export.  Will be: `building`, `ready`, or `failed`",
				Computed:    true,
			},
			"query_action": schema.StringAttribute{
				Description: "Filter results by this this action type. Valid values: `create`, `read`, `update`, `destroy`, `move`, `login`, `failedlogin`, `copy`, `user_create`, `user_update`, `user_destroy`, `group_create`, `group_update`, `group_destroy`, `permission_create`, `permission_destroy`, `api_key_create`, `api_key_update`, `api_key_destroy`",
				Computed:    true,
			},
			"query_interface": schema.StringAttribute{
				Description: "Filter results by this this interface type. Valid values: `web`, `ftp`, `robot`, `jsapi`, `webdesktopapi`, `sftp`, `dav`, `desktop`, `restapi`, `scim`, `office`, `mobile`, `as2`, `inbound_email`, `remote`",
				Computed:    true,
			},
			"query_user_id": schema.StringAttribute{
				Description: "Return results that are actions performed by the user indiciated by this User ID",
				Computed:    true,
			},
			"query_file_id": schema.StringAttribute{
				Description: "Return results that are file actions related to the file indicated by this File ID",
				Computed:    true,
			},
			"query_parent_id": schema.StringAttribute{
				Description: "Return results that are file actions inside the parent folder specified by this folder ID",
				Computed:    true,
			},
			"query_path": schema.StringAttribute{
				Description: "Return results that are file actions related to paths matching this pattern.",
				Computed:    true,
			},
			"query_folder": schema.StringAttribute{
				Description: "Return results that are file actions related to files or folders inside folder paths matching this pattern.",
				Computed:    true,
			},
			"query_src": schema.StringAttribute{
				Description: "Return results that are file moves originating from paths matching this pattern.",
				Computed:    true,
			},
			"query_destination": schema.StringAttribute{
				Description: "Return results that are file moves with paths matching this pattern as destination.",
				Computed:    true,
			},
			"query_ip": schema.StringAttribute{
				Description: "Filter results by this IP address.",
				Computed:    true,
			},
			"query_username": schema.StringAttribute{
				Description: "Filter results by this username.",
				Computed:    true,
			},
			"query_failure_type": schema.StringAttribute{
				Description: "If searching for Histories about login failures, this parameter restricts results to failures of this specific type.  Valid values: `expired_trial`, `account_overdue`, `locked_out`, `ip_mismatch`, `password_mismatch`, `site_mismatch`, `username_not_found`, `none`, `no_ftp_permission`, `no_web_permission`, `no_directory`, `errno_enoent`, `no_sftp_permission`, `no_dav_permission`, `no_restapi_permission`, `key_mismatch`, `region_mismatch`, `expired_access`, `desktop_ip_mismatch`, `desktop_api_key_not_used_quickly_enough`, `disabled`, `country_mismatch`, `insecure_ftp`, `insecure_cipher`, `rate_limited`",
				Computed:    true,
			},
			"query_target_id": schema.StringAttribute{
				Description: "If searching for Histories about specific objects (such as Users, or API Keys), this paremeter restricts results to objects that match this ID.",
				Computed:    true,
			},
			"query_target_name": schema.StringAttribute{
				Description: "If searching for Histories about Users, Groups or other objects with names, this parameter restricts results to objects with this name/username.",
				Computed:    true,
			},
			"query_target_permission": schema.StringAttribute{
				Description: "If searching for Histories about Permisisons, this parameter restricts results to permissions of this level.",
				Computed:    true,
			},
			"query_target_user_id": schema.StringAttribute{
				Description: "If searching for Histories about API keys, this parameter restricts results to API keys created by/for this user ID.",
				Computed:    true,
			},
			"query_target_username": schema.StringAttribute{
				Description: "If searching for Histories about API keys, this parameter restricts results to API keys created by/for this username.",
				Computed:    true,
			},
			"query_target_platform": schema.StringAttribute{
				Description: "If searching for Histories about API keys, this parameter restricts results to API keys associated with this platform.",
				Computed:    true,
			},
			"query_target_permission_set": schema.StringAttribute{
				Description: "If searching for Histories about API keys, this parameter restricts results to API keys with this permission set.",
				Computed:    true,
			},
			"results_url": schema.StringAttribute{
				Description: "If `status` is `ready`, this will be a URL where all the results can be downloaded at once as a CSV.",
				Computed:    true,
			},
		},
	}
}

func (r *historyExportDataSource) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
	var data historyExportDataSourceModel
	diags := req.Config.Get(ctx, &data)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	paramsHistoryExportFind := files_sdk.HistoryExportFindParams{}
	paramsHistoryExportFind.Id = data.Id.ValueInt64()

	historyExport, err := r.client.Find(paramsHistoryExportFind, files_sdk.WithContext(ctx))
	if err != nil {
		resp.Diagnostics.AddError(
			"Error Reading Files HistoryExport",
			"Could not read history_export id "+fmt.Sprint(data.Id.ValueInt64())+": "+err.Error(),
		)
		return
	}

	diags = r.populateDataSourceModel(ctx, historyExport, &data)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	diags = resp.State.Set(ctx, data)
	resp.Diagnostics.Append(diags...)
}

func (r *historyExportDataSource) populateDataSourceModel(ctx context.Context, historyExport files_sdk.HistoryExport, state *historyExportDataSourceModel) (diags diag.Diagnostics) {
	state.Id = types.Int64Value(historyExport.Id)
	state.HistoryVersion = types.StringValue(historyExport.HistoryVersion)
	if err := lib.TimeToStringType(ctx, path.Root("start_at"), historyExport.StartAt, &state.StartAt); err != nil {
		diags.AddError(
			"Error Creating Files HistoryExport",
			"Could not convert state start_at to string: "+err.Error(),
		)
	}
	if err := lib.TimeToStringType(ctx, path.Root("end_at"), historyExport.EndAt, &state.EndAt); err != nil {
		diags.AddError(
			"Error Creating Files HistoryExport",
			"Could not convert state end_at to string: "+err.Error(),
		)
	}
	state.Status = types.StringValue(historyExport.Status)
	state.QueryAction = types.StringValue(historyExport.QueryAction)
	state.QueryInterface = types.StringValue(historyExport.QueryInterface)
	state.QueryUserId = types.StringValue(historyExport.QueryUserId)
	state.QueryFileId = types.StringValue(historyExport.QueryFileId)
	state.QueryParentId = types.StringValue(historyExport.QueryParentId)
	state.QueryPath = types.StringValue(historyExport.QueryPath)
	state.QueryFolder = types.StringValue(historyExport.QueryFolder)
	state.QuerySrc = types.StringValue(historyExport.QuerySrc)
	state.QueryDestination = types.StringValue(historyExport.QueryDestination)
	state.QueryIp = types.StringValue(historyExport.QueryIp)
	state.QueryUsername = types.StringValue(historyExport.QueryUsername)
	state.QueryFailureType = types.StringValue(historyExport.QueryFailureType)
	state.QueryTargetId = types.StringValue(historyExport.QueryTargetId)
	state.QueryTargetName = types.StringValue(historyExport.QueryTargetName)
	state.QueryTargetPermission = types.StringValue(historyExport.QueryTargetPermission)
	state.QueryTargetUserId = types.StringValue(historyExport.QueryTargetUserId)
	state.QueryTargetUsername = types.StringValue(historyExport.QueryTargetUsername)
	state.QueryTargetPlatform = types.StringValue(historyExport.QueryTargetPlatform)
	state.QueryTargetPermissionSet = types.StringValue(historyExport.QueryTargetPermissionSet)
	state.ResultsUrl = types.StringValue(historyExport.ResultsUrl)

	return
}
