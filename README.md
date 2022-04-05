<p align="center">
  <img src="images/dog-segmented-green.network-400x400.png">
</p>

# dog Terraform Provider (Terraform Plugin Framework)

This provides the ability to manage the dog firewall management system via Terraform.

## Requirements

- [Terraform](https://www.terraform.io/downloads.html) >= 1.0
- [Go](https://golang.org/doc/install) >= 1.17
- [dog](https://relaypro-open.github.io/dog/) >= 1.3

## Building The Provider

1. Clone the repository
1. Enter the repository directory
1. Build the provider using the Go `install` command:

```shell
go install
```

## Using the provider

An example main.tf:
```
terraform {
  required_providers {
    dog = {
      source = "github.com/relaypro-open/dog"
    }
  }
}
  provider "dog" {
    api_key = "my-key"
    api_endpoint = "http://dog-server:7070/api/V2"
  }
```

## Developing the Provider

If you wish to work on the provider, you'll first need [Go](http://www.golang.org) installed on your machine (see [Requirements](#requirements) above).

To compile the provider, run `go install`. This will build the provider and put the provider binary in the `$GOPATH/bin` directory.

To generate or update documentation, run `go generate`.

In order to run the full suite of Acceptance tests, run `make testacc`.

*Note:* Acceptance tests create real resources, and often cost money to run.

```shell
make testacc
```
