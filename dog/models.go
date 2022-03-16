package dog

import "github.com/hashicorp/terraform-plugin-framework/types"

type HostList []Host

type Host struct {
	Active      types.String `tfsdk:"active"`
	Environment types.String `tfsdk:"environment"`
	Group       types.String `tfsdk:"group"`
	ID          types.String `tfsdk:"id"`
	HostKey     types.String `tfsdk:"hostkey"`
	Location    types.String `tfsdk:"location"`
	Name        types.String `tfsdk:"name"`
}

type GroupList []Group

type Group struct {
	//Created        int    `json:"created,omitempty"` //TODO: created has both int and string entries
	Description    types.String `tfsdk:"description"`
	ID             types.String `tfsdk:"id"`
	Name           types.String `tfsdk:"name"`
	ProfileName    types.String `tfsdk:"profile_name"`
	ProfileVersion types.String `tfsdk:"profile_version"`
}

type ServiceList []Service

type Service struct {
	//Created types.Int64  `tfsdk:"created"`
	ID       types.String   `tfsdk:"id"`
	Services []PortProtocol `tfsdk:"services"`
	Name     types.String   `tfsdk:"name"`
	Version  types.Int64    `tfsdk:"version"`
}

type Services []PortProtocol

type PortProtocol struct {
	Ports    []string     `tfsdk:"ports"`
	Protocol types.String `tfsdk:"protocol"`
}
