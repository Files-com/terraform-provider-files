package provider

import (
	"context"
	"fmt"

	files_sdk "github.com/Files-com/files-sdk-go/v3"
	group "github.com/Files-com/files-sdk-go/v3/group"
	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/datasource/schema"
	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

var (
	_ datasource.DataSource              = &groupDataSource{}
	_ datasource.DataSourceWithConfigure = &groupDataSource{}
)

func NewGroupDataSource() datasource.DataSource {
	return &groupDataSource{}
}

type groupDataSource struct {
	client *group.Client
}

type groupDataSourceModel struct {
	Id                types.Int64  `tfsdk:"id"`
	Name              types.String `tfsdk:"name"`
	AllowedIps        types.String `tfsdk:"allowed_ips"`
	AdminIds          types.String `tfsdk:"admin_ids"`
	Notes             types.String `tfsdk:"notes"`
	UserIds           types.String `tfsdk:"user_ids"`
	Usernames         types.String `tfsdk:"usernames"`
	FtpPermission     types.Bool   `tfsdk:"ftp_permission"`
	SftpPermission    types.Bool   `tfsdk:"sftp_permission"`
	DavPermission     types.Bool   `tfsdk:"dav_permission"`
	RestapiPermission types.Bool   `tfsdk:"restapi_permission"`
}

func (r *groupDataSource) Configure(_ context.Context, req datasource.ConfigureRequest, resp *datasource.ConfigureResponse) {
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

	r.client = &group.Client{Config: sdk_config}
}

func (r *groupDataSource) Metadata(_ context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_group"
}

func (r *groupDataSource) Schema(_ context.Context, _ datasource.SchemaRequest, resp *datasource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Description: "Groups are a powerful tool for permissions and user management on Files.com. Users can belong to multiple groups.\n\n\n\nAll permissions can be managed via Groups, and Groups can also be synced to your identity platform via LDAP or SCIM.\n\n\n\nFiles.com's Group Admin feature allows you to define Group Admins, who then have access to add and remove users within their groups.",
		Attributes: map[string]schema.Attribute{
			"id": schema.Int64Attribute{
				Description: "Group ID",
				Required:    true,
			},
			"name": schema.StringAttribute{
				Description: "Group name",
				Computed:    true,
			},
			"allowed_ips": schema.StringAttribute{
				Description: "A list of allowed IPs if applicable.  Newline delimited",
				Computed:    true,
			},
			"admin_ids": schema.StringAttribute{
				Description: "Comma-delimited list of user IDs who are group administrators (separated by commas)",
				Computed:    true,
			},
			"notes": schema.StringAttribute{
				Description: "Notes about this group",
				Computed:    true,
			},
			"user_ids": schema.StringAttribute{
				Description: "Comma-delimited list of user IDs who belong to this group (separated by commas)",
				Computed:    true,
			},
			"usernames": schema.StringAttribute{
				Description: "Comma-delimited list of usernames who belong to this group (separated by commas)",
				Computed:    true,
			},
			"ftp_permission": schema.BoolAttribute{
				Description: "If true, users in this group can use FTP to login.  This will override a false value of `ftp_permission` on the user level.",
				Computed:    true,
			},
			"sftp_permission": schema.BoolAttribute{
				Description: "If true, users in this group can use SFTP to login.  This will override a false value of `sftp_permission` on the user level.",
				Computed:    true,
			},
			"dav_permission": schema.BoolAttribute{
				Description: "If true, users in this group can use WebDAV to login.  This will override a false value of `dav_permission` on the user level.",
				Computed:    true,
			},
			"restapi_permission": schema.BoolAttribute{
				Description: "If true, users in this group can use the REST API to login.  This will override a false value of `restapi_permission` on the user level.",
				Computed:    true,
			},
		},
	}
}

func (r *groupDataSource) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
	var data groupDataSourceModel
	diags := req.Config.Get(ctx, &data)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	paramsGroupFind := files_sdk.GroupFindParams{}
	paramsGroupFind.Id = data.Id.ValueInt64()

	group, err := r.client.Find(paramsGroupFind, files_sdk.WithContext(ctx))
	if err != nil {
		resp.Diagnostics.AddError(
			"Error Reading Files Group",
			"Could not read group id "+fmt.Sprint(data.Id.ValueInt64())+": "+err.Error(),
		)
		return
	}

	diags = r.populateDataSourceModel(ctx, group, &data)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	diags = resp.State.Set(ctx, data)
	resp.Diagnostics.Append(diags...)
}

func (r *groupDataSource) populateDataSourceModel(ctx context.Context, group files_sdk.Group, state *groupDataSourceModel) (diags diag.Diagnostics) {
	state.Id = types.Int64Value(group.Id)
	state.Name = types.StringValue(group.Name)
	state.AllowedIps = types.StringValue(group.AllowedIps)
	state.AdminIds = types.StringValue(group.AdminIds)
	state.Notes = types.StringValue(group.Notes)
	state.UserIds = types.StringValue(group.UserIds)
	state.Usernames = types.StringValue(group.Usernames)
	state.FtpPermission = types.BoolPointerValue(group.FtpPermission)
	state.SftpPermission = types.BoolPointerValue(group.SftpPermission)
	state.DavPermission = types.BoolPointerValue(group.DavPermission)
	state.RestapiPermission = types.BoolPointerValue(group.RestapiPermission)

	return
}
