package provider

import (
	"context"
	"fmt"

	files_sdk "github.com/Files-com/files-sdk-go/v3"
	automation_version "github.com/Files-com/files-sdk-go/v3/automationversion"
	"github.com/Files-com/terraform-provider-files/lib"
	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/datasource/schema"
	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/path"

	"github.com/hashicorp/terraform-plugin-framework/types"
)

var (
	_ datasource.DataSource              = &automationVersionDataSource{}
	_ datasource.DataSourceWithConfigure = &automationVersionDataSource{}
)

func NewAutomationVersionDataSource() datasource.DataSource {
	return &automationVersionDataSource{}
}

type automationVersionDataSource struct {
	client *automation_version.Client
}

type automationVersionDataSourceModel struct {
	Id           types.Int64  `tfsdk:"id"`
	AutomationId types.Int64  `tfsdk:"automation_id"`
	WorkspaceId  types.Int64  `tfsdk:"workspace_id"`
	Version      types.Int64  `tfsdk:"version"`
	CreatedAt    types.String `tfsdk:"created_at"`
	UserId       types.Int64  `tfsdk:"user_id"`
	Username     types.String `tfsdk:"username"`
	ApiKeyId     types.Int64  `tfsdk:"api_key_id"`
	Definition   types.Object `tfsdk:"definition"`
}

func (r *automationVersionDataSource) Configure(_ context.Context, req datasource.ConfigureRequest, resp *datasource.ConfigureResponse) {
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

	r.client = &automation_version.Client{Config: sdk_config}
}

func (r *automationVersionDataSource) Metadata(_ context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_automation_version"
}

