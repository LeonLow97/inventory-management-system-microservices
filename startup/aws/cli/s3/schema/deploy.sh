#!/bin/bash

set -euo pipefail

ROOT_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")/../.." && pwd)"
cd "$ROOT_DIR"

source "$ROOT_DIR/lib/common.sh"
awscli_load_config "$ROOT_DIR"
awscli_load_state

if [[ "${AUTHIFY_SERVICE_SCHEMA_DIR}" = /* ]]; then
  SCHEMA_SOURCE="${AUTHIFY_SERVICE_SCHEMA_DIR}"
else
  SCHEMA_SOURCE="${ROOT_DIR}/${AUTHIFY_SERVICE_SCHEMA_DIR}"
fi

if aws s3api head-bucket --bucket "${AUTHIFY_S3_SCHEMA_BUCKET}" >/dev/null 2>&1; then
  echo "Schema bucket $AUTHIFY_S3_SCHEMA_BUCKET already exists, skipping creation."
else
  echo "Creating schema bucket $AUTHIFY_S3_SCHEMA_BUCKET ..."
  if [[ "${REGION}" == "us-east-1" ]]; then
    aws s3api create-bucket --bucket "${AUTHIFY_S3_SCHEMA_BUCKET}" >/dev/null
  else
    aws s3api create-bucket \
      --bucket "${AUTHIFY_S3_SCHEMA_BUCKET}" \
      --region "${REGION}" \
      --create-bucket-configuration "LocationConstraint=${REGION}" >/dev/null
  fi
fi

aws s3 cp "${SCHEMA_SOURCE}" "s3://${AUTHIFY_S3_SCHEMA_BUCKET}/users.sql" --region "$REGION" >/dev/null
aws s3 cp "$ROOT_DIR/rds/initialize-schema.sh" "s3://${AUTHIFY_S3_SCHEMA_BUCKET}/rds/initialize-schema.sh" --region "$REGION" >/dev/null

awscli_write_state \
  "AUTHIFY_S3_SCHEMA_BUCKET=$AUTHIFY_S3_SCHEMA_BUCKET"

echo "✅ Database schema and initializer uploaded to S3."
