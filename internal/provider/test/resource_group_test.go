// +build acceptance resource group

package dog_test

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

func TestAccDogGroup_Basic(t *testing.T) {
	name := "dog_group"
	randomName := "tf_test_group_" + acctest.RandString(5)
	resourceName := name + "." + randomName

	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { testAccPreCheck(t) },
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testAccDogGroupConfig_basic(name, randomName),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet(resourceName, "id"),
					resource.TestCheckResourceAttr(resourceName, "description", ""),
					resource.TestCheckResourceAttr(resourceName, "name", randomName),
					resource.TestCheckResourceAttr(resourceName, "profile_name", "resource_group"),
					resource.TestCheckResourceAttr(resourceName, "profile_version", "latest"),
					resource.TestCheckResourceAttr(resourceName, "vars", "{\"key\":\"value\",\"key2\":\"value2\"}"),
					resource.TestCheckResourceAttr(resourceName, "alert_enable", "true"),
				),
			},
			{
				Config: testAccDogGroupConfig_basic_remove_var(name, randomName),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet(resourceName, "id"),
					resource.TestCheckResourceAttr(resourceName, "description", ""),
					resource.TestCheckResourceAttr(resourceName, "name", randomName),
					resource.TestCheckResourceAttr(resourceName, "profile_name", "resource_group"),
					resource.TestCheckResourceAttr(resourceName, "profile_version", "latest"),
					resource.TestCheckResourceAttr(resourceName, "vars", "{\"key\":\"value\"}"),
					resource.TestCheckResourceAttr(resourceName, "alert_enable", "false"),
				),
			},
			{
				Config: testAccDogGroupConfig_basic_remove_esgs(name, randomName),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet(resourceName, "id"),
					resource.TestCheckResourceAttr(resourceName, "description", ""),
					resource.TestCheckResourceAttr(resourceName, "name", randomName),
					resource.TestCheckResourceAttr(resourceName, "profile_name", "resource_group"),
					resource.TestCheckResourceAttr(resourceName, "profile_version", "latest"),
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

func testAccDogGroupRulesetResourceConfig() string {
	return fmt.Sprintf(`
resource "dog_ruleset" "resource_group" {
  name = "resource_group"
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
        service = "any"
        states = []
        type = "BASIC"
      }
    ]
  }
}
`)
}

func testAccDogGroupProfileResourceConfig() string {
	return fmt.Sprintf(`
resource "dog_profile" "resource_group" {
  name = "resource_group"
  version = "1.0"
}
`)
}

func testAccDogGroupConfig_basic(name, randomName string) string {
	g := fmt.Sprintf(`
resource %[1]q %[2]q {
  description = ""
  name = %[2]q
  profile_id = dog_profile.resource_group.id
  profile_name = dog_profile.resource_group.name
  profile_version = "latest"
  ec2_security_group_ids = [
    { 
      region = "us-test-region"
      sgid = "sg-test"
    }
  ]
  vars = jsonencode({
	  key = "value"
	  key2 = "value2"
  })
  alert_enable = true
}
`, name, randomName)
	gr := testAccDogGroupRulesetResourceConfig()
	gp := testAccDogGroupProfileResourceConfig()

	to := gr + gp + g
	return to

}

func testAccDogGroupConfig_basic_remove_var(name, randomName string) string {
	g := fmt.Sprintf(`
resource %[1]q %[2]q {
  description = ""
  name = %[2]q
  profile_id = dog_profile.resource_group.id
  profile_name = dog_profile.resource_group.name
  profile_version = "latest"
  ec2_security_group_ids = [
    { 
      region = "us-test-region"
      sgid = "sg-test"
    }
  ]
  vars = jsonencode({
	  key = "value"
  })
  alert_enable = false
}
`, name, randomName)
	gr := testAccDogGroupRulesetResourceConfig()
	gp := testAccDogGroupProfileResourceConfig()

	to := gr + gp + g
	return to

}

func testAccDogGroupConfig_basic_remove_esgs(name, randomName string) string {
	g := fmt.Sprintf(`
resource %[1]q %[2]q {
  description = ""
  name = %[2]q
  profile_id = dog_profile.resource_group.id
  profile_name = dog_profile.resource_group.name
  profile_version = "latest"
  ec2_security_group_ids = [
  ]
  vars = jsonencode({
	  key = "value"
  })
  alert_enable = false
}
`, name, randomName)
	gr := testAccDogGroupRulesetResourceConfig()
	gp := testAccDogGroupProfileResourceConfig()

	to := gr + gp + g
	return to

}
