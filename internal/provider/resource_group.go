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
	groupResource struct {
		p dogProvider
	}
)

var (
	_ resource.Resource                = (*groupResource)(nil)
	_ resource.ResourceWithImportState = (*groupResource)(nil)
)

func NewGroupResource() resource.Resource {
	return &groupResource{}
}

func (*groupResource) Metadata(ctx context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_group"
}


func (*groupResource) GetSchema(ctx context.Context) (tfsdk.Schema, diag.Diagnostics) {
	return tfsdk.Schema{
		Attributes: map[string]tfsdk.Attribute{
			// This description is used by the documentation generator and the language server.
			"description": {
				MarkdownDescription: "group description",
				Optional:            true,
				Type:                types.StringType,
			},
			"name": {
				MarkdownDescription: "group name",
				Required:            true,
				Type:                types.StringType,
			},
			"profile_name": {
				MarkdownDescription: "group profile name",
				Required:            true,
				Type:                types.StringType,
			},
			"profile_version": {
				MarkdownDescription: "group profile version",
				Required:            true,
				Type:                types.StringType,
			},
			"id": {
				Computed:            true,
				MarkdownDescription: "group identifier",
				PlanModifiers: tfsdk.AttributePlanModifiers{
					resource.UseStateForUnknown(),
				},
				Type: types.StringType,
			},
		},
	}, nil
}

func (r *groupResource) Configure(ctx context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {
	// Prevent panic if the provider has not been configured
	if req.ProviderData == nil {
		return
	}

	client, ok := req.ProviderData.(*api.Client)
	if !ok {
		resp.Diagnostics.AddError(
			"Unexpected Resource Configure Type",
			fmt.Sprintf("Expected *dog.Client, got: %T. Please report this issue to the provider developers.", req.ProviderData),
		)

		return
	}

	r.p.dog = client
}

func (*groupResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	resource.ImportStatePassthroughID(ctx, path.Root("id"), req, resp)
}

type groupResourceData struct {
	//Created        int    `json:"created,omitempty"` //TODO: created has both int and string entries
	Description    string       `tfsdk:"description"`
	ID             types.String `tfsdk:"id"`
	Name           string       `tfsdk:"name"`
	ProfileName    string       `tfsdk:"profile_name"`
	ProfileVersion string       `tfsdk:"profile_version"`
}

func GroupToCreateRequest(plan groupResourceData) api.GroupCreateRequest {
	newGroup := api.GroupCreateRequest{
		Description:    plan.Description,
		Name:           plan.Name,
		ProfileName:    plan.ProfileName,
		ProfileVersion: plan.ProfileVersion,
	}
	return newGroup
}

func GroupToUpdateRequest(plan groupResourceData) api.GroupUpdateRequest {
	newGroup := api.GroupUpdateRequest{
		Description:    plan.Description,
		Name:           plan.Name,
		ProfileName:    plan.ProfileName,
		ProfileVersion: plan.ProfileVersion,
	}
	return newGroup
}

func ApiToGroup(group api.Group) Group {
	h := Group{
		Description:    types.String{Value: group.Description},
		ID:             types.String{Value: group.ID},
		Name:           types.String{Value: group.Name},
		ProfileName:    types.String{Value: group.ProfileName},
		ProfileVersion: types.String{Value: group.ProfileVersion},
	}
	return h
}

func (r *groupResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	var state Group

	var plan groupResourceData
	diags := req.Plan.Get(ctx, &plan)
	resp.Diagnostics.Append(diags...)
	//	resp.Diagnostics.AddError("Client Error", fmt.Sprintf("client: %+v\n", r.provider.client))
	if resp.Diagnostics.HasError() {
		return
	}

	newGroup := GroupToCreateRequest(plan)
	group, statusCode, err := r.p.dog.CreateGroup(newGroup, nil)
	log.Printf(fmt.Sprintf("group: %+v\n", group))
	tflog.Trace(ctx, fmt.Sprintf("group: %+v\n", group))
	state = ApiToGroup(group)
	if err != nil {
		resp.Diagnostics.AddError("Client Error", fmt.Sprintf("Unable to create group, got error: %s", err))
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

func (r *groupResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	var state Group

	diags := req.State.Get(ctx, &state)
	resp.Diagnostics.Append(diags...)

	if resp.Diagnostics.HasError() {
		return
	}

	groupID := state.ID.Value

	group, statusCode, err := r.p.dog.GetGroup(groupID, nil)
	if (statusCode < 200 || statusCode > 299) && statusCode != 404 {
		resp.Diagnostics.AddError("Client Unsuccesful", fmt.Sprintf("Status Code: %d", statusCode))
		return
	}
	if err != nil {
		resp.Diagnostics.AddError("Client Error", fmt.Sprintf("Unable to read group, got error: %s", err))
		return
	}
	state = ApiToGroup(group)
	diags = resp.State.Set(ctx, &state)
	resp.Diagnostics.Append(diags...)
}


func (r *groupResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	//var data groupResourceData

	//diags := req.Plan.Get(ctx, &data)
	//resp.Diagnostics.Append(diags...)

	//if resp.Diagnostics.HasError() {
	//	return
	//}
	var state Group

	diags := req.State.Get(ctx, &state)
	resp.Diagnostics.Append(diags...)

	if resp.Diagnostics.HasError() {
		return
	}

	groupID := state.ID.Value

	var plan groupResourceData
	diags = req.Plan.Get(ctx, &plan)
	resp.Diagnostics.Append(diags...)
	//	resp.Diagnostics.AddError("Client Error", fmt.Sprintf("client: %+v\n", r.provider.client))
	if resp.Diagnostics.HasError() {
		return
	}

	newGroup := GroupToUpdateRequest(plan)
	group, statusCode, err := r.p.dog.UpdateGroup(groupID, newGroup, nil)
	log.Printf(fmt.Sprintf("group: %+v\n", group))
	tflog.Trace(ctx, fmt.Sprintf("group: %+v\n", group))
	//resp.Diagnostics.AddError("group", fmt.Sprintf("group: %+v\n", group))
	state = ApiToGroup(group)
	if err != nil {
		resp.Diagnostics.AddError("Client Error", fmt.Sprintf("Unable to create group, got error: %s", err))
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

func (r *groupResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	var state Group

	diags := req.State.Get(ctx, &state)
	resp.Diagnostics.Append(diags...)

	if resp.Diagnostics.HasError() {
		return
	}

	groupID := state.ID.Value
	group, statusCode, err := r.p.dog.DeleteGroup(groupID, nil)
	if statusCode < 200 || statusCode > 299 {
		resp.Diagnostics.AddError("Client Unsuccesful", fmt.Sprintf("Status Code: %d", statusCode))
		return
	}
	if err != nil {
		resp.Diagnostics.AddError("Client Error", fmt.Sprintf("Unable to read group, got error: %s", err))
		return
	}
	tflog.Trace(ctx, fmt.Sprintf("group deleted: %+v\n", group))

	diags = resp.State.Set(ctx, &state)
	resp.Diagnostics.Append(diags...)
	resp.State.RemoveResource(ctx)
}
