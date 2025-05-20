# Project Current State

- Only API Gateway and Authentication Microservice is working.
- Previously, all services (authentication, order and inventory) are working but currently trying to deploy 2 services to AWS so I shut down the network communication to order and inventory microservices.
- Will be revamping the entire Microservices Architecture and deploy to AWS.

# Inventory Management System

## Table of Contents

- [Overview](#overview)
- [Technologies](#technologies)
- [Microservices Architecture](#microservices-architecture)
- [Project Setup](#project-setup)
  - [Docker](#docker)
  - [Kubernetes](#kubernetes)
  - [AWS Deployment](#aws-deployment)

## Overview

The Inventory Management System (IMS) is designed to assist individuals and business owners in efficiently tracking their inventory, whether for a physical store or an online shop. This system ensures that you can effortlessly manage your stock levels and make informed decisions.

When a customer places an order, the inventory count is updated once the order has been processed, providing real-time visibility into available stock. The system also alerts you when inventory levels are low, enabling timely reordering of items. This proactive approach ensures that you always have sufficient stock on hand to meet customer demand.

In addition to inventory tracking, the system records customer details during each purchase, along with cost and profit for each item sold. This comprehensive data allows you to analyze product performance by brand and category, giving you deeper insights into your inventory management or to perform other business analysis. Ultimately, this system empowers you to maintain optimal inventory levels and enhances your overall business efficiency.

This project was originally created to support my online business at [YourHypeStore](https://www.carousell.sg/u/yourhypestore)

## Technologies

| Technology   | Type                   | Version | Ports  |
| ------------ | ---------------------- | :-----: | ------ |
| Golang       | Language               |         |        |
| Gin          | Framework              |         |        |
| Gorilla Mux  | Framework              |         |        |
| PostgreSQL   | Database               |         | `5432` |
| MySQL        | Database               |         | `3306` |
| Apache Kafka | Message Queue          |         |        |
| Docker       | Containerization       |         |        |
| Kubernetes   | Orchestration          |         |        |
| REST         | Communication Protocol |         |        |
| gRPC         | Communication Protocol |         |        |

## Microservices Architecture

1. **API Gateway**

- An API Gateway is a reverse proxy that acts as a single entry point for all incoming client requests to backend services.
- In IMS, the API Gateway accepts REST-based requests (JSON format) and forwards them to microservices via gRPC (serialized binary format) for better performance and lower overhead.
- The **Aggregator Pattern** is utilized to collect data from multiple microservices and return a single aggregated response, reducing the number of network calls and improving response times.
- The API Gateway in IMS implements the following key features:
  - **Request Routing**: Directs client requests to the appropriate microservice.
  - **Rate Limiting & Throttling**: Token Bucket algorithm prevents DDoS attacks by controlling requests rates.
  - **IP Whitelisting**: Perform IP whitelisting for admin user routes.

2. **Authentication Microservice**

3. **Inventory Microservice**

4. **Order Microservice**

## Project Setup

- Ensure these ports are available on your localhost machine as stated in [this section](#technologies). You can also modify the host port if needed.

### Docker

1. Navigate to the `server/project` directory:

```bash
cd ./server/project
```

2. Run `init.sh` to setup `.proto` files and install gRPC. These `.proto` files are compiled using `protoc` which is the protocol buffers compiler, to generate code in Golang that enables efficient serialization and deserialization in gRPC communication.

```bash
sh init.sh
```

3. Start the Docker containers by pulling images from Docker Hub registry, creating the containers and running them on localhost:

```bash
docker-compose up --build -d
```

OR (if using Makefile)

```bash
make build
```

### Kubernetes

1. Create directories for persistent volumes on your local machine for the following microservices: Authentication, Inventory and Order. The persistent volumes will be created in the user's home directory:

```bash
mkdir -p $(echo $HOME)/persistent_volume/inventory-mysql $(echo $HOME)/persistent_volume/order-postgres $(echo $HOME)/persistent_volume/authentication-postgres
```

2. Navigate to and open the Kubernetes Persistent Volume (PV) config file `persistent-volume.yml`

```bash
cd ./server/project/k8s/
vi persistent-volume.yml
```

3. Edit the `hostPath` for all 3 persistent volumes - `inventory-mysql`, `order-postgres`, `authentication-postgres`

- For my MacOS, the path is `/Users/leonlow`.
- For your machine, run `echo $HOME` to get your home directory path, and use the result to update the `hostPath`.

```bash
hostPath:
  path: /Users/leonlow/persistent_volume/inventory-mysql
```

4. Run the following bash script `k8s.sh`. The script automates the setup of a **local Docker Registry** (TODO: probably should use Docker Hub Registry) and pushes container images for the IMS microservices (API Gateway, Authentication, Order and Inventory) to it. Then, it configures a **Kubernetes namespace** (`inventory-management-system`), sets the context, and verifies the namespace switch. Finally, it **deploys an NGINX Ingress Controller** to manage external access to the microservices within the Kubernetes Cluster.

```bash
sh k8s.sh
```

5. Deploy all Kubernetes Configuration files (YAML) which includes Ingress, ClusterIP, Deployments, Services, ConfigMaps, Secrets, StatefulSets, StorageClass.

```bash
kubectl apply -f ./k8s
```

Might throw an error because ingress resources takes a while to start as Kubernetes takes time to schedule and start the NGINX Ingress Controller pods. Wait for a few minutes and run the following command:

```sh
kubectl apply -f ./k8s/ingress-resource.yml
```

6. Ensure all pods are up and running:

```
kubectl get pods
```

7. Call the HealthCheck endpoint on API Gateway to ensure service is up and running:

```sh
curl http://localhost:80/healthcheck
```

# AWS Deployment

## AWS Elastic Beanstalk

- For deployment steps, check out [IMS AWS Elastic Beanstalk Deployment Guide](./server/project/docs/AWS/aws-elastic-beanstalk/deployment.md)
- For deployment steps, check out [IMS AWS EC2 Deployment Guide](./server/project/docs/AWS/aws-ec2/deployment.md)
