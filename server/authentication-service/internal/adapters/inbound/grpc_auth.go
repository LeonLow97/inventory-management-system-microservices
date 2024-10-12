package inbound

import (
	"context"
	"errors"

	"github.com/LeonLow97/internal/core/domain"
	"github.com/LeonLow97/internal/core/services"
	pb "github.com/LeonLow97/proto"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type authGRPCServer struct {
	service services.Service
	pb.AuthenticationServiceServer
}

func NewAuthGRPCServer(s services.Service) *authGRPCServer {
	return &authGRPCServer{
		service: s,
	}
}

func (s *authGRPCServer) Login(ctx context.Context, req *pb.AuthRequest) (*pb.AuthResponse, error) {
	user, token, err := s.service.Login(ctx, req.Username, req.Password)
	if err != nil {
		switch {
		case errors.Is(err, services.ErrInvalidCredentials) ||
			errors.Is(err, services.ErrInactiveUser):
			return nil, status.Error(codes.Unauthenticated, "Invalid credentials.")
		default:
			return nil, status.Error(codes.Internal, "Internal Server Error")
		}
	}

	return &pb.AuthResponse{
		FirstName: user.FirstName,
		LastName:  user.LastName,
		Username:  user.Username,
		Email:     user.Email,
		Active:    int32(user.Active),
		Admin:     int32(user.Admin),
		Token:     token,
	}, nil
}

func (s *authGRPCServer) SignUp(ctx context.Context, req *pb.SignUpRequest) (*pb.SignUpResponse, error) {
	user := &domain.User{
		FirstName: req.FirstName,
		LastName:  req.LastName,
		Username:  req.Username,
		Password:  req.Password,
		Email:     req.Email,
	}

	user.Sanitize()

	if err := s.service.SignUp(ctx, user); err != nil {
		switch {
		case errors.Is(err, services.ErrInvalidEmailFormat),
			errors.Is(err, services.ErrInvalidPasswordFormat):
			return nil, status.Error(codes.InvalidArgument, err.Error())
		case errors.Is(err, services.ErrUsernameTaken):
			return nil, status.Error(codes.AlreadyExists, err.Error())
		default:
			return nil, status.Error(codes.Internal, "Internal Server Error")
		}
	}
	return &pb.SignUpResponse{
		Username: req.Username,
	}, nil
}
