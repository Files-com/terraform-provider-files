package provider

import (
	"context"
	"fmt"

	files_sdk "github.com/Files-com/files-sdk-go/v3"
	pending_work_event "github.com/Files-com/files-sdk-go/v3/pendingworkevent"
	"github.com/Files-com/terraform-provider-files/lib"
	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/datasource/schema"
	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

var (
	_ datasource.DataSource              = &pendingWorkEventDataSource{}
	_ datasource.DataSourceWithConfigure = &pendingWorkEventDataSource{}
)

func NewPendingWorkEventDataSource() datasource.DataSource {
	return &pendingWorkEventDataSource{}
}

type pendingWorkEventDataSource struct {
	client *pending_work_event.Client
}

type pendingWorkEventDataSourceModel struct {
	Id               types.Int64  `tfsdk:"id"`
	EventType        types.String `tfsdk:"event_type"`
	Status           types.String `tfsdk:"status"`
	Body             types.String `tfsdk:"body"`
	EventErrors      types.List   `tfsdk:"event_errors"`
	CreatedAt        types.String `tfsdk:"created_at"`
	BodyUrl          types.String `tfsdk:"body_url"`
	FolderBehaviorId types.Int64  `tfsdk:"folder_behavior_id"`
}

func (r *pendingWorkEventDataSource) Configure(_ context.Context, req datasource.ConfigureRequest, resp *datasource.ConfigureResponse) {
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

	r.client = &pending_work_event.Client{Config: sdk_config}
}

func (r *pendingWorkEventDataSource) Metadata(_ context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_pending_work_event"
}

func (r *pendingWorkEventDataSource) Schema(_ context.Context, _ datasource.SchemaRequest, resp *datasource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Description: "A PendingWorkEvent is a log record for pending file work failures, such as GPG Encryption and GPG Decryption.",
		Attributes: map[string]schema.Attribute{
			"id": schema.Int64Attribute{
				Description: "Event ID",
				Required:    true,
			},
			"event_type": schema.StringAttribute{
				Description: "Type of pending work event being recorded.",
				Computed:    true,
			},
			"status": schema.StringAttribute{
				Description: "Status of event.",
				Computed:    true,
			},
			"body": schema.StringAttribute{
				Description: "Event body.",
				Computed:    true,
			},
			"event_errors": schema.ListAttribute{
				Description: "Event errors.",
				Computed:    true,
				ElementType: types.StringType,
			},
			"created_at": schema.StringAttribute{
				Description: "Event create date/time.",
				Computed:    true,
			},
			"body_url": schema.StringAttribute{
				Description: "Link to log file.",
				Computed:    true,
			},
			"folder_behavior_id": schema.Int64Attribute{
				Description: "Folder Behavior ID.",
				Computed:    true,
			},
		},
	}
}

func (r *pendingWorkEventDataSource) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
	var data pendingWorkEventDataSourceModel
	diags := req.Config.Get(ctx, &data)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	paramsPendingWorkEventFind := files_sdk.PendingWorkEventFindParams{}
	paramsPendingWorkEventFind.Id = data.Id.ValueInt64()

	pendingWorkEvent, err := r.client.Find(paramsPendingWorkEventFind, files_sdk.WithContext(ctx))
	if err != nil {
		resp.Diagnostics.AddError(
			"Error Reading Files PendingWorkEvent",
			"Could not read pending_work_event id "+fmt.Sprint(data.Id.ValueInt64())+": "+err.Error(),
		)
		return
	}

	diags = r.populateDataSourceModel(ctx, pendingWorkEvent, &data)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	diags = resp.State.Set(ctx, data)
	resp.Diagnostics.Append(diags...)
}

func (r *pendingWorkEventDataSource) populateDataSourceModel(ctx context.Context, pendingWorkEvent files_sdk.PendingWorkEvent, state *pendingWorkEventDataSourceModel) (diags diag.Diagnostics) {
	var propDiags diag.Diagnostics

	state.Id = types.Int64Value(pendingWorkEvent.Id)
	state.EventType = types.StringValue(pendingWorkEvent.EventType)
	state.Status = types.StringValue(pendingWorkEvent.Status)
	state.Body = types.StringValue(pendingWorkEvent.Body)
	state.EventErrors, propDiags = types.ListValueFrom(ctx, types.StringType, pendingWorkEvent.EventErrors)
	diags.Append(propDiags...)
	if err := lib.TimeToStringType(ctx, path.Root("created_at"), pendingWorkEvent.CreatedAt, &state.CreatedAt); err != nil {
		diags.AddError(
			"Error Creating Files PendingWorkEvent",
			"Could not convert state created_at to string: "+err.Error(),
		)
	}
	state.BodyUrl = types.StringValue(pendingWorkEvent.BodyUrl)
	state.FolderBehaviorId = types.Int64Value(pendingWorkEvent.FolderBehaviorId)

	return
}
