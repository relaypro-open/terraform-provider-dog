package dog

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

func TestAccHostResource(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { testAccPreCheck(t) },
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			// Create and Read testing
			{
				Config: testAccHostResourceConfig("one"),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("droid-qa-aws03", "active", "active"),
					resource.TestCheckResourceAttr("droid-qa-aws03", "environment", "*"),
				),
			},
			// ImportState testing
			{
				ResourceName:      "dog_host.droid-qa-aws03",
				ImportState:       true,
				ImportStateVerify: true,
				// This is not normally necessary, but is here because this
				// host code does not have an actual upstream service.
				// Once the Read method is able to refresh information from
				// the upstream service, this can be removed.
				// ImportStateVerifyIgnore: []string{"configurable_attribute"},
			},
			// Update and Read testing
			{
				Config: testAccHostResourceConfig("two"),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("dog_host.droid-qa-aws03", "active", "retired"),
				),
			},
			// Delete testing automatically occurs in TestCase
		},
	})
}

func testAccHostResourceConfig(configurableAttribute string) string {
	return fmt.Sprintf(`
resource "dog_host" "droid-qa-aws03" {
  name = %[1]q
  active = "active"
  environment = "*"
  group = "update_group"
  hostkey = "d4eb483c-99a1-11ec-bcff-03dfdfc9eeb8"
  location = "*"
}
`, configurableAttribute)
}