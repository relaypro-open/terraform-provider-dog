package dog_test

import (
	"fmt"
	"os"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/relaypro-open/dog_api_golang/api"
	//"github.com/davecgh/go-spew/spew"
)

func TestProvider_DogProfileNameAttribute(t *testing.T) {
	//rulesetResourceType := "dog_ruleset"
	//rulesetRandomName := "datasource_" + acctest.RandString(5)
	//rulesetResourceName := rulesetResourceType + "." + rulesetRandomName

	resourceType := "dog_profile"
	randomName := "datasource_" + acctest.RandString(5)
	resourceName := resourceType + "." + randomName

	resource.UnitTest(t, resource.TestCase{
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		PreCheck:                 func() { testAccPreCheck(t) },
		Steps: []resource.TestStep{
			//{
			//	Config: testAccDogProfileRulesetDataSourceConfig(t, rulesetResourceType, rulesetRandomName, &ruleset_id),
			//	Check: resource.ComposeTestCheckFunc(
			//		resource.TestCheckResourceAttr(rulesetResourceName, "name", rulesetRandomName),
			//	),
			//},
			{
				Config: testAccDogProfileDataSourceConfig(resourceType, randomName),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(resourceName, "name", randomName),
				),
			},
		},
	})
}

func testAccDogProfileRulesetDataSourceConfig(t *testing.T, rulesetResourceType, rulesetRandomName string, ruleset_id *string) string {
	//t.Log(fmt.Sprintf("ZZZtestAccDogProfileRulesetDataSourceConfig")
	c := api.NewClient(os.Getenv("DOG_API_TOKEN"), os.Getenv("DOG_API_ENDPOINT"))

	newRule := api.RulesetCreateRequest{
		Rules: &api.Rules{
			Inbound: []*api.Rule{
				&api.Rule{
					Action:       "ACCEPT",
					Active:       true,
					Comment:      "",
					Environments: []string{},
					Group:        "any",
					GroupType:    "ANY",
					Interface:    "",
					Log:          false,
					LogPrefix:    "",
					Order:        1,
					Service:      "any",
					States:       []string{},
					Type:         "BASIC",
				},
			},
			Outbound: []*api.Rule{
				&api.Rule{
					Action:       "DROP",
					Active:       true,
					Comment:      "",
					Environments: []string{},
					Group:        "any",
					GroupType:    "ANY",
					Interface:    "",
					Log:          false,
					LogPrefix:    "",
					Order:        1,
					Service:      "any",
					States:       []string{},
					Type:         "BASIC",
				},
			},
		},
		Name:    "name",
	}
	
	//t.Log(spew.Sprintf("ZZZZZZZZZZZZZZZZZZZnewRule: %v", newRule))

	res, _, _ := c.CreateRuleset(newRule, nil)
	//t.Log(fmt.Sprintf("ZZZres: %v", res))
	//str := fmt.Sprintf("%#v", res)
	//t.Log(fmt.Sprintf("ZZZstr: %s", str))
	*ruleset_id = res.ID

	//bogus return .tf - just need to return any valid resource syntax
	return fmt.Sprintf(`
resource %[1]q %[2]q {
  name = %[2]q
  rules = {
    inbound = [
    ]
    outbound = [
    ]
  }
}
`, rulesetResourceType, rulesetRandomName )
}

func testAccDogProfileDataSourceConfig(name, randomName string) string {
	return fmt.Sprintf(`
resource %[1]q %[2]q {
  name = %[2]q
  version = "1.0"
}
`, name, randomName)
} 
