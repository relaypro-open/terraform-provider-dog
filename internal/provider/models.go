package provider

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
