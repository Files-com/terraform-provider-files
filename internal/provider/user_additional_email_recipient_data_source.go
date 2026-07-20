package provider

import (
	"context"
	"fmt"

	files_sdk "github.com/Files-com/files-sdk-go/v3"
	user_additional_email_recipient "github.com/Files-com/files-sdk-go/v3/useradditionalemailrecipient"
	"github.com/Files-com/terraform-provider-files/lib"
	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/datasource/schema"
	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

var (
	_ datasource.DataSource              = &userAdditionalEmailRecipientDataSource{}
	_ datasource.DataSourceWithConfigure = &userAdditionalEmailRecipientDataSource{}
)

func NewUserAdditionalEmailRecipientDataSource() datasource.DataSource {
	return &userAdditionalEmailRecipientDataSource{}
}

type userAdditionalEmailRecipientDataSource struct {
	client *user_additional_email_recipient.Client
}

type userAdditionalEmailRecipientDataSourceModel struct {
	Id          types.Int64  `tfsdk:"id"`
	UserId      types.Int64  `tfsdk:"user_id"`
	WorkspaceId types.Int64  `tfsdk:"workspace_id"`
	Email       types.String `tfsdk:"email"`
	CreatedAt   types.String `tfsdk:"created_at"`
}

func (r *userAdditionalEmailRecipientDataSource) Configure(_ context.Context, req datasource.ConfigureRequest, resp *datasource.ConfigureResponse) {
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

	r.client = &user_additional_email_recipient.Client{Config: sdk_config}
}

func (r *userAdditionalEmailRecipientDataSource) Metadata(_ context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_user_additional_email_recipient"
}

func (r *userAdditionalEmailRecipientDataSource) Schema(_ context.Context, _ datasource.SchemaRequest, resp *datasource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Description: "",
		Attributes: map[string]schema.Attribute{
			"id": schema.Int64Attribute{
				Description: "User additional email recipient ID",
				Required:    true,
			},
			"user_id": schema.Int64Attribute{
				Description: "User ID",
				Computed:    true,
			},
			"workspace_id": schema.Int64Attribute{
				Description: "Workspace ID (0 for default workspace).",
				Computed:    true,
			},
			"email": schema.StringAttribute{
				Description: "Additional email recipient address",
				Computed:    true,
			},
			"created_at": schema.StringAttribute{
				Description: "Created at date/time",
				Computed:    true,
			},
		},
	}
}

func (r *userAdditionalEmailRecipientDataSource) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
	var data userAdditionalEmailRecipientDataSourceModel
	diags := req.Config.Get(ctx, &data)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	paramsUserAdditionalEmailRecipientFind := files_sdk.UserAdditionalEmailRecipientFindParams{}
	paramsUserAdditionalEmailRecipientFind.Id = data.Id.ValueInt64()

	userAdditionalEmailRecipient, err := r.client.Find(paramsUserAdditionalEmailRecipientFind, files_sdk.WithContext(ctx))
	if err != nil {
		resp.Diagnostics.AddError(
			"Error Reading Files UserAdditionalEmailRecipient",
			"Could not read user_additional_email_recipient id "+fmt.Sprint(data.Id.ValueInt64())+": "+err.Error(),
		)
		return
	}

	diags = r.populateDataSourceModel(ctx, userAdditionalEmailRecipient, &data)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	diags = resp.State.Set(ctx, data)
	resp.Diagnostics.Append(diags...)
}

func (r *userAdditionalEmailRecipientDataSource) populateDataSourceModel(ctx context.Context, userAdditionalEmailRecipient files_sdk.UserAdditionalEmailRecipient, state *userAdditionalEmailRecipientDataSourceModel) (diags diag.Diagnostics) {
	state.Id = types.Int64Value(userAdditionalEmailRecipient.Id)
	state.UserId = types.Int64Value(userAdditionalEmailRecipient.UserId)
	state.WorkspaceId = types.Int64Value(userAdditionalEmailRecipient.WorkspaceId)
	state.Email = types.StringValue(userAdditionalEmailRecipient.Email)
	if err := lib.TimeToStringType(ctx, path.Root("created_at"), userAdditionalEmailRecipient.CreatedAt, &state.CreatedAt); err != nil {
		diags.AddError(
			"Error Creating Files UserAdditionalEmailRecipient",
			"Could not convert state created_at to string: "+err.Error(),
		)
	}

	return
}
