package mongoadapter

import (
	"github.com/blendermux/common"
	"github.com/blendermux/server/models"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo/options"
)

func (db *MongoAdapter) CreateMigration(timestamp string) error {
	//validate timestamp and create the migration to save
	verr := models.ValidateMigrationTimestamp(timestamp)
	if verr.Status != models.ModelValid {
		return common.ChainError("error validating migration timestamp", verr)
	}
	migration := models.CreateNewMigration(timestamp)

	//insert the migration
	ctx, cancel := db.CreateStandardTimeoutContext()
	_, err := db.Migrations.InsertOne(ctx, migration)
	cancel()

	if err != nil {
		return common.ChainError("error inserting migration", err)
	}

	return nil
}

func (db *MongoAdapter) GetLatestTimestamp() (string, bool, error) {
	//set options to sort by timestamp desc and get max of 1 result
	opts := options.Find()
	opts.SetSort(bson.D{{"timestamp", -1}})
	opts.SetLimit(1)

	//run the query
	ctx, cancel := db.CreateStandardTimeoutContext()
	cursor, err := db.Migrations.Find(ctx, bson.D{}, opts)
	cancel()

	if err != nil {
		return "", false, err
	}

	//parse the results
	var results []models.Migration

	ctx, cancel = db.CreateStandardTimeoutContext()
	err = cursor.All(ctx, &results)
	cancel()

	if err != nil {
		return "", false, err
	}

	//return the latest timestamp from the results
	if len(results) > 0 {
		return results[0].Timestamp, true, nil
	}

	return "", false, nil
}
