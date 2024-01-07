package models

type CreateOrderRequest struct {
	CustomerName string `json:"customer_name" validate:"omitempty,min=10,max=255"`
	ProductName  string `json:"product_name" validate:"required,min=1,max=100"`
	BrandName    string `json:"brand_name" validate:"required,min=1,max=100"`
	CategoryName string `json:"category_name" validate:"required,min=1,max=100"`
	Color        string `json:"color" validate:"omitempty,min=1,max=50"`
	Size         string `json:"size" validate:"omitempty,min=1,max=50"`
	Quantity     int64  `json:"quantity" validate:"required,min=1"`
	Description  string `json:"description" validate:"omitempty"`
	Revenue      int64  `json:"revenue" validate:"required,min=1"`
	Cost         int64  `json:"cost" validate:"required,min=1"`
	Profit       int64  `json:"profit" validate:"required,min=1"`
	HasReviewed  bool   `json:"has_reviewed" validate:"required"`
}
