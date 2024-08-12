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
	Exclusive            types.Bool   `tfsdk:"exclusive"`
	AllowAccessByAnyUser types.Bool   `tfsdk:"allow_access_by_any_user"`
	Depth                types.String `tfsdk:"depth"`
	Owner                types.String `tfsdk:"owner"`
	Scope                types.String `tfsdk:"scope"`
	Token                types.String `tfsdk:"token"`
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
		Description: "A Lock is not used by Files.com's web interface, but can be used by your applications to\n\nimplement locking and concurrency features. Note that these locks are advisory in nature,\n\nand creating a lock does not prevent other API requests from being fulfilled.\n\n\n\nOur lock feature is designed to emulate the locking feature offered by WebDAV.\n\nYou can read the WebDAV spec and understand how all of the below endpoints work.\n\n\n\nFiles.com's WebDAV offering and desktop app does leverage this locking API.",
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
			"exclusive": schema.BoolAttribute{
				Description: "Is lock exclusive?",
				Computed:    true,
				Optional:    true,
				PlanModifiers: []planmodifier.Bool{
					boolplanmodifier.UseStateForUnknown(),
					boolplanmodifier.RequiresReplace(),
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
			"depth": schema.StringAttribute{
				Computed: true,
			},
			"owner": schema.StringAttribute{
				Description: "Owner of the lock.  This can be any arbitrary string.",
				Computed:    true,
			},
			"scope": schema.StringAttribute{
				Computed: true,
			},
			"token": schema.StringAttribute{
				Description: "Lock token.  Use to release lock.",
				Computed:    true,
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

	paramsLockCreate := files_sdk.LockCreateParams{}
	paramsLockCreate.Path = plan.Path.ValueString()
	if !plan.AllowAccessByAnyUser.IsNull() && !plan.AllowAccessByAnyUser.IsUnknown() {
		paramsLockCreate.AllowAccessByAnyUser = plan.AllowAccessByAnyUser.ValueBoolPointer()
	}
	if !plan.Exclusive.IsNull() && !plan.Exclusive.IsUnknown() {
		paramsLockCreate.Exclusive = plan.Exclusive.ValueBoolPointer()
	}
	if !plan.Recursive.IsNull() && !plan.Recursive.IsUnknown() {
		paramsLockCreate.Recursive = plan.Recursive.ValueBoolPointer()
	}
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
		"Error Updating Files Lock",
		"Update operation not implemented",
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
