package provider

import (
	"context"
	"fmt"
	"strconv"
	"strings"
	"time"

	files_sdk "github.com/Files-com/files-sdk-go/v3"
	gpg_key "github.com/Files-com/files-sdk-go/v3/gpgkey"
	"github.com/Files-com/terraform-provider-files/lib"
	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/boolplanmodifier"
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
	Name                  types.String `tfsdk:"name"`
	PartnerId             types.Int64  `tfsdk:"partner_id"`
	UserId                types.Int64  `tfsdk:"user_id"`
	PublicKey             types.String `tfsdk:"public_key"`
	PrivateKey            types.String `tfsdk:"private_key"`
	PrivateKeyPassword    types.String `tfsdk:"private_key_password"`
	GenerateExpiresAt     types.String `tfsdk:"generate_expires_at"`
	GenerateKeypair       types.Bool   `tfsdk:"generate_keypair"`
	GenerateFullName      types.String `tfsdk:"generate_full_name"`
	GenerateEmail         types.String `tfsdk:"generate_email"`
	Id                    types.Int64  `tfsdk:"id"`
	ExpiresAt             types.String `tfsdk:"expires_at"`
	PartnerName           types.String `tfsdk:"partner_name"`
	PublicKeyMd5          types.String `tfsdk:"public_key_md5"`
	PrivateKeyMd5         types.String `tfsdk:"private_key_md5"`
	GeneratedPublicKey    types.String `tfsdk:"generated_public_key"`
	GeneratedPrivateKey   types.String `tfsdk:"generated_private_key"`
	PrivateKeyPasswordMd5 types.String `tfsdk:"private_key_password_md5"`
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
		Description: "A GPGKey object on Files.com is used to securely store both the private and public keys associated with a GPG (GNU Privacy Guard) encryption key pair. This object enables the encryption and decryption of data using GPG, allowing you to protect sensitive information.\n\n\n\nThe private key is kept confidential and is used for decrypting data or signing messages to prove authenticity, while the public key is used to encrypt messages that only the owner of the private key can decrypt.\n\n\n\nBy storing both keys together in a GPGKey object, Files.com makes it easier to understand encryption operations, ensuring secure and efficient handling of encrypted data within the platform.",
		Attributes: map[string]schema.Attribute{
			"name": schema.StringAttribute{
				Description: "Your GPG key name.",
				Required:    true,
			},
			"partner_id": schema.Int64Attribute{
				Description: "Partner ID who owns this GPG Key, if applicable.",
				Computed:    true,
				Optional:    true,
				PlanModifiers: []planmodifier.Int64{
					int64planmodifier.UseStateForUnknown(),
				},
			},
			"user_id": schema.Int64Attribute{
				Description: "User ID who owns this GPG Key, if applicable.",
				Computed:    true,
				Optional:    true,
				PlanModifiers: []planmodifier.Int64{
					int64planmodifier.UseStateForUnknown(),
					int64planmodifier.RequiresReplace(),
				},
			},
			"public_key": schema.StringAttribute{
				Description: "MD5 hash of your GPG public key",
				Optional:    true,
				WriteOnly:   true,
			},
			"private_key": schema.StringAttribute{
				Description: "MD5 hash of your GPG private key.",
				Optional:    true,
				WriteOnly:   true,
			},
			"private_key_password": schema.StringAttribute{
				Description: "Your GPG private key password. Only required for password protected keys.",
				Optional:    true,
				WriteOnly:   true,
			},
			"generate_expires_at": schema.StringAttribute{
				Description: "Expiration date of the key. Used for the generation of the key. Will be ignored if `generate_keypair` is false.",
				Optional:    true,
				WriteOnly:   true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.RequiresReplace(),
				},
			},
			"generate_keypair": schema.BoolAttribute{
				Description: "If true, generate a new GPG key pair. Can not be used with `public_key`/`private_key`",
				Optional:    true,
				WriteOnly:   true,
				PlanModifiers: []planmodifier.Bool{
					boolplanmodifier.RequiresReplace(),
				},
			},
			"generate_full_name": schema.StringAttribute{
				Description: "Full name of the key owner. Used for the generation of the key. Will be ignored if `generate_keypair` is false.",
				Optional:    true,
				WriteOnly:   true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.RequiresReplace(),
				},
			},
			"generate_email": schema.StringAttribute{
				Description: "Email address of the key owner. Used for the generation of the key. Will be ignored if `generate_keypair` is false.",
				Optional:    true,
				WriteOnly:   true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.RequiresReplace(),
				},
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
			"partner_name": schema.StringAttribute{
				Description: "Name of the Partner who owns this GPG Key, if applicable.",
				Computed:    true,
			},
			"public_key_md5": schema.StringAttribute{
				Description: "MD5 hash of your GPG public key",
				Computed:    true,
			},
			"private_key_md5": schema.StringAttribute{
				Description: "MD5 hash of your GPG private key.",
				Computed:    true,
			},
			"generated_public_key": schema.StringAttribute{
				Description: "Your GPG public key",
				Computed:    true,
			},
			"generated_private_key": schema.StringAttribute{
				Description: "Your GPG private key.",
				Computed:    true,
			},
			"private_key_password_md5": schema.StringAttribute{
				Description: "Your GPG private key password. Only required for password protected keys.",
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
	var config gpgKeyResourceModel
	diags = req.Config.Get(ctx, &config)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	paramsGpgKeyCreate := files_sdk.GpgKeyCreateParams{}
	paramsGpgKeyCreate.UserId = plan.UserId.ValueInt64()
	paramsGpgKeyCreate.PartnerId = plan.PartnerId.ValueInt64()
	paramsGpgKeyCreate.PublicKey = config.PublicKey.ValueString()
	paramsGpgKeyCreate.PrivateKey = config.PrivateKey.ValueString()
	paramsGpgKeyCreate.PrivateKeyPassword = config.PrivateKeyPassword.ValueString()
	paramsGpgKeyCreate.Name = plan.Name.ValueString()
	if !config.GenerateExpiresAt.IsNull() {
		if config.GenerateExpiresAt.ValueString() == "" {
			paramsGpgKeyCreate.GenerateExpiresAt = new(time.Time)
		} else {
			createGenerateExpiresAt, err := time.Parse(time.RFC3339, config.GenerateExpiresAt.ValueString())
			if err != nil {
				resp.Diagnostics.AddAttributeError(
					path.Root("generate_expires_at"),
					"Error Parsing generate_expires_at Time",
					"Could not parse generate_expires_at time: "+err.Error(),
				)
			} else {
				paramsGpgKeyCreate.GenerateExpiresAt = &createGenerateExpiresAt
			}
		}
	}
	if !config.GenerateKeypair.IsNull() && !config.GenerateKeypair.IsUnknown() {
		paramsGpgKeyCreate.GenerateKeypair = config.GenerateKeypair.ValueBoolPointer()
	}
	paramsGpgKeyCreate.GenerateFullName = config.GenerateFullName.ValueString()
	paramsGpgKeyCreate.GenerateEmail = config.GenerateEmail.ValueString()

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
	var config gpgKeyResourceModel
	diags = req.Config.Get(ctx, &config)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	paramsGpgKeyUpdate := files_sdk.GpgKeyUpdateParams{}
	paramsGpgKeyUpdate.Id = plan.Id.ValueInt64()
	paramsGpgKeyUpdate.PartnerId = plan.PartnerId.ValueInt64()
	paramsGpgKeyUpdate.PublicKey = config.PublicKey.ValueString()
	paramsGpgKeyUpdate.PrivateKey = config.PrivateKey.ValueString()
	paramsGpgKeyUpdate.PrivateKeyPassword = config.PrivateKeyPassword.ValueString()
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
	state.PartnerId = types.Int64Value(gpgKey.PartnerId)
	state.PartnerName = types.StringValue(gpgKey.PartnerName)
	state.UserId = types.Int64Value(gpgKey.UserId)
	state.PublicKeyMd5 = types.StringValue(gpgKey.PublicKeyMd5)
	state.PrivateKeyMd5 = types.StringValue(gpgKey.PrivateKeyMd5)
	state.GeneratedPublicKey = types.StringValue(gpgKey.GeneratedPublicKey)
	state.GeneratedPrivateKey = types.StringValue(gpgKey.GeneratedPrivateKey)
	state.PrivateKeyPasswordMd5 = types.StringValue(gpgKey.PrivateKeyPasswordMd5)

	return
}
