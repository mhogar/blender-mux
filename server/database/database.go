package database

import (
	migrationrunner "github.com/blendermux/common/migration_runner"
	"github.com/blendermux/server/models"
	"github.com/google/uuid"
)

type Database interface {
	migrationrunner.MigrationCRUD
	//UserCRUD
	//SessionCRUD
	Destroy() error
	Ping() error
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
