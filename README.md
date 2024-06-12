# Terraform Provider DotEnv

[![Terraform](https://img.shields.io/badge/Terraform-844FBA.svg?style=for-the-badge&logo=Terraform&logoColor=white)](https://registry.terraform.io/providers/germanbrew/dotenv/latest)
[![OpenTofu](https://img.shields.io/badge/OpenTofu-FFDA18.svg?style=for-the-badge&logo=OpenTofu&logoColor=black)](https://github.com/opentofu/registry/blob/main/providers/g/germanbrew/dotenv.json)
[![GitHub Release](https://img.shields.io/github/v/release/germanbrew/terraform-provider-dotenv?sort=date&display_name=release&style=for-the-badge&logo=github&link=https%3A%2F%2Fgithub.com%2Fgermanbrew%2Fterraform-provider-dotenv%2Freleases%2Flatest)](https://github.com/germanbrew/terraform-provider-dotenv/releases/latest)
[![GitHub Actions Workflow Status](https://img.shields.io/github/actions/workflow/status/germanbrew/terraform-provider-dotenv/test.yaml?branch=main&style=for-the-badge&logo=github&label=Tests&link=https%3A%2F%2Fgithub.com%2Fgermanbrew%2Fterraform-provider-dotenv%2Factions%2Fworkflows%2Ftest.yaml)](https://github.com/germanbrew/terraform-provider-dotenv/actions/workflows/test.yaml)

A utility Terraform provider for dotfiles

## Requirements

- [Terraform](https://developer.hashicorp.com/terraform/downloads) >= 1.0

## Installing and Using this Plugin

You most likely want to download the provider from [Terraform Registry](https://registry.terraform.io/providers/germanbrew/dotenv/latest/docs).
The provider is also published in the [OpenTofu Registry](https://github.com/opentofu/registry/tree/main/providers/g/germanbrew).

Using Provider from Terraform Registry (TF >= 1.0)
This provider is published and available there. If you want to use it, just add the following to your terraform.tf:

```terraform
terraform {
  required_providers {
  dotenv = {
    source = "germanbrew/dotenv"
      version = "1.0.0"  # Adjust to latest version
    }
  }
  required_version = ">= 1.0"
}
```

Then run terraform init to download the provider.

## Development

### Requirements

- [Go](https://golang.org/) 1.21 (to build the provider plugin)
- [golangci-lint](https://github.com/golangci/golangci-lint) (to lint code)
- [terraform-plugin-docs](https://github.com/hashicorp/terraform-plugin-docs) (to generate registry documentation)

### Makefile Commands

Check the subcommands in our [Makefile](Makefile) for useful dev tools and scripts.

### Building The Provider

1. Clone the repository
1. Enter the repository directory
1. Build the provider using the Go `install` command:

```shell
go install
```

### Adding Dependencies

This provider uses [Go modules](https://github.com/golang/go/wiki/Modules).
Please see the Go documentation for the most up to date information about using Go modules.

To add a new dependency `github.com/author/dependency` to your Terraform provider:

```shell
go get github.com/author/dependency
go mod tidy
```

Then commit the changes to `go.mod` and `go.sum`.

### Developing the Provider

If you wish to work on the provider, you'll first need [Go](http://www.golang.org) installed on your machine (see [Requirements](#requirements) above).

To compile the provider, run `go install`. This will build the provider and put the provider binary in the `$GOPATH/bin` directory.

To generate or update documentation, run `go generate`.

In order to run the full suite of Acceptance tests, run `make testacc`.

*Note:* Acceptance tests create real resources, and often cost money to run.

```shell
make testacc
```

### Testing the provider locally

To test the provider locally:

1. Build the provider binary with `make build`
2. Create a new file `~/.terraform.rc` and point the provider to the absolute **directory** path of the binary file:
    ```hcl
    provider_installation {
      dev_overrides {
        "germanbrew/dotenv" = "/path/to/your/terraform-provider-dotenv/bin/"
      }
      direct {}
    }
    ```
3.  - Set the variable before running terraform commands:

```sh
TF_CLI_CONFIG_FILE=~/.terraform.rc terraform plan
```

    - Or set the env variable `TF_CLI_CONFIG_FILE` and point it to `~/.terraform.rc`: e.g.

```sh
export TF_CLI_CONFIG_FILE=~/.terraform.rc`
```

4. Now you can just use terraform normally. A warning will appear, that notifies you that you are using an provider override
    ```
    Warning: Provider development overrides are in effect
    ...
    ```
5. Unset the env variable if you don't want to use the local provider anymore:
    ```sh
    unset TF_CLI_CONFIG_FILE
    ```
