
package dog

import (
	"context"
	"fmt"
	"log"

	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/tfsdk"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-log/tflog"
	api "github.com/relaypro-open/dog_api_golang/api"
	"golang.org/x/exp/slices"
)


type (
	rulesetResource struct {
		p dogProvider
	}
)

var (
	_ resource.Resource                = (*rulesetResource)(nil)
	_ resource.ResourceWithImportState = (*rulesetResource)(nil)
)

func NewRulesetResource() resource.Resource {
	return &rulesetResource{}
}

func (*rulesetResource) Metadata(ctx context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_ruleset"
}


func (*rulesetResource) GetSchema(ctx context.Context) (tfsdk.Schema, diag.Diagnostics) {
	return tfsdk.Schema{
		Attributes: map[string]tfsdk.Attribute{
			"rules": {
				MarkdownDescription: "Rule rules",
				Required:            true,
				Type: types.ObjectType{
					AttrTypes: map[string]attr.Type{
						"inbound": types.ListType{
							ElemType: types.ObjectType{
								AttrTypes: map[string]attr.Type{
									"action":  types.StringType,
									"active":  types.BoolType,
									"comment": types.StringType,
									"environments": types.ListType{
										ElemType: types.StringType,
									},
									"group":      types.StringType,
									"group_type": types.StringType,
									"interface":  types.StringType,
									"log":        types.BoolType,
									"log_prefix": types.StringType,
									"order":      types.Int64Type,
									"service":    types.StringType,
									"states": types.ListType{
										ElemType: types.StringType,
									},
									"type": types.StringType,
								},
							},
						},
						"outbound": types.ListType{
							ElemType: types.ObjectType{
								AttrTypes: map[string]attr.Type{
									"action":  types.StringType,
									"active":  types.BoolType,
									"comment": types.StringType,
									"environments": types.ListType{
										ElemType: types.StringType,
									},
									"group":      types.StringType,
									"group_type": types.StringType,
									"interface":  types.StringType,
									"log":        types.BoolType,
									"log_prefix": types.StringType,
									"order":      types.Int64Type,
									"service":    types.StringType,
									"states": types.ListType{
										ElemType: types.StringType,
									},
									"type": types.StringType,
								},
							},
						},
					},
				},
			},
			"id": {
				Computed:            true,
				MarkdownDescription: "Rule identifier",
				PlanModifiers: tfsdk.AttributePlanModifiers{
					resource.UseStateForUnknown(),
				},
				Type: types.StringType,
			},
		},
	}, nil
}

