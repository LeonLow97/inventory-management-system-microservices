package inbound

import (
	"context"
	"errors"

	"github.com/LeonLow97/internal/core/services"
	"github.com/LeonLow97/internal/pkg/grpcerror"
	pb "github.com/LeonLow97/proto"
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
	updateUserInput := ToUpdateUserInput(req)
	if err := s.service.UpdateUser(ctx, req.UserID, updateUserInput); err != nil {
		switch {
		case errors.Is(err, services.ErrUserNotFound):
			return nil, grpcerror.ECUnauthorized.GRPCError(err)
		default:
			return nil, grpcerror.ECInternalServerError.GRPCError(err)
		}
	}
	return &empty.Empty{}, nil
}
