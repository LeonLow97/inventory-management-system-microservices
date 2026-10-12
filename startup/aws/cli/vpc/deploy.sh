#!/bin/bash

set -euo pipefail

ROOT_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")/.." && pwd)"
cd "$ROOT_DIR"

source "$ROOT_DIR/lib/common.sh"
awscli_load_config "$ROOT_DIR"
awscli_load_state

wait_for_nat_gateway() {
  local nat_id="$1"
  local state

  while true; do
    state="$(aws ec2 describe-nat-gateways \
      --nat-gateway-ids "$nat_id" \
      --query 'NatGateways[0].State' \
      --output text \
      --region "$REGION")"
    case "$state" in
      available)
        echo "✅ NAT Gateway $nat_id is now available."
        return 0
        ;;
      pending)
        sleep 10
        ;;
      *)
        echo "NAT Gateway $nat_id entered unexpected state '$state'." >&2
        return 1
        ;;
    esac
  done
}

ensure_route() {
  local route_table_id="$1"
  local destination_cidr="$2"
  local target_option="$3"
  local target_field="$4"
  local target_id="$5"
  local route_count existing_target existing_state operation

  route_count="$(aws ec2 describe-route-tables \
    --route-table-ids "$route_table_id" \
    --query "length(RouteTables[0].Routes[?DestinationCidrBlock=='${destination_cidr}'])" \
    --output text \
    --region "$REGION")"
  existing_target="$(aws ec2 describe-route-tables \
    --route-table-ids "$route_table_id" \
    --query "RouteTables[0].Routes[?DestinationCidrBlock=='${destination_cidr}'].${target_field} | [0]" \
    --output text \
    --region "$REGION")"
  existing_state="$(aws ec2 describe-route-tables \
    --route-table-ids "$route_table_id" \
    --query "RouteTables[0].Routes[?DestinationCidrBlock=='${destination_cidr}'].State | [0]" \
    --output text \
    --region "$REGION")"

  if [[ "$existing_target" == "$target_id" && "$existing_state" == "active" ]]; then
    return 0
  fi

  operation=create-route
  if [[ "$route_count" != "0" ]]; then
    operation=replace-route
  fi

  aws ec2 "$operation" \
    --route-table-id "$route_table_id" \
    --destination-cidr-block "$destination_cidr" \
    "$target_option" "$target_id" \
    --region "$REGION" >/dev/null
}

if aws ec2 describe-vpcs --vpc-ids "${VPC_ID:-missing}" --region "$REGION" >/dev/null 2>&1; then
  echo "VPC $VPC_NAME already exists, skipping creation."
else
  VPC_ID="$(aws ec2 describe-vpcs \
    --filters "Name=tag:Name,Values=${VPC_NAME}" \
    --query 'Vpcs[0].VpcId' \
    --output text)"
  [[ "$VPC_ID" == "None" || "$VPC_ID" == "null" ]] && VPC_ID=""
  if [[ -z "$VPC_ID" ]]; then
    VPC_ID="$(aws ec2 create-vpc \
      --cidr-block "$VPC_CIDR" \
      --region "$REGION" \
      --tag-specifications "ResourceType=vpc,Tags=[{Key=Name,Value=${VPC_NAME}}]" \
      --query 'Vpc.VpcId' \
      --output text)"
    echo "✅ Created VPC with ID: $VPC_ID"
  fi
fi

aws ec2 modify-vpc-attribute --vpc-id "$VPC_ID" --enable-dns-support '{"Value":true}' >/dev/null
aws ec2 modify-vpc-attribute --vpc-id "$VPC_ID" --enable-dns-hostnames '{"Value":true}' >/dev/null

