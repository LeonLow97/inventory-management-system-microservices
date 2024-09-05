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
	lis, err := net.Listen("tcp", fmt.Sprintf(":%s", inventoryServicePort))
	if err != nil {
		log.Fatalf("Failed to start the grpc server with error: %v", err)
	}

	// creates a new grpc server
	server := grpc.NewServer()
	inventoryGRPCServer := grpcserver.NewInventoryGRPCServer(app.service)

	pb.RegisterInventoryServiceServer(server, inventoryGRPCServer)
	log.Printf("Started inventory gRPC server at %v", lis.Addr())

	if err := server.Serve(lis); err != nil {
		log.Fatalf("Failed to start the inventory gRPC server with error %v", err)
	}
}
