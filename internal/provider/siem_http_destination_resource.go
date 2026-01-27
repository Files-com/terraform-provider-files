package provider

import (
	"context"
	"fmt"
	"strconv"
	"strings"

	files_sdk "github.com/Files-com/files-sdk-go/v3"
	siem_http_destination "github.com/Files-com/files-sdk-go/v3/siemhttpdestination"
	"github.com/Files-com/terraform-provider-files/lib"
	"github.com/hashicorp/terraform-plugin-framework-validators/stringvalidator"
	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/boolplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/dynamicplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/int64planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/schema/validator"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

var (
	_ resource.Resource                = &siemHttpDestinationResource{}
	_ resource.ResourceWithConfigure   = &siemHttpDestinationResource{}
	_ resource.ResourceWithImportState = &siemHttpDestinationResource{}
)

func NewSiemHttpDestinationResource() resource.Resource {
	return &siemHttpDestinationResource{}
}

type siemHttpDestinationResource struct {
	client *siem_http_destination.Client
}

type siemHttpDestinationResourceModel struct {
	DestinationType                               types.String  `tfsdk:"destination_type"`
	Name                                          types.String  `tfsdk:"name"`
	DestinationUrl                                types.String  `tfsdk:"destination_url"`
	FileDestinationPath                           types.String  `tfsdk:"file_destination_path"`
	FileFormat                                    types.String  `tfsdk:"file_format"`
	FileIntervalMinutes                           types.Int64   `tfsdk:"file_interval_minutes"`
	AdditionalHeaders                             types.Dynamic `tfsdk:"additional_headers"`
	SendingActive                                 types.Bool    `tfsdk:"sending_active"`
	GenericPayloadType                            types.String  `tfsdk:"generic_payload_type"`
	AzureDcrImmutableId                           types.String  `tfsdk:"azure_dcr_immutable_id"`
	AzureStreamName                               types.String  `tfsdk:"azure_stream_name"`
	AzureOauthClientCredentialsTenantId           types.String  `tfsdk:"azure_oauth_client_credentials_tenant_id"`
	AzureOauthClientCredentialsClientId           types.String  `tfsdk:"azure_oauth_client_credentials_client_id"`
	QradarUsername                                types.String  `tfsdk:"qradar_username"`
	ActionSendEnabled                             types.Bool    `tfsdk:"action_send_enabled"`
	SftpActionSendEnabled                         types.Bool    `tfsdk:"sftp_action_send_enabled"`
	FtpActionSendEnabled                          types.Bool    `tfsdk:"ftp_action_send_enabled"`
	WebDavActionSendEnabled                       types.Bool    `tfsdk:"web_dav_action_send_enabled"`
	SyncSendEnabled                               types.Bool    `tfsdk:"sync_send_enabled"`
	OutboundConnectionSendEnabled                 types.Bool    `tfsdk:"outbound_connection_send_enabled"`
	AutomationSendEnabled                         types.Bool    `tfsdk:"automation_send_enabled"`
	ApiRequestSendEnabled                         types.Bool    `tfsdk:"api_request_send_enabled"`
	PublicHostingRequestSendEnabled               types.Bool    `tfsdk:"public_hosting_request_send_enabled"`
	EmailSendEnabled                              types.Bool    `tfsdk:"email_send_enabled"`
	ExavaultApiRequestSendEnabled                 types.Bool    `tfsdk:"exavault_api_request_send_enabled"`
	SettingsChangeSendEnabled                     types.Bool    `tfsdk:"settings_change_send_enabled"`
	SplunkToken                                   types.String  `tfsdk:"splunk_token"`
	AzureOauthClientCredentialsClientSecret       types.String  `tfsdk:"azure_oauth_client_credentials_client_secret"`
	QradarPassword                                types.String  `tfsdk:"qradar_password"`
	SolarWindsToken                               types.String  `tfsdk:"solar_winds_token"`
	NewRelicApiKey                                types.String  `tfsdk:"new_relic_api_key"`
	DatadogApiKey                                 types.String  `tfsdk:"datadog_api_key"`
	Id                                            types.Int64   `tfsdk:"id"`
	SplunkTokenMasked                             types.String  `tfsdk:"splunk_token_masked"`
	AzureOauthClientCredentialsClientSecretMasked types.String  `tfsdk:"azure_oauth_client_credentials_client_secret_masked"`
	QradarPasswordMasked                          types.String  `tfsdk:"qradar_password_masked"`
	SolarWindsTokenMasked                         types.String  `tfsdk:"solar_winds_token_masked"`
	NewRelicApiKeyMasked                          types.String  `tfsdk:"new_relic_api_key_masked"`
	DatadogApiKeyMasked                           types.String  `tfsdk:"datadog_api_key_masked"`
	ActionEntriesSent                             types.Int64   `tfsdk:"action_entries_sent"`
	SftpActionEntriesSent                         types.Int64   `tfsdk:"sftp_action_entries_sent"`
	FtpActionEntriesSent                          types.Int64   `tfsdk:"ftp_action_entries_sent"`
	WebDavActionEntriesSent                       types.Int64   `tfsdk:"web_dav_action_entries_sent"`
	SyncEntriesSent                               types.Int64   `tfsdk:"sync_entries_sent"`
	OutboundConnectionEntriesSent                 types.Int64   `tfsdk:"outbound_connection_entries_sent"`
	AutomationEntriesSent                         types.Int64   `tfsdk:"automation_entries_sent"`
	ApiRequestEntriesSent                         types.Int64   `tfsdk:"api_request_entries_sent"`
	PublicHostingRequestEntriesSent               types.Int64   `tfsdk:"public_hosting_request_entries_sent"`
	EmailEntriesSent                              types.Int64   `tfsdk:"email_entries_sent"`
	ExavaultApiRequestEntriesSent                 types.Int64   `tfsdk:"exavault_api_request_entries_sent"`
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

func (r *siemHttpDestinationResource) Configure(_ context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {
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

func (r *siemHttpDestinationResource) Metadata(_ context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_siem_http_destination"
}

func (r *siemHttpDestinationResource) Schema(_ context.Context, _ resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Description: "",
		Attributes: map[string]schema.Attribute{
			"destination_type": schema.StringAttribute{
				Description: "Destination Type",
				Required:    true,
				Validators: []validator.String{
					stringvalidator.OneOf("generic", "splunk", "azure_legacy", "qradar", "sumo", "rapid7", "solar_winds", "new_relic", "datadog", "azure", "file"),
				},
			},
			"name": schema.StringAttribute{
				Description: "Name for this Destination",
				Computed:    true,
				Optional:    true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
			},
			"destination_url": schema.StringAttribute{
				Description: "Destination Url",
				Computed:    true,
				Optional:    true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
			},
			"file_destination_path": schema.StringAttribute{
				Description: "Applicable only for destination type: file. Destination folder path on Files.com.",
				Computed:    true,
				Optional:    true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
			},
			"file_format": schema.StringAttribute{
				Description: "Applicable only for destination type: file. Generated file format.",
				Computed:    true,
				Optional:    true,
				Validators: []validator.String{
					stringvalidator.OneOf("json", "csv"),
				},
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
			},
			"file_interval_minutes": schema.Int64Attribute{
				Description: "Applicable only for destination type: file. Interval, in minutes, between file deliveries.",
				Computed:    true,
				Optional:    true,
				PlanModifiers: []planmodifier.Int64{
					int64planmodifier.UseStateForUnknown(),
				},
			},
			"additional_headers": schema.DynamicAttribute{
				Description: "Additional HTTP Headers included in calls to the destination URL",
				Computed:    true,
				Optional:    true,
				PlanModifiers: []planmodifier.Dynamic{
					dynamicplanmodifier.UseStateForUnknown(),
				},
			},
			"sending_active": schema.BoolAttribute{
				Description: "Whether this SIEM HTTP Destination is currently being sent to or not",
				Computed:    true,
				Optional:    true,
				PlanModifiers: []planmodifier.Bool{
					boolplanmodifier.UseStateForUnknown(),
				},
			},
			"generic_payload_type": schema.StringAttribute{
				Description: "Applicable only for destination type: generic. Indicates the type of HTTP body. Can be json_newline or json_array. json_newline is multiple log entries as JSON separated by newlines. json_array is a single JSON array containing multiple log entries as JSON.",
				Computed:    true,
				Optional:    true,
				Validators: []validator.String{
					stringvalidator.OneOf("json_newline", "json_array"),
				},
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
			},
			"azure_dcr_immutable_id": schema.StringAttribute{
				Description: "Applicable only for destination types: azure, azure_legacy. Immutable ID of the Data Collection Rule.",
				Computed:    true,
				Optional:    true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
			},
			"azure_stream_name": schema.StringAttribute{
				Description: "Applicable only for destination type: azure. Name of the stream in the DCR that represents the destination table.",
				Computed:    true,
				Optional:    true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
			},
			"azure_oauth_client_credentials_tenant_id": schema.StringAttribute{
				Description: "Applicable only for destination types: azure, azure_legacy. Client Credentials OAuth Tenant ID.",
				Computed:    true,
				Optional:    true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
			},
			"azure_oauth_client_credentials_client_id": schema.StringAttribute{
				Description: "Applicable only for destination types: azure, azure_legacy. Client Credentials OAuth Client ID.",
				Computed:    true,
				Optional:    true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
			},
			"qradar_username": schema.StringAttribute{
				Description: "Applicable only for destination type: qradar. Basic auth username provided by QRadar.",
				Computed:    true,
				Optional:    true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
			},
			"action_send_enabled": schema.BoolAttribute{
				Description: "Whether or not sending is enabled for action logs.",
				Computed:    true,
				Optional:    true,
				PlanModifiers: []planmodifier.Bool{
					boolplanmodifier.UseStateForUnknown(),
				},
			},
			"sftp_action_send_enabled": schema.BoolAttribute{
				Description: "Whether or not sending is enabled for sftp_action logs.",
				Computed:    true,
				Optional:    true,
				PlanModifiers: []planmodifier.Bool{
					boolplanmodifier.UseStateForUnknown(),
				},
			},
			"ftp_action_send_enabled": schema.BoolAttribute{
				Description: "Whether or not sending is enabled for ftp_action logs.",
				Computed:    true,
				Optional:    true,
				PlanModifiers: []planmodifier.Bool{
					boolplanmodifier.UseStateForUnknown(),
				},
			},
			"web_dav_action_send_enabled": schema.BoolAttribute{
				Description: "Whether or not sending is enabled for web_dav_action logs.",
				Computed:    true,
				Optional:    true,
				PlanModifiers: []planmodifier.Bool{
					boolplanmodifier.UseStateForUnknown(),
				},
			},
			"sync_send_enabled": schema.BoolAttribute{
				Description: "Whether or not sending is enabled for sync logs.",
				Computed:    true,
				Optional:    true,
				PlanModifiers: []planmodifier.Bool{
					boolplanmodifier.UseStateForUnknown(),
				},
			},
			"outbound_connection_send_enabled": schema.BoolAttribute{
				Description: "Whether or not sending is enabled for outbound_connection logs.",
				Computed:    true,
				Optional:    true,
				PlanModifiers: []planmodifier.Bool{
					boolplanmodifier.UseStateForUnknown(),
				},
			},
			"automation_send_enabled": schema.BoolAttribute{
				Description: "Whether or not sending is enabled for automation logs.",
				Computed:    true,
				Optional:    true,
				PlanModifiers: []planmodifier.Bool{
					boolplanmodifier.UseStateForUnknown(),
				},
			},
			"api_request_send_enabled": schema.BoolAttribute{
				Description: "Whether or not sending is enabled for api_request logs.",
				Computed:    true,
				Optional:    true,
				PlanModifiers: []planmodifier.Bool{
					boolplanmodifier.UseStateForUnknown(),
				},
			},
			"public_hosting_request_send_enabled": schema.BoolAttribute{
				Description: "Whether or not sending is enabled for public_hosting_request logs.",
				Computed:    true,
				Optional:    true,
				PlanModifiers: []planmodifier.Bool{
					boolplanmodifier.UseStateForUnknown(),
				},
			},
			"email_send_enabled": schema.BoolAttribute{
				Description: "Whether or not sending is enabled for email logs.",
				Computed:    true,
				Optional:    true,
				PlanModifiers: []planmodifier.Bool{
					boolplanmodifier.UseStateForUnknown(),
				},
			},
			"exavault_api_request_send_enabled": schema.BoolAttribute{
				Description: "Whether or not sending is enabled for exavault_api_request logs.",
				Computed:    true,
				Optional:    true,
				PlanModifiers: []planmodifier.Bool{
					boolplanmodifier.UseStateForUnknown(),
				},
			},
			"settings_change_send_enabled": schema.BoolAttribute{
				Description: "Whether or not sending is enabled for settings_change logs.",
				Computed:    true,
				Optional:    true,
				PlanModifiers: []planmodifier.Bool{
					boolplanmodifier.UseStateForUnknown(),
				},
			},
			"splunk_token": schema.StringAttribute{
				Description: "Applicable only for destination type: splunk. Authentication token provided by Splunk.",
				Optional:    true,
				WriteOnly:   true,
			},
			"azure_oauth_client_credentials_client_secret": schema.StringAttribute{
				Description: "Applicable only for destination type: azure. Client Credentials OAuth Client Secret.",
				Optional:    true,
				WriteOnly:   true,
			},
			"qradar_password": schema.StringAttribute{
				Description: "Applicable only for destination type: qradar. Basic auth password provided by QRadar.",
				Optional:    true,
				WriteOnly:   true,
			},
			"solar_winds_token": schema.StringAttribute{
				Description: "Applicable only for destination type: solar_winds. Authentication token provided by Solar Winds.",
				Optional:    true,
				WriteOnly:   true,
			},
			"new_relic_api_key": schema.StringAttribute{
				Description: "Applicable only for destination type: new_relic. API key provided by New Relic.",
				Optional:    true,
				WriteOnly:   true,
			},
			"datadog_api_key": schema.StringAttribute{
				Description: "Applicable only for destination type: datadog. API key provided by Datadog.",
				Optional:    true,
				WriteOnly:   true,
			},
			"id": schema.Int64Attribute{
				Description: "SIEM HTTP Destination ID",
				Computed:    true,
				PlanModifiers: []planmodifier.Int64{
					int64planmodifier.UseStateForUnknown(),
				},
			},
			"splunk_token_masked": schema.StringAttribute{
				Description: "Applicable only for destination type: splunk. Authentication token provided by Splunk.",
				Computed:    true,
			},
			"azure_oauth_client_credentials_client_secret_masked": schema.StringAttribute{
				Description: "Applicable only for destination types: azure, azure_legacy. Client Credentials OAuth Client Secret.",
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
			"action_entries_sent": schema.Int64Attribute{
				Description: "Number of log entries sent for the lifetime of this destination.",
				Computed:    true,
			},
			"sftp_action_entries_sent": schema.Int64Attribute{
				Description: "Number of log entries sent for the lifetime of this destination.",
				Computed:    true,
			},
			"ftp_action_entries_sent": schema.Int64Attribute{
				Description: "Number of log entries sent for the lifetime of this destination.",
				Computed:    true,
			},
			"web_dav_action_entries_sent": schema.Int64Attribute{
				Description: "Number of log entries sent for the lifetime of this destination.",
				Computed:    true,
			},
			"sync_entries_sent": schema.Int64Attribute{
				Description: "Number of log entries sent for the lifetime of this destination.",
				Computed:    true,
			},
			"outbound_connection_entries_sent": schema.Int64Attribute{
				Description: "Number of log entries sent for the lifetime of this destination.",
				Computed:    true,
			},
			"automation_entries_sent": schema.Int64Attribute{
				Description: "Number of log entries sent for the lifetime of this destination.",
				Computed:    true,
			},
			"api_request_entries_sent": schema.Int64Attribute{
				Description: "Number of log entries sent for the lifetime of this destination.",
				Computed:    true,
			},
			"public_hosting_request_entries_sent": schema.Int64Attribute{
				Description: "Number of log entries sent for the lifetime of this destination.",
				Computed:    true,
			},
			"email_entries_sent": schema.Int64Attribute{
				Description: "Number of log entries sent for the lifetime of this destination.",
				Computed:    true,
			},
			"exavault_api_request_entries_sent": schema.Int64Attribute{
				Description: "Number of log entries sent for the lifetime of this destination.",
				Computed:    true,
			},
			"settings_change_entries_sent": schema.Int64Attribute{
				Description: "Number of log entries sent for the lifetime of this destination.",
				Computed:    true,
			},
			"last_http_call_target_type": schema.StringAttribute{
				Description: "Type of URL that was last called. Can be `destination_url` or `azure_oauth_client_credentials_url`",
				Computed:    true,
				Validators: []validator.String{
					stringvalidator.OneOf("destination_url", "azure_oauth_client_credentials_url", "file_destination"),
				},
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

func (r *siemHttpDestinationResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	var plan siemHttpDestinationResourceModel
	diags := req.Plan.Get(ctx, &plan)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}
	var config siemHttpDestinationResourceModel
	diags = req.Config.Get(ctx, &config)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	paramsSiemHttpDestinationCreate := files_sdk.SiemHttpDestinationCreateParams{}
	paramsSiemHttpDestinationCreate.Name = plan.Name.ValueString()
	createAdditionalHeaders, diags := lib.DynamicToInterface(ctx, path.Root("additional_headers"), plan.AdditionalHeaders)
	resp.Diagnostics.Append(diags...)
	paramsSiemHttpDestinationCreate.AdditionalHeaders = createAdditionalHeaders
	if !plan.SendingActive.IsNull() && !plan.SendingActive.IsUnknown() {
		paramsSiemHttpDestinationCreate.SendingActive = plan.SendingActive.ValueBoolPointer()
	}
	paramsSiemHttpDestinationCreate.GenericPayloadType = paramsSiemHttpDestinationCreate.GenericPayloadType.Enum()[plan.GenericPayloadType.ValueString()]
	paramsSiemHttpDestinationCreate.FileDestinationPath = plan.FileDestinationPath.ValueString()
	paramsSiemHttpDestinationCreate.FileFormat = paramsSiemHttpDestinationCreate.FileFormat.Enum()[plan.FileFormat.ValueString()]
	paramsSiemHttpDestinationCreate.FileIntervalMinutes = plan.FileIntervalMinutes.ValueInt64()
	paramsSiemHttpDestinationCreate.SplunkToken = config.SplunkToken.ValueString()
	paramsSiemHttpDestinationCreate.AzureDcrImmutableId = plan.AzureDcrImmutableId.ValueString()
	paramsSiemHttpDestinationCreate.AzureStreamName = plan.AzureStreamName.ValueString()
	paramsSiemHttpDestinationCreate.AzureOauthClientCredentialsTenantId = plan.AzureOauthClientCredentialsTenantId.ValueString()
	paramsSiemHttpDestinationCreate.AzureOauthClientCredentialsClientId = plan.AzureOauthClientCredentialsClientId.ValueString()
	paramsSiemHttpDestinationCreate.AzureOauthClientCredentialsClientSecret = config.AzureOauthClientCredentialsClientSecret.ValueString()
	paramsSiemHttpDestinationCreate.QradarUsername = plan.QradarUsername.ValueString()
	paramsSiemHttpDestinationCreate.QradarPassword = config.QradarPassword.ValueString()
	paramsSiemHttpDestinationCreate.SolarWindsToken = config.SolarWindsToken.ValueString()
	paramsSiemHttpDestinationCreate.NewRelicApiKey = config.NewRelicApiKey.ValueString()
	paramsSiemHttpDestinationCreate.DatadogApiKey = config.DatadogApiKey.ValueString()
	if !plan.ActionSendEnabled.IsNull() && !plan.ActionSendEnabled.IsUnknown() {
		paramsSiemHttpDestinationCreate.ActionSendEnabled = plan.ActionSendEnabled.ValueBoolPointer()
	}
	if !plan.SftpActionSendEnabled.IsNull() && !plan.SftpActionSendEnabled.IsUnknown() {
		paramsSiemHttpDestinationCreate.SftpActionSendEnabled = plan.SftpActionSendEnabled.ValueBoolPointer()
	}
	if !plan.FtpActionSendEnabled.IsNull() && !plan.FtpActionSendEnabled.IsUnknown() {
		paramsSiemHttpDestinationCreate.FtpActionSendEnabled = plan.FtpActionSendEnabled.ValueBoolPointer()
	}
	if !plan.WebDavActionSendEnabled.IsNull() && !plan.WebDavActionSendEnabled.IsUnknown() {
		paramsSiemHttpDestinationCreate.WebDavActionSendEnabled = plan.WebDavActionSendEnabled.ValueBoolPointer()
	}
	if !plan.SyncSendEnabled.IsNull() && !plan.SyncSendEnabled.IsUnknown() {
		paramsSiemHttpDestinationCreate.SyncSendEnabled = plan.SyncSendEnabled.ValueBoolPointer()
	}
	if !plan.OutboundConnectionSendEnabled.IsNull() && !plan.OutboundConnectionSendEnabled.IsUnknown() {
		paramsSiemHttpDestinationCreate.OutboundConnectionSendEnabled = plan.OutboundConnectionSendEnabled.ValueBoolPointer()
	}
	if !plan.AutomationSendEnabled.IsNull() && !plan.AutomationSendEnabled.IsUnknown() {
		paramsSiemHttpDestinationCreate.AutomationSendEnabled = plan.AutomationSendEnabled.ValueBoolPointer()
	}
	if !plan.ApiRequestSendEnabled.IsNull() && !plan.ApiRequestSendEnabled.IsUnknown() {
		paramsSiemHttpDestinationCreate.ApiRequestSendEnabled = plan.ApiRequestSendEnabled.ValueBoolPointer()
	}
	if !plan.PublicHostingRequestSendEnabled.IsNull() && !plan.PublicHostingRequestSendEnabled.IsUnknown() {
		paramsSiemHttpDestinationCreate.PublicHostingRequestSendEnabled = plan.PublicHostingRequestSendEnabled.ValueBoolPointer()
	}
	if !plan.EmailSendEnabled.IsNull() && !plan.EmailSendEnabled.IsUnknown() {
		paramsSiemHttpDestinationCreate.EmailSendEnabled = plan.EmailSendEnabled.ValueBoolPointer()
	}
	if !plan.ExavaultApiRequestSendEnabled.IsNull() && !plan.ExavaultApiRequestSendEnabled.IsUnknown() {
		paramsSiemHttpDestinationCreate.ExavaultApiRequestSendEnabled = plan.ExavaultApiRequestSendEnabled.ValueBoolPointer()
	}
	if !plan.SettingsChangeSendEnabled.IsNull() && !plan.SettingsChangeSendEnabled.IsUnknown() {
		paramsSiemHttpDestinationCreate.SettingsChangeSendEnabled = plan.SettingsChangeSendEnabled.ValueBoolPointer()
	}
	paramsSiemHttpDestinationCreate.DestinationType = paramsSiemHttpDestinationCreate.DestinationType.Enum()[plan.DestinationType.ValueString()]
	paramsSiemHttpDestinationCreate.DestinationUrl = plan.DestinationUrl.ValueString()

	if resp.Diagnostics.HasError() {
		return
	}

	siemHttpDestination, err := r.client.Create(paramsSiemHttpDestinationCreate, files_sdk.WithContext(ctx))
	if err != nil {
		resp.Diagnostics.AddError(
			"Error Creating Files SiemHttpDestination",
			"Could not create siem_http_destination, unexpected error: "+err.Error(),
		)
		return
	}

	diags = r.populateResourceModel(ctx, siemHttpDestination, &plan)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	diags = resp.State.Set(ctx, plan)
	resp.Diagnostics.Append(diags...)
}

func (r *siemHttpDestinationResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	var state siemHttpDestinationResourceModel
	diags := req.State.Get(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	paramsSiemHttpDestinationFind := files_sdk.SiemHttpDestinationFindParams{}
	paramsSiemHttpDestinationFind.Id = state.Id.ValueInt64()

	siemHttpDestination, err := r.client.Find(paramsSiemHttpDestinationFind, files_sdk.WithContext(ctx))
	if err != nil {
		if files_sdk.IsNotExist(err) {
			resp.State.RemoveResource(ctx)
			return
		}

		resp.Diagnostics.AddError(
			"Error Reading Files SiemHttpDestination",
			"Could not read siem_http_destination id "+fmt.Sprint(state.Id.ValueInt64())+": "+err.Error(),
		)
		return
	}

	diags = r.populateResourceModel(ctx, siemHttpDestination, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	diags = resp.State.Set(ctx, state)
	resp.Diagnostics.Append(diags...)
}

func (r *siemHttpDestinationResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	var plan siemHttpDestinationResourceModel
	diags := req.Plan.Get(ctx, &plan)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}
	var config siemHttpDestinationResourceModel
	diags = req.Config.Get(ctx, &config)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	paramsSiemHttpDestinationUpdate := map[string]interface{}{}
	if !plan.Id.IsNull() && !plan.Id.IsUnknown() {
		paramsSiemHttpDestinationUpdate["id"] = plan.Id.ValueInt64()
	}
	if !config.Name.IsNull() && !config.Name.IsUnknown() {
		paramsSiemHttpDestinationUpdate["name"] = config.Name.ValueString()
	}
	updateAdditionalHeaders, diags := lib.DynamicToInterface(ctx, path.Root("additional_headers"), config.AdditionalHeaders)
	resp.Diagnostics.Append(diags...)
	paramsSiemHttpDestinationUpdate["additional_headers"] = updateAdditionalHeaders
	if !config.SendingActive.IsNull() && !config.SendingActive.IsUnknown() {
		paramsSiemHttpDestinationUpdate["sending_active"] = config.SendingActive.ValueBool()
	}
	if !config.GenericPayloadType.IsNull() && !config.GenericPayloadType.IsUnknown() {
		paramsSiemHttpDestinationUpdate["generic_payload_type"] = config.GenericPayloadType.ValueString()
	}
	if !config.FileDestinationPath.IsNull() && !config.FileDestinationPath.IsUnknown() {
		paramsSiemHttpDestinationUpdate["file_destination_path"] = config.FileDestinationPath.ValueString()
	}
	if !config.FileFormat.IsNull() && !config.FileFormat.IsUnknown() {
		paramsSiemHttpDestinationUpdate["file_format"] = config.FileFormat.ValueString()
	}
	if !config.FileIntervalMinutes.IsNull() && !config.FileIntervalMinutes.IsUnknown() {
		paramsSiemHttpDestinationUpdate["file_interval_minutes"] = config.FileIntervalMinutes.ValueInt64()
	}
	if !config.SplunkToken.IsNull() && !config.SplunkToken.IsUnknown() {
		paramsSiemHttpDestinationUpdate["splunk_token"] = config.SplunkToken.ValueString()
	}
	if !config.AzureDcrImmutableId.IsNull() && !config.AzureDcrImmutableId.IsUnknown() {
		paramsSiemHttpDestinationUpdate["azure_dcr_immutable_id"] = config.AzureDcrImmutableId.ValueString()
	}
	if !config.AzureStreamName.IsNull() && !config.AzureStreamName.IsUnknown() {
		paramsSiemHttpDestinationUpdate["azure_stream_name"] = config.AzureStreamName.ValueString()
	}
	if !config.AzureOauthClientCredentialsTenantId.IsNull() && !config.AzureOauthClientCredentialsTenantId.IsUnknown() {
		paramsSiemHttpDestinationUpdate["azure_oauth_client_credentials_tenant_id"] = config.AzureOauthClientCredentialsTenantId.ValueString()
	}
	if !config.AzureOauthClientCredentialsClientId.IsNull() && !config.AzureOauthClientCredentialsClientId.IsUnknown() {
		paramsSiemHttpDestinationUpdate["azure_oauth_client_credentials_client_id"] = config.AzureOauthClientCredentialsClientId.ValueString()
	}
	if !config.AzureOauthClientCredentialsClientSecret.IsNull() && !config.AzureOauthClientCredentialsClientSecret.IsUnknown() {
		paramsSiemHttpDestinationUpdate["azure_oauth_client_credentials_client_secret"] = config.AzureOauthClientCredentialsClientSecret.ValueString()
	}
	if !config.QradarUsername.IsNull() && !config.QradarUsername.IsUnknown() {
		paramsSiemHttpDestinationUpdate["qradar_username"] = config.QradarUsername.ValueString()
	}
	if !config.QradarPassword.IsNull() && !config.QradarPassword.IsUnknown() {
		paramsSiemHttpDestinationUpdate["qradar_password"] = config.QradarPassword.ValueString()
	}
	if !config.SolarWindsToken.IsNull() && !config.SolarWindsToken.IsUnknown() {
		paramsSiemHttpDestinationUpdate["solar_winds_token"] = config.SolarWindsToken.ValueString()
	}
	if !config.NewRelicApiKey.IsNull() && !config.NewRelicApiKey.IsUnknown() {
		paramsSiemHttpDestinationUpdate["new_relic_api_key"] = config.NewRelicApiKey.ValueString()
	}
	if !config.DatadogApiKey.IsNull() && !config.DatadogApiKey.IsUnknown() {
		paramsSiemHttpDestinationUpdate["datadog_api_key"] = config.DatadogApiKey.ValueString()
	}
	if !config.ActionSendEnabled.IsNull() && !config.ActionSendEnabled.IsUnknown() {
		paramsSiemHttpDestinationUpdate["action_send_enabled"] = config.ActionSendEnabled.ValueBool()
	}
	if !config.SftpActionSendEnabled.IsNull() && !config.SftpActionSendEnabled.IsUnknown() {
		paramsSiemHttpDestinationUpdate["sftp_action_send_enabled"] = config.SftpActionSendEnabled.ValueBool()
	}
	if !config.FtpActionSendEnabled.IsNull() && !config.FtpActionSendEnabled.IsUnknown() {
		paramsSiemHttpDestinationUpdate["ftp_action_send_enabled"] = config.FtpActionSendEnabled.ValueBool()
	}
	if !config.WebDavActionSendEnabled.IsNull() && !config.WebDavActionSendEnabled.IsUnknown() {
		paramsSiemHttpDestinationUpdate["web_dav_action_send_enabled"] = config.WebDavActionSendEnabled.ValueBool()
	}
	if !config.SyncSendEnabled.IsNull() && !config.SyncSendEnabled.IsUnknown() {
		paramsSiemHttpDestinationUpdate["sync_send_enabled"] = config.SyncSendEnabled.ValueBool()
	}
	if !config.OutboundConnectionSendEnabled.IsNull() && !config.OutboundConnectionSendEnabled.IsUnknown() {
		paramsSiemHttpDestinationUpdate["outbound_connection_send_enabled"] = config.OutboundConnectionSendEnabled.ValueBool()
	}
	if !config.AutomationSendEnabled.IsNull() && !config.AutomationSendEnabled.IsUnknown() {
		paramsSiemHttpDestinationUpdate["automation_send_enabled"] = config.AutomationSendEnabled.ValueBool()
	}
	if !config.ApiRequestSendEnabled.IsNull() && !config.ApiRequestSendEnabled.IsUnknown() {
		paramsSiemHttpDestinationUpdate["api_request_send_enabled"] = config.ApiRequestSendEnabled.ValueBool()
	}
	if !config.PublicHostingRequestSendEnabled.IsNull() && !config.PublicHostingRequestSendEnabled.IsUnknown() {
		paramsSiemHttpDestinationUpdate["public_hosting_request_send_enabled"] = config.PublicHostingRequestSendEnabled.ValueBool()
	}
	if !config.EmailSendEnabled.IsNull() && !config.EmailSendEnabled.IsUnknown() {
		paramsSiemHttpDestinationUpdate["email_send_enabled"] = config.EmailSendEnabled.ValueBool()
	}
	if !config.ExavaultApiRequestSendEnabled.IsNull() && !config.ExavaultApiRequestSendEnabled.IsUnknown() {
		paramsSiemHttpDestinationUpdate["exavault_api_request_send_enabled"] = config.ExavaultApiRequestSendEnabled.ValueBool()
	}
	if !config.SettingsChangeSendEnabled.IsNull() && !config.SettingsChangeSendEnabled.IsUnknown() {
		paramsSiemHttpDestinationUpdate["settings_change_send_enabled"] = config.SettingsChangeSendEnabled.ValueBool()
	}
	if !config.DestinationType.IsNull() && !config.DestinationType.IsUnknown() {
		paramsSiemHttpDestinationUpdate["destination_type"] = config.DestinationType.ValueString()
	}
	if !config.DestinationUrl.IsNull() && !config.DestinationUrl.IsUnknown() {
		paramsSiemHttpDestinationUpdate["destination_url"] = config.DestinationUrl.ValueString()
	}

	if resp.Diagnostics.HasError() {
		return
	}

	siemHttpDestination, err := r.client.UpdateWithMap(paramsSiemHttpDestinationUpdate, files_sdk.WithContext(ctx))
	if err != nil {
		resp.Diagnostics.AddError(
			"Error Updating Files SiemHttpDestination",
			"Could not update siem_http_destination, unexpected error: "+err.Error(),
		)
		return
	}

	diags = r.populateResourceModel(ctx, siemHttpDestination, &plan)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	diags = resp.State.Set(ctx, plan)
	resp.Diagnostics.Append(diags...)
}

func (r *siemHttpDestinationResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	var state siemHttpDestinationResourceModel
	diags := req.State.Get(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	paramsSiemHttpDestinationDelete := files_sdk.SiemHttpDestinationDeleteParams{}
	paramsSiemHttpDestinationDelete.Id = state.Id.ValueInt64()

	err := r.client.Delete(paramsSiemHttpDestinationDelete, files_sdk.WithContext(ctx))
	if err != nil && !files_sdk.IsNotExist(err) {
		resp.Diagnostics.AddError(
			"Error Deleting Files SiemHttpDestination",
			"Could not delete siem_http_destination id "+fmt.Sprint(state.Id.ValueInt64())+": "+err.Error(),
		)
	}
}

func (r *siemHttpDestinationResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	idParts := strings.SplitN(req.ID, ",", 1)

	if len(idParts) != 1 || idParts[0] == "" {
		resp.Diagnostics.AddError(
			"Unexpected Import Identifier",
			fmt.Sprintf("Expected import identifier with format: id. Got: %q", req.ID),
		)
		return
	}

	idPart, err := strconv.ParseFloat(idParts[0], 64)
	if err != nil {
		resp.Diagnostics.AddError(
			"Error Parsing ID",
			"Could not parse id: "+err.Error(),
		)
		return
	}
	resp.Diagnostics.Append(resp.State.SetAttribute(ctx, path.Root("id"), idPart)...)

}

func (r *siemHttpDestinationResource) populateResourceModel(ctx context.Context, siemHttpDestination files_sdk.SiemHttpDestination, state *siemHttpDestinationResourceModel) (diags diag.Diagnostics) {
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
