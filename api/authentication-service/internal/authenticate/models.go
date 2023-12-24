package authenticate

import (
	"strings"
	"time"
)

type LoginRequestDTO struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

type User struct {
	FirstName string    `db:"first_name"`
	LastName  string    `db:"last_name"`
	Username  string    `db:"username"`
	Password  string    `db:"password"`
	Email     string    `db:"email"`
	Active    int       `db:"active"`
	Admin     int       `db:"admin"`
	UpdatedAt time.Time `db:"updated_at"`
	CreatedAt time.Time `db:"created_at"`
}

func loginSanitize(req *LoginRequestDTO) {
	req.Username = strings.TrimSpace(req.Username)
	req.Password = strings.TrimSpace(req.Password)
}

type SignUpRequest struct {
	FirstName string `json:"first_name"`
	LastName  string `json:"last_name"`
	Username  string `json:"username" validate:"required,min=5,max=50"`
	Password  string `json:"password" validate:"required,min=8,max=20"`
	Email     string `json:"email" validate:"required,min=10,max=100"`
}

func signUpSanitize(req *SignUpRequest) {
	req.FirstName = strings.TrimSpace(req.FirstName)
	req.LastName = strings.TrimSpace(req.LastName)
	req.Username = strings.TrimSpace(req.Username)
	req.Password = strings.TrimSpace(req.Password)
	req.Email = strings.TrimSpace(req.Email)
}
