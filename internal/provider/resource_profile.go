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

type profileResourceData struct {
	ID      types.String          `tfsdk:"id"`
	Name    string                `tfsdk:"name"`
	Version string                `tfsdk:"version"`
}

type (
	profileResource struct {
		p dogProvider
	}
)

var (
	_ resource.Resource                = (*profileResource)(nil)
	_ resource.ResourceWithImportState = (*profileResource)(nil)
)

func NewProfileResource() resource.Resource {
	return &profileResource{}
}

func (*profileResource) Metadata(ctx context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_profile"
}


func (*profileResource) GetSchema(ctx context.Context) (tfsdk.Schema, diag.Diagnostics) {
	return tfsdk.Schema{
		Attributes: map[string]tfsdk.Attribute{
			"name": {
				MarkdownDescription: "Profile name",
				Required:            true,
				Type:                types.StringType,
			},
			"version": {
				MarkdownDescription: "Profile version",
				Required:            true,
				Type:                types.StringType,
			},
			"id": {
				Computed:            true,
				MarkdownDescription: "Profile identifier",
				PlanModifiers: tfsdk.AttributePlanModifiers{
					resource.UseStateForUnknown(),
				},
				Type: types.StringType,
			},
		},
	}, nil
}

func (r *profileResource) Configure(ctx context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {
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

func (*profileResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	resource.ImportStatePassthroughID(ctx, path.Root("id"), req, resp)
}



func ProfileToCreateRequest(plan profileResourceData) api.ProfileCreateRequest {

	newProfile := api.ProfileCreateRequest{
		Name:    plan.Name,
		Version: plan.Version,
	}
	return newProfile
}

func ProfileToUpdateRequest(plan profileResourceData) api.ProfileUpdateRequest {
	newProfile := api.ProfileUpdateRequest{
		Name: plan.Name,
		Version: plan.Version,
	}
	return newProfile
}

func ApiToProfile(profile api.Profile) Profile {
	h := Profile{
		//Created:     types.Int64{Value: int64(profile.Created)},
		ID:   types.String{Value: profile.ID},
		Name: types.String{Value: profile.Name},
		Version: types.String{Value: profile.Version},
	}
	return h
}

func (r *profileResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	var state Profile

	var plan profileResourceData
	diags := req.Plan.Get(ctx, &plan)
	resp.Diagnostics.Append(diags...)
	//	resp.Diagnostics.AddError("Client Error", fmt.Sprintf("client: %+v\n", r.provider.client))
	if resp.Diagnostics.HasError() {
		return
	}

	newProfile := ProfileToCreateRequest(plan)
	log.Printf(fmt.Sprintf("r.p.dog: %+v\n", r.p.dog))
	profile, statusCode, err := r.p.dog.CreateProfile(newProfile, nil)
	log.Printf(fmt.Sprintf("profile: %+v\n", profile))
	tflog.Trace(ctx, fmt.Sprintf("profile: %+v\n", profile))
	if statusCode != 201 {
		resp.Diagnostics.AddError("Client Unsuccesful", fmt.Sprintf("Status Code: %d", statusCode))
	}
	if err != nil {
		resp.Diagnostics.AddError("Client Error", fmt.Sprintf("Unable to create profile, got error: %s", err))
	}
	if resp.Diagnostics.HasError() {
		return
	}
	state = ApiToProfile(profile)

	plan.ID = state.ID

	// write logs using the tflog package
	// see https://pkg.go.dev/github.com/hashicorp/terraform-plugin-log/tflog
	// for more information
	tflog.Trace(ctx, "created a resource")

	diags = resp.State.Set(ctx, &state)
	resp.Diagnostics.Append(diags...)
}

func (r *profileResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	var state Profile

	diags := req.State.Get(ctx, &state)
	resp.Diagnostics.Append(diags...)

	if resp.Diagnostics.HasError() {
		return
	}

	profileID := state.ID.Value

	log.Printf(fmt.Sprintf("r.p: %+v\n", r.p))
	log.Printf(fmt.Sprintf("r.p.dog: %+v\n", r.p.dog))
	profile, statusCode, err := r.p.dog.GetProfile(profileID, nil)
	if statusCode != 200 {
		resp.Diagnostics.AddError("Client Unsuccesful", fmt.Sprintf("Status Code: %d", statusCode))
	}
	if err != nil {
		resp.Diagnostics.AddError("Client Error", fmt.Sprintf("Unable to read profile, got error: %s", err))
	}
	if resp.Diagnostics.HasError() {
		return
	}
	state = ApiToProfile(profile)
	diags = resp.State.Set(ctx, &state)
	resp.Diagnostics.Append(diags...)
}


func (r *profileResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	var state Profile

	diags := req.State.Get(ctx, &state)
	resp.Diagnostics.Append(diags...)

	if resp.Diagnostics.HasError() {
		return
	}

	profileID := state.ID.Value

	var plan profileResourceData
	diags = req.Plan.Get(ctx, &plan)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	newProfile := ProfileToUpdateRequest(plan)
	profile, statusCode, err := r.p.dog.UpdateProfile(profileID, newProfile, nil)
	log.Printf(fmt.Sprintf("profile: %+v\n", profile))
	tflog.Trace(ctx, fmt.Sprintf("profile: %+v\n", profile))
	state = ApiToProfile(profile)
	if err != nil {
		resp.Diagnostics.AddError("Client Error", fmt.Sprintf("Unable to create profile, got error: %s", err))
	}
	ok := []int{303, 200, 201}
	if slices.Contains(ok, statusCode) != true {
		resp.Diagnostics.AddError("Client Unsuccesful", fmt.Sprintf("Status Code: %d", statusCode))
	}
	if resp.Diagnostics.HasError() {
		return
	}

	plan.ID = state.ID

	tflog.Trace(ctx, "created a resource")

	diags = resp.State.Set(ctx, &state)
	resp.Diagnostics.Append(diags...)

}

func (r *profileResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	var state Profile

	diags := req.State.Get(ctx, &state)
	resp.Diagnostics.Append(diags...)

	if resp.Diagnostics.HasError() {
		return
	}

	profileID := state.ID.Value
	profile, statusCode, err := r.p.dog.DeleteProfile(profileID, nil)
	if statusCode != 204 {
		resp.Diagnostics.AddError("Client Unsuccesful", fmt.Sprintf("Status Code: %d", statusCode))
	}
	if err != nil {
		resp.Diagnostics.AddError("Client Error", fmt.Sprintf("Unable to read profile, got error: %s", err))
	}
	if resp.Diagnostics.HasError() {
		return
	}
	tflog.Trace(ctx, fmt.Sprintf("profile deleted: %+v\n", profile))

	diags = resp.State.Set(ctx, &state)
	resp.Diagnostics.Append(diags...)
	resp.State.RemoveResource(ctx)
}
