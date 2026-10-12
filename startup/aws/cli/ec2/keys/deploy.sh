#!/bin/bash

set -euo pipefail

ROOT_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")/../.." && pwd)"
cd "$ROOT_DIR"

source "$ROOT_DIR/lib/common.sh"
awscli_load_config "$ROOT_DIR"
awscli_load_state

ensure_key_pair() {
  local key_name="$1"
  local key_path="$ROOT_DIR/${key_name}.pem"

  if aws ec2 describe-key-pairs --key-names "$key_name" >/dev/null 2>&1; then
    [[ -f "$key_path" ]] || { echo "Missing $key_path"; exit 1; }
  else
    echo "Creating key pair $key_name ..."
    aws ec2 create-key-pair --key-name "$key_name" --query 'KeyMaterial' --output text > "$key_path"
  fi

  chmod 400 "$key_path"
}

ensure_key_pair "$EC2_API_GATEWAY_KEY_NAME_AZ1"
ensure_key_pair "$EC2_API_GATEWAY_KEY_NAME_AZ2"
ensure_key_pair "$EC2_AUTHIFY_SERVICE_KEY_NAME_AZ1"
ensure_key_pair "$EC2_AUTHIFY_SERVICE_KEY_NAME_AZ2"

echo "✅ API Gateway and Authify Service key pairs ready."
