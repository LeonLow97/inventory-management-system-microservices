#!/bin/bash

source variables.txt
source .env

# `execute` function checks and executes AWS CLI commands safely
execute() {
    echo "🔹 $1..."
    shift
    if "$@" >/dev/null 2>&1; then
        echo "✅ Successfully completed: $1"
    else
        echo "❌ Failed: $1" >&2
    fi
}

############### START OF CLEANUP ###############
echo "Starting Cleanup of IMS to AWS EC2..."

# Step 14: Delete RDS Instance
execute "Deleting RDS Instance: $DB_INSTANCE_IDENTIFIER..." \
    aws rds delete-db-instance --db-instance-identifier $DB_INSTANCE_IDENTIFIER --skip-final-snapshot --region $REGION >/dev/null 2>&1
execute "Waiting for RDS instance deletion to complete..." \
aws rds wait db-instance-deleted --db-instance-identifier $DB_INSTANCE_IDENTIFIER --region $REGION >/dev/null 2>&1
echo "✅ RDS instance deleted successfully."

# Step 15: Delete DB Subnet Group
execute "Deleting DB Subnet Group: $DB_SUBNET_GROUP_NAME..." \
    aws rds delete-db-subnet-group --db-subnet-group-name $DB_SUBNET_GROUP_NAME --region $REGION >/dev/null 2>&1

# Step 16: Delete SSM
aws ssm delete-parameter --name "/ims/db/hostname" --region $REGION >/dev/null 2>&1
aws ssm delete-parameter --name "/ims/db/master-username" --region $REGION >/dev/null 2>&1
aws ssm delete-parameter --name "/ims/db/master-password" --region $REGION >/dev/null 2>&1
aws ssm delete-parameter --name "/ims/db/db-name" --region $REGION >/dev/null 2>&1

# Step 1: Deregister Targets from Target Group
execute "Deregistering API Gateway instance from Target Group" \
  aws elbv2 deregister-targets --target-group-arn "$API_GATEWAY_TG_ARN" --targets Id="$API_GATEWAY_INSTANCE_ID_AZ1" Id="$API_GATEWAY_INSTANCE_ID_AZ2"

# Step 2: Delete ALB Listener
execute "Deleting ALB Listener" \
    aws elbv2 delete-listener --listener-arn "$ALB_LISTENER_ARN"

# Step 3: Delete ALB
execute "Deleting ALB" \
    aws elbv2 delete-load-balancer --load-balancer-arn "$ALB_ARN"

echo "Waiting for ALB to be deleted..."
aws elbv2 wait load-balancers-deleted --load-balancer-arns "$ALB_ARN" >/dev/null 2>&1
echo "✅ Successfully deleted ALB."

# Step 4: Delete Target Group
execute "Deleting Target Group" \
    aws elbv2 delete-target-group --target-group-arn "$API_GATEWAY_TG_ARN"

# Step 5: Terminate EC2 Instances
execute "Terminating EC2 instances" \
    aws ec2 terminate-instances --instance-ids \
    $API_GATEWAY_INSTANCE_ID_AZ1 $API_GATEWAY_INSTANCE_ID_AZ2 \
    $AUTH_SERVICE_INSTANCE_ID_AZ1 $AUTH_SERVICE_INSTANCE_ID_AZ2 \
    $BASTION_INSTANCE_ID_AZ1 $BASTION_INSTANCE_ID_AZ2

echo "Waiting for EC2 instances to be terminated..."
aws ec2 wait instance-terminated --instance-ids \
    $API_GATEWAY_INSTANCE_ID_AZ1 $API_GATEWAY_INSTANCE_ID_AZ2 \
    $AUTH_SERVICE_INSTANCE_ID_AZ1 $AUTH_SERVICE_INSTANCE_ID_AZ2 \
    $BASTION_INSTANCE_ID_AZ1 $BASTION_INSTANCE_ID_AZ2
echo "✅ Successfully terminated EC2 instances."

