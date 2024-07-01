package models

import (
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type Order struct {
	ID        primitive.ObjectID `json:"_id" bson:"_id"`
	ProductId primitive.ObjectID `json:"ProductId" bson:"ProductId"`
	UserId    string             `json:"UserId" bson:"UserId"`
	NoOfUnits int                `json:"NoOfUnits" bson:"NoOfUnits"`
	Product   *Product           `json:"Product" bson:"Product"`
	TotalPrice int				`json:"TotalPrice" bson:"TotalPrice"`
}
