package models

import (
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
