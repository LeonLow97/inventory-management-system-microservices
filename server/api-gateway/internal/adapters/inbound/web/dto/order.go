package dto

type Order struct {
	OrderID      int64   `json:"order_id,omitempty"`
	ProductID    int64   `json:"product_id,omitempty"`
	ProductName  string  `json:"product_name,omitempty"`
	CustomerName string  `json:"customer_name,omitempty"`
	BrandName    string  `json:"brand_name,omitempty"`
	CategoryName string  `json:"category_name,omitempty"`
	Color        string  `json:"color,omitempty"`
	Size         string  `json:"size,omitempty"`
	Quantity     int64   `json:"quantity,omitempty"`
	Description  string  `json:"description,omitempty"`
	Revenue      float32 `json:"revenue,omitempty"`
	Cost         float32 `json:"cost,omitempty"`
	Profit       float32 `json:"profit,omitempty"`
	HasReviewed  bool    `json:"has_reviewed,omitempty"`
	Status       string  `json:"status,omitempty"`
	StatusReason string  `json:"status_reason,omitempty"`
	CreatedAt    string  `json:"created_at,omitempty"`
}

type GetOrdersResponse struct {
	Orders []Order `json:"orders"`
}

type GetOrderResponse struct {
	Order
}

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
