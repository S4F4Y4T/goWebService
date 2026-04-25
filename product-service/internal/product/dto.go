package product

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
