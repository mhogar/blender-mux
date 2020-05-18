package models

import "errors"

// Model is an interface for representing a model.
//
// Validate should check if a model if valid and return an appropriate ValidateError.
type Model interface {
	Validate() ValidateError
}

// ValidateError is a struct for encapsulating the return value of Model's Validate method.
//
// Status is an int that describes the type of error.
//
// error is internal error object.
type ValidateError struct {
	Status int
	error
}

// A ValidateError status.
const (
	ValidateErrorModelValid                = iota
	ValidateErrorMigrationInvalidID        = iota
	ValidateErrorMigrationInvalidTimestamp = iota
	ValidateErrorUserInvalidID             = iota
	ValidateErrorUserInvalidUsername       = iota
	ValidateErrorUserInvalidPasswordHash   = iota
)

// CreateModelValidValidateError creates a ValidateError with status ValidateErrorModelValid and no error.
func CreateModelValidValidateError() ValidateError {
	return ValidateError{ValidateErrorModelValid, nil}
}

// CreateValidateError creates a ValidateError with the provided status and an error with the provided message.
func CreateValidateError(status int, message string) ValidateError {
	return ValidateError{status, errors.New(message)}
}
