package users

import (
	"strings"
	"time"
)

type UpdateUserRequestDTO struct {
	FirstName string
	LastName  string
	Username  string
	Password  string
	Email     string
}

func updateUserSanitize(req *UpdateUserRequestDTO) {
	req.FirstName = strings.TrimSpace(req.FirstName)
	req.LastName = strings.TrimSpace(req.LastName)
	req.Username = strings.TrimSpace(req.Username)
	req.Password = strings.TrimSpace(req.Password)
	req.Email = strings.TrimSpace(req.Email)
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
