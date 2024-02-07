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

	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { testAccPreCheck(t) },
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testAccDogFactConfig_basic(resourceType, randomName),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet(resourceName, "id"),
					resource.TestCheckResourceAttr(resourceName, "name", randomName),
					resource.TestCheckResourceAttr(resourceName, "groups.all.vars", "{\"key\":\"value\",\"key2\":\"value2\"}"),
					resource.TestCheckResourceAttr(resourceName, "groups.app.hosts.host1.key", "value"),
				),
			},
			{
				Config: testAccDogFactConfig_add_vars(resourceType, randomName),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet(resourceName, "id"),
					resource.TestCheckResourceAttr(resourceName, "name", randomName),
					resource.TestCheckResourceAttr(resourceName, "groups.all.vars", "{\"key\":\"value\",\"key2\":\"value2\",\"key3\":\"value3 + 3\"}"),
					resource.TestCheckResourceAttr(resourceName, "groups.app.hosts.host1.key", "value"),
				),
			},
			{
				Config: testAccDogFactConfig_remove_vars(resourceType, randomName),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet(resourceName, "id"),
					resource.TestCheckResourceAttr(resourceName, "groups.all.vars", "{\"key2\":\"value2\",\"key3\":\"value3\"}"),
					resource.TestCheckResourceAttr(resourceName, "name", randomName),
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
       vars = jsonencode({
		key = "value"
		key2 = "value2"
	})
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
	vars = jsonencode({
		key = "value"
	})
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

func testAccDogFactConfig_add_vars(resourceType, name string) string {
	return fmt.Sprintf(`
resource %[1]q %[2]q {
  name = %[2]q 
  groups = {
     all = {
       vars = jsonencode({
		key = "value"
		key2 = "value2"
		key3 = "value3 + 3"
	})
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
	vars = jsonencode({
		key = "value"
		key2 = "value2"
	})
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

func testAccDogFactConfig_remove_vars(resourceType, name string) string {
	return fmt.Sprintf(`
resource %[1]q %[2]q {
  name = %[2]q 
  groups = {
     all = {
       vars = jsonencode({
		key2 = "value2"
		key3 = "value3"
	})
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
	vars = jsonencode({
		key = "value"
	})
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

func TestAccDogFact_Big(t *testing.T) {
	resourceType := "dog_fact"
	randomName := "tf_test_fact_big_" + acctest.RandString(5)
	resourceName := resourceType + "." + randomName

	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { testAccPreCheck(t) },
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testAccDogFactConfig_big(resourceType, randomName),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet(resourceName, "id"),
					resource.TestCheckResourceAttr(resourceName, "name", randomName),
					resource.TestCheckResourceAttr(resourceName, "groups.all.hosts.host1.key", "value"),
				),
			},
			{
				Config: testAccDogFactConfig_big_update(resourceType, randomName),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet(resourceName, "id"),
					resource.TestCheckResourceAttr(resourceName, "name", randomName),
					resource.TestCheckResourceAttr(resourceName, "groups.all.hosts.host1.key", "value"),
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

func testAccDogFactConfig_big(resourceType, name string) string {
	return fmt.Sprintf(`
resource %[1]q %[2]q {
  name = %[2]q 
    groups = {
      all = {
        children = ["dog_test"]
        hosts = {
          host2 = {
            key2 = "value2"
          }
          host1 = {
            key = "value"
            key2 = "value2"
          }
        }
        vars = jsonencode({
          subnet_qa_public_us_east_1a = "subnet-50f9cc09"
          ssh_ca_fingerprint = "69f6617f461ceb6d05cd57a0b9e05ee27f555cfaaae5b52f6760bd9bcc976920"
          product = "foo"
          peering_id = "pcx-52b2c651"
          logger_elb = "qa.logger.foo.io"
          provider = "ec2"
          subnet_qa_public_us_east_1e = "subnet-c5f2a9ee"
          consul_cluster = "{{ region }}"
          ec2_dns_suffix = "compute-1"
          scout_environment = "QA"
          cluster = "mob"
          foopro_domain = "foogo.info"
          email_alert_distro = []
          subnet_qa_public_us_east_1d = "subnet-6fc0d512"
          common_build_root_dir = "/home/{{ ansible_user }}"
          vpc_id = "vpc-2f0fd9eb"
          cluster_separator = "_"
          filebeat_version = "oss-2.2.0"
          region = "us-east-1"
          foogo_domain = "foogo.info"
          vm_mms_numbers = ["9199255565"]
          subnet_qa_public_us_east_2a = "subnet-67fc020e"
          subnet_qa_public_us_west_2c = "subnet-d255c121"
          nginx_config = "republicwireless-com"
          subnet_qa_public_us_west_2b = "subnet-1dabf772"
          foo_dash_company_id = "{{ foo_iris_account }}"
          credstash_table = "{{'credential-store_' + cluster + cluster_separator + env}}"
          service_domain = "foo.io"
          ssh_ca_provisioner = "site-reliability-engineering@foopro.com"
          riak_http_port = 2067
          dog_env = "qa"
          telegraf_interval = 500
          env = "qa"
          elasticsearch_domain = "vpc-logger-qa-m52rd66iwl6df5gieiylyaso4e.us-east-1.es.amazonaws.com"
          cert_name = "star-republicdev-info"
          cluster_env = "mob_qa"
          igw_id = "igw-1ad5e97f"
          monit_environment = "qa"
          ansible_python_interpreter = "python"
          domain = "phoneboothdev.info"
          keypair = "bandwidithaws"
          foo_iris_account = 9900004
          subnet_qa_public_us_east_2b = "subnet-12151b60"
          image = "ami-09150b5c79250e26c"
          subnet_qa_public_us_west_2a = "subnet-44b2ca55"
          subnet_qa_public_us_east_1c = "subnet-2ded22b0"
          dog_env_group = "www-data"
          metrics_agent = "collectd"
          subnet_qa_public_us_east_2c = "subnet-f25c12b2"
          ssh_ca_url = "https://api-qa.ca.phoneboothdev.info/"
      })
    }
  }
}
`, resourceType, name)
}

func testAccDogFactConfig_big_update(resourceType, name string) string {
	return fmt.Sprintf(`
resource %[1]q %[2]q {
  name = %[2]q 
    groups = {
      all = {
        children = ["dog_test"]
        hosts = {
          host2 = {
            key2 = "value2"
          }
          host1 = {
            key = "value"
            key2 = "value2"
          }
        }
        vars = jsonencode({
          subnet_qa_public_us_east_1a = "subnet-50f9cc09"
          ssh_ca_fingerprint = "69f6617f461ceb6d05cd57a0b9e05ee27f555cfaaae5b52f6760bd9bcc976920"
          product = "foo"
          peering_id = "pcx-52b2c651"
          logger_elb = "qa.logger.foo.io"
          provider = "ec2"
          subnet_qa_public_us_east_1e = "subnet-c5f2a9ee"
          consul_cluster = "{{ region }}"
          ec2_dns_suffix = "compute-1"
          scout_environment = "QA"
          cluster = "mob"
          foopro_domain = "foogo.info"
          email_alert_distro = []
          subnet_qa_public_us_east_1d = "subnet-6fc0d512"
          common_build_root_dir = "/home/{{ ansible_user }}"
          vpc_id = "vpc-2f0fd9eb"
          cluster_separator = "_"
          filebeat_version = "oss-2.2.0"
          region = "us-east-1"
          foogo_domain = "foogo.info"
          vm_mms_numbers = ["9199255565"]
          subnet_qa_public_us_east_2a = "subnet-67fc020e"
          subnet_qa_public_us_west_2c = "subnet-d255c121"
          nginx_config = "republicwireless-com"
          subnet_qa_public_us_west_2b = "subnet-1dabf772"
          foo_dash_company_id = "{{ foo_iris_account }}"
          credstash_table = "{{'credential-store_' + cluster + cluster_separator + env}}"
          service_domain = "foo.io"
          ssh_ca_provisioner = "site-reliability-engineering@foopro.com"
          riak_http_port = 2067
          dog_env = "qa"
          telegraf_interval = 500
          env = "qa"
          elasticsearch_domain = "vpc-logger-qa-m52rd66iwl6df5gieiylyaso4e.us-east-1.es.amazonaws.com"
          cert_name = "star-republicdev-info"
          cluster_env = "mob_qa"
          igw_id = "igw-1ad5e97f"
          monit_environment = "qa"
          ansible_python_interpreter = "python"
          domain = "phoneboothdev.info"
          keypair = "bandwidithaws"
          foo_iris_account = 9900004
          subnet_qa_public_us_east_2b = "subnet-12151b60"
          image = "ami-09150b5c79250e26c"
          subnet_qa_public_us_west_2a = "subnet-44b2ca55"
          subnet_qa_public_us_east_1c = "subnet-2ded22b0"
          dog_env_group = "www-data"
          metrics_agent = "collectd"
          subnet_qa_public_us_east_2c = "subnet-f25c12b2"
          ssh_ca_url = "https://api-qa.ca.phoneboothdev.info/"
		  test = "best"
      })
    }
  }
}
`, resourceType, name)
}
