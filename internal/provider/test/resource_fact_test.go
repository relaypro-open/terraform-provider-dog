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
					resource.TestCheckResourceAttr(resourceName, "groups.app.hosts", "{\"host1\":{\"key\":\"value\"}}"),
				),
			},
			{
				Config: testAccDogFactConfig_add_vars(resourceType, randomName),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet(resourceName, "id"),
					resource.TestCheckResourceAttr(resourceName, "name", randomName),
					resource.TestCheckResourceAttr(resourceName, "groups.all.vars", "{\"key\":\"value\",\"key2\":\"value2\",\"key3\":\"value3 + 3\"}"),
					resource.TestCheckResourceAttr(resourceName, "groups.app.hosts", "{\"host1\":{\"key\":\"value\"}}"),
				),
			},
			{
				Config: testAccDogFactConfig_remove_vars(resourceType, randomName),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet(resourceName, "id"),
					resource.TestCheckResourceAttr(resourceName, "groups.all.vars", "{\"key2\":\"value2\",\"key3\":\"value3\"}"),
					resource.TestCheckResourceAttr(resourceName, "name", randomName),
					resource.TestCheckResourceAttr(resourceName, "groups.app.hosts", "{\"host1\":{\"key\":\"value\"}}"),
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
	hosts = jsonencode({
	  host1 = {
	    key = "value",
	    key2 = "value2"
	  }
	  host2 = {
	    key2 = "value2"
	  }
	}),
	children = [
		"test"
	]
     },
     app = {
	vars = jsonencode({
		key = "value"
	})
	hosts = jsonencode({
	  host1 = {
	    key = "value"
	  }
	}),
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
	hosts = jsonencode({
	  host1 = {
	    key = "value",
	    key2 = "value2"
	  }
	  host2 = {
	    key2 = "value2"
	  }
	}),
	children = [
		"test"
	]
     },
     app = {
	vars = jsonencode({
		key = "value"
		key2 = "value2"
	})
	hosts = jsonencode({
	  host1 = {
	    key = "value"
	  }
	}),
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
	hosts = jsonencode({
	  host1 = {
	    key = "value",
	    key2 = "value2"
	  }
	  host2 = {
	    key2 = "value2"
	  }
	}),
	children = [
		"test"
	]
     },
     app = {
	vars = jsonencode({
		key = "value"
	})
	hosts = jsonencode({
	  host1 = {
	    key = "value"
	  }
	}),
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
					resource.TestCheckResourceAttr(resourceName, "groups.all.hosts", "{\"host1\":{\"key\":\"value\",\"key2\":\"value2\"},\"host2\":{\"key2\":\"value2\"}}"),
				),
			},
			{
				Config: testAccDogFactConfig_big_update(resourceType, randomName),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet(resourceName, "id"),
					resource.TestCheckResourceAttr(resourceName, "name", randomName),
					resource.TestCheckResourceAttr(resourceName, "groups.all.hosts", "{\"host1\":{\"key\":\"value\",\"key2\":\"value2\"},\"host2\":{\"key2\":\"value2\"}}"),
					resource.TestCheckResourceAttr(resourceName, "groups.all.vars", "{\"ansible_python_interpreter\":\"python\",\"cert_name\":\"star-republicdev-info\",\"cluster\":\"mob\",\"cluster_env\":\"mob_qa\",\"cluster_separator\":\"_\",\"common_build_root_dir\":\"/home/{{ ansible_user }}\",\"consul_cluster\":\"{{ region }}\",\"credstash_table\":\"{{'credential-store_' + cluster + cluster_separator + env}}\",\"dog_env\":\"qa\",\"dog_env_group\":\"www-data\",\"domain\":\"phoneboothdev.info\",\"ec2_dns_suffix\":\"compute-1\",\"elasticsearch_domain\":\"vpc-logger-qa-m52rd66iwl6df5gieiylyaso4e.us-east-1.es.amazonaws.com\",\"email_alert_distro\":[],\"env\":\"qa\",\"filebeat_version\":\"oss-2.2.0\",\"foo_dash_company_id\":\"{{ foo_iris_account }}\",\"foo_iris_account\":9900004,\"foogo_domain\":\"foogo.info\",\"foopro_domain\":\"foogo.info\",\"igw_id\":\"igw-1ad5e97f\",\"image\":\"ami-09150b5c79250e26c\",\"keypair\":\"bandwidithaws\",\"logger_elb\":\"qa.logger.foo.io\",\"metrics_agent\":\"collectd\",\"monit_environment\":\"qa\",\"nginx_config\":\"republicwireless-com\",\"peering_id\":\"pcx-52b2c651\",\"product\":\"foo\",\"provider\":\"ec2\",\"region\":\"us-east-1\",\"riak_http_port\":2067,\"scout_environment\":\"QA\",\"service_domain\":\"foo.io\",\"ssh_ca_fingerprint\":\"69f6617f461ceb6d05cd57a0b9e05ee27f555cfaaae5b52f6760bd9bcc976920\",\"ssh_ca_provisioner\":\"site-reliability-engineering@foopro.com\",\"ssh_ca_url\":\"https://api-qa.ca.phoneboothdev.info/\",\"subnet_qa_public_us_east_1a\":\"subnet-50f9cc09\",\"subnet_qa_public_us_east_1c\":\"subnet-2ded22b0\",\"subnet_qa_public_us_east_1d\":\"subnet-6fc0d512\",\"subnet_qa_public_us_east_1e\":\"subnet-c5f2a9ee\",\"subnet_qa_public_us_east_2a\":\"subnet-67fc020e\",\"subnet_qa_public_us_east_2b\":\"subnet-12151b60\",\"subnet_qa_public_us_east_2c\":\"subnet-f25c12b2\",\"subnet_qa_public_us_west_2a\":\"subnet-44b2ca55\",\"subnet_qa_public_us_west_2b\":\"subnet-1dabf772\",\"subnet_qa_public_us_west_2c\":\"subnet-d255c121\",\"telegraf_interval\":500,\"test\":\"best\",\"vm_mms_numbers\":[\"9199255565\"],\"vpc_id\":\"vpc-2f0fd9eb\"}"),
				),
			},
			{
				Config: testAccDogFactConfig_big_update_remove(resourceType, randomName),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet(resourceName, "id"),
					resource.TestCheckResourceAttr(resourceName, "name", randomName),
					resource.TestCheckResourceAttr(resourceName, "groups.all.hosts", "{\"host1\":{\"key\":\"value\",\"key2\":\"value2\"},\"host2\":{\"key2\":\"value2\"}}"),
					resource.TestCheckResourceAttr(resourceName, "groups.all.vars", "{\"ansible_python_interpreter\":\"python\",\"cert_name\":\"star-republicdev-info\",\"cluster\":\"mob\",\"cluster_env\":\"mob_qa\",\"cluster_separator\":\"_\",\"common_build_root_dir\":\"/home/{{ ansible_user }}\",\"consul_cluster\":\"{{ region }}\",\"credstash_table\":\"{{'credential-store_' + cluster + cluster_separator + env}}\",\"dog_env\":\"qa\",\"dog_env_group\":\"www-data\",\"domain\":\"phoneboothdev.info\",\"ec2_dns_suffix\":\"compute-1\",\"elasticsearch_domain\":\"vpc-logger-qa-m52rd66iwl6df5gieiylyaso4e.us-east-1.es.amazonaws.com\",\"email_alert_distro\":[],\"env\":\"qa\",\"filebeat_version\":\"oss-2.2.0\",\"foo_dash_company_id\":\"{{ foo_iris_account }}\",\"foo_iris_account\":9900004,\"foogo_domain\":\"foogo.info\",\"foopro_domain\":\"foogo.info\",\"igw_id\":\"igw-1ad5e97f\",\"image\":\"ami-09150b5c79250e26c\",\"keypair\":\"bandwidithaws\",\"logger_elb\":\"qa.logger.foo.io\",\"metrics_agent\":\"collectd\",\"monit_environment\":\"qa\",\"nginx_config\":\"republicwireless-com\",\"peering_id\":\"pcx-52b2c651\",\"product\":\"foo\",\"provider\":\"ec2\",\"region\":\"us-east-1\",\"riak_http_port\":2067,\"scout_environment\":\"QA\",\"service_domain\":\"foo.io\",\"ssh_ca_fingerprint\":\"69f6617f461ceb6d05cd57a0b9e05ee27f555cfaaae5b52f6760bd9bcc976920\",\"ssh_ca_provisioner\":\"site-reliability-engineering@foopro.com\",\"ssh_ca_url\":\"https://api-qa.ca.phoneboothdev.info/\",\"subnet_qa_public_us_east_1a\":\"subnet-50f9cc09\",\"subnet_qa_public_us_east_1c\":\"subnet-2ded22b0\",\"subnet_qa_public_us_east_1d\":\"subnet-6fc0d512\",\"subnet_qa_public_us_east_1e\":\"subnet-c5f2a9ee\",\"subnet_qa_public_us_east_2a\":\"subnet-67fc020e\",\"subnet_qa_public_us_east_2b\":\"subnet-12151b60\",\"subnet_qa_public_us_east_2c\":\"subnet-f25c12b2\",\"subnet_qa_public_us_west_2a\":\"subnet-44b2ca55\",\"subnet_qa_public_us_west_2b\":\"subnet-1dabf772\",\"subnet_qa_public_us_west_2c\":\"subnet-d255c121\",\"telegraf_interval\":500,\"vm_mms_numbers\":[\"9199255565\"],\"vpc_id\":\"vpc-2f0fd9eb\"}"),
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
        hosts = jsonencode({
          host2 = {
            key2 = "value2"
          }
          host1 = {
            key = "value"
            key2 = "value2"
          }
        })
        vars = jsonencode({
          ansible_python_interpreter = "python"
          cert_name = "star-republicdev-info"
          cluster = "mob"
          cluster_env = "mob_qa"
          cluster_separator = "_"
          common_build_root_dir = "/home/{{ ansible_user }}"
          consul_cluster = "{{ region }}"
          credstash_table = "{{'credential-store_' + cluster + cluster_separator + env}}"
          dog_env = "qa"
          dog_env_group = "www-data"
          domain = "phoneboothdev.info"
          ec2_dns_suffix = "compute-1"
          elasticsearch_domain = "vpc-logger-qa-m52rd66iwl6df5gieiylyaso4e.us-east-1.es.amazonaws.com"
          email_alert_distro = []
          env = "qa"
          filebeat_version = "oss-2.2.0"
          foo_dash_company_id = "{{ foo_iris_account }}"
          foo_iris_account = 9900004
          foogo_domain = "foogo.info"
          foopro_domain = "foogo.info"
          igw_id = "igw-1ad5e97f"
          image = "ami-09150b5c79250e26c"
          keypair = "bandwidithaws"
          logger_elb = "qa.logger.foo.io"
          metrics_agent = "collectd"
          monit_environment = "qa"
          nginx_config = "republicwireless-com"
          peering_id = "pcx-52b2c651"
          product = "foo"
          provider = "ec2"
          region = "us-east-1"
          riak_http_port = 2067
          scout_environment = "QA"
          service_domain = "foo.io"
          ssh_ca_fingerprint = "69f6617f461ceb6d05cd57a0b9e05ee27f555cfaaae5b52f6760bd9bcc976920"
          ssh_ca_provisioner = "site-reliability-engineering@foopro.com"
          ssh_ca_url = "https://api-qa.ca.phoneboothdev.info/"
          subnet_qa_public_us_east_1a = "subnet-50f9cc09"
          subnet_qa_public_us_east_1c = "subnet-2ded22b0"
          subnet_qa_public_us_east_1d = "subnet-6fc0d512"
          subnet_qa_public_us_east_1e = "subnet-c5f2a9ee"
          subnet_qa_public_us_east_2a = "subnet-67fc020e"
          subnet_qa_public_us_east_2b = "subnet-12151b60"
          subnet_qa_public_us_east_2c = "subnet-f25c12b2"
          subnet_qa_public_us_west_2a = "subnet-44b2ca55"
          subnet_qa_public_us_west_2b = "subnet-1dabf772"
          subnet_qa_public_us_west_2c = "subnet-d255c121"
          telegraf_interval = 500
          vm_mms_numbers = ["9199255565"]
          vpc_id = "vpc-2f0fd9eb"
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
        hosts = jsonencode({
          host2 = {
            key2 = "value2"
          }
          host1 = {
            key = "value"
            key2 = "value2"
          }
        })
        vars = jsonencode({
		  test = "best"
          ansible_python_interpreter = "python"
          cert_name = "star-republicdev-info"
          cluster = "mob"
          cluster_env = "mob_qa"
          cluster_separator = "_"
          common_build_root_dir = "/home/{{ ansible_user }}"
          consul_cluster = "{{ region }}"
          credstash_table = "{{'credential-store_' + cluster + cluster_separator + env}}"
          dog_env = "qa"
          dog_env_group = "www-data"
          domain = "phoneboothdev.info"
          ec2_dns_suffix = "compute-1"
          elasticsearch_domain = "vpc-logger-qa-m52rd66iwl6df5gieiylyaso4e.us-east-1.es.amazonaws.com"
          email_alert_distro = []
          env = "qa"
          filebeat_version = "oss-2.2.0"
          foo_dash_company_id = "{{ foo_iris_account }}"
          foo_iris_account = 9900004
          foogo_domain = "foogo.info"
          foopro_domain = "foogo.info"
          igw_id = "igw-1ad5e97f"
          image = "ami-09150b5c79250e26c"
          keypair = "bandwidithaws"
          logger_elb = "qa.logger.foo.io"
          metrics_agent = "collectd"
          monit_environment = "qa"
          nginx_config = "republicwireless-com"
          peering_id = "pcx-52b2c651"
          product = "foo"
          provider = "ec2"
          region = "us-east-1"
          riak_http_port = 2067
          scout_environment = "QA"
          service_domain = "foo.io"
          ssh_ca_fingerprint = "69f6617f461ceb6d05cd57a0b9e05ee27f555cfaaae5b52f6760bd9bcc976920"
          ssh_ca_provisioner = "site-reliability-engineering@foopro.com"
          ssh_ca_url = "https://api-qa.ca.phoneboothdev.info/"
          subnet_qa_public_us_east_1a = "subnet-50f9cc09"
          subnet_qa_public_us_east_1c = "subnet-2ded22b0"
          subnet_qa_public_us_east_1d = "subnet-6fc0d512"
          subnet_qa_public_us_east_1e = "subnet-c5f2a9ee"
          subnet_qa_public_us_east_2a = "subnet-67fc020e"
          subnet_qa_public_us_east_2b = "subnet-12151b60"
          subnet_qa_public_us_east_2c = "subnet-f25c12b2"
          subnet_qa_public_us_west_2a = "subnet-44b2ca55"
          subnet_qa_public_us_west_2b = "subnet-1dabf772"
          subnet_qa_public_us_west_2c = "subnet-d255c121"
          telegraf_interval = 500
          vm_mms_numbers = ["9199255565"]
          vpc_id = "vpc-2f0fd9eb"
      })
    }
  }
}
`, resourceType, name)
}
func testAccDogFactConfig_big_update_remove(resourceType, name string) string {
	return fmt.Sprintf(`
