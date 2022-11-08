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
	serviceResource struct {
		p dogProvider
	}
)

var (
	_ resource.Resource                = (*serviceResource)(nil)
	_ resource.ResourceWithImportState = (*serviceResource)(nil)
)

func NewServiceResource() resource.Resource {
	return &serviceResource{}
}

func (*serviceResource) Metadata(ctx context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_service"
}


func (*serviceResource) GetSchema(ctx context.Context) (tfsdk.Schema, diag.Diagnostics) {
	return tfsdk.Schema{
		Attributes: map[string]tfsdk.Attribute{
			// This description is used by the documentation generator and the language server.
			"ipv4_addresses": {
				MarkdownDescription: "List of Ipv4 Addresses",
				Required:            true,
				Type: types.ListType{
					ElemType: types.StringType,
				},
			},
			"ipv6_addresses": {
				MarkdownDescription: "List of Ipv6 Addresses",
				Required:            true,
				Type: types.ListType{
					ElemType: types.StringType,
				},
			},
			"name": {
				MarkdownDescription: "Service name",
				Required:            true,
				Type:                types.StringType,
			},
			"id": {
				Computed:            true,
				MarkdownDescription: "Service identifier",
				PlanModifiers: tfsdk.AttributePlanModifiers{
					resource.UseStateForUnknown(),
				},
				Type: types.StringType,
			},
		},
	}, nil
}

func (r *serviceResource) Configure(ctx context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {
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


func (*serviceResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	resource.ImportStatePassthroughID(ctx, path.Root("id"), req, resp)
}

type serviceResourceData struct {
	ID       types.String                `tfsdk:"id"`
	Services []*portProtocolResourceData `tfsdk:"services"`
	Name     string                      `tfsdk:"name"`
	Version  int                         `tfsdk:"version"`
}

type portProtocolResourceData struct {
	Ports    []string     `tfsdk:"ports"`
	Protocol types.String `tfsdk:"protocol"`
}

func ServiceToCreateRequest(plan serviceResourceData) api.ServiceCreateRequest {
	newServices := []*api.PortProtocol{}
	for _, port_protocol := range plan.Services {
		pp := &api.PortProtocol{
			Ports:    port_protocol.Ports,
			Protocol: port_protocol.Protocol.Value,
		}
		newServices = append(newServices, pp)
	}

	newService := api.ServiceCreateRequest{
		Version:  plan.Version,
		Name:     plan.Name,
		Services: newServices,
	}
	return newService
}

func ServiceToUpdateRequest(plan serviceResourceData) api.ServiceUpdateRequest {
	newServices := []*api.PortProtocol{}
	for _, port_protocol := range plan.Services {
		pp := &api.PortProtocol{
			Ports:    port_protocol.Ports,
			Protocol: port_protocol.Protocol.Value,
		}
		newServices = append(newServices, pp)
	}

	newService := api.ServiceUpdateRequest{
		Version:  plan.Version,
		Name:     plan.Name,
		Services: newServices,
	}
	return newService
}

func ApiToService(service api.Service) Service {
	newServices := []*PortProtocol{}
	for _, port_protocol := range service.Services {
		pp := &PortProtocol{
			Ports:    port_protocol.Ports,
			Protocol: types.String{Value: port_protocol.Protocol},
		}
		newServices = append(newServices, pp)
	}
	h := Service{
		ID:       types.String{Value: service.ID},
		Services: newServices,
		Name:     types.String{Value: service.Name},
		Version:  types.Int64{Value: int64(service.Version)},
	}
	return h
}

func (r *serviceResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	var state Service

	var plan serviceResourceData
	diags := req.Plan.Get(ctx, &plan)
	resp.Diagnostics.Append(diags...)
	//	resp.Diagnostics.AddError("Client Error", fmt.Sprintf("client: %+v\n", r.provider.client))
	if resp.Diagnostics.HasError() {
		return
	}

	newService := ServiceToCreateRequest(plan)
	log.Printf(fmt.Sprintf("r.p.dog: %+v\n", r.p.dog))
	service, statusCode, err := r.p.dog.CreateService(newService, nil)
	log.Printf(fmt.Sprintf("service: %+v\n", service))
	tflog.Trace(ctx, fmt.Sprintf("service: %+v\n", service))
	state = ApiToService(service)
	if err != nil {
		resp.Diagnostics.AddError("Client Error", fmt.Sprintf("Unable to create service, got error: %s", err))
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

func (r *serviceResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	var state Service

	diags := req.State.Get(ctx, &state)
	resp.Diagnostics.Append(diags...)

	if resp.Diagnostics.HasError() {
		return
	}

	serviceID := state.ID.Value

	log.Printf(fmt.Sprintf("r.p: %+v\n", r.p))
	log.Printf(fmt.Sprintf("r.p.dog: %+v\n", r.p.dog))
	service, statusCode, err := r.p.dog.GetService(serviceID, nil)
	if (statusCode < 200 || statusCode > 299) && statusCode != 404 {
		resp.Diagnostics.AddError("Client Unsuccesful", fmt.Sprintf("Status Code: %d", statusCode))
		return
	}
	if err != nil {
		resp.Diagnostics.AddError("Client Error", fmt.Sprintf("Unable to read service, got error: %s", err))
		return
	}
	state = ApiToService(service)
	diags = resp.State.Set(ctx, &state)
	resp.Diagnostics.Append(diags...)
}


func (r *serviceResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	//var data serviceResourceData

	//diags := req.Plan.Get(ctx, &data)
	//resp.Diagnostics.Append(diags...)

	//if resp.Diagnostics.HasError() {
	//	return
	//}
	var state Service

	diags := req.State.Get(ctx, &state)
	resp.Diagnostics.Append(diags...)

	if resp.Diagnostics.HasError() {
		return
	}

	serviceID := state.ID.Value

	var plan serviceResourceData
	diags = req.Plan.Get(ctx, &plan)
	resp.Diagnostics.Append(diags...)
	//	resp.Diagnostics.AddError("Client Error", fmt.Sprintf("client: %+v\n", r.provider.client))
	if resp.Diagnostics.HasError() {
		return
	}

	newService := ServiceToUpdateRequest(plan)
	service, statusCode, err := r.p.dog.UpdateService(serviceID, newService, nil)
	log.Printf(fmt.Sprintf("service: %+v\n", service))
	tflog.Trace(ctx, fmt.Sprintf("service: %+v\n", service))
	//resp.Diagnostics.AddError("service", fmt.Sprintf("service: %+v\n", service))
	state = ApiToService(service)
	if err != nil {
		resp.Diagnostics.AddError("Client Error", fmt.Sprintf("Unable to create service, got error: %s", err))
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

func (r *serviceResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	var state Service

	diags := req.State.Get(ctx, &state)
	resp.Diagnostics.Append(diags...)

	if resp.Diagnostics.HasError() {
		return
	}

	serviceID := state.ID.Value
	service, statusCode, err := r.p.dog.DeleteService(serviceID, nil)
	if statusCode < 200 || statusCode > 299 {
		resp.Diagnostics.AddError("Client Unsuccesful", fmt.Sprintf("Status Code: %d", statusCode))
		return
	}
	if err != nil {
		resp.Diagnostics.AddError("Client Error", fmt.Sprintf("Unable to read service, got error: %s", err))
		return
	}
	tflog.Trace(ctx, fmt.Sprintf("service deleted: %+v\n", service))

	diags = resp.State.Set(ctx, &state)
	resp.Diagnostics.Append(diags...)
	resp.State.RemoveResource(ctx)
}
