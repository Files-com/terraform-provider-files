package provider

import (
	"context"
	"fmt"
	"strconv"
	"strings"

	files_sdk "github.com/Files-com/files-sdk-go/v3"
	user_lifecycle_rule "github.com/Files-com/files-sdk-go/v3/userlifecyclerule"
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
	_ resource.Resource                = &userLifecycleRuleResource{}
	_ resource.ResourceWithConfigure   = &userLifecycleRuleResource{}
	_ resource.ResourceWithImportState = &userLifecycleRuleResource{}
)

func NewUserLifecycleRuleResource() resource.Resource {
	return &userLifecycleRuleResource{}
}

type userLifecycleRuleResource struct {
	client *user_lifecycle_rule.Client
}

type userLifecycleRuleResourceModel struct {
	AuthenticationMethod types.String `tfsdk:"authentication_method"`
	InactivityDays       types.Int64  `tfsdk:"inactivity_days"`
	IncludeFolderAdmins  types.Bool   `tfsdk:"include_folder_admins"`
	IncludeSiteAdmins    types.Bool   `tfsdk:"include_site_admins"`
	Action               types.String `tfsdk:"action"`
	UserState            types.String `tfsdk:"user_state"`
	Id                   types.Int64  `tfsdk:"id"`
	SiteId               types.Int64  `tfsdk:"site_id"`
}

func (r *userLifecycleRuleResource) Configure(_ context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {
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

	r.client = &user_lifecycle_rule.Client{Config: sdk_config}
}

func (r *userLifecycleRuleResource) Metadata(_ context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_user_lifecycle_rule"
}

func (r *userLifecycleRuleResource) Schema(_ context.Context, _ resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Description: "A UserLifecycleRule represents a rule that applies to users based on their inactivity, state and authentication method.\n\n\n\nThe rule either disable or delete users who have been inactive or disabled for a specified number of days.\n\n\n\nThe authentication_method property specifies the authentication method for the rule, which can be set to \"all\" or other specific methods.\n\n\n\nThe rule can also include or exclude site and folder admins from the action.",
		Attributes: map[string]schema.Attribute{
			"authentication_method": schema.StringAttribute{
				Description: "User authentication method for the rule",
				Computed:    true,
				Optional:    true,
				Validators: []validator.String{
					stringvalidator.OneOf("all", "password", "sso", "none", "email_signup", "password_with_imported_hash", "password_and_ssh_key"),
				},
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
			},
			"inactivity_days": schema.Int64Attribute{
				Description: "Number of days of inactivity before the rule applies",
				Computed:    true,
				Optional:    true,
				PlanModifiers: []planmodifier.Int64{
					int64planmodifier.UseStateForUnknown(),
				},
			},
			"include_folder_admins": schema.BoolAttribute{
				Description: "Include folder admins in the rule",
				Computed:    true,
				Optional:    true,
				PlanModifiers: []planmodifier.Bool{
					boolplanmodifier.UseStateForUnknown(),
				},
			},
			"include_site_admins": schema.BoolAttribute{
				Description: "Include site admins in the rule",
				Computed:    true,
				Optional:    true,
				PlanModifiers: []planmodifier.Bool{
					boolplanmodifier.UseStateForUnknown(),
				},
			},
			"action": schema.StringAttribute{
				Description: "Action to take on inactive users (disable or delete)",
				Computed:    true,
				Optional:    true,
				Validators: []validator.String{
					stringvalidator.OneOf("disable", "delete"),
				},
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
			},
			"user_state": schema.StringAttribute{
				Description: "State of the users to apply the rule to (inactive or disabled)",
				Computed:    true,
				Optional:    true,
				Validators: []validator.String{
					stringvalidator.OneOf("inactive", "disabled"),
				},
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
			},
			"id": schema.Int64Attribute{
				Description: "User Lifecycle Rule ID",
				Computed:    true,
				PlanModifiers: []planmodifier.Int64{
					int64planmodifier.UseStateForUnknown(),
				},
			},
			"site_id": schema.Int64Attribute{
				Description: "Site ID",
				Computed:    true,
			},
		},
	}
}

func (r *userLifecycleRuleResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	var plan userLifecycleRuleResourceModel
	diags := req.Plan.Get(ctx, &plan)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	paramsUserLifecycleRuleCreate := files_sdk.UserLifecycleRuleCreateParams{}
	paramsUserLifecycleRuleCreate.Action = paramsUserLifecycleRuleCreate.Action.Enum()[plan.Action.ValueString()]
	paramsUserLifecycleRuleCreate.AuthenticationMethod = paramsUserLifecycleRuleCreate.AuthenticationMethod.Enum()[plan.AuthenticationMethod.ValueString()]
	paramsUserLifecycleRuleCreate.InactivityDays = plan.InactivityDays.ValueInt64()
	if !plan.IncludeSiteAdmins.IsNull() && !plan.IncludeSiteAdmins.IsUnknown() {
		paramsUserLifecycleRuleCreate.IncludeSiteAdmins = plan.IncludeSiteAdmins.ValueBoolPointer()
	}
	if !plan.IncludeFolderAdmins.IsNull() && !plan.IncludeFolderAdmins.IsUnknown() {
		paramsUserLifecycleRuleCreate.IncludeFolderAdmins = plan.IncludeFolderAdmins.ValueBoolPointer()
	}
	paramsUserLifecycleRuleCreate.UserState = paramsUserLifecycleRuleCreate.UserState.Enum()[plan.UserState.ValueString()]

	if resp.Diagnostics.HasError() {
		return
	}

	userLifecycleRule, err := r.client.Create(paramsUserLifecycleRuleCreate, files_sdk.WithContext(ctx))
	if err != nil {
		resp.Diagnostics.AddError(
			"Error Creating Files UserLifecycleRule",
			"Could not create user_lifecycle_rule, unexpected error: "+err.Error(),
		)
		return
	}

	diags = r.populateResourceModel(ctx, userLifecycleRule, &plan)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	diags = resp.State.Set(ctx, plan)
	resp.Diagnostics.Append(diags...)
}

