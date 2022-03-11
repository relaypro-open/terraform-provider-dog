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

type groupResourceType struct{}

func (t groupResourceType) GetSchema(ctx context.Context) (tfsdk.Schema, diag.Diagnostics) {
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
					tfsdk.UseStateForUnknown(),
				},
				Type: types.StringType,
			},
		},
	}, nil
}

func (t groupResourceType) NewResource(ctx context.Context, in tfsdk.Provider) (tfsdk.Resource, diag.Diagnostics) {
	provider, diags := convertProviderType(in)

	return groupResource{
		provider: provider,
	}, diags
}

type groupResourceData struct {
	//Created        int    `json:"created,omitempty"` //TODO: created has both int and string entries
	Description    string       `tfsdk:"description"`
	ID             types.String `tfsdk:"id"`
	Name           string       `tfsdk:"name"`
	ProfileName    string       `tfsdk:"profile_name"`
	ProfileVersion string       `tfsdk:"profile_version"`
}

type groupResource struct {
	provider provider
}

func (r groupResource) Create(ctx context.Context, req tfsdk.CreateResourceRequest, resp *tfsdk.CreateResourceResponse) {
	var state Group

	var plan groupResourceData
	diags := req.Plan.Get(ctx, &plan)
	resp.Diagnostics.Append(diags...)
	//	resp.Diagnostics.AddError("Client Error", fmt.Sprintf("client: %+v\n", r.provider.client))
	if resp.Diagnostics.HasError() {
		return
	}

	newGroup := api.GroupCreateRequest{
		Description:    plan.Description,
		Name:           plan.Name,
		ProfileName:    plan.ProfileName,
		ProfileVersion: plan.ProfileVersion,
	}

	group, statusCode, err := r.provider.client.CreateGroup(newGroup, nil)
	log.Printf(fmt.Sprintf("group: %+v\n", group))
	tflog.Trace(ctx, fmt.Sprintf("group: %+v\n", group))
	//resp.Diagnostics.AddError("group", fmt.Sprintf("group: %+v\n", group))
	h := Group{
		Description:    types.String{Value: group.Description},
		ID:             types.String{Value: group.ID},
		Name:           types.String{Value: group.Name},
		ProfileName:    types.String{Value: group.ProfileName},
		ProfileVersion: types.String{Value: group.ProfileVersion},
	}
	state = h
	if err != nil {
		resp.Diagnostics.AddError("Client Error", fmt.Sprintf("Unable to create group, got error: %s", err))
		return
	}
	if statusCode < 200 && statusCode > 299 {
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

func (r groupResource) Read(ctx context.Context, req tfsdk.ReadResourceRequest, resp *tfsdk.ReadResourceResponse) {
	var state Group

	diags := req.State.Get(ctx, &state)
	resp.Diagnostics.Append(diags...)

	if resp.Diagnostics.HasError() {
		return
	}

	groupID := state.ID.Value

	group, statusCode, err := r.provider.client.GetGroup(groupID, nil)
	if statusCode < 200 && statusCode > 299 {
		resp.Diagnostics.AddError("Client Unsuccesful", fmt.Sprintf("Status Code: %d", statusCode))
		return
	}
	if err != nil {
		resp.Diagnostics.AddError("Client Error", fmt.Sprintf("Unable to read group, got error: %s", err))
		return
	}
	h := Group{
		Description:    types.String{Value: group.Description},
		ID:             types.String{Value: group.ID},
		Name:           types.String{Value: group.Name},
		ProfileName:    types.String{Value: group.ProfileName},
		ProfileVersion: types.String{Value: group.ProfileVersion},
	}

	state = h
	diags = resp.State.Set(ctx, &state)
	resp.Diagnostics.Append(diags...)
}

func (r groupResource) Update(ctx context.Context, req tfsdk.UpdateResourceRequest, resp *tfsdk.UpdateResourceResponse) {
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

	newGroup := api.GroupUpdateRequest{
		Description:    plan.Description,
		Name:           plan.Name,
		ProfileName:    plan.ProfileName,
		ProfileVersion: plan.ProfileVersion,
	}
	group, statusCode, err := r.provider.client.UpdateGroup(groupID, newGroup, nil)
	log.Printf(fmt.Sprintf("group: %+v\n", group))
	tflog.Trace(ctx, fmt.Sprintf("group: %+v\n", group))
	//resp.Diagnostics.AddError("group", fmt.Sprintf("group: %+v\n", group))
	h := Group{
		Description:    types.String{Value: group.Description},
		ID:             types.String{Value: group.ID},
		Name:           types.String{Value: group.Name},
		ProfileName:    types.String{Value: group.ProfileName},
		ProfileVersion: types.String{Value: group.ProfileVersion},
	}
	state = h
	if err != nil {
		resp.Diagnostics.AddError("Client Error", fmt.Sprintf("Unable to create group, got error: %s", err))
		return
	}
	if statusCode < 200 && statusCode > 299 {
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

func (r groupResource) Delete(ctx context.Context, req tfsdk.DeleteResourceRequest, resp *tfsdk.DeleteResourceResponse) {
	var state Group

	diags := req.State.Get(ctx, &state)
	resp.Diagnostics.Append(diags...)

	if resp.Diagnostics.HasError() {
		return
	}

	groupID := state.ID.Value
	group, statusCode, err := r.provider.client.DeleteGroup(groupID, nil)
	if statusCode < 200 && statusCode > 299 {
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

func (r groupResource) ImportState(ctx context.Context, req tfsdk.ImportResourceStateRequest, resp *tfsdk.ImportResourceStateResponse) {
	tfsdk.ResourceImportStatePassthroughID(ctx, tftypes.NewAttributePath().WithAttributeName("id"), req, resp)
}
