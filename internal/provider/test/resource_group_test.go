package dog_test

import (
	"testing"
	"fmt"

	//"github.com/hashicorp/terraform-plugin-sdk/v2/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

func TestAccDogGroup_Basic(t *testing.T) {
	name := "dog_group"
	//randomName := "tf_test_group_" + acctest.RandString(5)
	randomName := "tf_test_group_" + "1"
	resourceName := name + "." + randomName
	
	resource.ParallelTest(t, resource.TestCase{
		PreCheck:                 func() { testAccPreCheck(t) },
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testAccDogProfileConfig_basic("dog_profile", "terraform_test_profile"),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("dog_profile.terraform_test_profile", "name", "terraform_test_profile"),
				),
			},
			{
				Config: testAccDogGroupConfig_basic(name, randomName),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet(resourceName, "id"),
					resource.TestCheckResourceAttr(resourceName, "description", ""),
					resource.TestCheckResourceAttr(resourceName, "name", randomName),
					resource.TestCheckResourceAttr(resourceName, "profile_name", "terraform_test_profile"),
					resource.TestCheckResourceAttr(resourceName, "profile_version", "latest"),
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


func testAccDogGroupConfig_basic(name, randomName string) string {
	return fmt.Sprintf(`
resource %[1]q %[2]q {
  description = ""
  name = %[2]q
  profile_name = "terraform_test_profile"
  profile_version = "latest"
  ec2_security_group_ids = [
    { 
      region = "us-test-region"
      sgid = "sg-test"
    }
  ]
}
`, name, randomName)
}

func testAccDogProfileConfig_basic(resourceName, name string) string {
	return fmt.Sprintf(`
resource %[1]q %[2]q {
  name = %[2]q
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
`, resourceName, name)
}
