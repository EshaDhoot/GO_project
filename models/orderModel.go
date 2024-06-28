package models

import "go.mongodb.org/mongo-driver/bson/primitive"

type Order struct {
	ID        primitive.ObjectID `json:"_id" bson:"_id"`
	ProductId string             `json:"ProductId" bson:"ProductId"`
	UserId    string             `json:"UserId" bson:"UserId"`
	NoOfUnits int                `json:"NoOfUnits" bson:"NoOfUnits"`
}
