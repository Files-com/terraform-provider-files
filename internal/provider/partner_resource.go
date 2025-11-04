package provider

import (
	"context"
	"fmt"
	"strconv"
	"strings"

	files_sdk "github.com/Files-com/files-sdk-go/v3"
	partner "github.com/Files-com/files-sdk-go/v3/partner"
	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/boolplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/int64planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringplanmodifier"
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
	AllowBypassing2faPolicies types.Bool   `tfsdk:"allow_bypassing_2fa_policies"`
	AllowCredentialChanges    types.Bool   `tfsdk:"allow_credential_changes"`
	AllowProvidingGpgKeys     types.Bool   `tfsdk:"allow_providing_gpg_keys"`
	AllowUserCreation         types.Bool   `tfsdk:"allow_user_creation"`
	Name                      types.String `tfsdk:"name"`
	Notes                     types.String `tfsdk:"notes"`
	RootFolder                types.String `tfsdk:"root_folder"`
	Tags                      types.String `tfsdk:"tags"`
	Id                        types.Int64  `tfsdk:"id"`
	PartnerAdminIds           types.List   `tfsdk:"partner_admin_ids"`
	UserIds                   types.List   `tfsdk:"user_ids"`
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
			"allow_bypassing_2fa_policies": schema.BoolAttribute{
				Description: "Allow users created under this Partner to bypass Two-Factor Authentication policies.",
				Computed:    true,
				Optional:    true,
				PlanModifiers: []planmodifier.Bool{
					boolplanmodifier.UseStateForUnknown(),
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
			"name": schema.StringAttribute{
				Description: "The name of the Partner.",
				Computed:    true,
				Optional:    true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
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
			"root_folder": schema.StringAttribute{
				Description: "The root folder path for this Partner.",
				Computed:    true,
				Optional:    true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
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
	paramsPartnerCreate.Name = plan.Name.ValueString()
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
	paramsPartnerCreate.Notes = plan.Notes.ValueString()
	paramsPartnerCreate.RootFolder = plan.RootFolder.ValueString()
	paramsPartnerCreate.Tags = plan.Tags.ValueString()

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

	paramsPartnerUpdate := files_sdk.PartnerUpdateParams{}
	paramsPartnerUpdate.Id = plan.Id.ValueInt64()
	paramsPartnerUpdate.Name = plan.Name.ValueString()
	if !plan.AllowBypassing2faPolicies.IsNull() && !plan.AllowBypassing2faPolicies.IsUnknown() {
		paramsPartnerUpdate.AllowBypassing2faPolicies = plan.AllowBypassing2faPolicies.ValueBoolPointer()
	}
	if !plan.AllowCredentialChanges.IsNull() && !plan.AllowCredentialChanges.IsUnknown() {
		paramsPartnerUpdate.AllowCredentialChanges = plan.AllowCredentialChanges.ValueBoolPointer()
	}
	if !plan.AllowProvidingGpgKeys.IsNull() && !plan.AllowProvidingGpgKeys.IsUnknown() {
		paramsPartnerUpdate.AllowProvidingGpgKeys = plan.AllowProvidingGpgKeys.ValueBoolPointer()
	}
	if !plan.AllowUserCreation.IsNull() && !plan.AllowUserCreation.IsUnknown() {
		paramsPartnerUpdate.AllowUserCreation = plan.AllowUserCreation.ValueBoolPointer()
	}
	paramsPartnerUpdate.Notes = plan.Notes.ValueString()
	paramsPartnerUpdate.RootFolder = plan.RootFolder.ValueString()
	paramsPartnerUpdate.Tags = plan.Tags.ValueString()

	if resp.Diagnostics.HasError() {
		return
	}

	partner, err := r.client.Update(paramsPartnerUpdate, files_sdk.WithContext(ctx))
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
	state.AllowCredentialChanges = types.BoolPointerValue(partner.AllowCredentialChanges)
	state.AllowProvidingGpgKeys = types.BoolPointerValue(partner.AllowProvidingGpgKeys)
	state.AllowUserCreation = types.BoolPointerValue(partner.AllowUserCreation)
	state.Id = types.Int64Value(partner.Id)
	state.Name = types.StringValue(partner.Name)
	state.Notes = types.StringValue(partner.Notes)
	state.PartnerAdminIds, propDiags = types.ListValueFrom(ctx, types.Int64Type, partner.PartnerAdminIds)
	diags.Append(propDiags...)
	state.RootFolder = types.StringValue(partner.RootFolder)
	state.Tags = types.StringValue(partner.Tags)
	state.UserIds, propDiags = types.ListValueFrom(ctx, types.Int64Type, partner.UserIds)
	diags.Append(propDiags...)

	return
}
