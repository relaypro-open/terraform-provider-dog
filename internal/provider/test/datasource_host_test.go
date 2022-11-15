package dog_test

import (
	"testing"
	"fmt"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

func TestProvider_DogHostNameAttribute(t *testing.T) {
	resource.UnitTest(t, resource.TestCase{
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		PreCheck:                 func() { testAccPreCheck(t) },
		Steps: []resource.TestStep{
			//{
			//	Config: testAccDogHostProfileDataSourceConfig(),
			//	Check: resource.ComposeAggregateTestCheckFunc(
			//		resource.TestCheckResourceAttr("dog_profile.terraform_test_profile", "name", "terraform_test_profile"),
			//	),
			//},
			//{
			//	Config: testAccDogHostGroupDataSourceConfig(),
			//	Check: resource.ComposeTestCheckFunc(
			//		resource.TestCheckResourceAttr("dog_group.terraform_test_group", "name", "terraform_test_group"),
			//	),
			//},
			{
				Config: testAccDogHostDataSourceConfig("terraform_host_test"),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("dog_profile.terraform_test_profile", "name", "terraform_test_profile"),
					resource.TestCheckResourceAttr("dog_group.terraform_test_group", "name", "terraform_test_group"),
					resource.TestCheckResourceAttr("dog_host.terraform_host_test", "name", "terraform_host_test"),
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
  rules = {
    inbound = [
      {
        action = "DROP"
        active = "true"
        comment = ""
        environments = []
        group = "any"
        group_type = "ANY"
        interface = ""
        log = "false"
        log_prefix = ""
        order = "2"
        service = "any"
        states = []
        type = "BASIC"
      }
    ]
    outbound = [
      {
        action = "ACCEPT"
        active = "true"
        comment = ""
        environments = []
        group = "any"
        group_type = "ANY"
        interface = ""
        log = "false"
        log_prefix = ""
        order = "1"
        service = "any"
        states = []
        type = "BASIC"
      }
    ]
  }
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
  group = "terraform_test_group"
  hostkey = "1726819861d5245b0afcd25127a7b181a5365620"
  location = "*"
  name = %[1]q
}
`, configurableAttribute)
} 
