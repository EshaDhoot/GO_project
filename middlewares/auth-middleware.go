package middlewares

import (
	"fmt"
	"net/http"
	"os"
	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
	"github.com/joho/godotenv"
)

func AuthenticateMiddleware(ctx *gin.Context) {

	tokenString, err := ctx.Cookie("token")
	if err != nil {
		fmt.Println("Token missing in cookie")
		// c.Redirect(http.StatusSeeOther, "/login")
		ctx.AbortWithStatusJSON(http.StatusNotFound, gin.H{
			"message": "Token missing in cookie",
			"error":   err,
		})
		return
	}

	token, err := verifyToken(tokenString)
	if err != nil {
		fmt.Printf("Token verification failed: %v\\n", err)
		// c.Redirect(http.StatusSeeOther, "/login")
		ctx.AbortWithStatusJSON(http.StatusBadRequest, gin.H{
			"message": "Token verification failed",
			"error":   err,
		})
		return
	}

	fmt.Printf("Token verified successfully. Claims: %+v\\n", token.Claims)
	ctx.Next()
}

func verifyToken(tokenString string) (*jwt.Token, error) {
	godotenv.Load(".env")
	SECRET := os.Getenv("JWT_SECRET")
	SECRET_KEY := []byte(SECRET)

	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		return SECRET_KEY, nil
	})

	if err != nil {
		return nil, err
	}

	if !token.Valid {
		return nil, fmt.Errorf("invalid token")
	}

	return token, nil
}
