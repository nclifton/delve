# Creates peering connection between the cluster and Optus Proxy server

# Set the mtmo prod account as a provider 
# (default provider is set to the account the cluster is created in)
provider "aws" {
  alias   = "mtmo_prod"
  profile = var.mtmo_prod_aws_profile
  region  = "ap-southeast-2"
}

data "aws_caller_identity" "mtmo_prod" {
  provider = aws.mtmo_prod
}

# WARNING TODO: Optus Proxy is currently in the MTMO VPC and not in its own. 
# As a result the below connects to the MTMO VPC and not a standalone Optus VPC.
# Optus Proxy will be split out from MTMO VPC in the 2021
data "aws_vpc" "optus_proxy_vpc" {
  # Setting provider here to mtmo_prod so terraform connects to the mtmo production account
  provider = aws.mtmo_prod

  filter {
    name   = "tag:Name"
    values = ["mtmo"]
  }
}

resource "aws_vpc_peering_connection" "cluster_to_optus_proxy" {
  vpc_id        = data.aws_vpc.cluster_vpc.id
  peer_owner_id = data.aws_caller_identity.mtmo_prod.account_id
  peer_vpc_id   = data.aws_vpc.optus_proxy_vpc.id
  peer_region   = "ap-southeast-2"
  auto_accept   = false

  tags = {
    Name           = "${terraform.workspace} to optus-proxy"
    "managed-by"   = "terraform"
    "cluster-name" = "${terraform.workspace}"
  }
}

resource "aws_vpc_peering_connection_accepter" "optus_proxy" {
  # Setting provider here to mtmo_prod so terraform connects to the mtmo production account
  provider                  = aws.mtmo_prod
  vpc_peering_connection_id = aws_vpc_peering_connection.cluster_to_optus_proxy.id
  auto_accept               = true

  tags = {
    Name           = "${terraform.workspace} to optus-proxy"
    "managed-by"   = "terraform"
    "cluster-name" = "${terraform.workspace}"
  }
}
