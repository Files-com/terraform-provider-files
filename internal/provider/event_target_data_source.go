package provider

import (
	"context"
	"fmt"

	files_sdk "github.com/Files-com/files-sdk-go/v3"
	event_target "github.com/Files-com/files-sdk-go/v3/eventtarget"
	"github.com/Files-com/terraform-provider-files/lib"
	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/datasource/schema"
	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

var (
	_ datasource.DataSource              = &eventTargetDataSource{}
	_ datasource.DataSourceWithConfigure = &eventTargetDataSource{}
)

func NewEventTargetDataSource() datasource.DataSource {
	return &eventTargetDataSource{}
}

type eventTargetDataSource struct {
	client *event_target.Client
}

type eventTargetDataSourceModel struct {
	Id                   types.Int64   `tfsdk:"id"`
	Name                 types.String  `tfsdk:"name"`
	TargetType           types.String  `tfsdk:"target_type"`
	WorkspaceId          types.Int64   `tfsdk:"workspace_id"`
	ApplyToAllWorkspaces types.Bool    `tfsdk:"apply_to_all_workspaces"`
	Enabled              types.Bool    `tfsdk:"enabled"`
	Config               types.Dynamic `tfsdk:"config"`
	DeliveryPolicy       types.Dynamic `tfsdk:"delivery_policy"`
	CreatedAt            types.String  `tfsdk:"created_at"`
	UpdatedAt            types.String  `tfsdk:"updated_at"`
}

func (r *eventTargetDataSource) Configure(_ context.Context, req datasource.ConfigureRequest, resp *datasource.ConfigureResponse) {
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

func (r *eventTargetDataSource) Metadata(_ context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_event_target"
}

func (r *eventTargetDataSource) Schema(_ context.Context, _ datasource.SchemaRequest, resp *datasource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Description: "An EventTarget is a delivery destination for EventRecords.",
		Attributes: map[string]schema.Attribute{
			"id": schema.Int64Attribute{
				Description: "Event Target ID",
				Required:    true,
			},
			"name": schema.StringAttribute{
				Description: "Event Target name.",
				Computed:    true,
			},
			"target_type": schema.StringAttribute{
				Description: "Event Target type.",
				Computed:    true,
			},
			"workspace_id": schema.Int64Attribute{
				Description: "Workspace ID. 0 means the default workspace or site-wide.",
				Computed:    true,
			},
			"apply_to_all_workspaces": schema.BoolAttribute{
				Description: "If true, this default-workspace target can receive events from all workspaces.",
				Computed:    true,
			},
			"enabled": schema.BoolAttribute{
				Description: "Whether this Event Target can receive events.",
				Computed:    true,
			},
			"config": schema.DynamicAttribute{
				Description: "Event Target configuration.",
				Computed:    true,
			},
			"delivery_policy": schema.DynamicAttribute{
				Description: "Event Target delivery policy. Email targets support batch_interval in seconds, between 600 and 86400.",
				Computed:    true,
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

func (r *eventTargetDataSource) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
	var data eventTargetDataSourceModel
	diags := req.Config.Get(ctx, &data)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	paramsEventTargetFind := files_sdk.EventTargetFindParams{}
	paramsEventTargetFind.Id = data.Id.ValueInt64()

	eventTarget, err := r.client.Find(paramsEventTargetFind, files_sdk.WithContext(ctx))
	if err != nil {
		resp.Diagnostics.AddError(
			"Error Reading Files EventTarget",
			"Could not read event_target id "+fmt.Sprint(data.Id.ValueInt64())+": "+err.Error(),
		)
		return
	}

	diags = r.populateDataSourceModel(ctx, eventTarget, &data)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	diags = resp.State.Set(ctx, data)
	resp.Diagnostics.Append(diags...)
}

func (r *eventTargetDataSource) populateDataSourceModel(ctx context.Context, eventTarget files_sdk.EventTarget, state *eventTargetDataSourceModel) (diags diag.Diagnostics) {
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
