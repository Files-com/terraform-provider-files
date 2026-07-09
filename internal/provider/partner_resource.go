package provider

import (
	"context"
	"fmt"
	"strconv"
	"strings"

	files_sdk "github.com/Files-com/files-sdk-go/v3"
	partner "github.com/Files-com/files-sdk-go/v3/partner"
	"github.com/hashicorp/terraform-plugin-framework-validators/stringvalidator"
	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/boolplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/int64planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/schema/validator"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

var (
	_ resource.Resource                = &partnerResource{}
	_ resource.ResourceWithConfigure   = &partnerResource{}
	_ resource.ResourceWithImportState = &partnerResource{}
)

func NewPartnerResource() resource.Resource {
	return &partnerResource{}
}

type partnerResource struct {
	client *partner.Client
}

type partnerResourceModel struct {
	Name                       types.String `tfsdk:"name"`
	RootFolder                 types.String `tfsdk:"root_folder"`
	AllowBypassing2faPolicies  types.Bool   `tfsdk:"allow_bypassing_2fa_policies"`
	AllowedIps                 types.String `tfsdk:"allowed_ips"`
	AllowCredentialChanges     types.Bool   `tfsdk:"allow_credential_changes"`
	AllowProvidingGpgKeys      types.Bool   `tfsdk:"allow_providing_gpg_keys"`
	AllowUserCreation          types.Bool   `tfsdk:"allow_user_creation"`
	CcEmailsToResponsibleParty types.Bool   `tfsdk:"cc_emails_to_responsible_party"`
	AiAssistantPersonalityId   types.Int64  `tfsdk:"ai_assistant_personality_id"`
	WorkspaceId                types.Int64  `tfsdk:"workspace_id"`
	Notes                      types.String `tfsdk:"notes"`
	PartnerChannelTemplateId   types.Int64  `tfsdk:"partner_channel_template_id"`
	ResponsibleGroupId         types.Int64  `tfsdk:"responsible_group_id"`
	ResponsibleUserId          types.Int64  `tfsdk:"responsible_user_id"`
	Tags                       types.String `tfsdk:"tags"`
	Id                         types.Int64  `tfsdk:"id"`
	PartnerAdminIds            types.List   `tfsdk:"partner_admin_ids"`
	PartnershipRole            types.String `tfsdk:"partnership_role"`
	UserIds                    types.List   `tfsdk:"user_ids"`
}

func (r *partnerResource) Configure(_ context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {
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

	r.client = &partner.Client{Config: sdk_config}
}

func (r *partnerResource) Metadata(_ context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_partner"
}

func (r *partnerResource) Schema(_ context.Context, _ resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Description: "A Partner is a first-class entity that cleanly represents an external organization, enables delegated administration, and enforces strict boundaries.",
		Attributes: map[string]schema.Attribute{
			"name": schema.StringAttribute{
				Description: "The name of the Partner.",
				Required:    true,
			},
			"root_folder": schema.StringAttribute{
				Description: "The root folder path for this Partner.",
				Required:    true,
			},
			"allow_bypassing_2fa_policies": schema.BoolAttribute{
				Description: "Allow Partner Admins to change Two-Factor Authentication requirements for Partner Users.",
				Computed:    true,
				Optional:    true,
				PlanModifiers: []planmodifier.Bool{
					boolplanmodifier.UseStateForUnknown(),
				},
			},
			"allowed_ips": schema.StringAttribute{
				Description: "A list of allowed IPs for this Partner. Newline delimited. Partner User IP access is allowed when the IP matches the Partner, User, or Site allowed IP lists.",
				Computed:    true,
				Optional:    true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
			},
			"allow_credential_changes": schema.BoolAttribute{
				Description: "Allow Partner Admins to change or reset credentials for users belonging to this Partner.",
				Computed:    true,
				Optional:    true,
				PlanModifiers: []planmodifier.Bool{
					boolplanmodifier.UseStateForUnknown(),
				},
			},
			"allow_providing_gpg_keys": schema.BoolAttribute{
				Description: "Allow Partner Admins to provide GPG keys.",
				Computed:    true,
				Optional:    true,
				PlanModifiers: []planmodifier.Bool{
					boolplanmodifier.UseStateForUnknown(),
				},
			},
			"allow_user_creation": schema.BoolAttribute{
				Description: "Allow Partner Admins to create users.",
				Computed:    true,
				Optional:    true,
				PlanModifiers: []planmodifier.Bool{
					boolplanmodifier.UseStateForUnknown(),
				},
			},
			"cc_emails_to_responsible_party": schema.BoolAttribute{
				Description: "When `true`, emails sent to Partner users are copied to the responsible User or Group.",
				Computed:    true,
				Optional:    true,
				PlanModifiers: []planmodifier.Bool{
					boolplanmodifier.UseStateForUnknown(),
				},
			},
			"ai_assistant_personality_id": schema.Int64Attribute{
				Description: "AI Assistant Personality ID assigned to this Partner, if any. Users in the Partner inherit it unless a direct per-user assignment overrides it.",
				Computed:    true,
				Optional:    true,
				PlanModifiers: []planmodifier.Int64{
					int64planmodifier.UseStateForUnknown(),
				},
			},
			"workspace_id": schema.Int64Attribute{
				Description: "ID of the Workspace associated with this Partner.",
				Computed:    true,
				Optional:    true,
				PlanModifiers: []planmodifier.Int64{
					int64planmodifier.UseStateForUnknown(),
					int64planmodifier.RequiresReplace(),
				},
			},
			"notes": schema.StringAttribute{
				Description: "Notes about this Partner.",
				Computed:    true,
				Optional:    true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
			},
			"partner_channel_template_id": schema.Int64Attribute{
				Description: "ID of the Partner Channel Template assigned to this Partner.",
				Computed:    true,
				Optional:    true,
				PlanModifiers: []planmodifier.Int64{
					int64planmodifier.UseStateForUnknown(),
				},
			},
			"responsible_group_id": schema.Int64Attribute{
				Description: "ID of the Group responsible for this Partner.",
				Computed:    true,
				Optional:    true,
				PlanModifiers: []planmodifier.Int64{
					int64planmodifier.UseStateForUnknown(),
				},
			},
			"responsible_user_id": schema.Int64Attribute{
				Description: "ID of the User responsible for this Partner.",
				Computed:    true,
				Optional:    true,
				PlanModifiers: []planmodifier.Int64{
					int64planmodifier.UseStateForUnknown(),
				},
			},
			"tags": schema.StringAttribute{
				Description: "Comma-separated list of Tags for this Partner. Tags are used for other features, such as UserLifecycleRules, which can target specific tags.  Tags must only contain lowercase letters, numbers, and hyphens.",
				Computed:    true,
				Optional:    true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
			},
			"id": schema.Int64Attribute{
				Description: "The unique ID of the Partner.",
				Computed:    true,
				PlanModifiers: []planmodifier.Int64{
					int64planmodifier.UseStateForUnknown(),
				},
			},
			"partner_admin_ids": schema.ListAttribute{
				Description: "Array of User IDs that are Partner Admins for this Partner.",
				Computed:    true,
				ElementType: types.Int64Type,
			},
			"partnership_role": schema.StringAttribute{
				Description: "This site's role in Partner Site relationships for this Partner. Can be `host`, `guest`, `host_and_guest`, or null.",
				Computed:    true,
				Validators: []validator.String{
					stringvalidator.OneOf("host", "guest", "host_and_guest"),
				},
			},
			"user_ids": schema.ListAttribute{
				Description: "Array of User IDs that belong to this Partner.",
				Computed:    true,
				ElementType: types.Int64Type,
			},
		},
	}
}

func (r *partnerResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	var plan partnerResourceModel
	diags := req.Plan.Get(ctx, &plan)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}
	var config partnerResourceModel
	diags = req.Config.Get(ctx, &config)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	paramsPartnerCreate := files_sdk.PartnerCreateParams{}
	paramsPartnerCreate.AiAssistantPersonalityId = plan.AiAssistantPersonalityId.ValueInt64()
	paramsPartnerCreate.AllowedIps = plan.AllowedIps.ValueString()
	if !plan.AllowBypassing2faPolicies.IsNull() && !plan.AllowBypassing2faPolicies.IsUnknown() {
		paramsPartnerCreate.AllowBypassing2faPolicies = plan.AllowBypassing2faPolicies.ValueBoolPointer()
	}
	if !plan.AllowCredentialChanges.IsNull() && !plan.AllowCredentialChanges.IsUnknown() {
		paramsPartnerCreate.AllowCredentialChanges = plan.AllowCredentialChanges.ValueBoolPointer()
	}
	if !plan.AllowProvidingGpgKeys.IsNull() && !plan.AllowProvidingGpgKeys.IsUnknown() {
		paramsPartnerCreate.AllowProvidingGpgKeys = plan.AllowProvidingGpgKeys.ValueBoolPointer()
	}
	if !plan.AllowUserCreation.IsNull() && !plan.AllowUserCreation.IsUnknown() {
		paramsPartnerCreate.AllowUserCreation = plan.AllowUserCreation.ValueBoolPointer()
	}
	if !plan.CcEmailsToResponsibleParty.IsNull() && !plan.CcEmailsToResponsibleParty.IsUnknown() {
		paramsPartnerCreate.CcEmailsToResponsibleParty = plan.CcEmailsToResponsibleParty.ValueBoolPointer()
	}
	paramsPartnerCreate.Notes = plan.Notes.ValueString()
	paramsPartnerCreate.PartnerChannelTemplateId = plan.PartnerChannelTemplateId.ValueInt64()
	paramsPartnerCreate.ResponsibleGroupId = plan.ResponsibleGroupId.ValueInt64()
	paramsPartnerCreate.ResponsibleUserId = plan.ResponsibleUserId.ValueInt64()
	paramsPartnerCreate.Tags = plan.Tags.ValueString()
	paramsPartnerCreate.Name = plan.Name.ValueString()
	paramsPartnerCreate.RootFolder = plan.RootFolder.ValueString()
	paramsPartnerCreate.WorkspaceId = plan.WorkspaceId.ValueInt64()

	if resp.Diagnostics.HasError() {
		return
	}

	partner, err := r.client.Create(paramsPartnerCreate, files_sdk.WithContext(ctx))
	if err != nil {
		resp.Diagnostics.AddError(
			"Error Creating Files Partner",
			"Could not create partner, unexpected error: "+err.Error(),
		)
		return
	}

	diags = r.populateResourceModel(ctx, partner, &plan)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	diags = resp.State.Set(ctx, plan)
	resp.Diagnostics.Append(diags...)
}

