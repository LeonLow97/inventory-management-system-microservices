package order

type OrderEvent struct {
	OrderUUID string `json:"order_uuid"`
	UserID    int    `json:"user_id"`
	Quantity  int    `json:"quantity"`
}

type GetOrdersDTO struct {
	UserID int
}

type GetOrderDTO struct {
	UserID  int
	OrderID int
}

type CreateOrderDTO struct {
	UserID       int
	CustomerName string
	ProductName  string
	BrandName    string
	CategoryName string
	Color        string
	Size         string
	Quantity     int
	Description  string
	Revenue      int64
	Cost         int64
	Profit       int64
	HasReviewed  bool
	Status       string
	StatusReason string
	OrderUUID    string
}

type Order struct {
	OrderID      int    `db:"id"`
	ProductID    int    `db:"product_id"`
	CustomerName string `db:"customer_name"`
	BrandName    string `db:"brand_name"`
	CategoryName string `db:"category_name"`
	Color        string `db:"color"`
	Size         string `db:"size"`
	Quantity     int    `db:"quantity"`
	Description  string `db:"description"`
	Revenue      int64  `db:"revenue"`
	Cost         int64  `db:"cost"`
	Profit       int64  `db:"profit"`
	HasReviewed  bool   `db:"has_reviewed"`
	UpdatedAt    string `db:"updated_at"`
	CreatedAt    string `db:"created_at"`
}
