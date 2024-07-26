package provider

import (
	"context"
	"fmt"
	"strconv"
	"strings"

	files_sdk "github.com/Files-com/files-sdk-go/v3"
	bundle_notification "github.com/Files-com/files-sdk-go/v3/bundlenotification"
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
	_ resource.Resource                = &bundleNotificationResource{}
	_ resource.ResourceWithConfigure   = &bundleNotificationResource{}
	_ resource.ResourceWithImportState = &bundleNotificationResource{}
)

func NewBundleNotificationResource() resource.Resource {
	return &bundleNotificationResource{}
}

type bundleNotificationResource struct {
	client *bundle_notification.Client
}

type bundleNotificationResourceModel struct {
	BundleId             types.Int64 `tfsdk:"bundle_id"`
	NotifyOnRegistration types.Bool  `tfsdk:"notify_on_registration"`
	NotifyOnUpload       types.Bool  `tfsdk:"notify_on_upload"`
	UserId               types.Int64 `tfsdk:"user_id"`
	Id                   types.Int64 `tfsdk:"id"`
}

func (r *bundleNotificationResource) Configure(_ context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {
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

	r.client = &bundle_notification.Client{Config: sdk_config}
}

func (r *bundleNotificationResource) Metadata(_ context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_bundle_notification"
}

func (r *bundleNotificationResource) Schema(_ context.Context, _ resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Description: "Bundle notifications are emails sent out to users when certain actions are performed on or within a shared set of files and folders.",
		Attributes: map[string]schema.Attribute{
			"bundle_id": schema.Int64Attribute{
				Description: "Bundle ID to notify on",
				Required:    true,
				PlanModifiers: []planmodifier.Int64{
					int64planmodifier.RequiresReplace(),
				},
			},
			"notify_on_registration": schema.BoolAttribute{
				Description: "Triggers bundle notification when a registration action occurs for it.",
				Computed:    true,
				Optional:    true,
				PlanModifiers: []planmodifier.Bool{
					boolplanmodifier.UseStateForUnknown(),
				},
			},
			"notify_on_upload": schema.BoolAttribute{
				Description: "Triggers bundle notification when a upload action occurs for it.",
				Computed:    true,
				Optional:    true,
				PlanModifiers: []planmodifier.Bool{
					boolplanmodifier.UseStateForUnknown(),
				},
			},
			"user_id": schema.Int64Attribute{
				Description: "The id of the user to notify.",
				Computed:    true,
				Optional:    true,
				PlanModifiers: []planmodifier.Int64{
					int64planmodifier.UseStateForUnknown(),
					int64planmodifier.RequiresReplace(),
				},
			},
			"id": schema.Int64Attribute{
				Description: "Bundle Notification ID",
				Computed:    true,
				PlanModifiers: []planmodifier.Int64{
					int64planmodifier.UseStateForUnknown(),
				},
			},
		},
	}
}

func (r *bundleNotificationResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	var plan bundleNotificationResourceModel
	diags := req.Plan.Get(ctx, &plan)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	paramsBundleNotificationCreate := files_sdk.BundleNotificationCreateParams{}
	paramsBundleNotificationCreate.BundleId = plan.BundleId.ValueInt64()
	paramsBundleNotificationCreate.UserId = plan.UserId.ValueInt64()
	if !plan.NotifyOnRegistration.IsNull() && !plan.NotifyOnRegistration.IsUnknown() {
		paramsBundleNotificationCreate.NotifyOnRegistration = plan.NotifyOnRegistration.ValueBoolPointer()
	}
	if !plan.NotifyOnUpload.IsNull() && !plan.NotifyOnUpload.IsUnknown() {
		paramsBundleNotificationCreate.NotifyOnUpload = plan.NotifyOnUpload.ValueBoolPointer()
	}

	if resp.Diagnostics.HasError() {
		return
	}

	bundleNotification, err := r.client.Create(paramsBundleNotificationCreate, files_sdk.WithContext(ctx))
	if err != nil {
		resp.Diagnostics.AddError(
			"Error Creating Files BundleNotification",
			"Could not create bundle_notification, unexpected error: "+err.Error(),
		)
		return
	}

	diags = r.populateResourceModel(ctx, bundleNotification, &plan)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	diags = resp.State.Set(ctx, plan)
	resp.Diagnostics.Append(diags...)
}

