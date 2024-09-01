package services

import (
	"log"

	"github.com/LeonLow97/internal/core/domain"
	"github.com/LeonLow97/internal/ports"
)

type Service interface {
	GetProducts(userID int) (*[]domain.Product, error)
	GetProductByID(userID, productID int) (*domain.Product, error)
	GetProductByName(userID int, productName string) (*domain.Product, error)
	CreateProduct(req domain.Product, userID int) error
	UpdateProductByID(req domain.Product, userID, productID int) error
	DeleteProductByID(userID, productID int) error
}

type service struct {
	repo ports.Repository
}

func NewService(r ports.Repository) Service {
	return &service{
		repo: r,
	}
}

func (s service) GetProducts(userID int) (*[]domain.Product, error) {
	products, err := s.repo.GetProducts(userID)
	if err != nil {
		log.Println("error getting products", err)
		return nil, err
	}

	return products, nil
}
func (s service) GetProductByID(userID, productID int) (*domain.Product, error) {
	product, err := s.repo.GetProductByID(userID, productID)
	if err != nil {
		log.Println("error getting products in service", err)
		return nil, err
	}

	return product, nil
}

func (s service) GetProductByName(userID int, productName string) (*domain.Product, error) {
	product, err := s.repo.GetProductByName(userID, productName)
	if err != nil {
		log.Println("error getting products in service", err)
		return nil, err
	}

	return product, nil
}

func (s service) CreateProduct(req domain.Product, userID int) error {
	brand, err := s.repo.GetBrandByName(req.BrandName)
	if err != nil {
		log.Println(err)
		return err
	}

	category, err := s.repo.GetCategoryByName(req.CategoryName)
	if err != nil {
		log.Println(err)
		return err
	}

	if err := s.repo.CreateProduct(req, userID, brand.ID, category.ID); err != nil {
		log.Println(err)
		return err
	}

	return nil
}

func (s service) UpdateProductByID(req domain.Product, userID, productID int) error {
	brand, err := s.repo.GetBrandByName(req.BrandName)
	if err != nil {
		log.Println(err)
		return err
	}

	category, err := s.repo.GetCategoryByName(req.CategoryName)
	if err != nil {
		log.Println(err)
		return err
	}

	return s.repo.UpdateProductByID(req, brand.ID, category.ID, userID, productID)
}

func (s service) DeleteProductByID(userID, productID int) error {
	return s.repo.DeleteProductByID(userID, productID)
}
