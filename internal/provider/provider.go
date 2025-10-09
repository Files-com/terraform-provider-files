package provider

import (
	"context"
	"os"
	"strings"

	files_sdk "github.com/Files-com/files-sdk-go/v3"
	"github.com/Files-com/terraform-provider-files/lib"

	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/provider"
	"github.com/hashicorp/terraform-plugin-framework/provider/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-log/tflog"
)

var (
	_ provider.Provider = &filesProvider{}
)

func New(version string) func() provider.Provider {
	return func() provider.Provider {
		return &filesProvider{
			version: version,
		}
	}
}

type filesProvider struct {
	version string
}

type filesProviderModel struct {
	APIKey           types.String `tfsdk:"api_key"`
	EndpointOverride types.String `tfsdk:"endpoint_override"`
	Environment      types.String `tfsdk:"environment"`
	FeatureFlags     types.List   `tfsdk:"feature_flags"`
}

func (p *filesProvider) Metadata(_ context.Context, _ provider.MetadataRequest, resp *provider.MetadataResponse) {
	resp.TypeName = "files"
	resp.Version = p.version
}

func (p *filesProvider) Schema(_ context.Context, _ provider.SchemaRequest, resp *provider.SchemaResponse) {
	resp.Schema = schema.Schema{
		Attributes: map[string]schema.Attribute{
			"api_key": schema.StringAttribute{
				Description: "The API key used to authenticate with Files.com. It can also be sourced from the `FILES_API_KEY` environment variable.",
				Optional:    true,
				Sensitive:   true,
			},
			"endpoint_override": schema.StringAttribute{
				Description: "Required if your site is configured to disable global acceleration. This can also be set to use a mock server in development or CI.",
				Optional:    true,
			},
			"environment": schema.StringAttribute{
				Optional: true,
			},
			"feature_flags": schema.ListAttribute{
				ElementType: types.StringType,
				Optional:    true,
			},
		},
	}
}

func (p *filesProvider) Configure(ctx context.Context, req provider.ConfigureRequest, resp *provider.ConfigureResponse) {
	var config filesProviderModel
	diags := req.Config.Get(ctx, &config)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	if config.APIKey.IsUnknown() {
		resp.Diagnostics.AddAttributeError(
			path.Root("api_key"),
			"Unknown Files API Key",
			"The provider cannot create the Files API client as there is an unknown configuration value for the Files API key. "+
				"Either target apply the source of the value first, set the value statically in the configuration, or use the FILES_API_KEY environment variable.",
		)
	}

	if config.EndpointOverride.IsUnknown() {
		resp.Diagnostics.AddAttributeError(
			path.Root("endpoint_override"),
			"Unknown Files API Endpoint Override",
			"The provider cannot create the Files API client as there is an unknown configuration value for the Files API endpoint override. "+
				"Either target apply the source of the value first or set the value statically in the configuration.",
		)
	}

	if config.Environment.IsUnknown() {
		resp.Diagnostics.AddAttributeError(
			path.Root("environment"),
			"Unknown Files API Environment",
			"The provider cannot create the Files API client as there is an unknown configuration value for the Files API environment. "+
				"Either target apply the source of the value first or set the value statically in the configuration.",
		)
	}

	if config.FeatureFlags.IsUnknown() {
		resp.Diagnostics.AddAttributeError(
			path.Root("feature_flags"),
			"Unknown Files API Feature Flags",
			"The provider cannot create the Files API client as there is an unknown configuration value for the Files API feature flags. "+
				"Either target apply the source of the value first or set the value statically in the configuration.",
		)
	}

	if resp.Diagnostics.HasError() {
		return
	}

	apiKey := os.Getenv("FILES_API_KEY")
	endpointOverride := ""
	environment := ""
	featureFlags := []string{}

	if !config.APIKey.IsNull() {
		tflog.Info(ctx, "Using API key from configuration")
		apiKey = config.APIKey.ValueString()
	}
	if !config.EndpointOverride.IsNull() {
		tflog.Info(ctx, "Using endpoint override from configuration")
		endpointOverride = config.EndpointOverride.ValueString()
	}
	if !config.Environment.IsNull() {
		tflog.Info(ctx, "Using environment from configuration")
		environment = config.Environment.ValueString()
	}
	if !config.FeatureFlags.IsNull() {
		tflog.Info(ctx, "Using feature flags from configuration")
		diags = config.FeatureFlags.ElementsAs(ctx, &featureFlags, false)
		resp.Diagnostics.Append(diags...)
	}

	if apiKey == "" {
		resp.Diagnostics.AddAttributeError(
			path.Root("api_key"),
			"Missing Files API Key",
			"The provider cannot create the Files API client as there is a missing or empty value for the Files API key. "+
				"Set the API key value in the configuration or use the FILES_API_KEY environment variable. "+
				"If either is already set, ensure the value is not empty.",
		)
	}

	if resp.Diagnostics.HasError() {
		return
	}

	sdkConfig := files_sdk.Config{
		APIKey:           apiKey,
		EndpointOverride: endpointOverride,
		Environment:      files_sdk.NewEnvironment(environment),
		Logger:           lib.Logger{Ctx: ctx},
	}
	if len(featureFlags) > 0 {
		sdkConfig.FeatureFlags = map[string]bool{}
		for _, flag := range featureFlags {
			sdkConfig.FeatureFlags[flag] = true
		}
	}
	sdkConfig = sdkConfig.Init()
	sdkConfig.Client.Logger = sdkConfig.Logger
	sdkConfig.UserAgent = "Files.com Terraform " + strings.TrimSpace(p.version) // Set this after Init() to avoid overwriting.

	resp.DataSourceData = sdkConfig
	resp.ResourceData = sdkConfig
}

