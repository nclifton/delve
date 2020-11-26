# Creating a New EKS Cluster

To automatically deploy a new EKS cluster run:

```
./create_cluster.sh <CLUSTER_NAME>
```

E.g. to create "tp-staging" run:

```
./create_cluster.sh tp-staging
```

## Terraform

This directory is slowly being converted from using AWS Cli and eksctl to Terraform. As a result, there is a mix of Terraform modules and simple bash scripts. Over time, only the Terraform templates will remain.

The Terraform templates are called through the create_cluster script for now so no further steps are required to apply the Terraform templates when creating a cluster.

Terraform directory structure is explained at: https://www.terraform.io/docs/modules/index.html
