package provider

import (
	"context"
	"fmt"

	files_sdk "github.com/Files-com/files-sdk-go/v3"
	scheduled_export "github.com/Files-com/files-sdk-go/v3/scheduledexport"
	"github.com/Files-com/terraform-provider-files/lib"
	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/datasource/schema"
	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

var (
	_ datasource.DataSource              = &scheduledExportDataSource{}
	_ datasource.DataSourceWithConfigure = &scheduledExportDataSource{}
)

func NewScheduledExportDataSource() datasource.DataSource {
	return &scheduledExportDataSource{}
}

type scheduledExportDataSource struct {
	client *scheduled_export.Client
}

type scheduledExportDataSourceModel struct {
	Id                    types.Int64   `tfsdk:"id"`
	Name                  types.String  `tfsdk:"name"`
	ExportType            types.String  `tfsdk:"export_type"`
	ReportName            types.String  `tfsdk:"report_name"`
	ExportOptions         types.Dynamic `tfsdk:"export_options"`
	UserId                types.Int64   `tfsdk:"user_id"`
	Disabled              types.Bool    `tfsdk:"disabled"`
	Trigger               types.String  `tfsdk:"trigger"`
	Interval              types.String  `tfsdk:"interval"`
	RecurringDay          types.Int64   `tfsdk:"recurring_day"`
	ScheduleDaysOfWeek    types.List    `tfsdk:"schedule_days_of_week"`
	ScheduleTimesOfDay    types.List    `tfsdk:"schedule_times_of_day"`
	ScheduleTimeZone      types.String  `tfsdk:"schedule_time_zone"`
	HolidayRegion         types.String  `tfsdk:"holiday_region"`
	HumanReadableSchedule types.String  `tfsdk:"human_readable_schedule"`
	LastRunAt             types.String  `tfsdk:"last_run_at"`
	LastExportId          types.Int64   `tfsdk:"last_export_id"`
	CreatedAt             types.String  `tfsdk:"created_at"`
	UpdatedAt             types.String  `tfsdk:"updated_at"`
}

func (r *scheduledExportDataSource) Configure(_ context.Context, req datasource.ConfigureRequest, resp *datasource.ConfigureResponse) {
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

	r.client = &scheduled_export.Client{Config: sdk_config}
}

func (r *scheduledExportDataSource) Metadata(_ context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_scheduled_export"
}

func (r *scheduledExportDataSource) Schema(_ context.Context, _ datasource.SchemaRequest, resp *datasource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Description: "A Scheduled Export defines a recurring schedule for generating one of the built-in CSV exports and e-mailing it to a Site Admin recipient.",
		Attributes: map[string]schema.Attribute{
			"id": schema.Int64Attribute{
				Description: "Scheduled Export ID",
				Required:    true,
			},
			"name": schema.StringAttribute{
				Description: "Name for this scheduled export.",
				Computed:    true,
			},
			"export_type": schema.StringAttribute{
				Description: "Export report type. Valid values: folder_size_audit, group_membership_audit, permission_audit, share_link_audit",
				Computed:    true,
			},
			"report_name": schema.StringAttribute{
				Description: "Human-readable report name.",
				Computed:    true,
			},
			"export_options": schema.DynamicAttribute{
				Description: "Report-specific options. `permission_audit` supports `group_by` with `user` or `path`.",
				Computed:    true,
			},
			"user_id": schema.Int64Attribute{
				Description: "Site Admin user who receives the completed export e-mail.",
				Computed:    true,
			},
			"disabled": schema.BoolAttribute{
				Description: "If true, this scheduled export will not run.",
				Computed:    true,
			},
			"trigger": schema.StringAttribute{
				Description: "Schedule trigger type: `daily` or `custom_schedule`.",
				Computed:    true,
			},
			"interval": schema.StringAttribute{
				Description: "If trigger is `daily`, this specifies how often to run the scheduled export.",
				Computed:    true,
			},
			"recurring_day": schema.Int64Attribute{
				Description: "If trigger is `daily`, this selects the day number inside the chosen interval.",
				Computed:    true,
			},
			"schedule_days_of_week": schema.ListAttribute{
				Description: "If trigger is `custom_schedule`, the 0-based weekdays used by the schedule.",
				Computed:    true,
				ElementType: types.Int64Type,
			},
			"schedule_times_of_day": schema.ListAttribute{
				Description: "Times of day in HH:MM format for schedule-driven exports.",
				Computed:    true,
				ElementType: types.StringType,
			},
			"schedule_time_zone": schema.StringAttribute{
				Description: "Time zone used by the scheduled export.",
				Computed:    true,
			},
			"holiday_region": schema.StringAttribute{
				Description: "Optional holiday region used by schedule-driven exports.",
				Computed:    true,
			},
			"human_readable_schedule": schema.StringAttribute{
				Description: "Human-readable schedule description.",
				Computed:    true,
			},
			"last_run_at": schema.StringAttribute{
				Description: "Most recent scheduled run time.",
				Computed:    true,
			},
			"last_export_id": schema.Int64Attribute{
				Description: "Most recent Export ID created by this schedule.",
				Computed:    true,
			},
			"created_at": schema.StringAttribute{
				Description: "Creation time.",
				Computed:    true,
			},
			"updated_at": schema.StringAttribute{
				Description: "Last update time.",
				Computed:    true,
			},
		},
	}
}

