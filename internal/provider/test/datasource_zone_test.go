package dog_test

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
)

func TestProvider_DogZoneNameAttribute(t *testing.T) {
	resource.UnitTest(t, resource.TestCase{
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		PreCheck:                 func() { testAccPreCheck(t) },
		Steps: []resource.TestStep{
			{
				Config: testAccDogZoneDataSourceConfig("drew_test"),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("data.dog_zone.drew_test", "name", "drew_test"),
				),
			},
		},
	})
}

func testAccDogZoneDataSourceConfig(configurableAttribute string) string {
	return fmt.Sprintf(`
resource "dog_zone" "drew_test" {
  name = %[1]q
  ipv4_addresses = ["1.1.1.1"]
  ipv6_addresses = []
}

data "dog_zone" "drew_test" {
  name = dog_zone.%[1]s.name
}
`, configurableAttribute)
}
