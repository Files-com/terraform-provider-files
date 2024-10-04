package provider

import (
	"context"
	"fmt"
	"strconv"
	"strings"

	files_sdk "github.com/Files-com/files-sdk-go/v3"
	request "github.com/Files-com/files-sdk-go/v3/request"
	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/int64planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

var (
	_ resource.Resource                = &requestResource{}
	_ resource.ResourceWithConfigure   = &requestResource{}
	_ resource.ResourceWithImportState = &requestResource{}
)

func NewRequestResource() resource.Resource {
	return &requestResource{}
}

type requestResource struct {
	client *request.Client
}

type requestResourceModel struct {
	Path            types.String `tfsdk:"path"`
	Destination     types.String `tfsdk:"destination"`
	UserIds         types.String `tfsdk:"user_ids"`
	GroupIds        types.String `tfsdk:"group_ids"`
	Id              types.Int64  `tfsdk:"id"`
	Source          types.String `tfsdk:"source"`
	AutomationId    types.String `tfsdk:"automation_id"`
	UserDisplayName types.String `tfsdk:"user_display_name"`
}

func (r *requestResource) Configure(_ context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {
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

	r.client = &request.Client{Config: sdk_config}
}

func (r *requestResource) Metadata(_ context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_request"
}

func (r *requestResource) Schema(_ context.Context, _ resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Description: "A Request is a file that *should* be uploaded by a specific user or group.\n\n\n\nRequests can either be manually created and managed, or managed automatically by an Automation.",
		Attributes: map[string]schema.Attribute{
			"path": schema.StringAttribute{
				Description: "Folder path. This must be slash-delimited, but it must neither start nor end with a slash. Maximum of 5000 characters.",
				Required:    true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.RequiresReplace(),
				},
			},
			"destination": schema.StringAttribute{
				Description: "Destination filename",
				Required:    true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.RequiresReplace(),
				},
			},
			"user_ids": schema.StringAttribute{
				Description: "A list of user IDs to request the file from. If sent as a string, it should be comma-delimited.",
				Optional:    true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.RequiresReplace(),
				},
			},
			"group_ids": schema.StringAttribute{
				Description: "A list of group IDs to request the file from. If sent as a string, it should be comma-delimited.",
				Optional:    true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.RequiresReplace(),
				},
			},
			"id": schema.Int64Attribute{
				Description: "Request ID",
				Computed:    true,
				PlanModifiers: []planmodifier.Int64{
					int64planmodifier.UseStateForUnknown(),
				},
			},
			"source": schema.StringAttribute{
				Description: "Source filename, if applicable",
				Computed:    true,
			},
			"automation_id": schema.StringAttribute{
				Description: "ID of automation that created request",
				Computed:    true,
			},
			"user_display_name": schema.StringAttribute{
				Description: "User making the request (if applicable)",
				Computed:    true,
			},
		},
	}
}

func (r *requestResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	var plan requestResourceModel
	diags := req.Plan.Get(ctx, &plan)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	paramsRequestCreate := files_sdk.RequestCreateParams{}
	paramsRequestCreate.Path = plan.Path.ValueString()
	paramsRequestCreate.Destination = plan.Destination.ValueString()
	paramsRequestCreate.UserIds = plan.UserIds.ValueString()
	paramsRequestCreate.GroupIds = plan.GroupIds.ValueString()

	if resp.Diagnostics.HasError() {
		return
	}

	request, err := r.client.Create(paramsRequestCreate, files_sdk.WithContext(ctx))
	if err != nil {
		resp.Diagnostics.AddError(
			"Error Creating Files Request",
			"Could not create request, unexpected error: "+err.Error(),
		)
		return
	}

	diags = r.populateResourceModel(ctx, request, &plan)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	diags = resp.State.Set(ctx, plan)
	resp.Diagnostics.Append(diags...)
}

func (r *requestResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	var state requestResourceModel
	diags := req.State.Get(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	paramsRequestList := files_sdk.RequestListParams{}

	requestIt, err := r.client.List(paramsRequestList, files_sdk.WithContext(ctx))
	if err != nil {
		if files_sdk.IsNotExist(err) {
			resp.State.RemoveResource(ctx)
			return
		}

		resp.Diagnostics.AddError(
			"Error Reading Files Request",
			"Could not read request id "+fmt.Sprint(state.Id.ValueInt64())+": "+err.Error(),
		)
		return
	}

	var request *files_sdk.Request
	for requestIt.Next() {
		entry := requestIt.Request()
		if entry.Id == state.Id.ValueInt64() {
			request = &entry
			break
		}
	}

	if err = requestIt.Err(); err != nil {
		if files_sdk.IsNotExist(err) {
			resp.State.RemoveResource(ctx)
			return
		}

		resp.Diagnostics.AddError(
			"Error Reading Files Request",
			"Could not read request id "+fmt.Sprint(state.Id.ValueInt64())+": "+err.Error(),
		)
	}

	if request == nil {
		resp.State.RemoveResource(ctx)
		return
	}

	diags = r.populateResourceModel(ctx, *request, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	diags = resp.State.Set(ctx, state)
	resp.Diagnostics.Append(diags...)
}

func (r *requestResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	resp.Diagnostics.AddError(
		"Resource Update Not Implemented",
		"This resource does not support updates.",
	)
}

func (r *requestResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	var state requestResourceModel
	diags := req.State.Get(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	paramsRequestDelete := files_sdk.RequestDeleteParams{}
	paramsRequestDelete.Id = state.Id.ValueInt64()

	err := r.client.Delete(paramsRequestDelete, files_sdk.WithContext(ctx))
	if err != nil && !files_sdk.IsNotExist(err) {
		resp.Diagnostics.AddError(
			"Error Deleting Files Request",
			"Could not delete request id "+fmt.Sprint(state.Id.ValueInt64())+": "+err.Error(),
		)
	}
}

func (r *requestResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
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

func (r *requestResource) populateResourceModel(ctx context.Context, request files_sdk.Request, state *requestResourceModel) (diags diag.Diagnostics) {
	state.Id = types.Int64Value(request.Id)
	state.Path = types.StringValue(request.Path)
	state.Source = types.StringValue(request.Source)
	state.Destination = types.StringValue(request.Destination)
	state.AutomationId = types.StringValue(request.AutomationId)
	state.UserDisplayName = types.StringValue(request.UserDisplayName)

	return
}
