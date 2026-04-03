package model

type Product struct {
	ID   int    `json:"id"`
	Name string `json:"name"`
}

type CreateProductRequest struct {
	Name string `json:"name"`
}

type UpdateProductRequest struct {
	ID   int    `json:"id"`
	Name string `json:"name"`
}

type DeleteProductRequest struct {
	ID int `json:"id"`
}

type GetProductRequest struct {
	ID int `json:"id"`
}

type GetProductsRequest struct {
	Limit  int `json:"limit"`
	Offset int `json:"offset"`
}

type GetProductsResponse struct {
	Products []Product `json:"products"`
	Total    int       `json:"total"`
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
