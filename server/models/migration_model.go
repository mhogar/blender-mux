package models

import (
	"regexp"

	"github.com/google/uuid"
)

// Migration represents the migration model.
type Migration struct {
	ID        uuid.UUID `bson:"id"`
	Timestamp string    `bson:"timestamp"`
}

// CreateNewMigration creates a new migration with a new ID and the given timestamp.
func CreateNewMigration(timestamp string) *Migration {
	return &Migration{
		ID:        uuid.New(),
		Timestamp: timestamp,
	}
}

// ValidateMigrationTimestamp validates the given timestamp is in a valid format.
// Returns a ValidateError indicating its result.
func ValidateMigrationTimestamp(timestamp string) ValidateError {
	matched, _ := regexp.MatchString(`^\d{14}$`, timestamp)
	if !matched {
		return CreateValidateError(ValidateMigrationInvalidTimestamp, "timestamp is in invalid format")
	}

	return CreateValidateModelValid()
}

// Validate validates the migration is a valid migration model.
// Returns a ValidateError indicating its result.
func (m Migration) Validate() ValidateError {
	if m.ID == uuid.Nil {
		return CreateValidateError(ValidateMigrationInvalidID, "id cannot be nil")
	}

	err := ValidateMigrationTimestamp(m.Timestamp)
	if err.Status != ValidateModelValid {
		return err
	}

	return CreateValidateModelValid()
}
