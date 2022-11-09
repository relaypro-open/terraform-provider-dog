package dog_test

import (
	"testing"
	"fmt"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

func TestAccDogZone_Basic(t *testing.T) {
	name := "dog_zone"
	randomName := acctest.RandomWithPrefix("tf-test-zone")
	resourceName := name + "." + randomName
	
	resource.ParallelTest(t, resource.TestCase{
		PreCheck:                 func() { testAccPreCheck(t) },
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testAccDogZoneConfig_basic(name, randomName),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet(resourceName, "id"),
					resource.TestCheckResourceAttr(resourceName, "name", randomName),
					resource.TestCheckResourceAttr(resourceName, "ipv4_addresses.0", "1.1.1.1"),
					resource.TestCheckResourceAttr(resourceName, "ipv6_addresses.#", "0"),
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


func testAccDogZoneConfig_basic(resourceName, name string) string {
	return fmt.Sprintf(`
resource %[1]q %[2]q {
  name = %[2]q
  ipv4_addresses = ["1.1.1.1"]
  ipv6_addresses = []
}
`, resourceName, name)
}
