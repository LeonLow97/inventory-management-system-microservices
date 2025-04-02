#!/bin/bash
set -e  # Exit if any command fails

# Update system and install Docker
sudo yum update -y
sudo yum install -y docker
sudo systemctl start docker
sudo systemctl enable docker

# Add EC2 user to Docker group (avoids permission issues)
sudo usermod -aG docker ec2-user

aws ecr get-login-password --region ap-southeast-1 | sudo docker login --username AWS --password-stdin 931100716322.dkr.ecr.ap-southeast-1.amazonaws.com

sudo docker pull 931100716322.dkr.ecr.ap-southeast-1.amazonaws.com/ims-repository:api-gateway-latest
sudo docker run -d --name api-gateway -p 8080:8080 931100716322.dkr.ecr.ap-southeast-1.amazonaws.com/ims-repository:api-gateway-latest
