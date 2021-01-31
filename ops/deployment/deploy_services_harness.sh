#!/bin/bash
SCRIPTDIR=$(/usr/bin/dirname $0)
if [ "$SCRIPTDIR" != "." ]; then
  cd $SCRIPTDIR
fi
SCRIPTDIR=$(/bin/pwd)
ROOTDIR="$(cd ../../ >/dev/null 2>&1 && pwd)"

if ! [[ "${AWS_PROFILE}" =~ ^(sendsei-non-prod)$ ]]; then
  echo "Must set your AWS_PROFILE to the aws account Harness stores its secrets (sendsei-non-prod) even though this is mtmo-tp"
  exit 1
fi

ENV=$1
if ! [[ "${ENV}" =~ ^(tp-staging|tp-production)$ ]]; then
  echo "Environment name must be specified (tp-staging | tp-production)"
  echo "Usage: ./deploy_services.sh <ENV>"
  exit 1
fi

if [ "${ENV}" = "tp-production" ]; then
  echo "Are you sure you want to release TP v${RELEASE_VERSION} to PRODUCTION?"
  read -p "Type 'production' to proceed or any other character(s) to abort: "
  echo
  if [[ ! $REPLY = "production" ]]; then
    echo "Phew, that was a close one, exiting..."
    exit 1
  fi
fi

# Deploy using Harness by calling the Continuous Deployment Pipeline via the Harness API

ACCOUNT_ID="WLh9EopjQkCLAMp5z1pUVg"
APP_ID="F2xbNEw6SZGWbWqxTuj11w"
X_API_KEY=$(aws secretsmanager get-secret-value --secret-id "harness/HARNESS_X_API_KEY" | jq -r ".SecretString")

PIPELINE_NAME="continuous-deployment"
# Empty as DEPLOY_ALL_SERVICES flag ensures all services are deployed
SERVICES_TO_DEPLOY=()
DEPLOY_ALL_SERVICES="true"
DOCKER_TAG="latest"

# Blank as no ticket associated with a manual deployment
TICKET_URL=""

# Override defaults on pipeline so we can deploy to specified environment
# and not just staging.
INFRA_DEFINITION="${ENV}-infra"

# Get Pipeline ID
echo "Fetching Pipeline ID for ${PIPELINE_NAME}..."
pipelineIdQuery='{"query":"{pipelineByName(applicationId: \"'${APP_ID}'\", pipelineName: \"'${PIPELINE_NAME}'\") {id}}"}'
pipelineId=$(curl -L -X POST "https://app.harness.io/gateway/api/graphql?accountId=${ACCOUNT_ID}" -H "x-api-key: ${X_API_KEY}" -H "Content-Type: application/json" --data-raw "${pipelineIdQuery}" | jq -r ".data | .pipelineByName | .id")

echo "Pipeline Result: ${pipelineId}"
if [[ ${pipelineId} == *"error"* ]]; then
    exit 1
fi

# Fetch Execution Inputs
echo "Fetching Execution Inputs for pipeline..."
inputsQuery='{"query":"{executionInputs(applicationId: \"'${APP_ID}'\", entityId: \"'${pipelineId}'\", executionType: PIPELINE){serviceInputs{id name artifactType}}}"}'
executionInputs=$(curl -L -X POST "https://app.harness.io/gateway/api/graphql?accountId=${ACCOUNT_ID}" -H "x-api-key: ${X_API_KEY}" -H "Content-Type: application/json" --data-raw "${inputsQuery}" | jq -r ".data | .executionInputs")

echo "Execution Inputs Results: ${executionInputs}"
if [[ ${executionInputs} == *"error"* ]]; then
    exit 1
fi

# Next, we need to grab all artifact sources so the pipeline knows what image to deploy
echo "Fetching Service Inputs for pipeline..."
serviceInputs=()
for k in $(jq ".serviceInputs | keys | .[]" <<< "${executionInputs}"); do
    serviceName=$(jq -r ".serviceInputs[${k}] | .name" <<< "${executionInputs}")
    serviceId=$(jq -r ".serviceInputs[${k}] | .id" <<< "${executionInputs}")

    echo "Fetching ${serviceName} with image tag ${DOCKER_TAG}..."

    # Fetch Artifact related to Service
    artifactSourceQuery='{"query":"{service(serviceId: \"'${serviceId}'\"){artifactSources{name ...on ECRArtifactSource{name id}}}}"}'
    artifactSource=$(curl -L -X POST "https://app.harness.io/gateway/api/graphql?accountId=${ACCOUNT_ID}" -H "x-api-key: ${X_API_KEY}" -H "Content-Type: application/json" --data-raw "${artifactSourceQuery}" | jq -r ".data | .service | .artifactSources | .[0]")
    
    echo "Artifact Source Result: ${artifactSource}"
    if [[ ${artifactSource} == *"error"* ]]; then
        exit 1
    fi
    
    artifactSourceName=$(jq -r ".name" <<< "${artifactSource}")

    fmttd='{name: \"'${serviceName}'\", artifactValueInput: {valueType: BUILD_NUMBER, buildNumber: {buildNumber: \"'${DOCKER_TAG}'\", artifactSourceName: \"'${artifactSourceName}'\"}}}'
    serviceInputs+=("${fmttd}, ")
done

# Start Pipeline Execution
echo "Executing pipeline..."
serviceInputStr="${serviceInputs[@]}"
executionQuery='{"query":"mutation {startExecution(input: {applicationId: \"'${APP_ID}'\", entityId: \"'${pipelineId}'\", executionType: PIPELINE, variableInputs: [{name: \"Environment\" variableValue: {type: NAME value: \"'${ENV}'\"}},{name: \"InfraDefinition_KUBERNETES\" variableValue: {type: NAME value: \"'${INFRA_DEFINITION}'\"}},{name: \"servicesToDeploy\" variableValue: {type: NAME value: \"'${SERVICES_TO_DEPLOY[@]}'\"}},{name: \"deployAllServices\" variableValue: {type: NAME value: \"'${DEPLOY_ALL_SERVICES}'\"}},{name: \"dockerTag\" variableValue: {type: NAME value: \"'${DOCKER_TAG}'\"}},{name: \"ticketUrl\" variableValue: {type: NAME value: \"'${TICKET_URL}'\"}}], serviceInputs: ['${serviceInputStr%??}']}){clientMutationId execution{id status}}}"}'
result=$(curl -L -X POST "https://app.harness.io/gateway/api/graphql?accountId=${ACCOUNT_ID}" -H "x-api-key: ${X_API_KEY}" -H "Content-Type: application/json" --data-raw "${executionQuery}")

echo "Execution Result: ${result}"
if [[ ${result} == *"error"* ]]; then
    exit 1
fi

echo "Finished"
