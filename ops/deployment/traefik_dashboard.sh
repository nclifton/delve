#!/bin/bash
SCRIPTDIR=`/usr/bin/dirname $0`
if [ "$SCRIPTDIR" != "." ]; then
  cd $SCRIPTDIR
fi

# Set environment name
ENV=$1
if [ -z "${ENV}" ]; then
	echo "Environment name must be specified"
	echo "Usage: ./launch_dashboard.sh <ENV>"
	exit 1
fi

kube_config="${SCRIPTDIR}/${ENV}/connection_config.yaml"

RUNNING_PID=`ps -ef | egrep "kubectl.*port-forward.*9000:9000" | egrep -v "grep" | awk '{print $2}'`

if [ ! -z "${RUNNING_PID}" ]; then
  echo "Traefik Port Forward is already running, killing process (${RUNNING_PID})"
  kill ${RUNNING_PID}
fi

kubectl --kubeconfig ${kube_config} port-forward $(kubectl --kubeconfig ${kube_config} get pods --selector "app.kubernetes.io/name=traefik" -n traefik --output=name) 9000:9000 -n traefik &

echo "Please access http://localhost:9000/dashboard/ via browser"
