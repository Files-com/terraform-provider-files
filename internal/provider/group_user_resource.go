package provider

import (
	"context"
	"fmt"
	"strconv"
	"strings"

	files_sdk "github.com/Files-com/files-sdk-go/v3"
	group_user "github.com/Files-com/files-sdk-go/v3/groupuser"
	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/boolplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/int64planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

var (
	_ resource.Resource                = &groupUserResource{}
	_ resource.ResourceWithConfigure   = &groupUserResource{}
	_ resource.ResourceWithImportState = &groupUserResource{}
)

func NewGroupUserResource() resource.Resource {
	return &groupUserResource{}
}

type groupUserResource struct {
	client *group_user.Client
}

type groupUserResourceModel struct {
	GroupId   types.Int64  `tfsdk:"group_id"`
	UserId    types.Int64  `tfsdk:"user_id"`
	Admin     types.Bool   `tfsdk:"admin"`
	GroupName types.String `tfsdk:"group_name"`
	Usernames types.List   `tfsdk:"usernames"`
	Id        types.Int64  `tfsdk:"id"`
}

func (r *groupUserResource) Configure(_ context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {
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

	r.client = &group_user.Client{Config: sdk_config}
}

func (r *groupUserResource) Metadata(_ context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_group_user"
}

func (r *groupUserResource) Schema(_ context.Context, _ resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Description: "A GroupUser describes the membership of a User within a Group.\n\n\n\n## Creating GroupUsers\n\nGroupUsers can be created via the normal `create` action. When using the `update` action, if the\n\nGroupUser record does not exist for the given user/group IDs it will be created.",
		Attributes: map[string]schema.Attribute{
			"group_id": schema.Int64Attribute{
				Description: "Group ID",
				Required:    true,
			},
			"user_id": schema.Int64Attribute{
				Description: "User ID",
				Required:    true,
			},
			"admin": schema.BoolAttribute{
				Description: "Is this user an administrator of this group?",
				Computed:    true,
				Optional:    true,
				PlanModifiers: []planmodifier.Bool{
					boolplanmodifier.UseStateForUnknown(),
				},
			},
			"group_name": schema.StringAttribute{
				Description: "Group name",
				Computed:    true,
			},
			"usernames": schema.ListAttribute{
				Description: "A list of usernames for users in this group",
				Computed:    true,
				ElementType: types.StringType,
			},
			"id": schema.Int64Attribute{
				Description: "Group User ID.",
				Computed:    true,
				PlanModifiers: []planmodifier.Int64{
					int64planmodifier.UseStateForUnknown(),
				},
			},
		},
	}
}

func (r *groupUserResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	var plan groupUserResourceModel
	diags := req.Plan.Get(ctx, &plan)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	paramsGroupUserCreate := files_sdk.GroupUserCreateParams{}
	paramsGroupUserCreate.GroupId = plan.GroupId.ValueInt64()
	paramsGroupUserCreate.UserId = plan.UserId.ValueInt64()
	if !plan.Admin.IsNull() && !plan.Admin.IsUnknown() {
		paramsGroupUserCreate.Admin = plan.Admin.ValueBoolPointer()
	}

	if resp.Diagnostics.HasError() {
		return
	}

	groupUser, err := r.client.Create(paramsGroupUserCreate, files_sdk.WithContext(ctx))
	if err != nil {
		resp.Diagnostics.AddError(
			"Error Creating Files GroupUser",
			"Could not create group_user, unexpected error: "+err.Error(),
		)
		return
	}

	diags = r.populateResourceModel(ctx, groupUser, &plan)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	diags = resp.State.Set(ctx, plan)
	resp.Diagnostics.Append(diags...)
}

func (r *groupUserResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	var state groupUserResourceModel
	diags := req.State.Get(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	paramsGroupUserList := files_sdk.GroupUserListParams{}

	groupUserIt, err := r.client.List(paramsGroupUserList, files_sdk.WithContext(ctx))
	if err != nil {
		if files_sdk.IsNotExist(err) {
			resp.State.RemoveResource(ctx)
			return
		}

		resp.Diagnostics.AddError(
			"Error Reading Files GroupUser",
			"Could not read group_user id "+fmt.Sprint(state.Id.ValueInt64())+": "+err.Error(),
		)
		return
	}

	var groupUser *files_sdk.GroupUser
	for groupUserIt.Next() {
		entry := groupUserIt.GroupUser()
		if entry.Id == state.Id.ValueInt64() {
			groupUser = &entry
			break
		}
	}

	if groupUser == nil {
		resp.State.RemoveResource(ctx)
		return
	}

	diags = r.populateResourceModel(ctx, *groupUser, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	diags = resp.State.Set(ctx, state)
	resp.Diagnostics.Append(diags...)
}

func (r *groupUserResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	var plan groupUserResourceModel
	diags := req.Plan.Get(ctx, &plan)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	paramsGroupUserUpdate := files_sdk.GroupUserUpdateParams{}
	paramsGroupUserUpdate.Id = plan.Id.ValueInt64()
	paramsGroupUserUpdate.GroupId = plan.GroupId.ValueInt64()
	paramsGroupUserUpdate.UserId = plan.UserId.ValueInt64()
	if !plan.Admin.IsNull() && !plan.Admin.IsUnknown() {
		paramsGroupUserUpdate.Admin = plan.Admin.ValueBoolPointer()
	}

	if resp.Diagnostics.HasError() {
		return
	}

	groupUser, err := r.client.Update(paramsGroupUserUpdate, files_sdk.WithContext(ctx))
	if err != nil {
		resp.Diagnostics.AddError(
			"Error Updating Files GroupUser",
			"Could not update group_user, unexpected error: "+err.Error(),
		)
		return
	}

	diags = r.populateResourceModel(ctx, groupUser, &plan)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	diags = resp.State.Set(ctx, plan)
	resp.Diagnostics.Append(diags...)
}

func (r *groupUserResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	var state groupUserResourceModel
	diags := req.State.Get(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	paramsGroupUserDelete := files_sdk.GroupUserDeleteParams{}
	paramsGroupUserDelete.Id = state.Id.ValueInt64()
	paramsGroupUserDelete.GroupId = state.GroupId.ValueInt64()
	paramsGroupUserDelete.UserId = state.UserId.ValueInt64()

	err := r.client.Delete(paramsGroupUserDelete, files_sdk.WithContext(ctx))
	if err != nil && !files_sdk.IsNotExist(err) {
		resp.Diagnostics.AddError(
			"Error Deleting Files GroupUser",
			"Could not delete group_user id "+fmt.Sprint(state.Id.ValueInt64())+": "+err.Error(),
		)
	}
}

func (r *groupUserResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
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

func (r *groupUserResource) populateResourceModel(ctx context.Context, groupUser files_sdk.GroupUser, state *groupUserResourceModel) (diags diag.Diagnostics) {
	var propDiags diag.Diagnostics

	state.GroupName = types.StringValue(groupUser.GroupName)
	state.GroupId = types.Int64Value(groupUser.GroupId)
	state.UserId = types.Int64Value(groupUser.UserId)
	state.Admin = types.BoolPointerValue(groupUser.Admin)
	state.Usernames, propDiags = types.ListValueFrom(ctx, types.StringType, groupUser.Usernames)
	diags.Append(propDiags...)

	return
}
