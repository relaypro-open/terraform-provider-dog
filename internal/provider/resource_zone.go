package dog

import (
	"context"
	"fmt"
	"log"
	"regexp"

	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-framework/schema/validator"
	"github.com/hashicorp/terraform-plugin-framework-validators/stringvalidator"
	"github.com/hashicorp/terraform-plugin-framework-validators/listvalidator"
	"github.com/hashicorp/terraform-plugin-log/tflog"
	api "github.com/relaypro-open/dog_api_golang/api"
	"golang.org/x/exp/slices"
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

func (*zoneResource) Schema(ctx context.Context, req resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = schema.Schema{
		// This description is used by the documentation generator and the language server.
		MarkdownDescription: "Zone data source",

		Attributes: map[string]schema.Attribute{
			// This description is used by the documentation generator and the language server.
			"ipv4_addresses": schema.ListAttribute{
				MarkdownDescription: "List of Ipv4 Addresses",
				Optional:            true,
				ElementType: types.StringType,
				Validators: []validator.List{
					listvalidator.ValueStringsAre(stringvalidator.RegexMatches(
						regexp.MustCompile(`\b(?:(?:2(?:[0-4][0-9]|5[0-5])|[0-1]?[0-9]?[0-9])\.){3}(?:(?:2([0-4][0-9]|5[0-5])|[0-1]?[0-9]?[0-9]))\b`), "Must be valid IPv4 address"),
					),
				},
			},
			"ipv6_addresses": schema.ListAttribute{
				MarkdownDescription: "List of Ipv6 Addresses",
				Optional:            true,
				ElementType: types.StringType,
				Validators: []validator.List{
					listvalidator.ValueStringsAre(stringvalidator.RegexMatches(
						regexp.MustCompile(`((([0-9a-fA-F]{0,4})\:){2,7})([0-9a-fA-F]{0,4})`), "Must be valid IPv6 address"),
					),
				},
			},
			"name": schema.StringAttribute{
				MarkdownDescription: "Zone name",
				Optional:            true,
				Validators: []validator.String{
					stringvalidator.LengthBetween(1, 28),
					stringvalidator.RegexMatches(
						regexp.MustCompile(`^[A-Za-z0-9_.-](.*)$`),
						"must start with alphanumeric characters, %, _, ., -",
					),
				},
			},
			"id": schema.StringAttribute{
				Optional:            true,
				MarkdownDescription: "Zone identifier",
				Computed: true,
			},
		},
	}
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
		ID:            types.StringValue(zone.ID),
		IPv4Addresses: newIpv4Addresses,
		IPv6Addresses: newIpv6Addresses,
		Name:          types.StringValue(zone.Name),
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

	zoneID := state.ID.ValueString()

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

	zoneID := state.ID.ValueString()

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

func (r *zoneResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	var state Zone

	diags := req.State.Get(ctx, &state)
	resp.Diagnostics.Append(diags...)

	if resp.Diagnostics.HasError() {
		return
	}

	zoneID := state.ID.ValueString()
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
