package provider

import (
	"context"
	"fmt"

	files_sdk "github.com/Files-com/files-sdk-go/v3"
	as2_incoming_message "github.com/Files-com/files-sdk-go/v3/as2incomingmessage"
	"github.com/Files-com/terraform-provider-files/lib"

	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/datasource/schema"
	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/path"

	"github.com/hashicorp/terraform-plugin-framework/types"
)

var (
	_ datasource.DataSource              = &as2IncomingMessageDataSource{}
	_ datasource.DataSourceWithConfigure = &as2IncomingMessageDataSource{}
)

func NewAs2IncomingMessageDataSource() datasource.DataSource {
	return &as2IncomingMessageDataSource{}
}

type as2IncomingMessageDataSource struct {
	client *as2_incoming_message.Client
}

type as2IncomingMessageDataSourceModel struct {
	Id                          types.Int64   `tfsdk:"id"`
	ReplayedFromMessageId       types.Int64   `tfsdk:"replayed_from_message_id"`
	WorkspaceId                 types.Int64   `tfsdk:"workspace_id"`
	As2PartnerId                types.Int64   `tfsdk:"as2_partner_id"`
	As2StationId                types.Int64   `tfsdk:"as2_station_id"`
	Uuid                        types.String  `tfsdk:"uuid"`
	ContentType                 types.String  `tfsdk:"content_type"`
	HttpHeaders                 types.Dynamic `tfsdk:"http_headers"`
	ProcessingResult            types.String  `tfsdk:"processing_result"`
	ProcessingResultDescription types.String  `tfsdk:"processing_result_description"`
	Mic                         types.String  `tfsdk:"mic"`
	MicAlgo                     types.String  `tfsdk:"mic_algo"`
	As2To                       types.String  `tfsdk:"as2_to"`
	As2From                     types.String  `tfsdk:"as2_from"`
	MessageId                   types.String  `tfsdk:"message_id"`
	Subject                     types.String  `tfsdk:"subject"`
	Date                        types.String  `tfsdk:"date"`
	BodySize                    types.String  `tfsdk:"body_size"`
	AttachmentFilename          types.String  `tfsdk:"attachment_filename"`
	Ip                          types.String  `tfsdk:"ip"`
	CreatedAt                   types.String  `tfsdk:"created_at"`
	HttpResponseCode            types.String  `tfsdk:"http_response_code"`
	HttpResponseHeaders         types.Dynamic `tfsdk:"http_response_headers"`
	RecipientSerial             types.String  `tfsdk:"recipient_serial"`
	HexRecipientSerial          types.String  `tfsdk:"hex_recipient_serial"`
	RecipientIssuer             types.String  `tfsdk:"recipient_issuer"`
	MessageReceived             types.Bool    `tfsdk:"message_received"`
	MessageDecrypted            types.Bool    `tfsdk:"message_decrypted"`
	MessageSignatureVerified    types.Bool    `tfsdk:"message_signature_verified"`
	MessageProcessingSuccess    types.Bool    `tfsdk:"message_processing_success"`
	MessageMdnReturned          types.Bool    `tfsdk:"message_mdn_returned"`
	EncryptedUri                types.String  `tfsdk:"encrypted_uri"`
	SmimeSignedUri              types.String  `tfsdk:"smime_signed_uri"`
	SmimeUri                    types.String  `tfsdk:"smime_uri"`
	RawUri                      types.String  `tfsdk:"raw_uri"`
	MdnResponseUri              types.String  `tfsdk:"mdn_response_uri"`
}

func (r *as2IncomingMessageDataSource) Configure(_ context.Context, req datasource.ConfigureRequest, resp *datasource.ConfigureResponse) {
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

	r.client = &as2_incoming_message.Client{Config: sdk_config}
}

func (r *as2IncomingMessageDataSource) Metadata(_ context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_as2_incoming_message"
}

