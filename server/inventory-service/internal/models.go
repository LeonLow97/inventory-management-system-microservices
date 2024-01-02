package inventory

import "strings"

type GetProductByIdDTO struct {
	UserID    int
	ProductID int
}

type CreateProductDTO struct {
	UserID       int
	BrandName    string
	CategoryName string
	ProductName  string
	Description  string
	Size         string
	Color        string
	Quantity     int
}

type UpdateProductDTO struct {
	UserID       int
	ProductID    int
	BrandName    string
	BrandID      int
	CategoryName string
	CategoryID   int
	ProductName  string
	Description  string
	Size         string
	Color        string
	Quantity     int
}

type DeleteProductDTO struct {
	UserID    int
	ProductID int
}

func createProductSanitize(createProductDTO *CreateProductDTO) {
	createProductDTO.BrandName = strings.TrimSpace(createProductDTO.BrandName)
	createProductDTO.CategoryName = strings.TrimSpace(createProductDTO.CategoryName)
	createProductDTO.ProductName = strings.TrimSpace(createProductDTO.ProductName)
	createProductDTO.Description = strings.TrimSpace(createProductDTO.Description)
	createProductDTO.Size = strings.TrimSpace(createProductDTO.Size)
	createProductDTO.Color = strings.TrimSpace(createProductDTO.Color)
}

func updateProductSanitize(updateProductDTO *UpdateProductDTO) {
	updateProductDTO.BrandName = strings.TrimSpace(updateProductDTO.BrandName)
	updateProductDTO.CategoryName = strings.TrimSpace(updateProductDTO.CategoryName)
	updateProductDTO.ProductName = strings.TrimSpace(updateProductDTO.ProductName)
	updateProductDTO.Description = strings.TrimSpace(updateProductDTO.Description)
	updateProductDTO.Size = strings.TrimSpace(updateProductDTO.Size)
	updateProductDTO.Color = strings.TrimSpace(updateProductDTO.Color)
}

type Product struct {
	BrandName    string `db:"brand_name"`
	CategoryName string `db:"category_name"`
	ProductName  string `db:"product_name"`
	Description  string `db:"description"`
	Size         string `db:"size"`
	Color        string `db:"color"`
	Quantity     int    `db:"quantity"`
	CreatedAt    string `db:"created_at"`
	UpdatedAt    string `db:"updated_at"`
}

type Brand struct {
	ID        int    `db:"id"`
	Name      string `db:"name"`
	CreatedAt string `db:"created_at"`
}

type Category struct {
	ID        int    `db:"id"`
	Name      string `db:"name"`
	CreatedAt string `db:"created_at"`
}
