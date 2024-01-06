package dog_test

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

func TestAccDogFact_Basic(t *testing.T) {
	resourceType := "dog_fact"
	randomName := "tf_test_fact_" + acctest.RandString(5)
	resourceName := resourceType + "." + randomName

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:                 func() { testAccPreCheck(t) },
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testAccDogFactConfig_basic(resourceType, randomName),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet(resourceName, "id"),
					resource.TestCheckResourceAttr(resourceName, "name", randomName),
					resource.TestCheckResourceAttr(resourceName, "groups.all.vars.key", "value"),
					resource.TestCheckResourceAttr(resourceName, "groups.app.hosts.host1.key", "value"),
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

func testAccDogFactConfig_basic(resourceType, name string) string {
	return fmt.Sprintf(`
resource %[1]q %[2]q {
  name = %[2]q 
  groups = {
     all = {
       vars = {
		key = "value"
		key2 = "value2"
	}
	hosts = {
	  host1 = {
	    key = "value",
	    key2 = "value2"
	  }
	  host2 = {
	    key2 = "value2"
	  }
	},
	children = [
		"test"
	]
     },
     app = {
	vars = {
		key = "value"
	}
	hosts = {
	  host1 = {
	    key = "value"
	  }
	},
	children = [
		"test2"
	]
     }
  }
}
`, resourceType, name)
}
