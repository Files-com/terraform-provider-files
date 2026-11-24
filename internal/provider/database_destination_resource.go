package provider

import (
	"context"
	"fmt"
	"strconv"
	"strings"

	files_sdk "github.com/Files-com/files-sdk-go/v3"
	database_destination "github.com/Files-com/files-sdk-go/v3/databasedestination"
	"github.com/Files-com/terraform-provider-files/lib"

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
	_ resource.Resource                = &databaseDestinationResource{}
	_ resource.ResourceWithConfigure   = &databaseDestinationResource{}
	_ resource.ResourceWithImportState = &databaseDestinationResource{}
)

func NewDatabaseDestinationResource() resource.Resource {
	return &databaseDestinationResource{}
}

type databaseDestinationResource struct {
	client *database_destination.Client
}

type databaseDestinationResourceModel struct {
	Name        types.String `tfsdk:"name"`
	Provider_   types.String `tfsdk:"provider_"`
	Config      types.Object `tfsdk:"config"`
	SecretId    types.Int64  `tfsdk:"secret_id"`
	WorkspaceId types.Int64  `tfsdk:"workspace_id"`
	Id          types.Int64  `tfsdk:"id"`
	CreatedAt   types.String `tfsdk:"created_at"`
	UpdatedAt   types.String `tfsdk:"updated_at"`
}

func (r *databaseDestinationResource) Configure(_ context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {
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

	r.client = &database_destination.Client{Config: sdk_config}
}

func (r *databaseDestinationResource) Metadata(_ context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_database_destination"
}

func (r *databaseDestinationResource) Schema(_ context.Context, _ resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Description: "A Database Destination lets an Automation append rows to an existing table. Each destination references a credential in the Secret Store.",
		Attributes: map[string]schema.Attribute{
			"name": schema.StringAttribute{
				Description: "Name of the destination.",
				Required:    true,
			},
			"provider_": schema.StringAttribute{
				Description: "Database provider.",
				Required:    true,
				Validators: []validator.String{
					stringvalidator.OneOf("snowflake", "bigquery", "supabase", "postgresql", "mysql"),
				},
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.RequiresReplace(),
				},
			},
			"config": schema.SingleNestedAttribute{
				Description: "Provider settings for the destination. Use the schema matching provider. The connection address and table are fixed after creation.",
				Required:    true,
				Attributes:  map[string]schema.Attribute{},
			},
			"secret_id": schema.Int64Attribute{
				Description: "ID of a Secret in the same Workspace. Use a Token Secret for Snowflake, a service account Key/Value Secret for BigQuery, or a Username/Password Secret for a database login.",
				Required:    true,
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
			"id": schema.Int64Attribute{
				Description: "Database Destination ID.",
				Computed:    true,
				PlanModifiers: []planmodifier.Int64{
					int64planmodifier.UseStateForUnknown(),
				},
			},
			"created_at": schema.StringAttribute{
				Description: "When the destination was created.",
				Computed:    true,
			},
			"updated_at": schema.StringAttribute{
				Description: "When the destination was last updated.",
				Computed:    true,
			},
		},
	}
}

func (r *databaseDestinationResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	var plan databaseDestinationResourceModel
	diags := req.Plan.Get(ctx, &plan)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}
	var config databaseDestinationResourceModel
	diags = req.Config.Get(ctx, &config)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	paramsDatabaseDestinationCreate := files_sdk.DatabaseDestinationCreateParams{}
	paramsDatabaseDestinationCreate.Name = plan.Name.ValueString()
	paramsDatabaseDestinationCreate.Provider = paramsDatabaseDestinationCreate.Provider.Enum()[plan.Provider_.ValueString()]
	createConfig, diags := lib.SchemaAttributeToInterface(ctx, path.Root("config"), plan.Config)
	resp.Diagnostics.Append(diags...)

	paramsDatabaseDestinationCreate.Config = createConfig
	paramsDatabaseDestinationCreate.SecretId = plan.SecretId.ValueInt64()
	paramsDatabaseDestinationCreate.WorkspaceId = plan.WorkspaceId.ValueInt64()

	if resp.Diagnostics.HasError() {
		return
	}

	databaseDestination, err := r.client.Create(paramsDatabaseDestinationCreate, files_sdk.WithContext(ctx))
	if err != nil {
		resp.Diagnostics.AddError(
			"Error Creating Files DatabaseDestination",
			"Could not create database_destination, unexpected error: "+err.Error(),
		)
		return
	}

	diags = r.populateResourceModel(ctx, databaseDestination, &plan)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	diags = resp.State.Set(ctx, plan)
	resp.Diagnostics.Append(diags...)
}

