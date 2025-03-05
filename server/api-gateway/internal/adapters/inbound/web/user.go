package web

import (
	"net/http"
	"strconv"

	"github.com/LeonLow97/internal/adapters/inbound/web/dto"
	"github.com/LeonLow97/internal/core/domain"
	user "github.com/LeonLow97/internal/core/services/user"
	"github.com/LeonLow97/internal/pkg/apierror"
	"github.com/LeonLow97/internal/pkg/contextstore"
	"github.com/LeonLow97/internal/pkg/handler"
	"github.com/gin-gonic/gin"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/metadata"
	"google.golang.org/grpc/status"
)

type UserHandler struct {
	handler.Handler
	UserService user.User
}

func NewUserHandler(userService user.User) *UserHandler {
	return &UserHandler{
		Handler:     handler.NewHandler(),
		UserService: userService,
	}
}

func (h *UserHandler) GetUsers() gin.HandlerFunc {
	return func(c *gin.Context) {
		limitStr := h.GetQueryParam(c, "limit", "10")
		limit, err := strconv.Atoi(limitStr)
		if err != nil {
			apierror.ErrBadRequest.APIError(c, err)
			return
		}
		cursor := h.GetQueryParam(c, "cursor", "")

		md, err := contextstore.GRPCMetadataFromContext(c)
		if err != nil {
			apierror.ErrInternalServerError.APIError(c, err)
			return
		}
		grpcCtx := metadata.NewOutgoingContext(c, md)

		domainResp, nextCursor, err := h.UserService.GetUsers(grpcCtx, int64(limit), cursor)
		if err != nil {
			status, ok := status.FromError(err)
			if !ok {
				apierror.ErrInternalServerError.APIError(c, err)
				return
			}

			var customErr *apierror.CustomError
			switch status.Code() {
			case codes.Unauthenticated:
				customErr = apierror.ErrUnauthorized
			default:
				customErr = apierror.ErrInternalServerError
			}
			customErr.APIError(c, err)
			return
		}

		resp := &dto.GetUsersResponse{
			Users:      make([]dto.User, len(domainResp)),
			NextCursor: nextCursor,
		}
		for i, user := range domainResp {
			resp.Users[i] = dto.User{
				ID:        user.ID,
				FirstName: user.FirstName,
				LastName:  user.LastName,
				Email:     user.Email,
				Token:     user.Token,
				Active:    user.Active,
				Admin:     user.Admin,
			}
		}

		h.ResponseJSON(c, http.StatusOK, resp)
	}
}

func (h *UserHandler) UpdateUser() gin.HandlerFunc {
	return func(c *gin.Context) {
		md, err := contextstore.GRPCMetadataFromContext(c)
		if err != nil {
			apierror.ErrInternalServerError.APIError(c, err)
			return
		}
		grpcCtx := metadata.NewOutgoingContext(c, md)

		var req dto.UpdateUserRequest
		if err := c.BindJSON(&req); err != nil {
			apierror.ErrBadRequest.APIError(c, nil)
			return
		}

		if err := h.ValidateStruct(req); err != nil {
			apierror.ErrBadRequest.APIError(c, err)
			return
		}

		user := domain.User{
			Password:  req.Password,
			FirstName: req.FirstName,
			LastName:  req.LastName,
			Email:     req.Email,
		}

		if err := h.UserService.UpdateUser(grpcCtx, user); err != nil {
			if status, ok := status.FromError(err); ok {
				errorCode := status.Code()
				switch int32(errorCode) {
				case 5:
					apierror.ErrBadRequest.APIError(c, nil)
				default:
					apierror.ErrInternalServerError.APIError(c, nil)
				}
				return
			}
		}

		h.ResponseNoContent(c, http.StatusNoContent)
	}
}
