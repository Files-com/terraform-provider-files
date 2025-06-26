package provider

import (
	"context"
	"fmt"

	files_sdk "github.com/Files-com/files-sdk-go/v3"
	sync "github.com/Files-com/files-sdk-go/v3/sync"
	"github.com/Files-com/terraform-provider-files/lib"
	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/datasource/schema"
	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

var (
	_ datasource.DataSource              = &syncDataSource{}
	_ datasource.DataSourceWithConfigure = &syncDataSource{}
)

func NewSyncDataSource() datasource.DataSource {
	return &syncDataSource{}
}

type syncDataSource struct {
	client *sync.Client
}

type syncDataSourceModel struct {
	Id                  types.Int64  `tfsdk:"id"`
	Name                types.String `tfsdk:"name"`
	Description         types.String `tfsdk:"description"`
	SiteId              types.Int64  `tfsdk:"site_id"`
	UserId              types.Int64  `tfsdk:"user_id"`
	SrcPath             types.String `tfsdk:"src_path"`
	DestPath            types.String `tfsdk:"dest_path"`
	SrcRemoteServerId   types.Int64  `tfsdk:"src_remote_server_id"`
	DestRemoteServerId  types.Int64  `tfsdk:"dest_remote_server_id"`
	TwoWay              types.Bool   `tfsdk:"two_way"`
	KeepAfterCopy       types.Bool   `tfsdk:"keep_after_copy"`
	DeleteEmptyFolders  types.Bool   `tfsdk:"delete_empty_folders"`
	Disabled            types.Bool   `tfsdk:"disabled"`
	Trigger             types.String `tfsdk:"trigger"`
	TriggerFile         types.String `tfsdk:"trigger_file"`
	IncludePatterns     types.List   `tfsdk:"include_patterns"`
	ExcludePatterns     types.List   `tfsdk:"exclude_patterns"`
	CreatedAt           types.String `tfsdk:"created_at"`
	UpdatedAt           types.String `tfsdk:"updated_at"`
	SyncIntervalMinutes types.Int64  `tfsdk:"sync_interval_minutes"`
	Interval            types.String `tfsdk:"interval"`
	RecurringDay        types.Int64  `tfsdk:"recurring_day"`
	ScheduleDaysOfWeek  types.List   `tfsdk:"schedule_days_of_week"`
	ScheduleTimesOfDay  types.List   `tfsdk:"schedule_times_of_day"`
	ScheduleTimeZone    types.String `tfsdk:"schedule_time_zone"`
	HolidayRegion       types.String `tfsdk:"holiday_region"`
}

func (r *syncDataSource) Configure(_ context.Context, req datasource.ConfigureRequest, resp *datasource.ConfigureResponse) {
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

	r.client = &sync.Client{Config: sdk_config}
}

func (r *syncDataSource) Metadata(_ context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_sync"
}

func (r *syncDataSource) Schema(_ context.Context, _ datasource.SchemaRequest, resp *datasource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Description: "A Sync represents a file synchronization job between two locations (local-remote, remote-remote, local-child_site, etc). \n\nIt can be scheduled, run manually, or triggered by custom logic. \n\nSyncs track their runs, status, and configuration.",
		Attributes: map[string]schema.Attribute{
			"id": schema.Int64Attribute{
				Description: "Sync ID",
				Required:    true,
			},
			"name": schema.StringAttribute{
				Description: "Name for this sync job",
				Computed:    true,
			},
			"description": schema.StringAttribute{
				Description: "Description for this sync job",
				Computed:    true,
			},
			"site_id": schema.Int64Attribute{
				Description: "Site ID this sync belongs to",
				Computed:    true,
			},
			"user_id": schema.Int64Attribute{
				Description: "User who created or owns this sync",
				Computed:    true,
			},
			"src_path": schema.StringAttribute{
				Description: "Absolute source path for the sync",
				Computed:    true,
			},
			"dest_path": schema.StringAttribute{
				Description: "Absolute destination path for the sync",
				Computed:    true,
			},
			"src_remote_server_id": schema.Int64Attribute{
				Description: "Remote server ID for the source (if remote)",
				Computed:    true,
			},
			"dest_remote_server_id": schema.Int64Attribute{
				Description: "Remote server ID for the destination (if remote)",
				Computed:    true,
			},
			"two_way": schema.BoolAttribute{
				Description: "Is this a two-way sync?",
				Computed:    true,
			},
			"keep_after_copy": schema.BoolAttribute{
				Description: "Keep files after copying?",
				Computed:    true,
			},
			"delete_empty_folders": schema.BoolAttribute{
				Description: "Delete empty folders after sync?",
				Computed:    true,
			},
			"disabled": schema.BoolAttribute{
				Description: "Is this sync disabled?",
				Computed:    true,
			},
			"trigger": schema.StringAttribute{
				Description: "Trigger type: daily, custom_schedule, or manual",
				Computed:    true,
			},
			"trigger_file": schema.StringAttribute{
				Description: "Some MFT services request an empty file (known as a trigger file) to signal the sync is complete and they can begin further processing. If trigger_file is set, a zero-byte file will be sent at the end of the sync.",
				Computed:    true,
			},
			"include_patterns": schema.ListAttribute{
				Description: "Array of glob patterns to include",
				Computed:    true,
				ElementType: types.StringType,
			},
			"exclude_patterns": schema.ListAttribute{
				Description: "Array of glob patterns to exclude",
				Computed:    true,
				ElementType: types.StringType,
			},
			"created_at": schema.StringAttribute{
				Description: "When this sync was created",
				Computed:    true,
			},
			"updated_at": schema.StringAttribute{
				Description: "When this sync was last updated",
				Computed:    true,
			},
			"sync_interval_minutes": schema.Int64Attribute{
				Description: "Frequency in minutes between syncs. If set, this value must be greater than or equal to the `remote_sync_interval` value for the site's plan. If left blank, the plan's `remote_sync_interval` will be used. This setting is only used if `trigger` is empty.",
				Computed:    true,
			},
			"interval": schema.StringAttribute{
				Description: "If trigger is `daily`, this specifies how often to run this sync.  One of: `day`, `week`, `week_end`, `month`, `month_end`, `quarter`, `quarter_end`, `year`, `year_end`",
				Computed:    true,
			},
			"recurring_day": schema.Int64Attribute{
				Description: "If trigger type is `daily`, this specifies a day number to run in one of the supported intervals: `week`, `month`, `quarter`, `year`.",
				Computed:    true,
			},
			"schedule_days_of_week": schema.ListAttribute{
				Description: "If trigger is `custom_schedule`, Custom schedule description for when the sync should be run. 0-based days of the week. 0 is Sunday, 1 is Monday, etc.",
				Computed:    true,
				ElementType: types.Int64Type,
			},
			"schedule_times_of_day": schema.ListAttribute{
				Description: "If trigger is `custom_schedule`, Custom schedule description for when the sync should be run. Times of day in HH:MM format.",
				Computed:    true,
				ElementType: types.StringType,
			},
			"schedule_time_zone": schema.StringAttribute{
				Description: "If trigger is `custom_schedule`, Custom schedule Time Zone for when the sync should be run.",
				Computed:    true,
			},
			"holiday_region": schema.StringAttribute{
				Description: "If trigger is `custom_schedule`, the Automation will check if there is a formal, observed holiday for the region, and if so, it will not run.",
				Computed:    true,
			},
		},
	}
}

