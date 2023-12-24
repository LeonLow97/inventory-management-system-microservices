## API Gateway

- [x] Custom API Gateway
- [x] IP Whitelisting (using gin-gonic c.ClientIP() to retrieve IP address)
- [x] Rate Limiter (using `golang.org/x/time/rate` for rate limiter (Token Bucket Algorithm). limiting by client IP Address)
- [ ] Caching (Redis)
- [ ] Logging (ELK Stack)
- [ ] Monitoring (Prometheus, Grafana, Datadog)
- [x] gRPC Communication with Microservices

## Golang Framework and Resources

- `gin` golang framework
- `go get -u github.com/go-chi/chi/v5`
- [Write your own API Gateway](https://itnext.io/why-should-you-write-your-own-api-gateway-from-scratch-378074bfc49e)
- [Gin Framework](https://github.com/gin-gonic/gin)

## API Gateway

- **API Gateway**: acts as a single entry point for multiple backend services/APIs.
- **Function**: It manages the client's requests, handles authentication, traffic and security.
- **Key Features**
  - **Request Routing**: Directs incoming requests to the appropriate backend service.
  - **Protocol Translation**: Translates requests from one protocol to another if necessary. E.g., Translating from HTTP (Client request) to MQTT (Message Queuing Telemetry Transport backend)
  - **Load Balancing**: Distributes incoming requests across multiple servers to ensure optimal performance.
  - **Authentication and Authorization**: Validates user's identifies and controls access to APIs.
  - **Rate Limiting**: Controls the rate of incoming requests to prevent overwhelming the backend.
  - **Caching**: Stores frequently accessed data temporarily to improve response times and reduce the load on backend servers.
    - E.g., API Gateway can cache responses to GET requests for a specific duration, serving subsequent identical requests from the cache without involving the backend.
  - **Logging**: Records and stores information about incoming requests and their corresponding responses.
    - E.g., API Gateway logs can contain details such as timestamps, request headers, response codes, and backend service response times.
  - **Monitoring and Analytics**: Collects data on API usage, performance, and errors for analysis and optimization.
- **Benefits**:
  - **Centralized Management**: Simplifies API management by offering a unified interface for various services.
  - **Security**: Provides a secure layer to protect APIs from unauthorized access and attacks.
  - **Scalability**: Enables scaling of applications by managing traffic and distributing it efficiently among backend services.
  - **Analytics**: Offers insights into API usage patterns and performance metrics for better decision-making.
- **API Gateway Types**:
  - **Reverse Proxy**: Acts as an intermediary between clients and servers, receiving requests and forwarding them to the appropriate backend.
  - **Microservices API Gateway**: Designed specifically for managing microservices-based architectures, providing features to handle communication between microservices.
  - **Cloud-based API Gateway**: Offered as a service by cloud providers, allowing business to deploy and manage APIs without maintaining infrastructure. E.g., Amazon API Gateway, Google Cloud Endpoints, Microsoft Azure API Management
