package dog

import (
	"context"
	"fmt"
	"log"

	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/tfsdk"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-log/tflog"
	api "github.com/relaypro-open/dog_api_golang/api"
	"github.com/hashicorp/terraform-plugin-framework/path"
)


type (
	zoneResource struct {
		p dogProvider
	}
)

var (
	_ resource.Resource                = (*zoneResource)(nil)
	_ resource.ResourceWithImportState = (*zoneResource)(nil)
)

func NewZoneResource() resource.Resource {
	return &zoneResource{}
}

func (*zoneResource) Metadata(ctx context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_zone"
}


func (*zoneResource) GetSchema(ctx context.Context) (tfsdk.Schema, diag.Diagnostics) {
	return tfsdk.Schema{
		Attributes: map[string]tfsdk.Attribute{
			// This description is used by the documentation generator and the language server.
			"ipv4_addresses": {
				MarkdownDescription: "List of Ipv4 Addresses",
				Required:            true,
				Type: types.ListType{
					ElemType: types.StringType,
				},
			},
			"ipv6_addresses": {
				MarkdownDescription: "List of Ipv6 Addresses",
				Required:            true,
				Type: types.ListType{
					ElemType: types.StringType,
				},
			},
			"name": {
				MarkdownDescription: "Zone name",
				Required:            true,
				Type:                types.StringType,
			},
			"id": {
				Computed:            true,
				MarkdownDescription: "Zone identifier",
				PlanModifiers: tfsdk.AttributePlanModifiers{
					resource.UseStateForUnknown(),
				},
				Type: types.StringType,
			},
		},
	}, nil
}

