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

TOKEN=`kubectl --kubeconfig ${kube_config} -n kube-system describe secret $(kubectl --kubeconfig ${kube_config} -n kube-system get secret | grep eks-admin | awk '{print $1}') | grep "token\:" | awk '{print $2}'`

RUNNING_PID=`ps -ef | egrep "kubectl.*proxy.*port=8001" | egrep -v "grep" | awk '{print $2}'`

if [ ! -z "${RUNNING_PID}" ]; then
  echo "Proxy is already running, killing process (${RUNNING_PID})"
  kill ${RUNNING_PID}
fi

kubectl proxy --kubeconfig ${kube_config} --port=8001 &

echo "Please access http://localhost:8001/api/v1/namespaces/kubernetes-dashboard/services/https:kubernetes-dashboard:https/proxy/#/login via browser, then input token:"
echo "$TOKEN"
