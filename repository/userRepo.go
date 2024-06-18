package repository

import (
	"context"
	"go_project/models"
	"log"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
)

type UserRepository struct {
	MongoDB *mongo.Database
}

// Constructor for UserRepository
func NewUserRepository(mongoDB *mongo.Database) *UserRepository {
	return &UserRepository{
		MongoDB: mongoDB,
	}
}

func (c *UserRepository) InsertData(ctx context.Context, user *models.User) error {
	_, err := c.MongoDB.Collection("users").InsertOne(ctx, user)
	if err != nil {
		log.Printf("Error inserting data: %v", err)
		return err
	}
	return nil
}

func (c *UserRepository) FindUserByPhone(ctx context.Context, PhoneNumber string) (*models.User, error) {
    filter := bson.M{"phoneNumber": PhoneNumber}
    var result models.User
    
    err := c.MongoDB.Collection("users").FindOne(ctx, filter).Decode(&result)
    if err != nil {
        if err == mongo.ErrNoDocuments {
            return nil, err
        }
        return nil, err
    }
    
    return &result, nil
}

func (c *UserRepository) FindUserByEmail(ctx context.Context, EmailId string) (*models.User, error) {
	
	filter := bson.M{"emailId": EmailId}
    var result models.User
    
    err := c.MongoDB.Collection("users").FindOne(ctx, filter).Decode(&result)
    if err != nil {
        if err == mongo.ErrNoDocuments {
            return nil, err
        }
        return nil, err
    }
    
    return &result, nil
}