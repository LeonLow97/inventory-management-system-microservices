package inventory

import (
	"context"
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
