package provider

import (
	"context"
	"encoding/json"
	"fmt"
	"strings"

	files_sdk "github.com/Files-com/files-sdk-go/v3"
	style "github.com/Files-com/files-sdk-go/v3/style"
	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

var (
	_ resource.Resource                = &styleResource{}
	_ resource.ResourceWithConfigure   = &styleResource{}
	_ resource.ResourceWithImportState = &styleResource{}
)

func NewStyleResource() resource.Resource {
	return &styleResource{}
}

type styleResource struct {
	client *style.Client
}

type styleResourceModel struct {
	Path          types.String `tfsdk:"path"`
	LogoClickHref types.String `tfsdk:"logo_click_href"`
	Id            types.Int64  `tfsdk:"id"`
	Logo          types.String `tfsdk:"logo"`
	Thumbnail     types.String `tfsdk:"thumbnail"`
}

func (r *styleResource) Configure(_ context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {
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

	r.client = &style.Client{Config: sdk_config}
}

func (r *styleResource) Metadata(_ context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_style"
}

func (r *styleResource) Schema(_ context.Context, _ resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Description: "A Style is a custom set of branding that can be applied on a per-folder basis.\n\nCurrently these support a logo per folder and an optional click-through URL for public visitors.\n\nIn the future we may extend these to also support colors.\n\nIf you want to see that, please let us know so we can add your vote to the list.",
		Attributes: map[string]schema.Attribute{
			"path": schema.StringAttribute{
				Description: "Folder path. This must be slash-delimited, but it must neither start nor end with a slash. Maximum of 5000 characters.",
				Required:    true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
			},
			"logo_click_href": schema.StringAttribute{
				Description: "URL to open when a public visitor clicks the logo",
				Computed:    true,
				Optional:    true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
			},
			"id": schema.Int64Attribute{
				Description: "Style ID",
				Computed:    true,
			},
			"logo": schema.StringAttribute{
				Description: "Logo",
				Computed:    true,
			},
			"thumbnail": schema.StringAttribute{
				Description: "Logo thumbnail",
				Computed:    true,
			},
		},
	}
}

func (r *styleResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	resp.Diagnostics.AddError(
		"Resource Create Not Implemented",
		"This resource does not support creation. Please import the resource.",
	)
}

func (r *styleResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	var state styleResourceModel
	diags := req.State.Get(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	paramsStyleFind := files_sdk.StyleFindParams{}
	paramsStyleFind.Path = state.Path.ValueString()

	style, err := r.client.Find(paramsStyleFind, files_sdk.WithContext(ctx))
	if err != nil {
		if files_sdk.IsNotExist(err) {
			resp.State.RemoveResource(ctx)
			return
		}

		resp.Diagnostics.AddError(
			"Error Reading Files Style",
			"Could not read style path "+fmt.Sprint(state.Path.ValueString())+": "+err.Error(),
		)
		return
	}

	diags = r.populateResourceModel(ctx, style, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	diags = resp.State.Set(ctx, state)
	resp.Diagnostics.Append(diags...)
}

func (r *styleResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	var plan styleResourceModel
	diags := req.Plan.Get(ctx, &plan)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}
	var config styleResourceModel
	diags = req.Config.Get(ctx, &config)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	paramsStyleUpdate := map[string]interface{}{}
	if !plan.Path.IsNull() && !plan.Path.IsUnknown() {
		paramsStyleUpdate["path"] = plan.Path.ValueString()
	}
	if !config.LogoClickHref.IsNull() && !config.LogoClickHref.IsUnknown() {
		paramsStyleUpdate["logo_click_href"] = config.LogoClickHref.ValueString()
	}

	if resp.Diagnostics.HasError() {
		return
	}

	style, err := r.client.UpdateWithMap(paramsStyleUpdate, files_sdk.WithContext(ctx))
	if err != nil {
		resp.Diagnostics.AddError(
			"Error Updating Files Style",
			"Could not update style, unexpected error: "+err.Error(),
		)
		return
	}

	diags = r.populateResourceModel(ctx, style, &plan)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	diags = resp.State.Set(ctx, plan)
	resp.Diagnostics.Append(diags...)
}

func (r *styleResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	var state styleResourceModel
	diags := req.State.Get(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	paramsStyleDelete := files_sdk.StyleDeleteParams{}
	paramsStyleDelete.Path = state.Path.ValueString()

	err := r.client.Delete(paramsStyleDelete, files_sdk.WithContext(ctx))
	if err != nil && !files_sdk.IsNotExist(err) {
		resp.Diagnostics.AddError(
			"Error Deleting Files Style",
			"Could not delete style path "+fmt.Sprint(state.Path.ValueString())+": "+err.Error(),
		)
	}
}

func (r *styleResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	idParts := strings.SplitN(req.ID, ",", 1)

	if len(idParts) != 1 || idParts[0] == "" {
		resp.Diagnostics.AddError(
			"Unexpected Import Identifier",
			fmt.Sprintf("Expected import identifier with format: path. Got: %q", req.ID),
		)
		return
	}

	resp.Diagnostics.Append(resp.State.SetAttribute(ctx, path.Root("path"), idParts[0])...)

}

func (r *styleResource) populateResourceModel(ctx context.Context, style files_sdk.Style, state *styleResourceModel) (diags diag.Diagnostics) {
	state.Id = types.Int64Value(style.Id)
	state.Path = types.StringValue(style.Path)
	respLogo, err := json.Marshal(style.Logo)
	if err != nil {
		diags.AddError(
			"Error Creating Files Style",
			"Could not marshal logo to JSON: "+err.Error(),
		)
	}
	state.Logo = types.StringValue(string(respLogo))
	state.LogoClickHref = types.StringValue(style.LogoClickHref)
	respThumbnail, err := json.Marshal(style.Thumbnail)
	if err != nil {
		diags.AddError(
			"Error Creating Files Style",
			"Could not marshal thumbnail to JSON: "+err.Error(),
		)
	}
	state.Thumbnail = types.StringValue(string(respThumbnail))

	return
}
