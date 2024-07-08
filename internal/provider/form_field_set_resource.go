package provider

import (
	"context"
	"fmt"
	"strconv"
	"strings"

	files_sdk "github.com/Files-com/files-sdk-go/v3"
	form_field_set "github.com/Files-com/files-sdk-go/v3/formfieldset"
	"github.com/Files-com/terraform-provider-files/lib"
	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/boolplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/dynamicplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/int64planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

var (
	_ resource.Resource                = &formFieldSetResource{}
	_ resource.ResourceWithConfigure   = &formFieldSetResource{}
	_ resource.ResourceWithImportState = &formFieldSetResource{}
)

func NewFormFieldSetResource() resource.Resource {
	return &formFieldSetResource{}
}

type formFieldSetResource struct {
	client *form_field_set.Client
}

type formFieldSetResourceModel struct {
	Id          types.Int64   `tfsdk:"id"`
	Title       types.String  `tfsdk:"title"`
	FormLayout  types.List    `tfsdk:"form_layout"`
	FormFields  types.Dynamic `tfsdk:"form_fields"`
	SkipName    types.Bool    `tfsdk:"skip_name"`
	SkipEmail   types.Bool    `tfsdk:"skip_email"`
	SkipCompany types.Bool    `tfsdk:"skip_company"`
	UserId      types.Int64   `tfsdk:"user_id"`
}

func (r *formFieldSetResource) Configure(_ context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {
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

	r.client = &form_field_set.Client{Config: sdk_config}
}

func (r *formFieldSetResource) Metadata(_ context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_form_field_set"
}

func (r *formFieldSetResource) Schema(_ context.Context, _ resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Description: "Form Field Sets define a custom forms to be used for bundle and inbox registrations.\n\n \n\n Each Form Field Set contains one or more Form Fields. A form and all of its form fields are submitted in a single create request. The order of form fields in the array is the order they will be displayed.\n\n \n\n Once created, a form field set can then be associated with one or more bundle(s) and/or inbox(s). Once associated, you will be required to submit well-formatted form-data when creating a bundle-registration or inbox registration.",
		Attributes: map[string]schema.Attribute{
			"id": schema.Int64Attribute{
				Description: "Form field set id",
				Computed:    true,
				PlanModifiers: []planmodifier.Int64{
					int64planmodifier.UseStateForUnknown(),
				},
			},
			"title": schema.StringAttribute{
				Description: "Title to be displayed",
				Computed:    true,
				Optional:    true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
			},
			"form_layout": schema.ListAttribute{
				Description: "Layout of the form",
				Computed:    true,
				ElementType: types.Int64Type,
			},
			"form_fields": schema.DynamicAttribute{
				Description: "Associated form fields",
				Computed:    true,
				Optional:    true,
				PlanModifiers: []planmodifier.Dynamic{
					dynamicplanmodifier.UseStateForUnknown(),
				},
			},
			"skip_name": schema.BoolAttribute{
				Description: "Any associated InboxRegistrations or BundleRegistrations can be saved without providing name",
				Computed:    true,
				Optional:    true,
				PlanModifiers: []planmodifier.Bool{
					boolplanmodifier.UseStateForUnknown(),
				},
			},
			"skip_email": schema.BoolAttribute{
				Description: "Any associated InboxRegistrations or BundleRegistrations can be saved without providing email",
				Computed:    true,
				Optional:    true,
				PlanModifiers: []planmodifier.Bool{
					boolplanmodifier.UseStateForUnknown(),
				},
			},
			"skip_company": schema.BoolAttribute{
				Description: "Any associated InboxRegistrations or BundleRegistrations can be saved without providing company",
				Computed:    true,
				Optional:    true,
				PlanModifiers: []planmodifier.Bool{
					boolplanmodifier.UseStateForUnknown(),
				},
			},
			"user_id": schema.Int64Attribute{
				Description: "User ID.  Provide a value of `0` to operate the current session's user.",
				Optional:    true,
				PlanModifiers: []planmodifier.Int64{
					int64planmodifier.RequiresReplace(),
				},
			},
		},
	}
}

func (r *formFieldSetResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	var plan formFieldSetResourceModel
	diags := req.Plan.Get(ctx, &plan)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	paramsFormFieldSetCreate := files_sdk.FormFieldSetCreateParams{}
	paramsFormFieldSetCreate.UserId = plan.UserId.ValueInt64()
	paramsFormFieldSetCreate.Title = plan.Title.ValueString()
	paramsFormFieldSetCreate.SkipEmail = plan.SkipEmail.ValueBoolPointer()
	paramsFormFieldSetCreate.SkipName = plan.SkipName.ValueBoolPointer()
	paramsFormFieldSetCreate.SkipCompany = plan.SkipCompany.ValueBoolPointer()
	paramsFormFieldSetCreate.FormFields, diags = lib.DynamicToStringMapSlice(ctx, path.Root("form_fields"), plan.FormFields)
	resp.Diagnostics.Append(diags...)

	if resp.Diagnostics.HasError() {
		return
	}

	formFieldSet, err := r.client.Create(paramsFormFieldSetCreate, files_sdk.WithContext(ctx))
	if err != nil {
		resp.Diagnostics.AddError(
			"Error Creating Files FormFieldSet",
			"Could not create form_field_set, unexpected error: "+err.Error(),
		)
		return
	}

	diags = r.populateResourceModel(ctx, formFieldSet, &plan)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	diags = resp.State.Set(ctx, plan)
	resp.Diagnostics.Append(diags...)
}

func (r *formFieldSetResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	var state formFieldSetResourceModel
	diags := req.State.Get(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	paramsFormFieldSetFind := files_sdk.FormFieldSetFindParams{}
	paramsFormFieldSetFind.Id = state.Id.ValueInt64()

	formFieldSet, err := r.client.Find(paramsFormFieldSetFind, files_sdk.WithContext(ctx))
	if err != nil {
		resp.Diagnostics.AddError(
			"Error Reading Files FormFieldSet",
			"Could not read form_field_set id "+fmt.Sprint(state.Id.ValueInt64())+": "+err.Error(),
		)
		return
	}

	diags = r.populateResourceModel(ctx, formFieldSet, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	diags = resp.State.Set(ctx, state)
	resp.Diagnostics.Append(diags...)
}

func (r *formFieldSetResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	var plan formFieldSetResourceModel
	diags := req.Plan.Get(ctx, &plan)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	paramsFormFieldSetUpdate := files_sdk.FormFieldSetUpdateParams{}
	paramsFormFieldSetUpdate.Id = plan.Id.ValueInt64()
	paramsFormFieldSetUpdate.Title = plan.Title.ValueString()
	paramsFormFieldSetUpdate.SkipEmail = plan.SkipEmail.ValueBoolPointer()
	paramsFormFieldSetUpdate.SkipName = plan.SkipName.ValueBoolPointer()
	paramsFormFieldSetUpdate.SkipCompany = plan.SkipCompany.ValueBoolPointer()
	paramsFormFieldSetUpdate.FormFields, diags = lib.DynamicToStringMapSlice(ctx, path.Root("form_fields"), plan.FormFields)
	resp.Diagnostics.Append(diags...)

	if resp.Diagnostics.HasError() {
		return
	}

	formFieldSet, err := r.client.Update(paramsFormFieldSetUpdate, files_sdk.WithContext(ctx))
	if err != nil {
		resp.Diagnostics.AddError(
			"Error Updating Files FormFieldSet",
			"Could not update form_field_set, unexpected error: "+err.Error(),
		)
		return
	}

	diags = r.populateResourceModel(ctx, formFieldSet, &plan)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	diags = resp.State.Set(ctx, plan)
	resp.Diagnostics.Append(diags...)
}

func (r *formFieldSetResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	var state formFieldSetResourceModel
	diags := req.State.Get(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	paramsFormFieldSetDelete := files_sdk.FormFieldSetDeleteParams{}
	paramsFormFieldSetDelete.Id = state.Id.ValueInt64()

	err := r.client.Delete(paramsFormFieldSetDelete, files_sdk.WithContext(ctx))
	if err != nil {
		resp.Diagnostics.AddError(
			"Error Deleting Files FormFieldSet",
			"Could not delete form_field_set id "+fmt.Sprint(state.Id.ValueInt64())+": "+err.Error(),
		)
	}
}

func (r *formFieldSetResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
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

func (r *formFieldSetResource) populateResourceModel(ctx context.Context, formFieldSet files_sdk.FormFieldSet, state *formFieldSetResourceModel) (diags diag.Diagnostics) {
	var propDiags diag.Diagnostics

	state.Id = types.Int64Value(formFieldSet.Id)
	state.Title = types.StringValue(formFieldSet.Title)
	state.FormLayout, propDiags = types.ListValueFrom(ctx, types.Int64Type, formFieldSet.FormLayout)
	diags.Append(propDiags...)
	state.FormFields, propDiags = lib.ToDynamic(ctx, path.Root("form_fields"), formFieldSet.FormFields, state.FormFields.UnderlyingValue())
	diags.Append(propDiags...)
	state.SkipName = types.BoolPointerValue(formFieldSet.SkipName)
	state.SkipEmail = types.BoolPointerValue(formFieldSet.SkipEmail)
	state.SkipCompany = types.BoolPointerValue(formFieldSet.SkipCompany)

	return
}
