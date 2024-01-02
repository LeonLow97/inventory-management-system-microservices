package inventory

import (
	"context"
	"errors"
	"log"

	pb "github.com/LeonLow97/proto"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	empty "google.golang.org/protobuf/types/known/emptypb"
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

func (s *inventoryGRPCServer) CreateProduct(ctx context.Context, req *pb.CreateProductRequest) (*empty.Empty, error) {
	createProductDTO := &CreateProductDTO{
		UserID:       int(req.GetUserID()),
		BrandName:    req.GetBrandName(),
		CategoryName: req.GetCategoryName(),
		ProductName:  req.GetProductName(),
		Description:  req.GetDescription(),
		Size:         req.GetSize(),
		Color:        req.GetColor(),
		Quantity:     int(req.GetQuantity()),
	}

	// sanitize data
	createProductSanitize(createProductDTO)

	err := s.service.CreateProduct(*createProductDTO)
	switch {
	case errors.Is(err, ErrBrandNotFound):
		return &empty.Empty{}, status.Error(codes.NotFound, "Brand not found.")
	case errors.Is(err, ErrCategoryNotFound):
		return &empty.Empty{}, status.Error(codes.NotFound, "Category not found.")
	case err != nil:
		return &empty.Empty{}, status.Error(codes.Internal, "Internal Server Error")
	default:
		return &empty.Empty{}, nil
	}
}

func (s *inventoryGRPCServer) UpdateProduct(ctx context.Context, req *pb.UpdateProductRequest) (*empty.Empty, error) {
	updateProductDTO := &UpdateProductDTO{
		UserID:       int(req.UserID),
		ProductID:    int(req.ProductID),
		BrandName:    req.GetBrandName(),
		CategoryName: req.GetCategoryName(),
		ProductName:  req.GetProductName(),
		Description:  req.GetDescription(),
		Size:         req.GetSize(),
		Color:        req.GetColor(),
		Quantity:     int(req.GetQuantity()),
	}

	// sanitize data
	updateProductSanitize(updateProductDTO)

	err := s.service.UpdateProductByID(*updateProductDTO)
	switch {
	case errors.Is(err, ErrBrandNotFound):
		return &empty.Empty{}, status.Error(codes.NotFound, "Brand not found.")
	case errors.Is(err, ErrCategoryNotFound):
		return &empty.Empty{}, status.Error(codes.NotFound, "Category not found.")
	case errors.Is(err, ErrProductNotFound):
		return &empty.Empty{}, status.Error(codes.NotFound, "Product does not exist for the user.")
	case err != nil:
		return &empty.Empty{}, status.Error(codes.Internal, "Internal Server Error")
	default:
		return &empty.Empty{}, nil
	}
}

func (s *inventoryGRPCServer) DeleteProduct(ctx context.Context, req *pb.DeleteProductRequest) (*empty.Empty, error) {
	deleteProductDTO := &DeleteProductDTO{
		UserID:    int(req.UserID),
		ProductID: int(req.ProductID),
	}

	err := s.service.DeleteProductByID(*deleteProductDTO)
	switch {
	case errors.Is(err, ErrProductNotFound):
		return &empty.Empty{}, status.Error(codes.NotFound, "Product does not exist for the user.")
	case err != nil:
		return &empty.Empty{}, status.Error(codes.Internal, "Internal Server Error")
	default:
		return &empty.Empty{}, nil
	}
}
