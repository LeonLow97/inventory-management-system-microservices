package users

import (
	"context"
	"errors"
	"fmt"
	"log"

	pb "github.com/LeonLow97/proto"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	empty "google.golang.org/protobuf/types/known/emptypb"
)

type usersGRPCServer struct {
	service Service
	pb.UserServiceServer
}

func NewUsersGRPCHandler(s Service) *usersGRPCServer {
	return &usersGRPCServer{
		service: s,
	}
}

func (s *usersGRPCServer) UpdateUser(ctx context.Context, req *pb.UpdateUserRequest) (*empty.Empty, error) {
	// Validate the fields manually for gRPC requests, unable to validator golang package
	if req.UserID == 0 || req.UserID < 0 {
		return nil, status.Error(codes.InvalidArgument, fmt.Sprintf("UserID %d is invalid", req.UserID))
	}
	if len(req.FirstName) > 50 {
		return nil, status.Error(codes.InvalidArgument, "LastName length did not meet requirements.")
	}
	if len(req.LastName) > 50 {
		return nil, status.Error(codes.InvalidArgument, "FirstName length did not meet requirements.")
	}
	if len(req.Password) > 0 {
		if len(req.Password) < 8 || len(req.Password) > 20 {
			return nil, status.Error(codes.InvalidArgument, "Password length did not meet requirements.")
		}
	}
	if len(req.Email) > 0 {
		if len(req.Email) < 10 || len(req.Email) > 100 {
			return nil, status.Error(codes.InvalidArgument, "Email length did not meet requirements.")
		}
	}

	updateUserDTO := &UpdateUserRequestDTO{
		UserID:    int(req.UserID),
		FirstName: req.FirstName,
		LastName:  req.LastName,
		Password:  req.Password,
		Email:     req.Email,
	}

	// sanitize data
	updateUserSanitize(updateUserDTO)

	err := s.service.UpdateUser(*updateUserDTO)
	switch {
	case errors.Is(err, ErrSameValue):
		log.Println(err.Error())
		return nil, status.Error(codes.InvalidArgument, "Update of same value is not allowed.")
	case errors.Is(err, ErrInvalidPasswordFormat):
		log.Println(err.Error())
		return nil, status.Error(codes.InvalidArgument, "Invalid password format. Please try again.")
	case errors.Is(err, ErrNotFound):
		log.Println(err.Error())
		return nil, status.Error(codes.NotFound, "User does not exist.")
	case err != nil:
		log.Println(err)
		return nil, status.Error(codes.Internal, "Internal Server Error")
	default:
		return &empty.Empty{}, nil
	}
}

func (s *usersGRPCServer) GetUsers(ctx context.Context, empty *empty.Empty) (*pb.GetUsersResponse, error) {
	users, err := s.service.GetUsers()
	if err != nil {
		log.Println(err)
		return nil, status.Error(codes.Internal, "Internal Server Error")
	}

	// Convert []*User to []*pb.User directly
	var pbUsers []*pb.User
	for _, user := range *users {
		pbUsers = append(pbUsers, &pb.User{
			FirstName: user.FirstName,
			LastName:  user.LastName,
			Username:  user.Username,
			Email:     user.Email,
			Active:    int64(user.Active),
			Admin:     int64(user.Admin),
		})
	}

	return &pb.GetUsersResponse{
		Users: pbUsers,
	}, nil
}
