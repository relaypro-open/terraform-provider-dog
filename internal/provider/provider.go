package dog

import (
	"context"
	"os"
	"log"
	"fmt"

	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/tfsdk"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/provider"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	api "github.com/relaypro-open/dog_api_golang/api"
)

type (
	dogProvider struct {
		dog *api.Client
		configured bool

		version string
	}

	dogProviderModel struct {
		API_Key      types.String `tfsdk:"api_key"`
		API_Endpoint types.String `tfsdk:"api_endpoint"`
	}
)


var (
	_ provider.Provider             = (*dogProvider)(nil)
	_ provider.ProviderWithMetadata = (*dogProvider)(nil)
)

func New(version string) func() provider.Provider {
	return func() provider.Provider {
		return &dogProvider{
			version: version,
		}
	}
}

func (p *dogProvider) Metadata(_ context.Context, req provider.MetadataRequest, resp *provider.MetadataResponse) {
	resp.TypeName = "dog"
}

func (*dogProvider) GetSchema(ctx context.Context) (tfsdk.Schema, diag.Diagnostics) {
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
				Sensitive:           true,
			},
		},
	}, nil
}

func (p *dogProvider) Configure(ctx context.Context, req provider.ConfigureRequest, resp *provider.ConfigureResponse) {
	var data dogProviderModel
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

	c := api.NewClient(data.API_Key.Value, data.API_Endpoint.Value)
	//p.client = *api.NewClient(data.API_Key.Value)

	p.configured = true
	log.Printf(fmt.Sprintf("p.dog: %+v\n", p.dog))
	log.Printf(fmt.Sprintf("p.configured: %+v\n", p.configured))
	log.Printf(fmt.Sprintf("p.version: %+v\n", p.version))

	p.dog = c

	resp.DataSourceData = p.dog
	resp.ResourceData = p.dog
}


func (*dogProvider) Resources(ctx context.Context) []func() resource.Resource {
	return []func() resource.Resource{
		NewHostResource,
		NewGroupResource,
		NewServiceResource,
		NewZoneResource,
		NewLinkResource,
		NewProfileResource,
	}
}

func (*dogProvider) DataSources(ctx context.Context) []func() datasource.DataSource {
	return []func() datasource.DataSource{
		NewHostDataSource,
		NewGroupDataSource,
		NewServiceDataSource,
		NewZoneDataSource,
		NewLinkDataSource,
		NewProfileDataSource,
	}
}

