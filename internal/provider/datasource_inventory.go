package dog

import (
	"context"
	"fmt"

	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/tfsdk"
	"github.com/hashicorp/terraform-plugin-framework/types"
	api "github.com/relaypro-open/dog_api_golang/api"
)

type (
	inventoryDataSource struct {
		p dogProvider
	}

	InventoryList []Inventory

	Inventory struct {
		ID     types.String      `tfsdk:"id"`
		Name   types.String      `tfsdk:"name"`
		Groups []*InventoryGroup `tfsdk:"groups"`
	}

	InventoryGroup struct {
		Name  types.String           `tfsdk:"name"`
		Vars  map[string]string `tfsdk:"vars"`
		Hosts map[string]map[string]string `tfsdk:"hosts"`
	}
)

var (
	_ datasource.DataSource = (*inventoryDataSource)(nil)
)

func NewInventoryDataSource() datasource.DataSource {
	return &inventoryDataSource{}
}

func (*inventoryDataSource) Metadata(ctx context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_inventory"
}

func (*inventoryDataSource) GetSchema(ctx context.Context) (tfsdk.Schema, diag.Diagnostics) {
	return tfsdk.Schema{
		// This description is used by the documentation generator and the language server.
		MarkdownDescription: "Inventory data source",

		Attributes: map[string]tfsdk.Attribute{
			"api_token": {
				MarkdownDescription: "Inventory configurable attribute",
				Optional:            true,
				Type:                types.StringType,
			},
			"id": {
				MarkdownDescription: "Inventory identifier",
				Type:                types.StringType,
				Computed:            true,
			},
		},
	}, nil
}

func (d *inventoryDataSource) Configure(ctx context.Context, req datasource.ConfigureRequest, resp *datasource.ConfigureResponse) {
	// Prevent panic if the provider has not been configured.
	if req.ProviderData == nil {
		return
	}

	client, ok := req.ProviderData.(*api.Client)

	if !ok {
		resp.Diagnostics.AddError(
			"Unexpected Data Source Configure Type",
			fmt.Sprintf("Expected *dog.Client, got: %T. Please report this issue to the provider developers.", req.ProviderData),
		)

		return
	}

	d.p.dog = client
}

type inventoryDataSourceData struct {
	ApiToken types.String `tfsdk:"api_token"`
	Id       types.String `tfsdk:"id"`
}

//type inventoryDataSource struct {
//	provider provider
//}

func (d *inventoryDataSource) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
	var state InventoryList

	res, statusCode, err := d.p.dog.GetInventories(nil)
	if (statusCode < 200 || statusCode > 299) && statusCode != 404 {
		resp.Diagnostics.AddError("Client Unsuccesful", fmt.Sprintf("Status Code: %d", statusCode))
	}
	if err != nil {
		resp.Diagnostics.AddError("Client Error", fmt.Sprintf("Unable to read inventories, got error: %s", err))
	}
	if resp.Diagnostics.HasError() {
		return
	}

	// Set state
	for _, api_inventory := range res {
		inventory := ApiToInventory(api_inventory)
		state = append(state, inventory)
	}
	diags := resp.State.Set(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}
}
