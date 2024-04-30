//go:build acceptance || profile

package dog_test

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/hashicorp/terraform-plugin-testing/helper/acctest"
)

func TestAccDogProfile_Basic(t *testing.T) {
	name := "dog_profile"
	randomName := "resource_" + acctest.RandString(5)
	resourceName := name + "." + randomName

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:                 func() { testAccPreCheck(t) },
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testAccDogProfileConfig_basic(name, randomName),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet(resourceName, "id"),
					resource.TestCheckResourceAttr(resourceName, "version", "1.0"),
					resource.TestCheckResourceAttr(resourceName, "name", randomName),
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

func testAccDogProfileConfig_basic(name, resourceName string) string {
	return fmt.Sprintf(`
resource %[1]q %[2]q {
  name = %[2]q
  version = "1.0"
}
`, name, resourceName)
}
