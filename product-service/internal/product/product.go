package product

import (
	"context"
	"errors"
	"strings"
	"time"
)

var (
	ErrInvalidProductName = errors.New("product name must be between 2 and 255 characters")
)

type Product struct {
	ID        uint      `json:"id"`
	Name      string    `json:"name"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

func NewProduct(name string) (*Product, error) {
	p := &Product{}
	if err := p.UpdateName(name); err != nil {
		return nil, err
	}
	return p, nil
}

func (p *Product) UpdateName(name string) error {
	name = strings.TrimSpace(name)
	if len(name) < 2 || len(name) > 255 {
		return ErrInvalidProductName
	}
	p.Name = name
	return nil
}

type ProductRepository interface {
	Create(ctx context.Context, p *Product) error
	Update(ctx context.Context, p *Product) error
	Delete(ctx context.Context, id uint) error
	FindByID(ctx context.Context, id uint) (*Product, error)
	FindAll(ctx context.Context, limit, offset int) ([]Product, int64, error)
}
