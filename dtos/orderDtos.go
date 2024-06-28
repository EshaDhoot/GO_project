package dtos

type OrderRequest struct {
	ProductId	string	  `json:"ProductId" bson:"ProductId"`
	NoOfUnits   int     `json:"TotalUnits" bson:"TotalUnits"`	
}