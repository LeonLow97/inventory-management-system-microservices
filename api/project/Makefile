AUTHENTICATION_BINARY=authApp
APIGATEWAY_BINARY=apiGateway

up_build: build_authentication build_api_gateway
	docker-compose down
	@echo "Building and starting docker images..."
	docker-compose up --build -d
	@echo "Docker images built and started!"

build_authentication:
	@echo "Build authentication binary..."
	cd ../authentication-service && env GOOS=linux CGO_ENABLED=0 go build -o ${AUTHENTICATION_BINARY} ./cmd/api

build_api_gateway:
	@echo "Build api gateway binary..."
	cd ../api-gateway && env GOOS=linux CGO_ENABLED=0 go build -o ${APIGATEWAY_BINARY} .

down:
	docker-compose down
