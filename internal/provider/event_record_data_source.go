package provider

import (
	"context"
	"fmt"

	files_sdk "github.com/Files-com/files-sdk-go/v3"
	event_record "github.com/Files-com/files-sdk-go/v3/eventrecord"
	"github.com/Files-com/terraform-provider-files/lib"
	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/datasource/schema"
	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

var (
	_ datasource.DataSource              = &eventRecordDataSource{}
	_ datasource.DataSourceWithConfigure = &eventRecordDataSource{}
)

func NewEventRecordDataSource() datasource.DataSource {
	return &eventRecordDataSource{}
}

type eventRecordDataSource struct {
	client *event_record.Client
}

type eventRecordDataSourceModel struct {
	Id           types.Int64   `tfsdk:"id"`
	WorkspaceId  types.Int64   `tfsdk:"workspace_id"`
	EventUuid    types.String  `tfsdk:"event_uuid"`
	EventType    types.String  `tfsdk:"event_type"`
	Severity     types.String  `tfsdk:"severity"`
	SourceType   types.String  `tfsdk:"source_type"`
	SourceId     types.Int64   `tfsdk:"source_id"`
	OccurredAt   types.String  `tfsdk:"occurred_at"`
	HumanTitle   types.String  `tfsdk:"human_title"`
	HumanSummary types.String  `tfsdk:"human_summary"`
	HumanFields  types.Dynamic `tfsdk:"human_fields"`
	Actor        types.Dynamic `tfsdk:"actor"`
	Resources    types.Dynamic `tfsdk:"resources"`
	Payload      types.Dynamic `tfsdk:"payload"`
	CreatedAt    types.String  `tfsdk:"created_at"`
}

func (r *eventRecordDataSource) Configure(_ context.Context, req datasource.ConfigureRequest, resp *datasource.ConfigureResponse) {
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

	r.client = &event_record.Client{Config: sdk_config}
}

func (r *eventRecordDataSource) Metadata(_ context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_event_record"
}

func (r *eventRecordDataSource) Schema(_ context.Context, _ datasource.SchemaRequest, resp *datasource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Description: "An EventRecord is a durable event emitted by Files.com for routing through Event Channels.",
		Attributes: map[string]schema.Attribute{
			"id": schema.Int64Attribute{
				Description: "Event Record ID",
				Required:    true,
			},
			"workspace_id": schema.Int64Attribute{
				Description: "Workspace ID. 0 means the default workspace or site-wide.",
				Computed:    true,
			},
			"event_uuid": schema.StringAttribute{
				Description: "Stable event UUID.",
				Computed:    true,
			},
			"event_type": schema.StringAttribute{
				Description: "Versioned event type string.",
				Computed:    true,
			},
			"severity": schema.StringAttribute{
				Description: "Event severity.",
				Computed:    true,
			},
			"source_type": schema.StringAttribute{
				Description: "Source record type.",
				Computed:    true,
			},
			"source_id": schema.Int64Attribute{
				Description: "Source record ID.",
				Computed:    true,
			},
			"occurred_at": schema.StringAttribute{
				Description: "Event occurrence date/time.",
				Computed:    true,
			},
			"human_title": schema.StringAttribute{
				Description: "Human-readable event title.",
				Computed:    true,
			},
			"human_summary": schema.StringAttribute{
				Description: "Human-readable event summary.",
				Computed:    true,
			},
			"human_fields": schema.DynamicAttribute{
				Description: "Human-readable event detail fields.",
				Computed:    true,
			},
			"actor": schema.DynamicAttribute{
				Description: "Actor associated with the event.",
				Computed:    true,
			},
			"resources": schema.DynamicAttribute{
				Description: "Resources associated with the event.",
				Computed:    true,
			},
			"payload": schema.DynamicAttribute{
				Description: "Event payload.",
				Computed:    true,
			},
			"created_at": schema.StringAttribute{
				Description: "Event Record create date/time.",
				Computed:    true,
			},
		},
	}
}

func (r *eventRecordDataSource) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
	var data eventRecordDataSourceModel
	diags := req.Config.Get(ctx, &data)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	paramsEventRecordFind := files_sdk.EventRecordFindParams{}
	paramsEventRecordFind.Id = data.Id.ValueInt64()

	eventRecord, err := r.client.Find(paramsEventRecordFind, files_sdk.WithContext(ctx))
	if err != nil {
		resp.Diagnostics.AddError(
			"Error Reading Files EventRecord",
			"Could not read event_record id "+fmt.Sprint(data.Id.ValueInt64())+": "+err.Error(),
		)
		return
	}

	diags = r.populateDataSourceModel(ctx, eventRecord, &data)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	diags = resp.State.Set(ctx, data)
	resp.Diagnostics.Append(diags...)
}

func (r *eventRecordDataSource) populateDataSourceModel(ctx context.Context, eventRecord files_sdk.EventRecord, state *eventRecordDataSourceModel) (diags diag.Diagnostics) {
	var propDiags diag.Diagnostics

	state.Id = types.Int64Value(eventRecord.Id)
	state.WorkspaceId = types.Int64Value(eventRecord.WorkspaceId)
	state.EventUuid = types.StringValue(eventRecord.EventUuid)
	state.EventType = types.StringValue(eventRecord.EventType)
	state.Severity = types.StringValue(eventRecord.Severity)
	state.SourceType = types.StringValue(eventRecord.SourceType)
	state.SourceId = types.Int64Value(eventRecord.SourceId)
	if err := lib.TimeToStringType(ctx, path.Root("occurred_at"), eventRecord.OccurredAt, &state.OccurredAt); err != nil {
		diags.AddError(
			"Error Creating Files EventRecord",
			"Could not convert state occurred_at to string: "+err.Error(),
		)
	}
	state.HumanTitle = types.StringValue(eventRecord.HumanTitle)
	state.HumanSummary = types.StringValue(eventRecord.HumanSummary)
	state.HumanFields, propDiags = lib.ToDynamic(ctx, path.Root("human_fields"), eventRecord.HumanFields, state.HumanFields.UnderlyingValue())
	diags.Append(propDiags...)
	state.Actor, propDiags = lib.ToDynamic(ctx, path.Root("actor"), eventRecord.Actor, state.Actor.UnderlyingValue())
	diags.Append(propDiags...)
	state.Resources, propDiags = lib.ToDynamic(ctx, path.Root("resources"), eventRecord.Resources, state.Resources.UnderlyingValue())
	diags.Append(propDiags...)
	state.Payload, propDiags = lib.ToDynamic(ctx, path.Root("payload"), eventRecord.Payload, state.Payload.UnderlyingValue())
	diags.Append(propDiags...)
	if err := lib.TimeToStringType(ctx, path.Root("created_at"), eventRecord.CreatedAt, &state.CreatedAt); err != nil {
		diags.AddError(
			"Error Creating Files EventRecord",
			"Could not convert state created_at to string: "+err.Error(),
		)
	}

	return
}
