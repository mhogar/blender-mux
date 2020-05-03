package models

import (
	"errors"
	"regexp"

	"github.com/google/uuid"
)

type User struct {
	ID           uuid.UUID
	Email        string
	PasswordHash []byte
}

const (
	UserInvalidID           = iota
	UserInvalidEmail        = iota
	UserInvalidPasswordHash = iota
)

func (u *User) Validate() *ValidateError {
	if u.ID == uuid.Nil {
		return &ValidateError{UserInvalidID, errors.New("id cannot be nil")}
	}

	matched, _ := regexp.MatchString(`^[a-zA-Z0-9_\-\.]+@[a-zA-Z0-9_\-\.]+\.[a-zA-Z]{2,}$`, u.Email)
	if !matched {
		return &ValidateError{UserInvalidEmail, errors.New("email is in invalid format")}
	}

	if len(u.PasswordHash) == 0 {
		return &ValidateError{UserInvalidPasswordHash, errors.New("password hash cannot be null")}
	}

	return nil
}
