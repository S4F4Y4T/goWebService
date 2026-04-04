package model

import "time"

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
	Create(req *CreateProductRequest) (*Product, error)
	Update(req *UpdateProductRequest) (*Product, error)
	Delete(req *DeleteProductRequest) error
	FindByID(req *GetProductRequest) (*Product, error)
	FindAll(req *GetProductsRequest) (*GetProductsResponse, error)
}
