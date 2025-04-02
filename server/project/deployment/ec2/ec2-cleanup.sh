# Cleanup

source variables.txt

# Cleanup ALB Resources
echo "Cleaning ALB Resources..."
if [ -n "$ALB_ARN" ]; then
  # Delete the Listener
  echo "\tDeleting ALB Listener..."
  aws elbv2 delete-listener --listener-arn $ALB_LISTENER_ARN
  echo "\tALB Listener deleted successfully."

  # Deregister the target (API Gateway EC2 instance) from the target group
  echo "\tDeregistering Target from the Target Group..."
  aws elbv2 deregister-targets --target-group-arn $API_GATEWAY_TG_ARN --targets Id=$API_GATEWAY_INSTANCE_ID
  echo "\tTarget deregistered successfully."

  # Delete the Target Group
  echo "\tDeleting Target Group..."
  aws elbv2 delete-target-group --target-group-arn $API_GATEWAY_TG_ARN
  echo "\tTarget Group deleted successfully."

  # Delete the Load Balancer (ALB)
  echo "\tDeleting ALB..."
  aws elbv2 delete-load-balancer --load-balancer-arn $ALB_ARN
  echo "\tALB deleted successfully."
fi

# Cleanup EC2 Instances
echo "Cleaning EC2 Instances..."

terminate_instance() {
  local instance_id=$1
  local name=$2

  if [ -n "$instance_id" ]; then
    # Terminate the EC2 Instance
    echo "\tTerminating $name EC2 Instance..."
    aws ec2 terminate-instances --instance-ids $instance_id  >/dev/null 2>&1

    # Wait until the instance is terminated
    echo "\tWaiting for $name EC2 Instance to be terminated..."
    aws ec2 wait instance-terminated --instance-ids $instance_id
    echo "\t$name EC2 Instance terminated successfully."
  fi
}

# Run in parallel
terminate_instance "$API_GATEWAY_INSTANCE_ID" "API Gateway" &
terminate_instance "$AUTH_SERVICE_INSTANCE_ID" "Authentication Microservice" &
terminate_instance "$BASTION_INSTANCE_ID" "Bastion Host" &
wait # wait for all EC2 instances to be terminated
echo "\tAll EC2 instances have been terminated."

# Cleanup EC2 Key Pairs
echo "Cleaning EC2 Key Pairs..."

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
echo "\tAll EC2 Key Pairs deleted."

sleep 60

# Cleaning NAT Gateways
echo "Cleaning NAT Gateways..."
echo "\tDeleting NAT Gateway in AZ1..."
aws ec2 delete-nat-gateway --nat-gateway-id $NAT_GATEWAY_ID_AZ1 >/dev/null 2>&1
echo "\tDeleting NAT Gateway in AZ2..."
aws ec2 delete-nat-gateway --nat-gateway-id $NAT_GATEWAY_ID_AZ2 >/dev/null 2>&1
echo "\tNAT Gateways have been deleted."

# Release Elastic IPs
echo "Releasing Elastic IPs..."
echo "\tReleasing Elastic IP in AZ1..."
# ASSOCIATION_ID_AZ1=$(aws ec2 associate-address --instance-id $API_GATEWAY_INSTANCE_ID --allocation-id $EIP_AZ1 --query 'AssociationId' --output text)
# aws ec2 disassociate-address --association-id $ASSOCIATION_ID_AZ1
aws ec2 release-address --allocation-id $EIP_AZ1
echo "\tReleasing Elastic IP in AZ2..."
# ASSOCIATION_ID_AZ2=$(aws ec2 associate-address --instance-id $API_GATEWAY_INSTANCE_ID --allocation-id $EIP_AZ2 --query 'AssociationId' --output text)
# aws ec2 disassociate-address --association-id $ASSOCIATION_ID_AZ2
aws ec2 release-address --allocation-id $EIP_AZ2
echo "\tAll Elastic IPs have been released."

# VPC Cleanup
echo "Deleting Security Groups in VPC..."

delete_security_group() {
  local group_id="$1"
  if [ -n "$group_id" ]; then
    echo "\tDeleting Security Group: $group_id..."
    aws ec2 delete-security-group --group-id "$group_id" &
  fi
}
delete_security_group "$AUTH_SG_ID"
delete_security_group "$API_GATEWAY_SG_ID"
delete_security_group "$ALB_SG_ID"
delete_security_group "$BASTION_SG_ID"
wait
echo "\tAll Security Groups deleted in VPC."

echo "Deleting Public and Private Subnets in VPC..."

delete_subnet_parallel() {
  local subnet_id="$1"
  if [ -n "$subnet_id" ]; then
    echo "\tDeleting Subnet: $subnet_id..."
    aws ec2 delete-subnet --subnet-id "$subnet_id" &
  fi
}

delete_subnet_parallel "$PUBLIC_SUBNET_AZ1_ID"
delete_subnet_parallel "$PUBLIC_SUBNET_AZ2_ID"
delete_subnet_parallel "$PRIVATE_SUBNET_AZ1_ID"
delete_subnet_parallel "$PRIVATE_SUBNET_AZ2_ID"
wait
echo "\tAll Public and Private subnets deleted in VPC."

if [ -n "$IGW_ID" ]; then
  echo "Detaching and Deleting Internet Gateway..."
  aws ec2 detach-internet-gateway --internet-gateway-id $IGW_ID --vpc-id $VpcId
  aws ec2 delete-internet-gateway --internet-gateway-id $IGW_ID
fi

if [ -n "$RTB_ID" ]; then
  echo "Deleting Route Table..."
  aws ec2 delete-route-table --route-table-id $RTB_ID
fi

if [ -n "$VpcId" ]; then
  echo "Deleting VPC..."
  aws ec2 delete-vpc --vpc-id $VpcId
fi
