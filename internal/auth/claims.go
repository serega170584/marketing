package auth

import "github.com/golang-jwt/jwt/v5"

type Claims struct {
	Permissions []string `json:"permissions"`
	jwt.RegisteredClaims
}
