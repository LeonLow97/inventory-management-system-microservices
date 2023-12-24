package users

import (
	"context"
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
	log.Println(req)

	return &empty.Empty{}, nil
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
