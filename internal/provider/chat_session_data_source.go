package provider

import (
	"context"
	"fmt"

	files_sdk "github.com/Files-com/files-sdk-go/v3"
	chat_session "github.com/Files-com/files-sdk-go/v3/chatsession"
	"github.com/Files-com/terraform-provider-files/lib"
	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/datasource/schema"
	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

var (
	_ datasource.DataSource              = &chatSessionDataSource{}
	_ datasource.DataSourceWithConfigure = &chatSessionDataSource{}
)

func NewChatSessionDataSource() datasource.DataSource {
	return &chatSessionDataSource{}
}

type chatSessionDataSource struct {
	client *chat_session.Client
}

type chatSessionDataSourceModel struct {
	Id           types.String  `tfsdk:"id"`
	UserId       types.Int64   `tfsdk:"user_id"`
	WorkspaceId  types.Int64   `tfsdk:"workspace_id"`
	LastActiveAt types.String  `tfsdk:"last_active_at"`
	CreatedAt    types.String  `tfsdk:"created_at"`
	Messages     types.Dynamic `tfsdk:"messages"`
}

func (r *chatSessionDataSource) Configure(_ context.Context, req datasource.ConfigureRequest, resp *datasource.ConfigureResponse) {
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

	r.client = &chat_session.Client{Config: sdk_config}
}

func (r *chatSessionDataSource) Metadata(_ context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_chat_session"
}

func (r *chatSessionDataSource) Schema(_ context.Context, _ datasource.SchemaRequest, resp *datasource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Description: "A ChatSession represents one conversation with the Files.com AI Assistant.",
		Attributes: map[string]schema.Attribute{
			"id": schema.StringAttribute{
				Description: "Chat Session ID.",
				Required:    true,
			},
			"user_id": schema.Int64Attribute{
				Description: "User ID.",
				Computed:    true,
			},
			"workspace_id": schema.Int64Attribute{
				Description: "Workspace ID. `0` means the default workspace.",
				Computed:    true,
			},
			"last_active_at": schema.StringAttribute{
				Description: "Most recent chat activity date/time.",
				Computed:    true,
			},
			"created_at": schema.StringAttribute{
				Description: "Chat session creation date/time.",
				Computed:    true,
			},
			"messages": schema.DynamicAttribute{
				Description: "Visible conversation messages in this chat session. For performance reasons, this is not provided when listing chat sessions.",
				Computed:    true,
			},
		},
	}
}

func (r *chatSessionDataSource) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
	var data chatSessionDataSourceModel
	diags := req.Config.Get(ctx, &data)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	paramsChatSessionFind := files_sdk.ChatSessionFindParams{}
	paramsChatSessionFind.Id = data.Id.ValueString()

	chatSession, err := r.client.Find(paramsChatSessionFind, files_sdk.WithContext(ctx))
	if err != nil {
		resp.Diagnostics.AddError(
			"Error Reading Files ChatSession",
			"Could not read chat_session id "+fmt.Sprint(data.Id.ValueString())+": "+err.Error(),
		)
		return
	}

	diags = r.populateDataSourceModel(ctx, chatSession, &data)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	diags = resp.State.Set(ctx, data)
	resp.Diagnostics.Append(diags...)
}

func (r *chatSessionDataSource) populateDataSourceModel(ctx context.Context, chatSession files_sdk.ChatSession, state *chatSessionDataSourceModel) (diags diag.Diagnostics) {
	var propDiags diag.Diagnostics

	state.Id = types.StringValue(chatSession.Id)
	state.UserId = types.Int64Value(chatSession.UserId)
	state.WorkspaceId = types.Int64Value(chatSession.WorkspaceId)
	if err := lib.TimeToStringType(ctx, path.Root("last_active_at"), chatSession.LastActiveAt, &state.LastActiveAt); err != nil {
		diags.AddError(
			"Error Creating Files ChatSession",
			"Could not convert state last_active_at to string: "+err.Error(),
		)
	}
	if err := lib.TimeToStringType(ctx, path.Root("created_at"), chatSession.CreatedAt, &state.CreatedAt); err != nil {
		diags.AddError(
			"Error Creating Files ChatSession",
			"Could not convert state created_at to string: "+err.Error(),
		)
	}
	state.Messages, propDiags = lib.ToDynamic(ctx, path.Root("messages"), chatSession.Messages, state.Messages.UnderlyingValue())
	diags.Append(propDiags...)

	return
}
