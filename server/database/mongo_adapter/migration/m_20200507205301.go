package main

import (
	"context"
	"log"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
)

type M20200507205301 struct{}

func (m M20200507205301) GetTimestamp() string {
	return "20200507205301"
}

func (m M20200507205301) Up(db *mongo.Database) {
	//index users on email
	userIndexes := []mongo.IndexModel{
		{
			Keys: bson.D{{"email", 1}},
		},
	}

	//set user indexes
	_, err := db.Collection("users").Indexes().CreateMany(context.TODO(), userIndexes)
	if err != nil {
		log.Fatal(err)
	}
}

func (m M20200507205301) Down(db *mongo.Database) {
	//remove the created user indexes
	_, err := db.Collection("users").Indexes().DropOne(context.TODO(), "email_1")
	if err != nil {
		log.Fatal(err)
	}
}
