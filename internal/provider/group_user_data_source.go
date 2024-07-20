package provider

import (
	"context"
	"fmt"

	files_sdk "github.com/Files-com/files-sdk-go/v3"
	group_user "github.com/Files-com/files-sdk-go/v3/groupuser"
	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/datasource/schema"
	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

var (
	_ datasource.DataSource              = &groupUserDataSource{}
	_ datasource.DataSourceWithConfigure = &groupUserDataSource{}
)

func NewGroupUserDataSource() datasource.DataSource {
	return &groupUserDataSource{}
}

type groupUserDataSource struct {
	client *group_user.Client
}

type groupUserDataSourceModel struct {
	GroupName types.String `tfsdk:"group_name"`
	GroupId   types.Int64  `tfsdk:"group_id"`
	UserId    types.Int64  `tfsdk:"user_id"`
	Admin     types.Bool   `tfsdk:"admin"`
	Usernames types.List   `tfsdk:"usernames"`
	Id        types.Int64  `tfsdk:"id"`
}

func (r *groupUserDataSource) Configure(_ context.Context, req datasource.ConfigureRequest, resp *datasource.ConfigureResponse) {
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

func (r *groupUserDataSource) Metadata(_ context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_group_user"
}

func (r *groupUserDataSource) Schema(_ context.Context, _ datasource.SchemaRequest, resp *datasource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Description: "A GroupUser describes the membership of a User within a Group.\n\n\n\n## Creating GroupUsers\n\nGroupUsers can be created via the normal `create` action. When using the `update` action, if the\n\nGroupUser record does not exist for the given user/group IDs it will be created.",
		Attributes: map[string]schema.Attribute{
			"group_name": schema.StringAttribute{
				Description: "Group name",
				Computed:    true,
			},
			"group_id": schema.Int64Attribute{
				Description: "Group ID",
				Computed:    true,
			},
			"user_id": schema.Int64Attribute{
				Description: "User ID",
				Computed:    true,
			},
			"admin": schema.BoolAttribute{
				Description: "Is this user an administrator of this group?",
				Computed:    true,
			},
			"usernames": schema.ListAttribute{
				Description: "A list of usernames for users in this group",
				Computed:    true,
				ElementType: types.StringType,
			},
			"id": schema.Int64Attribute{
				Description: "Group User ID.",
				Required:    true,
			},
		},
	}
}

func (r *groupUserDataSource) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
	var data groupUserDataSourceModel
	diags := req.Config.Get(ctx, &data)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	paramsGroupUserList := files_sdk.GroupUserListParams{}

	groupUserIt, err := r.client.List(paramsGroupUserList, files_sdk.WithContext(ctx))
	if err != nil {
		resp.Diagnostics.AddError(
			"Error Reading Files GroupUser",
			"Could not read group_user id "+fmt.Sprint(data.Id.ValueInt64())+": "+err.Error(),
		)
		return
	}

	var groupUser *files_sdk.GroupUser
	for groupUserIt.Next() {
		entry := groupUserIt.GroupUser()
		if entry.Id == data.Id.ValueInt64() {
			groupUser = &entry
			break
		}
	}

	if groupUser == nil {
		resp.Diagnostics.AddError(
			"Error Reading Files GroupUser",
			"Could not find group_user id "+fmt.Sprint(data.Id.ValueInt64()),
		)
		return
	}

	diags = r.populateDataSourceModel(ctx, *groupUser, &data)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	diags = resp.State.Set(ctx, data)
	resp.Diagnostics.Append(diags...)
}

func (r *groupUserDataSource) populateDataSourceModel(ctx context.Context, groupUser files_sdk.GroupUser, state *groupUserDataSourceModel) (diags diag.Diagnostics) {
	var propDiags diag.Diagnostics

	state.GroupName = types.StringValue(groupUser.GroupName)
	state.GroupId = types.Int64Value(groupUser.GroupId)
	state.UserId = types.Int64Value(groupUser.UserId)
	state.Admin = types.BoolPointerValue(groupUser.Admin)
	state.Usernames, propDiags = types.ListValueFrom(ctx, types.StringType, groupUser.Usernames)
	diags.Append(propDiags...)

	return
}
