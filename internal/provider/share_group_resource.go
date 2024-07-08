package provider

import (
	"context"
	"fmt"
	"strconv"
	"strings"

	files_sdk "github.com/Files-com/files-sdk-go/v3"
	share_group "github.com/Files-com/files-sdk-go/v3/sharegroup"
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
	_ resource.Resource                = &shareGroupResource{}
	_ resource.ResourceWithConfigure   = &shareGroupResource{}
	_ resource.ResourceWithImportState = &shareGroupResource{}
)

func NewShareGroupResource() resource.Resource {
	return &shareGroupResource{}
}

type shareGroupResource struct {
	client *share_group.Client
}

type shareGroupResourceModel struct {
	Id      types.Int64   `tfsdk:"id"`
	Name    types.String  `tfsdk:"name"`
	Notes   types.String  `tfsdk:"notes"`
	UserId  types.Int64   `tfsdk:"user_id"`
	Members types.Dynamic `tfsdk:"members"`
}

func (r *shareGroupResource) Configure(_ context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {
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

func (r *shareGroupResource) Metadata(_ context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_share_group"
}

func (r *shareGroupResource) Schema(_ context.Context, _ resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Description: "Share groups allow you to store and name groups of email contacts to be used for sending share and inbox invitations.",
		Attributes: map[string]schema.Attribute{
			"id": schema.Int64Attribute{
				Description: "Share Group ID",
				Computed:    true,
				PlanModifiers: []planmodifier.Int64{
					int64planmodifier.UseStateForUnknown(),
				},
			},
			"name": schema.StringAttribute{
				Description: "Name of the share group",
				Required:    true,
			},
			"notes": schema.StringAttribute{
				Description: "Additional notes of the share group",
				Computed:    true,
				Optional:    true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
			},
			"user_id": schema.Int64Attribute{
				Description: "Owner User ID",
				Computed:    true,
				Optional:    true,
				PlanModifiers: []planmodifier.Int64{
					int64planmodifier.UseStateForUnknown(),
					int64planmodifier.RequiresReplace(),
				},
			},
			"members": schema.DynamicAttribute{
				Description: "A list of share group members",
				Required:    true,
			},
		},
	}
}

func (r *shareGroupResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	var plan shareGroupResourceModel
	diags := req.Plan.Get(ctx, &plan)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	paramsShareGroupCreate := files_sdk.ShareGroupCreateParams{}
	paramsShareGroupCreate.UserId = plan.UserId.ValueInt64()
	paramsShareGroupCreate.Notes = plan.Notes.ValueString()
	paramsShareGroupCreate.Name = plan.Name.ValueString()
	paramsShareGroupCreate.Members, diags = lib.DynamicToStringMapSlice(ctx, path.Root("members"), plan.Members)
	resp.Diagnostics.Append(diags...)

	if resp.Diagnostics.HasError() {
		return
	}

	shareGroup, err := r.client.Create(paramsShareGroupCreate, files_sdk.WithContext(ctx))
	if err != nil {
		resp.Diagnostics.AddError(
			"Error Creating Files ShareGroup",
			"Could not create share_group, unexpected error: "+err.Error(),
		)
		return
	}

	diags = r.populateResourceModel(ctx, shareGroup, &plan)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	diags = resp.State.Set(ctx, plan)
	resp.Diagnostics.Append(diags...)
}

func (r *shareGroupResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	var state shareGroupResourceModel
	diags := req.State.Get(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	paramsShareGroupFind := files_sdk.ShareGroupFindParams{}
	paramsShareGroupFind.Id = state.Id.ValueInt64()

	shareGroup, err := r.client.Find(paramsShareGroupFind, files_sdk.WithContext(ctx))
	if err != nil {
		resp.Diagnostics.AddError(
			"Error Reading Files ShareGroup",
			"Could not read share_group id "+fmt.Sprint(state.Id.ValueInt64())+": "+err.Error(),
		)
		return
	}

	diags = r.populateResourceModel(ctx, shareGroup, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	diags = resp.State.Set(ctx, state)
	resp.Diagnostics.Append(diags...)
}

func (r *shareGroupResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	var plan shareGroupResourceModel
	diags := req.Plan.Get(ctx, &plan)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	paramsShareGroupUpdate := files_sdk.ShareGroupUpdateParams{}
	paramsShareGroupUpdate.Id = plan.Id.ValueInt64()
	paramsShareGroupUpdate.Notes = plan.Notes.ValueString()
	paramsShareGroupUpdate.Name = plan.Name.ValueString()
	paramsShareGroupUpdate.Members, diags = lib.DynamicToStringMapSlice(ctx, path.Root("members"), plan.Members)
	resp.Diagnostics.Append(diags...)

	if resp.Diagnostics.HasError() {
		return
	}

	shareGroup, err := r.client.Update(paramsShareGroupUpdate, files_sdk.WithContext(ctx))
	if err != nil {
		resp.Diagnostics.AddError(
			"Error Updating Files ShareGroup",
			"Could not update share_group, unexpected error: "+err.Error(),
		)
		return
	}

	diags = r.populateResourceModel(ctx, shareGroup, &plan)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	diags = resp.State.Set(ctx, plan)
	resp.Diagnostics.Append(diags...)
}

func (r *shareGroupResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	var state shareGroupResourceModel
	diags := req.State.Get(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	paramsShareGroupDelete := files_sdk.ShareGroupDeleteParams{}
	paramsShareGroupDelete.Id = state.Id.ValueInt64()

	err := r.client.Delete(paramsShareGroupDelete, files_sdk.WithContext(ctx))
	if err != nil {
		resp.Diagnostics.AddError(
			"Error Deleting Files ShareGroup",
			"Could not delete share_group id "+fmt.Sprint(state.Id.ValueInt64())+": "+err.Error(),
		)
	}
}

func (r *shareGroupResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	idParts := strings.SplitN(req.ID, ",", 1)

	if len(idParts) != 1 || idParts[0] == "" {
		resp.Diagnostics.AddError(
			"Unexpected Import Identifier",
			fmt.Sprintf("Expected import identifier with format: id. Got: %q", req.ID),
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

}

func (r *shareGroupResource) populateResourceModel(ctx context.Context, shareGroup files_sdk.ShareGroup, state *shareGroupResourceModel) (diags diag.Diagnostics) {
	var propDiags diag.Diagnostics

	state.Id = types.Int64Value(shareGroup.Id)
	state.Name = types.StringValue(shareGroup.Name)
	state.Notes = types.StringValue(shareGroup.Notes)
	state.UserId = types.Int64Value(shareGroup.UserId)
	state.Members, propDiags = lib.ToDynamic(ctx, path.Root("members"), shareGroup.Members, state.Members.UnderlyingValue())
	diags.Append(propDiags...)

	return
}
