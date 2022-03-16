package dog

import (
	"context"
	"fmt"
	"log"

	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/tfsdk"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

type serviceDataSourceType struct{}

func (t serviceDataSourceType) GetSchema(ctx context.Context) (tfsdk.Schema, diag.Diagnostics) {
	return tfsdk.Schema{
		// This description is used by the documentation generator and the language server.
		MarkdownDescription: "Service data source",

		Attributes: map[string]tfsdk.Attribute{
			"api_key": {
				MarkdownDescription: "Service configurable attribute",
				Optional:            true,
				Type:                types.StringType,
			},
			"id": {
				MarkdownDescription: "Service identifier",
				Type:                types.StringType,
				Computed:            true,
			},
		},
	}, nil
}

func (t serviceDataSourceType) NewDataSource(ctx context.Context, in tfsdk.Provider) (tfsdk.DataSource, diag.Diagnostics) {
	provider, diags := convertProviderType(in)

	return serviceDataSource{
		provider: provider,
	}, diags
}

type serviceDataSourceData struct {
	ApiKey types.String `tfsdk:"api_key"`
	Id     types.String `tfsdk:"id"`
}

type serviceDataSource struct {
	provider provider
}

func (d serviceDataSource) Read(ctx context.Context, req tfsdk.ReadDataSourceRequest, resp *tfsdk.ReadDataSourceResponse) {
	var data serviceDataSourceData

	var resourceState struct {
		Services ServiceList `tfsdk:"services"`
	}

	diags := req.Config.Get(ctx, &data)
	resp.Diagnostics.Append(diags...)

	log.Printf("got here")

	if resp.Diagnostics.HasError() {
		return
	}

	log.Printf("got here")

	// If applicable, this is a great opportunity to initialize any necessary
	// provider client data and make a call using it.
	services, statusCode, err := d.provider.client.GetServices(nil)
	var h Service
	for _, service := range services {
		var s Services
		for _, port_protocol := range service.Services {
			var pp PortProtocol
			pp = PortProtocol{
				Ports:    port_protocol.Ports,
				Protocol: types.String{Value: port_protocol.Protocol},
			}
			s = append(s, pp)
		}
		h = Service{
			Created:  types.Int64{Value: int64(service.Created)},
			ID:       types.String{Value: service.ID},
			Services: s,
			Name:     types.String{Value: service.Name},
			Version:  types.Int64{Value: int64(service.Version)},
		}
		resourceState.Services = append(resourceState.Services, h)
	}
	if statusCode < 200 && statusCode > 299 {
		resp.Diagnostics.AddError("Client Unsuccesful", fmt.Sprintf("Status Code: %d", statusCode))
		return
	}
	if err != nil {
		resp.Diagnostics.AddError("Client Error", fmt.Sprintf("Unable to read service, got error: %s", err))
		return
	}

	// For the purposes of this service code, hardcoding a response value to
	// save into the Terraform state.
	data.Id = types.String{Value: "Service.ID"}
	diags = resp.State.Set(ctx, &resourceState)
	resp.Diagnostics.Append(diags...)
}
