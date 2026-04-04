package repository

import (
	"context"
	"errors"

	"github.com/S4F4Y4T/goWebService/internal/model"
	"gorm.io/gorm"
)

type productRepository struct {
	db *gorm.DB
}

func NewProductRepository(db *gorm.DB) model.ProductRepository {
	return &productRepository{
		db: db,
	}
}

func (r *productRepository) Create(ctx context.Context, req *model.CreateProductRequest) (*model.Product, error) {
	product := model.Product{
		Name: req.Name,
	}

	if err := r.db.WithContext(ctx).Create(&product).Error; err != nil {
		return nil, err
	}

	return &product, nil
}

func (r *productRepository) Update(ctx context.Context, req *model.UpdateProductRequest) (*model.Product, error) {
	var product model.Product
	if err := r.db.WithContext(ctx).First(&product, req.ID).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errors.New("product not found")
		}
		return nil, err
	}

	product.Name = req.Name
	if err := r.db.WithContext(ctx).Save(&product).Error; err != nil {
		return nil, err
	}

	return &product, nil
}

func (r *productRepository) Delete(ctx context.Context, req *model.DeleteProductRequest) error {
	var product model.Product
	if err := r.db.WithContext(ctx).First(&product, req.ID).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return errors.New("product not found")
		}
		return err
	}

	if err := r.db.WithContext(ctx).Delete(&product).Error; err != nil {
		return err
	}
	return nil
}

func (r *productRepository) FindByID(ctx context.Context, req *model.GetProductRequest) (*model.Product, error) {
	var product model.Product
	if err := r.db.WithContext(ctx).First(&product, req.ID).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errors.New("product not found")
		}
		return nil, err
	}

	return &product, nil
}

func (r *productRepository) FindAll(ctx context.Context, req *model.GetProductsRequest) (*model.GetProductsResponse, error) {
	var products []model.Product
	var total int64

	query := r.db.WithContext(ctx).Model(&model.Product{})
	if err := query.Count(&total).Error; err != nil {
		return nil, err
	}

	if req.Limit > 0 {
		query = query.Limit(req.Limit)
	}
	if req.Offset > 0 {
		query = query.Offset(req.Offset)
	}

	if err := query.Find(&products).Error; err != nil {
		return nil, err
	}

	return &model.GetProductsResponse{
		Products: products,
		Total:    total,
		Limit:    req.Limit,
		Offset:   req.Offset,
	}, nil
}
