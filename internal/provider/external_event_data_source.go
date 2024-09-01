package provider

import (
	"context"
	"fmt"

	files_sdk "github.com/Files-com/files-sdk-go/v3"
	external_event "github.com/Files-com/files-sdk-go/v3/externalevent"
	"github.com/Files-com/terraform-provider-files/lib"
	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/datasource/schema"
	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

var (
	_ datasource.DataSource              = &externalEventDataSource{}
	_ datasource.DataSourceWithConfigure = &externalEventDataSource{}
)

func NewExternalEventDataSource() datasource.DataSource {
	return &externalEventDataSource{}
}

type externalEventDataSource struct {
	client *external_event.Client
}

type externalEventDataSourceModel struct {
	Id                    types.Int64  `tfsdk:"id"`
	EventType             types.String `tfsdk:"event_type"`
	Status                types.String `tfsdk:"status"`
	Body                  types.String `tfsdk:"body"`
	CreatedAt             types.String `tfsdk:"created_at"`
	BodyUrl               types.String `tfsdk:"body_url"`
	FolderBehaviorId      types.Int64  `tfsdk:"folder_behavior_id"`
	SiemHttpDestinationId types.Int64  `tfsdk:"siem_http_destination_id"`
	SuccessfulFiles       types.Int64  `tfsdk:"successful_files"`
	ErroredFiles          types.Int64  `tfsdk:"errored_files"`
	BytesSynced           types.Int64  `tfsdk:"bytes_synced"`
	ComparedFiles         types.Int64  `tfsdk:"compared_files"`
	ComparedFolders       types.Int64  `tfsdk:"compared_folders"`
	RemoteServerType      types.String `tfsdk:"remote_server_type"`
}

func (r *externalEventDataSource) Configure(_ context.Context, req datasource.ConfigureRequest, resp *datasource.ConfigureResponse) {
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

	r.client = &external_event.Client{Config: sdk_config}
}

func (r *externalEventDataSource) Metadata(_ context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_external_event"
}

func (r *externalEventDataSource) Schema(_ context.Context, _ datasource.SchemaRequest, resp *datasource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Description: "An ExternalEvent is a log record with activity such as logins, credential syncs, and lockouts.",
		Attributes: map[string]schema.Attribute{
			"id": schema.Int64Attribute{
				Description: "Event ID",
				Required:    true,
			},
			"event_type": schema.StringAttribute{
				Description: "Type of event being recorded.",
				Computed:    true,
			},
			"status": schema.StringAttribute{
				Description: "Status of event.",
				Computed:    true,
			},
			"body": schema.StringAttribute{
				Description: "Event body",
				Computed:    true,
			},
			"created_at": schema.StringAttribute{
				Description: "External event create date/time",
				Computed:    true,
			},
			"body_url": schema.StringAttribute{
				Description: "Link to log file.",
				Computed:    true,
			},
			"folder_behavior_id": schema.Int64Attribute{
				Description: "Folder Behavior ID",
				Computed:    true,
			},
			"siem_http_destination_id": schema.Int64Attribute{
				Description: "SIEM HTTP Destination ID.",
				Computed:    true,
			},
			"successful_files": schema.Int64Attribute{
				Description: "For sync events, the number of files handled successfully.",
				Computed:    true,
			},
			"errored_files": schema.Int64Attribute{
				Description: "For sync events, the number of files that encountered errors.",
				Computed:    true,
			},
			"bytes_synced": schema.Int64Attribute{
				Description: "For sync events, the total number of bytes synced.",
				Computed:    true,
			},
			"compared_files": schema.Int64Attribute{
				Description: "For sync events, the number of files considered for the sync.",
				Computed:    true,
			},
			"compared_folders": schema.Int64Attribute{
				Description: "For sync events, the number of folders listed and considered for the sync.",
				Computed:    true,
			},
			"remote_server_type": schema.StringAttribute{
				Description: "Associated Remote Server type, if any",
				Computed:    true,
			},
		},
	}
}

func (r *externalEventDataSource) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
	var data externalEventDataSourceModel
	diags := req.Config.Get(ctx, &data)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	paramsExternalEventFind := files_sdk.ExternalEventFindParams{}
	paramsExternalEventFind.Id = data.Id.ValueInt64()

	externalEvent, err := r.client.Find(paramsExternalEventFind, files_sdk.WithContext(ctx))
	if err != nil {
		resp.Diagnostics.AddError(
			"Error Reading Files ExternalEvent",
			"Could not read external_event id "+fmt.Sprint(data.Id.ValueInt64())+": "+err.Error(),
		)
		return
	}

	diags = r.populateDataSourceModel(ctx, externalEvent, &data)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	diags = resp.State.Set(ctx, data)
	resp.Diagnostics.Append(diags...)
}

func (r *externalEventDataSource) populateDataSourceModel(ctx context.Context, externalEvent files_sdk.ExternalEvent, state *externalEventDataSourceModel) (diags diag.Diagnostics) {
	state.Id = types.Int64Value(externalEvent.Id)
	state.EventType = types.StringValue(externalEvent.EventType)
	state.Status = types.StringValue(externalEvent.Status)
	state.Body = types.StringValue(externalEvent.Body)
	if err := lib.TimeToStringType(ctx, path.Root("created_at"), externalEvent.CreatedAt, &state.CreatedAt); err != nil {
		diags.AddError(
			"Error Creating Files ExternalEvent",
			"Could not convert state created_at to string: "+err.Error(),
		)
	}
	state.BodyUrl = types.StringValue(externalEvent.BodyUrl)
	state.FolderBehaviorId = types.Int64Value(externalEvent.FolderBehaviorId)
	state.SiemHttpDestinationId = types.Int64Value(externalEvent.SiemHttpDestinationId)
	state.SuccessfulFiles = types.Int64Value(externalEvent.SuccessfulFiles)
	state.ErroredFiles = types.Int64Value(externalEvent.ErroredFiles)
	state.BytesSynced = types.Int64Value(externalEvent.BytesSynced)
	state.ComparedFiles = types.Int64Value(externalEvent.ComparedFiles)
	state.ComparedFolders = types.Int64Value(externalEvent.ComparedFolders)
	state.RemoteServerType = types.StringValue(externalEvent.RemoteServerType)

	return
}
