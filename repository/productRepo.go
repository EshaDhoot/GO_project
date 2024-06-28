package repository

import (
	"context"
	"go_project/models"
	"log"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
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

func (c *ProductRepository) GetAllProducts(ctx context.Context) ([]bson.M, error) {
	filter := bson.D{}
	opts := options.Find().SetLimit(2)
	cursor, err := c.MongoDB.Collection("products").Find(ctx, filter, opts)
	if err != nil {
		log.Fatal(err)
		return nil, err
	}
	var products []bson.M
	if err = cursor.All(ctx, &products); err != nil {
		log.Fatal(err)
		return nil, err
	}
	return products, nil
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

func (c *ProductRepository) FindByIdAndUpdate(ctx context.Context, ID primitive.ObjectID, update bson.M) (*models.Product, error) {
	filter := bson.M{"_id": ID}
	updateData := bson.M{"$set": update}

	var updatedProduct models.Product
	err := c.MongoDB.Collection("products").FindOneAndUpdate(ctx, filter, updateData).Decode(&updatedProduct)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			return nil, err
		}
		return nil, err
	}

	return &updatedProduct, nil
}