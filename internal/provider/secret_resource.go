package provider

import (
	"context"
	"fmt"
	"strconv"
	"strings"

	files_sdk "github.com/Files-com/files-sdk-go/v3"
	secret "github.com/Files-com/files-sdk-go/v3/secret"
	"github.com/Files-com/terraform-provider-files/lib"
	"github.com/hashicorp/terraform-plugin-framework-validators/stringvalidator"
	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/dynamicplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/int64planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/schema/validator"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

var (
	_ resource.Resource                = &secretResource{}
	_ resource.ResourceWithConfigure   = &secretResource{}
	_ resource.ResourceWithImportState = &secretResource{}
)

func NewSecretResource() resource.Resource {
	return &secretResource{}
}

type secretResource struct {
	client *secret.Client
}

type secretResourceModel struct {
	Name            types.String  `tfsdk:"name"`
	SecretType      types.String  `tfsdk:"secret_type"`
	WorkspaceId     types.Int64   `tfsdk:"workspace_id"`
	Description     types.String  `tfsdk:"description"`
	Metadata        types.Dynamic `tfsdk:"metadata"`
	Id              types.Int64   `tfsdk:"id"`
	ValueFieldNames types.List    `tfsdk:"value_field_names"`
	CreatedAt       types.String  `tfsdk:"created_at"`
	UpdatedAt       types.String  `tfsdk:"updated_at"`
}

func (r *secretResource) Configure(_ context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {
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

	r.client = &secret.Client{Config: sdk_config}
}

func (r *secretResource) Metadata(_ context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_secret"
}

func (r *secretResource) Schema(_ context.Context, _ resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Description: "A Secret stores named, typed secret material for later use by features that reference the Secret by ID.\n\n\n\nSecret values are encrypted at rest and are write-only. API responses include metadata and configured value field names, but never include the stored secret values.",
		Attributes: map[string]schema.Attribute{
			"name": schema.StringAttribute{
				Description: "Secret name.",
				Required:    true,
			},
			"secret_type": schema.StringAttribute{
				Description: "Secret type.",
				Required:    true,
				Validators: []validator.String{
					stringvalidator.OneOf("basic", "token", "headers", "certificate", "key_value"),
				},
			},
			"workspace_id": schema.Int64Attribute{
				Description: "Workspace ID. 0 means the default workspace.",
				Computed:    true,
				Optional:    true,
				PlanModifiers: []planmodifier.Int64{
					int64planmodifier.UseStateForUnknown(),
					int64planmodifier.RequiresReplace(),
				},
			},
			"description": schema.StringAttribute{
				Description: "Internal description for your reference.",
				Computed:    true,
				Optional:    true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
			},
			"metadata": schema.DynamicAttribute{
				Description: "Non-secret metadata for the Secret type.",
				Computed:    true,
				Optional:    true,
				PlanModifiers: []planmodifier.Dynamic{
					dynamicplanmodifier.UseStateForUnknown(),
				},
			},
			"id": schema.Int64Attribute{
				Description: "Secret ID.",
				Computed:    true,
				PlanModifiers: []planmodifier.Int64{
					int64planmodifier.UseStateForUnknown(),
				},
			},
			"value_field_names": schema.ListAttribute{
				Description: "Names of configured secret value fields. Secret values are never returned.",
				Computed:    true,
				ElementType: types.StringType,
			},
			"created_at": schema.StringAttribute{
				Description: "Secret create date/time.",
				Computed:    true,
			},
			"updated_at": schema.StringAttribute{
				Description: "Secret update date/time.",
				Computed:    true,
			},
		},
	}
}

func (r *secretResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	var plan secretResourceModel
	diags := req.Plan.Get(ctx, &plan)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}
	var config secretResourceModel
	diags = req.Config.Get(ctx, &config)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	paramsSecretCreate := files_sdk.SecretCreateParams{}
	paramsSecretCreate.Name = plan.Name.ValueString()
	paramsSecretCreate.Description = plan.Description.ValueString()
	paramsSecretCreate.SecretType = paramsSecretCreate.SecretType.Enum()[plan.SecretType.ValueString()]
	createMetadata, diags := lib.DynamicToInterface(ctx, path.Root("metadata"), plan.Metadata)
	resp.Diagnostics.Append(diags...)
	paramsSecretCreate.Metadata = createMetadata
	paramsSecretCreate.WorkspaceId = plan.WorkspaceId.ValueInt64()

	if resp.Diagnostics.HasError() {
		return
	}

	secret, err := r.client.Create(paramsSecretCreate, files_sdk.WithContext(ctx))
	if err != nil {
		resp.Diagnostics.AddError(
			"Error Creating Files Secret",
			"Could not create secret, unexpected error: "+err.Error(),
		)
		return
	}

	diags = r.populateResourceModel(ctx, secret, &plan)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	diags = resp.State.Set(ctx, plan)
	resp.Diagnostics.Append(diags...)
}