func (r *as2IncomingMessageDataSource) Schema(_ context.Context, _ datasource.SchemaRequest, resp *datasource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Description: "An AS2IncomingMessage is a single record created for each individual AS2 file transfer incoming from a Partner. It contains the message details and tracks the lifecycle status with milestones of the transfer.\n\nReplaying a retained message creates a new incoming message. Its `replayed_from_message_id` identifies the incoming-message record that was replayed. The original AS2 `message_id` is preserved. The source record ID remains available after that record is removed by message retention.",
		Attributes: map[string]schema.Attribute{
			"id": schema.Int64Attribute{
				Description: "ID of this AS2 incoming message record.",
				Required:    true,
			},
			"replayed_from_message_id": schema.Int64Attribute{
				Description: "ID of the incoming message record this message was replayed from. Empty for messages received directly from a trading partner.",
				Computed:    true,
			},
			"workspace_id": schema.Int64Attribute{
				Description: "ID of the Workspace associated with this AS2 Incoming Message.",
				Computed:    true,
			},
			"as2_partner_id": schema.Int64Attribute{
				Description: "Id of the AS2 Partner associated with this message.",
				Computed:    true,
			},
			"as2_station_id": schema.Int64Attribute{
				Description: "Id of the AS2 Station associated with this message.",
				Computed:    true,
			},
			"uuid": schema.StringAttribute{
				Description: "UUID assigned to this message.",
				Computed:    true,
			},
			"content_type": schema.StringAttribute{
				Description: "Content Type header of the incoming message.",
				Computed:    true,
			},
			"http_headers": schema.DynamicAttribute{
				Description: "HTTP Headers sent with this message.",
				Computed:    true,
			},
			"processing_result": schema.StringAttribute{
				Description: "Result of processing.",
				Computed:    true,
			},
			"processing_result_description": schema.StringAttribute{
				Description: "Result of processing description.",
				Computed:    true,
			},
			"mic": schema.StringAttribute{
				Description: "AS2 Message Integrity Check",
				Computed:    true,
			},
			"mic_algo": schema.StringAttribute{
				Description: "AS2 Message Integrity Check Algorithm Used",
				Computed:    true,
			},
			"as2_to": schema.StringAttribute{
				Description: "AS2 TO header of message",
				Computed:    true,
			},
			"as2_from": schema.StringAttribute{
				Description: "AS2 FROM header of message",
				Computed:    true,
			},
			"message_id": schema.StringAttribute{
				Description: "AS2 Message Id",
				Computed:    true,
			},
			"subject": schema.StringAttribute{
				Description: "AS2 Subject Header",
				Computed:    true,
			},
			"date": schema.StringAttribute{
				Description: "Date Header",
				Computed:    true,
			},
			"body_size": schema.StringAttribute{
				Description: "Encrypted Payload Body Size",
				Computed:    true,
			},
			"attachment_filename": schema.StringAttribute{
				Description: "Filename of the file being received.",
				Computed:    true,
			},
			"ip": schema.StringAttribute{
				Description: "IP Address of the Sender",
				Computed:    true,
			},
			"created_at": schema.StringAttribute{
				Description: "Message creation date/time",
				Computed:    true,
			},
			"http_response_code": schema.StringAttribute{
				Description: "HTTP Response Code sent for this message",
				Computed:    true,
			},
			"http_response_headers": schema.DynamicAttribute{
				Description: "HTTP Headers sent for this message.",
				Computed:    true,
			},
			"recipient_serial": schema.StringAttribute{
				Description: "Incoming Message Recipient(the Client Cert used to encrypt this message)'s serial",
				Computed:    true,
			},
			"hex_recipient_serial": schema.StringAttribute{
				Description: "Incoming Message Recipient(the Client Cert used to encrypt this message)'s serial in hex format.",
				Computed:    true,
			},
			"recipient_issuer": schema.StringAttribute{
				Description: "Incoming Message Recipient(the Client Cert used to encrypt this message)'s issuer",
				Computed:    true,
			},
			"message_received": schema.BoolAttribute{
				Description: "Message body received?",
				Computed:    true,
			},
			"message_decrypted": schema.BoolAttribute{
				Description: "Message decrypted successfully?",
				Computed:    true,
			},
			"message_signature_verified": schema.BoolAttribute{
				Description: "Message signature verified?",
				Computed:    true,
			},
			"message_processing_success": schema.BoolAttribute{
				Description: "Message processed successfully?",
				Computed:    true,
			},
			"message_mdn_returned": schema.BoolAttribute{
				Description: "MDN returned?",
				Computed:    true,
			},
			"encrypted_uri": schema.StringAttribute{
				Description: "URL to download the encrypted signed smime that is to sent as AS2 body",
				Computed:    true,
			},
			"smime_signed_uri": schema.StringAttribute{
				Description: "URL to download the file contents as smime with signature",
				Computed:    true,
			},
			"smime_uri": schema.StringAttribute{
				Description: "URL to download the file contents encoded as smime",
				Computed:    true,
			},
			"raw_uri": schema.StringAttribute{
				Description: "URL to download the original file contents",
				Computed:    true,
			},
			"mdn_response_uri": schema.StringAttribute{
				Description: "URL to download the http response body",
				Computed:    true,
			},
		},
	}
}

