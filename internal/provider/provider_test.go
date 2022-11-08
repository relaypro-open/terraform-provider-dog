package dog

import (
	"os"
	"testing"
	"fmt"

	"github.com/hashicorp/terraform-plugin-framework/providerserver"
	"github.com/hashicorp/terraform-plugin-go/tfprotov6"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

// testAccProtov6ProviderFactories are used to instantiate a provider during
// acceptance testing. The factory function will be invoked for every Terraform
// CLI command executed to create a provider server to which the CLI can
// reattach
var testAccProtoV6ProviderFactories = map[string]func() (tfprotov6.ProviderServer, error){
	"dog": providerserver.NewProtocol6WithError(New("dog")()),
}

func testAccPreCheck(t *testing.T) {
	if v := os.Getenv("DOG_API_ENDPOINT"); v == "" {
		t.Fatal("DOG_API_ENDPOINT must be set to run acceptance tests.")
	}

	if v := os.Getenv("DOG_API_KEY"); v == "" {
		t.Fatal("DOG_API_KEY must be set to run acceptance tests.")
	}
}

func TestProvider_ZoneNameAttribute(t *testing.T) {
	resource.UnitTest(t, resource.TestCase{
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		PreCheck:                 func() { testAccPreCheck(t) },
		Steps: []resource.TestStep{
			{
				Config: testAccExampleResourceConfig("drew_test"),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("dog_zone.drew_test", "name", "drew_test"),
				),
			},
		},
	})
}


func testAccExampleResourceConfig(configurableAttribute string) string {
	return fmt.Sprintf(`
resource "dog_zone" "drew_test" {
  name = %[1]q
  ipv4_addresses = ["1.1.1.1"]
  ipv6_addresses = []
}
`, configurableAttribute)
} 
