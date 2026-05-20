package provider

import (
	"context"
	"fmt"

	files_sdk "github.com/Files-com/files-sdk-go/v3"
	event_channel "github.com/Files-com/files-sdk-go/v3/eventchannel"
	"github.com/Files-com/terraform-provider-files/lib"
	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/datasource/schema"
	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

var (
	_ datasource.DataSource              = &eventChannelDataSource{}
	_ datasource.DataSourceWithConfigure = &eventChannelDataSource{}
)

func NewEventChannelDataSource() datasource.DataSource {
	return &eventChannelDataSource{}
}

type eventChannelDataSource struct {
	client *event_channel.Client
}

type eventChannelDataSourceModel struct {
	Id             types.Int64  `tfsdk:"id"`
	Name           types.String `tfsdk:"name"`
	Description    types.String `tfsdk:"description"`
	Enabled        types.Bool   `tfsdk:"enabled"`
	DefaultChannel types.Bool   `tfsdk:"default_channel"`
	CreatedAt      types.String `tfsdk:"created_at"`
	UpdatedAt      types.String `tfsdk:"updated_at"`
}

func (r *eventChannelDataSource) Configure(_ context.Context, req datasource.ConfigureRequest, resp *datasource.ConfigureResponse) {
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

	r.client = &event_channel.Client{Config: sdk_config}
}

func (r *eventChannelDataSource) Metadata(_ context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_event_channel"
}

func (r *eventChannelDataSource) Schema(_ context.Context, _ datasource.SchemaRequest, resp *datasource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Description: "An EventChannel is a named grouping of EventSubscriptions.",
		Attributes: map[string]schema.Attribute{
			"id": schema.Int64Attribute{
				Description: "Event Channel ID",
				Required:    true,
			},
			"name": schema.StringAttribute{
				Description: "Event Channel name.",
				Computed:    true,
			},
			"description": schema.StringAttribute{
				Description: "Event Channel description.",
				Computed:    true,
			},
			"enabled": schema.BoolAttribute{
				Description: "Whether this Event Channel can dispatch events.",
				Computed:    true,
			},
			"default_channel": schema.BoolAttribute{
				Description: "Whether this Event Channel is the default destination for newly published events.",
				Computed:    true,
			},
			"created_at": schema.StringAttribute{
				Description: "Event Channel create date/time.",
				Computed:    true,
			},
			"updated_at": schema.StringAttribute{
				Description: "Event Channel update date/time.",
				Computed:    true,
			},
		},
	}
}

func (r *eventChannelDataSource) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
	var data eventChannelDataSourceModel
	diags := req.Config.Get(ctx, &data)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	paramsEventChannelFind := files_sdk.EventChannelFindParams{}
	paramsEventChannelFind.Id = data.Id.ValueInt64()

	eventChannel, err := r.client.Find(paramsEventChannelFind, files_sdk.WithContext(ctx))
	if err != nil {
		resp.Diagnostics.AddError(
			"Error Reading Files EventChannel",
			"Could not read event_channel id "+fmt.Sprint(data.Id.ValueInt64())+": "+err.Error(),
		)
		return
	}

	diags = r.populateDataSourceModel(ctx, eventChannel, &data)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	diags = resp.State.Set(ctx, data)
	resp.Diagnostics.Append(diags...)
}

func (r *eventChannelDataSource) populateDataSourceModel(ctx context.Context, eventChannel files_sdk.EventChannel, state *eventChannelDataSourceModel) (diags diag.Diagnostics) {
	state.Id = types.Int64Value(eventChannel.Id)
	state.Name = types.StringValue(eventChannel.Name)
	state.Description = types.StringValue(eventChannel.Description)
	state.Enabled = types.BoolPointerValue(eventChannel.Enabled)
	state.DefaultChannel = types.BoolPointerValue(eventChannel.DefaultChannel)
	if err := lib.TimeToStringType(ctx, path.Root("created_at"), eventChannel.CreatedAt, &state.CreatedAt); err != nil {
		diags.AddError(
			"Error Creating Files EventChannel",
			"Could not convert state created_at to string: "+err.Error(),
		)
	}
	if err := lib.TimeToStringType(ctx, path.Root("updated_at"), eventChannel.UpdatedAt, &state.UpdatedAt); err != nil {
		diags.AddError(
			"Error Creating Files EventChannel",
			"Could not convert state updated_at to string: "+err.Error(),
		)
	}

	return
}
