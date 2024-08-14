// +build acceptance datasource host

package dog_test

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

func TestProvider_DogHostNameAttribute(t *testing.T) {
	resource.UnitTest(t, resource.TestCase{
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		PreCheck:                 func() { testAccPreCheck(t) },
		Steps: []resource.TestStep{
			{
				Config: testAccDogHostHostDataSourceConfig("terraform_host_test"),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("data.dog_host.terraform_host_test", "name", "terraform_host_test"),
				),
			},
		},
	})
}

func testAccDogHostDataSourceConfig(configurableAttribute string) string {
	profile := testAccDogHostProfileDataSourceConfig()
	group := testAccDogHostGroupDataSourceConfig()
	host := testAccDogHostHostDataSourceConfig(configurableAttribute)
	all := profile + group + host
	return all
}

func testAccDogHostProfileDataSourceConfig() string {
	return fmt.Sprintf(`
resource "dog_profile" "terraform_test_profile" {
  name = "terraform_test_profile"
  version = "1.0"
  ruleset_id = "123"
}
`)
}

func testAccDogHostGroupDataSourceConfig() string {
	return fmt.Sprintf(`
resource "dog_group" "terraform_test_group" {
  description = ""
  name = "terraform_test_group"
  profile_name = "terraform_test_profile"
  profile_version = "latest"
}
`)
}

func testAccDogHostHostDataSourceConfig(configurableAttribute string) string {
	return fmt.Sprintf(`
resource "dog_host" %[1]q {
  environment = "*"
  group = "dog_test"
  hostkey = "1726819861d5245b0afcd25127a7b181a5365620"
  location = "*"
  name = %[1]q
  vars = jsonencode({
      test = "dog_host"
  })
  alert_enable = false
}

data "dog_host" %[1]q {
  name = dog_host.%[1]s.name
  group = "dog_test"
  hostkey = "1726819861d5245b0afcd25127a7b181a5365620"
}
`, configurableAttribute)
}
