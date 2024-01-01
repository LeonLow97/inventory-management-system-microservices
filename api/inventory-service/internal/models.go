package inventory

type GetProductByIdDTO struct {
	UserID    int
	ProductID int
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
