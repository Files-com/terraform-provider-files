package provider

import (
	"context"
	"fmt"
	"strconv"
	"strings"

	files_sdk "github.com/Files-com/files-sdk-go/v3"
	message_comment_reaction "github.com/Files-com/files-sdk-go/v3/messagecommentreaction"
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
	_ resource.Resource                = &messageCommentReactionResource{}
	_ resource.ResourceWithConfigure   = &messageCommentReactionResource{}
	_ resource.ResourceWithImportState = &messageCommentReactionResource{}
)

func NewMessageCommentReactionResource() resource.Resource {
	return &messageCommentReactionResource{}
}

type messageCommentReactionResource struct {
	client *message_comment_reaction.Client
}

type messageCommentReactionResourceModel struct {
	Emoji  types.String `tfsdk:"emoji"`
	UserId types.Int64  `tfsdk:"user_id"`
	Id     types.Int64  `tfsdk:"id"`
}

func (r *messageCommentReactionResource) Configure(_ context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {
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

	r.client = &message_comment_reaction.Client{Config: sdk_config}
}

func (r *messageCommentReactionResource) Metadata(_ context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_message_comment_reaction"
}

func (r *messageCommentReactionResource) Schema(_ context.Context, _ resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Description: "A MessageCommentReaction is a reaction emoji made by a user on a message comment.",
		Attributes: map[string]schema.Attribute{
			"emoji": schema.StringAttribute{
				Description: "Emoji used in the reaction.",
				Required:    true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.RequiresReplace(),
				},
			},
			"user_id": schema.Int64Attribute{
				Description: "User ID.  Provide a value of `0` to operate the current session's user.",
				Optional:    true,
				PlanModifiers: []planmodifier.Int64{
					int64planmodifier.RequiresReplace(),
				},
			},
			"id": schema.Int64Attribute{
				Description: "Reaction ID",
				Computed:    true,
				PlanModifiers: []planmodifier.Int64{
					int64planmodifier.UseStateForUnknown(),
				},
			},
		},
	}
}

func (r *messageCommentReactionResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	var plan messageCommentReactionResourceModel
	diags := req.Plan.Get(ctx, &plan)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	paramsMessageCommentReactionCreate := files_sdk.MessageCommentReactionCreateParams{}
	paramsMessageCommentReactionCreate.UserId = plan.UserId.ValueInt64()
	paramsMessageCommentReactionCreate.Emoji = plan.Emoji.ValueString()

	if resp.Diagnostics.HasError() {
		return
	}

	messageCommentReaction, err := r.client.Create(paramsMessageCommentReactionCreate, files_sdk.WithContext(ctx))
	if err != nil {
		resp.Diagnostics.AddError(
			"Error Creating Files MessageCommentReaction",
			"Could not create message_comment_reaction, unexpected error: "+err.Error(),
		)
		return
	}

	diags = r.populateResourceModel(ctx, messageCommentReaction, &plan)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	diags = resp.State.Set(ctx, plan)
	resp.Diagnostics.Append(diags...)
}

func (r *messageCommentReactionResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	var state messageCommentReactionResourceModel
	diags := req.State.Get(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	paramsMessageCommentReactionFind := files_sdk.MessageCommentReactionFindParams{}
	paramsMessageCommentReactionFind.Id = state.Id.ValueInt64()

	messageCommentReaction, err := r.client.Find(paramsMessageCommentReactionFind, files_sdk.WithContext(ctx))
	if err != nil {
		if files_sdk.IsNotExist(err) {
			resp.State.RemoveResource(ctx)
			return
		}

		resp.Diagnostics.AddError(
			"Error Reading Files MessageCommentReaction",
			"Could not read message_comment_reaction id "+fmt.Sprint(state.Id.ValueInt64())+": "+err.Error(),
		)
		return
	}

	diags = r.populateResourceModel(ctx, messageCommentReaction, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	diags = resp.State.Set(ctx, state)
	resp.Diagnostics.Append(diags...)
}

func (r *messageCommentReactionResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	resp.Diagnostics.AddError(
		"Resource Update Not Implemented",
		"This resource does not support updates.",
	)
}

func (r *messageCommentReactionResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	var state messageCommentReactionResourceModel
	diags := req.State.Get(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	paramsMessageCommentReactionDelete := files_sdk.MessageCommentReactionDeleteParams{}
	paramsMessageCommentReactionDelete.Id = state.Id.ValueInt64()

	err := r.client.Delete(paramsMessageCommentReactionDelete, files_sdk.WithContext(ctx))
	if err != nil && !files_sdk.IsNotExist(err) {
		resp.Diagnostics.AddError(
			"Error Deleting Files MessageCommentReaction",
			"Could not delete message_comment_reaction id "+fmt.Sprint(state.Id.ValueInt64())+": "+err.Error(),
		)
	}
}

func (r *messageCommentReactionResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
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

func (r *messageCommentReactionResource) populateResourceModel(ctx context.Context, messageCommentReaction files_sdk.MessageCommentReaction, state *messageCommentReactionResourceModel) (diags diag.Diagnostics) {
	state.Id = types.Int64Value(messageCommentReaction.Id)
	state.Emoji = types.StringValue(messageCommentReaction.Emoji)

	return
}
