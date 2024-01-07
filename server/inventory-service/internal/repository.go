package inventory

import (
	"context"
	"database/sql"
	"log"
	"time"
)

type Repository interface {
	GetProducts(userID int) (*[]Product, error)
	GetProductByID(getProductByIdDTO GetProductByIdDTO) (*Product, error)
	GetProductByName(req GetProductDetailsDTO) (*Product, error)

	GetBrandByName(brandName string) (*Brand, error)
	GetCategoryByName(categoryName string) (*Category, error)

	CreateProduct(createProductDTO CreateProductDTO, brandID, categoryID int) error

	UpdateProductByID(updateProductDTO UpdateProductDTO) error
	DeleteProductByID(req DeleteProductDTO) error
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
		WHERE p.user_id = ? AND p.is_deleted = 0;
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
		if err == sql.ErrNoRows {
			return nil, ErrProductsNotFound
		}
		return nil, err
	}

	return &products, nil
}

func (r *MySQLRepo) GetProductByID(getProductByIdDTO GetProductByIdDTO) (*Product, error) {
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

	var product Product
	if err := r.db.QueryRowContext(
		ctx, query, getProductByIdDTO.UserID, getProductByIdDTO.ProductID,
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

func (r *MySQLRepo) GetProductByName(req GetProductDetailsDTO) (*Product, error) {
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

	var product Product
	if err := r.db.QueryRowContext(
		ctx, query, req.UserID, req.ProductName,
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

func (r *MySQLRepo) GetBrandByName(brandName string) (*Brand, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 15*time.Second)
	defer cancel()

	query := `
		SELECT id, name, created_at FROM brands WHERE name = ?;
	`

	var brand Brand
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

func (r *MySQLRepo) GetCategoryByName(categoryName string) (*Category, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 15*time.Second)
	defer cancel()

	query := `
		SELECT id, name, created_at FROM categories WHERE name = ?;
	`

	var category Category
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

func (r *MySQLRepo) CreateProduct(createProductDTO CreateProductDTO, brandID, categoryID int) error {
	ctx, cancel := context.WithTimeout(context.Background(), 15*time.Second)
	defer cancel()

	query := `
		INSERT INTO products (user_id, brand_id, category_id, name, description, size, color, quantity)
		VALUES (
			?, ?, ?, ?, ?, ?, ?, ?
		);
	`

	_, err := r.db.ExecContext(ctx, query,
		createProductDTO.UserID,
		brandID,
		categoryID,
		createProductDTO.ProductName,
		createProductDTO.Description,
		createProductDTO.Size,
		createProductDTO.Color,
		createProductDTO.Quantity)
	if err != nil {
		log.Println("Error creating product", err)
		return err
	}

	return nil
}

func (r *MySQLRepo) UpdateProductByID(updateProductDTO UpdateProductDTO) error {
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
		updateProductDTO.BrandID,
		updateProductDTO.CategoryID,
		updateProductDTO.ProductName,
		updateProductDTO.Description,
		updateProductDTO.Size,
		updateProductDTO.Color,
		updateProductDTO.Quantity,
		updateProductDTO.UserID,
		updateProductDTO.ProductID,
	)
	if err != nil {
		log.Println("Error updating product", err)
		return err
	}

	return nil
}

func (r *MySQLRepo) DeleteProductByID(req DeleteProductDTO) error {
	ctx, cancel := context.WithTimeout(context.Background(), 15*time.Second)
	defer cancel()

	query := `
		UPDATE products
		SET
			is_deleted = 1
		WHERE user_id = ? AND id = ?;
	`

	result, err := r.db.ExecContext(ctx, query,
		req.UserID,
		req.ProductID,
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