func (r *syncDataSource) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
	var data syncDataSourceModel
	diags := req.Config.Get(ctx, &data)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	paramsSyncFind := files_sdk.SyncFindParams{}
	paramsSyncFind.Id = data.Id.ValueInt64()

	sync, err := r.client.Find(paramsSyncFind, files_sdk.WithContext(ctx))
	if err != nil {
		resp.Diagnostics.AddError(
			"Error Reading Files Sync",
			"Could not read sync id "+fmt.Sprint(data.Id.ValueInt64())+": "+err.Error(),
		)
		return
	}

	diags = r.populateDataSourceModel(ctx, sync, &data)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	diags = resp.State.Set(ctx, data)
	resp.Diagnostics.Append(diags...)
}

func (r *syncDataSource) populateDataSourceModel(ctx context.Context, sync files_sdk.Sync, state *syncDataSourceModel) (diags diag.Diagnostics) {
	var propDiags diag.Diagnostics

	state.Id = types.Int64Value(sync.Id)
	state.Name = types.StringValue(sync.Name)
	state.Description = types.StringValue(sync.Description)
	state.SiteId = types.Int64Value(sync.SiteId)
	state.UserId = types.Int64Value(sync.UserId)
	state.SrcPath = types.StringValue(sync.SrcPath)
	state.DestPath = types.StringValue(sync.DestPath)
	state.SrcRemoteServerId = types.Int64Value(sync.SrcRemoteServerId)
	state.DestRemoteServerId = types.Int64Value(sync.DestRemoteServerId)
	state.TwoWay = types.BoolPointerValue(sync.TwoWay)
	state.KeepAfterCopy = types.BoolPointerValue(sync.KeepAfterCopy)
	state.DeleteEmptyFolders = types.BoolPointerValue(sync.DeleteEmptyFolders)
	state.Disabled = types.BoolPointerValue(sync.Disabled)
	state.Trigger = types.StringValue(sync.Trigger)
	state.TriggerFile = types.StringValue(sync.TriggerFile)
	state.IncludePatterns, propDiags = types.ListValueFrom(ctx, types.StringType, sync.IncludePatterns)
	diags.Append(propDiags...)
	state.ExcludePatterns, propDiags = types.ListValueFrom(ctx, types.StringType, sync.ExcludePatterns)
	diags.Append(propDiags...)
	if err := lib.TimeToStringType(ctx, path.Root("created_at"), sync.CreatedAt, &state.CreatedAt); err != nil {
		diags.AddError(
			"Error Creating Files Sync",
			"Could not convert state created_at to string: "+err.Error(),
		)
	}
	if err := lib.TimeToStringType(ctx, path.Root("updated_at"), sync.UpdatedAt, &state.UpdatedAt); err != nil {
		diags.AddError(
			"Error Creating Files Sync",
			"Could not convert state updated_at to string: "+err.Error(),
		)
	}
	state.SyncIntervalMinutes = types.Int64Value(sync.SyncIntervalMinutes)
	state.Interval = types.StringValue(sync.Interval)
	state.RecurringDay = types.Int64Value(sync.RecurringDay)
	state.ScheduleDaysOfWeek, propDiags = types.ListValueFrom(ctx, types.Int64Type, sync.ScheduleDaysOfWeek)
	diags.Append(propDiags...)
	state.ScheduleTimesOfDay, propDiags = types.ListValueFrom(ctx, types.StringType, sync.ScheduleTimesOfDay)
	diags.Append(propDiags...)
	state.ScheduleTimeZone = types.StringValue(sync.ScheduleTimeZone)
	state.HolidayRegion = types.StringValue(sync.HolidayRegion)

	return
}
