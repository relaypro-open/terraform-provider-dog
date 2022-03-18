package dog

import (
	"context"
	"fmt"
	"log"

	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/tfsdk"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

type groupDataSourceType struct{}

func (t groupDataSourceType) GetSchema(ctx context.Context) (tfsdk.Schema, diag.Diagnostics) {
	return tfsdk.Schema{
		// This description is used by the documentation generator and the language server.
		MarkdownDescription: "Group data source",

		Attributes: map[string]tfsdk.Attribute{
			"api_key": {
				MarkdownDescription: "Group configurable attribute",
				Optional:            true,
				Type:                types.StringType,
			},
			"id": {
				MarkdownDescription: "Group identifier",
				Type:                types.StringType,
				Computed:            true,
			},
		},
	}, nil
}

func (t groupDataSourceType) NewDataSource(ctx context.Context, in tfsdk.Provider) (tfsdk.DataSource, diag.Diagnostics) {
	provider, diags := convertProviderType(in)

	return groupDataSource{
		provider: provider,
	}, diags
}

type groupDataSourceData struct {
	ApiKey types.String `tfsdk:"api_key"`
	Id     types.String `tfsdk:"id"`
}

type groupDataSource struct {
	provider provider
}

func (d groupDataSource) Read(ctx context.Context, req tfsdk.ReadDataSourceRequest, resp *tfsdk.ReadDataSourceResponse) {
	var data groupDataSourceData

	var resourceState struct {
		Groups GroupList `tfsdk:"groups"`
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
	groups, statusCode, err := d.provider.client.GetGroups(nil)
	for _, group := range groups {
		h := Group{
			Description:    types.String{Value: group.Description},
			ID:             types.String{Value: group.ID},
			Name:           types.String{Value: group.Name},
			ProfileName:    types.String{Value: group.ProfileName},
			ProfileVersion: types.String{Value: group.ProfileVersion},
		}
		resourceState.Groups = append(resourceState.Groups, h)
	}
	if statusCode < 200 && statusCode > 299 {
		resp.Diagnostics.AddError("Client Unsuccesful", fmt.Sprintf("Status Code: %d", statusCode))
		return
	}
	if err != nil {
		resp.Diagnostics.AddError("Client Error", fmt.Sprintf("Unable to read group, got error: %s", err))
		return
	}

	// For the purposes of this group code, hardcoding a response value to
	// save into the Terraform state.
	data.Id = types.String{Value: "Group.ID"}
	diags = resp.State.Set(ctx, &resourceState)
	resp.Diagnostics.Append(diags...)
}
