package provider

import (
	"context"
	"fmt"
	"strconv"
	"strings"

	files_sdk "github.com/Files-com/files-sdk-go/v3"
	key_lifecycle_rule "github.com/Files-com/files-sdk-go/v3/keylifecyclerule"
	"github.com/hashicorp/terraform-plugin-framework-validators/stringvalidator"
	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/int64planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/schema/validator"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

var (
	_ resource.Resource                = &keyLifecycleRuleResource{}
	_ resource.ResourceWithConfigure   = &keyLifecycleRuleResource{}
	_ resource.ResourceWithImportState = &keyLifecycleRuleResource{}
)

func NewKeyLifecycleRuleResource() resource.Resource {
	return &keyLifecycleRuleResource{}
}

type keyLifecycleRuleResource struct {
	client *key_lifecycle_rule.Client
}

type keyLifecycleRuleResourceModel struct {
	KeyType        types.String `tfsdk:"key_type"`
	InactivityDays types.Int64  `tfsdk:"inactivity_days"`
	Name           types.String `tfsdk:"name"`
	Id             types.Int64  `tfsdk:"id"`
}

func (r *keyLifecycleRuleResource) Configure(_ context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {
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

	r.client = &key_lifecycle_rule.Client{Config: sdk_config}
}

func (r *keyLifecycleRuleResource) Metadata(_ context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_key_lifecycle_rule"
}

func (r *keyLifecycleRuleResource) Schema(_ context.Context, _ resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Description: "A KeyLifecycleRule represents a rule that applies to GPG keys and SSH keys (also called User Public Keys) based on their inactivity.\n\n\n\nKeys that have been unused for the specified number of days will be deleted.",
		Attributes: map[string]schema.Attribute{
			"key_type": schema.StringAttribute{
				Description: "Key type for which the rule will apply (gpg or ssh).",
				Computed:    true,
				Optional:    true,
				Validators: []validator.String{
					stringvalidator.OneOf("gpg", "ssh"),
				},
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
			},
			"inactivity_days": schema.Int64Attribute{
				Description: "Number of days of inactivity before the rule applies.",
				Computed:    true,
				Optional:    true,
				PlanModifiers: []planmodifier.Int64{
					int64planmodifier.UseStateForUnknown(),
				},
			},
			"name": schema.StringAttribute{
				Description: "Key Lifecycle Rule name",
				Computed:    true,
				Optional:    true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
			},
			"id": schema.Int64Attribute{
				Description: "Key Lifecycle Rule ID",
				Computed:    true,
				PlanModifiers: []planmodifier.Int64{
					int64planmodifier.UseStateForUnknown(),
				},
			},
		},
	}
}

func (r *keyLifecycleRuleResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	var plan keyLifecycleRuleResourceModel
	diags := req.Plan.Get(ctx, &plan)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}
	var config keyLifecycleRuleResourceModel
	diags = req.Config.Get(ctx, &config)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	paramsKeyLifecycleRuleCreate := files_sdk.KeyLifecycleRuleCreateParams{}
	paramsKeyLifecycleRuleCreate.KeyType = paramsKeyLifecycleRuleCreate.KeyType.Enum()[plan.KeyType.ValueString()]
	paramsKeyLifecycleRuleCreate.InactivityDays = plan.InactivityDays.ValueInt64()
	paramsKeyLifecycleRuleCreate.Name = plan.Name.ValueString()

	if resp.Diagnostics.HasError() {
		return
	}

	keyLifecycleRule, err := r.client.Create(paramsKeyLifecycleRuleCreate, files_sdk.WithContext(ctx))
	if err != nil {
		resp.Diagnostics.AddError(
			"Error Creating Files KeyLifecycleRule",
			"Could not create key_lifecycle_rule, unexpected error: "+err.Error(),
		)
		return
	}

	diags = r.populateResourceModel(ctx, keyLifecycleRule, &plan)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	diags = resp.State.Set(ctx, plan)
	resp.Diagnostics.Append(diags...)
}

