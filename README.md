# Inventory Management System

## Overview

The Inventory Management System (IMS) is designed to assist individuals and business owners in efficiently tracking their inventory, whether for a physical store or an online shop. This system ensures that you can effortlessly manage your stock levels and make informed decisions.

When a customer places an order, the inventory count is updated, providing real-time visibility into available stock. The system also alerts you when inventory levels are low, enabling timely reordering of items. This proactive approach ensures that you always have sufficient stock on hand to meet customer demand.

In addition to inventory tracking, the system records customer details during each purchase, along with cost and profit for each item sold. This comprehensive data allows you to analyze product performance by brand and category, giving you deeper insights into your inventory management. Ultimately, this system empowers you to maintain optimal inventory levels and enhances your overall business efficiency.

This project was originally created to support my online business at [YourHypeStore](https://www.carousell.sg/u/yourhypestore)

## Architecture

### Microservices

1. **API Gateway**

- An API Gateway is an API management tool that sits between a client and a collection of backend services.
- Accepts HTTP requests from client and routes them to respective microservices with gRPC communication.
- Utilised API Gateway **Aggregator Pattern** to collect data from various microservices and returns an aggregate for processing.
- Implemented **rate limiting with Leaky Bucket Algorithm** to avoid overloading a server from too many requests at the same time and prevent Denial of Service (DoS) attacks.

2. **Authentication**

- Implements JWT-based authentication for secure user access.
- Validates user credentials and issues access tokens.
- Handles user registration and management.

3. **Inventory**

- Manages inventory data and performs operations like:
  - Adding new inventory items
  - Updating existing items
  - Retrieving inventory status

4. **Order**

- Processes customer orders and tracks their status
- Interacts with the Inventory service to check stock levels.

## Technologies Used

- **Languages**: Golang
- **Frameworks**: Gin, Gorilla Mux
- **Databases**: PostgreSQL, MySQL
- **Messaging**: Apache Kafka
- **Containerization and Orchestration**: Docker, Kubernetes
- **Communication Protocols**: REST, gRPC

## Setup

There are 2 methods to setup IMS - Docker and Kubernetes, and this is only applicable to localhost environment.

---

### Docker Setup

- Ensure ports 5432 and 3306 are not taken as there will be used by PostgreSQL and MySQL respectively. You could also change the host port to another port number to accommodate the setup.

```
# Step 1: change directory into `project` folder
cd ./server/project

# Step 2: run init.sh to setup proto files and install grpc
sh init.sh

# Step 3: Either run the command in `Makefile` or manual docker-compose command
make build
docker-compose up # alternative command
```

### Kubernetes Setup

```
# Step 1: Create directories for persistent volumes in Authentication, Inventory and Order microservices, we are storing PVs in local directories (on production, we will use AWS elastic block store). On MacOS, go to your user home directory (cd ~, then perform ls to find our user directory).
cd /Users/leonlow/
mkdir persistent_volume
cd persistent_volume
mkdir inventory-mysql order-postgres authentication-postgres

# Step 2: Edit the hostPath for PV in k8s config file, `persistent-volume.yml`
cd ./server/project/k8s/
vi persistent-volume.yml

----- EDIT THIS LINE -----
# change `leonlow` to your own `home directory`, make sure to do for all 3 PV
hostPath:
  path: /Users/leonlow/persistent_volume/inventory-mysql

# Step 3: Run k8s.sh to set up local Docker Registry (in reality, we will use AWS ECR or Docker Hub)
sh k8s.sh

# Step 4: After Step 3 has completed (important to wait till it's done), setup k8s resources (this step will take about 5 - 10 minutes).
kubectl apply -f ./k8s

# If step 4 throws a warning saying that ingress controller still not started, wait for a few minutes then run
kubectl apply -f ./k8s/ingress-resource.yml

# Step 5: Check and ensure all pods are running
kubectl get pods

# Step 6: Perform healthcheck on API Gateway
curl http://localhost:80/healthcheck
```

---
