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
	groupDataSource struct {
		p dogProvider
	}

	GroupList []Group

	Group struct {
		Description    types.String `tfsdk:"description"`
		ID             types.String `tfsdk:"id"`
		Name           types.String `tfsdk:"name"`
		ProfileName    types.String `tfsdk:"profile_name"`
		ProfileVersion types.String `tfsdk:"profile_version"`
	}

)

var (
	_ datasource.DataSource = (*groupDataSource)(nil)
)

func NewGroupDataSource() datasource.DataSource {
	return &groupDataSource{}
}


func (*groupDataSource) Metadata(ctx context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_group"
}

func (*groupDataSource) GetSchema(ctx context.Context) (tfsdk.Schema, diag.Diagnostics) {
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
			"group_id": {
				Required:    true,
				Type:        types.StringType,
				Description: "The ID of the group.",
			},
		},
	}, nil
}

func (d *groupDataSource) Configure(ctx context.Context, req datasource.ConfigureRequest, resp *datasource.ConfigureResponse) {
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

type groupDataSourceData struct {
	ApiKey types.String `tfsdk:"api_key"`
	Id     types.String `tfsdk:"id"`
}

//type groupDataSource struct {
//	provider provider
//}

func (d *groupDataSource) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
	var state GroupList

	res, statusCode, err := d.p.dog.GetGroups(nil)
	if (statusCode < 200 || statusCode > 299) && statusCode != 404 {
		resp.Diagnostics.AddError("Client Unsuccesful", fmt.Sprintf("Status Code: %d", statusCode))
		return
	}
	if err != nil {
		resp.Diagnostics.AddError("Client Error", fmt.Sprintf("Unable to read groups, got error: %s", err))
		return
	}

	if resp.Diagnostics.HasError() {
		return
	}

	// Set state
	for _, api_group := range res {
		group := ApiToGroup(api_group)
		state = append(state, group)
	}
	diags := resp.State.Set(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}
}
