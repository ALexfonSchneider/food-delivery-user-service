package postgres

import (
	"github.com/ALexfonSchneider/food-delivery-user-service/internal/domain"
)

func UserModelToDomain(user User) *domain.User {
	return &domain.User{
		Id:        user.Id,
		FirstName: user.FirstName,
		LastName:  user.LastName,
		Email:     user.Email,
		Phone:     user.Phone,
		CreatedAt: user.CreatedAt,
	}
}

func AddressModelToDomain(address Address) *domain.Address {
	return &domain.Address{
		Id:       address.Id,
		UserId:   address.UserId,
		City:     address.City,
		Street:   address.Street,
		Building: address.Building,
		Apparent: address.Apparent,
		Notes:    address.Notes,
	}
}
