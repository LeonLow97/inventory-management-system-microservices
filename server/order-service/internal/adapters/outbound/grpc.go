package outbound

import (
	"context"

	pb "github.com/LeonLow97/proto"
)

func (r *Repository) GetProductID(ctx context.Context, userID int, brandName, categoryName, productName string) (int, error) {
	req := &pb.GetProductDetailsRequest{
		UserID:       int64(userID),
		BrandName:    brandName,
		CategoryName: categoryName,
		ProductName:  productName,
	}

	resp, err := r.grpcConn.GetProductDetails(ctx, req)
	if err != nil {
		return 0, err
	}
	if resp.ProductID == 0 {
		return 0, ErrProductNotFound
	}

	return int(resp.ProductID), nil
}