resource %[1]q %[2]q {
  name = %[2]q 
    groups = {
      all = {
        children = ["dog_test"]
        hosts = jsonencode({
          host2 = {
            key2 = "value2"
          }
          host1 = {
            key = "value"
            key2 = "value2"
          }
        })
        vars = jsonencode({
          ansible_python_interpreter = "python"
          cert_name = "star-republicdev-info"
          cluster = "mob"
          cluster_env = "mob_qa"
          cluster_separator = "_"
          common_build_root_dir = "/home/{{ ansible_user }}"
          consul_cluster = "{{ region }}"
          credstash_table = "{{'credential-store_' + cluster + cluster_separator + env}}"
          dog_env = "qa"
          dog_env_group = "www-data"
          domain = "phoneboothdev.info"
          ec2_dns_suffix = "compute-1"
          elasticsearch_domain = "vpc-logger-qa-m52rd66iwl6df5gieiylyaso4e.us-east-1.es.amazonaws.com"
          email_alert_distro = []
          env = "qa"
          filebeat_version = "oss-2.2.0"
          foo_dash_company_id = "{{ foo_iris_account }}"
          foo_iris_account = 9900004
          foogo_domain = "foogo.info"
          foopro_domain = "foogo.info"
          igw_id = "igw-1ad5e97f"
          image = "ami-09150b5c79250e26c"
          keypair = "bandwidithaws"
          logger_elb = "qa.logger.foo.io"
          metrics_agent = "collectd"
          monit_environment = "qa"
          nginx_config = "republicwireless-com"
          peering_id = "pcx-52b2c651"
          product = "foo"
          provider = "ec2"
          region = "us-east-1"
          riak_http_port = 2067
          scout_environment = "QA"
          service_domain = "foo.io"
          ssh_ca_fingerprint = "69f6617f461ceb6d05cd57a0b9e05ee27f555cfaaae5b52f6760bd9bcc976920"
          ssh_ca_provisioner = "site-reliability-engineering@foopro.com"
          ssh_ca_url = "https://api-qa.ca.phoneboothdev.info/"
          subnet_qa_public_us_east_1a = "subnet-50f9cc09"
          subnet_qa_public_us_east_1c = "subnet-2ded22b0"
          subnet_qa_public_us_east_1d = "subnet-6fc0d512"
          subnet_qa_public_us_east_1e = "subnet-c5f2a9ee"
          subnet_qa_public_us_east_2a = "subnet-67fc020e"
          subnet_qa_public_us_east_2b = "subnet-12151b60"
          subnet_qa_public_us_east_2c = "subnet-f25c12b2"
          subnet_qa_public_us_west_2a = "subnet-44b2ca55"
          subnet_qa_public_us_west_2b = "subnet-1dabf772"
          subnet_qa_public_us_west_2c = "subnet-d255c121"
          telegraf_interval = 500
          vm_mms_numbers = ["9199255565"]
          vpc_id = "vpc-2f0fd9eb"
      })
    }
  }
}
`, resourceType, name)
}
