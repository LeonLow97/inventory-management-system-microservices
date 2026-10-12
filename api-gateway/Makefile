.PHONY: run
MODE ?= development
run:
	MODE=$(MODE) go run ./cmd/api

.PHONY: test
test:
	go test ./...

.PHONY: lint
lint:
	golangci-lint run ./...
