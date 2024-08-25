package domain

type Order struct {
	OrderID      int64
	ProductID    int64
	ProductName  string
	CustomerName string
	BrandName    string
	CategoryName string
	Color        string
	Size         string
	Quantity     int64
	Description  string
	Revenue      float32
	Cost         float32
	Profit       float32
	HasReviewed  bool
	Status       string
	StatusReason string
	CreatedAt    string
}
