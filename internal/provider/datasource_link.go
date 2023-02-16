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
	linkDataSource struct {
		p dogProvider
	}
	LinkList []Link

	Link struct {
		ID              types.String `tfsdk:"id"`
		AddressHandling types.String `tfsdk:"address_handling"`
		Connection      *Connection  `tfsdk:"connection"`
		ConnectionType  types.String `tfsdk:"connection_type"`
		Direction       types.String `tfsdk:"direction"`
		Enabled         types.Bool   `tfsdk:"enabled"`
		Name            types.String `tfsdk:"name"`
	}
	Connection struct {
		ApiPort     types.Int64  `tfsdk:"api_port"`
		Host        types.String `tfsdk:"host"`
		Password    types.String `tfsdk:"password"`
		Port        types.Int64  `tfsdk:"port"`
		SSLOptions  *SSLOptions  `tfsdk:"ssl_options"`
		User        types.String `tfsdk:"user"`
		VirtualHost types.String `tfsdk:"virtual_host"`
	}

	SSLOptions struct {
		CaCertFile           types.String `tfsdk:"cacertfile"`
		CertFile             types.String `tfsdk:"certfile"`
		FailIfNoPeerCert     types.Bool   `tfsdk:"fail_if_no_peer_cert"`
		KeyFile              types.String `tfsdk:"keyfile"`
		ServerNameIndication types.String `tfsdk:"server_name_indication"`
		Verify               types.String `tfsdk:"verify"`
	}

)

var (
	_ datasource.DataSource = (*linkDataSource)(nil)
)

func NewLinkDataSource() datasource.DataSource {
	return &linkDataSource{}
}


func (*linkDataSource) Metadata(ctx context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_link"
}

func (*linkDataSource) GetSchema(ctx context.Context) (tfsdk.Schema, diag.Diagnostics) {
	return tfsdk.Schema{
		// This description is used by the documentation generator and the language server.
		MarkdownDescription: "Link data source",

		Attributes: map[string]tfsdk.Attribute{
			"api_token": {
				MarkdownDescription: "Link configurable attribute",
				Optional:            true,
				Type:                types.StringType,
			},
			"id": {
				MarkdownDescription: "Link identifier",
				Type:                types.StringType,
				Computed:            true,
			},
		},
	}, nil
}

func (d *linkDataSource) Configure(ctx context.Context, req datasource.ConfigureRequest, resp *datasource.ConfigureResponse) {
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

type linkDataSourceData struct {
	ApiToken types.String `tfsdk:"api_token"`
	Id     types.String `tfsdk:"id"`
}

//type linkDataSource struct {
//	provider provider
//}

func (d *linkDataSource) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
	var state LinkList

	res, statusCode, err := d.p.dog.GetLinks(nil)
	if (statusCode < 200 || statusCode > 299) && statusCode != 404 {
		resp.Diagnostics.AddError("Client Unsuccesful", fmt.Sprintf("Status Code: %d", statusCode))
	}
	if err != nil {
		resp.Diagnostics.AddError("Client Error", fmt.Sprintf("Unable to read links, got error: %s", err))
	}
	if resp.Diagnostics.HasError() {
		return
	}

	// Set state
	for _, api_link := range res {
		link := ApiToLink(api_link)
		state = append(state, link)
	}
	diags := resp.State.Set(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}
}
