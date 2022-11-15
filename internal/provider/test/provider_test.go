package dog_test

import (
	"os"
	"testing"

	"github.com/hashicorp/terraform-plugin-framework/providerserver"
	"github.com/hashicorp/terraform-plugin-go/tfprotov6"
	dog "terraform-provider-dog/internal/provider"
)

// testAccProtov6ProviderFactories are used to instantiate a provider during
// acceptance testing. The factory function will be invoked for every Terraform
// CLI command executed to create a provider server to which the CLI can
// reattach
var testAccProtoV6ProviderFactories = map[string]func() (tfprotov6.ProviderServer, error){
	"dog": providerserver.NewProtocol6WithError(dog.New("dog")()),
}

func testAccPreCheck(t *testing.T) {
	if v := os.Getenv("DOG_API_ENDPOINT"); v == "" {
		t.Fatal("DOG_API_ENDPOINT must be set to run acceptance tests.")
	}

	if v := os.Getenv("DOG_API_KEY"); v == "" {
		t.Fatal("DOG_API_KEY must be set to run acceptance tests.")
	}
}
