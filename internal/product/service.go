package product

import (
	"context"

	"github.com/S4F4Y4T/goWebService/internal/shared/event"
	"github.com/S4F4Y4T/goWebService/pkg/apperror"
)

type Service struct {
	repo       ProductRepository
	dispatcher *event.Dispatcher
}

func NewService(repo ProductRepository, dispatcher *event.Dispatcher) *Service {
	return &Service{
		repo:       repo,
		dispatcher: dispatcher,
	}
}

func (s *Service) Create(ctx context.Context, req *CreateProductRequest) (*Product, error) {
	p, err := NewProduct(req.Name)
	if err != nil {
		return nil, apperror.New(apperror.BadRequest, err.Error(), err)
	}

	if err := s.repo.Create(ctx, p); err != nil {
		return nil, apperror.New(apperror.Internal, "failed to create product", err)
	}
	return p, nil
}

func (s *Service) Update(ctx context.Context, req *UpdateProductRequest) (*Product, error) {
	p, err := s.repo.FindByID(ctx, req.ID)
	if err != nil {
		return nil, apperror.New(apperror.NotFound, "product not found", err)
	}

	if err := p.UpdateName(req.Name); err != nil {
		return nil, apperror.New(apperror.BadRequest, err.Error(), err)
	}

	if err := s.repo.Update(ctx, p); err != nil {
		return nil, apperror.New(apperror.Internal, "failed to update product", err)
	}
	return p, nil
}

func (s *Service) Delete(ctx context.Context, id uint) error {
	if err := s.repo.Delete(ctx, id); err != nil {
		return apperror.New(apperror.Internal, "failed to delete product", err)
	}
	return nil
}

func (s *Service) FindByID(ctx context.Context, id uint) (*Product, error) {
	p, err := s.repo.FindByID(ctx, id)
	if err != nil {
		return nil, apperror.New(apperror.NotFound, "product not found", err)
	}
	return p, nil
}

func (s *Service) FindAll(ctx context.Context, limit, offset int) (*GetProductsResponse, error) {
	products, total, err := s.repo.FindAll(ctx, limit, offset)
	if err != nil {
		return nil, apperror.New(apperror.Internal, "failed to list products", err)
	}

	return &GetProductsResponse{
		Products: products,
		Total:    total,
		Limit:    limit,
		Offset:   offset,
	}, nil
}
