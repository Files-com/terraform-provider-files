package provider

import (
	"context"
	"fmt"

	files_sdk "github.com/Files-com/files-sdk-go/v3"
	secret "github.com/Files-com/files-sdk-go/v3/secret"
	"github.com/Files-com/terraform-provider-files/lib"
	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/datasource/schema"
	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

var (
	_ datasource.DataSource              = &secretDataSource{}
	_ datasource.DataSourceWithConfigure = &secretDataSource{}
)

func NewSecretDataSource() datasource.DataSource {
	return &secretDataSource{}
}

type secretDataSource struct {
	client *secret.Client
}

type secretDataSourceModel struct {
	Id              types.Int64   `tfsdk:"id"`
	WorkspaceId     types.Int64   `tfsdk:"workspace_id"`
	Name            types.String  `tfsdk:"name"`
	Description     types.String  `tfsdk:"description"`
	SecretType      types.String  `tfsdk:"secret_type"`
	Metadata        types.Dynamic `tfsdk:"metadata"`
	ValueFieldNames types.List    `tfsdk:"value_field_names"`
	CreatedAt       types.String  `tfsdk:"created_at"`
	UpdatedAt       types.String  `tfsdk:"updated_at"`
}

func (r *secretDataSource) Configure(_ context.Context, req datasource.ConfigureRequest, resp *datasource.ConfigureResponse) {
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

func (r *secretDataSource) Metadata(_ context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_secret"
}

func (r *secretDataSource) Schema(_ context.Context, _ datasource.SchemaRequest, resp *datasource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Description: "A Secret stores named, typed secret material for later use by features that reference the Secret by ID.\n\n\n\nSecret values are encrypted at rest and are write-only. API responses include metadata and configured value field names, but never include the stored secret values.",
		Attributes: map[string]schema.Attribute{
			"id": schema.Int64Attribute{
				Description: "Secret ID.",
				Required:    true,
			},
			"workspace_id": schema.Int64Attribute{
				Description: "Workspace ID. 0 means the default workspace.",
				Computed:    true,
			},
			"name": schema.StringAttribute{
				Description: "Secret name.",
				Computed:    true,
			},
			"description": schema.StringAttribute{
				Description: "Internal description for your reference.",
				Computed:    true,
			},
			"secret_type": schema.StringAttribute{
				Description: "Secret type.",
				Computed:    true,
			},
			"metadata": schema.DynamicAttribute{
				Description: "Non-secret metadata for the Secret type.",
				Computed:    true,
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

func (r *secretDataSource) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
	var data secretDataSourceModel
	diags := req.Config.Get(ctx, &data)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	paramsSecretFind := files_sdk.SecretFindParams{}
	paramsSecretFind.Id = data.Id.ValueInt64()

	secret, err := r.client.Find(paramsSecretFind, files_sdk.WithContext(ctx))
	if err != nil {
		resp.Diagnostics.AddError(
			"Error Reading Files Secret",
			"Could not read secret id "+fmt.Sprint(data.Id.ValueInt64())+": "+err.Error(),
		)
		return
	}

	diags = r.populateDataSourceModel(ctx, secret, &data)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	diags = resp.State.Set(ctx, data)
	resp.Diagnostics.Append(diags...)
}

func (r *secretDataSource) populateDataSourceModel(ctx context.Context, secret files_sdk.Secret, state *secretDataSourceModel) (diags diag.Diagnostics) {
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
