#!/bin/bash

source .env

# Remove local Docker images
docker rmi $ECR_AUTHENTICATION_SERVICE_IMAGE_NAME $ECR_API_GATEWAY_IMAGE_NAME

# Delete the ECR repository (if no longer needed)
aws ecr delete-repository --repository-name $ECR_REPOSITORY_NAME --region $REGION --force >/dev/null 2>&1
