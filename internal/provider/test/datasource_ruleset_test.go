package dog_test

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

func TestProvider_DogRulesetNameAttribute(t *testing.T) {
	resourceType := "dog_ruleset"
	randomName := "datasource_" + acctest.RandString(5)
	resourceName := resourceType + "." + randomName

	resource.UnitTest(t, resource.TestCase{
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		PreCheck:                 func() { testAccPreCheck(t) },
		Steps: []resource.TestStep{
			{
				Config: testAccDogRulesetDataSourceConfig(resourceType, randomName),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(resourceName, "name", randomName),
				),
			},
		},
	})
}

func testAccDogRulesetDataSourceConfig(resourceType, randomName string) string {
	return fmt.Sprintf(`
resource %[1]q %[2]q {
  name = %[2]q
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
`, resourceType, randomName)
}
