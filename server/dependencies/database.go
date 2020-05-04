package dependencies

import (
	"github.com/blendermux/server/models"
	"github.com/google/uuid"
)

type Database interface {
	UserCRUD
	SessionCRUD
}

type UserCRUD interface {
	CreateUser(user *models.User) error
	GetUserByEmail(email string) *models.User
}

type SessionCRUD interface {
	CreateSession(session *models.Session)
	GetSessionByID(id uuid.UUID) *models.Session
	DeleteSession(session *models.Session)
}
