#!/bin/bash

set -euo pipefail

ROOT_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")/.." && pwd)"
cd "$ROOT_DIR"

source "$ROOT_DIR/lib/common.sh"
awscli_load_config "$ROOT_DIR"
awscli_load_state

if aws rds describe-db-instances --db-instance-identifier "${RDS_INSTANCE_IDENTIFIER}" >/dev/null 2>&1; then
  aws rds delete-db-instance \
    --db-instance-identifier "${RDS_INSTANCE_IDENTIFIER}" \
    --skip-final-snapshot \
    --region "${REGION}" >/dev/null 2>&1 || true
  aws rds wait db-instance-deleted \
    --db-instance-identifier "${RDS_INSTANCE_IDENTIFIER}" \
    --region "${REGION}"
fi

if aws rds describe-db-subnet-groups --db-subnet-group-name "${DB_SUBNET_GROUP_NAME}" >/dev/null 2>&1; then
  aws rds delete-db-subnet-group \
    --db-subnet-group-name "${DB_SUBNET_GROUP_NAME}" \
    --region "${REGION}" >/dev/null
fi

if aws rds describe-db-parameter-groups --db-parameter-group-name "${DB_PARAMETER_GROUP_NAME}" >/dev/null 2>&1; then
  aws rds delete-db-parameter-group \
    --db-parameter-group-name "${DB_PARAMETER_GROUP_NAME}" \
    --region "${REGION}" >/dev/null
fi

awscli_remove_state DB_SUBNET_GROUP_NAME DB_PARAMETER_GROUP_NAME RDS_INSTANCE_IDENTIFIER RDS_ENDPOINT MASTER_USERNAME DB_NAME AUTHIFY_USER

echo "✅ RDS cleanup completed."
