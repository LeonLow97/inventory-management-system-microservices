#!/bin/bash

set -euo pipefail

ROOT_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)"
cd "$ROOT_DIR"

source "$ROOT_DIR/lib/common.sh"
awscli_load_config "$ROOT_DIR"
awscli_load_state

ALB_CLEAN_SCRIPT="${ROOT_DIR}/alb/clean.sh"
EC2_API_GATEWAY_CLEAN_SCRIPT="${ROOT_DIR}/ec2/api-gateway/clean.sh"
EC2_AUTHIFY_SERVICE_CLEAN_SCRIPT="${ROOT_DIR}/ec2/authify-service/clean.sh"
EC2_BASTION_CLEAN_SCRIPT="${ROOT_DIR}/ec2/bastion/clean.sh"
S3_PEM_CLEAN_SCRIPT="${ROOT_DIR}/s3/pem/clean.sh"
SSM_CLEAN_SCRIPT="${ROOT_DIR}/ssm/clean.sh"
ELASTICACHE_CLEAN_SCRIPT="${ROOT_DIR}/elasticache/clean.sh"
RDS_CLEAN_SCRIPT="${ROOT_DIR}/rds/clean.sh"
SG_CLEAN_SCRIPT="${ROOT_DIR}/security-group/clean.sh"
S3_SCHEMA_CLEAN_SCRIPT="${ROOT_DIR}/s3/schema/clean.sh"
ECR_CLEAN_SCRIPT="${ROOT_DIR}/ecr/clean.sh"
IAM_CLEAN_SCRIPT="${ROOT_DIR}/iam/clean.sh"
VPC_CLEAN_SCRIPT="${ROOT_DIR}/vpc/clean.sh"

"${ALB_CLEAN_SCRIPT}"
"${EC2_API_GATEWAY_CLEAN_SCRIPT}"
"${EC2_AUTHIFY_SERVICE_CLEAN_SCRIPT}"
"${EC2_BASTION_CLEAN_SCRIPT}"
"${S3_PEM_CLEAN_SCRIPT}"
"${SSM_CLEAN_SCRIPT}"
"${ELASTICACHE_CLEAN_SCRIPT}"
"${RDS_CLEAN_SCRIPT}"
"${SG_CLEAN_SCRIPT}"
"${S3_SCHEMA_CLEAN_SCRIPT}"
"${ECR_CLEAN_SCRIPT}"
"${IAM_CLEAN_SCRIPT}"
"${VPC_CLEAN_SCRIPT}"

: > "$VARIABLES_FILE"
echo "✅ Cleanup completed."
