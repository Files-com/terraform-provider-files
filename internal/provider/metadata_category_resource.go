package provider

import (
	"context"
	"fmt"
	"strconv"
	"strings"

	files_sdk "github.com/Files-com/files-sdk-go/v3"
	metadata_category "github.com/Files-com/files-sdk-go/v3/metadatacategory"
	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/int64planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/listplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

var (
	_ resource.Resource                = &metadataCategoryResource{}
	_ resource.ResourceWithConfigure   = &metadataCategoryResource{}
	_ resource.ResourceWithImportState = &metadataCategoryResource{}
)

func NewMetadataCategoryResource() resource.Resource {
	return &metadataCategoryResource{}
}

type metadataCategoryResource struct {
	client *metadata_category.Client
}

type metadataCategoryResourceModel struct {
	Name           types.String `tfsdk:"name"`
	DefaultColumns types.List   `tfsdk:"default_columns"`
	Id             types.Int64  `tfsdk:"id"`
	Definitions    types.Map    `tfsdk:"definitions"`
}

func (r *metadataCategoryResource) Configure(_ context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {
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

	r.client = &metadata_category.Client{Config: sdk_config}
}

func (r *metadataCategoryResource) Metadata(_ context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_metadata_category"
}

func (r *metadataCategoryResource) Schema(_ context.Context, _ resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Description: "A MetadataCategory defines a reusable set of Custom Metadata rules that can be assigned to folders\n\nvia a folder behavior. Each category specifies named metadata keys with optional allowed-value\n\nconstraints, and a set of default columns to display in the UI.\n\n\n\nIf a key's `allowed_values` array is empty, it is treated as a free-form text field.\n\nIf the array is non-empty, the key is constrained to those values in the Web UI.",
		Attributes: map[string]schema.Attribute{
			"name": schema.StringAttribute{
				Description: "Name of the metadata category.",
				Required:    true,
			},
			"default_columns": schema.ListAttribute{
				Description: "Metadata keys that should appear as columns in the UI by default.",
				Computed:    true,
				Optional:    true,
				ElementType: types.StringType,
				PlanModifiers: []planmodifier.List{
					listplanmodifier.UseStateForUnknown(),
				},
			},
			"id": schema.Int64Attribute{
				Description: "Metadata Category ID",
				Computed:    true,
				PlanModifiers: []planmodifier.Int64{
					int64planmodifier.UseStateForUnknown(),
				},
			},
			"definitions": schema.MapAttribute{
				Description: "Map of key names to arrays of allowed values. An empty array means free-form text.",
				Computed:    true,
				ElementType: types.ListType{ElemType: types.StringType},
			},
		},
	}
}

func (r *metadataCategoryResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	var plan metadataCategoryResourceModel
	diags := req.Plan.Get(ctx, &plan)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}
	var config metadataCategoryResourceModel
	diags = req.Config.Get(ctx, &config)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	paramsMetadataCategoryCreate := files_sdk.MetadataCategoryCreateParams{}
	paramsMetadataCategoryCreate.Name = plan.Name.ValueString()
	if !plan.DefaultColumns.IsNull() && !plan.DefaultColumns.IsUnknown() {
		diags = plan.DefaultColumns.ElementsAs(ctx, &paramsMetadataCategoryCreate.DefaultColumns, false)
		resp.Diagnostics.Append(diags...)
	}

	if resp.Diagnostics.HasError() {
		return
	}

	metadataCategory, err := r.client.Create(paramsMetadataCategoryCreate, files_sdk.WithContext(ctx))
	if err != nil {
		resp.Diagnostics.AddError(
			"Error Creating Files MetadataCategory",
			"Could not create metadata_category, unexpected error: "+err.Error(),
		)
		return
	}

	diags = r.populateResourceModel(ctx, metadataCategory, &plan)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	diags = resp.State.Set(ctx, plan)
	resp.Diagnostics.Append(diags...)
}