func (r *secretResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	var state secretResourceModel
	diags := req.State.Get(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	paramsSecretFind := files_sdk.SecretFindParams{}
	paramsSecretFind.Id = state.Id.ValueInt64()

	secret, err := r.client.Find(paramsSecretFind, files_sdk.WithContext(ctx))
	if err != nil {
		if files_sdk.IsNotExist(err) {
			resp.State.RemoveResource(ctx)
			return
		}

		resp.Diagnostics.AddError(
			"Error Reading Files Secret",
			"Could not read secret id "+fmt.Sprint(state.Id.ValueInt64())+": "+err.Error(),
		)
		return
	}

	diags = r.populateResourceModel(ctx, secret, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	diags = resp.State.Set(ctx, state)
	resp.Diagnostics.Append(diags...)
}

func (r *secretResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	var plan secretResourceModel
	diags := req.Plan.Get(ctx, &plan)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}
	var config secretResourceModel
	diags = req.Config.Get(ctx, &config)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	paramsSecretUpdate := map[string]interface{}{}
	if !plan.Id.IsNull() && !plan.Id.IsUnknown() {
		paramsSecretUpdate["id"] = plan.Id.ValueInt64()
	}
	if !config.Name.IsNull() && !config.Name.IsUnknown() {
		paramsSecretUpdate["name"] = config.Name.ValueString()
	}
	if !config.Description.IsNull() && !config.Description.IsUnknown() {
		paramsSecretUpdate["description"] = config.Description.ValueString()
	}
	if !config.SecretType.IsNull() && !config.SecretType.IsUnknown() {
		paramsSecretUpdate["secret_type"] = config.SecretType.ValueString()
	}
	updateMetadata, diags := lib.DynamicToInterface(ctx, path.Root("metadata"), config.Metadata)
	resp.Diagnostics.Append(diags...)
	paramsSecretUpdate["metadata"] = updateMetadata

	if resp.Diagnostics.HasError() {
		return
	}

	secret, err := r.client.UpdateWithMap(paramsSecretUpdate, files_sdk.WithContext(ctx))
	if err != nil {
		resp.Diagnostics.AddError(
			"Error Updating Files Secret",
			"Could not update secret, unexpected error: "+err.Error(),
		)
		return
	}

	diags = r.populateResourceModel(ctx, secret, &plan)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	diags = resp.State.Set(ctx, plan)
	resp.Diagnostics.Append(diags...)
}

func (r *secretResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	var state secretResourceModel
	diags := req.State.Get(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	paramsSecretDelete := files_sdk.SecretDeleteParams{}
	paramsSecretDelete.Id = state.Id.ValueInt64()

	err := r.client.Delete(paramsSecretDelete, files_sdk.WithContext(ctx))
	if err != nil && !files_sdk.IsNotExist(err) {
		resp.Diagnostics.AddError(
			"Error Deleting Files Secret",
			"Could not delete secret id "+fmt.Sprint(state.Id.ValueInt64())+": "+err.Error(),
		)
	}
}

func (r *secretResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
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

func (r *secretResource) populateResourceModel(ctx context.Context, secret files_sdk.Secret, state *secretResourceModel) (diags diag.Diagnostics) {
	var propDiags diag.Diagnostics

	state.Id = types.Int64Value(secret.Id)
	state.WorkspaceId = types.Int64Value(secret.WorkspaceId)
	state.Name = types.StringValue(secret.Name)
	state.Description = types.StringValue(secret.Description)
	state.SecretType = types.StringValue(secret.SecretType)
	state.Metadata, propDiags = lib.ToDynamic(ctx, path.Root("metadata"), secret.Metadata, state.Metadata.UnderlyingValue())
	diags.Append(propDiags...)
	state.ValueFieldNames, propDiags = types.ListValueFrom(ctx, types.StringType, secret.ValueFieldNames)
	diags.Append(propDiags...)
	if err := lib.TimeToStringType(ctx, path.Root("created_at"), secret.CreatedAt, &state.CreatedAt); err != nil {
		diags.AddError(
			"Error Creating Files Secret",
			"Could not convert state created_at to string: "+err.Error(),
		)
	}
	if err := lib.TimeToStringType(ctx, path.Root("updated_at"), secret.UpdatedAt, &state.UpdatedAt); err != nil {
		diags.AddError(
			"Error Creating Files Secret",
			"Could not convert state updated_at to string: "+err.Error(),
		)
	}

	return
}
