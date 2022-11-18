package main

import (
	"bufio"
	"fmt"
	"log"
	"os"
	"strings"

	"github.com/relaypro-open/dog_api_golang/api"
)

func check(err error) {
	if err != nil {
		panic(err)
	}
}

func toTerraformName(name string) string {
	no_dots := strings.ReplaceAll(name, ".", "_")
	no_open_parenthesis := strings.ReplaceAll(no_dots, "(", "_")
	no_close_parenthesis := strings.ReplaceAll(no_open_parenthesis, ")", "_")
	no_forward_slash := strings.ReplaceAll(no_close_parenthesis, "/", "_")
	no_spaces := strings.ReplaceAll(no_forward_slash, " ", "_")
	no_colons := strings.ReplaceAll(no_spaces, ":", "_")
	return no_colons
}

func createDir(output_dir string, table string) {
	if err := os.MkdirAll(fmt.Sprintf("%s/%s", output_dir, table), os.ModePerm); err != nil {
		log.Fatal(err)
	}
}

func terraformOutputFile(output_dir string, table string) *bufio.Writer {
	tf_f, err := os.Create(fmt.Sprintf("%s/dog_%s.tf", output_dir, table))
	check(err)
	tf_w := bufio.NewWriter(tf_f)

	return tf_w
}

func importOutputFile(output_dir string, table string) *bufio.Writer {
	import_f, err := os.Create(fmt.Sprintf("%s/%s_import.sh", output_dir, table))
	check(err)
	import_w := bufio.NewWriter(import_f)
	fmt.Fprintf(import_w, "#!/bin/bash\n")

	return import_w
}

func link_export(output_dir string, environment string) {
	fmt.Printf("link_export\n")
	table := "link"
	import_w := importOutputFile(output_dir, table)

	c := api.NewClient(os.Getenv("DOG_API_KEY"), os.Getenv("DOG_API_ENDPOINT"))

	res, statusCode, err := c.GetLinks(nil)
	if err != nil {
		log.Fatalln(err)
	}
	if statusCode != 200 {
		log.Fatalln(err)
	}

	tf_w := terraformOutputFile(output_dir, table)
	for _, row := range res {
		terraformName := toTerraformName(row.Name)
		fmt.Fprintf(tf_w, "resource \"dog_link\" \"%s\" {\n", terraformName)
		fmt.Fprintf(tf_w, "  address_handling = \"%s\"\n", row.AddressHandling)
		fmt.Fprintf(tf_w, "  connection = {\n")
		fmt.Fprintf(tf_w, "    api_port = %d\n", row.Connection.ApiPort)
		fmt.Fprintf(tf_w, "    host = \"%s\"\n", row.Connection.Host)
		fmt.Fprintf(tf_w, "    password = \"%s\"\n", row.Connection.Password)
		fmt.Fprintf(tf_w, "    port = %d\n", row.Connection.Port)
		fmt.Fprintf(tf_w, "    ssl_options = {\n")
		fmt.Fprintf(tf_w, "        cacertfile = \"%s\"\n", row.Connection.SSLOptions.CaCertFile)
		fmt.Fprintf(tf_w, "        certfile = \"%s\"\n", row.Connection.SSLOptions.CertFile)
		fmt.Fprintf(tf_w, "        fail_if_no_peer_cert = %t\n", row.Connection.SSLOptions.FailIfNoPeerCert)
		fmt.Fprintf(tf_w, "        keyfile = \"%s\"\n", row.Connection.SSLOptions.KeyFile)
		fmt.Fprintf(tf_w, "        server_name_indication = \"%s\"\n", row.Connection.SSLOptions.ServerNameIndication)
		fmt.Fprintf(tf_w, "        verify = \"%s\"\n", row.Connection.SSLOptions.Verify)
		fmt.Fprintf(tf_w, "      },\n")
		fmt.Fprintf(tf_w, "    user = \"%s\"\n", row.Connection.User)
		fmt.Fprintf(tf_w, "    virtual_host = \"%s\"\n", row.Connection.VirtualHost)
		fmt.Fprintf(tf_w, "  }\n")
		fmt.Fprintf(tf_w, "  connection_type = \"%s\"\n", row.ConnectionType)
		fmt.Fprintf(tf_w, "  direction = \"%s\"\n", row.Direction)
		fmt.Fprintf(tf_w, "  enabled = %t\n", row.Enabled)
		fmt.Fprintf(tf_w, "  name = \"%s\"\n", row.Name)
		fmt.Fprintf(tf_w, "  provider = dog.%s\n", environment)
		fmt.Fprintf(tf_w, "}\n")
		fmt.Fprintf(tf_w, "\n")
		fmt.Fprintf(import_w, "terraform import dog_link.%s %s\n", terraformName, row.ID)

	}
	tf_w.Flush()
	import_w.Flush()
}

