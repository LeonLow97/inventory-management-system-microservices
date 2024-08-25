package grpcclient

import (
	"context"

	"github.com/LeonLow97/internal/core/domain"
	"github.com/LeonLow97/internal/ports"
	pb "github.com/LeonLow97/proto"
	"google.golang.org/grpc"
)

type OrderRepo struct {
	conn pb.OrderServiceClient
}

func NewOrderRepo(conn *grpc.ClientConn) ports.OrderRepo {
	return &OrderRepo{
		conn: pb.NewOrderServiceClient(conn),
	}
}

func (r *OrderRepo) GetOrders(ctx context.Context, userID int) (*[]domain.Order, error) {
	grpcReq := &pb.GetOrdersRequest{
		UserID: int64(userID),
	}

	grpcResp, err := r.conn.GetOrders(ctx, grpcReq)
	if err != nil {
		return nil, err
	}

	orders := make([]domain.Order, len(grpcResp.Orders))

	for i, grpcOrder := range grpcResp.Orders {
		orders[i] = domain.Order{
			OrderID:      grpcOrder.OrderId,
			ProductID:    grpcOrder.ProductId,
			CustomerName: grpcOrder.CustomerName,
			BrandName:    grpcOrder.BrandName,
			CategoryName: grpcOrder.CategoryName,
			Color:        grpcOrder.Color,
			Size:         grpcOrder.Size,
			Quantity:     grpcOrder.Quantity,
			Description:  grpcOrder.Description,
			Revenue:      grpcOrder.Revenue,
			Cost:         grpcOrder.Cost,
			Profit:       grpcOrder.Profit,
			HasReviewed:  grpcOrder.HasReviewed,
			Status:       grpcOrder.Status,
			StatusReason: grpcOrder.StatusReason,
			CreatedAt:    grpcOrder.CreatedAt,
		}
	}

	return &orders, nil
}

func (r *OrderRepo) GetOrder(ctx context.Context, userID, orderID int) (*domain.Order, error) {
	grpcReq := &pb.GetOrderRequest{
		UserID:  int64(userID),
		OrderID: int64(orderID),
	}

	grpcResp, err := r.conn.GetOrder(ctx, grpcReq)
	if err != nil {
		return nil, err
	}

	order := domain.Order{
		OrderID:      grpcResp.OrderId,
		ProductID:    grpcResp.ProductId,
		CustomerName: grpcResp.CustomerName,
		BrandName:    grpcResp.BrandName,
		CategoryName: grpcResp.CategoryName,
		Color:        grpcResp.Color,
		Size:         grpcResp.Size,
		Quantity:     grpcResp.Quantity,
		Description:  grpcResp.Description,
		Revenue:      grpcResp.Revenue,
		Cost:         grpcResp.Cost,
		Profit:       grpcResp.Profit,
		HasReviewed:  grpcResp.HasReviewed,
		Status:       grpcResp.Status,
		StatusReason: grpcResp.StatusReason,
		CreatedAt:    grpcResp.CreatedAt,
	}

	return &order, nil
}

func (r *OrderRepo) CreateOrder(ctx context.Context, req domain.Order, userID int) error {
	grpcReq := &pb.CreateOrderRequest{
		UserID:       int64(userID),
		CustomerName: req.CustomerName,
		ProductName:  req.ProductName,
		BrandName:    req.BrandName,
		CategoryName: req.CategoryName,
		Color:        req.Color,
		Size:         req.Size,
		Quantity:     req.Quantity,
		Description:  req.Description,
		Revenue:      int64(req.Revenue),
		Cost:         int64(req.Cost),
		Profit:       int64(req.Profit),
		HasReviewed:  req.HasReviewed,
	}

	_, err := r.conn.CreateOrder(ctx, grpcReq)
	return err
}
