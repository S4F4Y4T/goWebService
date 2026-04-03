package service

import (
	"github.com/S4F4Y4T/goWebService/internal/model"
)

type ProductService struct {
	repo model.ProductRepository
}

func NewProductService(repo model.ProductRepository) *ProductService {
	return &ProductService{repo: repo}
}

func (s *ProductService) Create(req *model.CreateProductRequest) (*model.Product, error) {
	return s.repo.Create(req)
}

func (s *ProductService) Update(req *model.UpdateProductRequest) (*model.Product, error) {
	return s.repo.Update(req)
}

func (s *ProductService) Delete(req *model.DeleteProductRequest) error {
	return s.repo.Delete(req)
}

func (s *ProductService) FindByID(req *model.GetProductRequest) (*model.Product, error) {
	return s.repo.FindByID(req)
}

func (s *ProductService) FindAll() (*model.GetProductsResponse, error) {
	return s.repo.FindAll(&model.GetProductsRequest{Limit: 10, Offset: 0})
}
