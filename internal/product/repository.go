package product

import (
	"context"
	"errors"

	"gorm.io/gorm"
)

type productRepository struct {
	db *gorm.DB
}

func NewProductRepository(db *gorm.DB) ProductRepository {
	return &productRepository{
		db: db,
	}
}

func (r *productRepository) Create(ctx context.Context, p *Product) error {
	return r.db.WithContext(ctx).Create(p).Error
}

func (r *productRepository) Update(ctx context.Context, p *Product) error {
	return r.db.WithContext(ctx).Save(p).Error
}

func (r *productRepository) Delete(ctx context.Context, id uint) error {
	return r.db.WithContext(ctx).Delete(&Product{}, id).Error
}

func (r *productRepository) FindByID(ctx context.Context, id uint) (*Product, error) {
	var p Product
	if err := r.db.WithContext(ctx).First(&p, id).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errors.New("product not found")
		}
		return nil, err
	}
	return &p, nil
}

func (r *productRepository) FindAll(ctx context.Context, limit, offset int) ([]Product, int64, error) {
	var products []Product
	var total int64

	query := r.db.WithContext(ctx).Model(&Product{})
	if err := query.Count(&total).Error; err != nil {
		return nil, 0, err
	}

	if limit > 0 {
		query = query.Limit(limit)
	}
	if offset > 0 {
		query = query.Offset(offset)
	}

	if err := query.Find(&products).Error; err != nil {
		return nil, 0, err
	}

	return products, total, nil
}