func (p *filesProvider) DataSources(_ context.Context) []func() datasource.DataSource {
	return []func() datasource.DataSource{
		NewActionNotificationExportDataSource,
		NewApiKeyDataSource,
		NewAs2PartnerDataSource,
		NewAs2StationDataSource,
		NewAutomationDataSource,
		NewAutomationRunDataSource,
		NewBehaviorDataSource,
		NewBundleDataSource,
		NewBundleNotificationDataSource,
		NewChildSiteManagementPolicyDataSource,
		NewClickwrapDataSource,
		NewExternalEventDataSource,
		NewFileDataSource,
		NewFileCommentDataSource,
		NewFileMigrationDataSource,
		NewFolderDataSource,
		NewFormFieldSetDataSource,
		NewGpgKeyDataSource,
		NewGroupDataSource,
		NewGroupUserDataSource,
		NewHistoryExportDataSource,
		NewInvoiceDataSource,
		NewLockDataSource,
		NewMessageDataSource,
		NewMessageCommentDataSource,
		NewMessageCommentReactionDataSource,
		NewMessageReactionDataSource,
		NewNotificationDataSource,
		NewPartnerDataSource,
		NewPaymentDataSource,
		NewPermissionDataSource,
		NewPriorityDataSource,
		NewProjectDataSource,
		NewPublicKeyDataSource,
		NewRemoteMountBackendDataSource,
		NewRemoteServerDataSource,
		NewRequestDataSource,
		NewSftpHostKeyDataSource,
		NewShareGroupDataSource,
		NewSiemHttpDestinationDataSource,
		NewSiteDataSource,
		NewSnapshotDataSource,
		NewSsoStrategyDataSource,
		NewStyleDataSource,
		NewSyncDataSource,
		NewSyncRunDataSource,
		NewUserDataSource,
		NewUserLifecycleRuleDataSource,
		NewUserRequestDataSource,
	}
}

func (p *filesProvider) Resources(_ context.Context) []func() resource.Resource {
	return []func() resource.Resource{
		NewApiKeyResource,
		NewAs2PartnerResource,
		NewAs2StationResource,
		NewAutomationResource,
		NewBehaviorResource,
		NewBundleResource,
		NewBundleNotificationResource,
		NewChildSiteManagementPolicyResource,
		NewClickwrapResource,
		NewFileResource,
		NewFileCommentResource,
		NewFolderResource,
		NewFormFieldSetResource,
		NewGpgKeyResource,
		NewGroupResource,
		NewGroupUserResource,
		NewLockResource,
		NewMessageResource,
		NewMessageCommentResource,
		NewMessageCommentReactionResource,
		NewMessageReactionResource,
		NewNotificationResource,
		NewPartnerResource,
		NewPermissionResource,
		NewProjectResource,
		NewPublicKeyResource,
		NewRemoteMountBackendResource,
		NewRemoteServerResource,
		NewRequestResource,
		NewSftpHostKeyResource,
		NewShareGroupResource,
		NewSiemHttpDestinationResource,
		NewSiteResource,
		NewSnapshotResource,
		NewSyncResource,
		NewUserResource,
		NewUserLifecycleRuleResource,
		NewUserRequestResource,
	}
}