# Step 6: Delete all EC2 Key Pairs
echo "Deleting EC2 Key Pairs"
delete_key_pair() {
  if [ -n "$key_name" ]; then
    # Delete EC2 Key Pair from EC2 and local machine
    echo "\tDeleting Key Pair: $key_name..."
    aws ec2 delete-key-pair --key-name "$key_name" >/dev/null 2>&1
    rm -f "$key_name.pem"
    echo "\tKey Pair $key_name deleted successfully."
  fi
}
KEY_NAMES=(
  "$EC2_API_GATEWAY_KEY_NAME_AZ1"
  "$EC2_API_GATEWAY_KEY_NAME_AZ2"
  "$EC2_AUTH_MICROSERVICE_KEY_NAME_AZ1"
  "$EC2_AUTH_MICROSERVICE_KEY_NAME_AZ2"
  "$EC2_BASTION_HOST_KEY_NAME_AZ1"
  "$EC2_BASTION_HOST_KEY_NAME_AZ2"
)
for key_name in "${KEY_NAMES[@]}"; do
  if [ -n "$key_name" ]; then
    delete_key_pair "$key_name" &
  fi
done
wait
echo "✅ All EC2 Key Pairs deleted."

# Step 10: Delete NAT Gateways
execute "Deleting NAT Gateways" \
    aws ec2 delete-nat-gateway --nat-gateway-id "$NAT_GATEWAY_ID_AZ1" > /dev/null 2>&1 && \
    aws ec2 delete-nat-gateway --nat-gateway-id "$NAT_GATEWAY_ID_AZ2" > /dev/null 2>&1

# Step 7: Delete Security Groups
execute "Deleting Security Groups" \
    aws ec2 delete-security-group --group-id "$BASTION_SG_ID" && \
    aws ec2 delete-security-group --group-id "$DB_SG_ID" && \
    aws ec2 delete-security-group --group-id "$AUTH_SG_ID" && \
    aws ec2 delete-security-group --group-id "$API_GATEWAY_SG_ID" && \
    aws ec2 delete-security-group --group-id "$ALB_SG_ID"

# Step 8: Delete Public and Private Subnets
# execute "Deleting Public and Private Subnets" \
#     aws ec2 delete-subnet --subnet-id "$PRIVATE_RDS_SUBNET_AZ1_ID" && \
#     aws ec2 delete-subnet --subnet-id "$PRIVATE_RDS_SUBNET_AZ2_ID" && \
#     aws ec2 delete-subnet --subnet-id "$PRIVATE_SUBNET_AZ1_ID" && \
#     aws ec2 delete-subnet --subnet-id "$PRIVATE_SUBNET_AZ2_ID" && \
#     aws ec2 delete-subnet --subnet-id "$PUBLIC_SUBNET_AZ1_ID" && \
#     aws ec2 delete-subnet --subnet-id "$PUBLIC_SUBNET_AZ2_ID"

# Step 9: Releasing Elastic IPs
execute "Releasing Elastic IPs" \
    aws ec2 release-address --allocation-id "$EIP_AZ1" && \
    aws ec2 release-address --allocation-id "$EIP_AZ2"

# Step 13: Delete VPC
execute "Deleting VPC" \
    aws ec2 delete-vpc --vpc-id "$VpcId"

# # Step 11: Delete Route Table
# execute "Deleting Public Route Table" \
#     aws ec2 delete-route-table --route-table-id "$PUBLIC_RTB_ID"
# execute "Deleting Private Route Table in AZ1" \
#     aws ec2 delete-route-table --route-table-id "$PRIVATE_RTB_AZ1_ID"
# execute "Deleting Private Route Table in AZ2" \
#     aws ec2 delete-route-table --route-table-id "$PRIVATE_RTB_AZ2_ID"

# # Step 12: Detach and Delete Internet Gateway
# execute "Detaching and Deleting IGW" \
#     aws ec2 detach-internet-gateway --internet-gateway-id "$IGW_ID" --vpc-id "$VpcId" && \
#     aws ec2 delete-internet-gateway --internet-gateway-id "$IGW_ID"

echo "✅ Cleanup Completed for AWS EC2 Instances!"
############### END OF CLEANUP ###############
