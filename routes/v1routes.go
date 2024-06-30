package routes

import (
	"go_project/controllers"
	"go_project/repository"
	"go_project/services"
"go_project/middlewares"
	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/mongo"
)

func Router(incomingRoutes *gin.Engine, mongoDB *mongo.Database) {

	authRepository := repository.NewUserRepository(mongoDB)
	authService := services.NewUserService(authRepository)
	authController := controllers.NewUserController(authService)

	productRepository := repository.NewProductRepository(mongoDB)
	productService := services.NewProductService(productRepository)
	productController := controllers.NewProductController(productService)

	orderRepository := repository.NewOrderRepository(mongoDB)
	orderService := services.NewOrderService(orderRepository)
	orderController := controllers.NewOrderController(orderService, productService)


	v1 := incomingRoutes.Group("/api/v1")
	{
		v1.POST("/users/signup", authController.SignUp)
		v1.POST("/users/signin", authController.SignIn)
		v1.POST("/users/verify-otp", authController.VerifyOtp)
		v1.POST("/users/refresh-token", authController.RefreshToken)

		v1.POST("/products", productController.CreateProduct)
		v1.GET("/products", middlewares.AuthenticateMiddleware, productController.FetchProducts)
		v1.GET("/product/:id", middlewares.AuthenticateMiddleware, productController.FetchProductById)
		v1.DELETE("/product/:id", productController.FetchProductByIdAndDelete)

		v1.POST("/order/details", middlewares.AuthenticateMiddleware, orderController.CalculatePrice)
		v1.POST("/order", middlewares.AuthenticateMiddleware, orderController.CreateOrder)

	}

}


