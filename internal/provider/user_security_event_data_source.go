package provider

import (
	"context"
	"fmt"

	files_sdk "github.com/Files-com/files-sdk-go/v3"
	user_security_event "github.com/Files-com/files-sdk-go/v3/usersecurityevent"
	"github.com/Files-com/terraform-provider-files/lib"
	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/datasource/schema"
	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

var (
	_ datasource.DataSource              = &userSecurityEventDataSource{}
	_ datasource.DataSourceWithConfigure = &userSecurityEventDataSource{}
)

func NewUserSecurityEventDataSource() datasource.DataSource {
	return &userSecurityEventDataSource{}
}

type userSecurityEventDataSource struct {
	client *user_security_event.Client
}

type userSecurityEventDataSourceModel struct {
	Id          types.Int64  `tfsdk:"id"`
	EventType   types.String `tfsdk:"event_type"`
	Body        types.String `tfsdk:"body"`
	EventErrors types.List   `tfsdk:"event_errors"`
	CreatedAt   types.String `tfsdk:"created_at"`
	BodyUrl     types.String `tfsdk:"body_url"`
	UserId      types.Int64  `tfsdk:"user_id"`
}

func (r *userSecurityEventDataSource) Configure(_ context.Context, req datasource.ConfigureRequest, resp *datasource.ConfigureResponse) {
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

	r.client = &user_security_event.Client{Config: sdk_config}
}

func (r *userSecurityEventDataSource) Metadata(_ context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_user_security_event"
}

func (r *userSecurityEventDataSource) Schema(_ context.Context, _ datasource.SchemaRequest, resp *datasource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Description: "A UserSecurityEvent is a log record for user security activity such as user lockouts.",
		Attributes: map[string]schema.Attribute{
			"id": schema.Int64Attribute{
				Description: "Event ID",
				Required:    true,
			},
			"event_type": schema.StringAttribute{
				Description: "Type of user security event being recorded.",
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
			"user_id": schema.Int64Attribute{
				Description: "User ID.",
				Computed:    true,
			},
		},
	}
}

func (r *userSecurityEventDataSource) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
	var data userSecurityEventDataSourceModel
	diags := req.Config.Get(ctx, &data)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	paramsUserSecurityEventFind := files_sdk.UserSecurityEventFindParams{}
	paramsUserSecurityEventFind.Id = data.Id.ValueInt64()

	userSecurityEvent, err := r.client.Find(paramsUserSecurityEventFind, files_sdk.WithContext(ctx))
	if err != nil {
		resp.Diagnostics.AddError(
			"Error Reading Files UserSecurityEvent",
			"Could not read user_security_event id "+fmt.Sprint(data.Id.ValueInt64())+": "+err.Error(),
		)
		return
	}

	diags = r.populateDataSourceModel(ctx, userSecurityEvent, &data)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	diags = resp.State.Set(ctx, data)
	resp.Diagnostics.Append(diags...)
}

func (r *userSecurityEventDataSource) populateDataSourceModel(ctx context.Context, userSecurityEvent files_sdk.UserSecurityEvent, state *userSecurityEventDataSourceModel) (diags diag.Diagnostics) {
	var propDiags diag.Diagnostics

	state.Id = types.Int64Value(userSecurityEvent.Id)
	state.EventType = types.StringValue(userSecurityEvent.EventType)
	state.Body = types.StringValue(userSecurityEvent.Body)
	state.EventErrors, propDiags = types.ListValueFrom(ctx, types.StringType, userSecurityEvent.EventErrors)
	diags.Append(propDiags...)
	if err := lib.TimeToStringType(ctx, path.Root("created_at"), userSecurityEvent.CreatedAt, &state.CreatedAt); err != nil {
		diags.AddError(
			"Error Creating Files UserSecurityEvent",
			"Could not convert state created_at to string: "+err.Error(),
		)
	}
	state.BodyUrl = types.StringValue(userSecurityEvent.BodyUrl)
	state.UserId = types.Int64Value(userSecurityEvent.UserId)

	return
}
