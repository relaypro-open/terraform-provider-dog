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
	zoneDataSource struct {
		p dogProvider
	}

	ZoneList []Zone

	Zone struct {
		ID            types.String `tfsdk:"id"`
		IPv4Addresses []string     `tfsdk:"ipv4_addresses"`
		IPv6Addresses []string     `tfsdk:"ipv6_addresses"`
		Name          types.String `tfsdk:"name"`
	}
)

var (
	_ datasource.DataSource = (*zoneDataSource)(nil)
)

func NewZoneDataSource() datasource.DataSource {
	return &zoneDataSource{}
}

func (*zoneDataSource) Metadata(ctx context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_zone"
}

func (*zoneDataSource) GetSchema(ctx context.Context) (tfsdk.Schema, diag.Diagnostics) {
	return tfsdk.Schema{
		// This description is used by the documentation generator and the language server.
		MarkdownDescription: "Zone data source",

		Attributes: map[string]tfsdk.Attribute{
			// This description is used by the documentation generator and the language server.
			"ipv4_addresses": {
				MarkdownDescription: "List of Ipv4 Addresses",
				Optional:            true,
				Type: types.ListType{
					ElemType: types.StringType,
				},
			},
			"ipv6_addresses": {
				MarkdownDescription: "List of Ipv6 Addresses",
				Optional:            true,
				Type: types.ListType{
					ElemType: types.StringType,
				},
			},
			"name": {
				MarkdownDescription: "Zone name",
				Optional:            true,
				Type:                types.StringType,
			},
			"id": {
				Optional:            true,
				MarkdownDescription: "Zone identifier",
				Type: types.StringType,
			},
		},
	}, nil
}

func (d *zoneDataSource) Configure(ctx context.Context, req datasource.ConfigureRequest, resp *datasource.ConfigureResponse) {
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

type zoneDataSourceData struct {
	ApiToken types.String `tfsdk:"api_token"`
	Id       types.String `tfsdk:"id"`
}

//type zoneDataSource struct {
//	provider provider
//}

func (d *zoneDataSource) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
	var state Zone
	var zoneName string

	req.Config.GetAttribute(ctx, path.Root("name"), &zoneName)

	res, statusCode, err := d.p.dog.GetZones(nil)
	if (statusCode < 200 || statusCode > 299) && statusCode != 404 {
		resp.Diagnostics.AddError("Client Unsuccesful", fmt.Sprintf("Status Code: %d", statusCode))
	}
	if err != nil {
		resp.Diagnostics.AddError("Client Error", fmt.Sprintf("Unable to read zones, got error: %s", err))
	}
	if resp.Diagnostics.HasError() {
		return
	}

	var filteredZonesName []api.Zone
	if zoneName != "" {
		filteredZonesName = goterators.Filter(res, func(zone api.Zone) bool {
			return zone.Name == zoneName
		})
	} else {
		filteredZonesName = res
	}

	filteredZones := filteredZonesName

	if filteredZones == nil {
		resp.Diagnostics.AddError("Data Error", fmt.Sprintf("dog_zone data source returned no results."))
	} 
	if len(filteredZones) > 1 {
		resp.Diagnostics.AddError("Data Error", fmt.Sprintf("dog_zone data source returned more than one result."))
	}
	if resp.Diagnostics.HasError() {
		return
	}

	zone := filteredZones[0] 
	// Set state
	state = ApiToZone(zone)
	diags := resp.State.Set(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}
}
