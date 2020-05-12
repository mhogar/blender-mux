package mongoadapter

import (
	"context"
	"time"

	"go.mongodb.org/mongo-driver/mongo"
)

type MongoAdapter struct {
	context    context.Context
	cancelFunc context.CancelFunc
	Client     *mongo.Client
	Migrations *mongo.Collection
	Users      *mongo.Collection
}

func (db *MongoAdapter) CreateStandardTimeoutContext() (context.Context, context.CancelFunc) {
	return context.WithTimeout(db.context, 3*time.Second)
}
