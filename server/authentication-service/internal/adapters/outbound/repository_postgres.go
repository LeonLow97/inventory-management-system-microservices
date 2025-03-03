package outbound

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"strings"

	"github.com/LeonLow97/internal/core/domain"
	"github.com/LeonLow97/internal/core/services"
)

func (r *Repository) GetUserByEmail(ctx context.Context, email string) (*domain.User, error) {
	query := `
        SELECT 
            u.id, u.email, u.hashed_password, u.first_name, u.last_name, u.active,
			COALESCE(au.active, FALSE) as "admin", u.updated_at, u.created_at
        FROM users u
		LEFT JOIN admin_users au
		ON
			au.user_id = u.id
        WHERE
			email = $1
    `

	var user domain.User
	if err := r.db.GetContext(ctx, &user, query, email); err != nil {
		switch {
		case errors.Is(err, sql.ErrNoRows):
			return nil, services.ErrUserNotFound
		default:
			return nil, err
		}
	}
	return &user, nil
}

func (r *Repository) EmailExists(ctx context.Context, email string) (bool, error) {
	query := `
		SELECT EXISTS (
			SELECT 1 FROM users WHERE email = $1
		)
	`

	var exists bool
	if err := r.db.QueryRowContext(ctx, query, email).Scan(&exists); err != nil {
		return false, err
	}
	return exists, nil
}

// TODO: Convert this to cursor pagination
func (r *Repository) GetUsers() (*[]domain.User, error) {
	// Return an empty slice and no error
	return &[]domain.User{}, nil
}

func (r *Repository) InsertUser(ctx context.Context, user *domain.User) error {
	query := `
		INSERT INTO users (email, hashed_password, first_name, last_name)
		VALUES (?, ?, ?, ?)
    `

	args := []interface{}{user.Email, user.HashedPassword, user.FirstName, user.LastName}
	_, err := r.db.ExecContext(ctx, r.db.Rebind(query), args...)
	return err
}

func (r *Repository) UpdateUserByID(ctx context.Context, user *domain.User) error {
	var setClauses []string
	var args []interface{}

	// Dynamically add fields to update only if they are provided
	if user.FirstName != nil && *user.FirstName != "" {
		setClauses = append(setClauses, "first_name = ?")
		args = append(args, *user.FirstName)
	}
	if user.LastName != nil && *user.LastName != "" {
		setClauses = append(setClauses, "last_name = ?")
		args = append(args, *user.LastName)
	}
	if user.HashedPassword != nil && *user.HashedPassword != "" {
		setClauses = append(setClauses, "password = ?")
		args = append(args, *user.HashedPassword)
	}

	// Return early if there are no updates
	if len(setClauses) == 0 {
		return nil
	}

	query := fmt.Sprintf(`
		UPDATE users
		SET
			%s, updated_at = NOW()
		WHERE id = ?
	`, strings.Join(setClauses, ", "))

	args = append(args, user.ID)
	_, err := r.db.ExecContext(ctx, r.db.Rebind(query), args...)
	return err
}
