package provider

import (
	"context"
	"fmt"
	"strconv"
	"strings"

	files_sdk "github.com/Files-com/files-sdk-go/v3"
	automation "github.com/Files-com/files-sdk-go/v3/automation"
	"github.com/Files-com/terraform-provider-files/lib"
	"github.com/hashicorp/terraform-plugin-framework-validators/stringvalidator"
	"github.com/hashicorp/terraform-plugin-framework/diag"
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

var (
	_ resource.Resource                = &automationResource{}
	_ resource.ResourceWithConfigure   = &automationResource{}
	_ resource.ResourceWithImportState = &automationResource{}
)

func NewAutomationResource() resource.Resource {
	return &automationResource{}
}

type automationResource struct {
	client *automation.Client
}

type automationResourceModel struct {
	Automation                       types.String  `tfsdk:"automation"`
	AlwaysOverwriteSizeMatchingFiles types.Bool    `tfsdk:"always_overwrite_size_matching_files"`
	Description                      types.String  `tfsdk:"description"`
	DestinationReplaceFrom           types.String  `tfsdk:"destination_replace_from"`
	DestinationReplaceTo             types.String  `tfsdk:"destination_replace_to"`
	Destinations                     types.List    `tfsdk:"destinations"`
	Disabled                         types.Bool    `tfsdk:"disabled"`
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
	Schedule                         types.Dynamic `tfsdk:"schedule"`
	ScheduleDaysOfWeek               types.List    `tfsdk:"schedule_days_of_week"`
	ScheduleTimesOfDay               types.List    `tfsdk:"schedule_times_of_day"`
	ScheduleTimeZone                 types.String  `tfsdk:"schedule_time_zone"`
	Source                           types.String  `tfsdk:"source"`
	SyncIds                          types.List    `tfsdk:"sync_ids"`
	TriggerActions                   types.List    `tfsdk:"trigger_actions"`
	Trigger                          types.String  `tfsdk:"trigger"`
	UserIds                          types.List    `tfsdk:"user_ids"`
	Value                            types.Dynamic `tfsdk:"value"`
	Destination                      types.String  `tfsdk:"destination"`
	Id                               types.Int64   `tfsdk:"id"`
	Deleted                          types.Bool    `tfsdk:"deleted"`
	LastModifiedAt                   types.String  `tfsdk:"last_modified_at"`
	HumanReadableSchedule            types.String  `tfsdk:"human_readable_schedule"`
	UserId                           types.Int64   `tfsdk:"user_id"`
	WebhookUrl                       types.String  `tfsdk:"webhook_url"`
}

func (r *automationResource) Configure(_ context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {
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

	r.client = &automation.Client{Config: sdk_config}
}

func (r *automationResource) Metadata(_ context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_automation"
}

func (r *automationResource) Schema(_ context.Context, _ resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Description: "Automations allow you to automate workflows on your Files.com site.\n\n\n\nAutomations are different from Behaviors because Behaviors are associated with a current folder, while Automations apply across your entire site.\n\n\n\nAutomations are never removed when folders are removed, while Behaviors are removed when the associated folder is removed.\n\n\n\n## Path Matching\n\n\n\nThe `path` attribute specifies which folders this automation applies to.\n\nIt gets combined with the `source` attribute to determine which files are actually affected by the automation.\n\nNote that the `path` attribute supports globs, and only refers to _folders_.\n\nIt's the `source` attribute combined with the `path` attribute that determines which files are affected, and automations only operate on the files themselves.\n\nAdditionally, paths in Automations can refer to folders which don't yet exist.\n\n\n\n### Path Globs\n\n\n\nAlthough Automations may have a `path` specified, it can be a glob (which includes wildcards), which affects multiple folders.\n\n\n\n`*` matches any folder at that level of the path, but not subfolders. For example, `path/to/*` matches `path/to/folder1` and `path/to/folder2`, but not `path/to/folder1/subfolder`.\n\n\n\n`**` matches subfolders recursively. For example, `path/to/**` matches `path/to/folder1`, `path/to/folder1/subfolder`, `path/to/folder2`, `path/to/folder2/subfolder`, etc.\n\n\n\n`?` matches any one character.\n\n\n\nUse square brackets `[]` to match any character from a set. This works like a regular expression, including negation using `^`.\n\n\n\nCurly brackets `{}` can be used to denote parts of a pattern which will accept a number of alternatives, separated by commas `,`.\n\nThese alternatives can either be literal text or include special characters including nested curly brackets.\n\nFor example `{Mon,Tue,Wed,Thu,Fri}` would match abbreviated weekdays, and `202{3-{0[7-9],1?},4-0[1-6]}-*` would match dates from `2023-07-01` through `2024-06-30`.\n\n\n\nTo match any of the special characters literally, precede it with a backslash and enclose that pair with square brackets. For example to match a literal `?`, use `[\\?]`.\n\n\n\nGlobs are not supported on remote paths of any kind.\n\n\n\nBy default, Copy and Move automations that use globs will implicitly replicate matched folder structures at the destination. If you want to flatten the folder structure, set `flatten_destination_structure` to `true`.\n\n\n\n## Automation Triggers\n\n\n\nAutomations can be triggered in the following ways:\n\n\n\n* `custom_schedule` : The automation will run according to the custom schedule parameters for `days_of_week` (0-based) and `times_of_day`. A time zone may be specified via `time_zone` in Rails TimeZone name format.\n\n* `daily` : The automation will run once in a picked `interval`. You can specify `recurring_day` to select day number inside a picked `interval` it should be run on.\n\n* `webhook` : the automation will run when a request is sent to the corresponding webhook URL.\n\n* `action` : The automation will run when a specific action happens, e.g. a file is created or downloaded.\n\n\n\nFuture enhancements will allow Automations to be triggered by an incoming email, or by other services.\n\n\n\nCurrently, all Automation types support all triggers, with the following exceptions: `Create Folder` and `Run Remote Server Sync` are not supported by the `action` trigger.\n\n\n\nAutomations can be triggered manually if trigger is not set to `action`.\n\n\n\n## Destinations\n\n\n\nThe `destinations` parameter is a list of paths where files will be copied, moved, or created. It may include formatting parameters to dynamically determine the destination at runtime.\n\n\n\n### Relative vs. Absolute Paths\n\n\n\nIn order to specify a relative path, it must start with either `./` or `../`. All other paths are considered absolute. In general, leading slashes should never be used on Files.com paths, including here. Paths are interpreted as absolute in all contexts, even without a leading slash.\n\n\n\n### Files vs. Folders\n\n\n\nIf the destination path ends with a `/`, the filename from the source path will be preserved and put into the folder of this name. If the destination path does not end with a `/`, it will be interpreted as a filename and will override the source file's filename entirely.\n\n\n\n### Formatting Parameters\n\n\n\n**Action-Triggered Automations**\n\n\n\n* `%tf` : The name of the file that triggered the automation.\n\n* `%tp` : The path of the file that triggered the automation.\n\n* `%td` : The directory of the file that triggered the automation.\n\n\n\nFor example, if the triggering file is at `path/to/file.txt`, then the automation destination `path/to/dest/incoming-%tf` will result in the actual destination being `path/to/dest/incoming-file.txt`.\n\n\n\n**Parent Folders**\n\n\n\nTo reference the parent folder of a source file, use `%p1`, `%p2`, `%p3`, etc. for the first, second, third, etc. parent folder, respectively.\n\n\n\nFor example, if the source file is at `accounts/file.txt`, then the automation destination `path/to/dest/%p1/some_file_name.txt` will result in the actual destination being `path/to/dest/accounts/some_file_name.txt`.\n\n\n\n**Dates and Times**\n\n\n\n* `%Y` : The current year (4 digits)\n\n* `%m` : The current month (2 digits)\n\n* `%B` : The current month (full name)\n\n* `%d` : The current day (2 digits)\n\n* `%H` : The current hour (2 digits, 24-hour clock)\n\n* `%M` : The current minute (2 digits)\n\n* `%S` : The current second (2 digits)\n\n* `%z` : UTC Time Zone (e.g. -0900)\n\n\n\nFor example, if the current date is June 23, 2023 and the source file is named `daily_sales.csv`, then the following automation destination `path/to/dest/%Y/%m/%d/` will result in the actual destination being `path/to/dest/2023/06/23/daily_sales.csv`.\n\n\n\n### Replacing Text\n\n\n\nTo replace text in the source filename, use the `destination_replace_from` and `destination_replace_to` parameters. This will perform a simple text replacement on the source filename before inserting it into the destination path.\n\n\n\nFor example, if the `destination_replace_from` is `incoming` and the `destination_replace_to` is `outgoing`, then `path/to/incoming.txt` will translate to `path/to/outgoing.txt`.\n\n\n\n\n\n## Automation Types\n\n\n\nThere are several types of automations: Create Folder, Copy File, Move File, Delete File and, Run Remote Server Sync.\n\n\n\n\n\n### Create Folder\n\n\n\nCreates the folder with named by `destinations` in the path named by `path`.\n\nDestination may include formatting parameters to insert the date/time into the destination name.\n\n\n\nExample Use case: Our business files sales tax for each division in 11 states every quarter.\n\nI want to create the folders where those sales tax forms and data will be collected.\n\n\n\nI could create a Create Folder automation as follows:\n\n\n\n* Trigger: `daily`\n\n* Interval: `quarter_end`\n\n* Path: `AccountingAndTax/SalesTax/State/*/`\n\n* Destinations: `%Y/Quarter-ending-%m-%d`\n\n\n\nNote this assumes you have folders in `AccountingAndTax/SalesTax/State/` already created for each state, e.g. `AccountingAndTax/SalesTax/State/CA/`.\n\n\n\n### Delete File\n\n\n\nDeletes the file with path matching `source` (wildcards allowed) in the path named by `path`.\n\n\n\n\n\n### Copy File\n\n\n\nCopies files in the folder named by `path` to the path specified in `destinations`.\n\nThe automation will only fire on files matching the `source` (wildcards allowed). In the case of an action-triggered automation, it will only operate on the actual file that triggered the automation.\n\nIf the parameter `limit` exists, the automation will only copy the newest `limit` files in each matching folder.\n\n\n\n\n\n### Move File\n\n\n\nMoves files in the folder named by `path` to the path specified in `destinations`.\n\nThe automation will only fire on files matching the `source` (wildcards allowed). In the case of an action-triggered automation, it will only operate on the actual file that triggered the automation.\n\nIf the parameter `limit` exists, the automation will only move the newest `limit` files in each matching folder.\n\nNote that for a move with multiple destinations, all but one destination is treated as a copy.\n\n\n\n### Run Remote Server Sync\n\n\n\nThe Run Remote Server Sync automation runs the remote server syncs specified by the `sync_ids`.\n\n\n\nTypically when this automation is used, the remote server syncs in question are set to the manual\n\nscheduling mode (`manual` to `true` via the API) to disable the built in sync scheduler.\n\n\n\n\n\n### Help us build the future of Automations\n\n\n\nDo you have an idea for something that would work well as a Files.com Automation? Let us know!\n\nWe are actively improving the types of automations offered on our platform.",
		Attributes: map[string]schema.Attribute{
			"automation": schema.StringAttribute{
				Description: "Automation type",
				Required:    true,
				Validators: []validator.String{
					stringvalidator.OneOf("create_folder", "delete_file", "copy_file", "move_file", "as2_send", "run_sync"),
				},
			},
			"always_overwrite_size_matching_files": schema.BoolAttribute{
				Description: "Ordinarily, files with identical size in the source and destination will be skipped from copy operations to prevent wasted transfer.  If this flag is `true` we will overwrite the destination file always.  Note that this may cause large amounts of wasted transfer usage.",
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
				Description: "If true, existing files will be overwritten with new files on Move/Copy automations.  Note: by default files will not be overwritten if they appear to be the same file size as the newly incoming file.  Use the `:always_overwrite_size_matching_files` option to override this.",
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
			"schedule": schema.DynamicAttribute{
				Description: "If trigger is `custom_schedule`, Custom schedule description for when the automation should be run in json format.",
				Computed:    true,
				Optional:    true,
				PlanModifiers: []planmodifier.Dynamic{
					dynamicplanmodifier.UseStateForUnknown(),
				},
			},
			"schedule_days_of_week": schema.ListAttribute{
				Description: "If trigger is `custom_schedule`, Custom schedule description for when the automation should be run. 0-based days of the week. 0 is Sunday, 1 is Monday, etc.",
				Computed:    true,
				Optional:    true,
				ElementType: types.Int64Type,
				PlanModifiers: []planmodifier.List{
					listplanmodifier.UseStateForUnknown(),
				},
			},
			"schedule_times_of_day": schema.ListAttribute{
				Description: "If trigger is `custom_schedule`, Custom schedule description for when the automation should be run. Times of day in HH:MM format.",
				Computed:    true,
				Optional:    true,
				ElementType: types.StringType,
				PlanModifiers: []planmodifier.List{
					listplanmodifier.UseStateForUnknown(),
				},
			},
			"schedule_time_zone": schema.StringAttribute{
				Description: "If trigger is `custom_schedule`, Custom schedule Time Zone for when the automation should be run.",
				Computed:    true,
				Optional:    true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
			},
			"source": schema.StringAttribute{
				Description: "Source Path",
				Computed:    true,
				Optional:    true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
			},
			"sync_ids": schema.ListAttribute{
				Description: "IDs of remote sync folder behaviors to run by this Automation",
				Computed:    true,
				Optional:    true,
				ElementType: types.Int64Type,
				PlanModifiers: []planmodifier.List{
					listplanmodifier.UseStateForUnknown(),
				},
			},
			"trigger_actions": schema.ListAttribute{
				Description: "If trigger is `action`, this is the list of action types on which to trigger the automation. Valid actions are create, read, update, destroy, move, copy",
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
					stringvalidator.OneOf("daily", "custom_schedule", "webhook", "email", "action"),
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
				PlanModifiers: []planmodifier.Dynamic{
					dynamicplanmodifier.UseStateForUnknown(),
				},
			},
			"destination": schema.StringAttribute{
				Optional: true,
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
			"last_modified_at": schema.StringAttribute{
				Description: "Time when automation was last modified. Does not change for name or description updates.",
				Computed:    true,
			},
			"human_readable_schedule": schema.StringAttribute{
				Description: "If trigger is `custom_schedule`, Human readable Custom schedule description for when the automation should be run.",
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
	}
}

func (r *automationResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	var plan automationResourceModel
	diags := req.Plan.Get(ctx, &plan)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	paramsAutomationCreate := files_sdk.AutomationCreateParams{}
	paramsAutomationCreate.Source = plan.Source.ValueString()
	paramsAutomationCreate.Destination = plan.Destination.ValueString()
	if !plan.Destinations.IsNull() && !plan.Destinations.IsUnknown() {
		diags = plan.Destinations.ElementsAs(ctx, &paramsAutomationCreate.Destinations, false)
		resp.Diagnostics.Append(diags...)
	}
	paramsAutomationCreate.DestinationReplaceFrom = plan.DestinationReplaceFrom.ValueString()
	paramsAutomationCreate.DestinationReplaceTo = plan.DestinationReplaceTo.ValueString()
	paramsAutomationCreate.Interval = plan.Interval.ValueString()
	paramsAutomationCreate.Path = plan.Path.ValueString()
	paramsAutomationCreate.SyncIds, diags = lib.ListValueToString(ctx, path.Root("sync_ids"), plan.SyncIds, ",")
	resp.Diagnostics.Append(diags...)
	paramsAutomationCreate.UserIds, diags = lib.ListValueToString(ctx, path.Root("user_ids"), plan.UserIds, ",")
	resp.Diagnostics.Append(diags...)
	paramsAutomationCreate.GroupIds, diags = lib.ListValueToString(ctx, path.Root("group_ids"), plan.GroupIds, ",")
	resp.Diagnostics.Append(diags...)
	createSchedule, diags := lib.DynamicToStringMap(ctx, path.Root("schedule"), plan.Schedule)
	resp.Diagnostics.Append(diags...)
	paramsAutomationCreate.Schedule = createSchedule
	if !plan.ScheduleDaysOfWeek.IsNull() && !plan.ScheduleDaysOfWeek.IsUnknown() {
		diags = plan.ScheduleDaysOfWeek.ElementsAs(ctx, &paramsAutomationCreate.ScheduleDaysOfWeek, false)
		resp.Diagnostics.Append(diags...)
	}
	if !plan.ScheduleTimesOfDay.IsNull() && !plan.ScheduleTimesOfDay.IsUnknown() {
		diags = plan.ScheduleTimesOfDay.ElementsAs(ctx, &paramsAutomationCreate.ScheduleTimesOfDay, false)
		resp.Diagnostics.Append(diags...)
	}
	paramsAutomationCreate.ScheduleTimeZone = plan.ScheduleTimeZone.ValueString()
	paramsAutomationCreate.AlwaysOverwriteSizeMatchingFiles = plan.AlwaysOverwriteSizeMatchingFiles.ValueBoolPointer()
	paramsAutomationCreate.Description = plan.Description.ValueString()
	paramsAutomationCreate.Disabled = plan.Disabled.ValueBoolPointer()
	paramsAutomationCreate.FlattenDestinationStructure = plan.FlattenDestinationStructure.ValueBoolPointer()
	paramsAutomationCreate.IgnoreLockedFolders = plan.IgnoreLockedFolders.ValueBoolPointer()
	paramsAutomationCreate.LegacyFolderMatching = plan.LegacyFolderMatching.ValueBoolPointer()
	paramsAutomationCreate.Name = plan.Name.ValueString()
	paramsAutomationCreate.OverwriteFiles = plan.OverwriteFiles.ValueBoolPointer()
	paramsAutomationCreate.PathTimeZone = plan.PathTimeZone.ValueString()
	paramsAutomationCreate.Trigger = paramsAutomationCreate.Trigger.Enum()[plan.Trigger.ValueString()]
	if !plan.TriggerActions.IsNull() && !plan.TriggerActions.IsUnknown() {
		diags = plan.TriggerActions.ElementsAs(ctx, &paramsAutomationCreate.TriggerActions, false)
		resp.Diagnostics.Append(diags...)
	}
	createValue, diags := lib.DynamicToStringMap(ctx, path.Root("value"), plan.Value)
	resp.Diagnostics.Append(diags...)
	paramsAutomationCreate.Value = createValue
	paramsAutomationCreate.RecurringDay = plan.RecurringDay.ValueInt64()
	paramsAutomationCreate.Automation = paramsAutomationCreate.Automation.Enum()[plan.Automation.ValueString()]

	if resp.Diagnostics.HasError() {
		return
	}

	automation, err := r.client.Create(paramsAutomationCreate, files_sdk.WithContext(ctx))
	if err != nil {
		resp.Diagnostics.AddError(
			"Error Creating Files Automation",
			"Could not create automation, unexpected error: "+err.Error(),
		)
		return
	}

	diags = r.populateResourceModel(ctx, automation, &plan)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	diags = resp.State.Set(ctx, plan)
	resp.Diagnostics.Append(diags...)
}

func (r *automationResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	var state automationResourceModel
	diags := req.State.Get(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	paramsAutomationFind := files_sdk.AutomationFindParams{}
	paramsAutomationFind.Id = state.Id.ValueInt64()

	automation, err := r.client.Find(paramsAutomationFind, files_sdk.WithContext(ctx))
	if err != nil {
		resp.Diagnostics.AddError(
			"Error Reading Files Automation",
			"Could not read automation id "+fmt.Sprint(state.Id.ValueInt64())+": "+err.Error(),
		)
		return
	}

	diags = r.populateResourceModel(ctx, automation, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	diags = resp.State.Set(ctx, state)
	resp.Diagnostics.Append(diags...)
}

func (r *automationResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	var plan automationResourceModel
	diags := req.Plan.Get(ctx, &plan)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	paramsAutomationUpdate := files_sdk.AutomationUpdateParams{}
	paramsAutomationUpdate.Id = plan.Id.ValueInt64()
	paramsAutomationUpdate.Source = plan.Source.ValueString()
	paramsAutomationUpdate.Destination = plan.Destination.ValueString()
	if !plan.Destinations.IsNull() && !plan.Destinations.IsUnknown() {
		diags = plan.Destinations.ElementsAs(ctx, &paramsAutomationUpdate.Destinations, false)
		resp.Diagnostics.Append(diags...)
	}
	paramsAutomationUpdate.DestinationReplaceFrom = plan.DestinationReplaceFrom.ValueString()
	paramsAutomationUpdate.DestinationReplaceTo = plan.DestinationReplaceTo.ValueString()
	paramsAutomationUpdate.Interval = plan.Interval.ValueString()
	paramsAutomationUpdate.Path = plan.Path.ValueString()
	paramsAutomationUpdate.SyncIds, diags = lib.ListValueToString(ctx, path.Root("sync_ids"), plan.SyncIds, ",")
	resp.Diagnostics.Append(diags...)
	paramsAutomationUpdate.UserIds, diags = lib.ListValueToString(ctx, path.Root("user_ids"), plan.UserIds, ",")
	resp.Diagnostics.Append(diags...)
	paramsAutomationUpdate.GroupIds, diags = lib.ListValueToString(ctx, path.Root("group_ids"), plan.GroupIds, ",")
	resp.Diagnostics.Append(diags...)
	updateSchedule, diags := lib.DynamicToStringMap(ctx, path.Root("schedule"), plan.Schedule)
	resp.Diagnostics.Append(diags...)
	paramsAutomationUpdate.Schedule = updateSchedule
	if !plan.ScheduleDaysOfWeek.IsNull() && !plan.ScheduleDaysOfWeek.IsUnknown() {
		diags = plan.ScheduleDaysOfWeek.ElementsAs(ctx, &paramsAutomationUpdate.ScheduleDaysOfWeek, false)
		resp.Diagnostics.Append(diags...)
	}
	if !plan.ScheduleTimesOfDay.IsNull() && !plan.ScheduleTimesOfDay.IsUnknown() {
		diags = plan.ScheduleTimesOfDay.ElementsAs(ctx, &paramsAutomationUpdate.ScheduleTimesOfDay, false)
		resp.Diagnostics.Append(diags...)
	}
	paramsAutomationUpdate.ScheduleTimeZone = plan.ScheduleTimeZone.ValueString()
	paramsAutomationUpdate.AlwaysOverwriteSizeMatchingFiles = plan.AlwaysOverwriteSizeMatchingFiles.ValueBoolPointer()
	paramsAutomationUpdate.Description = plan.Description.ValueString()
	paramsAutomationUpdate.Disabled = plan.Disabled.ValueBoolPointer()
	paramsAutomationUpdate.FlattenDestinationStructure = plan.FlattenDestinationStructure.ValueBoolPointer()
	paramsAutomationUpdate.IgnoreLockedFolders = plan.IgnoreLockedFolders.ValueBoolPointer()
	paramsAutomationUpdate.LegacyFolderMatching = plan.LegacyFolderMatching.ValueBoolPointer()
	paramsAutomationUpdate.Name = plan.Name.ValueString()
	paramsAutomationUpdate.OverwriteFiles = plan.OverwriteFiles.ValueBoolPointer()
	paramsAutomationUpdate.PathTimeZone = plan.PathTimeZone.ValueString()
	paramsAutomationUpdate.Trigger = paramsAutomationUpdate.Trigger.Enum()[plan.Trigger.ValueString()]
	if !plan.TriggerActions.IsNull() && !plan.TriggerActions.IsUnknown() {
		diags = plan.TriggerActions.ElementsAs(ctx, &paramsAutomationUpdate.TriggerActions, false)
		resp.Diagnostics.Append(diags...)
	}
	updateValue, diags := lib.DynamicToStringMap(ctx, path.Root("value"), plan.Value)
	resp.Diagnostics.Append(diags...)
	paramsAutomationUpdate.Value = updateValue
	paramsAutomationUpdate.RecurringDay = plan.RecurringDay.ValueInt64()
	paramsAutomationUpdate.Automation = paramsAutomationUpdate.Automation.Enum()[plan.Automation.ValueString()]

	if resp.Diagnostics.HasError() {
		return
	}

	automation, err := r.client.Update(paramsAutomationUpdate, files_sdk.WithContext(ctx))
	if err != nil {
		resp.Diagnostics.AddError(
			"Error Updating Files Automation",
			"Could not update automation, unexpected error: "+err.Error(),
		)
		return
	}

	diags = r.populateResourceModel(ctx, automation, &plan)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	diags = resp.State.Set(ctx, plan)
	resp.Diagnostics.Append(diags...)
}

func (r *automationResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	var state automationResourceModel
	diags := req.State.Get(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	paramsAutomationDelete := files_sdk.AutomationDeleteParams{}
	paramsAutomationDelete.Id = state.Id.ValueInt64()

	err := r.client.Delete(paramsAutomationDelete, files_sdk.WithContext(ctx))
	if err != nil {
		resp.Diagnostics.AddError(
			"Error Deleting Files Automation",
			"Could not delete automation id "+fmt.Sprint(state.Id.ValueInt64())+": "+err.Error(),
		)
	}
}

func (r *automationResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
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

func (r *automationResource) populateResourceModel(ctx context.Context, automation files_sdk.Automation, state *automationResourceModel) (diags diag.Diagnostics) {
	var propDiags diag.Diagnostics

	state.Id = types.Int64Value(automation.Id)
	state.AlwaysOverwriteSizeMatchingFiles = types.BoolPointerValue(automation.AlwaysOverwriteSizeMatchingFiles)
	state.Automation = types.StringValue(automation.Automation)
	state.Deleted = types.BoolPointerValue(automation.Deleted)
	state.Description = types.StringValue(automation.Description)
	state.DestinationReplaceFrom = types.StringValue(automation.DestinationReplaceFrom)
	state.DestinationReplaceTo = types.StringValue(automation.DestinationReplaceTo)
	state.Destinations, propDiags = types.ListValueFrom(ctx, types.StringType, automation.Destinations)
	diags.Append(propDiags...)
	state.Disabled = types.BoolPointerValue(automation.Disabled)
	state.FlattenDestinationStructure = types.BoolPointerValue(automation.FlattenDestinationStructure)
	state.GroupIds, propDiags = types.ListValueFrom(ctx, types.Int64Type, automation.GroupIds)
	diags.Append(propDiags...)
	state.IgnoreLockedFolders = types.BoolPointerValue(automation.IgnoreLockedFolders)
	state.Interval = types.StringValue(automation.Interval)
	if err := lib.TimeToStringType(ctx, path.Root("last_modified_at"), automation.LastModifiedAt, &state.LastModifiedAt); err != nil {
		diags.AddError(
			"Error Creating Files Automation",
			"Could not convert state last_modified_at to string: "+err.Error(),
		)
	}
	state.LegacyFolderMatching = types.BoolPointerValue(automation.LegacyFolderMatching)
	state.Name = types.StringValue(automation.Name)
	state.OverwriteFiles = types.BoolPointerValue(automation.OverwriteFiles)
	state.Path = types.StringValue(automation.Path)
	state.PathTimeZone = types.StringValue(automation.PathTimeZone)
	state.RecurringDay = types.Int64Value(automation.RecurringDay)
	state.Schedule, propDiags = lib.ToDynamic(ctx, path.Root("schedule"), automation.Schedule, state.Schedule.UnderlyingValue())
	diags.Append(propDiags...)
	state.HumanReadableSchedule = types.StringValue(automation.HumanReadableSchedule)
	state.ScheduleDaysOfWeek, propDiags = types.ListValueFrom(ctx, types.Int64Type, automation.ScheduleDaysOfWeek)
	diags.Append(propDiags...)
	state.ScheduleTimesOfDay, propDiags = types.ListValueFrom(ctx, types.StringType, automation.ScheduleTimesOfDay)
	diags.Append(propDiags...)
	state.ScheduleTimeZone = types.StringValue(automation.ScheduleTimeZone)
	state.Source = types.StringValue(automation.Source)
	state.SyncIds, propDiags = types.ListValueFrom(ctx, types.Int64Type, automation.SyncIds)
	diags.Append(propDiags...)
	state.TriggerActions, propDiags = types.ListValueFrom(ctx, types.StringType, automation.TriggerActions)
	diags.Append(propDiags...)
	state.Trigger = types.StringValue(automation.Trigger)
	state.UserId = types.Int64Value(automation.UserId)
	state.UserIds, propDiags = types.ListValueFrom(ctx, types.Int64Type, automation.UserIds)
	diags.Append(propDiags...)
	state.Value, propDiags = lib.ToDynamic(ctx, path.Root("value"), automation.Value, state.Value.UnderlyingValue())
	diags.Append(propDiags...)
	state.WebhookUrl = types.StringValue(automation.WebhookUrl)

	return
}
