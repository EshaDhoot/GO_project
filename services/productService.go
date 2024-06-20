package services

import (
	"context"
	"go_project/models"
	"go_project/repository"
	"log"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
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

	log.Println("ProductService: product created successfully")
	return nil
}

func (s *ProductService) GetProducts(ctx context.Context) ([]bson.M, error) {
	products, err := s.ProductRepo.GetAllProducts(ctx)
	if err != nil {
		log.Printf("ProductService: error fetching all products: %v", err)
		return nil, err
	}

	log.Println("ProductService: products fetched successfully")
	return products, nil
}

func (s *ProductService) FindProductById(ctx context.Context, ID primitive.ObjectID) (*models.Product, error) {
	product, err := s.ProductRepo.FindById(ctx, ID)
	if err != nil {
		log.Printf("UserService: unable to fetch product")
		return nil, err
	}

	log.Println("UserService: successfully fetched product")
	return product, nil
}


func (s *ProductService) FindProductByIdAndDelete(ctx context.Context, ID primitive.ObjectID) (*models.Product, error) {
	product, err := s.ProductRepo.FindByIdAndDelete(ctx, ID)
	if err != nil {
		log.Printf("UserService: unable to delete product")
		return nil, err
	}

	log.Println("UserService: Deleted product successfully")
	return product, nil
}