package dog

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

func TestAccHostDataSource(t *testing.T) {
	name := "test"
	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { testAccPreCheck(t) },
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			// Read testing
			{
				Config: testAccHostDataSourceConfig(name),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("dog_host.test", "active", "active"),
					resource.TestCheckResourceAttr("dog_host.test", "environment", "*"),
					resource.TestCheckResourceAttr("dog_host.test", "group", "build_slave"),
					resource.TestCheckResourceAttr("dog_host.test", "hostkey", "d4eb483c-99a1-11ec-bcff-03dfdfc9eeb8"),
					resource.TestCheckResourceAttr("dog_host.test", "location", "*"),
					resource.TestCheckResourceAttr("dog_host.test", "name", name),
				),
			},
		},
	})
}

func testAccHostDataSourceConfig(name string) string {
	return fmt.Sprintf(`
terraform {
  required_providers {
    dog = {
      source = "github.com/relaypro-open/dog"
    }
  }
}
  provider "dog" {
    api_key = "my-key"
    api_endpoint = "http://dog-ubuntu-server.lxd:7070/api/V2"
  }

resource "dog_host" "test" {
  active =  "active"
  environment = "*"
  group = "build_slave"
  hostkey = "d4eb483c-99a1-11ec-bcff-03dfdfc9eeb8"
  location = "*"
  name = "%s"
}
`, name)
}
