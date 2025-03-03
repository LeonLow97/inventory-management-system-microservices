package web

import (
	"net/http"

	"github.com/LeonLow97/internal/adapters/inbound/web/dto"
	"github.com/LeonLow97/internal/core/domain"
	user "github.com/LeonLow97/internal/core/services/user"
	"github.com/LeonLow97/internal/pkg/apierror"
	"github.com/LeonLow97/internal/pkg/handler"
	"github.com/LeonLow97/internal/pkg/utils"
	"github.com/gin-gonic/gin"
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
		domainResp, err := h.UserService.GetUsers(c)
		if err != nil {
			apierror.ErrInternalServerError.APIError(c, err)
			return
		}

		resp := &dto.GetUsersResponse{
			Users: make([]dto.User, len(*domainResp)),
		}
		for i, user := range *domainResp {
			resp.Users[i] = dto.User{
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
		userID, err := utils.GetUserIDFromContext(c)
		if err != nil {
			apierror.ErrUnauthorized.APIError(c, err)
			return
		}

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

		if err := h.UserService.UpdateUser(c, user, userID); err != nil {
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
