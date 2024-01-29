package dog

import (
	"context"
	"fmt"
	"log"

	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-log/tflog"
	api "github.com/relaypro-open/dog_api_golang/api"
	"golang.org/x/exp/slices"
)

type (
	factResource struct {
		p dogProvider
	}
)

var (
	_ resource.Resource                = (*factResource)(nil)
	_ resource.ResourceWithImportState = (*factResource)(nil)
)

func NewFactResource() resource.Resource {
	return &factResource{}
}

func (*factResource) Metadata(ctx context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_fact"
}

func (*factResource) Schema(ctx context.Context, req resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = schema.Schema{
		// This description is used by the documentation generator and the language server.

		Attributes: map[string]schema.Attribute{
			// This description is used by the documentation generator and the language server.
			//"groups": schema.MapAttribute{
			//ElementType: types.StringType,
			//"groups": schema.ListNestedAttribute{
			//NestedObject: schema.NestedAttributeObject{
			//Attributes: map[string]schema.Attribute{
			"groups": schema.MapNestedAttribute{
				NestedObject: schema.NestedAttributeObject{
					Attributes: map[string]schema.Attribute{
						"vars": schema.StringAttribute{
							MarkdownDescription: "json string of vars",
							Required:            true,
						},
						"hosts": schema.MapAttribute{
							Required:            true,
							ElementType:         types.MapType{ElemType: types.StringType},
						},
						"children": schema.ListAttribute{
							Required:            true,
							ElementType:         types.StringType,
						},
					},
				},
				Required:            true,
			},
			"name": schema.StringAttribute{
				Required:            true,
			},
			"id": schema.StringAttribute{
				Optional:            true,
				Computed: true,
			},
		},
	}
}

