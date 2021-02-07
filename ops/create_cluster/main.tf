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
  region  = var.aws_region
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
  postgres_endpoint = module.aurora_postgresql.endpoint
  env_dns           = var.env_dns
}

module "redis" {
  source  = "./modules/redis"
  env_dns = var.env_dns
}

module "rabbitmq" {
  source  = "./modules/rabbitmq"
  env_dns = var.env_dns
}

module "keda" {
  source = "./modules/keda"
}

module "harness" {
  # Only deploy to staging and production environments
  # https://stackoverflow.com/a/58193941
  count  = (length(regexall(".*staging.*", terraform.workspace)) > 0 || length(regexall(".*production.*", terraform.workspace)) > 0) ? 1 : 0
  source = "./modules/harness"
}

# TODO: Remove this and use peering connection once Optus Proxy has been
# re-architected. Due February 2021.
module "optus_vpn_connection" {
  source = "./modules/optus_vpn_connection"
}

module "optus_peering" {
  source                = "./modules/optus_peering"
  mtmo_prod_aws_profile = var.mtmo_prod_aws_profile
}
