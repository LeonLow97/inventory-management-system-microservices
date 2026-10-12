#!/bin/bash

set -euo pipefail

ROOT_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")/.." && pwd)"
cd "$ROOT_DIR"

source "$ROOT_DIR/lib/common.sh"
awscli_load_config "$ROOT_DIR"
awscli_load_state

if aws elasticache describe-cache-clusters \
  --cache-cluster-id "${ELASTICACHE_CLUSTER_ID}" \
  --region "${REGION}" >/dev/null 2>&1; then
  aws elasticache delete-cache-cluster \
    --cache-cluster-id "${ELASTICACHE_CLUSTER_ID}" \
    --region "${REGION}" >/dev/null 2>&1 || true
  aws elasticache wait cache-cluster-deleted \
    --cache-cluster-id "${ELASTICACHE_CLUSTER_ID}" \
    --region "${REGION}"
fi

if aws elasticache describe-cache-subnet-groups \
  --cache-subnet-group-name "${ELASTICACHE_SUBNET_GROUP_NAME}" \
  --region "${REGION}" >/dev/null 2>&1; then
  aws elasticache delete-cache-subnet-group \
    --cache-subnet-group-name "${ELASTICACHE_SUBNET_GROUP_NAME}" \
    --region "${REGION}" >/dev/null
fi

awscli_remove_state ELASTICACHE_CLUSTER_ID ELASTICACHE_SUBNET_GROUP_NAME ELASTICACHE_ENDPOINT ELASTICACHE_PORT ELASTICACHE_PASSWORD

echo "✅ ElastiCache cleanup completed."
