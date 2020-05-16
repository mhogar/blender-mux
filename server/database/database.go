package database

import (
	migrationrunner "blendermux/common/migration_runner"
	"blendermux/server/models"
	"github.com/google/uuid"
)

type Database interface {
	DBConnection
	migrationrunner.MigrationCRUD
	UserCRUD
	SessionCRUD
}

type DBConnection interface {
	OpenConnection() error
	CloseConnection() error
	Ping() error
}

type UserCRUD interface {
	CreateUser(user *models.User) error
	GetUserByEmail(email string) (*models.User, error)
}

type SessionCRUD interface {
	CreateSession(session *models.Session) error
	GetSessionByToken(token uuid.UUID) (*models.Session, error)
	DeleteSession(session *models.Session) error
}
