package main

import (
	"time"

	"github.com/golang-jwt/jwt/v5"
)

func generateJWT(config *Config, username string) (string, error) {
	// create claims
	claims := jwt.MapClaims{
		"username": username,
		"exp":      time.Now().Add(time.Duration(config.JWTExpirationMinutes) * time.Minute).Unix(),
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	return token.SignedString([]byte(config.JWTSecret))
}
