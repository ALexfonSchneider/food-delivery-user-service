package domain

import (
	"github.com/google/uuid"
	"strings"
	"time"
)

type User struct {
	Id           string
	FirstName    string
	LastName     *string
	Email        string
	Phone        string
	Password     string
	HashPassword string
	CreatedAt    time.Time
}

func isValidEmail(email string) bool {
	return strings.Contains(email, "@")
}

func generateID() string {
	return uuid.New().String()
}

func NewUser(email, password, firstName, phone string, lastName *string) (*User, error) {
	if !isValidEmail(email) {
		return nil, ErrInvalidEmail
	}

	if len(password) < 8 {
		return nil, ErrPasswordTooWeak
	}

	return &User{
		Id:        generateID(),
		Email:     email,
		Password:  password,
		FirstName: firstName,
		LastName:  lastName,
		Phone:     phone,
	}, nil
}

//type Address struct {
//	Id       uuid.UUID
//	UserId   uuid.UUID
//	City     string
//	Street   string
//	Building string
//	Apparent *string
//	Notes    *string
//}
