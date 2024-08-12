package provider

import (
	"context"
	"fmt"
	"strconv"
	"strings"
	"time"

	files_sdk "github.com/Files-com/files-sdk-go/v3"
	api_key "github.com/Files-com/files-sdk-go/v3/apikey"
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
	_ resource.Resource                = &apiKeyResource{}
	_ resource.ResourceWithConfigure   = &apiKeyResource{}
	_ resource.ResourceWithImportState = &apiKeyResource{}
)

func NewApiKeyResource() resource.Resource {
	return &apiKeyResource{}
}

type apiKeyResource struct {
	client *api_key.Client
}

type apiKeyResourceModel struct {
	Name             types.String `tfsdk:"name"`
	Description      types.String `tfsdk:"description"`
	ExpiresAt        types.String `tfsdk:"expires_at"`
	PermissionSet    types.String `tfsdk:"permission_set"`
	UserId           types.Int64  `tfsdk:"user_id"`
	Path             types.String `tfsdk:"path"`
	Id               types.Int64  `tfsdk:"id"`
	DescriptiveLabel types.String `tfsdk:"descriptive_label"`
	CreatedAt        types.String `tfsdk:"created_at"`
	Key              types.String `tfsdk:"key"`
	LastUseAt        types.String `tfsdk:"last_use_at"`
	Platform         types.String `tfsdk:"platform"`
	Url              types.String `tfsdk:"url"`
}

