package provider

import (
	"context"
	"fmt"
	"strings"

	files_sdk "github.com/Files-com/files-sdk-go/v3"
	lock "github.com/Files-com/files-sdk-go/v3/lock"

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
	_ resource.Resource                = &lockResource{}
	_ resource.ResourceWithConfigure   = &lockResource{}
	_ resource.ResourceWithImportState = &lockResource{}
)

func NewLockResource() resource.Resource {
	return &lockResource{}
}

type lockResource struct {
	client *lock.Client
}

type lockResourceModel struct {
	Path                 types.String `tfsdk:"path"`
	Timeout              types.Int64  `tfsdk:"timeout"`
	Recursive            types.Bool   `tfsdk:"recursive"`
	Owner                types.String `tfsdk:"owner"`
	Exclusive            types.Bool   `tfsdk:"exclusive"`
	Token                types.String `tfsdk:"token"`
	AllowAccessByAnyUser types.Bool   `tfsdk:"allow_access_by_any_user"`
	ExpectedToken        types.String `tfsdk:"expected_token"`
	Depth                types.String `tfsdk:"depth"`
	Scope                types.String `tfsdk:"scope"`
	Type                 types.String `tfsdk:"type"`
	UserId               types.Int64  `tfsdk:"user_id"`
	Username             types.String `tfsdk:"username"`
}

func (r *lockResource) Configure(_ context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {
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

func (r *lockResource) Metadata(_ context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_lock"
}

func (r *lockResource) Schema(_ context.Context, _ resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Description: "A Lock can be used by your custom-developed applications to implement file locking and concurrency features. These locks are advisory, meaning that while a lock can be created, it does not prevent other API requests from being processed concurrently.  You are responsible for checking locks prior to accessing a file.\n\nThe lock feature is designed to emulate the locking functionality provided by WebDAV. For a deeper understanding of how the lock mechanism works, refer to the WebDAV specification, which outlines how these endpoints function.\n\nFiles.com's WebDAV offering and desktop app leverage this locking API to manage concurrent file operations, ensuring consistency when multiple users or systems interact with the same files.  It is not used within the Files.com web interface.\n\nThe optional owner parameter is a descriptive label, not the lock creator or a grant of permission. It can be set when creating a lock; refreshing a lock or replacing its token preserves it.\n\nTo refresh only an existing lock or replace its token, send expected_token, token, and timeout to the create endpoint. Set token to expected_token to refresh, or to a different value to replace. The expected token must identify an existing, unexpired lock on that path, and the caller must have permission to modify it. The token check and update happen together; invalid replacement values leave the stored lock unchanged.\n\nA missing, expired, or mismatched expected token returns processing-failure/resource-locked with data.lock_token containing an active token on that path, or an empty string when none exists. Omitting expected_token retains the existing acquire-or-refresh behavior. Shared locks retain their existing semantics.",
		Attributes: map[string]schema.Attribute{
			"path": schema.StringAttribute{
				Description: "Path. This must be slash-delimited, but it must neither start nor end with a slash. Maximum of 5000 characters.",
				Required:    true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
					stringplanmodifier.RequiresReplace(),
				},
			},
			"timeout": schema.Int64Attribute{
				Description: "Lock timeout in seconds",
				Computed:    true,
				Optional:    true,
				PlanModifiers: []planmodifier.Int64{
					int64planmodifier.UseStateForUnknown(),
					int64planmodifier.RequiresReplace(),
				},
			},
			"recursive": schema.BoolAttribute{
				Description: "Does lock apply to subfolders?",
				Computed:    true,
				Optional:    true,
				PlanModifiers: []planmodifier.Bool{
					boolplanmodifier.UseStateForUnknown(),
					boolplanmodifier.RequiresReplace(),
				},
			},
			"owner": schema.StringAttribute{
				Description: "Arbitrary descriptive label for the lock. Does not change the lock creator or permissions.",
				Computed:    true,
				Optional:    true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
					stringplanmodifier.RequiresReplace(),
				},
			},
			"exclusive": schema.BoolAttribute{
				Description: "Is lock exclusive?",
				Computed:    true,
				Optional:    true,
				PlanModifiers: []planmodifier.Bool{
					boolplanmodifier.UseStateForUnknown(),
					boolplanmodifier.RequiresReplace(),
				},
			},
			"token": schema.StringAttribute{
				Description: "Lock token.  Use to release lock.",
				Computed:    true,
				Optional:    true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
					stringplanmodifier.RequiresReplace(),
				},
			},
			"allow_access_by_any_user": schema.BoolAttribute{
				Description: "Can lock be modified by users other than its creator?",
				Computed:    true,
				Optional:    true,
				PlanModifiers: []planmodifier.Bool{
					boolplanmodifier.UseStateForUnknown(),
					boolplanmodifier.RequiresReplace(),
				},
			},
			"expected_token": schema.StringAttribute{
				Description: "Require this existing, unexpired token before refreshing or replacing a lock. Set token to the same value to refresh, or a different value to replace.",
				Optional:    true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.RequiresReplace(),
				},
			},
			"depth": schema.StringAttribute{
				Computed: true,
			},
			"scope": schema.StringAttribute{
				Computed: true,
			},
			"type": schema.StringAttribute{
				Computed: true,
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

func (r *lockResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	var plan lockResourceModel
	diags := req.Plan.Get(ctx, &plan)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}
	var config lockResourceModel
	diags = req.Config.Get(ctx, &config)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	paramsLockCreate := files_sdk.LockCreateParams{}
	paramsLockCreate.Path = plan.Path.ValueString()
	paramsLockCreate.Token = plan.Token.ValueString()
	paramsLockCreate.ExpectedToken = plan.ExpectedToken.ValueString()
	if !plan.AllowAccessByAnyUser.IsNull() && !plan.AllowAccessByAnyUser.IsUnknown() {
		paramsLockCreate.AllowAccessByAnyUser = plan.AllowAccessByAnyUser.ValueBoolPointer()
	}
	if !plan.Exclusive.IsNull() && !plan.Exclusive.IsUnknown() {
		paramsLockCreate.Exclusive = plan.Exclusive.ValueBoolPointer()
	}
	if !plan.Recursive.IsNull() && !plan.Recursive.IsUnknown() {
		paramsLockCreate.Recursive = plan.Recursive.ValueBoolPointer()
	}
	paramsLockCreate.Owner = plan.Owner.ValueString()
	paramsLockCreate.Timeout = plan.Timeout.ValueInt64()

	if resp.Diagnostics.HasError() {
		return
	}

	lock, err := r.client.Create(paramsLockCreate, files_sdk.WithContext(ctx))
	if err != nil {
		resp.Diagnostics.AddError(
			"Error Creating Files Lock",
			"Could not create lock, unexpected error: "+err.Error(),
		)
		return
	}

	diags = r.populateResourceModel(ctx, lock, &plan)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	diags = resp.State.Set(ctx, plan)
	resp.Diagnostics.Append(diags...)
}

