package grpcserver

import (
	"context"
	"errors"
	"log"

	"github.com/LeonLow97/internal/adapters/outbound"
	"github.com/LeonLow97/internal/core/domain"
	"github.com/LeonLow97/internal/core/services"
	pb "github.com/LeonLow97/proto"
	"github.com/golang/protobuf/ptypes/empty"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type inventoryGRPCServer struct {
	service services.Service
	pb.InventoryServiceServer
}

func NewInventoryGRPCServer(s services.Service) *inventoryGRPCServer {
	return &inventoryGRPCServer{
		service: s,
	}
}

func (s *inventoryGRPCServer) GetProducts(ctx context.Context, req *pb.GetProductsRequest) (*pb.GetProductsResponse, error) {
	products, err := s.service.GetProducts(int(req.UserID))
	switch {
	case errors.Is(err, outbound.ErrProductsNotFound):
		return nil, status.Error(codes.NotFound, "No Products found for the given User ID")
	case err != nil:
		return nil, status.Error(codes.Internal, "Internal Server Error")
	}

	pbProducts := make([]*pb.Product, len(*products))
	for i, product := range *products {
		pbProducts[i] = &pb.Product{
			BrandName:    product.BrandName,
			CategoryName: product.CategoryName,
			ProductName:  product.ProductName,
			Description:  product.Description,
			Size:         product.Size,
			Color:        product.Color,
			Quantity:     int32(product.Quantity),
			CreatedAt:    product.CreatedAt,
			UpdatedAt:    product.UpdatedAt,
		}
	}

	return &pb.GetProductsResponse{
		Products: pbProducts,
	}, nil
}

func (s *inventoryGRPCServer) GetProductByID(ctx context.Context, req *pb.GetProductByIDRequest) (*pb.Product, error) {
	// validate the fields
	if req.UserID <= 0 {
		return nil, status.Error(codes.InvalidArgument, "UserID must be greater than 0")
	}
	if req.ProductID <= 0 {
		return nil, status.Error(codes.InvalidArgument, "ProductID must be greater than 0")
	}

	product, err := s.service.GetProductByID(int(req.UserID), int(req.ProductID))
	switch {
	case errors.Is(err, outbound.ErrProductNotFound):
		log.Println(err)
		return nil, status.Error(codes.NotFound, "Product does not exist.")
	case err != nil:
		log.Println(err)
		return nil, status.Error(codes.Internal, "Internal Server Error")
	default:
		return &pb.Product{
			BrandName:    product.BrandName,
			CategoryName: product.CategoryName,
			ProductName:  product.ProductName,
			Description:  product.Description,
			Size:         product.Size,
			Color:        product.Color,
			Quantity:     int32(product.Quantity),
			CreatedAt:    product.CreatedAt,
			UpdatedAt:    product.UpdatedAt,
		}, nil
	}
}

func (s *inventoryGRPCServer) CreateProduct(ctx context.Context, req *pb.CreateProductRequest) (*empty.Empty, error) {
	product := &domain.Product{
		BrandName:    req.GetBrandName(),
		CategoryName: req.GetCategoryName(),
		ProductName:  req.GetProductName(),
		Description:  req.GetDescription(),
		Size:         req.GetSize(),
		Color:        req.GetColor(),
		Quantity:     int(req.GetQuantity()),
	}

	// sanitize data
	product.Sanitize()

	err := s.service.CreateProduct(*product, int(req.GetUserID()))
	switch {
	case errors.Is(err, outbound.ErrBrandNotFound):
		return &empty.Empty{}, status.Error(codes.NotFound, "Brand not found.")
	case errors.Is(err, outbound.ErrCategoryNotFound):
		return &empty.Empty{}, status.Error(codes.NotFound, "Category not found.")
	case err != nil:
		return &empty.Empty{}, status.Error(codes.Internal, "Internal Server Error")
	default:
		return &empty.Empty{}, nil
	}
}

func (s *inventoryGRPCServer) UpdateProduct(ctx context.Context, req *pb.UpdateProductRequest) (*empty.Empty, error) {
	product := &domain.Product{
		BrandName:    req.GetBrandName(),
		CategoryName: req.GetCategoryName(),
		ProductName:  req.GetProductName(),
		Description:  req.GetDescription(),
		Size:         req.GetSize(),
		Color:        req.GetColor(),
		Quantity:     int(req.GetQuantity()),
	}

	// sanitize data
	product.Sanitize()

	if err := s.service.UpdateProductByID(*product, int(req.UserID), int(req.ProductID)); err != nil {
		switch {
		case errors.Is(err, outbound.ErrBrandNotFound):
			return &empty.Empty{}, status.Error(codes.NotFound, "Brand not found.")
		case errors.Is(err, outbound.ErrCategoryNotFound):
			return &empty.Empty{}, status.Error(codes.NotFound, "Category not found.")
		case errors.Is(err, outbound.ErrProductNotFound):
			return &empty.Empty{}, status.Error(codes.NotFound, "Product does not exist for the user.")
		default:
			return &empty.Empty{}, status.Error(codes.Internal, "Internal Server Error")
		}
	}
	return &empty.Empty{}, nil
}

func (s *inventoryGRPCServer) DeleteProduct(ctx context.Context, req *pb.DeleteProductRequest) (*empty.Empty, error) {
	err := s.service.DeleteProductByID(int(req.UserID), int(req.ProductID))
	switch {
	case errors.Is(err, outbound.ErrProductNotFound):
		return &empty.Empty{}, status.Error(codes.NotFound, "Product does not exist for the user.")
	case err != nil:
		return &empty.Empty{}, status.Error(codes.Internal, "Internal Server Error")
	default:
		return &empty.Empty{}, nil
	}
}

func (s *inventoryGRPCServer) GetProductDetails(ctx context.Context, req *pb.GetProductDetailsRequest) (*pb.GetProductDetailsResponse, error) {
	product, err := s.service.GetProductByName(int(req.UserID), req.ProductName)
	switch {
	case errors.Is(err, outbound.ErrProductNotFound):
		log.Println(err)
		return nil, status.Error(codes.NotFound, "Product does not exist.")
	case err != nil:
		log.Println(err)
		return nil, status.Error(codes.Internal, "Internal Server Error")
	default:
		return &pb.GetProductDetailsResponse{
			UserID:    req.UserID,
			ProductID: product.ID,
		}, nil
	}
}
