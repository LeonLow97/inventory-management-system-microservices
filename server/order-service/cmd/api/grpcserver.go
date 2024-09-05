package main

import (
	"fmt"
	"log"
	"net"

	"github.com/LeonLow97/internal/adapters/inbound/grpcserver"
	pb "github.com/LeonLow97/proto"
	"google.golang.org/grpc"
)

func (app application) InitiateGRPCServer() {
	lis, err := net.Listen("tcp", fmt.Sprintf(":%s", orderServicePort))
	if err != nil {
		log.Fatalf("Failed to start grpc server with error: %v\n", err)
	}

	// creates a new grpc server
	grpcServer := grpc.NewServer()
	orderGRPCServer := grpcserver.NewOrderGRPCServer(app.orderService)

	pb.RegisterOrderServiceServer(grpcServer, orderGRPCServer)
	log.Printf("Started order gRPC server at %v", lis.Addr())

	if err := grpcServer.Serve(lis); err != nil {
		log.Fatalf("Failed to start order gRPC server with error %v", err)
	}
}
