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
	linkResource struct {
		p dogProvider
	}
)

var (
	_ resource.Resource                = (*linkResource)(nil)
	_ resource.ResourceWithImportState = (*linkResource)(nil)
)

func NewLinkResource() resource.Resource {
	return &linkResource{}
}

func (*linkResource) Metadata(ctx context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_link"
}

func (*linkResource) Schema(ctx context.Context, req resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = schema.Schema{
		// This description is used by the documentation generator and the language server.
		MarkdownDescription: "Link data source",

		Attributes: map[string]schema.Attribute{
			// This description is used by the documentation generator and the language server.
			"address_handling": schema.StringAttribute{
				MarkdownDescription: "Type of address handling",
				Optional:            true,
			},
			"dog_connection": schema.SingleNestedAttribute{
				MarkdownDescription: "Connection specification",
				Optional:            true,
				Attributes: map[string]schema.Attribute{
					"api_port": schema.Int64Attribute{
						Required: true,
					},
					"host": schema.StringAttribute{
						Required: true,
					},
					"password": schema.StringAttribute{
						Required:  true,
						Sensitive: true,
					},
					"port": schema.Int64Attribute{
						Required: true,
					},
					"ssl_options": schema.SingleNestedAttribute{
						Required: true,
						Attributes: map[string]schema.Attribute{
							"cacertfile": schema.StringAttribute{
								Required: true,
							},
							"certfile": schema.StringAttribute{
								Required: true,
							},
							"fail_if_no_peer_cert": schema.BoolAttribute{
								Required: true,
							},
							"keyfile": schema.StringAttribute{
								Required: true,
							},
							"server_name_indication": schema.StringAttribute{
								Required: true,
							},
							"verify": schema.StringAttribute{
								Required: true,
							},
						},
					},
					"user": schema.StringAttribute{
						Required: true,
					},
					"virtual_host": schema.StringAttribute{
						Required: true,
					},
				},
			},
			"connection_type": schema.StringAttribute{
				MarkdownDescription: "Connection type",
				Optional:            true,
			},
			"direction": schema.StringAttribute{
				MarkdownDescription: "Connection direction",
				Optional:            true,
			},
			"enabled": schema.BoolAttribute{
				MarkdownDescription: "Connection enabled",
				Optional:            true,
			},
			"name": schema.StringAttribute{
				MarkdownDescription: "Link name",
				Optional:            true,
			},
			"id": schema.StringAttribute{
				Optional:            true,
				MarkdownDescription: "Link identifier",
				Computed:            true,
			},
		},
	}
}

