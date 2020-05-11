package models

import (
	"regexp"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type User struct {
	ID           primitive.ObjectID
	Email        string
	PasswordHash []byte
}

func CreateNewUser(email string, passwordHash []byte) *User {
	return &User{
		primitive.NewObjectID(),
		email,
		passwordHash,
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
	if u.ID == primitive.NilObjectID {
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
