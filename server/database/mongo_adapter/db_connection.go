package mongoadapter

import (
	"context"

	"github.com/blendermux/common"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"go.mongodb.org/mongo-driver/mongo/readpref"
)

func (db *MongoAdapter) OpenConnection() error {
	db.context, db.cancelFunc = context.WithCancel(context.Background())

	//connect to the db
	ctx, cancel := db.CreateStandardTimeoutContext()
	client, err := mongo.Connect(ctx, options.Client().ApplyURI("mongodb://localhost:27017"))
	cancel()

	if err != nil {
		return common.ChainError("error opening database connection", err)
	}

	//set the adapter fields
	core := client.Database("core")
	db.Client = client
	db.Migrations = core.Collection("migrations")
	db.Users = core.Collection("users")

	return nil
}

func (db *MongoAdapter) CloseConnection() error {
	defer db.cancelFunc() //cancel any remaining requests that may still be running

	ctx, cancel := db.CreateStandardTimeoutContext()
	err := db.Client.Disconnect(ctx)
	cancel()

	if err != nil {
		return common.ChainError("error closing database connection", err)
	}

	return nil
}

func (db *MongoAdapter) Ping() error {
	ctx, cancel := db.CreateStandardTimeoutContext()
	err := db.Client.Ping(ctx, readpref.Primary())
	cancel()

	if err != nil {
		return common.ChainError("error pinging database", err)
	}

	return nil
}
