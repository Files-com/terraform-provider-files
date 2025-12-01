package provider

import (
	"context"
	"fmt"
	"strconv"
	"strings"

	files_sdk "github.com/Files-com/files-sdk-go/v3"
	file_comment "github.com/Files-com/files-sdk-go/v3/filecomment"
	"github.com/Files-com/terraform-provider-files/lib"
	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/int64planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

var (
	_ resource.Resource                = &fileCommentResource{}
	_ resource.ResourceWithConfigure   = &fileCommentResource{}
	_ resource.ResourceWithImportState = &fileCommentResource{}
)

func NewFileCommentResource() resource.Resource {
	return &fileCommentResource{}
}

type fileCommentResource struct {
	client *file_comment.Client
}

type fileCommentResourceModel struct {
	Body      types.String  `tfsdk:"body"`
	Path      types.String  `tfsdk:"path"`
	Id        types.Int64   `tfsdk:"id"`
	Reactions types.Dynamic `tfsdk:"reactions"`
}

func (r *fileCommentResource) Configure(_ context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {
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

func (r *fileCommentResource) Metadata(_ context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_file_comment"
}

func (r *fileCommentResource) Schema(_ context.Context, _ resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Description: "A FileComment is a comment attached to a file by a user.",
		Attributes: map[string]schema.Attribute{
			"body": schema.StringAttribute{
				Description: "Comment body.",
				Required:    true,
			},
			"path": schema.StringAttribute{
				Description: "File path.",
				Required:    true,
				WriteOnly:   true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.RequiresReplace(),
				},
			},
			"id": schema.Int64Attribute{
				Description: "File Comment ID",
				Computed:    true,
				PlanModifiers: []planmodifier.Int64{
					int64planmodifier.UseStateForUnknown(),
				},
			},
			"reactions": schema.DynamicAttribute{
				Description: "Reactions to this comment.",
				Computed:    true,
			},
		},
	}
}

func (r *fileCommentResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	var plan fileCommentResourceModel
	diags := req.Plan.Get(ctx, &plan)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}
	var config fileCommentResourceModel
	diags = req.Config.Get(ctx, &config)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	paramsFileCommentCreate := files_sdk.FileCommentCreateParams{}
	paramsFileCommentCreate.Body = plan.Body.ValueString()
	paramsFileCommentCreate.Path = config.Path.ValueString()

	if resp.Diagnostics.HasError() {
		return
	}

	fileComment, err := r.client.Create(paramsFileCommentCreate, files_sdk.WithContext(ctx))
	if err != nil {
		resp.Diagnostics.AddError(
			"Error Creating Files FileComment",
			"Could not create file_comment, unexpected error: "+err.Error(),
		)
		return
	}

	diags = r.populateResourceModel(ctx, fileComment, &plan)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	diags = resp.State.Set(ctx, plan)
	resp.Diagnostics.Append(diags...)
}

func (r *fileCommentResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	var state fileCommentResourceModel
	diags := req.State.Get(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	paramsFileCommentListFor := files_sdk.FileCommentListForParams{}
	paramsFileCommentListFor.Path = state.Path.ValueString()

	fileCommentIt, err := r.client.ListFor(paramsFileCommentListFor, files_sdk.WithContext(ctx))
	if err != nil {
		if files_sdk.IsNotExist(err) {
			resp.State.RemoveResource(ctx)
			return
		}

		resp.Diagnostics.AddError(
			"Error Reading Files FileComment",
			"Could not read file_comment id "+fmt.Sprint(state.Id.ValueInt64())+": "+err.Error(),
		)
		return
	}

	var fileComment *files_sdk.FileComment
	for fileCommentIt.Next() {
		entry := fileCommentIt.FileComment()
		if entry.Id == state.Id.ValueInt64() {
			fileComment = &entry
			break
		}
	}

	if err = fileCommentIt.Err(); err != nil {
		if files_sdk.IsNotExist(err) {
			resp.State.RemoveResource(ctx)
			return
		}

		resp.Diagnostics.AddError(
			"Error Reading Files FileComment",
			"Could not read file_comment id "+fmt.Sprint(state.Id.ValueInt64())+": "+err.Error(),
		)
	}

	if fileComment == nil {
		resp.State.RemoveResource(ctx)
		return
	}

	diags = r.populateResourceModel(ctx, *fileComment, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	diags = resp.State.Set(ctx, state)
	resp.Diagnostics.Append(diags...)
}

func (r *fileCommentResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	var plan fileCommentResourceModel
	diags := req.Plan.Get(ctx, &plan)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}
	var config fileCommentResourceModel
	diags = req.Config.Get(ctx, &config)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	paramsFileCommentUpdate := map[string]interface{}{}
	if !plan.Id.IsNull() && !plan.Id.IsUnknown() {
		paramsFileCommentUpdate["id"] = plan.Id.ValueInt64()
	}
	if !config.Body.IsNull() && !config.Body.IsUnknown() {
		paramsFileCommentUpdate["body"] = config.Body.ValueString()
	}

	if resp.Diagnostics.HasError() {
		return
	}

	fileComment, err := r.client.UpdateWithMap(paramsFileCommentUpdate, files_sdk.WithContext(ctx))
	if err != nil {
		resp.Diagnostics.AddError(
			"Error Updating Files FileComment",
			"Could not update file_comment, unexpected error: "+err.Error(),
		)
		return
	}

	diags = r.populateResourceModel(ctx, fileComment, &plan)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	diags = resp.State.Set(ctx, plan)
	resp.Diagnostics.Append(diags...)
}

func (r *fileCommentResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	var state fileCommentResourceModel
	diags := req.State.Get(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	paramsFileCommentDelete := files_sdk.FileCommentDeleteParams{}
	paramsFileCommentDelete.Id = state.Id.ValueInt64()

	err := r.client.Delete(paramsFileCommentDelete, files_sdk.WithContext(ctx))
	if err != nil && !files_sdk.IsNotExist(err) {
		resp.Diagnostics.AddError(
			"Error Deleting Files FileComment",
			"Could not delete file_comment id "+fmt.Sprint(state.Id.ValueInt64())+": "+err.Error(),
		)
	}
}

func (r *fileCommentResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	idParts := strings.SplitN(req.ID, ",", 2)

	if len(idParts) != 2 || idParts[0] == "" || idParts[1] == "" {
		resp.Diagnostics.AddError(
			"Unexpected Import Identifier",
			fmt.Sprintf("Expected import identifier with format: id,path. Got: %q", req.ID),
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
	resp.Diagnostics.Append(resp.State.SetAttribute(ctx, path.Root("path"), idParts[1])...)

}

func (r *fileCommentResource) populateResourceModel(ctx context.Context, fileComment files_sdk.FileComment, state *fileCommentResourceModel) (diags diag.Diagnostics) {
	var propDiags diag.Diagnostics

	state.Id = types.Int64Value(fileComment.Id)
	state.Body = types.StringValue(fileComment.Body)
	state.Reactions, propDiags = lib.ToDynamic(ctx, path.Root("reactions"), fileComment.Reactions, state.Reactions.UnderlyingValue())
	diags.Append(propDiags...)

	return
}
