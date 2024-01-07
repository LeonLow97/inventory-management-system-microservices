package order

import (
	"context"
	"database/sql"
	"time"

	"github.com/jmoiron/sqlx"
)

type Repository interface {
	GetOrders(req GetOrdersDTO) (*[]Order, error)
	GetOrderByID(req GetOrderDTO) (*Order, error)
	CreateOrder(req CreateOrderDTO, productID int) error
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
		if err == sql.ErrNoRows {
			return nil, ErrNoOrdersFound
		}
		return nil, err
	}

	return &orders, nil
}

func (r repo) GetOrderByID(req GetOrderDTO) (*Order, error) {
	ctx, cancel := context.WithTimeout(context.Background(), time.Minute)
	defer cancel()

	query := `
		SELECT id, product_id, customer_name, brand_name, category_name, 
			color, size, quantity, description, revenue, cost, profit, has_reviewed, 
			created_at, updated_at
		FROM orders
		WHERE user_id = $1 AND id = $2;
	`

	var order Order
	if err := r.db.QueryRowContext(ctx, query, req.UserID, req.OrderID).Scan(
		&order.OrderID,
		&order.ProductID,
		&order.CustomerName,
		&order.BrandName,
		&order.CategoryName,
		&order.Color,
		&order.Size,
		&order.Quantity,
		&order.Description,
		&order.Revenue,
		&order.Cost,
		&order.Profit,
		&order.HasReviewed,
		&order.CreatedAt,
		&order.UpdatedAt,
	); err != nil {
		if err == sql.ErrNoRows {
			return nil, ErrNoOrderFound
		}
		return nil, err
	}

	return &order, nil
}

func (r repo) CreateOrder(req CreateOrderDTO, productID int) error {
	ctx, cancel := context.WithTimeout(context.Background(), time.Minute)
	defer cancel()

	query := `
		INSERT INTO orders (
			product_id, user_id, customer_name, brand_name, category_name, color, size, 
			quantity, description, revenue, cost, profit, has_reviewed, status, status_reason, order_uuid
		) VALUES (
			$1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12, $13, $14
		)
	`

	//cost, profit, has_reviewed, status
	_, err := r.db.ExecContext(ctx, query,
		productID,
		req.UserID,
		req.CustomerName,
		req.BrandName,
		req.CategoryName,
		req.Color,
		req.Size,
		req.Quantity,
		req.Description,
		req.Revenue,
		req.Cost,
		req.Profit,
		req.HasReviewed,
		req.Status,
		req.StatusReason,
		req.OrderUUID,
	)

	return err
}
