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

// UserCRUD is an interface for performing CRUD operations on a user
type UserCRUD interface {
	// CreateUser creates a new user and returns any errors.
	CreateUser(user *models.User) error

	// GetUserBySessionId fetches the user associated with session id.
	// If no users are found, should return nil. Also returns any errors.
	GetUserBySessionId(sID uuid.UUID) (*models.User, error)

	// GetUserByUsername fetches the user with the matching username.
	// If no users are found, should return nil. Also returns any errors.
	GetUserByUsername(username string) (*models.User, error)

	// CreateUser updates the user and returns any errors.
	UpdateUser(user *models.User) error

	// DeleteUser deletes the user with the matching id.
	// Result should be false if no user was deleted, true otherwise. Also returns any errors.
	DeleteUser(id uuid.UUID) (result bool, err error)
}

type SessionCRUD interface {
	CreateSession(session *models.Session) error
	GetSessionByID(ID uuid.UUID) (*models.Session, error)
	DeleteSession(session *models.Session) error
}
