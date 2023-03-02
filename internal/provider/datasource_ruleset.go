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
	rulesetDataSource struct {
		p dogProvider
	}

	RulesetList []Ruleset

	Ruleset struct {
		ID      types.String `tfsdk:"id"`
		Name    types.String `tfsdk:"name"`
		Rules   *rulesetResourceRules `tfsdk:"rules"`
		//ProfileId types.String `tfsdk:"profile_id" force:",omitempty"` 
		ProfileId types.String `tfsdk:"profile_id" json:"profile_id,omitempty"` 

	}

	Rules struct {
		Inbound  []*Rule `json:"inbound"`
		Outbound []*Rule `json:"outbound"`
	}

	Rule struct {
		Action       string   `json:"action"`
		Active       bool     `json:"active"`
		Comment      string   `json:"comment"`
		Environments []string `json:"environments"`
		Group        string   `json:"group"`
		GroupType    string   `json:"group_type"`
		Interface    string   `json:"interface"`
		Log          bool     `json:"log"`
		LogPrefix    string   `json:"log_prefix"`
		Order        int      `json:"order"`
		Service      string   `json:"service"`
		States       []string `json:"states"`
		Type         string   `json:"type"`
	}
)

var (
	_ datasource.DataSource = (*rulesetDataSource)(nil)
)

func NewRulesetDataSource() datasource.DataSource {
	return &rulesetDataSource{}
}

func (*rulesetDataSource) Metadata(ctx context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_ruleset"
}

func (*rulesetDataSource) GetSchema(ctx context.Context) (tfsdk.Schema, diag.Diagnostics) {
	return tfsdk.Schema{
		// This description is used by the documentation generator and the language server.
		MarkdownDescription: "Ruleset data source",

		Attributes: map[string]tfsdk.Attribute{
			"api_token": {
				MarkdownDescription: "Ruleset configurable attribute",
				Optional:            true,
				Type:                types.StringType,
			},
			"id": {
				MarkdownDescription: "Ruleset identifier",
				Type:                types.StringType,
				Computed:            true,
			},
		},
	}, nil
}

func (d *rulesetDataSource) Configure(ctx context.Context, req datasource.ConfigureRequest, resp *datasource.ConfigureResponse) {
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

type rulesetDataSourceData struct {
	ApiToken types.String `tfsdk:"api_token"`
	Id       types.String `tfsdk:"id"`
}

//type rulesetDataSource struct {
//	provider provider
//}

func (d *rulesetDataSource) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
	var state RulesetList

	res, statusCode, err := d.p.dog.GetRulesets(nil)
	if (statusCode < 200 || statusCode > 299) && statusCode != 404 {
		resp.Diagnostics.AddError("Client Unsuccesful", fmt.Sprintf("Status Code: %d", statusCode))
	}
	if err != nil {
		resp.Diagnostics.AddError("Client Error", fmt.Sprintf("Unable to read rulesets, got error: %s", err))
	}
	if resp.Diagnostics.HasError() {
		return
	}

	// Set state
	for _, api_ruleset := range res {
		ruleset := ApiToRuleset(ctx, api_ruleset)
		state = append(state, ruleset)
	}
	diags := resp.State.Set(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}
}
