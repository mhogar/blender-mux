package models

import "errors"

type Model interface {
	Validate() ValidateError
}

type ValidateError struct {
	Status int
	error
}

const (
	ModelValid                = iota
	MigrationInvalidID        = iota
	MigrationInvalidTimestamp = iota
	UserInvalidID             = iota
	UserInvalidEmail          = iota
	UserInvalidPasswordHash   = iota
)

func CreateModelValidValidateError() ValidateError {
	return ValidateError{ModelValid, nil}
}

func CreateValidateError(status int, message string) ValidateError {
	return ValidateError{status, errors.New(message)}
}
