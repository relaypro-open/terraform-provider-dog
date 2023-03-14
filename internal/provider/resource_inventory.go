package dog

import (
	"context"
	"fmt"
	"log"

	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/tfsdk"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-log/tflog"
	api "github.com/relaypro-open/dog_api_golang/api"
	"golang.org/x/exp/slices"
)

type (
	inventoryResource struct {
		p dogProvider
	}
)

var (
	_ resource.Resource                = (*inventoryResource)(nil)
	_ resource.ResourceWithImportState = (*inventoryResource)(nil)
)

func NewInventoryResource() resource.Resource {
	return &inventoryResource{}
}

func (*inventoryResource) Metadata(ctx context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_inventory"
}

func (*inventoryResource) GetSchema(ctx context.Context) (tfsdk.Schema, diag.Diagnostics) {
	return tfsdk.Schema{
		Attributes: map[string]tfsdk.Attribute{
			// This description is used by the documentation generator and the language server.
			"groups": {
				MarkdownDescription: "List of inventory groups",
				Required:            true,
				Attributes: tfsdk.ListNestedAttributes(map[string]tfsdk.Attribute{
					"name": {
						MarkdownDescription: "inventory group name",
						Required:            true,
						Type:                types.StringType,
					},
					"vars": {
						MarkdownDescription: "Arbitrary collection of variables used for inventory",
						Required:            true,
						Type:                types.MapType{ElemType: types.StringType},
					},
					"hosts": {
						MarkdownDescription: "Arbitrary collection of hosts used for inventory",
						Required:            true,
						Type:                types.MapType{ElemType: types.MapType{ElemType: types.StringType}},
					},
				}),
			},
			"name": {
				MarkdownDescription: "Inventory name",
				Required:            true,
				Type:                types.StringType,
			},
			"id": {
				Computed:            true,
				MarkdownDescription: "Inventory identifier",
				PlanModifiers: tfsdk.AttributePlanModifiers{
					resource.UseStateForUnknown(),
				},
				Type: types.StringType,
			},
		},
	}, nil
}