if ! aws ec2 describe-subnets --subnet-ids "${PUBLIC_SUBNET_AZ1_ID:-missing}" --region "$REGION" >/dev/null 2>&1; then
  PUBLIC_SUBNET_AZ1_ID="$(aws ec2 describe-subnets \
    --filters "Name=tag:Name,Values=${PUBLIC_SUBNET_AZ1_NAME}" "Name=vpc-id,Values=${VPC_ID}" \
    --query 'Subnets[0].SubnetId' \
    --output text)"
  [[ "$PUBLIC_SUBNET_AZ1_ID" == "None" || "$PUBLIC_SUBNET_AZ1_ID" == "null" ]] && PUBLIC_SUBNET_AZ1_ID=""
  if [[ -z "$PUBLIC_SUBNET_AZ1_ID" ]]; then
    PUBLIC_SUBNET_AZ1_ID="$(aws ec2 create-subnet \
      --vpc-id "$VPC_ID" \
      --cidr-block "$PUBLIC_SUBNET_AZ1_CIDR" \
      --availability-zone "$AZ1" \
      --tag-specifications "ResourceType=subnet,Tags=[{Key=Name,Value=${PUBLIC_SUBNET_AZ1_NAME}}]" \
      --query 'Subnet.SubnetId' \
      --output text)"
  fi
fi

if ! aws ec2 describe-subnets --subnet-ids "${PUBLIC_SUBNET_AZ2_ID:-missing}" --region "$REGION" >/dev/null 2>&1; then
  PUBLIC_SUBNET_AZ2_ID="$(aws ec2 describe-subnets \
    --filters "Name=tag:Name,Values=${PUBLIC_SUBNET_AZ2_NAME}" "Name=vpc-id,Values=${VPC_ID}" \
    --query 'Subnets[0].SubnetId' \
    --output text)"
  [[ "$PUBLIC_SUBNET_AZ2_ID" == "None" || "$PUBLIC_SUBNET_AZ2_ID" == "null" ]] && PUBLIC_SUBNET_AZ2_ID=""
  if [[ -z "$PUBLIC_SUBNET_AZ2_ID" ]]; then
    PUBLIC_SUBNET_AZ2_ID="$(aws ec2 create-subnet \
      --vpc-id "$VPC_ID" \
      --cidr-block "$PUBLIC_SUBNET_AZ2_CIDR" \
      --availability-zone "$AZ2" \
      --tag-specifications "ResourceType=subnet,Tags=[{Key=Name,Value=${PUBLIC_SUBNET_AZ2_NAME}}]" \
      --query 'Subnet.SubnetId' \
      --output text)"
  fi
fi

if ! aws ec2 describe-subnets --subnet-ids "${PRIVATE_SUBNET_AZ1_ID:-missing}" --region "$REGION" >/dev/null 2>&1; then
  PRIVATE_SUBNET_AZ1_ID="$(aws ec2 describe-subnets \
    --filters "Name=tag:Name,Values=${PRIVATE_SUBNET_AZ1_NAME}" "Name=vpc-id,Values=${VPC_ID}" \
    --query 'Subnets[0].SubnetId' \
    --output text)"
  [[ "$PRIVATE_SUBNET_AZ1_ID" == "None" || "$PRIVATE_SUBNET_AZ1_ID" == "null" ]] && PRIVATE_SUBNET_AZ1_ID=""
  if [[ -z "$PRIVATE_SUBNET_AZ1_ID" ]]; then
    PRIVATE_SUBNET_AZ1_ID="$(aws ec2 create-subnet \
      --vpc-id "$VPC_ID" \
      --cidr-block "$PRIVATE_SUBNET_AZ1_CIDR" \
      --availability-zone "$AZ1" \
      --tag-specifications "ResourceType=subnet,Tags=[{Key=Name,Value=${PRIVATE_SUBNET_AZ1_NAME}}]" \
      --query 'Subnet.SubnetId' \
      --output text)"
  fi
fi

