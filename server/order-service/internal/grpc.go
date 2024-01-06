package order

import (
	"context"

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

func (s *orderGRPCServer) GetOrders(ctx context.Context, req *pb.GetOrdersRequest) (*pb.GetOrdersResponse, error) {
	return &pb.GetOrdersResponse{}, nil
}
