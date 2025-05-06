package domain

import (
	"github.com/google/uuid"
	"time"
)

type User struct {
	Id        uuid.UUID
	FirstName string
	LastName  *string
	Email     string
	Phone     *string
	CreatedAt time.Time
}

type Address struct {
	Id       uuid.UUID
	UserId   uuid.UUID
	City     string
	Street   string
	Building string
	Apparent *string
	Notes    *string
}

type UserCreate struct {
	Id        uuid.UUID
	FirstName string
	LastName  *string
	Email     string
	Phone     *string
	CreatedAt time.Time
	Password  string
}

type UserUpdate struct {
	FirstName string
	LastName  *string
	Email     string
	Phone     *string
}

type AddressCreate struct {
	Id       uuid.UUID
	UserId   uuid.UUID
	City     string
	Street   string
	Building string
	Apparent *string
	Notes    *string
}

type AddressUpdate struct {
	City     string
	Street   string
	Building string
	Apparent *string
	Notes    *string
}
