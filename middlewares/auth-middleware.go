package middlewares

import (
	"fmt"
	"net/http"
	"os"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
	"github.com/joho/godotenv"
)

func AuthenticateMiddleware(ctx *gin.Context) {

	tokenString, err := ctx.Cookie("accessToken")
	if err != nil {
		fmt.Println("Token missing in cookie")
		ctx.AbortWithStatusJSON(http.StatusNotFound, gin.H{
			"message": "Token missing in cookie",
			"error":   err,
		})
		return
	}

	token, err := verifyToken(tokenString)
	if err != nil {
		fmt.Printf("Token verification failed: %v\\n", err)
		ctx.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{
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
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}
		return SECRET_KEY, nil
	})

	if err != nil {
		return nil, fmt.Errorf("error parsing token: %v", err)
	}

	if !token.Valid {
		return nil, fmt.Errorf("invalid token")
	}

	claims, ok := token.Claims.(jwt.MapClaims)
	if !ok {
		return nil, fmt.Errorf("invalid token claims")
	}

	if exp, ok := claims["exp"].(float64); ok {
		if time.Unix(int64(exp), 0).Before(time.Now()) {
			return nil, fmt.Errorf("access token has expired")
		}
	} else {
		return nil, fmt.Errorf("expiration time not found in token")
	}

	return token, nil
}
