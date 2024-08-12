package provider

import (
	"context"
	"fmt"

	files_sdk "github.com/Files-com/files-sdk-go/v3"
	invoice "github.com/Files-com/files-sdk-go/v3/invoice"
	"github.com/Files-com/terraform-provider-files/lib"
	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/datasource/schema"
	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

var (
	_ datasource.DataSource              = &invoiceDataSource{}
	_ datasource.DataSourceWithConfigure = &invoiceDataSource{}
)

func NewInvoiceDataSource() datasource.DataSource {
	return &invoiceDataSource{}
}

type invoiceDataSource struct {
	client *invoice.Client
}

type invoiceDataSourceModel struct {
	Id                types.Int64   `tfsdk:"id"`
	Amount            types.String  `tfsdk:"amount"`
	Balance           types.String  `tfsdk:"balance"`
	CreatedAt         types.String  `tfsdk:"created_at"`
	Currency          types.String  `tfsdk:"currency"`
	DownloadUri       types.String  `tfsdk:"download_uri"`
	InvoiceLineItems  types.Dynamic `tfsdk:"invoice_line_items"`
	Method            types.String  `tfsdk:"method"`
	PaymentLineItems  types.Dynamic `tfsdk:"payment_line_items"`
	PaymentReversedAt types.String  `tfsdk:"payment_reversed_at"`
	PaymentType       types.String  `tfsdk:"payment_type"`
	SiteName          types.String  `tfsdk:"site_name"`
	Type              types.String  `tfsdk:"type"`
}

func (r *invoiceDataSource) Configure(_ context.Context, req datasource.ConfigureRequest, resp *datasource.ConfigureResponse) {
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

	r.client = &invoice.Client{Config: sdk_config}
}

func (r *invoiceDataSource) Metadata(_ context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_invoice"
}

func (r *invoiceDataSource) Schema(_ context.Context, _ datasource.SchemaRequest, resp *datasource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Description: "An AccountLineItem is a single line item in the accounting ledger for a billing account. These include payments and invoices.",
		Attributes: map[string]schema.Attribute{
			"id": schema.Int64Attribute{
				Description: "Line item Id",
				Required:    true,
			},
			"amount": schema.StringAttribute{
				Description: "Line item amount",
				Computed:    true,
			},
			"balance": schema.StringAttribute{
				Description: "Line item balance",
				Computed:    true,
			},
			"created_at": schema.StringAttribute{
				Description: "Line item created at",
				Computed:    true,
			},
			"currency": schema.StringAttribute{
				Description: "Line item currency",
				Computed:    true,
			},
			"download_uri": schema.StringAttribute{
				Description: "Line item download uri",
				Computed:    true,
			},
			"invoice_line_items": schema.DynamicAttribute{
				Description: "Associated invoice line items",
				Computed:    true,
			},
			"method": schema.StringAttribute{
				Description: "Line item payment method",
				Computed:    true,
			},
			"payment_line_items": schema.DynamicAttribute{
				Description: "Associated payment line items",
				Computed:    true,
			},
			"payment_reversed_at": schema.StringAttribute{
				Description: "Date/time payment was reversed if applicable",
				Computed:    true,
			},
			"payment_type": schema.StringAttribute{
				Description: "Type of payment if applicable",
				Computed:    true,
			},
			"site_name": schema.StringAttribute{
				Description: "Site name this line item is for",
				Computed:    true,
			},
			"type": schema.StringAttribute{
				Description: "Type of line item, either payment or invoice",
				Computed:    true,
			},
		},
	}
}

func (r *invoiceDataSource) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
	var data invoiceDataSourceModel
	diags := req.Config.Get(ctx, &data)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	paramsInvoiceFind := files_sdk.InvoiceFindParams{}
	paramsInvoiceFind.Id = data.Id.ValueInt64()

	invoice, err := r.client.Find(paramsInvoiceFind, files_sdk.WithContext(ctx))
	if err != nil {
		resp.Diagnostics.AddError(
			"Error Reading Files Invoice",
			"Could not read invoice id "+fmt.Sprint(data.Id.ValueInt64())+": "+err.Error(),
		)
		return
	}

	diags = r.populateDataSourceModel(ctx, invoice, &data)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	diags = resp.State.Set(ctx, data)
	resp.Diagnostics.Append(diags...)
}

func (r *invoiceDataSource) populateDataSourceModel(ctx context.Context, invoice files_sdk.AccountLineItem, state *invoiceDataSourceModel) (diags diag.Diagnostics) {
	var propDiags diag.Diagnostics

	state.Id = types.Int64Value(invoice.Id)
	state.Amount = types.StringValue(invoice.Amount)
	state.Balance = types.StringValue(invoice.Balance)
	if err := lib.TimeToStringType(ctx, path.Root("created_at"), invoice.CreatedAt, &state.CreatedAt); err != nil {
		diags.AddError(
			"Error Creating Files Invoice",
			"Could not convert state created_at to string: "+err.Error(),
		)
	}
	state.Currency = types.StringValue(invoice.Currency)
	state.DownloadUri = types.StringValue(invoice.DownloadUri)
	state.InvoiceLineItems, propDiags = lib.ToDynamic(ctx, path.Root("invoice_line_items"), invoice.InvoiceLineItems, state.InvoiceLineItems.UnderlyingValue())
	diags.Append(propDiags...)
	state.Method = types.StringValue(invoice.Method)
	state.PaymentLineItems, propDiags = lib.ToDynamic(ctx, path.Root("payment_line_items"), invoice.PaymentLineItems, state.PaymentLineItems.UnderlyingValue())
	diags.Append(propDiags...)
	if err := lib.TimeToStringType(ctx, path.Root("payment_reversed_at"), invoice.PaymentReversedAt, &state.PaymentReversedAt); err != nil {
		diags.AddError(
			"Error Creating Files Invoice",
			"Could not convert state payment_reversed_at to string: "+err.Error(),
		)
	}
	state.PaymentType = types.StringValue(invoice.PaymentType)
	state.SiteName = types.StringValue(invoice.SiteName)
	state.Type = types.StringValue(invoice.Type)

	return
}
