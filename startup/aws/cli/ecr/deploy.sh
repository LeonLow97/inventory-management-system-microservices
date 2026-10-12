#!/bin/bash

set -euo pipefail

SCRIPT_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)"
ROOT_DIR="$(cd "$SCRIPT_DIR/.." && pwd)"
PROJECT_ROOT="$(cd "$ROOT_DIR/../.." && pwd)"
cd "$ROOT_DIR"

source "$ROOT_DIR/lib/common.sh"
awscli_load_config "$ROOT_DIR"
CONFIG_ECR_AUTHIFY_REPOSITORY="${ECR_AUTHIFY_REPOSITORY}"
CONFIG_DOCKER_AUTHIFY_API_GATEWAY_IMAGE_TAG="${DOCKER_AUTHIFY_API_GATEWAY_IMAGE_TAG}"
CONFIG_DOCKER_AUTHIFY_SERVICE_IMAGE_TAG="${DOCKER_AUTHIFY_SERVICE_IMAGE_TAG}"
awscli_load_state

ECR_AUTHIFY_REPOSITORY="${CONFIG_ECR_AUTHIFY_REPOSITORY}"
DOCKER_AUTHIFY_API_GATEWAY_IMAGE_TAG="${CONFIG_DOCKER_AUTHIFY_API_GATEWAY_IMAGE_TAG}"
DOCKER_AUTHIFY_SERVICE_IMAGE_TAG="${CONFIG_DOCKER_AUTHIFY_SERVICE_IMAGE_TAG}"

if [[ -z "${AWS_ACCOUNT_ID:-}" ]]; then
  AWS_ACCOUNT_ID="$(aws sts get-caller-identity --query 'Account' --output text)"
fi

if aws ecr describe-repositories --repository-names "${ECR_AUTHIFY_REPOSITORY}" --region "${REGION}" >/dev/null 2>&1; then
  echo "ECR repository $ECR_AUTHIFY_REPOSITORY already exists, skipping creation."
else
  echo "Creating ECR repository $ECR_AUTHIFY_REPOSITORY ..."
  aws ecr create-repository \
    --repository-name "${ECR_AUTHIFY_REPOSITORY}" \
    --region "${REGION}" >/dev/null
fi

aws ecr get-login-password --region "${REGION}" \
  | docker login --username AWS --password-stdin "${AWS_ACCOUNT_ID}.dkr.ecr.${REGION}.amazonaws.com"

ECR_API_GATEWAY_IMAGE="${AWS_ACCOUNT_ID}.dkr.ecr.${REGION}.amazonaws.com/${ECR_AUTHIFY_REPOSITORY}:${DOCKER_AUTHIFY_API_GATEWAY_IMAGE_TAG}"
ECR_AUTHIFY_SERVICE_IMAGE="${AWS_ACCOUNT_ID}.dkr.ecr.${REGION}.amazonaws.com/${ECR_AUTHIFY_REPOSITORY}:${DOCKER_AUTHIFY_SERVICE_IMAGE_TAG}"

cd "${PROJECT_ROOT}"
docker build -f api-gateway/Dockerfile --platform linux/amd64 -t "${DOCKER_AUTHIFY_API_GATEWAY_IMAGE}" .
docker tag "${DOCKER_AUTHIFY_API_GATEWAY_IMAGE}" "${ECR_API_GATEWAY_IMAGE}"
docker push "${ECR_API_GATEWAY_IMAGE}"

cd "${PROJECT_ROOT}"
docker build -f authify-service/Dockerfile --platform linux/amd64 -t "${DOCKER_AUTHIFY_SERVICE_IMAGE}" .
docker tag "${DOCKER_AUTHIFY_SERVICE_IMAGE}" "${ECR_AUTHIFY_SERVICE_IMAGE}"
docker push "${ECR_AUTHIFY_SERVICE_IMAGE}"

awscli_write_state \
  "AWS_ACCOUNT_ID=$AWS_ACCOUNT_ID" \
  "ECR_API_GATEWAY_IMAGE=$ECR_API_GATEWAY_IMAGE" \
  "ECR_AUTHIFY_SERVICE_IMAGE=$ECR_AUTHIFY_SERVICE_IMAGE"

echo "✅ ECR images pushed."
