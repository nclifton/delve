resource "aws_db_subnet_group" "dbsubnet" {
  name       = "${terraform.workspace}-dbsubnet"
  subnet_ids = data.aws_subnet_ids.privatesubnets.ids

  tags = {
    "managed-by"   = "terraform"
    "cluster-name" = "${terraform.workspace}"
    Name           = "${terraform.workspace}-dbsubnet"
  }
}

resource "aws_rds_cluster_instance" "cluster_instances" {
  count                = 3
  identifier           = "${terraform.workspace}-${count.index}"
  cluster_identifier   = aws_rds_cluster.postgresql.id
  db_subnet_group_name = aws_db_subnet_group.dbsubnet.id
  instance_class       = "db.t3.medium"
  engine               = "aurora-postgresql"
  engine_version       = "11.7"
  tags = {
    "managed-by"   = "terraform"
    "cluster-name" = "${terraform.workspace}"
  }
}

resource "aws_rds_cluster" "postgresql" {
  cluster_identifier        = terraform.workspace
  engine                    = "aurora-postgresql"
  engine_version            = "11.7"
  availability_zones        = ["ap-southeast-2b", "ap-southeast-2a", "ap-southeast-2c"]
  database_name             = "sendsei"
  master_username           = "foo"
  master_password           = "barbut8chars"
  skip_final_snapshot       = true
  final_snapshot_identifier = "sendseifinalsnapshot"
  db_subnet_group_name      = aws_db_subnet_group.dbsubnet.id
  vpc_security_group_ids    = [data.aws_security_group.clustersg.id]
  tags = {
    "managed-by"   = "terraform"
    "cluster-name" = "${terraform.workspace}"
  }
}
