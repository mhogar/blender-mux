package models

import (
	"github.com/google/uuid"
)

type User struct {
	ID           uuid.UUID
	Username     string
	PasswordHash []byte
}

func CreateNewUser(username string, passwordHash []byte) *User {
	return &User{
		ID:           uuid.New(),
		Username:     username,
		PasswordHash: passwordHash,
	}
}

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
