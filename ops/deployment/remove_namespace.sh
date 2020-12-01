#!/bin/bash
SCRIPTDIR=`/usr/bin/dirname $0`
if [ "$SCRIPTDIR" != "." ]; then
	  cd $SCRIPTDIR
fi
SCRIPTDIR=`/bin/pwd`
ROOTDIR="$( cd ../../ >/dev/null 2>&1 && pwd )"

# Set environment name
ENV=$1
if [ -z "${ENV}" ]; then
	echo "Environment name must be specified"
	echo "Usage: ./remove_namespace.sh <ENV>"
	exit 1
fi

if [ ${ENV} = "production" ]; then
  echo "Are you sure you want to delete TP in PRODUCTION?"
  read -p "Type 'production' to proceed or any other character(s) to abort: "
  echo
  if [[ ! $REPLY = "production" ]]; then
    echo "Phew, that was a close one, exiting..."
    exit 1
  fi
  read -p "Are you really sure? (y/N) "
  echo
  if [[ ! $REPLY = "y" ]]; then
    echo "Phew, that was a close one, exiting..."
    exit 1
  fi
fi

NAMESPACE=$2
if [ -z "${NAMESPACE}" ]; then
	NAMESPACE=tp
fi

kube_config="${SCRIPTDIR}/${ENV}/connection_config.yaml"

# Must try and delete before force deleting to ensure resources are cleaned up
# If you force delete straight away then services and deployments aren't actually deleted
echo "Deleting ${NAMESPACE} namespace..."
kubectl --kubeconfig $kube_config delete ns ${NAMESPACE} &
sleep 10

echo "Checking deletion..."
deleted=false
for i in {1..3}
do
	s=$(($i * 20))
	echo "Waiting ${s} more seconds..."
	sleep $s

	exists=`kubectl --kubeconfig $kube_config get ns | grep ${NAMESPACE}`
	if [ -z "$exists" ]; then
		deleted=true
		break
	fi
done

# If the namespace deleted then exit
if [ "$deleted" = true ]; then
	echo "${NAMESPACE} namespace deleted successfully!"
	exit
fi

# If the namespace didn't delete then force delete as there was an issue
echo "Unable to delete ${NAMESPACE} namespace in time, force deleting..."

echo "Starting Kubernetes Proxy..."
PROXY_PID=`ps -ef | egrep "kubectl.*proxy.*port=8001" | egrep -v "grep" | awk '{print $2}'`

if [ ! -z "$PROXY_PID" ]; then
	echo "Proxy is already running, killing process ($PROXY_PID)"
    	kill $PROXY_PID
fi

kubectl --kubeconfig $kube_config proxy --port=8001 &
PROXY_PID=$!

echo "Waiting for proxy to start..."
sleep 5

echo "Updating config to remove finalizers..."
kubectl --kubeconfig $kube_config get ns ${NAMESPACE} -o json | jq '.spec.finalizers=[]' > ns-without-finalizers.json
curl -X PUT http://localhost:8001/api/v1/namespaces/${NAMESPACE}/finalize -H "Content-Type: application/json" --data-binary @ns-without-finalizers.json

echo "Deleting ${NAMESPACE} namespace..."
kubectl --kubeconfig $kube_config delete ns ${NAMESPACE}

echo "Cleaning up..."
kill $PROXY_PID
rm ns-without-finalizers.json

echo "Below table should no longer display ${NAMESPACE} namespace..."
kubectl --kubeconfig $kube_config get ns
