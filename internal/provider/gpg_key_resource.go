package provider

import (
	"context"
	"fmt"
	"strconv"
	"strings"

	files_sdk "github.com/Files-com/files-sdk-go/v3"
	gpg_key "github.com/Files-com/files-sdk-go/v3/gpgkey"
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
	_ resource.Resource                = &gpgKeyResource{}
	_ resource.ResourceWithConfigure   = &gpgKeyResource{}
	_ resource.ResourceWithImportState = &gpgKeyResource{}
)

func NewGpgKeyResource() resource.Resource {
	return &gpgKeyResource{}
}

type gpgKeyResource struct {
	client *gpg_key.Client
}

type gpgKeyResourceModel struct {
	Name                   types.String `tfsdk:"name"`
	UserId                 types.Int64  `tfsdk:"user_id"`
	PublicKey              types.String `tfsdk:"public_key"`
	PublicKeyHash          types.String `tfsdk:"public_key_hash"`
	PrivateKey             types.String `tfsdk:"private_key"`
	PrivateKeyHash         types.String `tfsdk:"private_key_hash"`
	PrivateKeyPassword     types.String `tfsdk:"private_key_password"`
	PrivateKeyPasswordHash types.String `tfsdk:"private_key_password_hash"`
	Id                     types.Int64  `tfsdk:"id"`
	ExpiresAt              types.String `tfsdk:"expires_at"`
}

func (r *gpgKeyResource) Configure(_ context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {
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

	r.client = &gpg_key.Client{Config: sdk_config}
}

func (r *gpgKeyResource) Metadata(_ context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_gpg_key"
}

func (r *gpgKeyResource) Schema(_ context.Context, _ resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Description: "A GPGKey is a key record for decrypt or encrypt Behavior. It can hold both private and public key in a single record.",
		Attributes: map[string]schema.Attribute{
			"name": schema.StringAttribute{
				Description: "Your GPG key name.",
				Required:    true,
			},
			"user_id": schema.Int64Attribute{
				Description: "GPG owner's user id",
				Computed:    true,
				Optional:    true,
				PlanModifiers: []planmodifier.Int64{
					int64planmodifier.UseStateForUnknown(),
					int64planmodifier.RequiresReplace(),
				},
			},
			"public_key": schema.StringAttribute{
				Description: "Your GPG public key",
				Computed:    true,
				Optional:    true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
			},
			"public_key_hash": schema.StringAttribute{
				Computed: true,
			},
			"private_key": schema.StringAttribute{
				Description: "Your GPG private key.",
				Computed:    true,
				Optional:    true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
			},
			"private_key_hash": schema.StringAttribute{
				Computed: true,
			},
			"private_key_password": schema.StringAttribute{
				Description: "Your GPG private key password. Only required for password protected keys.",
				Computed:    true,
				Optional:    true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
			},
			"private_key_password_hash": schema.StringAttribute{
				Computed: true,
			},
			"id": schema.Int64Attribute{
				Description: "Your GPG key ID.",
				Computed:    true,
				PlanModifiers: []planmodifier.Int64{
					int64planmodifier.UseStateForUnknown(),
				},
			},
			"expires_at": schema.StringAttribute{
				Description: "Your GPG key expiration date.",
				Computed:    true,
			},
		},
	}
}

func (r *gpgKeyResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	var plan gpgKeyResourceModel
	diags := req.Plan.Get(ctx, &plan)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	paramsGpgKeyCreate := files_sdk.GpgKeyCreateParams{}
	paramsGpgKeyCreate.UserId = plan.UserId.ValueInt64()
	paramsGpgKeyCreate.PublicKey = plan.PublicKey.ValueString()
	paramsGpgKeyCreate.PrivateKey = plan.PrivateKey.ValueString()
	paramsGpgKeyCreate.PrivateKeyPassword = plan.PrivateKeyPassword.ValueString()
	paramsGpgKeyCreate.Name = plan.Name.ValueString()

	if resp.Diagnostics.HasError() {
		return
	}

	gpgKey, err := r.client.Create(paramsGpgKeyCreate, files_sdk.WithContext(ctx))
	if err != nil {
		resp.Diagnostics.AddError(
			"Error Creating Files GpgKey",
			"Could not create gpg_key, unexpected error: "+err.Error(),
		)
		return
	}

	diags = r.populateResourceModel(ctx, gpgKey, &plan)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	diags = resp.State.Set(ctx, plan)
	resp.Diagnostics.Append(diags...)
}

