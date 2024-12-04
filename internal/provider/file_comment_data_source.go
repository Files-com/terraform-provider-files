package provider

import (
	"context"
	"fmt"

	files_sdk "github.com/Files-com/files-sdk-go/v3"
	file_comment "github.com/Files-com/files-sdk-go/v3/filecomment"
	"github.com/Files-com/terraform-provider-files/lib"
	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/datasource/schema"
	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

var (
	_ datasource.DataSource              = &fileCommentDataSource{}
	_ datasource.DataSourceWithConfigure = &fileCommentDataSource{}
)

func NewFileCommentDataSource() datasource.DataSource {
	return &fileCommentDataSource{}
}

type fileCommentDataSource struct {
	client *file_comment.Client
}

type fileCommentDataSourceModel struct {
	Id        types.Int64   `tfsdk:"id"`
	Path      types.String  `tfsdk:"path"`
	Body      types.String  `tfsdk:"body"`
	Reactions types.Dynamic `tfsdk:"reactions"`
}

func (r *fileCommentDataSource) Configure(_ context.Context, req datasource.ConfigureRequest, resp *datasource.ConfigureResponse) {
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

	r.client = &file_comment.Client{Config: sdk_config}
}

func (r *fileCommentDataSource) Metadata(_ context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_file_comment"
}

func (r *fileCommentDataSource) Schema(_ context.Context, _ datasource.SchemaRequest, resp *datasource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Description: "A FileComment is a comment attached to a file by a user.",
		Attributes: map[string]schema.Attribute{
			"id": schema.Int64Attribute{
				Description: "File Comment ID",
				Required:    true,
			},
			"path": schema.StringAttribute{
				Description: "File path.",
				Required:    true,
			},
			"body": schema.StringAttribute{
				Description: "Comment body.",
				Computed:    true,
			},
			"reactions": schema.DynamicAttribute{
				Description: "Reactions to this comment.",
				Computed:    true,
			},
		},
	}
}

func (r *fileCommentDataSource) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
	var data fileCommentDataSourceModel
	diags := req.Config.Get(ctx, &data)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	paramsFileCommentListFor := files_sdk.FileCommentListForParams{}
	paramsFileCommentListFor.Path = data.Path.ValueString()

	fileCommentIt, err := r.client.ListFor(paramsFileCommentListFor, files_sdk.WithContext(ctx))
	if err != nil {
		resp.Diagnostics.AddError(
			"Error Reading Files FileComment",
			"Could not read file_comment id "+fmt.Sprint(data.Id.ValueInt64())+": "+err.Error(),
		)
		return
	}

	var fileComment *files_sdk.FileComment
	for fileCommentIt.Next() {
		entry := fileCommentIt.FileComment()
		if entry.Id == data.Id.ValueInt64() {
			fileComment = &entry
			break
		}
	}

	if err = fileCommentIt.Err(); err != nil {
		resp.Diagnostics.AddError(
			"Error Reading Files FileComment",
			"Could not read file_comment id "+fmt.Sprint(data.Id.ValueInt64())+": "+err.Error(),
		)
	}

	if fileComment == nil {
		resp.Diagnostics.AddError(
			"Error Reading Files FileComment",
			"Could not find file_comment id "+fmt.Sprint(data.Id.ValueInt64())+"",
		)
		return
	}

	diags = r.populateDataSourceModel(ctx, *fileComment, &data)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	diags = resp.State.Set(ctx, data)
	resp.Diagnostics.Append(diags...)
}

func (r *fileCommentDataSource) populateDataSourceModel(ctx context.Context, fileComment files_sdk.FileComment, state *fileCommentDataSourceModel) (diags diag.Diagnostics) {
	var propDiags diag.Diagnostics

	state.Id = types.Int64Value(fileComment.Id)
	state.Body = types.StringValue(fileComment.Body)
	state.Reactions, propDiags = lib.ToDynamic(ctx, path.Root("reactions"), fileComment.Reactions, state.Reactions.UnderlyingValue())
	diags.Append(propDiags...)

	return
}
