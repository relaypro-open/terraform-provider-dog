package main

import (
	"bufio"
	"fmt"
	"log"
	"os"
	"strings"

	"github.com/relaypro-open/dog_api_golang/api"
	"gopkg.in/rethinkdb/rethinkdb-go.v6"
	r "gopkg.in/rethinkdb/rethinkdb-go.v6"
)

func dbSession() *rethinkdb.Session {
	//SetTags("rethinkdb", "json")
	session, err := r.Connect(r.ConnectOpts{
		Address:  "dog-ubuntu-server.lxd:28015",
		Database: "dog",
		Username: "admin",
		Password: "",
	})
	if err != nil {
		log.Fatalln(err)
	}
	return session
}

func check(err error) {
	if err != nil {
		panic(err)
	}
}

func outputFiles(table string) (*bufio.Writer, *bufio.Writer) {
	tf_f, err := os.Create(fmt.Sprintf("/tmp/%s.tf", table))
	check(err)
	tf_w := bufio.NewWriter(tf_f)

	import_f, err := os.Create(fmt.Sprintf("/tmp/%s_import.sh", table))
	check(err)
	import_w := bufio.NewWriter(import_f)

	return tf_w, import_w
}

func link_export(session *rethinkdb.Session) {
	fmt.Printf("link_export\n")
	tf_w, import_w := outputFiles("link")

	c := api.NewClient(os.Getenv("DOG_API_KEY"), os.Getenv("DOG_API_ENDPOINT"))

	res, statusCode, err := c.GetLinks(nil)
	if err != nil {
		log.Fatalln(err)
	}
	if statusCode != 200 {
		log.Fatalln(err)
	}

	for _, row := range res {
		terraformName := strings.ReplaceAll(row.Name, ".", "_")
		fmt.Fprintf(tf_w, "resource \"dog_link\" \"%s\" {\n", terraformName)
		fmt.Fprintf(tf_w, "  address_handling = \"%s\"\n", row.AddressHandling)
		fmt.Fprintf(tf_w, "  connection = \n")
		fmt.Fprintf(tf_w, "  {\n")
		fmt.Fprintf(tf_w, "    api_port = \"%d\"\n", row.Connection.ApiPort)
		fmt.Fprintf(tf_w, "    host = \"%s\"\n", row.Connection.Host)
		fmt.Fprintf(tf_w, "    password = \"%s\"\n", row.Connection.Password)
		fmt.Fprintf(tf_w, "    port = \"%d\"\n", row.Connection.Port)
		fmt.Fprintf(tf_w, "    ssl_options = \n")
		fmt.Fprintf(tf_w, "      {\n")
		fmt.Fprintf(tf_w, "        ca_cert_file = \"%s\"\n", row.Connection.SSLOptions.CaCertFile)
		fmt.Fprintf(tf_w, "        cert_file = \"%s\"\n", row.Connection.SSLOptions.CertFile)
		fmt.Fprintf(tf_w, "        fail_if_no_peer_cert = \"%t\"\n", row.Connection.SSLOptions.FailIfNoPeerCert)
		fmt.Fprintf(tf_w, "        keyfile = \"%s\"\n", row.Connection.SSLOptions.KeyFile)
		fmt.Fprintf(tf_w, "        server_name_indication = \"%s\"\n", row.Connection.SSLOptions.ServerNameIndication)
		fmt.Fprintf(tf_w, "        verify = \"%s\"\n", row.Connection.SSLOptions.Verify)
		fmt.Fprintf(tf_w, "      }\n")
		fmt.Fprintf(tf_w, "  },\n")
		fmt.Fprintf(tf_w, "  connection_type = \"%s\"\n", row.ConnectionType)
		fmt.Fprintf(tf_w, "  direction = \"%s\"\n", row.Direction)
		fmt.Fprintf(tf_w, "  enabled = \"%t\"\n", row.Enabled)
		fmt.Fprintf(tf_w, "  name = \"%s\"\n", row.Name)
		fmt.Fprintf(tf_w, "}\n")
		fmt.Fprintf(tf_w, "\n")

		fmt.Fprintf(import_w, "terraform import dog_link.%s %s\n", row.Name, row.ID)
	}
	tf_w.Flush()
	import_w.Flush()
}

