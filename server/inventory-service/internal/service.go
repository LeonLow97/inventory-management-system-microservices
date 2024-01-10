package inventory

type Service interface {
	GetProducts(userID int) (*[]Product, error)
	GetProductByID(getProductByIdDTO GetProductByIdDTO) (*Product, error)
	GetProductByName(req GetProductDetailsDTO) (*Product, error)
	CreateProduct(createProductDTO CreateProductDTO) error
	UpdateProductByID(updateProductDTO UpdateProductDTO) error
	DeleteProductByID(req DeleteProductDTO) error
	ConsumeKafkaUpdateInventoryCount() error
}
