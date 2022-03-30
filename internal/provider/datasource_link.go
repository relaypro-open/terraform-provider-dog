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
		h := ApiToLink(ctx, link)
		resourceState.Links = append(resourceState.Links, h)
	}
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
