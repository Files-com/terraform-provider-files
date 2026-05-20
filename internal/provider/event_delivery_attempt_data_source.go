package provider

import (
	"context"
	"fmt"

	files_sdk "github.com/Files-com/files-sdk-go/v3"
	event_delivery_attempt "github.com/Files-com/files-sdk-go/v3/eventdeliveryattempt"
	"github.com/Files-com/terraform-provider-files/lib"
	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/datasource/schema"
	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

var (
	_ datasource.DataSource              = &eventDeliveryAttemptDataSource{}
	_ datasource.DataSourceWithConfigure = &eventDeliveryAttemptDataSource{}
)

func NewEventDeliveryAttemptDataSource() datasource.DataSource {
	return &eventDeliveryAttemptDataSource{}
}

type eventDeliveryAttemptDataSource struct {
	client *event_delivery_attempt.Client
}

type eventDeliveryAttemptDataSourceModel struct {
	Id                  types.Int64  `tfsdk:"id"`
	EventRecordId       types.Int64  `tfsdk:"event_record_id"`
	EventSubscriptionId types.Int64  `tfsdk:"event_subscription_id"`
	EventTargetId       types.Int64  `tfsdk:"event_target_id"`
	WorkspaceId         types.Int64  `tfsdk:"workspace_id"`
	Status              types.String `tfsdk:"status"`
	AttemptNumber       types.Int64  `tfsdk:"attempt_number"`
	ResponseCode        types.Int64  `tfsdk:"response_code"`
	ErrorMessage        types.String `tfsdk:"error_message"`
	ResponseBody        types.String `tfsdk:"response_body"`
	LatencyMs           types.Int64  `tfsdk:"latency_ms"`
	DeliveredAt         types.String `tfsdk:"delivered_at"`
	LastAttemptedAt     types.String `tfsdk:"last_attempted_at"`
	NextAttemptAt       types.String `tfsdk:"next_attempt_at"`
	CreatedAt           types.String `tfsdk:"created_at"`
}

func (r *eventDeliveryAttemptDataSource) Configure(_ context.Context, req datasource.ConfigureRequest, resp *datasource.ConfigureResponse) {
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

	r.client = &event_delivery_attempt.Client{Config: sdk_config}
}

func (r *eventDeliveryAttemptDataSource) Metadata(_ context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_event_delivery_attempt"
}

func (r *eventDeliveryAttemptDataSource) Schema(_ context.Context, _ datasource.SchemaRequest, resp *datasource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Description: "An EventDeliveryAttempt records delivery state for an EventRecord and EventTarget.",
		Attributes: map[string]schema.Attribute{
			"id": schema.Int64Attribute{
				Description: "Event Delivery Attempt ID",
				Required:    true,
			},
			"event_record_id": schema.Int64Attribute{
				Description: "Event Record ID",
				Computed:    true,
			},
			"event_subscription_id": schema.Int64Attribute{
				Description: "Event Subscription ID",
				Computed:    true,
			},
			"event_target_id": schema.Int64Attribute{
				Description: "Event Target ID",
				Computed:    true,
			},
			"workspace_id": schema.Int64Attribute{
				Description: "Workspace ID. 0 means the default workspace or site-wide.",
				Computed:    true,
			},
			"status": schema.StringAttribute{
				Description: "Delivery status.",
				Computed:    true,
			},
			"attempt_number": schema.Int64Attribute{
				Description: "Number of delivery attempts made.",
				Computed:    true,
			},
			"response_code": schema.Int64Attribute{
				Description: "HTTP response code, if applicable.",
				Computed:    true,
			},
			"error_message": schema.StringAttribute{
				Description: "Delivery error message, if applicable.",
				Computed:    true,
			},
			"response_body": schema.StringAttribute{
				Description: "Delivery response body, if applicable.",
				Computed:    true,
			},
			"latency_ms": schema.Int64Attribute{
				Description: "Delivery latency in milliseconds.",
				Computed:    true,
			},
			"delivered_at": schema.StringAttribute{
				Description: "Successful delivery date/time.",
				Computed:    true,
			},
			"last_attempted_at": schema.StringAttribute{
				Description: "Most recent attempt date/time.",
				Computed:    true,
			},
			"next_attempt_at": schema.StringAttribute{
				Description: "Next scheduled attempt date/time.",
				Computed:    true,
			},
			"created_at": schema.StringAttribute{
				Description: "Delivery Attempt create date/time.",
				Computed:    true,
			},
		},
	}
}

