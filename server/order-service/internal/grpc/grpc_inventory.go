package grpcclient

import (
	"context"
	"log"
	"time"

	pb "github.com/LeonLow97/proto"
	"google.golang.org/grpc"
)

type InventoryServiceClient interface {
	GRPCGetProductDetailsHandler(userID int, brandName, categoryName, productName string) (*pb.GetProductDetailsResponse, error)
}

type inventoryGRPCClient struct {
	client pb.InventoryServiceClient
}

func NewInventoryGRPCClient(conn *grpc.ClientConn) InventoryServiceClient {
	return &inventoryGRPCClient{
		client: pb.NewInventoryServiceClient(conn),
	}
}

func (i inventoryGRPCClient) GRPCGetProductDetailsHandler(userID int, brandName, categoryName, productName string) (*pb.GetProductDetailsResponse, error) {
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*10)
	defer cancel()

	getProductDetailsRequest := &pb.GetProductDetailsRequest{
		UserID:       int64(userID),
		BrandName:    brandName,
		CategoryName: categoryName,
		ProductName:  productName,
	}

	resp, err := i.client.GetProductDetails(ctx, getProductDetailsRequest)
	if err != nil {
		log.Println(err)
		return nil, err
	}

	return resp, nil
}
