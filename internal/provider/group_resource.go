package provider

import (
	"context"
	"fmt"
	"strconv"
	"strings"

	files_sdk "github.com/Files-com/files-sdk-go/v3"
	group "github.com/Files-com/files-sdk-go/v3/group"
	"github.com/Files-com/terraform-provider-files/lib"
	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/boolplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/int64planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

var (
	_ resource.Resource                = &groupResource{}
	_ resource.ResourceWithConfigure   = &groupResource{}
	_ resource.ResourceWithImportState = &groupResource{}
)

func NewGroupResource() resource.Resource {
	return &groupResource{}
}

type groupResource struct {
	client *group.Client
}

type groupResourceModel struct {
	Name              types.String            `tfsdk:"name"`
	AllowedIps        types.String            `tfsdk:"allowed_ips"`
	AdminIds          lib.SortedElementString `tfsdk:"admin_ids"`
	Notes             types.String            `tfsdk:"notes"`
	UserIds           lib.SortedElementString `tfsdk:"user_ids"`
	FtpPermission     types.Bool              `tfsdk:"ftp_permission"`
	SftpPermission    types.Bool              `tfsdk:"sftp_permission"`
	DavPermission     types.Bool              `tfsdk:"dav_permission"`
	RestapiPermission types.Bool              `tfsdk:"restapi_permission"`
	WorkspaceId       types.Int64             `tfsdk:"workspace_id"`
	Id                types.Int64             `tfsdk:"id"`
	Usernames         types.String            `tfsdk:"usernames"`
	SiteId            types.Int64             `tfsdk:"site_id"`
}

func (r *groupResource) Configure(_ context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {
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

func (r *groupResource) Metadata(_ context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_group"
}

func (r *groupResource) Schema(_ context.Context, _ resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Description: "A Group is a powerful tool for permissions and user management on Files.com. Users can belong to multiple groups.\n\n\n\nAll permissions can be managed via Groups, and Groups can also be synced to your identity platform via LDAP or SCIM.\n\n\n\nFiles.com's Group Admin feature allows you to define Group Admins, who then have access to add and remove users within their groups.",
		Attributes: map[string]schema.Attribute{
			"name": schema.StringAttribute{
				Description: "Group name",
				Required:    true,
			},
			"allowed_ips": schema.StringAttribute{
				Description: "A list of allowed IPs if applicable.  Newline delimited",
				Computed:    true,
				Optional:    true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
			},
			"admin_ids": schema.StringAttribute{
				Description: "Comma-delimited list of user IDs who are group administrators (separated by commas)",
				Computed:    true,
				Optional:    true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
				CustomType: lib.SortedElementStringType{},
			},
			"notes": schema.StringAttribute{
				Description: "Notes about this group",
				Computed:    true,
				Optional:    true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
			},
			"user_ids": schema.StringAttribute{
				Description: "Comma-delimited list of user IDs who belong to this group (separated by commas)",
				Computed:    true,
				Optional:    true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
				CustomType: lib.SortedElementStringType{},
			},
			"ftp_permission": schema.BoolAttribute{
				Description: "If true, users in this group can use FTP to login.  This will override a false value of `ftp_permission` on the user level.",
				Computed:    true,
				Optional:    true,
				PlanModifiers: []planmodifier.Bool{
					boolplanmodifier.UseStateForUnknown(),
				},
			},
			"sftp_permission": schema.BoolAttribute{
				Description: "If true, users in this group can use SFTP to login.  This will override a false value of `sftp_permission` on the user level.",
				Computed:    true,
				Optional:    true,
				PlanModifiers: []planmodifier.Bool{
					boolplanmodifier.UseStateForUnknown(),
				},
			},
			"dav_permission": schema.BoolAttribute{
				Description: "If true, users in this group can use WebDAV to login.  This will override a false value of `dav_permission` on the user level.",
				Computed:    true,
				Optional:    true,
				PlanModifiers: []planmodifier.Bool{
					boolplanmodifier.UseStateForUnknown(),
				},
			},
			"restapi_permission": schema.BoolAttribute{
				Description: "If true, users in this group can use the REST API to login.  This will override a false value of `restapi_permission` on the user level.",
				Computed:    true,
				Optional:    true,
				PlanModifiers: []planmodifier.Bool{
					boolplanmodifier.UseStateForUnknown(),
				},
			},
			"workspace_id": schema.Int64Attribute{
				Description: "Workspace ID",
				Computed:    true,
				Optional:    true,
				PlanModifiers: []planmodifier.Int64{
					int64planmodifier.UseStateForUnknown(),
				},
			},
			"id": schema.Int64Attribute{
				Description: "Group ID",
				Computed:    true,
				PlanModifiers: []planmodifier.Int64{
					int64planmodifier.UseStateForUnknown(),
				},
			},
			"usernames": schema.StringAttribute{
				Description: "Comma-delimited list of usernames who belong to this group (separated by commas)",
				Computed:    true,
			},
			"site_id": schema.Int64Attribute{
				Description: "Site ID",
				Computed:    true,
			},
		},
	}
}

func (r *groupResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	var plan groupResourceModel
	diags := req.Plan.Get(ctx, &plan)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}
	var config groupResourceModel
	diags = req.Config.Get(ctx, &config)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	paramsGroupCreate := files_sdk.GroupCreateParams{}
	paramsGroupCreate.Notes = plan.Notes.ValueString()
	paramsGroupCreate.UserIds = plan.UserIds.ValueString()
	paramsGroupCreate.AdminIds = plan.AdminIds.ValueString()
	paramsGroupCreate.WorkspaceId = plan.WorkspaceId.ValueInt64()
	if !plan.FtpPermission.IsNull() && !plan.FtpPermission.IsUnknown() {
		paramsGroupCreate.FtpPermission = plan.FtpPermission.ValueBoolPointer()
	}
	if !plan.SftpPermission.IsNull() && !plan.SftpPermission.IsUnknown() {
		paramsGroupCreate.SftpPermission = plan.SftpPermission.ValueBoolPointer()
	}
	if !plan.DavPermission.IsNull() && !plan.DavPermission.IsUnknown() {
		paramsGroupCreate.DavPermission = plan.DavPermission.ValueBoolPointer()
	}
	if !plan.RestapiPermission.IsNull() && !plan.RestapiPermission.IsUnknown() {
		paramsGroupCreate.RestapiPermission = plan.RestapiPermission.ValueBoolPointer()
	}
	paramsGroupCreate.AllowedIps = plan.AllowedIps.ValueString()
	paramsGroupCreate.Name = plan.Name.ValueString()

	if resp.Diagnostics.HasError() {
		return
	}

	group, err := r.client.Create(paramsGroupCreate, files_sdk.WithContext(ctx))
	if err != nil {
		resp.Diagnostics.AddError(
			"Error Creating Files Group",
			"Could not create group, unexpected error: "+err.Error(),
		)
		return
	}

	diags = r.populateResourceModel(ctx, group, &plan)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	diags = resp.State.Set(ctx, plan)
	resp.Diagnostics.Append(diags...)
}

func (r *groupResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	var state groupResourceModel
	diags := req.State.Get(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	paramsGroupFind := files_sdk.GroupFindParams{}
	paramsGroupFind.Id = state.Id.ValueInt64()

	group, err := r.client.Find(paramsGroupFind, files_sdk.WithContext(ctx))
	if err != nil {
		if files_sdk.IsNotExist(err) {
			resp.State.RemoveResource(ctx)
			return
		}

		resp.Diagnostics.AddError(
			"Error Reading Files Group",
			"Could not read group id "+fmt.Sprint(state.Id.ValueInt64())+": "+err.Error(),
		)
		return
	}

	diags = r.populateResourceModel(ctx, group, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	diags = resp.State.Set(ctx, state)
	resp.Diagnostics.Append(diags...)
}

func (r *groupResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	var plan groupResourceModel
	diags := req.Plan.Get(ctx, &plan)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}
	var config groupResourceModel
	diags = req.Config.Get(ctx, &config)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	paramsGroupUpdate := map[string]interface{}{}
	if !plan.Id.IsNull() && !plan.Id.IsUnknown() {
		paramsGroupUpdate["id"] = plan.Id.ValueInt64()
	}
	if !config.Notes.IsNull() && !config.Notes.IsUnknown() {
		paramsGroupUpdate["notes"] = config.Notes.ValueString()
	}
	if !config.UserIds.IsNull() && !config.UserIds.IsUnknown() {
		paramsGroupUpdate["user_ids"] = config.UserIds.ValueString()
	}
	if !config.AdminIds.IsNull() && !config.AdminIds.IsUnknown() {
		paramsGroupUpdate["admin_ids"] = config.AdminIds.ValueString()
	}
	if !config.WorkspaceId.IsNull() && !config.WorkspaceId.IsUnknown() {
		paramsGroupUpdate["workspace_id"] = config.WorkspaceId.ValueInt64()
	}
	if !config.FtpPermission.IsNull() && !config.FtpPermission.IsUnknown() {
		paramsGroupUpdate["ftp_permission"] = config.FtpPermission.ValueBool()
	}
	if !config.SftpPermission.IsNull() && !config.SftpPermission.IsUnknown() {
		paramsGroupUpdate["sftp_permission"] = config.SftpPermission.ValueBool()
	}
	if !config.DavPermission.IsNull() && !config.DavPermission.IsUnknown() {
		paramsGroupUpdate["dav_permission"] = config.DavPermission.ValueBool()
	}
	if !config.RestapiPermission.IsNull() && !config.RestapiPermission.IsUnknown() {
		paramsGroupUpdate["restapi_permission"] = config.RestapiPermission.ValueBool()
	}
	if !config.AllowedIps.IsNull() && !config.AllowedIps.IsUnknown() {
		paramsGroupUpdate["allowed_ips"] = config.AllowedIps.ValueString()
	}
	if !config.Name.IsNull() && !config.Name.IsUnknown() {
		paramsGroupUpdate["name"] = config.Name.ValueString()
	}

	if resp.Diagnostics.HasError() {
		return
	}

	group, err := r.client.UpdateWithMap(paramsGroupUpdate, files_sdk.WithContext(ctx))
	if err != nil {
		resp.Diagnostics.AddError(
			"Error Updating Files Group",
			"Could not update group, unexpected error: "+err.Error(),
		)
		return
	}

	diags = r.populateResourceModel(ctx, group, &plan)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	diags = resp.State.Set(ctx, plan)
	resp.Diagnostics.Append(diags...)
}

func (r *groupResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	var state groupResourceModel
	diags := req.State.Get(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	paramsGroupDelete := files_sdk.GroupDeleteParams{}
	paramsGroupDelete.Id = state.Id.ValueInt64()

	err := r.client.Delete(paramsGroupDelete, files_sdk.WithContext(ctx))
	if err != nil && !files_sdk.IsNotExist(err) {
		resp.Diagnostics.AddError(
			"Error Deleting Files Group",
			"Could not delete group id "+fmt.Sprint(state.Id.ValueInt64())+": "+err.Error(),
		)
	}
}

func (r *groupResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
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

func (r *groupResource) populateResourceModel(ctx context.Context, group files_sdk.Group, state *groupResourceModel) (diags diag.Diagnostics) {
	state.Id = types.Int64Value(group.Id)
	state.Name = types.StringValue(group.Name)
	state.AllowedIps = types.StringValue(group.AllowedIps)
	state.AdminIds = lib.SortedElementStringValue(group.AdminIds)
	state.Notes = types.StringValue(group.Notes)
	state.UserIds = lib.SortedElementStringValue(group.UserIds)
	state.Usernames = types.StringValue(group.Usernames)
	state.FtpPermission = types.BoolPointerValue(group.FtpPermission)
	state.SftpPermission = types.BoolPointerValue(group.SftpPermission)
	state.DavPermission = types.BoolPointerValue(group.DavPermission)
	state.RestapiPermission = types.BoolPointerValue(group.RestapiPermission)
	state.SiteId = types.Int64Value(group.SiteId)
	state.WorkspaceId = types.Int64Value(group.WorkspaceId)

	return
}
