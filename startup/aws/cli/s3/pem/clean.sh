#!/bin/bash

set -euo pipefail

ROOT_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")/../.." && pwd)"
cd "$ROOT_DIR"

source "$ROOT_DIR/lib/common.sh"
awscli_load_config "$ROOT_DIR"
awscli_load_state

if aws s3api head-bucket --bucket "${AUTHIFY_S3_PEM_BUCKET}" >/dev/null 2>&1; then
  aws s3 rm "s3://${AUTHIFY_S3_PEM_BUCKET}" --recursive >/dev/null 2>&1 || true
  aws s3 rb "s3://${AUTHIFY_S3_PEM_BUCKET}" --force >/dev/null 2>&1 || true
fi

awscli_remove_state AUTHIFY_S3_PEM_BUCKET

echo "✅ PEM bucket cleanup completed."
