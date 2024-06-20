package dog

import (
	"context"
	"fmt"
	"log"
	"regexp"

	"github.com/hashicorp/terraform-plugin-framework-validators/stringvalidator"
	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/schema/validator"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-log/tflog"
	api "github.com/relaypro-open/dog_api_golang/api"
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

func (*groupResource) Schema(ctx context.Context, req resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = schema.Schema{
		// This description is used by the documentation generator and the language server.
		Attributes: map[string]schema.Attribute{
			// This description is used by the documentation generator and the language server.
			"description": schema.StringAttribute{
				MarkdownDescription: "group description",
				Optional:            true,
			},
			"name": schema.StringAttribute{
				MarkdownDescription: "group name",
				Optional:            true,
				Validators: []validator.String{
					stringvalidator.LengthBetween(1, 28),
					stringvalidator.RegexMatches(
						regexp.MustCompile(`^[A-Za-z0-9_.-](.*)$`),
						"must start with 'sg-'",
					),
				},
			},
			"profile_id": schema.StringAttribute{
				MarkdownDescription: "group profile id",
				Optional:            true,
			},
			"profile_name": schema.StringAttribute{
				MarkdownDescription: "group profile name",
				Optional:            true,
			},
			"profile_version": schema.StringAttribute{
				MarkdownDescription: "group profile version",
				Optional:            true,
			},
			"ec2_security_group_ids": schema.ListNestedAttribute{
				MarkdownDescription: "List of EC2 Security Groups to control",
				Optional:            true,
				NestedObject: schema.NestedAttributeObject{
					Attributes: map[string]schema.Attribute{
						"region": schema.StringAttribute{
							MarkdownDescription: "EC2 Region",
							Required:            true,
							Validators: []validator.String{
								stringvalidator.LengthBetween(9, 256),
								stringvalidator.RegexMatches(
									regexp.MustCompile(`^(.*)-(.*)-(.*)$`),
									"must be valid region",
								),
							},
						},
						"sgid": schema.StringAttribute{
							MarkdownDescription: "EC2 Security Group ID",
							Required:            true,
							Validators: []validator.String{
								stringvalidator.LengthBetween(3, 256),
								stringvalidator.RegexMatches(
									regexp.MustCompile(`^sg-(.*)$`),
									"must start with 'sg-'",
								),
							},
						},
					},
				},
			},
			"vars": schema.StringAttribute{
				MarkdownDescription: "json string of vars",
				Optional:            true,
			},
			"alert_enable": schema.BoolAttribute{
				MarkdownDescription: "alert enable",
				Optional:            true,
			},
			"id": schema.StringAttribute{
				MarkdownDescription: "group identifier",
				Optional:            true,
				Computed:            true,
			},
		},
		Version: 1,
	}
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
	Description         string                             `tfsdk:"description"`
	ID                  types.String                       `tfsdk:"id"`
	Name                string                             `tfsdk:"name"`
	ProfileId           string                             `tfsdk:"profile_id"`
	ProfileName         string                             `tfsdk:"profile_name"`
	ProfileVersion      string                             `tfsdk:"profile_version"`
	Ec2SecurityGroupIds []*ec2SecurityGroupIdsResourceData `tfsdk:"ec2_security_group_ids"`
	Vars                *string                            `tfsdk:"vars"`
	AlertEnable         *bool                              `tfsdk:"alert_enable"`
}

type ec2SecurityGroupIdsResourceData struct {
	Region string `tfsdk:"region"`
	SgId   string `tfsdk:"sgid"`
}

func GroupToApiGroup(plan Group) api.Group {
	newEc2SecurityGroupIds := []*api.Ec2SecurityGroupIds{}
	for _, region_sgid := range plan.Ec2SecurityGroupIds {
		rs := &api.Ec2SecurityGroupIds{
			Region: region_sgid.Region.ValueString(),
			SgId:   region_sgid.SgId.ValueString(),
		}
		newEc2SecurityGroupIds = append(newEc2SecurityGroupIds, rs)
	}

	if plan.Vars.ValueString() != "" {
		if plan.AlertEnable.IsNull() {
			newGroup := api.Group{
				Description:         plan.Description.ValueString(),
				Name:                plan.Name.ValueString(),
				ProfileId:           plan.ProfileId.ValueString(),
				ProfileName:         plan.ProfileName.ValueString(),
				ProfileVersion:      plan.ProfileVersion.ValueString(),
				Ec2SecurityGroupIds: newEc2SecurityGroupIds,
				Vars:                plan.Vars.ValueString(),
			}
			return newGroup
		} else {
			newGroup := api.Group{
				Description:         plan.Description.ValueString(),
				Name:                plan.Name.ValueString(),
				ProfileId:           plan.ProfileId.ValueString(),
				ProfileName:         plan.ProfileName.ValueString(),
				ProfileVersion:      plan.ProfileVersion.ValueString(),
				Ec2SecurityGroupIds: newEc2SecurityGroupIds,
				Vars:                plan.Vars.ValueString(),
				AlertEnable:         plan.AlertEnable.ValueBoolPointer(),
			}
			return newGroup
		}
	} else {
		if plan.AlertEnable.IsNull() {
			newGroup := api.Group{
				Description:         plan.Description.ValueString(),
				Name:                plan.Name.ValueString(),
				ProfileId:           plan.ProfileId.ValueString(),
				ProfileName:         plan.ProfileName.ValueString(),
				ProfileVersion:      plan.ProfileVersion.ValueString(),
				Ec2SecurityGroupIds: newEc2SecurityGroupIds,
			}
			return newGroup
		} else {
			newGroup := api.Group{
				Description:         plan.Description.ValueString(),
				Name:                plan.Name.ValueString(),
				ProfileId:           plan.ProfileId.ValueString(),
				ProfileName:         plan.ProfileName.ValueString(),
				ProfileVersion:      plan.ProfileVersion.ValueString(),
				Ec2SecurityGroupIds: newEc2SecurityGroupIds,
				AlertEnable:         plan.AlertEnable.ValueBoolPointer(),
			}
			return newGroup
		}
	}
}

func ApiToGroup(group api.Group) Group {
	newEc2SecurityGroupIds := []*Ec2SecurityGroupIds{}
	for _, region_sgid := range group.Ec2SecurityGroupIds {
		rs := &Ec2SecurityGroupIds{
			Region: types.StringValue(region_sgid.Region),
			SgId:   types.StringValue(region_sgid.SgId),
		}
		newEc2SecurityGroupIds = append(newEc2SecurityGroupIds, rs)
	}
	if group.Vars != "" {
		if group.AlertEnable == nil {
			h := Group{
				Description:         types.StringValue(group.Description),
				ID:                  types.StringValue(group.ID),
				Name:                types.StringValue(group.Name),
				ProfileId:           types.StringValue(group.ProfileId),
				ProfileName:         types.StringValue(group.ProfileName),
				ProfileVersion:      types.StringValue(group.ProfileVersion),
				Ec2SecurityGroupIds: newEc2SecurityGroupIds,
				Vars:                types.StringValue(group.Vars),
			}
			return h
		} else {
			h := Group{
				Description:         types.StringValue(group.Description),
				ID:                  types.StringValue(group.ID),
				Name:                types.StringValue(group.Name),
				ProfileId:           types.StringValue(group.ProfileId),
				ProfileName:         types.StringValue(group.ProfileName),
				ProfileVersion:      types.StringValue(group.ProfileVersion),
				Ec2SecurityGroupIds: newEc2SecurityGroupIds,
				Vars:                types.StringValue(group.Vars),
				AlertEnable:         types.BoolValue(*group.AlertEnable),
			}
			return h

		}
	} else {
		if group.AlertEnable == nil {
			h := Group{
				Description:         types.StringValue(group.Description),
				ID:                  types.StringValue(group.ID),
				Name:                types.StringValue(group.Name),
				ProfileId:           types.StringValue(group.ProfileId),
				ProfileName:         types.StringValue(group.ProfileName),
				ProfileVersion:      types.StringValue(group.ProfileVersion),
				Ec2SecurityGroupIds: newEc2SecurityGroupIds,
			}
			return h
		} else {
			h := Group{
				Description:         types.StringValue(group.Description),
				ID:                  types.StringValue(group.ID),
				Name:                types.StringValue(group.Name),
				ProfileId:           types.StringValue(group.ProfileId),
				ProfileName:         types.StringValue(group.ProfileName),
				ProfileVersion:      types.StringValue(group.ProfileVersion),
				Ec2SecurityGroupIds: newEc2SecurityGroupIds,
				AlertEnable:         types.BoolValue(*group.AlertEnable),
			}
			return h
		}
	}
}

func GroupToCreateRequest(plan groupResourceData) api.GroupCreateRequest {
	newEc2SecurityGroupIds := []*api.Ec2SecurityGroupIds{}
	for _, region_sgid := range plan.Ec2SecurityGroupIds {
		rs := &api.Ec2SecurityGroupIds{
			Region: region_sgid.Region,
			SgId:   region_sgid.SgId,
		}
		newEc2SecurityGroupIds = append(newEc2SecurityGroupIds, rs)
	}

	newGroup := api.GroupCreateRequest{
		Description:         plan.Description,
		Name:                plan.Name,
		ProfileId:           plan.ProfileId,
		ProfileName:         plan.ProfileName,
		ProfileVersion:      plan.ProfileVersion,
		Ec2SecurityGroupIds: newEc2SecurityGroupIds,
		Vars:                *plan.Vars,
		AlertEnable:         plan.AlertEnable,
	}
	return newGroup
}

func GroupToUpdateRequest(plan groupResourceData) api.GroupUpdateRequest {
	newEc2SecurityGroupIds := []*api.Ec2SecurityGroupIds{}
	for _, region_sgid := range plan.Ec2SecurityGroupIds {
		rs := &api.Ec2SecurityGroupIds{
			Region: region_sgid.Region,
			SgId:   region_sgid.SgId,
		}
		newEc2SecurityGroupIds = append(newEc2SecurityGroupIds, rs)
	}

	newGroup := api.GroupUpdateRequest{
		Description:         plan.Description,
		Name:                plan.Name,
		ProfileId:           plan.ProfileId,
		ProfileName:         plan.ProfileName,
		ProfileVersion:      plan.ProfileVersion,
		Ec2SecurityGroupIds: newEc2SecurityGroupIds,
		Vars:                *plan.Vars,
		AlertEnable:         plan.AlertEnable,
	}
	return newGroup
}

func (r *groupResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	var state Group

	var plan Group
	diags := req.Plan.Get(ctx, &plan)
	resp.Diagnostics.Append(diags...)
	//	resp.Diagnostics.AddError("Client Error", fmt.Sprintf("client: %+v\n", r.provider.client))
	if resp.Diagnostics.HasError() {
		return
	}

	tflog.Debug(ctx, PrettyFmt("group create plan", plan))
	newGroup := GroupToApiGroup(plan)
	tflog.Debug(ctx, PrettyFmt("group create newGroup", newGroup))
	group, statusCode, err := r.p.dog.CreateGroupEncode(newGroup, nil)
	tflog.Debug(ctx, PrettyFmt("group create group", group))
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
	tflog.Debug(ctx, PrettyFmt("group create state", state))

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

	groupID := state.ID.ValueString()

	group, statusCode, err := r.p.dog.GetGroupEncode(groupID, nil)
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

	groupID := state.ID.ValueString()

	var plan Group
	diags = req.Plan.Get(ctx, &plan)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	newGroup := GroupToApiGroup(plan)
	group, statusCode, err := r.p.dog.UpdateGroupEncode(groupID, newGroup, nil)
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

	groupID := state.ID.ValueString()
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

type groupResourceModelV0 struct {
	Description         string                             `tfsdk:"description"`
	ID                  types.String                       `tfsdk:"id"`
	Name                string                             `tfsdk:"name"`
	ProfileId           string                             `tfsdk:"profile_id"`
	ProfileName         string                             `tfsdk:"profile_name"`
	ProfileVersion      string                             `tfsdk:"profile_version"`
	Ec2SecurityGroupIds []*ec2SecurityGroupIdsResourceData `tfsdk:"ec2_security_group_ids"`
	Vars                *string                            `tfsdk:"vars"`
}

type groupResourceModelV1 struct {
	Description         string                             `tfsdk:"description"`
	ID                  types.String                       `tfsdk:"id"`
	Name                string                             `tfsdk:"name"`
	ProfileId           string                             `tfsdk:"profile_id"`
	ProfileName         string                             `tfsdk:"profile_name"`
	ProfileVersion      string                             `tfsdk:"profile_version"`
	Ec2SecurityGroupIds []*ec2SecurityGroupIdsResourceData `tfsdk:"ec2_security_group_ids"`
	Vars                *string                            `tfsdk:"vars"`
	AlertEnable         *bool                              `tfsdk:"alert_enable"`
}

func (r *groupResource) UpgradeState(ctx context.Context) map[int64]resource.StateUpgrader {
	tflog.Debug(ctx, "UpgradeState")
	return map[int64]resource.StateUpgrader{
		// State upgrade implementation from 0 (prior state version) to 1 (Schema.Version)
		0: {
			PriorSchema: &schema.Schema{
				// This description is used by the documentation generator and the language server.
				Attributes: map[string]schema.Attribute{
					// This description is used by the documentation generator and the language server.
					"description": schema.StringAttribute{
						MarkdownDescription: "group description",
						Optional:            true,
					},
					"name": schema.StringAttribute{
						MarkdownDescription: "group name",
						Optional:            true,
						Validators: []validator.String{
							stringvalidator.LengthBetween(1, 28),
							stringvalidator.RegexMatches(
								regexp.MustCompile(`^[A-Za-z0-9_.-](.*)$`),
								"must start with 'sg-'",
							),
						},
					},
					"profile_id": schema.StringAttribute{
						MarkdownDescription: "group profile id",
						Optional:            true,
					},
					"profile_name": schema.StringAttribute{
						MarkdownDescription: "group profile name",
						Optional:            true,
					},
					"profile_version": schema.StringAttribute{
						MarkdownDescription: "group profile version",
						Optional:            true,
					},
					"ec2_security_group_ids": schema.ListNestedAttribute{
						MarkdownDescription: "List of EC2 Security Groups to control",
						Optional:            true,
						NestedObject: schema.NestedAttributeObject{
							Attributes: map[string]schema.Attribute{
								"region": schema.StringAttribute{
									MarkdownDescription: "EC2 Region",
									Required:            true,
									Validators: []validator.String{
										stringvalidator.LengthBetween(9, 256),
										stringvalidator.RegexMatches(
											regexp.MustCompile(`^(.*)-(.*)-(.*)$`),
											"must be valid region",
										),
									},
								},
								"sgid": schema.StringAttribute{
									MarkdownDescription: "EC2 Security Group ID",
									Required:            true,
									Validators: []validator.String{
										stringvalidator.LengthBetween(3, 256),
										stringvalidator.RegexMatches(
											regexp.MustCompile(`^sg-(.*)$`),
											"must start with 'sg-'",
										),
									},
								},
							},
						},
					},
					"vars": schema.StringAttribute{
						MarkdownDescription: "json string of vars",
						Optional:            true,
					},
					"id": schema.StringAttribute{
						MarkdownDescription: "group identifier",
						Optional:            true,
						Computed:            true,
					},
				},
			},
			StateUpgrader: func(ctx context.Context, req resource.UpgradeStateRequest, resp *resource.UpgradeStateResponse) {
				var priorStateData groupResourceModelV0

				//resp.Diagnostics.Append(
				req.State.Get(ctx, &priorStateData)
				//)

				if resp.Diagnostics.HasError() {
					return
				}

				var alertEnable *bool
				req.State.GetAttribute(ctx, path.Root("alertEnable"), alertEnable)

				upgradedStateData := groupResourceModelV1{
					Description:         priorStateData.Description,
					Name:                priorStateData.Name,
					ProfileId:           priorStateData.ProfileId,
					ProfileName:         priorStateData.ProfileName,
					ProfileVersion:      priorStateData.ProfileVersion,
					Ec2SecurityGroupIds: priorStateData.Ec2SecurityGroupIds,
					Vars:                priorStateData.Vars,
					AlertEnable:         alertEnable,
					ID:                  priorStateData.ID,
				}

				resp.State.Set(ctx, upgradedStateData)
			},
		},
	}
}
