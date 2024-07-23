package provider

import (
	"context"
	"fmt"
	"strconv"
	"strings"

	files_sdk "github.com/Files-com/files-sdk-go/v3"
	user_request "github.com/Files-com/files-sdk-go/v3/userrequest"
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
	_ resource.Resource                = &userRequestResource{}
	_ resource.ResourceWithConfigure   = &userRequestResource{}
	_ resource.ResourceWithImportState = &userRequestResource{}
)

func NewUserRequestResource() resource.Resource {
	return &userRequestResource{}
}

type userRequestResource struct {
	client *user_request.Client
}

type userRequestResourceModel struct {
	Name    types.String `tfsdk:"name"`
	Email   types.String `tfsdk:"email"`
	Details types.String `tfsdk:"details"`
	Company types.String `tfsdk:"company"`
	Id      types.Int64  `tfsdk:"id"`
}

func (r *userRequestResource) Configure(_ context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {
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

	r.client = &user_request.Client{Config: sdk_config}
}

func (r *userRequestResource) Metadata(_ context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_user_request"
}

func (r *userRequestResource) Schema(_ context.Context, _ resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Description: "User Requests allow anonymous users to place a request for access on the login screen to the site administrator.",
		Attributes: map[string]schema.Attribute{
			"name": schema.StringAttribute{
				Description: "User's full name",
				Required:    true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.RequiresReplace(),
				},
			},
			"email": schema.StringAttribute{
				Description: "User email address",
				Required:    true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.RequiresReplace(),
				},
			},
			"details": schema.StringAttribute{
				Description: "Details of the user's request",
				Required:    true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.RequiresReplace(),
				},
			},
			"company": schema.StringAttribute{
				Description: "User's company name",
				Computed:    true,
				Optional:    true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
					stringplanmodifier.RequiresReplace(),
				},
			},
			"id": schema.Int64Attribute{
				Description: "ID",
				Computed:    true,
				PlanModifiers: []planmodifier.Int64{
					int64planmodifier.UseStateForUnknown(),
				},
			},
		},
	}
}

func (r *userRequestResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	var plan userRequestResourceModel
	diags := req.Plan.Get(ctx, &plan)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	paramsUserRequestCreate := files_sdk.UserRequestCreateParams{}
	paramsUserRequestCreate.Name = plan.Name.ValueString()
	paramsUserRequestCreate.Email = plan.Email.ValueString()
	paramsUserRequestCreate.Details = plan.Details.ValueString()
	paramsUserRequestCreate.Company = plan.Company.ValueString()

	if resp.Diagnostics.HasError() {
		return
	}

	userRequest, err := r.client.Create(paramsUserRequestCreate, files_sdk.WithContext(ctx))
	if err != nil {
		resp.Diagnostics.AddError(
			"Error Creating Files UserRequest",
			"Could not create user_request, unexpected error: "+err.Error(),
		)
		return
	}

	diags = r.populateResourceModel(ctx, userRequest, &plan)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	diags = resp.State.Set(ctx, plan)
	resp.Diagnostics.Append(diags...)
}

func (r *userRequestResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	var state userRequestResourceModel
	diags := req.State.Get(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	paramsUserRequestFind := files_sdk.UserRequestFindParams{}
	paramsUserRequestFind.Id = state.Id.ValueInt64()

	userRequest, err := r.client.Find(paramsUserRequestFind, files_sdk.WithContext(ctx))
	if err != nil {
		if files_sdk.IsNotExist(err) {
			resp.State.RemoveResource(ctx)
			return
		}

		resp.Diagnostics.AddError(
			"Error Reading Files UserRequest",
			"Could not read user_request id "+fmt.Sprint(state.Id.ValueInt64())+": "+err.Error(),
		)
		return
	}

	diags = r.populateResourceModel(ctx, userRequest, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	diags = resp.State.Set(ctx, state)
	resp.Diagnostics.Append(diags...)
}

func (r *userRequestResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	resp.Diagnostics.AddError(
		"Error Updating Files UserRequest",
		"Update operation not implemented",
	)
}

func (r *userRequestResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	var state userRequestResourceModel
	diags := req.State.Get(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	paramsUserRequestDelete := files_sdk.UserRequestDeleteParams{}
	paramsUserRequestDelete.Id = state.Id.ValueInt64()

	err := r.client.Delete(paramsUserRequestDelete, files_sdk.WithContext(ctx))
	if err != nil && !files_sdk.IsNotExist(err) {
		resp.Diagnostics.AddError(
			"Error Deleting Files UserRequest",
			"Could not delete user_request id "+fmt.Sprint(state.Id.ValueInt64())+": "+err.Error(),
		)
	}
}

func (r *userRequestResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
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

func (r *userRequestResource) populateResourceModel(ctx context.Context, userRequest files_sdk.UserRequest, state *userRequestResourceModel) (diags diag.Diagnostics) {
	state.Id = types.Int64Value(userRequest.Id)
	state.Name = types.StringValue(userRequest.Name)
	state.Email = types.StringValue(userRequest.Email)
	state.Details = types.StringValue(userRequest.Details)
	state.Company = types.StringValue(userRequest.Company)

	return
}
