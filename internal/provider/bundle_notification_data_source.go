package provider

import (
	"context"
	"fmt"

	files_sdk "github.com/Files-com/files-sdk-go/v3"
	bundle_notification "github.com/Files-com/files-sdk-go/v3/bundlenotification"
	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/datasource/schema"
	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

var (
	_ datasource.DataSource              = &bundleNotificationDataSource{}
	_ datasource.DataSourceWithConfigure = &bundleNotificationDataSource{}
)

func NewBundleNotificationDataSource() datasource.DataSource {
	return &bundleNotificationDataSource{}
}

type bundleNotificationDataSource struct {
	client *bundle_notification.Client
}

type bundleNotificationDataSourceModel struct {
	Id                   types.Int64 `tfsdk:"id"`
	BundleId             types.Int64 `tfsdk:"bundle_id"`
	NotifyOnRegistration types.Bool  `tfsdk:"notify_on_registration"`
	NotifyOnUpload       types.Bool  `tfsdk:"notify_on_upload"`
	UserId               types.Int64 `tfsdk:"user_id"`
}

func (r *bundleNotificationDataSource) Configure(_ context.Context, req datasource.ConfigureRequest, resp *datasource.ConfigureResponse) {
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

	r.client = &bundle_notification.Client{Config: sdk_config}
}

func (r *bundleNotificationDataSource) Metadata(_ context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_bundle_notification"
}

func (r *bundleNotificationDataSource) Schema(_ context.Context, _ datasource.SchemaRequest, resp *datasource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Description: "A BundleNotification is an E-Mail sent out to users when certain actions are performed on or within a shared set of files and folders.",
		Attributes: map[string]schema.Attribute{
			"id": schema.Int64Attribute{
				Description: "Bundle Notification ID",
				Required:    true,
			},
			"bundle_id": schema.Int64Attribute{
				Description: "Bundle ID to notify on",
				Computed:    true,
			},
			"notify_on_registration": schema.BoolAttribute{
				Description: "Triggers bundle notification when a registration action occurs for it.",
				Computed:    true,
			},
			"notify_on_upload": schema.BoolAttribute{
				Description: "Triggers bundle notification when a upload action occurs for it.",
				Computed:    true,
			},
			"user_id": schema.Int64Attribute{
				Description: "The id of the user to notify.",
				Computed:    true,
			},
		},
	}
}

func (r *bundleNotificationDataSource) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
	var data bundleNotificationDataSourceModel
	diags := req.Config.Get(ctx, &data)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	paramsBundleNotificationFind := files_sdk.BundleNotificationFindParams{}
	paramsBundleNotificationFind.Id = data.Id.ValueInt64()

	bundleNotification, err := r.client.Find(paramsBundleNotificationFind, files_sdk.WithContext(ctx))
	if err != nil {
		resp.Diagnostics.AddError(
			"Error Reading Files BundleNotification",
			"Could not read bundle_notification id "+fmt.Sprint(data.Id.ValueInt64())+": "+err.Error(),
		)
		return
	}

	diags = r.populateDataSourceModel(ctx, bundleNotification, &data)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	diags = resp.State.Set(ctx, data)
	resp.Diagnostics.Append(diags...)
}

func (r *bundleNotificationDataSource) populateDataSourceModel(ctx context.Context, bundleNotification files_sdk.BundleNotification, state *bundleNotificationDataSourceModel) (diags diag.Diagnostics) {
	state.BundleId = types.Int64Value(bundleNotification.BundleId)
	state.Id = types.Int64Value(bundleNotification.Id)
	state.NotifyOnRegistration = types.BoolPointerValue(bundleNotification.NotifyOnRegistration)
	state.NotifyOnUpload = types.BoolPointerValue(bundleNotification.NotifyOnUpload)
	state.UserId = types.Int64Value(bundleNotification.UserId)

	return
}
