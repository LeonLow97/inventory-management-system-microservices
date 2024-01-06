package order

import "log"

type Service interface {
	GetOrders(req GetOrdersDTO) (*[]Order, error)
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