func host_export(output_dir string, environment string) {
	fmt.Printf("host_export\n")
	table := "host"
	import_w := importOutputFile(output_dir, table)

	c := api.NewClient(os.Getenv("DOG_API_KEY"), os.Getenv("DOG_API_ENDPOINT"))

	res, statusCode, err := c.GetHosts(nil)
	if err != nil {
		log.Fatalln(err)
	}
	if statusCode != 200 {
		log.Fatalln(err)
	}

	tf_w := terraformOutputFile(output_dir, table)
	for _, row := range res {
		HostKeyWithoutPlus := strings.ReplaceAll(row.HostKey, "+", "-")
		terraformName := toTerraformName(row.Name) + "_" + HostKeyWithoutPlus
		fmt.Fprintf(tf_w, "resource \"dog_host\" \"%s\" {\n", terraformName)
		fmt.Fprintf(tf_w, "  environment = \"%s\"\n", row.Environment)
		fmt.Fprintf(tf_w, "  group = \"%s\"\n", row.Group)
		fmt.Fprintf(tf_w, "  hostkey = \"%s\"\n", row.HostKey)
		fmt.Fprintf(tf_w, "  location = \"%s\"\n", row.Location)
		fmt.Fprintf(tf_w, "  name = \"%s\"\n", row.Name)
		fmt.Fprintf(tf_w, "  provider = dog.%s\n", environment)
		fmt.Fprintf(tf_w, "}\n")
		fmt.Fprintf(tf_w, "\n")

		fmt.Fprintf(import_w, "terraform import dog_host.%s %s\n", terraformName, row.ID)
	}
	tf_w.Flush()
	import_w.Flush()
}

func group_export(output_dir string, environment string) {
	fmt.Printf("group_export\n")
	table := "group"
	import_w := importOutputFile(output_dir, table)

	c := api.NewClient(os.Getenv("DOG_API_KEY"), os.Getenv("DOG_API_ENDPOINT"))

	res, statusCode, err := c.GetGroups(nil)
	if err != nil {
		log.Fatalln(err)
	}
	if statusCode != 200 {
		log.Fatalln(err)
	}

	tf_w := terraformOutputFile(output_dir, table)
	for _, row := range res {
		if row.ID != "all-active" {
			terraformName := toTerraformName(row.Name)
			fmt.Fprintf(tf_w, "resource \"dog_group\" \"%s\" {\n", terraformName)
			fmt.Fprintf(tf_w, "  description = \"%s\"\n", row.Description)
			fmt.Fprintf(tf_w, "  name = \"%s\"\n", row.Name)
			fmt.Fprintf(tf_w, "  profile_name = \"%s\"\n", row.ProfileName)
			fmt.Fprintf(tf_w, "  profile_version = \"%s\"\n", row.ProfileVersion)
			fmt.Fprintf(tf_w, "  ec2_security_group_ids = [\n")
			regionsgid_output(tf_w, row.Ec2SecurityGroupIds)
			fmt.Fprintf(tf_w, "  ]\n")
			fmt.Fprintf(tf_w, "  provider = dog.%s\n", environment)
			fmt.Fprintf(tf_w, "}\n")
			fmt.Fprintf(tf_w, "\n")
			fmt.Fprintf(import_w, "terraform import dog_group.%s %s\n", terraformName, row.ID)
		}
	}
	tf_w.Flush()
	import_w.Flush()
}

func regionsgid_output(tf_w *bufio.Writer, ec2SecurityGroupIds []*api.Ec2SecurityGroupIds) {
	for _, region_sgid := range ec2SecurityGroupIds {
		fmt.Fprintf(tf_w, "      {\n")
		fmt.Fprintf(tf_w, "        region = \"%s\"\n", region_sgid.Region)
		fmt.Fprintf(tf_w, "        sgid = \"%s\"\n", region_sgid.SgId)
		fmt.Fprintf(tf_w, "      },\n")
	}
}

