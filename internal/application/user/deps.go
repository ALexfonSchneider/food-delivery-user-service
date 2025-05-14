package user

import (
	"context"
	"github.com/ALexfonSchneider/food-delivery-user-service/internal/domain"
)

//go:generate mockgen -source=deps.go -destination deps_mock_test.go -package "${GOPACKAGE}_test"

type Repository interface {
	GetUserByEmail(ctx context.Context, email string) (*domain.User, error)
	GetUserById(ctx context.Context, id string) (*domain.User, error)
	CreateUser(ctx context.Context, user *domain.User) error

	// Exec Служит для выполнения нескольких действий в одном. Должен уметь откатывать изменения
	Exec(ctx context.Context, f func(ctx context.Context) error) error
}

type Service struct {
	repo Repository
}

func NewService(repo Repository) *Service {
	return &Service{repo}
}
