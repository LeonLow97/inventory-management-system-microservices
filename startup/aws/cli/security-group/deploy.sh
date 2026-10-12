#!/bin/bash

set -euo pipefail

ROOT_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")/.." && pwd)"
cd "$ROOT_DIR"

source "$ROOT_DIR/lib/common.sh"
awscli_load_config "$ROOT_DIR"
awscli_load_state

if [[ -z "${VPC_ID:-}" ]]; then
  echo "Missing VPC_ID. Run vpc/deploy.sh first."
  exit 1
fi

PUBLIC_IP="$(curl -sS https://checkip.amazonaws.com | tr -d '[:space:]')"

authorize_ingress_rule() {
  local group_id="$1"
  local port="$2"
  local source_type="$3"
  local source_value="$4"
  local permission

  case "$source_type" in
    cidr)
      permission="[{\"IpProtocol\":\"tcp\",\"FromPort\":${port},\"ToPort\":${port},\"IpRanges\":[{\"CidrIp\":\"${source_value}\"}]}]"
      ;;
    group)
      permission="[{\"IpProtocol\":\"tcp\",\"FromPort\":${port},\"ToPort\":${port},\"UserIdGroupPairs\":[{\"GroupId\":\"${source_value}\"}]}]"
      ;;
    *)
      echo "Unsupported ingress rule source type '$source_type'." >&2
      return 1
      ;;
  esac

  # Apply one rule at a time so a duplicate does not hide a missing rule.
  awscli_run_allowing_error_code "InvalidPermission.Duplicate" \
    aws ec2 authorize-security-group-ingress \
      --group-id "$group_id" \
      --ip-permissions "$permission" \
      --region "$REGION" >/dev/null
}

if aws ec2 describe-security-groups --group-ids "${BASTION_SG_ID:-missing}" --region "${REGION}" >/dev/null 2>&1; then
  echo "Bastion security group already exists, skipping creation."
else
  BASTION_SG_ID="$(aws ec2 describe-security-groups \
    --filters "Name=group-name,Values=${BASTION_SG_NAME}" "Name=vpc-id,Values=${VPC_ID}" \
    --query 'SecurityGroups[0].GroupId' \
    --output text)"
  [[ "$BASTION_SG_ID" == "None" || "$BASTION_SG_ID" == "null" ]] && BASTION_SG_ID=""
  if [[ -z "$BASTION_SG_ID" ]]; then
    BASTION_SG_ID="$(aws ec2 create-security-group \
      --group-name "${BASTION_SG_NAME}" \
      --description "Security group for Bastion hosts" \
      --vpc-id "${VPC_ID}" \
      --region "${REGION}" \
      --tag-specifications "ResourceType=security-group,Tags=[{Key=Name,Value=${BASTION_SG_NAME}}]" \
      --query 'GroupId' \
      --output text)"
  fi
fi

if aws ec2 describe-security-groups --group-ids "${ALB_SG_ID:-missing}" --region "${REGION}" >/dev/null 2>&1; then
  echo "ALB security group already exists, skipping creation."
else
  ALB_SG_ID="$(aws ec2 describe-security-groups \
    --filters "Name=group-name,Values=${ALB_SG_NAME}" "Name=vpc-id,Values=${VPC_ID}" \
    --query 'SecurityGroups[0].GroupId' \
    --output text)"
  [[ "$ALB_SG_ID" == "None" || "$ALB_SG_ID" == "null" ]] && ALB_SG_ID=""
  if [[ -z "$ALB_SG_ID" ]]; then
    ALB_SG_ID="$(aws ec2 create-security-group \
      --group-name "${ALB_SG_NAME}" \
      --description "Security group for Authify ALB" \
      --vpc-id "${VPC_ID}" \
      --region "${REGION}" \
      --tag-specifications "ResourceType=security-group,Tags=[{Key=Name,Value=${ALB_SG_NAME}}]" \
      --query 'GroupId' \
      --output text)"
  fi
fi

if aws ec2 describe-security-groups --group-ids "${API_GATEWAY_SG_ID:-missing}" --region "${REGION}" >/dev/null 2>&1; then
  echo "API Gateway security group already exists, skipping creation."
else
  API_GATEWAY_SG_ID="$(aws ec2 describe-security-groups \
    --filters "Name=group-name,Values=${API_GATEWAY_SG_NAME}" "Name=vpc-id,Values=${VPC_ID}" \
    --query 'SecurityGroups[0].GroupId' \
    --output text)"
  [[ "$API_GATEWAY_SG_ID" == "None" || "$API_GATEWAY_SG_ID" == "null" ]] && API_GATEWAY_SG_ID=""
  if [[ -z "$API_GATEWAY_SG_ID" ]]; then
    API_GATEWAY_SG_ID="$(aws ec2 create-security-group \
      --group-name "${API_GATEWAY_SG_NAME}" \
      --description "Security group for Authify API Gateway" \
      --vpc-id "${VPC_ID}" \
      --region "${REGION}" \
      --tag-specifications "ResourceType=security-group,Tags=[{Key=Name,Value=${API_GATEWAY_SG_NAME}}]" \
      --query 'GroupId' \
      --output text)"
  fi
fi

if aws ec2 describe-security-groups --group-ids "${AUTHIFY_SERVICE_SG_ID:-missing}" --region "${REGION}" >/dev/null 2>&1; then
  echo "Authify Service security group already exists, skipping creation."
