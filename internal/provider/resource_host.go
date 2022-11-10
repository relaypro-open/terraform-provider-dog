package dog

import (
	"context"
	"fmt"
	"log"

	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/tfsdk"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-log/tflog"
	api "github.com/relaypro-open/dog_api_golang/api"
	"github.com/hashicorp/terraform-plugin-framework/path"
)


type (
	hostResource struct {
		p dogProvider
	}
)

var (
	_ resource.Resource                = (*hostResource)(nil)
	_ resource.ResourceWithImportState = (*hostResource)(nil)
)

func NewHostResource() resource.Resource {
	return &hostResource{}
}

func (*hostResource) Metadata(ctx context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_host"
}


func (*hostResource) GetSchema(ctx context.Context) (tfsdk.Schema, diag.Diagnostics) {
	return tfsdk.Schema{
		Attributes: map[string]tfsdk.Attribute{
			// This description is used by the documentation generator and the language server.
			"environment": {
				MarkdownDescription: "Host environment",
				Required:            true,
				Type:                types.StringType,
			},
			"group": {
				MarkdownDescription: "Host group",
				Required:            true,
				Type:                types.StringType,
			},
			"hostkey": {
				MarkdownDescription: "Host key",
				Required:            true,
				Type:                types.StringType,
			},
			"location": {
				MarkdownDescription: "Host location",
				Required:            true,
				Type:                types.StringType,
			},
			"name": {
				MarkdownDescription: "Host name",
				Required:            true,
				Type:                types.StringType,
			},
			"id": {
				Computed:            true,
				MarkdownDescription: "Host identifier",
				PlanModifiers: tfsdk.AttributePlanModifiers{
					resource.UseStateForUnknown(),
				},
				Type: types.StringType,
			},
		},
	}, nil
}

