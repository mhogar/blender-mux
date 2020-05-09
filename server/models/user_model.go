package models

import (
	"errors"
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

func ValidateUserEmail(email string) bool {
	matched, _ := regexp.MatchString(`^[a-zA-Z0-9_\-\.]+@[a-zA-Z0-9_\-\.]+\.[a-zA-Z]{2,}$`, email)
	return matched
}

func (u User) Validate() ValidateError {
	if u.ID == primitive.NilObjectID {
		return ValidateError{UserInvalidID, errors.New("id cannot be nil")}
	}

	if !ValidateUserEmail(u.Email) {
		return ValidateError{UserInvalidEmail, errors.New("email is in invalid format")}
	}

	if len(u.PasswordHash) == 0 {
		return ValidateError{UserInvalidPasswordHash, errors.New("password hash cannot be nil")}
	}

	return GetModelValidValidateError()
}
