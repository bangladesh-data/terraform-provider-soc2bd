---
layout: ""
page_title: "Soc2bd Provider"
description: |-
  The Soc2bd Terraform provider can be used to manage your Soc2bd private network in Terraform
---

# Soc2bd Provider

The Soc2bd provider can be used with [Soc2bd](https://www.soc2bd.com) to manage your Remote Networks, Connectors, and Resources through Terraform.

~> **Warning** Using the Soc2bd Terraform provider will cause any secrets, such as Connector tokens or Soc2bd API keys, that are managed in Terraform to be persisted in both Terraform's state file and in any generated plan files. If malicious attackers obtain these credentials, they could intercept network traffic intended for your private network or cause a denial of service event. For any Terraform module that reads or writes Soc2bd secrets, these files should be treated as sensitive and protected accordingly.

## Soc2bd Setup

You need an API key to use Soc2bd's Terraform provider. See our [documentation](https://docs.soc2bd.com/docs/api-overview) for more details about creating an API key. You will also need your network ID, or the prefix of your Soc2bd URL that you use to sign into the Admin Console. For example, if your URL is `autoco.soc2bd.com` your network ID is `autoco`.

## Guidance and documentation

Visit our [documentation](https://docs.soc2bd.com/docs) for more information on configuring and using Soc2bd.

## Example Usage

```terraform
provider "soc2bd" {
  api_token = "1234567890abcdef"
  network   = "autoco"
}
```

<!-- schema generated by tfplugindocs -->

## Schema

### Optional

- `api_token` (String, Sensitive) The access key for API operations. You can retrieve this
  from the Soc2bd Admin Console ([documentation](https://docs.soc2bd.com/docs/api-overview)).
  Alternatively, this can be specified using the SOC2BD_API_TOKEN environment variable.
- `http_max_retry` (Number) Specifies a retry limit for the http requests made. The default value is 10.
  Alternatively, this can be specified using the SOC2BD_HTTP_MAX_RETRY environment variable
- `http_timeout` (Number) Specifies a time limit in seconds for the http requests made. The default value is 10 seconds.
  Alternatively, this can be specified using the SOC2BD_HTTP_TIMEOUT environment variable
- `network` (String) Your Soc2bd network ID for API operations.
  You can find it in the Admin Console URL, for example:
  `autoco.soc2bd.com`, where `autoco` is your network ID
  Alternatively, this can be specified using the SOC2BD_NETWORK environment variable.
- `url` (String) The default is 'soc2bd.com'
  This is optional and shouldn't be changed under normal circumstances.
