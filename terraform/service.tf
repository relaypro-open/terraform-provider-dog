resource "dog_service" "https" {
  version =         1
  name =          "new_https_name" 
  services = [
    {
    ports = ["443"]
    protocol = "tcp"
    }
  ]
}

#resource "dog_service" "https" {
#  version =         2
#  name =          "update_https_name" 
#  services = [
#  {
#    ports = ["8443"]
#    protocol = "tcp"
#  }
#  ]
#}
