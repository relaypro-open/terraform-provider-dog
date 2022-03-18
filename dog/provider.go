package dog

import (
	"context"
	"fmt"
	"os"

	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/tfsdk"
	"github.com/hashicorp/terraform-plugin-framework/types"
	api "github.com/relaypro-open/dog_api_golang/api"
)

// provider satisfies the tfsdk.Provider interface and usually is included
// with all Resource and DataSource implementations.
type provider struct {
	// client can contain the upstream provider SDK or HTTP client used to
	// communicate with the upstream service. Resource and DataSource
	// implementations can then make calls using this client.
	//
	// TODO: If appropriate, implement upstream provider SDK or HTTP client.
	// client vendorsdk.HostClient
	client api.Client

	// configured is set to true at the end of the Configure method.
	// This can be used in Resource and DataSource implementations to verify
	// that the provider was previously configured.
	configured bool

	// version is set to the provider version on release, "dev" when the
	// provider is built and ran locally, and "test" when running acceptance
	// testing.
	version string
}

// providerData can be used to store data from the Terraform configuration.
type providerData struct {
	API_Key      types.String `tfsdk:"api_key"`
	API_Endpoint types.String `tfsdk:"api_endpoint"`
}

func (p *provider) Configure(ctx context.Context, req tfsdk.ConfigureProviderRequest, resp *tfsdk.ConfigureProviderResponse) {
	var data providerData
	diags := req.Config.Get(ctx, &data)
	resp.Diagnostics.Append(diags...)

	if resp.Diagnostics.HasError() {
		return
	}
	//resp.Diagnostics.AddError("data", fmt.Sprintf("data: %+v\n", data))

	// Configuration values are now available.
	// if data.Host.Null { /* ... */ }

	// If the upstream provider SDK or HTTP client requires configuration, such
	// as authentication or logging, this is a great opportunity to do so.
	if data.API_Key.Unknown {
		resp.Diagnostics.AddError(
			"Unknown Provider Configuration Value",
			"API_Key not defined. Either define a terraform variable, or set `DOG_API_KEY` environment variable",
		)
		return
	}
	if data.API_Endpoint.Unknown {
		resp.Diagnostics.AddError(
			"Unknown Provider Configuration Value",
			"API_Endpoint not defined. Either define a terraform variable, or set `DOG_API_ENDPOINT` environment variable",
		)
		return
	}

	if data.API_Key.Null {
		data.API_Key.Value = os.Getenv("DOG_API_KEY")
	}
	if data.API_Endpoint.Null {
		data.API_Endpoint.Value = os.Getenv("DOG_API_ENDPOINT")
	}
	//resp.Diagnostics.AddError("data.API_Key.Value", fmt.Sprintf("data.API_Key.Value: %+v\n", data.API_Key.Value))
	//resp.Diagnostics.AddError("data.API_Endpoint.Value", fmt.Sprintf("data.API_Endpoint.Value: %+v\n", data.API_Endpoint.Value))

	p.client = *api.NewClient(data.API_Key.Value, data.API_Endpoint.Value)
	//p.client = *api.NewClient(data.API_Key.Value)

	p.configured = true
}

func (p *provider) GetResources(ctx context.Context) (map[string]tfsdk.ResourceType, diag.Diagnostics) {
	return map[string]tfsdk.ResourceType{
		"dog_host":    hostResourceType{},
		"dog_group":   groupResourceType{},
		"dog_service": serviceResourceType{},
		"dog_zone":    zoneResourceType{},
		"dog_link":    linkResourceType{},
	}, nil
}

func (p *provider) GetDataSources(ctx context.Context) (map[string]tfsdk.DataSourceType, diag.Diagnostics) {
	return map[string]tfsdk.DataSourceType{
		"dog_host":    hostDataSourceType{},
		"dog_group":   groupDataSourceType{},
		"dog_service": serviceDataSourceType{},
		"dog_zone":    zoneDataSourceType{},
		"dog_link":    linkDataSourceType{},
	}, nil
}

func (p *provider) GetSchema(ctx context.Context) (tfsdk.Schema, diag.Diagnostics) {
	return tfsdk.Schema{
		Attributes: map[string]tfsdk.Attribute{
			"api_endpoint": {
				MarkdownDescription: "API endpoint URL",
				Optional:            true,
				Type:                types.StringType,
			},
			"api_key": {
				MarkdownDescription: "API Key",
				Optional:            true,
				Type:                types.StringType,
			},
		},
	}, nil
}

func New(version string) func() tfsdk.Provider {
	return func() tfsdk.Provider {
		return &provider{
			version: version,
		}
	}
}

// convertProviderType is a helper function for NewResource and NewDataSource
// implementations to associate the concrete provider type. Alternatively,
// this helper can be skipped and the provider type can be directly type
// asserted (e.g. provider: in.(*provider)), however using this can prevent
// potential panics.
func convertProviderType(in tfsdk.Provider) (provider, diag.Diagnostics) {
	var diags diag.Diagnostics

	p, ok := in.(*provider)

	if !ok {
		diags.AddError(
			"Unexpected Provider Instance Type",
			fmt.Sprintf("While creating the data source or resource, an unexpected provider type (%T) was received. This is always a bug in the provider code and should be reported to the provider developers.", p),
		)
		return provider{}, diags
	}

	if p == nil {
		diags.AddError(
			"Unexpected Provider Instance Type",
			"While creating the data source or resource, an unexpected empty provider instance was received. This is always a bug in the provider code and should be reported to the provider developers.",
		)
		return provider{}, diags
	}

	return *p, diags
}
