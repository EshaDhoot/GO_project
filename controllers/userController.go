package controllers

import (
	"context"
	"fmt"
	"os"
	"time"

	// "fmt"
	"go_project/dtos"
	"go_project/models"
	"go_project/services"

	// "go_project/utils"
	"log"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
	"github.com/joho/godotenv"

	// "github.com/golang-jwt/jwt/v5"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type UserController struct {
	UserService *services.UserService
}

func NewUserController(userService *services.UserService) *UserController {
	return &UserController{
		UserService: userService,
	}
}
func (s *UserController) SignUp(ctx *gin.Context) {
	var payload dtos.SignUpRequest
	err := ctx.ShouldBindJSON(&payload)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if payload.FirstName == "" || payload.LastName == "" || payload.PhoneNumber == "" || payload.EmailId == " " {
		ctx.AbortWithStatusJSON(http.StatusBadRequest, gin.H{
			"data":    nil,
			"message": "Fields should not be empty, all the fields are required",
			"success": false,
		})
		return
	}
	existingUserbyPhone, err := s.UserService.FindUserByPhone(context.Background(), payload.PhoneNumber)
	if err == nil {
		ctx.AbortWithStatusJSON(http.StatusBadRequest, gin.H{
			"data":    existingUserbyPhone,
			"message": "User already registered by this phone number",
			"success": false,
		})
		return
	}
	existingUserbyEmail, err := s.UserService.FindUserByEmail(context.Background(), payload.EmailId)
	if err == nil {
		ctx.AbortWithStatusJSON(http.StatusBadRequest, gin.H{
			"data":    existingUserbyEmail,
			"message": "User already registered by this email",
			"success": false,
		})
		return
	}

	user := &models.User{
		ID:          primitive.NewObjectID(),
		FirstName:   payload.FirstName,
		LastName:    payload.LastName,
		PhoneNumber: payload.PhoneNumber,
		EmailId:     payload.EmailId,
	}

	errr := s.UserService.CreateUser(context.Background(), user)
	if errr != nil {
		log.Printf("UserController: error inserting user data: %v", err)
		ctx.AbortWithStatusJSON(http.StatusBadRequest, gin.H{
			"data":    nil,
			"message": "not able to create a new user",
			"success": false,
			"error":   errr.Error(),
		})
	} else {
		ctx.AbortWithStatusJSON(http.StatusCreated, gin.H{
			"data":    user,
			"message": "successfully created a new user",
			"success": true,
			"error":   nil,
		})
	}
	log.Println("UserController: user created successfully")

}

func (s *UserController) SignIn(ctx *gin.Context) {
	var payload dtos.SignInRequest
	err := ctx.ShouldBindJSON(&payload)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	if payload.EmailId == "" {
		ctx.AbortWithStatusJSON(http.StatusBadRequest, gin.H{
			"message": "Please enter email",
			"success": false,
		})
		return
	}
	existingUserbyEmail, err := s.UserService.FindUserByEmail(context.Background(), payload.EmailId)
	if err != nil {
		ctx.AbortWithStatusJSON(http.StatusBadRequest, gin.H{
			"message": "User not found by this email",
			"success": false,
		})
		return
	}
	ctx.AbortWithStatusJSON(http.StatusOK, gin.H{
		"data":    existingUserbyEmail,
		"message": "User found by this email",
		"success": true,
	})

}


func (s *UserController) VerifyOtp(ctx *gin.Context) {
	var otpPayload dtos.OtpRequest
	err := ctx.ShouldBindJSON(&otpPayload)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if otpPayload.OTP == "" {
		ctx.AbortWithStatusJSON(http.StatusBadRequest, gin.H{
			"data":    nil,
			"message": "Enter the OTP",
			"success": false,
		})
		return
	}

	if otpPayload.OTP != "123456" {
		ctx.AbortWithStatusJSON(http.StatusBadRequest, gin.H{
			"message": "Enter correct OTP",
			"success": false,
			"error":   "Incorrect OTP",
		})
		return
	}

	if otpPayload.OTP == "123456" {
		if otpPayload.EmailId == "" {
			ctx.AbortWithStatusJSON(http.StatusBadRequest, gin.H{
				"message": "Please enter email",
				"success": false,
			})
			return
		}
		existingUserbyEmail, err := s.UserService.FindUserByEmail(context.Background(), otpPayload.EmailId)
		
		if err != nil {
			ctx.JSON(http.StatusBadRequest, gin.H{
				"error":   err,
				"message": "user not found",
			})
			return
		}
		token, err := createToken(existingUserbyEmail)
		if err != nil {
			ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		ctx.SetCookie("token", token, 3600, "/", "localhost", false, true)
		ctx.AbortWithStatusJSON(http.StatusOK, gin.H{
			"message": "successfully verified user",
			"success": true,
			"error":   nil,
		})
		return
	}

	// ctx.JSON(http.StatusOK, gin.H{
	// 	"message": "successfully verified user",
	// 	"success": true,
	// 	"error":   nil,
	// })
}

func createToken(user *models.User) (string, error) {

	godotenv.Load(".env")
	SECRET := os.Getenv("JWT_SECRET")
	SECRET_KEY := []byte(SECRET)
	claims := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"firstname": user.FirstName,
		"lastname":  user.LastName,
		"email":     user.EmailId,
		"exp":       time.Now().Add(time.Minute * 30).Unix(),
	})

	tokenString, err := claims.SignedString(SECRET_KEY)
	if err != nil {
		return "", err
	}

	fmt.Printf("Token claims added: %+v\n", claims)
	return tokenString, nil
}
