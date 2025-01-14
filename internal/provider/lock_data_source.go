package provider

import (
	"context"
	"fmt"

	files_sdk "github.com/Files-com/files-sdk-go/v3"
	lock "github.com/Files-com/files-sdk-go/v3/lock"
	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/datasource/schema"
	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

var (
	_ datasource.DataSource              = &lockDataSource{}
	_ datasource.DataSourceWithConfigure = &lockDataSource{}
)

func NewLockDataSource() datasource.DataSource {
	return &lockDataSource{}
}

type lockDataSource struct {
	client *lock.Client
}

type lockDataSourceModel struct {
	Path                 types.String `tfsdk:"path"`
	Timeout              types.Int64  `tfsdk:"timeout"`
	Depth                types.String `tfsdk:"depth"`
	Recursive            types.Bool   `tfsdk:"recursive"`
	Owner                types.String `tfsdk:"owner"`
	Scope                types.String `tfsdk:"scope"`
	Exclusive            types.Bool   `tfsdk:"exclusive"`
	Token                types.String `tfsdk:"token"`
	Type                 types.String `tfsdk:"type"`
	AllowAccessByAnyUser types.Bool   `tfsdk:"allow_access_by_any_user"`
	UserId               types.Int64  `tfsdk:"user_id"`
	Username             types.String `tfsdk:"username"`
}

func (r *lockDataSource) Configure(_ context.Context, req datasource.ConfigureRequest, resp *datasource.ConfigureResponse) {
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

	r.client = &lock.Client{Config: sdk_config}
}

func (r *lockDataSource) Metadata(_ context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_lock"
}

func (r *lockDataSource) Schema(_ context.Context, _ datasource.SchemaRequest, resp *datasource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Description: "A Lock can be used by your custom-developed applications to implement file locking and concurrency features. These locks are advisory, meaning that while a lock can be created, it does not prevent other API requests from being processed concurrently. You are responsible for checking locks prior to accessing a file.\n\n\n\nThe lock feature is designed to emulate the locking functionality provided by WebDAV. For a deeper understanding of how the lock mechanism works, refer to the WebDAV specification, which outlines how these endpoints function.\n\n\n\nFiles.com's WebDAV offering and desktop app leverage this locking API to manage concurrent file operations, ensuring consistency when multiple users or systems interact with the same files. It is not used within the Files.com web interface.",
		Attributes: map[string]schema.Attribute{
			"path": schema.StringAttribute{
				Description: "Path. This must be slash-delimited, but it must neither start nor end with a slash. Maximum of 5000 characters.",
				Required:    true,
			},
			"timeout": schema.Int64Attribute{
				Description: "Lock timeout in seconds",
				Computed:    true,
			},
			"depth": schema.StringAttribute{
				Computed: true,
			},
			"recursive": schema.BoolAttribute{
				Description: "Does lock apply to subfolders?",
				Computed:    true,
			},
			"owner": schema.StringAttribute{
				Description: "Owner of the lock.  This can be any arbitrary string.",
				Computed:    true,
			},
			"scope": schema.StringAttribute{
				Computed: true,
			},
			"exclusive": schema.BoolAttribute{
				Description: "Is lock exclusive?",
				Computed:    true,
			},
			"token": schema.StringAttribute{
				Description: "Lock token.  Use to release lock.",
				Computed:    true,
			},
			"type": schema.StringAttribute{
				Computed: true,
			},
			"allow_access_by_any_user": schema.BoolAttribute{
				Description: "Can lock be modified by users other than its creator?",
				Computed:    true,
			},
			"user_id": schema.Int64Attribute{
				Description: "Lock creator user ID",
				Computed:    true,
			},
			"username": schema.StringAttribute{
				Description: "Lock creator username",
				Computed:    true,
			},
		},
	}
}

func (r *lockDataSource) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
	var data lockDataSourceModel
	diags := req.Config.Get(ctx, &data)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	paramsLockListFor := files_sdk.LockListForParams{}
	paramsLockListFor.Path = data.Path.ValueString()

	lockIt, err := r.client.ListFor(paramsLockListFor, files_sdk.WithContext(ctx))
	if err != nil {
		resp.Diagnostics.AddError(
			"Error Reading Files Lock",
			"Could not read lock path "+fmt.Sprint(data.Path.ValueString())+": "+err.Error(),
		)
		return
	}

	var lock *files_sdk.Lock
	for lockIt.Next() {
		entry := lockIt.Lock()
		if entry.Path == data.Path.ValueString() {
			lock = &entry
			break
		}
	}

	if err = lockIt.Err(); err != nil {
		resp.Diagnostics.AddError(
			"Error Reading Files Lock",
			"Could not read lock path "+fmt.Sprint(data.Path.ValueString())+": "+err.Error(),
		)
	}

	if lock == nil {
		resp.Diagnostics.AddError(
			"Error Reading Files Lock",
			"Could not find lock path "+fmt.Sprint(data.Path.ValueString())+"",
		)
		return
	}

	diags = r.populateDataSourceModel(ctx, *lock, &data)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	diags = resp.State.Set(ctx, data)
	resp.Diagnostics.Append(diags...)
}

func (r *lockDataSource) populateDataSourceModel(ctx context.Context, lock files_sdk.Lock, state *lockDataSourceModel) (diags diag.Diagnostics) {
	state.Path = types.StringValue(lock.Path)
	state.Timeout = types.Int64Value(lock.Timeout)
	state.Depth = types.StringValue(lock.Depth)
	state.Recursive = types.BoolPointerValue(lock.Recursive)
	state.Owner = types.StringValue(lock.Owner)
	state.Scope = types.StringValue(lock.Scope)
	state.Exclusive = types.BoolPointerValue(lock.Exclusive)
	state.Token = types.StringValue(lock.Token)
	state.Type = types.StringValue(lock.Type)
	state.AllowAccessByAnyUser = types.BoolPointerValue(lock.AllowAccessByAnyUser)
	state.UserId = types.Int64Value(lock.UserId)
	state.Username = types.StringValue(lock.Username)

	return
}
