package product

import (
	"context"
	"time"
)

// Product is the Aggregate Root for the Product domain
type Product struct {
	ID        uint      `json:"id" gorm:"primaryKey"`
	Name      string    `json:"name"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

// ProductRepository defines the persistence contract for Products
type ProductRepository interface {
	Create(ctx context.Context, p *Product) error
	Update(ctx context.Context, p *Product) error
	Delete(ctx context.Context, id uint) error
	FindByID(ctx context.Context, id uint) (*Product, error)
	FindAll(ctx context.Context, limit, offset int) ([]Product, int64, error)
}

// API DTOs
type CreateProductRequest struct {
	Name string `json:"name" validate:"required,min=2,max=255"`
}

type UpdateProductRequest struct {
	ID   uint   `json:"id"`
	Name string `json:"name" validate:"required,min=2,max=255"`
}

type GetProductsResponse struct {
	Products []Product `json:"products"`
	Total    int64     `json:"total"`
	Limit    int       `json:"limit"`
	Offset   int       `json:"offset"`
}
