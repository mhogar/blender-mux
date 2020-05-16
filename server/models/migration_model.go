package models

import (
	"regexp"

	"github.com/google/uuid"
)

type Migration struct {
	ID        uuid.UUID `bson:"id"`
	Timestamp string    `bson:"timestamp"`
}

func CreateNewMigration(timestamp string) *Migration {
	return &Migration{
		ID:        uuid.New(),
		Timestamp: timestamp,
	}
}

func ValidateMigrationTimestamp(timestamp string) ValidateError {
	matched, _ := regexp.MatchString(`^\d{14}$`, timestamp)
	if !matched {
		return CreateValidateError(MigrationInvalidTimestamp, "timestamp is in invalid format")
	}

	return CreateModelValidValidateError()
}

func (m Migration) Validate() ValidateError {
	if m.ID == uuid.Nil {
		return CreateValidateError(MigrationInvalidID, "id cannot be nil")
	}

	err := ValidateMigrationTimestamp(m.Timestamp)
	if err.Status != ModelValid {
		return err
	}

	return CreateModelValidValidateError()
}
