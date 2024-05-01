package dog

import (
	"context"
	"fmt"

	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/datasource/schema"
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
		Vars     *string                 `tfsdk:"vars"`
		Hosts    *string `tfsdk:"hosts"`
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

func (*factDataSource) Schema(ctx context.Context, req datasource.SchemaRequest, resp *datasource.SchemaResponse) {
	resp.Schema = schema.Schema{
		// This description is used by the documentation generator and the language server.

		Attributes: map[string]schema.Attribute{
			"groups": schema.MapNestedAttribute{
				NestedObject: schema.NestedAttributeObject{
					Attributes: map[string]schema.Attribute{
						"vars": schema.StringAttribute{
							MarkdownDescription: "json string of vars",
							Required:            true,
						},
						"hosts": schema.StringAttribute{
							MarkdownDescription: "json string of hosts",
							Required:            true,
						},
						//"hosts": schema.MapAttribute{
						//	Required:            true,
						//	ElementType:         types.MapType{ElemType: types.StringType},
						//},
						"children": schema.ListAttribute{
							Required:            true,
							ElementType:         types.StringType,
						},
					},
				},
				Optional: true,
			},
			"name": schema.StringAttribute{
				Optional:            true,
			},
			"id": schema.StringAttribute{
				Optional:            true,
				Computed: true,
			},
		},
	}
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

	res, statusCode, err := d.p.dog.GetFactsEncode(nil)
	if (statusCode < 200 || statusCode > 299) && statusCode != 404 {
		resp.Diagnostics.AddError("Client Unsuccessful", fmt.Sprintf("Status Code: %d", statusCode))
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
