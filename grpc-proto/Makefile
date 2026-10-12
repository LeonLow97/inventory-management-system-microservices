GO ?= go
PROTOC ?= protoc
PROTO_MODULE := github.com/LeonLow97/authify/grpc-proto
PROTO_FILE := proto/v1/users.proto
GO_BIN_DIR ?= $(shell $(GO) env GOPATH)/bin
PROTOC_INCLUDE ?= $(shell dirname "$$(command -v "$(PROTOC)")")/../include
PROTOC_GEN_GO_VERSION ?= v1.28.1
PROTOC_GEN_GO_GRPC_VERSION ?= v1.3.0

export PATH := $(GO_BIN_DIR):$(PATH)

.PHONY: all generate proto
all: generate

proto: generate

generate:
	@command -v "$(PROTOC)" >/dev/null 2>&1 || { echo "protoc is required; install the Protocol Buffers compiler and retry"; exit 1; }
	@test -f "$(PROTOC_INCLUDE)/google/protobuf/empty.proto" || { echo "Protocol Buffers well-known types were not found under $(PROTOC_INCLUDE); set PROTOC_INCLUDE to the directory containing google/protobuf"; exit 1; }
	GOBIN="$(GO_BIN_DIR)" $(GO) install google.golang.org/protobuf/cmd/protoc-gen-go@$(PROTOC_GEN_GO_VERSION)
	GOBIN="$(GO_BIN_DIR)" $(GO) install google.golang.org/grpc/cmd/protoc-gen-go-grpc@$(PROTOC_GEN_GO_GRPC_VERSION)
	$(PROTOC) --proto_path=. --proto_path="$(PROTOC_INCLUDE)" \
		--go_out=. --go_opt=module=$(PROTO_MODULE) \
		--go-grpc_out=. --go-grpc_opt=module=$(PROTO_MODULE) \
		$(PROTO_FILE)
