package order

import "github.com/jmoiron/sqlx"

type Repository interface {
}

type repo struct {
	db *sqlx.DB
}

func NewRepository(db *sqlx.DB) Repository {
	return &repo{
		db: db,
	}
}