func (r *databaseDestinationResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	var state databaseDestinationResourceModel
	diags := req.State.Get(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	paramsDatabaseDestinationFind := files_sdk.DatabaseDestinationFindParams{}
	paramsDatabaseDestinationFind.Id = state.Id.ValueInt64()

	databaseDestination, err := r.client.Find(paramsDatabaseDestinationFind, files_sdk.WithContext(ctx))
	if err != nil {
		if files_sdk.IsNotExist(err) {
			resp.State.RemoveResource(ctx)
			return
		}

		resp.Diagnostics.AddError(
			"Error Reading Files DatabaseDestination",
			"Could not read database_destination id "+fmt.Sprint(state.Id.ValueInt64())+": "+err.Error(),
		)
		return
	}

	diags = r.populateResourceModel(ctx, databaseDestination, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	diags = resp.State.Set(ctx, state)
	resp.Diagnostics.Append(diags...)
}

func (r *databaseDestinationResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	var plan databaseDestinationResourceModel
	diags := req.Plan.Get(ctx, &plan)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}
	var config databaseDestinationResourceModel
	diags = req.Config.Get(ctx, &config)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	paramsDatabaseDestinationUpdate := map[string]interface{}{}
	if !plan.Id.IsNull() && !plan.Id.IsUnknown() {
		paramsDatabaseDestinationUpdate["id"] = plan.Id.ValueInt64()
	}
	if !config.Name.IsNull() && !config.Name.IsUnknown() {
		paramsDatabaseDestinationUpdate["name"] = config.Name.ValueString()
	}
	if !config.SecretId.IsNull() && !config.SecretId.IsUnknown() {
		paramsDatabaseDestinationUpdate["secret_id"] = config.SecretId.ValueInt64()
	}
	if !config.Config.IsNull() && !config.Config.IsUnknown() {
		updateConfig, diags := lib.SchemaAttributeToInterface(ctx, path.Root("config"), config.Config)
		resp.Diagnostics.Append(diags...)

		paramsDatabaseDestinationUpdate["config"] = updateConfig
	}

	if resp.Diagnostics.HasError() {
		return
	}

	databaseDestination, err := r.client.UpdateWithMap(paramsDatabaseDestinationUpdate, files_sdk.WithContext(ctx))
	if err != nil {
		resp.Diagnostics.AddError(
			"Error Updating Files DatabaseDestination",
			"Could not update database_destination, unexpected error: "+err.Error(),
		)
		return
	}

	diags = r.populateResourceModel(ctx, databaseDestination, &plan)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	diags = resp.State.Set(ctx, plan)
	resp.Diagnostics.Append(diags...)
}

func (r *databaseDestinationResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	var state databaseDestinationResourceModel
	diags := req.State.Get(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	paramsDatabaseDestinationDelete := files_sdk.DatabaseDestinationDeleteParams{}
	paramsDatabaseDestinationDelete.Id = state.Id.ValueInt64()

	err := r.client.Delete(paramsDatabaseDestinationDelete, files_sdk.WithContext(ctx))
	if err != nil && !files_sdk.IsNotExist(err) {
		resp.Diagnostics.AddError(
			"Error Deleting Files DatabaseDestination",
			"Could not delete database_destination id "+fmt.Sprint(state.Id.ValueInt64())+": "+err.Error(),
		)
	}
}

func (r *databaseDestinationResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
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

func (r *databaseDestinationResource) populateResourceModel(ctx context.Context, databaseDestination files_sdk.DatabaseDestination, state *databaseDestinationResourceModel) (diags diag.Diagnostics) {
	var propDiags diag.Diagnostics

	state.Id = types.Int64Value(databaseDestination.Id)
	state.WorkspaceId = types.Int64Value(databaseDestination.WorkspaceId)
	state.Name = types.StringValue(databaseDestination.Name)
	state.Provider_ = types.StringValue(databaseDestination.Provider)
	configValue := interface{}(databaseDestination.Config)

	state.Config, propDiags = lib.ToObject(ctx, path.Root("config"), configValue, state.Config)
	diags.Append(propDiags...)
	state.SecretId = types.Int64Value(databaseDestination.SecretId)
	if err := lib.TimeToStringType(ctx, path.Root("created_at"), databaseDestination.CreatedAt, &state.CreatedAt); err != nil {
		diags.AddError(
			"Error Creating Files DatabaseDestination",
			"Could not convert state created_at to string: "+err.Error(),
		)
	}
	if err := lib.TimeToStringType(ctx, path.Root("updated_at"), databaseDestination.UpdatedAt, &state.UpdatedAt); err != nil {
		diags.AddError(
			"Error Creating Files DatabaseDestination",
			"Could not convert state updated_at to string: "+err.Error(),
		)
	}

	return
}
