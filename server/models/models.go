package models

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

func GetModelValidValidateError() ValidateError {
	return ValidateError{ModelValid, nil}
}
