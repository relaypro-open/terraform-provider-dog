package dog

import (
	"context"
	"fmt"

	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/tfsdk"
	"github.com/hashicorp/terraform-plugin-framework/types"
	api "github.com/relaypro-open/dog_api_golang/api"
)

type (
	serviceDataSource struct {
		p dogProvider
	}

	ServiceList []Service

	Service struct {
		ID       types.String    `tfsdk:"id"`
		Services []*PortProtocol `tfsdk:"services"`
		Name     types.String    `tfsdk:"name"`
		Version  types.Int64     `tfsdk:"version"`
	}

	Services []PortProtocol

	PortProtocol struct {
		Ports    []string     `tfsdk:"ports"`
		Protocol types.String `tfsdk:"protocol"`
	}
)

var (
	_ datasource.DataSource = (*serviceDataSource)(nil)
)

func NewServiceDataSource() datasource.DataSource {
	return &serviceDataSource{}
}

func (*serviceDataSource) Metadata(ctx context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_service"
}

func (*serviceDataSource) GetSchema(ctx context.Context) (tfsdk.Schema, diag.Diagnostics) {
	return tfsdk.Schema{
		// This description is used by the documentation generator and the language server.
		MarkdownDescription: "Service data source",

		Attributes: map[string]tfsdk.Attribute{
			// This description is used by the documentation generator and the language server.
			"services": {
				MarkdownDescription: "List of Services",
				Optional:            true,
				Attributes: tfsdk.ListNestedAttributes(map[string]tfsdk.Attribute{
					"protocol": {
						MarkdownDescription: "Service protocol",
						Required:            true,
						Type:                types.StringType,
					},
					"ports": {
						MarkdownDescription: "Service ports",
						Required:            true,
						Type: types.ListType{
							ElemType: types.StringType,
						},
					},
				}),
			},
			"name": {
				MarkdownDescription: "Service name",
				Optional:            true,
				Type:                types.StringType,
			},
			"version": {
				MarkdownDescription: "Service version",
				Optional:            true,
				Type:                types.Int64Type,
			},
			"id": {
				Optional:            true,
				MarkdownDescription: "Service identifier",
				Type: types.StringType,
			},
		},
	}, nil
}

func (d *serviceDataSource) Configure(ctx context.Context, req datasource.ConfigureRequest, resp *datasource.ConfigureResponse) {
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

type serviceDataSourceData struct {
	ApiToken types.String `tfsdk:"api_token"`
	Id       types.String `tfsdk:"id"`
}

//type serviceDataSource struct {
//	provider provider
//}

func (d *serviceDataSource) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
	var state ServiceList

	res, statusCode, err := d.p.dog.GetServices(nil)
	if (statusCode < 200 || statusCode > 299) && statusCode != 404 {
		resp.Diagnostics.AddError("Client Unsuccesful", fmt.Sprintf("Status Code: %d", statusCode))
	}
	if err != nil {
		resp.Diagnostics.AddError("Client Error", fmt.Sprintf("Unable to read services, got error: %s", err))
	}
	if resp.Diagnostics.HasError() {
		return
	}

	// Set state
	for _, api_service := range res {
		service := ApiToService(api_service)
		state = append(state, service)
	}
	diags := resp.State.Set(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}
}
