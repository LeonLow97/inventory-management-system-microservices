package domain

import (
	"strings"
	"time"
)

type User struct {
	ID        int       `db:"id"`
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

// Sanitize trims leading and trailing whitespace from string fields in the User struct.
func (u *User) Sanitize() {
	u.FirstName = strings.TrimSpace(u.FirstName)
	u.LastName = strings.TrimSpace(u.LastName)
	u.Username = strings.TrimSpace(u.Username)
	u.Password = strings.TrimSpace(u.Password) // Consider security implications
	u.Email = strings.TrimSpace(u.Email)
}
