package dto

type GetProductDetailsDTO struct {
	UserID       int
	BrandName    string
	CategoryName string
	ProductName  string
}

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

type Product struct {
	ID           int64
	BrandName    string
	CategoryName string
	ProductName  string
	Description  string
	Size         string
	Color        string
	Quantity     int
	CreatedAt    string
	UpdatedAt    string
}

type Brand struct {
	ID        int
	Name      string
	CreatedAt string
}

type Category struct {
	ID        int
	Name      string
	CreatedAt string
}
