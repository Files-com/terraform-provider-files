package provider

import (
	"context"
	"fmt"
	"strconv"
	"strings"

	files_sdk "github.com/Files-com/files-sdk-go/v3"
	public_key "github.com/Files-com/files-sdk-go/v3/publickey"
	"github.com/Files-com/terraform-provider-files/lib"
	"github.com/hashicorp/terraform-plugin-framework-validators/stringvalidator"
	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/boolplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/int64planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/schema/validator"
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
	Title                      types.String `tfsdk:"title"`
	PublicKey                  types.String `tfsdk:"public_key"`
	UserId                     types.Int64  `tfsdk:"user_id"`
	GenerateKeypair            types.Bool   `tfsdk:"generate_keypair"`
	GeneratePrivateKeyPassword types.String `tfsdk:"generate_private_key_password"`
	GenerateAlgorithm          types.String `tfsdk:"generate_algorithm"`
	GenerateLength             types.Int64  `tfsdk:"generate_length"`
	Id                         types.Int64  `tfsdk:"id"`
	CreatedAt                  types.String `tfsdk:"created_at"`
	Fingerprint                types.String `tfsdk:"fingerprint"`
	FingerprintSha256          types.String `tfsdk:"fingerprint_sha256"`
	Status                     types.String `tfsdk:"status"`
	LastLoginAt                types.String `tfsdk:"last_login_at"`
	PrivateKey                 types.String `tfsdk:"private_key"`
	Username                   types.String `tfsdk:"username"`
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
		Description: "A PublicKey is used to authenticate to Files.com via SFTP (SSH File Transfer Protocol). This method of authentication allows users to use their private key (which is never shared with Files.com) to authenticate themselves against the PublicKey stored on Files.com.\n\n\n\nWhen a user configures their PublicKey, it allows them to bypass traditional password-based authentication, leveraging the security of key-based authentication instead.\n\n\n\nNote that Files.com’s SSH support is limited to file operations only. While users can securely transfer files and manage their data via SFTP, they do not have access to a full shell environment for executing arbitrary commands.",
		Attributes: map[string]schema.Attribute{
			"title": schema.StringAttribute{
				Description: "Public key title",
				Required:    true,
			},
			"public_key": schema.StringAttribute{
				Description: "Public key generated for the user.",
				Computed:    true,
				Optional:    true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
					stringplanmodifier.RequiresReplace(),
				},
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
			"generate_keypair": schema.BoolAttribute{
				Description: "If true, generate a new SSH key pair. Can not be used with `public_key`",
				Optional:    true,
				PlanModifiers: []planmodifier.Bool{
					boolplanmodifier.RequiresReplace(),
				},
			},
			"generate_private_key_password": schema.StringAttribute{
				Description: "Password for the private key. Used for the generation of the key. Will be ignored if `generate_keypair` is false.",
				Optional:    true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.RequiresReplace(),
				},
			},
			"generate_algorithm": schema.StringAttribute{
				Description: "Type of key to generate.  One of rsa, dsa, ecdsa, ed25519. Used for the generation of the key. Will be ignored if `generate_keypair` is false.",
				Optional:    true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.RequiresReplace(),
				},
			},
			"generate_length": schema.Int64Attribute{
				Description: "Length of key to generate. If algorithm is ecdsa, this is the signature size. Used for the generation of the key. Will be ignored if `generate_keypair` is false.",
				Optional:    true,
				PlanModifiers: []planmodifier.Int64{
					int64planmodifier.RequiresReplace(),
				},
			},
			"id": schema.Int64Attribute{
				Description: "Public key ID",
				Computed:    true,
				PlanModifiers: []planmodifier.Int64{
					int64planmodifier.UseStateForUnknown(),
				},
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
			"status": schema.StringAttribute{
				Description: "Can be invalid, not_generated, generating, complete",
				Computed:    true,
				Validators: []validator.String{
					stringvalidator.OneOf("error", "not_set", "to_be_generated", "generating", "complete"),
				},
			},
			"last_login_at": schema.StringAttribute{
				Description: "Key's most recent login time via SFTP",
				Computed:    true,
			},
			"private_key": schema.StringAttribute{
				Description: "Private key generated for the user.",
				Computed:    true,
			},
			"username": schema.StringAttribute{
				Description: "Username of the user this public key is associated with",
				Computed:    true,
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
	if !plan.GenerateKeypair.IsNull() && !plan.GenerateKeypair.IsUnknown() {
		paramsPublicKeyCreate.GenerateKeypair = plan.GenerateKeypair.ValueBoolPointer()
	}
	paramsPublicKeyCreate.GeneratePrivateKeyPassword = plan.GeneratePrivateKeyPassword.ValueString()
	paramsPublicKeyCreate.GenerateAlgorithm = plan.GenerateAlgorithm.ValueString()
	paramsPublicKeyCreate.GenerateLength = plan.GenerateLength.ValueInt64()

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
		if files_sdk.IsNotExist(err) {
			resp.State.RemoveResource(ctx)
			return
		}

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
	if err != nil && !files_sdk.IsNotExist(err) {
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
	state.Status = types.StringValue(publicKey.Status)
	if err := lib.TimeToStringType(ctx, path.Root("last_login_at"), publicKey.LastLoginAt, &state.LastLoginAt); err != nil {
		diags.AddError(
			"Error Creating Files PublicKey",
			"Could not convert state last_login_at to string: "+err.Error(),
		)
	}
	state.PrivateKey = types.StringValue(publicKey.PrivateKey)
	state.PublicKey = types.StringValue(publicKey.PublicKey)
	state.Username = types.StringValue(publicKey.Username)
	state.UserId = types.Int64Value(publicKey.UserId)

	return
}
