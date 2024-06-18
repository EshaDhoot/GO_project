package dtos

type OtpRequest struct {
	EmailId     string    `json:"emailId" bson:"emailId"`
	OTP         string    `json:"otp" bson:"otp"`
}