func (r *rulesetResource) Configure(ctx context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {
	// Prevent panic if the provider has not been configured
	if req.ProviderData == nil {
		return
	}

	client, ok := req.ProviderData.(*api.Client)
	if !ok {
		resp.Diagnostics.AddError(
			"Unexpected Resource Configure Type",
			fmt.Sprintf("Expected *api.Client, got: %T. Please report this issue to the provider developers.", req.ProviderData),
		)

		return
	}

	r.p.dog = client
}

func (*rulesetResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	resource.ImportStatePassthroughID(ctx, path.Root("id"), req, resp)
}


type rulesetResourceData struct {
	ID      types.String          `tfsdk:"id"`
	Rules   *rulesetResourceRules `tfsdk:"rules"`
	Name    string                `tfsdk:"name"`
}

type rulesetResourceRules struct {
	Inbound  []*rulesetResourceRule `tfsdk:"inbound"`
	Outbound []*rulesetResourceRule `tfsdk:"outbound"`
}

type rulesetResourceRule struct {
	Action       types.String `tfsdk:"action"`
	Active       types.Bool   `tfsdk:"active"`
	Comment      types.String `tfsdk:"comment"`
	Environments []string     `tfsdk:"environments"`
	Group        types.String `tfsdk:"group"`
	GroupType    types.String `tfsdk:"group_type"`
	Interface    types.String `tfsdk:"interface"`
	Log          types.Bool   `tfsdk:"log"`
	LogPrefix    types.String `tfsdk:"log_prefix"`
	Order        types.Int64  `tfsdk:"order"`
	Service      types.String `tfsdk:"service"`
	States       []string     `tfsdk:"states"`
	Type         types.String `tfsdk:"type"`
}

func RuleToCreateRequest(plan rulesetResourceData) api.RulesetCreateRequest {
	inboundRules := []*api.Rule{}
	for _, inbound_rule := range plan.Rules.Inbound {
		rule := &api.Rule{
			Action:       inbound_rule.Action.Value,
			Active:       inbound_rule.Active.Value,
			Comment:      inbound_rule.Comment.Value,
			Environments: inbound_rule.Environments,
			Group:        inbound_rule.Group.Value,
			GroupType:    inbound_rule.GroupType.Value,
			Interface:    inbound_rule.Interface.Value,
			Log:          inbound_rule.Log.Value,
			LogPrefix:    inbound_rule.LogPrefix.Value,
			Order:        int(inbound_rule.Order.Value),
			Service:      inbound_rule.Service.Value,
			States:       inbound_rule.States,
			Type:         inbound_rule.Type.Value,
		}
		inboundRules = append(inboundRules, rule)
	}
	outboundRules := []*api.Rule{}
	for _, outbound_rule := range plan.Rules.Outbound {
		rule := &api.Rule{
			Action:       outbound_rule.Action.Value,
			Active:       outbound_rule.Active.Value,
			Comment:      outbound_rule.Comment.Value,
			Environments: outbound_rule.Environments,
			Group:        outbound_rule.Group.Value,
			GroupType:    outbound_rule.GroupType.Value,
			Interface:    outbound_rule.Interface.Value,
			Log:          outbound_rule.Log.Value,
			LogPrefix:    outbound_rule.LogPrefix.Value,
			Order:        int(outbound_rule.Order.Value),
			Service:      outbound_rule.Service.Value,
			States:       outbound_rule.States,
			Type:         outbound_rule.Type.Value,
		}
		outboundRules = append(outboundRules, rule)
	}

	newRuleset := api.RulesetCreateRequest{
		Rules: &api.Rules{
			Inbound:  inboundRules,
			Outbound: outboundRules,
		},
	}
	return newRuleset
}

func RuleToUpdateRequest(plan rulesetResourceData) api.RulesetUpdateRequest {
	inboundRules := []*api.Rule{}
	for _, inbound_rule := range plan.Rules.Inbound {
		rule := &api.Rule{
			Action:       inbound_rule.Action.Value,
			Active:       inbound_rule.Active.Value,
			Comment:      inbound_rule.Comment.Value,
			Environments: inbound_rule.Environments,
			Group:        inbound_rule.Group.Value,
			GroupType:    inbound_rule.GroupType.Value,
			Interface:    inbound_rule.Interface.Value,
			Log:          inbound_rule.Log.Value,
			LogPrefix:    inbound_rule.LogPrefix.Value,
			Order:        int(inbound_rule.Order.Value),
			Service:      inbound_rule.Service.Value,
			States:       inbound_rule.States,
			Type:         inbound_rule.Type.Value,
		}
		inboundRules = append(inboundRules, rule)
	}
	outboundRules := []*api.Rule{}
	for _, outbound_rule := range plan.Rules.Outbound {
		rule := &api.Rule{
			Action:       outbound_rule.Action.Value,
			Active:       outbound_rule.Active.Value,
			Comment:      outbound_rule.Comment.Value,
			Environments: outbound_rule.Environments,
			Group:        outbound_rule.Group.Value,
			GroupType:    outbound_rule.GroupType.Value,
			Interface:    outbound_rule.Interface.Value,
			Log:          outbound_rule.Log.Value,
			LogPrefix:    outbound_rule.LogPrefix.Value,
			Order:        int(outbound_rule.Order.Value),
			Service:      outbound_rule.Service.Value,
			States:       outbound_rule.States,
			Type:         outbound_rule.Type.Value,
		}
		outboundRules = append(outboundRules, rule)
	}

	newRuleset := api.RulesetUpdateRequest{
		Rules: &api.Rules{
			Inbound:  inboundRules,
			Outbound: outboundRules,
		},
	}
	return newRuleset
}


func ApiToRuleset(ruleset api.Ruleset) Ruleset {
	newInboundRules := []*rulesetResourceRule{}
	for _, inbound_rule := range ruleset.Rules.Inbound {
		rule := &rulesetResourceRule{
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
	newOutboundRules := []*rulesetResourceRule{}
	for _, outbound_rule := range ruleset.Rules.Outbound {
		rule := &rulesetResourceRule{
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
	h := Ruleset{
		//Created:     types.Int64{Value: int64(rule.Created)},
		ID: types.String{Value: ruleset.ID},
		Name: types.String{Value: ruleset.Name},
		//Rules: &Rules{
		Rules: &rulesetResourceRules{
			Inbound:  newInboundRules,
			Outbound: newOutboundRules,
		},
	}
	return h
}

func (r *rulesetResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	var state Ruleset

	var plan rulesetResourceData
	diags := req.Plan.Get(ctx, &plan)
	resp.Diagnostics.Append(diags...)
	//	resp.Diagnostics.AddError("Client Error", fmt.Sprintf("client: %+v\n", r.provider.client))
	if resp.Diagnostics.HasError() {
		return
	}

	newRule := RuleToCreateRequest(plan)
	log.Printf(fmt.Sprintf("r.p.dog: %+v\n", r.p.dog))
	rule, statusCode, err := r.p.dog.CreateRuleset(newRule, nil)
	log.Printf(fmt.Sprintf("rule: %+v\n", rule))
	tflog.Trace(ctx, fmt.Sprintf("rule: %+v\n", rule))
	if statusCode != 201 {
		resp.Diagnostics.AddError("Client Unsuccesful", fmt.Sprintf("Status Code: %d", statusCode))
	}
	if err != nil {
		resp.Diagnostics.AddError("Client Error", fmt.Sprintf("Unable to create rule, got error: %s", err))
	}
	if resp.Diagnostics.HasError() {
		return
	}
	state = ApiToRuleset(rule)

	plan.ID = state.ID

	// write logs using the tflog package
	// see https://pkg.go.dev/github.com/hashicorp/terraform-plugin-log/tflog
	// for more information
	tflog.Trace(ctx, "created a resource")

	diags = resp.State.Set(ctx, &state)
	resp.Diagnostics.Append(diags...)
}

func (r *rulesetResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	var state Ruleset

	diags := req.State.Get(ctx, &state)
	resp.Diagnostics.Append(diags...)

	if resp.Diagnostics.HasError() {
		return
	}

	rulesetID := state.ID.Value

	log.Printf(fmt.Sprintf("r.p: %+v\n", r.p))
	log.Printf(fmt.Sprintf("r.p.dog: %+v\n", r.p.dog))
	rule, statusCode, err := r.p.dog.GetRuleset(rulesetID, nil)
	if statusCode != 200 {
		resp.Diagnostics.AddError("Client Unsuccesful", fmt.Sprintf("Status Code: %d", statusCode))
	}
	if err != nil {
		resp.Diagnostics.AddError("Client Error", fmt.Sprintf("Unable to read rule, got error: %s", err))
	}
	if resp.Diagnostics.HasError() {
		return
	}
	state = ApiToRuleset(rule)
	diags = resp.State.Set(ctx, &state)
	resp.Diagnostics.Append(diags...)
}

func (r *rulesetResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	var state Ruleset

	diags := req.State.Get(ctx, &state)
	resp.Diagnostics.Append(diags...)

	if resp.Diagnostics.HasError() {
		return
	}

	rulesetID := state.ID.Value

	var plan rulesetResourceData
	diags = req.Plan.Get(ctx, &plan)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	newRule := RuleToUpdateRequest(plan)
	rule, statusCode, err := r.p.dog.UpdateRuleset(rulesetID, newRule, nil)
	log.Printf(fmt.Sprintf("rule: %+v\n", rule))
	tflog.Trace(ctx, fmt.Sprintf("rule: %+v\n", rule))
	state = ApiToRuleset(rule)
	if err != nil {
		resp.Diagnostics.AddError("Client Error", fmt.Sprintf("Unable to create rule, got error: %s", err))
	}
	ok := []int{303, 200, 201}
	if slices.Contains(ok, statusCode) != true {
		resp.Diagnostics.AddError("Client Unsuccesful", fmt.Sprintf("Status Code: %d", statusCode))
	}
	if resp.Diagnostics.HasError() {
		return
	}

	plan.ID = state.ID

	tflog.Trace(ctx, "created a resource")

	diags = resp.State.Set(ctx, &state)
	resp.Diagnostics.Append(diags...)

}

func (r *rulesetResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	var state Ruleset

	diags := req.State.Get(ctx, &state)
	resp.Diagnostics.Append(diags...)

	if resp.Diagnostics.HasError() {
		return
	}

	rulesetID := state.ID.Value
	rule, statusCode, err := r.p.dog.DeleteRuleset(rulesetID, nil)
	if statusCode != 204 {
		resp.Diagnostics.AddError("Client Unsuccesful", fmt.Sprintf("Status Code: %d", statusCode))
	}
	if err != nil {
		resp.Diagnostics.AddError("Client Error", fmt.Sprintf("Unable to read rule, got error: %s", err))
	}
	if resp.Diagnostics.HasError() {
		return
	}
	tflog.Trace(ctx, fmt.Sprintf("rule deleted: %+v\n", rule))

	diags = resp.State.Set(ctx, &state)
	resp.Diagnostics.Append(diags...)
	resp.State.RemoveResource(ctx)
}
