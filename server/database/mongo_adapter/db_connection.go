package mongoadapter

import (
	"context"
	"errors"
	"fmt"
	"log"

	"github.com/blendermux/common"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"go.mongodb.org/mongo-driver/mongo/readpref"
)

func (db *MongoAdapter) OpenConnection() error {
	dbConfig, err := db.GetDatabaseConfig()
	if err != nil {
		return common.ChainError("error getting database config", err)
	}

	db.context, db.cancelFunc = context.WithCancel(context.Background())
	db.timeout = dbConfig.Timeout

	//connect to the db
	connectionStr := fmt.Sprintf("mongodb://%s:%s", dbConfig.URL, dbConfig.Port)
	log.Println("Connecting to database with connection string", connectionStr)

	ctx, cancel := db.CreateStandardTimeoutContext()
	client, err := mongo.Connect(ctx, options.Client().ApplyURI(connectionStr))
	cancel()

	if err != nil {
		return common.ChainError("error opening database connection", err)
	}

	//get the database name
	dbName, ok := dbConfig.Dbs[db.DbKey]
	if !ok {
		return errors.New("could not find database name with key " + db.DbKey)
	}

	//set the adapter fields
	core := client.Database(dbName)
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
