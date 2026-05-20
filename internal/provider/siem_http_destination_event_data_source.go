package provider

import (
	"context"
	"fmt"

	files_sdk "github.com/Files-com/files-sdk-go/v3"
	siem_http_destination_event "github.com/Files-com/files-sdk-go/v3/siemhttpdestinationevent"
	"github.com/Files-com/terraform-provider-files/lib"
	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/datasource/schema"
	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

var (
	_ datasource.DataSource              = &siemHttpDestinationEventDataSource{}
	_ datasource.DataSourceWithConfigure = &siemHttpDestinationEventDataSource{}
)

func NewSiemHttpDestinationEventDataSource() datasource.DataSource {
	return &siemHttpDestinationEventDataSource{}
}

type siemHttpDestinationEventDataSource struct {
	client *siem_http_destination_event.Client
}

type siemHttpDestinationEventDataSourceModel struct {
	Id                    types.Int64  `tfsdk:"id"`
	EventType             types.String `tfsdk:"event_type"`
	Status                types.String `tfsdk:"status"`
	Body                  types.String `tfsdk:"body"`
	EventErrors           types.List   `tfsdk:"event_errors"`
	CreatedAt             types.String `tfsdk:"created_at"`
	BodyUrl               types.String `tfsdk:"body_url"`
	SiemHttpDestinationId types.Int64  `tfsdk:"siem_http_destination_id"`
}

func (r *siemHttpDestinationEventDataSource) Configure(_ context.Context, req datasource.ConfigureRequest, resp *datasource.ConfigureResponse) {
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

	r.client = &siem_http_destination_event.Client{Config: sdk_config}
}

func (r *siemHttpDestinationEventDataSource) Metadata(_ context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_siem_http_destination_event"
}

func (r *siemHttpDestinationEventDataSource) Schema(_ context.Context, _ datasource.SchemaRequest, resp *datasource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Description: "A SiemHttpDestinationEvent is a log record for SIEM publishing failures and recoveries.",
		Attributes: map[string]schema.Attribute{
			"id": schema.Int64Attribute{
				Description: "Event ID",
				Required:    true,
			},
			"event_type": schema.StringAttribute{
				Description: "Type of SIEM event being recorded.",
				Computed:    true,
			},
			"status": schema.StringAttribute{
				Description: "Status of event.",
				Computed:    true,
			},
			"body": schema.StringAttribute{
				Description: "Event body.",
				Computed:    true,
			},
			"event_errors": schema.ListAttribute{
				Description: "Event errors.",
				Computed:    true,
				ElementType: types.StringType,
			},
			"created_at": schema.StringAttribute{
				Description: "Event create date/time.",
				Computed:    true,
			},
			"body_url": schema.StringAttribute{
				Description: "Link to log file.",
				Computed:    true,
			},
			"siem_http_destination_id": schema.Int64Attribute{
				Description: "SIEM ID.",
				Computed:    true,
			},
		},
	}
}

func (r *siemHttpDestinationEventDataSource) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
	var data siemHttpDestinationEventDataSourceModel
	diags := req.Config.Get(ctx, &data)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	paramsSiemHttpDestinationEventFind := files_sdk.SiemHttpDestinationEventFindParams{}
	paramsSiemHttpDestinationEventFind.Id = data.Id.ValueInt64()

	siemHttpDestinationEvent, err := r.client.Find(paramsSiemHttpDestinationEventFind, files_sdk.WithContext(ctx))
	if err != nil {
		resp.Diagnostics.AddError(
			"Error Reading Files SiemHttpDestinationEvent",
			"Could not read siem_http_destination_event id "+fmt.Sprint(data.Id.ValueInt64())+": "+err.Error(),
		)
		return
	}

	diags = r.populateDataSourceModel(ctx, siemHttpDestinationEvent, &data)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	diags = resp.State.Set(ctx, data)
	resp.Diagnostics.Append(diags...)
}

func (r *siemHttpDestinationEventDataSource) populateDataSourceModel(ctx context.Context, siemHttpDestinationEvent files_sdk.SiemHttpDestinationEvent, state *siemHttpDestinationEventDataSourceModel) (diags diag.Diagnostics) {
	var propDiags diag.Diagnostics

	state.Id = types.Int64Value(siemHttpDestinationEvent.Id)
	state.EventType = types.StringValue(siemHttpDestinationEvent.EventType)
	state.Status = types.StringValue(siemHttpDestinationEvent.Status)
	state.Body = types.StringValue(siemHttpDestinationEvent.Body)
	state.EventErrors, propDiags = types.ListValueFrom(ctx, types.StringType, siemHttpDestinationEvent.EventErrors)
	diags.Append(propDiags...)
	if err := lib.TimeToStringType(ctx, path.Root("created_at"), siemHttpDestinationEvent.CreatedAt, &state.CreatedAt); err != nil {
		diags.AddError(
			"Error Creating Files SiemHttpDestinationEvent",
			"Could not convert state created_at to string: "+err.Error(),
		)
	}
	state.BodyUrl = types.StringValue(siemHttpDestinationEvent.BodyUrl)
	state.SiemHttpDestinationId = types.Int64Value(siemHttpDestinationEvent.SiemHttpDestinationId)

	return
}
