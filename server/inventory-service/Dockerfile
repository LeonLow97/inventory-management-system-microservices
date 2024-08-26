FROM golang:alpine

WORKDIR /app

COPY go.mod .
COPY go.sum .

RUN go mod download

COPY . .

ENV GOOS=linux CGO_ENABLED=0

RUN go build -o inventoryService ./cmd/api

EXPOSE 8001

CMD ["./inventoryService"]
