package provider

import (
	"context"
	"fmt"
	"strconv"
	"strings"

	files_sdk "github.com/Files-com/files-sdk-go/v3"
	permission "github.com/Files-com/files-sdk-go/v3/permission"
	"github.com/hashicorp/terraform-plugin-framework-validators/stringvalidator"
	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/boolplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/int64planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/schema/validator"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

var (
	_ resource.Resource                = &permissionResource{}
	_ resource.ResourceWithConfigure   = &permissionResource{}
	_ resource.ResourceWithImportState = &permissionResource{}
)

func NewPermissionResource() resource.Resource {
	return &permissionResource{}
}

type permissionResource struct {
	client *permission.Client
}

type permissionResourceModel struct {
	Path       types.String `tfsdk:"path"`
	UserId     types.Int64  `tfsdk:"user_id"`
	Username   types.String `tfsdk:"username"`
	GroupId    types.Int64  `tfsdk:"group_id"`
	Permission types.String `tfsdk:"permission"`
	Recursive  types.Bool   `tfsdk:"recursive"`
	Id         types.Int64  `tfsdk:"id"`
	GroupName  types.String `tfsdk:"group_name"`
}

func (r *permissionResource) Configure(_ context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {
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

func (r *permissionResource) Metadata(_ context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_permission"
}

func (r *permissionResource) Schema(_ context.Context, _ resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Description: "A Permission object represents a grant of access permission on a specific Path to a User or Group.\n\n\n\nThey can be optionally recursive or nonrecursive into the subfolders of that path.\n\n\n\nA Permission may be applied to a User *or* a Group, but not both at once.\n\n\n\nThe following table sets forth the available Permission types:\n\n\n\n| Permission | Access Level Granted | Automatically Also Includes/Implies Permissions |\n\n| --- | ----------- | --------------------- |\n\n| `admin` | Able to manage Folder Behaviors, Permissions, and Notifications for the folder. Also grants all other permissions. | `bundle`, `full`, `writeonly`, `readonly`, `list`, `history` |\n\n| `full` | Able to read, write, move, delete, and rename files and folders. Also grants the ability to overwrite files upon upload. | `writeonly`, `readonly`, `list` |\n\n| `readonly` | Able to list, preview, and download files and folders. | `list` |\n\n| `writeonly` | Able to upload files, create folders and list subfolders the user has write permission to. | none |\n\n| `list` | Able to list files and folders, but not download. | none |\n\n| `bundle` | Able to share files and folders via a Bundle (share link). | `readonly`, `list` |\n\n| `history` | Able to view the history of files and folders and to create email notifications for themselves. | `list` |",
		Attributes: map[string]schema.Attribute{
			"path": schema.StringAttribute{
				Description: "Path. This must be slash-delimited, but it must neither start nor end with a slash. Maximum of 5000 characters.",
				Computed:    true,
				Optional:    true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
					stringplanmodifier.RequiresReplace(),
				},
			},
			"user_id": schema.Int64Attribute{
				Description: "User ID",
				Computed:    true,
				Optional:    true,
				PlanModifiers: []planmodifier.Int64{
					int64planmodifier.UseStateForUnknown(),
					int64planmodifier.RequiresReplace(),
				},
			},
			"username": schema.StringAttribute{
				Description: "Username (if applicable)",
				Computed:    true,
				Optional:    true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
					stringplanmodifier.RequiresReplace(),
				},
			},
			"group_id": schema.Int64Attribute{
				Description: "Group ID",
				Computed:    true,
				Optional:    true,
				PlanModifiers: []planmodifier.Int64{
					int64planmodifier.UseStateForUnknown(),
					int64planmodifier.RequiresReplace(),
				},
			},
			"permission": schema.StringAttribute{
				Description: "Permission type.  See the table referenced in the documentation for an explanation of each permission.",
				Computed:    true,
				Optional:    true,
				Validators: []validator.String{
					stringvalidator.OneOf("full", "readonly", "writeonly", "list", "history", "admin", "bundle"),
				},
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
					stringplanmodifier.RequiresReplace(),
				},
			},
			"recursive": schema.BoolAttribute{
				Description: "Recursive: does this permission apply to subfolders?",
				Computed:    true,
				Optional:    true,
				PlanModifiers: []planmodifier.Bool{
					boolplanmodifier.UseStateForUnknown(),
					boolplanmodifier.RequiresReplace(),
				},
			},
			"id": schema.Int64Attribute{
				Description: "Permission ID",
				Computed:    true,
				PlanModifiers: []planmodifier.Int64{
					int64planmodifier.UseStateForUnknown(),
				},
			},
			"group_name": schema.StringAttribute{
				Description: "Group name (if applicable)",
				Computed:    true,
			},
		},
	}
}

