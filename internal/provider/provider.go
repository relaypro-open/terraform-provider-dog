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
		Api_Token      		types.String `tfsdk:"api_token"`
		API_Endpoint 		types.String `tfsdk:"api_endpoint"`
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
			"api_token": {
				MarkdownDescription: "API Key",
				Optional:            true,
				Type:                types.StringType,
				Sensitive:           true,
			},
		},
	}, nil
}

func (p *dogProvider) Configure(ctx context.Context, req provider.ConfigureRequest, resp *provider.ConfigureResponse) {
	var config dogProviderModel
	diags := req.Config.Get(ctx, &config)
	resp.Diagnostics.Append(diags...)

	if resp.Diagnostics.HasError() {
		return
	}

	// If practitioner provided a configuration value for any of the
	// attributes, it must be a known value.

	if config.API_Endpoint.Unknown {
		resp.Diagnostics.AddError(
			"Unknown Dog API Endpoint",
			"The provider cannot create the Dog API client as there is an unknown configuration value for the Dog API endpoint. "+
			"Either target apply the source of the value first, set the value statically in the configuration, or use the DOG_API_Endpoint environment variable.",
		)
	}

	if config.Api_Token.Unknown {
		resp.Diagnostics.AddError(
			"Unknown Dog API Key",
			"The provider cannot create the Dog API client as there is an unknown configuration value for the Dog API key. "+
			"Either target apply the source of the value first, set the value statically in the configuration, or use the DOG_API_TOKEN environment variable. ",
		)
	}

	if resp.Diagnostics.HasError() {
		return
	}

	// Default values to environment variables, but override
	// with Terraform configuration value if set.

	api_endpoint := os.Getenv("DOG_API_ENDPOINT")
	api_token := os.Getenv("DOG_API_TOKEN")

	if !config.API_Endpoint.IsNull() {
		api_endpoint = config.API_Endpoint.ValueString()
	}

	if !config.Api_Token.IsNull() {
		api_token = config.Api_Token.ValueString()
	}

	if (api_endpoint == "" || api_token == "") {
	        resp.Diagnostics.AddError(
	           "config values",
	    	fmt.Sprintf("config.Api_Token: %+v\n", config.Api_Token.Value)+
	    	fmt.Sprintf("config.API_Endpoint: %+v\n", config.API_Endpoint.Value)+
	    	fmt.Sprintf("api_endpoint: %+v\n", api_endpoint)+
	    	fmt.Sprintf("api_token: %+v\n", api_token),
	        )
	}


	// If any of the expected configurations are missing, return
	// errors with provider-specific guidance.

	if api_endpoint == "" {
		resp.Diagnostics.AddError(
			"Missing Dog API Endpoint",
			"The provider cannot create the Dog API client as there is a missing or empty value for the Dog API endpoint. "+
			"Set the API Endpoint value in the configuration or use the DOG_API_ENDPOINT environment variable. "+
			"If either is already set, ensure the value is not empty.",
		)
	}

	if api_token == "" {
		resp.Diagnostics.AddError(
			"Missing Dog API Key",
			"The provider cannot create the Dog API client as there is a missing or empty value for the Dog API key. "+
			"Set the API Key value in the configuration or use the DOG_API_TOKEN environment variable. "+
			"If either is already set, ensure the value is not empty.",
		)
	}

	if resp.Diagnostics.HasError() {
		return
	}

	c := api.NewClient(api_token, api_endpoint)

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

