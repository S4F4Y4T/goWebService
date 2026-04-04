package model

import (
	"context"
	"time"
)

type Product struct {
	ID        uint      `json:"id" gorm:"primaryKey"`
	Name      string    `json:"name"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

type CreateProductRequest struct {
	Name string `json:"name"`
}

type UpdateProductRequest struct {
	ID   uint   `json:"id"`
	Name string `json:"name"`
}

type DeleteProductRequest struct {
	ID uint `json:"id"`
}

type GetProductRequest struct {
	ID uint `json:"id"`
}

type GetProductsRequest struct {
	Limit  int `json:"limit"`
	Offset int `json:"offset"`
}

type GetProductsResponse struct {
	Products []Product `json:"products"`
	Total    int64     `json:"total"`
	Limit    int       `json:"limit"`
	Offset   int       `json:"offset"`
}

type ProductRepository interface {
	Create(ctx context.Context, req *CreateProductRequest) (*Product, error)
	Update(ctx context.Context, req *UpdateProductRequest) (*Product, error)
	Delete(ctx context.Context, req *DeleteProductRequest) error
	FindByID(ctx context.Context, req *GetProductRequest) (*Product, error)
	FindAll(ctx context.Context, req *GetProductsRequest) (*GetProductsResponse, error)
}