func (r *inventoryResource) Configure(ctx context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {
	// Prevent panic if the provider has not been configured
	if req.ProviderData == nil {
		return
	}

	client, ok := req.ProviderData.(*api.Client)
	if !ok {
		resp.Diagnostics.AddError(
			"Unexpected Resource Configure Type",
			fmt.Sprintf("Expected *api.Client, got: %T. Please report this issue to the provider developers.", req.ProviderData),
		)

		return
	}

	r.p.dog = client
}

func (*inventoryResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	resource.ImportStatePassthroughID(ctx, path.Root("id"), req, resp)
}

type inventoryResourceData struct {
	ID     types.String      `tfsdk:"id"`
	Groups []*InventoryGroup `tfsdk:"groups"`
	Name   string            `tfsdk:"name"`
}

func InventoryToCreateRequest(plan inventoryResourceData) api.InventoryCreateRequest {
	newGroups := []*api.InventoryGroup{}
	for _, group := range plan.Groups {
		g := &api.InventoryGroup{
			Name:  group.Name.Value,
			Vars:  group.Vars,
			Hosts: group.Hosts,
		}
		newGroups = append(newGroups, g)
	}
	newInventory := api.InventoryCreateRequest{
		Groups: newGroups,
		Name:   plan.Name,
	}
	return newInventory
}

func InventoryToUpdateRequest(plan inventoryResourceData) api.InventoryUpdateRequest {
	newGroups := []*api.InventoryGroup{}
	for _, group := range plan.Groups {
		g := &api.InventoryGroup{
			Name:  group.Name.Value,
			Vars:  group.Vars,
			Hosts: group.Hosts,
		}
		newGroups = append(newGroups, g)
	}
	newInventory := api.InventoryUpdateRequest{
		Groups: newGroups,
		Name:   plan.Name,
	}
	return newInventory
}

func ApiToInventory(inventory api.Inventory) Inventory {
	newGroups := []*InventoryGroup{}
	for _, group := range inventory.Groups {
		g := &InventoryGroup{
			Name:  types.String{Value: group.Name},
			Vars:  group.Vars,
			Hosts: group.Hosts,
		}
		newGroups = append(newGroups, g)
	}
	h := Inventory{
		ID:     types.String{Value: inventory.ID},
		Groups: newGroups,
		Name:   types.String{Value: inventory.Name},
	}

	return h
}

func (r *inventoryResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	var state Inventory

	var plan inventoryResourceData
	diags := req.Plan.Get(ctx, &plan)
	resp.Diagnostics.Append(diags...)
	//	resp.Diagnostics.AddError("Client Error", fmt.Sprintf("client: %+v\n", r.provider.client))
	if resp.Diagnostics.HasError() {
		return
	}

	newInventory := InventoryToCreateRequest(plan)
	log.Printf(fmt.Sprintf("r.p.dog: %+v\n", r.p.dog))
	inventory, statusCode, err := r.p.dog.CreateInventory(newInventory, nil)
	log.Printf(fmt.Sprintf("inventory: %+v\n", inventory))
	tflog.Trace(ctx, fmt.Sprintf("inventory: %+v\n", inventory))
	if err != nil {
		resp.Diagnostics.AddError("Client Error", fmt.Sprintf("Unable to create inventory, got error: %s", err))
	}
	if statusCode != 201 {
		resp.Diagnostics.AddError("Client Unsuccesful", fmt.Sprintf("Status Code: %d", statusCode))
	}
	if resp.Diagnostics.HasError() {
		return
	}
	state = ApiToInventory(inventory)

	plan.ID = state.ID

	// write logs using the tflog package
	// see https://pkg.go.dev/github.com/hashicorp/terraform-plugin-log/tflog
	// for more information
	tflog.Trace(ctx, "created a resource")

	diags = resp.State.Set(ctx, &state)
	resp.Diagnostics.Append(diags...)
}

func (r *inventoryResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	var state Inventory

	diags := req.State.Get(ctx, &state)
	resp.Diagnostics.Append(diags...)

	if resp.Diagnostics.HasError() {
		return
	}

	inventoryID := state.ID.Value

	log.Printf(fmt.Sprintf("r.p: %+v\n", r.p))
	log.Printf(fmt.Sprintf("r.p.dog: %+v\n", r.p.dog))
	inventory, statusCode, err := r.p.dog.GetInventory(inventoryID, nil)
	if statusCode != 200 {
		resp.Diagnostics.AddError("Client Unsuccesful", fmt.Sprintf("Status Code: %d", statusCode))
	}
	if err != nil {
		resp.Diagnostics.AddError("Client Error", fmt.Sprintf("Unable to read inventory, got error: %s", err))
	}
	if resp.Diagnostics.HasError() {
		return
	}
	state = ApiToInventory(inventory)
	diags = resp.State.Set(ctx, &state)
	resp.Diagnostics.Append(diags...)
}

func (r *inventoryResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	var state Inventory

	diags := req.State.Get(ctx, &state)
	resp.Diagnostics.Append(diags...)

	if resp.Diagnostics.HasError() {
		return
	}

	inventoryID := state.ID.Value

	var plan inventoryResourceData
	diags = req.Plan.Get(ctx, &plan)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	newInventory := InventoryToUpdateRequest(plan)
	inventory, statusCode, err := r.p.dog.UpdateInventory(inventoryID, newInventory, nil)
	log.Printf(fmt.Sprintf("inventory: %+v\n", inventory))
	tflog.Trace(ctx, fmt.Sprintf("inventory: %+v\n", inventory))
	state = ApiToInventory(inventory)
	if err != nil {
		resp.Diagnostics.AddError("Client Error", fmt.Sprintf("Unable to create inventory, got error: %s", err))
	}
	ok := []int{303, 200, 201}
	if slices.Contains(ok, statusCode) != true {
		resp.Diagnostics.AddError("Client Unsuccesful", fmt.Sprintf("Status Code: %d", statusCode))
	}
	if resp.Diagnostics.HasError() {
		return
	}

	plan.ID = state.ID

	// write logs using the tflog package
	// see https://pkg.go.dev/github.com/hashicorp/terraform-plugin-log/tflog
	// for more information
	tflog.Trace(ctx, "created a resource")

	diags = resp.State.Set(ctx, &state)
	resp.Diagnostics.Append(diags...)

}

func (r *inventoryResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	var state Inventory

	diags := req.State.Get(ctx, &state)
	resp.Diagnostics.Append(diags...)

	if resp.Diagnostics.HasError() {
		return
	}

	inventoryID := state.ID.Value
	inventory, statusCode, err := r.p.dog.DeleteInventory(inventoryID, nil)
	if statusCode != 204 {
		resp.Diagnostics.AddError("Client Unsuccesful", fmt.Sprintf("Status Code: %d", statusCode))
	}
	if err != nil {
		resp.Diagnostics.AddError("Client Error", fmt.Sprintf("Unable to read inventory, got error: %s", err))
	}
	if resp.Diagnostics.HasError() {
		return
	}
	tflog.Trace(ctx, fmt.Sprintf("inventory deleted: %+v\n", inventory))

	diags = resp.State.Set(ctx, &state)
	resp.Diagnostics.Append(diags...)
	resp.State.RemoveResource(ctx)
}
