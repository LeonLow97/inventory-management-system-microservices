package order

import (
	"context"
	"errors"

	pb "github.com/LeonLow97/proto"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
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
	dto := GetOrdersDTO{
		UserID: int(req.UserID),
	}

	orders, err := s.service.GetOrders(dto)
	switch {
	case errors.Is(err, ErrNoOrdersFound):
		return nil, status.Error(codes.NotFound, "No Orders Found")
	case err != nil:
		return nil, status.Error(codes.Internal, "Internal Server Error")
	default:
		var pbOrders []*pb.Order
		for _, order := range *orders {
			pbOrders = append(pbOrders, &pb.Order{
				OrderId:      int64(order.OrderID),
				ProductId:    int64(order.ProductID),
				CustomerName: order.CustomerName,
				BrandName:    order.BrandName,
				CategoryName: order.CategoryName,
				Color:        order.Color,
				Size:         order.Size,
				Quantity:     int64(order.Quantity),
				Description:  order.Description,
				Revenue:      float32(order.Revenue),
				Cost:         float32(order.Cost),
				Profit:       float32(order.Profit),
				HasReviewed:  order.HasReviewed,
				CreatedAt:    order.CreatedAt,
			})
		}
		return &pb.GetOrdersResponse{
			Orders: pbOrders,
		}, nil
	}
}
