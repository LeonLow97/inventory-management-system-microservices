package order

type Service interface {
	GetOrders(req GetOrdersDTO) (*[]Order, error)
	GetOrderByID(req GetOrderDTO) (*Order, error)
	CreateOrder(req CreateOrderDTO) error
	ConsumeKafkaUpdateOrderStatus() error
}
