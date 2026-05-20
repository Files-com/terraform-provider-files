package provider

import (
	"context"
	"fmt"
	"strconv"
	"strings"

	files_sdk "github.com/Files-com/files-sdk-go/v3"
	event_target "github.com/Files-com/files-sdk-go/v3/eventtarget"
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
	"github.com/hashicorp/terraform-plugin-framework/schema/validator"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

var (
	_ resource.Resource                = &eventTargetResource{}
	_ resource.ResourceWithConfigure   = &eventTargetResource{}
	_ resource.ResourceWithImportState = &eventTargetResource{}
)

func NewEventTargetResource() resource.Resource {
	return &eventTargetResource{}
}

type eventTargetResource struct {
	client *event_target.Client
}

type eventTargetResourceModel struct {
	Name                 types.String  `tfsdk:"name"`
	TargetType           types.String  `tfsdk:"target_type"`
	Config               types.Dynamic `tfsdk:"config"`
	WorkspaceId          types.Int64   `tfsdk:"workspace_id"`
	ApplyToAllWorkspaces types.Bool    `tfsdk:"apply_to_all_workspaces"`
	Enabled              types.Bool    `tfsdk:"enabled"`
	DeliveryPolicy       types.Dynamic `tfsdk:"delivery_policy"`
	Id                   types.Int64   `tfsdk:"id"`
	CreatedAt            types.String  `tfsdk:"created_at"`
	UpdatedAt            types.String  `tfsdk:"updated_at"`
}

func (r *eventTargetResource) Configure(_ context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {
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

	r.client = &event_target.Client{Config: sdk_config}
}

func (r *eventTargetResource) Metadata(_ context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_event_target"
}

func (r *eventTargetResource) Schema(_ context.Context, _ resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Description: "An EventTarget is a delivery destination for EventRecords.",
		Attributes: map[string]schema.Attribute{
			"name": schema.StringAttribute{
				Description: "Event Target name.",
				Required:    true,
			},
			"target_type": schema.StringAttribute{
				Description: "Event Target type.",
				Required:    true,
				Validators: []validator.String{
					stringvalidator.OneOf("email", "webhook", "slack_webhook", "teams_webhook", "amazon_sns", "google_pubsub"),
				},
			},
			"config": schema.DynamicAttribute{
				Description: "Event Target configuration.",
				Required:    true,
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
				Description: "If true, this default-workspace target can receive events from all workspaces.",
				Computed:    true,
				Optional:    true,
				PlanModifiers: []planmodifier.Bool{
					boolplanmodifier.UseStateForUnknown(),
				},
			},
			"enabled": schema.BoolAttribute{
				Description: "Whether this Event Target can receive events.",
				Computed:    true,
				Optional:    true,
				PlanModifiers: []planmodifier.Bool{
					boolplanmodifier.UseStateForUnknown(),
				},
			},
			"delivery_policy": schema.DynamicAttribute{
				Description: "Event Target delivery policy. Email targets support batch_interval in seconds, between 600 and 86400.",
				Computed:    true,
				Optional:    true,
				PlanModifiers: []planmodifier.Dynamic{
					dynamicplanmodifier.UseStateForUnknown(),
				},
			},
			"id": schema.Int64Attribute{
				Description: "Event Target ID",
				Computed:    true,
				PlanModifiers: []planmodifier.Int64{
					int64planmodifier.UseStateForUnknown(),
				},
			},
			"created_at": schema.StringAttribute{
				Description: "Event Target create date/time.",
				Computed:    true,
			},
			"updated_at": schema.StringAttribute{
				Description: "Event Target update date/time.",
				Computed:    true,
			},
		},
	}
}

func (r *eventTargetResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	var plan eventTargetResourceModel
	diags := req.Plan.Get(ctx, &plan)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}
	var config eventTargetResourceModel
	diags = req.Config.Get(ctx, &config)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	paramsEventTargetCreate := files_sdk.EventTargetCreateParams{}
	paramsEventTargetCreate.Name = plan.Name.ValueString()
	paramsEventTargetCreate.WorkspaceId = plan.WorkspaceId.ValueInt64()
	if !plan.ApplyToAllWorkspaces.IsNull() && !plan.ApplyToAllWorkspaces.IsUnknown() {
		paramsEventTargetCreate.ApplyToAllWorkspaces = plan.ApplyToAllWorkspaces.ValueBoolPointer()
	}
	paramsEventTargetCreate.TargetType = paramsEventTargetCreate.TargetType.Enum()[plan.TargetType.ValueString()]
	if !plan.Enabled.IsNull() && !plan.Enabled.IsUnknown() {
		paramsEventTargetCreate.Enabled = plan.Enabled.ValueBoolPointer()
	}
	createConfig, diags := lib.DynamicToInterface(ctx, path.Root("config"), plan.Config)
	resp.Diagnostics.Append(diags...)
	paramsEventTargetCreate.Config = createConfig
	createDeliveryPolicy, diags := lib.DynamicToInterface(ctx, path.Root("delivery_policy"), plan.DeliveryPolicy)
	resp.Diagnostics.Append(diags...)
	paramsEventTargetCreate.DeliveryPolicy = createDeliveryPolicy

	if resp.Diagnostics.HasError() {
		return
	}

	eventTarget, err := r.client.Create(paramsEventTargetCreate, files_sdk.WithContext(ctx))
	if err != nil {
		resp.Diagnostics.AddError(
			"Error Creating Files EventTarget",
			"Could not create event_target, unexpected error: "+err.Error(),
		)
		return
	}

	diags = r.populateResourceModel(ctx, eventTarget, &plan)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	diags = resp.State.Set(ctx, plan)
	resp.Diagnostics.Append(diags...)
}

