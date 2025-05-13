package config

import (
	"fmt"
	"net"
	"net/url"
	"time"
)

type AuthConfig struct {
	SecretKey       string         `koanf:"secretKey"`
	Issuer          *string        `koanf:"Issuer"`
	HashAlgorithm   *string        `koanf:"hashAlgorithm"`
	AccessTokenTTL  *time.Duration `koanf:"accessTokenTTL"`
	RefreshTokenTTL *time.Duration `koanf:"refreshTokenTTL"`
}

func (a *AuthConfig) GetSecretKey() string {
	return a.SecretKey
}

func (a *AuthConfig) GetIssuer() string {
	if a.Issuer != nil {
		return *a.Issuer
	}
	return "user-service"
}

func (a *AuthConfig) GetHashAlgorithm() string {
	if a.HashAlgorithm != nil {
		return *a.HashAlgorithm
	}
	return "HS256"
}

func (a *AuthConfig) GetAccessTokenTTL() time.Duration {
	if a.AccessTokenTTL != nil {
		return *a.AccessTokenTTL
	}
	return time.Hour
}

func (a *AuthConfig) GetRefreshTokenTTL() time.Duration {
	if a.RefreshTokenTTL != nil {
		return *a.RefreshTokenTTL
	}
	return time.Hour * 24
}

func (p *PostgresConfig) ConnectionString() string {
	return fmt.Sprintf(
		"postgres://%s:%s@%s/%s?sslmode=disable",
		p.GetUser(), url.PathEscape(p.GetPassword()),
		net.JoinHostPort(p.GetHost(), p.GetPort()),
		p.GetDatabase(),
	)
}
