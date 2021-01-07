#!/bin/bash
SCRIPTDIR=$(/usr/bin/dirname "$0")
if [ "$SCRIPTDIR" != "." ]; then
  cd "$SCRIPTDIR" || {
    echo "ERROR: Cannot cd to ${SCRIPTDIR} "
    exit
  }
fi

CLUSTER_NAME=$1
if [ -z "${CLUSTER_NAME}" ]; then
  echo "Cluster name must be specified"
  exit 1
fi

if ! [[ "${AWS_PROFILE}" =~ ^(mtmo-non-prod|mtmo-prod)$ ]]; then
  echo "Must set your AWS_PROFILE to the target aws account currently either (mtmo-non-prod or mtmo-prod)"
  exit 1
fi

CONTEXT=$(kubectl config current-context)
echo "Current kubectl cluster context is: ${CONTEXT}"
read -p "Is this the correct context for ${CLUSTER_NAME}? (y/N) "
echo
if [[ ! $REPLY =~ ^[Yy]$ ]]; then
    kubectl config get-contexts
    echo "Select the correct context above and change to it using: kubectl config use-context <CONTEXT_NAME>"
    exit 1
fi

ENV_DNS="mtmostaging.com"
if [[ "${CLUSTER_NAME}" == "tp-production" ]]; then
    ENV_DNS="tp.mtmo.io"
fi

# Make sure no previous local state is used
rm -rf .terraform || {
  echo "ERROR: Cannot rm terraform cache"
  exit
}

# Load remote state and cluster workspace
terraform init -backend-config="bucket=${AWS_PROFILE}-tfstate-bucket"

terraform workspace select "${CLUSTER_NAME}"
if [ $? == 1 ]; then
  terraform workspace new "${CLUSTER_NAME}"
fi

terraform apply -var "aws_profile=${AWS_PROFILE}" -var "env_dns=${ENV_DNS}"

# Get the cluster_endpoint and update the service-config file
ENDPOINT=$(terraform output postgresql_endpoint)
echo "******* POSTGRES CLUSTER ENDPOINT: ${ENDPOINT} **********"
