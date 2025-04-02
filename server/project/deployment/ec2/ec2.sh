#!/bin/bash

source .env

PUBLIC_IP=$(curl -s https://checkip.amazonaws.com)

if [ -f variables.txt ]; then
  > variables.txt
fi

# STEP 2: Create a VPC, Subnets, Internet Gateway and Security Groups
# Create VPC
echo "Creating VPC..."
VpcId=$(aws ec2 create-vpc --cidr-block 10.0.0.0/16 --region $REGION --query 'Vpc.VpcId' --output text)
echo "\tVPC created successfully: VpcId=$VpcId"

# Find Default Route Table ID (for cleanup later)
echo "Creating Default Route Table in VPC..."
DEFAULT_ROUTE_TABLE_ID=$(aws ec2 describe-route-tables --filters "Name=vpc-id,Values=$VpcId" --query 'RouteTables[?Associations[?Main==`true`]].RouteTableId' --output text)
echo "\tRoute Table created successfully: DEFAULT_ROUTE_TABLE_ID=$DEFAULT_ROUTE_TABLE_ID"

# Create 2 Public Subnets in VPC in different AZs for availability (for ALB to distribute incoming requests)
echo "Creating Private and Public Subnets in VPC..."
PUBLIC_SUBNET_AZ1_ID=$(aws ec2 create-subnet --vpc-id $VpcId --cidr-block 10.0.1.0/24 --availability-zone $AZ1 --query 'Subnet.SubnetId' --output text)
PUBLIC_SUBNET_AZ2_ID=$(aws ec2 create-subnet --vpc-id $VpcId --cidr-block 10.0.2.0/24 --availability-zone $AZ2 --query 'Subnet.SubnetId' --output text)
PRIVATE_SUBNET_AZ1_ID=$(aws ec2 create-subnet --vpc-id $VpcId --cidr-block 10.0.3.0/24 --availability-zone $AZ1 --query 'Subnet.SubnetId' --output text)
PRIVATE_SUBNET_AZ2_ID=$(aws ec2 create-subnet --vpc-id $VpcId --cidr-block 10.0.4.0/24 --availability-zone $AZ2 --query 'Subnet.SubnetId' --output text)
echo "\tAll Private and Public Subnets created."

# Create Internet Gateway (connects VPC to Internet) and attach it to VPC
echo "Creating Internet Gateway..."
IGW_ID=$(aws ec2 create-internet-gateway --region $REGION --query 'InternetGateway.InternetGatewayId' --output text)
aws ec2 attach-internet-gateway --internet-gateway-id $IGW_ID --vpc-id $VpcId
echo "\tSuccessfully created internet gateway."

# Create Route Table (connects public subnet to IGW) and associate with public subnet
echo "Creating Route Table..."
RTB_ID=$(aws ec2 create-route-table --vpc-id $VpcId --region $REGION --query 'RouteTable.RouteTableId' --output text)

# Adds a route in the route table that sends all traffic destined for the internet (0.0.0.0/0) to the Internet Gateway (IGW)
echo "\tCreating Route for Internet Gateway in Route Table..."
aws ec2 create-route --route-table-id $RTB_ID --destination-cidr-block 0.0.0.0/0 --gateway-id $IGW_ID --region $REGION >/dev/null 2>&1

# Associate the route table with a public subnet, enabling instances in that subnet to route internet-bound traffic to the IGW
echo "\tAssociated route table with public subnets..."
aws ec2 associate-route-table --route-table-id $RTB_ID --subnet-id $PUBLIC_SUBNET_AZ1_ID --region $REGION >/dev/null 2>&1
aws ec2 associate-route-table --route-table-id $RTB_ID --subnet-id $PUBLIC_SUBNET_AZ2_ID --region $REGION >/dev/null 2>&1
echo "\tSuccessfully created route table and associated routes."

# CREATE SECURITY GROUPS
echo "Creating Security Groups..."
# SECURITY GROUP FOR ALB
ALB_SG_ID=$(aws ec2 create-security-group --group-name ALB-SG \
  --description "Security group for ALB" --vpc-id $VpcId --region $REGION \
  --query 'GroupId' --output text)
# SECURITY GROUP FOR API GATEWAY
API_GATEWAY_SG_ID=$(aws ec2 create-security-group --group-name API-Gateway-SG \
  --description "Security group for API Gateway" --vpc-id $VpcId --region $REGION \
  --query 'GroupId' --output text)
# SECURITY GROUP FOR AUTHENTICATION MICROSERVICE
AUTH_SG_ID=$(aws ec2 create-security-group --group-name Auth-Service-SG \
  --description "Security group for authentication microservice" --vpc-id $VpcId --region $REGION \
  --query 'GroupId' --output text)
# SECURITY GROUP FOR BASTION HOST
BASTION_SG_ID=$(aws ec2 create-security-group --group-name Bastion-SG \
  --description "Security group for Bastion Host" --vpc-id $VpcId \
  --query 'GroupId' --output text)
echo "\tSuccessfully created all security groups."

# SETTING INBOUND RULES FOR SECURITY GROUPS
echo "Setting Inbound Rules for Security Groups..."
# INBOUND RULE FOR ALB
aws ec2 authorize-security-group-ingress --group-id $ALB_SG_ID --region $REGION \
  --ip-permissions '[
    {"IpProtocol": "tcp", "FromPort": 80, "ToPort": 80, "IpRanges": [{"CidrIp": "0.0.0.0/0"}]},
    {"IpProtocol": "tcp", "FromPort": 443, "ToPort": 443, "IpRanges": [{"CidrIp": "0.0.0.0/0"}]}
  ]' >/dev/null 2>&1 &
