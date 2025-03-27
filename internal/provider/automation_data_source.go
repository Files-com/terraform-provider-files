package provider

import (
	"context"
	"fmt"

	files_sdk "github.com/Files-com/files-sdk-go/v3"
	automation "github.com/Files-com/files-sdk-go/v3/automation"
	"github.com/Files-com/terraform-provider-files/lib"
	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/datasource/schema"
	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

var (
	_ datasource.DataSource              = &automationDataSource{}
	_ datasource.DataSourceWithConfigure = &automationDataSource{}
)

func NewAutomationDataSource() datasource.DataSource {
	return &automationDataSource{}
}

type automationDataSource struct {
	client *automation.Client
}

type automationDataSourceModel struct {
	Id                               types.Int64   `tfsdk:"id"`
	AlwaysOverwriteSizeMatchingFiles types.Bool    `tfsdk:"always_overwrite_size_matching_files"`
	Automation                       types.String  `tfsdk:"automation"`
	Deleted                          types.Bool    `tfsdk:"deleted"`
	Description                      types.String  `tfsdk:"description"`
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
	LastModifiedAt                   types.String  `tfsdk:"last_modified_at"`
	LegacyFolderMatching             types.Bool    `tfsdk:"legacy_folder_matching"`
	Name                             types.String  `tfsdk:"name"`
	OverwriteFiles                   types.Bool    `tfsdk:"overwrite_files"`
	Path                             types.String  `tfsdk:"path"`
	PathTimeZone                     types.String  `tfsdk:"path_time_zone"`
	RecurringDay                     types.Int64   `tfsdk:"recurring_day"`
	RetryOnFailureIntervalInMinutes  types.Int64   `tfsdk:"retry_on_failure_interval_in_minutes"`
	RetryOnFailureNumberOfAttempts   types.Int64   `tfsdk:"retry_on_failure_number_of_attempts"`
	Schedule                         types.Dynamic `tfsdk:"schedule"`
	HumanReadableSchedule            types.String  `tfsdk:"human_readable_schedule"`
	ScheduleDaysOfWeek               types.List    `tfsdk:"schedule_days_of_week"`
	ScheduleTimesOfDay               types.List    `tfsdk:"schedule_times_of_day"`
	ScheduleTimeZone                 types.String  `tfsdk:"schedule_time_zone"`
	Source                           types.String  `tfsdk:"source"`
	SyncIds                          types.List    `tfsdk:"sync_ids"`
	TriggerActions                   types.List    `tfsdk:"trigger_actions"`
	Trigger                          types.String  `tfsdk:"trigger"`
	UserId                           types.Int64   `tfsdk:"user_id"`
	UserIds                          types.List    `tfsdk:"user_ids"`
	Value                            types.Dynamic `tfsdk:"value"`
	WebhookUrl                       types.String  `tfsdk:"webhook_url"`
}

func (r *automationDataSource) Configure(_ context.Context, req datasource.ConfigureRequest, resp *datasource.ConfigureResponse) {
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

func (r *automationDataSource) Metadata(_ context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_automation"
}

func (r *automationDataSource) Schema(_ context.Context, _ datasource.SchemaRequest, resp *datasource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Description: "An Automation is an automated process of controlling workflows on your Files.com site.\n\n\n\nAutomations are different from Behaviors because Behaviors are associated with a current folder, while Automations apply across your entire site.\n\n\n\nAutomations are never removed when folders are removed, while Behaviors are removed when the associated folder is removed.\n\n\n\n## Path Matching\n\n\n\nThe `path` attribute specifies which folders this automation applies to.\n\nIt gets combined with the `source` attribute to determine which files are actually affected by the automation.\n\nNote that the `path` attribute supports globs, and only refers to _folders_.\n\nIt's the `source` attribute, which also supports globs, combined with the `path` attribute that determines which files are affected, and automations only operate on the files themselves.\n\nAdditionally, paths in Automations can refer to folders which don't yet exist.\n\n\n\n### Path Globs\n\n\n\nAlthough Automations may have a `path` specified, it can be a glob (which includes wildcards), which affects multiple folders.\n\n\n\n`*` matches any folder at that level of the path, but not subfolders. For example, `path/to/*` matches `path/to/folder1` and `path/to/folder2`, but not `path/to/folder1/subfolder`.\n\n\n\n`**` matches subfolders recursively. For example, `path/to/**` matches `path/to/folder1`, `path/to/folder1/subfolder`, `path/to/folder2`, `path/to/folder2/subfolder`, etc.\n\n\n\n`?` matches any one character.\n\n\n\nUse square brackets `[]` to match any character from a set. This works like a regular expression, including negation using `^`.\n\n\n\nCurly brackets `{}` can be used to denote parts of a pattern which will accept a number of alternatives, separated by commas `,`.\n\nThese alternatives can either be literal text or include special characters including nested curly brackets.\n\nFor example `{Mon,Tue,Wed,Thu,Fri}` would match abbreviated weekdays, and `202{3-{0[7-9],1?},4-0[1-6]}-*` would match dates from `2023-07-01` through `2024-06-30`.\n\n\n\nTo match any of the special characters literally, precede it with a backslash and enclose that pair with square brackets. For example to match a literal `?`, use `[\\?]`.\n\n\n\nGlobs are supported on `path`, `source`, and `exclude_pattern` fields. Globs are not supported on remote paths of any kind or for any field.\n\n\n\nBy default, Copy and Move automations that use globs will implicitly replicate matched folder structures at the destination. If you want to flatten the folder structure, set `flatten_destination_structure` to `true`.\n\n\n\n## Automation Triggers\n\n\n\nAutomations can be triggered in the following ways:\n\n\n\n* `custom_schedule` : The automation will run according to the custom schedule parameters for `days_of_week` (0-based) and `times_of_day`. A time zone may be specified via `time_zone` in Rails TimeZone name format.\n\n* `daily` : The automation will run once in a picked `interval`. You can specify `recurring_day` to select day number inside a picked `interval` it should be run on.\n\n* `webhook` : the automation will run when a request is sent to the corresponding webhook URL.\n\n* `action` : The automation will run when a specific action happens, e.g. a file is created or downloaded.\n\n\n\nFuture enhancements will allow Automations to be triggered by an incoming email, or by other services.\n\n\n\nCurrently, all Automation types support all triggers, with the following exceptions: `Create Folder` and `Run Remote Server Sync` are not supported by the `action` trigger.\n\n\n\nAutomations can be triggered manually if trigger is not set to `action`.\n\n\n\n## Destinations\n\n\n\nThe `destinations` parameter is a list of paths where files will be copied, moved, or created. It may include formatting parameters to dynamically determine the destination at runtime.\n\n\n\n### Relative vs. Absolute Paths\n\n\n\nIn order to specify a relative path, it must start with either `./` or `../`. All other paths are considered absolute. In general, leading slashes should never be used on Files.com paths, including here. Paths are interpreted as absolute in all contexts, even without a leading slash.\n\n\n\n### Files vs. Folders\n\n\n\nIf the destination path ends with a `/`, the filename from the source path will be preserved and put into the folder of this name. If the destination path does not end with a `/`, it will be interpreted as a filename and will override the source file's filename entirely.\n\n\n\n### Formatting Parameters\n\n\n\n**Action-Triggered Automations**\n\n\n\n* `%tf` : The name of the file that triggered the automation.\n\n* `%tp` : The path of the file that triggered the automation.\n\n* `%td` : The directory of the file that triggered the automation.\n\n* `%tb` : The name of the file (without extension) that triggered the automation.\n\n* `%te` : The extension of the file that triggered the automation.\n\n\n\nFor example, if the triggering file is at `path/to/file.txt`, then the automation destination `path/to/dest/incoming-%tf` will result in the actual destination being `path/to/dest/incoming-file.txt`.\n\n\n\n**Parent Folders**\n\n\n\nTo reference the parent folder of a source file, use `%p1`, `%p2`, `%p3`, etc. for the first, second, third, etc. parent folder, respectively.\n\n\n\nFor example, if the source file is at `accounts/file.txt`, then the automation destination `path/to/dest/%p1/some_file_name.txt` will result in the actual destination being `path/to/dest/accounts/some_file_name.txt`.\n\n\n\n**Dates and Times**\n\n\n\n* `%Y` : The current year (4 digits)\n\n* `%m` : The current month (2 digits)\n\n* `%B` : The current month (full name)\n\n* `%d` : The current day (2 digits)\n\n* `%H` : The current hour (2 digits, 24-hour clock)\n\n* `%M` : The current minute (2 digits)\n\n* `%S` : The current second (2 digits)\n\n* `%z` : UTC Time Zone (e.g. -0900)\n\n\n\nFor example, if the current date is June 23, 2023 and the source file is named `daily_sales.csv`, then the following automation destination `path/to/dest/%Y/%m/%d/` will result in the actual destination being `path/to/dest/2023/06/23/daily_sales.csv`.\n\n\n\n### Replacing Text\n\n\n\nTo replace text in the source filename, use the `destination_replace_from` and `destination_replace_to` parameters. This will perform a simple text replacement on the source filename before inserting it into the destination path.\n\n\n\nFor example, if the `destination_replace_from` is `incoming` and the `destination_replace_to` is `outgoing`, then `path/to/incoming.txt` will translate to `path/to/outgoing.txt`.\n\n\n\n\n\n## Automation Types\n\n\n\nThere are several types of automations: Create Folder, Copy File, Move File, Delete File and, Run Remote Server Sync.\n\n\n\n\n\n### Create Folder\n\n\n\nCreates the folder with named by `destinations` in the path named by `path`.\n\nDestination may include formatting parameters to insert the date/time into the destination name.\n\n\n\nExample Use case: Our business files sales tax for each division in 11 states every quarter.\n\nI want to create the folders where those sales tax forms and data will be collected.\n\n\n\nI could create a Create Folder automation as follows:\n\n\n\n* Trigger: `daily`\n\n* Interval: `quarter_end`\n\n* Path: `AccountingAndTax/SalesTax/State/*/`\n\n* Destinations: `%Y/Quarter-ending-%m-%d`\n\n\n\nNote this assumes you have folders in `AccountingAndTax/SalesTax/State/` already created for each state, e.g. `AccountingAndTax/SalesTax/State/CA/`.\n\n\n\n\n\n### Delete File\n\n\n\nDeletes the file with path matching `source` (wildcards allowed) in the path named by `path`.\n\n\n\n\n\n### Copy File\n\n\n\nCopies files in the folder named by `path` to the path specified in `destinations`.\n\nThe automation will only fire on files matching the `source` (wildcards allowed). In the case of an action-triggered automation, it will only operate on the actual file that triggered the automation.\n\nIf the parameter `limit` exists, the automation will only copy the newest `limit` files in each matching folder.\n\n\n\n\n\n### Move File\n\n\n\nMoves files in the folder named by `path` to the path specified in `destinations`.\n\nThe automation will only fire on files matching the `source` (wildcards allowed). In the case of an action-triggered automation, it will only operate on the actual file that triggered the automation.\n\nIf the parameter `limit` exists, the automation will only move the newest `limit` files in each matching folder.\n\nNote that for a move with multiple destinations, all but one destination is treated as a copy.\n\n\n\n\n\n### Run Remote Server Sync\n\n\n\nThe Run Remote Server Sync automation runs the remote server syncs specified by the `sync_ids`.\n\n\n\nTypically when this automation is used, the remote server syncs in question are set to the manual\n\nscheduling mode (`manual` to `true` via the API) to disable the built in sync scheduler.\n\n\n\n\n\n### Import File\n\n\n\nRetrieves files from one or more URLs and saves the results under the path specified in `destinations`.\n\n\n\nThe URLs to retrieve are specified as a JSON array in the `import_urls` property.\n\n\n\n```json\n\n[\n\n {\n\n \"name\": \"response.json\",\n\n \"url\": \"https://example.com/api\",\n\n \"method\": \"post\",\n\n \"headers\": {\n\n \"Content-Type\": \"application/json\"\n\n },\n\n \"content\": { \"trigger-file\": \"%tp\" }\n\n }\n\n]\n\n```\n\n\n\nThe recognized keys are:\n\n\n\n* `name`: The file name which will be used to save the returned content. Required. `%` tokens will be replaced as described under Formatting Parameters.\n\n* `url`: The URL which will be requested. Required.\n\n* `method`: The HTTP method to be used for the request. May be either `get` or `post` (case insensitive). Defaults to `get`.\n\n* `headers`: Optional headers to be included in the request. `%` tokens in the values will be replaced as described under Formatting Parameters.\n\n* `content`: Optional body to send for POST request. If supplied as a string, `%` tokens will be expanded. If supplied as a JSON Object, `%` tokens will be expanded for top-level values. Other JSON types will be sent as-is.\n\n\n\n\n\n### Help us build the future of Automations\n\n\n\nDo you have an idea for something that would work well as a Files.com Automation? Let us know!\n\nWe are actively improving the types of automations offered on our platform.\n\n\n\n\n\n## Retrying Failues\n\n\n\nAutomations will automatically retry individual action steps up to 3 times, with pauses between retries that increase from 15 seconds to 1 minute. If individual action steps fail after our 3rd attempt, that action will fail. If every action step in an Automation Run fails, that automation run will move to a `failure` status. If at least one step succeeds and one step fails, that automation run will move to a `partial_failure` status.\n\n\n\nAutomation Runs can be retried automatically when they enter a `failure` or `partial_failure` status as described above. A retry will re-run the automation from scratch, including the \"planning\" phase, which expands globs (wildcards) and identifies which files to transfer or skip.\n\n\n\nRetrying of entire Automation Runs must be explicitly enabled by setting the `retry_on_failure_interval_in_minutes` and `retry_on_failure_number_of_attempts` values on the Automation.\n\n\n\nWhen retrying entire Automation Runs, we currently do not skip action steps which were skipped or successfully completed in a previous version of the Automation Run. We will soon be adding this functionality, which should enhance the usefulness of the retry apparatus for Automations in most situations.",
		Attributes: map[string]schema.Attribute{
			"id": schema.Int64Attribute{
				Description: "Automation ID",
				Required:    true,
			},
			"always_overwrite_size_matching_files": schema.BoolAttribute{
				Description: "Ordinarily, files with identical size in the source and destination will be skipped from copy operations to prevent wasted transfer.  If this flag is `true` we will overwrite the destination file always.  Note that this may cause large amounts of wasted transfer usage.  This setting has no effect unless `overwrite_files` is also set to `true`.",
				Computed:    true,
			},
			"automation": schema.StringAttribute{
				Description: "Automation type",
				Computed:    true,
			},
			"deleted": schema.BoolAttribute{
				Description: "Indicates if the automation has been deleted.",
				Computed:    true,
			},
			"description": schema.StringAttribute{
				Description: "Description for the this Automation.",
				Computed:    true,
			},
			"destination_replace_from": schema.StringAttribute{
				Description: "If set, this string in the destination path will be replaced with the value in `destination_replace_to`.",
				Computed:    true,
			},
			"destination_replace_to": schema.StringAttribute{
				Description: "If set, this string will replace the value `destination_replace_from` in the destination filename. You can use special patterns here.",
				Computed:    true,
			},
			"destinations": schema.ListAttribute{
				Description: "Destination Paths",
				Computed:    true,
				ElementType: types.StringType,
			},
			"disabled": schema.BoolAttribute{
				Description: "If true, this automation will not run.",
				Computed:    true,
			},
			"exclude_pattern": schema.StringAttribute{
				Description: "If set, this glob pattern will exclude files from the automation. Supports globs, except on remote mounts.",
				Computed:    true,
			},
			"import_urls": schema.DynamicAttribute{
				Description: "List of URLs to be imported and names to be used.",
				Computed:    true,
			},
			"flatten_destination_structure": schema.BoolAttribute{
				Description: "Normally copy and move automations that use globs will implicitly preserve the source folder structure in the destination.  If this flag is `true`, the source folder structure will be flattened in the destination.  This is useful for copying or moving files from multiple folders into a single destination folder.",
				Computed:    true,
			},
			"group_ids": schema.ListAttribute{
				Description: "IDs of Groups for the Automation (i.e. who to Request File from)",
				Computed:    true,
				ElementType: types.Int64Type,
			},
			"ignore_locked_folders": schema.BoolAttribute{
				Description: "If true, the Lock Folders behavior will be disregarded for automated actions.",
				Computed:    true,
			},
			"interval": schema.StringAttribute{
				Description: "If trigger is `daily`, this specifies how often to run this automation.  One of: `day`, `week`, `week_end`, `month`, `month_end`, `quarter`, `quarter_end`, `year`, `year_end`",
				Computed:    true,
			},
			"last_modified_at": schema.StringAttribute{
				Description: "Time when automation was last modified. Does not change for name or description updates.",
				Computed:    true,
			},
			"legacy_folder_matching": schema.BoolAttribute{
				Description: "If `true`, use the legacy behavior for this automation, where it can operate on folders in addition to just files.  This behavior no longer works and should not be used.",
				Computed:    true,
			},
			"name": schema.StringAttribute{
				Description: "Name for this automation.",
				Computed:    true,
			},
			"overwrite_files": schema.BoolAttribute{
				Description: "If true, existing files will be overwritten with new files on Move/Copy automations.  Note: by default files will not be overwritten if they appear to be the same file size as the newly incoming file.  Use the `always_overwrite_size_matching_files` option in conjunction with `overwrite_files` to override this behavior and overwrite files no matter what.",
				Computed:    true,
			},
			"path": schema.StringAttribute{
				Description: "Path on which this Automation runs.  Supports globs, except on remote mounts. This must be slash-delimited, but it must neither start nor end with a slash. Maximum of 5000 characters.",
				Computed:    true,
			},
			"path_time_zone": schema.StringAttribute{
				Description: "Timezone to use when rendering timestamps in paths.",
				Computed:    true,
			},
			"recurring_day": schema.Int64Attribute{
				Description: "If trigger type is `daily`, this specifies a day number to run in one of the supported intervals: `week`, `month`, `quarter`, `year`.",
				Computed:    true,
			},
			"retry_on_failure_interval_in_minutes": schema.Int64Attribute{
				Description: "If the Automation fails, retry at this interval (in minutes).  Acceptable values are 5 through 1440 (one day).  Set to null to disable.",
				Computed:    true,
			},
			"retry_on_failure_number_of_attempts": schema.Int64Attribute{
				Description: "If the Automation fails, retry at most this many times.  Maximum allowed value: 10.  Set to null to disable.",
				Computed:    true,
			},
			"schedule": schema.DynamicAttribute{
				Description: "If trigger is `custom_schedule`, Custom schedule description for when the automation should be run in json format.",
				Computed:    true,
			},
			"human_readable_schedule": schema.StringAttribute{
				Description: "If trigger is `custom_schedule`, Human readable Custom schedule description for when the automation should be run.",
				Computed:    true,
			},
			"schedule_days_of_week": schema.ListAttribute{
				Description: "If trigger is `custom_schedule`, Custom schedule description for when the automation should be run. 0-based days of the week. 0 is Sunday, 1 is Monday, etc.",
				Computed:    true,
				ElementType: types.Int64Type,
			},
			"schedule_times_of_day": schema.ListAttribute{
				Description: "If trigger is `custom_schedule`, Custom schedule description for when the automation should be run. Times of day in HH:MM format.",
				Computed:    true,
				ElementType: types.StringType,
			},
			"schedule_time_zone": schema.StringAttribute{
				Description: "If trigger is `custom_schedule`, Custom schedule Time Zone for when the automation should be run.",
				Computed:    true,
			},
			"source": schema.StringAttribute{
				Description: "Source path/glob.  See Automation docs for exact description, but this is used to filter for files in the `path` to find files to operate on. Supports globs, except on remote mounts.",
				Computed:    true,
			},
			"sync_ids": schema.ListAttribute{
				Description: "IDs of remote sync folder behaviors to run by this Automation",
				Computed:    true,
				ElementType: types.Int64Type,
			},
			"trigger_actions": schema.ListAttribute{
				Description: "If trigger is `action`, this is the list of action types on which to trigger the automation. Valid actions are create, read, update, destroy, move, copy",
				Computed:    true,
				ElementType: types.StringType,
			},
			"trigger": schema.StringAttribute{
				Description: "How this automation is triggered to run.",
				Computed:    true,
			},
			"user_id": schema.Int64Attribute{
				Description: "User ID of the Automation's creator.",
				Computed:    true,
			},
			"user_ids": schema.ListAttribute{
				Description: "IDs of Users for the Automation (i.e. who to Request File from)",
				Computed:    true,
				ElementType: types.Int64Type,
			},
			"value": schema.DynamicAttribute{
				Description: "A Hash of attributes specific to the automation type.",
				Computed:    true,
			},
			"webhook_url": schema.StringAttribute{
				Description: "If trigger is `webhook`, this is the URL of the webhook to trigger the Automation.",
				Computed:    true,
			},
		},
	}
}

func (r *automationDataSource) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
	var data automationDataSourceModel
	diags := req.Config.Get(ctx, &data)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	paramsAutomationFind := files_sdk.AutomationFindParams{}
	paramsAutomationFind.Id = data.Id.ValueInt64()

	automation, err := r.client.Find(paramsAutomationFind, files_sdk.WithContext(ctx))
	if err != nil {
		resp.Diagnostics.AddError(
			"Error Reading Files Automation",
			"Could not read automation id "+fmt.Sprint(data.Id.ValueInt64())+": "+err.Error(),
		)
		return
	}

	diags = r.populateDataSourceModel(ctx, automation, &data)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	diags = resp.State.Set(ctx, data)
	resp.Diagnostics.Append(diags...)
}

func (r *automationDataSource) populateDataSourceModel(ctx context.Context, automation files_sdk.Automation, state *automationDataSourceModel) (diags diag.Diagnostics) {
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
	state.ExcludePattern = types.StringValue(automation.ExcludePattern)
	state.ImportUrls, propDiags = lib.ToDynamic(ctx, path.Root("import_urls"), automation.ImportUrls, state.ImportUrls.UnderlyingValue())
	diags.Append(propDiags...)
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
	state.RetryOnFailureIntervalInMinutes = types.Int64Value(automation.RetryOnFailureIntervalInMinutes)
	state.RetryOnFailureNumberOfAttempts = types.Int64Value(automation.RetryOnFailureNumberOfAttempts)
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
