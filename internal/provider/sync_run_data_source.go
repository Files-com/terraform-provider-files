package provider

import (
	"context"
	"fmt"

	files_sdk "github.com/Files-com/files-sdk-go/v3"
	sync_run "github.com/Files-com/files-sdk-go/v3/syncrun"
	"github.com/Files-com/terraform-provider-files/lib"
	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/datasource/schema"
	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

var (
	_ datasource.DataSource              = &syncRunDataSource{}
	_ datasource.DataSourceWithConfigure = &syncRunDataSource{}
)

func NewSyncRunDataSource() datasource.DataSource {
	return &syncRunDataSource{}
}

type syncRunDataSource struct {
	client *sync_run.Client
}

type syncRunDataSourceModel struct {
	Id                   types.Int64  `tfsdk:"id"`
	Body                 types.String `tfsdk:"body"`
	BytesSynced          types.Int64  `tfsdk:"bytes_synced"`
	ComparedFiles        types.Int64  `tfsdk:"compared_files"`
	ComparedFolders      types.Int64  `tfsdk:"compared_folders"`
	CompletedAt          types.String `tfsdk:"completed_at"`
	CreatedAt            types.String `tfsdk:"created_at"`
	DestRemoteServerType types.String `tfsdk:"dest_remote_server_type"`
	DryRun               types.Bool   `tfsdk:"dry_run"`
	ErroredFiles         types.Int64  `tfsdk:"errored_files"`
	EstimatedBytesCount  types.Int64  `tfsdk:"estimated_bytes_count"`
	EventErrors          types.List   `tfsdk:"event_errors"`
	LogUrl               types.String `tfsdk:"log_url"`
	Runtime              types.String `tfsdk:"runtime"`
	SiteId               types.Int64  `tfsdk:"site_id"`
	WorkspaceId          types.Int64  `tfsdk:"workspace_id"`
	SrcRemoteServerType  types.String `tfsdk:"src_remote_server_type"`
	Status               types.String `tfsdk:"status"`
	SuccessfulFiles      types.Int64  `tfsdk:"successful_files"`
	SyncId               types.Int64  `tfsdk:"sync_id"`
	SyncName             types.String `tfsdk:"sync_name"`
	UpdatedAt            types.String `tfsdk:"updated_at"`
}

func (r *syncRunDataSource) Configure(_ context.Context, req datasource.ConfigureRequest, resp *datasource.ConfigureResponse) {
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

	r.client = &sync_run.Client{Config: sdk_config}
}

func (r *syncRunDataSource) Metadata(_ context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_sync_run"
}

func (r *syncRunDataSource) Schema(_ context.Context, _ datasource.SchemaRequest, resp *datasource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Description: "A SyncRun represents a single execution (run) of a Sync job.\n\nIt tracks status, statistics, logs, and timing for each sync operation.",
		Attributes: map[string]schema.Attribute{
			"id": schema.Int64Attribute{
				Description: "SyncRun ID",
				Required:    true,
			},
			"body": schema.StringAttribute{
				Description: "Log or summary body for this run",
				Computed:    true,
			},
			"bytes_synced": schema.Int64Attribute{
				Description: "Total bytes synced in this run",
				Computed:    true,
			},
			"compared_files": schema.Int64Attribute{
				Description: "Number of files compared",
				Computed:    true,
			},
			"compared_folders": schema.Int64Attribute{
				Description: "Number of folders compared",
				Computed:    true,
			},
			"completed_at": schema.StringAttribute{
				Description: "When this run was completed",
				Computed:    true,
			},
			"created_at": schema.StringAttribute{
				Description: "When this run was created",
				Computed:    true,
			},
			"dest_remote_server_type": schema.StringAttribute{
				Description: "Destination remote server type, if any",
				Computed:    true,
			},
			"dry_run": schema.BoolAttribute{
				Description: "Whether this run was a dry run (no actual changes made)",
				Computed:    true,
			},
			"errored_files": schema.Int64Attribute{
				Description: "Number of files that errored",
				Computed:    true,
			},
			"estimated_bytes_count": schema.Int64Attribute{
				Description: "Estimated bytes count for this run",
				Computed:    true,
			},
			"event_errors": schema.ListAttribute{
				Description: "Array of errors encountered during the run",
				Computed:    true,
				ElementType: types.StringType,
			},
			"log_url": schema.StringAttribute{
				Description: "Link to external log file.",
				Computed:    true,
			},
			"runtime": schema.StringAttribute{
				Description: "Total runtime in seconds",
				Computed:    true,
			},
			"site_id": schema.Int64Attribute{
				Description: "Site ID",
				Computed:    true,
			},
			"workspace_id": schema.Int64Attribute{
				Description: "Workspace ID",
				Computed:    true,
			},
			"src_remote_server_type": schema.StringAttribute{
				Description: "Source remote server type, if any",
				Computed:    true,
			},
			"status": schema.StringAttribute{
				Description: "Status of the sync run (success, failure, partial_failure, in_progress, skipped)",
				Computed:    true,
			},
			"successful_files": schema.Int64Attribute{
				Description: "Number of files successfully synced",
				Computed:    true,
			},
			"sync_id": schema.Int64Attribute{
				Description: "ID of the Sync this run belongs to",
				Computed:    true,
			},
			"sync_name": schema.StringAttribute{
				Description: "Name of the Sync this run belongs to",
				Computed:    true,
			},
			"updated_at": schema.StringAttribute{
				Description: "When this run was last updated",
				Computed:    true,
			},
		},
	}
}

