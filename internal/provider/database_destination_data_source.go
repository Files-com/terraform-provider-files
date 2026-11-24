package provider

import (
	"context"
	"fmt"

	files_sdk "github.com/Files-com/files-sdk-go/v3"
	database_destination "github.com/Files-com/files-sdk-go/v3/databasedestination"
	"github.com/Files-com/terraform-provider-files/lib"
	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/datasource/schema"
	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/path"

	"github.com/hashicorp/terraform-plugin-framework/types"
)

var (
	_ datasource.DataSource              = &databaseDestinationDataSource{}
	_ datasource.DataSourceWithConfigure = &databaseDestinationDataSource{}
)

func NewDatabaseDestinationDataSource() datasource.DataSource {
	return &databaseDestinationDataSource{}
}

type databaseDestinationDataSource struct {
	client *database_destination.Client
}

type databaseDestinationDataSourceModel struct {
	Id          types.Int64  `tfsdk:"id"`
	WorkspaceId types.Int64  `tfsdk:"workspace_id"`
	Name        types.String `tfsdk:"name"`
	Provider_   types.String `tfsdk:"provider_"`
	Config      types.Object `tfsdk:"config"`
	SecretId    types.Int64  `tfsdk:"secret_id"`
	CreatedAt   types.String `tfsdk:"created_at"`
	UpdatedAt   types.String `tfsdk:"updated_at"`
}

func (r *databaseDestinationDataSource) Configure(_ context.Context, req datasource.ConfigureRequest, resp *datasource.ConfigureResponse) {
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

func (r *databaseDestinationDataSource) Metadata(_ context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_database_destination"
}

func (r *databaseDestinationDataSource) Schema(_ context.Context, _ datasource.SchemaRequest, resp *datasource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Description: "A Database Destination lets an Automation append rows to an existing table. Each destination references a credential in the Secret Store.",
		Attributes: map[string]schema.Attribute{
			"id": schema.Int64Attribute{
				Description: "Database Destination ID.",
				Required:    true,
			},
			"workspace_id": schema.Int64Attribute{
				Description: "Workspace ID. 0 means the default workspace.",
				Computed:    true,
			},
			"name": schema.StringAttribute{
				Description: "Name of the destination.",
				Computed:    true,
			},
			"provider_": schema.StringAttribute{
				Description: "Database provider.",
				Computed:    true,
			},
			"config": schema.SingleNestedAttribute{
				Description: "Provider settings for the destination. Use the schema matching provider. The connection address and table are fixed after creation.",
				Computed:    true,
				Attributes:  map[string]schema.Attribute{},
			},
			"secret_id": schema.Int64Attribute{
				Description: "ID of a Secret in the same Workspace. Use a Token Secret for Snowflake, a service account Key/Value Secret for BigQuery, or a Username/Password Secret for a database login.",
				Computed:    true,
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

func (r *databaseDestinationDataSource) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
	var data databaseDestinationDataSourceModel
	diags := req.Config.Get(ctx, &data)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	paramsDatabaseDestinationFind := files_sdk.DatabaseDestinationFindParams{}
	paramsDatabaseDestinationFind.Id = data.Id.ValueInt64()

	databaseDestination, err := r.client.Find(paramsDatabaseDestinationFind, files_sdk.WithContext(ctx))
	if err != nil {
		resp.Diagnostics.AddError(
			"Error Reading Files DatabaseDestination",
			"Could not read database_destination id "+fmt.Sprint(data.Id.ValueInt64())+": "+err.Error(),
		)
		return
	}

	diags = r.populateDataSourceModel(ctx, databaseDestination, &data)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	diags = resp.State.Set(ctx, data)
	resp.Diagnostics.Append(diags...)
}

func (r *databaseDestinationDataSource) populateDataSourceModel(ctx context.Context, databaseDestination files_sdk.DatabaseDestination, state *databaseDestinationDataSourceModel) (diags diag.Diagnostics) {
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
