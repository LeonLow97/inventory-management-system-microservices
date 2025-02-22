package inbound

import (
	"context"
	"errors"

	"github.com/LeonLow97/internal/core/domain"
	"github.com/LeonLow97/internal/core/services"
	pb "github.com/LeonLow97/proto"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	empty "google.golang.org/protobuf/types/known/emptypb"
)

type userGRPCServer struct {
	service services.Service
	pb.UserServiceServer
}

func NewUserGRPCServer(s services.Service) *userGRPCServer {
	return &userGRPCServer{
		service: s,
	}
}

func (s *userGRPCServer) UpdateUser(ctx context.Context, req *pb.UpdateUserRequest) (*empty.Empty, error) {
	user := &domain.User{
		ID:        int(req.UserID),
		FirstName: req.FirstName,
		LastName:  req.LastName,
		Password:  req.Password,
		Email:     req.Email,
	}

	user.Sanitize()

	if err := s.service.UpdateUser(ctx, user); err != nil {
		switch {
		case errors.Is(err, services.ErrUserNotFound):
			return nil, status.Error(codes.NotFound, "User does not exist.")
		case errors.Is(err, services.ErrInvalidPasswordFormat):
			return nil, status.Error(codes.InvalidArgument, "Invalid password format. Please try again.")
		default:
			return nil, status.Error(codes.Internal, "Internal Server Error")
		}
	}
	return &empty.Empty{}, nil
}