func (r *zoneResource) Configure(ctx context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {
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


//func (r *zoneResource) Configure(ctx context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {
//	resp.ResourceData = r
//	// Prevent panic if the provider has not been configured
//	if req.ProviderData == nil {
//		return
//	}
//
//	client, ok := req.ProviderData.(*api.Client)
//	if !ok {
//		resp.Diagnostics.AddError(
//			"Unexpected Resource Configure Type",
//			fmt.Sprintf("Expected *dog.Client, got: %T. Please report this issue to the provider developers.", req.ProviderData),
//		)
//
//		return
//	}
//
//	r.dog = client
//}

func (*zoneResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	resource.ImportStatePassthroughID(ctx, path.Root("id"), req, resp)
}


type zoneResourceData struct {
	ID            types.String `tfsdk:"id"`
	IPv4Addresses []string     `tfsdk:"ipv4_addresses"`
	IPv6Addresses []string     `tfsdk:"ipv6_addresses"`
	Name          string       `tfsdk:"name"`
}

func ZoneToCreateRequest(plan zoneResourceData) api.ZoneCreateRequest {
	ipv4Addresses := []string{}
	for _, ipv4 := range plan.IPv4Addresses {
		ipv4Addresses = append(ipv4Addresses, ipv4)
	}
	ipv6Addresses := []string{}
	for _, ipv6 := range plan.IPv6Addresses {
		ipv6Addresses = append(ipv6Addresses, ipv6)
	}
	newZone := api.ZoneCreateRequest{
		IPv4Addresses: ipv4Addresses,
		IPv6Addresses: ipv6Addresses,
		Name:          plan.Name,
	}
	return newZone
}

func ZoneToUpdateRequest(plan zoneResourceData) api.ZoneUpdateRequest {
	ipv4Addresses := []string{}
	for _, ipv4 := range plan.IPv4Addresses {
		ipv4Addresses = append(ipv4Addresses, ipv4)
	}
	ipv6Addresses := []string{}
	for _, ipv6 := range plan.IPv6Addresses {
		ipv6Addresses = append(ipv6Addresses, ipv6)
	}
	newZone := api.ZoneUpdateRequest{
		IPv4Addresses: ipv4Addresses,
		IPv6Addresses: ipv6Addresses,
		Name:          plan.Name,
	}
	return newZone
}

func ApiToZone(zone api.Zone) Zone {
	newIpv4Addresses := []string{}
	for _, ipv4 := range zone.IPv4Addresses {
		newIpv4Addresses = append(newIpv4Addresses, ipv4)
	}
	newIpv6Addresses := []string{}
	for _, ipv6 := range zone.IPv6Addresses {
		newIpv6Addresses = append(newIpv6Addresses, ipv6)
	}
	h := Zone{
		ID:            types.String{Value: zone.ID},
		IPv4Addresses: newIpv4Addresses,
		IPv6Addresses: newIpv6Addresses,
		Name:          types.String{Value: zone.Name},
	}
	return h
}

func (r *zoneResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	var state Zone

	var plan zoneResourceData
	diags := req.Plan.Get(ctx, &plan)
	resp.Diagnostics.Append(diags...)
	//	resp.Diagnostics.AddError("Client Error", fmt.Sprintf("client: %+v\n", r.provider.client))
	if resp.Diagnostics.HasError() {
		return
	}

	newZone := ZoneToCreateRequest(plan)
	log.Printf(fmt.Sprintf("r.p.dog: %+v\n", r.p.dog))
	zone, statusCode, err := r.p.dog.CreateZone(newZone, nil)
	log.Printf(fmt.Sprintf("zone: %+v\n", zone))
	tflog.Trace(ctx, fmt.Sprintf("zone: %+v\n", zone))
	if err != nil {
		resp.Diagnostics.AddError("Client Error", fmt.Sprintf("Unable to create zone, got error: %s", err))
	}
	if statusCode != 201 {
		resp.Diagnostics.AddError("Client Unsuccesful", fmt.Sprintf("Status Code: %d", statusCode))
	}
	if resp.Diagnostics.HasError() {
		return
	}
	state = ApiToZone(zone)

	plan.ID = state.ID

	// write logs using the tflog package
	// see https://pkg.go.dev/github.com/hashicorp/terraform-plugin-log/tflog
	// for more information
	tflog.Trace(ctx, "created a resource")

	diags = resp.State.Set(ctx, &state)
	resp.Diagnostics.Append(diags...)
}

func (r *zoneResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	var state Zone

	diags := req.State.Get(ctx, &state)
	resp.Diagnostics.Append(diags...)

	if resp.Diagnostics.HasError() {
		return
	}

	zoneID := state.ID.Value

	log.Printf(fmt.Sprintf("r.p: %+v\n", r.p))
	log.Printf(fmt.Sprintf("r.p.dog: %+v\n", r.p.dog))
	zone, statusCode, err := r.p.dog.GetZone(zoneID, nil)
	if statusCode != 200 {
		resp.Diagnostics.AddError("Client Unsuccesful", fmt.Sprintf("Status Code: %d", statusCode))
	}
	if err != nil {
		resp.Diagnostics.AddError("Client Error", fmt.Sprintf("Unable to read zone, got error: %s", err))
	}
	if resp.Diagnostics.HasError() {
		return
	}
	state = ApiToZone(zone)
	diags = resp.State.Set(ctx, &state)
	resp.Diagnostics.Append(diags...)
}


func (r *zoneResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	var state Zone

	diags := req.State.Get(ctx, &state)
	resp.Diagnostics.Append(diags...)

	if resp.Diagnostics.HasError() {
		return
	}

	zoneID := state.ID.Value

	var plan zoneResourceData
	diags = req.Plan.Get(ctx, &plan)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	newZone := ZoneToUpdateRequest(plan)
	zone, statusCode, err := r.p.dog.UpdateZone(zoneID, newZone, nil)
	log.Printf(fmt.Sprintf("zone: %+v\n", zone))
	tflog.Trace(ctx, fmt.Sprintf("zone: %+v\n", zone))
	state = ApiToZone(zone)
	if err != nil {
		resp.Diagnostics.AddError("Client Error", fmt.Sprintf("Unable to create zone, got error: %s", err))
	}
	if statusCode != 303 {
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

func (r *zoneResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	var state Zone

	diags := req.State.Get(ctx, &state)
	resp.Diagnostics.Append(diags...)

	if resp.Diagnostics.HasError() {
		return
	}

	zoneID := state.ID.Value
	zone, statusCode, err := r.p.dog.DeleteZone(zoneID, nil)
	if statusCode != 204 {
		resp.Diagnostics.AddError("Client Unsuccesful", fmt.Sprintf("Status Code: %d", statusCode))
	}
	if err != nil {
		resp.Diagnostics.AddError("Client Error", fmt.Sprintf("Unable to read zone, got error: %s", err))
	}
	if resp.Diagnostics.HasError() {
		return
	}
	tflog.Trace(ctx, fmt.Sprintf("zone deleted: %+v\n", zone))

	diags = resp.State.Set(ctx, &state)
	resp.Diagnostics.Append(diags...)
	resp.State.RemoveResource(ctx)
}
