package provider

import (
	"context"
	"fmt"
	"strconv"
	"strings"

	files_sdk "github.com/Files-com/files-sdk-go/v3"
	ai_assistant_personality "github.com/Files-com/files-sdk-go/v3/aiassistantpersonality"
	"github.com/Files-com/terraform-provider-files/lib"
	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/boolplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/int64planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

var (
	_ resource.Resource                = &aiAssistantPersonalityResource{}
	_ resource.ResourceWithConfigure   = &aiAssistantPersonalityResource{}
	_ resource.ResourceWithImportState = &aiAssistantPersonalityResource{}
)

func NewAiAssistantPersonalityResource() resource.Resource {
	return &aiAssistantPersonalityResource{}
}

type aiAssistantPersonalityResource struct {
	client *ai_assistant_personality.Client
}

type aiAssistantPersonalityResourceModel struct {
	Name                 types.String `tfsdk:"name"`
	SystemPrompt         types.String `tfsdk:"system_prompt"`
	WorkspaceId          types.Int64  `tfsdk:"workspace_id"`
	UseByDefault         types.Bool   `tfsdk:"use_by_default"`
	ApplyToAllWorkspaces types.Bool   `tfsdk:"apply_to_all_workspaces"`
	Id                   types.Int64  `tfsdk:"id"`
	CreatedAt            types.String `tfsdk:"created_at"`
	UpdatedAt            types.String `tfsdk:"updated_at"`
}

func (r *aiAssistantPersonalityResource) Configure(_ context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {
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

func (r *aiAssistantPersonalityResource) Metadata(_ context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_ai_assistant_personality"
}

func (r *aiAssistantPersonalityResource) Schema(_ context.Context, _ resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Description: "An AI Assistant Personality defines a system prompt used to customize the in-app AI Assistant.",
		Attributes: map[string]schema.Attribute{
			"name": schema.StringAttribute{
				Description: "AI Assistant Personality name.",
				Required:    true,
			},
			"system_prompt": schema.StringAttribute{
				Description: "System prompt injected into the in-app AI Assistant.",
				Required:    true,
			},
			"workspace_id": schema.Int64Attribute{
				Description: "Workspace ID. `0` means the default workspace.",
				Computed:    true,
				Optional:    true,
				PlanModifiers: []planmodifier.Int64{
					int64planmodifier.UseStateForUnknown(),
				},
			},
			"use_by_default": schema.BoolAttribute{
				Description: "Whether this personality is the default personality for the Workspace.",
				Computed:    true,
				Optional:    true,
				PlanModifiers: []planmodifier.Bool{
					boolplanmodifier.UseStateForUnknown(),
				},
			},
			"apply_to_all_workspaces": schema.BoolAttribute{
				Description: "If true, this default-workspace personality can apply to users in all workspaces.",
				Computed:    true,
				Optional:    true,
				PlanModifiers: []planmodifier.Bool{
					boolplanmodifier.UseStateForUnknown(),
				},
			},
			"id": schema.Int64Attribute{
				Description: "AI Assistant Personality ID.",
				Computed:    true,
				PlanModifiers: []planmodifier.Int64{
					int64planmodifier.UseStateForUnknown(),
				},
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

func (r *aiAssistantPersonalityResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	var plan aiAssistantPersonalityResourceModel
	diags := req.Plan.Get(ctx, &plan)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}
	var config aiAssistantPersonalityResourceModel
	diags = req.Config.Get(ctx, &config)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	paramsAiAssistantPersonalityCreate := files_sdk.AiAssistantPersonalityCreateParams{}
	if !plan.ApplyToAllWorkspaces.IsNull() && !plan.ApplyToAllWorkspaces.IsUnknown() {
		paramsAiAssistantPersonalityCreate.ApplyToAllWorkspaces = plan.ApplyToAllWorkspaces.ValueBoolPointer()
	}
	paramsAiAssistantPersonalityCreate.Name = plan.Name.ValueString()
	paramsAiAssistantPersonalityCreate.SystemPrompt = plan.SystemPrompt.ValueString()
	if !plan.UseByDefault.IsNull() && !plan.UseByDefault.IsUnknown() {
		paramsAiAssistantPersonalityCreate.UseByDefault = plan.UseByDefault.ValueBoolPointer()
	}
	paramsAiAssistantPersonalityCreate.WorkspaceId = plan.WorkspaceId.ValueInt64()

	if resp.Diagnostics.HasError() {
		return
	}

	aiAssistantPersonality, err := r.client.Create(paramsAiAssistantPersonalityCreate, files_sdk.WithContext(ctx))
	if err != nil {
		resp.Diagnostics.AddError(
			"Error Creating Files AiAssistantPersonality",
			"Could not create ai_assistant_personality, unexpected error: "+err.Error(),
		)
		return
	}

	diags = r.populateResourceModel(ctx, aiAssistantPersonality, &plan)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	diags = resp.State.Set(ctx, plan)
	resp.Diagnostics.Append(diags...)
}

func (r *aiAssistantPersonalityResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	var state aiAssistantPersonalityResourceModel
	diags := req.State.Get(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	paramsAiAssistantPersonalityFind := files_sdk.AiAssistantPersonalityFindParams{}
	paramsAiAssistantPersonalityFind.Id = state.Id.ValueInt64()

	aiAssistantPersonality, err := r.client.Find(paramsAiAssistantPersonalityFind, files_sdk.WithContext(ctx))
	if err != nil {
		if files_sdk.IsNotExist(err) {
			resp.State.RemoveResource(ctx)
			return
		}

		resp.Diagnostics.AddError(
			"Error Reading Files AiAssistantPersonality",
			"Could not read ai_assistant_personality id "+fmt.Sprint(state.Id.ValueInt64())+": "+err.Error(),
		)
		return
	}

	diags = r.populateResourceModel(ctx, aiAssistantPersonality, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	diags = resp.State.Set(ctx, state)
	resp.Diagnostics.Append(diags...)
}

func (r *aiAssistantPersonalityResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	var plan aiAssistantPersonalityResourceModel
	diags := req.Plan.Get(ctx, &plan)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}
	var config aiAssistantPersonalityResourceModel
	diags = req.Config.Get(ctx, &config)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	paramsAiAssistantPersonalityUpdate := map[string]interface{}{}
	if !plan.Id.IsNull() && !plan.Id.IsUnknown() {
		paramsAiAssistantPersonalityUpdate["id"] = plan.Id.ValueInt64()
	}
	if !config.ApplyToAllWorkspaces.IsNull() && !config.ApplyToAllWorkspaces.IsUnknown() {
		paramsAiAssistantPersonalityUpdate["apply_to_all_workspaces"] = config.ApplyToAllWorkspaces.ValueBool()
	}
	if !config.Name.IsNull() && !config.Name.IsUnknown() {
		paramsAiAssistantPersonalityUpdate["name"] = config.Name.ValueString()
	}
	if !config.SystemPrompt.IsNull() && !config.SystemPrompt.IsUnknown() {
		paramsAiAssistantPersonalityUpdate["system_prompt"] = config.SystemPrompt.ValueString()
	}
	if !config.UseByDefault.IsNull() && !config.UseByDefault.IsUnknown() {
		paramsAiAssistantPersonalityUpdate["use_by_default"] = config.UseByDefault.ValueBool()
	}
	if !config.WorkspaceId.IsNull() && !config.WorkspaceId.IsUnknown() {
		paramsAiAssistantPersonalityUpdate["workspace_id"] = config.WorkspaceId.ValueInt64()
	}

	if resp.Diagnostics.HasError() {
		return
	}

	aiAssistantPersonality, err := r.client.UpdateWithMap(paramsAiAssistantPersonalityUpdate, files_sdk.WithContext(ctx))
	if err != nil {
		resp.Diagnostics.AddError(
			"Error Updating Files AiAssistantPersonality",
			"Could not update ai_assistant_personality, unexpected error: "+err.Error(),
		)
		return
	}

	diags = r.populateResourceModel(ctx, aiAssistantPersonality, &plan)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	diags = resp.State.Set(ctx, plan)
	resp.Diagnostics.Append(diags...)
}

func (r *aiAssistantPersonalityResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	var state aiAssistantPersonalityResourceModel
	diags := req.State.Get(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	paramsAiAssistantPersonalityDelete := files_sdk.AiAssistantPersonalityDeleteParams{}
	paramsAiAssistantPersonalityDelete.Id = state.Id.ValueInt64()

	err := r.client.Delete(paramsAiAssistantPersonalityDelete, files_sdk.WithContext(ctx))
	if err != nil && !files_sdk.IsNotExist(err) {
		resp.Diagnostics.AddError(
			"Error Deleting Files AiAssistantPersonality",
			"Could not delete ai_assistant_personality id "+fmt.Sprint(state.Id.ValueInt64())+": "+err.Error(),
		)
	}
}

func (r *aiAssistantPersonalityResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
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

func (r *aiAssistantPersonalityResource) populateResourceModel(ctx context.Context, aiAssistantPersonality files_sdk.AiAssistantPersonality, state *aiAssistantPersonalityResourceModel) (diags diag.Diagnostics) {
	state.Id = types.Int64Value(aiAssistantPersonality.Id)
	state.WorkspaceId = types.Int64Value(aiAssistantPersonality.WorkspaceId)
	state.Name = types.StringValue(aiAssistantPersonality.Name)
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
