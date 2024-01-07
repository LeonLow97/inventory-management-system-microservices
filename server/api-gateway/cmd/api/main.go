package main

import (
	"fmt"
	"log"

	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

const gatewayPort = "80"

type application struct {
}

type grpcClientConn struct {
	orderConn     *grpc.ClientConn
	inventoryConn *grpc.ClientConn
	authConn      *grpc.ClientConn
}

func main() {
	// setup application config
	app := application{}

	grpcClients := app.initiateGRPCClients()

	// getting router with gin engine
	router := app.routes(grpcClients)

	// Using gin to start api gateway server, exit status 1 if fail to start server
	if err := router.Run(fmt.Sprintf(":%s", gatewayPort)); err != nil {
		log.Fatal(err)
	}
}

func (app *application) initiateGRPCClients() *grpcClientConn {
	orderConn, err := grpc.Dial(ORDER_SERVICE_URL, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		log.Println("error dialing order microservice grpc", err)
	}

	inventoryConn, err := grpc.Dial(INVENTORY_SERVICE_URL, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		log.Println("error dialing inventory microservice grpc", err)
	}

	authConn, err := grpc.Dial(AUTHENTICATION_SERVICE_URL, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		log.Println("error dialing authentication microservice grpc", err)
	}

	return &grpcClientConn{
		orderConn:     orderConn,
		inventoryConn: inventoryConn,
		authConn:      authConn,
	}
}
