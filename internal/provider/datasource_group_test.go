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
				Config: testAccDogGroupDataSourceConfig("drew_test"),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("dog_group.drew_test", "name", "drew_test"),
				),
			},
		},
	})
}


func testAccDogGroupDataSourceConfig(configurableAttribute string) string {
	return fmt.Sprintf(`
resource "dog_group" "drew_test" {
  description = ""
  name = %[1]q 
  profile_name = "test_qa"
  profile_version = "latest"
}
`, configurableAttribute)
} 
