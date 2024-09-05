package outbound

import (
	"context"
	"database/sql"
	"log"
	"time"

	"github.com/LeonLow97/internal/core/domain"
)

func (r *repository) GetProducts(userID int) (*[]domain.Product, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 15*time.Second)
	defer cancel()

	query := `
		SELECT b.name AS brand_name, c.name AS category_name, p.name AS product_name,
			p.description, p.size, p.color, p.quantity, p.created_at, p.updated_at
		FROM products p
		JOIN brands b ON b.id = p.brand_id
		JOIN categories c ON c.id = p.category_id
		WHERE p.user_id = ? AND p.is_deleted = 0;
	`

	rows, err := r.db.QueryContext(ctx, query, userID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var products []domain.Product
	for rows.Next() {
		var product domain.Product
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
		if err == sql.ErrNoRows {
			return nil, ErrProductsNotFound
		}
		return nil, err
	}

	return &products, nil
}

func (r *repository) GetProductByID(userID, productID int) (*domain.Product, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 15*time.Second)
	defer cancel()

	query := `
		SELECT b.name AS brand_name, c.name AS category_name, p.name AS product_name,
			p.description, p.size, p.color, p.quantity, p.created_at, p.updated_at
		FROM products p
		JOIN brands b ON b.id = p.brand_id
		JOIN categories c ON c.id = p.category_id
		WHERE p.user_id = ? AND p.id = ? AND p.is_deleted = 0;
	`

	var product domain.Product
	if err := r.db.QueryRowContext(
		ctx, query, userID, productID,
	).Scan(
		&product.BrandName,
		&product.CategoryName,
		&product.ProductName,
		&product.Description,
		&product.Size,
		&product.Color,
		&product.Quantity,
		&product.CreatedAt,
		&product.UpdatedAt,
	); err != nil {
		if err == sql.ErrNoRows {
			return nil, ErrProductNotFound
		}
		return nil, err
	}

	return &product, nil
}

func (r *repository) GetProductByName(userID int, productName string) (*domain.Product, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 15*time.Second)
	defer cancel()

	query := `
		SELECT p.id, b.name AS brand_name, c.name AS category_name, p.name AS product_name,
			p.description, p.size, p.color, p.quantity, p.created_at, p.updated_at
		FROM products p
		JOIN brands b ON b.id = p.brand_id
		JOIN categories c ON c.id = p.category_id
		WHERE p.user_id = ? AND p.name = ? AND p.is_deleted = 0;
	`

	var product domain.Product
	if err := r.db.QueryRowContext(
		ctx, query, userID, productName,
	).Scan(
		&product.ID,
		&product.BrandName,
		&product.CategoryName,
		&product.ProductName,
		&product.Description,
		&product.Size,
		&product.Color,
		&product.Quantity,
		&product.CreatedAt,
		&product.UpdatedAt,
	); err != nil {
		if err == sql.ErrNoRows {
			return nil, ErrProductNotFound
		}
		return nil, err
	}

	return &product, nil
}

func (r *repository) GetBrandByName(brandName string) (*domain.Brand, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 15*time.Second)
	defer cancel()

	query := `
		SELECT id, name, created_at FROM brands WHERE name = ?;
	`

	var brand domain.Brand
	if err := r.db.QueryRowContext(ctx, query, brandName).Scan(
		&brand.ID,
		&brand.Name,
		&brand.CreatedAt,
	); err != nil {
		if err == sql.ErrNoRows {
			return nil, ErrBrandNotFound
		}
		log.Println("error getting brand by name")
		return nil, err
	}

	return &brand, nil
}

func (r *repository) GetCategoryByName(categoryName string) (*domain.Category, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 15*time.Second)
	defer cancel()

	query := `
		SELECT id, name, created_at FROM categories WHERE name = ?;
	`

	var category domain.Category
	if err := r.db.QueryRowContext(ctx, query, categoryName).Scan(
		&category.ID,
		&category.Name,
		&category.CreatedAt,
	); err != nil {
		if err == sql.ErrNoRows {
			return nil, ErrCategoryNotFound
		}
		log.Println("error getting category by name")
		return nil, err
	}

	return &category, nil
}

func (r *repository) CreateProduct(req domain.Product, userID, brandID, categoryID int) error {
	ctx, cancel := context.WithTimeout(context.Background(), 15*time.Second)
	defer cancel()

	query := `
		INSERT INTO products (user_id, brand_id, category_id, name, description, size, color, quantity)
		VALUES (
			?, ?, ?, ?, ?, ?, ?, ?
		);
	`

	_, err := r.db.ExecContext(ctx, query,
		userID,
		brandID,
		categoryID,
		req.ProductName,
		req.Description,
		req.Size,
		req.Color,
		req.Quantity)
	if err != nil {
		log.Println("Error creating product", err)
		return err
	}

	return nil
}

func (r *repository) UpdateProductByID(req domain.Product, brandID, categoryID, userID, productID int) error {
	ctx, cancel := context.WithTimeout(context.Background(), 15*time.Second)
	defer cancel()

	query := `
		UPDATE products
		SET
			brand_id = COALESCE(NULLIF(?, 0), brand_id),
			category_id = COALESCE(NULLIF(?, 0), category_id),
			name = COALESCE(NULLIF(?, ''), name),
			description = COALESCE(NULLIF(?, ''), description),
			size = COALESCE(NULLIF(?, ''), size),
			color = COALESCE(NULLIF(?, ''), color),
			quantity = COALESCE(NULLIF(?, 0), quantity)
		WHERE user_id = ? AND id = ?;
	`

	_, err := r.db.ExecContext(ctx, query,
		brandID,
		categoryID,
		req.ProductName,
		req.Description,
		req.Size,
		req.Color,
		req.Quantity,
		userID,
		productID,
	)
	if err != nil {
		log.Println("Error updating product", err)
		return err
	}

	return nil
}

func (r *repository) UpdateProductQuantityByID(quantity, userID, productID int) error {
	ctx, cancel := context.WithTimeout(context.Background(), 15*time.Second)
	defer cancel()

	query := `
		UPDATE products
		SET
			quantity = ?
		WHERE user_id = ? AND id = ?
	`

	_, err := r.db.ExecContext(ctx, query, quantity, userID, productID)
	return err
}

func (r *repository) DeleteProductByID(userID, productID int) error {
	ctx, cancel := context.WithTimeout(context.Background(), 15*time.Second)
	defer cancel()

	query := `
		UPDATE products
		SET
			is_deleted = 1
		WHERE user_id = ? AND id = ?;
	`

	result, err := r.db.ExecContext(ctx, query,
		userID,
		productID,
	)
	if err != nil {
		log.Println("Error performing soft delete on product", err)
		return err
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return err
	}
	if rowsAffected != 1 {
		return ErrProductNotFound
	}

	return nil
}
