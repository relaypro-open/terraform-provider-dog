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

type linkResourceType struct{}

func (t linkResourceType) GetSchema(ctx context.Context) (tfsdk.Schema, diag.Diagnostics) {
	return tfsdk.Schema{
		Attributes: map[string]tfsdk.Attribute{
			// This description is used by the documentation generator and the language server.
			"address_handling": {
				MarkdownDescription: "Type of address handling",
				Required:            true,
				Type:                types.StringType,
			},
			"connection": {
				MarkdownDescription: "Connection specification",
				Required:            true,
				Attributes: tfsdk.SingleNestedAttributes(map[string]tfsdk.Attribute{
					"api_port": {
						Type:     types.Int64Type,
						Required: true,
					},
					"host": {
						Type:     types.StringType,
						Required: true,
					},
					"password": {
						Type:      types.StringType,
						Required:  true,
						Sensitive: true,
					},
					"port": {
						Type:     types.Int64Type,
						Required: true,
					},
					"ssl_options": {
						Required: true,
						Attributes: tfsdk.SingleNestedAttributes(map[string]tfsdk.Attribute{
							"cacertfile": {
								Type:     types.StringType,
								Required: true,
							},
							"certfile": {
								Type:     types.StringType,
								Required: true,
							},
							"fail_if_no_peer_cert": {
								Type:     types.BoolType,
								Required: true,
							},
							"keyfile": {
								Type:     types.StringType,
								Required: true,
							},
							"server_name_indication": {
								Type:     types.StringType,
								Required: true,
							},
							"verify": {
								Type:     types.StringType,
								Required: true,
							},
						}),
					},
					"user": {
						Type:     types.StringType,
						Required: true,
					},
					"virtual_host": {
						Type:     types.StringType,
						Required: true,
					},
				}),
			},
			"connection_type": {
				MarkdownDescription: "Connection type",
				Required:            true,
				Type:                types.StringType,
			},
			"direction": {
				MarkdownDescription: "Connection direction",
				Required:            true,
				Type:                types.StringType,
			},
			"enabled": {
				MarkdownDescription: "Connection enabled",
				Required:            true,
				Type:                types.BoolType,
			},
			"name": {
				MarkdownDescription: "Link name",
				Required:            true,
				Type:                types.StringType,
			},
			"id": {
				Computed:            true,
				MarkdownDescription: "Link identifier",
				PlanModifiers: tfsdk.AttributePlanModifiers{
					tfsdk.UseStateForUnknown(),
				},
				Type: types.StringType,
			},
		},
	}, nil
}

func (t linkResourceType) NewResource(ctx context.Context, in tfsdk.Provider) (tfsdk.Resource, diag.Diagnostics) {
	provider, diags := convertProviderType(in)

	return linkResource{
		provider: provider,
	}, diags
}

//type linkResourceData struct {
//	AddressHandling string       `tfsdk:"address_handling"`
//	Connection      Connection   `tfsdk:"connection"`
//	ConnectionType  string       `tfsdk:"connection_type"`
//	Direction       string       `tfsdk:"direction"`
//	Enabled         bool         `tfsdk:"enabled"`
//	ID              types.String `tfsdk:"id"`
//	Name            string       `tfsdk:"name"`
//}

type linkResourceData struct {
	AddressHandling types.String            `tfsdk:"address_handling"`
	Connection      *connectionResourceData `tfsdk:"connection"`
	ConnectionType  types.String            `tfsdk:"connection_type"`
	Direction       types.String            `tfsdk:"direction"`
	Enabled         types.Bool              `tfsdk:"enabled"`
	ID              types.String            `tfsdk:"id"`
	Name            types.String            `tfsdk:"name"`
}

type connectionResourceData struct {
	ApiPort     types.Int64            `tfsdk:"api_port"`
	Host        types.String           `tfsdk:"host"`
	Password    types.String           `tfsdk:"password"`
	Port        types.Int64            `tfsdk:"port"`
	SSLOptions  *sslOptionsResouceData `tfsdk:"ssl_options"`
	User        types.String           `tfsdk:"user"`
	VirtualHost types.String           `tfsdk:"virtual_host"`
}

type sslOptionsResouceData struct {
	CaCertFile           types.String `tfsdk:"cacertfile"`
	CertFile             types.String `tfsdk:"certfile"`
	FailIfNoPeerCert     types.Bool   `tfsdk:"fail_if_no_peer_cert"`
	KeyFile              types.String `tfsdk:"keyfile"`
	ServerNameIndication types.String `tfsdk:"server_name_indication"`
	Verify               types.String `tfsdk:"verify"`
}