func (r *partnerResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	var state partnerResourceModel
	diags := req.State.Get(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	paramsPartnerFind := files_sdk.PartnerFindParams{}
	paramsPartnerFind.Id = state.Id.ValueInt64()

	partner, err := r.client.Find(paramsPartnerFind, files_sdk.WithContext(ctx))
	if err != nil {
		if files_sdk.IsNotExist(err) {
			resp.State.RemoveResource(ctx)
			return
		}

		resp.Diagnostics.AddError(
			"Error Reading Files Partner",
			"Could not read partner id "+fmt.Sprint(state.Id.ValueInt64())+": "+err.Error(),
		)
		return
	}

	diags = r.populateResourceModel(ctx, partner, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	diags = resp.State.Set(ctx, state)
	resp.Diagnostics.Append(diags...)
}

func (r *partnerResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	var plan partnerResourceModel
	diags := req.Plan.Get(ctx, &plan)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}
	var config partnerResourceModel
	diags = req.Config.Get(ctx, &config)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	paramsPartnerUpdate := map[string]interface{}{}
	if !plan.Id.IsNull() && !plan.Id.IsUnknown() {
		paramsPartnerUpdate["id"] = plan.Id.ValueInt64()
	}
	if !config.AiAssistantPersonalityId.IsNull() && !config.AiAssistantPersonalityId.IsUnknown() {
		paramsPartnerUpdate["ai_assistant_personality_id"] = config.AiAssistantPersonalityId.ValueInt64()
	}
	if !config.AllowedIps.IsNull() && !config.AllowedIps.IsUnknown() {
		paramsPartnerUpdate["allowed_ips"] = config.AllowedIps.ValueString()
	}
	if !config.AllowBypassing2faPolicies.IsNull() && !config.AllowBypassing2faPolicies.IsUnknown() {
		paramsPartnerUpdate["allow_bypassing_2fa_policies"] = config.AllowBypassing2faPolicies.ValueBool()
	}
	if !config.AllowCredentialChanges.IsNull() && !config.AllowCredentialChanges.IsUnknown() {
		paramsPartnerUpdate["allow_credential_changes"] = config.AllowCredentialChanges.ValueBool()
	}
	if !config.AllowProvidingGpgKeys.IsNull() && !config.AllowProvidingGpgKeys.IsUnknown() {
		paramsPartnerUpdate["allow_providing_gpg_keys"] = config.AllowProvidingGpgKeys.ValueBool()
	}
	if !config.AllowUserCreation.IsNull() && !config.AllowUserCreation.IsUnknown() {
		paramsPartnerUpdate["allow_user_creation"] = config.AllowUserCreation.ValueBool()
	}
	if !config.CcEmailsToResponsibleParty.IsNull() && !config.CcEmailsToResponsibleParty.IsUnknown() {
		paramsPartnerUpdate["cc_emails_to_responsible_party"] = config.CcEmailsToResponsibleParty.ValueBool()
	}
	if !config.Notes.IsNull() && !config.Notes.IsUnknown() {
		paramsPartnerUpdate["notes"] = config.Notes.ValueString()
	}
	if !config.PartnerChannelTemplateId.IsNull() && !config.PartnerChannelTemplateId.IsUnknown() {
		paramsPartnerUpdate["partner_channel_template_id"] = config.PartnerChannelTemplateId.ValueInt64()
	}
	if !config.ResponsibleGroupId.IsNull() && !config.ResponsibleGroupId.IsUnknown() {
		paramsPartnerUpdate["responsible_group_id"] = config.ResponsibleGroupId.ValueInt64()
	}
	if !config.ResponsibleUserId.IsNull() && !config.ResponsibleUserId.IsUnknown() {
		paramsPartnerUpdate["responsible_user_id"] = config.ResponsibleUserId.ValueInt64()
	}
	if !config.Tags.IsNull() && !config.Tags.IsUnknown() {
		paramsPartnerUpdate["tags"] = config.Tags.ValueString()
	}
	if !config.Name.IsNull() && !config.Name.IsUnknown() {
		paramsPartnerUpdate["name"] = config.Name.ValueString()
	}
	if !config.RootFolder.IsNull() && !config.RootFolder.IsUnknown() {
		paramsPartnerUpdate["root_folder"] = config.RootFolder.ValueString()
	}

	if resp.Diagnostics.HasError() {
		return
	}

	partner, err := r.client.UpdateWithMap(paramsPartnerUpdate, files_sdk.WithContext(ctx))
	if err != nil {
		resp.Diagnostics.AddError(
			"Error Updating Files Partner",
			"Could not update partner, unexpected error: "+err.Error(),
		)
		return
	}

	diags = r.populateResourceModel(ctx, partner, &plan)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	diags = resp.State.Set(ctx, plan)
	resp.Diagnostics.Append(diags...)
}

func (r *partnerResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	var state partnerResourceModel
	diags := req.State.Get(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	paramsPartnerDelete := files_sdk.PartnerDeleteParams{}
	paramsPartnerDelete.Id = state.Id.ValueInt64()

	err := r.client.Delete(paramsPartnerDelete, files_sdk.WithContext(ctx))
	if err != nil && !files_sdk.IsNotExist(err) {
		resp.Diagnostics.AddError(
			"Error Deleting Files Partner",
			"Could not delete partner id "+fmt.Sprint(state.Id.ValueInt64())+": "+err.Error(),
		)
	}
}

func (r *partnerResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
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

func (r *partnerResource) populateResourceModel(ctx context.Context, partner files_sdk.Partner, state *partnerResourceModel) (diags diag.Diagnostics) {
	var propDiags diag.Diagnostics

	state.AllowBypassing2faPolicies = types.BoolPointerValue(partner.AllowBypassing2faPolicies)
	state.AllowedIps = types.StringValue(partner.AllowedIps)
	state.AllowCredentialChanges = types.BoolPointerValue(partner.AllowCredentialChanges)
	state.AllowProvidingGpgKeys = types.BoolPointerValue(partner.AllowProvidingGpgKeys)
	state.AllowUserCreation = types.BoolPointerValue(partner.AllowUserCreation)
	state.CcEmailsToResponsibleParty = types.BoolPointerValue(partner.CcEmailsToResponsibleParty)
	state.Id = types.Int64Value(partner.Id)
	state.AiAssistantPersonalityId = types.Int64Value(partner.AiAssistantPersonalityId)
	state.WorkspaceId = types.Int64Value(partner.WorkspaceId)
	state.Name = types.StringValue(partner.Name)
	state.Notes = types.StringValue(partner.Notes)
	state.PartnerAdminIds, propDiags = types.ListValueFrom(ctx, types.Int64Type, partner.PartnerAdminIds)
	diags.Append(propDiags...)
	state.PartnerChannelTemplateId = types.Int64Value(partner.PartnerChannelTemplateId)
	state.PartnershipRole = types.StringValue(partner.PartnershipRole)
	state.ResponsibleGroupId = types.Int64Value(partner.ResponsibleGroupId)
	state.ResponsibleUserId = types.Int64Value(partner.ResponsibleUserId)
	state.RootFolder = types.StringValue(partner.RootFolder)
	state.Tags = types.StringValue(partner.Tags)
	state.UserIds, propDiags = types.ListValueFrom(ctx, types.Int64Type, partner.UserIds)
	diags.Append(propDiags...)

	return
}
