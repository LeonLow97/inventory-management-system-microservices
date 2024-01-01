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

func createProductSanitize(createProductDTO *CreateProductDTO) {
	createProductDTO.BrandName = strings.TrimSpace(createProductDTO.BrandName)
	createProductDTO.CategoryName = strings.TrimSpace(createProductDTO.CategoryName)
	createProductDTO.ProductName = strings.TrimSpace(createProductDTO.ProductName)
	createProductDTO.Description = strings.TrimSpace(createProductDTO.Description)
	createProductDTO.Size = strings.TrimSpace(createProductDTO.Size)
	createProductDTO.Color = strings.TrimSpace(createProductDTO.Color)
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
