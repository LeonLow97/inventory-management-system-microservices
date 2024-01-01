package inventory

import "log"

type Service interface {
	GetProducts(userID int) (*[]Product, error)
}

type service struct {
	repo Repository
}

func NewService(r Repository) Service {
	return &service{
		repo: r,
	}
}

func (s service) GetProducts(userID int) (*[]Product, error) {
	products, err := s.repo.GetProducts(userID)
	if err != nil {
		log.Println("error getting products", err)
		return nil, err
	}

	return products, nil
}
