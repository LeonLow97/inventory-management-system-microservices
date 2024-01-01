package users

import (
	"strings"
	"time"
)

type UpdateUserRequestDTO struct {
	UserID    int
	FirstName string
	LastName  string
	Password  string
	Email     string
}

func updateUserSanitize(req *UpdateUserRequestDTO) {
	req.FirstName = strings.TrimSpace(req.FirstName)
	req.LastName = strings.TrimSpace(req.LastName)
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
