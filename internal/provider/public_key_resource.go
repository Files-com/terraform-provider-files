package provider

import (
	"context"
	"fmt"
	"strconv"
	"strings"

	files_sdk "github.com/Files-com/files-sdk-go/v3"
	public_key "github.com/Files-com/files-sdk-go/v3/publickey"
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
	_ resource.Resource                = &publicKeyResource{}
	_ resource.ResourceWithConfigure   = &publicKeyResource{}
	_ resource.ResourceWithImportState = &publicKeyResource{}
)

func NewPublicKeyResource() resource.Resource {
	return &publicKeyResource{}
}

type publicKeyResource struct {
	client *public_key.Client
}

type publicKeyResourceModel struct {
	Id                types.Int64  `tfsdk:"id"`
	Title             types.String `tfsdk:"title"`
	CreatedAt         types.String `tfsdk:"created_at"`
	Fingerprint       types.String `tfsdk:"fingerprint"`
	FingerprintSha256 types.String `tfsdk:"fingerprint_sha256"`
	Username          types.String `tfsdk:"username"`
	UserId            types.Int64  `tfsdk:"user_id"`
	PublicKey         types.String `tfsdk:"public_key"`
}

func (r *publicKeyResource) Configure(_ context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {
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

	r.client = &public_key.Client{Config: sdk_config}
}

func (r *publicKeyResource) Metadata(_ context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_public_key"
}

func (r *publicKeyResource) Schema(_ context.Context, _ resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Description: "Public keys are used by Users who want to connect via SFTP/SSH.\n\n(Note that our SSH support is limited to file operations only, no shell is provided.)",
		Attributes: map[string]schema.Attribute{
			"id": schema.Int64Attribute{
				Description: "Public key ID",
				Computed:    true,
				PlanModifiers: []planmodifier.Int64{
					int64planmodifier.UseStateForUnknown(),
				},
			},
			"title": schema.StringAttribute{
				Description: "Public key title",
				Required:    true,
			},
			"created_at": schema.StringAttribute{
				Description: "Public key created at date/time",
				Computed:    true,
			},
			"fingerprint": schema.StringAttribute{
				Description: "Public key fingerprint (MD5)",
				Computed:    true,
			},
			"fingerprint_sha256": schema.StringAttribute{
				Description: "Public key fingerprint (SHA256)",
				Computed:    true,
			},
			"username": schema.StringAttribute{
				Description: "Username of the user this public key is associated with",
				Computed:    true,
			},
			"user_id": schema.Int64Attribute{
				Description: "User ID this public key is associated with",
				Computed:    true,
				Optional:    true,
				PlanModifiers: []planmodifier.Int64{
					int64planmodifier.UseStateForUnknown(),
					int64planmodifier.RequiresReplace(),
				},
			},
			"public_key": schema.StringAttribute{
				Description: "Actual contents of SSH key.",
				Required:    true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.RequiresReplace(),
				},
			},
		},
	}
}

func (r *publicKeyResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	var plan publicKeyResourceModel
	diags := req.Plan.Get(ctx, &plan)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	paramsPublicKeyCreate := files_sdk.PublicKeyCreateParams{}
	paramsPublicKeyCreate.UserId = plan.UserId.ValueInt64()
	paramsPublicKeyCreate.Title = plan.Title.ValueString()
	paramsPublicKeyCreate.PublicKey = plan.PublicKey.ValueString()

	if resp.Diagnostics.HasError() {
		return
	}

	publicKey, err := r.client.Create(paramsPublicKeyCreate, files_sdk.WithContext(ctx))
	if err != nil {
		resp.Diagnostics.AddError(
			"Error Creating Files PublicKey",
			"Could not create public_key, unexpected error: "+err.Error(),
		)
		return
	}

	diags = r.populateResourceModel(ctx, publicKey, &plan)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	diags = resp.State.Set(ctx, plan)
	resp.Diagnostics.Append(diags...)
}

func (r *publicKeyResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	var state publicKeyResourceModel
	diags := req.State.Get(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	paramsPublicKeyFind := files_sdk.PublicKeyFindParams{}
	paramsPublicKeyFind.Id = state.Id.ValueInt64()

	publicKey, err := r.client.Find(paramsPublicKeyFind, files_sdk.WithContext(ctx))
	if err != nil {
		resp.Diagnostics.AddError(
			"Error Reading Files PublicKey",
			"Could not read public_key id "+fmt.Sprint(state.Id.ValueInt64())+": "+err.Error(),
		)
		return
	}

	diags = r.populateResourceModel(ctx, publicKey, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	diags = resp.State.Set(ctx, state)
	resp.Diagnostics.Append(diags...)
}

func (r *publicKeyResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	var plan publicKeyResourceModel
	diags := req.Plan.Get(ctx, &plan)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	paramsPublicKeyUpdate := files_sdk.PublicKeyUpdateParams{}
	paramsPublicKeyUpdate.Id = plan.Id.ValueInt64()
	paramsPublicKeyUpdate.Title = plan.Title.ValueString()

	if resp.Diagnostics.HasError() {
		return
	}

	publicKey, err := r.client.Update(paramsPublicKeyUpdate, files_sdk.WithContext(ctx))
	if err != nil {
		resp.Diagnostics.AddError(
			"Error Updating Files PublicKey",
			"Could not update public_key, unexpected error: "+err.Error(),
		)
		return
	}

	diags = r.populateResourceModel(ctx, publicKey, &plan)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	diags = resp.State.Set(ctx, plan)
	resp.Diagnostics.Append(diags...)
}

func (r *publicKeyResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	var state publicKeyResourceModel
	diags := req.State.Get(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	paramsPublicKeyDelete := files_sdk.PublicKeyDeleteParams{}
	paramsPublicKeyDelete.Id = state.Id.ValueInt64()

	err := r.client.Delete(paramsPublicKeyDelete, files_sdk.WithContext(ctx))
	if err != nil {
		resp.Diagnostics.AddError(
			"Error Deleting Files PublicKey",
			"Could not delete public_key id "+fmt.Sprint(state.Id.ValueInt64())+": "+err.Error(),
		)
	}
}

func (r *publicKeyResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
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

func (r *publicKeyResource) populateResourceModel(ctx context.Context, publicKey files_sdk.PublicKey, state *publicKeyResourceModel) (diags diag.Diagnostics) {
	state.Id = types.Int64Value(publicKey.Id)
	state.Title = types.StringValue(publicKey.Title)
	if err := lib.TimeToStringType(ctx, path.Root("created_at"), publicKey.CreatedAt, &state.CreatedAt); err != nil {
		diags.AddError(
			"Error Creating Files PublicKey",
			"Could not convert state created_at to string: "+err.Error(),
		)
	}
	state.Fingerprint = types.StringValue(publicKey.Fingerprint)
	state.FingerprintSha256 = types.StringValue(publicKey.FingerprintSha256)
	state.Username = types.StringValue(publicKey.Username)
	state.UserId = types.Int64Value(publicKey.UserId)

	return
}