func (r *syncRunDataSource) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
	var data syncRunDataSourceModel
	diags := req.Config.Get(ctx, &data)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	paramsSyncRunFind := files_sdk.SyncRunFindParams{}
	paramsSyncRunFind.Id = data.Id.ValueInt64()

	syncRun, err := r.client.Find(paramsSyncRunFind, files_sdk.WithContext(ctx))
	if err != nil {
		resp.Diagnostics.AddError(
			"Error Reading Files SyncRun",
			"Could not read sync_run id "+fmt.Sprint(data.Id.ValueInt64())+": "+err.Error(),
		)
		return
	}

	diags = r.populateDataSourceModel(ctx, syncRun, &data)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	diags = resp.State.Set(ctx, data)
	resp.Diagnostics.Append(diags...)
}

func (r *syncRunDataSource) populateDataSourceModel(ctx context.Context, syncRun files_sdk.SyncRun, state *syncRunDataSourceModel) (diags diag.Diagnostics) {
	var propDiags diag.Diagnostics

	state.Id = types.Int64Value(syncRun.Id)
	state.Body = types.StringValue(syncRun.Body)
	state.BytesSynced = types.Int64Value(syncRun.BytesSynced)
	state.ComparedFiles = types.Int64Value(syncRun.ComparedFiles)
	state.ComparedFolders = types.Int64Value(syncRun.ComparedFolders)
	if err := lib.TimeToStringType(ctx, path.Root("completed_at"), syncRun.CompletedAt, &state.CompletedAt); err != nil {
		diags.AddError(
			"Error Creating Files SyncRun",
			"Could not convert state completed_at to string: "+err.Error(),
		)
	}
	if err := lib.TimeToStringType(ctx, path.Root("created_at"), syncRun.CreatedAt, &state.CreatedAt); err != nil {
		diags.AddError(
			"Error Creating Files SyncRun",
			"Could not convert state created_at to string: "+err.Error(),
		)
	}
	state.DestRemoteServerType = types.StringValue(syncRun.DestRemoteServerType)
	state.DryRun = types.BoolPointerValue(syncRun.DryRun)
	state.ErroredFiles = types.Int64Value(syncRun.ErroredFiles)
	state.EstimatedBytesCount = types.Int64Value(syncRun.EstimatedBytesCount)
	state.EventErrors, propDiags = types.ListValueFrom(ctx, types.StringType, syncRun.EventErrors)
	diags.Append(propDiags...)
	state.LogUrl = types.StringValue(syncRun.LogUrl)
	state.Runtime = types.StringValue(syncRun.Runtime)
	state.SiteId = types.Int64Value(syncRun.SiteId)
	state.WorkspaceId = types.Int64Value(syncRun.WorkspaceId)
	state.SrcRemoteServerType = types.StringValue(syncRun.SrcRemoteServerType)
	state.Status = types.StringValue(syncRun.Status)
	state.SuccessfulFiles = types.Int64Value(syncRun.SuccessfulFiles)
	state.SyncId = types.Int64Value(syncRun.SyncId)
	state.SyncName = types.StringValue(syncRun.SyncName)
	if err := lib.TimeToStringType(ctx, path.Root("updated_at"), syncRun.UpdatedAt, &state.UpdatedAt); err != nil {
		diags.AddError(
			"Error Creating Files SyncRun",
			"Could not convert state updated_at to string: "+err.Error(),
		)
	}

	return
}
