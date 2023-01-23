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
				Config: testAccDogGroupDataSourceConfig("terraform_group_test"),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("dog_group.terraform_group_test", "name", "terraform_group_test"),
				),
			},
		},
	})
}


func testAccDogGroupDataSourceConfig(configurableAttribute string) string {
	return fmt.Sprintf(`
resource "dog_group" %[1]q {
  description = ""
  name = %[1]q 
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
`, configurableAttribute)
} 
