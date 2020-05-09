package models

import (
	"errors"
	"regexp"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type Migration struct {
	ID        primitive.ObjectID `bson:"_id"`
	Timestamp string             `bson:"timestamp"`
}

func CreateNewMigration(timestamp string) *Migration {
	return &Migration{
		primitive.NewObjectID(),
		timestamp,
	}
}

func ValidateMigrationTimestamp(timestamp string) ValidateError {
	matched, _ := regexp.MatchString(`^\d{14}$`, timestamp)
	if !matched {
		return ValidateError{MigrationInvalidTimestamp, errors.New("timestamp is in invalid format")}
	}

	return GetModelValidValidateError()
}

func (m Migration) Validate() ValidateError {
	if m.ID == primitive.NilObjectID {
		return ValidateError{MigrationInvalidID, errors.New("id cannot be nil")}
	}

	err := ValidateMigrationTimestamp(m.Timestamp)
	if err.Status != ModelValid {
		return err
	}

	return GetModelValidValidateError()
}
