# Soc2bd Terraform Provider

[![Coverage Status](https://coveralls.io/repos/github/bangladesh-data/terraform-provider-soc2bd/badge.svg?branch=main&t=rqgifB)](https://coveralls.io/github/bangladesh-data/terraform-provider-soc2bd?branch=main)

## Requirements

- Bash
- [Go](https://golang.org/doc/install) 1.19 (to build the provider plugin)
- [Terraform](https://www.terraform.io/downloads.html) 1.x

## Build

Run the following command to build the provider

```shell
make build
```

## Test

Run unit tests:

```shell
make test
```

To run acceptance tests against a real Soc2bd network you first need to define the following 3 environment variables:

```shell
export SOC2BD_URL=soc2bd.com
export SOC2BD_NETWORK=<your network slug - <slug>.soc2bd.com>
export SOC2BD_API_TOKEN=<API token with Read, Write & Provision permissions>
```

Then you can run by:

```shell
make testacc
```

## Install

Install the provider for local testing.

```shell
make install
```

## Documentation

To update the documentation edit the files in `templates/` and then run `make docs`. The files in `docs/` are auto-generated and should not be updated manually.

## Contributions

Contributions to this project are [released](https://help.github.com/articles/github-terms-of-service/#6-contributions-under-repository-license) under the [project's open source license](LICENSE).
