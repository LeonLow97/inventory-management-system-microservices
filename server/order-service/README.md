## Order Management Service

### Endpoints

| Method | Endpoint       | Description                          |
| :----: | -------------- | ------------------------------------ |
|  GET   | `/order`       | Retrieve a list of orders            |
|  GET   | `/orders/{id}` | Retrieve details of a specific order |
|  POST  | `/order`       | Create a new order                   |
|  PUT   | `/orders/{id}` | Update details of an order           |
| DELETE | `/orders/{id}` | Delete an order                      |

### gRPC CLI Commands

```
go install google.golang.org/protobuf/cmd/protoc-gen-go@v1.28
go install google.golang.org/grpc/cmd/protoc-gen-go-gr	pc@v1.2
export PATH="$PATH:$(go env GOPATH)/bin"
protoc --go_out=. --go-grpc_out=. proto/order.proto
protoc --go_out=. --go-grpc_out=. proto/inventory.proto
```

### Kafka CLI Commands

```
// Access docker container
docker-compose exec broker bash

// Create Topic
kafka-topics --create --topic test-topic --bootstrap-server localhost:9092 --replication-factor 1 --partitions 1

// List Topics
kafka-topics --list --bootstrap-server localhost:9092

// Write to Topic
kafka-console-producer --topic test-topic --bootstrap-server localhost:9092

// Read Topic
kafka-console-consumer --topic test-topic --from-beginning --bootstrap-server localhost:9092
```

### Folder Structure

- [Folder Structure](https://martengartner.medium.com/my-favourite-go-project-setup-479563662834)
