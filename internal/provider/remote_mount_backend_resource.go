package provider

import (
	"context"
	"fmt"
	"strconv"
	"strings"

	files_sdk "github.com/Files-com/files-sdk-go/v3"
	remote_mount_backend "github.com/Files-com/files-sdk-go/v3/remotemountbackend"
	"github.com/Files-com/terraform-provider-files/lib"
	"github.com/hashicorp/terraform-plugin-framework-validators/stringvalidator"
	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/boolplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/int64planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/schema/validator"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

var (
	_ resource.Resource                = &remoteMountBackendResource{}
	_ resource.ResourceWithConfigure   = &remoteMountBackendResource{}
	_ resource.ResourceWithImportState = &remoteMountBackendResource{}
)

func NewRemoteMountBackendResource() resource.Resource {
	return &remoteMountBackendResource{}
}

type remoteMountBackendResource struct {
	client *remote_mount_backend.Client
}

type remoteMountBackendResourceModel struct {
	CanaryFilePath        types.String  `tfsdk:"canary_file_path"`
	RemoteServerId        types.Int64   `tfsdk:"remote_server_id"`
	RemoteServerMountId   types.Int64   `tfsdk:"remote_server_mount_id"`
	Enabled               types.Bool    `tfsdk:"enabled"`
	Fall                  types.Int64   `tfsdk:"fall"`
	HealthCheckEnabled    types.Bool    `tfsdk:"health_check_enabled"`
	HealthCheckType       types.String  `tfsdk:"health_check_type"`
	Interval              types.Int64   `tfsdk:"interval"`
	MinFreeCpu            types.String  `tfsdk:"min_free_cpu"`
	MinFreeMem            types.String  `tfsdk:"min_free_mem"`
	Priority              types.Int64   `tfsdk:"priority"`
	RemotePath            types.String  `tfsdk:"remote_path"`
	Rise                  types.Int64   `tfsdk:"rise"`
	HealthCheckResults    types.Dynamic `tfsdk:"health_check_results"`
	Id                    types.Int64   `tfsdk:"id"`
	Status                types.String  `tfsdk:"status"`
	UndergoingMaintenance types.Bool    `tfsdk:"undergoing_maintenance"`
}

func (r *remoteMountBackendResource) Configure(_ context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {
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

	r.client = &remote_mount_backend.Client{Config: sdk_config}
}

func (r *remoteMountBackendResource) Metadata(_ context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_remote_mount_backend"
}

func (r *remoteMountBackendResource) Schema(_ context.Context, _ resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Description: "A Remote Mount Backend is used to provide high availability for a Remote Server Mount Folder Behavior.",
		Attributes: map[string]schema.Attribute{
			"canary_file_path": schema.StringAttribute{
				Description: "Path to the canary file used for health checks.",
				Required:    true,
			},
			"remote_server_id": schema.Int64Attribute{
				Description: "The remote server that this backend is associated with.",
				Required:    true,
			},
			"remote_server_mount_id": schema.Int64Attribute{
				Description: "The mount ID of the Remote Server Mount that this backend is associated with.",
				Required:    true,
				PlanModifiers: []planmodifier.Int64{
					int64planmodifier.RequiresReplace(),
				},
			},
			"enabled": schema.BoolAttribute{
				Description: "True if this backend is enabled.",
				Computed:    true,
				Optional:    true,
				PlanModifiers: []planmodifier.Bool{
					boolplanmodifier.UseStateForUnknown(),
				},
			},
			"fall": schema.Int64Attribute{
				Description: "Number of consecutive failures before considering the backend unhealthy.",
				Computed:    true,
				Optional:    true,
				PlanModifiers: []planmodifier.Int64{
					int64planmodifier.UseStateForUnknown(),
				},
			},
			"health_check_enabled": schema.BoolAttribute{
				Description: "True if health checks are enabled for this backend.",
				Computed:    true,
				Optional:    true,
				PlanModifiers: []planmodifier.Bool{
					boolplanmodifier.UseStateForUnknown(),
				},
			},
			"health_check_type": schema.StringAttribute{
				Description: "Type of health check to perform.",
				Computed:    true,
				Optional:    true,
				Validators: []validator.String{
					stringvalidator.OneOf("active", "passive"),
				},
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
			},
			"interval": schema.Int64Attribute{
				Description: "Interval in seconds between health checks.",
				Computed:    true,
				Optional:    true,
				PlanModifiers: []planmodifier.Int64{
					int64planmodifier.UseStateForUnknown(),
				},
			},
			"min_free_cpu": schema.StringAttribute{
				Description: "Minimum free CPU percentage required for this backend to be considered healthy.",
				Computed:    true,
				Optional:    true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
			},
			"min_free_mem": schema.StringAttribute{
				Description: "Minimum free memory percentage required for this backend to be considered healthy.",
				Computed:    true,
				Optional:    true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
			},
			"priority": schema.Int64Attribute{
				Description: "Priority of this backend.",
				Computed:    true,
				Optional:    true,
				PlanModifiers: []planmodifier.Int64{
					int64planmodifier.UseStateForUnknown(),
				},
			},
			"remote_path": schema.StringAttribute{
				Description: "Path on the remote server to treat as the root of this mount.",
				Computed:    true,
				Optional:    true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
			},
			"rise": schema.Int64Attribute{
				Description: "Number of consecutive successes before considering the backend healthy.",
				Computed:    true,
				Optional:    true,
				PlanModifiers: []planmodifier.Int64{
					int64planmodifier.UseStateForUnknown(),
				},
			},
			"health_check_results": schema.DynamicAttribute{
				Description: "Array of recent health check results.",
				Computed:    true,
			},
			"id": schema.Int64Attribute{
				Description: "Unique identifier for this backend.",
				Computed:    true,
				PlanModifiers: []planmodifier.Int64{
					int64planmodifier.UseStateForUnknown(),
				},
			},
			"status": schema.StringAttribute{
				Description: "Status of this backend.",
				Computed:    true,
				Validators: []validator.String{
					stringvalidator.OneOf("healthy", "degraded", "failed", "desynced"),
				},
			},
			"undergoing_maintenance": schema.BoolAttribute{
				Description: "True if this backend is undergoing maintenance.",
				Computed:    true,
			},
		},
	}
}

func (r *remoteMountBackendResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	var plan remoteMountBackendResourceModel
	diags := req.Plan.Get(ctx, &plan)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}
	var config remoteMountBackendResourceModel
	diags = req.Config.Get(ctx, &config)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	paramsRemoteMountBackendCreate := files_sdk.RemoteMountBackendCreateParams{}
	if !plan.Enabled.IsNull() && !plan.Enabled.IsUnknown() {
		paramsRemoteMountBackendCreate.Enabled = plan.Enabled.ValueBoolPointer()
	}
	paramsRemoteMountBackendCreate.Fall = plan.Fall.ValueInt64()
	if !plan.HealthCheckEnabled.IsNull() && !plan.HealthCheckEnabled.IsUnknown() {
		paramsRemoteMountBackendCreate.HealthCheckEnabled = plan.HealthCheckEnabled.ValueBoolPointer()
	}
	paramsRemoteMountBackendCreate.HealthCheckType = paramsRemoteMountBackendCreate.HealthCheckType.Enum()[plan.HealthCheckType.ValueString()]
	paramsRemoteMountBackendCreate.Interval = plan.Interval.ValueInt64()
	paramsRemoteMountBackendCreate.MinFreeCpu = plan.MinFreeCpu.ValueString()
	paramsRemoteMountBackendCreate.MinFreeMem = plan.MinFreeMem.ValueString()
	paramsRemoteMountBackendCreate.Priority = plan.Priority.ValueInt64()
	paramsRemoteMountBackendCreate.RemotePath = plan.RemotePath.ValueString()
	paramsRemoteMountBackendCreate.Rise = plan.Rise.ValueInt64()
	paramsRemoteMountBackendCreate.CanaryFilePath = plan.CanaryFilePath.ValueString()
	paramsRemoteMountBackendCreate.RemoteServerMountId = plan.RemoteServerMountId.ValueInt64()
	paramsRemoteMountBackendCreate.RemoteServerId = plan.RemoteServerId.ValueInt64()

	if resp.Diagnostics.HasError() {
		return
	}

	remoteMountBackend, err := r.client.Create(paramsRemoteMountBackendCreate, files_sdk.WithContext(ctx))
	if err != nil {
		resp.Diagnostics.AddError(
			"Error Creating Files RemoteMountBackend",
			"Could not create remote_mount_backend, unexpected error: "+err.Error(),
		)
		return
	}

	diags = r.populateResourceModel(ctx, remoteMountBackend, &plan)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	diags = resp.State.Set(ctx, plan)
	resp.Diagnostics.Append(diags...)
}

