package dog_test

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

func TestProvider_DogFactNameAttribute(t *testing.T) {
	resource.UnitTest(t, resource.TestCase{
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		PreCheck:                 func() { testAccPreCheck(t) },
		Steps: []resource.TestStep{
			{
				Config: testAccDogFactDataSourceConfig("drew_test"),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("data.dog_fact.drew_test", "name", "drew_test"),
				),
			},
		},
	})
}

func testAccDogFactDataSourceConfig(name string) string {
	return fmt.Sprintf(`
resource "dog_fact" %[1]q {
  name = %[1]q 
  groups = {
		all= {
			vars = jsonencode({
				key = "value",
				key2 = "value2"
			}),
			hosts = {
				host1 = {
					key = "value",
					key2 = "value2"
				},
				host2 = {
					key2 = "value2"
				}
			},
			children = [
				"test"
			]
		},
		app = {
			vars = jsonencode({
				key = "value"
			}),
			hosts = {
				host1 = {
					key = "value"
				}
			},
			children = [
				"test2"
			]
		}
	}
}


data "dog_fact" %[1]q {
  name = dog_fact.%[1]s.name 
}
`, name)
}
