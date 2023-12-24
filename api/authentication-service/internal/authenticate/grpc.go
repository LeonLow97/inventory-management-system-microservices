package authenticate

import (
	"context"
	"errors"
	"log"

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

	loginRequestDTO := &LoginRequestDTO{
		Username: req.Username,
		Password: req.Password,
	}

	// sanitize data
	loginSanitize(loginRequestDTO)

	// call Login service (business logic)
	user, token, err := s.service.Login(*loginRequestDTO)

	// Translate specific errors to gRPC status codes
	switch {
	case errors.Is(err, ErrInvalidCredentials), errors.Is(err, ErrInactiveUser), errors.Is(err, ErrNotFound):
		return nil, status.Error(codes.Unauthenticated, "Invalid credentials.")
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

func (s *authenticationGRPCServer) SignUp(ctx context.Context, req *pb.SignUpRequest) (*pb.SignUpResponse, error) {
	// Validate the fields manually for gRPC requests, unable to validator golang package
	if req.Username == "" || len(req.Username) < 5 || len(req.Username) > 50 {
		return nil, status.Error(codes.InvalidArgument, "Username length did not meet requirements.")
	}
	if req.Password == "" || len(req.Password) < 8 || len(req.Password) > 20 {
		return nil, status.Error(codes.InvalidArgument, "Password length did not meet requirements.")
	}
	if req.Email == "" || len(req.Email) < 10 || len(req.Email) > 100 {
		return nil, status.Error(codes.InvalidArgument, "Email length did not meet requirements.")
	}

	signUpRequestDTO := &SignUpRequestDTO{
		FirstName: req.FirstName,
		LastName:  req.LastName,
		Username:  req.Username,
		Password:  req.Password,
		Email:     req.Email,
	}

	// sanitize data
	signUpSanitize(signUpRequestDTO)

	// call sign up service (business logic)
	err := s.service.SignUp(*signUpRequestDTO)
	switch {
	case errors.Is(err, ErrInvalidEmailFormat), errors.Is(err, ErrInvalidPasswordFormat):
		log.Println(err.Error())
		return nil, status.Error(codes.InvalidArgument, err.Error())
	case errors.Is(err, ErrExistingUserFound):
		log.Println(err.Error())
		return nil, status.Error(codes.AlreadyExists, err.Error())
	case err != nil:
		log.Println(err)
		return nil, status.Error(codes.Internal, "Internal Server Error")
	default:
		response := &pb.SignUpResponse{
			Username: req.Username,
		}
		return response, nil
	}
}
