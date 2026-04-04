package service

import (
	"context"

	"github.com/S4F4Y4T/goWebService/internal/model"
)

type ProductService struct {
	repo model.ProductRepository
}

func NewProductService(repo model.ProductRepository) *ProductService {
	return &ProductService{repo: repo}
}

func (s *ProductService) Create(ctx context.Context, req *model.CreateProductRequest) (*model.Product, error) {
	return s.repo.Create(ctx, req)
}

func (s *ProductService) Update(ctx context.Context, req *model.UpdateProductRequest) (*model.Product, error) {
	return s.repo.Update(ctx, req)
}

func (s *ProductService) Delete(ctx context.Context, req *model.DeleteProductRequest) error {
	return s.repo.Delete(ctx, req)
}

func (s *ProductService) FindByID(ctx context.Context, req *model.GetProductRequest) (*model.Product, error) {
	return s.repo.FindByID(ctx, req)
}

func (s *ProductService) FindAll(ctx context.Context) (*model.GetProductsResponse, error) {
	return s.repo.FindAll(ctx, &model.GetProductsRequest{Limit: 10, Offset: 0})
}
