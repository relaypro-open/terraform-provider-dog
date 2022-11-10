package dog_test

import (
	"testing"
	"fmt"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

func TestAccDogLink_Basic(t *testing.T) {
	name := "dog_link"
	randomName := "d" + acctest.RandString(1)
	resourceName := name + "." + randomName
	
	resource.ParallelTest(t, resource.TestCase{
		PreCheck:                 func() { testAccPreCheck(t) },
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testAccDogLinkConfig_basic(name, randomName),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet(resourceName, "id"),
					resource.TestCheckResourceAttr(resourceName, "name", randomName),
					resource.TestCheckResourceAttr(resourceName, "address_handling", "union"),
					resource.TestCheckResourceAttr(resourceName, "conection.port", "5673"),
					resource.TestCheckResourceAttr(resourceName, "conection.%", "7"),
					resource.TestCheckResourceAttr(resourceName, "conection.ssl_options.%", "6"),
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


func testAccDogLinkConfig_basic(resourceName, name string) string {
	return fmt.Sprintf(`
resource %[1]q %[2]q {
  address_handling = "union"
  connection = {
    api_port = 15672
    host = "dog-broker.test.domain"
    password = "apassword"
    port = 5673
    ssl_options = {
        cacertfile = "certs/ca.crt"
        certfile = "certs/server.crt"
        fail_if_no_peer_cert = true
        keyfile = "private/server.key"
        server_name_indication = "disable"
        verify = "verify_peer"
      },
    user = "dog_trainer"
    virtual_host = "dog"
  }
  connection_type = "thumper"
  direction = "bidirectional"
  enabled = false
  name = %[2]q
}
`, resourceName, name)
}