func service_export(output_dir string, environment string) {
	fmt.Printf("service_export\n")
	table := "service"
	import_w := importOutputFile(output_dir, table)

	c := api.NewClient(os.Getenv("DOG_API_KEY"), os.Getenv("DOG_API_ENDPOINT"))

	res, statusCode, err := c.GetServices(nil)
	if err != nil {
		log.Fatalln(err)
	}
	if statusCode != 200 {
		log.Fatalln(err)
	}

	tf_w := terraformOutputFile(output_dir, table)
	for _, row := range res {
		terraformName := toTerraformName(row.Name)
		fmt.Fprintf(tf_w, "resource \"dog_service\" \"%s\" {\n", terraformName)
		fmt.Fprintf(tf_w, "  name = \"%s\"\n", row.Name)
		fmt.Fprintf(tf_w, "  version = \"%d\"\n", row.Version)
		fmt.Fprintf(tf_w, "  services = [\n")
		portprotocols_output(tf_w, row.Services)
		fmt.Fprintf(tf_w, "  ]\n")
		fmt.Fprintf(tf_w, "  provider = dog.%s\n", environment)
		fmt.Fprintf(tf_w, "}\n")
		fmt.Fprintf(tf_w, "\n")

		fmt.Fprintf(import_w, "terraform import dog_service.%s %s\n", terraformName, row.ID)
	}
	tf_w.Flush()
	import_w.Flush()
}

func portprotocols_output(tf_w *bufio.Writer, portProtocols []*api.PortProtocol) {
	for _, port_protocol := range portProtocols {
		fmt.Fprintf(tf_w, "      {\n")
		fmt.Fprintf(tf_w, "        protocol = \"%s\"\n", port_protocol.Protocol)
		fmt.Fprintf(tf_w, strings.ReplaceAll(fmt.Sprintf("        ports = %q\n", port_protocol.Ports), "\" \"", "\",\""))
		fmt.Fprintf(tf_w, "      },\n")
	}
}

func zone_export(output_dir string, environment string) {
	fmt.Printf("zone_export\n")
	table := "zone"
	import_w := importOutputFile(output_dir, table)

	c := api.NewClient(os.Getenv("DOG_API_KEY"), os.Getenv("DOG_API_ENDPOINT"))

	res, statusCode, err := c.GetZones(nil)
	if err != nil {
		log.Fatalln(err)
	}
	if statusCode != 200 {
		log.Fatalln(err)
	}

	tf_w := terraformOutputFile(output_dir, table)
	for _, row := range res {
		terraformName := toTerraformName(row.Name)
		fmt.Fprintf(tf_w, "resource \"dog_zone\" \"%s\" {\n", terraformName)
		fmt.Fprintf(tf_w, "  name = \"%s\"\n", row.Name)
		fmt.Fprintf(tf_w, strings.ReplaceAll(fmt.Sprintf("  ipv4_addresses = %q\n", row.IPv4Addresses), "\" \"", "\",\""))
		fmt.Fprintf(tf_w, strings.ReplaceAll(fmt.Sprintf("  ipv6_addresses = %q\n", row.IPv6Addresses), "\" \"", "\",\""))
		fmt.Fprintf(tf_w, "  provider = dog.%s\n", environment)
		fmt.Fprintf(tf_w, "}\n")
		fmt.Fprintf(tf_w, "\n")

		fmt.Fprintf(import_w, "terraform import dog_zone.%s %s\n", terraformName, row.ID)
	}
	tf_w.Flush()
	import_w.Flush()
}

func profile_export(output_dir string, environment string) {
	fmt.Printf("profile_export\n")
	table := "profile"
	import_w := importOutputFile(output_dir, table)

	c := api.NewClient(os.Getenv("DOG_API_KEY"), os.Getenv("DOG_API_ENDPOINT"))

	res, statusCode, err := c.GetProfiles(nil)
	if err != nil {
		log.Fatalln(err)
	}
	if statusCode != 200 {
		log.Fatalln(err)
	}

	tf_w := terraformOutputFile(output_dir, table)
	for _, row := range res {
		terraformName := toTerraformName(row.Name)
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
		fmt.Fprintf(tf_w, "  }\n")
		fmt.Fprintf(tf_w, "  provider = dog.%s\n", environment)
		fmt.Fprintf(tf_w, "}\n")

		fmt.Fprintf(import_w, "terraform import dog_profile.%s %s\n", terraformName, row.ID)
	}
	tf_w.Flush()
	import_w.Flush()
}

func rules_output(tf_w *bufio.Writer, rules []*api.Rule) {
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

		fmt.Fprintf(tf_w, "      },\n")
	}
}

func main() {
	environment := os.Args[1]
	output_dir := "/tmp/dog-import"
	group_export(output_dir, environment)
	host_export(output_dir, environment)
	link_export(output_dir, environment)
	profile_export(output_dir, environment)
	service_export(output_dir, environment)
	zone_export(output_dir, environment)
	fmt.Printf("check %s/ for output files\n", output_dir)
}
