#!/bin/bash

source .env
PUBLIC_IP=$(curl -s https://checkip.amazonaws.com)

set -e  # Exit on command failure
set -u  # Treat unset variables as an error
set -o pipefail  # Catch pipeline errors

############### START OF DEPLOYMENT ###############
echo "Starting Deployment of IMS to AWS EC2..."

##### STEP 1: Create a VPC, Public and Private Subnets, Internet Gateways, Route Table
## Create VPC
echo "Creating VPC..."
VpcId=$(aws ec2 create-vpc --cidr-block 10.0.0.0/16 \ 
  --region $REGION \
  --query 'Vpc.VpcId' \
  --output text) || {
  echo "❌ Failed to create VPC";
  exit 1;
}
echo "✅ VPC created: VpcId=$VpcId"

# Create 2 Public Subnets and 2 Private Subnets in VPC
# For high availability, we created them in multiple Availability Zones (for ALB to distribute incoming requests). 
echo "Creating Public & Private Subnets..."
PUBLIC_SUBNET_AZ1_ID=$(aws ec2 create-subnet --vpc-id "$VpcId" --cidr-block 10.0.1.0/24 --availability-zone "$AZ1" --query 'Subnet.SubnetId' --output text) || { echo "❌ Failed to create Public Subnet AZ1"; exit 1; }
PUBLIC_SUBNET_AZ2_ID=$(aws ec2 create-subnet --vpc-id "$VpcId" --cidr-block 10.0.2.0/24 --availability-zone "$AZ2" --query 'Subnet.SubnetId' --output text) || { echo "❌ Failed to create Public Subnet AZ2"; exit 1; }
PRIVATE_SUBNET_AZ1_ID=$(aws ec2 create-subnet --vpc-id "$VpcId" --cidr-block 10.0.3.0/24 --availability-zone "$AZ1" --query 'Subnet.SubnetId' --output text) || { echo "❌ Failed to create Private Subnet AZ1"; exit 1; }
PRIVATE_SUBNET_AZ2_ID=$(aws ec2 create-subnet --vpc-id "$VpcId" --cidr-block 10.0.4.0/24 --availability-zone "$AZ2" --query 'Subnet.SubnetId' --output text) || { echo "❌ Failed to create Private Subnet AZ2"; exit 1; }
echo "✅ Subnets created:"
echo "\t- PUBLIC_SUBNET_AZ1_ID=$PUBLIC_SUBNET_AZ1_ID"
echo "\t- PUBLIC_SUBNET_AZ2_ID=$PUBLIC_SUBNET_AZ2_ID"
echo "\t- PRIVATE_SUBNET_AZ1_ID=$PRIVATE_SUBNET_AZ1_ID"
echo "\t- PRIVATE_SUBNET_AZ2_ID=$PRIVATE_SUBNET_AZ2_ID"

# Create Internet Gateway (connects VPC to Internet) and attach it to the VPC
echo "Creating Internet Gateway..."
IGW_ID=$(aws ec2 create-internet-gateway \
  --region $REGION \
  --query 'InternetGateway.InternetGatewayId' \
  --output text) || {
  echo "❌ Failed to create Internet Gateway";
  exit 1;
}
echo "✅ Internet Gateway created: $IGW_ID"
aws ec2 attach-internet-gateway \
  --internet-gateway-id $IGW_ID \
  --vpc-id $VpcId || {
  echo "❌ Failed to attach Internet Gateway to VPC";
  exit 1;
}
echo "✅ Internet Gateway attached to VPC: $VpcId"

# Find Default Route Table ID
echo "🔍 Finding Default Route Table in VPC..."
DEFAULT_RTB_ID=$(aws ec2 describe-route-tables \
  --filters "Name=vpc-id,Values=$VpcId" \
  --query 'RouteTables[?Associations[?Main==`true`]].RouteTableId' \
  --output text) || {
  echo "❌ Failed to find Default Route Table";
  exit 1;
}
echo "✅ Default Route Table found: $DEFAULT_RTB_ID"

# Create a new Route Table (for Public Subnets) and associate it with the IGW
echo "Creating a new Route Table..."
RTB_ID=$(aws ec2 create-route-table \
  --vpc-id $VpcId \
  --region $REGION \
  --query 'RouteTable.RouteTableId' \
  --output text) || {
  echo "❌ Failed to create Route Table";
  exit 1;
}
echo "✅ Route Table created: $RTB_ID"

# Add a route to the Internet Gateway (0.0.0.0 <-> IGW)
echo "Creating Route for Internet Gateway in Route Table..."
aws ec2 create-route \
  --route-table-id $RTB_ID \
  --destination-cidr-block 0.0.0.0/0 \
  --gateway-id $IGW_ID \
  --region $REGION >/dev/null 2>&1 || {
  echo "❌ Failed to create route to Internet Gateway";
  exit 1;
}
echo "✅ Route to IGW added successfully."

# Associate the Route Table with Public Subnets
echo "Associating Route Table with Public Subnets..."
for SUBNET in "$PUBLIC_SUBNET_AZ1_ID" "$PUBLIC_SUBNET_AZ2_ID"; do
  aws ec2 associate-route-table \
    --route-table-id "$RTB_ID" \
    --subnet-id "$SUBNET" \
    --region "$REGION" >/dev/null 2>&1 || {
    echo "❌ Failed to associate Route Table $RTB_ID with Subnet $SUBNET";
    exit 1;
  }
done
echo "✅ Successfully associated Route Table with Public Subnets."

##### STEP 2: Create Security Groups and setting inbound rules in SG
echo "Creating Security Groups..."
ALB_SG_ID=$(aws ec2 create-security-group \
  --group-name ALB-SG \
  --description "Security group for ALB" \
  --vpc-id $VpcId --region $REGION \
  --query 'GroupId' --output text) || { echo "❌ Failed to create security group: ALB-SG"; exit 1; }
API_GATEWAY_SG_ID=$(aws ec2 create-security-group \
  --group-name API-Gateway-SG \
  --description "Security group for API Gateway" \
  --vpc-id $VpcId --region $REGION \
  --query 'GroupId' --output text) || { echo "❌ Failed to create security group: API-Gateway-SG"; exit 1; }
AUTH_SG_ID=$(aws ec2 create-security-group \
  --group-name Auth-Service-SG \
  --description "Security group for authentication microservice" \
  --vpc-id $VpcId --region $REGION \
  --query 'GroupId' --output text) || { echo "❌ Failed to create security group: Auth-Service-SG"; exit 1; }
BASTION_SG_ID=$(aws ec2 create-security-group \
  --group-name Bastion-SG \
  --description "Security group for Bastion Host" \
  --vpc-id $VpcId --region $REGION \
  --query 'GroupId' --output text) || { echo "❌ Failed to create security group: Bastion-SG"; exit 1; }
echo "✅ Successfully created all security groups."
echo "\t- ALB_SG_ID=$ALB_SG_ID"
echo "\t- API_GATEWAY_SG_ID=$API_GATEWAY_SG_ID"
echo "\t- AUTH_SG_ID=$AUTH_SG_ID"
echo "\t- BASTION_SG_ID=$BASTION_SG_ID"

# Setting Inbound Rules for Security Groups
echo "Setting Inbound Rules for Security Groups..."
echo "Setting inbound rules for Application Load Balancer ($ALB_SG_ID)"
aws ec2 authorize-security-group-ingress --group-id $ALB_SG_ID --region $REGION \
  --ip-permissions '[
    {"IpProtocol": "tcp", "FromPort": 80, "ToPort": 80, "IpRanges": [{"CidrIp": "0.0.0.0/0"}]},
    {"IpProtocol": "tcp", "FromPort": 443, "ToPort": 443, "IpRanges": [{"CidrIp": "0.0.0.0/0"}]}
  ]' >/dev/null 2>&1 || { echo "❌ Failed to set inbound rules for $ALB_SG_ID"; exit 1; } &

