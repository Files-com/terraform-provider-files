package provider

import (
	"context"

	"github.com/Files-com/terraform-provider-files/lib"
	"github.com/hashicorp/terraform-plugin-framework-validators/stringvalidator"
	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/boolplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/dynamicplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/int64planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/listplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/schema/validator"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

type automationResourceModelV0 struct {
	Automation                       types.String  `tfsdk:"automation"`
	WorkspaceId                      types.Int64   `tfsdk:"workspace_id"`
	AlwaysSerializeJobs              types.Bool    `tfsdk:"always_serialize_jobs"`
	AlwaysOverwriteSizeMatchingFiles types.Bool    `tfsdk:"always_overwrite_size_matching_files"`
	Description                      types.String  `tfsdk:"description"`
	Definition                       types.Dynamic `tfsdk:"definition"`
	DestinationReplaceFrom           types.String  `tfsdk:"destination_replace_from"`
	DestinationReplaceTo             types.String  `tfsdk:"destination_replace_to"`
	Destinations                     types.List    `tfsdk:"destinations"`
	Disabled                         types.Bool    `tfsdk:"disabled"`
	ExcludePattern                   types.String  `tfsdk:"exclude_pattern"`
	ImportUrls                       types.Dynamic `tfsdk:"import_urls"`
	FlattenDestinationStructure      types.Bool    `tfsdk:"flatten_destination_structure"`
	GroupIds                         types.List    `tfsdk:"group_ids"`
	IgnoreLockedFolders              types.Bool    `tfsdk:"ignore_locked_folders"`
	Interval                         types.String  `tfsdk:"interval"`
	LegacyFolderMatching             types.Bool    `tfsdk:"legacy_folder_matching"`
	Name                             types.String  `tfsdk:"name"`
	OverwriteFiles                   types.Bool    `tfsdk:"overwrite_files"`
	Path                             types.String  `tfsdk:"path"`
	PathTimeZone                     types.String  `tfsdk:"path_time_zone"`
	RecurringDay                     types.Int64   `tfsdk:"recurring_day"`
	RecurringDays                    types.List    `tfsdk:"recurring_days"`
	ScheduleId                       types.Int64   `tfsdk:"schedule_id"`
	RetryOnFailureIntervalInMinutes  types.Int64   `tfsdk:"retry_on_failure_interval_in_minutes"`
	RetryOnFailureNumberOfAttempts   types.Int64   `tfsdk:"retry_on_failure_number_of_attempts"`
	ScheduleDaysOfWeek               types.List    `tfsdk:"schedule_days_of_week"`
	ScheduleTimesOfDay               types.List    `tfsdk:"schedule_times_of_day"`
	ScheduleTimeZone                 types.String  `tfsdk:"schedule_time_zone"`
	Source                           types.String  `tfsdk:"source"`
	SyncIds                          types.List    `tfsdk:"sync_ids"`
	TriggerActions                   types.List    `tfsdk:"trigger_actions"`
	Trigger                          types.String  `tfsdk:"trigger"`
	UserIds                          types.List    `tfsdk:"user_ids"`
	Value                            types.Dynamic `tfsdk:"value"`
	HolidayRegion                    types.String  `tfsdk:"holiday_region"`
	Id                               types.Int64   `tfsdk:"id"`
	Deleted                          types.Bool    `tfsdk:"deleted"`
	InboundEmailAddress              types.String  `tfsdk:"inbound_email_address"`
	LastModifiedAt                   types.String  `tfsdk:"last_modified_at"`
	Version                          types.Int64   `tfsdk:"version"`
	Schedule                         types.Dynamic `tfsdk:"schedule"`
	HumanReadableSchedule            types.String  `tfsdk:"human_readable_schedule"`
	UserId                           types.Int64   `tfsdk:"user_id"`
	WebhookUrl                       types.String  `tfsdk:"webhook_url"`
}

func (r *automationResource) UpgradeState(ctx context.Context) map[int64]resource.StateUpgrader {
	return map[int64]resource.StateUpgrader{
		0: {
			PriorSchema: &schema.Schema{
				Attributes: map[string]schema.Attribute{

					"automation": schema.StringAttribute{
						Description: "Automation type",
						Required:    true,
						Validators: []validator.String{
							stringvalidator.OneOf("create_folder", "delete_file", "copy_file", "move_file", "as2_send", "run_sync", "import_file", "v2"),
						},
					},
					"workspace_id": schema.Int64Attribute{
						Description: "Workspace ID",
						Computed:    true,
						Optional:    true,
						PlanModifiers: []planmodifier.Int64{
							int64planmodifier.UseStateForUnknown(),
							int64planmodifier.RequiresReplace(),
						},
					},
					"always_serialize_jobs": schema.BoolAttribute{
						Description: "Ordinarily, we will allow automation runs to run in parallel for non-scheduled automations. If this flag is `true` we will force automation runs to be serialized (run one at a time, one after another). This can resolve some issues with race conditions on remote systems at the cost of some performance.",
						Computed:    true,
						Optional:    true,
						PlanModifiers: []planmodifier.Bool{
							boolplanmodifier.UseStateForUnknown(),
						},
					},
					"always_overwrite_size_matching_files": schema.BoolAttribute{
						Description: "Ordinarily, files with identical size in the source and destination will be skipped from copy operations to prevent wasted transfer.  If this flag is `true` we will overwrite the destination file always.  Note that this may cause large amounts of wasted transfer usage.  This setting has no effect unless `overwrite_files` is also set to `true`.",
						Computed:    true,
						Optional:    true,
						PlanModifiers: []planmodifier.Bool{
							boolplanmodifier.UseStateForUnknown(),
						},
					},
					"description": schema.StringAttribute{
						Description: "Description for the this Automation.",
						Computed:    true,
						Optional:    true,
						PlanModifiers: []planmodifier.String{
							stringplanmodifier.UseStateForUnknown(),
						},
					},
					"definition": schema.DynamicAttribute{
						Description: "Automation v2 graph definition.",
						Computed:    true,
					},
					"destination_replace_from": schema.StringAttribute{
						Description: "If set, this string in the destination path will be replaced with the value in `destination_replace_to`.",
						Computed:    true,
						Optional:    true,
						PlanModifiers: []planmodifier.String{
							stringplanmodifier.UseStateForUnknown(),
						},
					},
					"destination_replace_to": schema.StringAttribute{
						Description: "If set, this string will replace the value `destination_replace_from` in the destination filename. You can use special patterns here.",
						Computed:    true,
						Optional:    true,
						PlanModifiers: []planmodifier.String{
							stringplanmodifier.UseStateForUnknown(),
						},
					},
					"destinations": schema.ListAttribute{
						Description: "Destination Paths",
						Computed:    true,
						Optional:    true,
						ElementType: types.StringType,
						PlanModifiers: []planmodifier.List{
							listplanmodifier.UseStateForUnknown(),
						},
					},
					"disabled": schema.BoolAttribute{
						Description: "If true, this automation will not run.",
						Computed:    true,
						Optional:    true,
						PlanModifiers: []planmodifier.Bool{
							boolplanmodifier.UseStateForUnknown(),
						},
					},
					"exclude_pattern": schema.StringAttribute{
						Description: "If set, this glob pattern will exclude files from the automation. Supports globs, except on remote mounts.",
						Computed:    true,
						Optional:    true,
						PlanModifiers: []planmodifier.String{
							stringplanmodifier.UseStateForUnknown(),
						},
					},
					"import_urls": schema.DynamicAttribute{
						Description: "List of URLs to be imported and names to be used.",
						Computed:    true,
						Optional:    true,
						PlanModifiers: []planmodifier.Dynamic{
							dynamicplanmodifier.UseStateForUnknown(),
						},
					},
					"flatten_destination_structure": schema.BoolAttribute{
						Description: "Normally copy and move automations that use globs will implicitly preserve the source folder structure in the destination.  If this flag is `true`, the source folder structure will be flattened in the destination.  This is useful for copying or moving files from multiple folders into a single destination folder.",
						Computed:    true,
						Optional:    true,
						PlanModifiers: []planmodifier.Bool{
							boolplanmodifier.UseStateForUnknown(),
						},
					},
					"group_ids": schema.ListAttribute{
						Description: "IDs of Groups for the Automation (i.e. who to Request File from)",
						Computed:    true,
						Optional:    true,
						ElementType: types.Int64Type,
						PlanModifiers: []planmodifier.List{
							listplanmodifier.UseStateForUnknown(),
						},
					},
					"ignore_locked_folders": schema.BoolAttribute{
						Description: "If true, the Lock Folders behavior will be disregarded for automated actions.",
						Computed:    true,
						Optional:    true,
						PlanModifiers: []planmodifier.Bool{
							boolplanmodifier.UseStateForUnknown(),
						},
					},
					"interval": schema.StringAttribute{
						Description: "If trigger is `daily`, this specifies how often to run this automation.  One of: `day`, `week`, `week_end`, `month`, `month_end`, `quarter`, `quarter_end`, `year`, `year_end`",
						Computed:    true,
						Optional:    true,
						PlanModifiers: []planmodifier.String{
							stringplanmodifier.UseStateForUnknown(),
						},
					},
					"legacy_folder_matching": schema.BoolAttribute{
						Description: "If `true`, use the legacy behavior for this automation, where it can operate on folders in addition to just files.  This behavior no longer works and should not be used.",
						Computed:    true,
						Optional:    true,
						PlanModifiers: []planmodifier.Bool{
							boolplanmodifier.UseStateForUnknown(),
						},
					},
					"name": schema.StringAttribute{
						Description: "Name for this automation.",
						Computed:    true,
						Optional:    true,
						PlanModifiers: []planmodifier.String{
							stringplanmodifier.UseStateForUnknown(),
						},
					},
					"overwrite_files": schema.BoolAttribute{
						Description: "If true, existing files will be overwritten with new files on Move/Copy automations.  Note: by default files will not be overwritten on Copy automations if they appear to be the same file size as the newly incoming file.  Use the `always_overwrite_size_matching_files` option in conjunction with `overwrite_files` to override this behavior and overwrite files no matter what.",
						Computed:    true,
						Optional:    true,
						PlanModifiers: []planmodifier.Bool{
							boolplanmodifier.UseStateForUnknown(),
						},
					},
					"path": schema.StringAttribute{
						Description: "Path on which this Automation runs.  Supports globs, except on remote mounts. This must be slash-delimited, but it must neither start nor end with a slash. Maximum of 5000 characters.",
						Computed:    true,
						Optional:    true,
						PlanModifiers: []planmodifier.String{
							stringplanmodifier.UseStateForUnknown(),
						},
					},
					"path_time_zone": schema.StringAttribute{
						Description: "Timezone to use when rendering timestamps in paths.",
						Computed:    true,
						Optional:    true,
						PlanModifiers: []planmodifier.String{
							stringplanmodifier.UseStateForUnknown(),
						},
					},
					"recurring_day": schema.Int64Attribute{
						Description: "If trigger type is `daily`, this specifies a day number to run in one of the supported intervals: `week`, `month`, `quarter`, `year`.",
						Computed:    true,
						Optional:    true,
						PlanModifiers: []planmodifier.Int64{
							int64planmodifier.UseStateForUnknown(),
						},
					},
					"recurring_days": schema.ListAttribute{
						Description: "If trigger type is `daily`, this specifies one or more day numbers to run in one of the supported intervals: `week`, `month`, `quarter`, `year`.",
						Computed:    true,
						Optional:    true,
						ElementType: types.Int64Type,
						PlanModifiers: []planmodifier.List{
							listplanmodifier.UseStateForUnknown(),
						},
					},
					"schedule_id": schema.Int64Attribute{
						Description: "If trigger is `custom_schedule`, the reusable Schedule used instead of the automation's schedule fields.",
						Computed:    true,
						Optional:    true,
						PlanModifiers: []planmodifier.Int64{
							int64planmodifier.UseStateForUnknown(),
						},
					},
					"retry_on_failure_interval_in_minutes": schema.Int64Attribute{
						Description: "If the Automation fails, retry at this interval (in minutes).  Acceptable values are 5 through 1440 (one day).  Set to null to disable.",
						Computed:    true,
						Optional:    true,
						PlanModifiers: []planmodifier.Int64{
							int64planmodifier.UseStateForUnknown(),
						},
					},
					"retry_on_failure_number_of_attempts": schema.Int64Attribute{
						Description: "If the Automation fails, retry at most this many times.  Maximum allowed value: 10.  Set to null to disable.",
						Computed:    true,
						Optional:    true,
						PlanModifiers: []planmodifier.Int64{
							int64planmodifier.UseStateForUnknown(),
						},
					},
					"schedule_days_of_week": schema.ListAttribute{
						Description: "If trigger is `custom_schedule`, Custom schedule description for when the automation should be run. 0 is Sunday, 1 is Monday, etc.",
						Computed:    true,
						Optional:    true,
						ElementType: types.Int64Type,
						PlanModifiers: []planmodifier.List{
							listplanmodifier.UseStateForUnknown(),
						},
					},
					"schedule_times_of_day": schema.ListAttribute{
						Description: "Times of day to run in HH:MM format (24-hour). For `custom_schedule`, run at these times on specified days of week. For `daily`, run at these times on the scheduled interval date.",
						Computed:    true,
						Optional:    true,
						ElementType: types.StringType,
						PlanModifiers: []planmodifier.List{
							listplanmodifier.UseStateForUnknown(),
						},
					},
					"schedule_time_zone": schema.StringAttribute{
						Description: "Time zone for the schedule. If not set, times are interpreted as UTC.",
						Computed:    true,
						Optional:    true,
						PlanModifiers: []planmodifier.String{
							stringplanmodifier.UseStateForUnknown(),
						},
					},
					"source": schema.StringAttribute{
						Description: "Source path/glob.  See Automation docs for exact description, but this is used to filter for files in the `path` to find files to operate on. Supports globs, except on remote mounts.",
						Computed:    true,
						Optional:    true,
						PlanModifiers: []planmodifier.String{
							stringplanmodifier.UseStateForUnknown(),
						},
					},
					"sync_ids": schema.ListAttribute{
						Description: "IDs of Syncs to run by this Automation.",
						Computed:    true,
						Optional:    true,
						ElementType: types.Int64Type,
						PlanModifiers: []planmodifier.List{
							listplanmodifier.UseStateForUnknown(),
						},
					},
					"trigger_actions": schema.ListAttribute{
						Description: "If trigger is `action`, this is the list of action types on which to trigger the automation. Valid actions are create, copy, move, archived_delete, update, read, destroy",
						Computed:    true,
						Optional:    true,
						ElementType: types.StringType,
						PlanModifiers: []planmodifier.List{
							listplanmodifier.UseStateForUnknown(),
						},
					},
					"trigger": schema.StringAttribute{
						Description: "How this automation is triggered to run.",
						Computed:    true,
						Optional:    true,
						Validators: []validator.String{
							stringvalidator.OneOf("manual", "daily", "custom_schedule", "webhook", "email", "action"),
						},
						PlanModifiers: []planmodifier.String{
							stringplanmodifier.UseStateForUnknown(),
						},
					},
					"user_ids": schema.ListAttribute{
						Description: "IDs of Users for the Automation (i.e. who to Request File from)",
						Computed:    true,
						Optional:    true,
						ElementType: types.Int64Type,
						PlanModifiers: []planmodifier.List{
							listplanmodifier.UseStateForUnknown(),
						},
					},
					"value": schema.DynamicAttribute{
						Description: "A Hash of attributes specific to the automation type.",
						Computed:    true,
						Optional:    true,
						Validators: []validator.Dynamic{
							lib.DeprecatedJSONEncoding("files_automation.value", "value = {\n  limit = \"1\"\n}", "March 1, 2027"),
						},
						PlanModifiers: []planmodifier.Dynamic{
							dynamicplanmodifier.UseStateForUnknown(),
						},
					},
					"holiday_region": schema.StringAttribute{
						Description: "Skip the automation if there is a formal, observed holiday for this region.",
						Computed:    true,
						Optional:    true,
						PlanModifiers: []planmodifier.String{
							stringplanmodifier.UseStateForUnknown(),
						},
					},
					"id": schema.Int64Attribute{
						Description: "Automation ID",
						Computed:    true,
						PlanModifiers: []planmodifier.Int64{
							int64planmodifier.UseStateForUnknown(),
						},
					},
					"deleted": schema.BoolAttribute{
						Description: "Indicates if the automation has been deleted.",
						Computed:    true,
					},
					"inbound_email_address": schema.StringAttribute{
						Description: "If trigger is `email`, this is the address that triggers the Automation.",
						Computed:    true,
					},
					"last_modified_at": schema.StringAttribute{
						Description: "Time when automation was last modified. Does not change for name or description updates.",
						Computed:    true,
					},
					"version": schema.Int64Attribute{
						Description: "Current Automation v2 definition version.",
						Computed:    true,
					},
					"schedule": schema.DynamicAttribute{
						Description: "If trigger is `custom_schedule`, Custom schedule description for when the automation should be run in json format.",
						Computed:    true,
					},
					"human_readable_schedule": schema.StringAttribute{
						Description: "If trigger is `custom_schedule` or `daily` with times, Human readable schedule description for when the automation should be run.",
						Computed:    true,
					},
					"user_id": schema.Int64Attribute{
						Description: "User ID of the Automation's creator.",
						Computed:    true,
					},
					"webhook_url": schema.StringAttribute{
						Description: "If trigger is `webhook`, this is the URL of the webhook to trigger the Automation.",
						Computed:    true,
					},
				},
				Version: 0,
			},
			StateUpgrader: func(ctx context.Context, req resource.UpgradeStateRequest, resp *resource.UpgradeStateResponse) {
				var priorState automationResourceModelV0
				resp.Diagnostics.Append(req.State.Get(ctx, &priorState)...)
				if resp.Diagnostics.HasError() {
					return
				}
				upgradedState := automationResourceModel{

					Automation:                       priorState.Automation,
					WorkspaceId:                      priorState.WorkspaceId,
					AlwaysSerializeJobs:              priorState.AlwaysSerializeJobs,
					AlwaysOverwriteSizeMatchingFiles: priorState.AlwaysOverwriteSizeMatchingFiles,
					Description:                      priorState.Description,
					DestinationReplaceFrom:           priorState.DestinationReplaceFrom,
					DestinationReplaceTo:             priorState.DestinationReplaceTo,
					Destinations:                     priorState.Destinations,
					Disabled:                         priorState.Disabled,
					ExcludePattern:                   priorState.ExcludePattern,
					ImportUrls:                       priorState.ImportUrls,
					FlattenDestinationStructure:      priorState.FlattenDestinationStructure,
					GroupIds:                         priorState.GroupIds,
					IgnoreLockedFolders:              priorState.IgnoreLockedFolders,
					Interval:                         priorState.Interval,
					LegacyFolderMatching:             priorState.LegacyFolderMatching,
					Name:                             priorState.Name,
					OverwriteFiles:                   priorState.OverwriteFiles,
					Path:                             priorState.Path,
					PathTimeZone:                     priorState.PathTimeZone,
					RecurringDay:                     priorState.RecurringDay,
					RecurringDays:                    priorState.RecurringDays,
					ScheduleId:                       priorState.ScheduleId,
					RetryOnFailureIntervalInMinutes:  priorState.RetryOnFailureIntervalInMinutes,
					RetryOnFailureNumberOfAttempts:   priorState.RetryOnFailureNumberOfAttempts,
					ScheduleDaysOfWeek:               priorState.ScheduleDaysOfWeek,
					ScheduleTimesOfDay:               priorState.ScheduleTimesOfDay,
					ScheduleTimeZone:                 priorState.ScheduleTimeZone,
					Source:                           priorState.Source,
					SyncIds:                          priorState.SyncIds,
					TriggerActions:                   priorState.TriggerActions,
					Trigger:                          priorState.Trigger,
					UserIds:                          priorState.UserIds,
					Value:                            priorState.Value,
					HolidayRegion:                    priorState.HolidayRegion,
					Id:                               priorState.Id,
					Deleted:                          priorState.Deleted,
					InboundEmailAddress:              priorState.InboundEmailAddress,
					LastModifiedAt:                   priorState.LastModifiedAt,
					Version:                          priorState.Version,
					Schedule:                         priorState.Schedule,
					HumanReadableSchedule:            priorState.HumanReadableSchedule,
					UserId:                           priorState.UserId,
					WebhookUrl:                       priorState.WebhookUrl,
				}
				currentSchema := r.resourceSchema()
				definitionValue, conversionDiags := lib.JSONValueToAPI(ctx, path.Root("definition"), priorState.Definition)
				resp.Diagnostics.Append(conversionDiags...)
				if resp.Diagnostics.HasError() {
					return
				}
				definitionValue, transformDiags0 := lib.WrapDiscriminatedUnionAtPath(ctx, path.Root("definition"), definitionValue, []string{"nodes"}, "type", []lib.JSONSchemaVariant{{Name: "trigger_scheduled", Value: "trigger_scheduled"}, {Name: "trigger_manual", Value: "trigger_manual"}, {Name: "trigger_action", Value: "trigger_action"}, {Name: "trigger_webhook", Value: "trigger_webhook"}, {Name: "trigger_email", Value: "trigger_email"}, {Name: "create_folder", Value: "create_folder"}, {Name: "copy_file", Value: "copy_file"}, {Name: "move_file", Value: "move_file"}, {Name: "delete_file", Value: "delete_file"}, {Name: "import_file", Value: "import_file"}, {Name: "run_sync", Value: "run_sync"}, {Name: "as2_send", Value: "as2_send"}, {Name: "send_email", Value: "send_email"}, {Name: "agent_compute", Value: "agent_compute"}, {Name: "set_metadata", Value: "set_metadata"}, {Name: "extract", Value: "extract"}, {Name: "document_convert", Value: "document_convert"}, {Name: "image_convert", Value: "image_convert"}, {Name: "zip", Value: "zip"}, {Name: "unzip", Value: "unzip"}, {Name: "gpg_encrypt", Value: "gpg_encrypt"}, {Name: "gpg_decrypt", Value: "gpg_decrypt"}, {Name: "if", Value: "if"}, {Name: "switch", Value: "switch"}, {Name: "filter", Value: "filter"}, {Name: "join", Value: "join"}, {Name: "aggregate", Value: "aggregate"}, {Name: "wait", Value: "wait"}, {Name: "transform", Value: "transform"}, {Name: "run_automation", Value: "run_automation"}})
				resp.Diagnostics.Append(transformDiags0...)
				definitionType := currentSchema.Attributes["definition"].GetType().(types.ObjectType)
				upgradedState.Definition, conversionDiags = lib.ToObject(ctx, path.Root("definition"), definitionValue, types.ObjectNull(definitionType.AttrTypes))
				resp.Diagnostics.Append(conversionDiags...)
				if resp.Diagnostics.HasError() {
					return
				}
				resp.Diagnostics.Append(resp.State.Set(ctx, upgradedState)...)
			},
		},
	}
}
