package outbound

import (
	"context"
	"strconv"

	pb "github.com/LeonLow97/authify/grpc-proto/pb"
	"github.com/LeonLow97/internal/contextstore"
	"github.com/LeonLow97/internal/core/domain"
	"github.com/LeonLow97/internal/core/ports"
	"google.golang.org/grpc/metadata"
)

const requestIDMetadataKey = "x-request-id"

type userOutbound struct {
	conn pb.AuthifyServiceClient
}

func NewUserOutbound(conn pb.AuthifyServiceClient) ports.UserOutbound {
	return &userOutbound{
		conn: conn,
	}
}

func (r *userOutbound) Login(ctx context.Context, req domain.User) (*domain.User, string, error) {
	grpcReq := &pb.LoginRequest{
		Email:    req.Email,
		Password: req.Password,
	}

	grpcResp, err := r.conn.Login(withRequestIDMetadata(ctx), grpcReq)
	if err != nil {
		return nil, "", mapGRPCError(err)
	}

	return &domain.User{
		FirstName: grpcResp.FirstName,
		LastName:  grpcResp.LastName,
		Email:     grpcResp.Email,
		Active:    grpcResp.Active,
		Admin:     grpcResp.Admin,
	}, grpcResp.Token, nil
}

func (r *userOutbound) SignUp(ctx context.Context, req domain.User) error {
	grpcReq := &pb.SignUpRequest{
		Email:     req.Email,
		Password:  req.Password,
		FirstName: req.FirstName,
		LastName:  req.LastName,
	}

	_, err := r.conn.SignUp(withRequestIDMetadata(ctx), grpcReq)
	return mapGRPCError(err)
}

func (r *userOutbound) GetUsers(ctx context.Context, limit int64, cursor string) ([]domain.User, string, error) {
	grpcReq := &pb.GetUsersRequest{
		Limit:  limit,
		Cursor: cursor,
	}

	grpcCtx := withUserIDMetadata(withRequestIDMetadata(ctx))
	grpcResp, err := r.conn.GetUsers(grpcCtx, grpcReq)
	if err != nil {
		return nil, "", mapGRPCError(err)
	}

	users := make([]domain.User, len(grpcResp.Users))
	for i, grpcUser := range grpcResp.Users {
		users[i] = domain.User{
			ID:        grpcUser.ID,
			FirstName: grpcUser.FirstName,
			LastName:  grpcUser.LastName,
			Email:     grpcUser.Email,
			Active:    grpcUser.Active,
			Admin:     grpcUser.Admin,
		}
	}

	return users, grpcResp.NextCursor, nil
}

func (r *userOutbound) UpdateUser(ctx context.Context, req domain.User) error {
	grpcReq := &pb.UpdateUserRequest{
		FirstName: req.FirstName,
		LastName:  req.LastName,
		Password:  req.Password,
	}

	grpcCtx := withUserIDMetadata(withRequestIDMetadata(ctx))
	_, err := r.conn.UpdateUser(grpcCtx, grpcReq)
	return mapGRPCError(err)
}

func withUserIDMetadata(ctx context.Context) context.Context {
	userID, ok := contextstore.UserIDFromContext(ctx)
	if !ok {
		return ctx
	}
	return metadata.AppendToOutgoingContext(ctx, "user_id", strconv.FormatInt(userID, 10))
}

func withRequestIDMetadata(ctx context.Context) context.Context {
	requestID, ok := contextstore.RequestIDFromContext(ctx)
	if !ok || requestID == "" {
		return ctx
	}
	return metadata.AppendToOutgoingContext(ctx, requestIDMetadataKey, requestID)
}
