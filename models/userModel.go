package models

import "go.mongodb.org/mongo-driver/bson/primitive"

type User struct {
	ID          primitive.ObjectID `bson:"_id"`
	FirstName   string             `json:"firstName" bson:"firstName" `
	LastName    string             `json:"lastName" bson:"lastName"`
	PhoneNumber string             `json:"phoneNumber" bson:"phoneNumber"`
	EmailId     string             `json:"emailId" bson:"emailId"`
}
