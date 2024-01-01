package inventory

import "log"

type Service interface {
	GetProducts(userID int) (*[]Product, error)
	GetProductByID(getProductByIdDTO GetProductByIdDTO) (*Product, error)
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

func (s service) GetProductByID(getProductByIdDTO GetProductByIdDTO) (*Product, error) {
	product, err := s.repo.GetProductByID(getProductByIdDTO)
	if err != nil {
		log.Println("error getting products in service", err)
		return nil, err
	}

	return product, nil
}
