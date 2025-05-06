package domain

import "github.com/pkg/errors"

var (
	RecordNotFoundError      = errors.New("record not found")
	RecordAlreadyExistsError = errors.New("record already exists")
)
