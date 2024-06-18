package main

import (
	"fmt"
	"log"
	"os"

	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
	"go_project/configs"
	"go_project/routes"
)

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

	routes.UserRouter(router, db)

	address := fmt.Sprintf(":%s", port)
	router.Run(address)
}
