package main

import (
	"fmt"
	"log"
	"net/http"
	"os"
	"strings"

	"go_project/configs"
	"go_project/routes"

	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
)

func CORSMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		for _, domain := range []string{

			"http://localhost:3000",
		} {
			if strings.Contains(domain, c.Request.Header.Get("Origin")) {
				c.Writer.Header().Set("Access-Control-Allow-Origin", c.Request.Header.Get("Origin"))
			}
		}
		c.Writer.Header().Set("Content-Type", "application/json")
		c.Writer.Header().Set("Access-Control-Max-Age", "86400")
		c.Writer.Header().Set("Access-Control-Allow-Methods", "POST, GET, OPTIONS, PUT, DELETE, UPDATE, PATCH")
		c.Writer.Header().Set("Access-Control-Allow-Headers", "Content-Type, Content-Length, Accept-Encoding, X-CSRF-Token, Authorization, accept, origin, Cache-Control, X-Requested-With, Authentication")
		c.Writer.Header().Set("Access-Control-Allow-Credentials", "true")
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
