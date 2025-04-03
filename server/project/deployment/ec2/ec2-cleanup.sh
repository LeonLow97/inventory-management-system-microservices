#!/bin/bash

source variables.txt

# `execute` function checks and executes AWS CLI commands safely
execute() {
    echo "🔹 $1..."
    shift
    if "$@" >/dev/null 2>&1; then
        echo "✅ Successfully completed: $1"
    else
        echo "❌ Failed: $1" >&2
        exit 1
    fi
}

############### START OF CLEANUP ###############
echo "Starting Cleanup of IMS to AWS EC2..."

# Step 1: Deregister Targets from Target Group
execute "Deregistering API Gateway instance from Target Group" \
    aws elbv2 deregister-targets --target-group-arn "$API_GATEWAY_TG_ARN" --targets Id="$API_GATEWAY_INSTANCE_ID"

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
    aws ec2 terminate-instances --instance-ids "$API_GATEWAY_INSTANCE_ID" "$AUTH_SERVICE_INSTANCE_ID" "$BASTION_INSTANCE_ID"

echo "Waiting for EC2 instances to be terminated..."
aws ec2 wait instance-terminated --instance-ids "$API_GATEWAY_INSTANCE_ID" "$AUTH_SERVICE_INSTANCE_ID" "$BASTION_INSTANCE_ID"
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
  "$EC2_API_GATEWAY_KEY_NAME"
  "$EC2_AUTH_MICROSERVICE_KEY_NAME"
  "$EC2_BASTION_HOST_KEY_NAME"
)
for key_name in "${KEY_NAMES[@]}"; do
  if [ -n "$key_name" ]; then
    delete_key_pair "$key_name" &
  fi
done
wait
echo "✅ All EC2 Key Pairs deleted."

# Step 7: Delete Security Groups
execute "Deleting Security Groups" \
    aws ec2 delete-security-group --group-id "$ALB_SG_ID" && \
    aws ec2 delete-security-group --group-id "$API_GATEWAY_SG_ID" && \
    aws ec2 delete-security-group --group-id "$AUTH_SG_ID" && \
    aws ec2 delete-security-group --group-id "$BASTION_SG_ID"

# Step 8: Delete Public and Private Subnets
execute "Deleting Public and Private Subnets" \
    aws ec2 delete-subnet --subnet-id "$PUBLIC_SUBNET_AZ1_ID" && \
    aws ec2 delete-subnet --subnet-id "$PUBLIC_SUBNET_AZ2_ID" && \
    aws ec2 delete-subnet --subnet-id "$PRIVATE_SUBNET_AZ1_ID" && \
    aws ec2 delete-subnet --subnet-id "$PRIVATE_SUBNET_AZ2_ID"

# Step 9: Releasing Elastic IPs
execute "Releasing Elastic IPs" \
    aws ec2 release-address --allocation-id "$EIP_AZ1" && \
    aws ec2 release-address --allocation-id "$EIP_AZ2"

# Step 10: Delete NAT Gateways
execute "Deleting NAT Gateways" \
    aws ec2 delete-nat-gateway --nat-gateway-id "$NAT_GATEWAY_ID_AZ1" && \
    aws ec2 delete-nat-gateway --nat-gateway-id "$NAT_GATEWAY_ID_AZ2"

# Step 11: Delete Route Table
execute "Deleting Route Table" \
    aws ec2 delete-route-table --route-table-id "$RTB_ID"

# Step 12: Detach and Delete Internet Gateway
execute "Detaching and Deleting IGW" \
    aws ec2 detach-internet-gateway --internet-gateway-id "$IGW_ID" --vpc-id "$VpcId" && \
    aws ec2 delete-internet-gateway --internet-gateway-id "$IGW_ID"

# Step 13: Delete VPC
execute "Deleting VPC" \
    aws ec2 delete-vpc --vpc-id "$VpcId"

echo "✅ Cleanup Completed for AWS EC2 Instances!"
############### END OF CLEANUP ###############
