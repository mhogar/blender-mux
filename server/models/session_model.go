package models

import (
	"github.com/google/uuid"
)

type Session struct {
	ID     uuid.UUID
	UserID uuid.UUID
}

func CreateNewSession(userID uuid.UUID) *Session {
	return &Session{
		ID:     uuid.New(),
		UserID: userID,
	}
}

func (s Session) Validate() ValidateError {
	return CreateValidateModelValid()
}
