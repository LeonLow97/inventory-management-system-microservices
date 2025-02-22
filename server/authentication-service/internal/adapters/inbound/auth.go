package inbound

import (
	"context"
	"errors"

	"github.com/LeonLow97/internal/core/services"
	"github.com/LeonLow97/internal/pkg/grpcerror"
	pb "github.com/LeonLow97/proto"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type authInbound struct {
	service services.Service
	pb.AuthenticationServiceServer
}

func NewAuthGRPCServer(s services.Service) *authInbound {
	return &authInbound{
		service: s,
	}
}

func (s *authInbound) Login(ctx context.Context, req *pb.LoginRequest) (*pb.LoginResponse, error) {
	loginInput := ToLoginInput(req)
	user, token, err := s.service.Login(ctx, loginInput)
	if err != nil {
		switch {
		case errors.Is(err, services.ErrInvalidCredentials) ||
			errors.Is(err, services.ErrInactiveUser):
			return nil, grpcerror.ECInactiveUser.GRPCError(err)
		default:
			return nil, grpcerror.ECInternalServerError.GRPCError(err)
		}
	}

	return &pb.LoginResponse{
		FirstName: user.FirstName,
		LastName:  user.LastName,
		Username:  user.Username,
		Email:     user.Email,
		Active:    int32(user.Active),
		Admin:     int32(user.Admin),
		Token:     token,
	}, nil
}

func (s *authInbound) SignUp(ctx context.Context, req *pb.SignUpRequest) (*pb.SignUpResponse, error) {
	signupInput := ToSignUpInput(req)
	if err := s.service.SignUp(ctx, signupInput); err != nil {
		switch {
		case errors.Is(err, services.ErrInvalidEmailFormat),
			errors.Is(err, services.ErrInvalidPasswordFormat):
			return nil, status.Error(codes.InvalidArgument, err.Error())
		case errors.Is(err, services.ErrUsernameTaken):
			return nil, grpcerror.ECUsernameTaken.GRPCError(err)
		default:
			return nil, grpcerror.ECInternalServerError.GRPCError(err)
		}
	}

	return &pb.SignUpResponse{
		Username: req.Username,
	}, nil
}
