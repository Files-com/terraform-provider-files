package provider

import (
	"context"
	"fmt"

	files_sdk "github.com/Files-com/files-sdk-go/v3"
	siem_http_destination "github.com/Files-com/files-sdk-go/v3/siemhttpdestination"
	"github.com/Files-com/terraform-provider-files/lib"
	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/datasource/schema"
	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

var (
	_ datasource.DataSource              = &siemHttpDestinationDataSource{}
	_ datasource.DataSourceWithConfigure = &siemHttpDestinationDataSource{}
)

func NewSiemHttpDestinationDataSource() datasource.DataSource {
	return &siemHttpDestinationDataSource{}
}

type siemHttpDestinationDataSource struct {
	client *siem_http_destination.Client
}

type siemHttpDestinationDataSourceModel struct {
	Id                                            types.Int64   `tfsdk:"id"`
	Name                                          types.String  `tfsdk:"name"`
	DestinationType                               types.String  `tfsdk:"destination_type"`
	DestinationUrl                                types.String  `tfsdk:"destination_url"`
	FileDestinationPath                           types.String  `tfsdk:"file_destination_path"`
	FileFormat                                    types.String  `tfsdk:"file_format"`
	FileIntervalMinutes                           types.Int64   `tfsdk:"file_interval_minutes"`
	AdditionalHeaders                             types.Dynamic `tfsdk:"additional_headers"`
	SendingActive                                 types.Bool    `tfsdk:"sending_active"`
	GenericPayloadType                            types.String  `tfsdk:"generic_payload_type"`
	SplunkTokenMasked                             types.String  `tfsdk:"splunk_token_masked"`
	AzureDcrImmutableId                           types.String  `tfsdk:"azure_dcr_immutable_id"`
	AzureStreamName                               types.String  `tfsdk:"azure_stream_name"`
	AzureOauthClientCredentialsTenantId           types.String  `tfsdk:"azure_oauth_client_credentials_tenant_id"`
	AzureOauthClientCredentialsClientId           types.String  `tfsdk:"azure_oauth_client_credentials_client_id"`
	AzureOauthClientCredentialsClientSecretMasked types.String  `tfsdk:"azure_oauth_client_credentials_client_secret_masked"`
	QradarUsername                                types.String  `tfsdk:"qradar_username"`
	QradarPasswordMasked                          types.String  `tfsdk:"qradar_password_masked"`
	SolarWindsTokenMasked                         types.String  `tfsdk:"solar_winds_token_masked"`
	NewRelicApiKeyMasked                          types.String  `tfsdk:"new_relic_api_key_masked"`
	DatadogApiKeyMasked                           types.String  `tfsdk:"datadog_api_key_masked"`
	ActionSendEnabled                             types.Bool    `tfsdk:"action_send_enabled"`
	ActionEntriesSent                             types.Int64   `tfsdk:"action_entries_sent"`
	SftpActionSendEnabled                         types.Bool    `tfsdk:"sftp_action_send_enabled"`
	SftpActionEntriesSent                         types.Int64   `tfsdk:"sftp_action_entries_sent"`
	FtpActionSendEnabled                          types.Bool    `tfsdk:"ftp_action_send_enabled"`
	FtpActionEntriesSent                          types.Int64   `tfsdk:"ftp_action_entries_sent"`
	WebDavActionSendEnabled                       types.Bool    `tfsdk:"web_dav_action_send_enabled"`
	WebDavActionEntriesSent                       types.Int64   `tfsdk:"web_dav_action_entries_sent"`
	SyncSendEnabled                               types.Bool    `tfsdk:"sync_send_enabled"`
	SyncEntriesSent                               types.Int64   `tfsdk:"sync_entries_sent"`
	OutboundConnectionSendEnabled                 types.Bool    `tfsdk:"outbound_connection_send_enabled"`
	OutboundConnectionEntriesSent                 types.Int64   `tfsdk:"outbound_connection_entries_sent"`
	AutomationSendEnabled                         types.Bool    `tfsdk:"automation_send_enabled"`
	AutomationEntriesSent                         types.Int64   `tfsdk:"automation_entries_sent"`
	ApiRequestSendEnabled                         types.Bool    `tfsdk:"api_request_send_enabled"`
	ApiRequestEntriesSent                         types.Int64   `tfsdk:"api_request_entries_sent"`
	PublicHostingRequestSendEnabled               types.Bool    `tfsdk:"public_hosting_request_send_enabled"`
	PublicHostingRequestEntriesSent               types.Int64   `tfsdk:"public_hosting_request_entries_sent"`
	EmailSendEnabled                              types.Bool    `tfsdk:"email_send_enabled"`
	EmailEntriesSent                              types.Int64   `tfsdk:"email_entries_sent"`
	ExavaultApiRequestSendEnabled                 types.Bool    `tfsdk:"exavault_api_request_send_enabled"`
	ExavaultApiRequestEntriesSent                 types.Int64   `tfsdk:"exavault_api_request_entries_sent"`
	SettingsChangeSendEnabled                     types.Bool    `tfsdk:"settings_change_send_enabled"`
	SettingsChangeEntriesSent                     types.Int64   `tfsdk:"settings_change_entries_sent"`
	LastHttpCallTargetType                        types.String  `tfsdk:"last_http_call_target_type"`
	LastHttpCallSuccess                           types.Bool    `tfsdk:"last_http_call_success"`
	LastHttpCallResponseCode                      types.Int64   `tfsdk:"last_http_call_response_code"`
	LastHttpCallResponseBody                      types.String  `tfsdk:"last_http_call_response_body"`
	LastHttpCallErrorMessage                      types.String  `tfsdk:"last_http_call_error_message"`
	LastHttpCallTime                              types.String  `tfsdk:"last_http_call_time"`
	LastHttpCallDurationMs                        types.Int64   `tfsdk:"last_http_call_duration_ms"`
	MostRecentHttpCallSuccessTime                 types.String  `tfsdk:"most_recent_http_call_success_time"`
	ConnectionTestEntry                           types.String  `tfsdk:"connection_test_entry"`
}

func (r *siemHttpDestinationDataSource) Configure(_ context.Context, req datasource.ConfigureRequest, resp *datasource.ConfigureResponse) {
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

	r.client = &siem_http_destination.Client{Config: sdk_config}
}

func (r *siemHttpDestinationDataSource) Metadata(_ context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_siem_http_destination"
}

func (r *siemHttpDestinationDataSource) Schema(_ context.Context, _ datasource.SchemaRequest, resp *datasource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Description: "",
		Attributes: map[string]schema.Attribute{
			"id": schema.Int64Attribute{
				Description: "SIEM HTTP Destination ID",
				Required:    true,
			},
			"name": schema.StringAttribute{
				Description: "Name for this Destination",
				Computed:    true,
			},
			"destination_type": schema.StringAttribute{
				Description: "Destination Type",
				Computed:    true,
			},
			"destination_url": schema.StringAttribute{
				Description: "Destination Url",
				Computed:    true,
			},
			"file_destination_path": schema.StringAttribute{
				Description: "Applicable only for destination type: file. Destination folder path on Files.com.",
				Computed:    true,
			},
			"file_format": schema.StringAttribute{
				Description: "Applicable only for destination type: file. Generated file format.",
				Computed:    true,
			},
			"file_interval_minutes": schema.Int64Attribute{
				Description: "Applicable only for destination type: file. Interval, in minutes, between file deliveries.",
				Computed:    true,
			},
			"additional_headers": schema.DynamicAttribute{
				Description: "Additional HTTP Headers included in calls to the destination URL",
				Computed:    true,
			},
			"sending_active": schema.BoolAttribute{
				Description: "Whether this SIEM HTTP Destination is currently being sent to or not",
				Computed:    true,
			},
			"generic_payload_type": schema.StringAttribute{
				Description: "Applicable only for destination type: generic. Indicates the type of HTTP body. Can be json_newline or json_array. json_newline is multiple log entries as JSON separated by newlines. json_array is a single JSON array containing multiple log entries as JSON.",
				Computed:    true,
			},
			"splunk_token_masked": schema.StringAttribute{
				Description: "Applicable only for destination type: splunk. Authentication token provided by Splunk.",
				Computed:    true,
			},
			"azure_dcr_immutable_id": schema.StringAttribute{
				Description: "Applicable only for destination types: azure, azure_legacy. Immutable ID of the Data Collection Rule.",
				Computed:    true,
			},
			"azure_stream_name": schema.StringAttribute{
				Description: "Applicable only for destination type: azure. Name of the stream in the DCR that represents the destination table.",
				Computed:    true,
			},
			"azure_oauth_client_credentials_tenant_id": schema.StringAttribute{
				Description: "Applicable only for destination types: azure, azure_legacy. Client Credentials OAuth Tenant ID.",
				Computed:    true,
			},
			"azure_oauth_client_credentials_client_id": schema.StringAttribute{
				Description: "Applicable only for destination types: azure, azure_legacy. Client Credentials OAuth Client ID.",
				Computed:    true,
			},
			"azure_oauth_client_credentials_client_secret_masked": schema.StringAttribute{
				Description: "Applicable only for destination types: azure, azure_legacy. Client Credentials OAuth Client Secret.",
				Computed:    true,
			},
			"qradar_username": schema.StringAttribute{
				Description: "Applicable only for destination type: qradar. Basic auth username provided by QRadar.",
				Computed:    true,
			},
			"qradar_password_masked": schema.StringAttribute{
				Description: "Applicable only for destination type: qradar. Basic auth password provided by QRadar.",
				Computed:    true,
			},
			"solar_winds_token_masked": schema.StringAttribute{
				Description: "Applicable only for destination type: solar_winds. Authentication token provided by Solar Winds.",
				Computed:    true,
			},
			"new_relic_api_key_masked": schema.StringAttribute{
				Description: "Applicable only for destination type: new_relic. API key provided by New Relic.",
				Computed:    true,
			},
			"datadog_api_key_masked": schema.StringAttribute{
				Description: "Applicable only for destination type: datadog. API key provided by Datadog.",
				Computed:    true,
			},
			"action_send_enabled": schema.BoolAttribute{
				Description: "Whether or not sending is enabled for action logs.",
				Computed:    true,
			},
			"action_entries_sent": schema.Int64Attribute{
				Description: "Number of log entries sent for the lifetime of this destination.",
				Computed:    true,
			},
			"sftp_action_send_enabled": schema.BoolAttribute{
				Description: "Whether or not sending is enabled for sftp_action logs.",
				Computed:    true,
			},
			"sftp_action_entries_sent": schema.Int64Attribute{
				Description: "Number of log entries sent for the lifetime of this destination.",
				Computed:    true,
			},
			"ftp_action_send_enabled": schema.BoolAttribute{
				Description: "Whether or not sending is enabled for ftp_action logs.",
				Computed:    true,
			},
			"ftp_action_entries_sent": schema.Int64Attribute{
				Description: "Number of log entries sent for the lifetime of this destination.",
				Computed:    true,
			},
			"web_dav_action_send_enabled": schema.BoolAttribute{
				Description: "Whether or not sending is enabled for web_dav_action logs.",
				Computed:    true,
			},
			"web_dav_action_entries_sent": schema.Int64Attribute{
				Description: "Number of log entries sent for the lifetime of this destination.",
				Computed:    true,
			},
			"sync_send_enabled": schema.BoolAttribute{
				Description: "Whether or not sending is enabled for sync logs.",
				Computed:    true,
			},
			"sync_entries_sent": schema.Int64Attribute{
				Description: "Number of log entries sent for the lifetime of this destination.",
				Computed:    true,
			},
			"outbound_connection_send_enabled": schema.BoolAttribute{
				Description: "Whether or not sending is enabled for outbound_connection logs.",
				Computed:    true,
			},
			"outbound_connection_entries_sent": schema.Int64Attribute{
				Description: "Number of log entries sent for the lifetime of this destination.",
				Computed:    true,
			},
			"automation_send_enabled": schema.BoolAttribute{
				Description: "Whether or not sending is enabled for automation logs.",
				Computed:    true,
			},
			"automation_entries_sent": schema.Int64Attribute{
				Description: "Number of log entries sent for the lifetime of this destination.",
				Computed:    true,
			},
			"api_request_send_enabled": schema.BoolAttribute{
				Description: "Whether or not sending is enabled for api_request logs.",
				Computed:    true,
			},
			"api_request_entries_sent": schema.Int64Attribute{
				Description: "Number of log entries sent for the lifetime of this destination.",
				Computed:    true,
			},
			"public_hosting_request_send_enabled": schema.BoolAttribute{
				Description: "Whether or not sending is enabled for public_hosting_request logs.",
				Computed:    true,
			},
			"public_hosting_request_entries_sent": schema.Int64Attribute{
				Description: "Number of log entries sent for the lifetime of this destination.",
				Computed:    true,
			},
			"email_send_enabled": schema.BoolAttribute{
				Description: "Whether or not sending is enabled for email logs.",
				Computed:    true,
			},
			"email_entries_sent": schema.Int64Attribute{
				Description: "Number of log entries sent for the lifetime of this destination.",
				Computed:    true,
			},
			"exavault_api_request_send_enabled": schema.BoolAttribute{
				Description: "Whether or not sending is enabled for exavault_api_request logs.",
				Computed:    true,
			},
			"exavault_api_request_entries_sent": schema.Int64Attribute{
				Description: "Number of log entries sent for the lifetime of this destination.",
				Computed:    true,
			},
			"settings_change_send_enabled": schema.BoolAttribute{
				Description: "Whether or not sending is enabled for settings_change logs.",
				Computed:    true,
			},
			"settings_change_entries_sent": schema.Int64Attribute{
				Description: "Number of log entries sent for the lifetime of this destination.",
				Computed:    true,
			},
			"last_http_call_target_type": schema.StringAttribute{
				Description: "Type of URL that was last called. Can be `destination_url` or `azure_oauth_client_credentials_url`",
				Computed:    true,
			},
			"last_http_call_success": schema.BoolAttribute{
				Description: "Was the last HTTP call made successful?",
				Computed:    true,
			},
			"last_http_call_response_code": schema.Int64Attribute{
				Description: "Last HTTP Call Response Code",
				Computed:    true,
			},
			"last_http_call_response_body": schema.StringAttribute{
				Description: "Last HTTP Call Response Body. Large responses are truncated.",
				Computed:    true,
			},
			"last_http_call_error_message": schema.StringAttribute{
				Description: "Last HTTP Call Error Message if applicable",
				Computed:    true,
			},
			"last_http_call_time": schema.StringAttribute{
				Description: "Time of Last HTTP Call",
				Computed:    true,
			},
			"last_http_call_duration_ms": schema.Int64Attribute{
				Description: "Duration of the last HTTP Call in milliseconds",
				Computed:    true,
			},
			"most_recent_http_call_success_time": schema.StringAttribute{
				Description: "Time of Most Recent Successful HTTP Call",
				Computed:    true,
			},
			"connection_test_entry": schema.StringAttribute{
				Description: "Connection Test Entry",
				Computed:    true,
			},
		},
	}
}

