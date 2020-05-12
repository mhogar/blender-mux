package main

import (
	"log"

	migrationrunner "github.com/blendermux/common/migration_runner"
	"github.com/blendermux/server/dependencies"
)

func main() {
	resolver := dependencies.CreateDependencyResolver()
	db := resolver.Database

	//open the db connection
	err := db.OpenConnection()
	if err != nil {
		log.Fatal("Could not create database connection:", err)
	}

	defer db.CloseConnection()

	//check db is connected
	err = db.Ping()
	if err != nil {
		log.Fatal("Could not reach database:", err)
	}

	//run the migrations
	err = migrationrunner.RunMigrations(resolver.MigrationRepository, db)
	if err != nil {
		log.Fatal("Error running migrations:", err)
	}
}
