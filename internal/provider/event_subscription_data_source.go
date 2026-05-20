package provider

import (
	"context"
	"fmt"

	files_sdk "github.com/Files-com/files-sdk-go/v3"
	event_subscription "github.com/Files-com/files-sdk-go/v3/eventsubscription"
	"github.com/Files-com/terraform-provider-files/lib"
	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/datasource/schema"
	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

var (
	_ datasource.DataSource              = &eventSubscriptionDataSource{}
	_ datasource.DataSourceWithConfigure = &eventSubscriptionDataSource{}
)

func NewEventSubscriptionDataSource() datasource.DataSource {
	return &eventSubscriptionDataSource{}
}

type eventSubscriptionDataSource struct {
	client *event_subscription.Client
}

type eventSubscriptionDataSourceModel struct {
	Id                   types.Int64   `tfsdk:"id"`
	EventChannelId       types.Int64   `tfsdk:"event_channel_id"`
	WorkspaceId          types.Int64   `tfsdk:"workspace_id"`
	ApplyToAllWorkspaces types.Bool    `tfsdk:"apply_to_all_workspaces"`
	Name                 types.String  `tfsdk:"name"`
	Enabled              types.Bool    `tfsdk:"enabled"`
	EventTypes           types.List    `tfsdk:"event_types"`
	Filter               types.Dynamic `tfsdk:"filter"`
	DeliveryPolicy       types.Dynamic `tfsdk:"delivery_policy"`
	EventTargetIds       types.List    `tfsdk:"event_target_ids"`
	CreatedAt            types.String  `tfsdk:"created_at"`
	UpdatedAt            types.String  `tfsdk:"updated_at"`
}

func (r *eventSubscriptionDataSource) Configure(_ context.Context, req datasource.ConfigureRequest, resp *datasource.ConfigureResponse) {
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

	r.client = &event_subscription.Client{Config: sdk_config}
}

func (r *eventSubscriptionDataSource) Metadata(_ context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_event_subscription"
}

func (r *eventSubscriptionDataSource) Schema(_ context.Context, _ datasource.SchemaRequest, resp *datasource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Description: "An EventSubscription selects EventRecords for an EventChannel and sends them to one or more EventTargets.",
		Attributes: map[string]schema.Attribute{
			"id": schema.Int64Attribute{
				Description: "Event Subscription ID",
				Required:    true,
			},
			"event_channel_id": schema.Int64Attribute{
				Description: "Event Channel ID",
				Computed:    true,
			},
			"workspace_id": schema.Int64Attribute{
				Description: "Workspace ID. 0 means the default workspace or site-wide.",
				Computed:    true,
			},
			"apply_to_all_workspaces": schema.BoolAttribute{
				Description: "If true, this default-workspace subscription applies to events from all workspaces.",
				Computed:    true,
			},
			"name": schema.StringAttribute{
				Description: "Event Subscription name.",
				Computed:    true,
			},
			"enabled": schema.BoolAttribute{
				Description: "Whether this Event Subscription can dispatch events.",
				Computed:    true,
			},
			"event_types": schema.ListAttribute{
				Description: "Event type strings matched by this subscription. Blank means all event types.",
				Computed:    true,
				ElementType: types.StringType,
			},
			"filter": schema.DynamicAttribute{
				Description: "Structured event payload filter.",
				Computed:    true,
			},
			"delivery_policy": schema.DynamicAttribute{
				Description: "Event Subscription delivery policy.",
				Computed:    true,
			},
			"event_target_ids": schema.ListAttribute{
				Description: "Event Target IDs this subscription sends to.",
				Computed:    true,
				ElementType: types.Int64Type,
			},
			"created_at": schema.StringAttribute{
				Description: "Event Subscription create date/time.",
				Computed:    true,
			},
			"updated_at": schema.StringAttribute{
				Description: "Event Subscription update date/time.",
				Computed:    true,
			},
		},
	}
}

func (r *eventSubscriptionDataSource) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
	var data eventSubscriptionDataSourceModel
	diags := req.Config.Get(ctx, &data)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	paramsEventSubscriptionFind := files_sdk.EventSubscriptionFindParams{}
	paramsEventSubscriptionFind.Id = data.Id.ValueInt64()

	eventSubscription, err := r.client.Find(paramsEventSubscriptionFind, files_sdk.WithContext(ctx))
	if err != nil {
		resp.Diagnostics.AddError(
			"Error Reading Files EventSubscription",
			"Could not read event_subscription id "+fmt.Sprint(data.Id.ValueInt64())+": "+err.Error(),
		)
		return
	}

	diags = r.populateDataSourceModel(ctx, eventSubscription, &data)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	diags = resp.State.Set(ctx, data)
	resp.Diagnostics.Append(diags...)
}

func (r *eventSubscriptionDataSource) populateDataSourceModel(ctx context.Context, eventSubscription files_sdk.EventSubscription, state *eventSubscriptionDataSourceModel) (diags diag.Diagnostics) {
	var propDiags diag.Diagnostics

	state.Id = types.Int64Value(eventSubscription.Id)
	state.EventChannelId = types.Int64Value(eventSubscription.EventChannelId)
	state.WorkspaceId = types.Int64Value(eventSubscription.WorkspaceId)
	state.ApplyToAllWorkspaces = types.BoolPointerValue(eventSubscription.ApplyToAllWorkspaces)
	state.Name = types.StringValue(eventSubscription.Name)
	state.Enabled = types.BoolPointerValue(eventSubscription.Enabled)
	state.EventTypes, propDiags = types.ListValueFrom(ctx, types.StringType, eventSubscription.EventTypes)
	diags.Append(propDiags...)
	state.Filter, propDiags = lib.ToDynamic(ctx, path.Root("filter"), eventSubscription.Filter, state.Filter.UnderlyingValue())
	diags.Append(propDiags...)
	state.DeliveryPolicy, propDiags = lib.ToDynamic(ctx, path.Root("delivery_policy"), eventSubscription.DeliveryPolicy, state.DeliveryPolicy.UnderlyingValue())
	diags.Append(propDiags...)
	state.EventTargetIds, propDiags = types.ListValueFrom(ctx, types.Int64Type, eventSubscription.EventTargetIds)
	diags.Append(propDiags...)
	if err := lib.TimeToStringType(ctx, path.Root("created_at"), eventSubscription.CreatedAt, &state.CreatedAt); err != nil {
		diags.AddError(
			"Error Creating Files EventSubscription",
			"Could not convert state created_at to string: "+err.Error(),
		)
	}
	if err := lib.TimeToStringType(ctx, path.Root("updated_at"), eventSubscription.UpdatedAt, &state.UpdatedAt); err != nil {
		diags.AddError(
			"Error Creating Files EventSubscription",
			"Could not convert state updated_at to string: "+err.Error(),
		)
	}

	return
}