func (r *automationVersionDataSource) Schema(_ context.Context, _ datasource.SchemaRequest, resp *datasource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Description: "An AutomationVersion is an immutable revision of an Automation v2 graph. Each saved graph change creates a revision. Restoring a revision saves its definition as a new current revision and keeps existing runs and Run Automation nodes pinned to their selected revisions.",
		Attributes: map[string]schema.Attribute{
			"id": schema.Int64Attribute{
				Description: "Automation revision ID.",
				Required:    true,
			},
			"automation_id": schema.Int64Attribute{
				Description: "ID of the Automation that owns this revision.",
				Computed:    true,
			},
			"workspace_id": schema.Int64Attribute{
				Description: "Workspace that owns the Automation. Zero is the Default Workspace.",
				Computed:    true,
			},
			"version": schema.Int64Attribute{
				Description: "Revision number within the Automation.",
				Computed:    true,
			},
			"created_at": schema.StringAttribute{
				Description: "Time the revision was saved.",
				Computed:    true,
			},
			"user_id": schema.Int64Attribute{
				Description: "ID of the user who saved the revision, when available.",
				Computed:    true,
			},
			"username": schema.StringAttribute{
				Description: "Username recorded when the revision was saved, when available.",
				Computed:    true,
			},
			"api_key_id": schema.Int64Attribute{
				Description: "ID of the API key used to save the revision, when available.",
				Computed:    true,
			},
			"definition": schema.SingleNestedAttribute{
				Description: "Saved Automation v2 graph definition. Included when retrieving a single revision.",
				Computed:    true,
				Attributes: map[string]schema.Attribute{
					"schema_version": schema.Int64Attribute{
						Description: "Automation definition schema version.",
						Computed:    true,
					},
					"nodes": schema.ListNestedAttribute{
						Description: "Flat list of nodes. Every node requires a unique id and a type, accepts the node-specific config fields documented below, and action or control-flow nodes may set return to true to include their output in the Automation result.",
						Computed:    true,
						NestedObject: schema.NestedAttributeObject{Attributes: map[string]schema.Attribute{
							"trigger_scheduled": schema.SingleNestedAttribute{
								Computed: true,
								Attributes: map[string]schema.Attribute{
									"id": schema.StringAttribute{
										Description: "Unique node ID within the graph.",
										Computed:    true,
									},
									"config": schema.SingleNestedAttribute{
										Description: "Configuration fields for this node type.",
										Computed:    true,
										Attributes: map[string]schema.Attribute{
											"interval": schema.StringAttribute{
												Computed: true,
											},
											"recurring_day": schema.Int64Attribute{
												Computed: true,
											},
											"recurring_days": schema.ListAttribute{
												Computed:    true,
												ElementType: types.Int64Type,
											},
											"schedule_id": schema.Int64Attribute{
												Computed: true,
											},
											"schedule_days_of_week": schema.ListAttribute{
												Computed:    true,
												ElementType: types.Int64Type,
											},
											"schedule_interval_minutes": schema.Int64Attribute{
												Computed: true,
											},
											"schedule_times_of_day": schema.ListAttribute{
												Computed:    true,
												ElementType: types.StringType,
											},
											"schedule_time_zone": schema.StringAttribute{
												Computed: true,
											},
											"holiday_region": schema.StringAttribute{
												Computed: true,
											},
											"path": schema.StringAttribute{
												Computed: true,
											},
											"source": schema.StringAttribute{
												Description: "A string value. Supports {{ fts }} interpolation, evaluated per item at execution time. In path fields, {{ fts }} supplies one file or folder name; {{{ fts }}} supplies a path fragment.",
												Computed:    true,
											},
											"limit": schema.Int64Attribute{
												Computed: true,
											},
											"exclude_pattern": schema.StringAttribute{
												Computed: true,
											},
											"path_time_zone": schema.StringAttribute{
												Computed: true,
											},
										},
									},
								},
							},
							"trigger_manual": schema.SingleNestedAttribute{
								Computed: true,
								Attributes: map[string]schema.Attribute{
									"id": schema.StringAttribute{
										Description: "Unique node ID within the graph.",
										Computed:    true,
									},
									"config": schema.SingleNestedAttribute{
										Description: "Configuration fields for this node type.",
										Computed:    true,
										Attributes: map[string]schema.Attribute{
											"path": schema.StringAttribute{
												Computed: true,
											},
											"source": schema.StringAttribute{
												Description: "A string value. Supports {{ fts }} interpolation, evaluated per item at execution time. In path fields, {{ fts }} supplies one file or folder name; {{{ fts }}} supplies a path fragment.",
												Computed:    true,
											},
											"limit": schema.Int64Attribute{
												Computed: true,
											},
											"exclude_pattern": schema.StringAttribute{
												Computed: true,
											},
											"path_time_zone": schema.StringAttribute{
												Computed: true,
											},
										},
									},
								},
							},
							"trigger_shared": schema.SingleNestedAttribute{
								Computed: true,
								Attributes: map[string]schema.Attribute{
									"id": schema.StringAttribute{
										Description: "Unique node ID within the graph.",
										Computed:    true,
									},
									"config": schema.SingleNestedAttribute{
										Description: "Configuration fields for this node type.",
										Computed:    true,
										Attributes: map[string]schema.Attribute{
											"form_field_set_id": schema.Int64Attribute{
												Computed: true,
											},
										},
									},
								},
							},
							"trigger_action": schema.SingleNestedAttribute{
								Computed: true,
								Attributes: map[string]schema.Attribute{
									"id": schema.StringAttribute{
										Description: "Unique node ID within the graph.",
										Computed:    true,
									},
									"config": schema.SingleNestedAttribute{
										Description: "Configuration fields for this node type.",
										Computed:    true,
										Attributes: map[string]schema.Attribute{
											"kind": schema.StringAttribute{
												Computed: true,
											},
											"trigger_actions": schema.ListAttribute{
												Computed:    true,
												ElementType: types.StringType,
											},
											"path": schema.StringAttribute{
												Computed: true,
											},
											"source": schema.StringAttribute{
												Computed: true,
											},
											"exclude_pattern": schema.StringAttribute{
												Computed: true,
											},
											"path_time_zone": schema.StringAttribute{
												Computed: true,
											},
										},
									},
								},
							},
							"trigger_folder_creation": schema.SingleNestedAttribute{
								Computed: true,
								Attributes: map[string]schema.Attribute{
									"id": schema.StringAttribute{
										Description: "Unique node ID within the graph.",
										Computed:    true,
									},
									"config": schema.SingleNestedAttribute{
										Description: "Configuration fields for this node type.",
										Computed:    true,
										Attributes: map[string]schema.Attribute{
											"path": schema.StringAttribute{
												Computed: true,
											},
											"source": schema.StringAttribute{
												Computed: true,
											},
											"exclude_pattern": schema.StringAttribute{
												Computed: true,
											},
											"path_time_zone": schema.StringAttribute{
												Computed: true,
											},
										},
									},
								},
							},
							"trigger_webhook": schema.SingleNestedAttribute{
								Computed: true,
								Attributes: map[string]schema.Attribute{
									"id": schema.StringAttribute{
										Description: "Unique node ID within the graph.",
										Computed:    true,
									},
									"config": schema.SingleNestedAttribute{
										Description: "Configuration fields for this node type.",
										Computed:    true,
										Attributes: map[string]schema.Attribute{
											"path": schema.StringAttribute{
												Computed: true,
											},
											"source": schema.StringAttribute{
												Description: "A string value. Supports {{ fts }} interpolation, evaluated per item at execution time. In path fields, {{ fts }} supplies one file or folder name; {{{ fts }}} supplies a path fragment.",
												Computed:    true,
											},
											"limit": schema.Int64Attribute{
												Computed: true,
											},
											"exclude_pattern": schema.StringAttribute{
												Computed: true,
											},
											"path_time_zone": schema.StringAttribute{
												Computed: true,
											},
											"upload_file": schema.BoolAttribute{
												Computed: true,
											},
											"destination": schema.StringAttribute{
												Computed: true,
											},
											"bearer_token_sha256": schema.StringAttribute{
												Computed: true,
											},
											"allowed_ips": schema.ListAttribute{
												Computed:    true,
												ElementType: types.StringType,
											},
										},
									},
								},
							},
							"trigger_email": schema.SingleNestedAttribute{
								Computed: true,
								Attributes: map[string]schema.Attribute{
									"id": schema.StringAttribute{
										Description: "Unique node ID within the graph.",
										Computed:    true,
									},
									"config": schema.SingleNestedAttribute{
										Description: "Configuration fields for this node type.",
										Computed:    true,
										Attributes: map[string]schema.Attribute{
											"allowed_senders": schema.ListAttribute{
												Computed:    true,
												ElementType: types.StringType,
											},
											"save_body": schema.BoolAttribute{
												Computed: true,
											},
											"save_html_body": schema.BoolAttribute{
												Computed: true,
											},
											"save_raw_message": schema.BoolAttribute{
												Computed: true,
											},
											"source": schema.StringAttribute{
												Computed: true,
											},
											"limit": schema.Int64Attribute{
												Computed: true,
											},
											"exclude_pattern": schema.StringAttribute{
												Computed: true,
											},
											"path_time_zone": schema.StringAttribute{
												Computed: true,
											},
										},
									},
								},
							},
							"create_folder": schema.SingleNestedAttribute{
								Computed: true,
								Attributes: map[string]schema.Attribute{
									"id": schema.StringAttribute{
										Description: "Unique node ID within the graph.",
										Computed:    true,
									},
									"config": schema.SingleNestedAttribute{
										Description: "Configuration fields for this node type.",
										Computed:    true,
										Attributes: map[string]schema.Attribute{
											"destinations": schema.ListAttribute{
												Description: "One or more destination path templates.",
												Computed:    true,
												ElementType: types.StringType,
											},
											"on_error": schema.ListNestedAttribute{
												Description: "Error handlers. Each item selects an error pattern and whether to continue or propagate it.",
												Computed:    true,
												NestedObject: schema.NestedAttributeObject{Attributes: map[string]schema.Attribute{
													"error": schema.StringAttribute{
														Description: "A typed error family (not-found, processing-failure, service-unavailable), family/specific-error, or * for all.",
														Computed:    true,
													},
													"action": schema.StringAttribute{
														Computed: true,
													},
												},
												},
											},
										},
									},
									"return": schema.BoolAttribute{
										Description: "Whether this node's output is included in the Automation result.",
										Computed:    true,
									},
								},
							},
							"create_user": schema.SingleNestedAttribute{
								Computed: true,
								Attributes: map[string]schema.Attribute{
									"id": schema.StringAttribute{
										Description: "Unique node ID within the graph.",
										Computed:    true,
									},
									"config": schema.SingleNestedAttribute{
										Description: "Configuration fields for this node type.",
										Computed:    true,
										Attributes: map[string]schema.Attribute{
											"username": schema.StringAttribute{
												Description: "A string value. Supports {{ fts }} interpolation, evaluated per item at execution time. In path fields, {{ fts }} supplies one file or folder name; {{{ fts }}} supplies a path fragment.",
												Computed:    true,
											},
											"email": schema.StringAttribute{
												Description: "A string value. Supports {{ fts }} interpolation, evaluated per item at execution time. In path fields, {{ fts }} supplies one file or folder name; {{{ fts }}} supplies a path fragment.",
												Computed:    true,
											},
											"name": schema.StringAttribute{
												Description: "A string value. Supports {{ fts }} interpolation, evaluated per item at execution time. In path fields, {{ fts }} supplies one file or folder name; {{{ fts }}} supplies a path fragment.",
												Computed:    true,
											},
											"company": schema.StringAttribute{
												Description: "A string value. Supports {{ fts }} interpolation, evaluated per item at execution time. In path fields, {{ fts }} supplies one file or folder name; {{{ fts }}} supplies a path fragment.",
												Computed:    true,
											},
											"notes": schema.StringAttribute{
												Description: "A string value. Supports {{ fts }} interpolation, evaluated per item at execution time. In path fields, {{ fts }} supplies one file or folder name; {{{ fts }}} supplies a path fragment.",
												Computed:    true,
											},
											"on_error": schema.ListNestedAttribute{
												Description: "Error handlers. Each item selects an error pattern and whether to continue or propagate it.",
												Computed:    true,
												NestedObject: schema.NestedAttributeObject{Attributes: map[string]schema.Attribute{
													"error": schema.StringAttribute{
														Description: "A typed error family (not-found, processing-failure, service-unavailable), family/specific-error, or * for all.",
														Computed:    true,
													},
													"action": schema.StringAttribute{
														Computed: true,
													},
												},
												},
											},
										},
									},
									"return": schema.BoolAttribute{
										Description: "Whether this node's output is included in the Automation result.",
										Computed:    true,
									},
								},
							},
							"create_group": schema.SingleNestedAttribute{
								Computed: true,
								Attributes: map[string]schema.Attribute{
									"id": schema.StringAttribute{
										Description: "Unique node ID within the graph.",
										Computed:    true,
									},
									"config": schema.SingleNestedAttribute{
										Description: "Configuration fields for this node type.",
										Computed:    true,
										Attributes: map[string]schema.Attribute{
											"name": schema.StringAttribute{
												Description: "A string value. Supports {{ fts }} interpolation, evaluated per item at execution time. In path fields, {{ fts }} supplies one file or folder name; {{{ fts }}} supplies a path fragment.",
												Computed:    true,
											},
											"notes": schema.StringAttribute{
												Description: "A string value. Supports {{ fts }} interpolation, evaluated per item at execution time. In path fields, {{ fts }} supplies one file or folder name; {{{ fts }}} supplies a path fragment.",
												Computed:    true,
											},
											"on_error": schema.ListNestedAttribute{
												Description: "Error handlers. Each item selects an error pattern and whether to continue or propagate it.",
												Computed:    true,
												NestedObject: schema.NestedAttributeObject{Attributes: map[string]schema.Attribute{
													"error": schema.StringAttribute{
														Description: "A typed error family (not-found, processing-failure, service-unavailable), family/specific-error, or * for all.",
														Computed:    true,
													},
													"action": schema.StringAttribute{
														Computed: true,
													},
												},
												},
											},
										},
									},
									"return": schema.BoolAttribute{
										Description: "Whether this node's output is included in the Automation result.",
										Computed:    true,
									},
								},
							},
							"add_user_to_group": schema.SingleNestedAttribute{
								Computed: true,
								Attributes: map[string]schema.Attribute{
									"id": schema.StringAttribute{
										Description: "Unique node ID within the graph.",
										Computed:    true,
									},
									"config": schema.SingleNestedAttribute{
										Description: "Configuration fields for this node type.",
										Computed:    true,
										Attributes: map[string]schema.Attribute{
											"user_id": schema.Int64Attribute{
												Computed: true,
											},
											"user_id_expression": schema.StringAttribute{
												Description: "A bare Files TransformScript expression, no braces.",
												Computed:    true,
											},
											"group_id": schema.Int64Attribute{
												Computed: true,
											},
											"group_id_expression": schema.StringAttribute{
												Description: "A bare Files TransformScript expression, no braces.",
												Computed:    true,
											},
											"on_error": schema.ListNestedAttribute{
												Description: "Error handlers. Each item selects an error pattern and whether to continue or propagate it.",
												Computed:    true,
												NestedObject: schema.NestedAttributeObject{Attributes: map[string]schema.Attribute{
													"error": schema.StringAttribute{
														Description: "A typed error family (not-found, processing-failure, service-unavailable), family/specific-error, or * for all.",
														Computed:    true,
													},
													"action": schema.StringAttribute{
														Computed: true,
													},
												},
												},
											},
										},
									},
									"return": schema.BoolAttribute{
										Description: "Whether this node's output is included in the Automation result.",
										Computed:    true,
									},
								},
							},
							"create_permission": schema.SingleNestedAttribute{
								Computed: true,
								Attributes: map[string]schema.Attribute{
									"id": schema.StringAttribute{
										Description: "Unique node ID within the graph.",
										Computed:    true,
									},
									"config": schema.SingleNestedAttribute{
										Description: "Configuration fields for this node type.",
										Computed:    true,
										Attributes: map[string]schema.Attribute{
											"folder_path": schema.StringAttribute{
												Description: "A string value. Supports {{ fts }} interpolation, evaluated per item at execution time. In path fields, {{ fts }} supplies one file or folder name; {{{ fts }}} supplies a path fragment.",
												Computed:    true,
											},
											"user_id": schema.Int64Attribute{
												Computed: true,
											},
											"user_id_expression": schema.StringAttribute{
												Description: "A bare Files TransformScript expression, no braces.",
												Computed:    true,
											},
											"group_id": schema.Int64Attribute{
												Computed: true,
											},
											"group_id_expression": schema.StringAttribute{
												Description: "A bare Files TransformScript expression, no braces.",
												Computed:    true,
											},
											"permission": schema.StringAttribute{
												Computed: true,
											},
											"recursive": schema.BoolAttribute{
												Computed: true,
											},
											"on_error": schema.ListNestedAttribute{
												Description: "Error handlers. Each item selects an error pattern and whether to continue or propagate it.",
												Computed:    true,
												NestedObject: schema.NestedAttributeObject{Attributes: map[string]schema.Attribute{
													"error": schema.StringAttribute{
														Description: "A typed error family (not-found, processing-failure, service-unavailable), family/specific-error, or * for all.",
														Computed:    true,
													},
													"action": schema.StringAttribute{
														Computed: true,
													},
												},
												},
											},
										},
									},
									"return": schema.BoolAttribute{
										Description: "Whether this node's output is included in the Automation result.",
										Computed:    true,
									},
								},
							},
							"create_share_link": schema.SingleNestedAttribute{
								Computed: true,
								Attributes: map[string]schema.Attribute{
									"id": schema.StringAttribute{
										Description: "Unique node ID within the graph.",
										Computed:    true,
									},
									"config": schema.SingleNestedAttribute{
										Description: "Configuration fields for this node type.",
										Computed:    true,
										Attributes: map[string]schema.Attribute{
											"description": schema.StringAttribute{
												Description: "A string value. Supports {{ fts }} interpolation, evaluated per item at execution time. In path fields, {{ fts }} supplies one file or folder name; {{{ fts }}} supplies a path fragment.",
												Computed:    true,
											},
											"internal_name": schema.StringAttribute{
												Description: "A string value. Supports {{ fts }} interpolation, evaluated per item at execution time. In path fields, {{ fts }} supplies one file or folder name; {{{ fts }}} supplies a path fragment.",
												Computed:    true,
											},
											"note": schema.StringAttribute{
												Description: "A string value. Supports {{ fts }} interpolation, evaluated per item at execution time. In path fields, {{ fts }} supplies one file or folder name; {{{ fts }}} supplies a path fragment.",
												Computed:    true,
											},
											"permissions": schema.StringAttribute{
												Computed: true,
											},
											"expires_at": schema.StringAttribute{
												Description: "A string value. Supports {{ fts }} interpolation, evaluated per item at execution time. In path fields, {{ fts }} supplies one file or folder name; {{{ fts }}} supplies a path fragment.",
												Computed:    true,
											},
											"expires_in_minutes": schema.Int64Attribute{
												Computed: true,
											},
											"start_access_on_date": schema.StringAttribute{
												Description: "A string value. Supports {{ fts }} interpolation, evaluated per item at execution time. In path fields, {{ fts }} supplies one file or folder name; {{{ fts }}} supplies a path fragment.",
												Computed:    true,
											},
											"max_uses": schema.Int64Attribute{
												Computed: true,
											},
											"password_secret_id": schema.Int64Attribute{
												Computed: true,
											},
											"clickwrap_id": schema.Int64Attribute{
												Computed: true,
											},
											"require_registration": schema.BoolAttribute{
												Computed: true,
											},
											"require_share_recipient": schema.BoolAttribute{
												Computed: true,
											},
											"send_one_time_password_to_recipient_at_registration": schema.BoolAttribute{
												Computed: true,
											},
											"recipients": schema.ListAttribute{
												Computed:    true,
												ElementType: types.StringType,
											},
											"recipients_expression": schema.StringAttribute{
												Description: "A bare Files TransformScript expression, no braces.",
												Computed:    true,
											},
											"recipient_note": schema.StringAttribute{
												Description: "A string value. Supports {{ fts }} interpolation, evaluated per item at execution time. In path fields, {{ fts }} supplies one file or folder name; {{{ fts }}} supplies a path fragment.",
												Computed:    true,
											},
											"on_error": schema.ListNestedAttribute{
												Description: "Error handlers. Each item selects an error pattern and whether to continue or propagate it.",
												Computed:    true,
												NestedObject: schema.NestedAttributeObject{Attributes: map[string]schema.Attribute{
													"error": schema.StringAttribute{
														Description: "A typed error family (not-found, processing-failure, service-unavailable), family/specific-error, or * for all.",
														Computed:    true,
													},
													"action": schema.StringAttribute{
														Computed: true,
													},
												},
												},
											},
										},
									},
									"return": schema.BoolAttribute{
										Description: "Whether this node's output is included in the Automation result.",
										Computed:    true,
									},
								},
							},
							"copy_file": schema.SingleNestedAttribute{
								Computed: true,
								Attributes: map[string]schema.Attribute{
									"id": schema.StringAttribute{
										Description: "Unique node ID within the graph.",
										Computed:    true,
									},
									"config": schema.SingleNestedAttribute{
										Description: "Configuration fields for this node type.",
										Computed:    true,
										Attributes: map[string]schema.Attribute{
											"destinations": schema.ListAttribute{
												Description: "One or more destination path templates.",
												Computed:    true,
												ElementType: types.StringType,
											},
											"destination_replace_from": schema.StringAttribute{
												Computed: true,
											},
											"destination_replace_to": schema.StringAttribute{
												Computed: true,
											},
											"overwrite_files": schema.BoolAttribute{
												Computed: true,
											},
											"always_overwrite_size_matching_files": schema.BoolAttribute{
												Computed: true,
											},
											"flatten_destination_structure": schema.BoolAttribute{
												Computed: true,
											},
											"exclude_pattern": schema.StringAttribute{
												Computed: true,
											},
											"ignore_locked_folders": schema.BoolAttribute{
												Computed: true,
											},
											"on_error": schema.ListNestedAttribute{
												Description: "Error handlers. Each item selects an error pattern and whether to continue or propagate it.",
												Computed:    true,
												NestedObject: schema.NestedAttributeObject{Attributes: map[string]schema.Attribute{
													"error": schema.StringAttribute{
														Description: "A typed error family (not-found, processing-failure, service-unavailable), family/specific-error, or * for all.",
														Computed:    true,
													},
													"action": schema.StringAttribute{
														Computed: true,
													},
												},
												},
											},
										},
									},
									"return": schema.BoolAttribute{
										Description: "Whether this node's output is included in the Automation result.",
										Computed:    true,
									},
								},
							},
							"move_file": schema.SingleNestedAttribute{
								Computed: true,
								Attributes: map[string]schema.Attribute{
									"id": schema.StringAttribute{
										Description: "Unique node ID within the graph.",
										Computed:    true,
									},
									"config": schema.SingleNestedAttribute{
										Description: "Configuration fields for this node type.",
										Computed:    true,
										Attributes: map[string]schema.Attribute{
											"destinations": schema.ListAttribute{
												Description: "One or more destination path templates.",
												Computed:    true,
												ElementType: types.StringType,
											},
											"destination_replace_from": schema.StringAttribute{
												Computed: true,
											},
											"destination_replace_to": schema.StringAttribute{
												Computed: true,
											},
											"overwrite_files": schema.BoolAttribute{
												Computed: true,
											},
											"flatten_destination_structure": schema.BoolAttribute{
												Computed: true,
											},
											"exclude_pattern": schema.StringAttribute{
												Computed: true,
											},
											"ignore_locked_folders": schema.BoolAttribute{
												Computed: true,
											},
											"on_error": schema.ListNestedAttribute{
												Description: "Error handlers. Each item selects an error pattern and whether to continue or propagate it.",
												Computed:    true,
												NestedObject: schema.NestedAttributeObject{Attributes: map[string]schema.Attribute{
													"error": schema.StringAttribute{
														Description: "A typed error family (not-found, processing-failure, service-unavailable), family/specific-error, or * for all.",
														Computed:    true,
													},
													"action": schema.StringAttribute{
														Computed: true,
													},
												},
												},
											},
										},
									},
									"return": schema.BoolAttribute{
										Description: "Whether this node's output is included in the Automation result.",
										Computed:    true,
									},
								},
							},
							"delete_file": schema.SingleNestedAttribute{
								Computed: true,
								Attributes: map[string]schema.Attribute{
									"id": schema.StringAttribute{
										Description: "Unique node ID within the graph.",
										Computed:    true,
									},
									"config": schema.SingleNestedAttribute{
										Description: "Configuration fields for this node type.",
										Computed:    true,
										Attributes: map[string]schema.Attribute{
											"exclude_pattern": schema.StringAttribute{
												Computed: true,
											},
											"on_error": schema.ListNestedAttribute{
												Description: "Error handlers. Each item selects an error pattern and whether to continue or propagate it.",
												Computed:    true,
												NestedObject: schema.NestedAttributeObject{Attributes: map[string]schema.Attribute{
													"error": schema.StringAttribute{
														Description: "A typed error family (not-found, processing-failure, service-unavailable), family/specific-error, or * for all.",
														Computed:    true,
													},
													"action": schema.StringAttribute{
														Computed: true,
													},
												},
												},
											},
										},
									},
									"return": schema.BoolAttribute{
										Description: "Whether this node's output is included in the Automation result.",
										Computed:    true,
									},
								},
							},
							"import_file": schema.SingleNestedAttribute{
								Computed: true,
								Attributes: map[string]schema.Attribute{
									"id": schema.StringAttribute{
										Description: "Unique node ID within the graph.",
										Computed:    true,
									},
									"config": schema.SingleNestedAttribute{
										Description: "Configuration fields for this node type.",
										Computed:    true,
										Attributes: map[string]schema.Attribute{
											"url": schema.StringAttribute{
												Description: "A string value. Supports {{ fts }} interpolation, evaluated per item at execution time. In path fields, {{ fts }} supplies one file or folder name; {{{ fts }}} supplies a path fragment.",
												Computed:    true,
											},
											"destinations": schema.ListAttribute{
												Description: "One or more destination path templates.",
												Computed:    true,
												ElementType: types.StringType,
											},
											"method": schema.StringAttribute{
												Computed: true,
											},
											"headers": schema.MapAttribute{
												Computed:    true,
												ElementType: types.StringType,
											},
											"content": schema.StringAttribute{
												Description: "A string value. Supports {{ fts }} interpolation, evaluated per item at execution time. In path fields, {{ fts }} supplies one file or folder name; {{{ fts }}} supplies a path fragment.",
												Computed:    true,
											},
											"on_error": schema.ListNestedAttribute{
												Description: "Error handlers. Each item selects an error pattern and whether to continue or propagate it.",
												Computed:    true,
												NestedObject: schema.NestedAttributeObject{Attributes: map[string]schema.Attribute{
													"error": schema.StringAttribute{
														Description: "A typed error family (not-found, processing-failure, service-unavailable), family/specific-error, or * for all.",
														Computed:    true,
													},
													"action": schema.StringAttribute{
														Computed: true,
													},
												},
												},
											},
										},
									},
									"return": schema.BoolAttribute{
										Description: "Whether this node's output is included in the Automation result.",
										Computed:    true,
									},
								},
							},
							"import_into_database": schema.SingleNestedAttribute{
								Computed: true,
								Attributes: map[string]schema.Attribute{
									"id": schema.StringAttribute{
										Description: "Unique node ID within the graph.",
										Computed:    true,
									},
									"config": schema.SingleNestedAttribute{
										Description: "Configuration fields for this node type.",
										Computed:    true,
										Attributes: map[string]schema.Attribute{
											"database_destination_id": schema.Int64Attribute{
												Computed: true,
											},
											"rows": schema.StringAttribute{
												Description: "A bare Files TransformScript expression, no braces.",
												Computed:    true,
											},
											"on_error": schema.ListNestedAttribute{
												Description: "Error handlers. Each item selects an error pattern and whether to continue or propagate it.",
												Computed:    true,
												NestedObject: schema.NestedAttributeObject{Attributes: map[string]schema.Attribute{
													"error": schema.StringAttribute{
														Description: "A typed error family (not-found, processing-failure, service-unavailable), family/specific-error, or * for all.",
														Computed:    true,
													},
													"action": schema.StringAttribute{
														Computed: true,
													},
												},
												},
											},
										},
									},
									"return": schema.BoolAttribute{
										Description: "Whether this node's output is included in the Automation result.",
										Computed:    true,
									},
								},
							},
							"run_sync": schema.SingleNestedAttribute{
								Computed: true,
								Attributes: map[string]schema.Attribute{
									"id": schema.StringAttribute{
										Description: "Unique node ID within the graph.",
										Computed:    true,
									},
									"config": schema.SingleNestedAttribute{
										Description: "Configuration fields for this node type.",
										Computed:    true,
										Attributes: map[string]schema.Attribute{
											"sync_ids": schema.ListAttribute{
												Computed:    true,
												ElementType: types.Int64Type,
											},
											"on_error": schema.ListNestedAttribute{
												Description: "Error handlers. Each item selects an error pattern and whether to continue or propagate it.",
												Computed:    true,
												NestedObject: schema.NestedAttributeObject{Attributes: map[string]schema.Attribute{
													"error": schema.StringAttribute{
														Description: "A typed error family (not-found, processing-failure, service-unavailable), family/specific-error, or * for all.",
														Computed:    true,
													},
													"action": schema.StringAttribute{
														Computed: true,
													},
												},
												},
											},
										},
									},
									"return": schema.BoolAttribute{
										Description: "Whether this node's output is included in the Automation result.",
										Computed:    true,
									},
								},
							},
							"as2_send": schema.SingleNestedAttribute{
								Computed: true,
								Attributes: map[string]schema.Attribute{
									"id": schema.StringAttribute{
										Description: "Unique node ID within the graph.",
										Computed:    true,
									},
									"config": schema.SingleNestedAttribute{
										Description: "Configuration fields for this node type.",
										Computed:    true,
										Attributes: map[string]schema.Attribute{
											"as2_station_id": schema.Int64Attribute{
												Computed: true,
											},
											"as2_partner_id": schema.Int64Attribute{
												Computed: true,
											},
											"as2_subject": schema.StringAttribute{
												Computed: true,
											},
											"on_error": schema.ListNestedAttribute{
												Description: "Error handlers. Each item selects an error pattern and whether to continue or propagate it.",
												Computed:    true,
												NestedObject: schema.NestedAttributeObject{Attributes: map[string]schema.Attribute{
													"error": schema.StringAttribute{
														Description: "A typed error family (not-found, processing-failure, service-unavailable), family/specific-error, or * for all.",
														Computed:    true,
													},
													"action": schema.StringAttribute{
														Computed: true,
													},
												},
												},
											},
										},
									},
									"return": schema.BoolAttribute{
										Description: "Whether this node's output is included in the Automation result.",
										Computed:    true,
									},
								},
							},
							"send_email": schema.SingleNestedAttribute{
								Computed: true,
								Attributes: map[string]schema.Attribute{
									"id": schema.StringAttribute{
										Description: "Unique node ID within the graph.",
										Computed:    true,
									},
									"config": schema.SingleNestedAttribute{
										Description: "Configuration fields for this node type.",
										Computed:    true,
										Attributes: map[string]schema.Attribute{
											"from": schema.StringAttribute{
												Description: "A string value. Supports {{ fts }} interpolation, evaluated per item at execution time. In path fields, {{ fts }} supplies one file or folder name; {{{ fts }}} supplies a path fragment.",
												Computed:    true,
											},
											"to": schema.ListAttribute{
												Computed:    true,
												ElementType: types.StringType,
											},
											"cc": schema.ListAttribute{
												Computed:    true,
												ElementType: types.StringType,
											},
											"bcc": schema.ListAttribute{
												Computed:    true,
												ElementType: types.StringType,
											},
											"subject": schema.StringAttribute{
												Description: "A string value. Supports {{ fts }} interpolation, evaluated per item at execution time. In path fields, {{ fts }} supplies one file or folder name; {{{ fts }}} supplies a path fragment.",
												Computed:    true,
											},
											"body": schema.StringAttribute{
												Description: "A string value. Supports {{ fts }} interpolation, evaluated per item at execution time. In path fields, {{ fts }} supplies one file or folder name; {{{ fts }}} supplies a path fragment.",
												Computed:    true,
											},
											"attachments": schema.ListAttribute{
												Computed:    true,
												ElementType: types.StringType,
											},
											"on_error": schema.ListNestedAttribute{
												Description: "Error handlers. Each item selects an error pattern and whether to continue or propagate it.",
												Computed:    true,
												NestedObject: schema.NestedAttributeObject{Attributes: map[string]schema.Attribute{
													"error": schema.StringAttribute{
														Description: "A typed error family (not-found, processing-failure, service-unavailable), family/specific-error, or * for all.",
														Computed:    true,
													},
													"action": schema.StringAttribute{
														Computed: true,
													},
												},
												},
											},
										},
									},
									"return": schema.BoolAttribute{
										Description: "Whether this node's output is included in the Automation result.",
										Computed:    true,
									},
								},
							},
							"send_chat_message": schema.SingleNestedAttribute{
								Computed: true,
								Attributes: map[string]schema.Attribute{
									"id": schema.StringAttribute{
										Description: "Unique node ID within the graph.",
										Computed:    true,
									},
									"config": schema.SingleNestedAttribute{
										Description: "Configuration fields for this node type.",
										Computed:    true,
										Attributes: map[string]schema.Attribute{
											"provider": schema.StringAttribute{
												Computed: true,
											},
											"webhook_url_secret_id": schema.Int64Attribute{
												Computed: true,
											},
											"message": schema.StringAttribute{
												Description: "A string value. Supports {{ fts }} interpolation, evaluated per item at execution time. In path fields, {{ fts }} supplies one file or folder name; {{{ fts }}} supplies a path fragment.",
												Computed:    true,
											},
											"channel": schema.StringAttribute{
												Computed: true,
											},
											"username": schema.StringAttribute{
												Computed: true,
											},
											"icon_emoji": schema.StringAttribute{
												Computed: true,
											},
											"on_error": schema.ListNestedAttribute{
												Description: "Error handlers. Each item selects an error pattern and whether to continue or propagate it.",
												Computed:    true,
												NestedObject: schema.NestedAttributeObject{Attributes: map[string]schema.Attribute{
													"error": schema.StringAttribute{
														Description: "A typed error family (not-found, processing-failure, service-unavailable), family/specific-error, or * for all.",
														Computed:    true,
													},
													"action": schema.StringAttribute{
														Computed: true,
													},
												},
												},
											},
										},
									},
									"return": schema.BoolAttribute{
										Description: "Whether this node's output is included in the Automation result.",
										Computed:    true,
									},
								},
							},
							"send_webhook": schema.SingleNestedAttribute{
								Computed: true,
								Attributes: map[string]schema.Attribute{
									"id": schema.StringAttribute{
										Description: "Unique node ID within the graph.",
										Computed:    true,
									},
									"config": schema.SingleNestedAttribute{
										Description: "Configuration fields for this node type.",
										Computed:    true,
										Attributes: map[string]schema.Attribute{
											"url": schema.StringAttribute{
												Computed: true,
											},
											"url_secret_id": schema.Int64Attribute{
												Computed: true,
											},
											"headers_secret_id": schema.Int64Attribute{
												Computed: true,
											},
											"signing_secret_id": schema.Int64Attribute{
												Computed: true,
											},
											"payload": schema.StringAttribute{
												Description: "A bare Files TransformScript expression, no braces.",
												Computed:    true,
											},
											"use_dedicated_ips": schema.BoolAttribute{
												Computed: true,
											},
											"on_error": schema.ListNestedAttribute{
												Description: "Error handlers. Each item selects an error pattern and whether to continue or propagate it.",
												Computed:    true,
												NestedObject: schema.NestedAttributeObject{Attributes: map[string]schema.Attribute{
													"error": schema.StringAttribute{
														Description: "A typed error family (not-found, processing-failure, service-unavailable), family/specific-error, or * for all.",
														Computed:    true,
													},
													"action": schema.StringAttribute{
														Computed: true,
													},
												},
												},
											},
										},
									},
									"return": schema.BoolAttribute{
										Description: "Whether this node's output is included in the Automation result.",
										Computed:    true,
									},
								},
							},
							"http_request": schema.SingleNestedAttribute{
								Computed: true,
								Attributes: map[string]schema.Attribute{
									"id": schema.StringAttribute{
										Description: "Unique node ID within the graph.",
										Computed:    true,
									},
									"config": schema.SingleNestedAttribute{
										Description: "Configuration fields for this node type.",
										Computed:    true,
										Attributes: map[string]schema.Attribute{
											"method": schema.StringAttribute{
												Computed: true,
											},
											"url": schema.StringAttribute{
												Description: "A string value. Supports {{ fts }} interpolation, evaluated per item at execution time. In path fields, {{ fts }} supplies one file or folder name; {{{ fts }}} supplies a path fragment.",
												Computed:    true,
											},
											"query": schema.MapAttribute{
												Computed:    true,
												ElementType: types.StringType,
											},
											"headers": schema.MapAttribute{
												Computed:    true,
												ElementType: types.StringType,
											},
											"body": schema.StringAttribute{
												Description: "A string value. Supports {{ fts }} interpolation, evaluated per item at execution time. In path fields, {{ fts }} supplies one file or folder name; {{{ fts }}} supplies a path fragment.",
												Computed:    true,
											},
											"auth_secret_id": schema.Int64Attribute{
												Computed: true,
											},
											"on_error": schema.ListNestedAttribute{
												Description: "Error handlers. Each item selects an error pattern and whether to continue or propagate it.",
												Computed:    true,
												NestedObject: schema.NestedAttributeObject{Attributes: map[string]schema.Attribute{
													"error": schema.StringAttribute{
														Description: "A typed error family (not-found, processing-failure, service-unavailable), family/specific-error, or * for all.",
														Computed:    true,
													},
													"action": schema.StringAttribute{
														Computed: true,
													},
												},
												},
											},
										},
									},
									"return": schema.BoolAttribute{
										Description: "Whether this node's output is included in the Automation result.",
										Computed:    true,
									},
								},
							},
							"publish_amazon_sns": schema.SingleNestedAttribute{
								Computed: true,
								Attributes: map[string]schema.Attribute{
									"id": schema.StringAttribute{
										Description: "Unique node ID within the graph.",
										Computed:    true,
									},
									"config": schema.SingleNestedAttribute{
										Description: "Configuration fields for this node type.",
										Computed:    true,
										Attributes: map[string]schema.Attribute{
											"arn": schema.StringAttribute{
												Computed: true,
											},
											"message_group_id": schema.StringAttribute{
												Computed: true,
											},
											"aws_credentials_secret_id": schema.Int64Attribute{
												Computed: true,
											},
											"payload": schema.StringAttribute{
												Description: "A bare Files TransformScript expression, no braces.",
												Computed:    true,
											},
											"on_error": schema.ListNestedAttribute{
												Description: "Error handlers. Each item selects an error pattern and whether to continue or propagate it.",
												Computed:    true,
												NestedObject: schema.NestedAttributeObject{Attributes: map[string]schema.Attribute{
													"error": schema.StringAttribute{
														Description: "A typed error family (not-found, processing-failure, service-unavailable), family/specific-error, or * for all.",
														Computed:    true,
													},
													"action": schema.StringAttribute{
														Computed: true,
													},
												},
												},
											},
										},
									},
									"return": schema.BoolAttribute{
										Description: "Whether this node's output is included in the Automation result.",
										Computed:    true,
									},
								},
							},
							"publish_amazon_sqs": schema.SingleNestedAttribute{
								Computed: true,
								Attributes: map[string]schema.Attribute{
									"id": schema.StringAttribute{
										Description: "Unique node ID within the graph.",
										Computed:    true,
									},
									"config": schema.SingleNestedAttribute{
										Description: "Configuration fields for this node type.",
										Computed:    true,
										Attributes: map[string]schema.Attribute{
											"arn": schema.StringAttribute{
												Computed: true,
											},
											"message_group_id": schema.StringAttribute{
												Computed: true,
											},
											"aws_credentials_secret_id": schema.Int64Attribute{
												Computed: true,
											},
											"payload": schema.StringAttribute{
												Description: "A bare Files TransformScript expression, no braces.",
												Computed:    true,
											},
											"on_error": schema.ListNestedAttribute{
												Description: "Error handlers. Each item selects an error pattern and whether to continue or propagate it.",
												Computed:    true,
												NestedObject: schema.NestedAttributeObject{Attributes: map[string]schema.Attribute{
													"error": schema.StringAttribute{
														Description: "A typed error family (not-found, processing-failure, service-unavailable), family/specific-error, or * for all.",
														Computed:    true,
													},
													"action": schema.StringAttribute{
														Computed: true,
													},
												},
												},
											},
										},
									},
									"return": schema.BoolAttribute{
										Description: "Whether this node's output is included in the Automation result.",
										Computed:    true,
									},
								},
							},
							"publish_google_pubsub": schema.SingleNestedAttribute{
								Computed: true,
								Attributes: map[string]schema.Attribute{
									"id": schema.StringAttribute{
										Description: "Unique node ID within the graph.",
										Computed:    true,
									},
									"config": schema.SingleNestedAttribute{
										Description: "Configuration fields for this node type.",
										Computed:    true,
										Attributes: map[string]schema.Attribute{
											"project_id": schema.StringAttribute{
												Computed: true,
											},
											"topic_id": schema.StringAttribute{
												Computed: true,
											},
											"google_credentials_secret_id": schema.Int64Attribute{
												Computed: true,
											},
											"payload": schema.StringAttribute{
												Description: "A bare Files TransformScript expression, no braces.",
												Computed:    true,
											},
											"on_error": schema.ListNestedAttribute{
												Description: "Error handlers. Each item selects an error pattern and whether to continue or propagate it.",
												Computed:    true,
												NestedObject: schema.NestedAttributeObject{Attributes: map[string]schema.Attribute{
													"error": schema.StringAttribute{
														Description: "A typed error family (not-found, processing-failure, service-unavailable), family/specific-error, or * for all.",
														Computed:    true,
													},
													"action": schema.StringAttribute{
														Computed: true,
													},
												},
												},
											},
										},
									},
									"return": schema.BoolAttribute{
										Description: "Whether this node's output is included in the Automation result.",
										Computed:    true,
									},
								},
							},
							"publish_azure_service_bus": schema.SingleNestedAttribute{
								Computed: true,
								Attributes: map[string]schema.Attribute{
									"id": schema.StringAttribute{
										Description: "Unique node ID within the graph.",
										Computed:    true,
									},
									"config": schema.SingleNestedAttribute{
										Description: "Configuration fields for this node type.",
										Computed:    true,
										Attributes: map[string]schema.Attribute{
											"namespace": schema.StringAttribute{
												Computed: true,
											},
											"entity_path": schema.StringAttribute{
												Computed: true,
											},
											"azure_credentials_secret_id": schema.Int64Attribute{
												Computed: true,
											},
											"payload": schema.StringAttribute{
												Description: "A bare Files TransformScript expression, no braces.",
												Computed:    true,
											},
											"on_error": schema.ListNestedAttribute{
												Description: "Error handlers. Each item selects an error pattern and whether to continue or propagate it.",
												Computed:    true,
												NestedObject: schema.NestedAttributeObject{Attributes: map[string]schema.Attribute{
													"error": schema.StringAttribute{
														Description: "A typed error family (not-found, processing-failure, service-unavailable), family/specific-error, or * for all.",
														Computed:    true,
													},
													"action": schema.StringAttribute{
														Computed: true,
													},
												},
												},
											},
										},
									},
									"return": schema.BoolAttribute{
										Description: "Whether this node's output is included in the Automation result.",
										Computed:    true,
									},
								},
							},
							"agent_compute": schema.SingleNestedAttribute{
								Computed: true,
								Attributes: map[string]schema.Attribute{
									"id": schema.StringAttribute{
										Description: "Unique node ID within the graph.",
										Computed:    true,
									},
									"config": schema.SingleNestedAttribute{
										Description: "Configuration fields for this node type.",
										Computed:    true,
										Attributes: map[string]schema.Attribute{
											"remote_server_id": schema.Int64Attribute{
												Computed: true,
											},
											"command": schema.StringAttribute{
												Computed: true,
											},
											"arguments": schema.MapAttribute{
												Computed:    true,
												ElementType: types.StringType,
											},
											"on_error": schema.ListNestedAttribute{
												Description: "Error handlers. Each item selects an error pattern and whether to continue or propagate it.",
												Computed:    true,
												NestedObject: schema.NestedAttributeObject{Attributes: map[string]schema.Attribute{
													"error": schema.StringAttribute{
														Description: "A typed error family (not-found, processing-failure, service-unavailable), family/specific-error, or * for all.",
														Computed:    true,
													},
													"action": schema.StringAttribute{
														Computed: true,
													},
												},
												},
											},
										},
									},
									"return": schema.BoolAttribute{
										Description: "Whether this node's output is included in the Automation result.",
										Computed:    true,
									},
								},
							},
							"set_metadata": schema.SingleNestedAttribute{
								Computed: true,
								Attributes: map[string]schema.Attribute{
									"id": schema.StringAttribute{
										Description: "Unique node ID within the graph.",
										Computed:    true,
									},
									"config": schema.SingleNestedAttribute{
										Description: "Configuration fields for this node type.",
										Computed:    true,
										Attributes: map[string]schema.Attribute{
											"metadata": schema.MapAttribute{
												Computed:    true,
												ElementType: types.StringType,
											},
											"on_error": schema.ListNestedAttribute{
												Description: "Error handlers. Each item selects an error pattern and whether to continue or propagate it.",
												Computed:    true,
												NestedObject: schema.NestedAttributeObject{Attributes: map[string]schema.Attribute{
													"error": schema.StringAttribute{
														Description: "A typed error family (not-found, processing-failure, service-unavailable), family/specific-error, or * for all.",
														Computed:    true,
													},
													"action": schema.StringAttribute{
														Computed: true,
													},
												},
												},
											},
										},
									},
									"return": schema.BoolAttribute{
										Description: "Whether this node's output is included in the Automation result.",
										Computed:    true,
									},
								},
							},
							"extract": schema.SingleNestedAttribute{
								Computed: true,
								Attributes: map[string]schema.Attribute{
									"id": schema.StringAttribute{
										Description: "Unique node ID within the graph.",
										Computed:    true,
									},
									"config": schema.SingleNestedAttribute{
										Description: "Configuration fields for this node type.",
										Computed:    true,
										Attributes: map[string]schema.Attribute{
											"include_metadata": schema.BoolAttribute{
												Computed: true,
											},
											"include_content": schema.BoolAttribute{
												Computed: true,
											},
											"content_mode": schema.StringAttribute{
												Computed: true,
											},
											"max_chars": schema.Int64Attribute{
												Computed: true,
											},
											"max_pages": schema.Int64Attribute{
												Computed: true,
											},
											"on_error": schema.ListNestedAttribute{
												Description: "Error handlers. Each item selects an error pattern and whether to continue or propagate it.",
												Computed:    true,
												NestedObject: schema.NestedAttributeObject{Attributes: map[string]schema.Attribute{
													"error": schema.StringAttribute{
														Description: "A typed error family (not-found, processing-failure, service-unavailable), family/specific-error, or * for all.",
														Computed:    true,
													},
													"action": schema.StringAttribute{
														Computed: true,
													},
												},
												},
											},
										},
									},
									"return": schema.BoolAttribute{
										Description: "Whether this node's output is included in the Automation result.",
										Computed:    true,
									},
								},
							},
							"document_convert": schema.SingleNestedAttribute{
								Computed: true,
								Attributes: map[string]schema.Attribute{
									"id": schema.StringAttribute{
										Description: "Unique node ID within the graph.",
										Computed:    true,
									},
									"config": schema.SingleNestedAttribute{
										Description: "Configuration fields for this node type.",
										Computed:    true,
										Attributes: map[string]schema.Attribute{
											"target_format": schema.StringAttribute{
												Computed: true,
											},
											"destination": schema.StringAttribute{
												Description: "A string value. Supports {{ fts }} interpolation, evaluated per item at execution time. In path fields, {{ fts }} supplies one file or folder name; {{{ fts }}} supplies a path fragment.",
												Computed:    true,
											},
											"overwrite_files": schema.BoolAttribute{
												Computed: true,
											},
											"on_error": schema.ListNestedAttribute{
												Description: "Error handlers. Each item selects an error pattern and whether to continue or propagate it.",
												Computed:    true,
												NestedObject: schema.NestedAttributeObject{Attributes: map[string]schema.Attribute{
													"error": schema.StringAttribute{
														Description: "A typed error family (not-found, processing-failure, service-unavailable), family/specific-error, or * for all.",
														Computed:    true,
													},
													"action": schema.StringAttribute{
														Computed: true,
													},
												},
												},
											},
										},
									},
									"return": schema.BoolAttribute{
										Description: "Whether this node's output is included in the Automation result.",
										Computed:    true,
									},
								},
							},
							"image_convert": schema.SingleNestedAttribute{
								Computed: true,
								Attributes: map[string]schema.Attribute{
									"id": schema.StringAttribute{
										Description: "Unique node ID within the graph.",
										Computed:    true,
									},
									"config": schema.SingleNestedAttribute{
										Description: "Configuration fields for this node type.",
										Computed:    true,
										Attributes: map[string]schema.Attribute{
											"target_format": schema.StringAttribute{
												Computed: true,
											},
											"width": schema.Int64Attribute{
												Computed: true,
											},
											"height": schema.Int64Attribute{
												Computed: true,
											},
											"destination": schema.StringAttribute{
												Description: "A string value. Supports {{ fts }} interpolation, evaluated per item at execution time. In path fields, {{ fts }} supplies one file or folder name; {{{ fts }}} supplies a path fragment.",
												Computed:    true,
											},
											"overwrite_files": schema.BoolAttribute{
												Computed: true,
											},
											"on_error": schema.ListNestedAttribute{
												Description: "Error handlers. Each item selects an error pattern and whether to continue or propagate it.",
												Computed:    true,
												NestedObject: schema.NestedAttributeObject{Attributes: map[string]schema.Attribute{
													"error": schema.StringAttribute{
														Description: "A typed error family (not-found, processing-failure, service-unavailable), family/specific-error, or * for all.",
														Computed:    true,
													},
													"action": schema.StringAttribute{
														Computed: true,
													},
												},
												},
											},
										},
									},
									"return": schema.BoolAttribute{
										Description: "Whether this node's output is included in the Automation result.",
										Computed:    true,
									},
								},
							},
							"zip": schema.SingleNestedAttribute{
								Computed: true,
								Attributes: map[string]schema.Attribute{
									"id": schema.StringAttribute{
										Description: "Unique node ID within the graph.",
										Computed:    true,
									},
									"config": schema.SingleNestedAttribute{
										Description: "Configuration fields for this node type.",
										Computed:    true,
										Attributes: map[string]schema.Attribute{
											"destination": schema.StringAttribute{
												Description: "A string value. Supports {{ fts }} interpolation, evaluated per item at execution time. In path fields, {{ fts }} supplies one file or folder name; {{{ fts }}} supplies a path fragment.",
												Computed:    true,
											},
											"overwrite_files": schema.BoolAttribute{
												Computed: true,
											},
											"archive_format": schema.StringAttribute{
												Computed: true,
											},
											"on_error": schema.ListNestedAttribute{
												Description: "Error handlers. Each item selects an error pattern and whether to continue or propagate it.",
												Computed:    true,
												NestedObject: schema.NestedAttributeObject{Attributes: map[string]schema.Attribute{
													"error": schema.StringAttribute{
														Description: "A typed error family (not-found, processing-failure, service-unavailable), family/specific-error, or * for all.",
														Computed:    true,
													},
													"action": schema.StringAttribute{
														Computed: true,
													},
												},
												},
											},
										},
									},
									"return": schema.BoolAttribute{
										Description: "Whether this node's output is included in the Automation result.",
										Computed:    true,
									},
								},
							},
							"unzip": schema.SingleNestedAttribute{
								Computed: true,
								Attributes: map[string]schema.Attribute{
									"id": schema.StringAttribute{
										Description: "Unique node ID within the graph.",
										Computed:    true,
									},
									"config": schema.SingleNestedAttribute{
										Description: "Configuration fields for this node type.",
										Computed:    true,
										Attributes: map[string]schema.Attribute{
											"destination": schema.StringAttribute{
												Description: "A string value. Supports {{ fts }} interpolation, evaluated per item at execution time. In path fields, {{ fts }} supplies one file or folder name; {{{ fts }}} supplies a path fragment.",
												Computed:    true,
											},
											"overwrite_files": schema.BoolAttribute{
												Computed: true,
											},
											"on_error": schema.ListNestedAttribute{
												Description: "Error handlers. Each item selects an error pattern and whether to continue or propagate it.",
												Computed:    true,
												NestedObject: schema.NestedAttributeObject{Attributes: map[string]schema.Attribute{
													"error": schema.StringAttribute{
														Description: "A typed error family (not-found, processing-failure, service-unavailable), family/specific-error, or * for all.",
														Computed:    true,
													},
													"action": schema.StringAttribute{
														Computed: true,
													},
												},
												},
											},
										},
									},
									"return": schema.BoolAttribute{
										Description: "Whether this node's output is included in the Automation result.",
										Computed:    true,
									},
								},
							},
							"gpg_encrypt": schema.SingleNestedAttribute{
								Computed: true,
								Attributes: map[string]schema.Attribute{
									"id": schema.StringAttribute{
										Description: "Unique node ID within the graph.",
										Computed:    true,
									},
									"config": schema.SingleNestedAttribute{
										Description: "Configuration fields for this node type.",
										Computed:    true,
										Attributes: map[string]schema.Attribute{
											"gpg_key_id": schema.Int64Attribute{
												Computed: true,
											},
											"destination": schema.StringAttribute{
												Description: "A string value. Supports {{ fts }} interpolation, evaluated per item at execution time. In path fields, {{ fts }} supplies one file or folder name; {{{ fts }}} supplies a path fragment.",
												Computed:    true,
											},
											"overwrite_files": schema.BoolAttribute{
												Computed: true,
											},
											"on_error": schema.ListNestedAttribute{
												Description: "Error handlers. Each item selects an error pattern and whether to continue or propagate it.",
												Computed:    true,
												NestedObject: schema.NestedAttributeObject{Attributes: map[string]schema.Attribute{
													"error": schema.StringAttribute{
														Description: "A typed error family (not-found, processing-failure, service-unavailable), family/specific-error, or * for all.",
														Computed:    true,
													},
													"action": schema.StringAttribute{
														Computed: true,
													},
												},
												},
											},
										},
									},
									"return": schema.BoolAttribute{
										Description: "Whether this node's output is included in the Automation result.",
										Computed:    true,
									},
								},
							},
							"gpg_decrypt": schema.SingleNestedAttribute{
								Computed: true,
								Attributes: map[string]schema.Attribute{
									"id": schema.StringAttribute{
										Description: "Unique node ID within the graph.",
										Computed:    true,
									},
									"config": schema.SingleNestedAttribute{
										Description: "Configuration fields for this node type.",
										Computed:    true,
										Attributes: map[string]schema.Attribute{
											"gpg_key_id": schema.Int64Attribute{
												Computed: true,
											},
											"destination": schema.StringAttribute{
												Description: "A string value. Supports {{ fts }} interpolation, evaluated per item at execution time. In path fields, {{ fts }} supplies one file or folder name; {{{ fts }}} supplies a path fragment.",
												Computed:    true,
											},
											"overwrite_files": schema.BoolAttribute{
												Computed: true,
											},
											"on_error": schema.ListNestedAttribute{
												Description: "Error handlers. Each item selects an error pattern and whether to continue or propagate it.",
												Computed:    true,
												NestedObject: schema.NestedAttributeObject{Attributes: map[string]schema.Attribute{
													"error": schema.StringAttribute{
														Description: "A typed error family (not-found, processing-failure, service-unavailable), family/specific-error, or * for all.",
														Computed:    true,
													},
													"action": schema.StringAttribute{
														Computed: true,
													},
												},
												},
											},
										},
									},
									"return": schema.BoolAttribute{
										Description: "Whether this node's output is included in the Automation result.",
										Computed:    true,
									},
								},
							},
							"if": schema.SingleNestedAttribute{
								Computed: true,
								Attributes: map[string]schema.Attribute{
									"id": schema.StringAttribute{
										Description: "Unique node ID within the graph.",
										Computed:    true,
									},
									"config": schema.SingleNestedAttribute{
										Description: "Configuration fields for this node type.",
										Computed:    true,
										Attributes: map[string]schema.Attribute{
											"condition": schema.StringAttribute{
												Description: "A bare Files TransformScript expression, no braces.",
												Computed:    true,
											},
											"on_error": schema.ListNestedAttribute{
												Description: "Error handlers. Each item selects an error pattern and whether to continue or propagate it.",
												Computed:    true,
												NestedObject: schema.NestedAttributeObject{Attributes: map[string]schema.Attribute{
													"error": schema.StringAttribute{
														Description: "A typed error family (not-found, processing-failure, service-unavailable), family/specific-error, or * for all.",
														Computed:    true,
													},
													"action": schema.StringAttribute{
														Computed: true,
													},
												},
												},
											},
										},
									},
									"return": schema.BoolAttribute{
										Description: "Whether this node's output is included in the Automation result.",
										Computed:    true,
									},
								},
							},
							"switch": schema.SingleNestedAttribute{
								Computed: true,
								Attributes: map[string]schema.Attribute{
									"id": schema.StringAttribute{
										Description: "Unique node ID within the graph.",
										Computed:    true,
									},
									"config": schema.SingleNestedAttribute{
										Description: "Configuration fields for this node type.",
										Computed:    true,
										Attributes: map[string]schema.Attribute{
											"cases": schema.ListNestedAttribute{
												Description: "Switch cases. Each item contains a unique output name and the condition that sends items to it.",
												Computed:    true,
												NestedObject: schema.NestedAttributeObject{Attributes: map[string]schema.Attribute{
													"name": schema.StringAttribute{
														Description: "Case-sensitive output name. Use 1-64 letters, numbers, underscores, or hyphens.",
														Computed:    true,
													},
													"condition": schema.StringAttribute{
														Description: "A bare Files TransformScript expression, no braces.",
														Computed:    true,
													},
												},
												},
											},
											"on_error": schema.ListNestedAttribute{
												Description: "Error handlers. Each item selects an error pattern and whether to continue or propagate it.",
												Computed:    true,
												NestedObject: schema.NestedAttributeObject{Attributes: map[string]schema.Attribute{
													"error": schema.StringAttribute{
														Description: "A typed error family (not-found, processing-failure, service-unavailable), family/specific-error, or * for all.",
														Computed:    true,
													},
													"action": schema.StringAttribute{
														Computed: true,
													},
												},
												},
											},
										},
									},
									"return": schema.BoolAttribute{
										Description: "Whether this node's output is included in the Automation result.",
										Computed:    true,
									},
								},
							},
							"filter": schema.SingleNestedAttribute{
								Computed: true,
								Attributes: map[string]schema.Attribute{
									"id": schema.StringAttribute{
										Description: "Unique node ID within the graph.",
										Computed:    true,
									},
									"config": schema.SingleNestedAttribute{
										Description: "Configuration fields for this node type.",
										Computed:    true,
										Attributes: map[string]schema.Attribute{
											"kind": schema.StringAttribute{
												Computed: true,
											},
											"condition": schema.StringAttribute{
												Description: "A bare Files TransformScript expression, no braces.",
												Computed:    true,
											},
											"on_error": schema.ListNestedAttribute{
												Description: "Error handlers. Each item selects an error pattern and whether to continue or propagate it.",
												Computed:    true,
												NestedObject: schema.NestedAttributeObject{Attributes: map[string]schema.Attribute{
													"error": schema.StringAttribute{
														Description: "A typed error family (not-found, processing-failure, service-unavailable), family/specific-error, or * for all.",
														Computed:    true,
													},
													"action": schema.StringAttribute{
														Computed: true,
													},
												},
												},
											},
										},
									},
									"return": schema.BoolAttribute{
										Description: "Whether this node's output is included in the Automation result.",
										Computed:    true,
									},
								},
							},
							"for_each": schema.SingleNestedAttribute{
								Computed: true,
								Attributes: map[string]schema.Attribute{
									"id": schema.StringAttribute{
										Description: "Unique node ID within the graph.",
										Computed:    true,
									},
									"config": schema.SingleNestedAttribute{
										Description: "Configuration fields for this node type.",
										Computed:    true,
										Attributes: map[string]schema.Attribute{
											"list": schema.StringAttribute{
												Description: "A bare Files TransformScript expression, no braces.",
												Computed:    true,
											},
											"on_error": schema.ListNestedAttribute{
												Description: "Error handlers. Each item selects an error pattern and whether to continue or propagate it.",
												Computed:    true,
												NestedObject: schema.NestedAttributeObject{Attributes: map[string]schema.Attribute{
													"error": schema.StringAttribute{
														Description: "A typed error family (not-found, processing-failure, service-unavailable), family/specific-error, or * for all.",
														Computed:    true,
													},
													"action": schema.StringAttribute{
														Computed: true,
													},
												},
												},
											},
										},
									},
									"return": schema.BoolAttribute{
										Description: "Whether this node's output is included in the Automation result.",
										Computed:    true,
									},
								},
							},
							"join": schema.SingleNestedAttribute{
								Computed: true,
								Attributes: map[string]schema.Attribute{
									"id": schema.StringAttribute{
										Description: "Unique node ID within the graph.",
										Computed:    true,
									},
									"config": schema.SingleNestedAttribute{
										Description: "Configuration fields for this node type.",
										Computed:    true,
										Attributes: map[string]schema.Attribute{
											"join_type": schema.StringAttribute{
												Computed: true,
											},
											"left_key": schema.StringAttribute{
												Description: "A bare Files TransformScript expression, no braces.",
												Computed:    true,
											},
											"right_key": schema.StringAttribute{
												Description: "A bare Files TransformScript expression, no braces.",
												Computed:    true,
											},
											"on_error": schema.ListNestedAttribute{
												Description: "Error handlers. Each item selects an error pattern and whether to continue or propagate it.",
												Computed:    true,
												NestedObject: schema.NestedAttributeObject{Attributes: map[string]schema.Attribute{
													"error": schema.StringAttribute{
														Description: "A typed error family (not-found, processing-failure, service-unavailable), family/specific-error, or * for all.",
														Computed:    true,
													},
													"action": schema.StringAttribute{
														Computed: true,
													},
												},
												},
											},
										},
									},
									"return": schema.BoolAttribute{
										Description: "Whether this node's output is included in the Automation result.",
										Computed:    true,
									},
								},
							},
							"aggregate": schema.SingleNestedAttribute{
								Computed: true,
								Attributes: map[string]schema.Attribute{
									"id": schema.StringAttribute{
										Description: "Unique node ID within the graph.",
										Computed:    true,
									},
									"config": schema.SingleNestedAttribute{
										Description: "Configuration fields for this node type.",
										Computed:    true,
										Attributes: map[string]schema.Attribute{
											"on_error": schema.ListNestedAttribute{
												Description: "Error handlers. Each item selects an error pattern and whether to continue or propagate it.",
												Computed:    true,
												NestedObject: schema.NestedAttributeObject{Attributes: map[string]schema.Attribute{
													"error": schema.StringAttribute{
														Description: "A typed error family (not-found, processing-failure, service-unavailable), family/specific-error, or * for all.",
														Computed:    true,
													},
													"action": schema.StringAttribute{
														Computed: true,
													},
												},
												},
											},
										},
									},
									"return": schema.BoolAttribute{
										Description: "Whether this node's output is included in the Automation result.",
										Computed:    true,
									},
								},
							},
							"wait": schema.SingleNestedAttribute{
								Computed: true,
								Attributes: map[string]schema.Attribute{
									"id": schema.StringAttribute{
										Description: "Unique node ID within the graph.",
										Computed:    true,
									},
									"config": schema.SingleNestedAttribute{
										Description: "Configuration fields for this node type.",
										Computed:    true,
										Attributes: map[string]schema.Attribute{
											"delay_seconds": schema.Int64Attribute{
												Computed: true,
											},
											"on_error": schema.ListNestedAttribute{
												Description: "Error handlers. Each item selects an error pattern and whether to continue or propagate it.",
												Computed:    true,
												NestedObject: schema.NestedAttributeObject{Attributes: map[string]schema.Attribute{
													"error": schema.StringAttribute{
														Description: "A typed error family (not-found, processing-failure, service-unavailable), family/specific-error, or * for all.",
														Computed:    true,
													},
													"action": schema.StringAttribute{
														Computed: true,
													},
												},
												},
											},
										},
									},
									"return": schema.BoolAttribute{
										Description: "Whether this node's output is included in the Automation result.",
										Computed:    true,
									},
								},
							},
							"transform": schema.SingleNestedAttribute{
								Computed: true,
								Attributes: map[string]schema.Attribute{
									"id": schema.StringAttribute{
										Description: "Unique node ID within the graph.",
										Computed:    true,
									},
									"config": schema.SingleNestedAttribute{
										Description: "Configuration fields for this node type.",
										Computed:    true,
										Attributes: map[string]schema.Attribute{
											"fts": schema.StringAttribute{
												Computed: true,
											},
											"destination": schema.StringAttribute{
												Description: "A string value. Supports {{ fts }} interpolation, evaluated per item at execution time. In path fields, {{ fts }} supplies one file or folder name; {{{ fts }}} supplies a path fragment.",
												Computed:    true,
											},
											"overwrite_files": schema.BoolAttribute{
												Computed: true,
											},
											"on_error": schema.ListNestedAttribute{
												Description: "Error handlers. Each item selects an error pattern and whether to continue or propagate it.",
												Computed:    true,
												NestedObject: schema.NestedAttributeObject{Attributes: map[string]schema.Attribute{
													"error": schema.StringAttribute{
														Description: "A typed error family (not-found, processing-failure, service-unavailable), family/specific-error, or * for all.",
														Computed:    true,
													},
													"action": schema.StringAttribute{
														Computed: true,
													},
												},
												},
											},
										},
									},
									"return": schema.BoolAttribute{
										Description: "Whether this node's output is included in the Automation result.",
										Computed:    true,
									},
								},
							},
							"run_automation": schema.SingleNestedAttribute{
								Computed: true,
								Attributes: map[string]schema.Attribute{
									"id": schema.StringAttribute{
										Description: "Unique node ID within the graph.",
										Computed:    true,
									},
									"config": schema.SingleNestedAttribute{
										Description: "Configuration fields for this node type.",
										Computed:    true,
										Attributes: map[string]schema.Attribute{
											"automation_id": schema.Int64Attribute{
												Computed: true,
											},
											"automation_version": schema.Int64Attribute{
												Computed: true,
											},
											"on_error": schema.ListNestedAttribute{
												Description: "Error handlers. Each item selects an error pattern and whether to continue or propagate it.",
												Computed:    true,
												NestedObject: schema.NestedAttributeObject{Attributes: map[string]schema.Attribute{
													"error": schema.StringAttribute{
														Description: "A typed error family (not-found, processing-failure, service-unavailable), family/specific-error, or * for all.",
														Computed:    true,
													},
													"action": schema.StringAttribute{
														Computed: true,
													},
												},
												},
											},
										},
									},
									"return": schema.BoolAttribute{
										Description: "Whether this node's output is included in the Automation result.",
										Computed:    true,
									},
								},
							},
						},
						},
					},
					"edges": schema.ListNestedAttribute{
						Description: "Directed connections between nodes.",
						Computed:    true,
						NestedObject: schema.NestedAttributeObject{Attributes: map[string]schema.Attribute{
							"from": schema.StringAttribute{
								Description: "Source node ID.",
								Computed:    true,
							},
							"to": schema.StringAttribute{
								Description: "Destination node ID.",
								Computed:    true,
							},
							"output": schema.StringAttribute{
								Description: "Named source output. Omit for the main output.",
								Computed:    true,
							},
							"input": schema.StringAttribute{
								Description: "Named destination input. Omit for the main input.",
								Computed:    true,
							},
						},
						},
					},
				},
			},
		},
	}
}

