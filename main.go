package main

import (
	"fmt"
	"log"
	"net/http"
	"os"

	"go_project/configs"
	"go_project/routes"

	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
)

func CORSMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Writer.Header().Set("Access-Control-Allow-Origin", "*") // Change this as needed
		c.Writer.Header().Set("Access-Control-Allow-Methods", "GET, POST, PUT, DELETE, OPTIONS")
		c.Writer.Header().Set("Access-Control-Allow-Headers", "Content-Type, Authorization")

		if c.Request.Method == "OPTIONS" {
			c.AbortWithStatus(http.StatusOK)
			return
		}

		c.Next()
	}
}

func main() {
	err := godotenv.Load(".env")
	if err != nil {
		log.Fatalf("Failed to load .env file: %v", err)
	}

	port := os.Getenv("PORT")
	fmt.Println("PORT => ", port)

	// Connect to MongoDB
	db, err := configs.Connect()
	if err != nil {
		log.Fatalf("Failed to connect to the database: %v", err)
	} else {
		fmt.Println("Connected to the database:", db.Name())
	}

	router := gin.Default()

	// Apply CORS middleware
	router.Use(CORSMiddleware())

	// Initialize routes
	routes.Router(router, db)

	address := fmt.Sprintf(":%s", port)
	router.Run(address)
}
