package dog

import (
	"context"
	"fmt"
	"regexp"

	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-framework/schema/validator"
	"github.com/hashicorp/terraform-plugin-framework-validators/stringvalidator"
	"github.com/hashicorp/terraform-plugin-log/tflog"
	api "github.com/relaypro-open/dog_api_golang/api"
	"golang.org/x/exp/slices"
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

func (*hostResource) Schema(ctx context.Context, req resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = schema.Schema{
		// This description is used by the documentation generator and the language server.
		MarkdownDescription: "Host data source",

		Attributes: map[string]schema.Attribute{
			// This description is used by the documentation generator and the language server.
			"environment": schema.StringAttribute{
				MarkdownDescription: "Host environment",
				Optional:            true,
			},
			"group": schema.StringAttribute{
				MarkdownDescription: "Host group",
				Optional:            true,
			},
			"hostkey": schema.StringAttribute{
				MarkdownDescription: "Host key",
				Optional:            true,
				Validators: []validator.String{
					stringvalidator.LengthBetween(10, 256),
					stringvalidator.RegexMatches(
						regexp.MustCompile(`^[A-Za-z0-9+%_.-](.*)$`),
						"must start with alphanumeric characters, %, _, ., -",
					),
				},
			},
			"location": schema.StringAttribute{
				MarkdownDescription: "Host location",
				Optional:            true,
			},
			"name": schema.StringAttribute{
				MarkdownDescription: "Host name",
				Optional:            true,
			},
			"vars": schema.StringAttribute{
				MarkdownDescription: "json string of vars",
				Optional:            true,
				//Required:            true,
			},
			"id": schema.StringAttribute{
				Optional:            true,
				MarkdownDescription: "Host identifier",
				Computed: true,
			},
		},
	}
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
	Environment string            `tfsdk:"environment"`
	Group       string            `tfsdk:"group"`
	ID          types.String      `tfsdk:"id"`
	HostKey     string            `tfsdk:"hostkey"`
	Location    string            `tfsdk:"location"`
	Name        string            `tfsdk:"name"`
	Vars        *string           `tfsdk:"vars"`
}

func HostToApiHost(plan Host) api.Host {
	if plan.Vars == nil {
		newHost := api.Host{
			Environment: plan.Environment.ValueString(),
			Group:       plan.Group.ValueString(),
			HostKey:     plan.HostKey.ValueString(),
			Location:    plan.Location.ValueString(),
			Name:        plan.Name.ValueString(),
		}
		return newHost
	} else {
		newHost := api.Host{
			Environment: plan.Environment.ValueString(),
			Group:       plan.Group.ValueString(),
			HostKey:     plan.HostKey.ValueString(),
			Location:    plan.Location.ValueString(),
			Name:        plan.Name.ValueString(),
			Vars:        *plan.Vars,
		}
		return newHost
	}
}

func ApiToHost(host api.Host) Host {
	if host.Vars == "" {
		h := Host{
			Environment: types.StringValue(host.Environment),
			Group:       types.StringValue(host.Group),
			ID:          types.StringValue(host.ID),
			HostKey:     types.StringValue(host.HostKey),
			Location:    types.StringValue(host.Location),
			Name:        types.StringValue(host.Name),
		}
		return h
	} else {
		h := Host{
			Environment: types.StringValue(host.Environment),
			Group:       types.StringValue(host.Group),
			ID:          types.StringValue(host.ID),
			HostKey:     types.StringValue(host.HostKey),
			Location:    types.StringValue(host.Location),
			Name:        types.StringValue(host.Name),
			Vars:        &host.Vars,
		}
		return h
	}
}

func HostToCreateRequest(plan hostResourceData) api.HostCreateRequest {
	newHost := api.HostCreateRequest{
		Environment: plan.Environment,
		Group:       plan.Group,
		HostKey:     plan.HostKey,
		Location:    plan.Location,
		Name:        plan.Name,
		Vars:        *plan.Vars,
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
		Vars:        *plan.Vars,
	}
	return newHost
}

func (r *hostResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	var state Host

	//var plan hostResourceData
	var plan Host
	diags := req.Plan.Get(ctx, &plan)
	resp.Diagnostics.Append(diags...)
	//	resp.Diagnostics.AddError("Client Error", fmt.Sprintf("client: %+v\n", r.provider.client))
	if resp.Diagnostics.HasError() {
		return
	}

	tflog.Trace(ctx,fmt.Sprintf("plan: %+v\n", plan))
	newHost := HostToApiHost(plan)
	tflog.Trace(ctx,fmt.Sprintf("newHost: %+v\n", newHost))
	tflog.Trace(ctx,fmt.Sprintf("r.p.dog: %+v\n", r.p.dog))
	host, statusCode, err := r.p.dog.CreateHostEncode(newHost, nil)
	tflog.Trace(ctx,fmt.Sprintf("host: %+v\n", host))
	if err != nil {
		resp.Diagnostics.AddError("Client Error", fmt.Sprintf("Unable to create host, got error: %s", err))
	}
	if statusCode != 201 {
		resp.Diagnostics.AddError("Client Unsuccessful", fmt.Sprintf("Status Code: %d", statusCode))
	}
	if resp.Diagnostics.HasError() {
		return
	}
	state = ApiToHost(host)
	tflog.Trace(ctx,fmt.Sprintf("state: %+v\n", state))

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

	hostID := state.ID.ValueString()

	tflog.Trace(ctx, fmt.Sprintf("r.p: %+v\n", r.p))
	tflog.Trace(ctx, fmt.Sprintf("r.p.dog: %+v\n", r.p.dog))
	host, statusCode, err := r.p.dog.GetHostEncode(hostID, nil)
	if statusCode != 200 {
		resp.Diagnostics.AddError("Client Unsuccessful", fmt.Sprintf("Status Code: %d", statusCode))
	}
	if err != nil {
		resp.Diagnostics.AddError("Client Error", fmt.Sprintf("Unable to read host, got error: %s", err))
	}
	if resp.Diagnostics.HasError() {
		return
	}
	state = ApiToHost(host)
	diags = resp.State.Set(ctx, &state)
	resp.Diagnostics.Append(diags...)
}

func (r *hostResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	var state Host

	diags := req.State.Get(ctx, &state)
	resp.Diagnostics.Append(diags...)

	if resp.Diagnostics.HasError() {
		return
	}

	hostID := state.ID.ValueString()

	var plan Host
	diags = req.Plan.Get(ctx, &plan)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	newHost := HostToApiHost(plan)
	host, statusCode, err := r.p.dog.UpdateHostEncode(hostID, newHost, nil)
	tflog.Trace(ctx, fmt.Sprintf("host: %+v\n", host))
	state = ApiToHost(host)
	if err != nil {
		resp.Diagnostics.AddError("Client Error", fmt.Sprintf("Unable to create host, got error: %s", err))
	}
	ok := []int{303, 200, 201}
	if slices.Contains(ok, statusCode) != true {
		resp.Diagnostics.AddError("Client Unsuccessful", fmt.Sprintf("Status Code: %d", statusCode))
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

func (r *hostResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	var state Host

	diags := req.State.Get(ctx, &state)
	resp.Diagnostics.Append(diags...)

	if resp.Diagnostics.HasError() {
		return
	}

	hostID := state.ID.ValueString()
	host, statusCode, err := r.p.dog.DeleteHost(hostID, nil)
	if statusCode != 204 {
		resp.Diagnostics.AddError("Client Unsuccessful", fmt.Sprintf("Status Code: %d", statusCode))
	}
	if err != nil {
		resp.Diagnostics.AddError("Client Error", fmt.Sprintf("Unable to read host, got error: %s", err))
	}
	if resp.Diagnostics.HasError() {
		return
	}
	tflog.Trace(ctx, fmt.Sprintf("host deleted: %+v\n", host))

	diags = resp.State.Set(ctx, &state)
	resp.Diagnostics.Append(diags...)
	resp.State.RemoveResource(ctx)
}