if ! aws ec2 describe-subnets --subnet-ids "${PRIVATE_SUBNET_AZ2_ID:-missing}" --region "$REGION" >/dev/null 2>&1; then
  PRIVATE_SUBNET_AZ2_ID="$(aws ec2 describe-subnets \
    --filters "Name=tag:Name,Values=${PRIVATE_SUBNET_AZ2_NAME}" "Name=vpc-id,Values=${VPC_ID}" \
    --query 'Subnets[0].SubnetId' \
    --output text)"
  [[ "$PRIVATE_SUBNET_AZ2_ID" == "None" || "$PRIVATE_SUBNET_AZ2_ID" == "null" ]] && PRIVATE_SUBNET_AZ2_ID=""
  if [[ -z "$PRIVATE_SUBNET_AZ2_ID" ]]; then
    PRIVATE_SUBNET_AZ2_ID="$(aws ec2 create-subnet \
      --vpc-id "$VPC_ID" \
      --cidr-block "$PRIVATE_SUBNET_AZ2_CIDR" \
      --availability-zone "$AZ2" \
      --tag-specifications "ResourceType=subnet,Tags=[{Key=Name,Value=${PRIVATE_SUBNET_AZ2_NAME}}]" \
      --query 'Subnet.SubnetId' \
      --output text)"
  fi
fi

if ! aws ec2 describe-subnets --subnet-ids "${PRIVATE_RDS_SUBNET_AZ1_ID:-missing}" --region "$REGION" >/dev/null 2>&1; then
  PRIVATE_RDS_SUBNET_AZ1_ID="$(aws ec2 describe-subnets \
    --filters "Name=tag:Name,Values=${PRIVATE_RDS_SUBNET_AZ1_NAME}" "Name=vpc-id,Values=${VPC_ID}" \
    --query 'Subnets[0].SubnetId' \
    --output text)"
  [[ "$PRIVATE_RDS_SUBNET_AZ1_ID" == "None" || "$PRIVATE_RDS_SUBNET_AZ1_ID" == "null" ]] && PRIVATE_RDS_SUBNET_AZ1_ID=""
  if [[ -z "$PRIVATE_RDS_SUBNET_AZ1_ID" ]]; then
    PRIVATE_RDS_SUBNET_AZ1_ID="$(aws ec2 create-subnet \
      --vpc-id "$VPC_ID" \
      --cidr-block "$PRIVATE_RDS_SUBNET_AZ1_CIDR" \
      --availability-zone "$AZ1" \
      --tag-specifications "ResourceType=subnet,Tags=[{Key=Name,Value=${PRIVATE_RDS_SUBNET_AZ1_NAME}}]" \
      --query 'Subnet.SubnetId' \
      --output text)"
  fi
fi

if ! aws ec2 describe-subnets --subnet-ids "${PRIVATE_RDS_SUBNET_AZ2_ID:-missing}" --region "$REGION" >/dev/null 2>&1; then
  PRIVATE_RDS_SUBNET_AZ2_ID="$(aws ec2 describe-subnets \
    --filters "Name=tag:Name,Values=${PRIVATE_RDS_SUBNET_AZ2_NAME}" "Name=vpc-id,Values=${VPC_ID}" \
    --query 'Subnets[0].SubnetId' \
    --output text)"
  [[ "$PRIVATE_RDS_SUBNET_AZ2_ID" == "None" || "$PRIVATE_RDS_SUBNET_AZ2_ID" == "null" ]] && PRIVATE_RDS_SUBNET_AZ2_ID=""
  if [[ -z "$PRIVATE_RDS_SUBNET_AZ2_ID" ]]; then
    PRIVATE_RDS_SUBNET_AZ2_ID="$(aws ec2 create-subnet \
      --vpc-id "$VPC_ID" \
      --cidr-block "$PRIVATE_RDS_SUBNET_AZ2_CIDR" \
      --availability-zone "$AZ2" \
      --tag-specifications "ResourceType=subnet,Tags=[{Key=Name,Value=${PRIVATE_RDS_SUBNET_AZ2_NAME}}]" \
      --query 'Subnet.SubnetId' \
      --output text)"
  fi
fi

if aws ec2 describe-internet-gateways --internet-gateway-ids "${IGW_ID:-missing}" --region "$REGION" >/dev/null 2>&1; then
  echo "Internet gateway $IGW_NAME already exists, skipping creation."
else
  IGW_ID="$(aws ec2 describe-internet-gateways \
    --filters "Name=tag:Name,Values=${IGW_NAME}" "Name=attachment.vpc-id,Values=${VPC_ID}" \
    --query 'InternetGateways[0].InternetGatewayId' \
    --output text)"
  [[ "$IGW_ID" == "None" || "$IGW_ID" == "null" ]] && IGW_ID=""
  if [[ -z "$IGW_ID" ]]; then
    IGW_ID="$(aws ec2 create-internet-gateway \
      --region "$REGION" \
      --tag-specifications "ResourceType=internet-gateway,Tags=[{Key=Name,Value=${IGW_NAME}}]" \
      --query 'InternetGateway.InternetGatewayId' \
      --output text)"
  fi
