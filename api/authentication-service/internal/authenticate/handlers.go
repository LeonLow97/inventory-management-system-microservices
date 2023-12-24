package authenticate

import (
	"errors"
	"log"
	"net/http"

	"github.com/LeonLow97/utils"
	"github.com/go-playground/validator/v10"
)

type authenticateHandler struct {
	service Service
}

func NewAuthenticateHandler(s Service) *authenticateHandler {
	return &authenticateHandler{
		service: s,
	}
}

func (h authenticateHandler) Login(w http.ResponseWriter, r *http.Request) {
	var loginRequest LoginRequestDTO
	if err := utils.ReadJSON(w, r, &loginRequest); err != nil {
		log.Println(err)
		utils.ErrorJSON(w, http.StatusInternalServerError)
		return
	}

	// validate json
	// create an instance of the validator
	validate := validator.New()
	if err := validate.Struct(loginRequest); err != nil {
		for _, err := range err.(validator.ValidationErrors) {
			log.Println("Validation Error:", err)
		}
		utils.ErrorJSON(w, http.StatusUnauthorized)
		return
	}

	// sanitize data
	loginSanitize(&loginRequest)

	// call Login service (business logic)
	user, _, err := h.service.Login(loginRequest)
	switch {
	case errors.Is(err, ErrInvalidCredentials), errors.Is(err, ErrInactiveUser), errors.Is(err, ErrNotFound):
		utils.ErrorJSON(w, http.StatusUnauthorized)
	case err != nil:
		utils.ErrorJSON(w, http.StatusInternalServerError)
	default:
		// set cookie with the signed jwt token in browser
		// http.SetCookie(w, cookie)
		_ = utils.WriteJSON(w, http.StatusOK, user)
	}
}

func (h authenticateHandler) SignUp(w http.ResponseWriter, r *http.Request) {
	var signUpRequest SignUpRequest
	if err := utils.ReadJSON(w, r, &signUpRequest); err != nil {
		log.Println(err)
		utils.ErrorJSON(w, http.StatusInternalServerError)
		return
	}

	// validate json
	// create an instance of the validator
	validate := validator.New()
	if err := validate.Struct(signUpRequest); err != nil {
		for _, err := range err.(validator.ValidationErrors) {
			log.Println("Validation Error:", err)
		}
		utils.ErrorJSON(w, http.StatusBadRequest)
		return
	}

	// sanitize data
	signUpSanitize(&signUpRequest)

	// call SignUp service
	err := h.service.SignUp(signUpRequest)
	switch {
	case errors.Is(err, ErrInvalidEmailFormat), errors.Is(err, ErrInvalidPasswordFormat), errors.Is(err, ErrExistingUserFound):
		utils.ErrorJSON(w, http.StatusBadRequest)
	case err != nil:
		utils.ErrorJSON(w, http.StatusInternalServerError)
	default:
		utils.WriteJSON(w, http.StatusCreated, nil)
	}
}
