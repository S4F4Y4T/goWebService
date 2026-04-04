package repository

import (
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

func (r *productRepository) Create(req *model.CreateProductRequest) (*model.Product, error) {
	product := model.Product{
		Name: req.Name,
	}

	if err := r.db.Create(&product).Error; err != nil {
		return nil, err
	}

	return &product, nil
}

func (r *productRepository) Update(req *model.UpdateProductRequest) (*model.Product, error) {
	var product model.Product
	if err := r.db.First(&product, req.ID).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errors.New("product not found")
		}
		return nil, err
	}

	product.Name = req.Name
	if err := r.db.Save(&product).Error; err != nil {
		return nil, err
	}

	return &product, nil
}

func (r *productRepository) Delete(req *model.DeleteProductRequest) error {
	var product model.Product
	if err := r.db.First(&product, req.ID).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return errors.New("product not found")
		}
		return err
	}

	if err := r.db.Delete(&product).Error; err != nil {
		return err
	}
	return nil
}

func (r *productRepository) FindByID(req *model.GetProductRequest) (*model.Product, error) {
	var product model.Product
	if err := r.db.First(&product, req.ID).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errors.New("product not found")
		}
		return nil, err
	}

	return &product, nil
}

func (r *productRepository) FindAll(req *model.GetProductsRequest) (*model.GetProductsResponse, error) {
	var products []model.Product
	var total int64

	query := r.db.Model(&model.Product{})
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
