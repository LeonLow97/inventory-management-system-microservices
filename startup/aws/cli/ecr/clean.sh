#!/bin/bash

set -euo pipefail

ROOT_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")/.." && pwd)"
cd "$ROOT_DIR"

source "$ROOT_DIR/lib/common.sh"
awscli_load_config "$ROOT_DIR"
awscli_load_state

docker rmi -f "${DOCKER_AUTHIFY_API_GATEWAY_IMAGE}" >/dev/null 2>&1 || true
docker rmi -f "${DOCKER_AUTHIFY_SERVICE_IMAGE}" >/dev/null 2>&1 || true

if aws ecr describe-repositories --repository-names "${ECR_AUTHIFY_REPOSITORY}" --region "${REGION}" >/dev/null 2>&1; then
  aws ecr delete-repository \
    --repository-name "${ECR_AUTHIFY_REPOSITORY}" \
    --region "${REGION}" \
    --force >/dev/null 2>&1 || true
fi

awscli_remove_state AWS_ACCOUNT_ID ECR_API_GATEWAY_IMAGE ECR_AUTHIFY_SERVICE_IMAGE DOCKER_AUTHIFY_API_GATEWAY_IMAGE_TAG DOCKER_AUTHIFY_SERVICE_IMAGE_TAG

echo "✅ ECR cleanup completed."
