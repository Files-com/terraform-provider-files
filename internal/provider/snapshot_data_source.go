package provider

import (
	"context"
	"fmt"

	files_sdk "github.com/Files-com/files-sdk-go/v3"
	snapshot "github.com/Files-com/files-sdk-go/v3/snapshot"
	"github.com/Files-com/terraform-provider-files/lib"
	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/datasource/schema"
	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

var (
	_ datasource.DataSource              = &snapshotDataSource{}
	_ datasource.DataSourceWithConfigure = &snapshotDataSource{}
)

func NewSnapshotDataSource() datasource.DataSource {
	return &snapshotDataSource{}
}

type snapshotDataSource struct {
	client *snapshot.Client
}

type snapshotDataSourceModel struct {
	Id          types.Int64  `tfsdk:"id"`
	ExpiresAt   types.String `tfsdk:"expires_at"`
	FinalizedAt types.String `tfsdk:"finalized_at"`
	Name        types.String `tfsdk:"name"`
	UserId      types.Int64  `tfsdk:"user_id"`
	BundleId    types.Int64  `tfsdk:"bundle_id"`
}

func (r *snapshotDataSource) Configure(_ context.Context, req datasource.ConfigureRequest, resp *datasource.ConfigureResponse) {
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

func (r *snapshotDataSource) Metadata(_ context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_snapshot"
}

func (r *snapshotDataSource) Schema(_ context.Context, _ datasource.SchemaRequest, resp *datasource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Description: "Snapshots allow you to create a read-only archive of files at a specific point in time. You can define a snapshot, add files to it, and then finalize it. Once finalized, the snapshot’s contents are immutable.\n\n\n\nEach snapshot may have an expiration date. When the expiration date is reached, the snapshot is automatically deleted from the Files.com platform.",
		Attributes: map[string]schema.Attribute{
			"id": schema.Int64Attribute{
				Description: "The snapshot's unique ID.",
				Required:    true,
			},
			"expires_at": schema.StringAttribute{
				Description: "When the snapshot expires.",
				Computed:    true,
			},
			"finalized_at": schema.StringAttribute{
				Description: "When the snapshot was finalized.",
				Computed:    true,
			},
			"name": schema.StringAttribute{
				Description: "A name for the snapshot.",
				Computed:    true,
			},
			"user_id": schema.Int64Attribute{
				Description: "The user that created this snapshot, if applicable.",
				Computed:    true,
			},
			"bundle_id": schema.Int64Attribute{
				Description: "The bundle using this snapshot, if applicable.",
				Computed:    true,
			},
		},
	}
}

func (r *snapshotDataSource) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
	var data snapshotDataSourceModel
	diags := req.Config.Get(ctx, &data)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	paramsSnapshotFind := files_sdk.SnapshotFindParams{}
	paramsSnapshotFind.Id = data.Id.ValueInt64()

	snapshot, err := r.client.Find(paramsSnapshotFind, files_sdk.WithContext(ctx))
	if err != nil {
		resp.Diagnostics.AddError(
			"Error Reading Files Snapshot",
			"Could not read snapshot id "+fmt.Sprint(data.Id.ValueInt64())+": "+err.Error(),
		)
		return
	}

	diags = r.populateDataSourceModel(ctx, snapshot, &data)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	diags = resp.State.Set(ctx, data)
	resp.Diagnostics.Append(diags...)
}

func (r *snapshotDataSource) populateDataSourceModel(ctx context.Context, snapshot files_sdk.Snapshot, state *snapshotDataSourceModel) (diags diag.Diagnostics) {
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
