#!/bin/bash

set -euo pipefail

ROOT_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")/.." && pwd)"
cd "$ROOT_DIR"

source "$ROOT_DIR/lib/common.sh"
awscli_load_config "$ROOT_DIR"
awscli_load_state

POLICIES=(
  "arn:aws:iam::aws:policy/AmazonEC2ContainerRegistryReadOnly"
  "arn:aws:iam::aws:policy/AmazonEC2FullAccess"
  "arn:aws:iam::aws:policy/AmazonS3ReadOnlyAccess"
  "arn:aws:iam::aws:policy/AmazonSSMReadOnlyAccess"
  "arn:aws:iam::aws:policy/AmazonSSMManagedInstanceCore"
)

for policy in "${POLICIES[@]}"; do
  aws iam detach-role-policy --role-name "${EC2_ROLE_NAME}" --policy-arn "${policy}" >/dev/null 2>&1 || true
done

aws iam remove-role-from-instance-profile \
  --instance-profile-name "${EC2_INSTANCE_PROFILE_NAME}" \
  --role-name "${EC2_ROLE_NAME}" >/dev/null 2>&1 || true
aws iam delete-instance-profile --instance-profile-name "${EC2_INSTANCE_PROFILE_NAME}" >/dev/null 2>&1 || true
aws iam delete-role --role-name "${EC2_ROLE_NAME}" >/dev/null 2>&1 || true

awscli_remove_state IAM_ROLE_ARN

echo "✅ IAM cleanup completed."
