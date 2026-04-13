package product

import (
	"context"
	"errors"
	"time"

	"gorm.io/gorm"
)

// productSchema is the Persistence Model (Database Schema)
type productSchema struct {
	ID        uint      `gorm:"primaryKey"`
	Name      string    `gorm:"size:255"`
	CreatedAt time.Time `gorm:"autoCreateTime"`
	UpdatedAt time.Time `gorm:"autoUpdateTime"`
}

func (productSchema) TableName() string {
	return "products"
}

// fromDomain converts a Domain Entity to a Persistence Model
func fromDomain(p *Product) productSchema {
	return productSchema{
		ID:        p.ID,
		Name:      p.Name,
		CreatedAt: p.CreatedAt,
		UpdatedAt: p.UpdatedAt,
	}
}

// toDomain converts a Persistence Model back to a Domain Entity
func (s productSchema) toDomain() *Product {
	return &Product{
		ID:        s.ID,
		Name:      s.Name,
		CreatedAt: s.CreatedAt,
		UpdatedAt: s.UpdatedAt,
	}
}

type productRepository struct {
	db *gorm.DB
}

func NewProductRepository(db *gorm.DB) ProductRepository {
	return &productRepository{
		db: db,
	}
}

func (r *productRepository) Create(ctx context.Context, p *Product) error {
	schema := fromDomain(p)
	if err := r.db.WithContext(ctx).Create(&schema).Error; err != nil {
		return err
	}
	p.ID = schema.ID
	p.CreatedAt = schema.CreatedAt
	p.UpdatedAt = schema.UpdatedAt
	return nil
}

func (r *productRepository) Update(ctx context.Context, p *Product) error {
	schema := fromDomain(p)
	return r.db.WithContext(ctx).Save(&schema).Error
}

func (r *productRepository) Delete(ctx context.Context, id uint) error {
	return r.db.WithContext(ctx).Delete(&productSchema{}, id).Error
}

func (r *productRepository) FindByID(ctx context.Context, id uint) (*Product, error) {
	var s productSchema
	if err := r.db.WithContext(ctx).First(&s, id).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errors.New("product not found")
		}
		return nil, err
	}
	return s.toDomain(), nil
}

func (r *productRepository) FindAll(ctx context.Context, limit, offset int) ([]Product, int64, error) {
	var schemas []productSchema
	var total int64

	query := r.db.WithContext(ctx).Model(&productSchema{})
	if err := query.Count(&total).Error; err != nil {
		return nil, 0, err
	}

	if limit > 0 {
		query = query.Limit(limit)
	}
	if offset > 0 {
		query = query.Offset(offset)
	}

	if err := query.Find(&schemas).Error; err != nil {
		return nil, 0, err
	}

	products := make([]Product, len(schemas))
	for i, s := range schemas {
		products[i] = *s.toDomain()
	}

	return products, total, nil
}
