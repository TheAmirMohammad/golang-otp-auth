package jwtutil

import (
	"time"
	
	"github.com/golang-jwt/jwt/v5"
)

func Generate(secret, userID string, ttl time.Duration) (string, error) {
	claims := jwt.MapClaims{"sub": userID, "exp": time.Now().Add(ttl).Unix(), "iat": time.Now().Unix()}
	t := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return t.SignedString([]byte(secret))
}
