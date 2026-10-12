# API Gateway

The API Gateway exposes REST endpoints and forwards user operations to Authify Service over gRPC. It also handles HTTP request validation, JWT cookies, login rate limiting, and readiness checks.

## API Endpoints

| Method | Endpoint | Description |
|:---:|:---|:---|
| GET | `/health/liveness` | Confirms that the gateway process is running. |
| GET | `/health/readiness` | Checks whether Redis is reachable. |
| GET | `/metrics` | Exposes Prometheus metrics. |
| POST | `/api/v1/login` | Authenticates a user and sets the JWT cookie. |
| POST | `/api/v1/signup` | Creates a user account. |
| POST | `/api/v1/logout` | Expires the JWT cookie in the browser. |
| GET | `/api/v1/users` | Lists users. Requires a valid JWT cookie and an administrator account. |
| PATCH | `/api/v1/users/me` | Updates the authenticated user's supported fields. Requires a valid JWT cookie. |

All requests pass through the gateway's global request logging, panic recovery, and HTTP metrics middleware. The diagrams below focus on each endpoint's main flow.

## Endpoint Flows

### `POST /api/v1/login`

Login validates the JSON request, applies burst and failed-login limits using Redis, then calls Authify Service over gRPC. On success, the gateway sets the JWT cookie and returns the user response.

```mermaid
flowchart TD
    A[POST /api/v1/login] --> B[Global HTTP middleware]
    B --> C[Bind JSON and validate fields]
    C -->|invalid request| D[Return 400 Bad Request]
    C -->|valid request| E[Map request DTO to domain input]
    E --> F[Check burst and failed-login limits in Redis]
    F -->|limit exceeded| G[Return 429 Too Many Requests]
    F -.->|Redis error: log and continue| H[Use case calls outbound adapter]
    F -->|within limits| H
    H --> I[Map domain input to protobuf request]
    I --> J[Authify Service authenticates user]
    J -->|unauthenticated| K[Best-effort increment failed-login counter]
    K --> L[Map gRPC error to HTTP response]
    J -->|other error| L
    J -->|success| M[Best-effort reset failed-login counter]
    M --> N[Set JWT cookie and return 200 OK with JSON]
```

### `POST /api/v1/signup`

Signup validates the JSON request and forwards the account creation operation to Authify Service. This endpoint does not currently use the login rate limiter.

```mermaid
flowchart TD
    A[POST /api/v1/signup] --> B[Global HTTP middleware]
    B --> C[Bind JSON and validate fields]
    C -->|invalid request| D[Return 400 Bad Request]
    C -->|valid request| E[Map request DTO to domain input]
    E --> F[Call signup use case]
    F --> G[Map domain input to protobuf request]
    G --> H[Authify Service creates account]
    H -->|error| I[Map gRPC error to HTTP response]
    H -->|success| J[Return 201 Created]
```

### `POST /api/v1/logout`

Logout does not require JWT verification and does not call Authify Service. It expires the browser's cookie.

```mermaid
flowchart TD
    A[POST /api/v1/logout] --> B[Global HTTP middleware]
    B --> C[Set JWT cookie Max-Age to -1]
    C --> D[Return 200 OK]
```

### `GET /api/v1/users`

The authentication middleware verifies the JWT cookie before the handler runs. The authenticated user ID is passed to Authify Service as gRPC metadata, where the service checks administrator access and returns a page of users.

```mermaid
flowchart TD
    A[GET /api/v1/users] --> B[Global HTTP middleware]
    B --> C[Read and verify JWT cookie]
    C -->|missing or invalid| D[Return 401 Unauthorized]
    C -->|valid| E[Store user ID in request context]
    E --> F[Parse limit and cursor]
    F -->|invalid limit| G[Return 400 Bad Request]
    F -->|valid query| H[Call get-users use case]
    H --> I[Map request to protobuf and attach user ID metadata]
    I --> J[Authify Service checks admin access and lists users]
    J -->|error| K[Map gRPC error to HTTP response]
    J -->|success| L[Map users to JSON response]
    L --> M[Return 200 OK]
```

### `PATCH /api/v1/users/me`

The authentication middleware verifies the JWT cookie. The gateway validates the update body and forwards the update with the authenticated user ID in gRPC metadata.

```mermaid
flowchart TD
    A[PATCH /api/v1/users/me] --> B[Global HTTP middleware]
    B --> C[Read and verify JWT cookie]
    C -->|missing or invalid| D[Return 401 Unauthorized]
    C -->|valid| E[Store user ID in request context]
    E --> F[Bind JSON and validate fields]
    F -->|invalid request| G[Return 400 Bad Request]
    F -->|valid request| H[Map request DTO to domain input]
    H --> I[Call update-user use case]
    I --> J[Map input to protobuf and attach user ID metadata]
    J --> K[Authify Service updates the user]
    K -->|error| L[Map gRPC error to HTTP response]
    K -->|success| M[Return 204 No Content]
```
