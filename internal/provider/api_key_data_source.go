package provider

import (
	"context"
	"fmt"

	files_sdk "github.com/Files-com/files-sdk-go/v3"
	api_key "github.com/Files-com/files-sdk-go/v3/apikey"
	"github.com/Files-com/terraform-provider-files/lib"
	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/datasource/schema"
	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

var (
	_ datasource.DataSource              = &apiKeyDataSource{}
	_ datasource.DataSourceWithConfigure = &apiKeyDataSource{}
)

func NewApiKeyDataSource() datasource.DataSource {
	return &apiKeyDataSource{}
}

type apiKeyDataSource struct {
	client *api_key.Client
}

type apiKeyDataSourceModel struct {
	Id               types.Int64  `tfsdk:"id"`
	DescriptiveLabel types.String `tfsdk:"descriptive_label"`
	Description      types.String `tfsdk:"description"`
	CreatedAt        types.String `tfsdk:"created_at"`
	ExpiresAt        types.String `tfsdk:"expires_at"`
	Key              types.String `tfsdk:"key"`
	LastUseAt        types.String `tfsdk:"last_use_at"`
	Name             types.String `tfsdk:"name"`
	PermissionSet    types.String `tfsdk:"permission_set"`
	Platform         types.String `tfsdk:"platform"`
	Url              types.String `tfsdk:"url"`
	UserId           types.Int64  `tfsdk:"user_id"`
}

func (r *apiKeyDataSource) Configure(_ context.Context, req datasource.ConfigureRequest, resp *datasource.ConfigureResponse) {
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

	r.client = &api_key.Client{Config: sdk_config}
}

func (r *apiKeyDataSource) Metadata(_ context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_api_key"
}

func (r *apiKeyDataSource) Schema(_ context.Context, _ datasource.SchemaRequest, resp *datasource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Description: "An APIKey is a key that allows programmatic access to your Site.\n\n\n\nAPI keys confer all the permissions of the user who owns them.\n\nIf an API key is created without a user owner, it is considered a site-wide API key, which has full permissions to do anything on the Site.",
		Attributes: map[string]schema.Attribute{
			"id": schema.Int64Attribute{
				Description: "API Key ID",
				Required:    true,
			},
			"descriptive_label": schema.StringAttribute{
				Description: "Unique label that describes this API key.  Useful for external systems where you may have API keys from multiple accounts and want a human-readable label for each key.",
				Computed:    true,
			},
			"description": schema.StringAttribute{
				Description: "User-supplied description of API key.",
				Computed:    true,
			},
			"created_at": schema.StringAttribute{
				Description: "Time which API Key was created",
				Computed:    true,
			},
			"expires_at": schema.StringAttribute{
				Description: "API Key expiration date",
				Computed:    true,
			},
			"key": schema.StringAttribute{
				Description: "API Key actual key string",
				Computed:    true,
			},
			"last_use_at": schema.StringAttribute{
				Description: "API Key last used - note this value is only updated once per 3 hour period, so the 'actual' time of last use may be up to 3 hours later than this timestamp.",
				Computed:    true,
			},
			"name": schema.StringAttribute{
				Description: "Internal name for the API Key.  For your use.",
				Computed:    true,
			},
			"permission_set": schema.StringAttribute{
				Description: "Permissions for this API Key. It must be full for site-wide API Keys.  Keys with the `desktop_app` permission set only have the ability to do the functions provided in our Desktop App (File and Share Link operations).  Additional permission sets may become available in the future, such as for a Site Admin to give a key with no administrator privileges.  If you have ideas for permission sets, please let us know.",
				Computed:    true,
			},
			"platform": schema.StringAttribute{
				Description: "If this API key represents a Desktop app, what platform was it created on?",
				Computed:    true,
			},
			"url": schema.StringAttribute{
				Description: "URL for API host.",
				Computed:    true,
			},
			"user_id": schema.Int64Attribute{
				Description: "User ID for the owner of this API Key.  May be blank for Site-wide API Keys.",
				Computed:    true,
			},
		},
	}
}

func (r *apiKeyDataSource) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
	var data apiKeyDataSourceModel
	diags := req.Config.Get(ctx, &data)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	paramsApiKeyFind := files_sdk.ApiKeyFindParams{}
	paramsApiKeyFind.Id = data.Id.ValueInt64()

	apiKey, err := r.client.Find(paramsApiKeyFind, files_sdk.WithContext(ctx))
	if err != nil {
		resp.Diagnostics.AddError(
			"Error Reading Files ApiKey",
			"Could not read api_key id "+fmt.Sprint(data.Id.ValueInt64())+": "+err.Error(),
		)
		return
	}

	diags = r.populateDataSourceModel(ctx, apiKey, &data)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	diags = resp.State.Set(ctx, data)
	resp.Diagnostics.Append(diags...)
}

func (r *apiKeyDataSource) populateDataSourceModel(ctx context.Context, apiKey files_sdk.ApiKey, state *apiKeyDataSourceModel) (diags diag.Diagnostics) {
	state.Id = types.Int64Value(apiKey.Id)
	state.DescriptiveLabel = types.StringValue(apiKey.DescriptiveLabel)
	state.Description = types.StringValue(apiKey.Description)
	if err := lib.TimeToStringType(ctx, path.Root("created_at"), apiKey.CreatedAt, &state.CreatedAt); err != nil {
		diags.AddError(
			"Error Creating Files ApiKey",
			"Could not convert state created_at to string: "+err.Error(),
		)
	}
	if err := lib.TimeToStringType(ctx, path.Root("expires_at"), apiKey.ExpiresAt, &state.ExpiresAt); err != nil {
		diags.AddError(
			"Error Creating Files ApiKey",
			"Could not convert state expires_at to string: "+err.Error(),
		)
	}
	state.Key = types.StringValue(apiKey.Key)
	if err := lib.TimeToStringType(ctx, path.Root("last_use_at"), apiKey.LastUseAt, &state.LastUseAt); err != nil {
		diags.AddError(
			"Error Creating Files ApiKey",
			"Could not convert state last_use_at to string: "+err.Error(),
		)
	}
	state.Name = types.StringValue(apiKey.Name)
	state.PermissionSet = types.StringValue(apiKey.PermissionSet)
	state.Platform = types.StringValue(apiKey.Platform)
	state.Url = types.StringValue(apiKey.Url)
	state.UserId = types.Int64Value(apiKey.UserId)

	return
}