echo "Setting inbound rules for API Gateway ($API_GATEWAY_SG_ID)"
aws ec2 authorize-security-group-ingress --group-id $API_GATEWAY_SG_ID \
  --ip-permissions '[
    {"IpProtocol": "tcp", "FromPort": '"$PORT_API_GATEWAY"', "ToPort": '"$PORT_API_GATEWAY"', "UserIdGroupPairs": [{"GroupId": "'"$ALB_SG_ID"'"}]},
    {"IpProtocol": "tcp", "FromPort": 22, "ToPort": 22, "UserIdGroupPairs": [{"GroupId": "'"$BASTION_SG_ID"'"}]},
    {"IpProtocol": "tcp", "FromPort": 80, "ToPort": 80, "UserIdGroupPairs": [{"GroupId": "'"$BASTION_SG_ID"'"}]},
    {"IpProtocol": "tcp", "FromPort": 443, "ToPort": 443, "UserIdGroupPairs": [{"GroupId": "'"$BASTION_SG_ID"'"}]}
  ]' >/dev/null 2>&1 || { echo "❌ Failed to set rules for API Gateway"; exit 1; } &

echo "Setting inbound rules for Authentication Microservice ($AUTH_SG_ID)"
aws ec2 authorize-security-group-ingress --group-id $AUTH_SG_ID \
  --ip-permissions '[
    {"IpProtocol": "tcp", "FromPort": '"$PORT_AUTH_SERVICE"',"ToPort": '"$PORT_AUTH_SERVICE"',"UserIdGroupPairs": [{"GroupId": "'"$API_GATEWAY_SG_ID"'"}]},
    {"IpProtocol": "tcp", "FromPort": '"$PORT_POSTGRESQL"',"ToPort": '"$PORT_POSTGRESQL"',"UserIdGroupPairs": [{"GroupId": "'"$AUTH_SG_ID"'"}]}
  ]' >/dev/null 2>&1 || { echo "❌ Failed to set rules for Authentication Microservice"; exit 1; } &

