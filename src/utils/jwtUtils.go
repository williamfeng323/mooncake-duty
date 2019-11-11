package utils

import (
	"time"

	"github.com/dgrijalva/jwt-go"
)

var mySigningKey = []byte(GetConf().JWTKey)

// SignToken is a function to help sign a jwt token for login
func SignToken(audience string) (string, error) {
	claims := &jwt.StandardClaims{
		Audience:  audience,
		IssuedAt:  time.Now().Unix(),
		ExpiresAt: 15000,
		Issuer:    "MoonCakeDuty",
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	tokenString, err := token.SignedString(mySigningKey)
	return tokenString, err
}