# INBOUND RULE FOR API GATEWAY
aws ec2 authorize-security-group-ingress --group-id $API_GATEWAY_SG_ID \
  --ip-permissions '[
    {"IpProtocol": "tcp", "FromPort": '"$PORT_API_GATEWAY"', "ToPort": '"$PORT_API_GATEWAY"', "UserIdGroupPairs": [{"GroupId": "'"$ALB_SG_ID"'"}]},
    {"IpProtocol": "tcp", "FromPort": 22, "ToPort": 22, "UserIdGroupPairs": [{"GroupId": "'"$BASTION_SG_ID"'"}]},
    {"IpProtocol": "tcp", "FromPort": 80, "ToPort": 80, "UserIdGroupPairs": [{"GroupId": "'"$BASTION_SG_ID"'"}]},
    {"IpProtocol": "tcp", "FromPort": 443, "ToPort": 443, "UserIdGroupPairs": [{"GroupId": "'"$BASTION_SG_ID"'"}]}
  ]' >/dev/null 2>&1 &
# INBOUND RULE FOR AUTHENTICATION MICROSERVICE
aws ec2 authorize-security-group-ingress --group-id $AUTH_SG_ID \
  --ip-permissions '[
    {"IpProtocol": "tcp", "FromPort": '"$PORT_AUTH_SERVICE"',"ToPort": '"$PORT_AUTH_SERVICE"',"UserIdGroupPairs": [{"GroupId": "'"$API_GATEWAY_SG_ID"'"}]},
    {"IpProtocol": "tcp", "FromPort": '"$PORT_POSTGRESQL"',"ToPort": '"$PORT_POSTGRESQL"',"UserIdGroupPairs": [{"GroupId": "'"$AUTH_SG_ID"'"}]}
  ]' >/dev/null 2>&1 &
# INBOUND RULE FOR BASTION HOST 
aws ec2 authorize-security-group-ingress --group-id $BASTION_SG_ID \
  --protocol tcp --port 22 --cidr $PUBLIC_IP/32 >/dev/null 2>&1 &
wait
echo "\tSuccessfully set inbound rules for all security groups."

# CREATING EC2 INSTANCES KEY PAIRS FOR SSH
echo "Creating Key Pairs for EC2 Instances..."
# KEY PAIR FOR API GATEWAY
aws ec2 create-key-pair --key-name $EC2_API_GATEWAY_KEY_NAME --query 'KeyMaterial' --output text > $EC2_API_GATEWAY_KEY_NAME.pem &
# KEY PAIR FOR AUTHENTICATION MICROSERVICE
aws ec2 create-key-pair --key-name $EC2_AUTH_MICROSERVICE_KEY_NAME --query 'KeyMaterial' --output text > $EC2_AUTH_MICROSERVICE_KEY_NAME.pem &
# KEY PAIR FOR BASTION HOST
aws ec2 create-key-pair --key-name $EC2_BASTION_HOST_KEY_NAME --query 'KeyMaterial' --output text > $EC2_BASTION_HOST_KEY_NAME.pem &
wait
echo "\tSuccessfully created key pairs for EC2 instances."

# CREATE NAT GATEWAY
# Allocate two Elastic IPs (one for each NAT Gateway)
echo "Allocating address for Elastic IPs..."
EIP_AZ1=$(aws ec2 allocate-address --domain vpc --query 'AllocationId' --output text)
EIP_AZ2=$(aws ec2 allocate-address --domain vpc --query 'AllocationId' --output text)
echo "\tSuccessfully allocated address for Elastic IPs."