echo "Setting inbound rules for Bastion Host ($BASTION_SG_ID)"
aws ec2 authorize-security-group-ingress --group-id $BASTION_SG_ID \
  --protocol tcp --port 22 --cidr $PUBLIC_IP/32 >/dev/null 2>&1 || { echo "❌ Failed to set rules for Bastion Host"; exit 1; } &
wait
echo "✅ Successfully set inbound rules for all security groups."

# Creating EC2 Instances Key Pairs for SSH
echo "Creating Key Pairs for EC2 Instances..."
echo "Creating key pair for API Gateway ($EC2_API_GATEWAY_KEY_NAME)"
aws ec2 create-key-pair --key-name $EC2_API_GATEWAY_KEY_NAME --query 'KeyMaterial' --output text > $EC2_API_GATEWAY_KEY_NAME.pem || { echo "❌ Failed to create key pair for API Gateway"; exit 1; } &

echo "Creating key pair for Authentication Microservice ($EC2_AUTH_MICROSERVICE_KEY_NAME)"
aws ec2 create-key-pair --key-name $EC2_AUTH_MICROSERVICE_KEY_NAME --query 'KeyMaterial' --output text > $EC2_AUTH_MICROSERVICE_KEY_NAME.pem || { echo "❌ Failed to create key pair for Authentication Microservice"; exit 1; } &

echo "Creating key pair for Bastion Host ($EC2_BASTION_HOST_KEY_NAME)"
aws ec2 create-key-pair --key-name $EC2_BASTION_HOST_KEY_NAME --query 'KeyMaterial' --output text > $EC2_BASTION_HOST_KEY_NAME.pem || { echo "❌ Failed to create key pair for Bastion Host"; exit 1; } &
wait
echo "✅ Successfully created key pairs for EC2 instances."

##### STEP 3: Create NAT Gateway
# Allocate two Elastic IPs (one for each NAT Gateway)
echo "Allocating Elastic IPs for NAT Gateways..."
EIP_AZ1=$(aws ec2 allocate-address --domain vpc --query 'AllocationId' --output text) || { echo "❌ Failed to allocate Elastic IP for AZ1"; exit 1; } &
EIP_AZ2=$(aws ec2 allocate-address --domain vpc --query 'AllocationId' --output text) || { echo "❌ Failed to allocate Elastic IP for AZ2"; exit 1; } &
wait
echo "✅ Successfully allocated Elastic IPs."

