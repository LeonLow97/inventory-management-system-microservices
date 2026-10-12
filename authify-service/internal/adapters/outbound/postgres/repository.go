package postgres

import (
	"github.com/LeonLow97/internal/core/ports"
	"github.com/jmoiron/sqlx"
)

type UserRepository struct {
	db *sqlx.DB
}

func NewUserRepository(db *sqlx.DB) ports.UserRepository {
	return &UserRepository{
		db: db,
	}
}
