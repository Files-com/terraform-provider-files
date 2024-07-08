package provider

import (
	"context"
	"fmt"
	"strconv"
	"strings"
	"time"

	files_sdk "github.com/Files-com/files-sdk-go/v3"
	snapshot "github.com/Files-com/files-sdk-go/v3/snapshot"
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
	_ resource.Resource                = &snapshotResource{}
	_ resource.ResourceWithConfigure   = &snapshotResource{}
	_ resource.ResourceWithImportState = &snapshotResource{}
)

func NewSnapshotResource() resource.Resource {
	return &snapshotResource{}
}

type snapshotResource struct {
	client *snapshot.Client
}

type snapshotResourceModel struct {
	Id          types.Int64  `tfsdk:"id"`
	ExpiresAt   types.String `tfsdk:"expires_at"`
	FinalizedAt types.String `tfsdk:"finalized_at"`
	Name        types.String `tfsdk:"name"`
	UserId      types.Int64  `tfsdk:"user_id"`
	BundleId    types.Int64  `tfsdk:"bundle_id"`
	Paths       types.List   `tfsdk:"paths"`
}

func (r *snapshotResource) Configure(_ context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {
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

	r.client = &snapshot.Client{Config: sdk_config}
}

func (r *snapshotResource) Metadata(_ context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_snapshot"
}

func (r *snapshotResource) Schema(_ context.Context, _ resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Description: "Snapshots are frozen groups of files in your site's hidden folder.",
		Attributes: map[string]schema.Attribute{
			"id": schema.Int64Attribute{
				Description: "The snapshot's unique ID.",
				Computed:    true,
				PlanModifiers: []planmodifier.Int64{
					int64planmodifier.UseStateForUnknown(),
				},
			},
			"expires_at": schema.StringAttribute{
				Description: "When the snapshot expires.",
				Computed:    true,
				Optional:    true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
			},
			"finalized_at": schema.StringAttribute{
				Description: "When the snapshot was finalized.",
				Computed:    true,
			},
			"name": schema.StringAttribute{
				Description: "A name for the snapshot.",
				Computed:    true,
				Optional:    true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
			},
			"user_id": schema.Int64Attribute{
				Description: "The user that created this snapshot, if applicable.",
				Computed:    true,
			},
			"bundle_id": schema.Int64Attribute{
				Description: "The bundle using this snapshot, if applicable.",
				Computed:    true,
			},
			"paths": schema.ListAttribute{
				Description: "An array of paths to add to the snapshot.",
				Optional:    true,
				ElementType: types.StringType,
			},
		},
	}
}

func (r *snapshotResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	var plan snapshotResourceModel
	diags := req.Plan.Get(ctx, &plan)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	paramsSnapshotCreate := files_sdk.SnapshotCreateParams{}
	if !plan.ExpiresAt.IsNull() && plan.ExpiresAt.ValueString() != "" {
		createExpiresAt, err := time.Parse(time.RFC3339, plan.ExpiresAt.ValueString())
		if err != nil {
			resp.Diagnostics.AddAttributeError(
				path.Root("expires_at"),
				"Error Parsing expires_at Time",
				"Could not parse expires_at time: "+err.Error(),
			)
		} else {
			paramsSnapshotCreate.ExpiresAt = &createExpiresAt
		}
	}
	paramsSnapshotCreate.Name = plan.Name.ValueString()
	if !plan.Paths.IsNull() && !plan.Paths.IsUnknown() {
		diags = plan.Paths.ElementsAs(ctx, &paramsSnapshotCreate.Paths, false)
		resp.Diagnostics.Append(diags...)
	}

	if resp.Diagnostics.HasError() {
		return
	}

	snapshot, err := r.client.Create(paramsSnapshotCreate, files_sdk.WithContext(ctx))
	if err != nil {
		resp.Diagnostics.AddError(
			"Error Creating Files Snapshot",
			"Could not create snapshot, unexpected error: "+err.Error(),
		)
		return
	}

	diags = r.populateResourceModel(ctx, snapshot, &plan)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	diags = resp.State.Set(ctx, plan)
	resp.Diagnostics.Append(diags...)
}