func host_export(session *rethinkdb.Session) {
	fmt.Printf("host_export\n")
	tf_w, import_w := outputFiles("host")

	c := api.NewClient(os.Getenv("DOG_API_KEY"), os.Getenv("DOG_API_ENDPOINT"))

	res, statusCode, err := c.GetHosts(nil)
	if err != nil {
		log.Fatalln(err)
	}
	if statusCode != 200 {
		log.Fatalln(err)
	}

	for _, row := range res {
		terraformName := strings.ReplaceAll(row.Name, ".", "_")
		fmt.Fprintf(tf_w, "resource \"dog_host\" \"%s\" {\n", terraformName)
		fmt.Fprintf(tf_w, "  active = \"%s\"\n", row.Active)
		fmt.Fprintf(tf_w, "  environment = \"%s\"\n", row.Environment)
		fmt.Fprintf(tf_w, "  group = \"%s\"\n", row.Group)
		fmt.Fprintf(tf_w, "  hostkey = \"%s\"\n", row.HostKey)
		fmt.Fprintf(tf_w, "  location = \"%s\"\n", row.Location)
		fmt.Fprintf(tf_w, "  name = \"%s\"\n", row.Name)
		fmt.Fprintf(tf_w, "}\n")
		fmt.Fprintf(tf_w, "\n")

		fmt.Fprintf(import_w, "terraform import dog_host.%s %s\n", row.Name, row.ID)
	}
	tf_w.Flush()
	import_w.Flush()
}

func portprotocols_output(tf_w *bufio.Writer, portProtocols []api.PortProtocol) {
	for _, port_protocol := range portProtocols {
		fmt.Fprintf(tf_w, "      {\n")
		fmt.Fprintf(tf_w, "        protocol = \"%s\"\n", port_protocol.Protocol)
		fmt.Fprintf(tf_w, strings.ReplaceAll(fmt.Sprintf("        ports = %q\n", port_protocol.Ports), "\" \"", "\",\""))
		fmt.Fprintf(tf_w, "      }\n")
	}
}

func group_export(session *rethinkdb.Session) {
	fmt.Printf("group_export\n")
	tf_w, import_w := outputFiles("group")

	c := api.NewClient(os.Getenv("DOG_API_KEY"), os.Getenv("DOG_API_ENDPOINT"))

	res, statusCode, err := c.GetGroups(nil)
	if err != nil {
		log.Fatalln(err)
	}
	if statusCode != 200 {
		log.Fatalln(err)
	}

	for _, row := range res {
		terraformName := strings.ReplaceAll(row.Name, ".", "_")
		fmt.Fprintf(tf_w, "resource \"dog_group\" \"%s\" {\n", terraformName)
		fmt.Fprintf(tf_w, "  description = \"%s\"\n", row.Description)
		fmt.Fprintf(tf_w, "  name = \"%s\"\n", row.Name)
		fmt.Fprintf(tf_w, "  profile_name = \"%s\"\n", row.ProfileName)
		fmt.Fprintf(tf_w, "  profile_version = \"%s\"\n", row.ProfileVersion)
		fmt.Fprintf(tf_w, "}\n")
		fmt.Fprintf(tf_w, "\n")

		fmt.Fprintf(import_w, "terraform import dog_group.%s %s\n", row.Name, row.ID)
	}
	tf_w.Flush()
	import_w.Flush()
}

func service_export(session *rethinkdb.Session) {
	fmt.Printf("service_export\n")
	tf_w, import_w := outputFiles("service")

	c := api.NewClient(os.Getenv("DOG_API_KEY"), os.Getenv("DOG_API_ENDPOINT"))

	res, statusCode, err := c.GetServices(nil)
	if err != nil {
		log.Fatalln(err)
	}
	if statusCode != 200 {
		log.Fatalln(err)
	}

	for _, row := range res {
		terraformName := strings.ReplaceAll(row.Name, ".", "_")
		fmt.Fprintf(tf_w, "resource \"dog_service\" \"%s\" {\n", terraformName)
		fmt.Fprintf(tf_w, "  name = \"%s\"\n", row.Name)
		fmt.Fprintf(tf_w, "  version = \"%d\"\n", row.Version)
		portprotocols_output(tf_w, row.Services)
		fmt.Fprintf(tf_w, "}\n")
		fmt.Fprintf(tf_w, "\n")

		fmt.Fprintf(import_w, "terraform import dog_service.%s %s\n", row.Name, row.ID)
	}
	tf_w.Flush()
	import_w.Flush()
}

