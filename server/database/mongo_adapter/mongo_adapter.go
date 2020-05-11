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
	Client     *mongo.Client
	Migrations *mongo.Collection
	Users      *mongo.Collection
}

func (db *MongoAdapter) Initialize() error {
	//connect to the db
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	client, err := mongo.Connect(ctx, options.Client().ApplyURI("mongodb://localhost:27017"))
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
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	err := db.Client.Disconnect(ctx)
	if err != nil {
		return common.ChainError("error disconnecting from database", err)
	}

	return nil
}

func (db MongoAdapter) Ping() error {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	err := db.Client.Ping(ctx, readpref.Primary())
	if err != nil {
		return common.ChainError("error pinging database", err)
	}

	return nil
}
