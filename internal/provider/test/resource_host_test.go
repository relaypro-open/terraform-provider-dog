package dog_test

import (
	"testing"
	"fmt"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

func TestAccDogHost_Basic(t *testing.T) {
	resourceType := "dog_host"
	randomName := "tf-test-host-" + acctest.RandString(5)
	resourceName := resourceType + "." + randomName
	
	resource.ParallelTest(t, resource.TestCase{
		PreCheck:                 func() { testAccPreCheck(t) },
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testAccDogHostConfig_basic(resourceType, randomName),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet(resourceName, "id"),
					resource.TestCheckResourceAttr(resourceName, "name", randomName),
					resource.TestCheckResourceAttr(resourceName, "environment", "*"),
					resource.TestCheckResourceAttr(resourceName, "hostkey", "1726819861d5245b0afcd25127a7b181a5365620"),
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


func testAccDogHostConfig_basic(resourceType, resourceName string) string {
	return fmt.Sprintf(`
resource %[1]q %[2]q {
  name = %[2]q
  environment = "*"
  group = "dog_test"
  hostkey = "1726819861d5245b0afcd25127a7b181a5365620"
  location = "*"
  vars = {
	  test = "dog_host"
  }
}
`, resourceType, resourceName)
}
