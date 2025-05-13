package auth

import (
	"github.com/golang-jwt/jwt/v5"
)

type UserCredentials struct {
	UserID string
}

type TokenType string

const (
	AccessToken  TokenType = "access_token"
	RefreshToken TokenType = "refresh_token"
)

type Claims struct {
	jwt.RegisteredClaims
	Use TokenType
	UserCredentials
}