# Create NAT Gateways in both public subnets with different Elastic IPs
echo "Creating NAT Gateways..."
NAT_GATEWAY_ID_AZ1=$(aws ec2 create-nat-gateway --subnet-id $PUBLIC_SUBNET_AZ1_ID --allocation-id $EIP_AZ1 --query "NatGateway.NatGatewayId" --output text)
NAT_GATEWAY_ID_AZ2=$(aws ec2 create-nat-gateway --subnet-id $PUBLIC_SUBNET_AZ2_ID --allocation-id $EIP_AZ2 --query "NatGateway.NatGatewayId" --output text)
echo "\tSuccessfully created NAT Gateways."

# Wait for both NAT Gateways to become available
until [ "$(aws ec2 describe-nat-gateways --nat-gateway-ids $NAT_GATEWAY_ID_AZ1 --query "NatGateways[0].State" --output text)" == "available" ]; do
    sleep 10
done
until [ "$(aws ec2 describe-nat-gateways --nat-gateway-ids $NAT_GATEWAY_ID_AZ2 --query "NatGateways[0].State" --output text)" == "available" ]; do
    sleep 10
done
echo "Both NAT Gateways are available. Proceeding with the next step."

# Update Private Subnet Route Tables to use their respective NAT Gateway
aws ec2 create-route --route-table-id $DEFAULT_ROUTE_TABLE_ID --destination-cidr-block 0.0.0.0/0 --nat-gateway-id $NAT_GATEWAY_ID_AZ1 >/dev/null 2>&1
aws ec2 create-route --route-table-id $DEFAULT_ROUTE_TABLE_ID --destination-cidr-block 0.0.0.0/0 --nat-gateway-id $NAT_GATEWAY_ID_AZ2 >/dev/null 2>&1

# LAUNCHING EC2 INSTANCES
# API GATEWAY
API_GATEWAY_INSTANCE_ID=$(aws ec2 run-instances \
  --image-id ami-06661384e66f2da0e \
  --count 1 \
  --instance-type t3.micro \
  --key-name $EC2_API_GATEWAY_KEY_NAME \
  --security-group-ids $API_GATEWAY_SG_ID \
  --subnet-id $PRIVATE_SUBNET_AZ1_ID \
  --user-data file://deployment/ec2/api-gateway-script.sh \
  --query 'Instances[0].InstanceId' \
  --output text)
# AUTHENTICATION MICROSERVICE
AUTH_SERVICE_INSTANCE_ID=$(aws ec2 run-instances \
  --image-id ami-06661384e66f2da0e \
  --count 1 \
  --instance-type t3.micro \
  --key-name $EC2_AUTH_MICROSERVICE_KEY_NAME \
  --security-group-ids $AUTH_SG_ID \
  --subnet-id $PRIVATE_SUBNET_AZ1_ID \
  --user-data file://deployment/ec2/authentication-script.sh \
  --query 'Instances[0].InstanceId' \
  --output text)
# BASTION HOST
BASTION_INSTANCE_ID=$(aws ec2 run-instances \
  --image-id ami-06661384e66f2da0e \
  --count 1 \
  --instance-type t3.micro \
  --key-name $EC2_BASTION_HOST_KEY_NAME \
  --security-group-ids $BASTION_SG_ID \
  --subnet-id $PUBLIC_SUBNET_AZ1_ID \
  --associate-public-ip-address \
  --query 'Instances[0].InstanceId' \
  --output text)

echo "API_GATEWAY_INSTANCE_ID=$API_GATEWAY_INSTANCE_ID"
echo "AUTH_SERVICE_INSTANCE_ID=$AUTH_SERVICE_INSTANCE_ID"
echo "BASTION_INSTANCE_ID=$BASTION_INSTANCE_ID"

# Attach IAM Role (IMS-EC2-Role) to EC2 Instance so they can pull ECR images
aws ec2 associate-iam-instance-profile --instance-id $API_GATEWAY_INSTANCE_ID \
    --iam-instance-profile Name=IMS-EC2-Role >/dev/null 2>&1
aws ec2 associate-iam-instance-profile --instance-id $AUTH_SERVICE_INSTANCE_ID \
    --iam-instance-profile Name=IMS-EC2-Role >/dev/null 2>&1

# Wait until EC2 instance is running and check the status of the EC2 instance
aws ec2 wait instance-running --instance-ids $API_GATEWAY_INSTANCE_ID
aws ec2 wait instance-running --instance-ids $AUTH_SERVICE_INSTANCE_ID
aws ec2 wait instance-status-ok --instance-ids $API_GATEWAY_INSTANCE_ID
aws ec2 wait instance-status-ok --instance-ids $AUTH_SERVICE_INSTANCE_ID

# Get the private IPv4 address of the newly launched EC2 instance
API_GATEWAY_PRIVATE_IPV4=$(aws ec2 describe-instances \
  --instance-ids $API_GATEWAY_INSTANCE_ID \
  --query 'Reservations[0].Instances[0].PrivateIpAddress' \
  --output text)
