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

if [ "${ENV}" = "production" ]; then
  echo "Are you sure you want to release TP v${RELEASE_VERSION} to PRODUCTION?"
  read -p "Type 'production' to proceed or any other character(s) to abort: "
  echo
  if [[ ! $REPLY = "production" ]]; then
    echo "Phew, that was a close one, exiting..."
    exit 1
  fi
fi

SLACK_NOTIFICATION_URL=https://hooks.slack.com/services/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
NOTIFICATION_ENVS=("qa")

if [[ " ${NOTIFICATION_ENVS[*]} " =~ ${ENV} ]]; then
  release_message="MTMO-TP deployment to ${ENV} environment started"
  curl -X POST -H 'Content-type: application/json' --data '{"text":"'"${release_message}"'"}' $SLACK_NOTIFICATION_URL
fi

kube_config="${SCRIPTDIR}/${ENV}/connection_config.yaml"

# Load in service paths based on whether a directory contains an infra folder or not
cd ${ROOTDIR}
service_paths=$(find ./ -type d -name infra | sed -e 's/\/infra$//' | sed -e 's/\.\///')
cd ${SCRIPTDIR}

# Create namespace if it doesn't exist
kubectl --kubeconfig ${kube_config} apply -f ./namespace.yaml 

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
    #/bin/bash ./helm_apply.sh ${chart_dir} ${ENV} "mtmostaging.com" "latest"
    echo "Environment managed by Harness, skipping service deployment..."
  elif [ ${ENV} = "qa" ]; then
    /bin/bash ./helm_apply.sh ${chart_dir} ${ENV} "qa.mtmostaging.com" "latest-release"
  elif [ ${ENV} = "production" ]; then
    /bin/bash ./helm_apply.sh ${chart_dir} ${ENV} "tp.mtmo.io" "latest-release"
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

if [[ " ${NOTIFICATION_ENVS[*]} " =~ ${ENV} ]]; then
  release_message="MTMO-TP deployment to ${ENV} environment complete!"
  curl -X POST -H "Content-type: application/json" --data "{\"text\":\"${release_message}\"}" $SLACK_NOTIFICATION_URL
fi
