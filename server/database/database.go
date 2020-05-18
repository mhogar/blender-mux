package database

import (
	migrationrunner "blendermux/common/migration_runner"
	"blendermux/server/models"

	"github.com/google/uuid"
)

// Database is an interface that encapsulates the other database intefaces.
type Database interface {
	DBConnection
	migrationrunner.MigrationCRUD
	UserCRUD
	SessionCRUD
}

// DBConnection is an interface for controlling the connection to the database.
type DBConnection interface {
	// OpenConnection should open the connection to the database. Returns any errors.
	OpenConnection() error

	// CloseConnection should close the connection to the database and cleanup associated resources. Returns any errors.
	CloseConnection() error

	// Ping should ping the database to verify it can still be reached.
	// Returns an error if the database can't be reached or if any other errors occur.
	Ping() error
}

type UserCRUD interface {
	CreateUser(user *models.User) error
	GetUserByUsername(username string) (*models.User, error)
}

type SessionCRUD interface {
	CreateSession(session *models.Session) error
	GetSessionByID(ID uuid.UUID) (*models.Session, error)
	DeleteSession(session *models.Session) error
}
