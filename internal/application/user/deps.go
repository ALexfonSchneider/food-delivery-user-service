package user

import (
	"context"
	"github.com/ALexfonSchneider/food-delivery-user-service/internal/domain"
	"github.com/google/uuid"
)

type Repository interface {
	GetUserByEmail(ctx context.Context, email string) (*domain.User, error)
	GetUserById(ctx context.Context, id uuid.UUID) (*domain.User, error)
	CreateUser(ctx context.Context, user *domain.UserCreate) error

	// Exec Служит для выполнения нескольких действий в одном. Должен уметь откатывать изменения
	Exec(ctx context.Context, f func(ctx context.Context) error) error
}

type Service struct {
	repo Repository
}

func NewService(repo Repository) *Service {
	return &Service{repo}
}
