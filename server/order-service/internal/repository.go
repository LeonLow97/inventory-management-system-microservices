package order

import (
	"context"
	"time"

	"github.com/jmoiron/sqlx"
)

type Repository interface {
	GetOrders(req GetOrdersDTO) (*[]Order, error)
}

type repo struct {
	db *sqlx.DB
}

func NewRepository(db *sqlx.DB) Repository {
	return &repo{
		db: db,
	}
}

func (r repo) GetOrders(req GetOrdersDTO) (*[]Order, error) {
	ctx, cancel := context.WithTimeout(context.Background(), time.Minute)
	defer cancel()

	query := `
		SELECT id, product_id, customer_name, brand_name, category_name, 
			color, size, quantity, description, revenue, cost, profit, has_reviewed, 
			created_at, updated_at
		FROM orders
		WHERE user_id = $1;
	`

	var orders []Order
	if err := r.db.SelectContext(ctx, &orders, query, req.UserID); err != nil {
		return nil, err
	}

	// If no rows are returned, return an empty slice instead of nil
	if len(orders) == 0 {
		return nil, ErrNoOrdersFound
	}

	return &orders, nil
}
