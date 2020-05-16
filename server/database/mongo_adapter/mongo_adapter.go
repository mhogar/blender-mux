package mongoadapter

import (
	"context"
	"time"

	"go.mongodb.org/mongo-driver/mongo"
)

type MongoAdapter struct {
	context    context.Context
	cancelFunc context.CancelFunc
	timeout    int
	DbKey      string
	Client     *mongo.Client
	Database   *mongo.Database
	Migrations *mongo.Collection
	Users      *mongo.Collection
}

func (db *MongoAdapter) CreateStandardTimeoutContext() (context.Context, context.CancelFunc) {
	return context.WithTimeout(db.context, time.Duration(db.timeout)*time.Millisecond)
}
