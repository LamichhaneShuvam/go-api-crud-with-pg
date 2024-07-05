package utils

import (
	"errors"
	"log"
	"os"
	"strings"
	"time"

	"github.com/golang-jwt/jwt/v4"
)

type Claims struct {
	ID int `json:"id"`
	jwt.RegisteredClaims
}

func GenerateJwt(id int) (string, error) {
	jwtSecret := []byte(os.Getenv("JWT_SECRET"))

	expirationTime := time.Now().Add(10 * time.Hour)

	claims := &Claims{
		ID: id,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(expirationTime),
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	tokenString, err := token.SignedString(jwtSecret)

	if err != nil {
		log.Print("I am here")
		return "", err
	}

	return tokenString, nil
}

func ValidateJWT(tokenString string) (*Claims, error) {
	claims := &Claims{}
	jwtSecret := []byte(os.Getenv("JWT_SECRET"))

	//* Seperated by " " instead of bearer because it seperated bearer word and the token
	normalizedTokenString := strings.Split(tokenString, " ")[1]

	// Parse the token
	token, err := jwt.ParseWithClaims(normalizedTokenString, claims, func(token *jwt.Token) (interface{}, error) {
		return jwtSecret, nil
	})

	if err != nil {
		if err == jwt.ErrSignatureInvalid {
			return nil, errors.New("invalid token signature")
		}
		return nil, err
	}

	if !token.Valid {
		return nil, errors.New("invalid token")
	}

	return claims, nil
}
