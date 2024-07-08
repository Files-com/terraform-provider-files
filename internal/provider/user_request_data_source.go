package provider

import (
	"context"
	"fmt"

	files_sdk "github.com/Files-com/files-sdk-go/v3"
	user_request "github.com/Files-com/files-sdk-go/v3/userrequest"
	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/datasource/schema"
	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

var (
	_ datasource.DataSource              = &userRequestDataSource{}
	_ datasource.DataSourceWithConfigure = &userRequestDataSource{}
)

func NewUserRequestDataSource() datasource.DataSource {
	return &userRequestDataSource{}
}

type userRequestDataSource struct {
	client *user_request.Client
}

type userRequestDataSourceModel struct {
	Id      types.Int64  `tfsdk:"id"`
	Name    types.String `tfsdk:"name"`
	Email   types.String `tfsdk:"email"`
	Details types.String `tfsdk:"details"`
	Company types.String `tfsdk:"company"`
}

func (r *userRequestDataSource) Configure(_ context.Context, req datasource.ConfigureRequest, resp *datasource.ConfigureResponse) {
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

	r.client = &user_request.Client{Config: sdk_config}
}

func (r *userRequestDataSource) Metadata(_ context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_user_request"
}

func (r *userRequestDataSource) Schema(_ context.Context, _ datasource.SchemaRequest, resp *datasource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Description: "User Requests allow anonymous users to place a request for access on the login screen to the site administrator.",
		Attributes: map[string]schema.Attribute{
			"id": schema.Int64Attribute{
				Description: "ID",
				Required:    true,
			},
			"name": schema.StringAttribute{
				Description: "User's full name",
				Computed:    true,
			},
			"email": schema.StringAttribute{
				Description: "User email address",
				Computed:    true,
			},
			"details": schema.StringAttribute{
				Description: "Details of the user's request",
				Computed:    true,
			},
			"company": schema.StringAttribute{
				Description: "User's company name",
				Computed:    true,
			},
		},
	}
}

func (r *userRequestDataSource) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
	var data userRequestDataSourceModel
	diags := req.Config.Get(ctx, &data)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	paramsUserRequestFind := files_sdk.UserRequestFindParams{}
	paramsUserRequestFind.Id = data.Id.ValueInt64()

	userRequest, err := r.client.Find(paramsUserRequestFind, files_sdk.WithContext(ctx))
	if err != nil {
		resp.Diagnostics.AddError(
			"Error Reading Files UserRequest",
			"Could not read user_request id "+fmt.Sprint(data.Id.ValueInt64())+": "+err.Error(),
		)
		return
	}

	diags = r.populateDataSourceModel(ctx, userRequest, &data)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	diags = resp.State.Set(ctx, data)
	resp.Diagnostics.Append(diags...)
}

func (r *userRequestDataSource) populateDataSourceModel(ctx context.Context, userRequest files_sdk.UserRequest, state *userRequestDataSourceModel) (diags diag.Diagnostics) {
	state.Id = types.Int64Value(userRequest.Id)
	state.Name = types.StringValue(userRequest.Name)
	state.Email = types.StringValue(userRequest.Email)
	state.Details = types.StringValue(userRequest.Details)
	state.Company = types.StringValue(userRequest.Company)

	return
}
