package provider

import (
	"context"
	"fmt"

	files_sdk "github.com/Files-com/files-sdk-go/v3"
	sso_event "github.com/Files-com/files-sdk-go/v3/ssoevent"
	"github.com/Files-com/terraform-provider-files/lib"
	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/datasource/schema"
	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

var (
	_ datasource.DataSource              = &ssoEventDataSource{}
	_ datasource.DataSourceWithConfigure = &ssoEventDataSource{}
)

func NewSsoEventDataSource() datasource.DataSource {
	return &ssoEventDataSource{}
}

type ssoEventDataSource struct {
	client *sso_event.Client
}

type ssoEventDataSourceModel struct {
	Id            types.Int64  `tfsdk:"id"`
	EventType     types.String `tfsdk:"event_type"`
	Status        types.String `tfsdk:"status"`
	Body          types.String `tfsdk:"body"`
	EventErrors   types.List   `tfsdk:"event_errors"`
	CreatedAt     types.String `tfsdk:"created_at"`
	BodyUrl       types.String `tfsdk:"body_url"`
	UserId        types.Int64  `tfsdk:"user_id"`
	Username      types.String `tfsdk:"username"`
	IdpUid        types.String `tfsdk:"idp_uid"`
	Provider_     types.String `tfsdk:"provider_"`
	ProviderLabel types.String `tfsdk:"provider_label"`
	Ip            types.String `tfsdk:"ip"`
	Region        types.String `tfsdk:"region"`
}

func (r *ssoEventDataSource) Configure(_ context.Context, req datasource.ConfigureRequest, resp *datasource.ConfigureResponse) {
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

	r.client = &sso_event.Client{Config: sdk_config}
}

func (r *ssoEventDataSource) Metadata(_ context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_sso_event"
}

func (r *ssoEventDataSource) Schema(_ context.Context, _ datasource.SchemaRequest, resp *datasource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Description: "An SsoEvent is a log record for SSO-related activity such as LDAP syncs and SSO login attempts.",
		Attributes: map[string]schema.Attribute{
			"id": schema.Int64Attribute{
				Description: "Event ID",
				Required:    true,
			},
			"event_type": schema.StringAttribute{
				Description: "Type of SSO event being recorded.",
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
			"user_id": schema.Int64Attribute{
				Description: "User ID.",
				Computed:    true,
			},
			"username": schema.StringAttribute{
				Description: "Username on Files.com for the SSO login attempt.",
				Computed:    true,
			},
			"idp_uid": schema.StringAttribute{
				Description: "Identity Provider UID for the SSO login attempt.",
				Computed:    true,
			},
			"provider_": schema.StringAttribute{
				Description: "SSO provider for the SSO login attempt.",
				Computed:    true,
			},
			"provider_label": schema.StringAttribute{
				Description: "SSO provider label for the SSO login attempt.",
				Computed:    true,
			},
			"ip": schema.StringAttribute{
				Description: "IP address for the SSO login attempt.",
				Computed:    true,
			},
			"region": schema.StringAttribute{
				Description: "Region for the SSO login attempt.",
				Computed:    true,
			},
		},
	}
}

func (r *ssoEventDataSource) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
	var data ssoEventDataSourceModel
	diags := req.Config.Get(ctx, &data)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	paramsSsoEventFind := files_sdk.SsoEventFindParams{}
	paramsSsoEventFind.Id = data.Id.ValueInt64()

	ssoEvent, err := r.client.Find(paramsSsoEventFind, files_sdk.WithContext(ctx))
	if err != nil {
		resp.Diagnostics.AddError(
			"Error Reading Files SsoEvent",
			"Could not read sso_event id "+fmt.Sprint(data.Id.ValueInt64())+": "+err.Error(),
		)
		return
	}

	diags = r.populateDataSourceModel(ctx, ssoEvent, &data)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	diags = resp.State.Set(ctx, data)
	resp.Diagnostics.Append(diags...)
}

func (r *ssoEventDataSource) populateDataSourceModel(ctx context.Context, ssoEvent files_sdk.SsoEvent, state *ssoEventDataSourceModel) (diags diag.Diagnostics) {
	var propDiags diag.Diagnostics

	state.Id = types.Int64Value(ssoEvent.Id)
	state.EventType = types.StringValue(ssoEvent.EventType)
	state.Status = types.StringValue(ssoEvent.Status)
	state.Body = types.StringValue(ssoEvent.Body)
	state.EventErrors, propDiags = types.ListValueFrom(ctx, types.StringType, ssoEvent.EventErrors)
	diags.Append(propDiags...)
	if err := lib.TimeToStringType(ctx, path.Root("created_at"), ssoEvent.CreatedAt, &state.CreatedAt); err != nil {
		diags.AddError(
			"Error Creating Files SsoEvent",
			"Could not convert state created_at to string: "+err.Error(),
		)
	}
	state.BodyUrl = types.StringValue(ssoEvent.BodyUrl)
	state.UserId = types.Int64Value(ssoEvent.UserId)
	state.Username = types.StringValue(ssoEvent.Username)
	state.IdpUid = types.StringValue(ssoEvent.IdpUid)
	state.Provider_ = types.StringValue(ssoEvent.Provider)
	state.ProviderLabel = types.StringValue(ssoEvent.ProviderLabel)
	state.Ip = types.StringValue(ssoEvent.Ip)
	state.Region = types.StringValue(ssoEvent.Region)

	return
}
