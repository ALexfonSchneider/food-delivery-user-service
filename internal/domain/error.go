package domain

import "github.com/pkg/errors"

var (
	ErrInvalidEmail          = errors.New("invalid email")
	ErrPasswordTooWeak       = errors.New("password is too weak")
	RecordNotFoundError      = errors.New("record not found")
	RecordAlreadyExistsError = errors.New("record already exists")
	UserNotFoundError        = errors.New("user not found")
	InvalidCredentialsError  = errors.New("invalid credentials")
)
