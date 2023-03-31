package dog_test

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

func TestProvider_DogLinkNameAttribute(t *testing.T) {
	resource.UnitTest(t, resource.TestCase{
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		PreCheck:                 func() { testAccPreCheck(t) },
		Steps: []resource.TestStep{
			{
				Config: testAccDogLinkDataSourceConfig("d1"),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("data.dog_link.d1", "name", "d1"),
				),
			},
		},
	})
}

func testAccDogLinkDataSourceConfig(configurableAttribute string) string {
	return fmt.Sprintf(`
resource "dog_link" %[1]q {
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
  name = %[1]q
}

data "dog_link" %[1]q {
  name = dog_link.%[1]s.name
}
`, configurableAttribute)
}
