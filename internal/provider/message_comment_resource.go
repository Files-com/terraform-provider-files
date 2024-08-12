package provider

import (
	"context"
	"fmt"
	"strconv"
	"strings"

	files_sdk "github.com/Files-com/files-sdk-go/v3"
	message_comment "github.com/Files-com/files-sdk-go/v3/messagecomment"
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
	_ resource.Resource                = &messageCommentResource{}
	_ resource.ResourceWithConfigure   = &messageCommentResource{}
	_ resource.ResourceWithImportState = &messageCommentResource{}
)

func NewMessageCommentResource() resource.Resource {
	return &messageCommentResource{}
}

type messageCommentResource struct {
	client *message_comment.Client
}

type messageCommentResourceModel struct {
	Body      types.String  `tfsdk:"body"`
	UserId    types.Int64   `tfsdk:"user_id"`
	Id        types.Int64   `tfsdk:"id"`
	Reactions types.Dynamic `tfsdk:"reactions"`
}

func (r *messageCommentResource) Configure(_ context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {
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

	r.client = &message_comment.Client{Config: sdk_config}
}

func (r *messageCommentResource) Metadata(_ context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_message_comment"
}

func (r *messageCommentResource) Schema(_ context.Context, _ resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Description: "A MessageComment is a comment made by a user on a message.",
		Attributes: map[string]schema.Attribute{
			"body": schema.StringAttribute{
				Description: "Comment body.",
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
				Description: "Message Comment ID",
				Computed:    true,
				PlanModifiers: []planmodifier.Int64{
					int64planmodifier.UseStateForUnknown(),
				},
			},
			"reactions": schema.DynamicAttribute{
				Description: "Reactions to this comment.",
				Computed:    true,
			},
		},
	}
}

func (r *messageCommentResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	var plan messageCommentResourceModel
	diags := req.Plan.Get(ctx, &plan)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	paramsMessageCommentCreate := files_sdk.MessageCommentCreateParams{}
	paramsMessageCommentCreate.UserId = plan.UserId.ValueInt64()
	paramsMessageCommentCreate.Body = plan.Body.ValueString()

	if resp.Diagnostics.HasError() {
		return
	}

	messageComment, err := r.client.Create(paramsMessageCommentCreate, files_sdk.WithContext(ctx))
	if err != nil {
		resp.Diagnostics.AddError(
			"Error Creating Files MessageComment",
			"Could not create message_comment, unexpected error: "+err.Error(),
		)
		return
	}

	diags = r.populateResourceModel(ctx, messageComment, &plan)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	diags = resp.State.Set(ctx, plan)
	resp.Diagnostics.Append(diags...)
}

func (r *messageCommentResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	var state messageCommentResourceModel
	diags := req.State.Get(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	paramsMessageCommentFind := files_sdk.MessageCommentFindParams{}
	paramsMessageCommentFind.Id = state.Id.ValueInt64()

	messageComment, err := r.client.Find(paramsMessageCommentFind, files_sdk.WithContext(ctx))
	if err != nil {
		if files_sdk.IsNotExist(err) {
			resp.State.RemoveResource(ctx)
			return
		}

		resp.Diagnostics.AddError(
			"Error Reading Files MessageComment",
			"Could not read message_comment id "+fmt.Sprint(state.Id.ValueInt64())+": "+err.Error(),
		)
		return
	}

	diags = r.populateResourceModel(ctx, messageComment, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	diags = resp.State.Set(ctx, state)
	resp.Diagnostics.Append(diags...)
}

func (r *messageCommentResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	var plan messageCommentResourceModel
	diags := req.Plan.Get(ctx, &plan)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	paramsMessageCommentUpdate := files_sdk.MessageCommentUpdateParams{}
	paramsMessageCommentUpdate.Id = plan.Id.ValueInt64()
	paramsMessageCommentUpdate.Body = plan.Body.ValueString()

	if resp.Diagnostics.HasError() {
		return
	}

	messageComment, err := r.client.Update(paramsMessageCommentUpdate, files_sdk.WithContext(ctx))
	if err != nil {
		resp.Diagnostics.AddError(
			"Error Updating Files MessageComment",
			"Could not update message_comment, unexpected error: "+err.Error(),
		)
		return
	}

	diags = r.populateResourceModel(ctx, messageComment, &plan)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	diags = resp.State.Set(ctx, plan)
	resp.Diagnostics.Append(diags...)
}

func (r *messageCommentResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	var state messageCommentResourceModel
	diags := req.State.Get(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	paramsMessageCommentDelete := files_sdk.MessageCommentDeleteParams{}
	paramsMessageCommentDelete.Id = state.Id.ValueInt64()

	err := r.client.Delete(paramsMessageCommentDelete, files_sdk.WithContext(ctx))
	if err != nil && !files_sdk.IsNotExist(err) {
		resp.Diagnostics.AddError(
			"Error Deleting Files MessageComment",
			"Could not delete message_comment id "+fmt.Sprint(state.Id.ValueInt64())+": "+err.Error(),
		)
	}
}

func (r *messageCommentResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
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

func (r *messageCommentResource) populateResourceModel(ctx context.Context, messageComment files_sdk.MessageComment, state *messageCommentResourceModel) (diags diag.Diagnostics) {
	var propDiags diag.Diagnostics

	state.Id = types.Int64Value(messageComment.Id)
	state.Body = types.StringValue(messageComment.Body)
	state.Reactions, propDiags = lib.ToDynamic(ctx, path.Root("reactions"), messageComment.Reactions, state.Reactions.UnderlyingValue())
	diags.Append(propDiags...)

	return
}
