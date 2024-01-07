## Inventory Service

### Endpoints

| Method | Endpoint                          | Description                              |
| :----: | --------------------------------- | ---------------------------------------- |
|  GET   | `/inventory/products`             | Retrieve a list of products              |
|  GET   | `/inventory/products/{id}`        | Retrieve details of a specific product   |
|  POST  | `/inventory/products`             | Create a new product                     |
|  PUT   | `/inventory/products/{id}`        | Update details of a product              |
| DELETE | `/inventory/products/{id}`        | Delete a product (Soft Delete)           |
|  POST  | `/inventory/products/{id}/adjust` | Adjust the inventory count for a product |

### Database Tables

- **Products Table**
  - Stores the product information including name, brand, size, color, quantity, category, and potentially additional details.
- **Brands Table**
  - Contains brand details associated with products. Each product can have a relationship with a brand through a foreign key reference.
- **Categories Table**
  - Holds category information for products. This table establishes the category relationship for each product.

### gRPC Protocol Buffer Compiler Commands

```
go install google.golang.org/protobuf/cmd/protoc-gen-go@v1.28
go install google.golang.org/grpc/cmd/protoc-gen-go-grpc@v1.2
export PATH="$PATH:$(go env GOPATH)/bin"
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

### Access MySQL in Docker Container command

```
// retrieve container name
docker ps

docker exec -it project-inventory-mysql-1 bash

mysql -u inventory-mysql -p password
```