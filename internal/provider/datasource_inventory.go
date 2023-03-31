package dog

import (
	"context"
	"fmt"

	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/tfsdk"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/ledongthuc/goterators"
	api "github.com/relaypro-open/dog_api_golang/api"
)

type (
	inventoryDataSource struct {
		p dogProvider
	}

	InventoryList []Inventory

	Inventory struct {
		ID     types.String               `tfsdk:"id"`
		Name   types.String               `tfsdk:"name"`
		Groups map[string]*InventoryGroup `tfsdk:"groups"`
	}

	InventoryGroup struct {
		Vars     map[string]string            `tfsdk:"vars"`
		Hosts    map[string]map[string]string `tfsdk:"hosts"`
		Children []string                     `tfsdk:"children"`
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
			// This description is used by the documentation generator and the language server.
			"groups": {
				MarkdownDescription: "List of inventory groups",
				Optional:            true,
				Attributes: tfsdk.MapNestedAttributes(map[string]tfsdk.Attribute{
					"vars": {
						MarkdownDescription: "Arbitrary collection of variables used for inventory",
						Optional:            true,
						Type:                types.MapType{ElemType: types.StringType},
					},
					"hosts": {
						MarkdownDescription: "Arbitrary collection of hosts used for inventory",
						Optional:            true,
						Type:                types.MapType{ElemType: types.MapType{ElemType: types.StringType}},
					},
					"children": {
						MarkdownDescription: "inventory group children",
						Optional:            true,
						Type:                types.ListType{ElemType: types.StringType},
					},
				}),
			},
			"name": {
				MarkdownDescription: "Inventory name",
				Optional:            true,
				Type:                types.StringType,
			},
			"id": {
				Computed:            true,
				MarkdownDescription: "Inventory identifier",
				Type: types.StringType,
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
	var state Inventory
	var inventoryName string

	req.Config.GetAttribute(ctx, path.Root("name"), &inventoryName)

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

	var filteredInventorysName []api.Inventory
	if inventoryName != "" {
		filteredInventorysName = goterators.Filter(res, func(inventory api.Inventory) bool {
			return inventory.Name == inventoryName
		})
	} else {
		filteredInventorysName = res
	}

	filteredInventorys := filteredInventorysName

	if filteredInventorys == nil {
		resp.Diagnostics.AddError("Data Error", fmt.Sprintf("dog_inventory data source returned no results."))
	} 
	if len(filteredInventorys) > 1 {
		resp.Diagnostics.AddError("Data Error", fmt.Sprintf("dog_inventory data source returned more than one result."))
	}
	if resp.Diagnostics.HasError() {
		return
	}

	inventory := filteredInventorys[0] 
	// Set state
	state = ApiToInventory(inventory)
	diags := resp.State.Set(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}
}
