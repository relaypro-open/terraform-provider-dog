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
	//Created        int    `tfsdk:"created,omitempty"` //TODO: created has both int and string entries
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

type ZoneList []Zone

type Zone struct {
	//Created       int          `tfsdk:"created"`
	ID            types.String `tfsdk:"id"`
	IPv4Addresses []string     `tfsdk:"ipv4_addresses"`
	IPv6Addresses []string     `tfsdk:"ipv6_addresses"`
	Name          types.String `tfsdk:"name"`
}

type LinkList []Link

type Link struct {
	ID              types.String `tfsdk:"id"`
	AddressHandling types.String `tfsdk:"address_handling"`
	Connection      Connection   `tfsdk:"connection"`
	ConnectionType  types.String `tfsdk:"connection_type"`
	Direction       types.String `tfsdk:"direction"`
	Enabled         types.Bool   `tfsdk:"enabled"`
	Name            types.String `tfsdk:"name"`
}

type Connection struct {
	ApiPort     types.Int64  `tfsdk:"api_port"`
	Host        types.String `tfsdk:"host"`
	Password    types.String `tfsdk:"password"`
	Port        types.Int64  `tfsdk:"port"`
	SSLOptions  SSLOptions   `tfsdk:"ssl_options"`
	User        types.String `tfsdk:"user"`
	VirtualHost types.String `tfsdk:"virtual_host"`
}

type SSLOptions struct {
	CaCertFile           types.String `tfsdk:"cacertfile"`
	CertFile             types.String `tfsdk:"certfile"`
	FailIfNoPeerCert     types.Bool   `tfsdk:"fail_if_no_peer_cert"`
	KeyFile              types.String `tfsdk:"keyfile"`
	ServerNameIndication types.String `tfsdk:"server_name_indication"`
	Verify               types.String `tfsdk:"verify"`
}
