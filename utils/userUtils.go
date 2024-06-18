package utils

// import (
// 	"context"
// 	"go_project/models"
// 	"os"

// 	"github.com/joho/godotenv"
// 	"go.mongodb.org/mongo-driver/bson"
// 	"go.mongodb.org/mongo-driver/mongo"
// )

// var client *mongo.Client
// var dB_name *mongo.Database
// func FindUserByPhone(PhoneNumber string) (*models.User) {
// 	godotenv.Load(".env")
// 	dB_name := client.Database(os.Getenv("DB_NAME"))
// 	collection := dB_name.Collection("users")

// 	ctx := context.TODO()
// 	filter := bson.M{"phoneNumber": PhoneNumber}

// 	var result models.User

// 	err := collection.FindOne(ctx, filter).Decode(&result)
// 	if err != nil {
// 		if err == mongo.ErrNoDocuments {
// 			return nil
// 		}
// 		return nil
// 	}

// 	return &result
// }

// func FindUserByEmail(EmailId string) (*models.User) {
// 	dB_name := client.Database(os.Getenv("DB_NAME"))
// 	collection := dB_name.Collection("users")

// 	ctx := context.TODO()
// 	filter := bson.M{"emailId": EmailId}

// 	var result models.User

// 	err := collection.FindOne(ctx, filter).Decode(&result)
// 	if err != nil {
// 		if err == mongo.ErrNoDocuments {
// 			return nil
// 		}
// 		return nil
// 	}

// 	return &result
// }
