package dog

import (
	"context"
	"fmt"

	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/tfsdk"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-framework/datasource"
	api "github.com/relaypro-open/dog_api_golang/api"
)

type (
	hostDataSource struct {
		p dogProvider
	}


	HostList []Host

	Host struct {
		Environment types.String `tfsdk:"environment"`
		Group       types.String `tfsdk:"group"`
		ID          types.String `tfsdk:"id"`
		HostKey     types.String `tfsdk:"hostkey"`
		Location    types.String `tfsdk:"location"`
		Name        types.String `tfsdk:"name"`
		Vars	    map[string]string `tfsdk:"vars"`
	}
)

var (
	_ datasource.DataSource = (*hostDataSource)(nil)
)

func NewHostDataSource() datasource.DataSource {
	return &hostDataSource{}
}


func (*hostDataSource) Metadata(ctx context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_host"
}

func (*hostDataSource) GetSchema(ctx context.Context) (tfsdk.Schema, diag.Diagnostics) {
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

func (d *hostDataSource) Configure(ctx context.Context, req datasource.ConfigureRequest, resp *datasource.ConfigureResponse) {
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

type hostDataSourceData struct {
	ApiKey types.String `tfsdk:"api_key"`
	Id     types.String `tfsdk:"id"`
}

func (d *hostDataSource) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
	var state HostList

	res, statusCode, err := d.p.dog.GetHosts(nil)
	if err != nil {
		resp.Diagnostics.AddError("Client Error", fmt.Sprintf("Unable to read hosts, got error: %s", err))
	}
	if (statusCode < 200 || statusCode > 299) && statusCode != 404 {
		resp.Diagnostics.AddError("Client Unsuccesful", fmt.Sprintf("Status Code: %d", statusCode))
	}
	if resp.Diagnostics.HasError() {
		return
	}

	// Set state
	for _, api_host := range res {
		host := ApiToHost(api_host)
		state = append(state, host)
	}
	diags := resp.State.Set(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}
}
