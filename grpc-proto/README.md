# Shared gRPC contracts

`proto/v1/users.proto` is the canonical v1 gRPC contract. Generated Go client and server code is written to `pb/` and imported by both `api-gateway` and `authify-service` from `github.com/LeonLow97/authify/grpc-proto/pb`.

The services use the module path `github.com/LeonLow97/authify/grpc-proto`; their `go.mod` files replace it with this local folder for repository builds. Update that path if the repository's canonical import path differs before publishing the module. Install `protoc`, then run `make -C grpc-proto` from the repository root to generate the Go files. The Makefile installs the pinned Go plugins and writes generated files into `pb/`. If the well-known protobuf definitions are outside the compiler's adjacent `include` directory, pass their parent directory with `PROTOC_INCLUDE`. `make -C startup proto` delegates to the same generator.

Keep future incompatible contract versions under a new source directory such as `proto/v2`. Each version should use a versioned protobuf package and should be generated into its own Go package when both versions need to coexist.