func (r *eventTargetResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	var state eventTargetResourceModel
	diags := req.State.Get(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	paramsEventTargetFind := files_sdk.EventTargetFindParams{}
	paramsEventTargetFind.Id = state.Id.ValueInt64()

	eventTarget, err := r.client.Find(paramsEventTargetFind, files_sdk.WithContext(ctx))
	if err != nil {
		if files_sdk.IsNotExist(err) {
			resp.State.RemoveResource(ctx)
			return
		}

		resp.Diagnostics.AddError(
			"Error Reading Files EventTarget",
			"Could not read event_target id "+fmt.Sprint(state.Id.ValueInt64())+": "+err.Error(),
		)
		return
	}

	diags = r.populateResourceModel(ctx, eventTarget, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	diags = resp.State.Set(ctx, state)
	resp.Diagnostics.Append(diags...)
}

func (r *eventTargetResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	var plan eventTargetResourceModel
	diags := req.Plan.Get(ctx, &plan)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}
	var config eventTargetResourceModel
	diags = req.Config.Get(ctx, &config)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	paramsEventTargetUpdate := map[string]interface{}{}
	if !plan.Id.IsNull() && !plan.Id.IsUnknown() {
		paramsEventTargetUpdate["id"] = plan.Id.ValueInt64()
	}
	if !config.Name.IsNull() && !config.Name.IsUnknown() {
		paramsEventTargetUpdate["name"] = config.Name.ValueString()
	}
	if !config.WorkspaceId.IsNull() && !config.WorkspaceId.IsUnknown() {
		paramsEventTargetUpdate["workspace_id"] = config.WorkspaceId.ValueInt64()
	}
	if !config.ApplyToAllWorkspaces.IsNull() && !config.ApplyToAllWorkspaces.IsUnknown() {
		paramsEventTargetUpdate["apply_to_all_workspaces"] = config.ApplyToAllWorkspaces.ValueBool()
	}
	if !config.TargetType.IsNull() && !config.TargetType.IsUnknown() {
		paramsEventTargetUpdate["target_type"] = config.TargetType.ValueString()
	}
	if !config.Enabled.IsNull() && !config.Enabled.IsUnknown() {
		paramsEventTargetUpdate["enabled"] = config.Enabled.ValueBool()
	}
	updateConfig, diags := lib.DynamicToInterface(ctx, path.Root("config"), config.Config)
	resp.Diagnostics.Append(diags...)
	paramsEventTargetUpdate["config"] = updateConfig
	updateDeliveryPolicy, diags := lib.DynamicToInterface(ctx, path.Root("delivery_policy"), config.DeliveryPolicy)
	resp.Diagnostics.Append(diags...)
	paramsEventTargetUpdate["delivery_policy"] = updateDeliveryPolicy

	if resp.Diagnostics.HasError() {
		return
	}

	eventTarget, err := r.client.UpdateWithMap(paramsEventTargetUpdate, files_sdk.WithContext(ctx))
	if err != nil {
		resp.Diagnostics.AddError(
			"Error Updating Files EventTarget",
			"Could not update event_target, unexpected error: "+err.Error(),
		)
		return
	}

	diags = r.populateResourceModel(ctx, eventTarget, &plan)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	diags = resp.State.Set(ctx, plan)
	resp.Diagnostics.Append(diags...)
}

func (r *eventTargetResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	var state eventTargetResourceModel
	diags := req.State.Get(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	paramsEventTargetDelete := files_sdk.EventTargetDeleteParams{}
	paramsEventTargetDelete.Id = state.Id.ValueInt64()

	err := r.client.Delete(paramsEventTargetDelete, files_sdk.WithContext(ctx))
	if err != nil && !files_sdk.IsNotExist(err) {
		resp.Diagnostics.AddError(
			"Error Deleting Files EventTarget",
			"Could not delete event_target id "+fmt.Sprint(state.Id.ValueInt64())+": "+err.Error(),
		)
	}
}

func (r *eventTargetResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
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

func (r *eventTargetResource) populateResourceModel(ctx context.Context, eventTarget files_sdk.EventTarget, state *eventTargetResourceModel) (diags diag.Diagnostics) {
	var propDiags diag.Diagnostics

	state.Id = types.Int64Value(eventTarget.Id)
	state.Name = types.StringValue(eventTarget.Name)
	state.TargetType = types.StringValue(eventTarget.TargetType)
	state.WorkspaceId = types.Int64Value(eventTarget.WorkspaceId)
	state.ApplyToAllWorkspaces = types.BoolPointerValue(eventTarget.ApplyToAllWorkspaces)
	state.Enabled = types.BoolPointerValue(eventTarget.Enabled)
	state.Config, propDiags = lib.ToDynamic(ctx, path.Root("config"), eventTarget.Config, state.Config.UnderlyingValue())
	diags.Append(propDiags...)
	state.DeliveryPolicy, propDiags = lib.ToDynamic(ctx, path.Root("delivery_policy"), eventTarget.DeliveryPolicy, state.DeliveryPolicy.UnderlyingValue())
	diags.Append(propDiags...)
	if err := lib.TimeToStringType(ctx, path.Root("created_at"), eventTarget.CreatedAt, &state.CreatedAt); err != nil {
		diags.AddError(
			"Error Creating Files EventTarget",
			"Could not convert state created_at to string: "+err.Error(),
		)
	}
	if err := lib.TimeToStringType(ctx, path.Root("updated_at"), eventTarget.UpdatedAt, &state.UpdatedAt); err != nil {
		diags.AddError(
			"Error Creating Files EventTarget",
			"Could not convert state updated_at to string: "+err.Error(),
		)
	}

	return
}
