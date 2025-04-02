build:
	docker-compose -f docker-compose.dev.yml down
	@echo "Building and starting docker images..."
	docker-compose -f docker-compose.dev.yml up --build -d
	@echo "Docker images built and started!"

down:
	docker-compose down -f docker-compose.dev.yml

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
