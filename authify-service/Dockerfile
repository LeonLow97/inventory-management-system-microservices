# ---- Builder Stage ----
FROM golang:1.27.1-alpine AS builder

WORKDIR /workspace

COPY grpc-proto/go.mod ./grpc-proto/go.mod
COPY authify-service/go.mod authify-service/go.sum ./authify-service/
WORKDIR /workspace/authify-service
RUN go mod download

WORKDIR /workspace
COPY grpc-proto/ ./grpc-proto/
COPY authify-service/ ./authify-service/
WORKDIR /workspace/authify-service

ENV CGO_ENABLED=0 GOOS=linux GOARCH=amd64
RUN go build -o authifyService ./cmd

# ---- Lightweight Image ----
FROM alpine:3.22.2

WORKDIR /app

COPY --from=builder /workspace/authify-service/authifyService .
COPY --from=builder /workspace/authify-service/config ./config

EXPOSE 50051

CMD ["./authifyService"]
