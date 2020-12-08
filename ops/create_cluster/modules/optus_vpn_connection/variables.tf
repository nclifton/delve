# VPC has been created by eksctl so grab already defined resources
# (may be module variables in the future)

data "aws_vpc" "cluster_vpc" {
  filter {
    name   = "tag:Name"
    values = ["eksctl-${terraform.workspace}-cluster/VPC"]
  }
}
