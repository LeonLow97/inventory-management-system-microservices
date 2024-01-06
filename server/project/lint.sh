echo "linting api-gateway service"
cd ../api-gateway
gofmt -w .

echo "linting authentication microservice"
cd ../authentication-service
gofmt -w .

echo "linting inventory microservice"
cd ../inventory-service
gofmt -w .

echo "linting order microservice"
cd ../order-service
gofmt -w .