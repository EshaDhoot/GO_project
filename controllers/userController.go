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
		ctx.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"error": err.Error()})
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
		log.Printf("UserController: error inserting user data: %v", errr)
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
		ctx.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	if payload.EmailId == "" {
		ctx.AbortWithStatusJSON(http.StatusBadRequest, gin.H{
			"message": "Please enter email",
			"success": false,
			"error":   err,
		})
		return
	}
	existingUserbyEmail, err := s.UserService.FindUserByEmail(context.Background(), payload.EmailId)
	if err != nil {
		ctx.AbortWithStatusJSON(http.StatusBadRequest, gin.H{
			"message": "User not found by this email",
			"success": false,
			"error":   err,
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
		ctx.AbortWithStatusJSON(http.StatusBadRequest, gin.H{
			"error":   err.Error(),
			"success": false,
		})
	}
	existingUser, err := s.UserService.FindUserByEmail(context.Background(), otpPayload.EmailId)

	if err != nil {
		ctx.AbortWithStatusJSON(http.StatusBadRequest, gin.H{
			"error":   err,
			"message": "user not found",
			"data":    nil,
			"success": false,
		})
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
		// existingUserbyEmail, err := s.UserService.FindUserByEmail(context.Background(), otpPayload.EmailId)

		// if err != nil {
		// 	ctx.JSON(http.StatusBadRequest, gin.H{
		// 		"error":   err,
		// 		"message": "user not found",
		// 	})
		// 	return
		// }
		accessToken, refreshToken, err := createToken(existingUser)
		if err != nil {
			ctx.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		ctx.SetCookie("accessToken", accessToken, 600, "/", "localhost", false, false)

		ctx.SetCookie("refreshToken", refreshToken, 7*24*3600, "/", "localhost", false, false)

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

func createToken(user *models.User) (string, string, error) {

	godotenv.Load(".env")
	SECRET := os.Getenv("JWT_SECRET")
	SECRET_KEY := []byte(SECRET)
	access_token_claims := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"_id":    user.ID,
		"firstname": user.FirstName,
		"lastname":  user.LastName,
		"email":     user.EmailId,
		"exp":       time.Now().Add(time.Minute * 10).Unix(),
	})

	refresh_token_claims := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"_id":    user.ID,
		"firstname": user.FirstName,
		"lastname":  user.LastName,
		"email":     user.EmailId,
		"exp":       time.Now().Add(30 * 24 * time.Hour).Unix(),
	})

	accessTokenString, err := access_token_claims.SignedString(SECRET_KEY)
	if err != nil {
		return "", "", err
	}

	refreshTokenString, err := refresh_token_claims.SignedString(SECRET_KEY)
	if err != nil {
		return "", "", err
	}

	fmt.Printf("AccessToken claims added: %+v\n", access_token_claims)
	fmt.Printf("RefreshToken claims added: %+v\n", refresh_token_claims)
	return accessTokenString, refreshTokenString, nil
}

func (s *UserController) RefreshToken(ctx *gin.Context) {
	refreshTokenString, err := ctx.Cookie("refreshToken")
	if err != nil {
		ctx.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"error": err, "message": "Refresh token is missing"})
		return
	}
	godotenv.Load(".env")
	SECRET := os.Getenv("JWT_SECRET")
	SECRET_KEY := []byte(SECRET)

	refreshToken, err := jwt.Parse(refreshTokenString, func(refreshToken *jwt.Token) (interface{}, error) {
		if _, ok := refreshToken.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", refreshToken.Header["alg"])
		}
		return SECRET_KEY, nil
	})

	if err != nil {
		ctx.AbortWithStatusJSON(http.StatusBadRequest, gin.H{
			"error":   err,
			"message": "Not able to parse refresh token",
			"success": false,
		})
		return

	}

	if !refreshToken.Valid {
		ctx.AbortWithStatusJSON(http.StatusBadRequest, gin.H{
			"error":   err,
			"message": "Invalid refresh token",
			"success": false,
		})
		return
	}

	claims, ok := refreshToken.Claims.(jwt.MapClaims)
	if !ok {
		ctx.AbortWithStatusJSON(http.StatusBadRequest, gin.H{
			"error":   err,
			"message": "invalid token claims",
			"success": false,
		})
		return
	}

	if exp, ok := claims["exp"].(float64); ok {
		if time.Unix(int64(exp), 0).Before(time.Now()) {
			ctx.AbortWithStatusJSON(http.StatusBadRequest, gin.H{
				"error":   err,
				"message": "refresh token has expired",
				"success": false,
			})
			return
		}
	} else {
		ctx.AbortWithStatusJSON(http.StatusBadRequest, gin.H{
			"error":   err,
			"message": "expiration time not found in token",
			"success": false,
		})
		return
	}

	userEmail := refreshToken.Claims.(jwt.MapClaims)["email"].(string)
	existingUserbyEmail, err := s.UserService.FindUserByEmail(context.Background(), userEmail)
	if err != nil {
		ctx.AbortWithStatusJSON(http.StatusBadRequest, gin.H{
			"error":   err,
			"message": "user not found",
			"success": false,
		})
		return
	}
	newAccessToken, newRefreshToken, err := createToken(existingUserbyEmail)
	if err != nil {
		ctx.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"error": err, "message": "token not found"})
		return
	}

	ctx.SetCookie("accessToken", newAccessToken, 3600, "/", "localhost", false, true)

	ctx.SetCookie("refreshToken", newRefreshToken, 30*24*3600, "/", "localhost", false, true)

	ctx.AbortWithStatusJSON(http.StatusOK, gin.H{
		"message": "successfully verified user",
		"success": true,
		"error":   nil,
	})

}
