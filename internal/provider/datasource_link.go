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
			// This description is used by the documentation generator and the language server.
			"address_handling": {
				MarkdownDescription: "Type of address handling",
				Optional:            true,
				Type:                types.StringType,
			},
			"connection": {
				MarkdownDescription: "Connection specification",
				Optional:            true,
				Attributes: tfsdk.SingleNestedAttributes(map[string]tfsdk.Attribute{
					"api_port": {
						Type:     types.Int64Type,
						Required: true,
					},
					"host": {
						Type:     types.StringType,
						Required: true,
					},
					"password": {
						Type:      types.StringType,
						Required:  true,
						Sensitive: true,
					},
					"port": {
						Type:     types.Int64Type,
						Required: true,
					},
					"ssl_options": {
						Required: true,
						Attributes: tfsdk.SingleNestedAttributes(map[string]tfsdk.Attribute{
							"cacertfile": {
								Type:     types.StringType,
								Required: true,
							},
							"certfile": {
								Type:     types.StringType,
								Required: true,
							},
							"fail_if_no_peer_cert": {
								Type:     types.BoolType,
								Required: true,
							},
							"keyfile": {
								Type:     types.StringType,
								Required: true,
							},
							"server_name_indication": {
								Type:     types.StringType,
								Required: true,
							},
							"verify": {
								Type:     types.StringType,
								Required: true,
							},
						}),
					},
					"user": {
						Type:     types.StringType,
						Required: true,
					},
					"virtual_host": {
						Type:     types.StringType,
						Required: true,
					},
				}),
			},
			"connection_type": {
				MarkdownDescription: "Connection type",
				Optional:            true,
				Type:                types.StringType,
			},
			"direction": {
				MarkdownDescription: "Connection direction",
				Optional:            true,
				Type:                types.StringType,
			},
			"enabled": {
				MarkdownDescription: "Connection enabled",
				Optional:            true,
				Type:                types.BoolType,
			},
			"name": {
				MarkdownDescription: "Link name",
				Optional:            true,
				Type:                types.StringType,
			},
			"id": {
				Optional:            true,
				MarkdownDescription: "Link identifier",
				Type: types.StringType,
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
	Id       types.String `tfsdk:"id"`
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
