package provider

import (
	"context"
	"fmt"
	"strconv"
	"strings"

	files_sdk "github.com/Files-com/files-sdk-go/v3"
	message "github.com/Files-com/files-sdk-go/v3/message"
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
	_ resource.Resource                = &messageResource{}
	_ resource.ResourceWithConfigure   = &messageResource{}
	_ resource.ResourceWithImportState = &messageResource{}
)

func NewMessageResource() resource.Resource {
	return &messageResource{}
}

type messageResource struct {
	client *message.Client
}

type messageResourceModel struct {
	Subject   types.String  `tfsdk:"subject"`
	Body      types.String  `tfsdk:"body"`
	ProjectId types.Int64   `tfsdk:"project_id"`
	UserId    types.Int64   `tfsdk:"user_id"`
	Id        types.Int64   `tfsdk:"id"`
	Comments  types.Dynamic `tfsdk:"comments"`
}

func (r *messageResource) Configure(_ context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {
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

	r.client = &message.Client{Config: sdk_config}
}

func (r *messageResource) Metadata(_ context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_message"
}

func (r *messageResource) Schema(_ context.Context, _ resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Description: "A Messages is a part of Files.com's project management features and represent a message posted by a user to a project.",
		Attributes: map[string]schema.Attribute{
			"subject": schema.StringAttribute{
				Description: "Message subject.",
				Required:    true,
			},
			"body": schema.StringAttribute{
				Description: "Message body.",
				Required:    true,
			},
			"project_id": schema.Int64Attribute{
				Description: "Project to which the message should be attached.",
				Required:    true,
			},
			"user_id": schema.Int64Attribute{
				Description: "User ID.  Provide a value of `0` to operate the current session's user.",
				Optional:    true,
				PlanModifiers: []planmodifier.Int64{
					int64planmodifier.RequiresReplace(),
				},
			},
			"id": schema.Int64Attribute{
				Description: "Message ID",
				Computed:    true,
				PlanModifiers: []planmodifier.Int64{
					int64planmodifier.UseStateForUnknown(),
				},
			},
			"comments": schema.DynamicAttribute{
				Description: "Comments.",
				Computed:    true,
			},
		},
	}
}

func (r *messageResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	var plan messageResourceModel
	diags := req.Plan.Get(ctx, &plan)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	paramsMessageCreate := files_sdk.MessageCreateParams{}
	paramsMessageCreate.UserId = plan.UserId.ValueInt64()
	paramsMessageCreate.ProjectId = plan.ProjectId.ValueInt64()
	paramsMessageCreate.Subject = plan.Subject.ValueString()
	paramsMessageCreate.Body = plan.Body.ValueString()

	if resp.Diagnostics.HasError() {
		return
	}

	message, err := r.client.Create(paramsMessageCreate, files_sdk.WithContext(ctx))
	if err != nil {
		resp.Diagnostics.AddError(
			"Error Creating Files Message",
			"Could not create message, unexpected error: "+err.Error(),
		)
		return
	}

	diags = r.populateResourceModel(ctx, message, &plan)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	diags = resp.State.Set(ctx, plan)
	resp.Diagnostics.Append(diags...)
}

func (r *messageResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	var state messageResourceModel
	diags := req.State.Get(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	paramsMessageFind := files_sdk.MessageFindParams{}
	paramsMessageFind.Id = state.Id.ValueInt64()

	message, err := r.client.Find(paramsMessageFind, files_sdk.WithContext(ctx))
	if err != nil {
		if files_sdk.IsNotExist(err) {
			resp.State.RemoveResource(ctx)
			return
		}

		resp.Diagnostics.AddError(
			"Error Reading Files Message",
			"Could not read message id "+fmt.Sprint(state.Id.ValueInt64())+": "+err.Error(),
		)
		return
	}

	diags = r.populateResourceModel(ctx, message, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	diags = resp.State.Set(ctx, state)
	resp.Diagnostics.Append(diags...)
}

func (r *messageResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	var plan messageResourceModel
	diags := req.Plan.Get(ctx, &plan)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	paramsMessageUpdate := files_sdk.MessageUpdateParams{}
	paramsMessageUpdate.Id = plan.Id.ValueInt64()
	paramsMessageUpdate.ProjectId = plan.ProjectId.ValueInt64()
	paramsMessageUpdate.Subject = plan.Subject.ValueString()
	paramsMessageUpdate.Body = plan.Body.ValueString()

	if resp.Diagnostics.HasError() {
		return
	}

	message, err := r.client.Update(paramsMessageUpdate, files_sdk.WithContext(ctx))
	if err != nil {
		resp.Diagnostics.AddError(
			"Error Updating Files Message",
			"Could not update message, unexpected error: "+err.Error(),
		)
		return
	}

	diags = r.populateResourceModel(ctx, message, &plan)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	diags = resp.State.Set(ctx, plan)
	resp.Diagnostics.Append(diags...)
}

func (r *messageResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	var state messageResourceModel
	diags := req.State.Get(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	paramsMessageDelete := files_sdk.MessageDeleteParams{}
	paramsMessageDelete.Id = state.Id.ValueInt64()

	err := r.client.Delete(paramsMessageDelete, files_sdk.WithContext(ctx))
	if err != nil && !files_sdk.IsNotExist(err) {
		resp.Diagnostics.AddError(
			"Error Deleting Files Message",
			"Could not delete message id "+fmt.Sprint(state.Id.ValueInt64())+": "+err.Error(),
		)
	}
}

func (r *messageResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
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

func (r *messageResource) populateResourceModel(ctx context.Context, message files_sdk.Message, state *messageResourceModel) (diags diag.Diagnostics) {
	var propDiags diag.Diagnostics

	state.Id = types.Int64Value(message.Id)
	state.Subject = types.StringValue(message.Subject)
	state.Body = types.StringValue(message.Body)
	state.Comments, propDiags = lib.ToDynamic(ctx, path.Root("comments"), message.Comments, state.Comments.UnderlyingValue())
	diags.Append(propDiags...)

	return
}
