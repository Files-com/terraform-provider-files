package provider

import (
	"context"
	"encoding/json"
	"fmt"
	"strconv"
	"strings"

	files_sdk "github.com/Files-com/files-sdk-go/v3"
	behavior "github.com/Files-com/files-sdk-go/v3/behavior"
	"github.com/Files-com/terraform-provider-files/lib"
	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/boolplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/dynamicplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/int64planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

var (
	_ resource.Resource                = &behaviorResource{}
	_ resource.ResourceWithConfigure   = &behaviorResource{}
	_ resource.ResourceWithImportState = &behaviorResource{}
)

func NewBehaviorResource() resource.Resource {
	return &behaviorResource{}
}

type behaviorResource struct {
	client *behavior.Client
}

type behaviorResourceModel struct {
	Id                          types.Int64   `tfsdk:"id"`
	Path                        types.String  `tfsdk:"path"`
	AttachmentUrl               types.String  `tfsdk:"attachment_url"`
	Behavior                    types.String  `tfsdk:"behavior"`
	Name                        types.String  `tfsdk:"name"`
	Description                 types.String  `tfsdk:"description"`
	Value                       types.Dynamic `tfsdk:"value"`
	DisableParentFolderBehavior types.Bool    `tfsdk:"disable_parent_folder_behavior"`
	Recursive                   types.Bool    `tfsdk:"recursive"`
}

func (r *behaviorResource) Configure(_ context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {
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

	r.client = &behavior.Client{Config: sdk_config}
}

func (r *behaviorResource) Metadata(_ context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_behavior"
}

func (r *behaviorResource) Schema(_ context.Context, _ resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Description: "Behaviors are the API resource for what are also known as Folder Settings. Every behavior is associated with a folder.\n\n\n\nDepending on the behavior, it may also operate on child folders. It may be overridable at the child folder level or maybe can be added to at the child folder level. The exact options for each behavior type are explained in the table below.\n\n\n\nAdditionally, some behaviors are visible to non-admins, and others are even settable by non-admins. All the details are below.\n\n\n\nEach behavior uses a different format for storing its settings value. Next to each behavior type is an example value. Our API and SDKs currently require that the value for behaviors be sent as raw JSON within the `value` field. Our SDK generator and API documentation generator doesn't fully keep up with this requirement, so if you need any help finding the exact syntax to use for your language or use case, just reach out.\n\n\n\nNote: Append Timestamp behavior removed. Check [Override Upload Filename](#override-upload-filename-behaviors) behavior which have even more functionality to modify name on upload.",
		Attributes: map[string]schema.Attribute{
			"id": schema.Int64Attribute{
				Description: "Folder behavior ID",
				Computed:    true,
				PlanModifiers: []planmodifier.Int64{
					int64planmodifier.UseStateForUnknown(),
				},
			},
			"path": schema.StringAttribute{
				Description: "Folder path.  Note that Behavior paths cannot be updated once initially set.  You will need to remove and re-create the behavior on the new path. This must be slash-delimited, but it must neither start nor end with a slash. Maximum of 5000 characters.",
				Required:    true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.RequiresReplace(),
				},
			},
			"attachment_url": schema.StringAttribute{
				Description: "URL for attached file",
				Computed:    true,
			},
			"behavior": schema.StringAttribute{
				Description: "Behavior type.",
				Required:    true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.RequiresReplace(),
				},
			},
			"name": schema.StringAttribute{
				Description: "Name for this behavior.",
				Computed:    true,
				Optional:    true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
			},
			"description": schema.StringAttribute{
				Description: "Description for this behavior.",
				Computed:    true,
				Optional:    true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
			},
			"value": schema.DynamicAttribute{
				Description: "Settings for this behavior.  See the section above for an example value to provide here.  Formatting is different for each Behavior type.  May be sent as nested JSON or a single JSON-encoded string.  If using XML encoding for the API call, this data must be sent as a JSON-encoded string.",
				Computed:    true,
				Optional:    true,
				PlanModifiers: []planmodifier.Dynamic{
					dynamicplanmodifier.UseStateForUnknown(),
				},
			},
			"disable_parent_folder_behavior": schema.BoolAttribute{
				Description: "If true, the parent folder's behavior will be disabled for this folder and its children.",
				Computed:    true,
				Optional:    true,
				PlanModifiers: []planmodifier.Bool{
					boolplanmodifier.UseStateForUnknown(),
				},
			},
			"recursive": schema.BoolAttribute{
				Description: "Is behavior recursive?",
				Computed:    true,
				Optional:    true,
				PlanModifiers: []planmodifier.Bool{
					boolplanmodifier.UseStateForUnknown(),
				},
			},
		},
	}
}

