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
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type OrderController struct {
	OrderService   *services.OrderService
	ProductService *services.ProductService
	
}

func NewOrderController(orderService *services.OrderService,  productService *services.ProductService) *OrderController {
	return &OrderController{
		OrderService: orderService,
		ProductService: productService,
		
	}
}

func (s *OrderController) CalculatePrice(ctx *gin.Context) {
    log.Println("CalculatePrice called")

    if s.ProductService == nil {
        log.Println("ProductService is nil")
        ctx.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{
            "error": "ProductService is not initialized",
        })
        return
    }

    var pricePayload dtos.OrderRequest
    err := ctx.ShouldBindJSON(&pricePayload)
    if err != nil {
        log.Printf("Error binding JSON: %v\n", err)
        ctx.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"error": err.Error()})
        return
    }

    log.Printf("Received OrderRequest: %+v\n", pricePayload)

    productObjID, err := primitive.ObjectIDFromHex(pricePayload.ProductId)
    if err != nil {
        log.Printf("Invalid product ID: %v\n", err)
        ctx.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"error": "Invalid product ID"})
        return
    }

    log.Printf("Fetching product with ID: %v\n", productObjID)
    product, err := s.ProductService.FindProductById(context.Background(), productObjID)
    if err != nil {
        log.Printf("Error finding product: %v\n", err)
        ctx.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{
            "data":    nil,
            "message": "Unable to fetch product",
            "success": false,
            "error":   err.Error(),
        })
        return
    }

    if product == nil {
        log.Println("Product not found")
        ctx.AbortWithStatusJSON(http.StatusNotFound, gin.H{
            "data":    nil,
            "message": "Product not found",
            "success": false,
        })
        return
    }

    log.Printf("Product found: %+v\n", product)

    remainingUnits := product.TotalUnits

    if remainingUnits < pricePayload.NoOfUnits {
        log.Println("Insufficient units")
        ctx.AbortWithStatusJSON(http.StatusBadRequest, gin.H{
            "data":    nil,
            "message": "These number of units are unavailable",
            "success": false,
        })
        return
    }

    log.Printf("NoOfUnits: %d, UnitPrice: %d\n", pricePayload.NoOfUnits, product.UnitPrice)
    totalPrice := pricePayload.NoOfUnits * product.UnitPrice
    log.Printf("TotalPrice: %d\n", totalPrice)

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

    log.Println("Price calculated successfully")
    ctx.JSON(http.StatusOK, gin.H{
        "data":    responsePayload,
        "message": "Successfully calculated price",
        "success": true,
        "error":   nil,
    })
}


func (s *OrderController) CreateOrder(ctx *gin.Context) {
	var orderPayload dtos.OrderRequest

	if err := ctx.ShouldBindJSON(&orderPayload); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request payload", "details": err.Error()})
		return
	}

	accessTokenString, err := ctx.Cookie("accessToken")
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "Access token is missing", "details": err.Error()})
		return
	}

	godotenv.Load(".env")
	SECRET := os.Getenv("JWT_SECRET")
	SECRET_KEY := []byte(SECRET)

	accessToken, err := jwt.Parse(accessTokenString, func(accessToken *jwt.Token) (interface{}, error) {
		if _, ok := accessToken.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", accessToken.Header["alg"])
		}
		return SECRET_KEY, nil
	})
	if err != nil || !accessToken.Valid {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "Invalid access token", "details": err.Error()})
		return
	}

	claims, ok := accessToken.Claims.(jwt.MapClaims)
	if !ok {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "Invalid token claims"})
		return
	}

	if exp, ok := claims["exp"].(float64); ok {
		if time.Unix(int64(exp), 0).Before(time.Now()) {
			ctx.JSON(http.StatusBadRequest, gin.H{"error": "Access token has expired"})
			return
		}
	} else {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "Expiration time not found in token"})
		return
	}

	userId, ok := claims["_id"].(string)
	if !ok {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "User ID not found in token"})
		return
	}

	productObjID, err := primitive.ObjectIDFromHex(orderPayload.ProductId)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "Invalid product ID"})
		return
	}

	product, err := s.ProductService.FindProductById(context.Background(), productObjID)
	if err != nil {
		log.Printf("OrderController: Unable to fetch product: %v", err)
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "Unable to fetch product", "details": err.Error()})
		return
	}

	
	remainingUnits := product.TotalUnits - orderPayload.NoOfUnits
	update := bson.M{"$set": bson.M{"TotalUnits": remainingUnits}}
	updatedProduct, err := s.ProductService.FindProductByIdAndUpdate(context.Background(), productObjID, update)
	if err != nil {
		log.Printf("OrderController: Error updating product data: %v", err)
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "Unable to update product", "details": err.Error()})
		return
	}
	log.Printf("OrderController: Updated product data: %v", updatedProduct)

	
	order := &models.Order{
		ID:        primitive.NewObjectID(),
		ProductId: productObjID,
		UserId:    userId,
		NoOfUnits: orderPayload.NoOfUnits,
		Product:   product,
	}

	if err := s.OrderService.CreateOrder(context.Background(), order); err != nil {
		log.Printf("OrderController: Error creating order: %v", err)
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "Unable to create order", "details": err.Error()})
		return
	}

	ctx.JSON(http.StatusCreated, gin.H{"message": "Successfully created a new order", "data": order})
	log.Println("OrderController: Order created successfully")
}

func (s *OrderController) FetchOrderByUserId(ctx *gin.Context) {
	userid := ctx.Param("userid")

	// objID, err := primitive.ObjectIDFromHex(id)
	// if err != nil {
	// 	ctx.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"error": "Invalid User ID"})
	// 	return
	// }

	orders, err := s.OrderService.FindOrderByUserId(context.Background(), userid)
	if err != nil {
		ctx.AbortWithStatusJSON(http.StatusBadRequest, gin.H{
			"data":    nil,
			"message": "unable to fetch order",
			"success": false,
			"error":   err,
		})
		return
	}
	ctx.AbortWithStatusJSON(http.StatusOK, gin.H{
		"data":    orders,
		"message": "fetched order successfully",
		"success": true,
		"error":   nil,
	})

}
