package database

import (
	migrationrunner "blendermux/common/migration_runner"
	"blendermux/server/models"

	"github.com/google/uuid"
)

// Database is an interface that wraps the other database intefaces.
type Database interface {
	DBConnection
	migrationrunner.MigrationCRUD
	UserCRUD
	SessionCRUD
}

// DBConnection is an interface for controlling the connection to the database.
//
// OpenConnection should open the connection to the database. Returns any errors.
//
// CloseConnection should close the connection to the database and cleanup associated resources. Returns any errors.
//
// Ping should ping the database to verify it can still be reached.
// Returns an error if the database can't be reached or if any other errors occur.
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
	GetSessionByID(ID uuid.UUID) (*models.Session, error)
	DeleteSession(session *models.Session) error
}
