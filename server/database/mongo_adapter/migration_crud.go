package mongoadapter

import (
	"blendermux/common"
	"blendermux/server/models"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo/options"
)

// CreateMigration creates a new migration model using the given timestamp and inserts it as a document into the migrations collection.
// Returns any errors.
func (db *MongoAdapter) CreateMigration(timestamp string) error {
	//validate timestamp and create the migration to save
	verr := models.ValidateMigrationTimestamp(timestamp)
	if verr.Status != models.ModelValid {
		return common.ChainError("migration timestamp is invalid", verr)
	}
	migration := models.CreateNewMigration(timestamp)

	//insert the migration
	ctx, cancel := db.CreateStandardTimeoutContext()
	_, err := db.Migrations.InsertOne(ctx, migration)
	cancel()

	if err != nil {
		return common.ChainError("error running insert query", err)
	}

	return nil
}

// GetLatestTimestamp returns the latest timestamp of all migrations in the migrations collection.
// If the migrations collection is empty, hasLatest will be false, else it will be true.
// Will also return any errors.
func (db *MongoAdapter) GetLatestTimestamp() (timestamp string, hasLatest bool, err error) {
	//set options to sort by timestamp desc and get max of 1 result
	opts := options.Find()
	opts.SetSort(bson.D{{Key: "timestamp", Value: -1}})
	opts.SetLimit(1)

	//run the query
	ctx, cancel := db.CreateStandardTimeoutContext()
	cursor, err := db.Migrations.Find(ctx, bson.D{}, opts)
	cancel()

	if err != nil {
		return "", false, common.ChainError("error running find query", err)
	}

	//parse the results
	var results []models.Migration

	ctx, cancel = db.CreateStandardTimeoutContext()
	err = cursor.All(ctx, &results)
	cancel()

	if err != nil {
		return "", false, common.ChainError("error iterating over results", err)
	}

	//return the latest timestamp from the results
	if len(results) > 0 {
		return results[0].Timestamp, true, nil
	}

	return "", false, nil
}

// DeleteMigrationByTimestamp deletes up to one migration from the migartions collections with a matching timestamp.
// Returns any errors.
func (db *MongoAdapter) DeleteMigrationByTimestamp(timestamp string) error {
	ctx, cancel := db.CreateStandardTimeoutContext()
	_, err := db.Migrations.DeleteOne(ctx, bson.D{{Key: "timestamp", Value: timestamp}})
	cancel()

	if err != nil {
		return common.ChainError("error running delete query", err)
	}

	return nil
}
