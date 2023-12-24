package users

import (
	"context"
	"database/sql"
	"log"
	"time"

	"github.com/jmoiron/sqlx"
)

type Repository interface {
	GetUserByUsername(username string) (User, error)
	UpdateUserByUsername(req UpdateUserRequest) error
	GetUsers() (*[]User, error)
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

func (r PostgresRepo) UpdateUserByUsername(req UpdateUserRequest) error {
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*15)
	defer cancel()

	query := `
		UPDATE users
		SET first_name = $1, last_name = $2, password = $3, email = $4, updated_at = $5
		WHERE username = $6;
	`

	if _, err := r.db.ExecContext(ctx, query,
		req.FirstName,
		req.LastName,
		req.Password,
		req.Email,
		time.Now(),
		req.Username,
	); err != nil {
		return err
	}

	return nil
}

func (r PostgresRepo) GetUsers() (*[]User, error) {
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*15)
	defer cancel()

	query := `
		SELECT first_name, last_name, username, email, active, admin, updated_at, created_at
		FROM users;
	`

	var users []User

	if err := r.db.SelectContext(ctx, &users, query); err != nil {
		log.Println("Error in get users", err)
		return nil, err
	}

	return &users, nil
}
