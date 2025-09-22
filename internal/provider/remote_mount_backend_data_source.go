package provider

import (
	"context"
	"fmt"

	files_sdk "github.com/Files-com/files-sdk-go/v3"
	remote_mount_backend "github.com/Files-com/files-sdk-go/v3/remotemountbackend"
	"github.com/Files-com/terraform-provider-files/lib"
	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/datasource/schema"
	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

var (
	_ datasource.DataSource              = &remoteMountBackendDataSource{}
	_ datasource.DataSourceWithConfigure = &remoteMountBackendDataSource{}
)

func NewRemoteMountBackendDataSource() datasource.DataSource {
	return &remoteMountBackendDataSource{}
}

type remoteMountBackendDataSource struct {
	client *remote_mount_backend.Client
}

type remoteMountBackendDataSourceModel struct {
	Id                    types.Int64   `tfsdk:"id"`
	CanaryFilePath        types.String  `tfsdk:"canary_file_path"`
	Enabled               types.Bool    `tfsdk:"enabled"`
	Fall                  types.Int64   `tfsdk:"fall"`
	HealthCheckEnabled    types.Bool    `tfsdk:"health_check_enabled"`
	HealthCheckResults    types.Dynamic `tfsdk:"health_check_results"`
	HealthCheckType       types.String  `tfsdk:"health_check_type"`
	Interval              types.Int64   `tfsdk:"interval"`
	MinFreeCpu            types.String  `tfsdk:"min_free_cpu"`
	MinFreeMem            types.String  `tfsdk:"min_free_mem"`
	Priority              types.Int64   `tfsdk:"priority"`
	RemotePath            types.String  `tfsdk:"remote_path"`
	RemoteServerId        types.Int64   `tfsdk:"remote_server_id"`
	RemoteServerMountId   types.Int64   `tfsdk:"remote_server_mount_id"`
	Rise                  types.Int64   `tfsdk:"rise"`
	Status                types.String  `tfsdk:"status"`
	UndergoingMaintenance types.Bool    `tfsdk:"undergoing_maintenance"`
}

func (r *remoteMountBackendDataSource) Configure(_ context.Context, req datasource.ConfigureRequest, resp *datasource.ConfigureResponse) {
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

func (r *remoteMountBackendDataSource) Metadata(_ context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_remote_mount_backend"
}

func (r *remoteMountBackendDataSource) Schema(_ context.Context, _ datasource.SchemaRequest, resp *datasource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Description: "A Remote Mount Backend is used to provide high availability for a Remote Server Mount Folder Behavior.",
		Attributes: map[string]schema.Attribute{
			"id": schema.Int64Attribute{
				Description: "Unique identifier for this backend.",
				Required:    true,
			},
			"canary_file_path": schema.StringAttribute{
				Description: "Path to the canary file used for health checks.",
				Computed:    true,
			},
			"enabled": schema.BoolAttribute{
				Description: "True if this backend is enabled.",
				Computed:    true,
			},
			"fall": schema.Int64Attribute{
				Description: "Number of consecutive failures before considering the backend unhealthy.",
				Computed:    true,
			},
			"health_check_enabled": schema.BoolAttribute{
				Description: "True if health checks are enabled for this backend.",
				Computed:    true,
			},
			"health_check_results": schema.DynamicAttribute{
				Description: "Array of recent health check results.",
				Computed:    true,
			},
			"health_check_type": schema.StringAttribute{
				Description: "Type of health check to perform.",
				Computed:    true,
			},
			"interval": schema.Int64Attribute{
				Description: "Interval in seconds between health checks.",
				Computed:    true,
			},
			"min_free_cpu": schema.StringAttribute{
				Description: "Minimum free CPU percentage required for this backend to be considered healthy.",
				Computed:    true,
			},
			"min_free_mem": schema.StringAttribute{
				Description: "Minimum free memory percentage required for this backend to be considered healthy.",
				Computed:    true,
			},
			"priority": schema.Int64Attribute{
				Description: "Priority of this backend.",
				Computed:    true,
			},
			"remote_path": schema.StringAttribute{
				Description: "Path on the remote server to treat as the root of this mount.",
				Computed:    true,
			},
			"remote_server_id": schema.Int64Attribute{
				Description: "The remote server that this backend is associated with.",
				Computed:    true,
			},
			"remote_server_mount_id": schema.Int64Attribute{
				Description: "The mount ID of the Remote Server Mount that this backend is associated with.",
				Computed:    true,
			},
			"rise": schema.Int64Attribute{
				Description: "Number of consecutive successes before considering the backend healthy.",
				Computed:    true,
			},
			"status": schema.StringAttribute{
				Description: "Status of this backend.",
				Computed:    true,
			},
			"undergoing_maintenance": schema.BoolAttribute{
				Description: "True if this backend is undergoing maintenance.",
				Computed:    true,
			},
		},
	}
}

func (r *remoteMountBackendDataSource) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
	var data remoteMountBackendDataSourceModel
	diags := req.Config.Get(ctx, &data)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	paramsRemoteMountBackendFind := files_sdk.RemoteMountBackendFindParams{}
	paramsRemoteMountBackendFind.Id = data.Id.ValueInt64()

	remoteMountBackend, err := r.client.Find(paramsRemoteMountBackendFind, files_sdk.WithContext(ctx))
	if err != nil {
		resp.Diagnostics.AddError(
			"Error Reading Files RemoteMountBackend",
			"Could not read remote_mount_backend id "+fmt.Sprint(data.Id.ValueInt64())+": "+err.Error(),
		)
		return
	}

	diags = r.populateDataSourceModel(ctx, remoteMountBackend, &data)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	diags = resp.State.Set(ctx, data)
	resp.Diagnostics.Append(diags...)
}

func (r *remoteMountBackendDataSource) populateDataSourceModel(ctx context.Context, remoteMountBackend files_sdk.RemoteMountBackend, state *remoteMountBackendDataSourceModel) (diags diag.Diagnostics) {
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
