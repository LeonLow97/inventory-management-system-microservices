APIGATEWAY_BINARY=apiGateway
AUTHENTICATION_BINARY=authApp
INVENTORY_BINARY=inventoryApp
ORDER_BINARY=orderApp

up_build: build_api_gateway build_authentication build_inventory build_order
	docker-compose down
	@echo "Building and starting docker images..."
	docker-compose up --build -d
	@echo "Docker images built and started!"

build_api_gateway:
	@echo "Build api gateway binary..."
	cd ../api-gateway && env GOOS=linux CGO_ENABLED=0 go build -o ${APIGATEWAY_BINARY} .

build_authentication:
	@echo "Build authentication binary..."
	cd ../authentication-service && env GOOS=linux CGO_ENABLED=0 go build -o ${AUTHENTICATION_BINARY} ./cmd/api

build_inventory:
	@echo "Build inventory binary..."
	cd ../inventory-service && env GOOS=linux CGO_ENABLED=0 go build -o ${INVENTORY_BINARY} ./cmd/api

build_order:
	@echo "Build order binary..."
	cd ../order-service && env GOOS=linux CGO_ENABLED=0 go build -o ${ORDER_BINARY} ./cmd/api

down:
	docker-compose down
