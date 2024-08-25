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

func (r *UserRepo) GetUsers(ctx context.Context) (*[]domain.User, error) {
	grpcResp, err := r.conn.GetUsers(ctx, nil)
	if err != nil {
		return nil, err
	}

	users := make([]domain.User, len(grpcResp.Users))
	for i, grpcUser := range grpcResp.Users {
		users[i] = domain.User{
			FirstName: grpcUser.FirstName,
			LastName:  grpcUser.LastName,
			Username:  grpcUser.Username,
			Email:     grpcUser.Email,
			Active:    int(grpcUser.Active),
			Admin:     int(grpcUser.Admin),
		}
	}

	return &users, nil
}

func (r *UserRepo) UpdateUser(ctx context.Context, req domain.User, userID int) error {
	grpcReq := &pb.UpdateUserRequest{
		UserID:    int64(userID),
		FirstName: req.FirstName,
		LastName:  req.LastName,
		Password:  req.Password,
		Email:     req.Email,
	}

	_, err := r.conn.UpdateUser(ctx, grpcReq)
	return err
}
