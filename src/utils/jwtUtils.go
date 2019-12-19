package utils

import (
	"fmt"
	"time"

	"github.com/dgrijalva/jwt-go"
)

var mySigningKey = []byte(GetConf().JWTKey)

// SignToken is a function to help sign a jwt token for login
func SignToken(audience string) (string, error) {
	claims := jwt.StandardClaims{
		Audience:  audience,
		IssuedAt:  time.Now().Unix(),
		ExpiresAt: time.Now().Unix() + 36000000,
		Issuer:    "MoonCakeDuty",
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	tokenString, err := token.SignedString(mySigningKey)
	return tokenString, err
}

// VerifyToken verify if the token is correct.
func VerifyToken(tokenString string) (*jwt.StandardClaims, error) {
	token, err := jwt.ParseWithClaims(tokenString, &jwt.StandardClaims{}, func(token *jwt.Token) (interface{}, error) {
		// Don't forget to validate the alg is what you expect:
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("Unexpected signing method: %v", token.Header["alg"])
		}
		return mySigningKey, nil
	})
	if err != nil {
		return nil, err
	}
	if claims, ok := token.Claims.(*jwt.StandardClaims); ok && token.Valid {
		return claims, nil
	}
	return nil, fmt.Errorf("Token %s is invalid", tokenString)
}
