package dog_test

import (
	"testing"
	"fmt"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

func TestProvider_DogInventoryNameAttribute(t *testing.T) {
	resource.UnitTest(t, resource.TestCase{
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		PreCheck:                 func() { testAccPreCheck(t) },
		Steps: []resource.TestStep{
			{
				Config: testAccDogInventoryDataSourceConfig("drew_test"),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("dog_inventory.drew_test", "name", "drew_test"),
				),
			},
		},
	})
}


func testAccDogInventoryDataSourceConfig(name string) string {
	return fmt.Sprintf(`
resource "dog_inventory" %[1]q {
  name = %[1]q 
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
`, name)
}
