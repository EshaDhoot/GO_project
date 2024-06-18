package routes

import (
	"go_project/controllers"
	"go_project/repository"
	"go_project/services"

	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/mongo"
)

func UserRouter(incomingRoutes *gin.Engine, mongoDB *mongo.Database) {

	authRepository := repository.NewUserRepository(mongoDB)
	authService := services.NewUserService(authRepository)
	authController := controllers.NewUserController(authService)

	v1 := incomingRoutes.Group("/api/v1")
	{
		v1.POST("/users/signup", authController.SignUp)
		v1.POST("/users/verify-otp", authController.VerifyOtp)

	}

}
