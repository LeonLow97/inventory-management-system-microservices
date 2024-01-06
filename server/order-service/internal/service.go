package order

import "log"

type Service interface {
	GetOrders(req GetOrdersDTO) (*[]Order, error)
	GetOrderByID(req GetOrderDTO) (*Order, error)
}

type service struct {
	repo Repository
}

func NewService(repo Repository) Service {
	return &service{
		repo: repo,
	}
}

func (s service) GetOrders(req GetOrdersDTO) (*[]Order, error) {
	orders, err := s.repo.GetOrders(req)
	if err != nil {
		log.Println("error getting orders", err)
		return nil, err
	}

	return orders, err
}

func (s service) GetOrderByID(req GetOrderDTO) (*Order, error) {
	order, err := s.repo.GetOrderByID(req)
	if err != nil {
		log.Printf("error getting 1 order by order_id %d with error %v\n", req.OrderID, err)
		return nil, err
	}

	return order, nil
}
