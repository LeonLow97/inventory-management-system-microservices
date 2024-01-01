package inventory

import (
	"context"
	"errors"
	"log"

	pb "github.com/LeonLow97/proto"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type inventoryGRPCServer struct {
	service Service
	pb.InventoryServiceServer
}

func NewInventoryGRPCHandler(s Service) *inventoryGRPCServer {
	return &inventoryGRPCServer{
		service: s,
	}
}

func (s *inventoryGRPCServer) GetProducts(ctx context.Context, req *pb.GetProductsRequest) (*pb.GetProductsResponse, error) {
	products, err := s.service.GetProducts(int(req.UserID))
	if err != nil {
		log.Println(err)
		return nil, status.Error(codes.Internal, "Internal Server Error")
	}

	var pbProducts []*pb.Product
	for _, product := range *products {
		pbProducts = append(pbProducts, &pb.Product{
			BrandName:    product.BrandName,
			CategoryName: product.CategoryName,
			ProductName:  product.ProductName,
			Description:  product.Description,
			Size:         product.Size,
			Color:        product.Color,
			Quantity:     int32(product.Quantity),
			CreatedAt:    product.CreatedAt,
			UpdatedAt:    product.UpdatedAt,
		})
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

	getProductByIdDTO := &GetProductByIdDTO{
		UserID:    int(req.UserID),
		ProductID: int(req.ProductID),
	}

	product, err := s.service.GetProductByID(*getProductByIdDTO)
	switch {
	case errors.Is(err, ErrProductNotFound):
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
