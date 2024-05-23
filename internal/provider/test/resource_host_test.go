package dog_test

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

func TestAccDogHost_Basic(t *testing.T) {
	resourceType := "dog_host"
	randomName := "tf-test-host-" + acctest.RandString(5)
	resourceName := resourceType + "." + randomName

	resource.Test(t, resource.TestCase{
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
					resource.TestCheckResourceAttr(resourceName, "vars", "{\"key\":\"value\",\"key2\":\"value2\"}"),
					resource.TestCheckResourceAttr(resourceName, "alert_enable", "true"),
				),
			},
			{
				Config: testAccDogHostConfig_basic_remove_var(resourceType, randomName),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet(resourceName, "id"),
					resource.TestCheckResourceAttr(resourceName, "name", randomName),
					resource.TestCheckResourceAttr(resourceName, "environment", "*"),
					resource.TestCheckResourceAttr(resourceName, "hostkey", "1726819861d5245b0afcd25127a7b181a5365620"),
					resource.TestCheckResourceAttr(resourceName, "vars", "{\"key\":\"value\"}"),
					resource.TestCheckResourceAttr(resourceName, "alert_enable", "false"),
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
  environment = "*"
  group = "dog_test"
  hostkey = "1726819861d5245b0afcd25127a7b181a5365620"
  location = "*"
  name = %[2]q
  vars = jsonencode({
	  key = "value"
	  key2 = "value2"
  })
  alert_enable = true
}
`, resourceType, resourceName)
}

func testAccDogHostConfig_basic_remove_var(resourceType, resourceName string) string {
	return fmt.Sprintf(`
resource %[1]q %[2]q {
  environment = "*"
  group = "dog_test"
  hostkey = "1726819861d5245b0afcd25127a7b181a5365620"
  location = "*"
  name = %[2]q
  vars = jsonencode({
	  key = "value"
  })
  alert_enable = false
}
`, resourceType, resourceName)
}
