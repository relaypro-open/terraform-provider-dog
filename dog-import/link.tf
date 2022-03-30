resource "dog_link" "d3" {
  address_handling = "union"
  connection = {
    api_port = 15672
    host = "dog-broker-dev3.nocell.io"
    password = "dog_trainer2"
    port = 5673
    ssl_options = {
        cacertfile = "/var/consul/data/pki/certs/ca.crt"
        certfile = "/var/consul/data/pki/certs/server.crt"
        fail_if_no_peer_cert = true
        keyfile = "/var/consul/data/pki/private/server.key"
        server_name_indication = "disable"
        verify = "verify_peer"
      },
    user = "dog_trainer"
    virtual_host = "dog"
  }
  connection_type = "thumper"
  direction = "bidirectional"
  enabled = true
  name = "d3"
}

