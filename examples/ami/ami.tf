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
