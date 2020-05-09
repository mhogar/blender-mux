package mongoadapter

import (
	"context"
	"errors"
	"time"

	"github.com/blendermux/server/models"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo/options"
)

func (db MongoAdapter) CreateMigration(timestamp string) error {
	//validate timestamp and create the migration to save
	verr := models.ValidateMigrationTimestamp(timestamp)
	if verr.Status != models.ModelValid {
		return errors.New("error validating migration timestamp" + verr.Error())
	}
	migration := models.CreateNewMigration(timestamp)

	//insert the migration
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	_, err := db.Migrations.InsertOne(ctx, migration)
	if err != nil {
		return errors.New("error inserting migration: " + err.Error())
	}

	return nil
}

func (db MongoAdapter) GetLatestTimestamp() (string, error) {
	//set options to sort by timestamp desc and get max of 1 result
	opts := options.Find()
	opts.SetSort(bson.D{{"timestamp", -1}})
	opts.SetLimit(1)

	//run the query
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	cursor, err := db.Migrations.Find(ctx, bson.D{}, opts)
	if err != nil {
		return "", err
	}

	//parse the results
	var results []models.Migration
	err = cursor.All(context.TODO(), &results)
	if err != nil {
		return "", err
	}

	//return the latest timestamp from the results
	if len(results) > 0 {
		return results[0].Timestamp, nil
	}

	return "", nil
}