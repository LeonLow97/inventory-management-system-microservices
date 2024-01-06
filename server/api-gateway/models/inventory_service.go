package models

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
