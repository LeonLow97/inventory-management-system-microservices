package inventory

import (
	"log"

	"github.com/LeonLow97/internal/kafkago"
)

type service struct {
	repo              Repository
	segmentioInstance *kafkago.Segmentio
}

func NewService(r Repository, segmentioInstance *kafkago.Segmentio) Service {
	return &service{
		repo:              r,
		segmentioInstance: segmentioInstance,
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

func (s service) GetProductByName(req GetProductDetailsDTO) (*Product, error) {
	product, err := s.repo.GetProductByName(req)
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

func (s service) UpdateProductByID(updateProductDTO UpdateProductDTO) error {
	brand, err := s.repo.GetBrandByName(updateProductDTO.BrandName)
	if err != nil {
		log.Println(err)
		return err
	}

	category, err := s.repo.GetCategoryByName(updateProductDTO.CategoryName)
	if err != nil {
		log.Println(err)
		return err
	}

	updateProductDTO.BrandID = brand.ID
	updateProductDTO.CategoryID = category.ID

	return s.repo.UpdateProductByID(updateProductDTO)
}

func (s service) DeleteProductByID(req DeleteProductDTO) error {
	return s.repo.DeleteProductByID(req)
}