func (r *eventDeliveryAttemptDataSource) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
	var data eventDeliveryAttemptDataSourceModel
	diags := req.Config.Get(ctx, &data)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	paramsEventDeliveryAttemptFind := files_sdk.EventDeliveryAttemptFindParams{}
	paramsEventDeliveryAttemptFind.Id = data.Id.ValueInt64()

	eventDeliveryAttempt, err := r.client.Find(paramsEventDeliveryAttemptFind, files_sdk.WithContext(ctx))
	if err != nil {
		resp.Diagnostics.AddError(
			"Error Reading Files EventDeliveryAttempt",
			"Could not read event_delivery_attempt id "+fmt.Sprint(data.Id.ValueInt64())+": "+err.Error(),
		)
		return
	}

	diags = r.populateDataSourceModel(ctx, eventDeliveryAttempt, &data)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	diags = resp.State.Set(ctx, data)
	resp.Diagnostics.Append(diags...)
}

func (r *eventDeliveryAttemptDataSource) populateDataSourceModel(ctx context.Context, eventDeliveryAttempt files_sdk.EventDeliveryAttempt, state *eventDeliveryAttemptDataSourceModel) (diags diag.Diagnostics) {
	state.Id = types.Int64Value(eventDeliveryAttempt.Id)
	state.EventRecordId = types.Int64Value(eventDeliveryAttempt.EventRecordId)
	state.EventSubscriptionId = types.Int64Value(eventDeliveryAttempt.EventSubscriptionId)
	state.EventTargetId = types.Int64Value(eventDeliveryAttempt.EventTargetId)
	state.WorkspaceId = types.Int64Value(eventDeliveryAttempt.WorkspaceId)
	state.Status = types.StringValue(eventDeliveryAttempt.Status)
	state.AttemptNumber = types.Int64Value(eventDeliveryAttempt.AttemptNumber)
	state.ResponseCode = types.Int64Value(eventDeliveryAttempt.ResponseCode)
	state.ErrorMessage = types.StringValue(eventDeliveryAttempt.ErrorMessage)
	state.ResponseBody = types.StringValue(eventDeliveryAttempt.ResponseBody)
	state.LatencyMs = types.Int64Value(eventDeliveryAttempt.LatencyMs)
	if err := lib.TimeToStringType(ctx, path.Root("delivered_at"), eventDeliveryAttempt.DeliveredAt, &state.DeliveredAt); err != nil {
		diags.AddError(
			"Error Creating Files EventDeliveryAttempt",
			"Could not convert state delivered_at to string: "+err.Error(),
		)
	}
	if err := lib.TimeToStringType(ctx, path.Root("last_attempted_at"), eventDeliveryAttempt.LastAttemptedAt, &state.LastAttemptedAt); err != nil {
		diags.AddError(
			"Error Creating Files EventDeliveryAttempt",
			"Could not convert state last_attempted_at to string: "+err.Error(),
		)
	}
	if err := lib.TimeToStringType(ctx, path.Root("next_attempt_at"), eventDeliveryAttempt.NextAttemptAt, &state.NextAttemptAt); err != nil {
		diags.AddError(
			"Error Creating Files EventDeliveryAttempt",
			"Could not convert state next_attempt_at to string: "+err.Error(),
		)
	}
	if err := lib.TimeToStringType(ctx, path.Root("created_at"), eventDeliveryAttempt.CreatedAt, &state.CreatedAt); err != nil {
		diags.AddError(
			"Error Creating Files EventDeliveryAttempt",
			"Could not convert state created_at to string: "+err.Error(),
		)
	}

	return
}
