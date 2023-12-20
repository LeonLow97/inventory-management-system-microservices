package authenticate

import (
	"context"
	"database/sql"
	"time"

	"github.com/jmoiron/sqlx"
)

type Repository interface {
	GetUserByUsername(username string) (User, error)
	GetUserCountByUsername(username string) (int, error)
	InsertOneUser(user SignUpRequest) error
}

type PostgresRepo struct {
	db *sqlx.DB
}

func NewRepo(db *sqlx.DB) Repository {
	return &PostgresRepo{
		db: db,
	}
}

func (r PostgresRepo) GetUserByUsername(username string) (User, error) {
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*15)
	defer cancel()

	var user User
	query := `
		SELECT first_name, last_name, username, password, email, active, admin, updated_at, created_at
		FROM users
		WHERE username = $1;
	`

	if err := r.db.GetContext(ctx, &user, query, username); err != nil {
		if err == sql.ErrNoRows {
			// User with the specified username was not found
			return user, ErrNotFound
		}
		// Return other errors encountered during the query execution
		return user, err
	}

	// User found
	return user, nil
}

func (r PostgresRepo) GetUserCountByUsername(username string) (int, error) {
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*15)
	defer cancel()

	var count int
	query := `
		SELECT COUNT(username) 
		FROM users
		WHERE username = $1;
	`

	if err := r.db.QueryRowContext(ctx, query, username).Scan(&count); err != nil {
		return 0, err
	}

	return count, nil
}

func (r PostgresRepo) InsertOneUser(user SignUpRequest) error {
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*15)
	defer cancel()

	query := `
		INSERT INTO users (first_name, last_name, username, password, email)
		VALUES ($1, $2, $3, $4, $5);
	`

	if _, err := r.db.ExecContext(ctx, query,
		user.FirstName,
		user.LastName,
		user.Username,
		user.Password,
		user.Email,
	); err != nil {
		return err
	}

	return nil
}