else
  AUTHIFY_SERVICE_SG_ID="$(aws ec2 describe-security-groups \
    --filters "Name=group-name,Values=${AUTHIFY_SERVICE_SG_NAME}" "Name=vpc-id,Values=${VPC_ID}" \
    --query 'SecurityGroups[0].GroupId' \
    --output text)"
  [[ "$AUTHIFY_SERVICE_SG_ID" == "None" || "$AUTHIFY_SERVICE_SG_ID" == "null" ]] && AUTHIFY_SERVICE_SG_ID=""
  if [[ -z "$AUTHIFY_SERVICE_SG_ID" ]]; then
    AUTHIFY_SERVICE_SG_ID="$(aws ec2 create-security-group \
      --group-name "${AUTHIFY_SERVICE_SG_NAME}" \
      --description "Security group for Authify Service" \
      --vpc-id "${VPC_ID}" \
      --region "${REGION}" \
      --tag-specifications "ResourceType=security-group,Tags=[{Key=Name,Value=${AUTHIFY_SERVICE_SG_NAME}}]" \
      --query 'GroupId' \
      --output text)"
  fi
fi

if aws ec2 describe-security-groups --group-ids "${RDS_SG_ID:-missing}" --region "${REGION}" >/dev/null 2>&1; then
  echo "RDS security group already exists, skipping creation."
else
  RDS_SG_ID="$(aws ec2 describe-security-groups \
    --filters "Name=group-name,Values=${RDS_SG_NAME}" "Name=vpc-id,Values=${VPC_ID}" \
    --query 'SecurityGroups[0].GroupId' \
    --output text)"
  [[ "$RDS_SG_ID" == "None" || "$RDS_SG_ID" == "null" ]] && RDS_SG_ID=""
  if [[ -z "$RDS_SG_ID" ]]; then
    RDS_SG_ID="$(aws ec2 create-security-group \
      --group-name "${RDS_SG_NAME}" \
      --description "Security group for Authify RDS" \
      --vpc-id "${VPC_ID}" \
      --region "${REGION}" \
      --tag-specifications "ResourceType=security-group,Tags=[{Key=Name,Value=${RDS_SG_NAME}}]" \
      --query 'GroupId' \
      --output text)"
  fi
fi

if aws ec2 describe-security-groups --group-ids "${ELASTICACHE_SG_ID:-missing}" --region "${REGION}" >/dev/null 2>&1; then
  echo "ElastiCache security group already exists, skipping creation."
else
  ELASTICACHE_SG_ID="$(aws ec2 describe-security-groups \
    --filters "Name=group-name,Values=${ELASTICACHE_SG_NAME}" "Name=vpc-id,Values=${VPC_ID}" \
    --query 'SecurityGroups[0].GroupId' \
    --output text)"
  [[ "$ELASTICACHE_SG_ID" == "None" || "$ELASTICACHE_SG_ID" == "null" ]] && ELASTICACHE_SG_ID=""
  if [[ -z "$ELASTICACHE_SG_ID" ]]; then
    ELASTICACHE_SG_ID="$(aws ec2 create-security-group \
      --group-name "${ELASTICACHE_SG_NAME}" \
      --description "Security group for Authify ElastiCache" \
      --vpc-id "${VPC_ID}" \
      --region "${REGION}" \
      --tag-specifications "ResourceType=security-group,Tags=[{Key=Name,Value=${ELASTICACHE_SG_NAME}}]" \
      --query 'GroupId' \
      --output text)"
  fi
fi

authorize_ingress_rule "$BASTION_SG_ID" 22 cidr "${PUBLIC_IP}/32"
authorize_ingress_rule "$ALB_SG_ID" 80 cidr "0.0.0.0/0"
authorize_ingress_rule "$ALB_SG_ID" 443 cidr "0.0.0.0/0"
authorize_ingress_rule "$API_GATEWAY_SG_ID" "$PORT_API_GATEWAY" group "$ALB_SG_ID"
authorize_ingress_rule "$API_GATEWAY_SG_ID" 22 group "$BASTION_SG_ID"
authorize_ingress_rule "$API_GATEWAY_SG_ID" "$PORT_API_GATEWAY" group "$BASTION_SG_ID"
authorize_ingress_rule "$AUTHIFY_SERVICE_SG_ID" 22 group "$BASTION_SG_ID"
authorize_ingress_rule "$AUTHIFY_SERVICE_SG_ID" "$PORT_AUTHIFY_SERVICE" group "$API_GATEWAY_SG_ID"
authorize_ingress_rule "$AUTHIFY_SERVICE_SG_ID" "$PORT_AUTHIFY_SERVICE" group "$BASTION_SG_ID"
authorize_ingress_rule "$RDS_SG_ID" "$PORT_POSTGRESQL" group "$AUTHIFY_SERVICE_SG_ID"
authorize_ingress_rule "$RDS_SG_ID" "$PORT_POSTGRESQL" group "$BASTION_SG_ID"
authorize_ingress_rule "$ELASTICACHE_SG_ID" 6379 group "$AUTHIFY_SERVICE_SG_ID"
authorize_ingress_rule "$ELASTICACHE_SG_ID" 6379 group "$API_GATEWAY_SG_ID"
authorize_ingress_rule "$ELASTICACHE_SG_ID" 6379 group "$BASTION_SG_ID"

awscli_write_state \
  "PUBLIC_IP=$PUBLIC_IP" \
  "BASTION_SG_ID=$BASTION_SG_ID" \
  "ALB_SG_ID=$ALB_SG_ID" \
  "API_GATEWAY_SG_ID=$API_GATEWAY_SG_ID" \
  "AUTHIFY_SERVICE_SG_ID=$AUTHIFY_SERVICE_SG_ID" \
  "RDS_SG_ID=$RDS_SG_ID" \
  "ELASTICACHE_SG_ID=$ELASTICACHE_SG_ID"

echo "✅ Security groups ready."
