package services

import (
	"context"
	"go_project/models"
	"go_project/repository"
	"log"
)

type ProductService struct {
	ProductRepo *repository.ProductRepository
}

func NewProductService(productRepo *repository.ProductRepository) *ProductService {
	return &ProductService{
		ProductRepo: productRepo,
	}
}

func (s *ProductService) CreateProduct(ctx context.Context, product *models.Product) error {
	err := s.ProductRepo.InsertProductData(ctx, product)
	if err != nil {
		log.Printf("ProductService: error inserting product data: %v", err)
		return err
	}

	log.Println("ProductService: user product successfully")
	return nil
}
