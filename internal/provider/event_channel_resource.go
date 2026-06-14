package provider

import (
	"context"
	"fmt"
	"strconv"
	"strings"

	files_sdk "github.com/Files-com/files-sdk-go/v3"
	event_channel "github.com/Files-com/files-sdk-go/v3/eventchannel"
	"github.com/Files-com/terraform-provider-files/lib"
	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/boolplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/int64planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

var (
	_ resource.Resource                = &eventChannelResource{}
	_ resource.ResourceWithConfigure   = &eventChannelResource{}
	_ resource.ResourceWithImportState = &eventChannelResource{}
)

func NewEventChannelResource() resource.Resource {
	return &eventChannelResource{}
}

type eventChannelResource struct {
	client *event_channel.Client
}

type eventChannelResourceModel struct {
	Name           types.String `tfsdk:"name"`
	WorkspaceId    types.Int64  `tfsdk:"workspace_id"`
	Description    types.String `tfsdk:"description"`
	Enabled        types.Bool   `tfsdk:"enabled"`
	DefaultChannel types.Bool   `tfsdk:"default_channel"`
	Id             types.Int64  `tfsdk:"id"`
	CreatedAt      types.String `tfsdk:"created_at"`
	UpdatedAt      types.String `tfsdk:"updated_at"`
}

func (r *eventChannelResource) Configure(_ context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {
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

func (r *eventChannelResource) Metadata(_ context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_event_channel"
}

func (r *eventChannelResource) Schema(_ context.Context, _ resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Description: "An EventChannel is a named grouping of EventSubscriptions.",
		Attributes: map[string]schema.Attribute{
			"name": schema.StringAttribute{
				Description: "Event Channel name.",
				Required:    true,
			},
			"workspace_id": schema.Int64Attribute{
				Description: "Workspace ID. 0 means the default workspace.",
				Computed:    true,
				Optional:    true,
				PlanModifiers: []planmodifier.Int64{
					int64planmodifier.UseStateForUnknown(),
				},
			},
			"description": schema.StringAttribute{
				Description: "Event Channel description.",
				Computed:    true,
				Optional:    true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
			},
			"enabled": schema.BoolAttribute{
				Description: "Whether this Event Channel can dispatch events.",
				Computed:    true,
				Optional:    true,
				PlanModifiers: []planmodifier.Bool{
					boolplanmodifier.UseStateForUnknown(),
				},
			},
			"default_channel": schema.BoolAttribute{
				Description: "Whether this Event Channel is the default destination for newly published events.",
				Computed:    true,
				Optional:    true,
				PlanModifiers: []planmodifier.Bool{
					boolplanmodifier.UseStateForUnknown(),
				},
			},
			"id": schema.Int64Attribute{
				Description: "Event Channel ID",
				Computed:    true,
				PlanModifiers: []planmodifier.Int64{
					int64planmodifier.UseStateForUnknown(),
				},
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

func (r *eventChannelResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	var plan eventChannelResourceModel
	diags := req.Plan.Get(ctx, &plan)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}
	var config eventChannelResourceModel
	diags = req.Config.Get(ctx, &config)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	paramsEventChannelCreate := files_sdk.EventChannelCreateParams{}
	paramsEventChannelCreate.Name = plan.Name.ValueString()
	paramsEventChannelCreate.WorkspaceId = plan.WorkspaceId.ValueInt64()
	paramsEventChannelCreate.Description = plan.Description.ValueString()
	if !plan.Enabled.IsNull() && !plan.Enabled.IsUnknown() {
		paramsEventChannelCreate.Enabled = plan.Enabled.ValueBoolPointer()
	}
	if !plan.DefaultChannel.IsNull() && !plan.DefaultChannel.IsUnknown() {
		paramsEventChannelCreate.DefaultChannel = plan.DefaultChannel.ValueBoolPointer()
	}

	if resp.Diagnostics.HasError() {
		return
	}

	eventChannel, err := r.client.Create(paramsEventChannelCreate, files_sdk.WithContext(ctx))
	if err != nil {
		resp.Diagnostics.AddError(
			"Error Creating Files EventChannel",
			"Could not create event_channel, unexpected error: "+err.Error(),
		)
		return
	}

	diags = r.populateResourceModel(ctx, eventChannel, &plan)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	diags = resp.State.Set(ctx, plan)
	resp.Diagnostics.Append(diags...)
}

func (r *eventChannelResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	var state eventChannelResourceModel
	diags := req.State.Get(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	paramsEventChannelFind := files_sdk.EventChannelFindParams{}
	paramsEventChannelFind.Id = state.Id.ValueInt64()

	eventChannel, err := r.client.Find(paramsEventChannelFind, files_sdk.WithContext(ctx))
	if err != nil {
		if files_sdk.IsNotExist(err) {
			resp.State.RemoveResource(ctx)
			return
		}

		resp.Diagnostics.AddError(
			"Error Reading Files EventChannel",
			"Could not read event_channel id "+fmt.Sprint(state.Id.ValueInt64())+": "+err.Error(),
		)
		return
	}

	diags = r.populateResourceModel(ctx, eventChannel, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	diags = resp.State.Set(ctx, state)
	resp.Diagnostics.Append(diags...)
}

func (r *eventChannelResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	var plan eventChannelResourceModel
	diags := req.Plan.Get(ctx, &plan)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}
	var config eventChannelResourceModel
	diags = req.Config.Get(ctx, &config)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	paramsEventChannelUpdate := map[string]interface{}{}
	if !plan.Id.IsNull() && !plan.Id.IsUnknown() {
		paramsEventChannelUpdate["id"] = plan.Id.ValueInt64()
	}
	if !config.Name.IsNull() && !config.Name.IsUnknown() {
		paramsEventChannelUpdate["name"] = config.Name.ValueString()
	}
	if !config.WorkspaceId.IsNull() && !config.WorkspaceId.IsUnknown() {
		paramsEventChannelUpdate["workspace_id"] = config.WorkspaceId.ValueInt64()
	}
	if !config.Description.IsNull() && !config.Description.IsUnknown() {
		paramsEventChannelUpdate["description"] = config.Description.ValueString()
	}
	if !config.Enabled.IsNull() && !config.Enabled.IsUnknown() {
		paramsEventChannelUpdate["enabled"] = config.Enabled.ValueBool()
	}
	if !config.DefaultChannel.IsNull() && !config.DefaultChannel.IsUnknown() {
		paramsEventChannelUpdate["default_channel"] = config.DefaultChannel.ValueBool()
	}

	if resp.Diagnostics.HasError() {
		return
	}

	eventChannel, err := r.client.UpdateWithMap(paramsEventChannelUpdate, files_sdk.WithContext(ctx))
	if err != nil {
		resp.Diagnostics.AddError(
			"Error Updating Files EventChannel",
			"Could not update event_channel, unexpected error: "+err.Error(),
		)
		return
	}

	diags = r.populateResourceModel(ctx, eventChannel, &plan)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	diags = resp.State.Set(ctx, plan)
	resp.Diagnostics.Append(diags...)
}

func (r *eventChannelResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	var state eventChannelResourceModel
	diags := req.State.Get(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	paramsEventChannelDelete := files_sdk.EventChannelDeleteParams{}
	paramsEventChannelDelete.Id = state.Id.ValueInt64()

	err := r.client.Delete(paramsEventChannelDelete, files_sdk.WithContext(ctx))
	if err != nil && !files_sdk.IsNotExist(err) {
		resp.Diagnostics.AddError(
			"Error Deleting Files EventChannel",
			"Could not delete event_channel id "+fmt.Sprint(state.Id.ValueInt64())+": "+err.Error(),
		)
	}
}

func (r *eventChannelResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
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

func (r *eventChannelResource) populateResourceModel(ctx context.Context, eventChannel files_sdk.EventChannel, state *eventChannelResourceModel) (diags diag.Diagnostics) {
	state.Id = types.Int64Value(eventChannel.Id)
	state.Name = types.StringValue(eventChannel.Name)
	state.WorkspaceId = types.Int64Value(eventChannel.WorkspaceId)
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