func (r *automationVersionDataSource) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
	var data automationVersionDataSourceModel
	diags := req.Config.Get(ctx, &data)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	paramsAutomationVersionFind := files_sdk.AutomationVersionFindParams{}
	paramsAutomationVersionFind.Id = data.Id.ValueInt64()

	automationVersion, err := r.client.Find(paramsAutomationVersionFind, files_sdk.WithContext(ctx))
	if err != nil {
		resp.Diagnostics.AddError(
			"Error Reading Files AutomationVersion",
			"Could not read automation_version id "+fmt.Sprint(data.Id.ValueInt64())+": "+err.Error(),
		)
		return
	}

	diags = r.populateDataSourceModel(ctx, automationVersion, &data)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	diags = resp.State.Set(ctx, data)
	resp.Diagnostics.Append(diags...)
}

func (r *automationVersionDataSource) populateDataSourceModel(ctx context.Context, automationVersion files_sdk.AutomationVersion, state *automationVersionDataSourceModel) (diags diag.Diagnostics) {
	var propDiags diag.Diagnostics

	state.Id = types.Int64Value(automationVersion.Id)
	state.AutomationId = types.Int64Value(automationVersion.AutomationId)
	state.WorkspaceId = types.Int64Value(automationVersion.WorkspaceId)
	state.Version = types.Int64Value(automationVersion.Version)
	if err := lib.TimeToStringType(ctx, path.Root("created_at"), automationVersion.CreatedAt, &state.CreatedAt); err != nil {
		diags.AddError(
			"Error Creating Files AutomationVersion",
			"Could not convert state created_at to string: "+err.Error(),
		)
	}
	state.UserId = types.Int64Value(automationVersion.UserId)
	state.Username = types.StringValue(automationVersion.Username)
	state.ApiKeyId = types.Int64Value(automationVersion.ApiKeyId)
	definitionValue := interface{}(automationVersion.Definition)
	definitionValue, definitionValueTransformDiags0 := lib.WrapDiscriminatedUnionAtPath(ctx, path.Root("definition"), definitionValue, []string{"nodes"}, "type", []lib.JSONSchemaVariant{{Name: "trigger_scheduled", Value: "trigger_scheduled"}, {Name: "trigger_manual", Value: "trigger_manual"}, {Name: "trigger_shared", Value: "trigger_shared"}, {Name: "trigger_action", Value: "trigger_action"}, {Name: "trigger_folder_creation", Value: "trigger_folder_creation"}, {Name: "trigger_webhook", Value: "trigger_webhook"}, {Name: "trigger_email", Value: "trigger_email"}, {Name: "create_folder", Value: "create_folder"}, {Name: "create_user", Value: "create_user"}, {Name: "create_group", Value: "create_group"}, {Name: "add_user_to_group", Value: "add_user_to_group"}, {Name: "create_permission", Value: "create_permission"}, {Name: "create_share_link", Value: "create_share_link"}, {Name: "copy_file", Value: "copy_file"}, {Name: "move_file", Value: "move_file"}, {Name: "delete_file", Value: "delete_file"}, {Name: "import_file", Value: "import_file"}, {Name: "import_into_database", Value: "import_into_database"}, {Name: "run_sync", Value: "run_sync"}, {Name: "as2_send", Value: "as2_send"}, {Name: "send_email", Value: "send_email"}, {Name: "send_chat_message", Value: "send_chat_message"}, {Name: "send_webhook", Value: "send_webhook"}, {Name: "http_request", Value: "http_request"}, {Name: "publish_amazon_sns", Value: "publish_amazon_sns"}, {Name: "publish_amazon_sqs", Value: "publish_amazon_sqs"}, {Name: "publish_google_pubsub", Value: "publish_google_pubsub"}, {Name: "publish_azure_service_bus", Value: "publish_azure_service_bus"}, {Name: "agent_compute", Value: "agent_compute"}, {Name: "set_metadata", Value: "set_metadata"}, {Name: "extract", Value: "extract"}, {Name: "document_convert", Value: "document_convert"}, {Name: "image_convert", Value: "image_convert"}, {Name: "zip", Value: "zip"}, {Name: "unzip", Value: "unzip"}, {Name: "gpg_encrypt", Value: "gpg_encrypt"}, {Name: "gpg_decrypt", Value: "gpg_decrypt"}, {Name: "if", Value: "if"}, {Name: "switch", Value: "switch"}, {Name: "filter", Value: "filter"}, {Name: "for_each", Value: "for_each"}, {Name: "join", Value: "join"}, {Name: "aggregate", Value: "aggregate"}, {Name: "wait", Value: "wait"}, {Name: "transform", Value: "transform"}, {Name: "run_automation", Value: "run_automation"}})
	diags.Append(definitionValueTransformDiags0...)
	state.Definition, propDiags = lib.ToObject(ctx, path.Root("definition"), definitionValue, state.Definition)
	diags.Append(propDiags...)

	return
}
