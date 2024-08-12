package provider

import (
	"context"
	"fmt"

	files_sdk "github.com/Files-com/files-sdk-go/v3"
	message_comment_reaction "github.com/Files-com/files-sdk-go/v3/messagecommentreaction"
	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/datasource/schema"
	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

var (
	_ datasource.DataSource              = &messageCommentReactionDataSource{}
	_ datasource.DataSourceWithConfigure = &messageCommentReactionDataSource{}
)

func NewMessageCommentReactionDataSource() datasource.DataSource {
	return &messageCommentReactionDataSource{}
}

type messageCommentReactionDataSource struct {
	client *message_comment_reaction.Client
}

type messageCommentReactionDataSourceModel struct {
	Id    types.Int64  `tfsdk:"id"`
	Emoji types.String `tfsdk:"emoji"`
}

func (r *messageCommentReactionDataSource) Configure(_ context.Context, req datasource.ConfigureRequest, resp *datasource.ConfigureResponse) {
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

func (r *messageCommentReactionDataSource) Metadata(_ context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_message_comment_reaction"
}

func (r *messageCommentReactionDataSource) Schema(_ context.Context, _ datasource.SchemaRequest, resp *datasource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Description: "A MessageCommentReaction is a reaction emoji made by a user on a message comment.",
		Attributes: map[string]schema.Attribute{
			"id": schema.Int64Attribute{
				Description: "Reaction ID",
				Required:    true,
			},
			"emoji": schema.StringAttribute{
				Description: "Emoji used in the reaction.",
				Computed:    true,
			},
		},
	}
}

func (r *messageCommentReactionDataSource) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
	var data messageCommentReactionDataSourceModel
	diags := req.Config.Get(ctx, &data)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	paramsMessageCommentReactionFind := files_sdk.MessageCommentReactionFindParams{}
	paramsMessageCommentReactionFind.Id = data.Id.ValueInt64()

	messageCommentReaction, err := r.client.Find(paramsMessageCommentReactionFind, files_sdk.WithContext(ctx))
	if err != nil {
		resp.Diagnostics.AddError(
			"Error Reading Files MessageCommentReaction",
			"Could not read message_comment_reaction id "+fmt.Sprint(data.Id.ValueInt64())+": "+err.Error(),
		)
		return
	}

	diags = r.populateDataSourceModel(ctx, messageCommentReaction, &data)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	diags = resp.State.Set(ctx, data)
	resp.Diagnostics.Append(diags...)
}

func (r *messageCommentReactionDataSource) populateDataSourceModel(ctx context.Context, messageCommentReaction files_sdk.MessageCommentReaction, state *messageCommentReactionDataSourceModel) (diags diag.Diagnostics) {
	state.Id = types.Int64Value(messageCommentReaction.Id)
	state.Emoji = types.StringValue(messageCommentReaction.Emoji)

	return
}
