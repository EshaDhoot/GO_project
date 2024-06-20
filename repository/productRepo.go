package repository

import (
	"context"
	"go_project/models"
	"log"

	"go.mongodb.org/mongo-driver/mongo"
)

type ProductRepository struct {
	MongoDB *mongo.Database
}

func NewProductRepository(mongoDB *mongo.Database) *ProductRepository {
	return &ProductRepository{
		MongoDB: mongoDB,
	}
}

func (c *ProductRepository) InsertProductData(ctx context.Context, product *models.Product) error {
	_, err := c.MongoDB.Collection("products").InsertOne(ctx, product)
	if err != nil {
		log.Printf("Error inserting product data: %v", err)
		return err
	}
	return nil
}
