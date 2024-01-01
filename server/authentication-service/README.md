# Authentication Microservice

## Golang Packages

```
go get -u github.com/gorilla/mux
go get -u github.com/go-playground/validator/v10
go get github.com/jackc/pgx/v5
go get github.com/jackc/pgx/v5/stdlib
go get github.com/spf13/viper
go get github.com/golang-migrate/migrate/v4
```

## Commands for protobuf

```
go install google.golang.org/protobuf/cmd/protoc-gen-go@v1.28
go install google.golang.org/grpc/cmd/protoc-gen-go-grpc@v1.2
export PATH="$PATH:$(go env GOPATH)/bin"
protoc --go_out=. --go-grpc_out=. proto/authentication.proto
protoc --go_out=. --go-grpc_out=. proto/users.proto
```