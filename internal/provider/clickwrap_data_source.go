package provider

import (
	"context"
	"fmt"

	files_sdk "github.com/Files-com/files-sdk-go/v3"
	clickwrap "github.com/Files-com/files-sdk-go/v3/clickwrap"
	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/datasource/schema"
	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

var (
	_ datasource.DataSource              = &clickwrapDataSource{}
	_ datasource.DataSourceWithConfigure = &clickwrapDataSource{}
)

func NewClickwrapDataSource() datasource.DataSource {
	return &clickwrapDataSource{}
}

type clickwrapDataSource struct {
	client *clickwrap.Client
}

type clickwrapDataSourceModel struct {
	Id             types.Int64  `tfsdk:"id"`
	Name           types.String `tfsdk:"name"`
	Body           types.String `tfsdk:"body"`
	UseWithUsers   types.String `tfsdk:"use_with_users"`
	UseWithBundles types.String `tfsdk:"use_with_bundles"`
	UseWithInboxes types.String `tfsdk:"use_with_inboxes"`
}

func (r *clickwrapDataSource) Configure(_ context.Context, req datasource.ConfigureRequest, resp *datasource.ConfigureResponse) {
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

	r.client = &clickwrap.Client{Config: sdk_config}
}

func (r *clickwrapDataSource) Metadata(_ context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_clickwrap"
}

func (r *clickwrapDataSource) Schema(_ context.Context, _ datasource.SchemaRequest, resp *datasource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Description: "A Clickwrap is a legal agreement (such as an NDA or Terms of Use) that your Users and/or Bundle/Inbox participants will need to agree to via a \"Clickwrap\" UI before accessing the site, bundle, or inbox.\n\n\n\nThe values for `use_with_users`, `use_with_bundles`, `use_with_inboxes` are explained as follows:\n\n\n\n* `none` - This Clickwrap may not be used in this context.\n\n* `available_to_all_users` - This Clickwrap may be assigned in this context by any user.\n\n* `available` - This Clickwrap may be assigned in this context, but only by Site Admins. We recognize that the name of this setting is somewhat ambiguous, but we maintain it for legacy reasons.\n\n* `required` - This Clickwrap will always be used in this context, and may not be overridden.",
		Attributes: map[string]schema.Attribute{
			"id": schema.Int64Attribute{
				Description: "Clickwrap ID",
				Required:    true,
			},
			"name": schema.StringAttribute{
				Description: "Name of the Clickwrap agreement (used when selecting from multiple Clickwrap agreements.)",
				Computed:    true,
			},
			"body": schema.StringAttribute{
				Description: "Body text of Clickwrap (supports Markdown formatting).",
				Computed:    true,
			},
			"use_with_users": schema.StringAttribute{
				Description: "Use this Clickwrap for User Registrations?  Note: This only applies to User Registrations where the User is invited to your Files.com site using an E-Mail invitation process where they then set their own password.",
				Computed:    true,
			},
			"use_with_bundles": schema.StringAttribute{
				Description: "Use this Clickwrap for Bundles?",
				Computed:    true,
			},
			"use_with_inboxes": schema.StringAttribute{
				Description: "Use this Clickwrap for Inboxes?",
				Computed:    true,
			},
		},
	}
}

func (r *clickwrapDataSource) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
	var data clickwrapDataSourceModel
	diags := req.Config.Get(ctx, &data)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	paramsClickwrapFind := files_sdk.ClickwrapFindParams{}
	paramsClickwrapFind.Id = data.Id.ValueInt64()

	clickwrap, err := r.client.Find(paramsClickwrapFind, files_sdk.WithContext(ctx))
	if err != nil {
		resp.Diagnostics.AddError(
			"Error Reading Files Clickwrap",
			"Could not read clickwrap id "+fmt.Sprint(data.Id.ValueInt64())+": "+err.Error(),
		)
		return
	}

	diags = r.populateDataSourceModel(ctx, clickwrap, &data)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	diags = resp.State.Set(ctx, data)
	resp.Diagnostics.Append(diags...)
}

func (r *clickwrapDataSource) populateDataSourceModel(ctx context.Context, clickwrap files_sdk.Clickwrap, state *clickwrapDataSourceModel) (diags diag.Diagnostics) {
	state.Id = types.Int64Value(clickwrap.Id)
	state.Name = types.StringValue(clickwrap.Name)
	state.Body = types.StringValue(clickwrap.Body)
	state.UseWithUsers = types.StringValue(clickwrap.UseWithUsers)
	state.UseWithBundles = types.StringValue(clickwrap.UseWithBundles)
	state.UseWithInboxes = types.StringValue(clickwrap.UseWithInboxes)

	return
}
