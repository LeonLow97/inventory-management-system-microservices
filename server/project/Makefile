build:
	docker-compose -f docker-compose.dev.yml down
	@echo "Building and starting docker images..."
	docker-compose -f docker-compose.dev.yml up --build -d
	@echo "Docker images built and started!"

down:
	docker-compose down -f docker-compose.dev.yml
