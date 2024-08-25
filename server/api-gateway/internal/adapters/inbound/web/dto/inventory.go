package dto

type Product struct {
	BrandName    string `json:"brand_name,omitempty"`
	CategoryName string `json:"category_name,omitempty"`
	ProductName  string `json:"product_name,omitempty"`
	Description  string `json:"description,omitempty"`
	Size         string `json:"size,omitempty"`
	Color        string `json:"color,omitempty"`
	Quantity     int32  `json:"quantity,omitempty"`
	CreatedAt    string `json:"created_at,omitempty"`
	UpdatedAt    string `json:"updated_at,omitempty"`
}

type GetProductsResponse struct {
	Products []Product `json:"products"`
}

type GetProductByIDResponse struct {
	Product
}

type CreateProductRequest struct {
	BrandName    string `json:"brand_name" validate:"required,min=1,max=100"`
	CategoryName string `json:"category_name" validate:"required,min=1,max=100"`
	ProductName  string `json:"product_name" validate:"required,min=1,max=100"`
	Description  string `json:"description" validate:"max=500"`
	Size         string `json:"size" validate:"required,max=50"`
	Color        string `json:"color" validate:"required,max=50"`
	Quantity     int32  `json:"quantity" validate:"required,min=1"`
}

type UpdateProductRequest struct {
	BrandName    string `json:"brand_name" validate:"required,min=1,max=100"`
	CategoryName string `json:"category_name" validate:"required,min=1,max=100"`
	ProductName  string `json:"product_name" validate:"omitempty,min=1,max=100"`
	Description  string `json:"description" validate:"omitempty,max=500"`
	Size         string `json:"size" validate:"omitempty,max=50"`
	Color        string `json:"color" validate:"omitempty,max=50"`
	Quantity     int32  `json:"quantity" validate:"omitempty,min=1"`
}
