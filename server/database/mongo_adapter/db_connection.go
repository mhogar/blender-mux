package mongoadapter

import (
	"context"
	"errors"
	"fmt"
	"log"

	"blendermux/common"
	"blendermux/server/config"

	"github.com/spf13/viper"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"go.mongodb.org/mongo-driver/mongo/readpref"
)

// OpenConnection opens the connection to MongoDB server using the fields from the database config.
// Initializes the adapter's context and cancel function, as well as its client, database, and collection pointers.
// Returns any errors.
func (db *MongoAdapter) OpenConnection() error {
	env := viper.GetString("env")
	mapResult, ok := viper.GetStringMap("database")[env]
	if !ok {
		return errors.New("no database config found for environment " + env)
	}

	dbConfig := mapResult.(config.DatabaseConfig)

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
	dbName, ok := dbConfig.DBs[db.DbKey]
	if !ok {
		return errors.New("could not find database name with key " + db.DbKey)
	}

	//set the adapter fields
	db.Client = client
	db.Database = client.Database(dbName)
	db.Migrations = db.Database.Collection("migrations")

	return nil
}

// CloseConnection closes the connection to the MongoDB server and resets its client, database, and connection pointers.
// The adapter also calls its cancel function to cancel any child requests that may still be running.
// None of the adapter's pointers or context should be used after calling this function.
// Returns any errors.
func (db *MongoAdapter) CloseConnection() error {
	ctx, cancel := db.CreateStandardTimeoutContext()
	err := db.Client.Disconnect(ctx)
	cancel()

	if err != nil {
		return common.ChainError("error closing database connection", err)
	}

	db.cancelFunc() //cancel any remaining requests that may still be running

	//clean up the resources
	db.Client = nil
	db.Database = nil
	db.Migrations = nil

	return nil
}

// Ping pings the MongoDB server to verify it can still be reached.
// Returns an error if it cannot, or if any other errors are encountered.
func (db *MongoAdapter) Ping() error {
	ctx, cancel := db.CreateStandardTimeoutContext()
	err := db.Client.Ping(ctx, readpref.Primary())
	cancel()

	if err != nil {
		return common.ChainError("error pinging database", err)
	}

	return nil
}
