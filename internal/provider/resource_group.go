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
	"golang.org/x/exp/slices"
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
			"profile_id": {
				MarkdownDescription: "group profile id",
				Optional:            true,
				Type:                types.StringType,
			},
			"profile_name": {
				MarkdownDescription: "group profile name",
				Optional:            true,
				Type:                types.StringType,
			},
			"profile_version": {
				MarkdownDescription: "group profile version",
				Optional:            true,
				Type:                types.StringType,
			},
		   "ec2_security_group_ids": {
				   MarkdownDescription: "List of EC2 Security Groups to control",
				   Optional:            true,
				   Attributes: tfsdk.ListNestedAttributes(map[string]tfsdk.Attribute{
						   "region": {
								   MarkdownDescription: "EC2 Region",
								   Required:            true,
								   Type:                types.StringType,
						   },
						   "sgid": {
								   MarkdownDescription: "EC2 Security Group ID",
								   Required:            true,
								   Type:                types.StringType,
						   },
				   }),
			},
			"vars": {
				MarkdownDescription: "Arbitrary collection of variables used for inventory",
				Type:        types.MapType{ElemType: types.StringType},
				Optional:    true,
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
	ProfileId      string       `tfsdk:"profile_id"`
	ProfileName    string       `tfsdk:"profile_name"`
	ProfileVersion string       `tfsdk:"profile_version"`
	Ec2SecurityGroupIds []*ec2SecurityGroupIdsResourceData `tfsdk:"ec2_security_group_ids"`
	Vars           map[string]string	`tfsdk:"vars"`
}

type ec2SecurityGroupIdsResourceData struct {
	Region	string `tfsdk:"region"`
	SgId	string `tfsdk:"sgid"`
}

func GroupToCreateRequest(plan groupResourceData) api.GroupCreateRequest {
	newEc2SecurityGroupIds := []*api.Ec2SecurityGroupIds{}
	for _, region_sgid := range plan.Ec2SecurityGroupIds {
		rs := &api.Ec2SecurityGroupIds{
			Region:    region_sgid.Region,
			SgId:      region_sgid.SgId,
		}
		newEc2SecurityGroupIds = append(newEc2SecurityGroupIds, rs)
	}

	newGroup := api.GroupCreateRequest{
		Description:    plan.Description,
		Name:           plan.Name,
		ProfileName:    plan.ProfileName,
		ProfileVersion: plan.ProfileVersion,
		Ec2SecurityGroupIds: newEc2SecurityGroupIds,
		Vars:		    plan.Vars,
	}
	return newGroup
}

func GroupToUpdateRequest(plan groupResourceData) api.GroupUpdateRequest {
	newEc2SecurityGroupIds := []*api.Ec2SecurityGroupIds{}
	for _, region_sgid := range plan.Ec2SecurityGroupIds {
		rs := &api.Ec2SecurityGroupIds{
			Region:    region_sgid.Region,
			SgId:      region_sgid.SgId,
		}
		newEc2SecurityGroupIds = append(newEc2SecurityGroupIds, rs)
	}

	newGroup := api.GroupUpdateRequest{
		Description:    plan.Description,
		Name:           plan.Name,
		ProfileName:    plan.ProfileName,
		ProfileVersion: plan.ProfileVersion,
		Ec2SecurityGroupIds: newEc2SecurityGroupIds,
		Vars:		    plan.Vars,
	}
	return newGroup
}

func ApiToGroup(group api.Group) Group {
	newEc2SecurityGroupIds := []*Ec2SecurityGroupIds{}
	for _, region_sgid := range group.Ec2SecurityGroupIds {
		rs := &Ec2SecurityGroupIds{
			Region:    types.String{Value: region_sgid.Region},
			SgId:      types.String{Value: region_sgid.SgId},
		}
		newEc2SecurityGroupIds = append(newEc2SecurityGroupIds, rs)
	}

	newVars := map[string]string{}
	for k, v := range group.Vars {
		newVars[k] = v
	}

	h := Group{
		Description:    types.String{Value: group.Description},
		ID:             types.String{Value: group.ID},
		Name:           types.String{Value: group.Name},
		ProfileName:    types.String{Value: group.ProfileName},
		ProfileVersion: types.String{Value: group.ProfileVersion},
		Ec2SecurityGroupIds: newEc2SecurityGroupIds,
		Vars:		newVars,
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
	if err != nil {
		resp.Diagnostics.AddError("Client Error", fmt.Sprintf("Unable to create group, got error: %s", err))
	}
	if statusCode != 201 {
		resp.Diagnostics.AddError("Client Unsuccesful", fmt.Sprintf("Status Code: %d", statusCode))
	}
	if resp.Diagnostics.HasError() {
		return
	}
	state = ApiToGroup(group)

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
	if statusCode != 200 {
		resp.Diagnostics.AddError("Client Unsuccesful", fmt.Sprintf("Status Code: %d", statusCode))
	}
	if err != nil {
		resp.Diagnostics.AddError("Client Error", fmt.Sprintf("Unable to read group, got error: %s", err))
	}
	if resp.Diagnostics.HasError() {
		return
	}
	state = ApiToGroup(group)
	diags = resp.State.Set(ctx, &state)
	resp.Diagnostics.Append(diags...)
}


func (r *groupResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
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
	if resp.Diagnostics.HasError() {
		return
	}

	newGroup := GroupToUpdateRequest(plan)
	group, statusCode, err := r.p.dog.UpdateGroup(groupID, newGroup, nil)
	log.Printf(fmt.Sprintf("group: %+v\n", group))
	tflog.Trace(ctx, fmt.Sprintf("group: %+v\n", group))
	state = ApiToGroup(group)
	if err != nil {
		resp.Diagnostics.AddError("Client Error", fmt.Sprintf("Unable to create group, got error: %s", err))
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

func (r *groupResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	var state Group

	diags := req.State.Get(ctx, &state)
	resp.Diagnostics.Append(diags...)

	if resp.Diagnostics.HasError() {
		return
	}

	groupID := state.ID.Value
	group, statusCode, err := r.p.dog.DeleteGroup(groupID, nil)
	if statusCode != 204 {
		resp.Diagnostics.AddError("Client Unsuccesful", fmt.Sprintf("Status Code: %d", statusCode))
	}
	if err != nil {
		resp.Diagnostics.AddError("Client Error", fmt.Sprintf("Unable to read group, got error: %s", err))
	}
	if resp.Diagnostics.HasError() {
		return
	}
	tflog.Trace(ctx, fmt.Sprintf("group deleted: %+v\n", group))

	diags = resp.State.Set(ctx, &state)
	resp.Diagnostics.Append(diags...)
	resp.State.RemoveResource(ctx)
}
