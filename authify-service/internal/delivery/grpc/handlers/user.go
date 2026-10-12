package handlers

import (
	"context"
	"errors"

	pb "github.com/LeonLow97/authify/grpc-proto/pb"
	coreerrors "github.com/LeonLow97/internal/core/errors"
	"github.com/LeonLow97/internal/core/usecases"
	"github.com/LeonLow97/internal/delivery/grpc/contextstore"
	grpcerror "github.com/LeonLow97/internal/delivery/grpc/errors"
	empty "google.golang.org/protobuf/types/known/emptypb"
)

type userHandler struct {
	service usecases.Service
	pb.AuthifyServiceServer
}

func NewUserHandler(s usecases.Service) *userHandler {
	return &userHandler{
		service: s,
	}
}

func (h *userHandler) Login(ctx context.Context, req *pb.LoginRequest) (*pb.LoginResponse, error) {
	user, token, err := h.service.Login(ctx, toLoginInput(req))
	if err != nil {
		switch {
		case errors.Is(err, coreerrors.ErrInvalidCredentials) ||
			errors.Is(err, coreerrors.ErrInactiveUser) ||
			errors.Is(err, coreerrors.ErrUserNotFound):
			return nil, grpcerror.ECUnauthorized.GRPCError(err)
		default:
			return nil, grpcerror.ECInternalServerError.GRPCError(err)
		}
	}
	return toLoginResponse(user, token), nil
}

func (h *userHandler) SignUp(ctx context.Context, req *pb.SignUpRequest) (*empty.Empty, error) {
	if err := h.service.SignUp(ctx, toSignUpInput(req)); err != nil {
		switch {
		case errors.Is(err, coreerrors.ErrEmailAlreadyExists):
			return nil, grpcerror.ECEmailAlreadyExists.GRPCError(err)
		default:
			return nil, grpcerror.ECInternalServerError.GRPCError(err)
		}
	}

	return &empty.Empty{}, nil
}

func (h *userHandler) GetUsers(ctx context.Context, req *pb.GetUsersRequest) (*pb.GetUsersResponse, error) {
	// Retrieve admin user ID from grpc request metadata
	adminUserID, err := contextstore.UserIDFromContext(ctx)
	if err != nil {
		return nil, grpcerror.ECUnauthorized.GRPCError(err)
	}

	req.Limit = SanitizePaginationCursor(req.Limit)
	users, nextCursor, err := h.service.GetUsers(ctx, adminUserID, req.Limit, req.Cursor)
	if err != nil {
		switch {
		case errors.Is(err, coreerrors.ErrUnauthorizedAdminAccess):
			return nil, grpcerror.ECUnauthorized.GRPCError(err)
		default:
			return nil, grpcerror.ECInternalServerError.GRPCError(err)
		}
	}

	return toGetUsersResponse(users, nextCursor), nil
}

func (h *userHandler) UpdateUser(ctx context.Context, req *pb.UpdateUserRequest) (*empty.Empty, error) {
	// Retrieve user ID from grpc request metadata
	userID, err := contextstore.UserIDFromContext(ctx)
	if err != nil {
		return nil, grpcerror.ECUnauthorized.GRPCError(err)
	}

	updateUserInput := ToUpdateUserInput(req)
	if err := h.service.UpdateUserByEmail(ctx, userID, updateUserInput); err != nil {
		switch {
		case errors.Is(err, coreerrors.ErrUserNotFound):
			return nil, grpcerror.ECUnauthorized.GRPCError(err)
		default:
			return nil, grpcerror.ECInternalServerError.GRPCError(err)
		}
	}

	return &empty.Empty{}, nil
}