fi

IGW_ATTACHMENT_VPC_ID="$(aws ec2 describe-internet-gateways \
  --internet-gateway-ids "$IGW_ID" \
  --query 'InternetGateways[0].Attachments[0].VpcId' \
  --output text \
  --region "$REGION")"
if [[ "$IGW_ATTACHMENT_VPC_ID" != "$VPC_ID" ]]; then
  aws ec2 attach-internet-gateway \
    --internet-gateway-id "$IGW_ID" \
    --vpc-id "$VPC_ID" \
    --region "$REGION" >/dev/null
fi

if ! aws ec2 describe-route-tables --route-table-ids "${PUBLIC_RTB_ID:-missing}" --region "$REGION" >/dev/null 2>&1; then
  PUBLIC_RTB_ID="$(aws ec2 describe-route-tables \
    --filters "Name=tag:Name,Values=${PUBLIC_RTB_NAME}" "Name=vpc-id,Values=${VPC_ID}" \
    --query 'RouteTables[0].RouteTableId' \
    --output text)"
  [[ "$PUBLIC_RTB_ID" == "None" || "$PUBLIC_RTB_ID" == "null" ]] && PUBLIC_RTB_ID=""
  if [[ -z "$PUBLIC_RTB_ID" ]]; then
    PUBLIC_RTB_ID="$(aws ec2 create-route-table \
      --vpc-id "$VPC_ID" \
      --region "$REGION" \
      --tag-specifications "ResourceType=route-table,Tags=[{Key=Name,Value=${PUBLIC_RTB_NAME}}]" \
      --query 'RouteTable.RouteTableId' \
      --output text)"
  fi
fi

ensure_route "$PUBLIC_RTB_ID" "0.0.0.0/0" "--gateway-id" "GatewayId" "$IGW_ID"

PUBLIC_ASSOC_AZ1="$(aws ec2 describe-route-tables \
  --route-table-ids "$PUBLIC_RTB_ID" \
  --query "RouteTables[0].Associations[?SubnetId=='${PUBLIC_SUBNET_AZ1_ID}'].RouteTableAssociationId | [0]" \
  --output text \
  --region "$REGION")"
[[ "$PUBLIC_ASSOC_AZ1" == "None" || "$PUBLIC_ASSOC_AZ1" == "null" ]] && PUBLIC_ASSOC_AZ1=""
if [[ -z "$PUBLIC_ASSOC_AZ1" ]]; then
  aws ec2 associate-route-table \
    --route-table-id "$PUBLIC_RTB_ID" \
    --subnet-id "$PUBLIC_SUBNET_AZ1_ID" \
    --region "$REGION" >/dev/null
fi

PUBLIC_ASSOC_AZ2="$(aws ec2 describe-route-tables \
  --route-table-ids "$PUBLIC_RTB_ID" \
  --query "RouteTables[0].Associations[?SubnetId=='${PUBLIC_SUBNET_AZ2_ID}'].RouteTableAssociationId | [0]" \
  --output text \
  --region "$REGION")"
[[ "$PUBLIC_ASSOC_AZ2" == "None" || "$PUBLIC_ASSOC_AZ2" == "null" ]] && PUBLIC_ASSOC_AZ2=""
if [[ -z "$PUBLIC_ASSOC_AZ2" ]]; then
  aws ec2 associate-route-table \
    --route-table-id "$PUBLIC_RTB_ID" \
    --subnet-id "$PUBLIC_SUBNET_AZ2_ID" \
    --region "$REGION" >/dev/null
fi

