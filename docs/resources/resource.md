---
# generated by https://github.com/hashicorp/terraform-plugin-docs
page_title: "soc2bd_resource Resource - terraform-provider-soc2bd"
subcategory: ""
description: |-
  Resources in Soc2bd represent servers on the private network that clients can connect to. Resources can be defined by IP, CIDR range, FQDN, or DNS zone. For more information, see the Soc2bd documentation https://docs.soc2bd.com/docs/resources-and-access-nodes.
---

# soc2bd_resource (Resource)

Resources in Soc2bd represent servers on the private network that clients can connect to. Resources can be defined by IP, CIDR range, FQDN, or DNS zone. For more information, see the Soc2bd [documentation](https://docs.soc2bd.com/docs/resources-and-access-nodes).

## Example Usage

```terraform
provider "soc2bd" {
  api_token = "1234567890abcdef"
  network   = "mynetwork"
}

resource "soc2bd_remote_network" "aws_network" {
  name = "aws_remote_network"
}

resource "soc2bd_group" "aws" {
  name = "aws_group"
}

resource "soc2bd_service_account" "github_actions_prod" {
  name = "Github Actions PROD"
}

resource "soc2bd_resource" "resource" {
  name = "network"
  address = "internal.int"
  remote_network_id = soc2bd_remote_network.aws_network.id

  protocols {
    allow_icmp = true
    tcp  {
      policy = "RESTRICTED"
      ports = ["80", "82-83"]
    }
    udp {
      policy = "ALLOW_ALL"
    }
  }

  access {
    group_ids = [soc2bd_group.aws.id]
    service_account_ids = [soc2bd_service_account.github_actions_prod.id]
  }
}
```

<!-- schema generated by tfplugindocs -->

## Schema

### Required

- `address` (String) The Resource's IP/CIDR or FQDN/DNS zone
- `name` (String) The name of the Resource
- `remote_network_id` (String) Remote Network ID where the Resource lives

### Optional

- `access` (Block List, Max: 1) Restrict access to certain groups or service accounts (see [below for nested schema](#nestedblock--access))
- `alias` (String) Set a DNS alias address for the Resource. Must be a DNS-valid name string.
- `is_authoritative` (Boolean) Determines whether assignments in the access block will override any existing assignments. Default is `true`. If set to `false`, assignments made outside of Terraform will be ignored.
- `is_browser_shortcut_enabled` (Boolean) Controls whether an "Open in Browser" shortcut will be shown for this Resource in the Soc2bd Client.
- `is_visible` (Boolean) Controls whether this Resource will be visible in the main Resource list in the Soc2bd Client.
- `protocols` (Block List, Max: 1) Restrict access to certain protocols and ports. By default or when this argument is not defined, there is no restriction, and all protocols and ports are allowed. (see [below for nested schema](#nestedblock--protocols))

### Read-Only

- `id` (String) Autogenerated ID of the Resource, encoded in base64

<a id="nestedblock--access"></a>

### Nested Schema for `access`

Optional:

- `group_ids` (Set of String) List of Group IDs that will have permission to access the Resource.
- `service_account_ids` (Set of String) List of Service Account IDs that will have permission to access the Resource.

<a id="nestedblock--protocols"></a>

### Nested Schema for `protocols`

Required:

- `tcp` (Block List, Min: 1, Max: 1) (see [below for nested schema](#nestedblock--protocols--tcp))
- `udp` (Block List, Min: 1, Max: 1) (see [below for nested schema](#nestedblock--protocols--udp))

Optional:

- `allow_icmp` (Boolean) Whether to allow ICMP (ping) traffic

<a id="nestedblock--protocols--tcp"></a>

### Nested Schema for `protocols.tcp`

Required:

- `policy` (String) Whether to allow or deny all ports, or restrict protocol access within certain port ranges: Can be `RESTRICTED` (only listed ports are allowed), `ALLOW_ALL`, or `DENY_ALL`

Optional:

- `ports` (List of String) List of port ranges between 1 and 65535 inclusive, in the format `100-200` for a range, or `8080` for a single port

<a id="nestedblock--protocols--udp"></a>

### Nested Schema for `protocols.udp`

Required:

- `policy` (String) Whether to allow or deny all ports, or restrict protocol access within certain port ranges: Can be `RESTRICTED` (only listed ports are allowed), `ALLOW_ALL`, or `DENY_ALL`

Optional:

- `ports` (List of String) List of port ranges between 1 and 65535 inclusive, in the format `100-200` for a range, or `8080` for a single port

## Import

Import is supported using the following syntax:

```shell
terraform import soc2bd_resource.resource UmVzb3VyY2U6MzQwNDQ3
```