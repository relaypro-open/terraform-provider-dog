package dog

import "github.com/hashicorp/terraform-plugin-framework/types"

type Host struct {
	Active      types.String `tfsdk:"active"`
	Environment types.String `tfsdk:"environment"`
	Group       types.String `tfsdk:"group"`
	ID          types.String `tfsdk:"id"`
	HostKey     types.String `tfsdk:"hostkey"`
	Location    types.String `tfsdk:"location"`
	Name        types.String `tfsdk:"name"`
}

type HostList []Host

type Group struct {
	//Created        int    `json:"created,omitempty"` //TODO: created has both int and string entries
	Description    types.String `tfsdk:"description"`
	ID             types.String `tfsdk:"id"`
	Name           types.String `tfsdk:"name"`
	ProfileName    types.String `tfsdk:"profile_name"`
	ProfileVersion types.String `tfsdk:"profile_version"`
}

type GroupList []Group