if ! aws ec2 describe-nat-gateways --nat-gateway-ids "${NAT_GATEWAY_ID_AZ1:-missing}" --region "$REGION" >/dev/null 2>&1; then
  NAT_GATEWAY_ID_AZ1="$(aws ec2 describe-nat-gateways \
    --filter "Name=tag:Name,Values=${NAT_GATEWAY_NAME_AZ1}" "Name=state,Values=pending,available" \
    --query 'NatGateways[0].NatGatewayId' \
    --output text)"
  [[ "$NAT_GATEWAY_ID_AZ1" == "None" || "$NAT_GATEWAY_ID_AZ1" == "null" ]] && NAT_GATEWAY_ID_AZ1=""
fi

EIP_AZ1_CHECK=""
if [[ -n "${EIP_AZ1:-}" ]]; then
  EIP_AZ1_CHECK="$(awscli_run_allowing_error_code "InvalidAllocationID.NotFound" \
    aws ec2 describe-addresses \
      --allocation-ids "$EIP_AZ1" \
      --query 'Addresses[0].AllocationId' \
      --output text \
      --region "$REGION")"
fi
[[ "$EIP_AZ1_CHECK" == "None" || "$EIP_AZ1_CHECK" == "null" ]] && EIP_AZ1_CHECK=""
if [[ -z "${EIP_AZ1:-}" || -z "$EIP_AZ1_CHECK" ]]; then
  EIP_AZ1=""
  if [[ -n "${NAT_GATEWAY_ID_AZ1:-}" ]]; then
    EIP_AZ1="$(aws ec2 describe-nat-gateways \
      --nat-gateway-ids "$NAT_GATEWAY_ID_AZ1" \
      --query 'NatGateways[0].NatGatewayAddresses[0].AllocationId' \
      --output text \
      --region "$REGION")"
    [[ "$EIP_AZ1" == "None" || "$EIP_AZ1" == "null" ]] && EIP_AZ1=""
  fi
fi

if [[ -z "${NAT_GATEWAY_ID_AZ1:-}" ]]; then
  if [[ -z "${EIP_AZ1:-}" ]]; then
    EIP_AZ1="$(aws ec2 allocate-address \
      --domain vpc \
      --tag-specifications "ResourceType=elastic-ip,Tags=[{Key=Name,Value=${NAT_GATEWAY_NAME_AZ1}}]" \
      --query 'AllocationId' \
      --output text)"
  fi
  NAT_GATEWAY_ID_AZ1="$(aws ec2 create-nat-gateway \
    --subnet-id "$PUBLIC_SUBNET_AZ1_ID" \
    --allocation-id "$EIP_AZ1" \
    --tag-specifications "ResourceType=natgateway,Tags=[{Key=Name,Value=${NAT_GATEWAY_NAME_AZ1}}]" \
    --query 'NatGateway.NatGatewayId' \
    --output text)"
fi

if ! aws ec2 describe-nat-gateways --nat-gateway-ids "${NAT_GATEWAY_ID_AZ2:-missing}" --region "$REGION" >/dev/null 2>&1; then
  NAT_GATEWAY_ID_AZ2="$(aws ec2 describe-nat-gateways \
    --filter "Name=tag:Name,Values=${NAT_GATEWAY_NAME_AZ2}" "Name=state,Values=pending,available" \
    --query 'NatGateways[0].NatGatewayId' \
    --output text)"
  [[ "$NAT_GATEWAY_ID_AZ2" == "None" || "$NAT_GATEWAY_ID_AZ2" == "null" ]] && NAT_GATEWAY_ID_AZ2=""
fi

EIP_AZ2_CHECK=""
if [[ -n "${EIP_AZ2:-}" ]]; then
  EIP_AZ2_CHECK="$(awscli_run_allowing_error_code "InvalidAllocationID.NotFound" \
    aws ec2 describe-addresses \
      --allocation-ids "$EIP_AZ2" \
      --query 'Addresses[0].AllocationId' \
      --output text \
      --region "$REGION")"
