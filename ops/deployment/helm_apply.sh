#!/bin/bash
SCRIPTDIR=$(/usr/bin/dirname $0)
if [ "$SCRIPTDIR" != "." ]; then
  cd $SCRIPTDIR
fi
SCRIPTDIR=$(/bin/pwd)
ROOTDIR="$(cd ../../ >/dev/null 2>&1 && pwd)"

CHART_DIR=$1
if [ -z "${CHART_DIR}" ]; then
  echo "Chart directory must be specified"
  echo "Usage: ./helm_apply.sh <CHART_DIR> <ENV> <ENV_DNS> <IMAGE_TAG>"
  exit 1
fi

ENV=$2
if [ -z "${ENV}" ]; then
  echo "Environment name must be specified"
  echo "Usage: ./helm_apply.sh <CHART_DIR> <ENV> <ENV_DNS> <IMAGE_TAG>"
  exit 1
fi

ENV_DNS=$3
if [ -z "${ENV_DNS}" ]; then
  echo "Environment DNS must be specified"
  echo "Usage: ./helm_apply.sh <CHART_DIR> <ENV> <ENV_DNS> <IMAGE_TAG>"
  exit 1
fi

IMAGE_TAG=$4
if [ -z "${IMAGE_TAG}" ]; then
  echo "Image tag must be specified"
  echo "Usage: ./helm_apply.sh <CHART_DIR> <ENV> <ENV_DNS> <IMAGE_TAG>"
  exit 1
fi

kube_config="${SCRIPTDIR}/${ENV}/connection_config.yaml"
service_name=$(basename ${CHART_DIR})
env_values_file="${CHART_DIR}/../values/${ENV}-values.yaml"

# This script is intended to mimic the functionality in Harness for our scripted deployments.
# Helm doesn't support local environment variables but Harness has the ability to replace env vars in values.yaml.
# This script replaces the most common env vars used between services.
# A list of built in variables this script can be used to replace can be found at: https://docs.harness.io/article/aza65y4af6-built-in-variables-list

# Create the temporary values file
cp ${CHART_DIR}/values.yaml ${CHART_DIR}/.generated-values.yaml

# Set Environment Name (pre-fixed by tp as a workaround for cluster names being prefixed with tp)
sed -i -e 's/${env.name}/tp-'"${ENV}"'/g' ${CHART_DIR}/.generated-values.yaml

# Set Environment URL
sed -i -e 's/${environmentVariable.dns}/'"${ENV_DNS}"'/g' ${CHART_DIR}/.generated-values.yaml

# Set Create Namespace
sed -i -e 's/${environmentVariable.createNamespace}/false/g' ${CHART_DIR}/.generated-values.yaml

# Set Image Tag
sed -i -e 's/${artifact.buildNo}/'"${IMAGE_TAG}"'/g' ${CHART_DIR}/.generated-values.yaml

# Harness does a rolling deploy but our manual scripts are currently uninstall/install
helm --kubeconfig ${kube_config} uninstall ${service_name} --namespace tp

# Install chart with generated values file and env specific values if exists
if [ -f ${env_values_file} ]; then
  helm --kubeconfig ${kube_config} install ${service_name} ${CHART_DIR} -f ${CHART_DIR}/.generated-values.yaml -f ${env_values_file} --namespace tp
else
  helm --kubeconfig ${kube_config} install ${service_name} ${CHART_DIR} -f ${CHART_DIR}/.generated-values.yaml --namespace tp
fi

rm -f "${CHART_DIR}/.generated-values.yaml"