func zone_export(session *rethinkdb.Session) {
	fmt.Printf("zone_export\n")
	tf_w, import_w := outputFiles("zone")

	c := api.NewClient(os.Getenv("DOG_API_KEY"), os.Getenv("DOG_API_ENDPOINT"))

	res, statusCode, err := c.GetZones(nil)
	if err != nil {
		log.Fatalln(err)
	}
	if statusCode != 200 {
		log.Fatalln(err)
	}

	for _, row := range res {
		terraformName := strings.ReplaceAll(row.Name, ".", "_")
		fmt.Fprintf(tf_w, "resource \"dog_zone\" \"%s\" {\n", terraformName)
		fmt.Fprintf(tf_w, "  name = \"%s\"\n", row.Name)
		fmt.Fprintf(tf_w, strings.ReplaceAll(fmt.Sprintf("  ipv4_addresses = %q\n", row.IPv4Addresses), "\" \"", "\",\""))
		fmt.Fprintf(tf_w, strings.ReplaceAll(fmt.Sprintf("  ipv6_addresses = %q\n", row.IPv6Addresses), "\" \"", "\",\""))
		fmt.Fprintf(tf_w, "}\n")
		fmt.Fprintf(tf_w, "\n")

		fmt.Fprintf(import_w, "terraform import dog_zone.%s %s\n", row.Name, row.ID)
	}
	tf_w.Flush()
	import_w.Flush()
}

func rules_output(tf_w *bufio.Writer, rules []api.Rule) {
	for _, rule := range rules {
		fmt.Fprintf(tf_w, "      {\n")
		fmt.Fprintf(tf_w, "        action = \"%s\"\n", rule.Action)
		fmt.Fprintf(tf_w, "        active = \"%t\"\n", rule.Active)
		fmt.Fprintf(tf_w, "        comment = \"%s\"\n", rule.Comment)
		fmt.Fprintf(tf_w, strings.ReplaceAll(fmt.Sprintf("        environments = %q\n", rule.Environments), "\" \"", "\",\""))
		fmt.Fprintf(tf_w, "        group = \"%s\"\n", rule.Group)
		fmt.Fprintf(tf_w, "        group_type = \"%s\"\n", rule.GroupType)
		fmt.Fprintf(tf_w, "        interface = \"%s\"\n", rule.Interface)
		fmt.Fprintf(tf_w, "        log = \"%t\"\n", rule.Log)
		fmt.Fprintf(tf_w, "        log_prefix = \"%s\"\n", rule.LogPrefix)
		fmt.Fprintf(tf_w, "        order = \"%d\"\n", rule.Order)
		fmt.Fprintf(tf_w, "        service = \"%s\"\n", rule.Service)
		fmt.Fprintf(tf_w, strings.ReplaceAll(fmt.Sprintf("        states = %q\n", rule.States), "\" \"", "\",\""))
		fmt.Fprintf(tf_w, "        type = \"%s\"\n", rule.Type)

		fmt.Fprintf(tf_w, "      }\n")
	}
}

func profile_export(session *rethinkdb.Session) {
	fmt.Printf("profile_export\n")
	tf_w, import_w := outputFiles("profile")

	c := api.NewClient(os.Getenv("DOG_API_KEY"), os.Getenv("DOG_API_ENDPOINT"))

	res, statusCode, err := c.GetProfiles(nil)
	if err != nil {
		log.Fatalln(err)
	}
	if statusCode != 200 {
		log.Fatalln(err)
	}

	for _, row := range res {
		terraformName := strings.ReplaceAll(row.Name, ".", "_")
		fmt.Fprintf(tf_w, "resource \"dog_profile\" \"%s\" {\n", terraformName)
		fmt.Fprintf(tf_w, "  name = \"%s\"\n", row.Name)
		fmt.Fprintf(tf_w, "  version = \"%s\"\n", row.Version)
		fmt.Fprintf(tf_w, "  rules = {\n")
		inbound := row.Rules.Inbound
		fmt.Fprintf(tf_w, "    inbound = [\n")
		rules_output(tf_w, inbound)
		fmt.Fprintf(tf_w, "    ]\n")
		fmt.Fprintf(tf_w, "    outbound = [\n")
		outbound := row.Rules.Outbound
		rules_output(tf_w, outbound)
		fmt.Fprintf(tf_w, "    ]\n")
		fmt.Fprintf(tf_w, "}\n")
		fmt.Fprintf(tf_w, "\n")

		fmt.Fprintf(import_w, "terraform import dog_profile.%s %s\n", row.Name, row.ID)
	}
	tf_w.Flush()
	import_w.Flush()
}

func main() {
	session := dbSession()
	group_export(session)
	host_export(session)
	link_export(session)
	profile_export(session)
	service_export(session)
	zone_export(session)
}