func (r *as2IncomingMessageDataSource) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
	var data as2IncomingMessageDataSourceModel
	diags := req.Config.Get(ctx, &data)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	paramsAs2IncomingMessageFind := files_sdk.As2IncomingMessageFindParams{}
	paramsAs2IncomingMessageFind.Id = data.Id.ValueInt64()

	as2IncomingMessage, err := r.client.Find(paramsAs2IncomingMessageFind, files_sdk.WithContext(ctx))
	if err != nil {
		resp.Diagnostics.AddError(
			"Error Reading Files As2IncomingMessage",
			"Could not read as2_incoming_message id "+fmt.Sprint(data.Id.ValueInt64())+": "+err.Error(),
		)
		return
	}

	diags = r.populateDataSourceModel(ctx, as2IncomingMessage, &data)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	diags = resp.State.Set(ctx, data)
	resp.Diagnostics.Append(diags...)
}

func (r *as2IncomingMessageDataSource) populateDataSourceModel(ctx context.Context, as2IncomingMessage files_sdk.As2IncomingMessage, state *as2IncomingMessageDataSourceModel) (diags diag.Diagnostics) {
	var propDiags diag.Diagnostics

	state.Id = types.Int64Value(as2IncomingMessage.Id)
	state.ReplayedFromMessageId = types.Int64Value(as2IncomingMessage.ReplayedFromMessageId)
	state.WorkspaceId = types.Int64Value(as2IncomingMessage.WorkspaceId)
	state.As2PartnerId = types.Int64Value(as2IncomingMessage.As2PartnerId)
	state.As2StationId = types.Int64Value(as2IncomingMessage.As2StationId)
	state.Uuid = types.StringValue(as2IncomingMessage.Uuid)
	state.ContentType = types.StringValue(as2IncomingMessage.ContentType)
	state.HttpHeaders, propDiags = lib.ToDynamic(ctx, path.Root("http_headers"), as2IncomingMessage.HttpHeaders, state.HttpHeaders.UnderlyingValue())
	diags.Append(propDiags...)
	state.ProcessingResult = types.StringValue(as2IncomingMessage.ProcessingResult)
	state.ProcessingResultDescription = types.StringValue(as2IncomingMessage.ProcessingResultDescription)
	state.Mic = types.StringValue(as2IncomingMessage.Mic)
	state.MicAlgo = types.StringValue(as2IncomingMessage.MicAlgo)
	state.As2To = types.StringValue(as2IncomingMessage.As2To)
	state.As2From = types.StringValue(as2IncomingMessage.As2From)
	state.MessageId = types.StringValue(as2IncomingMessage.MessageId)
	state.Subject = types.StringValue(as2IncomingMessage.Subject)
	state.Date = types.StringValue(as2IncomingMessage.Date)
	state.BodySize = types.StringValue(as2IncomingMessage.BodySize)
	state.AttachmentFilename = types.StringValue(as2IncomingMessage.AttachmentFilename)
	state.Ip = types.StringValue(as2IncomingMessage.Ip)
	if err := lib.TimeToStringType(ctx, path.Root("created_at"), as2IncomingMessage.CreatedAt, &state.CreatedAt); err != nil {
		diags.AddError(
			"Error Creating Files As2IncomingMessage",
			"Could not convert state created_at to string: "+err.Error(),
		)
	}
	state.HttpResponseCode = types.StringValue(as2IncomingMessage.HttpResponseCode)
	state.HttpResponseHeaders, propDiags = lib.ToDynamic(ctx, path.Root("http_response_headers"), as2IncomingMessage.HttpResponseHeaders, state.HttpResponseHeaders.UnderlyingValue())
	diags.Append(propDiags...)
	state.RecipientSerial = types.StringValue(as2IncomingMessage.RecipientSerial)
	state.HexRecipientSerial = types.StringValue(as2IncomingMessage.HexRecipientSerial)
	state.RecipientIssuer = types.StringValue(as2IncomingMessage.RecipientIssuer)
	state.MessageReceived = types.BoolPointerValue(as2IncomingMessage.MessageReceived)
	state.MessageDecrypted = types.BoolPointerValue(as2IncomingMessage.MessageDecrypted)
	state.MessageSignatureVerified = types.BoolPointerValue(as2IncomingMessage.MessageSignatureVerified)
	state.MessageProcessingSuccess = types.BoolPointerValue(as2IncomingMessage.MessageProcessingSuccess)
	state.MessageMdnReturned = types.BoolPointerValue(as2IncomingMessage.MessageMdnReturned)
	state.EncryptedUri = types.StringValue(as2IncomingMessage.EncryptedUri)
	state.SmimeSignedUri = types.StringValue(as2IncomingMessage.SmimeSignedUri)
	state.SmimeUri = types.StringValue(as2IncomingMessage.SmimeUri)
	state.RawUri = types.StringValue(as2IncomingMessage.RawUri)
	state.MdnResponseUri = types.StringValue(as2IncomingMessage.MdnResponseUri)

	return
}
