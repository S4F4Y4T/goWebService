package product

import (
	"time"

	"github.com/S4F4Y4T/goWebService/internal/shared/domain"
)

const ProductCreatedTopic = "product.created"

// ProductCreated is emitted when a new product is created
type ProductCreated struct {
	domain.BaseEvent
	ProductID uint
	Name      string
}

func NewProductCreated(productID uint, name string) ProductCreated {
	return ProductCreated{
		BaseEvent: domain.BaseEvent{Timestamp: time.Now()},
		ProductID: productID,
		Name:      name,
	}
}

func (e ProductCreated) Topic() string {
	return ProductCreatedTopic
}
