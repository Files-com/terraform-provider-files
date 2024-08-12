package provider

import (
	"context"
	"fmt"

	files_sdk "github.com/Files-com/files-sdk-go/v3"
	behavior "github.com/Files-com/files-sdk-go/v3/behavior"
	"github.com/Files-com/terraform-provider-files/lib"
	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/datasource/schema"
	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

var (
	_ datasource.DataSource              = &behaviorDataSource{}
	_ datasource.DataSourceWithConfigure = &behaviorDataSource{}
)

func NewBehaviorDataSource() datasource.DataSource {
	return &behaviorDataSource{}
}

type behaviorDataSource struct {
	client *behavior.Client
}

type behaviorDataSourceModel struct {
	Id                          types.Int64   `tfsdk:"id"`
	Path                        types.String  `tfsdk:"path"`
	AttachmentUrl               types.String  `tfsdk:"attachment_url"`
	Behavior                    types.String  `tfsdk:"behavior"`
	Name                        types.String  `tfsdk:"name"`
	Description                 types.String  `tfsdk:"description"`
	Value                       types.Dynamic `tfsdk:"value"`
	DisableParentFolderBehavior types.Bool    `tfsdk:"disable_parent_folder_behavior"`
	Recursive                   types.Bool    `tfsdk:"recursive"`
}

func (r *behaviorDataSource) Configure(_ context.Context, req datasource.ConfigureRequest, resp *datasource.ConfigureResponse) {
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

	r.client = &behavior.Client{Config: sdk_config}
}

func (r *behaviorDataSource) Metadata(_ context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_behavior"
}

func (r *behaviorDataSource) Schema(_ context.Context, _ datasource.SchemaRequest, resp *datasource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Description: "A Behavior is an API resource for what are also known as Folder Settings. Every behavior is associated with a folder.\n\n\n\nDepending on the behavior, it may also operate on child folders. It may be overridable at the child folder level or maybe can be added to at the child folder level. The exact options for each behavior type are explained in the table below.\n\n\n\nAdditionally, some behaviors are visible to non-admins, and others are even settable by non-admins. All the details are below.\n\n\n\nEach behavior uses a different format for storing its settings value. Next to each behavior type is an example value. Our API and SDKs currently require that the value for behaviors be sent as raw JSON within the `value` field. Our SDK generator and API documentation generator doesn't fully keep up with this requirement, so if you need any help finding the exact syntax to use for your language or use case, just reach out.\n\n\n\nNote: Append Timestamp behavior removed. Check [Override Upload Filename](#override-upload-filename-behaviors) behavior which have even more functionality to modify name on upload.",
		Attributes: map[string]schema.Attribute{
			"id": schema.Int64Attribute{
				Description: "Folder behavior ID",
				Required:    true,
			},
			"path": schema.StringAttribute{
				Description: "Folder path.  Note that Behavior paths cannot be updated once initially set.  You will need to remove and re-create the behavior on the new path. This must be slash-delimited, but it must neither start nor end with a slash. Maximum of 5000 characters.",
				Computed:    true,
			},
			"attachment_url": schema.StringAttribute{
				Description: "URL for attached file",
				Computed:    true,
			},
			"behavior": schema.StringAttribute{
				Description: "Behavior type.",
				Computed:    true,
			},
			"name": schema.StringAttribute{
				Description: "Name for this behavior.",
				Computed:    true,
			},
			"description": schema.StringAttribute{
				Description: "Description for this behavior.",
				Computed:    true,
			},
			"value": schema.DynamicAttribute{
				Description: "Settings for this behavior.  See the section above for an example value to provide here.  Formatting is different for each Behavior type.  May be sent as nested JSON or a single JSON-encoded string.  If using XML encoding for the API call, this data must be sent as a JSON-encoded string.",
				Computed:    true,
			},
			"disable_parent_folder_behavior": schema.BoolAttribute{
				Description: "If true, the parent folder's behavior will be disabled for this folder and its children.",
				Computed:    true,
			},
			"recursive": schema.BoolAttribute{
				Description: "Is behavior recursive?",
				Computed:    true,
			},
		},
	}
}

func (r *behaviorDataSource) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
	var data behaviorDataSourceModel
	diags := req.Config.Get(ctx, &data)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	paramsBehaviorFind := files_sdk.BehaviorFindParams{}
	paramsBehaviorFind.Id = data.Id.ValueInt64()

	behavior, err := r.client.Find(paramsBehaviorFind, files_sdk.WithContext(ctx))
	if err != nil {
		resp.Diagnostics.AddError(
			"Error Reading Files Behavior",
			"Could not read behavior id "+fmt.Sprint(data.Id.ValueInt64())+": "+err.Error(),
		)
		return
	}

	diags = r.populateDataSourceModel(ctx, behavior, &data)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	diags = resp.State.Set(ctx, data)
	resp.Diagnostics.Append(diags...)
}

func (r *behaviorDataSource) populateDataSourceModel(ctx context.Context, behavior files_sdk.Behavior, state *behaviorDataSourceModel) (diags diag.Diagnostics) {
	var propDiags diag.Diagnostics

	state.Id = types.Int64Value(behavior.Id)
	state.Path = types.StringValue(behavior.Path)
	state.AttachmentUrl = types.StringValue(behavior.AttachmentUrl)
	state.Behavior = types.StringValue(behavior.Behavior)
	state.Name = types.StringValue(behavior.Name)
	state.Description = types.StringValue(behavior.Description)
	state.Value, propDiags = lib.ToDynamic(ctx, path.Root("value"), behavior.Value, state.Value.UnderlyingValue())
	diags.Append(propDiags...)
	state.DisableParentFolderBehavior = types.BoolPointerValue(behavior.DisableParentFolderBehavior)
	state.Recursive = types.BoolPointerValue(behavior.Recursive)

	return
}
