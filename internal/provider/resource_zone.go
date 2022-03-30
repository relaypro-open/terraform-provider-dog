package dog

import (
	"context"
	"fmt"
	"log"

	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/tfsdk"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-go/tftypes"
	"github.com/hashicorp/terraform-plugin-log/tflog"
	api "github.com/relaypro-open/dog_api_golang/api"
)

//// stringListDefaultModifier is a plan modifier that sets a default value for a
//// types.StringType attribute when it is not configured. The attribute must be
//// marked as Optional and Computed. When setting the state during the resource
//// Create, Read, or Update methods, this default value must also be included or
//// the Terraform CLI will generate an error.
//type stringListDefaultModifier struct {
//	Default string
//}
//
//// Description returns a plain text description of the validator's behavior, suitable for a practitioner to understand its impact.
//func (m stringListDefaultModifier) Description(ctx context.Context) string {
//	return fmt.Sprintf("If value is not configured, defaults to %s", m.Default)
//}
//
//// MarkdownDescription returns a markdown formatted description of the validator's behavior, suitable for a practitioner to understand its impact.
//func (m stringListDefaultModifier) MarkdownDescription(ctx context.Context) string {
//	return fmt.Sprintf("If value is not configured, defaults to `%s`", m.Default)
//}
//
//// Modify runs the logic of the plan modifier.
//// Access to the configuration, plan, and state is available in `req`, while
//// `resp` contains fields for updating the planned value, triggering resource
//// replacement, and returning diagnostics.
//func (m stringListDefaultModifier) Modify(ctx context.Context, req tfsdk.ModifyAttributePlanRequest, resp *tfsdk.ModifyAttributePlanResponse) {
//	// types.String must be the attr.Value produced by the attr.Type in the schema for this attribute
//	// for generic plan modifiers, use
//	// https://pkg.go.dev/github.com/hashicorp/terraform-plugin-framework/tfsdk#ConvertValue
//	// to convert into a known type.
//	var str types.String
//	diags := tfsdk.ValueAs(ctx, req.AttributePlan, &str)
//	resp.Diagnostics.Append(diags...)
//	if diags.HasError() {
//		return
//	}
//
//	if !str.Null {
//		return
//	}
//
//	resp.AttributePlan = types.List{
//		Elem: []types.String{Value: m.Default},
//	}
//}
//
//func stringListDefault(defaultValue string) stringListDefaultModifier {
//	return stringListDefaultModifier{
//		Default: defaultValue,
//	}
//}

type zoneResourceType struct{}

func (t zoneResourceType) GetSchema(ctx context.Context) (tfsdk.Schema, diag.Diagnostics) {
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
					tfsdk.UseStateForUnknown(),
				},
				Type: types.StringType,
			},
		},
	}, nil
}

func (t zoneResourceType) NewResource(ctx context.Context, in tfsdk.Provider) (tfsdk.Resource, diag.Diagnostics) {
	provider, diags := convertProviderType(in)

	return zoneResource{
		provider: provider,
	}, diags
}

type zoneResourceData struct {
	ID            types.String `tfsdk:"id"`
	IPv4Addresses []string     `tfsdk:"ipv4_addresses"`
	IPv6Addresses []string     `tfsdk:"ipv6_addresses"`
	Name          string       `tfsdk:"name"`
}

