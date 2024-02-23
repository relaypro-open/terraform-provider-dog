package dog_test

import (
        "fmt"
        "testing"

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
                                        resource.TestCheckResourceAttr("data.dog_group.terraform_group_test", "name", "terraform_group_test"),
                                ),
                        },
                },
        })
}

func testAccDogGroupDataSourceConfig(t *testing.T, configurableAttribute string) string {
        g := fmt.Sprintf(`
resource "dog_group" %[1]q {
  name = %[1]q
  description = ""
  profile_id = "1234"
  profile_name = "test_qa"
  profile_version = "latest"
  ec2_security_group_ids = [
    {
      region = "us-test-region"
      sgid = "sg-test"
    }
  ]
  vars = jsonencode({
      test = "dog_group"
  })
}

data "dog_group" %[1]q {
      name = dog_group.%[1]s.name
      profile_id = "1234"
}
`, configurableAttribute)
      return g
}