fi
[[ "$EIP_AZ2_CHECK" == "None" || "$EIP_AZ2_CHECK" == "null" ]] && EIP_AZ2_CHECK=""
if [[ -z "${EIP_AZ2:-}" || -z "$EIP_AZ2_CHECK" ]]; then
  EIP_AZ2=""
  if [[ -n "${NAT_GATEWAY_ID_AZ2:-}" ]]; then
    EIP_AZ2="$(aws ec2 describe-nat-gateways \
      --nat-gateway-ids "$NAT_GATEWAY_ID_AZ2" \
      --query 'NatGateways[0].NatGatewayAddresses[0].AllocationId' \
      --output text \
      --region "$REGION")"
    [[ "$EIP_AZ2" == "None" || "$EIP_AZ2" == "null" ]] && EIP_AZ2=""
  fi
fi

if [[ -z "${NAT_GATEWAY_ID_AZ2:-}" ]]; then
  if [[ -z "${EIP_AZ2:-}" ]]; then
    EIP_AZ2="$(aws ec2 allocate-address \
      --domain vpc \
      --tag-specifications "ResourceType=elastic-ip,Tags=[{Key=Name,Value=${NAT_GATEWAY_NAME_AZ2}}]" \
      --query 'AllocationId' \
      --output text)"
  fi
  NAT_GATEWAY_ID_AZ2="$(aws ec2 create-nat-gateway \
    --subnet-id "$PUBLIC_SUBNET_AZ2_ID" \
    --allocation-id "$EIP_AZ2" \
    --tag-specifications "ResourceType=natgateway,Tags=[{Key=Name,Value=${NAT_GATEWAY_NAME_AZ2}}]" \
    --query 'NatGateway.NatGatewayId' \
    --output text)"
fi

wait_for_nat_gateway "$NAT_GATEWAY_ID_AZ1" &
NAT_GATEWAY_AZ1_WAIT_PID=$!
wait_for_nat_gateway "$NAT_GATEWAY_ID_AZ2" &
NAT_GATEWAY_AZ2_WAIT_PID=$!
NAT_GATEWAY_WAIT_FAILED=0
wait "$NAT_GATEWAY_AZ1_WAIT_PID" || NAT_GATEWAY_WAIT_FAILED=1
wait "$NAT_GATEWAY_AZ2_WAIT_PID" || NAT_GATEWAY_WAIT_FAILED=1
if [[ "$NAT_GATEWAY_WAIT_FAILED" -ne 0 ]]; then
  exit 1
fi

if ! aws ec2 describe-route-tables --route-table-ids "${PRIVATE_RTB_AZ1_ID:-missing}" --region "$REGION" >/dev/null 2>&1; then
  PRIVATE_RTB_AZ1_ID="$(aws ec2 describe-route-tables \
    --filters "Name=tag:Name,Values=${PRIVATE_RTB_AZ1_NAME}" "Name=vpc-id,Values=${VPC_ID}" \
    --query 'RouteTables[0].RouteTableId' \
    --output text)"
  [[ "$PRIVATE_RTB_AZ1_ID" == "None" || "$PRIVATE_RTB_AZ1_ID" == "null" ]] && PRIVATE_RTB_AZ1_ID=""
  if [[ -z "$PRIVATE_RTB_AZ1_ID" ]]; then
    PRIVATE_RTB_AZ1_ID="$(aws ec2 create-route-table \
      --vpc-id "$VPC_ID" \
      --region "$REGION" \
      --tag-specifications "ResourceType=route-table,Tags=[{Key=Name,Value=${PRIVATE_RTB_AZ1_NAME}}]" \
      --query 'RouteTable.RouteTableId' \
      --output text)"
  fi
fi

if ! aws ec2 describe-route-tables --route-table-ids "${PRIVATE_RTB_AZ2_ID:-missing}" --region "$REGION" >/dev/null 2>&1; then
  PRIVATE_RTB_AZ2_ID="$(aws ec2 describe-route-tables \
    --filters "Name=tag:Name,Values=${PRIVATE_RTB_AZ2_NAME}" "Name=vpc-id,Values=${VPC_ID}" \
    --query 'RouteTables[0].RouteTableId' \
    --output text)"
  [[ "$PRIVATE_RTB_AZ2_ID" == "None" || "$PRIVATE_RTB_AZ2_ID" == "null" ]] && PRIVATE_RTB_AZ2_ID=""
  if [[ -z "$PRIVATE_RTB_AZ2_ID" ]]; then
    PRIVATE_RTB_AZ2_ID="$(aws ec2 create-route-table \
      --vpc-id "$VPC_ID" \
      --region "$REGION" \
      --tag-specifications "ResourceType=route-table,Tags=[{Key=Name,Value=${PRIVATE_RTB_AZ2_NAME}}]" \
      --query 'RouteTable.RouteTableId' \
      --output text)"
  fi
