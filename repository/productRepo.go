package repository

import (
	"context"
	"fmt"
	"go_project/models"
	"log"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
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

func (c *ProductRepository) GetAllProducts(ctx context.Context) error {
	cursor, err := c.MongoDB.Collection("products").Find(ctx, bson.M{})
	if err != nil {
		log.Fatal(err)
		return err
	}
	var products []bson.M
	if err = cursor.All(ctx, &products); err != nil {
		log.Fatal(err)
		return err
	}
	fmt.Println(products)
	return nil
}

func (c *ProductRepository) FindById(ctx context.Context, ID primitive.ObjectID) (*models.Product, error) {

	filter := bson.M{"_id": ID}
	var result models.Product
	
	err := c.MongoDB.Collection("products").FindOne(ctx, filter).Decode(&result)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			return nil, err
		}
		return nil, err
	}

	return &result, nil
}

func (c *ProductRepository) FindByIdAndDelete(ctx context.Context, ID primitive.ObjectID) (*models.Product, error) {

	filter := bson.M{"_id": ID}
	var result models.Product
	
	err := c.MongoDB.Collection("products").FindOneAndDelete(ctx, filter).Decode(&result)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			return nil, err
		}
		return nil, err
	}

	return &result, nil
}
