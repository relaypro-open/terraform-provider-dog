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
	profileDataSource struct {
		p dogProvider
	}

	ProfileList []Profile

	Profile struct {
		ID      types.String `tfsdk:"id"`
		Name    types.String `tfsdk:"name"`
		Version types.String `tfsdk:"version"`
	}
)

var (
	_ datasource.DataSource = (*profileDataSource)(nil)
)

func NewProfileDataSource() datasource.DataSource {
	return &profileDataSource{}
}

func (*profileDataSource) Metadata(ctx context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_profile"
}

func (*profileDataSource) GetSchema(ctx context.Context) (tfsdk.Schema, diag.Diagnostics) {
	return tfsdk.Schema{
		// This description is used by the documentation generator and the language server.
		MarkdownDescription: "Profile data source",

		Attributes: map[string]tfsdk.Attribute{
			"name": {
				MarkdownDescription: "Profile name",
				Optional:            true,
				Type:                types.StringType,
			},
			"version": {
				MarkdownDescription: "Profile version",
				Optional:            true,
				Type:                types.StringType,
			},
			"id": {
				Optional:            true,
				MarkdownDescription: "Profile identifier",
				Type: types.StringType,
			},
		},
	}, nil
}

func (d *profileDataSource) Configure(ctx context.Context, req datasource.ConfigureRequest, resp *datasource.ConfigureResponse) {
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

type profileDataSourceData struct {
	ApiToken types.String `tfsdk:"api_token"`
	Id       types.String `tfsdk:"id"`
}

//type profileDataSource struct {
//	provider provider
//}

func (d *profileDataSource) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
	var state ProfileList

	res, statusCode, err := d.p.dog.GetProfiles(nil)
	if (statusCode < 200 || statusCode > 299) && statusCode != 404 {
		resp.Diagnostics.AddError("Client Unsuccesful", fmt.Sprintf("Status Code: %d", statusCode))
	}
	if err != nil {
		resp.Diagnostics.AddError("Client Error", fmt.Sprintf("Unable to read profiles, got error: %s", err))
	}
	if resp.Diagnostics.HasError() {
		return
	}

	// Set state
	for _, api_profile := range res {
		profile := ApiToProfile(api_profile)
		state = append(state, profile)
	}
	diags := resp.State.Set(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}
}
