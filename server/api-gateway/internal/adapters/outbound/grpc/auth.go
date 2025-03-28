package grpcclient

import (
	"context"

	"github.com/LeonLow97/internal/core/domain"
	"github.com/LeonLow97/internal/ports"
	pb "github.com/LeonLow97/proto"
	"google.golang.org/grpc"
)

type AuthRepo struct {
	conn pb.AuthenticationServiceClient
}

func NewAuthRepo(conn *grpc.ClientConn) ports.AuthRepo {
	return &AuthRepo{
		conn: pb.NewAuthenticationServiceClient(conn),
	}
}

func (r *AuthRepo) Login(ctx context.Context, req domain.User) (*domain.User, error) {
	grpcReq := &pb.LoginRequest{
		Email:    req.Email,
		Password: req.Password,
	}

	grpcResp, err := r.conn.Login(ctx, grpcReq)
	if err != nil {
		return nil, err
	}

	user := &domain.User{
		FirstName: grpcResp.FirstName,
		LastName:  grpcResp.LastName,
		Email:     grpcResp.Email,
		Active:    grpcResp.Active,
		Admin:     grpcResp.Admin,
		Token:     grpcResp.Token,
	}

	return user, nil
}

func (r *AuthRepo) SignUp(ctx context.Context, req domain.User) error {
	grpcReq := &pb.SignUpRequest{
		Email:     req.Email,
		Password:  req.Password,
		FirstName: req.FirstName,
		LastName:  req.LastName,
	}

	_, err := r.conn.SignUp(ctx, grpcReq)
	return err
}
