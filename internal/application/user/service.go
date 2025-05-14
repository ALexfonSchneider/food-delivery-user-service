package user

import (
	"context"
	"github.com/ALexfonSchneider/food-delivery-user-service/internal/domain"
)

func (s *Service) GetUserByEmail(ctx context.Context, email string) (*domain.User, error) {
	user, err := s.repo.GetUserByEmail(ctx, email)
	if err != nil {
		return nil, err
	}
	return user, nil
}

func (s *Service) GetUserById(ctx context.Context, id string) (*domain.User, error) {
	user, err := s.repo.GetUserById(ctx, id)
	if err != nil {
		return nil, err
	}
	return user, nil
}
