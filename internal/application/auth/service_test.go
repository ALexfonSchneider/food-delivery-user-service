package auth_test

import (
	"context"
	"errors"
	"github.com/ALexfonSchneider/food-delivery-user-service/internal/application/auth"
	"github.com/ALexfonSchneider/food-delivery-user-service/internal/domain"
	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"
	"log/slog"
	"os"
	"testing"
)

func getMocked(t *testing.T) (*MockRepository, *MockPasswordHasher, *MockJwtProvider, *slog.Logger) {
	ctrl := gomock.NewController(t)

	mockedRepository := NewMockRepository(ctrl)
	mockedHasher := NewMockPasswordHasher(ctrl)
	mockedJwt := NewMockJwtProvider(ctrl)

	l := slog.New(slog.NewTextHandler(os.Stdout, nil))

	return mockedRepository, mockedHasher, mockedJwt, l
}

var expectedError = errors.New("internal")

func TestService_LoginUser_RegisterSuccess(t *testing.T) {
	mockedRepository, mockedHasher, mockedJwt, l := getMocked(t)

	createdUser := domain.User{
		Email:        "alex@example.com",
		Password:     "123456789",
		HashPassword: "123456789",
	}

	mockedRepository.EXPECT().GetUserByEmail(gomock.Any(), createdUser.Email).Return(nil, nil)
	mockedHasher.EXPECT().Hash(createdUser.Password).Return("123456789", nil)
	mockedRepository.EXPECT().CreateUser(gomock.Any(), gomock.Any()).Return(nil)

	service := auth.NewService(l, mockedRepository, mockedHasher, mockedJwt)

	ctx := context.Background()

	usr, err := service.RegisterUser(ctx, &auth.RegisterUserRequest{
		Email:    createdUser.Email,
		Password: createdUser.Password,
	})

	assert.NoError(t, err)
	assert.NotNil(t, usr)
}

func TestService_LoginUser_RegisterUser_GetUserFail(t *testing.T) {
	mockedRepository, _, _, l := getMocked(t)

	expectedUser := domain.User{
		Email:        "alex@example.com",
		Password:     "123456789",
		HashPassword: "123456789",
	}

	mockedRepository.EXPECT().GetUserByEmail(gomock.Any(), expectedUser.Email).Return(nil, expectedError)

	service := auth.NewService(l, mockedRepository, nil, nil)

	ctx := context.Background()
	_, err := service.RegisterUser(ctx, &auth.RegisterUserRequest{
		Email: expectedUser.Email,
	})

	assert.ErrorIs(t, err, expectedError)
}

func TestService_LoginUser_RegisterFail_UserExists(t *testing.T) {
	mockedRepository, _, _, l := getMocked(t)

	expectedUser := domain.User{
		Email:        "alex@example.com",
		Password:     "123456789",
		HashPassword: "123456789",
	}

	mockedRepository.EXPECT().GetUserByEmail(gomock.Any(), expectedUser.Email).Return(&expectedUser, nil)

	service := auth.NewService(l, mockedRepository, nil, nil)
	ctx := context.Background()
	_, err := service.RegisterUser(ctx, &auth.RegisterUserRequest{
		Email:    expectedUser.Email,
		Password: expectedUser.Password,
	})

	assert.ErrorIs(t, err, domain.RecordAlreadyExistsError)
}

func TestService_RegisterUser_FailCreateUser(t *testing.T) {
	mockedRepository, mockedHasher, mockedJwt, l := getMocked(t)

	createdUser := domain.User{
		Email:        "alex@example.com",
		Password:     "123456789",
		HashPassword: "123456789",
	}

	mockedRepository.EXPECT().GetUserByEmail(gomock.Any(), createdUser.Email).Return(nil, nil)
	mockedHasher.EXPECT().Hash(createdUser.Password).Return("123456789", nil)
	mockedRepository.EXPECT().CreateUser(gomock.Any(), gomock.Any()).Return(expectedError)

	service := auth.NewService(l, mockedRepository, mockedHasher, mockedJwt)
	ctx := context.Background()
	_, err := service.RegisterUser(ctx, &auth.RegisterUserRequest{
		Email:    createdUser.Email,
		Password: createdUser.Password,
	})

	assert.ErrorIs(t, err, expectedError)
}

func TestService_LoginUser_LoginSuccess(t *testing.T) {
	mockedRepository, mockedHasher, mockedJwt, l := getMocked(t)

	existingUser := domain.User{
		Email:        "alex@example.com",
		Password:     "123456789",
		HashPassword: "123456789",
	}

	mockedRepository.EXPECT().GetUserByEmail(gomock.Any(), existingUser.Email).Return(&existingUser, nil)
	mockedHasher.EXPECT().Hash(existingUser.Password).Return("123456789", nil)
	mockedHasher.EXPECT().Compare(existingUser.HashPassword, "123456789").Return(nil)

	accessToken := "accessToken"
	refreshToken := "refreshToken"
	creditals := auth.UserCredentials{}

	mockedJwt.EXPECT().CreateToken(&creditals, auth.AccessToken).Return(accessToken, nil)
	mockedJwt.EXPECT().CreateToken(&creditals, auth.RefreshToken).Return(refreshToken, nil)

	service := auth.NewService(l, mockedRepository, mockedHasher, mockedJwt)
	ctx := context.Background()
	usr, err := service.LoginUser(ctx, &auth.LoginUserRequest{
		Email:    existingUser.Email,
		Password: existingUser.Password,
	})

	assert.NoError(t, err)
	assert.NotNil(t, usr)
	assert.Equal(t, *usr, auth.LoginUserResponse{AccessToken: accessToken, RefreshToken: refreshToken})
}

func TestService_LoginUser_UserDoesNotExists(t *testing.T) {
	mockedRepository, _, _, l := getMocked(t)

	mockedRepository.EXPECT().GetUserByEmail(gomock.Any(), gomock.Any()).Return(nil, domain.RecordNotFoundError)

	service := auth.NewService(l, mockedRepository, nil, nil)
	ctx := context.Background()
	_, err := service.LoginUser(ctx, &auth.LoginUserRequest{
		Email:    "alex@example.com",
		Password: "123456789",
	})

	assert.ErrorIs(t, err, domain.RecordNotFoundError)
}
