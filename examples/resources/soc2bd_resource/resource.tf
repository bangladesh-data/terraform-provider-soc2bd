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

