package inventory

import (
	"context"
	"database/sql"
	"time"
)

type Repository interface {
	GetProducts(userID int) (*[]Product, error)
}

type MySQLRepo struct {
	db *sql.DB
}

func NewRepository(db *sql.DB) Repository {
	return &MySQLRepo{
		db: db,
	}
}

func (r *MySQLRepo) GetProducts(userID int) (*[]Product, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 15*time.Second)
	defer cancel()

	query := `
		SELECT b.name AS brand_name, c.name AS category_name, p.name AS product_name,
			p.description, p.size, p.color, p.quantity, p.created_at, p.updated_at
		FROM products p
		JOIN brands b ON b.id = p.brand_id
		JOIN categories c ON c.id = p.category_id
		WHERE p.user_id = ?;
	`

	rows, err := r.db.QueryContext(ctx, query, userID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var products []Product
	for rows.Next() {
		var product Product
		err := rows.Scan(
			&product.BrandName,
			&product.CategoryName,
			&product.ProductName,
			&product.Description,
			&product.Size,
			&product.Color,
			&product.Quantity,
			&product.CreatedAt,
			&product.UpdatedAt,
		)
		if err != nil {
			return nil, err
		}
		products = append(products, product)
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}

	return &products, nil
}
