#!/bin/bash
SCRIPTDIR=$(/usr/bin/dirname "$0")
if [ "${SCRIPTDIR}" != "." ]; then
  cd "$SCRIPTDIR" || {
    echo "ERROR: Cannot cd to ${SCRIPTDIR} "
    exit
  }
fi

# Check required tools are installed and exit if not
command -v aws >/dev/null 2>&1 || {
  echo >&2 "Command 'aws' is not installed or is not executable, aborting."
  exit 1
}
command -v eksctl >/dev/null 2>&1 || {
  echo >&2 "Command 'eksctl' is not installed or is not executable, aborting."
  exit 1
}
command -v kubectl >/dev/null 2>&1 || {
  echo >&2 "Command 'kubectl' is not installed or is not executable, aborting."
  exit 1
}
command -v helm >/dev/null 2>&1 || {
  echo >&2 "Command 'helm' is not installed or is not executable, aborting."
  exit 1
}

command -v terraform >/dev/null 2>&1 || {
  echo >&2 "Command 'terraform' is not installed or is not executable, aborting."
  exit 1
}

# Set cluster name
CLUSTER_NAME=$1
VIPC_CIDR=$2
if [ -z "${CLUSTER_NAME}" ]; then
  echo "Cluster name must be specified"
  echo "Usage: ./create-cluster.sh <CLUSTER_NAME> [ VIPC_CIDR ]"
  exit 1
fi

if ! [[ "${AWS_PROFILE}" =~ ^(mtmo-non-prod|mtmo-prod)$ ]]; then
  echo "Must set your AWS_PROFILE to the target aws account currently either (mtmo-non-prod or mtmo-prod)"
  exit 1
fi

VERSION=1.17

# Create key pair for nodes to talk to EKS
aws ec2 delete-key-pair --key-name ${CLUSTER_NAME}
aws ec2 create-key-pair --key-name ${CLUSTER_NAME} | jq -r '.KeyMaterial' >${CLUSTER_NAME}.pem
chmod 400 ${CLUSTER_NAME}.pem
ssh-keygen -y -f ${CLUSTER_NAME}.pem >${CLUSTER_NAME}.pub

# this implicitly sets the kubectl context to use the correct connection config at the same time
echo "Creating cluster ${CLUSTER_NAME}..."
if [ -z "${VIPC_CIDR}" ]; then
  # Create cluster with default CIDR
  eksctl create cluster \
    --name ${CLUSTER_NAME} \
    --version ${VERSION} \
    --region ap-southeast-2 \
    --nodegroup-name sendsei \
    --node-volume-size 100 \
    --node-type m5.xlarge \
    --nodes 9 \
    --nodes-min 9 \
    --nodes-max 30 \
    --ssh-access \
    --ssh-public-key ${CLUSTER_NAME}.pub \
    --asg-access \
    --external-dns-access \
    --alb-ingress-access \
    --managed
else
  # Create cluster with passed CIDR
  eksctl create cluster \
    --name ${CLUSTER_NAME} \
    --version ${VERSION} \
    --region ap-southeast-2 \
    --nodegroup-name sendsei \
    --node-type m5.xlarge \
    --node-volume-size 100 \
    --nodes 9 \
    --nodes-min 9 \
    --nodes-max 30 \
    --ssh-access \
    --ssh-public-key ${CLUSTER_NAME}.pub \
    --asg-access \
    --external-dns-access \
    --alb-ingress-access \
    --vpc-cidr ${VIPC_CIDR} \
    --managed
fi

# Commit key pair to AWS Secret Manager
aws secretsmanager create-secret --name ${CLUSTER_NAME}-ssh-access.pub --description "SSH Public Key for ${CLUSTER_NAME} Cluster Nodes" --secret-string file://${CLUSTER_NAME}.pub
aws secretsmanager create-secret --name ${CLUSTER_NAME}-ssh-access.pem --description "SSH Access Key for ${CLUSTER_NAME} Cluster Nodes" --secret-string file://${CLUSTER_NAME}.pem
rm -f ${CLUSTER_NAME}.pub ${CLUSTER_NAME}.pem

# Wait for cluster to become active before trying to add/modify it
aws eks wait nodegroup-active --cluster-name ${CLUSTER_NAME} --nodegroup-name sendsei
aws eks wait cluster-active --name ${CLUSTER_NAME}

# Wait for cluster to become active before trying to add/modify it
aws eks wait cluster-active --name ${CLUSTER_NAME}

# Enable private vpc endpoint
echo "Enabling private VPC endpoint..."
aws eks update-cluster-config --name ${CLUSTER_NAME} --resources-vpc-config endpointPublicAccess=true,endpointPrivateAccess=true
aws eks wait cluster-active --name ${CLUSTER_NAME}

# Turn on cluster logging
# We shouldn't need this and it currently doesn't work anyways but need more investigation before removing
#eksctl utils update-cluster-logging --region=ap-southeast-2 --cluster=${CLUSTER_NAME} --enable-types=api,controllerManager --approve
#aws eks wait cluster-active --name ${CLUSTER_NAME}

# Allow eks-ops-user to update the cluster in the future
# (this allows anyone to deploy using the eks-ops-user creds instead of the user who created the cluster)
echo "Deploying eks-ops-user access..."
NODEGROUP_ROLES=$(eksctl get iamidentitymapping --cluster ${CLUSTER_NAME} | grep ${CLUSTER_NAME}-nodegroup | awk '{print $1}')
ACCOUNT_ID=$(aws sts get-caller-identity | jq -r '.Account')

cat <<EOF >./aws-auth-cm.yaml
apiVersion: v1
kind: ConfigMap
metadata:
  name: aws-auth
  namespace: kube-system
data:
  mapRoles: |
EOF

for role in ${NODEGROUP_ROLES}; do
  cat <<EOF >>./aws-auth-cm.yaml
    - rolearn: ${role}
      username: system:node:{{EC2PrivateDNSName}}
      groups:
        - system:bootstrappers
        - system:nodes
EOF
done

cat <<EOF >>./aws-auth-cm.yaml
  mapUsers: |
    - userarn: arn:aws:iam::${ACCOUNT_ID}:user/eks-ops-user
      username: eks-ops-user
      groups:
        - system:masters

EOF

kubectl apply -f ./aws-auth-cm.yaml
rm -f ./aws-auth-cm.yaml

# Deploy Terraform
/bin/bash ./terraform_apply.sh ${CLUSTER_NAME}
