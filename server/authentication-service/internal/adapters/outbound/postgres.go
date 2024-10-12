package outbound

import (
	"context"
	"database/sql"

	"github.com/LeonLow97/internal/core/domain"
	"github.com/LeonLow97/internal/core/services"
)

func (r *Repository) GetUserByUsername(ctx context.Context, username string) (*domain.User, error) {
	query := `
        SELECT
            id,
            first_name,
            last_name,
            username,
            password,
            email,
            active,
            admin,
            updated_at,
            created_at
        FROM users
        WHERE username = $1
    `

	var user domain.User
	if err := r.db.GetContext(ctx, &user, query, username); err != nil {
		if err == sql.ErrNoRows {
			return nil, services.ErrUserNotFound
		}
		return nil, err
	}
	return &user, nil
}

func (r *Repository) GetUserByID(ctx context.Context, userID int) (*domain.User, error) {
	query := `
        SELECT
            id,
            first_name,
            last_name,
            username,
            password,
            email,
            active,
            admin,
            updated_at,
            created_at
        FROM users
        WHERE id = $1
    `

	var user domain.User
	if err := r.db.GetContext(ctx, &user, query, userID); err != nil {
		if err == sql.ErrNoRows {
			return nil, services.ErrUserNotFound
		}
		return nil, err
	}
	return &user, nil
}

func (r *Repository) IsUsernameTaken(ctx context.Context, username string) (bool, error) {
	query := `
        SELECT EXISTS (
            SELECT 1
            FROM users
            WHERE username = $1
        )
    `

	var exists bool
	if err := r.db.GetContext(ctx, &exists, query, username); err != nil {
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
        INSERT INTO users
            (first_name, last_name, username, password, email)
        VALUES
            (?, ?, ?, ?, ?)
    `

	var args []interface{}
	args = append(args, user.FirstName, user.LastName, user.Username, user.Password, user.Email)

	_, err := r.db.ExecContext(ctx, r.db.Rebind(query), args...)
	return err
}

func (r *Repository) UpdateUserByID(ctx context.Context, user *domain.User) error {
	query := `
        UPDATE users
        SET
            first_name = COALESCE(NULLIF(?, ''), first_name),
            last_name = COALESCE(NULLIF(?, ''), last_name),
            password = COALESCE(NULLIF(?, ''), password),
            email = COALESCE(NULLIF(?, ''), email),
            updated_at = now()
        WHERE id = ?
    `

	var args []interface{}
	args = append(args, user.FirstName, user.LastName, user.Password, user.Email, user.ID)

	_, err := r.db.ExecContext(ctx, r.db.Rebind(query), args...)
	return err
}
