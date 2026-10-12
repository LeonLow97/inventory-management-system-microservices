#!/bin/bash

set -euo pipefail

ROOT_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")/../.." && pwd)"
cd "$ROOT_DIR"

source "$ROOT_DIR/lib/common.sh"
awscli_load_config "$ROOT_DIR"
awscli_load_state

if ! aws ec2 describe-instances --instance-ids "${AUTHIFY_SERVICE_INSTANCE_ID_AZ1:-missing}" >/dev/null 2>&1; then
  AUTHIFY_SERVICE_INSTANCE_ID_AZ1="$(aws ec2 describe-instances --filters "Name=tag:Name,Values=${AUTHIFY_SERVICE_INSTANCE_NAME_AZ1}" "Name=instance-state-name,Values=pending,running,stopping,stopped" --query 'Reservations[0].Instances[0].InstanceId' --output text 2>/dev/null || true)"
  [[ "$AUTHIFY_SERVICE_INSTANCE_ID_AZ1" == "None" || "$AUTHIFY_SERVICE_INSTANCE_ID_AZ1" == "null" ]] && AUTHIFY_SERVICE_INSTANCE_ID_AZ1=""
fi

if ! aws ec2 describe-instances --instance-ids "${AUTHIFY_SERVICE_INSTANCE_ID_AZ2:-missing}" >/dev/null 2>&1; then
  AUTHIFY_SERVICE_INSTANCE_ID_AZ2="$(aws ec2 describe-instances --filters "Name=tag:Name,Values=${AUTHIFY_SERVICE_INSTANCE_NAME_AZ2}" "Name=instance-state-name,Values=pending,running,stopping,stopped" --query 'Reservations[0].Instances[0].InstanceId' --output text 2>/dev/null || true)"
  [[ "$AUTHIFY_SERVICE_INSTANCE_ID_AZ2" == "None" || "$AUTHIFY_SERVICE_INSTANCE_ID_AZ2" == "null" ]] && AUTHIFY_SERVICE_INSTANCE_ID_AZ2=""
fi

if [[ -n "${AUTHIFY_SERVICE_INSTANCE_ID_AZ1:-}" ]]; then
  aws ec2 terminate-instances --instance-ids "$AUTHIFY_SERVICE_INSTANCE_ID_AZ1" >/dev/null 2>&1 || true
  aws ec2 wait instance-terminated --instance-ids "$AUTHIFY_SERVICE_INSTANCE_ID_AZ1" >/dev/null 2>&1 || true
fi

if [[ -n "${AUTHIFY_SERVICE_INSTANCE_ID_AZ2:-}" ]]; then
  aws ec2 terminate-instances --instance-ids "$AUTHIFY_SERVICE_INSTANCE_ID_AZ2" >/dev/null 2>&1 || true
  aws ec2 wait instance-terminated --instance-ids "$AUTHIFY_SERVICE_INSTANCE_ID_AZ2" >/dev/null 2>&1 || true
fi

aws ec2 delete-key-pair --key-name "${EC2_AUTHIFY_SERVICE_KEY_NAME_AZ1}" >/dev/null 2>&1 || true
aws ec2 delete-key-pair --key-name "${EC2_AUTHIFY_SERVICE_KEY_NAME_AZ2}" >/dev/null 2>&1 || true
rm -f "$ROOT_DIR/${EC2_AUTHIFY_SERVICE_KEY_NAME_AZ1}.pem" "$ROOT_DIR/${EC2_AUTHIFY_SERVICE_KEY_NAME_AZ2}.pem"

awscli_remove_state AUTHIFY_SERVICE_INSTANCE_ID_AZ1 AUTHIFY_SERVICE_INSTANCE_ID_AZ2 AUTHIFY_SERVICE_PRIVATE_IP_AZ1 AUTHIFY_SERVICE_PRIVATE_IP_AZ2

echo "✅ Authify Service cleanup completed."
