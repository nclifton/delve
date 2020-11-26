# VPC has been created by eksctl so grab already defined resources
# (may be module variables in the future)

data "aws_vpc" "selected" {
  filter {
    name   = "tag:Name"
    values = ["eksctl-${terraform.workspace}-cluster/VPC"]
  }
}

data "aws_subnet_ids" "privatesubnets" {
  vpc_id = data.aws_vpc.selected.id

  filter {
    name   = "tag:Name"
    values = ["eksctl-${terraform.workspace}-cluster/SubnetPrivate*"]
  }
}

data "aws_security_group" "clustersg" {
  vpc_id = data.aws_vpc.selected.id

  filter {
    name   = "tag:Name"
    values = ["eks-cluster-sg-${terraform.workspace}-*"]
  }
}
