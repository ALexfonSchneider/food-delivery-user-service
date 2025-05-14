package user_test

import (
	"context"
	user "github.com/ALexfonSchneider/food-delivery-user-service/internal/application/user"
	"github.com/ALexfonSchneider/food-delivery-user-service/internal/domain"
	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"
	"log/slog"
	"os"
	"testing"
)

func getMocked(t *testing.T) (*MockRepository, *slog.Logger) {
	ctrl := gomock.NewController(t)

	mockRepository := NewMockRepository(ctrl)

	l := slog.New(slog.NewTextHandler(os.Stdout, nil))

	return mockRepository, l
}

func TestService_GetUserById_Success(t *testing.T) {
	mockRepository, _ := getMocked(t)

	expectedUsr := domain.User{
		Id:    "12345",
		Email: "alex@example.com",
	}

	mockRepository.EXPECT().GetUserById(gomock.Any(), expectedUsr.Id).Return(&expectedUsr, nil)

	service := user.NewService(mockRepository)
	ctx := context.Background()

	usr, err := service.GetUserById(ctx, expectedUsr.Id)

	assert.Nil(t, err)
	assert.Equal(t, expectedUsr, *usr)
}

func TestService_GetUserById_UserNotExists(t *testing.T) {
	mockRepository, _ := getMocked(t)

	expectedUsr := domain.User{
		Id:    "12345",
		Email: "alex@example.com",
	}

	mockRepository.EXPECT().GetUserById(gomock.Any(), expectedUsr.Id).Return(nil, domain.RecordNotFoundError)

	service := user.NewService(mockRepository)
	ctx := context.Background()

	_, err := service.GetUserById(ctx, expectedUsr.Id)

	assert.ErrorIs(t, err, domain.RecordNotFoundError)
}

func TestService_GetUserByEmail_UserExists(t *testing.T) {
	mockRepository, _ := getMocked(t)

	expectedUsr := domain.User{
		Id:    "12345",
		Email: "alex@example.com",
	}

	mockRepository.EXPECT().GetUserByEmail(gomock.Any(), expectedUsr.Email).Return(&expectedUsr, nil)

	service := user.NewService(mockRepository)
	ctx := context.Background()

	usr, err := service.GetUserByEmail(ctx, expectedUsr.Email)

	assert.Nil(t, err)
	assert.Equal(t, expectedUsr, *usr)
}

func TestService_GetUserByEmail_UserNotExists(t *testing.T) {
	mockRepository, _ := getMocked(t)

	expectedUsr := domain.User{
		Id:    "12345",
		Email: "alex@example.com",
	}

	mockRepository.EXPECT().GetUserByEmail(gomock.Any(), expectedUsr.Email).Return(nil, domain.RecordNotFoundError)

	service := user.NewService(mockRepository)
	ctx := context.Background()

	_, err := service.GetUserByEmail(ctx, expectedUsr.Email)

	assert.ErrorIs(t, err, domain.RecordNotFoundError)
}