type linkResource struct {
	provider provider
}

func LinkToCreateRequest(plan linkResourceData) api.LinkCreateRequest {
	newLink := api.LinkCreateRequest{
		AddressHandling: plan.AddressHandling.Value,
		Connection: &api.Connection{
			ApiPort:  int(plan.Connection.ApiPort.Value),
			Host:     plan.Connection.Host.Value,
			Password: plan.Connection.Password.Value,
			Port:     int(plan.Connection.Port.Value),
			SSLOptions: &api.SSLOptions{
				CaCertFile:           plan.Connection.SSLOptions.CaCertFile.Value,
				CertFile:             plan.Connection.SSLOptions.CertFile.Value,
				FailIfNoPeerCert:     plan.Connection.SSLOptions.FailIfNoPeerCert.Value,
				KeyFile:              plan.Connection.SSLOptions.KeyFile.Value,
				ServerNameIndication: plan.Connection.SSLOptions.ServerNameIndication.Value,
				Verify:               plan.Connection.SSLOptions.Verify.Value,
			},
			User:        plan.Connection.User.Value,
			VirtualHost: plan.Connection.VirtualHost.Value,
		},
		ConnectionType: plan.ConnectionType.Value,
		Direction:      plan.Direction.Value,
		Enabled:        plan.Enabled.Value,
		Name:           plan.Name.Value,
	}
	return newLink
}

func LinkToUpdateRequest(plan linkResourceData) api.LinkUpdateRequest {
	newLink := api.LinkUpdateRequest{
		AddressHandling: plan.AddressHandling.Value,
		Connection: &api.Connection{
			ApiPort:  int(plan.Connection.ApiPort.Value),
			Host:     plan.Connection.Host.Value,
			Password: plan.Connection.Password.Value,
			Port:     int(plan.Connection.Port.Value),
			SSLOptions: &api.SSLOptions{
				CaCertFile:           plan.Connection.SSLOptions.CaCertFile.Value,
				CertFile:             plan.Connection.SSLOptions.CertFile.Value,
				FailIfNoPeerCert:     plan.Connection.SSLOptions.FailIfNoPeerCert.Value,
				KeyFile:              plan.Connection.SSLOptions.KeyFile.Value,
				ServerNameIndication: plan.Connection.SSLOptions.ServerNameIndication.Value,
				Verify:               plan.Connection.SSLOptions.Verify.Value,
			},
			User:        plan.Connection.User.Value,
			VirtualHost: plan.Connection.VirtualHost.Value,
		},
		ConnectionType: plan.ConnectionType.Value,
		Direction:      plan.Direction.Value,
		Enabled:        plan.Enabled.Value,
		Name:           plan.Name.Value,
	}
	return newLink
}

func ApiToLink(link api.Link) Link {
	newLink := Link{
		AddressHandling: types.String{Value: link.AddressHandling},
		Connection: &Connection{
			ApiPort:  types.Int64{Value: int64(link.Connection.ApiPort)},
			Host:     types.String{Value: link.Connection.Host},
			Password: types.String{Value: link.Connection.Password},
			Port:     types.Int64{Value: int64(link.Connection.Port)},
			SSLOptions: &SSLOptions{
				CaCertFile:           types.String{Value: link.Connection.SSLOptions.CaCertFile},
				CertFile:             types.String{Value: link.Connection.SSLOptions.CertFile},
				FailIfNoPeerCert:     types.Bool{Value: link.Connection.SSLOptions.FailIfNoPeerCert},
				KeyFile:              types.String{Value: link.Connection.SSLOptions.KeyFile},
				ServerNameIndication: types.String{Value: link.Connection.SSLOptions.ServerNameIndication},
				Verify:               types.String{Value: link.Connection.SSLOptions.Verify},
			},
			User:        types.String{Value: link.Connection.User},
			VirtualHost: types.String{Value: link.Connection.VirtualHost},
		},
		ConnectionType: types.String{Value: link.ConnectionType},
		Direction:      types.String{Value: link.Direction},
		Enabled:        types.Bool{Value: link.Enabled},
		Name:           types.String{Value: link.Name},
	}
	return newLink
}

