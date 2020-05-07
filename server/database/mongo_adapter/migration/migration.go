package main

import (
	"context"
	"fmt"
	"log"
	"time"

	"go.mongodb.org/mongo-driver/mongo/readpref"

	"github.com/blendermux/server/models"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type Migration interface {
	GetTimestamp() string
	Up(db *mongo.Database)
	Down(db *mongo.Database)
}

func main() {
	//connect to the db
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	client, err := mongo.Connect(ctx, options.Client().ApplyURI("mongodb://localhost:27017"))

	defer cancel()
	if err != nil {
		log.Fatal(err)
	}

	//verify the database has been conected to
	ctx, cancel = context.WithTimeout(context.Background(), 3*time.Second)
	err = client.Ping(ctx, readpref.Primary())

	defer cancel()
	if err != nil {
		log.Fatal(err)
	}

	db := client.Database("core")

	fmt.Println("connected to db")

	//init the migrations
	migrations := []Migration{
		M20200507205301{},
	}

	newestTimestamp := getNewestTimestamp(db)
	fmt.Println(newestTimestamp)

	for _, migration := range migrations {
		//TODO: log when migrations are run and which ones are skipped

		//if migration is newer than most recent in db
		if migration.GetTimestamp() > newestTimestamp {
			migration.Up(db)
			saveTimestamp(db, migration.GetTimestamp())
		}
	}
}

func getNewestTimestamp(db *mongo.Database) string {
	//set options to sort by timestamp desc and get max of 1 result
	opts := options.Find()
	opts.SetSort(bson.D{{"timestamp", -1}})
	opts.SetLimit(1)

	//run the query
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	cursor, err := db.Collection("migrations").Find(ctx, bson.D{}, opts)

	defer cancel()
	if err != nil {
		log.Fatal(err)
	}

	//parse the results
	var results []models.Migration
	err = cursor.All(context.TODO(), &results)
	if err != nil {
		log.Fatal(err)
	}

	//return the newest timestamp from the results
	if len(results) > 0 {
		return results[0].Timestamp
	}

	return ""
}

func saveTimestamp(db *mongo.Database, timestamp string) {
	//save the timestamp to the migrations collection
	_, err := db.Collection("migrations").InsertOne(context.TODO(), bson.D{{"timestamp", timestamp}})
	if err != nil {
		log.Fatal(err)
	}
}
