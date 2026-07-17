package provider

import (
	"context"
	"fmt"
	"strconv"
	"strings"

	files_sdk "github.com/Files-com/files-sdk-go/v3"
	holiday_calendar "github.com/Files-com/files-sdk-go/v3/holidaycalendar"
	"github.com/Files-com/terraform-provider-files/lib"
	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/int64planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

var (
	_ resource.Resource                = &holidayCalendarResource{}
	_ resource.ResourceWithConfigure   = &holidayCalendarResource{}
	_ resource.ResourceWithImportState = &holidayCalendarResource{}
)

func NewHolidayCalendarResource() resource.Resource {
	return &holidayCalendarResource{}
}

type holidayCalendarResource struct {
	client *holiday_calendar.Client
}

type holidayCalendarResourceModel struct {
	Name       types.String  `tfsdk:"name"`
	Id         types.Int64   `tfsdk:"id"`
	Definition types.Dynamic `tfsdk:"definition"`
	CreatedAt  types.String  `tfsdk:"created_at"`
	UpdatedAt  types.String  `tfsdk:"updated_at"`
}

func (r *holidayCalendarResource) Configure(_ context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {
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

func (r *holidayCalendarResource) Metadata(_ context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_holiday_calendar"
}

func (r *holidayCalendarResource) Schema(_ context.Context, _ resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Description: "A Holiday Calendar defines site-wide holiday dates and optional partial-day windows that scheduled resources skip.",
		Attributes: map[string]schema.Attribute{
			"name": schema.StringAttribute{
				Description: "Holiday Calendar name.",
				Required:    true,
			},
			"id": schema.Int64Attribute{
				Description: "Holiday Calendar ID. Use `custom_<id>` as a scheduled resource's `holiday_region`.",
				Computed:    true,
				PlanModifiers: []planmodifier.Int64{
					int64planmodifier.UseStateForUnknown(),
				},
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

func (r *holidayCalendarResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	var plan holidayCalendarResourceModel
	diags := req.Plan.Get(ctx, &plan)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}
	var config holidayCalendarResourceModel
	diags = req.Config.Get(ctx, &config)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	paramsHolidayCalendarCreate := files_sdk.HolidayCalendarCreateParams{}
	paramsHolidayCalendarCreate.Name = plan.Name.ValueString()

	if resp.Diagnostics.HasError() {
		return
	}

	holidayCalendar, err := r.client.Create(paramsHolidayCalendarCreate, files_sdk.WithContext(ctx))
	if err != nil {
		resp.Diagnostics.AddError(
			"Error Creating Files HolidayCalendar",
			"Could not create holiday_calendar, unexpected error: "+err.Error(),
		)
		return
	}

	diags = r.populateResourceModel(ctx, holidayCalendar, &plan)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	diags = resp.State.Set(ctx, plan)
	resp.Diagnostics.Append(diags...)
}

func (r *holidayCalendarResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	var state holidayCalendarResourceModel
	diags := req.State.Get(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	paramsHolidayCalendarFind := files_sdk.HolidayCalendarFindParams{}
	paramsHolidayCalendarFind.Id = state.Id.ValueInt64()

	holidayCalendar, err := r.client.Find(paramsHolidayCalendarFind, files_sdk.WithContext(ctx))
	if err != nil {
		if files_sdk.IsNotExist(err) {
			resp.State.RemoveResource(ctx)
			return
		}

		resp.Diagnostics.AddError(
			"Error Reading Files HolidayCalendar",
			"Could not read holiday_calendar id "+fmt.Sprint(state.Id.ValueInt64())+": "+err.Error(),
		)
		return
	}

	diags = r.populateResourceModel(ctx, holidayCalendar, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	diags = resp.State.Set(ctx, state)
	resp.Diagnostics.Append(diags...)
}

func (r *holidayCalendarResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	var plan holidayCalendarResourceModel
	diags := req.Plan.Get(ctx, &plan)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}
	var config holidayCalendarResourceModel
	diags = req.Config.Get(ctx, &config)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	paramsHolidayCalendarUpdate := map[string]interface{}{}
	if !plan.Id.IsNull() && !plan.Id.IsUnknown() {
		paramsHolidayCalendarUpdate["id"] = plan.Id.ValueInt64()
	}
	if !config.Name.IsNull() && !config.Name.IsUnknown() {
		paramsHolidayCalendarUpdate["name"] = config.Name.ValueString()
	}

	if resp.Diagnostics.HasError() {
		return
	}

	holidayCalendar, err := r.client.UpdateWithMap(paramsHolidayCalendarUpdate, files_sdk.WithContext(ctx))
	if err != nil {
		resp.Diagnostics.AddError(
			"Error Updating Files HolidayCalendar",
			"Could not update holiday_calendar, unexpected error: "+err.Error(),
		)
		return
	}

	diags = r.populateResourceModel(ctx, holidayCalendar, &plan)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	diags = resp.State.Set(ctx, plan)
	resp.Diagnostics.Append(diags...)
}

func (r *holidayCalendarResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	var state holidayCalendarResourceModel
	diags := req.State.Get(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	paramsHolidayCalendarDelete := files_sdk.HolidayCalendarDeleteParams{}
	paramsHolidayCalendarDelete.Id = state.Id.ValueInt64()

	err := r.client.Delete(paramsHolidayCalendarDelete, files_sdk.WithContext(ctx))
	if err != nil && !files_sdk.IsNotExist(err) {
		resp.Diagnostics.AddError(
			"Error Deleting Files HolidayCalendar",
			"Could not delete holiday_calendar id "+fmt.Sprint(state.Id.ValueInt64())+": "+err.Error(),
		)
	}
}

func (r *holidayCalendarResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
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

func (r *holidayCalendarResource) populateResourceModel(ctx context.Context, holidayCalendar files_sdk.HolidayCalendar, state *holidayCalendarResourceModel) (diags diag.Diagnostics) {
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
