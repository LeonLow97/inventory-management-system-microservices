package ports

import (
	"github.com/LeonLow97/internal/core/domain"
	"github.com/segmentio/kafka-go"
)

type Repository interface {
	// database calls
	GetProducts(userID int) (*[]domain.Product, error)
	GetProductByID(userID, productID int) (*domain.Product, error)
	GetProductByName(userID int, productName string) (*domain.Product, error)
	CreateProduct(req domain.Product, userID, brandID, categoryID int) error
	UpdateProductByID(req domain.Product, brandID, categoryID, userID, productID int) error
	UpdateProductQuantityByID(quantity, userID, productID int) error
	DeleteProductByID(userID, productID int) error
	GetBrandByName(brandName string) (*domain.Brand, error)
	GetCategoryByName(categoryName string) (*domain.Category, error)

	// events
	ConsumeOrderMessage(brokerAddress, topic string, messageChan chan interface{}, errorChan chan error)
	ProduceOrderMessage(brokerAddress, topic string, message []kafka.Message) error
}
