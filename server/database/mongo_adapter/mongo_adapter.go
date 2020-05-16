package mongoadapter

import (
	"context"
	"time"

	"go.mongodb.org/mongo-driver/mongo"
)

// MongoAdapter is a MongoDB implementation of the Database interface.
//
// DbKey is the key that will be used to resolve the database's name.
//
// Client is a pointer to mongo driver client.
//
// Database is a pointer to the database resolved by the db key.
//
// Migrations is a pointer to the Migrations collection.
type MongoAdapter struct {
	context    context.Context
	cancelFunc context.CancelFunc
	timeout    int
	DbKey      string
	Client     *mongo.Client
	Database   *mongo.Database
	Migrations *mongo.Collection
}

// CreateStandardTimeoutContext creates a context with the timeout loaded from the database config.
// It is a child of the adapter's context and can be canceled by the adapter's cancel function.
// Returns the created context and cancel function.
func (db *MongoAdapter) CreateStandardTimeoutContext() (context.Context, context.CancelFunc) {
	return context.WithTimeout(db.context, time.Duration(db.timeout)*time.Millisecond)
}
