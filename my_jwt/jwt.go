package my_jwt

import (
	"github.com/golang-jwt/jwt/v5"
)

type JwtCustomClaims struct {
	Name  string `json:"name"`
	Admin bool   `json:"admin"`
	jwt.RegisteredClaims
}

func IsAdmin(t *jwt.Token) bool {
	claims := t.Claims.(*JwtCustomClaims)
	return claims.Admin
}