func (r *behaviorResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	var plan behaviorResourceModel
	diags := req.Plan.Get(ctx, &plan)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	paramsBehaviorCreate := files_sdk.BehaviorCreateParams{}
	createValue, diags := lib.DynamicToStringMap(ctx, path.Root("value"), plan.Value)
	resp.Diagnostics.Append(diags...)
	createValueBytes, err := json.Marshal(createValue)
	if err != nil {
		resp.Diagnostics.AddAttributeError(
			path.Root("value"),
			"Error Creating Files Behavior",
			"Could not marshal value to JSON: "+err.Error(),
		)
	} else {
		paramsBehaviorCreate.Value = string(createValueBytes)
	}
	paramsBehaviorCreate.DisableParentFolderBehavior = plan.DisableParentFolderBehavior.ValueBoolPointer()
	paramsBehaviorCreate.Recursive = plan.Recursive.ValueBoolPointer()
	paramsBehaviorCreate.Name = plan.Name.ValueString()
	paramsBehaviorCreate.Description = plan.Description.ValueString()
	paramsBehaviorCreate.Path = plan.Path.ValueString()
	paramsBehaviorCreate.Behavior = plan.Behavior.ValueString()

	if resp.Diagnostics.HasError() {
		return
	}

	behavior, err := r.client.Create(paramsBehaviorCreate, files_sdk.WithContext(ctx))
	if err != nil {
		resp.Diagnostics.AddError(
			"Error Creating Files Behavior",
			"Could not create behavior, unexpected error: "+err.Error(),
		)
		return
	}

	diags = r.populateResourceModel(ctx, behavior, &plan)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	diags = resp.State.Set(ctx, plan)
	resp.Diagnostics.Append(diags...)
}

func (r *behaviorResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	var state behaviorResourceModel
	diags := req.State.Get(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	paramsBehaviorFind := files_sdk.BehaviorFindParams{}
	paramsBehaviorFind.Id = state.Id.ValueInt64()

	behavior, err := r.client.Find(paramsBehaviorFind, files_sdk.WithContext(ctx))
	if err != nil {
		resp.Diagnostics.AddError(
			"Error Reading Files Behavior",
			"Could not read behavior id "+fmt.Sprint(state.Id.ValueInt64())+": "+err.Error(),
		)
		return
	}

	diags = r.populateResourceModel(ctx, behavior, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	diags = resp.State.Set(ctx, state)
	resp.Diagnostics.Append(diags...)
}

func (r *behaviorResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	var plan behaviorResourceModel
	diags := req.Plan.Get(ctx, &plan)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	paramsBehaviorUpdate := files_sdk.BehaviorUpdateParams{}
	paramsBehaviorUpdate.Id = plan.Id.ValueInt64()
	updateValue, diags := lib.DynamicToStringMap(ctx, path.Root("value"), plan.Value)
	resp.Diagnostics.Append(diags...)
	updateValueBytes, err := json.Marshal(updateValue)
	if err != nil {
		resp.Diagnostics.AddAttributeError(
			path.Root("value"),
			"Error Creating Files Behavior",
			"Could not marshal value to JSON: "+err.Error(),
		)
	} else {
		paramsBehaviorUpdate.Value = string(updateValueBytes)
	}
	paramsBehaviorUpdate.DisableParentFolderBehavior = plan.DisableParentFolderBehavior.ValueBoolPointer()
	paramsBehaviorUpdate.Recursive = plan.Recursive.ValueBoolPointer()
	paramsBehaviorUpdate.Name = plan.Name.ValueString()
	paramsBehaviorUpdate.Description = plan.Description.ValueString()

	if resp.Diagnostics.HasError() {
		return
	}

	behavior, err := r.client.Update(paramsBehaviorUpdate, files_sdk.WithContext(ctx))
	if err != nil {
		resp.Diagnostics.AddError(
			"Error Updating Files Behavior",
			"Could not update behavior, unexpected error: "+err.Error(),
		)
		return
	}

	diags = r.populateResourceModel(ctx, behavior, &plan)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	diags = resp.State.Set(ctx, plan)
	resp.Diagnostics.Append(diags...)
}

func (r *behaviorResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	var state behaviorResourceModel
	diags := req.State.Get(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	paramsBehaviorDelete := files_sdk.BehaviorDeleteParams{}
	paramsBehaviorDelete.Id = state.Id.ValueInt64()

	err := r.client.Delete(paramsBehaviorDelete, files_sdk.WithContext(ctx))
	if err != nil {
		resp.Diagnostics.AddError(
			"Error Deleting Files Behavior",
			"Could not delete behavior id "+fmt.Sprint(state.Id.ValueInt64())+": "+err.Error(),
		)
	}
}

func (r *behaviorResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
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

func (r *behaviorResource) populateResourceModel(ctx context.Context, behavior files_sdk.Behavior, state *behaviorResourceModel) (diags diag.Diagnostics) {
	var propDiags diag.Diagnostics

	state.Id = types.Int64Value(behavior.Id)
	state.Path = types.StringValue(behavior.Path)
	state.AttachmentUrl = types.StringValue(behavior.AttachmentUrl)
	state.Behavior = types.StringValue(behavior.Behavior)
	state.Name = types.StringValue(behavior.Name)
	state.Description = types.StringValue(behavior.Description)
	state.Value, propDiags = lib.ToDynamic(ctx, path.Root("value"), behavior.Value, state.Value.UnderlyingValue())
	diags.Append(propDiags...)
	state.DisableParentFolderBehavior = types.BoolPointerValue(behavior.DisableParentFolderBehavior)
	state.Recursive = types.BoolPointerValue(behavior.Recursive)

	return
}