func (r *remoteMountBackendResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	var state remoteMountBackendResourceModel
	diags := req.State.Get(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	paramsRemoteMountBackendFind := files_sdk.RemoteMountBackendFindParams{}
	paramsRemoteMountBackendFind.Id = state.Id.ValueInt64()

	remoteMountBackend, err := r.client.Find(paramsRemoteMountBackendFind, files_sdk.WithContext(ctx))
	if err != nil {
		if files_sdk.IsNotExist(err) {
			resp.State.RemoveResource(ctx)
			return
		}

		resp.Diagnostics.AddError(
			"Error Reading Files RemoteMountBackend",
			"Could not read remote_mount_backend id "+fmt.Sprint(state.Id.ValueInt64())+": "+err.Error(),
		)
		return
	}

	diags = r.populateResourceModel(ctx, remoteMountBackend, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	diags = resp.State.Set(ctx, state)
	resp.Diagnostics.Append(diags...)
}

func (r *remoteMountBackendResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	var plan remoteMountBackendResourceModel
	diags := req.Plan.Get(ctx, &plan)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}
	var config remoteMountBackendResourceModel
	diags = req.Config.Get(ctx, &config)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	paramsRemoteMountBackendUpdate := files_sdk.RemoteMountBackendUpdateParams{}
	paramsRemoteMountBackendUpdate.Id = plan.Id.ValueInt64()
	if !plan.Enabled.IsNull() && !plan.Enabled.IsUnknown() {
		paramsRemoteMountBackendUpdate.Enabled = plan.Enabled.ValueBoolPointer()
	}
	paramsRemoteMountBackendUpdate.Fall = plan.Fall.ValueInt64()
	if !plan.HealthCheckEnabled.IsNull() && !plan.HealthCheckEnabled.IsUnknown() {
		paramsRemoteMountBackendUpdate.HealthCheckEnabled = plan.HealthCheckEnabled.ValueBoolPointer()
	}
	paramsRemoteMountBackendUpdate.HealthCheckType = paramsRemoteMountBackendUpdate.HealthCheckType.Enum()[plan.HealthCheckType.ValueString()]
	paramsRemoteMountBackendUpdate.Interval = plan.Interval.ValueInt64()
	paramsRemoteMountBackendUpdate.MinFreeCpu = plan.MinFreeCpu.ValueString()
	paramsRemoteMountBackendUpdate.MinFreeMem = plan.MinFreeMem.ValueString()
	paramsRemoteMountBackendUpdate.Priority = plan.Priority.ValueInt64()
	paramsRemoteMountBackendUpdate.RemotePath = plan.RemotePath.ValueString()
	paramsRemoteMountBackendUpdate.Rise = plan.Rise.ValueInt64()
	paramsRemoteMountBackendUpdate.CanaryFilePath = plan.CanaryFilePath.ValueString()
	paramsRemoteMountBackendUpdate.RemoteServerId = plan.RemoteServerId.ValueInt64()

	if resp.Diagnostics.HasError() {
		return
	}

	remoteMountBackend, err := r.client.Update(paramsRemoteMountBackendUpdate, files_sdk.WithContext(ctx))
	if err != nil {
		resp.Diagnostics.AddError(
			"Error Updating Files RemoteMountBackend",
			"Could not update remote_mount_backend, unexpected error: "+err.Error(),
		)
		return
	}

	diags = r.populateResourceModel(ctx, remoteMountBackend, &plan)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	diags = resp.State.Set(ctx, plan)
	resp.Diagnostics.Append(diags...)
}

func (r *remoteMountBackendResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	var state remoteMountBackendResourceModel
	diags := req.State.Get(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	paramsRemoteMountBackendDelete := files_sdk.RemoteMountBackendDeleteParams{}
	paramsRemoteMountBackendDelete.Id = state.Id.ValueInt64()

	err := r.client.Delete(paramsRemoteMountBackendDelete, files_sdk.WithContext(ctx))
	if err != nil && !files_sdk.IsNotExist(err) {
		resp.Diagnostics.AddError(
			"Error Deleting Files RemoteMountBackend",
			"Could not delete remote_mount_backend id "+fmt.Sprint(state.Id.ValueInt64())+": "+err.Error(),
		)
	}
}

func (r *remoteMountBackendResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
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

func (r *remoteMountBackendResource) populateResourceModel(ctx context.Context, remoteMountBackend files_sdk.RemoteMountBackend, state *remoteMountBackendResourceModel) (diags diag.Diagnostics) {
	var propDiags diag.Diagnostics

	state.CanaryFilePath = types.StringValue(remoteMountBackend.CanaryFilePath)
	state.Enabled = types.BoolPointerValue(remoteMountBackend.Enabled)
	state.Fall = types.Int64Value(remoteMountBackend.Fall)
	state.HealthCheckEnabled = types.BoolPointerValue(remoteMountBackend.HealthCheckEnabled)
	state.HealthCheckResults, propDiags = lib.ToDynamic(ctx, path.Root("health_check_results"), remoteMountBackend.HealthCheckResults, state.HealthCheckResults.UnderlyingValue())
	diags.Append(propDiags...)
	state.HealthCheckType = types.StringValue(remoteMountBackend.HealthCheckType)
	state.Id = types.Int64Value(remoteMountBackend.Id)
	state.Interval = types.Int64Value(remoteMountBackend.Interval)
	state.MinFreeCpu = types.StringValue(remoteMountBackend.MinFreeCpu)
	state.MinFreeMem = types.StringValue(remoteMountBackend.MinFreeMem)
	state.Priority = types.Int64Value(remoteMountBackend.Priority)
	state.RemotePath = types.StringValue(remoteMountBackend.RemotePath)
	state.RemoteServerId = types.Int64Value(remoteMountBackend.RemoteServerId)
	state.RemoteServerMountId = types.Int64Value(remoteMountBackend.RemoteServerMountId)
	state.Rise = types.Int64Value(remoteMountBackend.Rise)
	state.Status = types.StringValue(remoteMountBackend.Status)
	state.UndergoingMaintenance = types.BoolPointerValue(remoteMountBackend.UndergoingMaintenance)

	return
}
