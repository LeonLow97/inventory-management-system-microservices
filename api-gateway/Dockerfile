# ---- Builder Stage ----
FROM golang:1.27.1-alpine AS builder

WORKDIR /workspace

COPY grpc-proto/go.mod ./grpc-proto/go.mod
COPY api-gateway/go.mod api-gateway/go.sum ./api-gateway/
WORKDIR /workspace/api-gateway
RUN go mod download

WORKDIR /workspace
COPY grpc-proto/ ./grpc-proto/
COPY api-gateway/ ./api-gateway/
WORKDIR /workspace/api-gateway

ENV CGO_ENABLED=0 GOOS=linux GOARCH=amd64
RUN go build -o apiGateway ./cmd/api

# ---- Lightweight Image ----
FROM alpine:3.22.2

WORKDIR /app

COPY --from=builder /workspace/api-gateway/apiGateway .
COPY --from=builder /workspace/api-gateway/config ./config

EXPOSE 8080

CMD ["./apiGateway"]