func (r *scheduledExportDataSource) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
	var data scheduledExportDataSourceModel
	diags := req.Config.Get(ctx, &data)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	paramsScheduledExportFind := files_sdk.ScheduledExportFindParams{}
	paramsScheduledExportFind.Id = data.Id.ValueInt64()

	scheduledExport, err := r.client.Find(paramsScheduledExportFind, files_sdk.WithContext(ctx))
	if err != nil {
		resp.Diagnostics.AddError(
			"Error Reading Files ScheduledExport",
			"Could not read scheduled_export id "+fmt.Sprint(data.Id.ValueInt64())+": "+err.Error(),
		)
		return
	}

	diags = r.populateDataSourceModel(ctx, scheduledExport, &data)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	diags = resp.State.Set(ctx, data)
	resp.Diagnostics.Append(diags...)
}

func (r *scheduledExportDataSource) populateDataSourceModel(ctx context.Context, scheduledExport files_sdk.ScheduledExport, state *scheduledExportDataSourceModel) (diags diag.Diagnostics) {
	var propDiags diag.Diagnostics

	state.Id = types.Int64Value(scheduledExport.Id)
	state.Name = types.StringValue(scheduledExport.Name)
	state.ExportType = types.StringValue(scheduledExport.ExportType)
	state.ReportName = types.StringValue(scheduledExport.ReportName)
	state.ExportOptions, propDiags = lib.ToDynamic(ctx, path.Root("export_options"), scheduledExport.ExportOptions, state.ExportOptions.UnderlyingValue())
	diags.Append(propDiags...)
	state.UserId = types.Int64Value(scheduledExport.UserId)
	state.Disabled = types.BoolPointerValue(scheduledExport.Disabled)
	state.Trigger = types.StringValue(scheduledExport.Trigger)
	state.Interval = types.StringValue(scheduledExport.Interval)
	state.RecurringDay = types.Int64Value(scheduledExport.RecurringDay)
	state.ScheduleDaysOfWeek, propDiags = types.ListValueFrom(ctx, types.Int64Type, scheduledExport.ScheduleDaysOfWeek)
	diags.Append(propDiags...)
	state.ScheduleTimesOfDay, propDiags = types.ListValueFrom(ctx, types.StringType, scheduledExport.ScheduleTimesOfDay)
	diags.Append(propDiags...)
	state.ScheduleTimeZone = types.StringValue(scheduledExport.ScheduleTimeZone)
	state.HolidayRegion = types.StringValue(scheduledExport.HolidayRegion)
	state.HumanReadableSchedule = types.StringValue(scheduledExport.HumanReadableSchedule)
	if err := lib.TimeToStringType(ctx, path.Root("last_run_at"), scheduledExport.LastRunAt, &state.LastRunAt); err != nil {
		diags.AddError(
			"Error Creating Files ScheduledExport",
			"Could not convert state last_run_at to string: "+err.Error(),
		)
	}
	state.LastExportId = types.Int64Value(scheduledExport.LastExportId)
	if err := lib.TimeToStringType(ctx, path.Root("created_at"), scheduledExport.CreatedAt, &state.CreatedAt); err != nil {
		diags.AddError(
			"Error Creating Files ScheduledExport",
			"Could not convert state created_at to string: "+err.Error(),
		)
	}
	if err := lib.TimeToStringType(ctx, path.Root("updated_at"), scheduledExport.UpdatedAt, &state.UpdatedAt); err != nil {
		diags.AddError(
			"Error Creating Files ScheduledExport",
			"Could not convert state updated_at to string: "+err.Error(),
		)
	}

	return
}
