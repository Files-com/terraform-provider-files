package provider

import (
	"context"
	"fmt"
	"strconv"
	"strings"

	files_sdk "github.com/Files-com/files-sdk-go/v3"
	sftp_host_key "github.com/Files-com/files-sdk-go/v3/sftphostkey"
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
	_ resource.Resource                = &sftpHostKeyResource{}
	_ resource.ResourceWithConfigure   = &sftpHostKeyResource{}
	_ resource.ResourceWithImportState = &sftpHostKeyResource{}
)

func NewSftpHostKeyResource() resource.Resource {
	return &sftpHostKeyResource{}
}

type sftpHostKeyResource struct {
	client *sftp_host_key.Client
}

type sftpHostKeyResourceModel struct {
	Id                types.Int64  `tfsdk:"id"`
	Name              types.String `tfsdk:"name"`
	FingerprintMd5    types.String `tfsdk:"fingerprint_md5"`
	FingerprintSha256 types.String `tfsdk:"fingerprint_sha256"`
	PrivateKey        types.String `tfsdk:"private_key"`
}

func (r *sftpHostKeyResource) Configure(_ context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {
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

	r.client = &sftp_host_key.Client{Config: sdk_config}
}

func (r *sftpHostKeyResource) Metadata(_ context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_sftp_host_key"
}

func (r *sftpHostKeyResource) Schema(_ context.Context, _ resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Description: "",
		Attributes: map[string]schema.Attribute{
			"id": schema.Int64Attribute{
				Description: "Sftp Host Key ID",
				Computed:    true,
				PlanModifiers: []planmodifier.Int64{
					int64planmodifier.UseStateForUnknown(),
				},
			},
			"name": schema.StringAttribute{
				Description: "The friendly name of this SFTP Host Key.",
				Computed:    true,
				Optional:    true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
			},
			"fingerprint_md5": schema.StringAttribute{
				Description: "MD5 Fingerpint of the public key",
				Computed:    true,
			},
			"fingerprint_sha256": schema.StringAttribute{
				Description: "SHA256 Fingerpint of the public key",
				Computed:    true,
			},
			"private_key": schema.StringAttribute{
				Description: "The private key data.",
				Optional:    true,
			},
		},
	}
}

func (r *sftpHostKeyResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	var plan sftpHostKeyResourceModel
	diags := req.Plan.Get(ctx, &plan)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	paramsSftpHostKeyCreate := files_sdk.SftpHostKeyCreateParams{}
	paramsSftpHostKeyCreate.Name = plan.Name.ValueString()
	paramsSftpHostKeyCreate.PrivateKey = plan.PrivateKey.ValueString()

	if resp.Diagnostics.HasError() {
		return
	}

	sftpHostKey, err := r.client.Create(paramsSftpHostKeyCreate, files_sdk.WithContext(ctx))
	if err != nil {
		resp.Diagnostics.AddError(
			"Error Creating Files SftpHostKey",
			"Could not create sftp_host_key, unexpected error: "+err.Error(),
		)
		return
	}

	diags = r.populateResourceModel(ctx, sftpHostKey, &plan)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	diags = resp.State.Set(ctx, plan)
	resp.Diagnostics.Append(diags...)
}

func (r *sftpHostKeyResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	var state sftpHostKeyResourceModel
	diags := req.State.Get(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	paramsSftpHostKeyFind := files_sdk.SftpHostKeyFindParams{}
	paramsSftpHostKeyFind.Id = state.Id.ValueInt64()

	sftpHostKey, err := r.client.Find(paramsSftpHostKeyFind, files_sdk.WithContext(ctx))
	if err != nil {
		resp.Diagnostics.AddError(
			"Error Reading Files SftpHostKey",
			"Could not read sftp_host_key id "+fmt.Sprint(state.Id.ValueInt64())+": "+err.Error(),
		)
		return
	}

	diags = r.populateResourceModel(ctx, sftpHostKey, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	diags = resp.State.Set(ctx, state)
	resp.Diagnostics.Append(diags...)
}

func (r *sftpHostKeyResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	var plan sftpHostKeyResourceModel
	diags := req.Plan.Get(ctx, &plan)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	paramsSftpHostKeyUpdate := files_sdk.SftpHostKeyUpdateParams{}
	paramsSftpHostKeyUpdate.Id = plan.Id.ValueInt64()
	paramsSftpHostKeyUpdate.Name = plan.Name.ValueString()
	paramsSftpHostKeyUpdate.PrivateKey = plan.PrivateKey.ValueString()

	if resp.Diagnostics.HasError() {
		return
	}

	sftpHostKey, err := r.client.Update(paramsSftpHostKeyUpdate, files_sdk.WithContext(ctx))
	if err != nil {
		resp.Diagnostics.AddError(
			"Error Updating Files SftpHostKey",
			"Could not update sftp_host_key, unexpected error: "+err.Error(),
		)
		return
	}

	diags = r.populateResourceModel(ctx, sftpHostKey, &plan)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	diags = resp.State.Set(ctx, plan)
	resp.Diagnostics.Append(diags...)
}

func (r *sftpHostKeyResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	var state sftpHostKeyResourceModel
	diags := req.State.Get(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	paramsSftpHostKeyDelete := files_sdk.SftpHostKeyDeleteParams{}
	paramsSftpHostKeyDelete.Id = state.Id.ValueInt64()

	err := r.client.Delete(paramsSftpHostKeyDelete, files_sdk.WithContext(ctx))
	if err != nil {
		resp.Diagnostics.AddError(
			"Error Deleting Files SftpHostKey",
			"Could not delete sftp_host_key id "+fmt.Sprint(state.Id.ValueInt64())+": "+err.Error(),
		)
	}
}

func (r *sftpHostKeyResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
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

func (r *sftpHostKeyResource) populateResourceModel(ctx context.Context, sftpHostKey files_sdk.SftpHostKey, state *sftpHostKeyResourceModel) (diags diag.Diagnostics) {
	state.Id = types.Int64Value(sftpHostKey.Id)
	state.Name = types.StringValue(sftpHostKey.Name)
	state.FingerprintMd5 = types.StringValue(sftpHostKey.FingerprintMd5)
	state.FingerprintSha256 = types.StringValue(sftpHostKey.FingerprintSha256)

	return
}
