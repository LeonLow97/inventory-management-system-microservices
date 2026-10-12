#!/usr/bin/env bash

set -euo pipefail

SCRIPT_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)"
PROJECT_ROOT="$(cd "${SCRIPT_DIR}/../../../.." && pwd)"
TERRAFORM_DIR="${PROJECT_ROOT}/startup/aws/terraform/environments/ec2"
IMAGE_TAG="${1:?Pass the image tag used for this Terraform apply.}"

for executable in aws docker terraform; do
  if ! command -v "$executable" >/dev/null 2>&1; then
    echo "Required command '$executable' was not found on PATH." >&2
    exit 1
  fi
done

REGION="$(terraform -chdir="$TERRAFORM_DIR" output -raw aws_region)"
REPOSITORY_URL="$(terraform -chdir="$TERRAFORM_DIR" output -raw ecr_repository_url)"
REGISTRY_HOST="${REPOSITORY_URL%%/*}"
API_GATEWAY_IMAGE="${REPOSITORY_URL}:api-gateway-${IMAGE_TAG}"
AUTHIFY_SERVICE_IMAGE="${REPOSITORY_URL}:authify-service-${IMAGE_TAG}"

aws ecr get-login-password --region "$REGION" \
  | docker login --username AWS --password-stdin "$REGISTRY_HOST"

docker build --platform linux/amd64 \
  -f "${PROJECT_ROOT}/api-gateway/Dockerfile" \
  -t "$API_GATEWAY_IMAGE" \
  "$PROJECT_ROOT"
docker push "$API_GATEWAY_IMAGE"

docker build --platform linux/amd64 \
  -f "${PROJECT_ROOT}/authify-service/Dockerfile" \
  -t "$AUTHIFY_SERVICE_IMAGE" \
  "$PROJECT_ROOT"
docker push "$AUTHIFY_SERVICE_IMAGE"

echo "Pushed application images tagged ${IMAGE_TAG}."
