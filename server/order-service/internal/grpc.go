package order

import (
	pb "github.com/LeonLow97/proto"
)

type orderGRPCServer struct {
	service Service
	pb.OrderServiceServer
}

func NewOrderGRPCHandler(service Service) *orderGRPCServer {
	return &orderGRPCServer{
		service: service,
	}
}
