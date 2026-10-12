#!/usr/bin/env bash

set -euo pipefail

SCRIPT_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)"
PROJECT_ROOT="$(cd "${SCRIPT_DIR}/../../../.." && pwd)"
TERRAFORM_DIR="${PROJECT_ROOT}/startup/aws/terraform/environments/ec2"
API_BASE_URL="$(terraform -chdir="$TERRAFORM_DIR" output -raw api_base_url)"
TARGET_GROUP_ARN="$(terraform -chdir="$TERRAFORM_DIR" output -raw api_gateway_target_group_arn)"
REGION="$(terraform -chdir="$TERRAFORM_DIR" output -raw aws_region)"
ATTEMPTS="${READINESS_ATTEMPTS:-120}"
INTERVAL_SECONDS="${READINESS_INTERVAL_SECONDS:-15}"

for attempt in $(seq 1 "$ATTEMPTS"); do
  HEALTHY_TARGETS="$(aws elbv2 describe-target-health \
    --target-group-arn "$TARGET_GROUP_ARN" \
    --region "$REGION" \
    --query 'TargetHealthDescriptions[?TargetHealth.State==`healthy`].Target.Id' \
    --output text | wc -w | tr -d '[:space:]')"
  HTTP_STATUS="$(curl --connect-timeout 5 --max-time 10 -sS -o /dev/null -w '%{http_code}' \
    "${API_BASE_URL}/health/readiness" 2>/dev/null || true)"

  echo "Readiness attempt ${attempt}/${ATTEMPTS}: healthy_targets=${HEALTHY_TARGETS}, http_status=${HTTP_STATUS:-curl_failed}"
  if [[ "$HEALTHY_TARGETS" == "2" && "$HTTP_STATUS" == "200" ]]; then
    echo "Authify is ready at ${API_BASE_URL}."
    exit 0
  fi

  sleep "$INTERVAL_SECONDS"
done

echo "Authify did not become ready at ${API_BASE_URL}. Inspect EC2 cloud-init logs through Systems Manager." >&2
exit 1
