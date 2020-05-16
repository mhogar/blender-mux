package models

import (
	"regexp"

	"github.com/google/uuid"
)

type User struct {
	ID           uuid.UUID
	Email        string
	PasswordHash []byte
}

func CreateNewUser(email string, passwordHash []byte) *User {
	return &User{
		ID:           uuid.New(),
		Email:        email,
		PasswordHash: passwordHash,
	}
}

func ValidateUserEmail(email string) ValidateError {
	matched, _ := regexp.MatchString(`^[a-zA-Z0-9_\-\.]+@[a-zA-Z0-9_\-\.]+\.[a-zA-Z]{2,}$`, email)
	if !matched {
		return CreateValidateError(UserInvalidEmail, "email is in invalid format")
	}

	return CreateModelValidValidateError()
}

func (u User) Validate() ValidateError {
	if u.ID == uuid.Nil {
		return CreateValidateError(UserInvalidID, "id cannot be nil")
	}

	err := ValidateUserEmail(u.Email)
	if err.Status != ModelValid {
		return err
	}

	if len(u.PasswordHash) == 0 {
		return CreateValidateError(UserInvalidPasswordHash, "password hash cannot be nil")
	}

	return CreateModelValidValidateError()
}