func (r linkResource) Create(ctx context.Context, req tfsdk.CreateResourceRequest, resp *tfsdk.CreateResourceResponse) {
	tflog.Debug(ctx, "Create 1\n")
	var state Link

	var plan linkResourceData
	diags := req.Plan.Get(ctx, &plan)
	tflog.Debug(ctx, fmt.Sprintf("plan: %+v\n", plan))
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}
	newLink := LinkToCreateRequest(plan)
	tflog.Debug(ctx, fmt.Sprintf("ZZZZZZZZZZZZZZZZZZZZZZZZZ NewLink: %+v\n", newLink))
	link, statusCode, err := r.provider.client.CreateLink(newLink, nil)
	log.Printf(fmt.Sprintf("ZZZZZZZZZZZZZZZZZZZZZZZZZZZZZZ link: %+v\n", link))
	tflog.Debug(ctx, fmt.Sprintf("ZZZZZZZZZZZZZZZZZZZZZZZZZ link: %+v\n", link))
	state = ApiToLink(link)
	if err != nil {
		resp.Diagnostics.AddError("Client Error", fmt.Sprintf("Unable to create link, got error: %s", err))
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

func (r linkResource) Read(ctx context.Context, req tfsdk.ReadResourceRequest, resp *tfsdk.ReadResourceResponse) {
	tflog.Debug(ctx, "Read 1\n")
	var state Link

	diags := req.State.Get(ctx, &state)
	tflog.Debug(ctx, fmt.Sprintf("state: %+v\n", state))
	resp.Diagnostics.Append(diags...)

	if resp.Diagnostics.HasError() {
		return
	}

	linkID := state.ID.Value

	link, statusCode, err := r.provider.client.GetLink(linkID, nil)
	if (statusCode < 200 || statusCode > 299) && statusCode != 404 {
		resp.Diagnostics.AddError("Client Unsuccesful", fmt.Sprintf("Status Code: %d", statusCode))
		return
	}
	if err != nil {
		resp.Diagnostics.AddError("Client Error", fmt.Sprintf("Unable to read link, got error: %s", err))
		return
	}
	state = ApiToLink(link)
	diags = resp.State.Set(ctx, &state)
	resp.Diagnostics.Append(diags...)
}

func (r linkResource) Update(ctx context.Context, req tfsdk.UpdateResourceRequest, resp *tfsdk.UpdateResourceResponse) {
	tflog.Debug(ctx, "Update 1\n")
	var state Link

	diags := req.State.Get(ctx, &state)
	resp.Diagnostics.Append(diags...)

	if resp.Diagnostics.HasError() {
		return
	}

	linkID := state.ID.Value

	var plan linkResourceData
	diags = req.Plan.Get(ctx, &plan)
	resp.Diagnostics.Append(diags...)
	//	resp.Diagnostics.AddError("Client Error", fmt.Sprintf("client: %+v\n", r.provider.client))
	if resp.Diagnostics.HasError() {
		return
	}
	newLink := LinkToUpdateRequest(plan)

	link, statusCode, err := r.provider.client.UpdateLink(linkID, newLink, nil)
	log.Printf(fmt.Sprintf("link: %+v\n", link))
	tflog.Debug(ctx, fmt.Sprintf("link: %+v\n", link))
	//resp.Diagnostics.AddError("link", fmt.Sprintf("link: %+v\n", link))
	state = ApiToLink(link)
	if err != nil {
		resp.Diagnostics.AddError("Client Error", fmt.Sprintf("Unable to create link, got error: %s", err))
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

func (r linkResource) Delete(ctx context.Context, req tfsdk.DeleteResourceRequest, resp *tfsdk.DeleteResourceResponse) {
	tflog.Debug(ctx, "Delete 1\n")
	var state Link

	diags := req.State.Get(ctx, &state)
	resp.Diagnostics.Append(diags...)

	if resp.Diagnostics.HasError() {
		return
	}

	linkID := state.ID.Value
	link, statusCode, err := r.provider.client.DeleteLink(linkID, nil)
	tflog.Debug(ctx, fmt.Sprintf("type of statusCode is %T\n", statusCode))
	tflog.Debug(ctx, fmt.Sprintf("statusCode, err: %d, %+v\n", statusCode, err))
	if statusCode < 200 || statusCode > 299 {
		resp.Diagnostics.AddError("Client Unsuccesful", fmt.Sprintf("Status Code: %d", statusCode))
		return
	}
	if err != nil {
		resp.Diagnostics.AddError("Client Error", fmt.Sprintf("Unable to read link, got error: %s", err))
		return
	}
	tflog.Debug(ctx, fmt.Sprintf("Link deleted: %+v\n", link))

	diags = resp.State.Set(ctx, &state)
	resp.Diagnostics.Append(diags...)
	resp.State.RemoveResource(ctx)
}

func (r linkResource) ImportState(ctx context.Context, req tfsdk.ImportResourceStateRequest, resp *tfsdk.ImportResourceStateResponse) {
	tfsdk.ResourceImportStatePassthroughID(ctx, tftypes.NewAttributePath().WithAttributeName("id"), req, resp)
}
