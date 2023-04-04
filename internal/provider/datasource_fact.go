package dog

import (
	"context"
	"fmt"

	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/tfsdk"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/ledongthuc/goterators"
	api "github.com/relaypro-open/dog_api_golang/api"
)

type (
	factDataSource struct {
		p dogProvider
	}

	FactList []Fact

	Fact struct {
		ID     types.String               `tfsdk:"id"`
		Name   types.String               `tfsdk:"name"`
		Groups map[string]*FactGroup `tfsdk:"groups"`
	}

	FactGroup struct {
		Vars     map[string]string            `tfsdk:"vars"`
		Hosts    map[string]map[string]string `tfsdk:"hosts"`
		Children []string                     `tfsdk:"children"`
	}
)

var (
	_ datasource.DataSource = (*factDataSource)(nil)
)

func NewFactDataSource() datasource.DataSource {
	return &factDataSource{}
}

func (*factDataSource) Metadata(ctx context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_fact"
}

func (*factDataSource) GetSchema(ctx context.Context) (tfsdk.Schema, diag.Diagnostics) {
	return tfsdk.Schema{
		// This description is used by the documentation generator and the language server.
		MarkdownDescription: "Fact data source",

		Attributes: map[string]tfsdk.Attribute{
			// This description is used by the documentation generator and the language server.
			"groups": {
				MarkdownDescription: "List of fact groups",
				Optional:            true,
				Attributes: tfsdk.MapNestedAttributes(map[string]tfsdk.Attribute{
					"vars": {
						MarkdownDescription: "Arbitrary collection of variables used for fact",
						Optional:            true,
						Type:                types.MapType{ElemType: types.StringType},
					},
					"hosts": {
						MarkdownDescription: "Arbitrary collection of hosts used for fact",
						Optional:            true,
						Type:                types.MapType{ElemType: types.MapType{ElemType: types.StringType}},
					},
					"children": {
						MarkdownDescription: "fact group children",
						Optional:            true,
						Type:                types.ListType{ElemType: types.StringType},
					},
				}),
			},
			"name": {
				MarkdownDescription: "Fact name",
				Optional:            true,
				Type:                types.StringType,
			},
			"id": {
				Computed:            true,
				MarkdownDescription: "Fact identifier",
				Type: types.StringType,
			},
		},
	}, nil
}

func (d *factDataSource) Configure(ctx context.Context, req datasource.ConfigureRequest, resp *datasource.ConfigureResponse) {
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

type factDataSourceData struct {
	ApiToken types.String `tfsdk:"api_token"`
	Id       types.String `tfsdk:"id"`
}

//type factDataSource struct {
//	provider provider
//}

func (d *factDataSource) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
	var state Fact
	var factName string

	req.Config.GetAttribute(ctx, path.Root("name"), &factName)

	res, statusCode, err := d.p.dog.GetFacts(nil)
	if (statusCode < 200 || statusCode > 299) && statusCode != 404 {
		resp.Diagnostics.AddError("Client Unsuccesful", fmt.Sprintf("Status Code: %d", statusCode))
	}
	if err != nil {
		resp.Diagnostics.AddError("Client Error", fmt.Sprintf("Unable to read facts, got error: %s", err))
	}
	if resp.Diagnostics.HasError() {
		return
	}

	var filteredFactsName []api.Fact
	if factName != "" {
		filteredFactsName = goterators.Filter(res, func(fact api.Fact) bool {
			return fact.Name == factName
		})
	} else {
		filteredFactsName = res
	}

	filteredFacts := filteredFactsName

	if filteredFacts == nil {
		resp.Diagnostics.AddError("Data Error", fmt.Sprintf("dog_fact data source returned no results."))
	} 
	if len(filteredFacts) > 1 {
		resp.Diagnostics.AddError("Data Error", fmt.Sprintf("dog_fact data source returned more than one result."))
	}
	if resp.Diagnostics.HasError() {
		return
	}

	fact := filteredFacts[0] 
	// Set state
	state = ApiToFact(fact)
	diags := resp.State.Set(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}
}