# Create NAT Gateways in both public subnets with different Elastic IPs
echo "Creating NAT Gateways..."
NAT_GATEWAY_ID_AZ1=$(aws ec2 create-nat-gateway --subnet-id $PUBLIC_SUBNET_AZ1_ID --allocation-id $EIP_AZ1 --query "NatGateway.NatGatewayId" --output text) || { echo "❌ Failed to create NAT Gateway in AZ1"; exit 1; } &
NAT_GATEWAY_ID_AZ2=$(aws ec2 create-nat-gateway --subnet-id $PUBLIC_SUBNET_AZ2_ID --allocation-id $EIP_AZ2 --query "NatGateway.NatGatewayId" --output text) || { echo "❌ Failed to create NAT Gateway in AZ2"; exit 1; } &
wait
echo "✅ Successfully created NAT Gateways."

# Wait for both NAT Gateways to become available
echo "Waiting for NAT Gateways to become available..."
wait_for_nat_gateway() {
    local NAT_ID=$1
    while [ "$(aws ec2 describe-nat-gateways --nat-gateway-ids $NAT_ID --query "NatGateways[0].State" --output text)" != "available" ]; do
        sleep 10
    done
    echo "✅ NAT Gateway $NAT_ID is now available."
}
wait_for_nat_gateway $NAT_GATEWAY_ID_AZ1 &
wait_for_nat_gateway $NAT_GATEWAY_ID_AZ2 &
wait
echo "✅ Both NAT Gateways are available. Proceeding with the next step."

# Update Private Subnet Route Tables to use their respective NAT Gateways
echo "Updating route tables for private subnets..."
aws ec2 create-route --route-table-id $DEFAULT_RTB_ID --destination-cidr-block 0.0.0.0/0 --nat-gateway-id $NAT_GATEWAY_ID_AZ1 >/dev/null 2>&1 || { echo "❌ Failed to update route table for AZ1"; exit 1; } &
aws ec2 create-route --route-table-id $DEFAULT_RTB_ID --destination-cidr-block 0.0.0.0/0 --nat-gateway-id $NAT_GATEWAY_ID_AZ2 >/dev/null 2>&1 || { echo "❌ Failed to update route table for AZ2"; exit 1; } &
wait
echo "✅ Successfully updated route tables for private subnets."

##### STEP 4: Launch EC2 Instances
echo "Launching API Gateway EC2 instance..."
API_GATEWAY_INSTANCE_ID=$(aws ec2 run-instances \
  --image-id ami-06661384e66f2da0e \
  --count 1 \
  --instance-type t3.micro \
  --key-name $EC2_API_GATEWAY_KEY_NAME \
  --security-group-ids $API_GATEWAY_SG_ID \
  --subnet-id $PRIVATE_SUBNET_AZ1_ID \
  --user-data file://deployment/ec2/api-gateway-script.sh \
  --query 'Instances[0].InstanceId' \
  --output text) || { echo "❌ Failed to launch API Gateway instance"; exit 1; } &

echo "Launching Authentication Microservice EC2 instance..."
AUTH_SERVICE_INSTANCE_ID=$(aws ec2 run-instances \
  --image-id ami-06661384e66f2da0e \
  --count 1 \
  --instance-type t3.micro \
  --key-name $EC2_AUTH_MICROSERVICE_KEY_NAME \
  --security-group-ids $AUTH_SG_ID \
  --subnet-id $PRIVATE_SUBNET_AZ1_ID \
  --user-data file://deployment/ec2/authentication-script.sh \
  --query 'Instances[0].InstanceId' \
  --output text) || { echo "❌ Failed to launch Authentication Microservice instance"; exit 1; } &

echo "Launching Bastion Host EC2 instance..."
BASTION_INSTANCE_ID=$(aws ec2 run-instances \
  --image-id ami-06661384e66f2da0e \
  --count 1 \
  --instance-type t3.micro \
  --key-name $EC2_BASTION_HOST_KEY_NAME \
  --security-group-ids $BASTION_SG_ID \
  --subnet-id $PUBLIC_SUBNET_AZ1_ID \
  --associate-public-ip-address \
  --query 'Instances[0].InstanceId' \
  --output text) || { echo "❌ Failed to launch Bastion Host instance"; exit 1; } &
