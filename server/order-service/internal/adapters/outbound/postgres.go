package outbound

import (
	"context"
	"database/sql"
	"time"

	"github.com/LeonLow97/internal/core/domain"
)

func (r *Repository) GetOrders(userID int) (*[]domain.Order, error) {
	ctx, cancel := context.WithTimeout(context.Background(), time.Minute)
	defer cancel()

	query := `
		SELECT 
			id, 
			product_id, 
			customer_name, 
			brand_name, 
			category_name, 
			color, 
			size, 
			quantity, 
			description, 
			revenue, 
			cost, 
			profit, 
			has_reviewed, 
			status, 
			status_reason, 
			created_at, 
			updated_at
		FROM orders
		WHERE user_id = $1;
	`

	var orders []domain.Order
	if err := r.db.SelectContext(ctx, &orders, query, userID); err != nil {
		if err == sql.ErrNoRows {
			return nil, ErrOrdersNotFound
		}
		return nil, err
	}
	return &orders, nil
}

func (r *Repository) GetOrderByID(userID, orderID int) (*domain.Order, error) {
	ctx, cancel := context.WithTimeout(context.Background(), time.Minute)
	defer cancel()

	query := `
		SELECT 
			id, 
			product_id, 
			customer_name, 
			brand_name, 
			category_name, 
			color, 
			size, 
			quantity, 
			description, 
			revenue, 
			cost, 
			profit, 
			has_reviewed, 
			status, 
			status_reason, 
			created_at, 
			updated_at
		FROM orders
		WHERE user_id = $1 AND id = $2;
	`

	var order domain.Order
	if err := r.db.GetContext(ctx, &order, query, userID, orderID); err != nil {
		if err == sql.ErrNoRows {
			return nil, ErrOrderNotFound
		}
		return nil, err
	}
	return &order, nil
}

func (r *Repository) CreateOrder(req domain.Order, userID, productID int) (int, error) {
	ctx, cancel := context.WithTimeout(context.Background(), time.Minute)
	defer cancel()

	query := `
		INSERT INTO orders (
			product_id, 
			user_id, 
			customer_name, 
			brand_name, 
			category_name, 
			color, 
			size, 
			quantity, 
			description, 
			revenue, 
			cost, 
			profit, 
			has_reviewed, 
			status, 
			status_reason
		) VALUES (
			?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?
		) RETURNING id;
	`

	var orderID int
	if err := r.db.QueryRowContext(ctx, r.db.Rebind(query),
		productID,
		userID,
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
	).Scan(&orderID); err != nil {
		return 0, err
	}
	return orderID, nil
}

func (r *Repository) UpdateOrderByID(req domain.Order) error {
	ctx, cancel := context.WithTimeout(context.Background(), time.Minute)
	defer cancel()

	query := `
		UPDATE orders
		SET 	
			status = $1, 
			status_reason = $2
		WHERE id = $3;
	`

	_, err := r.db.ExecContext(ctx, query,
		req.Status,
		req.StatusReason,
		req.ID,
	)
	return err
}
