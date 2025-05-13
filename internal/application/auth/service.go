package auth

import (
	"context"
	"github.com/ALexfonSchneider/food-delivery-user-service/internal/domain"
	"github.com/pkg/errors"
	"time"
)

type RegisterUserRequest struct {
	Email     string
	Password  string
	FirstName string
	LastName  *string
	Phone     string
	CreatedAt time.Time
}

func (s *Service) RegisterUser(ctx context.Context, request *RegisterUserRequest) (*domain.User, error) {
	usrStored, err := s.repo.GetUserByEmail(ctx, request.Email)

	if err != nil {
		if !errors.Is(err, domain.RecordNotFoundError) {
			return nil, err
		}
	}
	if usrStored != nil {
		return nil, domain.RecordAlreadyExistsError
	}

	usr, err := domain.NewUser(request.Email, request.Password, request.FirstName, request.Phone, request.LastName)
	if err != nil {
		return nil, err
	}

	passwordHashed, err := s.hasher.Hash(usr.Password)
	if err != nil {
		return nil, err
	}
	usr.HashPassword = passwordHashed

	if err := s.repo.CreateUser(ctx, usr); err != nil {
		return nil, err
	}

	s.log.Info("Successfully registered new user", "email", usr.Email)

	return usr, nil
}

type LoginUserRequest struct {
	Email    string
	Password string
}

type LoginUserResponse struct {
	AccessToken  string
	RefreshToken string
}

func (s *Service) LoginUser(ctx context.Context, request *LoginUserRequest) (*LoginUserResponse, error) {
	usr, err := s.repo.GetUserByEmail(ctx, request.Email)
	if err != nil {
		return nil, err
	}

	passwordHashed, err := s.hasher.Hash(request.Password)
	if err != nil {
		return nil, err
	}
	if s.hasher.Compare(passwordHashed, request.Password) != err {
		return nil, domain.InvalidCredentialsError
	}

	accessToken, err := s.jwt.CreateToken(&UserCredentials{UserID: usr.Id}, AccessToken)
	if err != nil {
		return nil, err
	}

	refreshToken, err := s.jwt.CreateToken(&UserCredentials{UserID: usr.Id}, RefreshToken)
	if err != nil {
		return nil, err
	}

	return &LoginUserResponse{
		AccessToken:  accessToken,
		RefreshToken: refreshToken,
	}, nil
}
