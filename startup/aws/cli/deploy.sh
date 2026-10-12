#!/bin/bash

set -euo pipefail

ROOT_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)"
cd "$ROOT_DIR"

source "$ROOT_DIR/lib/common.sh"
awscli_load_config "$ROOT_DIR"
awscli_load_state

for executable in aws docker jq envsubst curl; do
  if ! command -v "$executable" >/dev/null 2>&1; then
    echo "Required command '$executable' was not found on PATH." >&2
    exit 1
  fi
done

aws sts get-caller-identity --region "$REGION" >/dev/null

require_non_empty() {
  local var_name="$1"
  if [[ -z "${!var_name:-}" ]]; then
    echo "Missing $var_name. Set it in $ROOT_DIR/.env before running deploy.sh."
    exit 1
  fi
}

require_non_empty MASTER_PASSWORD
require_non_empty AUTHIFY_USER_PASSWORD
require_non_empty AUTH_JWT_TOKEN_SECRET

VPC_SCRIPT="${ROOT_DIR}/vpc/deploy.sh"
IAM_SCRIPT="${ROOT_DIR}/iam/deploy.sh"
S3_SCHEMA_SCRIPT="${ROOT_DIR}/s3/schema/deploy.sh"
EC2_KEYS_SCRIPT="${ROOT_DIR}/ec2/keys/deploy.sh"
SG_SCRIPT="${ROOT_DIR}/security-group/deploy.sh"
RDS_SCRIPT="${ROOT_DIR}/rds/deploy.sh"
ECR_SCRIPT="${ROOT_DIR}/ecr/deploy.sh"
ELASTICACHE_SCRIPT="${ROOT_DIR}/elasticache/deploy.sh"
SSM_SCRIPT="${ROOT_DIR}/ssm/deploy.sh"
EC2_AUTHIFY_SERVICE_SCRIPT="${ROOT_DIR}/ec2/authify-service/deploy.sh"
EC2_API_GATEWAY_SCRIPT="${ROOT_DIR}/ec2/api-gateway/deploy.sh"
S3_PEM_SCRIPT="${ROOT_DIR}/s3/pem/deploy.sh"
EC2_BASTION_SCRIPT="${ROOT_DIR}/ec2/bastion/deploy.sh"
RDS_INITIALIZE_SCRIPT="${ROOT_DIR}/rds/initialize.sh"
ALB_SCRIPT="${ROOT_DIR}/alb/deploy.sh"

run_step() {
  local description="$1"
  local script="$2"

  printf '  %s\n' "$description"
  "$script"
}

deploy_foundation() {
  printf '\n==> Foundation\n'
  run_step "VPC and subnets" "$VPC_SCRIPT"
  run_step "IAM role and instance profile" "$IAM_SCRIPT"
  run_step "Database schema bucket" "$S3_SCHEMA_SCRIPT"
  run_step "Build and push application images" "$ECR_SCRIPT"
  run_step "Security groups" "$SG_SCRIPT"
}

deploy_data_services() {
  printf '\n==> Data services\n'
  run_step "RDS" "$RDS_SCRIPT"
  run_step "ElastiCache" "$ELASTICACHE_SCRIPT"
  run_step "SSM parameters" "$SSM_SCRIPT"
}

deploy_access() {
  printf '\n==> Bastion access\n'
  run_step "Create EC2 key pairs" "$EC2_KEYS_SCRIPT"
  run_step "Upload private keys to S3" "$S3_PEM_SCRIPT"
  run_step "Deploy bastion hosts" "$EC2_BASTION_SCRIPT"
}

initialize_database() {
  printf '\n==> Database initialization\n'
  run_step "Apply database user, grants, and schema through SSM" "$RDS_INITIALIZE_SCRIPT"
}

deploy_applications() {
  printf '\n==> Application services\n'
  run_step "Deploy Authify Service" "$EC2_AUTHIFY_SERVICE_SCRIPT"
  run_step "Deploy API Gateway" "$EC2_API_GATEWAY_SCRIPT"
}

deploy_edge() {
  printf '\n==> Public entry point\n'
  run_step "Deploy ALB and register application targets" "$ALB_SCRIPT"
}

deploy_foundation
deploy_data_services
deploy_access
initialize_database
deploy_applications
deploy_edge

echo "🎉 Infrastructure provisioning completed successfully"
