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
echo
if [[ "${CONTEXT}" != *"${CLUSTER_NAME}"* ]]; then
    echo "Context does not contain environment name, aborting."
    kubectl config get-contexts
    echo "Select the correct context above and change to it using: kubectl config use-context <CONTEXT_NAME>"
    exit 1
fi

ENV_DNS="mtmostaging.com"
if [[ "${CLUSTER_NAME}" == "tp-qa" ]]; then
    ENV_DNS="qa.mtmostaging.com"
elif [[ "${CLUSTER_NAME}" == "tp-sre" ]]; then
    ENV_DNS="sre.mtmostaging.com"
elif [[ "${CLUSTER_NAME}" == "tp-production" ]]; then
    ENV_DNS="tp.mtmo.io"
fi

# Make sure no previous local state is used
rm -rf .terraform || {
  echo "ERROR: Cannot rm terraform cache"
  exit 1
}

# Load remote state and cluster workspace
terraform init -backend-config="bucket=${AWS_PROFILE}-tfstate-bucket"

terraform workspace select "${CLUSTER_NAME}"
if [ $? == 1 ]; then
  echo "Workspace ${CLUSTER_NAME} does not exist!"
  exit 1
fi

terraform destroy -var "aws_profile=${AWS_PROFILE}" -var "env_dns=${ENV_DNS}"
