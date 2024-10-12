package domain

import "strings"

type Order struct {
	ID           int    `db:"id"`
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
	Status       string `db:"status"`
	StatusReason string `db:"status_reason"`
	HasReviewed  bool   `db:"has_reviewed"`
	UpdatedAt    string `db:"updated_at"`
	CreatedAt    string `db:"created_at"`
}

// Sanitize trims leading and trailing whitespace from string fields in the Order struct.
func (o *Order) Sanitize() {
	o.CustomerName = strings.TrimSpace(o.CustomerName)
	o.BrandName = strings.TrimSpace(o.BrandName)
	o.CategoryName = strings.TrimSpace(o.CategoryName)
	o.Color = strings.TrimSpace(o.Color)
	o.Size = strings.TrimSpace(o.Size)
	o.Description = strings.TrimSpace(o.Description)
	o.Status = strings.TrimSpace(o.Status)
	o.StatusReason = strings.TrimSpace(o.StatusReason)
	o.UpdatedAt = strings.TrimSpace(o.UpdatedAt)
	o.CreatedAt = strings.TrimSpace(o.CreatedAt)
}
