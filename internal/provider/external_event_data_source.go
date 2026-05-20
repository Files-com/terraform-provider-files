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
	Id        types.Int64  `tfsdk:"id"`
	EventType types.String `tfsdk:"event_type"`
	Status    types.String `tfsdk:"status"`
	Body      types.String `tfsdk:"body"`
	CreatedAt types.String `tfsdk:"created_at"`
	BodyUrl   types.String `tfsdk:"body_url"`
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
		Description: "An ExternalEvent is a log that is sent to the cloud from a client application such as the Files.com CLI.",
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

	return
}
