package inventory

import "log"

type Service interface {
	GetProducts(userID int) (*[]Product, error)
	GetProductByID(getProductByIdDTO GetProductByIdDTO) (*Product, error)

	CreateProduct(createProductDTO CreateProductDTO) error
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

func (s service) CreateProduct(createProductDTO CreateProductDTO) error {
	brand, err := s.repo.GetBrandByName(createProductDTO.BrandName)
	if err != nil {
		log.Println(err)
		return err
	}

	category, err := s.repo.GetCategoryByName(createProductDTO.CategoryName)
	if err != nil {
		log.Println(err)
		return err
	}

	if err := s.repo.CreateProduct(createProductDTO, brand.ID, category.ID); err != nil {
		log.Println(err)
		return err
	}

	return nil
}
