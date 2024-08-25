package grpcclient

import (
	"context"

	"github.com/LeonLow97/internal/core/domain"
	"github.com/LeonLow97/internal/ports"
	pb "github.com/LeonLow97/proto"
	"google.golang.org/grpc"
)

type InventoryRepo struct {
	conn pb.InventoryServiceClient
}

func NewInventoryRepo(conn *grpc.ClientConn) ports.InventoryRepo {
	return &InventoryRepo{
		conn: pb.NewInventoryServiceClient(conn),
	}
}

func (r *InventoryRepo) GetProducts(ctx context.Context, userID int) (*[]domain.Product, error) {
	grpcReq := &pb.GetProductsRequest{
		UserID: int32(userID),
	}

	grpcResp, err := r.conn.GetProducts(ctx, grpcReq)
	if err != nil {
		return nil, err
	}

	products := make([]domain.Product, len(grpcResp.Products))
	for i, grpcProduct := range grpcResp.Products {
		products[i] = domain.Product{
			BrandName:    grpcProduct.BrandName,
			CategoryName: grpcProduct.CategoryName,
			ProductName:  grpcProduct.ProductName,
			Description:  grpcProduct.Description,
			Size:         grpcProduct.Size,
			Color:        grpcProduct.Color,
			Quantity:     grpcProduct.Quantity,
			CreatedAt:    grpcProduct.CreatedAt,
			UpdatedAt:    grpcProduct.UpdatedAt,
		}
	}

	return &products, nil
}

func (r *InventoryRepo) GetProductByID(ctx context.Context, userID, productID int) (*domain.Product, error) {
	grpcReq := &pb.GetProductByIDRequest{
		UserID:    int32(userID),
		ProductID: int32(productID),
	}

	grpcResp, err := r.conn.GetProductByID(ctx, grpcReq)
	if err != nil {
		return nil, err
	}

	product := domain.Product{
		BrandName:    grpcResp.BrandName,
		CategoryName: grpcResp.CategoryName,
		ProductName:  grpcResp.ProductName,
		Description:  grpcResp.Description,
		Size:         grpcResp.Size,
		Color:        grpcResp.Color,
		Quantity:     grpcResp.Quantity,
		CreatedAt:    grpcResp.CreatedAt,
		UpdatedAt:    grpcResp.UpdatedAt,
	}

	return &product, nil
}

func (r *InventoryRepo) CreateProduct(ctx context.Context, req domain.Product, userID int) error {
	grpcReq := &pb.CreateProductRequest{
		UserID:       int32(userID),
		BrandName:    req.BrandName,
		CategoryName: req.CategoryName,
		ProductName:  req.ProductName,
		Description:  req.Description,
		Size:         req.Size,
		Color:        req.Color,
		Quantity:     req.Quantity,
	}

	_, err := r.conn.CreateProduct(ctx, grpcReq)
	return err
}

func (r *InventoryRepo) UpdateProduct(ctx context.Context, req domain.Product, userID, productID int) error {
	grpcReq := &pb.UpdateProductRequest{
		UserID:       int32(userID),
		ProductID:    int32(productID),
		BrandName:    req.BrandName,
		CategoryName: req.CategoryName,
		ProductName:  req.ProductName,
		Description:  req.Description,
		Size:         req.Size,
		Color:        req.Color,
		Quantity:     req.Quantity,
	}

	_, err := r.conn.UpdateProduct(ctx, grpcReq)
	return err
}

func (r *InventoryRepo) DeleteProduct(ctx context.Context, userID, productID int) error {
	_, err := r.conn.DeleteProduct(ctx, nil)
	return err
}