wait
echo "API_GATEWAY_INSTANCE_ID=$API_GATEWAY_INSTANCE_ID"
echo "AUTH_SERVICE_INSTANCE_ID=$AUTH_SERVICE_INSTANCE_ID"
echo "BASTION_INSTANCE_ID=$BASTION_INSTANCE_ID"
echo "✅ Launched EC2 Instances."

# Attach IAM Role (IMS-EC2-Role) to EC2 Instances so they can pull ECR images
echo "Associating IAM role with EC2 instances..."
aws ec2 associate-iam-instance-profile --instance-id $API_GATEWAY_INSTANCE_ID \
    --iam-instance-profile Name=IMS-EC2-Role >/dev/null 2>&1 || { echo "❌ Failed to associate IAM role with API Gateway instance"; exit 1; } &
aws ec2 associate-iam-instance-profile --instance-id $AUTH_SERVICE_INSTANCE_ID \
    --iam-instance-profile Name=IMS-EC2-Role >/dev/null 2>&1 || { echo "❌ Failed to associate IAM role with Authentication Microservice instance"; exit 1; } &
wait
echo "✅ Associated EC2 instances with IAM Role."

# Wait until EC2 instances are running and check their statuses
echo "Waiting for EC2 instances to become running..."
aws ec2 wait instance-running --instance-ids $API_GATEWAY_INSTANCE_ID || { echo "❌ API Gateway instance failed to start"; exit 1; }
aws ec2 wait instance-running --instance-ids $AUTH_SERVICE_INSTANCE_ID || { echo "❌ Authentication Microservice instance failed to start"; exit 1; }
aws ec2 wait instance-status-ok --instance-ids $API_GATEWAY_INSTANCE_ID || { echo "❌ API Gateway instance failed to reach status OK"; exit 1; }
aws ec2 wait instance-status-ok --instance-ids $AUTH_SERVICE_INSTANCE_ID || { echo "❌ Authentication Microservice instance failed to reach status OK"; exit 1; }
echo "✅ EC2 instances are running and ready."

# Get the private IPv4 address of the newly launched EC2 instance
echo "Retrieving the private IPv4 address of the API Gateway EC2 instance..."
API_GATEWAY_PRIVATE_IP=$(aws ec2 describe-instances \
  --instance-ids $API_GATEWAY_INSTANCE_ID \
  --query 'Reservations[0].Instances[0].PrivateIpAddress' \
  --output text) || { echo "❌ Failed to retrieve private IP for API Gateway instance"; exit 1; } &
echo "Retrieving the private IPv4 address of the Authentication Microservice EC2 instance..."
AUTH_SERVICE_PRIVATE_IP=$(aws ec2 describe-instances \
  --instance-ids $AUTH_SERVICE_INSTANCE_ID \
  --query 'Reservations[0].Instances[0].PrivateIpAddress' \
  --output text) || { echo "❌ Failed to retrieve private IP for API Gateway instance"; exit 1; } &

echo "Retrieving the public IPv4 address of the Bastion Host EC2 instance..."
BASTION_PUBLIC_IP=$(aws ec2 describe-instances \
  --instance-ids $BASTION_INSTANCE_ID \
  --query 'Reservations[0].Instances[0].PublicIpAddress' --output text) || { echo "❌ Failed to retrieve public IP for Bastion instance"; exit 1; } &
wait

echo "API_GATEWAY_PRIVATE_IP=$API_GATEWAY_PRIVATE_IP"
echo "API_GATEWAY_PRIVATE_IP=$API_GATEWAY_PRIVATE_IP"
echo "BASTION_PUBLIC_IP=$BASTION_PUBLIC_IP"

echo "\nRun the following commands in order to SSH into API Gateway EC2 Instance"
echo "ssh -i IMS_BASTION.pem ec2-user@$BASTION_PUBLIC_IP"
echo "scp -i IMS_BASTION.pem IMS_API_GATEWAY_KEY_PAIR.pem ec2-user@$BASTION_PUBLIC_IP:~"
echo "ssh -i IMS_BASTION.pem ec2-user@$BASTION_PUBLIC_IP"
echo "ssh -i IMS_API_GATEWAY_KEY_PAIR.pem ec2-user@$API_GATEWAY_PRIVATE_IP"

