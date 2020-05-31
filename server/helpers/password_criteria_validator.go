package helpers

import "errors"

// ValidatePasswordCriteria statuses
const (
	ValidatePasswordCriteriaValid                  = iota
	ValidatePasswordCriteriaTooShort               = iota
	ValidatePasswordCriteriaMissingLowerCaseLetter = iota
	ValidatePasswordCriteriaMissingUpperCaseLetter = iota
	ValidatePasswordCriteriaMissingDigit           = iota
	ValidatePasswordCriteriaMissingSymbol          = iota
)

// ValidatePasswordCriteriaError is a struct for encapsulating a validate password criteria status and internal error.
type ValidatePasswordCriteriaError struct {
	// Status is an int that describes the type of error.
	Status int

	// error is the internal error object.
	error
}

// PasswordCriteriaValidator is an interface for validating a password against criteria.
type PasswordCriteriaValidator interface {
	// ValidatePasswordCriteria validates the password meets the minimum complexity criteria.
	ValidatePasswordCriteria(password string) ValidatePasswordCriteriaError
}

// CreateValidatePasswordCriteriaValid creates a ValidatePasswordCriteriaError with a ValidatePasswordCriteriaValid status and nil err.
func CreateValidatePasswordCriteriaValid() ValidatePasswordCriteriaError {
	return ValidatePasswordCriteriaError{
		Status: ValidatePasswordCriteriaValid,
		error:  nil,
	}
}

// CreateValidatePasswordCriteriaError creates a ValidatePasswordCriteriaError with the provided status and error message.
func CreateValidatePasswordCriteriaError(status int, message string) ValidatePasswordCriteriaError {
	return ValidatePasswordCriteriaError{
		Status: status,
		error:  errors.New(message),
	}
}
