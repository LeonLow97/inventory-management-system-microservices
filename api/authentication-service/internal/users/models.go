package users

import (
	"strings"
	"time"
)

type UpdateUserRequest struct {
	FirstName string `json:"first_name"`
	LastName  string `json:"last_name"`
	Username  string `json:"username" validate:"required,min=5,max=50"`
	Password  string `json:"password" validate:"required,min=8,max=20"`
	Email     string `json:"email" validate:"required,min=10,max=100"`
}

func updateUserSanitize(req *UpdateUserRequest) {
	req.FirstName = strings.TrimSpace(req.FirstName)
	req.LastName = strings.TrimSpace(req.LastName)
	req.Username = strings.TrimSpace(req.Username)
	req.Password = strings.TrimSpace(req.Password)
	req.Email = strings.TrimSpace(req.Email)
}

type User struct {
	FirstName string    `json:"firstName" db:"first_name"`
	LastName  string    `json:"lastName" db:"last_name"`
	Username  string    `json:"username" db:"username"`
	Password  string    `json:"-" db:"password"`
	Email     string    `json:"email" db:"email"`
	Active    int       `json:"active" db:"active"`
	Admin     int       `json:"admin" db:"admin"`
	UpdatedAt time.Time `json:"-" db:"updated_at"`
	CreatedAt time.Time `json:"-" db:"created_at"`
}