fi

ensure_route "$PRIVATE_RTB_AZ1_ID" "0.0.0.0/0" "--nat-gateway-id" "NatGatewayId" "$NAT_GATEWAY_ID_AZ1"
ensure_route "$PRIVATE_RTB_AZ2_ID" "0.0.0.0/0" "--nat-gateway-id" "NatGatewayId" "$NAT_GATEWAY_ID_AZ2"

PRIVATE_ASSOC_AZ1="$(aws ec2 describe-route-tables \
  --route-table-ids "$PRIVATE_RTB_AZ1_ID" \
  --query "RouteTables[0].Associations[?SubnetId=='${PRIVATE_SUBNET_AZ1_ID}'].RouteTableAssociationId | [0]" \
  --output text \
  --region "$REGION")"
[[ "$PRIVATE_ASSOC_AZ1" == "None" || "$PRIVATE_ASSOC_AZ1" == "null" ]] && PRIVATE_ASSOC_AZ1=""
if [[ -z "$PRIVATE_ASSOC_AZ1" ]]; then
  aws ec2 associate-route-table \
    --route-table-id "$PRIVATE_RTB_AZ1_ID" \
    --subnet-id "$PRIVATE_SUBNET_AZ1_ID" \
    --region "$REGION" >/dev/null
fi

PRIVATE_ASSOC_AZ2="$(aws ec2 describe-route-tables \
  --route-table-ids "$PRIVATE_RTB_AZ2_ID" \
  --query "RouteTables[0].Associations[?SubnetId=='${PRIVATE_SUBNET_AZ2_ID}'].RouteTableAssociationId | [0]" \
  --output text \
  --region "$REGION")"
[[ "$PRIVATE_ASSOC_AZ2" == "None" || "$PRIVATE_ASSOC_AZ2" == "null" ]] && PRIVATE_ASSOC_AZ2=""
if [[ -z "$PRIVATE_ASSOC_AZ2" ]]; then
  aws ec2 associate-route-table \
    --route-table-id "$PRIVATE_RTB_AZ2_ID" \
    --subnet-id "$PRIVATE_SUBNET_AZ2_ID" \
    --region "$REGION" >/dev/null
fi

awscli_write_state \
  "VPC_ID=$VPC_ID" \
  "PUBLIC_SUBNET_AZ1_ID=$PUBLIC_SUBNET_AZ1_ID" \
  "PUBLIC_SUBNET_AZ2_ID=$PUBLIC_SUBNET_AZ2_ID" \
  "PRIVATE_SUBNET_AZ1_ID=$PRIVATE_SUBNET_AZ1_ID" \
  "PRIVATE_SUBNET_AZ2_ID=$PRIVATE_SUBNET_AZ2_ID" \
  "PRIVATE_RDS_SUBNET_AZ1_ID=$PRIVATE_RDS_SUBNET_AZ1_ID" \
  "PRIVATE_RDS_SUBNET_AZ2_ID=$PRIVATE_RDS_SUBNET_AZ2_ID" \
  "IGW_ID=$IGW_ID" \
  "PUBLIC_RTB_ID=$PUBLIC_RTB_ID" \
  "PRIVATE_RTB_AZ1_ID=$PRIVATE_RTB_AZ1_ID" \
  "PRIVATE_RTB_AZ2_ID=$PRIVATE_RTB_AZ2_ID" \
  "EIP_AZ1=$EIP_AZ1" \
  "EIP_AZ2=$EIP_AZ2" \
  "NAT_GATEWAY_ID_AZ1=$NAT_GATEWAY_ID_AZ1" \
  "NAT_GATEWAY_ID_AZ2=$NAT_GATEWAY_ID_AZ2"

echo "✅ VPC networking ready."
