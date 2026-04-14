package product

// CreateProductRequest is the DTO for creating a new product
type CreateProductRequest struct {
	Name string `json:"name" validate:"required,min=2,max=255"`
}

// UpdateProductRequest is the DTO for updating an existing product
type UpdateProductRequest struct {
	ID   uint   `json:"id"`
	Name string `json:"name" validate:"required,min=2,max=255"`
}

// GetProductsResponse is the DTO for paginated product list responses
type GetProductsResponse struct {
	Products []Product `json:"products"`
	Total    int64     `json:"total"`
	Limit    int       `json:"limit"`
	Offset   int       `json:"offset"`
}
