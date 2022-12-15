package services

import (
	"fmt"
	jwt "github.com/dgrijalva/jwt-go"
	"github.com/spf13/viper"
)

type JWTService interface {
	ValidateToken(token string) (*jwt.Token, error)
}

type jwtAuthService struct {
	secretKey string
}

func JWTAuthService() JWTService {
	return &jwtAuthService{secretKey: viper.GetString("JWTSECRETKEY")}
}

func (s *jwtAuthService) ValidateToken(token string) (*jwt.Token, error) {
	return jwt.Parse(token, func(token *jwt.Token) (interface{}, error) {
		if _, isvalid := token.Method.(*jwt.SigningMethodHMAC); !isvalid {
			return nil, fmt.Errorf("Invalid token", token.Header["alg"])
		}
		return []byte(s.secretKey), nil
	})
}
