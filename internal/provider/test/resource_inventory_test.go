package dog_test

import (
	"testing"
	"fmt"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

func TestAccDogInventory_Basic(t *testing.T) {
	resourceType := "dog_inventory"
	randomName := "tf_test_inventory_" + acctest.RandString(5)
	resourceName := resourceType + "." + randomName
	
	resource.ParallelTest(t, resource.TestCase{
		PreCheck:                 func() { testAccPreCheck(t) },
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testAccDogInventoryConfig_basic(resourceType, randomName),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet(resourceName, "id"),
					resource.TestCheckResourceAttr(resourceName, "name", randomName),
					resource.TestCheckResourceAttr(resourceName, "groups.0.name", "all"),
					resource.TestCheckResourceAttr(resourceName, "groups.0.vars.key", "value"),
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


func testAccDogInventoryConfig_basic(resourceType, name string) string {
	return fmt.Sprintf(`
resource %[1]q %[2]q {
  name = %[2]q 
  groups = [
     {
       name = "all"
       vars = {
        	key = "value"
          	key2 = "value2"
        }
     	hosts = {
          host1 = {
            key = "value",
            key2 = "value2"
          }
          host2 = {
            key2 = "value2"
          }
        },
	children = [
		"test"
	]
     },
     {
       name = "app"
     	vars = {
        	key = "value"
        }
     	hosts = {
          host1 = {
            key = "value"
          }
        },
	children = [
		"test2"
	]
     }
  ]
}
`, resourceType, name)
}
