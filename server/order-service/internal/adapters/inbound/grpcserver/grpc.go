package inbound

import (
	"context"
	"errors"

	"github.com/LeonLow97/internal/adapters/outbound"
	"github.com/LeonLow97/internal/core/domain"
	"github.com/LeonLow97/internal/core/services"
	pb "github.com/LeonLow97/proto"
	"github.com/golang/protobuf/ptypes/empty"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type orderGRPCServer struct {
	service services.Service
	pb.OrderServiceServer
}

func NewOrderGRPCServer(s services.Service) *orderGRPCServer {
	return &orderGRPCServer{
		service: s,
	}
}

func (s *orderGRPCServer) GetOrders(ctx context.Context, req *pb.GetOrdersRequest) (*pb.GetOrdersResponse, error) {
	orders, err := s.service.GetOrders(int(req.UserID))
	if err != nil {
		switch {
		case errors.Is(err, outbound.ErrNoOrdersFound):
			return nil, status.Error(codes.NotFound, "No Orders Found")
		default:
			return nil, status.Error(codes.Internal, "Internal Server Error")
		}
	}

	pbOrders := make([]*pb.Order, len(*orders))
	for i, order := range *orders {
		pbOrders[i] = &pb.Order{
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
			Status:       order.Status,
			StatusReason: order.StatusReason,
			CreatedAt:    order.CreatedAt,
		}
	}

	return &pb.GetOrdersResponse{Orders: pbOrders}, nil
}

func (s *orderGRPCServer) GetOrder(ctx context.Context, req *pb.GetOrderRequest) (*pb.Order, error) {
	order, err := s.service.GetOrderByID(int(req.UserID), int(req.OrderID))
	if err != nil {
		switch {
		case errors.Is(err, outbound.ErrNoOrderFound):
			return nil, status.Error(codes.NotFound, "No Order Found")
		default:
			return nil, status.Error(codes.Internal, "Internal Server Error")
		}
	}

	return &pb.Order{
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
		Status:       order.Status,
		StatusReason: order.StatusReason,
		CreatedAt:    order.CreatedAt,
	}, nil
}

func (s *orderGRPCServer) CreateOrder(ctx context.Context, req *pb.CreateOrderRequest) (*empty.Empty, error) {
	order := domain.Order{
		CustomerName: req.CustomerName,
		BrandName:    req.BrandName,
		CategoryName: req.CategoryName,
		Color:        req.Color,
		Size:         req.Size,
		Quantity:     int(req.Quantity),
		Description:  req.Description,
		Revenue:      req.Revenue,
		Cost:         req.Cost,
		Profit:       req.Profit,
		HasReviewed:  req.HasReviewed,
	}

	if err := s.service.CreateOrder(ctx, order, int(req.UserID), req.ProductName); err != nil {
		switch {
		case errors.Is(err, outbound.ErrProductNotFound):
			return &empty.Empty{}, status.Error(codes.NotFound, "Product not found for the given user.")
		default:
			return &empty.Empty{}, status.Error(codes.Internal, "Internal Server Error")
		}
	}
	return &empty.Empty{}, nil
}
