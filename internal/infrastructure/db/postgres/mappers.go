package postgres

import (
	"github.com/ALexfonSchneider/food-delivery-user-service/internal/domain"
)

func UserModelToDomain(user User) *domain.User {
	return &domain.User{
		Id:           user.Id,
		FirstName:    user.FirstName,
		LastName:     user.LastName,
		Email:        user.Email,
		Phone:        user.Phone,
		CreatedAt:    user.CreatedAt,
		HashPassword: user.HashPassword,
	}
}