func (r *bundleNotificationResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	var state bundleNotificationResourceModel
	diags := req.State.Get(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	paramsBundleNotificationFind := files_sdk.BundleNotificationFindParams{}
	paramsBundleNotificationFind.Id = state.Id.ValueInt64()

	bundleNotification, err := r.client.Find(paramsBundleNotificationFind, files_sdk.WithContext(ctx))
	if err != nil {
		if files_sdk.IsNotExist(err) {
			resp.State.RemoveResource(ctx)
			return
		}

		resp.Diagnostics.AddError(
			"Error Reading Files BundleNotification",
			"Could not read bundle_notification id "+fmt.Sprint(state.Id.ValueInt64())+": "+err.Error(),
		)
		return
	}

	diags = r.populateResourceModel(ctx, bundleNotification, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	diags = resp.State.Set(ctx, state)
	resp.Diagnostics.Append(diags...)
}

func (r *bundleNotificationResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	var plan bundleNotificationResourceModel
	diags := req.Plan.Get(ctx, &plan)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	paramsBundleNotificationUpdate := files_sdk.BundleNotificationUpdateParams{}
	paramsBundleNotificationUpdate.Id = plan.Id.ValueInt64()
	if !plan.NotifyOnRegistration.IsNull() && !plan.NotifyOnRegistration.IsUnknown() {
		paramsBundleNotificationUpdate.NotifyOnRegistration = plan.NotifyOnRegistration.ValueBoolPointer()
	}
	if !plan.NotifyOnUpload.IsNull() && !plan.NotifyOnUpload.IsUnknown() {
		paramsBundleNotificationUpdate.NotifyOnUpload = plan.NotifyOnUpload.ValueBoolPointer()
	}

	if resp.Diagnostics.HasError() {
		return
	}

	bundleNotification, err := r.client.Update(paramsBundleNotificationUpdate, files_sdk.WithContext(ctx))
	if err != nil {
		resp.Diagnostics.AddError(
			"Error Updating Files BundleNotification",
			"Could not update bundle_notification, unexpected error: "+err.Error(),
		)
		return
	}

	diags = r.populateResourceModel(ctx, bundleNotification, &plan)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	diags = resp.State.Set(ctx, plan)
	resp.Diagnostics.Append(diags...)
}

func (r *bundleNotificationResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	var state bundleNotificationResourceModel
	diags := req.State.Get(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	paramsBundleNotificationDelete := files_sdk.BundleNotificationDeleteParams{}
	paramsBundleNotificationDelete.Id = state.Id.ValueInt64()

	err := r.client.Delete(paramsBundleNotificationDelete, files_sdk.WithContext(ctx))
	if err != nil && !files_sdk.IsNotExist(err) {
		resp.Diagnostics.AddError(
			"Error Deleting Files BundleNotification",
			"Could not delete bundle_notification id "+fmt.Sprint(state.Id.ValueInt64())+": "+err.Error(),
		)
	}
}

func (r *bundleNotificationResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
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

func (r *bundleNotificationResource) populateResourceModel(ctx context.Context, bundleNotification files_sdk.BundleNotification, state *bundleNotificationResourceModel) (diags diag.Diagnostics) {
	state.BundleId = types.Int64Value(bundleNotification.BundleId)
	state.Id = types.Int64Value(bundleNotification.Id)
	state.NotifyOnRegistration = types.BoolPointerValue(bundleNotification.NotifyOnRegistration)
	state.NotifyOnUpload = types.BoolPointerValue(bundleNotification.NotifyOnUpload)
	state.UserId = types.Int64Value(bundleNotification.UserId)

	return
}
