package provider

import (
	"context"
	"fmt"

	files_sdk "github.com/Files-com/files-sdk-go/v3"
	permission "github.com/Files-com/files-sdk-go/v3/permission"
	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/datasource/schema"
	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

var (
	_ datasource.DataSource              = &permissionDataSource{}
	_ datasource.DataSourceWithConfigure = &permissionDataSource{}
)

func NewPermissionDataSource() datasource.DataSource {
	return &permissionDataSource{}
}

type permissionDataSource struct {
	client *permission.Client
}

type permissionDataSourceModel struct {
	Id         types.Int64  `tfsdk:"id"`
	Path       types.String `tfsdk:"path"`
	UserId     types.Int64  `tfsdk:"user_id"`
	Username   types.String `tfsdk:"username"`
	GroupId    types.Int64  `tfsdk:"group_id"`
	GroupName  types.String `tfsdk:"group_name"`
	Permission types.String `tfsdk:"permission"`
	Recursive  types.Bool   `tfsdk:"recursive"`
}

func (r *permissionDataSource) Configure(_ context.Context, req datasource.ConfigureRequest, resp *datasource.ConfigureResponse) {
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

	r.client = &permission.Client{Config: sdk_config}
}

func (r *permissionDataSource) Metadata(_ context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_permission"
}

func (r *permissionDataSource) Schema(_ context.Context, _ datasource.SchemaRequest, resp *datasource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Description: "A Permission object represents a grant of access permission on a specific Path to a User or Group.\n\n\n\nThey can be optionally recursive or nonrecursive into the subfolders of that path.\n\n\n\nA Permission may be applied to a User *or* a Group, but not both at once.\n\n\n\nThe following table sets forth the available Permission types:\n\n\n\n| Permission | Access Level Granted | Automatically Also Includes/Implies Permissions |\n\n| --- | ----------- | --------------------- |\n\n| `admin` | Able to manage Folder Behaviors, Permissions, and Notifications for the folder. Also grants all other permissions. | `bundle`, `full`, `writeonly`, `readonly`, `list`, `history` |\n\n| `full` | Able to read, write, move, delete, and rename files and folders. Also grants the ability to overwrite files upon upload. | `writeonly`, `readonly`, `list` |\n\n| `readonly` | Able to list, preview, and download files and folders. | `list` |\n\n| `writeonly` | Able to upload files, create folders and list subfolders the user has write permission to. | none |\n\n| `list` | Able to list files and folders, but not download. | none |\n\n| `bundle` | Able to share files and folders via a Bundle (share link). | `readonly`, `list` |\n\n| `history` | Able to view the history of files and folders and to create email notifications for themselves. | `list` |",
		Attributes: map[string]schema.Attribute{
			"id": schema.Int64Attribute{
				Description: "Permission ID",
				Required:    true,
			},
			"path": schema.StringAttribute{
				Description: "Path. This must be slash-delimited, but it must neither start nor end with a slash. Maximum of 5000 characters.",
				Computed:    true,
			},
			"user_id": schema.Int64Attribute{
				Description: "User ID",
				Computed:    true,
			},
			"username": schema.StringAttribute{
				Description: "Username (if applicable)",
				Computed:    true,
			},
			"group_id": schema.Int64Attribute{
				Description: "Group ID",
				Computed:    true,
			},
			"group_name": schema.StringAttribute{
				Description: "Group name (if applicable)",
				Computed:    true,
			},
			"permission": schema.StringAttribute{
				Description: "Permission type.  See the table referenced in the documentation for an explanation of each permission.",
				Computed:    true,
			},
			"recursive": schema.BoolAttribute{
				Description: "Recursive: does this permission apply to subfolders?",
				Computed:    true,
			},
		},
	}
}

func (r *permissionDataSource) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
	var data permissionDataSourceModel
	diags := req.Config.Get(ctx, &data)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	paramsPermissionList := files_sdk.PermissionListParams{}

	permissionIt, err := r.client.List(paramsPermissionList, files_sdk.WithContext(ctx))
	if err != nil {
		resp.Diagnostics.AddError(
			"Error Reading Files Permission",
			"Could not read permission id "+fmt.Sprint(data.Id.ValueInt64())+": "+err.Error(),
		)
		return
	}

	var permission *files_sdk.Permission
	for permissionIt.Next() {
		entry := permissionIt.Permission()
		if entry.Id == data.Id.ValueInt64() {
			permission = &entry
			break
		}
	}

	if err = permissionIt.Err(); err != nil {
		resp.Diagnostics.AddError(
			"Error Reading Files Permission",
			"Could not read permission id "+fmt.Sprint(data.Id.ValueInt64())+": "+err.Error(),
		)
	}

	if permission == nil {
		resp.Diagnostics.AddError(
			"Error Reading Files Permission",
			"Could not find permission id "+fmt.Sprint(data.Id.ValueInt64())+"",
		)
		return
	}

	diags = r.populateDataSourceModel(ctx, *permission, &data)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	diags = resp.State.Set(ctx, data)
	resp.Diagnostics.Append(diags...)
}

func (r *permissionDataSource) populateDataSourceModel(ctx context.Context, permission files_sdk.Permission, state *permissionDataSourceModel) (diags diag.Diagnostics) {
	state.Id = types.Int64Value(permission.Id)
	state.Path = types.StringValue(permission.Path)
	state.UserId = types.Int64Value(permission.UserId)
	state.Username = types.StringValue(permission.Username)
	state.GroupId = types.Int64Value(permission.GroupId)
	state.GroupName = types.StringValue(permission.GroupName)
	state.Permission = types.StringValue(permission.Permission)
	state.Recursive = types.BoolPointerValue(permission.Recursive)

	return
}
