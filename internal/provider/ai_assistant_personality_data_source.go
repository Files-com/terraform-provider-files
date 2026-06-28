package provider

import (
	"context"
	"fmt"

	files_sdk "github.com/Files-com/files-sdk-go/v3"
	ai_assistant_personality "github.com/Files-com/files-sdk-go/v3/aiassistantpersonality"
	"github.com/Files-com/terraform-provider-files/lib"
	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/datasource/schema"
	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

var (
	_ datasource.DataSource              = &aiAssistantPersonalityDataSource{}
	_ datasource.DataSourceWithConfigure = &aiAssistantPersonalityDataSource{}
)

func NewAiAssistantPersonalityDataSource() datasource.DataSource {
	return &aiAssistantPersonalityDataSource{}
}

type aiAssistantPersonalityDataSource struct {
	client *ai_assistant_personality.Client
}

type aiAssistantPersonalityDataSourceModel struct {
	Id                   types.Int64  `tfsdk:"id"`
	WorkspaceId          types.Int64  `tfsdk:"workspace_id"`
	SystemPrompt         types.String `tfsdk:"system_prompt"`
	UseByDefault         types.Bool   `tfsdk:"use_by_default"`
	ApplyToAllWorkspaces types.Bool   `tfsdk:"apply_to_all_workspaces"`
	CreatedAt            types.String `tfsdk:"created_at"`
	UpdatedAt            types.String `tfsdk:"updated_at"`
}

func (r *aiAssistantPersonalityDataSource) Configure(_ context.Context, req datasource.ConfigureRequest, resp *datasource.ConfigureResponse) {
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

	r.client = &ai_assistant_personality.Client{Config: sdk_config}
}

func (r *aiAssistantPersonalityDataSource) Metadata(_ context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_ai_assistant_personality"
}

func (r *aiAssistantPersonalityDataSource) Schema(_ context.Context, _ datasource.SchemaRequest, resp *datasource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Description: "An AI Assistant Personality defines a system prompt used to customize the in-app AI Assistant.",
		Attributes: map[string]schema.Attribute{
			"id": schema.Int64Attribute{
				Description: "AI Assistant Personality ID.",
				Required:    true,
			},
			"workspace_id": schema.Int64Attribute{
				Description: "Workspace ID. `0` means the default workspace.",
				Computed:    true,
			},
			"system_prompt": schema.StringAttribute{
				Description: "System prompt injected into the in-app AI Assistant.",
				Computed:    true,
			},
			"use_by_default": schema.BoolAttribute{
				Description: "Whether this personality is the default personality for the Workspace.",
				Computed:    true,
			},
			"apply_to_all_workspaces": schema.BoolAttribute{
				Description: "If true, this default-workspace personality can apply to users in all workspaces.",
				Computed:    true,
			},
			"created_at": schema.StringAttribute{
				Description: "Creation time.",
				Computed:    true,
			},
			"updated_at": schema.StringAttribute{
				Description: "Last update time.",
				Computed:    true,
			},
		},
	}
}

func (r *aiAssistantPersonalityDataSource) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
	var data aiAssistantPersonalityDataSourceModel
	diags := req.Config.Get(ctx, &data)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	paramsAiAssistantPersonalityFind := files_sdk.AiAssistantPersonalityFindParams{}
	paramsAiAssistantPersonalityFind.Id = data.Id.ValueInt64()

	aiAssistantPersonality, err := r.client.Find(paramsAiAssistantPersonalityFind, files_sdk.WithContext(ctx))
	if err != nil {
		resp.Diagnostics.AddError(
			"Error Reading Files AiAssistantPersonality",
			"Could not read ai_assistant_personality id "+fmt.Sprint(data.Id.ValueInt64())+": "+err.Error(),
		)
		return
	}

	diags = r.populateDataSourceModel(ctx, aiAssistantPersonality, &data)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	diags = resp.State.Set(ctx, data)
	resp.Diagnostics.Append(diags...)
}

func (r *aiAssistantPersonalityDataSource) populateDataSourceModel(ctx context.Context, aiAssistantPersonality files_sdk.AiAssistantPersonality, state *aiAssistantPersonalityDataSourceModel) (diags diag.Diagnostics) {
	state.Id = types.Int64Value(aiAssistantPersonality.Id)
	state.WorkspaceId = types.Int64Value(aiAssistantPersonality.WorkspaceId)
	state.SystemPrompt = types.StringValue(aiAssistantPersonality.SystemPrompt)
	state.UseByDefault = types.BoolPointerValue(aiAssistantPersonality.UseByDefault)
	state.ApplyToAllWorkspaces = types.BoolPointerValue(aiAssistantPersonality.ApplyToAllWorkspaces)
	if err := lib.TimeToStringType(ctx, path.Root("created_at"), aiAssistantPersonality.CreatedAt, &state.CreatedAt); err != nil {
		diags.AddError(
			"Error Creating Files AiAssistantPersonality",
			"Could not convert state created_at to string: "+err.Error(),
		)
	}
	if err := lib.TimeToStringType(ctx, path.Root("updated_at"), aiAssistantPersonality.UpdatedAt, &state.UpdatedAt); err != nil {
		diags.AddError(
			"Error Creating Files AiAssistantPersonality",
			"Could not convert state updated_at to string: "+err.Error(),
		)
	}

	return
}
