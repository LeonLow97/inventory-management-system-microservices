package grpcclient

import (
	"context"

	"github.com/LeonLow97/internal/core/domain"
	"github.com/LeonLow97/internal/ports"
	pb "github.com/LeonLow97/proto"
	"google.golang.org/grpc"
)

type UserRepo struct {
	conn pb.UserServiceClient
}

func NewUserRepo(conn *grpc.ClientConn) ports.UserRepo {
	return &UserRepo{
		conn: pb.NewUserServiceClient(conn),
	}
}

func (r *UserRepo) GetUsers(ctx context.Context, limit int64, cursor string) ([]domain.User, string, error) {
	grpcReq := &pb.GetUsersRequest{
		Limit:  limit,
		Cursor: cursor,
	}

	grpcResp, err := r.conn.GetUsers(ctx, grpcReq)
	if err != nil {
		return nil, "", err
	}

	users := make([]domain.User, len(grpcResp.Users))
	for i, grpcUser := range grpcResp.Users {
		users[i] = domain.User{
			ID:        grpcUser.ID,
			FirstName: grpcUser.FirstName,
			LastName:  grpcUser.LastName,
			Email:     grpcUser.Email,
			Active:    grpcUser.Active,
			Admin:     grpcUser.Admin,
		}
	}

	return users, grpcResp.NextCursor, nil
}

func (r *UserRepo) UpdateUser(ctx context.Context, req domain.User) error {
	grpcReq := &pb.UpdateUserRequest{
		FirstName: req.FirstName,
		LastName:  req.LastName,
		Password:  req.Password,
	}

	_, err := r.conn.UpdateUser(ctx, grpcReq)
	return err
}
