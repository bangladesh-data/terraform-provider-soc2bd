---
# generated by https://github.com/hashicorp/terraform-plugin-docs
page_title: "soc2bd_resource Data Source - terraform-provider-soc2bd"
subcategory: ""
description: |-
  Resources in Soc2bd represent any network destination address that you wish to provide private access to for users authorized via the Soc2bd Client application. Resources can be defined by either IP or DNS address, and all private DNS addresses will be automatically resolved with no client configuration changes. For more information, see the Soc2bd documentation https://docs.soc2bd.com/docs/resources-and-access-nodes.
---

# soc2bd_resource (Data Source)

Resources in Soc2bd represent any network destination address that you wish to provide private access to for users authorized via the Soc2bd Client application. Resources can be defined by either IP or DNS address, and all private DNS addresses will be automatically resolved with no client configuration changes. For more information, see the Soc2bd [documentation](https://docs.soc2bd.com/docs/resources-and-access-nodes).

## Example Usage

```terraform
data "soc2bd_resource" "foo" {
  id = "<your resource's id>"
}
```

<!-- schema generated by tfplugindocs -->

## Schema

### Required

- `id` (String) The ID of the Resource. The ID for the Resource can be obtained from the Admin API or the URL string in the Admin Console.

### Read-Only

- `address` (String) The Resource's address, which may be an IP address, CIDR range, or DNS address
- `name` (String) The name of the Resource
- `protocols` (Block List) By default (when this argument is not defined) no restriction is applied, and all protocols and ports are allowed. (see [below for nested schema](#nestedblock--protocols))
- `remote_network_id` (String) The Remote Network ID that the Resource is associated with. Resources may only be associated with a single Remote Network.

<a id="nestedblock--protocols"></a>

### Nested Schema for `protocols`

Read-Only:

- `allow_icmp` (Boolean) Whether to allow ICMP (ping) traffic
- `tcp` (Block List) (see [below for nested schema](#nestedblock--protocols--tcp))
- `udp` (Block List) (see [below for nested schema](#nestedblock--protocols--udp))

<a id="nestedblock--protocols--tcp"></a>

### Nested Schema for `protocols.tcp`

Read-Only:

- `policy` (String) Whether to allow or deny all ports, or restrict protocol access within certain port ranges: Can be `RESTRICTED` (only listed ports are allowed), `ALLOW_ALL`, or `DENY_ALL`
- `ports` (List of String) List of port ranges between 1 and 65535 inclusive, in the format `100-200` for a range, or `8080` for a single port

<a id="nestedblock--protocols--udp"></a>

### Nested Schema for `protocols.udp`

Read-Only:

- `policy` (String) Whether to allow or deny all ports, or restrict protocol access within certain port ranges: Can be `RESTRICTED` (only listed ports are allowed), `ALLOW_ALL`, or `DENY_ALL`
- `ports` (List of String) List of port ranges between 1 and 65535 inclusive, in the format `100-200` for a range, or `8080` for a single port
