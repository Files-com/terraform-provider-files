package provider

import (
	"context"
	"fmt"

	files_sdk "github.com/Files-com/files-sdk-go/v3"
	form_field_set "github.com/Files-com/files-sdk-go/v3/formfieldset"
	"github.com/Files-com/terraform-provider-files/lib"
	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/datasource/schema"
	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/path"

	"github.com/hashicorp/terraform-plugin-framework/types"
)

var (
	_ datasource.DataSource              = &formFieldSetDataSource{}
	_ datasource.DataSourceWithConfigure = &formFieldSetDataSource{}
)

func NewFormFieldSetDataSource() datasource.DataSource {
	return &formFieldSetDataSource{}
}

type formFieldSetDataSource struct {
	client *form_field_set.Client
}

type formFieldSetDataSourceModel struct {
	Id                             types.Int64   `tfsdk:"id"`
	Title                          types.String  `tfsdk:"title"`
	RequireCompleteMetadataCapture types.Bool    `tfsdk:"require_complete_metadata_capture"`
	FormLayout                     types.List    `tfsdk:"form_layout"`
	FormFields                     types.Dynamic `tfsdk:"form_fields"`
	SkipName                       types.Bool    `tfsdk:"skip_name"`
	SkipEmail                      types.Bool    `tfsdk:"skip_email"`
	SkipCompany                    types.Bool    `tfsdk:"skip_company"`
	InUse                          types.Bool    `tfsdk:"in_use"`
}

func (r *formFieldSetDataSource) Configure(_ context.Context, req datasource.ConfigureRequest, resp *datasource.ConfigureResponse) {
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

func (r *formFieldSetDataSource) Metadata(_ context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_form_field_set"
}

func (r *formFieldSetDataSource) Schema(_ context.Context, _ datasource.SchemaRequest, resp *datasource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Description: "A Form Field Set is a custom form to be used for bundle and inbox registrations.\n\nEach Form Field Set contains one or more Form Fields. A form and all of its form fields are submitted in a single create request. The order of form fields in the array is the order they will be displayed.\n\nOnce created, a form field set can then be associated with one or more bundle(s) and/or inbox(s). Once associated, you will be required to submit well-formatted form-data when creating a bundle-registration or inbox registration.\n\nCustom answers saved on a registration are saved as Custom Metadata on each file uploaded through that registration, with field labels as keys. This applies to existing registrations as well as new ones. Unanswered optional fields and values exceeding 1024 characters are omitted; the first 32 remaining answers in saved form-answer order are captured. Dates use YYYY-MM-DD. Existing answers with unavailable field definitions or unparseable dates are omitted. Built-in name, email, and company fields are excluded. Form answers take precedence over explicit upload metadata with matching keys, including keeping omitted answer labels absent.\n\nEnable require_complete_metadata_capture to reject answers exceeding these value or count limits instead of omitting them. Uploaded files use the labels and answers accepted at registration even if the form is edited later.\n\nEditing a field definition creates a replacement field with a new ID and preserves the original definition for existing registration answers. Unchanged fields and fields moved within the layout retain their IDs. Use the returned form_layout when submitting new answers or editing the form again.\n\nAuthenticated form field set responses include historical definitions in form_fields so existing answers can still be interpreted. Only IDs in form_layout are current fields. Forms embedded in Share Link and Inbox responses contain current definitions only.",
		Attributes: map[string]schema.Attribute{
			"id": schema.Int64Attribute{
				Description: "Form field set id",
				Required:    true,
			},
			"title": schema.StringAttribute{
				Description: "Title to be displayed",
				Computed:    true,
			},
			"require_complete_metadata_capture": schema.BoolAttribute{
				Description: "Reject answers over 1024 characters or more than 32 answered custom fields to ensure complete file metadata capture",
				Computed:    true,
			},
			"form_layout": schema.ListAttribute{
				Description: "Layout of the form",
				Computed:    true,
				ElementType: types.Int64Type,
			},
			"form_fields": schema.DynamicAttribute{
				Description: "Form fields in display order. Managed state contains the current fields listed in form_layout. The data source also returns historical definitions.",
				Computed:    true,
			},
			"skip_name": schema.BoolAttribute{
				Description: "Any associated InboxRegistrations or BundleRegistrations can be saved without providing name",
				Computed:    true,
			},
			"skip_email": schema.BoolAttribute{
				Description: "Any associated InboxRegistrations or BundleRegistrations can be saved without providing email",
				Computed:    true,
			},
			"skip_company": schema.BoolAttribute{
				Description: "Any associated InboxRegistrations or BundleRegistrations can be saved without providing company",
				Computed:    true,
			},
			"in_use": schema.BoolAttribute{
				Description: "Form Field Set is in use by an active Inbox / Bundle / Inbox Registration / Bundle Registration",
				Computed:    true,
			},
		},
	}
}

func (r *formFieldSetDataSource) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
	var data formFieldSetDataSourceModel
	diags := req.Config.Get(ctx, &data)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	paramsFormFieldSetFind := files_sdk.FormFieldSetFindParams{}
	paramsFormFieldSetFind.Id = data.Id.ValueInt64()

	var responseBody []byte
	formFieldSet, err := r.client.Find(paramsFormFieldSetFind, files_sdk.WithContext(ctx), lib.CaptureResponseBody(&responseBody))
	if err != nil {
		resp.Diagnostics.AddError(
			"Error Reading Files FormFieldSet",
			"Could not read form_field_set id "+fmt.Sprint(data.Id.ValueInt64())+": "+err.Error(),
		)
		return
	}

	diags = r.populateDataSourceModel(ctx, formFieldSet, &data, responseBody)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	diags = resp.State.Set(ctx, data)
	resp.Diagnostics.Append(diags...)
}

func (r *formFieldSetDataSource) populateDataSourceModel(ctx context.Context, formFieldSet files_sdk.FormFieldSet, state *formFieldSetDataSourceModel, responseBody []byte) (diags diag.Diagnostics) {
	var propDiags diag.Diagnostics

	state.Id = types.Int64Value(formFieldSet.Id)
	state.Title = types.StringValue(formFieldSet.Title)
	state.RequireCompleteMetadataCapture = types.BoolPointerValue(formFieldSet.RequireCompleteMetadataCapture)
	state.FormLayout, propDiags = types.ListValueFrom(ctx, types.Int64Type, formFieldSet.FormLayout)
	diags.Append(propDiags...)
	formFieldsResponse, formFieldsResponseDiags := lib.ResponseProperty(path.Root("form_fields"), responseBody, "form_fields")
	diags.Append(formFieldsResponseDiags...)
	state.FormFields, propDiags = lib.ToDynamic(ctx, path.Root("form_fields"), formFieldsResponse, nil)
	diags.Append(propDiags...)
	state.SkipName = types.BoolPointerValue(formFieldSet.SkipName)
	state.SkipEmail = types.BoolPointerValue(formFieldSet.SkipEmail)
	state.SkipCompany = types.BoolPointerValue(formFieldSet.SkipCompany)
	state.InUse = types.BoolPointerValue(formFieldSet.InUse)

	return
}