type zoneResource struct {
	provider provider
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

func (r zoneResource) Create(ctx context.Context, req tfsdk.CreateResourceRequest, resp *tfsdk.CreateResourceResponse) {
	tflog.Debug(ctx, "Create 1\n")
	var state Zone

	var plan zoneResourceData
	diags := req.Plan.Get(ctx, &plan)
	tflog.Debug(ctx, fmt.Sprintf("plan: %+v\n", plan))
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}
	newZone := ZoneToCreateRequest(plan)
	tflog.Debug(ctx, fmt.Sprintf("ZZZZZZZZZZZZZZZZZZZZZZZZZ NewZone: %+v\n", newZone))
	zone, statusCode, err := r.provider.client.CreateZone(newZone, nil)
	log.Printf(fmt.Sprintf("ZZZZZZZZZZZZZZZZZZZZZZZZZZZZZZ zone: %+v\n", zone))
	tflog.Debug(ctx, fmt.Sprintf("ZZZZZZZZZZZZZZZZZZZZZZZZZ zone: %+v\n", zone))
	state = ApiToZone(zone)
	if err != nil {
		resp.Diagnostics.AddError("Client Error", fmt.Sprintf("Unable to create zone, got error: %s", err))
		return
	}
	if statusCode < 200 || statusCode > 299 {
		resp.Diagnostics.AddError("Client Unsuccesful", fmt.Sprintf("Status Code: %d", statusCode))
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

func (r zoneResource) Read(ctx context.Context, req tfsdk.ReadResourceRequest, resp *tfsdk.ReadResourceResponse) {
	tflog.Debug(ctx, "Read 1\n")
	var state Zone

	diags := req.State.Get(ctx, &state)
	resp.Diagnostics.Append(diags...)

	if resp.Diagnostics.HasError() {
		return
	}

	zoneID := state.ID.Value

	zone, statusCode, err := r.provider.client.GetZone(zoneID, nil)
	if (statusCode < 200 || statusCode > 299) && statusCode != 404 {
		resp.Diagnostics.AddError("Client Unsuccesful", fmt.Sprintf("Status Code: %d", statusCode))
		return
	}
	if err != nil {
		resp.Diagnostics.AddError("Client Error", fmt.Sprintf("Unable to read zone, got error: %s", err))
		return
	}
	state = ApiToZone(zone)
	diags = resp.State.Set(ctx, &state)
	resp.Diagnostics.Append(diags...)
}

func (r zoneResource) Update(ctx context.Context, req tfsdk.UpdateResourceRequest, resp *tfsdk.UpdateResourceResponse) {
	tflog.Debug(ctx, "Update 1\n")
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
	//	resp.Diagnostics.AddError("Client Error", fmt.Sprintf("client: %+v\n", r.provider.client))
	if resp.Diagnostics.HasError() {
		return
	}
	newZone := ZoneToUpdateRequest(plan)
	zone, statusCode, err := r.provider.client.UpdateZone(zoneID, newZone, nil)
	log.Printf(fmt.Sprintf("zone: %+v\n", zone))
	tflog.Debug(ctx, fmt.Sprintf("zone: %+v\n", zone))
	state = ApiToZone(zone)
	if err != nil {
		resp.Diagnostics.AddError("Client Error", fmt.Sprintf("Unable to create zone, got error: %s", err))
		return
	}
	if statusCode < 200 || statusCode > 299 {
		resp.Diagnostics.AddError("Client Unsuccesful", fmt.Sprintf("Status Code: %d", statusCode))
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

func (r zoneResource) Delete(ctx context.Context, req tfsdk.DeleteResourceRequest, resp *tfsdk.DeleteResourceResponse) {
	tflog.Debug(ctx, "Delete 1\n")
	var state Zone

	diags := req.State.Get(ctx, &state)
	resp.Diagnostics.Append(diags...)

	if resp.Diagnostics.HasError() {
		return
	}

	zoneID := state.ID.Value
	zone, statusCode, err := r.provider.client.DeleteZone(zoneID, nil)
	if statusCode < 200 || statusCode > 299 {
		resp.Diagnostics.AddError("Client Unsuccesful", fmt.Sprintf("Status Code: %d", statusCode))
		return
	}
	if err != nil {
		resp.Diagnostics.AddError("Client Error", fmt.Sprintf("Unable to read zone, got error: %s", err))
		return
	}
	tflog.Debug(ctx, fmt.Sprintf("Zone deleted: %+v\n", zone))

	diags = resp.State.Set(ctx, &state)
	resp.Diagnostics.Append(diags...)
	resp.State.RemoveResource(ctx)
}

func (r zoneResource) ImportState(ctx context.Context, req tfsdk.ImportResourceStateRequest, resp *tfsdk.ImportResourceStateResponse) {
	tfsdk.ResourceImportStatePassthroughID(ctx, tftypes.NewAttributePath().WithAttributeName("id"), req, resp)
}
