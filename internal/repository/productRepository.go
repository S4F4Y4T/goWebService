package repository

import (
	"errors"
	"sync"

	"github.com/S4F4Y4T/goWebService/internal/model"
)

type productRepository struct {
	mu       sync.RWMutex
	products map[int]model.Product
	nextID   int
}

func NewProductRepository() model.ProductRepository {
	return &productRepository{
		products: make(map[int]model.Product),
		nextID:   1,
	}
}

func (r *productRepository) Create(req *model.CreateProductRequest) (*model.Product, error) {
	r.mu.Lock()
	defer r.mu.Unlock()

	product := model.Product{
		ID:   r.nextID,
		Name: req.Name,
	}
	r.products[r.nextID] = product
	r.nextID++

	return &product, nil
}

func (r *productRepository) Update(req *model.UpdateProductRequest) (*model.Product, error) {
	r.mu.Lock()
	defer r.mu.Unlock()

	product, ok := r.products[req.ID]
	if !ok {
		return nil, errors.New("product not found")
	}

	product.Name = req.Name
	r.products[req.ID] = product

	return &product, nil
}

func (r *productRepository) Delete(req *model.DeleteProductRequest) error {
	r.mu.Lock()
	defer r.mu.Unlock()

	if _, ok := r.products[req.ID]; !ok {
		return errors.New("product not found")
	}

	delete(r.products, req.ID)
	return nil
}

func (r *productRepository) FindByID(req *model.GetProductRequest) (*model.Product, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()

	product, ok := r.products[req.ID]
	if !ok {
		return nil, errors.New("product not found")
	}

	return &product, nil
}

func (r *productRepository) FindAll(req *model.GetProductsRequest) (*model.GetProductsResponse, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()

	products := make([]model.Product, 0, len(r.products))
	for _, product := range r.products {
		products = append(products, product)
	}

	return &model.GetProductsResponse{
		Products: products,
		Total:    len(products),
		Limit:    req.Limit,
		Offset:   req.Offset,
	}, nil
}
