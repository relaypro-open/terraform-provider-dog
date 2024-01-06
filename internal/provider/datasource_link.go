package dog

import (
	"context"
	"fmt"

	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/datasource/schema"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-log/tflog"
	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/ledongthuc/goterators"
	"github.com/davecgh/go-spew/spew"
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
		Connection      *Connection  `tfsdk:"dog_connection"`
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

func (*linkDataSource) Schema(ctx context.Context, req datasource.SchemaRequest, resp *datasource.SchemaResponse) {
	resp.Schema = schema.Schema{
		// This description is used by the documentation generator and the language server.
		MarkdownDescription: "Link data source",

		Attributes: map[string]schema.Attribute{
			// This description is used by the documentation generator and the language server.
			"address_handling": schema.StringAttribute{
				MarkdownDescription: "Type of address handling",
				Optional:            true,
			},
			"dog_connection": schema.SingleNestedAttribute{
				MarkdownDescription: "Connection specification",
				Optional:            true,
				Attributes: map[string]schema.Attribute{
					"api_port": schema.Int64Attribute{
						Required: true,
					},
					"host": schema.StringAttribute{
						Required: true,
					},
					"password": schema.StringAttribute{
						Required:  true,
						Sensitive: true,
					},
					"port": schema.Int64Attribute{
						Required: true,
					},
					"ssl_options": schema.SingleNestedAttribute{
						Required: true,
						Attributes: map[string]schema.Attribute{
							"cacertfile": schema.StringAttribute{
								Required: true,
							},
							"certfile": schema.StringAttribute{
								Required: true,
							},
							"fail_if_no_peer_cert": schema.BoolAttribute{
								Required: true,
							},
							"keyfile": schema.StringAttribute{
								Required: true,
							},
							"server_name_indication": schema.StringAttribute{
								Required: true,
							},
							"verify": schema.StringAttribute{
								Required: true,
							},
						},
					},
					"user": schema.StringAttribute{
						Required: true,
					},
					"virtual_host": schema.StringAttribute{
						Required: true,
					},
				},
			},
			"connection_type": schema.StringAttribute{
				MarkdownDescription: "Connection type",
				Optional:            true,
			},
			"direction": schema.StringAttribute{
				MarkdownDescription: "Connection direction",
				Optional:            true,
			},
			"enabled": schema.BoolAttribute{
				MarkdownDescription: "Connection enabled",
				Optional:            true,
			},
			"name": schema.StringAttribute{
				MarkdownDescription: "Link name",
				Optional:            true,
			},
			"id": schema.StringAttribute{
				Optional:            true,
				MarkdownDescription: "Link identifier",
				Computed: true,
			},
		},
	}
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
	var state Link
	var linkName string

	req.Config.GetAttribute(ctx, path.Root("name"), &linkName)

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

	var filteredLinksName []api.Link
	if linkName != "" {
		filteredLinksName = goterators.Filter(res, func(link api.Link) bool {
			return link.Name == linkName
		})
	} else {
		filteredLinksName = res
	}

	filteredLinks := filteredLinksName

	tflog.Debug(ctx, spew.Sprint("ZZZfilteredLinks: %#v", filteredLinks))
	if filteredLinks == nil {
		resp.Diagnostics.AddError("Data Error", fmt.Sprintf("dog_link data source returned no results."))
	} 
	if len(filteredLinks) > 1 {
		resp.Diagnostics.AddError("Data Error", fmt.Sprintf("dog_link data source returned more than one result."))
	}
	if resp.Diagnostics.HasError() {
		return
	}

	link := filteredLinks[0] 
	// Set state
	state = ApiToLink(link)
	diags := resp.State.Set(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}
}