echo "API_GATEWAY_PRIVATE_IPV4=$API_GATEWAY_PRIVATE_IPV4"

BASTION_PUBLIC_IP=$(aws ec2 describe-instances \
  --instance-ids $BASTION_INSTANCE_ID \
  --query 'Reservations[0].Instances[0].PublicIpAddress' --output text)
echo "BASTION_PUBLIC_IP=$BASTION_PUBLIC_IP"

echo "Run the following commands in order to SSH into API Gateway EC2 Instance"
echo "ssh -i IMS_BASTION.pem ec2-user@$BASTION_PUBLIC_IP"
echo "scp -i IMS_BASTION.pem IMS_API_GATEWAY_KEY_PAIR.pem ec2-user@$BASTION_PUBLIC_IP:~"
echo "ssh -i IMS_BASTION.pem ec2-user@$BASTION_PUBLIC_IP"
echo "ssh -i IMS_API_GATEWAY_KEY_PAIR.pem ec2-user@$API_GATEWAY_PRIVATE_IPV4"

# Step 6: Set up ALB to route traffic
# Create ALB --> Set up a public Application Load Balancer in a specified subnet
ALB_ARN=$(aws elbv2 create-load-balancer --name ims-alb --subnets $PUBLIC_SUBNET_AZ1_ID $PUBLIC_SUBNET_AZ2_ID --security-groups $ALB_SG_ID --query "LoadBalancers[0].LoadBalancerArn" --output text)
echo "ALB_ARN=$ALB_ARN"

# Create Target Group --> Defines where ALB should forward traffic (API Gateway on Port 8080)
API_GATEWAY_TG_ARN=$(aws elbv2 create-target-group --name "$API_GATEWAY_SERVICE-tg" --protocol HTTP --port 8080 --vpc-id $VpcId --query "TargetGroups[0].TargetGroupArn" --output text)
echo "API_GATEWAY_TG_ARN=$API_GATEWAY_TG_ARN"

# Modify API Gateway target group health check endpoint
aws elbv2 modify-target-group \
    --target-group-arn $API_GATEWAY_TG_ARN \
    --health-check-path "/healthcheck" \
    --health-check-port "8080" \
    --health-check-protocol "HTTP" \
    --health-check-interval-seconds 30 \
    --health-check-timeout-seconds 5 \
    --healthy-threshold-count 5 \
    --unhealthy-threshold-count 2 >/dev/null 2>&1

# Register Target --> Links the API Gateway EC2 Instance to the target group
echo "Waiting for API Gateway EC2 instance $API_GATEWAY_INSTANCE_ID to be running..."
aws ec2 wait instance-running --instance-ids $API_GATEWAY_INSTANCE_ID
aws elbv2 register-targets --target-group-arn $API_GATEWAY_TG_ARN --targets Id=$API_GATEWAY_INSTANCE_ID

# Create Listener --> Configured ALB to listen on port 80 and forward requests to Target Group
ALB_LISTENER_ARN=$(aws elbv2 create-listener --load-balancer-arn $ALB_ARN --protocol HTTP --port 80 --default-actions Type=forward,TargetGroupArn=$API_GATEWAY_TG_ARN --query "Listeners[0].ListenerArn" --output text)
echo "ALB_LISTENER_ARN=$ALB_LISTENER_ARN"

# Wait for the ALB to be in available state
echo "Waiting for the Load Balancer $ALB_ARN to be available..."
aws elbv2 wait load-balancer-available --load-balancer-arn $ALB_ARN

# Get the DNS name of the Load Balancer
ALB_DNS_NAME=$(aws elbv2 describe-load-balancers --load-balancer-arn $ALB_ARN --query "LoadBalancers[0].DNSName" --output text)
echo "ALB_DNS_NAME=$ALB_DNS_NAME"

# Test API Gateway HealthCheck endpoint
echo "curl http://$ALB_DNS_NAME/healthcheck"
curl http://$ALB_DNS_NAME/healthcheck

# Add variables to variables.txt
echo "VpcId=$VpcId" >> variables.txt
echo "DEFAULT_ROUTE_TABLE_ID=$DEFAULT_ROUTE_TABLE_ID" >> variables.txt
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
echo "API_GATEWAY_PRIVATE_IPV4=$API_GATEWAY_PRIVATE_IPV4" >> variables.txt
echo "BASTION_PUBLIC_IP=$BASTION_PUBLIC_IP" >> variables.txt
echo "ALB_ARN=$ALB_ARN" >> variables.txt
echo "API_GATEWAY_TG_ARN=$API_GATEWAY_TG_ARN" >> variables.txt
echo "ALB_LISTENER_ARN=$ALB_LISTENER_ARN" >> variables.txt
echo "ALB_DNS_NAME=$ALB_DNS_NAME" >> variables.txt
