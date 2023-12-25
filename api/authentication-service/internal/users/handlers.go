package users

import (
	"log"
	"net/http"

	"github.com/LeonLow97/utils"
	"github.com/go-playground/validator/v10"
)

type userHandler struct {
	service Service
}

func NewUserHandler(s Service) *userHandler {
	return &userHandler{
		service: s,
	}
}

func (h userHandler) UpdateUser(w http.ResponseWriter, r *http.Request) {
	var updateUserReq UpdateUserRequestDTO
	if err := utils.ReadJSON(w, r, &updateUserReq); err != nil {
		log.Println(err)
		utils.ErrorJSON(w, http.StatusInternalServerError)
		return
	}

	// validate json
	// create an instance of the validator
	validate := validator.New()
	if err := validate.Struct(updateUserReq); err != nil {
		for _, err := range err.(validator.ValidationErrors) {
			log.Println("Validation Error:", err)
		}
		utils.ErrorJSON(w, http.StatusBadRequest)
		return
	}

	// sanitize request
	updateUserSanitize(&updateUserReq)

	err := h.service.UpdateUser(updateUserReq)
	switch {
	case err != nil:
		utils.ErrorJSON(w, http.StatusInternalServerError)
	}
}

func (h userHandler) GetUsers(w http.ResponseWriter, r *http.Request) {
	users, err := h.service.GetUsers()
	if err != nil {
		log.Println(err)
		return
	}

	_ = utils.WriteJSON(w, http.StatusOK, users)
}
