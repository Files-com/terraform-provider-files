package provider

import (
	"context"
	"fmt"

	files_sdk "github.com/Files-com/files-sdk-go/v3"
	share_group "github.com/Files-com/files-sdk-go/v3/sharegroup"
	"github.com/Files-com/terraform-provider-files/lib"
	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/datasource/schema"
	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

var (
	_ datasource.DataSource              = &shareGroupDataSource{}
	_ datasource.DataSourceWithConfigure = &shareGroupDataSource{}
)

func NewShareGroupDataSource() datasource.DataSource {
	return &shareGroupDataSource{}
}

type shareGroupDataSource struct {
	client *share_group.Client
}

type shareGroupDataSourceModel struct {
	Id      types.Int64   `tfsdk:"id"`
	Name    types.String  `tfsdk:"name"`
	Notes   types.String  `tfsdk:"notes"`
	UserId  types.Int64   `tfsdk:"user_id"`
	Members types.Dynamic `tfsdk:"members"`
}

func (r *shareGroupDataSource) Configure(_ context.Context, req datasource.ConfigureRequest, resp *datasource.ConfigureResponse) {
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

	r.client = &share_group.Client{Config: sdk_config}
}

func (r *shareGroupDataSource) Metadata(_ context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_share_group"
}

func (r *shareGroupDataSource) Schema(_ context.Context, _ datasource.SchemaRequest, resp *datasource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Description: "A ShareGroup is a way for you to store and name groups of email contacts to be used for sending share and inbox invitations.",
		Attributes: map[string]schema.Attribute{
			"id": schema.Int64Attribute{
				Description: "Share Group ID",
				Required:    true,
			},
			"name": schema.StringAttribute{
				Description: "Name of the share group",
				Computed:    true,
			},
			"notes": schema.StringAttribute{
				Description: "Additional notes of the share group",
				Computed:    true,
			},
			"user_id": schema.Int64Attribute{
				Description: "Owner User ID",
				Computed:    true,
			},
			"members": schema.DynamicAttribute{
				Description: "A list of share group members",
				Computed:    true,
			},
		},
	}
}

func (r *shareGroupDataSource) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
	var data shareGroupDataSourceModel
	diags := req.Config.Get(ctx, &data)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	paramsShareGroupFind := files_sdk.ShareGroupFindParams{}
	paramsShareGroupFind.Id = data.Id.ValueInt64()

	shareGroup, err := r.client.Find(paramsShareGroupFind, files_sdk.WithContext(ctx))
	if err != nil {
		resp.Diagnostics.AddError(
			"Error Reading Files ShareGroup",
			"Could not read share_group id "+fmt.Sprint(data.Id.ValueInt64())+": "+err.Error(),
		)
		return
	}

	diags = r.populateDataSourceModel(ctx, shareGroup, &data)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	diags = resp.State.Set(ctx, data)
	resp.Diagnostics.Append(diags...)
}

func (r *shareGroupDataSource) populateDataSourceModel(ctx context.Context, shareGroup files_sdk.ShareGroup, state *shareGroupDataSourceModel) (diags diag.Diagnostics) {
	var propDiags diag.Diagnostics

	state.Id = types.Int64Value(shareGroup.Id)
	state.Name = types.StringValue(shareGroup.Name)
	state.Notes = types.StringValue(shareGroup.Notes)
	state.UserId = types.Int64Value(shareGroup.UserId)
	state.Members, propDiags = lib.ToDynamic(ctx, path.Root("members"), shareGroup.Members, state.Members.UnderlyingValue())
	diags.Append(propDiags...)

	return
}
