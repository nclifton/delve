# Creates peering connection between the cluster and Optus Proxy server

# Create VPN Connection
resource "aws_vpn_gateway" "vpn_gateway" {
  vpc_id = data.aws_vpc.cluster_vpc.id

  tags = {
    Name           = "optus-vpn-${terraform.workspace}"
    "managed-by"   = "terraform"
    "cluster-name" = terraform.workspace
  }
}

# VPN Master (Optus Mascot) - 52.62.164.168
#resource "aws_customer_gateway" "customer_gateway_mascot" {
#  bgp_asn    = 65000
#  ip_address = "52.62.164.168"
#  type       = "ipsec.1"

#  tags = {
#    Name = "optus-vpn-mascot-${terraform.workspace}"
#    "managed-by" = "terraform"
#    "cluster-name" = terraform.workspace
#  }
#}

# VPN Slave (Optus Hume) - 3.106.16.225
#resource "aws_customer_gateway" "customer_gateway_hume" {
#  bgp_asn    = 65000
#  ip_address = "3.106.16.225"
#  type       = "ipsec.1"

#  tags = {
#    Name = "optus-vpn-hume-${terraform.workspace}"
#    "managed-by" = "terraform"
#    "cluster-name" = terraform.workspace
#  }
#}

data "aws_customer_gateway" "customer_gateway_mascot" {
  filter {
    name   = "tag:Name"
    values = ["optus-vpn-mascot"]
  }
}

data "aws_customer_gateway" "customer_gateway_hume" {
  filter {
    name   = "tag:Name"
    values = ["optus-vpn-hume"]
  }
}

resource "aws_vpn_connection" "main_mascot" {
  vpn_gateway_id      = aws_vpn_gateway.vpn_gateway.id
  customer_gateway_id = data.aws_customer_gateway.customer_gateway_mascot.id
  type                = "ipsec.1"
  static_routes_only  = true

  tags = {
    Name           = "optus-vpn-mascot-${terraform.workspace}"
    "managed-by"   = "terraform"
    "cluster-name" = terraform.workspace
  }
}

resource "aws_vpn_connection" "main_hume" {
  vpn_gateway_id      = aws_vpn_gateway.vpn_gateway.id
  customer_gateway_id = data.aws_customer_gateway.customer_gateway_hume.id
  type                = "ipsec.1"
  static_routes_only  = true

  tags = {
    Name           = "optus-vpn-hume-${terraform.workspace}"
    "managed-by"   = "terraform"
    "cluster-name" = terraform.workspace
  }
}

resource "aws_vpn_connection_route" "optus_mascot" {
  destination_cidr_block = "210.49.127.33/32"
  vpn_connection_id      = aws_vpn_connection.main_mascot.id
}

resource "aws_vpn_connection_route" "optus_hume" {
  destination_cidr_block = "210.49.127.126/32"
  vpn_connection_id      = aws_vpn_connection.main_hume.id
}

# Update VPC Routing to use VPN Connection
data "aws_route_tables" "rts" {
  filter {
    name   = "tag:Name"
    values = ["eksctl-${terraform.workspace}-cluster*"]
  }
}

# VPN Master (Optus Mascot) - 210.49.127.33
resource "aws_route" "vpn_mascot" {
  count = length(data.aws_route_tables.rts.ids)

  route_table_id = sort(data.aws_route_tables.rts.ids)[count.index]

  destination_cidr_block = "210.49.127.33/32"
  gateway_id             = aws_vpn_connection.main_mascot.vpn_gateway_id

  depends_on = [aws_vpn_connection.main_mascot]
}

# VPN Slave (Optus Hume) - 210.49.127.126
resource "aws_route" "vpn_hume" {
  count = length(data.aws_route_tables.rts.ids)

  route_table_id = sort(data.aws_route_tables.rts.ids)[count.index]

  destination_cidr_block = "210.49.127.126/32"
  gateway_id             = aws_vpn_connection.main_hume.vpn_gateway_id

  depends_on = [aws_vpn_connection.main_hume]
}
