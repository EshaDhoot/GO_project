package dtos

type SignInRequest struct {
	EmailId string `json:"emailId" bson:"emailId"`
}
