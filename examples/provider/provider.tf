terraform {
  required_providers {
    dog = {
      source = "github.com/relaypro-open/dog"
    }
  }
}

provider "dog" {
    api_key = "my-key"
    api_endpoint = "http://dog-ubuntu-server.lxd:7070/api/V2"
}

