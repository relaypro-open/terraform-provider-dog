// +build acceptance datasource service

package dog_test

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

func TestProvider_DogServiceNameAttribute(t *testing.T) {
	resource.UnitTest(t, resource.TestCase{
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		PreCheck:                 func() { testAccPreCheck(t) },
		Steps: []resource.TestStep{
			{
				Config: testAccDogServiceDataSourceConfig("drew_test"),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("data.dog_service.drew_test", "name", "drew_test"),
				),
			},
		},
	})
}

func testAccDogServiceDataSourceConfig(configurableAttribute string) string {
	return fmt.Sprintf(`
resource "dog_service" %[1]q {
  name = %[1]q
  version = "1"
  services = [
      {
        protocol = "tcp"
        ports = ["22"]
      },
  ]
}

data "dog_service" %[1]q {
	name = dog_service.%[1]s.name
	services = [
	    {
	      protocol = "tcp"
	      ports = ["22"]
	    },
	]
}
`, configurableAttribute)
}
