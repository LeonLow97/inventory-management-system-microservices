package authenticate

import (
	"context"
	"errors"

	pb "github.com/LeonLow97/proto"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type authenticationGRPCServer struct {
	service Service
	pb.AuthenticationServiceServer
}

func NewAuthenticateGRPCHandler(s Service) *authenticationGRPCServer {
	return &authenticationGRPCServer{
		service: s,
	}
}

func (s *authenticationGRPCServer) Authenticate(ctx context.Context, req *pb.AuthRequest) (*pb.AuthResponse, error) {
	// Validate the fields manually for gRPC requests, unable to validator golang package
	if req.Username == "" || len(req.Username) < 5 || len(req.Username) > 50 {
		return nil, status.Error(codes.InvalidArgument, "Invalid username format.")
	}
	if req.Password == "" || len(req.Password) < 8 || len(req.Password) > 20 {
		return nil, status.Error(codes.InvalidArgument, "Invalid password format.")
	}

	loginRequestDTO := convertGRPCRequestToDTO(req)

	// sanitize data
	loginSanitize(loginRequestDTO)

	// call Login service (business logic)
	user, token, err := s.service.Login(*loginRequestDTO)

	// Translate specific errors to gRPC status codes
	switch {
	case errors.Is(err, ErrInvalidCredentials), errors.Is(err, ErrInactiveUser), errors.Is(err, ErrNotFound):
		return nil, status.Error(codes.Unauthenticated, "Invalid credentials")
	case err != nil:
		return nil, status.Error(codes.Internal, "Internal server error")
	default:
		response := &pb.AuthResponse{
			FirstName: user.FirstName,
			LastName:  user.LastName,
			Username:  user.Username,
			Email:     user.Email,
			Active:    int32(user.Active),
			Admin:     int32(user.Admin),
			Token:     token,
		}
		return response, nil
	}
}

func convertGRPCRequestToDTO(req *pb.AuthRequest) *LoginRequestDTO {
	return &LoginRequestDTO{
		Username: req.Username,
		Password: req.Password,
	}
}
