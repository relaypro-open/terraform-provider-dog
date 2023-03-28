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
	groupDataSource struct {
		p dogProvider
	}

	GroupList []Group

	Group struct {
		Description         types.String           `tfsdk:"description"`
		ID                  types.String           `tfsdk:"id"`
		Name                types.String           `tfsdk:"name"`
		ProfileId           types.String           `tfsdk:"profile_id"`
		ProfileName         types.String           `tfsdk:"profile_name"`
		ProfileVersion      types.String           `tfsdk:"profile_version"`
		Ec2SecurityGroupIds []*Ec2SecurityGroupIds `tfsdk:"ec2_security_group_ids"`
		Vars                map[string]string      `tfsdk:"vars"`
	}

	Ec2SecurityGroupIds struct {
		Region types.String `tfsdk:"region"`
		SgId   types.String `tfsdk:"sgid"`
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
			// This description is used by the documentation generator and the language server.
			"description": {
				MarkdownDescription: "group description",
				Optional:            true,
				Type:                types.StringType,
			},
			"name": {
				MarkdownDescription: "group name",
				Required:            true,
				Type:                types.StringType,
			},
			"profile_id": {
				MarkdownDescription: "group profile id",
				Required:            true,
				Type:                types.StringType,
			},
			"profile_name": {
				MarkdownDescription: "group profile name",
				Optional:            true,
				Type:                types.StringType,
			},
			"profile_version": {
				MarkdownDescription: "group profile version",
				Optional:            true,
				Type:                types.StringType,
			},
			"ec2_security_group_ids": {
				MarkdownDescription: "List of EC2 Security Groups to control",
				Optional:            true,
				Attributes: tfsdk.ListNestedAttributes(map[string]tfsdk.Attribute{
					"region": {
						MarkdownDescription: "EC2 Region",
						Required:            true,
						Type:                types.StringType,
					},
					"sgid": {
						MarkdownDescription: "EC2 Security Group ID",
						Required:            true,
						Type:                types.StringType,
					},
				}),
			},
			"vars": {
				MarkdownDescription: "Arbitrary collection of variables used for inventory",
				Type:                types.MapType{ElemType: types.StringType},
				Optional:            true,
			},
			"id": {
				MarkdownDescription: "group identifier",
				Type: types.StringType,
				Required: true,
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
	ApiToken types.String `tfsdk:"api_token"`
	Id       types.String `tfsdk:"id"`
}

//type groupDataSource struct {
//	provider provider
//}

func (d *groupDataSource) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
	var state GroupList

	res, statusCode, err := d.p.dog.GetGroups(nil)
	if (statusCode < 200 || statusCode > 299) && statusCode != 404 {
		resp.Diagnostics.AddError("Client Unsuccesful", fmt.Sprintf("Status Code: %d", statusCode))
	}
	if err != nil {
		resp.Diagnostics.AddError("Client Error", fmt.Sprintf("Unable to read groups, got error: %s", err))
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
