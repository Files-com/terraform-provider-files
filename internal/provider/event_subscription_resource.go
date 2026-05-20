package provider

import (
	"context"
	"fmt"
	"strconv"
	"strings"

	files_sdk "github.com/Files-com/files-sdk-go/v3"
	event_subscription "github.com/Files-com/files-sdk-go/v3/eventsubscription"
	"github.com/Files-com/terraform-provider-files/lib"
	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/boolplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/dynamicplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/int64planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/listplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

var (
	_ resource.Resource                = &eventSubscriptionResource{}
	_ resource.ResourceWithConfigure   = &eventSubscriptionResource{}
	_ resource.ResourceWithImportState = &eventSubscriptionResource{}
)

func NewEventSubscriptionResource() resource.Resource {
	return &eventSubscriptionResource{}
}

type eventSubscriptionResource struct {
	client *event_subscription.Client
}

type eventSubscriptionResourceModel struct {
	Name                 types.String  `tfsdk:"name"`
	EventChannelId       types.Int64   `tfsdk:"event_channel_id"`
	WorkspaceId          types.Int64   `tfsdk:"workspace_id"`
	ApplyToAllWorkspaces types.Bool    `tfsdk:"apply_to_all_workspaces"`
	Enabled              types.Bool    `tfsdk:"enabled"`
	EventTypes           types.List    `tfsdk:"event_types"`
	Filter               types.Dynamic `tfsdk:"filter"`
	DeliveryPolicy       types.Dynamic `tfsdk:"delivery_policy"`
	EventTargetIds       types.List    `tfsdk:"event_target_ids"`
	Id                   types.Int64   `tfsdk:"id"`
	CreatedAt            types.String  `tfsdk:"created_at"`
	UpdatedAt            types.String  `tfsdk:"updated_at"`
}

func (r *eventSubscriptionResource) Configure(_ context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {
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

func (r *eventSubscriptionResource) Metadata(_ context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_event_subscription"
}

func (r *eventSubscriptionResource) Schema(_ context.Context, _ resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Description: "An EventSubscription selects EventRecords for an EventChannel and sends them to one or more EventTargets.",
		Attributes: map[string]schema.Attribute{
			"name": schema.StringAttribute{
				Description: "Event Subscription name.",
				Required:    true,
			},
			"event_channel_id": schema.Int64Attribute{
				Description: "Event Channel ID",
				Computed:    true,
				Optional:    true,
				PlanModifiers: []planmodifier.Int64{
					int64planmodifier.UseStateForUnknown(),
				},
			},
			"workspace_id": schema.Int64Attribute{
				Description: "Workspace ID. 0 means the default workspace or site-wide.",
				Computed:    true,
				Optional:    true,
				PlanModifiers: []planmodifier.Int64{
					int64planmodifier.UseStateForUnknown(),
				},
			},
			"apply_to_all_workspaces": schema.BoolAttribute{
				Description: "If true, this default-workspace subscription applies to events from all workspaces.",
				Computed:    true,
				Optional:    true,
				PlanModifiers: []planmodifier.Bool{
					boolplanmodifier.UseStateForUnknown(),
				},
			},
			"enabled": schema.BoolAttribute{
				Description: "Whether this Event Subscription can dispatch events.",
				Computed:    true,
				Optional:    true,
				PlanModifiers: []planmodifier.Bool{
					boolplanmodifier.UseStateForUnknown(),
				},
			},
			"event_types": schema.ListAttribute{
				Description: "Event type strings matched by this subscription. Blank means all event types.",
				Computed:    true,
				Optional:    true,
				ElementType: types.StringType,
				PlanModifiers: []planmodifier.List{
					listplanmodifier.UseStateForUnknown(),
				},
			},
			"filter": schema.DynamicAttribute{
				Description: "Structured event payload filter.",
				Computed:    true,
				Optional:    true,
				PlanModifiers: []planmodifier.Dynamic{
					dynamicplanmodifier.UseStateForUnknown(),
				},
			},
			"delivery_policy": schema.DynamicAttribute{
				Description: "Event Subscription delivery policy.",
				Computed:    true,
				Optional:    true,
				PlanModifiers: []planmodifier.Dynamic{
					dynamicplanmodifier.UseStateForUnknown(),
				},
			},
			"event_target_ids": schema.ListAttribute{
				Description: "Event Target IDs this subscription sends to.",
				Computed:    true,
				Optional:    true,
				ElementType: types.Int64Type,
				PlanModifiers: []planmodifier.List{
					listplanmodifier.UseStateForUnknown(),
				},
			},
			"id": schema.Int64Attribute{
				Description: "Event Subscription ID",
				Computed:    true,
				PlanModifiers: []planmodifier.Int64{
					int64planmodifier.UseStateForUnknown(),
				},
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

func (r *eventSubscriptionResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	var plan eventSubscriptionResourceModel
	diags := req.Plan.Get(ctx, &plan)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}
	var config eventSubscriptionResourceModel
	diags = req.Config.Get(ctx, &config)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	paramsEventSubscriptionCreate := files_sdk.EventSubscriptionCreateParams{}
	paramsEventSubscriptionCreate.EventChannelId = plan.EventChannelId.ValueInt64()
	paramsEventSubscriptionCreate.WorkspaceId = plan.WorkspaceId.ValueInt64()
	if !plan.ApplyToAllWorkspaces.IsNull() && !plan.ApplyToAllWorkspaces.IsUnknown() {
		paramsEventSubscriptionCreate.ApplyToAllWorkspaces = plan.ApplyToAllWorkspaces.ValueBoolPointer()
	}
	paramsEventSubscriptionCreate.Name = plan.Name.ValueString()
	if !plan.Enabled.IsNull() && !plan.Enabled.IsUnknown() {
		paramsEventSubscriptionCreate.Enabled = plan.Enabled.ValueBoolPointer()
	}
	if !plan.EventTypes.IsNull() && !plan.EventTypes.IsUnknown() {
		diags = plan.EventTypes.ElementsAs(ctx, &paramsEventSubscriptionCreate.EventTypes, false)
		resp.Diagnostics.Append(diags...)
	}
	createFilter, diags := lib.DynamicToInterface(ctx, path.Root("filter"), plan.Filter)
	resp.Diagnostics.Append(diags...)
	paramsEventSubscriptionCreate.Filter = createFilter
	createDeliveryPolicy, diags := lib.DynamicToInterface(ctx, path.Root("delivery_policy"), plan.DeliveryPolicy)
	resp.Diagnostics.Append(diags...)
	paramsEventSubscriptionCreate.DeliveryPolicy = createDeliveryPolicy
	if !plan.EventTargetIds.IsNull() && !plan.EventTargetIds.IsUnknown() {
		diags = plan.EventTargetIds.ElementsAs(ctx, &paramsEventSubscriptionCreate.EventTargetIds, false)
		resp.Diagnostics.Append(diags...)
	}

	if resp.Diagnostics.HasError() {
		return
	}

	eventSubscription, err := r.client.Create(paramsEventSubscriptionCreate, files_sdk.WithContext(ctx))
	if err != nil {
		resp.Diagnostics.AddError(
			"Error Creating Files EventSubscription",
			"Could not create event_subscription, unexpected error: "+err.Error(),
		)
		return
	}

	diags = r.populateResourceModel(ctx, eventSubscription, &plan)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	diags = resp.State.Set(ctx, plan)
	resp.Diagnostics.Append(diags...)
}

func (r *eventSubscriptionResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	var state eventSubscriptionResourceModel
	diags := req.State.Get(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	paramsEventSubscriptionFind := files_sdk.EventSubscriptionFindParams{}
	paramsEventSubscriptionFind.Id = state.Id.ValueInt64()

	eventSubscription, err := r.client.Find(paramsEventSubscriptionFind, files_sdk.WithContext(ctx))
	if err != nil {
		if files_sdk.IsNotExist(err) {
			resp.State.RemoveResource(ctx)
			return
		}

		resp.Diagnostics.AddError(
			"Error Reading Files EventSubscription",
			"Could not read event_subscription id "+fmt.Sprint(state.Id.ValueInt64())+": "+err.Error(),
		)
		return
	}

	diags = r.populateResourceModel(ctx, eventSubscription, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	diags = resp.State.Set(ctx, state)
	resp.Diagnostics.Append(diags...)
}

func (r *eventSubscriptionResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	var plan eventSubscriptionResourceModel
	diags := req.Plan.Get(ctx, &plan)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}
	var config eventSubscriptionResourceModel
	diags = req.Config.Get(ctx, &config)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	paramsEventSubscriptionUpdate := map[string]interface{}{}
	if !plan.Id.IsNull() && !plan.Id.IsUnknown() {
		paramsEventSubscriptionUpdate["id"] = plan.Id.ValueInt64()
	}
	if !config.EventChannelId.IsNull() && !config.EventChannelId.IsUnknown() {
		paramsEventSubscriptionUpdate["event_channel_id"] = config.EventChannelId.ValueInt64()
	}
	if !config.WorkspaceId.IsNull() && !config.WorkspaceId.IsUnknown() {
		paramsEventSubscriptionUpdate["workspace_id"] = config.WorkspaceId.ValueInt64()
	}
	if !config.ApplyToAllWorkspaces.IsNull() && !config.ApplyToAllWorkspaces.IsUnknown() {
		paramsEventSubscriptionUpdate["apply_to_all_workspaces"] = config.ApplyToAllWorkspaces.ValueBool()
	}
	if !config.Name.IsNull() && !config.Name.IsUnknown() {
		paramsEventSubscriptionUpdate["name"] = config.Name.ValueString()
	}
	if !config.Enabled.IsNull() && !config.Enabled.IsUnknown() {
		paramsEventSubscriptionUpdate["enabled"] = config.Enabled.ValueBool()
	}
	if !config.EventTypes.IsNull() && !config.EventTypes.IsUnknown() {
		var updateEventTypes []string
		diags = config.EventTypes.ElementsAs(ctx, &updateEventTypes, false)
		resp.Diagnostics.Append(diags...)
		paramsEventSubscriptionUpdate["event_types"] = updateEventTypes
	}
	updateFilter, diags := lib.DynamicToInterface(ctx, path.Root("filter"), config.Filter)
	resp.Diagnostics.Append(diags...)
	paramsEventSubscriptionUpdate["filter"] = updateFilter
	updateDeliveryPolicy, diags := lib.DynamicToInterface(ctx, path.Root("delivery_policy"), config.DeliveryPolicy)
	resp.Diagnostics.Append(diags...)
	paramsEventSubscriptionUpdate["delivery_policy"] = updateDeliveryPolicy
	if !config.EventTargetIds.IsNull() && !config.EventTargetIds.IsUnknown() {
		var updateEventTargetIds []int64
		diags = config.EventTargetIds.ElementsAs(ctx, &updateEventTargetIds, false)
		resp.Diagnostics.Append(diags...)
		paramsEventSubscriptionUpdate["event_target_ids"] = updateEventTargetIds
	}

	if resp.Diagnostics.HasError() {
		return
	}

	eventSubscription, err := r.client.UpdateWithMap(paramsEventSubscriptionUpdate, files_sdk.WithContext(ctx))
	if err != nil {
		resp.Diagnostics.AddError(
			"Error Updating Files EventSubscription",
			"Could not update event_subscription, unexpected error: "+err.Error(),
		)
		return
	}

	diags = r.populateResourceModel(ctx, eventSubscription, &plan)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	diags = resp.State.Set(ctx, plan)
	resp.Diagnostics.Append(diags...)
}

func (r *eventSubscriptionResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	var state eventSubscriptionResourceModel
	diags := req.State.Get(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	paramsEventSubscriptionDelete := files_sdk.EventSubscriptionDeleteParams{}
	paramsEventSubscriptionDelete.Id = state.Id.ValueInt64()

	err := r.client.Delete(paramsEventSubscriptionDelete, files_sdk.WithContext(ctx))
	if err != nil && !files_sdk.IsNotExist(err) {
		resp.Diagnostics.AddError(
			"Error Deleting Files EventSubscription",
			"Could not delete event_subscription id "+fmt.Sprint(state.Id.ValueInt64())+": "+err.Error(),
		)
	}
}

func (r *eventSubscriptionResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
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

func (r *eventSubscriptionResource) populateResourceModel(ctx context.Context, eventSubscription files_sdk.EventSubscription, state *eventSubscriptionResourceModel) (diags diag.Diagnostics) {
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
