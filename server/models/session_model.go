package models

import (
	"github.com/google/uuid"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type Session struct {
	ID     primitive.ObjectID
	Token  uuid.UUID
	UserID primitive.ObjectID
}

func CreateNewSession(token uuid.UUID, userID primitive.ObjectID) *Session {
	return &Session{
		primitive.NewObjectID(),
		token,
		userID,
	}
}

func (s Session) Validate() ValidateError {
	return CreateModelValidValidateError()
}
