package models

import (
	"github.com/google/uuid"
)

// User represents the user model.
type User struct {
	ID           uuid.UUID
	Username     string
	PasswordHash []byte
}

// CreateNewUser creates a usel model with new id and the provided fields.
func CreateNewUser(username string, passwordHash []byte) *User {
	return &User{
		ID:           uuid.New(),
		Username:     username,
		PasswordHash: passwordHash,
	}
}

// ValidatePassword validates the password meets the minimum complexity criteria
func ValidatePassword(password string) ValidateError {
	//TODO: implement
	return CreateModelValidValidateError()
}

// Validate validates the the user model has valid fields.
// Returns a ValidateError indicating its result.
func (u User) Validate() ValidateError {
	if u.ID == uuid.Nil {
		return CreateValidateError(ValidateErrorUserInvalidID, "id cannot be nil")
	}

	if u.Username == "" {
		return CreateValidateError(ValidateErrorUserInvalidUsername, "username cannot be empty")
	}

	if len(u.PasswordHash) == 0 {
		return CreateValidateError(ValidateErrorUserInvalidPasswordHash, "password hash cannot be nil")
	}

	return CreateModelValidValidateError()
}
