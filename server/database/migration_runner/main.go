package main

import (
	"flag"
	"log"

	"github.com/blendermux/common"

	migrationrunner "github.com/blendermux/common/migration_runner"
	"github.com/blendermux/server/dependencies"
)

func main() {
	//parse flags
	down := flag.Bool("down", false, "Run the most recent migration down")
	flag.Parse()

	resolver := dependencies.CreateDependencyResolver()
	db := resolver.Database

	//open the db connection
	err := db.OpenConnection()
	if err != nil {
		log.Fatal(common.ChainError("Could not create database connection", err))
	}

	defer db.CloseConnection()

	//check db is connected
	err = db.Ping()
	if err != nil {
		log.Fatal(common.ChainError("Could not reach database", err))
	}

	//run the migrations
	if *down {
		err = migrationrunner.MigrateDown(resolver.MigrationRepository, db)
	} else {
		err = migrationrunner.MigrateUp(resolver.MigrationRepository, db)
	}

	if err != nil {
		log.Fatal(common.ChainError("Error running migrations", err))
	}
}
