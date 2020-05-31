package models

import (
	"github.com/google/uuid"
)

// Session ValidateError statuses.
const (
	ValidateSessionValid = iota
)

// Session represents the session model.
type Session struct {
	ID     uuid.UUID
	UserID uuid.UUID
}

// CreateNewSession creates a session model with a new id and the provided fields.
func CreateNewSession(userID uuid.UUID) *Session {
	return &Session{
		ID:     uuid.New(),
		UserID: userID,
	}
}

// CreateValidateSessionValid creates a ValidateError with status ValidateSessionValid and nil error.
func CreateValidateSessionValid() ValidateError {
	return ValidateError{ValidateSessionValid, nil}
}

// Validate validates the the session model has valid fields.
// Returns a ValidateError indicating its result.
func (s Session) Validate() ValidateError {
	//TODO: implement
	return CreateValidateSessionValid()
}