func (r *keyLifecycleRuleResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	var state keyLifecycleRuleResourceModel
	diags := req.State.Get(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	paramsKeyLifecycleRuleFind := files_sdk.KeyLifecycleRuleFindParams{}
	paramsKeyLifecycleRuleFind.Id = state.Id.ValueInt64()

	keyLifecycleRule, err := r.client.Find(paramsKeyLifecycleRuleFind, files_sdk.WithContext(ctx))
	if err != nil {
		if files_sdk.IsNotExist(err) {
			resp.State.RemoveResource(ctx)
			return
		}

		resp.Diagnostics.AddError(
			"Error Reading Files KeyLifecycleRule",
			"Could not read key_lifecycle_rule id "+fmt.Sprint(state.Id.ValueInt64())+": "+err.Error(),
		)
		return
	}

	diags = r.populateResourceModel(ctx, keyLifecycleRule, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	diags = resp.State.Set(ctx, state)
	resp.Diagnostics.Append(diags...)
}

func (r *keyLifecycleRuleResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	var plan keyLifecycleRuleResourceModel
	diags := req.Plan.Get(ctx, &plan)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}
	var config keyLifecycleRuleResourceModel
	diags = req.Config.Get(ctx, &config)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	paramsKeyLifecycleRuleUpdate := map[string]interface{}{}
	if !plan.Id.IsNull() && !plan.Id.IsUnknown() {
		paramsKeyLifecycleRuleUpdate["id"] = plan.Id.ValueInt64()
	}
	if !config.KeyType.IsNull() && !config.KeyType.IsUnknown() {
		paramsKeyLifecycleRuleUpdate["key_type"] = config.KeyType.ValueString()
	}
	if !config.InactivityDays.IsNull() && !config.InactivityDays.IsUnknown() {
		paramsKeyLifecycleRuleUpdate["inactivity_days"] = config.InactivityDays.ValueInt64()
	}
	if !config.Name.IsNull() && !config.Name.IsUnknown() {
		paramsKeyLifecycleRuleUpdate["name"] = config.Name.ValueString()
	}

	if resp.Diagnostics.HasError() {
		return
	}

	keyLifecycleRule, err := r.client.UpdateWithMap(paramsKeyLifecycleRuleUpdate, files_sdk.WithContext(ctx))
	if err != nil {
		resp.Diagnostics.AddError(
			"Error Updating Files KeyLifecycleRule",
			"Could not update key_lifecycle_rule, unexpected error: "+err.Error(),
		)
		return
	}

	diags = r.populateResourceModel(ctx, keyLifecycleRule, &plan)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	diags = resp.State.Set(ctx, plan)
	resp.Diagnostics.Append(diags...)
}

func (r *keyLifecycleRuleResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	var state keyLifecycleRuleResourceModel
	diags := req.State.Get(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	paramsKeyLifecycleRuleDelete := files_sdk.KeyLifecycleRuleDeleteParams{}
	paramsKeyLifecycleRuleDelete.Id = state.Id.ValueInt64()

	err := r.client.Delete(paramsKeyLifecycleRuleDelete, files_sdk.WithContext(ctx))
	if err != nil && !files_sdk.IsNotExist(err) {
		resp.Diagnostics.AddError(
			"Error Deleting Files KeyLifecycleRule",
			"Could not delete key_lifecycle_rule id "+fmt.Sprint(state.Id.ValueInt64())+": "+err.Error(),
		)
	}
}

func (r *keyLifecycleRuleResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
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

func (r *keyLifecycleRuleResource) populateResourceModel(ctx context.Context, keyLifecycleRule files_sdk.KeyLifecycleRule, state *keyLifecycleRuleResourceModel) (diags diag.Diagnostics) {
	state.Id = types.Int64Value(keyLifecycleRule.Id)
	state.KeyType = types.StringValue(keyLifecycleRule.KeyType)
	state.InactivityDays = types.Int64Value(keyLifecycleRule.InactivityDays)
	state.Name = types.StringValue(keyLifecycleRule.Name)

	return
}
