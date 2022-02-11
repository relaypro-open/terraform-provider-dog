package provider

import (
	"context"
	"fmt"

	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/tfsdk"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-go/tftypes"
	"github.com/hashicorp/terraform-plugin-log/tflog"
)

type hostResourceType struct{}

func (t hostResourceType) GetSchema(ctx context.Context) (tfsdk.Schema, diag.Diagnostics) {
	return tfsdk.Schema{
		Attributes: map[string]tfsdk.Attribute{
			// This description is used by the documentation generator and the language server.
			"active": {
				MarkdownDescription: "Host active state",
				Optional:            true,
				Type:                types.StringType,
			},
			"environment": {
				MarkdownDescription: "Host environment",
				Optional:            false,
				Type:                types.StringType,
			},
			"group": {
				MarkdownDescription: "Host group",
				Optional:            false,
				Type:                types.StringType,
			},
			"hostkey": {
				MarkdownDescription: "Host key",
				Optional:            false,
				Type:                types.StringType,
			},
			"location": {
				MarkdownDescription: "Host location",
				Optional:            false,
				Type:                types.StringType,
			},
			"name": {
				MarkdownDescription: "Host name",
				Optional:            false,
				Type:                types.StringType,
			},
			"id": {
				Computed:            true,
				MarkdownDescription: "Host identifier",
				PlanModifiers: tfsdk.AttributePlanModifiers{
					tfsdk.UseStateForUnknown(),
				},
				Type: types.StringType,
			},
		},
	}, nil
}

func (t hostResourceType) NewResource(ctx context.Context, in tfsdk.Provider) (tfsdk.Resource, diag.Diagnostics) {
	provider, diags := convertProviderType(in)

	return hostResource{
		provider: provider,
	}, diags
}

type hostResourceData struct {
	Active      *string      `tfsdk:"active"`
	Environment string       `tfsdk:"environment"`
	Group       string       `tfsdk:"group"`
	ID          types.String `tfsdk:"id"`
	HostKey     string       `tfsdk:"hostkey"`
	Location    string       `tfsdk:"location"`
	Name        string       `tfsdk:"name"`
}

type hostResource struct {
	provider provider
}

func (r hostResource) Create(ctx context.Context, req tfsdk.CreateResourceRequest, resp *tfsdk.CreateResourceResponse) {
	var data hostResourceData

	diags := req.Config.Get(ctx, &data)
	resp.Diagnostics.Append(diags...)

	if resp.Diagnostics.HasError() {
		return
	}

	// If applicable, this is a great opportunity to initialize any necessary
	// provider client data and make a call using it.
	// host, err := d.provider.client.CreateHost(...)
	// if err != nil {
	//     resp.Diagnostics.AddError("Client Error", fmt.Sprintf("Unable to create host, got error: %s", err))
	//     return
	// }
	hostID := data.ID.Value
	host, statusCode, err := r.provider.client.GetHost(hostID, nil)
	if err != nil {
		resp.Diagnostics.AddError("Client Error", fmt.Sprintf("Unable to create host, got error: %s", err))
		return
	}
	if statusCode < 200 && statusCode > 299 {
		resp.Diagnostics.AddError("Client Unsuccesful", fmt.Sprintf("Status Code: %d", statusCode))
		return
	}

	// For the purposes of this host code, hardcoding a response value to
	// save into the Terraform state.
	data.ID = types.String{Value: host.ID}

	// write logs using the tflog package
	// see https://pkg.go.dev/github.com/hashicorp/terraform-plugin-log/tflog
	// for more information
	tflog.Trace(ctx, "created a resource")

	diags = resp.State.Set(ctx, &data)
	resp.Diagnostics.Append(diags...)
}

func (r hostResource) Read(ctx context.Context, req tfsdk.ReadResourceRequest, resp *tfsdk.ReadResourceResponse) {
	var state Host

	diags := req.State.Get(ctx, &state)
	resp.Diagnostics.Append(diags...)

	if resp.Diagnostics.HasError() {
		return
	}

	// If applicable, this is a great opportunity to initialize any necessary
	// provider client data and make a call using it.
	// host, err := d.provider.client.ReadHost(...)
	// if err != nil {
	//     resp.Diagnostics.AddError("Client Error", fmt.Sprintf("Unable to read host, got error: %s", err))
	//     return
	// }
	hostID := state.ID.Value

	host, statusCode, err := r.provider.client.GetHost(hostID, nil)
	if statusCode < 200 && statusCode > 299 {
		resp.Diagnostics.AddError("Client Unsuccesful", fmt.Sprintf("Status Code: %d", statusCode))
		return
	}
	if err != nil {
		resp.Diagnostics.AddError("Client Error", fmt.Sprintf("Unable to read host, got error: %s", err))
		return
	}
	h := Host{
		Active:      types.String{Value: host.Active},
		Environment: types.String{Value: host.Environment},
		Group:       types.String{Value: host.Group},
		ID:          types.String{Value: host.ID},
		HostKey:     types.String{Value: host.HostKey},
		Location:    types.String{Value: host.Location},
		Name:        types.String{Value: host.Name},
	}

	state = h
	diags = resp.State.Set(ctx, &state)
	resp.Diagnostics.Append(diags...)
}

func (r hostResource) Update(ctx context.Context, req tfsdk.UpdateResourceRequest, resp *tfsdk.UpdateResourceResponse) {
	var data hostResourceData

	diags := req.Plan.Get(ctx, &data)
	resp.Diagnostics.Append(diags...)

	if resp.Diagnostics.HasError() {
		return
	}

	// If applicable, this is a great opportunity to initialize any necessary
	// provider client data and make a call using it.
	// host, err := d.provider.client.UpdateHost(...)
	// if err != nil {
	//     resp.Diagnostics.AddError("Client Error", fmt.Sprintf("Unable to update host, got error: %s", err))
	//     return
	// }

	diags = resp.State.Set(ctx, &data)
	resp.Diagnostics.Append(diags...)
}

func (r hostResource) Delete(ctx context.Context, req tfsdk.DeleteResourceRequest, resp *tfsdk.DeleteResourceResponse) {
	var data hostResourceData

	diags := req.State.Get(ctx, &data)
	resp.Diagnostics.Append(diags...)

	if resp.Diagnostics.HasError() {
		return
	}

	// If applicable, this is a great opportunity to initialize any necessary
	// provider client data and make a call using it.
	// host, err := d.provider.client.DeleteHost(...)
	// if err != nil {
	//     resp.Diagnostics.AddError("Client Error", fmt.Sprintf("Unable to delete host, got error: %s", err))
	//     return
	// }

	resp.State.RemoveResource(ctx)
}

func (r hostResource) ImportState(ctx context.Context, req tfsdk.ImportResourceStateRequest, resp *tfsdk.ImportResourceStateResponse) {
	tfsdk.ResourceImportStatePassthroughID(ctx, tftypes.NewAttributePath().WithAttributeName("id"), req, resp)
}
