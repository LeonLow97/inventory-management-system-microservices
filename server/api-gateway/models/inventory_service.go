package models

type CreateProductRequest struct {
	BrandName    string `json:"brand_name"`
	CategoryName string `json:"category_name"`
	ProductName  string `json:"product_name"`
	Description  string `json:"description"`
	Size         string `json:"size"`
	Color        string `json:"color"`
	Quantity     int32  `json:"quantity"`
}

type UpdateProductRequest struct {
	BrandName    string `json:"brand_name"`
	CategoryName string `json:"category_name"`
	ProductName  string `json:"product_name"`
	Description  string `json:"description"`
	Size         string `json:"size"`
	Color        string `json:"color"`
	Quantity     int32  `json:"quantity"`
}