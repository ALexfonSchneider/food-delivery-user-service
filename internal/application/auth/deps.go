package auth

import (
	"context"
	"github.com/ALexfonSchneider/food-delivery-user-service/internal/domain"
	"log/slog"
)

//go:generate mockgen -source=deps.go -destination deps_mock_test.go -package "${GOPACKAGE}_test"

type Repository interface {
	GetUserByEmail(ctx context.Context, email string) (*domain.User, error)
	GetUserById(ctx context.Context, id string) (*domain.User, error)
	CreateUser(ctx context.Context, user *domain.User) error

	// Exec Служит для выполнения нескольких действий в одном. Должен уметь откатывать изменения
	Exec(ctx context.Context, f func(ctx context.Context) error) error
}

type PasswordHasher interface {
	Hash(password string) (string, error)
	Compare(hash, password string) error
}

type JwtProvider interface {
	CreateToken(user *UserCredentials, Use TokenType) (string, error)
	ValidateToken(tokenString string) (*Claims, error)
}

type Service struct {
	log    *slog.Logger
	repo   Repository
	hasher PasswordHasher
	jwt    JwtProvider
}

func NewService(log *slog.Logger, repo Repository, hasher PasswordHasher, jwt JwtProvider) *Service {
	return &Service{log, repo, hasher, jwt}
}
