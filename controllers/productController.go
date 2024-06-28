package controllers

import (
	"context"
	"go_project/dtos"
	"go_project/models"
	"go_project/services"
	"log"
	"net/http"

	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type ProductController struct {
	ProductService *services.ProductService
}

func NewProductController(productService *services.ProductService) *ProductController {
	return &ProductController{
		ProductService: productService,
	}
}

func (s *ProductController) CreateProduct(ctx *gin.Context) {
	var payload dtos.ProductRequest
	err := ctx.ShouldBindJSON(&payload)
	if err != nil {
		ctx.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	product := &models.Product{
		ID:           primitive.NewObjectID(),
		BuyerName:    payload.BuyerName,
		SellerName:   payload.SellerName,
		UnitPrice:    payload.UnitPrice,
		TotalUnits:   payload.TotalUnits,
		Tenure:       payload.Tenure,
		DiscountRate: payload.DiscountRate,
		Xirr:         payload.Xirr,
	}

	errr := s.ProductService.CreateProduct(context.Background(), product)
	if errr != nil {
		log.Printf("ProductController: error inserting product data: %v", errr)
		ctx.AbortWithStatusJSON(http.StatusBadRequest, gin.H{
			"data":    nil,
			"message": "unable to create a new product",
			"success": false,
			"error":   errr.Error(),
		})
		return
	} else {
		ctx.AbortWithStatusJSON(http.StatusCreated, gin.H{
			"data":    product,
			"message": "successfully created a new product",
			"success": true,
			"error":   nil,
		})
	}
	log.Println("ProductController: product created successfully")
}

func (s *ProductController) FetchProducts(ctx *gin.Context) {
	products, err := s.ProductService.GetProducts(ctx)
	if err != nil {
		ctx.AbortWithStatusJSON(http.StatusBadRequest, gin.H{
			"data":    nil,
			"message": "unable to fetch products",
			"success": false,
			"error":   err,
		})
		return
	} else {
		ctx.AbortWithStatusJSON(http.StatusOK, gin.H{
			"data":    products,
			"message": "successfully fetched products",
			"success": true,
			"error":   nil,
		})
	}
	log.Println("ProductController: products fetched successfully")
}

func (s *ProductController) FetchProductById(ctx *gin.Context) {
	id := ctx.Param("id")

	objID, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		ctx.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"error": "Invalid product ID"})
		return
	}

	product, err := s.ProductService.FindProductById(context.Background(), objID)
	if err != nil {
		ctx.AbortWithStatusJSON(http.StatusBadRequest, gin.H{
			"data":    nil,
			"message": "unable to fetch product",
			"success": false,
			"error":   err,
		})
		return
	}
	ctx.AbortWithStatusJSON(http.StatusOK, gin.H{
		"data":    product,
		"message": "fetched product successfully",
		"success": true,
		"error":   nil,
	})

}

func (s *ProductController) FetchProductByIdAndDelete(ctx *gin.Context) {
	id := ctx.Param("id")

	objID, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		ctx.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"error": "Invalid product ID"})
		return
	}

	product, err := s.ProductService.FindProductByIdAndDelete(context.Background(), objID)
	if err != nil {
		ctx.AbortWithStatusJSON(http.StatusBadRequest, gin.H{
			"data":    nil,
			"message": "unable to delete product",
			"success": false,
			"error":   err,
		})
		return
	}
	ctx.AbortWithStatusJSON(http.StatusOK, gin.H{
		"data":    product,
		"message": "deleted product successfully",
		"success": true,
		"error":   nil,
	})

}
