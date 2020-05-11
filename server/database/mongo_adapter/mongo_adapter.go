package mongoadapter

import (
	"context"
	"time"

	"github.com/blendermux/common"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"go.mongodb.org/mongo-driver/mongo/readpref"
)

type MongoAdapter struct {
	context    context.Context
	cancelFunc context.CancelFunc
	Client     *mongo.Client
	Migrations *mongo.Collection
	Users      *mongo.Collection
}

func (db MongoAdapter) CreateStandardTimeoutContext() (context.Context, context.CancelFunc) {
	return context.WithTimeout(db.context, 3*time.Second)
}

func (db *MongoAdapter) Initialize() error {
	db.context, db.cancelFunc = context.WithCancel(context.Background())

	//connect to the db
	ctx, cancel := db.CreateStandardTimeoutContext()
	client, err := mongo.Connect(ctx, options.Client().ApplyURI("mongodb://localhost:27017"))
	cancel()

	if err != nil {
		return err
	}

	//set the adapter fields
	core := client.Database("core")
	db.Client = client
	db.Migrations = core.Collection("migrations")
	db.Users = core.Collection("users")

	return nil
}

func (db MongoAdapter) Destroy() error {
	defer db.cancelFunc() //cancel any remaining requests that may still be running

	ctx, cancel := db.CreateStandardTimeoutContext()
	err := db.Client.Disconnect(ctx)
	cancel()

	if err != nil {
		return common.ChainError("error disconnecting from database", err)
	}

	return nil
}

func (db MongoAdapter) Ping() error {
	ctx, cancel := db.CreateStandardTimeoutContext()
	err := db.Client.Ping(ctx, readpref.Primary())
	cancel()

	if err != nil {
		return common.ChainError("error pinging database", err)
	}

	return nil
}
