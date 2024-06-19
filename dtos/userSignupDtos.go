package dtos

type SignUpRequest struct {
	FirstName   string `json:"firstName" bson:"firstName"`
	LastName    string `json:"lastName" bson:"lastName"`
	PhoneNumber string `json:"phoneNumber" bson:"phoneNumber"`
	EmailId     string `json:"emailId" bson:"emailId"`
}
