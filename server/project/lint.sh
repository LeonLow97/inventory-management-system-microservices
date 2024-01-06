echo "linting api-gateway service"
cd ../api-gateway
gofmt -w .
go mod tidy
grep -rnw './' -e 'fmt.Println'

echo "linting authentication microservice"
cd ../authentication-service
gofmt -w .
go mod tidy
grep -rnw './' -e 'fmt.Println'

echo "linting inventory microservice"
cd ../inventory-service
gofmt -w .
go mod tidy
grep -rnw './' -e 'fmt.Println'

echo "linting order microservice"
cd ../order-service
gofmt -w .
go mod tidy
grep -rnw './' -e 'fmt.Println'
