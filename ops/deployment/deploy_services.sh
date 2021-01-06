#!/bin/bash
SCRIPTDIR=$(/usr/bin/dirname $0)
if [ "$SCRIPTDIR" != "." ]; then
  cd $SCRIPTDIR
fi
SCRIPTDIR=$(/bin/pwd)
ROOTDIR="$(cd ../../ >/dev/null 2>&1 && pwd)"

ENV=$1
if [ -z "${ENV}" ]; then
  echo "Environment name must be specified"
  echo "Usage: ./deploy_services.sh <ENV>"
  exit 1
fi

#RELEASE_VERSION=$(<${SCRIPTDIR}/../releases/RELEASE_VERSION)

if [ "${ENV}" = "production" ]; then
  echo "Are you sure you want to release TP v${RELEASE_VERSION} to PRODUCTION?"
  read -p "Type 'production' to proceed or any other character(s) to abort: "
  echo
  if [[ ! $REPLY = "production" ]]; then
    echo "Phew, that was a close one, exiting..."
    exit 1
  fi
fi

kube_config="${SCRIPTDIR}/${ENV}/connection_config.yaml"
alb_controller_file="${SCRIPTDIR}/${ENV}/alb-ingress-controller.yaml"
ingress_file="${SCRIPTDIR}/${ENV}/ingress.yaml"

# Load in service paths based on whether a directory contains an infra folder or not
cd ${ROOTDIR}
service_paths=$(find ./ -type d -name infra | sed -e 's/\/infra$//' | sed -e 's/\.\///')
cd ${SCRIPTDIR}

echo "Creating namespace..."
kubectl --kubeconfig ${kube_config} create namespace tp
# TODO: Annotation will be templated in Terraform once we move to using Terraform for deploying
# This annotation allows linkerd to auto inject proxies into each pod in the namespace
kubectl --kubeconfig ${kube_config} annotate namespace tp linkerd.io/inject=enabled
# These annotations allow tracing to be auto injected into all our services
kubectl --kubeconfig ${kube_config} annotate namespace tp config.linkerd.io/trace-collector=linkerd-collector.linkerd:55678
kubectl --kubeconfig ${kube_config} annotate namespace tp config.alpha.linkerd.io/trace-collector-service-account=linkerd-collector

echo "Deploying services..."
# Keep service names so we can delete unknown helm charts later
service_names=()

for service_path in ${service_paths}; do
  # Hack to get helm charts working, should host and use a helm chart repository
  infra="${ROOTDIR}/${service_path}/infra"
  # Find chart dir without knowing name but knowing there is another dir called values in infra
  chart_dir=$(find ${infra} -mindepth 1 -maxdepth 1 -type d | egrep -v "values")
  service_name=$(basename ${chart_dir})
  service_names+=($service_name)

  if [ ${ENV} = "staging" ]; then
    /bin/bash ./helm_apply.sh ${chart_dir} ${ENV} "mtmostaging.com" "latest"
    #echo "Environment managed by Harness, skipping Helm Deploy..."
  elif [ ${ENV} = "production" ]; then
    /bin/bash ./helm_apply.sh ${chart_dir} ${ENV} "tp.mtmo.io" "latest"
  else
    echo "Environment not supported!"
  fi
done

# Delete unknown Helm Charts from TP namespace
echo "Deleting old Helm Charts..."
helm_charts=$(helm --kubeconfig ${kube_config} list --short --namespace tp)

for helm_chart in ${helm_charts}; do
  # If helm chart isn't in service names array then uninstall
  if [[ ! " ${service_names[@]} " =~ " ${helm_chart} " ]]; then
    helm --kubeconfig ${kube_config} uninstall ${helm_chart} --namespace tp
  fi
done

echo "Creating ingress..."
kubectl --kubeconfig ${kube_config} apply -f ${alb_controller_file}
kubectl --kubeconfig ${kube_config} apply -f ${ingress_file}
