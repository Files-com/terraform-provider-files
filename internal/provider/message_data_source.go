package provider

import (
	"context"
	"fmt"

	files_sdk "github.com/Files-com/files-sdk-go/v3"
	message "github.com/Files-com/files-sdk-go/v3/message"
	"github.com/Files-com/terraform-provider-files/lib"
	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/datasource/schema"
	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

var (
	_ datasource.DataSource              = &messageDataSource{}
	_ datasource.DataSourceWithConfigure = &messageDataSource{}
)

func NewMessageDataSource() datasource.DataSource {
	return &messageDataSource{}
}

type messageDataSource struct {
	client *message.Client
}

type messageDataSourceModel struct {
	Id       types.Int64   `tfsdk:"id"`
	Subject  types.String  `tfsdk:"subject"`
	Body     types.String  `tfsdk:"body"`
	Comments types.Dynamic `tfsdk:"comments"`
}

func (r *messageDataSource) Configure(_ context.Context, req datasource.ConfigureRequest, resp *datasource.ConfigureResponse) {
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

	r.client = &message.Client{Config: sdk_config}
}

func (r *messageDataSource) Metadata(_ context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_message"
}

func (r *messageDataSource) Schema(_ context.Context, _ datasource.SchemaRequest, resp *datasource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Description: "A Messages is a part of Files.com's project management features and represent a message posted by a user to a project.",
		Attributes: map[string]schema.Attribute{
			"id": schema.Int64Attribute{
				Description: "Message ID",
				Required:    true,
			},
			"subject": schema.StringAttribute{
				Description: "Message subject.",
				Computed:    true,
			},
			"body": schema.StringAttribute{
				Description: "Message body.",
				Computed:    true,
			},
			"comments": schema.DynamicAttribute{
				Description: "Comments.",
				Computed:    true,
			},
		},
	}
}

func (r *messageDataSource) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
	var data messageDataSourceModel
	diags := req.Config.Get(ctx, &data)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	paramsMessageFind := files_sdk.MessageFindParams{}
	paramsMessageFind.Id = data.Id.ValueInt64()

	message, err := r.client.Find(paramsMessageFind, files_sdk.WithContext(ctx))
	if err != nil {
		resp.Diagnostics.AddError(
			"Error Reading Files Message",
			"Could not read message id "+fmt.Sprint(data.Id.ValueInt64())+": "+err.Error(),
		)
		return
	}

	diags = r.populateDataSourceModel(ctx, message, &data)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	diags = resp.State.Set(ctx, data)
	resp.Diagnostics.Append(diags...)
}

func (r *messageDataSource) populateDataSourceModel(ctx context.Context, message files_sdk.Message, state *messageDataSourceModel) (diags diag.Diagnostics) {
	var propDiags diag.Diagnostics

	state.Id = types.Int64Value(message.Id)
	state.Subject = types.StringValue(message.Subject)
	state.Body = types.StringValue(message.Body)
	state.Comments, propDiags = lib.ToDynamic(ctx, path.Root("comments"), message.Comments, state.Comments.UnderlyingValue())
	diags.Append(propDiags...)

	return
}
