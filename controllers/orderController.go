package controllers

import (
	"context"
	"fmt"
	"go_project/dtos"
	"go_project/models"
	"go_project/services"
	"log"
	"net/http"
	"os"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
	"github.com/joho/godotenv"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type OrderController struct {
	OrderService   *services.OrderService
	ProductService *services.ProductService
}

func NewOrderController(orderService *services.OrderService) *OrderController {
	return &OrderController{
		OrderService: orderService,
	}
}

func (s *OrderController) CalculatePrice(ctx *gin.Context) {
	var pricePayload dtos.OrderRequest
	err := ctx.ShouldBindJSON(&pricePayload)
	if err != nil {
		ctx.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	productObjID, err := primitive.ObjectIDFromHex(pricePayload.ProductId)
	if err != nil {
		ctx.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"error": "Invalid product ID"})
		return
	}

	product, err := s.ProductService.FindProductById(context.Background(), productObjID)
	if err != nil {
		ctx.AbortWithStatusJSON(http.StatusBadRequest, gin.H{
			"data":    nil,
			"message": "unable to fetch product",
			"success": false,
			"error":   err,
		})
		return
	}

	remainingUnits := product.TotalUnits

	if remainingUnits < pricePayload.NoOfUnits {
		ctx.AbortWithStatusJSON(http.StatusBadRequest, gin.H{
			"data":    nil,
			"message": "These number of units are unavailable",
			"success": false,
		})
		return
	}

	totalPrice := pricePayload.NoOfUnits * product.UnitPrice

	responsePayload := struct {
		TotalPrice   int
		BuyerName    string
		SellerName   string
		UnitPrice    int
		TotalUnits   int
		Tenure       int
		DiscountRate float32
		Xirr         float32
	}{
		TotalPrice:   totalPrice,
		BuyerName:    product.BuyerName,
		SellerName:   product.SellerName,
		UnitPrice:    product.UnitPrice,
		TotalUnits:   product.TotalUnits,
		Tenure:       product.Tenure,
		DiscountRate: product.DiscountRate,
		Xirr:         product.Xirr,
	}

	ctx.AbortWithStatusJSON(http.StatusOK, gin.H{
		"data":    responsePayload,
		"message": "successfully calculated price",
		"success": true,
		"error":   nil,
	})

}

func (s *OrderController) CreateOrder(ctx *gin.Context) {
	var orderPayload dtos.OrderRequest

	err := ctx.ShouldBindJSON(&orderPayload)
	if err != nil {
		ctx.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

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

	userId := refreshToken.Claims.(jwt.MapClaims)["_id"].(string)
	order := &models.Order{
		ID:        primitive.NewObjectID(),
		ProductId: orderPayload.ProductId,
		UserId:    userId,
		NoOfUnits: orderPayload.NoOfUnits,
	}
	errr := s.OrderService.CreateOrder(context.Background(), order)
	if errr != nil {
		log.Printf("OrderController: error inserting order data: %v", errr)
		ctx.AbortWithStatusJSON(http.StatusBadRequest, gin.H{
			"data":    nil,
			"message": "unable to create a new order",
			"success": false,
			"error":   errr.Error(),
		})
		return
	} else {
		ctx.AbortWithStatusJSON(http.StatusCreated, gin.H{
			"data":    order,
			"message": "successfully created a new order",
			"success": true,
			"error":   nil,
		})
	}
	log.Println("OrderController: order created successfully")
}