func (r *lockResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	var state lockResourceModel
	diags := req.State.Get(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	paramsLockListFor := files_sdk.LockListForParams{}
	paramsLockListFor.Path = state.Path.ValueString()

	lockIt, err := r.client.ListFor(paramsLockListFor, files_sdk.WithContext(ctx))
	if err != nil {
		if files_sdk.IsNotExist(err) {
			resp.State.RemoveResource(ctx)
			return
		}

		resp.Diagnostics.AddError(
			"Error Reading Files Lock",
			"Could not read lock path "+fmt.Sprint(state.Path.ValueString())+": "+err.Error(),
		)
		return
	}

	var lock *files_sdk.Lock
	for lockIt.Next() {
		entry := lockIt.Lock()
		if entry.Path == state.Path.ValueString() {
			lock = &entry
			break
		}
	}

	if err = lockIt.Err(); err != nil {
		if files_sdk.IsNotExist(err) {
			resp.State.RemoveResource(ctx)
			return
		}

		resp.Diagnostics.AddError(
			"Error Reading Files Lock",
			"Could not read lock path "+fmt.Sprint(state.Path.ValueString())+": "+err.Error(),
		)
	}

	if lock == nil {
		resp.State.RemoveResource(ctx)
		return
	}

	diags = r.populateResourceModel(ctx, *lock, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	diags = resp.State.Set(ctx, state)
	resp.Diagnostics.Append(diags...)
}

func (r *lockResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	resp.Diagnostics.AddError(
		"Resource Update Not Implemented",
		"This resource does not support updates.",
	)
}

func (r *lockResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	var state lockResourceModel
	diags := req.State.Get(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	paramsLockDelete := files_sdk.LockDeleteParams{}
	paramsLockDelete.Path = state.Path.ValueString()
	paramsLockDelete.Token = state.Token.ValueString()

	err := r.client.Delete(paramsLockDelete, files_sdk.WithContext(ctx))
	if err != nil && !files_sdk.IsNotExist(err) {
		resp.Diagnostics.AddError(
			"Error Deleting Files Lock",
			"Could not delete lock path "+fmt.Sprint(state.Path.ValueString())+": "+err.Error(),
		)
	}
}

func (r *lockResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	idParts := strings.SplitN(req.ID, ",", 1)

	if len(idParts) != 1 || idParts[0] == "" {
		resp.Diagnostics.AddError(
			"Unexpected Import Identifier",
			fmt.Sprintf("Expected import identifier with format: path. Got: %q", req.ID),
		)
		return
	}

	resp.Diagnostics.Append(resp.State.SetAttribute(ctx, path.Root("path"), idParts[0])...)

}

func (r *lockResource) populateResourceModel(ctx context.Context, lock files_sdk.Lock, state *lockResourceModel) (diags diag.Diagnostics) {
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
