package domain

import "strings"

type Product struct {
	ID           int64
	BrandName    string
	CategoryName string
	ProductName  string
	Description  string
	Size         string
	Color        string
	Quantity     int
	CreatedAt    string
	UpdatedAt    string
}

type Brand struct {
	ID        int
	Name      string
	CreatedAt string
}

type Category struct {
	ID        int
	Name      string
	CreatedAt string
}

func (p *Product) Sanitize() {
	p.BrandName = strings.TrimSpace(p.BrandName)
	p.CategoryName = strings.TrimSpace(p.CategoryName)
	p.ProductName = strings.TrimSpace(p.ProductName)
	p.Description = strings.TrimSpace(p.Description)
	p.Size = strings.TrimSpace(p.Size)
	p.Color = strings.TrimSpace(p.Color)
}

func (b *Brand) Sanitize() {
	b.Name = strings.TrimSpace(b.Name)
}

func (c *Category) Sanitize() {
	c.Name = strings.TrimSpace(c.Name)
}
