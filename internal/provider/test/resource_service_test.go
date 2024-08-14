// +build acceptance resource service

package dog_test

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

func TestAccDogService_Basic(t *testing.T) {
	name := "dog_service"
	randomName := "tf_test_service_" + acctest.RandString(5)
	resourceName := name + "." + randomName

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:                 func() { testAccPreCheck(t) },
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testAccDogServiceConfig_basic(name, randomName),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet(resourceName, "id"),
					resource.TestCheckResourceAttr(resourceName, "name", randomName),
					resource.TestCheckResourceAttr(resourceName, "version", "1"),
					resource.TestCheckResourceAttr(resourceName, "services.#", "1"),
					resource.TestCheckResourceAttr(resourceName, "services.0.protocol", "tcp"),
					resource.TestCheckResourceAttr(resourceName, "services.0.ports.0", "22"),
				),
			},
			{
				ResourceName:      resourceName,
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

func testAccDogServiceConfig_basic(resourceName, name string) string {
	return fmt.Sprintf(`
resource %[1]q %[2]q {
  name = %[2]q
  version = "1"
  services = [
      {
        protocol = "tcp"
        ports = ["22"]
      },
  ]
}
`, resourceName, name)
}
