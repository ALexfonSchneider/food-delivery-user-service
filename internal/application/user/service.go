package user

import (
	"context"
	"github.com/ALexfonSchneider/food-delivery-user-service/internal/domain"
	"github.com/google/uuid"
	"github.com/pkg/errors"
)

func (s *Service) GetUserByEmail(ctx context.Context, email string) (*domain.User, error) {
	user, err := s.repo.GetUserByEmail(ctx, email)
	if err != nil {
		if errors.Is(err, domain.RecordNotFoundError) {
			return nil, nil
		}
		return nil, err
	}
	return user, nil
}

func (s *Service) GetUserById(ctx context.Context, id uuid.UUID) (*domain.User, error) {
	user, err := s.repo.GetUserById(ctx, id)
	if err != nil {
		if errors.Is(err, domain.RecordNotFoundError) {
			return nil, nil
		}
		return nil, err
	}
	return user, nil
}

func (s *Service) CreateUser(ctx context.Context, user *domain.UserCreate) error {
	return s.repo.Exec(ctx, func(ctx context.Context) error {
		userStored, err := s.repo.GetUserById(ctx, user.Id)
		if err != nil {
			if !errors.Is(err, domain.RecordNotFoundError) {
				return err
			}
		}

		if userStored != nil {
			return errors.Wrap(domain.RecordAlreadyExistsError, "user already exists")
		}

		err = s.repo.CreateUser(ctx, user)
		if err != nil {
			return err
		}

		return nil
	})
}
