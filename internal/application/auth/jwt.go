package auth

import (
	"fmt"
	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
	"github.com/pkg/errors"
	"time"
)

type Config struct {
	SecretKey       string
	Issuer          string
	HashAlgorithm   string
	AccessTokenTTL  time.Duration
	RefreshTokenTTL time.Duration
}

type TokenProvider struct {
	Config
}

func NewTokenProvider(config Config) *TokenProvider {
	obj := TokenProvider{config}

	return &obj
}

func (t *TokenProvider) CreateToken(user *UserCredentials, Use TokenType) (string, error) {
	if user == nil {
		return "", errors.New("user is nil")
	}

	claims := Claims{
		UserCredentials: UserCredentials{
			UserID: user.UserID,
		},
		Use: Use,
		RegisteredClaims: jwt.RegisteredClaims{
			Issuer:    t.Issuer,
			IssuedAt:  jwt.NewNumericDate(time.Now()),
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(t.AccessTokenTTL)),
			ID:        uuid.New().String(),
			Subject:   user.UserID,
		},
	}

	token := jwt.NewWithClaims(jwt.GetSigningMethod(t.HashAlgorithm), claims)

	return token.SignedString([]byte(t.SecretKey))
}

func (t *TokenProvider) ValidateToken(tokenString string) (*Claims, error) {
	claims := Claims{}

	if _, err := jwt.ParseWithClaims(tokenString, &claims, func(token *jwt.Token) (interface{}, error) {
		if token.Method != jwt.GetSigningMethod(t.HashAlgorithm) {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}
		return []byte(t.SecretKey), nil
	}); err != nil {
		return nil, err
	}

	return &claims, nil
}
