package provider

import (
	"context"
	"fmt"
	"strconv"
	"strings"

	files_sdk "github.com/Files-com/files-sdk-go/v3"
	clickwrap "github.com/Files-com/files-sdk-go/v3/clickwrap"
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
	_ resource.Resource                = &clickwrapResource{}
	_ resource.ResourceWithConfigure   = &clickwrapResource{}
	_ resource.ResourceWithImportState = &clickwrapResource{}
)

func NewClickwrapResource() resource.Resource {
	return &clickwrapResource{}
}

type clickwrapResource struct {
	client *clickwrap.Client
}

type clickwrapResourceModel struct {
	Name           types.String `tfsdk:"name"`
	Body           types.String `tfsdk:"body"`
	UseWithUsers   types.String `tfsdk:"use_with_users"`
	UseWithBundles types.String `tfsdk:"use_with_bundles"`
	UseWithInboxes types.String `tfsdk:"use_with_inboxes"`
	Id             types.Int64  `tfsdk:"id"`
}

func (r *clickwrapResource) Configure(_ context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {
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

	r.client = &clickwrap.Client{Config: sdk_config}
}

func (r *clickwrapResource) Metadata(_ context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_clickwrap"
}

func (r *clickwrapResource) Schema(_ context.Context, _ resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Description: "A Clickwrap is a legal agreement (such as an NDA or Terms of Use) that your Users and/or Bundle/Inbox participants will need to agree to via a \"Clickwrap\" UI before accessing the site, bundle, or inbox.\n\n\n\nThe values for `use_with_users`, `use_with_bundles`, `use_with_inboxes` are explained as follows:\n\n\n\n* `none` - This Clickwrap may not be used in this context.\n\n* `available_to_all_users` - This Clickwrap may be assigned in this context by any user.\n\n* `available` - This Clickwrap may be assigned in this context, but only by Site Admins. We recognize that the name of this setting is somewhat ambiguous, but we maintain it for legacy reasons.\n\n* `required` - This Clickwrap will always be used in this context, and may not be overridden.",
		Attributes: map[string]schema.Attribute{
			"name": schema.StringAttribute{
				Description: "Name of the Clickwrap agreement (used when selecting from multiple Clickwrap agreements.)",
				Computed:    true,
				Optional:    true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
			},
			"body": schema.StringAttribute{
				Description: "Body text of Clickwrap (supports Markdown formatting).",
				Computed:    true,
				Optional:    true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
			},
			"use_with_users": schema.StringAttribute{
				Description: "Use this Clickwrap for User Registrations?  Note: This only applies to User Registrations where the User is invited to your Files.com site using an E-Mail invitation process where they then set their own password.",
				Computed:    true,
				Optional:    true,
				Validators: []validator.String{
					stringvalidator.OneOf("none", "require"),
				},
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
			},
			"use_with_bundles": schema.StringAttribute{
				Description: "Use this Clickwrap for Bundles?",
				Computed:    true,
				Optional:    true,
				Validators: []validator.String{
					stringvalidator.OneOf("none", "available", "require", "available_to_all_users"),
				},
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
			},
			"use_with_inboxes": schema.StringAttribute{
				Description: "Use this Clickwrap for Inboxes?",
				Computed:    true,
				Optional:    true,
				Validators: []validator.String{
					stringvalidator.OneOf("none", "available", "require", "available_to_all_users"),
				},
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
			},
			"id": schema.Int64Attribute{
				Description: "Clickwrap ID",
				Computed:    true,
				PlanModifiers: []planmodifier.Int64{
					int64planmodifier.UseStateForUnknown(),
				},
			},
		},
	}
}