echo "\nRun the following commands in order to SSH into Authentication Microservice EC2 Instance"
echo "ssh -i IMS_BASTION.pem ec2-user@$BASTION_PUBLIC_IP"
echo "scp -i IMS_BASTION.pem IMS_AUTH_SERVICE_KEY_PAIR.pem ec2-user@$BASTION_PUBLIC_IP:~"
echo "ssh -i IMS_BASTION.pem ec2-user@$BASTION_PUBLIC_IP"
echo "ssh -i IMS_AUTH_SERVICE_KEY_PAIR.pem ec2-user@$AUTH_SERVICE_PRIVATE_IP"

echo "\n✅ Retrieved Private IPv4 Addresses for EC2 instances and prepared SSH commands"

##### STEP 5: Launch Application Load Balancer (ALB) in Public Subnet
# Create ALB --> Set up a public Application Load Balancer in a specified subnet
echo "Creating the Application Load Balancer..."
ALB_ARN=$(aws elbv2 create-load-balancer --name ims-alb --subnets $PUBLIC_SUBNET_AZ1_ID $PUBLIC_SUBNET_AZ2_ID --security-groups $ALB_SG_ID --query "LoadBalancers[0].LoadBalancerArn" --output text) || { echo "❌ Failed to create ALB"; exit 1; } &

# Create Target Group --> Defines where ALB should forward traffic (API Gateway on Port 8080)
echo "Creating the Target Group for API Gateway..."
API_GATEWAY_TG_ARN=$(aws elbv2 create-target-group --name "$API_GATEWAY_SERVICE-tg" --protocol HTTP --port 8080 --vpc-id $VpcId --query "TargetGroups[0].TargetGroupArn" --output text) || { echo "❌ Failed to create Target Group"; exit 1; } &
wait
echo "✅ Created ALB and Target Group for API Gateway."

# Modify API Gateway target group health check endpoint
echo "Modifying health check settings for the API Gateway Target Group..."
aws elbv2 modify-target-group \
    --target-group-arn $API_GATEWAY_TG_ARN \
    --health-check-path "/healthcheck" \
    --health-check-port "8080" \
    --health-check-protocol "HTTP" \
    --health-check-interval-seconds 30 \
    --health-check-timeout-seconds 5 \
    --healthy-threshold-count 5 \
    --unhealthy-threshold-count 2 >/dev/null 2>&1 || { echo "❌ Failed to modify health check"; exit 1; }
echo "✅ Successfully configured API Gateway Target Group."

# Register Target --> Links the API Gateway EC2 Instance to the target group
echo "Registering API Gateway EC2 instance to the Target Group..."
aws elbv2 register-targets --target-group-arn $API_GATEWAY_TG_ARN --targets Id=$API_GATEWAY_INSTANCE_ID || { echo "❌ Failed to register targets"; exit 1; }
echo "✅ Successfully Registered API Gateway EC2 instance to target group."

# Create Listener --> Configured ALB to listen on port 80 and forward requests to Target Group
echo "Creating ALB Listener..."
ALB_LISTENER_ARN=$(aws elbv2 create-listener --load-balancer-arn $ALB_ARN --protocol HTTP --port 80 --default-actions Type=forward,TargetGroupArn=$API_GATEWAY_TG_ARN --query "Listeners[0].ListenerArn" --output text) || { echo "❌ Failed to create ALB listener"; exit 1; }
echo "✅ Successfully created ALB Listener."

# Wait for the ALB to be in available state
echo "Waiting for the Load Balancer $ALB_ARN to be available..."
aws elbv2 wait load-balancer-available --load-balancer-arn $ALB_ARN || { echo "❌ ALB is not available"; exit 1; }
echo "✅ ALB $ALB_ARN is AVAILABLE."