func (r *permissionResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	var plan permissionResourceModel
	diags := req.Plan.Get(ctx, &plan)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	paramsPermissionCreate := files_sdk.PermissionCreateParams{}
	paramsPermissionCreate.GroupId = plan.GroupId.ValueInt64()
	paramsPermissionCreate.Path = plan.Path.ValueString()
	paramsPermissionCreate.Permission = plan.Permission.ValueString()
	if !plan.Recursive.IsNull() && !plan.Recursive.IsUnknown() {
		paramsPermissionCreate.Recursive = plan.Recursive.ValueBoolPointer()
	}
	paramsPermissionCreate.UserId = plan.UserId.ValueInt64()
	paramsPermissionCreate.Username = plan.Username.ValueString()

	if resp.Diagnostics.HasError() {
		return
	}

	permission, err := r.client.Create(paramsPermissionCreate, files_sdk.WithContext(ctx))
	if err != nil {
		resp.Diagnostics.AddError(
			"Error Creating Files Permission",
			"Could not create permission, unexpected error: "+err.Error(),
		)
		return
	}

	diags = r.populateResourceModel(ctx, permission, &plan)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	diags = resp.State.Set(ctx, plan)
	resp.Diagnostics.Append(diags...)
}

func (r *permissionResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	var state permissionResourceModel
	diags := req.State.Get(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	paramsPermissionList := files_sdk.PermissionListParams{}

	permissionIt, err := r.client.List(paramsPermissionList, files_sdk.WithContext(ctx))
	if err != nil {
		if files_sdk.IsNotExist(err) {
			resp.State.RemoveResource(ctx)
			return
		}

		resp.Diagnostics.AddError(
			"Error Reading Files Permission",
			"Could not read permission id "+fmt.Sprint(state.Id.ValueInt64())+": "+err.Error(),
		)
		return
	}

	var permission *files_sdk.Permission
	for permissionIt.Next() {
		entry := permissionIt.Permission()
		if entry.Id == state.Id.ValueInt64() {
			permission = &entry
			break
		}
	}

	if err = permissionIt.Err(); err != nil {
		if files_sdk.IsNotExist(err) {
			resp.State.RemoveResource(ctx)
			return
		}

		resp.Diagnostics.AddError(
			"Error Reading Files Permission",
			"Could not read permission id "+fmt.Sprint(state.Id.ValueInt64())+": "+err.Error(),
		)
	}

	if permission == nil {
		resp.State.RemoveResource(ctx)
		return
	}

	diags = r.populateResourceModel(ctx, *permission, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	diags = resp.State.Set(ctx, state)
	resp.Diagnostics.Append(diags...)
}

func (r *permissionResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	resp.Diagnostics.AddError(
		"Resource Update Not Implemented",
		"This resource does not support updates.",
	)
}

func (r *permissionResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	var state permissionResourceModel
	diags := req.State.Get(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	paramsPermissionDelete := files_sdk.PermissionDeleteParams{}
	paramsPermissionDelete.Id = state.Id.ValueInt64()

	err := r.client.Delete(paramsPermissionDelete, files_sdk.WithContext(ctx))
	if err != nil && !files_sdk.IsNotExist(err) {
		resp.Diagnostics.AddError(
			"Error Deleting Files Permission",
			"Could not delete permission id "+fmt.Sprint(state.Id.ValueInt64())+": "+err.Error(),
		)
	}
}

func (r *permissionResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
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

func (r *permissionResource) populateResourceModel(ctx context.Context, permission files_sdk.Permission, state *permissionResourceModel) (diags diag.Diagnostics) {
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