func (r *apiKeyResource) Configure(_ context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {
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

func (r *apiKeyResource) Metadata(_ context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_api_key"
}

func (r *apiKeyResource) Schema(_ context.Context, _ resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Description: "An APIKey is a key that allows programmatic access to your Site.\n\n\n\nAPI keys confer all the permissions of the user who owns them.\n\nIf an API key is created without a user owner, it is considered a site-wide API key, which has full permissions to do anything on the Site.",
		Attributes: map[string]schema.Attribute{
			"name": schema.StringAttribute{
				Description: "Internal name for the API Key.  For your use.",
				Required:    true,
			},
			"description": schema.StringAttribute{
				Description: "User-supplied description of API key.",
				Computed:    true,
				Optional:    true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
			},
			"expires_at": schema.StringAttribute{
				Description: "API Key expiration date",
				Computed:    true,
				Optional:    true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
			},
			"permission_set": schema.StringAttribute{
				Description: "Permissions for this API Key. It must be full for site-wide API Keys.  Keys with the `desktop_app` permission set only have the ability to do the functions provided in our Desktop App (File and Share Link operations).  Additional permission sets may become available in the future, such as for a Site Admin to give a key with no administrator privileges.  If you have ideas for permission sets, please let us know.",
				Computed:    true,
				Optional:    true,
				Validators: []validator.String{
					stringvalidator.OneOf("none", "full", "desktop_app", "sync_app", "office_integration", "mobile_app"),
				},
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
			},
			"user_id": schema.Int64Attribute{
				Description: "User ID for the owner of this API Key.  May be blank for Site-wide API Keys.",
				Computed:    true,
				Optional:    true,
				PlanModifiers: []planmodifier.Int64{
					int64planmodifier.UseStateForUnknown(),
					int64planmodifier.RequiresReplace(),
				},
			},
			"path": schema.StringAttribute{
				Description: "Folder path restriction for this API key.",
				Optional:    true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.RequiresReplace(),
				},
			},
			"id": schema.Int64Attribute{
				Description: "API Key ID",
				Computed:    true,
				PlanModifiers: []planmodifier.Int64{
					int64planmodifier.UseStateForUnknown(),
				},
			},
			"descriptive_label": schema.StringAttribute{
				Description: "Unique label that describes this API key.  Useful for external systems where you may have API keys from multiple accounts and want a human-readable label for each key.",
				Computed:    true,
			},
			"created_at": schema.StringAttribute{
				Description: "Time which API Key was created",
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
			"platform": schema.StringAttribute{
				Description: "If this API key represents a Desktop app, what platform was it created on?",
				Computed:    true,
			},
			"url": schema.StringAttribute{
				Description: "URL for API host.",
				Computed:    true,
			},
		},
	}
}

func (r *apiKeyResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	var plan apiKeyResourceModel
	diags := req.Plan.Get(ctx, &plan)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	paramsApiKeyCreate := files_sdk.ApiKeyCreateParams{}
	paramsApiKeyCreate.UserId = plan.UserId.ValueInt64()
	paramsApiKeyCreate.Description = plan.Description.ValueString()
	if !plan.ExpiresAt.IsNull() && plan.ExpiresAt.ValueString() != "" {
		createExpiresAt, err := time.Parse(time.RFC3339, plan.ExpiresAt.ValueString())
		if err != nil {
			resp.Diagnostics.AddAttributeError(
				path.Root("expires_at"),
				"Error Parsing expires_at Time",
				"Could not parse expires_at time: "+err.Error(),
			)
		} else {
			paramsApiKeyCreate.ExpiresAt = &createExpiresAt
		}
	}
	paramsApiKeyCreate.PermissionSet = paramsApiKeyCreate.PermissionSet.Enum()[plan.PermissionSet.ValueString()]
	paramsApiKeyCreate.Name = plan.Name.ValueString()
	paramsApiKeyCreate.Path = plan.Path.ValueString()

	if resp.Diagnostics.HasError() {
		return
	}

	apiKey, err := r.client.Create(paramsApiKeyCreate, files_sdk.WithContext(ctx))
	if err != nil {
		resp.Diagnostics.AddError(
			"Error Creating Files ApiKey",
			"Could not create api_key, unexpected error: "+err.Error(),
		)
		return
	}

	diags = r.populateResourceModel(ctx, apiKey, &plan)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	diags = resp.State.Set(ctx, plan)
	resp.Diagnostics.Append(diags...)
}

func (r *apiKeyResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	var state apiKeyResourceModel
	diags := req.State.Get(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	paramsApiKeyFind := files_sdk.ApiKeyFindParams{}
	paramsApiKeyFind.Id = state.Id.ValueInt64()

	apiKey, err := r.client.Find(paramsApiKeyFind, files_sdk.WithContext(ctx))
	if err != nil {
		if files_sdk.IsNotExist(err) {
			resp.State.RemoveResource(ctx)
			return
		}

		resp.Diagnostics.AddError(
			"Error Reading Files ApiKey",
			"Could not read api_key id "+fmt.Sprint(state.Id.ValueInt64())+": "+err.Error(),
		)
		return
	}

	diags = r.populateResourceModel(ctx, apiKey, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	diags = resp.State.Set(ctx, state)
	resp.Diagnostics.Append(diags...)
}

func (r *apiKeyResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	var plan apiKeyResourceModel
	diags := req.Plan.Get(ctx, &plan)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	paramsApiKeyUpdate := files_sdk.ApiKeyUpdateParams{}
	paramsApiKeyUpdate.Id = plan.Id.ValueInt64()
	paramsApiKeyUpdate.Description = plan.Description.ValueString()
	if !plan.ExpiresAt.IsNull() && plan.ExpiresAt.ValueString() != "" {
		updateExpiresAt, err := time.Parse(time.RFC3339, plan.ExpiresAt.ValueString())
		if err != nil {
			resp.Diagnostics.AddAttributeError(
				path.Root("expires_at"),
				"Error Parsing expires_at Time",
				"Could not parse expires_at time: "+err.Error(),
			)
		} else {
			paramsApiKeyUpdate.ExpiresAt = &updateExpiresAt
		}
	}
	paramsApiKeyUpdate.PermissionSet = paramsApiKeyUpdate.PermissionSet.Enum()[plan.PermissionSet.ValueString()]
	paramsApiKeyUpdate.Name = plan.Name.ValueString()

	if resp.Diagnostics.HasError() {
		return
	}

	apiKey, err := r.client.Update(paramsApiKeyUpdate, files_sdk.WithContext(ctx))
	if err != nil {
		resp.Diagnostics.AddError(
			"Error Updating Files ApiKey",
			"Could not update api_key, unexpected error: "+err.Error(),
		)
		return
	}

	diags = r.populateResourceModel(ctx, apiKey, &plan)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	diags = resp.State.Set(ctx, plan)
	resp.Diagnostics.Append(diags...)
}

func (r *apiKeyResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	var state apiKeyResourceModel
	diags := req.State.Get(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	paramsApiKeyDelete := files_sdk.ApiKeyDeleteParams{}
	paramsApiKeyDelete.Id = state.Id.ValueInt64()

	err := r.client.Delete(paramsApiKeyDelete, files_sdk.WithContext(ctx))
	if err != nil && !files_sdk.IsNotExist(err) {
		resp.Diagnostics.AddError(
			"Error Deleting Files ApiKey",
			"Could not delete api_key id "+fmt.Sprint(state.Id.ValueInt64())+": "+err.Error(),
		)
	}
}

func (r *apiKeyResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
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

func (r *apiKeyResource) populateResourceModel(ctx context.Context, apiKey files_sdk.ApiKey, state *apiKeyResourceModel) (diags diag.Diagnostics) {
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
