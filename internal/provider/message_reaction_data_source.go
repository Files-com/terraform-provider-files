package provider

import (
	"context"
	"fmt"

	files_sdk "github.com/Files-com/files-sdk-go/v3"
	message_reaction "github.com/Files-com/files-sdk-go/v3/messagereaction"
	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/datasource/schema"
	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

var (
	_ datasource.DataSource              = &messageReactionDataSource{}
	_ datasource.DataSourceWithConfigure = &messageReactionDataSource{}
)

func NewMessageReactionDataSource() datasource.DataSource {
	return &messageReactionDataSource{}
}

type messageReactionDataSource struct {
	client *message_reaction.Client
}

type messageReactionDataSourceModel struct {
	Id    types.Int64  `tfsdk:"id"`
	Emoji types.String `tfsdk:"emoji"`
}

func (r *messageReactionDataSource) Configure(_ context.Context, req datasource.ConfigureRequest, resp *datasource.ConfigureResponse) {
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

func (r *messageReactionDataSource) Metadata(_ context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_message_reaction"
}

func (r *messageReactionDataSource) Schema(_ context.Context, _ datasource.SchemaRequest, resp *datasource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Description: "A message reaction represents a reaction emoji made by a user on a message.",
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

func (r *messageReactionDataSource) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
	var data messageReactionDataSourceModel
	diags := req.Config.Get(ctx, &data)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	paramsMessageReactionFind := files_sdk.MessageReactionFindParams{}
	paramsMessageReactionFind.Id = data.Id.ValueInt64()

	messageReaction, err := r.client.Find(paramsMessageReactionFind, files_sdk.WithContext(ctx))
	if err != nil {
		resp.Diagnostics.AddError(
			"Error Reading Files MessageReaction",
			"Could not read message_reaction id "+fmt.Sprint(data.Id.ValueInt64())+": "+err.Error(),
		)
		return
	}

	diags = r.populateDataSourceModel(ctx, messageReaction, &data)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	diags = resp.State.Set(ctx, data)
	resp.Diagnostics.Append(diags...)
}

func (r *messageReactionDataSource) populateDataSourceModel(ctx context.Context, messageReaction files_sdk.MessageReaction, state *messageReactionDataSourceModel) (diags diag.Diagnostics) {
	state.Id = types.Int64Value(messageReaction.Id)
	state.Emoji = types.StringValue(messageReaction.Emoji)

	return
}
