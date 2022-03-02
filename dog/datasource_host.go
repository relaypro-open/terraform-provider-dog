package dog

import (
	"context"
	"fmt"
	"log"

	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/tfsdk"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

type hostDataSourceType struct{}

func (t hostDataSourceType) GetSchema(ctx context.Context) (tfsdk.Schema, diag.Diagnostics) {
	return tfsdk.Schema{
		// This description is used by the documentation generator and the language server.
		MarkdownDescription: "Host data source",

		Attributes: map[string]tfsdk.Attribute{
			"api_key": {
				MarkdownDescription: "Host configurable attribute",
				Optional:            true,
				Type:                types.StringType,
			},
			"id": {
				MarkdownDescription: "Host identifier",
				Type:                types.StringType,
				Computed:            true,
			},
		},
	}, nil
}

func (t hostDataSourceType) NewDataSource(ctx context.Context, in tfsdk.Provider) (tfsdk.DataSource, diag.Diagnostics) {
	provider, diags := convertProviderType(in)

	return hostDataSource{
		provider: provider,
	}, diags
}

type hostDataSourceData struct {
	ApiKey types.String `tfsdk:"api_key"`
	Id     types.String `tfsdk:"id"`
}

type hostDataSource struct {
	provider provider
}

func (d hostDataSource) Read(ctx context.Context, req tfsdk.ReadDataSourceRequest, resp *tfsdk.ReadDataSourceResponse) {
	var data hostDataSourceData

	var resourceState struct {
		Hosts HostList `tfsdk:"hosts"`
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
	hosts, statusCode, err := d.provider.client.GetHosts(nil)
	for _, host := range hosts {
		h := Host{
			Active:      types.String{Value: host.Active},
			Environment: types.String{Value: host.Environment},
			Group:       types.String{Value: host.Group},
			ID:          types.String{Value: host.ID},
			HostKey:     types.String{Value: host.HostKey},
			Location:    types.String{Value: host.Location},
			Name:        types.String{Value: host.Name},
		}
		resourceState.Hosts = append(resourceState.Hosts, h)
	}
	if statusCode < 200 && statusCode > 299 {
		resp.Diagnostics.AddError("Client Unsuccesful", fmt.Sprintf("Status Code: %d", statusCode))
		return
	}
	if err != nil {
		resp.Diagnostics.AddError("Client Error", fmt.Sprintf("Unable to read host, got error: %s", err))
		return
	}

	// For the purposes of this host code, hardcoding a response value to
	// save into the Terraform state.
	data.Id = types.String{Value: "Host.ID"}
	diags = resp.State.Set(ctx, &resourceState)
	resp.Diagnostics.Append(diags...)
}