func (r *clickwrapResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	var plan clickwrapResourceModel
	diags := req.Plan.Get(ctx, &plan)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	paramsClickwrapCreate := files_sdk.ClickwrapCreateParams{}
	paramsClickwrapCreate.Name = plan.Name.ValueString()
	paramsClickwrapCreate.Body = plan.Body.ValueString()
	paramsClickwrapCreate.UseWithBundles = paramsClickwrapCreate.UseWithBundles.Enum()[plan.UseWithBundles.ValueString()]
	paramsClickwrapCreate.UseWithInboxes = paramsClickwrapCreate.UseWithInboxes.Enum()[plan.UseWithInboxes.ValueString()]
	paramsClickwrapCreate.UseWithUsers = paramsClickwrapCreate.UseWithUsers.Enum()[plan.UseWithUsers.ValueString()]

	if resp.Diagnostics.HasError() {
		return
	}

	clickwrap, err := r.client.Create(paramsClickwrapCreate, files_sdk.WithContext(ctx))
	if err != nil {
		resp.Diagnostics.AddError(
			"Error Creating Files Clickwrap",
			"Could not create clickwrap, unexpected error: "+err.Error(),
		)
		return
	}

	diags = r.populateResourceModel(ctx, clickwrap, &plan)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	diags = resp.State.Set(ctx, plan)
	resp.Diagnostics.Append(diags...)
}

func (r *clickwrapResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	var state clickwrapResourceModel
	diags := req.State.Get(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	paramsClickwrapFind := files_sdk.ClickwrapFindParams{}
	paramsClickwrapFind.Id = state.Id.ValueInt64()

	clickwrap, err := r.client.Find(paramsClickwrapFind, files_sdk.WithContext(ctx))
	if err != nil {
		resp.Diagnostics.AddError(
			"Error Reading Files Clickwrap",
			"Could not read clickwrap id "+fmt.Sprint(state.Id.ValueInt64())+": "+err.Error(),
		)
		return
	}

	diags = r.populateResourceModel(ctx, clickwrap, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	diags = resp.State.Set(ctx, state)
	resp.Diagnostics.Append(diags...)
}

func (r *clickwrapResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	var plan clickwrapResourceModel
	diags := req.Plan.Get(ctx, &plan)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	paramsClickwrapUpdate := files_sdk.ClickwrapUpdateParams{}
	paramsClickwrapUpdate.Id = plan.Id.ValueInt64()
	paramsClickwrapUpdate.Name = plan.Name.ValueString()
	paramsClickwrapUpdate.Body = plan.Body.ValueString()
	paramsClickwrapUpdate.UseWithBundles = paramsClickwrapUpdate.UseWithBundles.Enum()[plan.UseWithBundles.ValueString()]
	paramsClickwrapUpdate.UseWithInboxes = paramsClickwrapUpdate.UseWithInboxes.Enum()[plan.UseWithInboxes.ValueString()]
	paramsClickwrapUpdate.UseWithUsers = paramsClickwrapUpdate.UseWithUsers.Enum()[plan.UseWithUsers.ValueString()]

	if resp.Diagnostics.HasError() {
		return
	}

	clickwrap, err := r.client.Update(paramsClickwrapUpdate, files_sdk.WithContext(ctx))
	if err != nil {
		resp.Diagnostics.AddError(
			"Error Updating Files Clickwrap",
			"Could not update clickwrap, unexpected error: "+err.Error(),
		)
		return
	}

	diags = r.populateResourceModel(ctx, clickwrap, &plan)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	diags = resp.State.Set(ctx, plan)
	resp.Diagnostics.Append(diags...)
}

func (r *clickwrapResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	var state clickwrapResourceModel
	diags := req.State.Get(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	paramsClickwrapDelete := files_sdk.ClickwrapDeleteParams{}
	paramsClickwrapDelete.Id = state.Id.ValueInt64()

	err := r.client.Delete(paramsClickwrapDelete, files_sdk.WithContext(ctx))
	if err != nil {
		resp.Diagnostics.AddError(
			"Error Deleting Files Clickwrap",
			"Could not delete clickwrap id "+fmt.Sprint(state.Id.ValueInt64())+": "+err.Error(),
		)
	}
}

func (r *clickwrapResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
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

func (r *clickwrapResource) populateResourceModel(ctx context.Context, clickwrap files_sdk.Clickwrap, state *clickwrapResourceModel) (diags diag.Diagnostics) {
	state.Id = types.Int64Value(clickwrap.Id)
	state.Name = types.StringValue(clickwrap.Name)
	state.Body = types.StringValue(clickwrap.Body)
	state.UseWithUsers = types.StringValue(clickwrap.UseWithUsers)
	state.UseWithBundles = types.StringValue(clickwrap.UseWithBundles)
	state.UseWithInboxes = types.StringValue(clickwrap.UseWithInboxes)

	return
}
