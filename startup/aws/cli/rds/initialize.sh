#!/bin/bash

set -euo pipefail

ROOT_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")/.." && pwd)"
cd "$ROOT_DIR"

source "$ROOT_DIR/lib/common.sh"
awscli_load_config "$ROOT_DIR"
awscli_load_state

if [[ -z "${BASTION_INSTANCE_ID_AZ1:-}" || -z "${AUTHIFY_S3_SCHEMA_BUCKET:-}" ]]; then
  echo "Missing the AZ1 bastion instance or schema bucket. Run the required deploy steps first." >&2
  exit 1
fi

awscli_run_ssm_script \
  "$BASTION_INSTANCE_ID_AZ1" \
  "s3://${AUTHIFY_S3_SCHEMA_BUCKET}/rds/initialize-schema.sh" \
  "Initialize Authify database"

echo "✅ Authify database is initialized."
