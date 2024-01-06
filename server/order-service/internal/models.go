package order

type GetOrdersDTO struct {
	UserID int
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