func (r *metadataCategoryResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	var state metadataCategoryResourceModel
	diags := req.State.Get(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	paramsMetadataCategoryFind := files_sdk.MetadataCategoryFindParams{}
	paramsMetadataCategoryFind.Id = state.Id.ValueInt64()

	metadataCategory, err := r.client.Find(paramsMetadataCategoryFind, files_sdk.WithContext(ctx))
	if err != nil {
		if files_sdk.IsNotExist(err) {
			resp.State.RemoveResource(ctx)
			return
		}

		resp.Diagnostics.AddError(
			"Error Reading Files MetadataCategory",
			"Could not read metadata_category id "+fmt.Sprint(state.Id.ValueInt64())+": "+err.Error(),
		)
		return
	}

	diags = r.populateResourceModel(ctx, metadataCategory, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	diags = resp.State.Set(ctx, state)
	resp.Diagnostics.Append(diags...)
}

func (r *metadataCategoryResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	var plan metadataCategoryResourceModel
	diags := req.Plan.Get(ctx, &plan)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}
	var config metadataCategoryResourceModel
	diags = req.Config.Get(ctx, &config)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	paramsMetadataCategoryUpdate := map[string]interface{}{}
	if !plan.Id.IsNull() && !plan.Id.IsUnknown() {
		paramsMetadataCategoryUpdate["id"] = plan.Id.ValueInt64()
	}
	if !config.Name.IsNull() && !config.Name.IsUnknown() {
		paramsMetadataCategoryUpdate["name"] = config.Name.ValueString()
	}
	if !config.DefaultColumns.IsNull() && !config.DefaultColumns.IsUnknown() {
		var updateDefaultColumns []string
		diags = config.DefaultColumns.ElementsAs(ctx, &updateDefaultColumns, false)
		resp.Diagnostics.Append(diags...)
		paramsMetadataCategoryUpdate["default_columns"] = updateDefaultColumns
	}

	if resp.Diagnostics.HasError() {
		return
	}

	metadataCategory, err := r.client.UpdateWithMap(paramsMetadataCategoryUpdate, files_sdk.WithContext(ctx))
	if err != nil {
		resp.Diagnostics.AddError(
			"Error Updating Files MetadataCategory",
			"Could not update metadata_category, unexpected error: "+err.Error(),
		)
		return
	}

	diags = r.populateResourceModel(ctx, metadataCategory, &plan)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	diags = resp.State.Set(ctx, plan)
	resp.Diagnostics.Append(diags...)
}

func (r *metadataCategoryResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	var state metadataCategoryResourceModel
	diags := req.State.Get(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	paramsMetadataCategoryDelete := files_sdk.MetadataCategoryDeleteParams{}
	paramsMetadataCategoryDelete.Id = state.Id.ValueInt64()

	err := r.client.Delete(paramsMetadataCategoryDelete, files_sdk.WithContext(ctx))
	if err != nil && !files_sdk.IsNotExist(err) {
		resp.Diagnostics.AddError(
			"Error Deleting Files MetadataCategory",
			"Could not delete metadata_category id "+fmt.Sprint(state.Id.ValueInt64())+": "+err.Error(),
		)
	}
}

func (r *metadataCategoryResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
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

func (r *metadataCategoryResource) populateResourceModel(ctx context.Context, metadataCategory files_sdk.MetadataCategory, state *metadataCategoryResourceModel) (diags diag.Diagnostics) {
	var propDiags diag.Diagnostics

	state.Id = types.Int64Value(metadataCategory.Id)
	state.Name = types.StringValue(metadataCategory.Name)
	state.Definitions, propDiags = types.MapValueFrom(ctx, types.ListType{ElemType: types.StringType}, metadataCategory.Definitions)
	diags.Append(propDiags...)
	state.DefaultColumns, propDiags = types.ListValueFrom(ctx, types.StringType, metadataCategory.DefaultColumns)
	diags.Append(propDiags...)

	return
}
