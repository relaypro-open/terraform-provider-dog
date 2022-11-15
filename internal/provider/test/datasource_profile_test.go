package dog_test

import (
	"testing"
	"fmt"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

func TestProvider_DogProfileNameAttribute(t *testing.T) {
	resource.UnitTest(t, resource.TestCase{
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		PreCheck:                 func() { testAccPreCheck(t) },
		Steps: []resource.TestStep{
			{
				Config: testAccDogProfileDataSourceConfig("drew_test"),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("dog_profile.drew_test", "name", "drew_test"),
				),
			},
		},
	})
}


func testAccDogProfileDataSourceConfig(configurableAttribute string) string {
	return fmt.Sprintf(`
resource dog_profile %[1]q {
  name = %[1]q
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
      },
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
      },
    ]
  }
}
`, configurableAttribute)
} 