func (r *linkResource) Configure(ctx context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {
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

func (*linkResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	resource.ImportStatePassthroughID(ctx, path.Root("id"), req, resp)
}

type linkResourceData struct {
	AddressHandling types.String            `tfsdk:"address_handling"`
	Connection      *connectionResourceData `tfsdk:"dog_connection"`
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

func LinkToCreateRequest(plan linkResourceData) api.LinkCreateRequest {
	newLink := api.LinkCreateRequest{
		AddressHandling: plan.AddressHandling.ValueString(),
		Connection: &api.Connection{
			ApiPort:  int(plan.Connection.ApiPort.ValueInt64()),
			Host:     plan.Connection.Host.ValueString(),
			Password: plan.Connection.Password.ValueString(),
			Port:     int(plan.Connection.Port.ValueInt64()),
			SSLOptions: &api.SSLOptions{
				CaCertFile:           plan.Connection.SSLOptions.CaCertFile.ValueString(),
				CertFile:             plan.Connection.SSLOptions.CertFile.ValueString(),
				FailIfNoPeerCert:     plan.Connection.SSLOptions.FailIfNoPeerCert.ValueBool(),
				KeyFile:              plan.Connection.SSLOptions.KeyFile.ValueString(),
				ServerNameIndication: plan.Connection.SSLOptions.ServerNameIndication.ValueString(),
				Verify:               plan.Connection.SSLOptions.Verify.ValueString(),
			},
			User:        plan.Connection.User.ValueString(),
			VirtualHost: plan.Connection.VirtualHost.ValueString(),
		},
		ConnectionType: plan.ConnectionType.ValueString(),
		Direction:      plan.Direction.ValueString(),
		Enabled:        plan.Enabled.ValueBool(),
		Name:           plan.Name.ValueString(),
	}
	return newLink
}

func LinkToUpdateRequest(plan linkResourceData) api.LinkUpdateRequest {
	newLink := api.LinkUpdateRequest{
		AddressHandling: plan.AddressHandling.ValueString(),
		Connection: &api.Connection{
			ApiPort:  int(plan.Connection.ApiPort.ValueInt64()),
			Host:     plan.Connection.Host.ValueString(),
			Password: plan.Connection.Password.ValueString(),
			Port:     int(plan.Connection.Port.ValueInt64()),
			SSLOptions: &api.SSLOptions{
				CaCertFile:           plan.Connection.SSLOptions.CaCertFile.ValueString(),
				CertFile:             plan.Connection.SSLOptions.CertFile.ValueString(),
				FailIfNoPeerCert:     plan.Connection.SSLOptions.FailIfNoPeerCert.ValueBool(),
				KeyFile:              plan.Connection.SSLOptions.KeyFile.ValueString(),
				ServerNameIndication: plan.Connection.SSLOptions.ServerNameIndication.ValueString(),
				Verify:               plan.Connection.SSLOptions.Verify.ValueString(),
			},
			User:        plan.Connection.User.ValueString(),
			VirtualHost: plan.Connection.VirtualHost.ValueString(),
		},
		ConnectionType: plan.ConnectionType.ValueString(),
		Direction:      plan.Direction.ValueString(),
		Enabled:        plan.Enabled.ValueBool(),
		Name:           plan.Name.ValueString(),
	}
	return newLink
}

func ApiToLink(link api.Link) Link {
	newLink := Link{
		AddressHandling: types.StringValue(link.AddressHandling),
		Connection: &Connection{
			ApiPort:  types.Int64Value(int64(link.Connection.ApiPort)),
			Host:     types.StringValue(link.Connection.Host),
			Password: types.StringValue(link.Connection.Password),
			Port:     types.Int64Value(int64(link.Connection.Port)),
			SSLOptions: &SSLOptions{
				CaCertFile:           types.StringValue(link.Connection.SSLOptions.CaCertFile),
				CertFile:             types.StringValue(link.Connection.SSLOptions.CertFile),
				FailIfNoPeerCert:     types.BoolValue(link.Connection.SSLOptions.FailIfNoPeerCert),
				KeyFile:              types.StringValue(link.Connection.SSLOptions.KeyFile),
				ServerNameIndication: types.StringValue(link.Connection.SSLOptions.ServerNameIndication),
				Verify:               types.StringValue(link.Connection.SSLOptions.Verify),
			},
			User:        types.StringValue(link.Connection.User),
			VirtualHost: types.StringValue(link.Connection.VirtualHost),
		},
		ConnectionType: types.StringValue(link.ConnectionType),
		Direction:      types.StringValue(link.Direction),
		Enabled:        types.BoolValue(link.Enabled),
		Name:           types.StringValue(link.Name),
		ID:             types.StringValue(link.ID),
	}
	return newLink
}

func (r *linkResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	var state Link

	var plan linkResourceData
	diags := req.Plan.Get(ctx, &plan)
	resp.Diagnostics.Append(diags...)
	//	resp.Diagnostics.AddError("Client Error", fmt.Sprintf("client: %+v\n", r.provider.client))
	if resp.Diagnostics.HasError() {
		return
	}

	newLink := LinkToCreateRequest(plan)
	log.Printf(fmt.Sprintf("r.p.dog: %+v\n", r.p.dog))
	link, statusCode, err := r.p.dog.CreateLink(newLink, nil)
	log.Printf(fmt.Sprintf("link: %+v\n", link))
	tflog.Trace(ctx, fmt.Sprintf("link: %+v\n", link))
	if err != nil {
		resp.Diagnostics.AddError("Client Error", fmt.Sprintf("Unable to create link, got error: %s", err))
	}
	if statusCode != 201 {
		resp.Diagnostics.AddError("Client Unsuccesful", fmt.Sprintf("Status Code: %d", statusCode))
	}
	if resp.Diagnostics.HasError() {
		return
	}
	state = ApiToLink(link)

	plan.ID = state.ID

	// write logs using the tflog package
	// see https://pkg.go.dev/github.com/hashicorp/terraform-plugin-log/tflog
	// for more information
	tflog.Trace(ctx, "created a resource")

	diags = resp.State.Set(ctx, &state)
	resp.Diagnostics.Append(diags...)
}

func (r *linkResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	var state Link

	diags := req.State.Get(ctx, &state)
	resp.Diagnostics.Append(diags...)

	if resp.Diagnostics.HasError() {
		return
	}

	linkID := state.ID.ValueString()

	log.Printf(fmt.Sprintf("r.p: %+v\n", r.p))
	log.Printf(fmt.Sprintf("r.p.dog: %+v\n", r.p.dog))
	link, statusCode, err := r.p.dog.GetLink(linkID, nil)
	if statusCode != 200 {
		resp.Diagnostics.AddError("Client Unsuccesful", fmt.Sprintf("Status Code: %d", statusCode))
	}
	if err != nil {
		resp.Diagnostics.AddError("Client Error", fmt.Sprintf("Unable to read link, got error: %s", err))
	}
	if resp.Diagnostics.HasError() {
		return
	}
	state = ApiToLink(link)
	diags = resp.State.Set(ctx, &state)
	resp.Diagnostics.Append(diags...)
}

func (r *linkResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	var state Link

	diags := req.State.Get(ctx, &state)
	resp.Diagnostics.Append(diags...)

	if resp.Diagnostics.HasError() {
		return
	}

	linkID := state.ID.ValueString()

	var plan linkResourceData
	diags = req.Plan.Get(ctx, &plan)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	newLink := LinkToUpdateRequest(plan)
	link, statusCode, err := r.p.dog.UpdateLink(linkID, newLink, nil)
	log.Printf(fmt.Sprintf("link: %+v\n", link))
	tflog.Trace(ctx, fmt.Sprintf("link: %+v\n", link))
	state = ApiToLink(link)
	if err != nil {
		resp.Diagnostics.AddError("Client Error", fmt.Sprintf("Unable to create link, got error: %s", err))
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

func (r *linkResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	var state Link

	diags := req.State.Get(ctx, &state)
	resp.Diagnostics.Append(diags...)

	if resp.Diagnostics.HasError() {
		return
	}

	linkID := state.ID.ValueString()
	link, statusCode, err := r.p.dog.DeleteLink(linkID, nil)
	if statusCode != 204 {
		resp.Diagnostics.AddError("Client Unsuccesful", fmt.Sprintf("Status Code: %d", statusCode))
	}
	if err != nil {
		resp.Diagnostics.AddError("Client Error", fmt.Sprintf("Unable to read link, got error: %s", err))
	}
	if resp.Diagnostics.HasError() {
		return
	}
	tflog.Trace(ctx, fmt.Sprintf("link deleted: %+v\n", link))

	diags = resp.State.Set(ctx, &state)
	resp.Diagnostics.Append(diags...)
	resp.State.RemoveResource(ctx)
}
