package provider

import (
	"context"
	"fmt"

	files_sdk "github.com/Files-com/files-sdk-go/v3"
	message_comment "github.com/Files-com/files-sdk-go/v3/messagecomment"
	"github.com/Files-com/terraform-provider-files/lib"
	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/datasource/schema"
	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

var (
	_ datasource.DataSource              = &messageCommentDataSource{}
	_ datasource.DataSourceWithConfigure = &messageCommentDataSource{}
)

func NewMessageCommentDataSource() datasource.DataSource {
	return &messageCommentDataSource{}
}

type messageCommentDataSource struct {
	client *message_comment.Client
}

type messageCommentDataSourceModel struct {
	Id        types.Int64   `tfsdk:"id"`
	Body      types.String  `tfsdk:"body"`
	Reactions types.Dynamic `tfsdk:"reactions"`
}

func (r *messageCommentDataSource) Configure(_ context.Context, req datasource.ConfigureRequest, resp *datasource.ConfigureResponse) {
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

func (r *messageCommentDataSource) Metadata(_ context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_message_comment"
}

func (r *messageCommentDataSource) Schema(_ context.Context, _ datasource.SchemaRequest, resp *datasource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Description: "A message comment represents a comment made by a user on a message.",
		Attributes: map[string]schema.Attribute{
			"id": schema.Int64Attribute{
				Description: "Message Comment ID",
				Required:    true,
			},
			"body": schema.StringAttribute{
				Description: "Comment body.",
				Computed:    true,
			},
			"reactions": schema.DynamicAttribute{
				Description: "Reactions to this comment.",
				Computed:    true,
			},
		},
	}
}

func (r *messageCommentDataSource) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
	var data messageCommentDataSourceModel
	diags := req.Config.Get(ctx, &data)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	paramsMessageCommentFind := files_sdk.MessageCommentFindParams{}
	paramsMessageCommentFind.Id = data.Id.ValueInt64()

	messageComment, err := r.client.Find(paramsMessageCommentFind, files_sdk.WithContext(ctx))
	if err != nil {
		resp.Diagnostics.AddError(
			"Error Reading Files MessageComment",
			"Could not read message_comment id "+fmt.Sprint(data.Id.ValueInt64())+": "+err.Error(),
		)
		return
	}

	diags = r.populateDataSourceModel(ctx, messageComment, &data)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	diags = resp.State.Set(ctx, data)
	resp.Diagnostics.Append(diags...)
}

func (r *messageCommentDataSource) populateDataSourceModel(ctx context.Context, messageComment files_sdk.MessageComment, state *messageCommentDataSourceModel) (diags diag.Diagnostics) {
	var propDiags diag.Diagnostics

	state.Id = types.Int64Value(messageComment.Id)
	state.Body = types.StringValue(messageComment.Body)
	state.Reactions, propDiags = lib.ToDynamic(ctx, path.Root("reactions"), messageComment.Reactions, state.Reactions.UnderlyingValue())
	diags.Append(propDiags...)

	return
}