func (r *userLifecycleRuleResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	var state userLifecycleRuleResourceModel
	diags := req.State.Get(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	paramsUserLifecycleRuleFind := files_sdk.UserLifecycleRuleFindParams{}
	paramsUserLifecycleRuleFind.Id = state.Id.ValueInt64()

	userLifecycleRule, err := r.client.Find(paramsUserLifecycleRuleFind, files_sdk.WithContext(ctx))
	if err != nil {
		if files_sdk.IsNotExist(err) {
			resp.State.RemoveResource(ctx)
			return
		}

		resp.Diagnostics.AddError(
			"Error Reading Files UserLifecycleRule",
			"Could not read user_lifecycle_rule id "+fmt.Sprint(state.Id.ValueInt64())+": "+err.Error(),
		)
		return
	}

	diags = r.populateResourceModel(ctx, userLifecycleRule, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	diags = resp.State.Set(ctx, state)
	resp.Diagnostics.Append(diags...)
}

func (r *userLifecycleRuleResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	var plan userLifecycleRuleResourceModel
	diags := req.Plan.Get(ctx, &plan)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	paramsUserLifecycleRuleUpdate := files_sdk.UserLifecycleRuleUpdateParams{}
	paramsUserLifecycleRuleUpdate.Id = plan.Id.ValueInt64()
	paramsUserLifecycleRuleUpdate.Action = paramsUserLifecycleRuleUpdate.Action.Enum()[plan.Action.ValueString()]
	paramsUserLifecycleRuleUpdate.AuthenticationMethod = paramsUserLifecycleRuleUpdate.AuthenticationMethod.Enum()[plan.AuthenticationMethod.ValueString()]
	paramsUserLifecycleRuleUpdate.InactivityDays = plan.InactivityDays.ValueInt64()
	if !plan.IncludeSiteAdmins.IsNull() && !plan.IncludeSiteAdmins.IsUnknown() {
		paramsUserLifecycleRuleUpdate.IncludeSiteAdmins = plan.IncludeSiteAdmins.ValueBoolPointer()
	}
	if !plan.IncludeFolderAdmins.IsNull() && !plan.IncludeFolderAdmins.IsUnknown() {
		paramsUserLifecycleRuleUpdate.IncludeFolderAdmins = plan.IncludeFolderAdmins.ValueBoolPointer()
	}
	paramsUserLifecycleRuleUpdate.UserState = paramsUserLifecycleRuleUpdate.UserState.Enum()[plan.UserState.ValueString()]

	if resp.Diagnostics.HasError() {
		return
	}

	userLifecycleRule, err := r.client.Update(paramsUserLifecycleRuleUpdate, files_sdk.WithContext(ctx))
	if err != nil {
		resp.Diagnostics.AddError(
			"Error Updating Files UserLifecycleRule",
			"Could not update user_lifecycle_rule, unexpected error: "+err.Error(),
		)
		return
	}

	diags = r.populateResourceModel(ctx, userLifecycleRule, &plan)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	diags = resp.State.Set(ctx, plan)
	resp.Diagnostics.Append(diags...)
}

func (r *userLifecycleRuleResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	var state userLifecycleRuleResourceModel
	diags := req.State.Get(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	paramsUserLifecycleRuleDelete := files_sdk.UserLifecycleRuleDeleteParams{}
	paramsUserLifecycleRuleDelete.Id = state.Id.ValueInt64()

	err := r.client.Delete(paramsUserLifecycleRuleDelete, files_sdk.WithContext(ctx))
	if err != nil && !files_sdk.IsNotExist(err) {
		resp.Diagnostics.AddError(
			"Error Deleting Files UserLifecycleRule",
			"Could not delete user_lifecycle_rule id "+fmt.Sprint(state.Id.ValueInt64())+": "+err.Error(),
		)
	}
}

func (r *userLifecycleRuleResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
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

func (r *userLifecycleRuleResource) populateResourceModel(ctx context.Context, userLifecycleRule files_sdk.UserLifecycleRule, state *userLifecycleRuleResourceModel) (diags diag.Diagnostics) {
	state.Id = types.Int64Value(userLifecycleRule.Id)
	state.AuthenticationMethod = types.StringValue(userLifecycleRule.AuthenticationMethod)
	state.InactivityDays = types.Int64Value(userLifecycleRule.InactivityDays)
	state.IncludeFolderAdmins = types.BoolPointerValue(userLifecycleRule.IncludeFolderAdmins)
	state.IncludeSiteAdmins = types.BoolPointerValue(userLifecycleRule.IncludeSiteAdmins)
	state.Action = types.StringValue(userLifecycleRule.Action)
	state.UserState = types.StringValue(userLifecycleRule.UserState)
	state.SiteId = types.Int64Value(userLifecycleRule.SiteId)

	return
}