func (r *snapshotResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	var state snapshotResourceModel
	diags := req.State.Get(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	paramsSnapshotFind := files_sdk.SnapshotFindParams{}
	paramsSnapshotFind.Id = state.Id.ValueInt64()

	snapshot, err := r.client.Find(paramsSnapshotFind, files_sdk.WithContext(ctx))
	if err != nil {
		resp.Diagnostics.AddError(
			"Error Reading Files Snapshot",
			"Could not read snapshot id "+fmt.Sprint(state.Id.ValueInt64())+": "+err.Error(),
		)
		return
	}

	diags = r.populateResourceModel(ctx, snapshot, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	diags = resp.State.Set(ctx, state)
	resp.Diagnostics.Append(diags...)
}

func (r *snapshotResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	var plan snapshotResourceModel
	diags := req.Plan.Get(ctx, &plan)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	paramsSnapshotUpdate := files_sdk.SnapshotUpdateParams{}
	paramsSnapshotUpdate.Id = plan.Id.ValueInt64()
	if !plan.ExpiresAt.IsNull() && plan.ExpiresAt.ValueString() != "" {
		updateExpiresAt, err := time.Parse(time.RFC3339, plan.ExpiresAt.ValueString())
		if err != nil {
			resp.Diagnostics.AddAttributeError(
				path.Root("expires_at"),
				"Error Parsing expires_at Time",
				"Could not parse expires_at time: "+err.Error(),
			)
		} else {
			paramsSnapshotUpdate.ExpiresAt = &updateExpiresAt
		}
	}
	paramsSnapshotUpdate.Name = plan.Name.ValueString()
	if !plan.Paths.IsNull() && !plan.Paths.IsUnknown() {
		diags = plan.Paths.ElementsAs(ctx, &paramsSnapshotUpdate.Paths, false)
		resp.Diagnostics.Append(diags...)
	}

	if resp.Diagnostics.HasError() {
		return
	}

	snapshot, err := r.client.Update(paramsSnapshotUpdate, files_sdk.WithContext(ctx))
	if err != nil {
		resp.Diagnostics.AddError(
			"Error Updating Files Snapshot",
			"Could not update snapshot, unexpected error: "+err.Error(),
		)
		return
	}

	diags = r.populateResourceModel(ctx, snapshot, &plan)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	diags = resp.State.Set(ctx, plan)
	resp.Diagnostics.Append(diags...)
}

func (r *snapshotResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	var state snapshotResourceModel
	diags := req.State.Get(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	paramsSnapshotDelete := files_sdk.SnapshotDeleteParams{}
	paramsSnapshotDelete.Id = state.Id.ValueInt64()

	err := r.client.Delete(paramsSnapshotDelete, files_sdk.WithContext(ctx))
	if err != nil {
		resp.Diagnostics.AddError(
			"Error Deleting Files Snapshot",
			"Could not delete snapshot id "+fmt.Sprint(state.Id.ValueInt64())+": "+err.Error(),
		)
	}
}

func (r *snapshotResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
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

func (r *snapshotResource) populateResourceModel(ctx context.Context, snapshot files_sdk.Snapshot, state *snapshotResourceModel) (diags diag.Diagnostics) {
	state.Id = types.Int64Value(snapshot.Id)
	if err := lib.TimeToStringType(ctx, path.Root("expires_at"), snapshot.ExpiresAt, &state.ExpiresAt); err != nil {
		diags.AddError(
			"Error Creating Files Snapshot",
			"Could not convert state expires_at to string: "+err.Error(),
		)
	}
	if err := lib.TimeToStringType(ctx, path.Root("finalized_at"), snapshot.FinalizedAt, &state.FinalizedAt); err != nil {
		diags.AddError(
			"Error Creating Files Snapshot",
			"Could not convert state finalized_at to string: "+err.Error(),
		)
	}
	state.Name = types.StringValue(snapshot.Name)
	state.UserId = types.Int64Value(snapshot.UserId)
	state.BundleId = types.Int64Value(snapshot.BundleId)

	return
}
