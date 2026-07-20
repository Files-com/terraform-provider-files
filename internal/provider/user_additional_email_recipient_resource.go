package provider

import (
	"context"
	"fmt"
	"strconv"
	"strings"

	files_sdk "github.com/Files-com/files-sdk-go/v3"
	user_additional_email_recipient "github.com/Files-com/files-sdk-go/v3/useradditionalemailrecipient"
	"github.com/Files-com/terraform-provider-files/lib"
	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/int64planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

var (
	_ resource.Resource                = &userAdditionalEmailRecipientResource{}
	_ resource.ResourceWithConfigure   = &userAdditionalEmailRecipientResource{}
	_ resource.ResourceWithImportState = &userAdditionalEmailRecipientResource{}
)

func NewUserAdditionalEmailRecipientResource() resource.Resource {
	return &userAdditionalEmailRecipientResource{}
}

type userAdditionalEmailRecipientResource struct {
	client *user_additional_email_recipient.Client
}

type userAdditionalEmailRecipientResourceModel struct {
	Email       types.String `tfsdk:"email"`
	UserId      types.Int64  `tfsdk:"user_id"`
	Id          types.Int64  `tfsdk:"id"`
	WorkspaceId types.Int64  `tfsdk:"workspace_id"`
	CreatedAt   types.String `tfsdk:"created_at"`
}

func (r *userAdditionalEmailRecipientResource) Configure(_ context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {
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

	r.client = &user_additional_email_recipient.Client{Config: sdk_config}
}

func (r *userAdditionalEmailRecipientResource) Metadata(_ context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_user_additional_email_recipient"
}

func (r *userAdditionalEmailRecipientResource) Schema(_ context.Context, _ resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Description: "",
		Attributes: map[string]schema.Attribute{
			"email": schema.StringAttribute{
				Description: "Additional email recipient address",
				Required:    true,
			},
			"user_id": schema.Int64Attribute{
				Description: "User ID",
				Computed:    true,
				Optional:    true,
				PlanModifiers: []planmodifier.Int64{
					int64planmodifier.UseStateForUnknown(),
					int64planmodifier.RequiresReplace(),
				},
			},
			"id": schema.Int64Attribute{
				Description: "User additional email recipient ID",
				Computed:    true,
				PlanModifiers: []planmodifier.Int64{
					int64planmodifier.UseStateForUnknown(),
				},
			},
			"workspace_id": schema.Int64Attribute{
				Description: "Workspace ID (0 for default workspace).",
				Computed:    true,
			},
			"created_at": schema.StringAttribute{
				Description: "Created at date/time",
				Computed:    true,
			},
		},
	}
}

func (r *userAdditionalEmailRecipientResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	var plan userAdditionalEmailRecipientResourceModel
	diags := req.Plan.Get(ctx, &plan)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}
	var config userAdditionalEmailRecipientResourceModel
	diags = req.Config.Get(ctx, &config)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	paramsUserAdditionalEmailRecipientCreate := files_sdk.UserAdditionalEmailRecipientCreateParams{}
	paramsUserAdditionalEmailRecipientCreate.UserId = plan.UserId.ValueInt64()
	paramsUserAdditionalEmailRecipientCreate.Email = plan.Email.ValueString()

	if resp.Diagnostics.HasError() {
		return
	}

	userAdditionalEmailRecipient, err := r.client.Create(paramsUserAdditionalEmailRecipientCreate, files_sdk.WithContext(ctx))
	if err != nil {
		resp.Diagnostics.AddError(
			"Error Creating Files UserAdditionalEmailRecipient",
			"Could not create user_additional_email_recipient, unexpected error: "+err.Error(),
		)
		return
	}

	diags = r.populateResourceModel(ctx, userAdditionalEmailRecipient, &plan)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	diags = resp.State.Set(ctx, plan)
	resp.Diagnostics.Append(diags...)
}

func (r *userAdditionalEmailRecipientResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	var state userAdditionalEmailRecipientResourceModel
	diags := req.State.Get(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	paramsUserAdditionalEmailRecipientFind := files_sdk.UserAdditionalEmailRecipientFindParams{}
	paramsUserAdditionalEmailRecipientFind.Id = state.Id.ValueInt64()

	userAdditionalEmailRecipient, err := r.client.Find(paramsUserAdditionalEmailRecipientFind, files_sdk.WithContext(ctx))
	if err != nil {
		if files_sdk.IsNotExist(err) {
			resp.State.RemoveResource(ctx)
			return
		}

		resp.Diagnostics.AddError(
			"Error Reading Files UserAdditionalEmailRecipient",
			"Could not read user_additional_email_recipient id "+fmt.Sprint(state.Id.ValueInt64())+": "+err.Error(),
		)
		return
	}

	diags = r.populateResourceModel(ctx, userAdditionalEmailRecipient, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	diags = resp.State.Set(ctx, state)
	resp.Diagnostics.Append(diags...)
}

func (r *userAdditionalEmailRecipientResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	var plan userAdditionalEmailRecipientResourceModel
	diags := req.Plan.Get(ctx, &plan)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}
	var config userAdditionalEmailRecipientResourceModel
	diags = req.Config.Get(ctx, &config)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	paramsUserAdditionalEmailRecipientUpdate := map[string]interface{}{}
	if !plan.Id.IsNull() && !plan.Id.IsUnknown() {
		paramsUserAdditionalEmailRecipientUpdate["id"] = plan.Id.ValueInt64()
	}
	if !config.Email.IsNull() && !config.Email.IsUnknown() {
		paramsUserAdditionalEmailRecipientUpdate["email"] = config.Email.ValueString()
	}

	if resp.Diagnostics.HasError() {
		return
	}

	userAdditionalEmailRecipient, err := r.client.UpdateWithMap(paramsUserAdditionalEmailRecipientUpdate, files_sdk.WithContext(ctx))
	if err != nil {
		resp.Diagnostics.AddError(
			"Error Updating Files UserAdditionalEmailRecipient",
			"Could not update user_additional_email_recipient, unexpected error: "+err.Error(),
		)
		return
	}

	diags = r.populateResourceModel(ctx, userAdditionalEmailRecipient, &plan)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	diags = resp.State.Set(ctx, plan)
	resp.Diagnostics.Append(diags...)
}

func (r *userAdditionalEmailRecipientResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	var state userAdditionalEmailRecipientResourceModel
	diags := req.State.Get(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	paramsUserAdditionalEmailRecipientDelete := files_sdk.UserAdditionalEmailRecipientDeleteParams{}
	paramsUserAdditionalEmailRecipientDelete.Id = state.Id.ValueInt64()

	err := r.client.Delete(paramsUserAdditionalEmailRecipientDelete, files_sdk.WithContext(ctx))
	if err != nil && !files_sdk.IsNotExist(err) {
		resp.Diagnostics.AddError(
			"Error Deleting Files UserAdditionalEmailRecipient",
			"Could not delete user_additional_email_recipient id "+fmt.Sprint(state.Id.ValueInt64())+": "+err.Error(),
		)
	}
}

func (r *userAdditionalEmailRecipientResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
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

func (r *userAdditionalEmailRecipientResource) populateResourceModel(ctx context.Context, userAdditionalEmailRecipient files_sdk.UserAdditionalEmailRecipient, state *userAdditionalEmailRecipientResourceModel) (diags diag.Diagnostics) {
	state.Id = types.Int64Value(userAdditionalEmailRecipient.Id)
	state.UserId = types.Int64Value(userAdditionalEmailRecipient.UserId)
	state.WorkspaceId = types.Int64Value(userAdditionalEmailRecipient.WorkspaceId)
	state.Email = types.StringValue(userAdditionalEmailRecipient.Email)
	if err := lib.TimeToStringType(ctx, path.Root("created_at"), userAdditionalEmailRecipient.CreatedAt, &state.CreatedAt); err != nil {
		diags.AddError(
			"Error Creating Files UserAdditionalEmailRecipient",
			"Could not convert state created_at to string: "+err.Error(),
		)
	}

	return
}
