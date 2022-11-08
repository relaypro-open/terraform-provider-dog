package main

import (
	"context"
	"flag"
	"log"

	"github.com/hashicorp/terraform-plugin-framework/providerserver"
	dog "terraform-provider-dog/internal/provider"
)

var (
	// these will be set by the goreleaser configuration
	// to appropriate values for the compiled binary
	version string = "qa"
)

func main() {
	var debug bool

	flag.BoolVar(&debug, "debug", false, "set to true to run the provider with support for debuggers like delve")
	flag.Parse()

	opts := providerserver.ServeOpts{
		Address: "registry.terraform.io/relaypro-open/dog",
		Debug:   debug,
	}

	err := providerserver.Serve(context.Background(), dog.New(version), opts)

	if err != nil {
		log.Fatal(err.Error())
	}
}
