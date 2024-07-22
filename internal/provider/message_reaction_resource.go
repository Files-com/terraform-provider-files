package provider

import (
	"context"
	"fmt"
	"strconv"
	"strings"

	files_sdk "github.com/Files-com/files-sdk-go/v3"
	message_reaction "github.com/Files-com/files-sdk-go/v3/messagereaction"
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
	_ resource.Resource                = &messageReactionResource{}
	_ resource.ResourceWithConfigure   = &messageReactionResource{}
	_ resource.ResourceWithImportState = &messageReactionResource{}
)

func NewMessageReactionResource() resource.Resource {
	return &messageReactionResource{}
}

type messageReactionResource struct {
	client *message_reaction.Client
}

type messageReactionResourceModel struct {
	Emoji  types.String `tfsdk:"emoji"`
	UserId types.Int64  `tfsdk:"user_id"`
	Id     types.Int64  `tfsdk:"id"`
}

func (r *messageReactionResource) Configure(_ context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {
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

	r.client = &message_reaction.Client{Config: sdk_config}
}

func (r *messageReactionResource) Metadata(_ context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_message_reaction"
}

func (r *messageReactionResource) Schema(_ context.Context, _ resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Description: "A message reaction represents a reaction emoji made by a user on a message.",
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

func (r *messageReactionResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	var plan messageReactionResourceModel
	diags := req.Plan.Get(ctx, &plan)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	paramsMessageReactionCreate := files_sdk.MessageReactionCreateParams{}
	paramsMessageReactionCreate.UserId = plan.UserId.ValueInt64()
	paramsMessageReactionCreate.Emoji = plan.Emoji.ValueString()

	if resp.Diagnostics.HasError() {
		return
	}

	messageReaction, err := r.client.Create(paramsMessageReactionCreate, files_sdk.WithContext(ctx))
	if err != nil {
		resp.Diagnostics.AddError(
			"Error Creating Files MessageReaction",
			"Could not create message_reaction, unexpected error: "+err.Error(),
		)
		return
	}

	diags = r.populateResourceModel(ctx, messageReaction, &plan)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	diags = resp.State.Set(ctx, plan)
	resp.Diagnostics.Append(diags...)
}

func (r *messageReactionResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	var state messageReactionResourceModel
	diags := req.State.Get(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	paramsMessageReactionFind := files_sdk.MessageReactionFindParams{}
	paramsMessageReactionFind.Id = state.Id.ValueInt64()

	messageReaction, err := r.client.Find(paramsMessageReactionFind, files_sdk.WithContext(ctx))
	if err != nil {
		resp.Diagnostics.AddError(
			"Error Reading Files MessageReaction",
			"Could not read message_reaction id "+fmt.Sprint(state.Id.ValueInt64())+": "+err.Error(),
		)
		return
	}

	diags = r.populateResourceModel(ctx, messageReaction, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	diags = resp.State.Set(ctx, state)
	resp.Diagnostics.Append(diags...)
}

func (r *messageReactionResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	resp.Diagnostics.AddError(
		"Error Updating Files MessageReaction",
		"Update operation not implemented",
	)
}

func (r *messageReactionResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	var state messageReactionResourceModel
	diags := req.State.Get(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	paramsMessageReactionDelete := files_sdk.MessageReactionDeleteParams{}
	paramsMessageReactionDelete.Id = state.Id.ValueInt64()

	err := r.client.Delete(paramsMessageReactionDelete, files_sdk.WithContext(ctx))
	if err != nil {
		resp.Diagnostics.AddError(
			"Error Deleting Files MessageReaction",
			"Could not delete message_reaction id "+fmt.Sprint(state.Id.ValueInt64())+": "+err.Error(),
		)
	}
}

func (r *messageReactionResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
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

func (r *messageReactionResource) populateResourceModel(ctx context.Context, messageReaction files_sdk.MessageReaction, state *messageReactionResourceModel) (diags diag.Diagnostics) {
	state.Id = types.Int64Value(messageReaction.Id)
	state.Emoji = types.StringValue(messageReaction.Emoji)

	return
}
