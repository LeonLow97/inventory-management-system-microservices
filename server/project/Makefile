build:
	docker-compose -f docker-compose.dev.yml down
	@echo "Building and starting docker images..."
	docker-compose -f docker-compose.dev.yml up --build -d
	@echo "Docker images built and started!"

down:
	docker-compose -f docker-compose.dev.yml down

ec2_deploy:
	sh ./deployment/common/ecr.sh
	# sh ./deployment/common/rds.sh
	sh ./deployment/ec2/ec2.sh
	@echo "Successfully Deployed IMS to AWS EC2!"

ec2_clean:
	sh ./deployment/ec2/ec2-cleanup.sh
	# sh ./deployment/common/rds-cleanup.sh
	sh ./deployment/common/ecr-cleanup.sh
	@echo "Cleanup Completed!"

ecr_push:
	sh ./deployment/common/ecr.sh

ecr_clean:
	sh ./deployment/common/ecr-cleanup.sh
