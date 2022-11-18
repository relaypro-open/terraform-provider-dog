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
		API_Key      		types.String `tfsdk:"api_key"`
		API_Endpoint 		types.String `tfsdk:"api_endpoint"`
		API_Key_Variable_Name 	types.String `tfsdk:"api_key_variable_name"`
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
			"api_key_variable_name": {
				MarkdownDescription: "Name of ENVIRONMENT variable that contains api_key",
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
			"Unknown dog API Endpoint",
			"The provider cannot create the Dog API client as there is an unknown configuration value for the Dog API endpoint. "+
			"Either target apply the source of the value first, set the value statically in the configuration, or use the DOG_API_Endpoint environment variable.",
		)
	}

	if config.API_Key.Unknown && config.API_Key_Variable_Name.Unknown {
		resp.Diagnostics.AddError(
			"Unknown Dog API Key and API Key Variable Name",
			"The provider cannot create the Dog API client as there is an unknown configuration value for the Dog API key. "+
			"Either target apply the source of the value first, set the value statically in the configuration, or use the DOG_API_KEY environment variable. "+
			"Or set 'api_key_variable_name terraform variable and set the API Key as the value of that ENV variable",
		)
	}

	if resp.Diagnostics.HasError() {
		return
	}

	// Default values to environment variables, but override
	// with Terraform configuration value if set.

	api_endpoint := os.Getenv("DOG_API_ENDPOINT")
	api_key := os.Getenv("DOG_API_KEY")
	dog_qa_api_key := os.Getenv("DOG_QA_API_KEY")
	if os.Getenv(config.API_Key_Variable_Name.Value) != "" {
		api_key = os.Getenv(config.API_Key_Variable_Name.Value)
	}

	if !config.API_Endpoint.IsNull() {
		api_endpoint = config.API_Endpoint.ValueString()
	}

	if !config.API_Key.IsNull() {
		api_key = config.API_Key.ValueString()
	}

	if (api_endpoint == "" || api_key == "") {
	        resp.Diagnostics.AddError(
	           "config values",
	    	fmt.Sprintf("config.API_Key: %+v\n", config.API_Key.Value)+
	    	fmt.Sprintf("config.API_Endpoint: %+v\n", config.API_Endpoint.Value)+
	    	fmt.Sprintf("config.API_Key_Variable_Name: %+v\n", config.API_Key_Variable_Name)+
	    	fmt.Sprintf("api_endpoint: %+v\n", api_endpoint)+
	    	fmt.Sprintf("api_key: %+v\n", api_key)+
	    	fmt.Sprintf("dog_qa_api_key: %+v\n", dog_qa_api_key),
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

	if api_key == "" {
		resp.Diagnostics.AddError(
			"Missing Dog API Key",
			"The provider cannot create the Dog API client as there is a missing or empty value for the Dog API key. "+
			"Set the API Key value in the configuration or use the DOG_API_KEY environment variable. "+
			"Or set 'api_key_variable_name terraform variable and set the API Key as the value of that ENV variable"+
			"If either is already set, ensure the value is not empty.",
		)
	}

	if resp.Diagnostics.HasError() {
		return
	}

	c := api.NewClient(api_key, api_endpoint)

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

