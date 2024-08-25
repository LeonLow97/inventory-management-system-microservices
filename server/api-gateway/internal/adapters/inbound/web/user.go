package web

import (
	"fmt"
	"log"
	"net/http"

	"github.com/LeonLow97/internal/adapters/inbound/web/dto"
	"github.com/LeonLow97/internal/core/domain"
	user "github.com/LeonLow97/internal/core/services/user"
	"github.com/LeonLow97/internal/pkg/utils"
	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
	"google.golang.org/grpc/status"
)

type UserHandler struct {
	UserService user.User
}

func NewUserHandler(userService user.User) *UserHandler {
	return &UserHandler{
		UserService: userService,
	}
}

func (h *UserHandler) GetUsers() gin.HandlerFunc {
	return func(c *gin.Context) {
		domainResp, err := h.UserService.GetUsers(c)
		if err != nil {
			log.Println(err)
			c.JSON(http.StatusInternalServerError, gin.H{"status": http.StatusInternalServerError, "message": "Internal Server Error"})
			return
		}

		resp := &dto.GetUsersResponse{
			Users: make([]dto.User, len(*domainResp)),
		}
		for i, user := range *domainResp {
			resp.Users[i] = dto.User{
				Username:  user.Username,
				FirstName: user.FirstName,
				LastName:  user.LastName,
				Email:     user.Email,
				Token:     user.Token,
				Active:    user.Active,
				Admin:     user.Admin,
			}
		}

		c.JSON(http.StatusOK, gin.H{"status": http.StatusOK, "users": resp})
	}
}

func (h *UserHandler) UpdateUser() gin.HandlerFunc {
	return func(c *gin.Context) {
		userID, err := utils.GetUserIDFromContext(c)
		if err != nil {
			c.JSON(http.StatusUnauthorized, gin.H{"status": http.StatusUnauthorized, "message": "Unauthorized"})
			return
		}

		var req dto.UpdateUserRequest
		if err := c.BindJSON(&req); err != nil {
			log.Println(err)
			c.JSON(http.StatusBadRequest, gin.H{"status": http.StatusBadRequest, "message": "Bad Request"})
			return
		}

		validate := validator.New()
		if err := validate.Struct(req); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"status": http.StatusBadRequest, "message": fmt.Sprintf("Bad Request: %s", err.Error())})
			return
		}

		user := domain.User{
			Username:  req.Username,
			Password:  req.Password,
			FirstName: req.FirstName,
			LastName:  req.LastName,
			Email:     req.Email,
		}

		if err := h.UserService.UpdateUser(c, user, userID); err != nil {
			if status, ok := status.FromError(err); ok {
				errorCode := status.Code()
				switch int32(errorCode) {
				case 3:
					log.Println(err)
					c.JSON(http.StatusBadRequest, gin.H{"status": http.StatusBadRequest, "message": fmt.Sprintf("Bad Request: %s", status.Message())})
				case 5:
					log.Println(err)
					c.JSON(http.StatusBadRequest, gin.H{"status": http.StatusBadRequest, "message": "Bad Request"})
				default:
					log.Println(err)
					c.JSON(http.StatusInternalServerError, gin.H{"status": http.StatusInternalServerError, "message": "Internal Server Error"})
				}
				return
			}
		}
		c.JSON(http.StatusNoContent, gin.H{})
	}
}
