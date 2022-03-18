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

type profileDataSourceType struct{}

func (t profileDataSourceType) GetSchema(ctx context.Context) (tfsdk.Schema, diag.Diagnostics) {
	tflog.Debug(ctx, "GetSchema 1\n")
	return tfsdk.Schema{
		// This description is used by the documentation generator and the language server.
		MarkdownDescription: "Profile data source",

		Attributes: map[string]tfsdk.Attribute{
			"api_key": {
				MarkdownDescription: "Profile configurable attribute",
				Optional:            true,
				Type:                types.StringType,
			},
			"id": {
				MarkdownDescription: "Profile identifier",
				Type:                types.StringType,
				Computed:            true,
			},
		},
	}, nil
}

func (t profileDataSourceType) NewDataSource(ctx context.Context, in tfsdk.Provider) (tfsdk.DataSource, diag.Diagnostics) {
	provider, diags := convertProviderType(in)

	return profileDataSource{
		provider: provider,
	}, diags
}

type profileDataSourceData struct {
	ApiKey types.String `tfsdk:"api_key"`
	Id     types.String `tfsdk:"id"`
}

type profileDataSource struct {
	provider provider
}

func (d profileDataSource) Read(ctx context.Context, req tfsdk.ReadDataSourceRequest, resp *tfsdk.ReadDataSourceResponse) {
	tflog.Debug(ctx, "Read 1\n")
	var data profileDataSourceData

	var resourceState struct {
		Profiles ProfileList `tfsdk:"profiles"`
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
	profiles, statusCode, err := d.provider.client.GetProfiles(nil)
	for _, profile := range profiles {
		var newInboundRules []Rule
		for _, inbound_rule := range profile.Rules.Inbound {
			rule := Rule{
				Action:       types.String{Value: inbound_rule.Action},
				Active:       types.Bool{Value: inbound_rule.Active},
				Comment:      types.String{Value: inbound_rule.Comment},
				Environments: inbound_rule.Environments,
				Group:        types.String{Value: inbound_rule.Group},
				GroupType:    types.String{Value: inbound_rule.GroupType},
				Interface:    types.String{Value: inbound_rule.Interface},
				Log:          types.Bool{Value: inbound_rule.Log},
				LogPrefix:    types.String{Value: inbound_rule.LogPrefix},
				Order:        types.Int64{Value: int64(inbound_rule.Order)},
				Service:      types.String{Value: inbound_rule.Service},
				States:       inbound_rule.States,
				Type:         types.String{Value: inbound_rule.Type},
			}
			newInboundRules = append(newInboundRules, rule)
		}
		var newOutboundRules []Rule
		for _, outbound_rule := range profile.Rules.Inbound {
			rule := Rule{
				Action:       types.String{Value: outbound_rule.Action},
				Active:       types.Bool{Value: outbound_rule.Active},
				Comment:      types.String{Value: outbound_rule.Comment},
				Environments: outbound_rule.Environments,
				Group:        types.String{Value: outbound_rule.Group},
				GroupType:    types.String{Value: outbound_rule.GroupType},
				Interface:    types.String{Value: outbound_rule.Interface},
				Log:          types.Bool{Value: outbound_rule.Log},
				LogPrefix:    types.String{Value: outbound_rule.LogPrefix},
				Order:        types.Int64{Value: int64(outbound_rule.Order)},
				Service:      types.String{Value: outbound_rule.Service},
				States:       outbound_rule.States,
				Type:         types.String{Value: outbound_rule.Type},
			}
			newOutboundRules = append(newOutboundRules, rule)
		}
		h := Profile{
			//Created:     types.Int64{Value: int64(profile.Created)},
			Description: types.String{Value: profile.Description},
			ID:          types.String{Value: profile.ID},
			Name:        types.String{Value: profile.Name},
			Rules: Rules{
				Inbound:  newInboundRules,
				Outbound: newOutboundRules,
			},
			Version: types.String{Value: profile.Version},
		}
		resourceState.Profiles = append(resourceState.Profiles, h)
	}
	if statusCode < 200 && statusCode > 299 {
		resp.Diagnostics.AddError("Client Unsuccesful", fmt.Sprintf("Status Code: %d", statusCode))
		return
	}
	if err != nil {
		resp.Diagnostics.AddError("Client Error", fmt.Sprintf("Unable to read profile, got error: %s", err))
		return
	}

	// For the purposes of this profile code, hardcoding a response value to
	// save into the Terraform state.
	data.Id = types.String{Value: "Profile.ID"}
	diags = resp.State.Set(ctx, &resourceState)
	resp.Diagnostics.Append(diags...)
}
