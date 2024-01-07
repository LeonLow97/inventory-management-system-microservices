package order

import (
	"context"
	"errors"

	pb "github.com/LeonLow97/proto"
	empty "google.golang.org/protobuf/types/known/emptypb"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type OrderGRPCServer struct {
	service Service
	pb.OrderServiceServer
}

func NewOrderGRPCHandler(service Service) *OrderGRPCServer {
	return &OrderGRPCServer{
		service: service,
	}
}

func (s *OrderGRPCServer) GetOrders(ctx context.Context, req *pb.GetOrdersRequest) (*pb.GetOrdersResponse, error) {
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

func (s *OrderGRPCServer) GetOrder(ctx context.Context, req *pb.GetOrderRequest) (*pb.Order, error) {
	dto := GetOrderDTO{
		UserID:  int(req.UserID),
		OrderID: int(req.OrderID),
	}

	order, err := s.service.GetOrderByID(dto)
	switch {
	case errors.Is(err, ErrNoOrderFound):
		return nil, status.Error(codes.NotFound, "No Order Found")
	case err != nil:
		return nil, status.Error(codes.Internal, "Internal Server Error")
	default:
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
			CreatedAt:    order.CreatedAt,
		}, nil
	}
}

func (s *OrderGRPCServer) CreateOrder(ctx context.Context, req *pb.CreateOrderRequest) (*empty.Empty, error) {
	createOrderDTO := &CreateOrderDTO{
		UserID:       int(req.UserID),
		CustomerName: req.CustomerName,
		ProductName:  req.ProductName,
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

	err := s.service.CreateOrder(*createOrderDTO)
	switch {
	case errors.Is(err, ErrProductNotFound):
		return &empty.Empty{}, status.Error(codes.NotFound, "Product not found for the given user.")
	default:
		return &empty.Empty{}, err
	}
}
