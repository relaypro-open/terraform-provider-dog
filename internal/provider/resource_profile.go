package dog

import (
	"context"
	"fmt"
	"log"

	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/tfsdk"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-log/tflog"
	api "github.com/relaypro-open/dog_api_golang/api"
	"github.com/hashicorp/terraform-plugin-framework/path"
)


type (
	profileResource struct {
		p dogProvider
	}
)

var (
	_ resource.Resource                = (*profileResource)(nil)
	_ resource.ResourceWithImportState = (*profileResource)(nil)
)

func NewProfileResource() resource.Resource {
	return &profileResource{}
}

func (*profileResource) Metadata(ctx context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_profile"
}


func (*profileResource) GetSchema(ctx context.Context) (tfsdk.Schema, diag.Diagnostics) {
	return tfsdk.Schema{
		Attributes: map[string]tfsdk.Attribute{
			"name": {
				MarkdownDescription: "Profile name",
				Required:            true,
				Type:                types.StringType,
			},
			"rules": {
				MarkdownDescription: "Profile rules",
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
			"version": {
				MarkdownDescription: "Profile version",
				Required:            true,
				Type:                types.StringType,
			},
			"id": {
				Computed:            true,
				MarkdownDescription: "Profile identifier",
				PlanModifiers: tfsdk.AttributePlanModifiers{
					resource.UseStateForUnknown(),
				},
				Type: types.StringType,
			},
		},
	}, nil
}

func (r *profileResource) Configure(ctx context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {
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

func (*profileResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	resource.ImportStatePassthroughID(ctx, path.Root("id"), req, resp)
}


type profileResourceData struct {
	ID      types.String          `tfsdk:"id"`
	Name    string                `tfsdk:"name"`
	Rules   *profileResourceRules `tfsdk:"rules"`
	Version string                `tfsdk:"version"`
}

type profileResourceRules struct {
	Inbound  []*profileReseourceRule `tfsdk:"inbound"`
	Outbound []*profileReseourceRule `tfsdk:"outbound"`
}

type profileReseourceRule struct {
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

func ProfileToCreateRequest(plan profileResourceData) api.ProfileCreateRequest {
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

	newProfile := api.ProfileCreateRequest{
		Name: plan.Name,
		Rules: &api.Rules{
			Inbound:  inboundRules,
			Outbound: outboundRules,
		},
		Version: plan.Version,
	}
	return newProfile
}

func ProfileToUpdateRequest(plan profileResourceData) api.ProfileUpdateRequest {
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

	newProfile := api.ProfileUpdateRequest{
		Name: plan.Name,
		Rules: &api.Rules{
			Inbound:  inboundRules,
			Outbound: outboundRules,
		},
		Version: plan.Version,
	}
	return newProfile
}

func ApiToProfile(profile api.Profile) Profile {
	newInboundRules := []*Rule{}
	for _, inbound_rule := range profile.Rules.Inbound {
		rule := &Rule{
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
	newOutboundRules := []*Rule{}
	for _, outbound_rule := range profile.Rules.Outbound {
		rule := &Rule{
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
		ID:   types.String{Value: profile.ID},
		Name: types.String{Value: profile.Name},
		Rules: &Rules{
			Inbound:  newInboundRules,
			Outbound: newOutboundRules,
		},
		Version: types.String{Value: profile.Version},
	}
	return h
}

func (r *profileResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	var state Profile

	var plan profileResourceData
	diags := req.Plan.Get(ctx, &plan)
	resp.Diagnostics.Append(diags...)
	//	resp.Diagnostics.AddError("Client Error", fmt.Sprintf("client: %+v\n", r.provider.client))
	if resp.Diagnostics.HasError() {
		return
	}

	newProfile := ProfileToCreateRequest(plan)
	log.Printf(fmt.Sprintf("r.p.dog: %+v\n", r.p.dog))
	profile, statusCode, err := r.p.dog.CreateProfile(newProfile, nil)
	log.Printf(fmt.Sprintf("profile: %+v\n", profile))
	tflog.Trace(ctx, fmt.Sprintf("profile: %+v\n", profile))
	if statusCode < 200 || statusCode > 299 {
		resp.Diagnostics.AddError("Client Unsuccesful", fmt.Sprintf("Status Code: %d", statusCode))
		return
	}
	if err != nil {
		resp.Diagnostics.AddError("Client Error", fmt.Sprintf("Unable to create profile, got error: %s", err))
		return
	}
	state = ApiToProfile(profile)

	plan.ID = state.ID

	// write logs using the tflog package
	// see https://pkg.go.dev/github.com/hashicorp/terraform-plugin-log/tflog
	// for more information
	tflog.Trace(ctx, "created a resource")

	diags = resp.State.Set(ctx, &state)
	resp.Diagnostics.Append(diags...)
}

func (r *profileResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	var state Profile

	diags := req.State.Get(ctx, &state)
	resp.Diagnostics.Append(diags...)

	if resp.Diagnostics.HasError() {
		return
	}

	profileID := state.ID.Value

	log.Printf(fmt.Sprintf("r.p: %+v\n", r.p))
	log.Printf(fmt.Sprintf("r.p.dog: %+v\n", r.p.dog))
	profile, statusCode, err := r.p.dog.GetProfile(profileID, nil)
	if (statusCode < 200 || statusCode > 299) && statusCode != 404 {
		resp.Diagnostics.AddError("Client Unsuccesful", fmt.Sprintf("Status Code: %d", statusCode))
		return
	}
	if err != nil {
		resp.Diagnostics.AddError("Client Error", fmt.Sprintf("Unable to read profile, got error: %s", err))
		return
	}
	state = ApiToProfile(profile)
	diags = resp.State.Set(ctx, &state)
	resp.Diagnostics.Append(diags...)
}


func (r *profileResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	//var data profileResourceData

	//diags := req.Plan.Get(ctx, &data)
	//resp.Diagnostics.Append(diags...)

	//if resp.Diagnostics.HasError() {
	//	return
	//}
	var state Profile

	diags := req.State.Get(ctx, &state)
	resp.Diagnostics.Append(diags...)

	if resp.Diagnostics.HasError() {
		return
	}

	profileID := state.ID.Value

	var plan profileResourceData
	diags = req.Plan.Get(ctx, &plan)
	resp.Diagnostics.Append(diags...)
	//	resp.Diagnostics.AddError("Client Error", fmt.Sprintf("client: %+v\n", r.provider.client))
	if resp.Diagnostics.HasError() {
		return
	}

	newProfile := ProfileToUpdateRequest(plan)
	profile, statusCode, err := r.p.dog.UpdateProfile(profileID, newProfile, nil)
	log.Printf(fmt.Sprintf("profile: %+v\n", profile))
	tflog.Trace(ctx, fmt.Sprintf("profile: %+v\n", profile))
	//resp.Diagnostics.AddError("profile", fmt.Sprintf("profile: %+v\n", profile))
	state = ApiToProfile(profile)
	if err != nil {
		resp.Diagnostics.AddError("Client Error", fmt.Sprintf("Unable to create profile, got error: %s", err))
		return
	}
	if statusCode < 200 || statusCode > 299 {
		resp.Diagnostics.AddError("Client Unsuccesful", fmt.Sprintf("Status Code: %d", statusCode))
		return
	}

	plan.ID = state.ID

	// write logs using the tflog package
	// see https://pkg.go.dev/github.com/hashicorp/terraform-plugin-log/tflog
	// for more information
	tflog.Trace(ctx, "created a resource")

	diags = resp.State.Set(ctx, &state)
	resp.Diagnostics.Append(diags...)

}

func (r *profileResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	var state Profile

	diags := req.State.Get(ctx, &state)
	resp.Diagnostics.Append(diags...)

	if resp.Diagnostics.HasError() {
		return
	}

	profileID := state.ID.Value
	profile, statusCode, err := r.p.dog.DeleteProfile(profileID, nil)
	if statusCode < 200 || statusCode > 299 {
		resp.Diagnostics.AddError("Client Unsuccesful", fmt.Sprintf("Status Code: %d", statusCode))
		return
	}
	if err != nil {
		resp.Diagnostics.AddError("Client Error", fmt.Sprintf("Unable to read profile, got error: %s", err))
		return
	}
	tflog.Trace(ctx, fmt.Sprintf("profile deleted: %+v\n", profile))

	diags = resp.State.Set(ctx, &state)
	resp.Diagnostics.Append(diags...)
	resp.State.RemoveResource(ctx)
}
