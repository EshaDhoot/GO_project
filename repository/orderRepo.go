package repository

import (
	"context"
	"go_project/models"
	"log"

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