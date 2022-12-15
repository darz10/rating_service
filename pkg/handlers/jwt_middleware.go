package handlers

import (
	"fmt"
	"github.com/dgrijalva/jwt-go"
	"github.com/gin-gonic/gin"
	"net/http"
	"ratingBookingService/pkg/services"
	"strings"
)

func middlewareAuth(c *gin.Context) {
	const (
		BearerSchema = "Bearer"
		AuthHeader   = "Authorization"
	)
	valueAuthHeader := c.GetHeader(AuthHeader)
	if valueAuthHeader == "" {
		c.AbortWithStatus(http.StatusUnauthorized)
	}

	tokenString := strings.Split(valueAuthHeader, " ")
	if len(tokenString) != 2 {
		c.AbortWithStatus(http.StatusUnauthorized)
	}
	isBearer := false
	for _, val := range tokenString {
		if val == BearerSchema {
			isBearer = true
		}
	}
	if !isBearer {
		c.JSON(http.StatusUnauthorized, "Token must be Bearer type")
	}

	token, err := services.JWTAuthService().ValidateToken(tokenString[1])
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": err.Error()})
	}
	if token.Valid {
		claims, err := token.Claims.(jwt.MapClaims)
		if err != true {
			c.AbortWithStatus(http.StatusUnauthorized)
		}
		c.Set("User", claims["user_id"].(float64))
	} else {
		fmt.Println(err)
		c.JSON(http.StatusUnauthorized, "Token is not valid")
	}

}

func AuthJWTMiddleware() gin.HandlerFunc {
	return middlewareAuth
}
