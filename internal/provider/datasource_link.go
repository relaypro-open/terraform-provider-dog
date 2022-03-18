package dog

import (
	"context"
	"fmt"
	"log"

	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/tfsdk"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-log/tflog"
)

type linkDataSourceType struct{}

func (t linkDataSourceType) GetSchema(ctx context.Context) (tfsdk.Schema, diag.Diagnostics) {
	tflog.Debug(ctx, "GetSchema 1\n")
	return tfsdk.Schema{
		// This description is used by the documentation generator and the language server.
		MarkdownDescription: "Link data source",

		Attributes: map[string]tfsdk.Attribute{
			"api_key": {
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

func (t linkDataSourceType) NewDataSource(ctx context.Context, in tfsdk.Provider) (tfsdk.DataSource, diag.Diagnostics) {
	provider, diags := convertProviderType(in)

	return linkDataSource{
		provider: provider,
	}, diags
}

type linkDataSourceData struct {
	ApiKey types.String `tfsdk:"api_key"`
	Id     types.String `tfsdk:"id"`
}

type linkDataSource struct {
	provider provider
}

func (d linkDataSource) Read(ctx context.Context, req tfsdk.ReadDataSourceRequest, resp *tfsdk.ReadDataSourceResponse) {
	tflog.Debug(ctx, "Read 1\n")
	var data linkDataSourceData

	var resourceState struct {
		Links LinkList `tfsdk:"links"`
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
	links, statusCode, err := d.provider.client.GetLinks(nil)
	for _, link := range links {
		h := Link{
			ID:              types.String{Value: link.ID},
			AddressHandling: types.String{Value: link.AddressHandling},
			Connection: Connection{
				ApiPort:  types.Int64{Value: int64(link.Connection.ApiPort)},
				Host:     types.String{Value: link.Connection.Host},
				Password: types.String{Value: link.Connection.Password},
				Port:     types.Int64{Value: int64(link.Connection.Port)},
				SSLOptions: SSLOptions{
					CaCertFile:           types.String{Value: link.Connection.SSLOptions.CaCertFile},
					CertFile:             types.String{Value: link.Connection.SSLOptions.CaCertFile},
					FailIfNoPeerCert:     types.Bool{Value: link.Connection.SSLOptions.FailIfNoPeerCert},
					KeyFile:              types.String{Value: link.Connection.SSLOptions.KeyFile},
					ServerNameIndication: types.String{Value: link.Connection.SSLOptions.ServerNameIndication},
					Verify:               types.String{Value: link.Connection.SSLOptions.Verify},
				},
				User:        types.String{Value: link.Connection.User},
				VirtualHost: types.String{Value: link.Connection.VirtualHost},
			},
			ConnectionType: types.String{Value: link.ConnectionType},
			Direction:      types.String{Value: link.Direction},
			Enabled:        types.Bool{Value: link.Enabled},
			Name:           types.String{Value: link.Name},
		}
		resourceState.Links = append(resourceState.Links, h)
	}
	//Connection: types.Object{
	//	AttrTypes: map[string]attr.Type{
	//		"api_port":     types.Int64Type,
	//		"host":         types.StringType,
	//		"password":     types.StringType,
	//		"port":         types.Int64Type,
	//		"ssl_options":  types.ObjectType,
	//		"user":         types.StringType,
	//		"virtual_host": types.StringType,
	//	},
	//	Attrs: map[string]attr.Value{
	//		"api_port": types.Int64{Value: link.Connection.ApiPort},
	//		"host":     types.String{Value: link.Connection.Host},
	//		"password": types.String{Value: link.Connection.Password},
	//		"port":     types.Int64{Value: link.Connection.Port},
	//		"ssl_options": types.Object{
	//			AttrTypes: map[string]attr.Type{
	//				"cacertfile":             types.StringType,
	//				"certfile":               types.StringType,
	//				"fail_if_no_peer_cert":   types.BoolType,
	//				"keyfile":                types.StringType,
	//				"server_name_indication": types.StringType,
	//				"verify":                 types.StringType,
	//			},
	//			Attrs: map[string]attr.Value{
	//				"cacertfile":             types.String{Value: link.Connection.SSLOptions.CaCertFile},
	//				"certfile":               types.String{Value: link.Connection.SSLOptions.CertFile},
	//				"fail_if_no_peer_cert":   types.Bool{Value: link.Connection.SSLOptions.FailIfNoPeerCert},
	//				"keyfile":                types.String{Value: link.Connection.SSLOptions.KeyFile},
	//				"server_name_indication": types.String{Value: link.Connection.SSLOptions.ServerNameIndication},
	//				"verify":                 types.String{Value: link.Connection.SSLOptions.Verify},
	//			},
	//		},
	//		User:        types.String{Value: link.Connection.User},
	//		VirtualHost: types.String{Value: link.Connection.VirtualHost},
	//	},
	//},
	if statusCode < 200 && statusCode > 299 {
		resp.Diagnostics.AddError("Client Unsuccesful", fmt.Sprintf("Status Code: %d", statusCode))
		return
	}
	if err != nil {
		resp.Diagnostics.AddError("Client Error", fmt.Sprintf("Unable to read link, got error: %s", err))
		return
	}

	// For the purposes of this link code, hardcoding a response value to
	// save into the Terraform state.
	data.Id = types.String{Value: "Link.ID"}
	diags = resp.State.Set(ctx, &resourceState)
	resp.Diagnostics.Append(diags...)
}
