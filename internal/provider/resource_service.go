package dog

import (
	"context"
	"fmt"
	"log"

	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/tfsdk"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-go/tftypes"
	"github.com/hashicorp/terraform-plugin-log/tflog"
	api "github.com/relaypro-open/dog_api_golang/api"
)

type serviceResourceType struct{}

func (t serviceResourceType) GetSchema(ctx context.Context) (tfsdk.Schema, diag.Diagnostics) {
	return tfsdk.Schema{
		Attributes: map[string]tfsdk.Attribute{
			// This description is used by the documentation generator and the language server.
			//"created": {
			//	MarkdownDescription: "Service created timestamp",
			//	Optional:            true,
			//	Type:                types.Int64Type,
			//},
			"services": {
				MarkdownDescription: "List of Services",
				Required:            true,
				Attributes: tfsdk.ListNestedAttributes(map[string]tfsdk.Attribute{
					"protocol": {
						MarkdownDescription: "Service protocol",
						Required:            true,
						Type:                types.StringType,
					},
					"ports": {
						MarkdownDescription: "Service ports",
						Required:            true,
						Type: types.ListType{
							ElemType: types.StringType,
						},
					},
				}, tfsdk.ListNestedAttributesOptions{}),
			},
			"name": {
				MarkdownDescription: "Service name",
				Required:            true,
				Type:                types.StringType,
			},
			"version": {
				MarkdownDescription: "Service version",
				Optional:            true,
				Type:                types.Int64Type,
			},
			"id": {
				Computed:            true,
				MarkdownDescription: "Service identifier",
				PlanModifiers: tfsdk.AttributePlanModifiers{
					tfsdk.UseStateForUnknown(),
				},
				Type: types.StringType,
			},
		},
	}, nil
}

func (t serviceResourceType) NewResource(ctx context.Context, in tfsdk.Provider) (tfsdk.Resource, diag.Diagnostics) {
	provider, diags := convertProviderType(in)

	return serviceResource{
		provider: provider,
	}, diags
}

type serviceResourceData struct {
	//Created int          `tfsdk:"created"`
	ID       types.String   `tfsdk:"id"`
	Services []PortProtocol `tfsdk:"services"`
	Name     string         `tfsdk:"name"`
	Version  int            `tfsdk:"version"`
}

type serviceResource struct {
	provider provider
}

func ServiceToCreateRequest(plan serviceResourceData) api.ServiceCreateRequest {
	newServices := []api.PortProtocol{}
	for _, port_protocol := range plan.Services {
		var pp api.PortProtocol
		pp = api.PortProtocol{
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
	newServices := []api.PortProtocol{}
	for _, port_protocol := range plan.Services {
		var pp api.PortProtocol
		pp = api.PortProtocol{
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
	s := Services{}
	for _, port_protocol := range service.Services {
		var pp PortProtocol
		pp = PortProtocol{
			Ports:    port_protocol.Ports,
			Protocol: types.String{Value: port_protocol.Protocol},
		}
		s = append(s, pp)
	}
	h := Service{
		ID:       types.String{Value: service.ID},
		Services: s,
		Name:     types.String{Value: service.Name},
		Version:  types.Int64{Value: int64(service.Version)},
	}
	return h
}

func (r serviceResource) Create(ctx context.Context, req tfsdk.CreateResourceRequest, resp *tfsdk.CreateResourceResponse) {
	tflog.Debug(ctx, fmt.Sprintf("ZZZZZZZZZZZZZZZZZZ r: %+v\n", r))
	tflog.Debug(ctx, "Create 1\n")
	var state Service

	var plan serviceResourceData
	diags := req.Plan.Get(ctx, &plan)
	tflog.Debug(ctx, fmt.Sprintf("ZZZZZZZZZZZZZZZZZZ plan: %+v\n", plan))
	resp.Diagnostics.Append(diags...)
	//resp.Diagnostics.AddError("Client Error", fmt.Sprintf("plan: %+v\n", plan))
	if resp.Diagnostics.HasError() {
		return
	}
	tflog.Debug(ctx, "Create 2\n")
	newService := ServiceToCreateRequest(plan)
	tflog.Debug(ctx, fmt.Sprintf("ZZZZZZZZZZZZZZZZZZZZZZZZZ NewService: %+v\n", newService))
	service, statusCode, err := r.provider.client.CreateService(newService, nil)
	log.Printf(fmt.Sprintf("ZZZZZZZZZZZZZZZZZZZZZZZZZZZZZZ service: %+v\n", service))
	tflog.Debug(ctx, fmt.Sprintf("ZZZZZZZZZZZZZZZZZZZZZZZZZ service: %+v\n", service))

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

func (r serviceResource) Read(ctx context.Context, req tfsdk.ReadResourceRequest, resp *tfsdk.ReadResourceResponse) {
	tflog.Debug(ctx, "Read 1\n")
	var state Service

	diags := req.State.Get(ctx, &state)
	resp.Diagnostics.Append(diags...)

	if resp.Diagnostics.HasError() {
		return
	}

	serviceID := state.ID.Value

	service, statusCode, err := r.provider.client.GetService(serviceID, nil)
	if statusCode < 200 || statusCode > 299 {
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

func (r serviceResource) Update(ctx context.Context, req tfsdk.UpdateResourceRequest, resp *tfsdk.UpdateResourceResponse) {
	tflog.Debug(ctx, "Update 1\n")
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
	service, statusCode, err := r.provider.client.UpdateService(serviceID, newService, nil)
	log.Printf(fmt.Sprintf("service: %+v\n", service))
	tflog.Debug(ctx, fmt.Sprintf("service: %+v\n", service))
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

func (r serviceResource) Delete(ctx context.Context, req tfsdk.DeleteResourceRequest, resp *tfsdk.DeleteResourceResponse) {
	tflog.Debug(ctx, "Delete 1\n")
	var state Service

	diags := req.State.Get(ctx, &state)
	resp.Diagnostics.Append(diags...)

	if resp.Diagnostics.HasError() {
		return
	}

	serviceID := state.ID.Value
	service, statusCode, err := r.provider.client.DeleteService(serviceID, nil)
	if statusCode < 200 || statusCode > 299 {
		resp.Diagnostics.AddError("Client Unsuccesful", fmt.Sprintf("Status Code: %d", statusCode))
		return
	}
	if err != nil {
		resp.Diagnostics.AddError("Client Error", fmt.Sprintf("Unable to read service, got error: %s", err))
		return
	}
	tflog.Debug(ctx, fmt.Sprintf("Service deleted: %+v\n", service))

	diags = resp.State.Set(ctx, &state)
	resp.Diagnostics.Append(diags...)
	resp.State.RemoveResource(ctx)
}

func (r serviceResource) ImportState(ctx context.Context, req tfsdk.ImportResourceStateRequest, resp *tfsdk.ImportResourceStateResponse) {
	tfsdk.ResourceImportStatePassthroughID(ctx, tftypes.NewAttributePath().WithAttributeName("id"), req, resp)
}