# Get the DNS name of the Load Balancer
ALB_DNS_NAME=$(aws elbv2 describe-load-balancers --load-balancer-arn $ALB_ARN --query "LoadBalancers[0].DNSName" --output text) || { echo "❌ Failed to get ALB DNS name"; exit 1; }
echo "\t- ALB_DNS_NAME=$ALB_DNS_NAME"

# Test API Gateway HealthCheck endpoint
echo "Testing API Gateway HealthCheck endpoint..."
HTTP_STATUS=$(curl -o /dev/null -s -w "%{http_code}" http://$ALB_DNS_NAME/healthcheck)
if [ "$HTTP_STATUS" -eq 200 ]; then
  echo -e "✅ Health check passed! (HTTP 200)"
else
  echo -e "❌ Health check failed! (HTTP $HTTP_STATUS)"
fi

echo "✅ Launched IMS to AWS EC2!"

############### END OF DEPLOYMENT ###############

# Clean variables file (if exists)
if [ -f variables.txt ]; then
  > variables.txt
fi

# Add variables to variables.txt for cleanup script
echo "VpcId=$VpcId" >> variables.txt
echo "DEFAULT_RTB_ID=$DEFAULT_RTB_ID" >> variables.txt
echo "PUBLIC_SUBNET_AZ1_ID=$PUBLIC_SUBNET_AZ1_ID" >> variables.txt
echo "PUBLIC_SUBNET_AZ2_ID=$PUBLIC_SUBNET_AZ2_ID" >> variables.txt
echo "PRIVATE_SUBNET_AZ1_ID=$PRIVATE_SUBNET_AZ1_ID" >> variables.txt
echo "PRIVATE_SUBNET_AZ2_ID=$PRIVATE_SUBNET_AZ2_ID" >> variables.txt
echo "IGW_ID=$IGW_ID" >> variables.txt
echo "RTB_ID=$RTB_ID" >> variables.txt
echo "EC2_API_GATEWAY_KEY_NAME=$EC2_API_GATEWAY_KEY_NAME" >> variables.txt
echo "EC2_AUTH_MICROSERVICE_KEY_NAME=$EC2_AUTH_MICROSERVICE_KEY_NAME" >> variables.txt
echo "EC2_BASTION_HOST_KEY_NAME=$EC2_BASTION_HOST_KEY_NAME" >> variables.txt
echo "ALB_SG_ID=$ALB_SG_ID" >> variables.txt
echo "API_GATEWAY_SG_ID=$API_GATEWAY_SG_ID" >> variables.txt
echo "AUTH_SG_ID=$AUTH_SG_ID" >> variables.txt
echo "BASTION_SG_ID=$BASTION_SG_ID" >> variables.txt
echo "EIP_AZ1=$EIP_AZ1" >> variables.txt
echo "EIP_AZ2=$EIP_AZ2" >> variables.txt
echo "NAT_GATEWAY_ID_AZ1=$NAT_GATEWAY_ID_AZ1" >> variables.txt
echo "NAT_GATEWAY_ID_AZ2=$NAT_GATEWAY_ID_AZ2" >> variables.txt
echo "API_GATEWAY_INSTANCE_ID=$API_GATEWAY_INSTANCE_ID" >> variables.txt
echo "AUTH_SERVICE_INSTANCE_ID=$AUTH_SERVICE_INSTANCE_ID" >> variables.txt
echo "BASTION_INSTANCE_ID=$BASTION_INSTANCE_ID" >> variables.txt
echo "API_GATEWAY_PRIVATE_IP=$API_GATEWAY_PRIVATE_IP" >> variables.txt
echo "AUTH_SERVICE_PRIVATE_IP=$API_GATEWAY_PRIVATE_IP" >> variables.txt
echo "BASTION_PUBLIC_IP=$BASTION_PUBLIC_IP" >> variables.txt
echo "ALB_ARN=$ALB_ARN" >> variables.txt
echo "API_GATEWAY_TG_ARN=$API_GATEWAY_TG_ARN" >> variables.txt
echo "ALB_LISTENER_ARN=$ALB_LISTENER_ARN" >> variables.txt
echo "ALB_DNS_NAME=$ALB_DNS_NAME" >> variables.txt

echo "Variables saved in variables.txt file for cleanup process later."