func (r *factResource) Configure(ctx context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {
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

func (*factResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	resource.ImportStatePassthroughID(ctx, path.Root("id"), req, resp)
}

type factResourceData struct {
	ID     types.String               `tfsdk:"id"`
	Groups map[string]*FactGroup `tfsdk:"groups"`
	Name   string                     `tfsdk:"name"`
}

func FactToCreateRequest(plan factResourceData) api.FactCreateRequest {
	newGroups := map[string]*api.FactGroup{}
	for name, group := range plan.Groups {
		g := &api.FactGroup{
			Vars:     group.Vars.ValueString(),
			Hosts:    group.Hosts,
			Children: group.Children,
		}
		newGroups[name] = g
	}
	newFact := api.FactCreateRequest{
		Groups: newGroups,
		Name:   plan.Name,
	}
	return newFact
}

func FactToApiFact(plan Fact) api.Fact {
	newGroups := map[string]*api.FactGroup{}
	for name, group := range plan.Groups {
		g := &api.FactGroup{
			Vars:     group.Vars.ValueString(),
			Hosts:    group.Hosts,
			Children: group.Children,
		}
		newGroups[name] = g
	}
	newFact := api.Fact{
		Groups: newGroups,
		Name:   plan.Name.ValueString(),
	}
	return newFact
}

func FactToUpdateRequest(plan factResourceData) api.FactUpdateRequest {
	newGroups := map[string]*api.FactGroup{}
	for name, group := range plan.Groups {
		g := &api.FactGroup{
			Vars:     group.Vars.ValueString(),
			Hosts:    group.Hosts,
			Children: group.Children,
		}
		newGroups[name] = g
	}
	newFact := api.FactUpdateRequest{
		Groups: newGroups,
		Name:   plan.Name,
	}
	return newFact
}

func ApiToFact(fact api.Fact) Fact {
	newGroups := map[string]*FactGroup{}
	for name, group := range fact.Groups {
		g := &FactGroup{
			Vars:     types.StringValue(group.Vars),
			Hosts:    group.Hosts,
			Children: group.Children,
		}
		newGroups[name] = g
	}
	h := Fact{
		ID:     types.StringValue(fact.ID),
		Groups: newGroups,
		Name:   types.StringValue(fact.Name),
	}

	return h
}

func (r *factResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	var state Fact

	var plan Fact
	diags := req.Plan.Get(ctx, &plan)
	resp.Diagnostics.Append(diags...)
	//	resp.Diagnostics.AddError("Client Error", fmt.Sprintf("client: %+v\n", r.provider.client))
	if resp.Diagnostics.HasError() {
		return
	}

	newFact := FactToApiFact(plan)
	log.Printf(fmt.Sprintf("r.p.dog: %+v\n", r.p.dog))
	fact, statusCode, err := r.p.dog.CreateFactEncode(newFact, nil)
	log.Printf(fmt.Sprintf("fact: %+v\n", fact))
	tflog.Trace(ctx, fmt.Sprintf("fact: %+v\n", fact))
	if err != nil {
		resp.Diagnostics.AddError("Client Error", fmt.Sprintf("Unable to create fact, got error: %s", err))
	}
	if statusCode != 201 {
		resp.Diagnostics.AddError("Client Unsuccesful", fmt.Sprintf("Status Code: %d", statusCode))
	}
	if resp.Diagnostics.HasError() {
		return
	}
	state = ApiToFact(fact)

	plan.ID = state.ID

	// write logs using the tflog package
	// see https://pkg.go.dev/github.com/hashicorp/terraform-plugin-log/tflog
	// for more information
	tflog.Trace(ctx, "created a resource")

	diags = resp.State.Set(ctx, &state)
	resp.Diagnostics.Append(diags...)
}

func (r *factResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	var state Fact

	diags := req.State.Get(ctx, &state)
	resp.Diagnostics.Append(diags...)

	if resp.Diagnostics.HasError() {
		return
	}

	factID := state.ID.ValueString()

	log.Printf(fmt.Sprintf("r.p: %+v\n", r.p))
	log.Printf(fmt.Sprintf("r.p.dog: %+v\n", r.p.dog))
	fact, statusCode, err := r.p.dog.GetFactEncode(factID, nil)
	if statusCode != 200 {
		resp.Diagnostics.AddError("Client Unsuccesful", fmt.Sprintf("Status Code: %d", statusCode))
	}
	if err != nil {
		resp.Diagnostics.AddError("Client Error", fmt.Sprintf("Unable to read fact, got error: %s", err))
	}
	if resp.Diagnostics.HasError() {
		return
	}
	state = ApiToFact(fact)
	diags = resp.State.Set(ctx, &state)
	resp.Diagnostics.Append(diags...)
}

func (r *factResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	var state Fact

	diags := req.State.Get(ctx, &state)
	resp.Diagnostics.Append(diags...)

	if resp.Diagnostics.HasError() {
		return
	}

	factID := state.ID.ValueString()

	var plan Fact
	diags = req.Plan.Get(ctx, &plan)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	newFact := FactToApiFact(plan)
	fact, statusCode, err := r.p.dog.UpdateFactEncode(factID, newFact, nil)
	log.Printf(fmt.Sprintf("fact: %+v\n", fact))
	tflog.Trace(ctx, fmt.Sprintf("fact: %+v\n", fact))
	state = ApiToFact(fact)
	if err != nil {
		resp.Diagnostics.AddError("Client Error", fmt.Sprintf("Unable to create fact, got error: %s", err))
	}
	ok := []int{303, 200, 201}
	if slices.Contains(ok, statusCode) != true {
		resp.Diagnostics.AddError("Client Unsuccesful", fmt.Sprintf("Status Code: %d", statusCode))
	}
	if resp.Diagnostics.HasError() {
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

func (r *factResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	var state Fact

	diags := req.State.Get(ctx, &state)
	resp.Diagnostics.Append(diags...)

	if resp.Diagnostics.HasError() {
		return
	}

	factID := state.ID.ValueString()
	fact, statusCode, err := r.p.dog.DeleteFact(factID, nil)
	if statusCode != 204 {
		resp.Diagnostics.AddError("Client Unsuccesful", fmt.Sprintf("Status Code: %d", statusCode))
	}
	if err != nil {
		resp.Diagnostics.AddError("Client Error", fmt.Sprintf("Unable to read fact, got error: %s", err))
	}
	if resp.Diagnostics.HasError() {
		return
	}
	tflog.Trace(ctx, fmt.Sprintf("fact deleted: %+v\n", fact))

	diags = resp.State.Set(ctx, &state)
	resp.Diagnostics.Append(diags...)
	resp.State.RemoveResource(ctx)
}
