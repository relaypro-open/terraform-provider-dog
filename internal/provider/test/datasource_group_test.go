package dog_test

import (
	"testing"
	"fmt"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

func TestProvider_DogGroupNameAttribute(t *testing.T) {
	resource.UnitTest(t, resource.TestCase{
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		PreCheck:                 func() { testAccPreCheck(t) },
		Steps: []resource.TestStep{
			{
				Config: testAccDogGroupDataSourceConfig(t, "terraform_group_test"),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("dog_group.terraform_group_test", "name", "terraform_group_test"),
				),
			},
		},
	})
}

func testAccDogGroupRulesetDataSourceConfig() string {
	return fmt.Sprintf(`
resource "dog_ruleset" "datasource_group" {
  name = "datasource_group"
  rules = {
    inbound = [
      {
        action = "ACCEPT"
        active = "true"
        comment = "test_zone"
        environments = []
        group = "dog_test"
        group_type = "ROLE"
        interface = ""
        log = "false"
        log_prefix = ""
        order = "1"
        service = "ssh-tcp-22"
        states = []
        type = "BASIC"
      },
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

func testAccDogGroupProfileDataSourceConfig() string {
	return fmt.Sprintf(`
resource "dog_profile" "datasource_group" {
  name = "datasource_group"
  version = "1.0"
}
`)
} 

func testAccDogGroupDataSourceConfig(t *testing.T, configurableAttribute string) string {
	g := fmt.Sprintf(`
resource "dog_group" %[1]q {
  description = ""
  name = %[1]q 
  profile_id = dog_profile.datasource_group.id
  profile_name = "test_qa"
  profile_version = "latest"
  ec2_security_group_ids = [
    { 
      region = "us-test-region"
      sgid = "sg-test"
    }
  ]
  vars = {
	  test = "dog_group"
  }
}
` , configurableAttribute)
  gr := testAccDogGroupRulesetDataSourceConfig() 
  gp := testAccDogGroupProfileDataSourceConfig()

  to := gr + gp + g
  //t.Log(fmt.Sprintf("to: %s", to))
  return to
}
