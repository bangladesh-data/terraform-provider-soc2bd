---
subcategory: "aws"
page_title: "AWS EC2 Deployment Guide"
description: |-
This document walks you through a basic deployment using Soc2bd's Terraform provider on AWS
---

# Deployment Guide

This deployment guide walks you through a basic AWS deployment of Soc2bd. For more information about Soc2bd, please reference the Soc2bd [documentation](https://docs.soc2bd.com/docs). It assumes basic knowledge of Soc2bd's service, the AWS Terraform provider, and a pre-existing AWS deployment in Terraform.

## Before you begin

- Sign up for an account on the [Soc2bd website](https://www.soc2bd.com).
- Create a Soc2bd [API key](https://docs.soc2bd.com/docs/api-overview). The key will need to have full permissions to Read, Write, & Provision, in order to deploy Connectors through Terraform.

## Setting up the Provider

First, we need to set up the Soc2bd Terraform provider by providing your network ID and the API key you provisioned earlier.

```terraform
provider "soc2bd" {
  api_token = "1234567890abcdef"
  network   = "autoco"
}

variable "network" {
  default = "autoco"
}
```

In general, we recommend that you use [environment variables](https://www.terraform.io/language/values/variables#environment-variables) to set sensitive variables such as the API key and mark such variables as [`sensitive`](https://www.terraform.io/language/values/variables#suppressing-values-in-cli-output).

## Creating the Remote Network and Connectors in Soc2bd

Next, we'll create the objects in Soc2bd that correspond to the AWS network that we're deploying Soc2bd into: a Remote Network to represent the AWS VPC, and a Connector to be deployed in that VPC. We'll use these objects when we're deploying the Connector image and creating Resources to access through Soc2bd.

```terraform
resource "soc2bd_remote_network" "aws_network" {
  name = "AWS Network"
}

resource "soc2bd_connector" "aws_connector" {
  remote_network_id = soc2bd_remote_network.aws_network.id
}

resource "soc2bd_connector_tokens" "aws_connector_tokens" {
  connector_id = soc2bd_connector.aws_connector.id
}
```

## Deploying the Connector

Now that we have the data types created in Soc2bd, we need to deploy a Connector into the AWS VPC to handle Soc2bd traffic. We'll use the pre-existing AWS AMI image for the Soc2bd Connector. First, we need to look up the latest AMI ID.

```terraform
data "aws_ami" "latest" {
  most_recent = true
  filter {
    name = "name"
    values = [
      "soc2bd/images/hvm-ssd/soc2bd-amd64-*",
    ]
  }
  owners = ["617935088040"]
}
```

Now, let's go ahead and deploy the AMI. For this example, we're creating a new VPC and security group, but you can use an existing one too. We'll deploy the Connector on a private subnet, because it doesn't need and shouldn't have a public IP address. Note the shell script that we use to configure the Connector tokens when the AMI launches.

```terraform
# define or use an existing VPC
module "demo_vpc" {
  source  = "terraform-aws-modules/vpc/aws"
  version = "2.64.0"

  name = "demo_vpc"
  cidr = "10.0.0.0/16"

  azs                            = ["us-east-1a"]
  private_subnets                = ["10.0.1.0/24"]
  public_subnets                 = ["10.0.2.0/24"]
  enable_classiclink_dns_support = true
  enable_dns_hostnames           = true
  enable_nat_gateway             = true

}

# If you use an existing Security group, the Connector requires egress traffic enabled but does not require ingress
module "demo_sg" {
  source  = "terraform-aws-modules/security-group/aws"
  version = "3.17.0"
  vpc_id  = module.demo_vpc.vpc_id
  name    = "demo_security_group"
  egress_cidr_blocks = ["0.0.0.0/0"]
  egress_rules = ["all-tcp", "all-udp", "all-icmp"]
}

# spin off a ec2 instance from Soc2bd AMI and configure tokens in user_data
module "ec2_tenant_connector" {
  source  = "terraform-aws-modules/ec2-instance/aws"
  version = "2.19.0"

  name = "demo_connector"
  user_data = <<-EOT
    #!/bin/bash
    set -e
    mkdir -p /etc/soc2bd/
    {
      echo SOC2BD_URL="https://${var.network}.soc2bd.com"
      echo SOC2BD_ACCESS_TOKEN="${soc2bd_connector_tokens.aws_connector_tokens.access_token}"
      echo SOC2BD_REFRESH_TOKEN="${soc2bd_connector_tokens.aws_connector_tokens.refresh_token}"
    } > /etc/soc2bd/connector.conf
    sudo systemctl enable --now soc2bd-connector
  EOT
  ami                    = data.aws_ami.latest.id
  instance_type          = "t3a.micro"
  vpc_security_group_ids = [module.demo_sg.this_security_group_id]
  subnet_id              = module.demo_vpc.private_subnets[0]
}
```

## Creating Resources

Now that you've deployed the Connector, we can create Resources on the same Remote Network that can be accessed through Soc2bd. For this example, we'll assume you already have an `aws_instance` defined. You'll need to define the Group ID explicitly, which you can pull from the [Admin API](https://docs.soc2bd.com/docs/api-overview).

```terraform
resource "soc2bd_resource" "tg_instance" {
  name = "My AWS Instance"
  address = aws_instance.my_instance.private_dns
  remote_network_id = soc2bd_remote_network.my_aws_network.id
  access {
    group_ids = ["R3JvdXG6OGky"]
  }
}
```
