package utils

import (
	"fmt"
	"log"
	"time"

	"github.com/LeonLow97/internal/core/domain"
	"github.com/golang-jwt/jwt/v4"
)

func GenerateJWTToken(user *domain.User, tokenExpiry time.Duration, secretKey string) (string, error) {
	// generate token with claims
	tokenExpireTime := time.Now().Add(tokenExpiry)
	generateToken := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.RegisteredClaims{
		Issuer:    fmt.Sprintf("%d", user.ID),
		ExpiresAt: jwt.NewNumericDate(tokenExpireTime), // 1 hour
	})

	signedToken, err := generateToken.SignedString([]byte(secretKey))
	if err != nil {
		log.Println("error generating jwt token", err)
		return "", err
	}

	return signedToken, nil
}