func (r *hostResource) Configure(ctx context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {
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

func (*hostResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	resource.ImportStatePassthroughID(ctx, path.Root("id"), req, resp)
}


type hostResourceData struct {
	Environment string       `tfsdk:"environment"`
	Group       string       `tfsdk:"group"`
	ID          types.String `tfsdk:"id"`
	HostKey     string       `tfsdk:"hostkey"`
	Location    string       `tfsdk:"location"`
	Name        string       `tfsdk:"name"`
}

func HostToCreateRequest(plan hostResourceData) api.HostCreateRequest {
	newHost := api.HostCreateRequest{
		Environment: plan.Environment,
		Group:       plan.Group,
		HostKey:     plan.HostKey,
		Location:    plan.Location,
		Name:        plan.Name,
	}
	return newHost
}

func HostToUpdateRequest(plan hostResourceData) api.HostUpdateRequest {
	newHost := api.HostUpdateRequest{
		Environment: plan.Environment,
		Group:       plan.Group,
		HostKey:     plan.HostKey,
		Location:    plan.Location,
		Name:        plan.Name,
	}
	return newHost
}

func ApiToHost(host api.Host) Host {
	h := Host{
		Environment: types.String{Value: host.Environment},
		Group:       types.String{Value: host.Group},
		ID:          types.String{Value: host.ID},
		HostKey:     types.String{Value: host.HostKey},
		Location:    types.String{Value: host.Location},
		Name:        types.String{Value: host.Name},
	}
	return h
}

func (r *hostResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	var state Host

	var plan hostResourceData
	diags := req.Plan.Get(ctx, &plan)
	resp.Diagnostics.Append(diags...)
	//	resp.Diagnostics.AddError("Client Error", fmt.Sprintf("client: %+v\n", r.provider.client))
	if resp.Diagnostics.HasError() {
		return
	}

	newHost := HostToCreateRequest(plan)
	log.Printf(fmt.Sprintf("r.p.dog: %+v\n", r.p.dog))
	host, statusCode, err := r.p.dog.CreateHost(newHost, nil)
	log.Printf(fmt.Sprintf("host: %+v\n", host))
	tflog.Trace(ctx, fmt.Sprintf("host: %+v\n", host))
	if err != nil {
		resp.Diagnostics.AddError("Client Error", fmt.Sprintf("Unable to create host, got error: %s", err))
		return
	}
	if statusCode < 200 || statusCode > 299 {
		resp.Diagnostics.AddError("Client Unsuccesful", fmt.Sprintf("Status Code: %d", statusCode))
		return
	}
	state = ApiToHost(host)

	plan.ID = state.ID

	// write logs using the tflog package
	// see https://pkg.go.dev/github.com/hashicorp/terraform-plugin-log/tflog
	// for more information
	tflog.Trace(ctx, "created a resource")

	diags = resp.State.Set(ctx, &state)
	resp.Diagnostics.Append(diags...)
}

func (r *hostResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	var state Host

	diags := req.State.Get(ctx, &state)
	resp.Diagnostics.Append(diags...)

	if resp.Diagnostics.HasError() {
		return
	}

	hostID := state.ID.Value

	log.Printf(fmt.Sprintf("r.p: %+v\n", r.p))
	log.Printf(fmt.Sprintf("r.p.dog: %+v\n", r.p.dog))
	host, statusCode, err := r.p.dog.GetHost(hostID, nil)
	if (statusCode < 200 || statusCode > 299) && statusCode != 404 {
		resp.Diagnostics.AddError("Client Unsuccesful", fmt.Sprintf("Status Code: %d", statusCode))
		return
	}
	if err != nil {
		resp.Diagnostics.AddError("Client Error", fmt.Sprintf("Unable to read host, got error: %s", err))
		return
	}
	state = ApiToHost(host)
	diags = resp.State.Set(ctx, &state)
	resp.Diagnostics.Append(diags...)
}


func (r *hostResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	//var data hostResourceData

	//diags := req.Plan.Get(ctx, &data)
	//resp.Diagnostics.Append(diags...)

	//if resp.Diagnostics.HasError() {
	//	return
	//}
	var state Host

	diags := req.State.Get(ctx, &state)
	resp.Diagnostics.Append(diags...)

	if resp.Diagnostics.HasError() {
		return
	}

	hostID := state.ID.Value

	var plan hostResourceData
	diags = req.Plan.Get(ctx, &plan)
	resp.Diagnostics.Append(diags...)
	//	resp.Diagnostics.AddError("Client Error", fmt.Sprintf("client: %+v\n", r.provider.client))
	if resp.Diagnostics.HasError() {
		return
	}

	newHost := HostToUpdateRequest(plan)
	host, statusCode, err := r.p.dog.UpdateHost(hostID, newHost, nil)
	log.Printf(fmt.Sprintf("host: %+v\n", host))
	tflog.Trace(ctx, fmt.Sprintf("host: %+v\n", host))
	//resp.Diagnostics.AddError("host", fmt.Sprintf("host: %+v\n", host))
	state = ApiToHost(host)
	if err != nil {
		resp.Diagnostics.AddError("Client Error", fmt.Sprintf("Unable to create host, got error: %s", err))
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

func (r *hostResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	var state Host

	diags := req.State.Get(ctx, &state)
	resp.Diagnostics.Append(diags...)

	if resp.Diagnostics.HasError() {
		return
	}

	hostID := state.ID.Value
	host, statusCode, err := r.p.dog.DeleteHost(hostID, nil)
	if statusCode < 200 || statusCode > 299 {
		resp.Diagnostics.AddError("Client Unsuccesful", fmt.Sprintf("Status Code: %d", statusCode))
		return
	}
	if err != nil {
		resp.Diagnostics.AddError("Client Error", fmt.Sprintf("Unable to read host, got error: %s", err))
		return
	}
	tflog.Trace(ctx, fmt.Sprintf("host deleted: %+v\n", host))

	diags = resp.State.Set(ctx, &state)
	resp.Diagnostics.Append(diags...)
	resp.State.RemoveResource(ctx)
}