func (r *siemHttpDestinationDataSource) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
	var data siemHttpDestinationDataSourceModel
	diags := req.Config.Get(ctx, &data)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	paramsSiemHttpDestinationFind := files_sdk.SiemHttpDestinationFindParams{}
	paramsSiemHttpDestinationFind.Id = data.Id.ValueInt64()

	siemHttpDestination, err := r.client.Find(paramsSiemHttpDestinationFind, files_sdk.WithContext(ctx))
	if err != nil {
		resp.Diagnostics.AddError(
			"Error Reading Files SiemHttpDestination",
			"Could not read siem_http_destination id "+fmt.Sprint(data.Id.ValueInt64())+": "+err.Error(),
		)
		return
	}

	diags = r.populateDataSourceModel(ctx, siemHttpDestination, &data)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	diags = resp.State.Set(ctx, data)
	resp.Diagnostics.Append(diags...)
}

func (r *siemHttpDestinationDataSource) populateDataSourceModel(ctx context.Context, siemHttpDestination files_sdk.SiemHttpDestination, state *siemHttpDestinationDataSourceModel) (diags diag.Diagnostics) {
	var propDiags diag.Diagnostics

	state.Id = types.Int64Value(siemHttpDestination.Id)
	state.Name = types.StringValue(siemHttpDestination.Name)
	state.DestinationType = types.StringValue(siemHttpDestination.DestinationType)
	state.DestinationUrl = types.StringValue(siemHttpDestination.DestinationUrl)
	state.FileDestinationPath = types.StringValue(siemHttpDestination.FileDestinationPath)
	state.FileFormat = types.StringValue(siemHttpDestination.FileFormat)
	state.FileIntervalMinutes = types.Int64Value(siemHttpDestination.FileIntervalMinutes)
	state.AdditionalHeaders, propDiags = lib.ToDynamic(ctx, path.Root("additional_headers"), siemHttpDestination.AdditionalHeaders, state.AdditionalHeaders.UnderlyingValue())
	diags.Append(propDiags...)
	state.SendingActive = types.BoolPointerValue(siemHttpDestination.SendingActive)
	state.GenericPayloadType = types.StringValue(siemHttpDestination.GenericPayloadType)
	state.SplunkTokenMasked = types.StringValue(siemHttpDestination.SplunkTokenMasked)
	state.AzureDcrImmutableId = types.StringValue(siemHttpDestination.AzureDcrImmutableId)
	state.AzureStreamName = types.StringValue(siemHttpDestination.AzureStreamName)
	state.AzureOauthClientCredentialsTenantId = types.StringValue(siemHttpDestination.AzureOauthClientCredentialsTenantId)
	state.AzureOauthClientCredentialsClientId = types.StringValue(siemHttpDestination.AzureOauthClientCredentialsClientId)
	state.AzureOauthClientCredentialsClientSecretMasked = types.StringValue(siemHttpDestination.AzureOauthClientCredentialsClientSecretMasked)
	state.QradarUsername = types.StringValue(siemHttpDestination.QradarUsername)
	state.QradarPasswordMasked = types.StringValue(siemHttpDestination.QradarPasswordMasked)
	state.SolarWindsTokenMasked = types.StringValue(siemHttpDestination.SolarWindsTokenMasked)
	state.NewRelicApiKeyMasked = types.StringValue(siemHttpDestination.NewRelicApiKeyMasked)
	state.DatadogApiKeyMasked = types.StringValue(siemHttpDestination.DatadogApiKeyMasked)
	state.ActionSendEnabled = types.BoolPointerValue(siemHttpDestination.ActionSendEnabled)
	state.ActionEntriesSent = types.Int64Value(siemHttpDestination.ActionEntriesSent)
	state.SftpActionSendEnabled = types.BoolPointerValue(siemHttpDestination.SftpActionSendEnabled)
	state.SftpActionEntriesSent = types.Int64Value(siemHttpDestination.SftpActionEntriesSent)
	state.FtpActionSendEnabled = types.BoolPointerValue(siemHttpDestination.FtpActionSendEnabled)
	state.FtpActionEntriesSent = types.Int64Value(siemHttpDestination.FtpActionEntriesSent)
	state.WebDavActionSendEnabled = types.BoolPointerValue(siemHttpDestination.WebDavActionSendEnabled)
	state.WebDavActionEntriesSent = types.Int64Value(siemHttpDestination.WebDavActionEntriesSent)
	state.SyncSendEnabled = types.BoolPointerValue(siemHttpDestination.SyncSendEnabled)
	state.SyncEntriesSent = types.Int64Value(siemHttpDestination.SyncEntriesSent)
	state.OutboundConnectionSendEnabled = types.BoolPointerValue(siemHttpDestination.OutboundConnectionSendEnabled)
	state.OutboundConnectionEntriesSent = types.Int64Value(siemHttpDestination.OutboundConnectionEntriesSent)
	state.AutomationSendEnabled = types.BoolPointerValue(siemHttpDestination.AutomationSendEnabled)
	state.AutomationEntriesSent = types.Int64Value(siemHttpDestination.AutomationEntriesSent)
	state.ApiRequestSendEnabled = types.BoolPointerValue(siemHttpDestination.ApiRequestSendEnabled)
	state.ApiRequestEntriesSent = types.Int64Value(siemHttpDestination.ApiRequestEntriesSent)
	state.PublicHostingRequestSendEnabled = types.BoolPointerValue(siemHttpDestination.PublicHostingRequestSendEnabled)
	state.PublicHostingRequestEntriesSent = types.Int64Value(siemHttpDestination.PublicHostingRequestEntriesSent)
	state.EmailSendEnabled = types.BoolPointerValue(siemHttpDestination.EmailSendEnabled)
	state.EmailEntriesSent = types.Int64Value(siemHttpDestination.EmailEntriesSent)
	state.ExavaultApiRequestSendEnabled = types.BoolPointerValue(siemHttpDestination.ExavaultApiRequestSendEnabled)
	state.ExavaultApiRequestEntriesSent = types.Int64Value(siemHttpDestination.ExavaultApiRequestEntriesSent)
	state.SettingsChangeSendEnabled = types.BoolPointerValue(siemHttpDestination.SettingsChangeSendEnabled)
	state.SettingsChangeEntriesSent = types.Int64Value(siemHttpDestination.SettingsChangeEntriesSent)
	state.LastHttpCallTargetType = types.StringValue(siemHttpDestination.LastHttpCallTargetType)
	state.LastHttpCallSuccess = types.BoolPointerValue(siemHttpDestination.LastHttpCallSuccess)
	state.LastHttpCallResponseCode = types.Int64Value(siemHttpDestination.LastHttpCallResponseCode)
	state.LastHttpCallResponseBody = types.StringValue(siemHttpDestination.LastHttpCallResponseBody)
	state.LastHttpCallErrorMessage = types.StringValue(siemHttpDestination.LastHttpCallErrorMessage)
	state.LastHttpCallTime = types.StringValue(siemHttpDestination.LastHttpCallTime)
	state.LastHttpCallDurationMs = types.Int64Value(siemHttpDestination.LastHttpCallDurationMs)
	state.MostRecentHttpCallSuccessTime = types.StringValue(siemHttpDestination.MostRecentHttpCallSuccessTime)
	state.ConnectionTestEntry = types.StringValue(siemHttpDestination.ConnectionTestEntry)

	return
}
