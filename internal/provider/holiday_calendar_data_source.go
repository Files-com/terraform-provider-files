package provider

import (
	"context"
	"fmt"

	files_sdk "github.com/Files-com/files-sdk-go/v3"
	holiday_calendar "github.com/Files-com/files-sdk-go/v3/holidaycalendar"
	"github.com/Files-com/terraform-provider-files/lib"
	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/datasource/schema"
	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

var (
	_ datasource.DataSource              = &holidayCalendarDataSource{}
	_ datasource.DataSourceWithConfigure = &holidayCalendarDataSource{}
)

func NewHolidayCalendarDataSource() datasource.DataSource {
	return &holidayCalendarDataSource{}
}

type holidayCalendarDataSource struct {
	client *holiday_calendar.Client
}

type holidayCalendarDataSourceModel struct {
	Id         types.Int64   `tfsdk:"id"`
	Name       types.String  `tfsdk:"name"`
	Definition types.Dynamic `tfsdk:"definition"`
	CreatedAt  types.String  `tfsdk:"created_at"`
	UpdatedAt  types.String  `tfsdk:"updated_at"`
}

func (r *holidayCalendarDataSource) Configure(_ context.Context, req datasource.ConfigureRequest, resp *datasource.ConfigureResponse) {
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

	r.client = &holiday_calendar.Client{Config: sdk_config}
}

func (r *holidayCalendarDataSource) Metadata(_ context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_holiday_calendar"
}

func (r *holidayCalendarDataSource) Schema(_ context.Context, _ datasource.SchemaRequest, resp *datasource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Description: "A Holiday Calendar defines site-wide holiday dates and optional partial-day windows that scheduled resources skip.",
		Attributes: map[string]schema.Attribute{
			"id": schema.Int64Attribute{
				Description: "Holiday Calendar ID. Use `custom_<id>` as a scheduled resource's `holiday_region`.",
				Required:    true,
			},
			"name": schema.StringAttribute{
				Description: "Holiday Calendar name.",
				Computed:    true,
			},
			"definition": schema.DynamicAttribute{
				Description: "Holiday rules for the calendar. For more information, refer to the Holiday Calendars section of the Files.com documentation.",
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

func (r *holidayCalendarDataSource) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
	var data holidayCalendarDataSourceModel
	diags := req.Config.Get(ctx, &data)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	paramsHolidayCalendarFind := files_sdk.HolidayCalendarFindParams{}
	paramsHolidayCalendarFind.Id = data.Id.ValueInt64()

	holidayCalendar, err := r.client.Find(paramsHolidayCalendarFind, files_sdk.WithContext(ctx))
	if err != nil {
		resp.Diagnostics.AddError(
			"Error Reading Files HolidayCalendar",
			"Could not read holiday_calendar id "+fmt.Sprint(data.Id.ValueInt64())+": "+err.Error(),
		)
		return
	}

	diags = r.populateDataSourceModel(ctx, holidayCalendar, &data)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	diags = resp.State.Set(ctx, data)
	resp.Diagnostics.Append(diags...)
}

func (r *holidayCalendarDataSource) populateDataSourceModel(ctx context.Context, holidayCalendar files_sdk.HolidayCalendar, state *holidayCalendarDataSourceModel) (diags diag.Diagnostics) {
	var propDiags diag.Diagnostics

	state.Id = types.Int64Value(holidayCalendar.Id)
	state.Name = types.StringValue(holidayCalendar.Name)
	state.Definition, propDiags = lib.ToDynamic(ctx, path.Root("definition"), holidayCalendar.Definition, state.Definition.UnderlyingValue())
	diags.Append(propDiags...)
	if err := lib.TimeToStringType(ctx, path.Root("created_at"), holidayCalendar.CreatedAt, &state.CreatedAt); err != nil {
		diags.AddError(
			"Error Creating Files HolidayCalendar",
			"Could not convert state created_at to string: "+err.Error(),
		)
	}
	if err := lib.TimeToStringType(ctx, path.Root("updated_at"), holidayCalendar.UpdatedAt, &state.UpdatedAt); err != nil {
		diags.AddError(
			"Error Creating Files HolidayCalendar",
			"Could not convert state updated_at to string: "+err.Error(),
		)
	}

	return
}
