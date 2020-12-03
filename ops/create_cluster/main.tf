terraform {
  required_version = "~> 0.13"

  backend "s3" {
    key    = "terraform.tfstate"
    region = "ap-southeast-2"
  }

  required_providers {
    aws = {
      source  = "hashicorp/aws"
      version = "~> 2.70"
    }
    helm = {
      source  = "hashicorp/helm"
      version = "~> 1.2"
    }
    kubernetes = {
      source  = "hashicorp/kubernetes"
      version = "~> 1.11"
    }
    tls = {
      source  = "hashicorp/tls"
      version = "~> 2.1"
    }
    null = {
      source  = "hashicorp/null"
      version = "~> 2.1"
    }
    random = {
      source  = "hashicorp/random"
      version = "~> 2.3"
    }
  }
}

provider "aws" {
  profile = var.aws_profile
  region  = "ap-southeast-2"
}

# Call each child module

module "kubernetes_dashboard" {
  source = "./modules/kubernetes_dashboard"
}

module "metrics_server" {
  source = "./modules/metrics_server"
}

module "linkerd" {
  source = "./modules/linkerd"
}

module "alb_ingress" {
  source = "./modules/alb_ingress"
}

module "external_dns" {
  source = "./modules/external_dns"
}

module "cluster_autoscaler" {
  source = "./modules/cluster_autoscaler"
}

module "newrelic" {
  source = "./modules/newrelic"
}

module "aurora_postgresql" {
  source = "./modules/aurora_postgresql"
}

module "pgadmin" {
  source            = "./modules/pgadmin"
  postgres_endpoint = "${module.aurora_postgresql.endpoint}"
}

module "redis" {
  source = "./modules/redis"
}

module "rabbitmq" {
  source = "./modules/rabbitmq"
}

module "keda" {
  source = "./modules/keda"
}

module "harness" {
    source = "./modules/harness"
}

module "optus_peering" {
    source = "./modules/optus_peering"
    mtmo_prod_aws_profile = var.mtmo_prod_aws_profile
}