func (r *gpgKeyResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	var state gpgKeyResourceModel
	diags := req.State.Get(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	paramsGpgKeyFind := files_sdk.GpgKeyFindParams{}
	paramsGpgKeyFind.Id = state.Id.ValueInt64()

	gpgKey, err := r.client.Find(paramsGpgKeyFind, files_sdk.WithContext(ctx))
	if err != nil {
		if files_sdk.IsNotExist(err) {
			resp.State.RemoveResource(ctx)
			return
		}

		resp.Diagnostics.AddError(
			"Error Reading Files GpgKey",
			"Could not read gpg_key id "+fmt.Sprint(state.Id.ValueInt64())+": "+err.Error(),
		)
		return
	}

	diags = r.populateResourceModel(ctx, gpgKey, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	diags = resp.State.Set(ctx, state)
	resp.Diagnostics.Append(diags...)
}

func (r *gpgKeyResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	var plan gpgKeyResourceModel
	diags := req.Plan.Get(ctx, &plan)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	paramsGpgKeyUpdate := files_sdk.GpgKeyUpdateParams{}
	paramsGpgKeyUpdate.Id = plan.Id.ValueInt64()
	paramsGpgKeyUpdate.PublicKey = plan.PublicKey.ValueString()
	paramsGpgKeyUpdate.PrivateKey = plan.PrivateKey.ValueString()
	paramsGpgKeyUpdate.PrivateKeyPassword = plan.PrivateKeyPassword.ValueString()
	paramsGpgKeyUpdate.Name = plan.Name.ValueString()

	if resp.Diagnostics.HasError() {
		return
	}

	gpgKey, err := r.client.Update(paramsGpgKeyUpdate, files_sdk.WithContext(ctx))
	if err != nil {
		resp.Diagnostics.AddError(
			"Error Updating Files GpgKey",
			"Could not update gpg_key, unexpected error: "+err.Error(),
		)
		return
	}

	diags = r.populateResourceModel(ctx, gpgKey, &plan)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	diags = resp.State.Set(ctx, plan)
	resp.Diagnostics.Append(diags...)
}

func (r *gpgKeyResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	var state gpgKeyResourceModel
	diags := req.State.Get(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	paramsGpgKeyDelete := files_sdk.GpgKeyDeleteParams{}
	paramsGpgKeyDelete.Id = state.Id.ValueInt64()

	err := r.client.Delete(paramsGpgKeyDelete, files_sdk.WithContext(ctx))
	if err != nil && !files_sdk.IsNotExist(err) {
		resp.Diagnostics.AddError(
			"Error Deleting Files GpgKey",
			"Could not delete gpg_key id "+fmt.Sprint(state.Id.ValueInt64())+": "+err.Error(),
		)
	}
}

func (r *gpgKeyResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
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

func (r *gpgKeyResource) populateResourceModel(ctx context.Context, gpgKey files_sdk.GpgKey, state *gpgKeyResourceModel) (diags diag.Diagnostics) {
	state.Id = types.Int64Value(gpgKey.Id)
	if err := lib.TimeToStringType(ctx, path.Root("expires_at"), gpgKey.ExpiresAt, &state.ExpiresAt); err != nil {
		diags.AddError(
			"Error Creating Files GpgKey",
			"Could not convert state expires_at to string: "+err.Error(),
		)
	}
	state.Name = types.StringValue(gpgKey.Name)
	state.UserId = types.Int64Value(gpgKey.UserId)
	state.PublicKeyHash = types.StringValue(gpgKey.PublicKey)
	state.PrivateKeyHash = types.StringValue(gpgKey.PrivateKey)
	state.PrivateKeyPasswordHash = types.StringValue(gpgKey.PrivateKeyPassword)

	return
}
