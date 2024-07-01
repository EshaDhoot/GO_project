package repository

import (
	"context"
	"go_project/models"
	"log"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
)

type OrderRepository struct {
	MongoDB *mongo.Database
}

func NewOrderRepository(mongoDB *mongo.Database) *OrderRepository {
	return &OrderRepository {
		MongoDB: mongoDB,
	}
}

func (c *OrderRepository) InsertOrderData(ctx context.Context, order *models.Order) error {
	_, err := c.MongoDB.Collection("orders").InsertOne(ctx, order)
	if err != nil {
	log.Printf("Error inserting order data: %v", err)
		return err
	}
	return nil
}

func (c *OrderRepository) FindById(ctx context.Context, ID string) ([]models.Order, error) {

	filter := bson.M{"UserId": ID}
    cursor, err := c.MongoDB.Collection("orders").Find(ctx, filter)
    if err != nil {
        return nil, err
    }
    defer cursor.Close(ctx)

    var orders []models.Order
    if err := cursor.All(ctx, &orders); err != nil {
        return nil, err
    }
    return orders, nil
}