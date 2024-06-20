package models

import "go.mongodb.org/mongo-driver/bson/primitive"

type Product struct {
	ID           primitive.ObjectID `bson:"_id"`
	BuyerName    string             `json:"BuyerName" bson:"BuyerName"`
	SellerName   string             `json:"SellerName" bson:"SellerName"`
	UnitPrice    int                `json:"UnitPrice" bson:"UnitPrice"`
	TotalUnits   int                `json:"TotalUnits" bson:"TotalUnits"`
	Tenure       int             `json:"Tenure" bson:"Tenure"`
	DiscountRate float32            `json:"DiscountRate" bson:"DiscountRate"`
	Xirr         float32            `json:"Xirr" bson:"Xirr"`
}
